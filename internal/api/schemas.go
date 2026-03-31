package api

import (
	"net/http"
	"path"
	"sort"
	"strings"

	"github.com/invopop/gobl/data"
	"github.com/invopop/gobl/internal/cli"
	"github.com/invopop/gobl/schema"
)

func handleSchemaList(w http.ResponseWriter, _ *http.Request) {
	list := schema.List()
	items := make([]string, len(list))
	for i, v := range list {
		items[i] = v.String()
	}
	sort.Strings(items)
	writeJSON(w, http.StatusOK, map[string]any{"schemas": items})
}

func handleSchema(w http.ResponseWriter, r *http.Request) {
	p := r.PathValue("path")
	if p == "" {
		writeError(w, &cli.Error{Code: http.StatusBadRequest, Message: "missing schema path"})
		return
	}
	if !strings.HasSuffix(p, ".json") {
		p = p + ".json"
	}
	p = path.Join("schemas", p)

	d, err := data.Content.ReadFile(p)
	if err != nil {
		writeError(w, &cli.Error{Code: http.StatusNotFound, Message: "schema not found"})
		return
	}
	writeRawJSON(w, http.StatusOK, d)
}
