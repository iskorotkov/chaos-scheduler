package main

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/argo"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	host         = os.Getenv("EXECUTOR_HOST")
	port         = os.Getenv("EXECUTOR_PORT")
	failuresPath = os.Getenv("FAILURES_PATH")
	templatePath = os.Getenv("TEMPLATE_PATH")
	executor     argo.Executor
)

func main() {
	p, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		log.Fatalf("couldn't parse argo port")
	}
	executor = argo.NewExecutor(host, int(p))

	if failuresPath == "" {
		log.Fatalf("path to failures isn't set")
	}

	if templatePath == "" {
		log.Fatalf("path to scenario template isn't set")
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/home", homePage)
	http.HandleFunc("/test", test)
	log.Fatal(http.ListenAndServe(":8811", nil))
}

func homePage(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		returnHTMLPage(rw, "templates/html/home.html", nil)
	} else {
		methodNotAvailable(rw, r)
	}
}

func test(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			badRequest(rw, fmt.Sprintf("couldn't parse form data: %v\n", err))
			return
		}

		stages, err := strconv.ParseInt(r.FormValue("stages"), 10, 32)
		if err != nil {
			badRequest(rw, fmt.Sprintf("couldn't parse number of stages: %v\n", err))
			return
		}

		scenario, err := argo.NewScenario(argo.ScenarioConfig{Path: failuresPath, Stages: int(stages)})
		if err != nil {
			internalError(rw, fmt.Sprintf("couldn't create test scenario: %v", err))
		}

		output, err := argo.Format(argo.FormatConfig{TemplatePath: templatePath, Scenario: scenario})
		if err != nil {
			internalError(rw, fmt.Sprintf("couldn't convert scenario to given format: %v\n", err))
		}

		returnHTMLPage(rw, "templates/html/test.html", output)
	} else {
		methodNotAvailable(rw, r)
	}
}

func returnHTMLPage(rw http.ResponseWriter, path string, data interface{}) {
	t, err := template.ParseFiles(path)
	if err != nil {
		log.Println(fmt.Println(err))
		return
	}

	err = t.Execute(rw, data)
	if err != nil {
		log.Println(err)
		return
	}
}

func methodNotAvailable(rw http.ResponseWriter, r *http.Request) {
	err := fmt.Sprintf("method %v is not supported", r.Method)
	log.Println(err)
	http.Error(rw, err, http.StatusBadRequest)
}

func badRequest(rw http.ResponseWriter, msg string) {
	log.Println(msg)
	http.Error(rw, msg, http.StatusBadRequest)
}

func internalError(rw http.ResponseWriter, msg string) {
	log.Println(msg)
	http.Error(rw, msg, http.StatusInternalServerError)
}
