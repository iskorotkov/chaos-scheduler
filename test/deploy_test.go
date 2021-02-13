package test

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/shell"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path"
	"strings"
	"testing"
	"text/template"
	"time"
)

const (
	serverApp = "https://raw.githubusercontent.com/iskorotkov/chaos-server/master/deploy/counter.yaml"
	clientApp = "https://raw.githubusercontent.com/iskorotkov/chaos-client/master/deploy/counter.yaml"
)

func TestDeploy(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test skipped")
	}

	timestamp := time.Now().Format("20060102-150405")
	image := fmt.Sprintf("iskorotkov/chaos-scheduler:test-%s", timestamp)
	nsPrefix := fmt.Sprintf("test-scheduler-%s", timestamp)

	defaultOptions := k8s.NewKubectlOptions("", "", "default")
	appsOptions := k8s.NewKubectlOptions("", "", fmt.Sprintf("%s-apps", nsPrefix))
	serviceOptions := k8s.NewKubectlOptions("", "", fmt.Sprintf("%s-service", nsPrefix))

	t.Logf("Building image: %s", image)
	shell.RunCommand(t, shell.Command{
		WorkingDir: "../.",
		Command:    "docker",
		Args:       []string{"build", "-f", "./build/scheduler.dockerfile", "-t", image, "."},
	})

	t.Logf("Pushing image: %s", image)
	shell.RunCommand(t, shell.Command{
		Command: "docker",
		Args:    []string{"push", image},
	})

	// Assume that Litmus is available in litmus namespace
	// Assume that Argo is available in argo namespace

	t.Logf("Preparing YAML manifest file")
	tpl, err := template.ParseFiles(path.Join("../", "./deploy/scheduler-test.yaml"))
	if err != nil {
		panic(err)
	}

	var builder strings.Builder
	if err := tpl.Execute(&builder, struct {
		Image     string
		Namespace string
		AppNS     string
		ChaosNS   string
	}{
		Image:     image,
		Namespace: nsPrefix,
		AppNS:     appsOptions.Namespace,
		ChaosNS:   serviceOptions.Namespace,
	}); err != nil {
		panic(err)
	}

	t.Logf("Create namespaces: %s and %s", serviceOptions.Namespace, appsOptions.Namespace)
	k8s.CreateNamespace(t, defaultOptions, serviceOptions.Namespace)
	defer k8s.DeleteNamespace(t, defaultOptions, serviceOptions.Namespace)
	k8s.CreateNamespace(t, defaultOptions, appsOptions.Namespace)
	defer k8s.DeleteNamespace(t, defaultOptions, appsOptions.Namespace)

	t.Log("Deploy service to Kubernetes")
	k8s.KubectlApplyFromString(t, serviceOptions, builder.String())
	defer k8s.KubectlDeleteFromString(t, serviceOptions, builder.String())

	t.Log("Check service deployment")
	k8s.WaitUntilNumPodsCreated(t, serviceOptions, v1.ListOptions{}, 1, 10, time.Second)
	k8s.WaitUntilServiceAvailable(t, serviceOptions, "scheduler", 10, time.Second)

	t.Log("Deploy test apps to Kubernetes")
	k8s.KubectlApply(t, appsOptions, serverApp)
	defer k8s.KubectlDelete(t, appsOptions, serverApp)
	k8s.KubectlApply(t, appsOptions, clientApp)
	defer k8s.KubectlDelete(t, appsOptions, clientApp)

	t.Log("Check test apps deployment")
	k8s.WaitUntilNumPodsCreated(t, appsOptions, v1.ListOptions{}, 5, 10, time.Second)
	k8s.WaitUntilServiceAvailable(t, appsOptions, "server", 10, time.Second)

	// service := k8s.GetService(t, serviceOptions, "scheduler")
	// endpoint := k8s.GetServiceEndpoint(t, serviceOptions, service, 8811)

	// var body []byte
	// headers := make(map[string]string)

	//goland:noinspection HttpUrlsUsage
	// url := fmt.Sprintf("http://%s/api/v1/workflows/", endpoint)
	// http_helper.HTTPDoWithRetry(t, "GET", url, body, headers, 200, 10, time.Second, nil)
	// http_helper.HTTPDoWithRetry(t, "POST", url, body, headers, 200, 10, time.Second, nil)
}
