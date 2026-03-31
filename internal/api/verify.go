package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/invopop/gobl/internal/cli"
)

func handleVerify(w http.ResponseWriter, r *http.Request) {
	req := new(verifyRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeError(w, &cli.Error{Code: http.StatusBadRequest, Message: "invalid JSON: " + err.Error()})
		return
	}
	if len(req.Data) == 0 {
		writeError(w, &cli.Error{Code: http.StatusBadRequest, Message: "no payload"})
		return
	}

	if err := cli.Verify(r.Context(), bytes.NewReader(req.Data), req.PublicKey); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cli.VerifyResponse{OK: true})
}
