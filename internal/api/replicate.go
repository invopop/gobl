package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/internal/cli"
)

func handleReplicate(w http.ResponseWriter, r *http.Request) {
	req := new(replicateRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeError(w, gobl.ErrInput.WithCause(fmt.Errorf("invalid JSON: %w", err)))
		return
	}
	if len(req.Data) == 0 {
		writeError(w, gobl.ErrInput.WithReason("no payload"))
		return
	}

	opts := &cli.ReplicateOptions{
		ParseOptions: &cli.ParseOptions{
			Input: bytes.NewReader(req.Data),
		},
	}

	result, err := cli.Replicate(r.Context(), opts)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, result)
}
