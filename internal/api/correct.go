package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/invopop/gobl/internal/cli"
)

func handleCorrect(w http.ResponseWriter, r *http.Request) {
	req := new(correctRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeError(w, &cli.Error{Code: http.StatusBadRequest, Message: "invalid JSON: " + err.Error()})
		return
	}
	if len(req.Data) == 0 {
		writeError(w, &cli.Error{Code: http.StatusBadRequest, Message: "no payload"})
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
