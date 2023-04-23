package k8s

import (
	"os"

	"github.com/google/wire"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var ProviderSet = wire.NewSet(
	NewClientset,
)

// NewClientset returns a new kubernetes clientset.
func NewClientset() (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
