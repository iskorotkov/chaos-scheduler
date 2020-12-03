package kubernetes

import (
	"context"
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"time"
)

var (
	ConfigError    = errors.New("couldn't read config")
	ClientsetError = errors.New("couldn't create clientset from config")
	FetchError     = errors.New("couldn't fetch info from Kubernetes")
)

type Pod struct {
	Name   string
	Labels map[string]string
}

type Observer struct {
	namespace string
	clientset *kubernetes.Clientset
}

func (o Observer) Pods() ([]Pod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	pods, err := o.clientset.CoreV1().Pods(o.namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		logger.Error(err)
		return nil, FetchError
	}

	res := make([]Pod, 0)
	for _, pod := range pods.Items {
		p := Pod{Name: pod.Name, Labels: pod.Labels}
		res = append(res, p)
	}

	return res, nil
}

func NewObserver(namespace string, isKubernetes bool) (Observer, error) {
	var config *rest.Config
	var err error

	if isKubernetes {
		config, err = rest.InClusterConfig()
	} else {
		configFile := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", configFile)
	}

	if err != nil {
		logger.Error(err)
		return Observer{}, ConfigError
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error(err)
		return Observer{}, ClientsetError
	}

	return Observer{
		namespace: namespace,
		clientset: clientset,
	}, nil
}
