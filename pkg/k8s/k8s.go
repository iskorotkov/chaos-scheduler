package k8s

import (
	"context"
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type finder struct {
	clientset *kubernetes.Clientset
	logger    *zap.SugaredLogger
}

func NewFinder(logger *zap.SugaredLogger) (targets.TargetFinder, error) {
	clientset, err := newClient(logger)
	if err != nil {
		return nil, err
	}

	return &finder{
		clientset: clientset,
		logger:    logger,
	}, nil
}

func (k finder) List(namespace string, label string) ([]targets.Target, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	pods, err := k.clientset.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		k.logger.Error(err)
		return nil, targets.FetchError
	}

	res := make([]targets.Target, 0)
	for _, pod := range pods.Items {
		parts := strings.Split(pod.Name, "-")
		deployment := strings.Join(parts[0:len(parts)-2], "-")

		containers := make([]string, 0)
		for _, container := range pod.Spec.Containers {
			containers = append(containers, container.Name)
		}

		p := targets.Target{
			Pod:           pod.Name,
			Node:          pod.Spec.NodeName,
			Deployment:    deployment,
			MainContainer: containers[0],
			Containers:    containers,
			AppLabel:      fmt.Sprintf("%s=%s", label, pod.Labels[label]),
			Labels:        pod.Labels,
			Annotations:   pod.Annotations,
		}
		res = append(res, p)
	}

	return res, nil
}

func newClient(logger *zap.SugaredLogger) (*kubernetes.Clientset, error) {
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
		return nil, targets.ConfigError
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error(err.Error())
		return nil, targets.ClientError
	}

	return clientset, nil
}

func Available() bool {
	var config *rest.Config
	var err error

	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		config, err = rest.InClusterConfig()
	} else {
		configFile := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", configFile)
	}

	if err != nil {
		return false
	}

	_, err = kubernetes.NewForConfig(config)
	if err != nil {
		return false
	}

	return true
}
