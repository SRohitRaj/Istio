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

package util

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/sync/errgroup"

	"istio.io/istio/pkg/log"
)

const (
	podRunning   = "Running"
	podFailedGet = "Failed_Get"
	// The index of STATUS field in kubectl CLI output.
	statusField = 2
)

// Fill complete a template with given values and generate a new output file
func Fill(outFile, inFile string, values interface{}) error {
	var bytes bytes.Buffer
	w := bufio.NewWriter(&bytes)
	tmpl, err := template.ParseFiles(inFile)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(w, values); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}

	if err := ioutil.WriteFile(outFile, bytes.Bytes(), 0644); err != nil {
		return err
	}
	log.Infof("Created %s from template %s", outFile, inFile)
	return nil
}

// CreateNamespace create a kubernetes namespace
func CreateNamespace(n string) error {
	if _, err := Shell("kubectl create namespace %s", n); err != nil {
		return err
	}
	log.Infof("namespace %s created\n", n)
	return nil
}

// DeleteNamespace delete a kubernetes namespace
func DeleteNamespace(n string) error {
	_, err := Shell("kubectl delete namespace %s", n)
	return err
}

// NamespaceDeleted check if a kubernete namespace is deleted
func NamespaceDeleted(n string) (bool, error) {
	output, err := ShellSilent("kubectl get namespace %s -o name", n)
	if strings.Contains(output, "NotFound") {
		return true, nil
	}
	return false, err
}

// KubeApplyContents kubectl apply from contents
func KubeApplyContents(namespace, yamlContents string) error {
	tmpfile, err := WriteTempfile(os.TempDir(), "kubeapply", ".yaml", yamlContents)
	if err != nil {
		return err
	}
	defer removeFile(tmpfile)
	return KubeApply(namespace, tmpfile)
}

// KubeApply kubectl apply from file
func KubeApply(namespace, yamlFileName string) error {
	_, err := Shell("kubectl apply -n %s -f %s", namespace, yamlFileName)
	return err
}

// KubeDeleteContents kubectl apply from contents
func KubeDeleteContents(namespace, yamlContents string) error {
	tmpfile, err := WriteTempfile(os.TempDir(), "kubedelete", ".yaml", yamlContents)
	if err != nil {
		return err
	}
	defer removeFile(tmpfile)
	return KubeDelete(namespace, tmpfile)
}

func removeFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Errorf("Unable to remove %s: %v", path, err)
	}
}

// KubeDelete kubectl delete from file
func KubeDelete(namespace, yamlFileName string) error {
	_, err := Shell("kubectl delete -n %s -f %s", namespace, yamlFileName)
	return err
}

// GetIngress get istio ingress ip
func GetIngress(n string) (string, error) {
	retry := Retrier{
		BaseDelay: 1 * time.Second,
		MaxDelay:  1 * time.Second,
		Retries:   300, // ~5 minutes
	}
	ri := regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)
	//rp := regexp.MustCompile(`^[0-9]{1,5}$`) # Uncomment for minikube
	var ingress string
	retryFn := func(i int) error {
		ip, err := ShellSilent("kubectl get svc istio-ingress -n %s -o jsonpath='{.status.loadBalancer.ingress[*].ip}'", n)
		// For minikube, comment out the previous line and uncomment the following line
		//ip, err := Shell("kubectl get po -l istio=ingress -n %s -o jsonpath='{.items[0].status.hostIP}'", n)
		if err != nil {
			return err
		}
		ip = strings.Trim(ip, "'")
		if ri.FindString(ip) == "" {
			return errors.New("ingress ip not available yet")
		}
		ingress = ip
		// For minikube, comment out the previous line and uncomment the following lines
		//port, e := Shell("kubectl get svc istio-ingress -n %s -o jsonpath='{.spec.ports[0].nodePort}'", n)
		//if e != nil {
		//	return e
		//}
		//port = strings.Trim(port, "'")
		//if rp.FindString(port) == "" {
		//	err = fmt.Errorf("unable to find ingress port")
		//	log.Warn(err)
		//	return err
		//}
		//ingress = ip + ":" + port
		log.Infof("Istio ingress: %s", ingress)

		return nil
	}
	log.Info("Waiting for istio-ingress to get external IP")
	if _, err := retry.Retry(retryFn); err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	client := &http.Client{Timeout: 5 * time.Second}
	for {
		select {
		case <-ctx.Done():
			return "", errors.New("istio-ingress readiness check timed out")
		default:
			url := fmt.Sprintf("http://%s", ingress)
			log.Infof("Checking %v", url)
			response, err := ctxhttp.Get(ctx, client, url)
			if err == nil {
				log.Infof("Response %v %q received from istio-ingress", response.StatusCode, response.Status)
				return ingress, nil
			}
			log.Warnf("Error: %v, checking again", err)
		}
	}
}

// GetIngressPod get istio ingress ip
func GetIngressPod(n string) (string, error) {
	retry := Retrier{
		BaseDelay: 5 * time.Second,
		MaxDelay:  5 * time.Minute,
		Retries:   20,
	}
	ipRegex := regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)
	portRegex := regexp.MustCompile(`^[0-9]+$`)
	var ingress string
	retryFn := func(i int) error {
		podIP, err := Shell("kubectl get pod -l istio=ingress "+
			"-n %s -o jsonpath='{.items[0].status.hostIP}'", n)
		if err != nil {
			return err
		}
		podPort, err := Shell("kubectl get svc istio-ingress "+
			"-n %s -o jsonpath='{.spec.ports[0].nodePort}'", n)
		if err != nil {
			return err
		}
		podIP = strings.Trim(podIP, "'")
		podPort = strings.Trim(podPort, "'")
		if ipRegex.FindString(podIP) == "" {
			err = errors.New("unable to find ingress pod ip")
			log.Warna(err)
			return err
		}
		if portRegex.FindString(podPort) == "" {
			err = errors.New("unable to find ingress pod port")
			log.Warna(err)
			return err
		}
		ingress = fmt.Sprintf("%s:%s", podIP, podPort)
		log.Infof("Istio ingress: %s\n", ingress)
		return nil
	}
	_, err := retry.Retry(retryFn)
	return ingress, err
}

// GetPodsName gets names of all pods in specific namespace and return in a slice
func GetPodsName(n string) (pods []string) {
	res, err := Shell("kubectl -n %s get pods -o jsonpath='{.items[*].metadata.name}'", n)
	if err != nil {
		log.Infof("Failed to get pods name in namespace %s: %s", n, err)
		return
	}
	res = strings.Trim(res, "'")
	pods = strings.Split(res, " ")
	log.Infof("Existing pods: %v", pods)
	return
}

// GetPodStatus gets status of a pod from a namespace
// Note: It is not enough to check pod phase, which only implies there is at
// least one container running. Use kubectl CLI to get status so that we can
// ensure that all containers are running.
func GetPodStatus(n, pod string) string {
	status, err := Shell("kubectl -n %s get pods %s --no-headers", n, pod)
	if err != nil {
		log.Infof("Failed to get status of pod %s in namespace %s: %s", pod, n, err)
		status = podFailedGet
	}
	f := strings.Fields(status)
	if len(f) > statusField {
		return f[statusField]
	}
	return ""
}

// CheckPodsRunning return if all pods in a namespace are in "Running" status
// Also check container status to be running.
func CheckPodsRunning(n string) (ready bool) {
	retry := Retrier{
		BaseDelay: 30 * time.Second,
		MaxDelay:  30 * time.Second,
		Retries:   6,
	}

	retryFn := func(i int) error {
		pods := GetPodsName(n)
		ready = true
		for _, p := range pods {
			if status := GetPodStatus(n, p); status != podRunning {
				log.Infof("%s in namespace %s is not running: %s", p, n, status)
				if desc, err := ShellMuteOutput("kubectl describe pods -n %s %s", n, p); err != nil {
					log.Infof("Pod description: %s", desc)
				}
				ready = false
			}
		}
		if !ready {
			_, err := Shell("kubectl -n %s get pods -o wide", n)
			if err != nil {
				log.Infof("Cannot get pods: %s", err)
			}
			return fmt.Errorf("some pods are not ready")
		}
		return nil
	}
	_, err := retry.Retry(retryFn)
	if err != nil {
		return false
	}
	log.Info("Get all pods running!")
	return true
}

// CheckDeployment gets status of a deployment from a namespace
func CheckDeployment(ctx context.Context, namespace, deployment string) error {
	errc := make(chan error)
	go func() {
		if _, err := ShellMuteOutput("kubectl -n %s rollout status %s", namespace, deployment); err != nil {
			errc <- fmt.Errorf("%s in namespace %s failed", deployment, namespace)
		}
		errc <- nil
	}()
	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// CheckDeployments checks whether all deployment in a given namespace
func CheckDeployments(namespace string, timeout time.Duration) error {
	// wait for istio-system deployments to be fully rolled out before proceeding
	deployments, err := ShellMuteOutput("kubectl -n %s get deployment -o name", namespace)
	if err != nil {
		return fmt.Errorf("could not list deployments in namespace %q", namespace)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	for _, deployment := range strings.Fields(deployments) {
		g.Go(func() error { return CheckDeployment(ctx, namespace, deployment) })
	}
	return g.Wait()
}
