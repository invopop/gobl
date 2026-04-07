package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/invopop/gobl/internal/cli"
)

func handleValidate(w http.ResponseWriter, r *http.Request) {
	req := new(validateRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeError(w, &cli.Error{Code: http.StatusBadRequest, Message: "invalid JSON: " + err.Error()})
		return
	}
	if len(req.Data) == 0 {
		writeError(w, &cli.Error{Code: http.StatusBadRequest, Message: "no payload"})
		return
	}

	if err := cli.Validate(r.Context(), bytes.NewReader(req.Data)); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, cli.ValidateResponse{OK: true})
}
