// Package main is entry application point
package main

import (
	"log"
	"net/http"

	"github.com/RAprogramm/wsGolangChat/internal/handlers"
)

func main() {
	mux := routes()

	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()

	log.Println("Starting web server on port 8080")

	_ = http.ListenAndServe(":8080", mux)
}
