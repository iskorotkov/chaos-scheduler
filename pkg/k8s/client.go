package k8s

import (
	"errors"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

var (
	ConfigError    = errors.New("couldn't read config")
	ClientsetError = errors.New("couldn't create clientset from config")
)

func NewClient(logger *zap.SugaredLogger) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		config, err = rest.InClusterConfig()
	} else {
		configFile := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", configFile)
	}

	if err != nil {
		logger.Error(err.Error())
		return nil, ConfigError
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error(err.Error())
		return nil, ClientsetError
	}

	return clientset, nil
}
