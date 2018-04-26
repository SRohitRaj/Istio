// Copyright 2018 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	authn "istio.io/api/authentication/v1alpha1"
	"istio.io/istio/pkg/cache"
	"istio.io/istio/pkg/log"
)

const (
	// https://openid.net/specs/openid-connect-discovery-1_0.html
	// OpenID Providers supporting Discovery MUST make a JSON document available at the path
	// formed by concatenating the string /.well-known/openid-configuration to the Issuer.
	openIDDiscoveryCfgURLSuffix = "/.well-known/openid-configuration"

	// OpenID Discovery web request timeout.
	jwksHTTPTimeOutInSec = 5

	// JwksURI Cache expiration time duration, individual cached JwksURI item will be removed
	// from cache after its duration expires.
	jwksURICacheExpiration = time.Hour * 24

	// JwksURI Cache eviction time duration, cache eviction is done on a periodic basis,
	// jwksURICacheEviction specifies the frequency at which eviction activities take place.
	jwksURICacheEviction = time.Minute * 30

	// JwtPubKeyExpireDuration is the expire duration for JWT public key in the cache.
	// After this duration expire, refresher job will fetch key for the cached item again.
	JwtPubKeyExpireDuration = time.Hour

	// JwtPubKeyEvictionDuration is the life duration for cached item.
	// Cached item will be removed from the cache if it hasn't been used longer than JwtPubKeyEvictionDuration.
	JwtPubKeyEvictionDuration = 24 * 7 * time.Hour

	// JwtPubKeyRefreshInterval is the running interval of JWT pubKey refresh job.
	JwtPubKeyRefreshInterval = time.Minute * 20
)

// jwtPubKeyEntry is a single cached entry for jwt public key.
type jwtPubKeyEntry struct {
	pubKey string

	// Cached item will be fetched again by refresher job if (time.now >= expireTime).
	expireTime time.Time

	// Cached item's last used time, which is set in GetPublicKey.
	lastUsedTime time.Time
}

// jwksResolver is resolver for jwksURI and jwt public key.
type jwksResolver struct {
	// cache for jwksURI.
	JwksURICache cache.ExpiringCache

	// cache for JWT public key.
	// map key is jwksURI, map value is jwtPubKeyEntry.
	keyEntries sync.Map

	client        *http.Client
	closing       chan bool
	refreshTicker *time.Ticker

	expireDuration time.Duration

	// Cached key will be removed from cache if (time.now - cachedItem.lastUsedTime >= evictionDuration), this prevents key cache growing indefinitely.
	evictionDuration time.Duration

	// Refresher job running interval.
	refreshInterval time.Duration

	// How may times refresh job has detected JWT public key change happened, used in unit test.
	keyChangedCount uint64
}

// newJwksResolver creates new instance of jwksResolver.
func newJwksResolver(expireDuration, evictionDuration, refreshInterval time.Duration) *jwksResolver {
	ret := &jwksResolver{
		JwksURICache:     cache.NewTTL(jwksURICacheExpiration, jwksURICacheEviction),
		closing:          make(chan bool, 1),
		expireDuration:   expireDuration,
		evictionDuration: evictionDuration,
		refreshInterval:  refreshInterval,
		client: &http.Client{
			Timeout: jwksHTTPTimeOutInSec * time.Second,
		},
	}

	if false {
		log.Infof("*******************newJwksResolver secure")
		certDir := os.Getenv("ISTIO_CERT_DIR")
		if certDir == "" {
			certDir = "/etc/certs"
		}

		chainCertFile := fmt.Sprintf("%v/%v", certDir, CertChainFilename)
		chainKeyFile := fmt.Sprintf("%v/%v", certDir, KeyFilename)
		chainCert, err := tls.LoadX509KeyPair(chainCertFile, chainKeyFile)
		if err != nil {
			log.Infof("Unable to load certs to talk to control plane: %v", err)
		}
		caCertFile := fmt.Sprintf("%v/%v", certDir, RootCertFilename)
		caCert, err := ioutil.ReadFile(caCertFile)
		if err != nil {
			log.Infof("Unable to load ca root cert to talk to control plane: %v", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{chainCert},
			RootCAs:      caCertPool,
		}
		transport := &http.Transport{TLSClientConfig: tlsConfig}
		ret.client = &http.Client{
			Transport: transport,
			Timeout:   jwksHTTPTimeOutInSec * time.Second,
		}
	} else {
		if false {
			ret.client = &http.Client{
				Timeout: jwksHTTPTimeOutInSec * time.Second,
			}
		} else {
			log.Infof("*******************newJwksResolver InsecureSkipVerify")

			transport := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			ret.client = &http.Client{
				Transport: transport,
				Timeout:   jwksHTTPTimeOutInSec * time.Second,
			}
		}
	}

	atomic.StoreUint64(&ret.keyChangedCount, 0)
	go ret.refresher()

	return ret
}

// Set jwks_uri through openID discovery if it's not set in auth policy.
func (r *jwksResolver) SetAuthenticationPolicyJwksURIs(policy *authn.Policy) error {
	if policy == nil {
		return fmt.Errorf("invalid nil policy")
	}

	for _, method := range policy.Peers {
		switch method.GetParams().(type) {
		case *authn.PeerAuthenticationMethod_Jwt:
			policyJwt := method.GetJwt()
			if policyJwt.JwksUri == "" {
				uri, err := r.resolveJwksURIUsingOpenID(policyJwt.Issuer)
				if err != nil {
					log.Warnf("Failed to get jwks_uri for issuer %q: %v", policyJwt.Issuer, err)
					return err
				}
				policyJwt.JwksUri = uri
			}
		}
	}
	for _, method := range policy.Origins {
		// JWT is only allowed authentication method type for Origin.
		policyJwt := method.GetJwt()
		if policyJwt.JwksUri == "" {
			uri, err := r.resolveJwksURIUsingOpenID(policyJwt.Issuer)
			if err != nil {
				log.Warnf("Failed to get jwks_uri for issuer %q: %v", policyJwt.Issuer, err)
				return err
			}
			policyJwt.JwksUri = uri
		}
	}

	return nil
}

// GetPublicKey gets JWT public key and cache the key for future use.
func (r *jwksResolver) GetPublicKey(jwksURI string) (string, error) {
	now := time.Now()
	if val, found := r.keyEntries.Load(jwksURI); found {
		e := val.(jwtPubKeyEntry)

		log.Infof("******************GetPublicKey found %+v", e)

		// Return from cache if it's not expired.
		if e.expireTime.After(now) {
			// Update cached key's last used time.
			e.lastUsedTime = now
			r.keyEntries.Store(jwksURI, e)

			log.Infof("******************GetPublicKey found unexpire %+v", e)

			return e.pubKey, nil
		}
	}

	// Fetch key if it's not cached, or cached item is expired.
	resp, err := r.getRemoteContent(jwksURI)
	if err != nil {
		log.Errorf("**************failed to get key from %q", jwksURI)
		return "", err
	}

	pubKey := string(resp)
	r.keyEntries.Store(jwksURI, jwtPubKeyEntry{
		pubKey:       pubKey,
		expireTime:   now.Add(r.expireDuration),
		lastUsedTime: now,
	})

	return pubKey, nil
}

// Resolve jwks_uri through openID discovery and cache the jwks_uri for future use.
func (r *jwksResolver) resolveJwksURIUsingOpenID(issuer string) (string, error) {
	// Set policyJwt.JwksUri if the JwksUri could be found in cache.
	if uri, found := r.JwksURICache.Get(issuer); found {
		return uri.(string), nil
	}

	// Try to get jwks_uri through OpenID Discovery.
	body, err := r.getRemoteContent(issuer + openIDDiscoveryCfgURLSuffix)
	if err != nil {
		return "", err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	jwksURI, ok := data["jwks_uri"].(string)
	if !ok {
		return "", fmt.Errorf("invalid jwks_uri %v in openID discovery configuration", data["jwks_uri"])
	}

	// Set JwksUri in cache.
	r.JwksURICache.Set(issuer, jwksURI)

	return jwksURI, nil
}

func (r *jwksResolver) getRemoteContent(url string) ([]byte, error) {
	//resp, err := r.client.Get(url)
	/*
		resp, err := http.Get(url)
		if err != nil {
			log.Errorf("*************getRemoteContent get from %q get failed, %v", url, err)
			return nil, err
		}
		defer func() {
			_ = resp.Body.Close()
		}() */

	resp, err := r.client.Get(url)
	if err != nil {
		log.Errorf("*************resptest getRemoteContent get from %q get failed, %v", url, err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	log.Infof("***************************resptest getRemoteContent %+v", resp)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Errorf("*************getRemoteContent get from %q response code %d", url, resp.StatusCode)
		return nil, fmt.Errorf("unsuccessful response from %q", url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("*************getRemoteContent failed to read body %v", err)
		return nil, err
	}

	log.Errorf("*************getRemoteContent body %v", body)
	return body, nil
}

func (r *jwksResolver) refresher() {
	// Wake up once in a while and refresh stale items.
	r.refreshTicker = time.NewTicker(r.refreshInterval)
	for {
		select {
		case now := <-r.refreshTicker.C:
			r.refresh(now)
		case <-r.closing:
			r.refreshTicker.Stop()
			return
		}
	}
}

func (r *jwksResolver) refresh(t time.Time) {
	var wg sync.WaitGroup
	hasChange := false

	r.keyEntries.Range(func(key interface{}, value interface{}) bool {
		now := time.Now()
		jwksURI := key.(string)
		e := value.(jwtPubKeyEntry)

		// Remove cached item if it hasn't been used for a while.
		if now.Sub(e.lastUsedTime) >= r.evictionDuration {
			r.keyEntries.Delete(jwksURI)
			return true
		}

		oldPubKey := e.pubKey

		// key rotation: fetch JWT public key again if it's expired.
		if e.expireTime.Before(t) {
			// Increment the WaitGroup counter.
			wg.Add(1)

			go func() {
				// Decrement the counter when the goroutine completes.
				defer wg.Done()

				resp, err := r.getRemoteContent(jwksURI)
				if err != nil {
					log.Errorf("Cannot fetch JWT public key from %q", jwksURI)
					r.keyEntries.Delete(jwksURI)
					return
				}
				newPubKey := string(resp)

				r.keyEntries.Store(jwksURI, jwtPubKeyEntry{
					pubKey:       newPubKey,
					expireTime:   now.Add(r.expireDuration), // Update expireTime even if prev/current keys are the same.
					lastUsedTime: e.lastUsedTime,            // keep original lastUsedTime.
				})

				if oldPubKey != newPubKey {
					hasChange = true
				}
			}()
		}

		return true
	})

	// Wait for all go routine to complete.
	wg.Wait()

	if hasChange {
		atomic.AddUint64(&r.keyChangedCount, 1)
		// TODO(quanlin): send notification to update config and push config to sidecar.
	}
}

// Shut down the refresher job.
// TODO: may need to figure out the right place to call this function.
// (right now calls it from initDiscoveryService in pkg/bootstrap/server.go).
func (r *jwksResolver) Close() {
	r.closing <- true
}
