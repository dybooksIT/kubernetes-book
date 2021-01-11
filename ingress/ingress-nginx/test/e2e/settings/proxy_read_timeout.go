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

package settings

import (
	"fmt"
	"strings"

	"github.com/onsi/ginkgo"

	"k8s.io/ingress-nginx/test/e2e/framework"
)

var _ = framework.DescribeSetting("proxy-read-timeout", func() {
	f := framework.NewDefaultFramework("proxy")
	host := "proxy-read-timeout.com"

	ginkgo.BeforeEach(func() {
		f.NewEchoDeployment()
	})

	ginkgo.It("should set valid proxy read timeouts using configmap values", func() {
		ing := framework.NewSingleIngress(host, "/", host, f.Namespace, framework.EchoService, 80, nil)
		f.EnsureIngress(ing)

		proxyReadtimeout := "20"

		cm := make(map[string]string)
		cm["proxy-read-timeout"] = proxyReadtimeout
		f.SetNginxConfigMapData(cm)

		f.WaitForNginxServer(host,
			func(server string) bool {
				return strings.Contains(server, fmt.Sprintf("proxy_read_timeout %ss;", proxyReadtimeout))
			})
	})

	ginkgo.It("should not set invalid proxy read timeouts using configmap values", func() {
		ing := framework.NewSingleIngress(host, "/", host, f.Namespace, framework.EchoService, 80, nil)
		f.EnsureIngress(ing)

		proxyReadtimeout := "20k"

		cm := make(map[string]string)
		cm["proxy-read-timeout"] = proxyReadtimeout
		f.SetNginxConfigMapData(cm)

		f.WaitForNginxServer(host,
			func(server string) bool {
				return !strings.Contains(server, fmt.Sprintf("proxy_read_timeout %ss;", proxyReadtimeout))
			})
	})

})
