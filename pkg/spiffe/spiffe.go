package spiffe

import (
	"fmt"
	"strings"

	"istio.io/istio/pkg/log"
)

const (
	Scheme = "spiffe"

	URIPrefix = Scheme + "://"

	// The default SPIFFE URL value for trust domain
	defaultTrustDomain = "cluster.local"
	uriPrefix          = Scheme + "://"
)

var (
	trustDomain = defaultTrustDomain
)

func SetTrustDomain(value string) {
	// Replace special characters in spiffe
	trustDomain = strings.Replace(value, "@", ".", -1)
}

func GetTrustDomain() string {
	return trustDomain
}

func DetermineTrustDomain(commandLineTrustDomain string, isKubernetes bool) string {

	if len(commandLineTrustDomain) != 0 {
		return commandLineTrustDomain
	}
	if isKubernetes {
		return defaultTrustDomain
	}
	return ""
}

// GenSpiffeURI returns the formatted uri(SPIFFEE format for now) for the certificate.
func GenSpiffeURI(ns, serviceAccount string) (string, error) {
	var err error
	if ns == "" || serviceAccount == "" {
		err = fmt.Errorf(
			"namespace or service account can't be empty ns=%v serviceAccount=%v", ns, serviceAccount)
	}
	return URIPrefix + trustDomain + "/ns/" + ns + "/sa/" + serviceAccount, err
}

// MustGenSpiffeURI returns the formatted uri(SPIFFEE format for now) for the certificate and logs if there was an error.
func MustGenSpiffeURI(ns, serviceAccount string) string {
	uri, err := GenSpiffeURI(ns, serviceAccount)
	if err != nil {
		log.Error(err.Error())
	}
	return uri
}

// GenCustomSpiffe returns the  spiffe string that can have a custom structure
func GenCustomSpiffe(identity string) string {

	if identity == "" {
		log.Error("spiffe identity can't be empty")
		return ""
	}

	// replace specifial character in spiffe
	trustDomain = strings.Replace(trustDomain, "@", ".", -1)
	return uriPrefix + trustDomain + "/" + identity
}
