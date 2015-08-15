package server

import (
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == rootPath {
		renderTemplate(w, "index", nil)
	} else {
		http.NotFound(w, r)
	}
}
