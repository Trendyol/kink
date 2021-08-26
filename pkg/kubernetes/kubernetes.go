/*
Copyright © 2021 pe.container <pe.container@trendyol.com>

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

package kubernetes

import (
	"fmt"

	"k8s.io/client-go/kubernetes"

	// Initialize all known client auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func ManagedAnnotations() map[string]string {
	return map[string]string{
		"a8r.io/owner":      "@kink",
		"a8r.io/repository": "https://gitlab.trendyol.com:platform/base/poc/kink.git",
	}
}

func DefaultClientConfig() clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
}

func RestClientConfig() (*rest.Config, error) {
	kubeCfg := DefaultClientConfig()

	restConfig, err := kubeCfg.ClientConfig()
	if clientcmd.IsEmptyConfig(err) {
		restConfig, err := rest.InClusterConfig()
		if err != nil {
			return restConfig, fmt.Errorf("error creating REST client config in-cluster: %v", err)
		}

		return restConfig, nil
	}
	if err != nil {
		return restConfig, fmt.Errorf("error creating REST client config: %v", err)
	}

	return restConfig, nil
}

func Client() (kubernetes.Interface, error) {
	config, err := RestClientConfig()
	if err != nil {
		return nil, fmt.Errorf("getting client config for Kubernetes client: %w", err)
	}
	return kubernetes.NewForConfig(config)
}
