package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
)

func TestHttpStatusForKey(t *testing.T) {
	tests := []struct {
		name   string
		key    cbc.Key
		status int
	}{
		{name: "input", key: "input", status: http.StatusBadRequest},
		{name: "not-found", key: "not-found", status: http.StatusNotFound},
		{name: "internal", key: "internal", status: http.StatusInternalServerError},
		{name: "validation", key: "validation", status: http.StatusUnprocessableEntity},
		{name: "calculation", key: "calculation", status: http.StatusUnprocessableEntity},
		{name: "unknown key", key: "something-else", status: http.StatusUnprocessableEntity},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.status, httpStatusForKey(tt.key))
		})
	}
}

func TestWriteErrorWithGoblError(t *testing.T) {
	w := httptest.NewRecorder()
	WriteError(w, gobl.ErrInput.WithReason("bad input"))
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"key":"input"`)
	assert.Contains(t, w.Body.String(), `"message":"bad input"`)
}

func TestWriteErrorWithPlainError(t *testing.T) {
	w := httptest.NewRecorder()
	WriteError(w, errors.New("something broke"))
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), `"key":"internal"`)
	assert.Contains(t, w.Body.String(), `"message":"something broke"`)
}
