package main

import (
	"net/http"

	"github.com/RAprogramm/wsGolangChat/internal/handlers"
	"github.com/bmizerany/pat"
)

func routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndpoint))

	return mux
}
