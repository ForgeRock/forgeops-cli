package k8s

import (
	"errors"
	"os"
	"path/filepath"

	// Load Auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/ForgeRock/forgeops-cli/internal/printer"
)

// GetClient setup kubernetes client
func GetClient(path string) (*kubernetes.Clientset, error) {
	kubeconfig, present := os.LookupEnv("KUBECONFIG")
	if !present && path == "" {
		home := ""
		if home = homedir.HomeDir(); home == "" {
			printer.Errorln("couldn't find home dir")
			return &kubernetes.Clientset{}, errors.New("home couldn't establish kubeconfig in home or environment")
		}
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	if path != "" {
		kubeconfig = path
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		printer.Errorln("couldn't load conf")
		return &kubernetes.Clientset{}, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		printer.Errorln("couldn't load conf")
		return &kubernetes.Clientset{}, err
	}
	return clientset, nil
}
