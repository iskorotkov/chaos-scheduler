package main

import (
	"fmt"
	"github.com/iskorotkov/chaos-scheduler/internal/config"
	"github.com/iskorotkov/chaos-scheduler/internal/pages"
	"github.com/iskorotkov/chaos-scheduler/pkg/server"
	"log"
	"net/http"
)

func main() {
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./web/js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))

	cfg := config.ParseConfigFromEnv()

	http.HandleFunc("/", pages.Homepage)
	http.HandleFunc("/scenarios", server.WithConfig(pages.Scenarios, cfg))

	fmt.Println("Open http://127.0.0.1:8811 to work with scheduler")
	fmt.Println("Open http://127.0.0.1:2746 to see workflow progress in Argo UI")

	log.Fatal(http.ListenAndServe(":8811", nil))
}
