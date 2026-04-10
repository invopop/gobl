// Package editor provides the GOBL web editor UI using PopUI and Templ.
package editor

import (
	"embed"
	"io/fs"
	"net/http"

	popui "github.com/invopop/popui.go"
)

//go:embed assets/*
var editorAssetsEmbed embed.FS

// editorAssets is the sub-filesystem without the "assets" prefix, used
// both for serving files and for generating versioned paths.
var editorAssets fs.FS

func init() {
	editorAssets, _ = fs.Sub(editorAssetsEmbed, "assets")
}

// AssetPath is the URL prefix for serving editor assets.
const AssetPath = "/_editor"

// Handler returns an http.HandlerFunc that renders the editor page.
func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = Page().Render(r.Context(), w)
	}
}

// RegisterAssets registers the editor and popui static asset handlers
// onto the given ServeMux.
func RegisterAssets(mux *http.ServeMux) {
	// PopUI assets at /_popui/ (embed FS contains assets/ subdirectory)
	mux.Handle(popui.AssetPath+"/", http.StripPrefix(
		popui.AssetPath, http.FileServerFS(popui.Assets),
	))
	// Editor assets at /_editor/
	mux.Handle(AssetPath+"/", http.StripPrefix(
		AssetPath, http.FileServerFS(editorAssets),
	))
}
