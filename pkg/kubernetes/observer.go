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
	Name string
}

type Deployment struct {
	Name              string
	AvailableReplicas int
	DesiredReplicas   int
}

type Observer struct {
	namespace string
	clientset *kubernetes.Clientset
}

func (o Observer) Deployments() ([]Deployment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	deployments, err := o.clientset.AppsV1().Deployments(o.namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		logger.Error(err)
		return nil, FetchError
	}

	res := make([]Deployment, 0)
	for _, deploy := range deployments.Items {
		res = append(res, Deployment{
			Name:              deploy.Name,
			AvailableReplicas: int(deploy.Status.AvailableReplicas),
			DesiredReplicas:   int(*deploy.Spec.Replicas),
		})
	}

	return res, nil
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
		res = append(res, Pod{Name: pod.Name})
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
