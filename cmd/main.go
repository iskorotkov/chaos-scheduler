package main

import (
	"fmt"
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

	fmt.Println("Open http://127.0.0.1:8811 to work with scheduler")
	fmt.Println("Open http://127.0.0.1:2746 to see workflow progress in Argo UI")

	log.Fatal(http.ListenAndServe(":8811", nil))
}
