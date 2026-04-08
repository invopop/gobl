package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/internal/cli"
)

func handleCorrect(w http.ResponseWriter, r *http.Request) {
	req := new(correctRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeError(w, gobl.ErrInput.WithCause(fmt.Errorf("invalid JSON: %w", err)))
		return
	}
	if len(req.Data) == 0 {
		writeError(w, gobl.ErrInput.WithReason("no payload"))
		return
	}

	opts := &cli.CorrectOptions{
		ParseOptions: &cli.ParseOptions{
			Input: bytes.NewReader(req.Data),
		},
		OptionsSchema: req.Schema,
		Data:          req.Options,
	}

	result, err := cli.Correct(r.Context(), opts)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, result)
}
