package api

import (
	"net/http"
	"path"
	"sort"
	"strings"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/data"
	"github.com/invopop/gobl/schema"
)

func handleSchemaList(w http.ResponseWriter, _ *http.Request) {
	list := schema.List()
	items := make([]string, len(list))
	for i, v := range list {
		items[i] = v.String()
	}
	sort.Strings(items)
	WriteJSON(w, map[string]any{"schemas": items})
}

func handleSchema(w http.ResponseWriter, r *http.Request) {
	p := r.PathValue("path")
	if p == "" {
		WriteError(w, gobl.ErrInput.WithReason("missing schema path"))
		return
	}
	if !strings.HasSuffix(p, ".json") {
		p = p + ".json"
	}
	p = path.Join("schemas", p)

	if _, ok := r.URL.Query()["bundle"]; ok {
		d, err := schema.BundleSchema(p)
		if err != nil {
			WriteError(w, gobl.ErrNotFound.WithReason("schema not found"))
			return
		}
		WriteRawJSON(w, d)
		return
	}

	d, err := data.Content.ReadFile(p)
	if err != nil {
		WriteError(w, gobl.ErrNotFound.WithReason("schema not found"))
		return
	}
	WriteRawJSON(w, d)
}
