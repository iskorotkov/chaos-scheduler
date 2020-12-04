package main

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/targets"
)

func main() {
	observer, _ := targets.NewSeeker("chaos-app", false)
	pods, _ := observer.Targets()
	for i, pod := range pods {
		fmt.Printf("%d: %#v\n", i, pod)
	}
}
