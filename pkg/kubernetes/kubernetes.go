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

func defaultClientConfig() clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
}

func RestClientConfig() (*rest.Config, error) {
	kubeCfg := defaultClientConfig()

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
