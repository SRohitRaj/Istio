// Copyright 2019 Istio Authors.
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

package multicluster

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/google/go-cmp/cmp"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/clientcmd/api"

	"istio.io/istio/pkg/kube/secretcontroller"
)

var secretTemplate = `# Remote credentials for cluster "{{ .Name }}"
apiVersion: v1
kind: Secret
metadata:
  creationTimestamp: null
  labels:
    {{ .MultiClusterSecretLabel }}: "true"
  testName: istio-remote-secret-{{ .Name }}
stringData:
  {{ .Name }}: |
    apiVersion: v1
    clusters:
    - cluster:
        certificate-authority-data: {{ .CADataBase64 }}
        server: {{ .Server }}
      testName: {{ .Name }}
    contexts:
    - context:
        cluster: {{ .Name }}
        user: {{ .Name }}
      testName: {{ .Name }}
    current-context: {{ .Name }}
    kind: Config
    preferences: {}
    users:
    - testName: {{ .Name }}
      user:
        token: {{ .Token }}
---
`

type testClusterData struct {
	// SA data
	Name                    string
	CAData                  string
	CADataBase64            string
	Token                   string
	MultiClusterSecretLabel string

	// Secret with SA encoded in Kubeconfig
	kubeconfigSecretYaml string

	Server string
}

func makeTestClusterData(name string) *testClusterData {
	caData := fmt.Sprintf("caData-%v", name)
	token := fmt.Sprintf("token-%v", name)
	d := &testClusterData{
		Name:                    name,
		CAData:                  caData,
		CADataBase64:            base64.StdEncoding.EncodeToString([]byte(caData)),
		Token:                   token,
		Server:                  fmt.Sprintf("server-%v", name),
		MultiClusterSecretLabel: secretcontroller.MultiClusterSecretLabel,
	}

	var out bytes.Buffer
	tmpl, err := template.New("secretTemplate").Parse(secretTemplate)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(&out, d); err != nil {
		panic(err)
	}

	d.kubeconfigSecretYaml = out.String()

	return d
}

func makeServiceAccount(name string, secrets ...string) *v1.ServiceAccount {
	sa := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	}

	for _, secret := range secrets {
		sa.Secrets = append(sa.Secrets, v1.ObjectReference{
			Name:      secret,
			Namespace: testNamespace,
		})
	}

	return sa
}

func makeSecret(name, caData, token string) *v1.Secret {
	out := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
		Data: map[string][]byte{},
	}
	if len(caData) > 0 {
		out.Data[v1.ServiceAccountRootCAKey] = []byte(caData)
	}
	if len(token) > 0 {
		out.Data[v1.ServiceAccountTokenKey] = []byte(token)
	}
	return out
}

var (
	testNamespace = "istio-system-test"
	c0            = makeTestClusterData("c0")
)

type fakeOutputWriter struct {
	b           bytes.Buffer
	injectError error
	failAfter   int
}

func (w *fakeOutputWriter) Write(p []byte) (n int, err error) {
	w.failAfter--
	if w.failAfter <= 0 && w.injectError != nil {
		return 0, w.injectError
	}
	return w.b.Write(p)
}
func (w *fakeOutputWriter) String() string { return w.b.String() }

func TestCreateRemoteSecrets(t *testing.T) {
	prevStartingConfig := newStartingConfig
	defer func() { newStartingConfig = prevStartingConfig }()

	prevKubernetesInteface := newKubernetesInterface
	defer func() { newKubernetesInterface = prevKubernetesInteface }()

	prevOutputWriterStub := makeOutputWriterTestHook
	defer func() { makeOutputWriterTestHook = prevOutputWriterStub }()

	sa := makeServiceAccount(DefaultServiceAccountName, "saSecret")
	saSecret := makeSecret("saSecret", "caData", "token")
	saSecretMissingToken := makeSecret("saSecret", "caData", "")

	wantOutput := `# This file is autogenerated, do not edit.
apiVersion: v1
kind: Secret
metadata:
  creationTimestamp: null
  labels:
    istio/multiCluster: "true"
  name: istio-remote-secret-cluster-foo
stringData:
  cluster-foo: |
    apiVersion: v1
    clusters:
    - cluster:
        certificate-authority-data: Y2FEYXRh
        server: server
      name: cluster-foo
    contexts:
    - context:
        cluster: cluster-foo
        user: cluster-foo
      name: cluster-foo
    current-context: cluster-foo
    kind: Config
    preferences: {}
    users:
    - name: cluster-foo
      user:
        token: token
---
`

	badStartingConfigErrStr := "bad starting config"
	testKubeconfig := "test-kubeconfig"
	testContext := "test-context"

	cases := []struct {
		testName string

		// test input
		config *api.Config
		objs   []runtime.Object
		name   string

		// inject errors
		badStartingConfig bool
		outputWriterError error

		want       string
		wantErrStr string
	}{
		{
			testName:   "fail to get service account secret token",
			objs:       []runtime.Object{sa},
			wantErrStr: fmt.Sprintf("secrets %q not found", saSecret.Name),
		},
		{
			testName:          "fail to create starting config",
			objs:              []runtime.Object{sa, saSecret},
			badStartingConfig: true,
			wantErrStr:        badStartingConfigErrStr,
		},
		{
			testName: "fail to find cluster server in local kubeconfig",
			objs:     []runtime.Object{sa, saSecret},
			config: &api.Config{
				CurrentContext: "current",
				Clusters:       map[string]*api.Cluster{ /* missing cluster */ },
			},
			wantErrStr: `could not find server for context "current"`,
		},
		{
			testName: "fail to create remote secret token",
			objs:     []runtime.Object{sa, saSecretMissingToken},
			config: &api.Config{
				CurrentContext: "current",
				Clusters: map[string]*api.Cluster{
					"current": {Server: "server"},
				},
			},
			wantErrStr: `no "token" data found`,
		},
		{
			testName: "fail to encode secret",
			objs:     []runtime.Object{sa, saSecret},
			config: &api.Config{
				CurrentContext: "current",
				Clusters: map[string]*api.Cluster{
					"current": {Server: "server"},
				},
			},
			outputWriterError: errors.New("injected encode error"),
			wantErrStr:        "injected encode error",
		},
		{
			testName: "success",
			objs:     []runtime.Object{sa, saSecret},
			config: &api.Config{
				CurrentContext: "current",
				Clusters: map[string]*api.Cluster{
					"current": {Server: "server"},
				},
			},
			name: "cluster-foo",
			want: wantOutput,
		},
	}

	for i := range cases {
		c := &cases[i]
		t.Run(fmt.Sprintf("[%v] %v", i, c.testName), func(tt *testing.T) {
			newStartingConfig = func(kubeconfig, context string) (*api.Config, error) {
				if kubeconfig != testKubeconfig {
					t.Fatalf("newStartingConfig invoked with wrong kubeconfig: got %v want %v",
						kubeconfig, testKubeconfig)
				}
				if context != testContext {
					t.Fatalf("newStartingConfig invoked with wrong context: got %v want %v",
						context, testContext)
				}
				if c.badStartingConfig {
					return nil, errors.New(badStartingConfigErrStr)
				}
				return c.config, nil
			}

			newKubernetesInterface = func(kubeconfig, context string) (kubernetes.Interface, error) {
				if kubeconfig != testKubeconfig {
					t.Fatalf("newKubernetesInterface invoked with wrong kubeconfig: got %v want %v",
						kubeconfig, testKubeconfig)
				}
				if context != testContext {
					t.Fatalf("newKubernetesInterface invoked invoked with wrong context: got %v want %v",
						context, testContext)
				}
				return fake.NewSimpleClientset(c.objs...), nil
			}
			makeOutputWriterTestHook = func() writer {
				return &fakeOutputWriter{injectError: c.outputWriterError}
			}

			got, err := CreateRemoteSecret(testKubeconfig, testContext, testNamespace, DefaultServiceAccountName, c.name)
			if c.wantErrStr != "" {
				if err == nil {
					tt.Fatalf("wanted error including %q but got none", c.wantErrStr)
				} else if !strings.Contains(err.Error(), c.wantErrStr) {
					tt.Fatalf("wanted error including %q but got %v", c.wantErrStr, err)
				}
			} else if c.wantErrStr == "" && err != nil {
				tt.Fatalf("wanted non-error but got %q", err)
			} else if diff := cmp.Diff(got, c.want); diff != "" {
				tt.Errorf("got\n%v\nwant\n%vdiff %v", got, c.want, diff)
			}
		})
	}
}

func TestGetServiceAccountSecretToken(t *testing.T) {
	secret := makeSecret("secret", "caData", "token")

	cases := []struct {
		name string

		saNamespace string
		saName      string
		objs        []runtime.Object

		want       *v1.Secret
		wantErrStr string
	}{
		{
			name:        "missing service account",
			saName:      DefaultServiceAccountName,
			saNamespace: testNamespace,
			wantErrStr:  fmt.Sprintf("serviceaccounts %q not found", DefaultServiceAccountName),
		},
		{
			name:        "wrong number of secrets",
			saName:      DefaultServiceAccountName,
			saNamespace: testNamespace,
			objs: []runtime.Object{
				makeServiceAccount(DefaultServiceAccountName, "secret", "extra-secret"),
			},
			wantErrStr: "wrong number of secrets",
		},
		{
			name:        "missing service account token secret",
			saName:      DefaultServiceAccountName,
			saNamespace: testNamespace,
			objs: []runtime.Object{
				makeServiceAccount(DefaultServiceAccountName, "wrong-secret"),
				secret,
			},
			wantErrStr: `secrets "wrong-secret" not found`,
		},
		{
			name:        "success",
			saName:      DefaultServiceAccountName,
			saNamespace: testNamespace,
			objs: []runtime.Object{
				makeServiceAccount(DefaultServiceAccountName, "secret"),
				secret,
			},
			want: secret,
		},
	}

	for i := range cases {
		c := &cases[i]
		t.Run(fmt.Sprintf("[%v] %v", i, c.name), func(tt *testing.T) {
			kube := fake.NewSimpleClientset(c.objs...)

			got, err := getServiceAccountSecretToken(kube, c.saName, c.saNamespace)
			if c.wantErrStr != "" {
				if err == nil {
					tt.Fatalf("wanted error including %q but got none", c.wantErrStr)
				} else if !strings.Contains(err.Error(), c.wantErrStr) {
					tt.Fatalf("wanted error including %q but got %v", c.wantErrStr, err)
				}
			} else if c.wantErrStr == "" && err != nil {
				tt.Fatalf("wanted non-error but got %q", err)
			} else if diff := cmp.Diff(got, c.want); diff != "" {
				tt.Errorf("got\n%v\nwant\n%vdiff %v", got, c.want, diff)
			}
		})
	}
}

func TestGetClusterServerFromKubeconfig(t *testing.T) {
	prev := newStartingConfig
	defer func() { newStartingConfig = prev }()

	wantServer := "server0"
	context := "context0"

	cases := []struct {
		name              string
		badStartingConfig bool
		config            *api.Config
		wantErrStr        string
	}{
		{
			name:              "bad starting config",
			badStartingConfig: true,
			config: &api.Config{
				CurrentContext: context,
				Clusters: map[string]*api.Cluster{
					context: {Server: wantServer},
				},
			},
			wantErrStr: "bad starting config",
		},
		{
			name: "missing server",
			config: &api.Config{
				CurrentContext: context,
				Clusters:       map[string]*api.Cluster{},
			},
			wantErrStr: "could not find server for context",
		},
		{
			name: "success",
			config: &api.Config{
				CurrentContext: context,
				Clusters: map[string]*api.Cluster{
					context: {Server: wantServer},
				},
			},
		},
	}

	for i := range cases {
		c := &cases[i]
		t.Run(fmt.Sprintf("[%v] %v", i, c.name), func(tt *testing.T) {
			newStartingConfig = func(_, _ string) (*api.Config, error) {
				if c.badStartingConfig {
					return nil, errors.New("bad starting config")
				}
				return c.config, nil
			}

			gotServer, err := getClusterServerFromKubeconfig("foo", "bar")

			if c.wantErrStr != "" {
				if err == nil {
					tt.Fatalf("wanted error including %q but got none", c.wantErrStr)
				} else if !strings.Contains(err.Error(), c.wantErrStr) {
					tt.Fatalf("wanted error including %q but got %v", c.wantErrStr, err)
				}
			} else if c.wantErrStr == "" && err != nil {
				tt.Fatalf("wanted non-error but got %q", err)
			} else if gotServer != "server0" {
				t.Fatalf("got server %v want %v", gotServer, wantServer)
			}
		})
	}
}

func TestCreateRemoteKubeconfig(t *testing.T) {
	kubeconfig := `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: Y2FEYXRhLWMw
    server: ""
  name: c0
contexts:
- context:
    cluster: c0
    user: c0
  name: c0
current-context: c0
kind: Config
preferences: {}
users:
- name: c0
  user:
    token: token-c0
`

	cases := []struct {
		name        string
		in          *v1.Secret
		clusterName string
		server      string
		want        *v1.Secret
		wantErrStr  string
		context     string
	}{
		{
			name:       "missing caData",
			in:         makeSecret("", "", c0.Token),
			context:    "c0",
			wantErrStr: errMissingRootCAKey.Error(),
		},
		{
			name:       "missing token",
			in:         makeSecret("", c0.CAData, ""),
			context:    "c0",
			wantErrStr: errMissingTokenKey.Error(),
		},
		{
			name:    "success",
			in:      makeSecret("", c0.CAData, c0.Token),
			context: "c0",
			want: &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name: "istio-remote-secret-c0",
					Labels: map[string]string{
						secretcontroller.MultiClusterSecretLabel: "true",
					},
				},
				StringData: map[string]string{
					"c0": kubeconfig,
				},
			},
		},
	}

	for i := range cases {
		c := &cases[i]
		t.Run(fmt.Sprintf("[%v] %v", i, c.name), func(tt *testing.T) {
			got, err := createRemoteSecretFromTokenAndServer(c.in, c.context, c.server)
			if c.wantErrStr != "" {
				if err == nil {
					tt.Fatalf("wanted error including %q but none", c.wantErrStr)
				} else if !strings.Contains(err.Error(), c.wantErrStr) {
					tt.Fatalf("wanted error including %q but %v", c.wantErrStr, err)
				}
			} else if c.wantErrStr == "" && err != nil {
				tt.Fatalf("wanted non-error but got %q", err)
			} else if diff := cmp.Diff(got, c.want); diff != "" {
				tt.Fatalf(" got %v\nwant %v\ndiff %v", got, c.want, diff)
			}
		})
	}
}

func TestWriteEncodedSecret(t *testing.T) {
	s := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
		},
	}

	w := &fakeOutputWriter{failAfter: 0, injectError: errors.New("error")}
	if err := writeEncodedSecret(w, s); err == nil {
		t.Error("want error on first write failure")
	}

	w = &fakeOutputWriter{failAfter: 1, injectError: errors.New("error")}
	if err := writeEncodedSecret(w, s); err == nil {
		t.Error("want error on second write failure")
	}

	w = &fakeOutputWriter{failAfter: 2, injectError: errors.New("error")}
	if err := writeEncodedSecret(w, s); err == nil {
		t.Error("want error on third write failure")
	}

	w = &fakeOutputWriter{}
	if err := writeEncodedSecret(w, s); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	want := `# This file is autogenerated, do not edit.
apiVersion: v1
kind: Secret
metadata:
  creationTimestamp: null
  name: foo
---
`
	if w.String() != want {
		t.Errorf("got\n%q\nwant\n%q", w.String(), want)
	}

}
