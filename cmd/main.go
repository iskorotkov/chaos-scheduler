package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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
		returnHTMLPage(rw, "templates/html/test.html", nil)
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
