// Package main sets up the routing and starts the HTTP server for the web application.
package main

import (
	"net/http"

	// Importing the handlers package that contains the application's logic for handling requests.
	"github.com/RAprogramm/wsGolangChat/internal/handlers"

	// Importing the pat package for HTTP request routing.
	"github.com/bmizerany/pat"
)

// routes sets up the HTTP routes for the application using the pat router.
func routes() http.Handler {
	// Creating a new router instance.
	mux := pat.New()

	// Associating the root URL path "/" with the Home handler from the handlers package.
	mux.Get("/", http.HandlerFunc(handlers.Home))

	// Associating the "/ws" URL path with the WsEndpoint handler from the handlers package to handle websocket connections.
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndpoint))

	// Setting up a file server to serve static files from the "./static/" directory.
	fileServer := http.FileServer(http.Dir("./static/"))

	// Stripping the "/static" prefix before the file server handles the request so that
	// the file paths are relative to the "./static/" directory.
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	// Returning the configured router.
	return mux
}
