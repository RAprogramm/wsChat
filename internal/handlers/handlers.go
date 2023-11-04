// Package handlers provides application handlers.
package handlers

import (
	"net/http"

	"github.com/CloudyKit/jet/v6"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(), // auto refresh template when changed
)

// Home is handler for home page.
func Home(w http.ResponseWriter, r *http.Request) {

}

func renderPage(w http.ResponseWriter, r *http.Request) {

}
