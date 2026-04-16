package api

import (
	"net/http"
	"path"
	"strings"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/data"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

type addonSummary struct {
	Key         string      `json:"key"`
	Name        i18n.String `json:"name"`
	Description i18n.String `json:"description,omitempty"`
	Requires    []cbc.Key   `json:"requires,omitempty"`
}

func handleAddonList(w http.ResponseWriter, _ *http.Request) {
	defs := tax.AllAddonDefs()
	items := make([]addonSummary, len(defs))
	for i, a := range defs {
		items[i] = addonSummary{
			Key:         string(a.Key),
			Name:        a.Name,
			Description: a.Description,
			Requires:    a.Requires,
		}
	}
	writeJSON(w, map[string]any{"addons": items})
}

func handleAddon(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		writeError(w, gobl.ErrInput.WithReason("missing addon key"))
		return
	}

	if !strings.HasSuffix(key, ".json") {
		key = key + ".json"
	}
	p := path.Join("addons", key)

	d, err := data.Content.ReadFile(p)
	if err != nil {
		writeError(w, gobl.ErrNotFound.WithReason("addon not found"))
		return
	}
	writeRawJSON(w, d)
}
