package k8s

import (
	"k8s.io/client-go/dynamic"
	// Load Auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClientMgr create a new instance of a client manager
func NewClientMgr(kubeConfig string, configOverrides *clientcmd.ConfigOverrides) ClientMgr {
	return &clientMgr{
		kubeConfig,
		configOverrides,
	}
}

// ClientMgr a container for creating kubernetes rest clients
type ClientMgr interface {
	StaticClient() (*kubernetes.Clientset, error)
	DynamicClient() (dynamic.Interface, error)
	RestConfig() (*rest.Config, error)
}

type clientMgr struct {
	kubeConfig      string
	configOverrides *clientcmd.ConfigOverrides
}

func (cm *clientMgr) clientConfig() (*rest.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if cm.kubeConfig != "" {
		loadingRules.ExplicitPath = cm.kubeConfig
	}
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, cm.configOverrides)
	restConfig, err := config.ClientConfig()
	if err != nil {
		return &rest.Config{}, err
	}
	return restConfig, nil

}

// RestConfig generate a new rest config
func (cm *clientMgr) RestConfig() (*rest.Config, error) {
	return cm.clientConfig()
}

// StaticClient setup kubernetes client
func (cm *clientMgr) StaticClient() (*kubernetes.Clientset, error) {
	conf, err := cm.clientConfig()
	if err != nil {
		return &kubernetes.Clientset{}, err
	}
	return kubernetes.NewForConfig(conf)
}

// DynamicClient setup kubernetes client
func (cm *clientMgr) DynamicClient() (dynamic.Interface, error) {
	conf, err := cm.clientConfig()
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(conf)
}
