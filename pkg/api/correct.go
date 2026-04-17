package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/internal/ops"
)

func handleCorrect(w http.ResponseWriter, r *http.Request) {
	req := new(correctRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		WriteError(w, gobl.ErrInput.WithCause(fmt.Errorf("invalid JSON: %w", err)))
		return
	}
	if len(req.Data) == 0 {
		WriteError(w, gobl.ErrInput.WithReason("no payload"))
		return
	}

	opts := &ops.CorrectOptions{
		ParseOptions: &ops.ParseOptions{
			Input: bytes.NewReader(req.Data),
		},
		OptionsSchema: req.Schema,
		Data:          req.Options,
	}

	result, err := ops.Correct(r.Context(), opts)
	if err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, result)
}
