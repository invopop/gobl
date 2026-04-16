package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/internal/ops"
)

func handleBuild(w http.ResponseWriter, r *http.Request) {
	req := new(buildRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		WriteError(w, gobl.ErrInput.WithCause(fmt.Errorf("invalid JSON: %w", err)))
		return
	}
	if len(req.Data) == 0 {
		WriteError(w, gobl.ErrInput.WithReason("no payload"))
		return
	}

	opts := &ops.BuildOptions{
		ParseOptions: &ops.ParseOptions{
			Input:   bytes.NewReader(req.Data),
			DocType: req.DocType,
			Envelop: req.Envelop,
		},
	}
	if len(req.Template) != 0 {
		opts.Template = bytes.NewReader(req.Template)
	}

	result, err := ops.Build(r.Context(), opts)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, result)
}
