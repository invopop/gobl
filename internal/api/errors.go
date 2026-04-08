package api

import (
	"encoding/json"
	"net/http"

	"github.com/invopop/gobl/internal/cli"
)

// writeError writes a JSON error response. If the error is a *cli.Error,
// its Code is used as the HTTP status; otherwise 500 is used.
func writeError(w http.ResponseWriter, err error) {
	cliErr, ok := err.(*cli.Error)
	if !ok {
		cliErr = &cli.Error{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(cliErr.Code)
	_ = json.NewEncoder(w).Encode(cliErr) //nolint:errcheck
}

// writeJSON writes a JSON response with 200 OK status.
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v) //nolint:errcheck
}

// writeRawJSON writes pre-encoded JSON bytes as a response.
func writeRawJSON(w http.ResponseWriter, d []byte) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(d) //nolint:errcheck
}
