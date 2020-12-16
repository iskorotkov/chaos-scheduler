package targets

import (
	"context"
	"errors"
	"github.com/iskorotkov/chaos-scheduler/pkg/k8s"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
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
	logger    *zap.SugaredLogger
}

func (s Seeker) Targets() ([]Target, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	pods, err := s.clientset.CoreV1().Pods(s.namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		s.logger.Error(err)
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
			Node:          pod.Spec.NodeName,
			Deployment:    deployment,
			Containers:    containers,
			Labels:        pod.Labels,
			SelectorLabel: s.appLabel,
			Annotations:   pod.Annotations,
		}
		res = append(res, p)
	}

	return res, nil
}

func NewSeeker(namespace string, appLabel string, logger *zap.SugaredLogger) (Seeker, error) {
	isKubernetes := os.Getenv("KUBERNETES_SERVICE_HOST") != ""
	clientset, err := k8s.NewClient(isKubernetes, logger.Named("k8s client"))
	if err != nil {
		logger.Error(err)
		return Seeker{}, ClientsetError
	}

	return Seeker{
		namespace: namespace,
		appLabel:  appLabel,
		clientset: clientset,
		logger:    logger,
	}, nil
}
