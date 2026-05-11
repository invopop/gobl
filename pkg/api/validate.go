package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/internal/ops"
)

func handleValidate(w http.ResponseWriter, r *http.Request) {
	req := new(validateRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		WriteError(w, gobl.ErrInput.WithCause(fmt.Errorf("invalid JSON: %w", err)))
		return
	}
	if len(req.Data) == 0 {
		WriteError(w, gobl.ErrInput.WithReason("no payload"))
		return
	}

	if err := ops.Validate(r.Context(), bytes.NewReader(req.Data)); err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, ops.ValidateResponse{OK: true})
}
