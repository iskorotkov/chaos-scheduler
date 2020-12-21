package targets

import (
	"context"
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/k8s"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"strings"
	"time"
)

type KubernetesSeeker struct {
	namespace string
	appLabel  string
	clientset *kubernetes.Clientset
	logger    *zap.SugaredLogger
}

func (s KubernetesSeeker) Targets() ([]Target, error) {
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
			MainContainer: containers[0],
			Containers:    containers,
			AppLabel:      fmt.Sprintf("%s=%s", s.appLabel, pod.Labels[s.appLabel]),
			Labels:        pod.Labels,
			Annotations:   pod.Annotations,
		}
		res = append(res, p)
	}

	return res, nil
}

func NewSeeker(namespace string, appLabel string, logger *zap.SugaredLogger) (KubernetesSeeker, error) {
	isKubernetes := os.Getenv("KUBERNETES_SERVICE_HOST") != ""
	clientset, err := k8s.NewClient(isKubernetes, logger.Named("k8s"))
	if err != nil {
		logger.Error(err)
		return KubernetesSeeker{}, ClientsetError
	}

	return KubernetesSeeker{
		namespace: namespace,
		appLabel:  appLabel,
		clientset: clientset,
		logger:    logger,
	}, nil
}
