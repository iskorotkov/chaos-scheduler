package targets

import (
	"context"
	"errors"
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/k8s"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

var (
	ClientsetError = errors.New("couldn't create clientset")
	FetchError     = errors.New("couldn't fetch info from Kubernetes")
)

type Target struct {
	Pod           string            `json:"pod"`
	Deployment    string            `json:"deployment"`
	Node          string            `json:"node"`
	MainContainer string            `json:"mainContainer"`
	Containers    []string          `json:"containers"`
	AppLabel      string            `json:"appLabel"`
	Labels        map[string]string `json:"labels"`
	Annotations   map[string]string `json:"annotations"`
}

func (t Target) Generate(r *rand.Rand, _ int) reflect.Value {
	randomStr := func(prefix string) string {
		return fmt.Sprintf("%s-%d", prefix, r.Int())
	}

	containers := make([]string, 0)
	for i := 0; i < 1+r.Intn(10); i++ {
		containers = append(containers, randomStr("container"))
	}

	return reflect.ValueOf(Target{
		Pod:           randomStr("pod"),
		Deployment:    randomStr("deploy"),
		Node:          randomStr("node"),
		MainContainer: containers[r.Intn(len(containers))],
		Containers:    containers,
		AppLabel:      randomStr("label"),
		Labels: map[string]string{
			randomStr("label1"): randomStr("value"),
			randomStr("label2"): randomStr("value"),
		},
		Annotations: map[string]string{
			randomStr("annotation1"): randomStr("value"),
			randomStr("annotation2"): randomStr("value"),
		},
	})
}

func List(namespace string, label string, logger *zap.SugaredLogger) ([]Target, error) {
	clientset, err := k8s.NewClient(logger.Named("k8s"))
	if err != nil {
		logger.Error(err)
		return nil, ClientsetError
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	pods, err := clientset.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{})
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
