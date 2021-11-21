// Package k8s allows to fetch a list of targets from Kubernetes.
package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/iskorotkov/chaos-scheduler/pkg/workflows/targets"
	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
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

// List returns list of targets in a specified namespace.
func (k finder) List(namespace, label string) ([]targets.Target, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	pods, err := k.clientset.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		k.logger.Error(err)
		return nil, targets.ErrFetch
	}

	res := make([]targets.Target, 0)
	for _, pod := range pods.Items {
		containers := make([]string, 0)
		for _, container := range pod.Spec.Containers {
			containers = append(containers, container.Name)
		}

		appLabel, ok := pod.Labels[label]
		if !ok {
			k.logger.Error("pods doesn't have a required label", "pod", pod)
			return nil, targets.ErrInvalidTarget
		}

		p := targets.Target{
			Pod:           pod.Name,
			Node:          pod.Spec.NodeName,
			MainContainer: containers[0],
			Containers:    containers,
			AppLabel:      fmt.Sprintf("%s=%s", label, appLabel),
			AppLabelValue: appLabel,
			Labels:        pod.Labels,
			Annotations:   pod.Annotations,
		}
		res = append(res, p)
	}

	return res, nil
}

// newClient returns configured Kubernetes clientset.
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
		return nil, targets.ErrClient
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error(err.Error())
		return nil, targets.ErrClient
	}

	return clientset, nil
}
