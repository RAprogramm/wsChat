// Package main serves as the entry point of the web application.
// It sets up the HTTP server and starts listening for requests.
package main

import (
	"log" // Importing the log package to log messages for debugging and logging.
	"net/http" // Importing the net/http package to work with HTTP.

	// Importing the handlers package which contains the application's logic for handling different routes.
	"github.com/RAprogramm/wsGolangChat/internal/handlers"
)

// The main function is the entry point of the application.
func main() {
	// Calling the routes function to get the route handler.
	mux := routes()

	// Logging the start of the channel listener, which will listen for websocket connections.
	log.Println("Starting channel listener")
	// Starting the channel listener in a new go routine to run concurrently with the HTTP server.
	go handlers.ListenToWsChannel()

	// Logging that the web server is starting.
	log.Println("Starting web server on port 8080")
	// Starting the HTTP server on port 8080 and passing in the mux (route handler).
	// The "_" is used to ignore the error returned by ListenAndServe.
	// In production code, it's important to handle this error properly.
	_ = http.ListenAndServe(":8080", mux)
}
