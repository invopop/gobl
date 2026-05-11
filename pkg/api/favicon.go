package api

import (
	_ "embed"
	"net/http"
)

//go:embed favicon.svg
var faviconSVG []byte

func handleFavicon(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Header().Set("ETag", etag)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(faviconSVG) //nolint:errcheck
}
