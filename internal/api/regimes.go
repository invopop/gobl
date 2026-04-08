package api

import (
	"net/http"
	"path"
	"strings"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/data"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

type regimeSummary struct {
	Country     string      `json:"country"`
	Name        i18n.String `json:"name"`
	Description i18n.String `json:"description,omitempty"`
	Currency    string      `json:"currency"`
}

func handleRegimeList(w http.ResponseWriter, _ *http.Request) {
	defs := tax.AllRegimeDefs()
	items := make([]regimeSummary, len(defs))
	for i, r := range defs {
		items[i] = regimeSummary{
			Country:     string(r.Country),
			Name:        r.Name,
			Description: r.Description,
			Currency:    string(r.Currency),
		}
	}
	writeJSON(w, map[string]any{"regimes": items})
}

func handleRegime(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		writeError(w, gobl.ErrInput.WithReason("missing regime code"))
		return
	}

	p := path.Join("regimes", strings.ToLower(code)+".json")
	d, err := data.Content.ReadFile(p)
	if err != nil {
		writeError(w, gobl.ErrNotFound.WithReason("regime not found"))
		return
	}
	writeRawJSON(w, d)
}
