package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/internal/ops"
)

func handleVerify(w http.ResponseWriter, r *http.Request) {
	req := new(verifyRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		WriteError(w, gobl.ErrInput.WithCause(fmt.Errorf("invalid JSON: %w", err)))
		return
	}
	if len(req.Data) == 0 {
		WriteError(w, gobl.ErrInput.WithReason("no payload"))
		return
	}

	if err := ops.Verify(r.Context(), bytes.NewReader(req.Data), req.PublicKey); err != nil {
		WriteError(w, err)
		return
	}
	WriteJSON(w, ops.VerifyResponse{OK: true})
}
