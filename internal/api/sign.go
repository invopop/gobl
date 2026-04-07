package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/invopop/gobl/internal/cli"
)

func handleSign(w http.ResponseWriter, r *http.Request) {
	req := new(signRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		writeError(w, &cli.Error{Code: http.StatusBadRequest, Message: "invalid JSON: " + err.Error()})
		return
	}
	if len(req.Data) == 0 {
		writeError(w, &cli.Error{Code: http.StatusBadRequest, Message: "no payload"})
		return
	}

	opts := &cli.SignOptions{
		ParseOptions: &cli.ParseOptions{
			DocType: req.DocType,
			Input:   bytes.NewReader(req.Data),
			Envelop: req.Envelop,
		},
		PrivateKey: req.PrivateKey,
	}
	if len(req.Template) != 0 {
		opts.Template = bytes.NewReader(req.Template)
	}

	result, err := cli.Sign(r.Context(), opts)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, result)
}
