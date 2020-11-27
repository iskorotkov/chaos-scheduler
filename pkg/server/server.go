package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func WithConfig(f func(http.ResponseWriter, *http.Request, Config), cfg Config) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		f(rw, r, cfg)
	}
}

func ReturnHTMLPage(rw http.ResponseWriter, path string, data interface{}) {
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

func MethodNotAvailable(rw http.ResponseWriter, r *http.Request) {
	err := fmt.Sprintf("method %v is not supported", r.Method)
	log.Println(err)
	http.Error(rw, err, http.StatusBadRequest)
}

func BadRequest(rw http.ResponseWriter, msg string) {
	log.Println(msg)
	http.Error(rw, msg, http.StatusBadRequest)
}

func InternalError(rw http.ResponseWriter, msg string) {
	log.Println(msg)
	http.Error(rw, msg, http.StatusInternalServerError)
}
