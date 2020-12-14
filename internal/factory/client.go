package factory

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// Load Auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// NewFactory create a new instance of a client manager
func NewFactory(kubeConfigFlags *genericclioptions.ConfigFlags) Factory {
	return &factory{
		kubeConfigFlags,
	}
}

// Factory a container for creating kubernetes rest clients
type Factory interface {
	StaticClient() (*kubernetes.Clientset, error)
	DynamicClient() (dynamic.Interface, error)
	RestConfig() (*rest.Config, error)
	GetOverrideFlags() (*genericclioptions.ConfigFlags, error)
	Builder() *resource.Builder
}

type factory struct {
	kubeConfigFlags *genericclioptions.ConfigFlags
}

// GetOverrideFlags returns the command flags
func (f *factory) GetOverrideFlags() (*genericclioptions.ConfigFlags, error) {
	return f.kubeConfigFlags, nil
}

// RestConfig generate a new rest config
func (f *factory) RestConfig() (*rest.Config, error) {
	a, err := f.kubeConfigFlags.ToRESTConfig()
	return a, err
}

// StaticClient setup kubernetes client
func (f *factory) StaticClient() (*kubernetes.Clientset, error) {
	conf, err := f.RestConfig()
	if err != nil {
		return &kubernetes.Clientset{}, err
	}
	return kubernetes.NewForConfig(conf)
}

// DynamicClient setup kubernetes client
func (f *factory) DynamicClient() (dynamic.Interface, error) {
	conf, err := f.RestConfig()
	if err != nil {
		return nil, err
	}
	return dynamic.NewForConfig(conf)
}

// NewBuilder creates a new builder
func (f *factory) Builder() *resource.Builder {
	return resource.NewBuilder(f.kubeConfigFlags)
}
