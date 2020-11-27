package main

import (
	"github.com/iskorotkov/chaos-scheduler/api"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", api.Homepage)

	cfg := server.ParseConfig()
	http.HandleFunc("/scenarios", server.WithConfig(api.Scenarios, cfg))

	log.Fatal(http.ListenAndServe(":8811", nil))
}
