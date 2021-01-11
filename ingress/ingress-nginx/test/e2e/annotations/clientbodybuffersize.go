/*
Copyright 2018 The Kubernetes Authors.

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

package annotations

import (
	"fmt"
	"strings"

	"github.com/onsi/ginkgo"

	"k8s.io/ingress-nginx/test/e2e/framework"
)

var _ = framework.DescribeAnnotation("client-body-buffer-size", func() {
	f := framework.NewDefaultFramework("clientbodybuffersize")

	ginkgo.BeforeEach(func() {
		f.NewEchoDeployment()
	})

	ginkgo.It("should set client_body_buffer_size to 1000", func() {
		host := "client-body-buffer-size.com"

		clientBodyBufferSize := "1000"
		annotations := make(map[string]string)
		annotations["nginx.ingress.kubernetes.io/client-body-buffer-size"] = clientBodyBufferSize

		ing := framework.NewSingleIngress(host, "/", host, f.Namespace, framework.EchoService, 80, annotations)
		f.EnsureIngress(ing)

		f.WaitForNginxServer(host,
			func(server string) bool {
				return strings.Contains(server, fmt.Sprintf("client_body_buffer_size %s;", clientBodyBufferSize))
			})
	})

	ginkgo.It("should set client_body_buffer_size to 1K", func() {
		host := "client-body-buffer-size.com"

		clientBodyBufferSize := "1K"
		annotations := make(map[string]string)
		annotations["nginx.ingress.kubernetes.io/client-body-buffer-size"] = clientBodyBufferSize

		ing := framework.NewSingleIngress(host, "/", host, f.Namespace, framework.EchoService, 80, annotations)
		f.EnsureIngress(ing)

		f.WaitForNginxServer(host,
			func(server string) bool {
				return strings.Contains(server, fmt.Sprintf("client_body_buffer_size %s;", clientBodyBufferSize))
			})
	})

	ginkgo.It("should set client_body_buffer_size to 1k", func() {
		host := "client-body-buffer-size.com"

		clientBodyBufferSize := "1k"
		annotations := make(map[string]string)
		annotations["nginx.ingress.kubernetes.io/client-body-buffer-size"] = clientBodyBufferSize

		ing := framework.NewSingleIngress(host, "/", host, f.Namespace, framework.EchoService, 80, annotations)
		f.EnsureIngress(ing)

		f.WaitForNginxServer(host,
			func(server string) bool {
				return strings.Contains(server, fmt.Sprintf("client_body_buffer_size %s;", clientBodyBufferSize))
			})
	})

	ginkgo.It("should set client_body_buffer_size to 1m", func() {
		host := "client-body-buffer-size.com"

		clientBodyBufferSize := "1m"
		annotations := make(map[string]string)
		annotations["nginx.ingress.kubernetes.io/client-body-buffer-size"] = clientBodyBufferSize

		ing := framework.NewSingleIngress(host, "/", host, f.Namespace, framework.EchoService, 80, annotations)
		f.EnsureIngress(ing)

		f.WaitForNginxServer(host,
			func(server string) bool {
				return strings.Contains(server, fmt.Sprintf("client_body_buffer_size %s;", clientBodyBufferSize))
			})
	})

	ginkgo.It("should set client_body_buffer_size to 1M", func() {
		host := "client-body-buffer-size.com"

		clientBodyBufferSize := "1M"
		annotations := make(map[string]string)
		annotations["nginx.ingress.kubernetes.io/client-body-buffer-size"] = clientBodyBufferSize

		ing := framework.NewSingleIngress(host, "/", host, f.Namespace, framework.EchoService, 80, annotations)
		f.EnsureIngress(ing)

		f.WaitForNginxServer(host,
			func(server string) bool {
				return strings.Contains(server, fmt.Sprintf("client_body_buffer_size %s;", clientBodyBufferSize))
			})
	})

	ginkgo.It("should not set client_body_buffer_size to invalid 1b", func() {
		host := "client-body-buffer-size.com"

		clientBodyBufferSize := "1b"
		annotations := make(map[string]string)
		annotations["nginx.ingress.kubernetes.io/client-body-buffer-size"] = clientBodyBufferSize

		ing := framework.NewSingleIngress(host, "/", host, f.Namespace, framework.EchoService, 80, annotations)
		f.EnsureIngress(ing)

		f.WaitForNginxServer(host,
			func(server string) bool {
				return !strings.Contains(server, fmt.Sprintf("client_body_buffer_size %s;", clientBodyBufferSize))
			})
	})
})
