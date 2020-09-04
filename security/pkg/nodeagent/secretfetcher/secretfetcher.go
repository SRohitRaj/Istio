// Copyright Istio Authors
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

package secretfetcher

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"

	"istio.io/istio/pkg/security"
	nodeagentutil "istio.io/istio/security/pkg/nodeagent/util"
	"istio.io/pkg/env"
	"istio.io/pkg/log"
)

const (
	// The ID/name for the certificate chain in kubernetes generic secret.
	genericScrtCert = "cert"
	// The ID/name for the private key in kubernetes generic secret.
	genericScrtKey = "key"
	// The ID/name for the CA certificate in kubernetes generic secret.
	genericScrtCaCert = "cacert"

	// The ID/name for the certificate chain in kubernetes tls secret.
	tlsScrtCert = "tls.crt"
	// The ID/name for the k8sKey in kubernetes tls secret.
	tlsScrtKey = "tls.key"
	// The ID/name for the CA certificate in kubernetes tls secret
	tlsScrtCaCert = "ca.crt"

	// IngressSecretNamespace the namespace of kubernetes secrets to watch.
	ingressSecretNamespace = "INGRESS_GATEWAY_NAMESPACE"

	// GatewaySdsCaSuffix is the suffix of the sds resource name for root CA. All resource
	// names for gateway root certs end with "-cacert".
	GatewaySdsCaSuffix = "-cacert"

	// scrtTokenField is the token field in secret generated by istio.
	scrtTokenField = "token"

	// istioPrefix and prometheusPrefix are prefix in secrets generated by istio.
	istioPrefix      = "istio"
	prometheusPrefix = "prometheus"
)

var (
	// TODO(JimmyCYJ): Configure these two env variables in Helm
	// secretControllerResyncPeriod specifies the time period in seconds that secret controller
	// resyncs to API server.
	// example value format like "30s"
	secretControllerResyncPeriod = env.RegisterStringVar("SECRET_WATCHER_RESYNC_PERIOD", "", "").Get()
	// ingressFallbackSecret specifies the name of fallback secret for ingress gateway.
	secretFetcherLog = log.RegisterScope("secretfetcher", "secret fetcher debugging", 0)
)

// SecretFetcher fetches secret via watching k8s secrets or sending CSR to CA.
type SecretFetcher struct {
	// If CaClient is set, use caClient to send CSR to CA.
	CaClient security.Client

	// Controller and store for secret objects.
	scrtController cache.Controller
	scrtStore      cache.Store

	// secrets maps k8sKey to secrets
	secrets sync.Map

	// Add all entries containing secretName in SecretCache. Called when K8S secret is added.
	AddCache func(secretName string, ns security.SecretItem)
	// Delete all entries containing secretName in SecretCache. Called when K8S secret is deleted.
	DeleteCache func(secretName string)
	// Update all entries containing secretName in SecretCache. Called when K8S secret is updated.
	UpdateCache func(secretName string, ns security.SecretItem)

	// FallbackSecretName stores the name of fallback secret which is set at env variable
	// INGRESS_GATEWAY_FALLBACK_SECRET. If INGRESS_GATEWAY_FALLBACK_SECRET is empty, then use
	// gateway-fallback as default name of fallback secret. If a fallback secret exists,
	// FindGatewaySecret returns this fallback secret when expected secret is not available.
	FallbackSecretName string

	secretNamespace string
	coreV1          corev1.CoreV1Interface
}

// Run starts the SecretFetcher until a value is sent to ch.
// Only used when watching kubernetes gateway secrets.
func (sf *SecretFetcher) Run(ch chan struct{}) {
	go sf.scrtController.Run(ch)
	cache.WaitForCacheSync(ch, sf.scrtController.HasSynced)
}

var namespaceVar = env.RegisterStringVar(ingressSecretNamespace, "", "")

// InitWithKubeClient initializes SecretFetcher to watch kubernetes secrets.
func (sf *SecretFetcher) InitWithKubeClient(core corev1.CoreV1Interface) { // nolint:interfacer
	sf.InitWithKubeClientAndNs(core, namespaceVar.Get())
}

// InitWithKubeClientAndNs initializes SecretFetcher to watch kubernetes secrets.
func (sf *SecretFetcher) InitWithKubeClientAndNs(core corev1.CoreV1Interface, namespace string) { // nolint:interfacer
	istioSecretSelector := fields.SelectorFromSet(nil).String()
	scrtLW := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = istioSecretSelector
			return core.Secrets(namespace).List(context.TODO(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = istioSecretSelector
			return core.Secrets(namespace).Watch(context.TODO(), options)
		},
	}

	resyncPeriod := 0 * time.Second
	if e, err := time.ParseDuration(secretControllerResyncPeriod); err == nil {
		resyncPeriod = e
	}

	sf.scrtStore, sf.scrtController =
		cache.NewInformer(scrtLW, &v1.Secret{}, resyncPeriod, cache.ResourceEventHandlerFuncs{
			AddFunc:    sf.scrtAdded,
			DeleteFunc: sf.scrtDeleted,
			UpdateFunc: sf.scrtUpdated,
		})

	sf.secretNamespace = namespace
	sf.coreV1 = core

}

// isGatewaySecret checks secret and decides whether this is a secret generated for ingress
// gateway. For secrets with prefix "istio" and "prometheus", they are generated by istio system and
// they are not gateway secrets. Other secrets generated by istio could have "token" field.
// isGatewaySecret returns false if a secret has name prefix "istio", or "prometheus", or has
// "token" field.
func isGatewaySecret(scrt *v1.Secret) bool {
	secretName := scrt.GetName()
	if strings.HasPrefix(secretName, istioPrefix) || strings.HasPrefix(secretName, prometheusPrefix) {
		return false
	}
	if len(scrt.Data[scrtTokenField]) > 0 {
		return false
	}
	return true
}

// extractCertAndKey extracts server key, certificate, and indicates whether key and cert exist.
func extractCertAndKey(scrt *v1.Secret) (cert, key []byte, exist bool) {
	certAndKeyExist := false
	if len(scrt.Data[genericScrtCert]) > 0 {
		cert = scrt.Data[genericScrtCert]
		key = scrt.Data[genericScrtKey]
	} else {
		cert = scrt.Data[tlsScrtCert]
		key = scrt.Data[tlsScrtKey]
	}
	if len(cert) > 0 && len(key) > 0 {
		certAndKeyExist = true
	}

	return cert, key, certAndKeyExist
}

// extractCACert extracts the client CA certificate from either the Compound
// Secret, or from a separate Kubernetes TLS secret that has CA cert in `tls.crt` field.
func extractCACert(scrt *v1.Secret, fromCompoundSecret bool) (caCert []byte, exist bool) {
	if len(scrt.Data[genericScrtCaCert]) > 0 {
		caCert = scrt.Data[genericScrtCaCert]
	} else if len(scrt.Data[tlsScrtCaCert]) > 0 {
		caCert = scrt.Data[tlsScrtCaCert]
	} else if !fromCompoundSecret {
		caCert = scrt.Data[tlsScrtCert]
	}

	return caCert, len(caCert) > 0
}

// extractK8sSecretIntoSecretItem extracts a server cert/key pair and a client CA
// certificate from the k8s Secret into a pair of SecretItems. Returns SecretItems and a boolean
// indicating whether this is a CA only k8s Secret.
// A CA only k8s secret has name suffix `-cacert`, and is ONLY considered for a client CA;
// either a `cacert` or `tls.crt` must be provided.
// Otherwise the Secret can hold a server cert/key pair in `tls.crt`/`tls.key`,
// or a server cert/key pair in `cert`/`key` and an optional client CA cert in
// `-cacert`. A Secret with server cert/key and client CA cert is considered as a compound secret.
func extractK8sSecretIntoSecretItem(scrt *v1.Secret, t time.Time) (serverItem, clientCAItem *security.SecretItem, isCAOnlySecret bool) {
	resourceName := scrt.GetName()
	isCAOnlySecret = strings.HasSuffix(resourceName, GatewaySdsCaSuffix)

	// Extract CA cert from CA only k8s secret.
	if isCAOnlySecret {
		caCert, exist := extractCACert(scrt, false /* fromCompoundSecret */)
		if !exist {
			secretFetcherLog.Warnf("failed load CA only secret from %s: no 'cacert' or 'tls.crt' key in the secret", resourceName)
			return nil, nil, isCAOnlySecret
		}
		rootCertExpireTime, err := nodeagentutil.ParseCertAndGetExpiryTimestamp(caCert)
		if err != nil {
			secretFetcherLog.Warnf("skip loading secret. Kubernetes secret %v contains a root "+
				"certificate that fails to parse: %v", resourceName, err)
			return nil, nil, isCAOnlySecret
		}

		certificateAuthorityNewSecret := &security.SecretItem{
			ResourceName:                  resourceName,
			CreatedTime:                   t,
			Version:                       t.String(),
			RootCertOwnedByCompoundSecret: false,
			RootCert:                      caCert,
			ExpireTime:                    rootCertExpireTime,
		}

		return nil, certificateAuthorityNewSecret, isCAOnlySecret
	}

	// Extract server key/cert from k8s secret.
	cert, key, keyCertExist := extractCertAndKey(scrt)
	if !keyCertExist {
		secretFetcherLog.Warnf("failed load server cert/key pair from secret %s: server cert or private key is empty", resourceName)
		return nil, nil, isCAOnlySecret
	}
	certExpireTime, err := nodeagentutil.ParseCertAndGetExpiryTimestamp(cert)
	if err != nil {
		secretFetcherLog.Warnf("skip loading secret. Kubernetes secret %v contains a server "+
			"certificate that fails to parse: %v", resourceName, err)
		return nil, nil, isCAOnlySecret
	}
	newSecret := &security.SecretItem{
		ResourceName:     resourceName,
		CreatedTime:      t,
		Version:          t.String(),
		CertificateChain: cert,
		ExpireTime:       certExpireTime,
		PrivateKey:       key,
	}

	// Try to extract CA cert from k8s secret.
	caCert, caCertExist := extractCACert(scrt, true /* fromCompoundSecret */)
	if caCertExist {
		rootCertExpireTime, err := nodeagentutil.ParseCertAndGetExpiryTimestamp(caCert)
		if err != nil {
			secretFetcherLog.Warnf("skip loading secret. Kubernetes secret %v contains a root "+
				"certificate that fails to parse: %v", resourceName, err)
			return nil, nil, isCAOnlySecret
		}
		certificateAuthorityNewSecret := &security.SecretItem{
			ResourceName:                  resourceName + GatewaySdsCaSuffix,
			CreatedTime:                   t,
			Version:                       t.String(),
			RootCert:                      caCert,
			ExpireTime:                    rootCertExpireTime,
			RootCertOwnedByCompoundSecret: true,
		}
		return newSecret, certificateAuthorityNewSecret, isCAOnlySecret
	}

	return newSecret, nil, isCAOnlySecret
}

func (sf *SecretFetcher) scrtAdded(obj interface{}) {
	scrt, ok := obj.(*v1.Secret)
	if !ok {
		secretFetcherLog.Warnf("Failed to convert to secret object: %v", obj)
		return
	}

	resourceName := scrt.GetName()
	if !isGatewaySecret(scrt) {
		secretFetcherLog.Debugf("secret %s is not a gateway secret, skip adding secret", resourceName)
		return
	}

	t := time.Now()
	newSecret, certificateAuthorityNewSecret, isCaOnly := extractK8sSecretIntoSecretItem(scrt, t)

	// Load CA cert from CA only k8s secret and update cache.
	if isCaOnly && certificateAuthorityNewSecret != nil {
		sf.secrets.Delete(certificateAuthorityNewSecret.ResourceName)
		sf.secrets.Store(certificateAuthorityNewSecret.ResourceName, *certificateAuthorityNewSecret)
		secretFetcherLog.Debugf("secret %s is added as a client CA cert", certificateAuthorityNewSecret.ResourceName)
		if sf.AddCache != nil {
			sf.AddCache(certificateAuthorityNewSecret.ResourceName, *certificateAuthorityNewSecret)
		}
		return
	}

	if newSecret != nil {
		// Load server key/cert from k8s secret and update cache.
		sf.secrets.Delete(newSecret.ResourceName)
		sf.secrets.Store(newSecret.ResourceName, *newSecret)
		secretFetcherLog.Debugf("secret %s is added as a server certificate", newSecret.ResourceName)
		if sf.AddCache != nil {
			sf.AddCache(newSecret.ResourceName, *newSecret)
		}
		if certificateAuthorityNewSecret != nil {
			// Load client CA cert from compound k8s secret and update cache.
			sf.secrets.Delete(certificateAuthorityNewSecret.ResourceName)
			sf.secrets.Store(certificateAuthorityNewSecret.ResourceName, *certificateAuthorityNewSecret)
			secretFetcherLog.Debugf("secret %s is added as a client CA cert (from a compound Secret)", certificateAuthorityNewSecret.ResourceName)
			if sf.AddCache != nil {
				sf.AddCache(certificateAuthorityNewSecret.ResourceName, *certificateAuthorityNewSecret)
			}
		}
	}
}

func (sf *SecretFetcher) scrtDeleted(obj interface{}) {
	scrt, ok := obj.(*v1.Secret)
	if !ok {
		secretFetcherLog.Warnf("Failed to convert to secret object: %v", obj)
		return
	}

	key := scrt.GetName()
	sf.secrets.Delete(key)
	secretFetcherLog.Infof("secret %s is deleted", key)
	// Delete all cache entries that match the deleted key.
	if sf.DeleteCache != nil {
		sf.DeleteCache(key)
	}

	rootCertResourceName := key + GatewaySdsCaSuffix
	rootSecret, exists := sf.secrets.Load(rootCertResourceName)
	// If there is a root cert secret with the same resource name and it's owned
	// by the compound K8S secret, delete it now.
	if exists && rootSecret.(security.SecretItem).RootCertOwnedByCompoundSecret {
		// If there is root cert secret with the same resource name, delete that secret now.
		sf.secrets.Delete(rootCertResourceName)
		secretFetcherLog.Infof("secret %s is deleted", rootCertResourceName)
		// Delete all cache entries that match the deleted key.
		if sf.DeleteCache != nil {
			sf.DeleteCache(rootCertResourceName)
		}
	}
}

func (sf *SecretFetcher) scrtUpdated(oldObj, newObj interface{}) {
	oscrt, ok := oldObj.(*v1.Secret)
	if !ok {
		secretFetcherLog.Warnf("Failed to convert to old secret object: %v", oldObj)
		return
	}
	nscrt, ok := newObj.(*v1.Secret)
	if !ok {
		secretFetcherLog.Warnf("Failed to convert to new secret object: %v", newObj)
		return
	}

	oldScrtName := oscrt.GetName()
	newScrtName := nscrt.GetName()
	if oldScrtName != newScrtName {
		secretFetcherLog.Warnf("Failed to update secret: name does not match (%s vs %s).", oldScrtName, newScrtName)
		return
	}

	if !isGatewaySecret(nscrt) {
		secretFetcherLog.Debugf("kubernetes secret %s is not an gateway secret, skip update", newScrtName)
		return
	}

	secretFetcherLog.Infof("scrtUpdated is called on kubernetes secret %s", newScrtName)
	// Kubernetes secret update is done by deleting first and creating a new one with the same name.
	// Accordingly scrtDeleted and scrtAdded are called. When scrtUpdated is called, secret should remain unchanged.
	t := time.Now()
	oldScrt, oldCaScrt, _ := extractK8sSecretIntoSecretItem(oscrt, t)
	newScrt, newCaScrt, isCaOnlyNew := extractK8sSecretIntoSecretItem(nscrt, t)
	updateSecret := shouldUpdateSecret(oldScrt, oldCaScrt, newScrt, newCaScrt)
	if !updateSecret {
		secretFetcherLog.Infof("secret %s does not change, skip update", newScrtName)
		return
	}

	if oldScrt != nil || newScrt != nil {
		sf.updateSecretInCache(oldScrt, newScrt)
		secretFetcherLog.Debugf("secret %s is updated as a server certificate", newScrt.ResourceName)
	}
	if oldCaScrt != nil || newCaScrt != nil {
		sf.updateSecretInCache(oldCaScrt, newCaScrt)
		if isCaOnlyNew {
			secretFetcherLog.Debugf("secret %s is updated as a client CA cert (from a CA only secret)",
				newCaScrt.ResourceName)
		} else {
			secretFetcherLog.Debugf("secret %s is updated as a client CA cert (from a compound Secret)",
				newCaScrt.ResourceName)
		}
	}
}

// updateSecretInCache updates secret in cache, and pushes to client when new certs
// are reloaded from secret.
func (sf *SecretFetcher) updateSecretInCache(oldScrt, newScrt *security.SecretItem) {
	if oldScrt != nil {
		sf.secrets.Delete(oldScrt.ResourceName)
	}
	if newScrt != nil {
		sf.secrets.Store(newScrt.ResourceName, *newScrt)
		if sf.UpdateCache != nil {
			sf.UpdateCache(newScrt.ResourceName, *newScrt)
		}
	} else if oldScrt != nil {
		if sf.DeleteCache != nil {
			sf.DeleteCache(oldScrt.ResourceName)
		}
	}
}

// shouldUpdateSecret indicates whether secret update is required to reload new secret.
func shouldUpdateSecret(oldScrt, oldCaScrt, newScrt, newCaScrt *security.SecretItem) bool {
	if newScrt == nil && newCaScrt == nil {
		return false
	}

	if (oldScrt != nil && newScrt == nil) || (oldScrt == nil && newScrt != nil) {
		return true
	}
	if newScrt != nil && oldScrt != nil {
		if !bytes.Equal(oldScrt.CertificateChain, newScrt.CertificateChain) ||
			!bytes.Equal(oldScrt.PrivateKey, newScrt.PrivateKey) {
			return true
		}
	}

	if (oldCaScrt != nil && newCaScrt == nil) || (oldCaScrt == nil && newCaScrt != nil) {
		return true
	}
	if oldCaScrt != nil && newCaScrt != nil {
		if !bytes.Equal(oldCaScrt.RootCert, newCaScrt.RootCert) {
			return true
		}
	}
	return false
}

// FindGatewaySecret returns the secret whose name matches the key, or empty secret if no
// secret is present. The ok result indicates whether secret was found.
// If there is a fallback secret named FallbackSecretName, return the fall back secret.
func (sf *SecretFetcher) FindGatewaySecret(key string) (secret security.SecretItem, ok bool) {
	secretFetcherLog.Debugf("SecretFetcher search for secret %s", key)
	val, exist := sf.secrets.Load(key)
	secretFetcherLog.Debugf("load secret %s from secret fetcher: %v", key, exist)
	if !exist {
		// Sometimes we see that a secret in installed but not in cache because watcher is in an
		// obsolete state and wasn't reset promptly. We bail this case out by trying fetching
		// the secret from API call. Since this is a rare case, to avoid complication, we don't add
		// the secret back to cache as it is not a normal codepath. When watcher recovers, those secret
		// shall be added back. Note that this approach only covers the TLS server key/cert fetching.
		if sf.coreV1 != nil {
			if secret, err := sf.coreV1.Secrets(sf.secretNamespace).Get(context.TODO(), key, metav1.GetOptions{}); err == nil {
				secretItem, _, _ := extractK8sSecretIntoSecretItem(secret, time.Now())
				if secretItem != nil {
					secretFetcherLog.Infof("Return secret %s found by direct api call", key)
					return *secretItem, true
				}
				secretFetcherLog.Infof("Fail to extract secret %s found by direct api call", key)
			}
		}

		// Expected secret does not exist, try to find the fallback secret.
		// TODO(JimmyCYJ): Add metrics to node agent to imply usage of fallback secret
		secretFetcherLog.Warnf("Cannot find secret %s, searching for fallback secret %s", key, sf.FallbackSecretName)
		fallbackVal, fallbackExist := sf.secrets.Load(sf.FallbackSecretName)
		if fallbackExist {
			secretFetcherLog.Debugf("Return fallback secret %s for gateway secret %s", sf.FallbackSecretName, key)
			return fallbackVal.(security.SecretItem), true
		}

		secretFetcherLog.Errorf("cannot find secret %s and cannot find fallback secret %s", key, sf.FallbackSecretName)
		return security.SecretItem{}, false
	}
	e := val.(security.SecretItem)
	secretFetcherLog.Debugf("SecretFetcher return secret %s", key)
	return e, true
}

// AddSecret adds obj into local store. Only used for testing.
func (sf *SecretFetcher) AddSecret(obj interface{}) {
	sf.scrtAdded(obj)
}

// DeleteSecret deletes obj from local store. Only used for testing.
func (sf *SecretFetcher) DeleteSecret(obj interface{}) {
	sf.scrtDeleted(obj)
}
