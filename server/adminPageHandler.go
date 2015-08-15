package server

import (
	"net/http"
	"strings"
	// "strconv"
	//
	// "github.com/dragonrider23/inca/devices"
)

func adminPageHandler(w http.ResponseWriter, r *http.Request) {
	urlPieces := strings.Split(r.URL.Path, "/")[2:]

	if urlPieces[0] == "partial" {
		adminRenderPartial(w, r, urlPieces[1:])
	} else {
		renderTemplate(w, "adminpage", nil)
	}
}

func adminRenderPartial(w http.ResponseWriter, r *http.Request, p []string) {
	if len(p) == 0 {
		httpErrorPage(w, "Failed to load admin partial, incorrect url", false)
	}

	switch p[0] {
	case "cp":
		renderTemplate(w, "admin-cppartial", nil)
	case "dt":
		renderTemplate(w, "admin-dtpartial", nil)
	default:
		w.Write([]byte(""))
	}
}
