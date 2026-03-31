package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/invopop/gobl/internal/cli"
)

func handleReplicate(w http.ResponseWriter, r *http.Request) {
	req := new(replicateRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeError(w, &cli.Error{Code: http.StatusBadRequest, Message: "invalid JSON: " + err.Error()})
		return
	}
	if len(req.Data) == 0 {
		writeError(w, &cli.Error{Code: http.StatusBadRequest, Message: "no payload"})
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
	writeJSON(w, http.StatusOK, result)
}
