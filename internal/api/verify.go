package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/internal/cli"
)

func handleVerify(w http.ResponseWriter, r *http.Request) {
	req := new(verifyRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeError(w, gobl.ErrInput.WithCause(fmt.Errorf("invalid JSON: %w", err)))
		return
	}
	if len(req.Data) == 0 {
		writeError(w, gobl.ErrInput.WithReason("no payload"))
		return
	}

	if err := cli.Verify(r.Context(), bytes.NewReader(req.Data), req.PublicKey); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, cli.VerifyResponse{OK: true})
}
