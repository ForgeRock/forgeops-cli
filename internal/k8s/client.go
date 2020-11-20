package k8s

import (

	// Load Auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// GetClient setup kubernetes client
func GetClient(kubeConfig string, configOverrides *clientcmd.ConfigOverrides) (*kubernetes.Clientset, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeConfig != "" {
		loadingRules.ExplicitPath = kubeConfig
	}
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	restConfig, err := config.ClientConfig()
	if err != nil {
		return &kubernetes.Clientset{}, err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return &kubernetes.Clientset{}, err
	}

	return clientset, nil
}
