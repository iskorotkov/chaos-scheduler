package targets

import (
	"context"
	"errors"
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/logger"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"strings"
	"time"
)

var (
	ConfigError    = errors.New("couldn't read config")
	ClientsetError = errors.New("couldn't create clientset from config")
	FetchError     = errors.New("couldn't fetch info from Kubernetes")
)

type Target struct {
	Pod         string
	Deployment  string
	Containers  []string
	Labels      map[string]string
	Annotations map[string]string
}

func (t Target) MainContainer() string {
	return t.Containers[0]
}

func (t Target) Selector() string {
	return fmt.Sprintf("app=%s", t.Labels["app"])
}

type Seeker struct {
	namespace string
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
			Pod:         pod.Name,
			Deployment:  deployment,
			Containers:  containers,
			Labels:      pod.Labels,
			Annotations: pod.Annotations,
		}
		res = append(res, p)
	}

	return res, nil
}

func NewSeeker(namespace string, isKubernetes bool) (Seeker, error) {
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
		return Seeker{}, ConfigError
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error(err)
		return Seeker{}, ClientsetError
	}

	return Seeker{
		namespace: namespace,
		clientset: clientset,
	}, nil
}
