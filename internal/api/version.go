package api

import (
	"net/http"

	"github.com/invopop/gobl"
)

func handleVersion(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{
		"gobl":    "Welcome",
		"version": gobl.VERSION,
	})
}
