package api

import (
	"encoding/json"
	"net/http"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
)

// httpStatusForKey maps domain error keys to HTTP status codes.
func httpStatusForKey(key cbc.Key) int {
	switch key {
	case gobl.ErrInput.Key():
		return http.StatusBadRequest // 400
	case gobl.ErrNotFound.Key():
		return http.StatusNotFound // 404
	case gobl.ErrInternal.Key():
		return http.StatusInternalServerError // 500
	default:
		return http.StatusUnprocessableEntity // 422 for all domain errors
	}
}

// writeError writes a JSON error response. The HTTP status is derived
// from the error's Key using httpStatusForKey.
func writeError(w http.ResponseWriter, err error) {
	ge, ok := err.(*gobl.Error)
	if !ok {
		ge = gobl.ErrInternal.WithCause(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusForKey(ge.Key()))
	_ = json.NewEncoder(w).Encode(ge) //nolint:errcheck
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
