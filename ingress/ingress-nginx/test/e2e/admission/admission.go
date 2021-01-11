/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package admission

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/ingress-nginx/test/e2e/framework"
)

var _ = framework.IngressNginxDescribe("[Serial] admission controller", func() {
	f := framework.NewDefaultFramework("admission")

	ginkgo.BeforeEach(func() {
		f.NewEchoDeployment()
		f.NewSlowEchoDeployment()
	})

	ginkgo.AfterEach(func() {
		err := uninstallChart(f)
		assert.Nil(ginkgo.GinkgoT(), err, "uninstalling helm chart")
	})

	ginkgo.It("reject ingress with global-rate-limit annotations when memcached is not configured", func() {
		host := "admission-test"

		annotations := map[string]string{
			"nginx.ingress.kubernetes.io/global-rate-limit":        "100",
			"nginx.ingress.kubernetes.io/global-rate-limit-window": "1m",
		}
		ing := framework.NewSingleIngress("first-ingress", "/", host, f.Namespace, framework.EchoService, 80, annotations)

		ginkgo.By("rejects ingress when memcached is not configured")

		_, err := f.KubeClientSet.NetworkingV1beta1().Ingresses(f.Namespace).Create(context.TODO(), ing, metav1.CreateOptions{})
		assert.NotNil(ginkgo.GinkgoT(), err, "creating ingress with global throttle annotations when memcached is not configured")

		ginkgo.By("accepts ingress when memcached is not configured")

		f.UpdateNginxConfigMapData("global-rate-limit-memcached-host", "memc.default.svc.cluster.local")

		_, err = f.KubeClientSet.NetworkingV1beta1().Ingresses(f.Namespace).Create(context.TODO(), ing, metav1.CreateOptions{})
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress with global throttle annotations when memcached is configured")

		f.WaitForNginxServer(host,
			func(server string) bool {
				return strings.Contains(server, fmt.Sprintf("server_name %v", host))
			})
	})

	ginkgo.It("should not allow overlaps of host and paths without canary annotations", func() {
		host := "admission-test"

		firstIngress := framework.NewSingleIngress("first-ingress", "/", host, f.Namespace, framework.EchoService, 80, nil)
		_, err := f.KubeClientSet.NetworkingV1beta1().Ingresses(f.Namespace).Create(context.TODO(), firstIngress, metav1.CreateOptions{})
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")

		f.WaitForNginxServer(host,
			func(server string) bool {
				return strings.Contains(server, fmt.Sprintf("server_name %v", host))
			})

		secondIngress := framework.NewSingleIngress("second-ingress", "/", host, f.Namespace, framework.EchoService, 80, nil)
		_, err = f.KubeClientSet.NetworkingV1beta1().Ingresses(f.Namespace).Create(context.TODO(), secondIngress, metav1.CreateOptions{})
		assert.NotNil(ginkgo.GinkgoT(), err, "creating an ingress with the same host and path should return an error")
	})

	ginkgo.It("should allow overlaps of host and paths with canary annotation", func() {
		host := "admission-test"

		firstIngress := framework.NewSingleIngress("first-ingress", "/", host, f.Namespace, framework.EchoService, 80, nil)
		_, err := f.KubeClientSet.NetworkingV1beta1().Ingresses(f.Namespace).Create(context.TODO(), firstIngress, metav1.CreateOptions{})
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")

		f.WaitForNginxServer(host,
			func(server string) bool {
				return strings.Contains(server, fmt.Sprintf("server_name %v", host))
			})

		canaryAnnotations := map[string]string{
			"nginx.ingress.kubernetes.io/canary":           "true",
			"nginx.ingress.kubernetes.io/canary-by-header": "CanaryByHeader",
		}
		secondIngress := framework.NewSingleIngress("second-ingress", "/", host, f.Namespace, framework.SlowEchoService, 80, canaryAnnotations)
		_, err = f.KubeClientSet.NetworkingV1beta1().Ingresses(f.Namespace).Create(context.TODO(), secondIngress, metav1.CreateOptions{})
		assert.Nil(ginkgo.GinkgoT(), err, "creating an ingress with the same host and path should not return an error using a canary annotation")
	})

	ginkgo.It("should return an error if there is an error validating the ingress definition", func() {
		host := "admission-test"

		annotations := map[string]string{
			"nginx.ingress.kubernetes.io/configuration-snippet": "something invalid",
		}
		firstIngress := framework.NewSingleIngress("first-ingress", "/", host, f.Namespace, framework.EchoService, 80, annotations)
		_, err := f.KubeClientSet.NetworkingV1beta1().Ingresses(f.Namespace).Create(context.TODO(), firstIngress, metav1.CreateOptions{})
		assert.NotNil(ginkgo.GinkgoT(), err, "creating an ingress with invalid configuration should return an error")
	})

	ginkgo.It("should not return an error the ingress definition uses the deprecated extensions package", func() {
		err := createIngress(f.Namespace, validIngress)
		assert.Nil(ginkgo.GinkgoT(), err, "creating an ingress using kubectl")

		f.WaitForNginxConfiguration(func(cfg string) bool {
			return strings.Contains(cfg, "extensions")
		})

		f.HTTPTestClient().
			GET("/").
			WithHeader("Host", "extensions").
			Expect().
			Status(http.StatusOK)
	})

	ginkgo.It("should return an error if the ingress definition uses the deprecated extensions package and invalid annotations", func() {
		err := createIngress(f.Namespace, invalidIngress)
		assert.NotNil(ginkgo.GinkgoT(), err, "creating an ingress using kubectl")

		_, err = f.KubeClientSet.NetworkingV1beta1().Ingresses(f.Namespace).Get(context.TODO(), "extensions", metav1.GetOptions{})
		if !apierrors.IsNotFound(err) {
			assert.NotNil(ginkgo.GinkgoT(), err, "creating an ingress with invalid configuration should return an error")
		}
	})

	ginkgo.It("should not return an error if the Ingress V1 definition is valid", func() {
		if !f.IsIngressV1Ready {
			ginkgo.Skip("Test requires Kubernetes v1.19 or higher")
		}

		err := createIngress(f.Namespace, validV1Ingress)
		assert.Nil(ginkgo.GinkgoT(), err, "creating an ingress using kubectl")

		f.WaitForNginxConfiguration(func(cfg string) bool {
			return strings.Contains(cfg, "extensions")
		})

		f.HTTPTestClient().
			GET("/").
			WithHeader("Host", "extensions").
			Expect().
			Status(http.StatusOK)
	})

	ginkgo.It("should return an error if the Ingress V1 definition contains invalid annotations", func() {
		if !f.IsIngressV1Ready {
			ginkgo.Skip("Test requires Kubernetes v1.19 or higher")
		}

		err := createIngress(f.Namespace, invalidV1Ingress)
		assert.NotNil(ginkgo.GinkgoT(), err, "creating an ingress using kubectl")

		_, err = f.KubeClientSet.NetworkingV1beta1().Ingresses(f.Namespace).Get(context.TODO(), "extensions", metav1.GetOptions{})
		if !apierrors.IsNotFound(err) {
			assert.NotNil(ginkgo.GinkgoT(), err, "creating an ingress with invalid configuration should return an error")
		}
	})
})

func uninstallChart(f *framework.Framework) error {
	cmd := exec.Command("helm", "uninstall", "--namespace", f.Namespace, "nginx-ingress")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unexpected error uninstalling ingress-nginx release: %v", err)
	}

	return nil
}

const (
	validIngress = `
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: extensions
spec:
  rules:
  - host: extensions
    http:
      paths:
      - path: /
        backend:
          serviceName: echo
          servicePort: 80
---
`

	invalidIngress = `
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: extensions
  annotations:
    nginx.ingress.kubernetes.io/configuration-snippet: |
      invalid directive
spec:
  rules:
  - host: extensions
    http:
      paths:
      - path: /
        backend:
          serviceName: echo
          servicePort: 80
---
`

	validV1Ingress = `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: extensions
spec:
  rules:
  - host: extensions
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: echo
            port:
              number: 80

---
`

	invalidV1Ingress = `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: extensions
  annotations:
    nginx.ingress.kubernetes.io/configuration-snippet: |
      invalid directive
spec:
  rules:
  - host: extensions
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: echo
            port:
              number: 80
---
`
)

func createIngress(namespace, ingressDefinition string) error {
	var (
		execErr bytes.Buffer
	)

	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("%v --warnings-as-errors=false apply --namespace %s -f -", framework.KubectlPath, namespace))
	cmd.Stdin = strings.NewReader(ingressDefinition)
	cmd.Stderr = &execErr

	err := cmd.Run()
	if err != nil {
		stderr := strings.TrimSpace(execErr.String())
		return fmt.Errorf("Kubectl error: %v\n%v", err, stderr)
	}

	return nil
}
