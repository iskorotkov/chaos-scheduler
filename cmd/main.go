package main

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/pkg/loading"
	"github.com/iskorotkov/chaos-scheduler/pkg/scenario"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
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
		returnMethodNotAvailable(rw, r)
	}
}

func test(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		path := os.Getenv("FAILURES_PATH")
		if path == "" {
			returnBadRequest(rw, "FAILURES_PATH env var isn't set")
			return
		}

		failures, err := loading.Load(path)
		if err != nil {
			returnBadRequest(rw, fmt.Sprintf("couldn't load failures: %v\n", err))
			return
		}

		err = r.ParseForm()
		if err != nil {
			returnBadRequest(rw, fmt.Sprintf("couldn't parse form data: %v\n", err))
			return
		}

		stages, err := strconv.ParseInt(r.FormValue("stages"), 10, 32)
		if err != nil {
			returnBadRequest(rw, fmt.Sprintf("couldn't parse number of stages: %v\n", err))
			return
		}

		config := scenario.Config{Failures: failures, Stages: int(stages)}
		sc, err := scenario.NewScenario(config)
		if err != nil {
			returnBadRequest(rw, fmt.Sprintf("couldn't create test scenario: %v\n", err))
			return
		}

		returnHTMLPage(rw, "templates/html/test.html", sc)
	} else {
		returnMethodNotAvailable(rw, r)
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

func returnMethodNotAvailable(rw http.ResponseWriter, r *http.Request) {
	err := fmt.Sprintf("method %v is not supported", r.Method)
	log.Println(err)
	http.Error(rw, err, http.StatusBadRequest)
}

func returnBadRequest(rw http.ResponseWriter, msg string) {
	log.Println(msg)
	http.Error(rw, msg, http.StatusBadRequest)
}
