/*
Copyright 2017 The Kubernetes Authors.

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

package process

import (
	"os/exec"
	"syscall"

	"k8s.io/klog/v2"
)

// IsRespawnIfRequired checks if error type is exec.ExitError or not
func IsRespawnIfRequired(err error) bool {
	exitError, ok := err.(*exec.ExitError)
	if !ok {
		return false
	}

	waitStatus := exitError.Sys().(syscall.WaitStatus)
	klog.Warningf(`
-------------------------------------------------------------------------------
NGINX master process died (%v): %v
-------------------------------------------------------------------------------
`, waitStatus.ExitStatus(), err)
	return true
}
