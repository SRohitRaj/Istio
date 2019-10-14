// Copyright 2017 Istio Authors
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

package controller

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"

	"istio.io/istio/pkg/log"
	"istio.io/istio/pkg/spiffe"
	k8ssecret "istio.io/istio/security/pkg/k8s/secret"
	"istio.io/istio/security/pkg/listwatch"
	caerror "istio.io/istio/security/pkg/pki/error"
	"istio.io/istio/security/pkg/pki/util"
	certutil "istio.io/istio/security/pkg/util"
)

/* #nosec: disable gas linter */
const (
	// The Istio secret annotation type
	IstioSecretType = "istio.io/key-and-cert"

	// The ID/name for the certificate chain file.
	CertChainID = "cert-chain.pem"
	// The ID/name for the private key file.
	PrivateKeyID = "key.pem"
	// The ID/name for the CA root certificate file.
	RootCertID = "root-cert.pem"
	// The key to specify corresponding service account in the annotation of K8s secrets.
	ServiceAccountNameAnnotationKey = "istio.io/service-account.name"

	secretNamePrefix   = "istio."
	secretResyncPeriod = time.Minute

	recommendedMinGracePeriodRatio = 0.2
	recommendedMaxGracePeriodRatio = 0.8

	// The size of a private key for a leaf certificate.
	keySize = 2048

	// The number of retries when requesting to create secret.
	secretCreationRetry = 3

	// CASecret stores the key/cert of self-signed CA for persistency purpose.
	CASecret = "istio-ca-secret"
	// caCertID is the CA certificate chain file.
	caCertID = "ca-cert.pem"
	// caPrivateKeyID is the private key file of CA.
	caPrivateKeyID = "ca-key.pem"
)

// DNSNameEntry stores the service name and namespace to construct the DNS id.
// Service accounts matching the ServiceName and Namespace will have additional DNS SANs:
// ServiceName.Namespace.svc, ServiceName.Namespace and optionall CustomDomain.
// This is intended for control plane and trusted services.
type DNSNameEntry struct {
	// ServiceName is the name of the service account to match
	ServiceName string

	// Namespace restricts to a specific namespace.
	Namespace string

	// CustomDomain allows adding a user-defined domain.
	CustomDomains []string
}

// certificateAuthority contains methods to be supported by a CA.
type certificateAuthority interface {
	// Sign generates a certificate for a workload or CA, from the given CSR and TTL.
	// TODO(myidpt): simplify this interface and pass a struct with cert field values instead.
	Sign(csrPEM []byte, subjectIDs []string, ttl time.Duration, forCA bool) ([]byte, error)
	// GetCAKeyCertBundle returns the KeyCertBundle used by CA.
	GetCAKeyCertBundle() util.KeyCertBundle
}

// SecretController manages the service accounts' secrets that contains Istio keys and certificates.
type SecretController struct {
	monitoring monitoringMetrics
	ca         certificateAuthority
	core       corev1.CoreV1Interface
	certUtil   certutil.CertUtil

	// Controller and store for service account objects.
	saController cache.Controller
	saStore      cache.Store

	// Controller and store for secret objects.
	scrtController cache.Controller
	scrtStore      cache.Store

	caSecretController *CaSecretController

	// Controller and store for namespace objects
	namespaceController cache.Controller
	namespaceStore      cache.Store

	// Used to coordinate with label and check if this instance of Citadel should create secret
	istioCaStorageNamespace string

	certTTL time.Duration

	minGracePeriod time.Duration

	// The set of namespaces explicitly set for monitoring via commandline (an entry could be metav1.NamespaceAll)
	namespaces map[string]struct{}

	// DNS-enabled serviceAccount.namespace to service pair
	dnsNames map[string]*DNSNameEntry

	// Length of the grace period for the certificate rotation.
	gracePeriodRatio float32

	// Whether the certificates are for dual-use clients (SAN+CN).
	dualUse bool

	// Whether the certificates are for CAs.
	forCA bool

	// The most recent time when root cert in keycertbundle is synced with root
	// cert in istio-ca-secret.
	lastKCBSyncTime time.Time
}

// NewSecretController returns a pointer to a newly constructed SecretController instance.
func NewSecretController(ca certificateAuthority, certTTL time.Duration, gracePeriodRatio float32,
	minGracePeriod time.Duration, dualUse bool, core corev1.CoreV1Interface, forCA bool, namespaces []string,
	dnsNames map[string]*DNSNameEntry, istioCaStorageNamespace string) (*SecretController, error) {
	if gracePeriodRatio < 0 || gracePeriodRatio > 1 {
		return nil, fmt.Errorf("grace period ratio %f should be within [0, 1]", gracePeriodRatio)
	}
	if gracePeriodRatio < recommendedMinGracePeriodRatio || gracePeriodRatio > recommendedMaxGracePeriodRatio {
		log.Warnf("grace period ratio %f is out of the recommended window [%.2f, %.2f]",
			gracePeriodRatio, recommendedMinGracePeriodRatio, recommendedMaxGracePeriodRatio)
	}

	c := &SecretController{
		ca:                      ca,
		certTTL:                 certTTL,
		istioCaStorageNamespace: istioCaStorageNamespace,
		gracePeriodRatio:        gracePeriodRatio,
		certUtil:                certutil.NewCertUtil(int(gracePeriodRatio * 100)),
		caSecretController:      NewCaSecretController(core),
		minGracePeriod:          minGracePeriod,
		dualUse:                 dualUse,
		core:                    core,
		forCA:                   forCA,
		namespaces:              make(map[string]struct{}),
		dnsNames:                dnsNames,
		monitoring:              newMonitoringMetrics(),
		lastKCBSyncTime:         time.Time{},
	}

	for _, ns := range namespaces {
		c.namespaces[ns] = struct{}{}
	}

	saLW := listwatch.MultiNamespaceListerWatcher(namespaces, func(namespace string) cache.ListerWatcher {
		return &cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return core.ServiceAccounts(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return core.ServiceAccounts(namespace).Watch(options)
			},
		}
	})

	rehf := cache.ResourceEventHandlerFuncs{
		AddFunc:    c.saAdded,
		DeleteFunc: c.saDeleted,
		UpdateFunc: c.saUpdated,
	}
	c.saStore, c.saController = cache.NewInformer(saLW, &v1.ServiceAccount{}, time.Minute, rehf)

	istioSecretSelector := fields.SelectorFromSet(map[string]string{"type": IstioSecretType}).String()
	scrtLW := listwatch.MultiNamespaceListerWatcher(namespaces, func(namespace string) cache.ListerWatcher {
		return &cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				options.FieldSelector = istioSecretSelector
				return core.Secrets(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.FieldSelector = istioSecretSelector
				return core.Secrets(namespace).Watch(options)
			},
		}
	})
	c.scrtStore, c.scrtController =
		cache.NewInformer(scrtLW, &v1.Secret{}, secretResyncPeriod, cache.ResourceEventHandlerFuncs{
			DeleteFunc: c.scrtDeleted,
			UpdateFunc: c.scrtUpdated,
		})

	return c, nil
}

// Run starts the SecretController until a value is sent to stopCh.
func (sc *SecretController) Run(stopCh chan struct{}) {
	go sc.scrtController.Run(stopCh)

	// saAdded calls upsertSecret to update and insert secret
	// it throws error if the secret cache is not synchronized, but the secret exists in the system
	cache.WaitForCacheSync(stopCh, sc.scrtController.HasSynced)

	go sc.saController.Run(stopCh)
}

// GetSecretName returns the secret name for a given service account name.
func GetSecretName(saName string) string {
	return secretNamePrefix + saName
}

// Handles the event where a service account is added.
func (sc *SecretController) saAdded(obj interface{}) {
	acct := obj.(*v1.ServiceAccount)
	sc.upsertSecret(acct.GetName(), acct.GetNamespace())
	sc.monitoring.ServiceAccountCreation.Inc()
}

// Handles the event where a service account is deleted.
func (sc *SecretController) saDeleted(obj interface{}) {
	acct := obj.(*v1.ServiceAccount)
	sc.deleteSecret(acct.GetName(), acct.GetNamespace())
	sc.monitoring.ServiceAccountDeletion.Inc()
}

// Handles the event where a service account is updated.
func (sc *SecretController) saUpdated(oldObj, curObj interface{}) {
	if reflect.DeepEqual(oldObj, curObj) {
		// Nothing is changed. The method is invoked by periodical re-sync with the apiserver.
		return
	}
	oldSa := oldObj.(*v1.ServiceAccount)
	curSa := curObj.(*v1.ServiceAccount)

	curName := curSa.GetName()
	curNamespace := curSa.GetNamespace()
	oldName := oldSa.GetName()
	oldNamespace := oldSa.GetNamespace()

	// We only care the name and namespace of a service account.
	if curName != oldName || curNamespace != oldNamespace {
		sc.deleteSecret(oldName, oldNamespace)
		sc.upsertSecret(curName, curNamespace)

		log.Infof("Service account \"%s\" in namespace \"%s\" has been updated to \"%s\" in namespace \"%s\"",
			oldName, oldNamespace, curName, curNamespace)
	}
}

func (sc *SecretController) upsertSecret(saName, saNamespace string) {
	secret := k8ssecret.BuildSecret(saName, GetSecretName(saName), saNamespace, nil,
		nil, nil, nil, nil, IstioSecretType)

	_, exists, err := sc.scrtStore.Get(secret)
	if err != nil {
		log.Errorf("Failed to get secret from the store (error %v)", err)
	}

	if exists {
		// Do nothing for existing secrets. Rotating expiring certs are handled by the `scrtUpdated` method.
		return
	}

	// Now we know the secret does not exist yet. So we create a new one.
	chain, key, err := sc.generateKeyAndCert(saName, saNamespace)
	if err != nil {
		log.Errorf("Failed to generate key and certificate for service account %q in namespace %q (error %v)",
			saName, saNamespace, err)

		return
	}
	rootCert := sc.ca.GetCAKeyCertBundle().GetRootCertPem()
	secret.Data = map[string][]byte{
		CertChainID:  chain,
		PrivateKeyID: key,
		RootCertID:   rootCert,
	}

	// We retry several times when create secret to mitigate transient network failures.
	for i := 0; i < secretCreationRetry; i++ {
		_, err = sc.core.Secrets(saNamespace).Create(secret)
		if err == nil || errors.IsAlreadyExists(err) {
			if errors.IsAlreadyExists(err) {
				log.Infof("Istio secret for service account \"%s\" in namespace \"%s\" already exists", saName, saNamespace)
			}
			break
		} else {
			log.Errorf("Failed to create secret in attempt %v/%v, (error: %s)", i+1, secretCreationRetry, err)
		}
		time.Sleep(time.Second)
	}

	if err != nil && !errors.IsAlreadyExists(err) {
		log.Errorf("Failed to create secret for service account \"%s\"  (error: %s), retries %v times",
			saName, err, secretCreationRetry)
		return
	}

	log.Infof("Istio secret for service account \"%s\" in namespace \"%s\" has been created", saName, saNamespace)
}

func (sc *SecretController) deleteSecret(saName, saNamespace string) {
	err := sc.core.Secrets(saNamespace).Delete(GetSecretName(saName), nil)
	// kube-apiserver returns NotFound error when the secret is successfully deleted.
	if err == nil || errors.IsNotFound(err) {
		log.Infof("Istio secret for service account \"%s\" in namespace \"%s\" has been deleted", saName, saNamespace)
		return
	}

	log.Errorf("Failed to delete Istio secret for service account \"%s\" in namespace \"%s\" (error: %s)",
		saName, saNamespace, err)
}

func (sc *SecretController) scrtDeleted(obj interface{}) {
	scrt, ok := obj.(*v1.Secret)
	if !ok {
		log.Warnf("Failed to convert to secret object: %v", obj)
		return
	}

	saName := scrt.Annotations[ServiceAccountNameAnnotationKey]
	if _, error := sc.core.ServiceAccounts(scrt.GetNamespace()).Get(saName, metav1.GetOptions{}); error == nil {
		log.Errorf("Re-create deleted Istio secret for existing service account %s.", saName)
		sc.upsertSecret(saName, scrt.GetNamespace())
		sc.monitoring.SecretDeletion.Inc()
	}
}

func (sc *SecretController) generateKeyAndCert(saName string, saNamespace string) ([]byte, []byte, error) {
	id := spiffe.MustGenSpiffeURI(saNamespace, saName)
	if sc.dnsNames != nil {
		// Control plane components in same namespace.
		if e, ok := sc.dnsNames[saName]; ok {
			if e.Namespace == saNamespace {
				// Example: istio-pilot.istio-system.svc, istio-pilot.istio-system
				id += "," + fmt.Sprintf("%s.%s.svc", e.ServiceName, e.Namespace)
				id += "," + fmt.Sprintf("%s.%s", e.ServiceName, e.Namespace)
			}
		}
		// Custom adds more DNS entries using CLI
		if e, ok := sc.dnsNames[saName+"."+saNamespace]; ok {
			for _, d := range e.CustomDomains {
				id += "," + d
			}
		}
	}

	options := util.CertOptions{
		Host:       id,
		RSAKeySize: keySize,
		IsDualUse:  sc.dualUse,
	}

	csrPEM, keyPEM, err := util.GenCSR(options)
	if err != nil {
		log.Errorf("CSR generation error (%v)", err)
		sc.monitoring.CSRError.Inc()
		return nil, nil, err
	}

	certChainPEM := sc.ca.GetCAKeyCertBundle().GetCertChainPem()
	certPEM, signErr := sc.ca.Sign(csrPEM, strings.Split(id, ","), sc.certTTL, sc.forCA)
	if signErr != nil {
		log.Errorf("CSR signing error (%v)", signErr.Error())
		sc.monitoring.GetCertSignError(signErr.(*caerror.Error).ErrorType()).Inc()
		return nil, nil, fmt.Errorf("CSR signing error (%v)", signErr.(*caerror.Error))
	}
	certPEM = append(certPEM, certChainPEM...)

	return certPEM, keyPEM, nil
}

func (sc *SecretController) scrtUpdated(oldObj, newObj interface{}) {
	scrt, ok := newObj.(*v1.Secret)
	if !ok {
		log.Warnf("Failed to convert to secret object: %v", newObj)
		return
	}
	namespace := scrt.GetNamespace()
	name := scrt.GetName()

	_, waitErr := sc.certUtil.GetWaitTime(scrt.Data[CertChainID], time.Now(), sc.minGracePeriod)

	caCert, _, _, rootCertificate := sc.ca.GetCAKeyCertBundle().GetAllPem()

	// Check if root certificate in key cert bundle is not up-to-date. With mutiple
	// Citadel deployed in Istio, and Citadels are in self signed mode, the root
	// certificate in istio-ca-secret could be rotated by any Citadel and become newer
	// than the one in local key cert bundle.
	// When Citadel is in plugged cert mode, the root cert in keycertbundle and root
	// cert in istio-ca-secret are always identical.
	if sc.lastKCBSyncTime.IsZero() || time.Since(sc.lastKCBSyncTime) > 30*time.Second {
		if !bytes.Equal(rootCertificate, scrt.Data[RootCertID]) {
			caSecret, scrtErr := sc.caSecretController.LoadCASecretWithRetry(CASecret,
				sc.istioCaStorageNamespace, 100*time.Millisecond, 5*time.Second)
			if scrtErr != nil {
				log.Errorf("Fail to load CA secret %s:%s (error: %s), skip updating secret %s",
					sc.istioCaStorageNamespace, CASecret, scrtErr.Error(), name)
				return
			}
			// The CA cert from istio-ca-secret is the source of truth. If CA cert
			// in local keycertbundle does not match the CA cert in istio-ca-secret,
			// reload root cert into keycertbundle.
			if !bytes.Equal(caCert, caSecret.Data[caCertID]) {
				log.Warn("CA cert in KeyCertBundle does not match CA cert in " +
					"istio-ca-secret. Start to reload root cert into KeyCertBundle")
				// In self signed cert mode, no root cert file is appended, the root cert and ca cert
				// are the same.
				rootCertificate = caSecret.Data[caCertID]
				if err := sc.ca.GetCAKeyCertBundle().VerifyAndSetAll(caSecret.Data[caCertID],
					caSecret.Data[caPrivateKeyID], nil, rootCertificate); err != nil {
					log.Errorf("failed to reload root cert into KeyCertBundle (%v)", err)
					return
				}
				log.Info("Successfully reloaded root cert into KeyCertBundle.")
				sc.lastKCBSyncTime = time.Now()
			} else {
				log.Info("CA cert in KeyCertBundle matches CA cert in " +
					"istio-ca-secret. Skip reloading root cert into KeyCertBundle")
			}
		}
	}

	// Refresh the secret if 1) the certificate contained in the secret is about
	// to expire, or 2) the root certificate in the secret is different than the
	// one held by the ca (this may happen when the CA is restarted and
	// a new self-signed CA cert is generated).
	if waitErr != nil || !bytes.Equal(rootCertificate, scrt.Data[RootCertID]) {
		if waitErr != nil {
			log.Infof("Refreshing about to expire secret %s/%s: %s", namespace, GetSecretName(name), waitErr.Error())
		} else {
			log.Infof("Refreshing secret %s/%s (outdated root cert)", namespace, GetSecretName(name))
		}

		if err := sc.refreshSecret(scrt); err != nil {
			log.Errorf("Failed to update secret %s/%s (error: %s)", namespace, name, err)
		}
	}
}

// refreshSecret is an inner func to refresh cert secrets when necessary
func (sc *SecretController) refreshSecret(scrt *v1.Secret) error {
	namespace := scrt.GetNamespace()
	saName := scrt.Annotations[ServiceAccountNameAnnotationKey]

	chain, key, err := sc.generateKeyAndCert(saName, namespace)
	if err != nil {
		return err
	}

	scrt.Data[CertChainID] = chain
	scrt.Data[PrivateKeyID] = key
	scrt.Data[RootCertID] = sc.ca.GetCAKeyCertBundle().GetRootCertPem()

	_, err = sc.core.Secrets(namespace).Update(scrt)
	return err
}
