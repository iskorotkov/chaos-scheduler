package targets

import (
	"context"
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/k8s"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"
)

var (
	ClientsetError = errors.New("couldn't create clientset")
	FetchError     = errors.New("couldn't fetch info from Kubernetes")
)

type Seeker struct {
	namespace string
	appLabel  string
	clientset *kubernetes.Clientset
}

func (o Seeker) Targets() ([]Target, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	pods, err := o.clientset.CoreV1().Pods(o.namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		logger.Error(err)
		return nil, FetchError
	}

	res := make([]Target, 0)
	for _, pod := range pods.Items {
		parts := strings.Split(pod.Name, "-")
		deployment := strings.Join(parts[0:len(parts)-2], "-")

		containers := make([]string, 0)
		for _, container := range pod.Spec.Containers {
			containers = append(containers, container.Name)
		}

		p := Target{
			Pod:           pod.Name,
			Deployment:    deployment,
			Containers:    containers,
			Labels:        pod.Labels,
			SelectorLabel: o.appLabel,
			Annotations:   pod.Annotations,
		}
		res = append(res, p)
	}

	return res, nil
}

func NewSeeker(namespace string, appLabel string, isKubernetes bool) (Seeker, error) {
	clientset, err := k8s.NewClient(isKubernetes)
	if err != nil {
		logger.Error(err)
		return Seeker{}, ClientsetError
	}

	return Seeker{
		namespace: namespace,
		appLabel:  appLabel,
		clientset: clientset,
	}, nil
}
