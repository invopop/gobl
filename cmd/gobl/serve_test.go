package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const prefix = "/v0"

func TestServeVersion(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + prefix + "/")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, string(gobl.VERSION), resp.Header.Get("GOBL-Version"))

	var body map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "Welcome", body["gobl"])
	assert.NotEmpty(t, body["version"])
}

func TestServeBuild(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	invoice := json.RawMessage(`{
		"$schema": "https://gobl.org/draft-0/bill/invoice",
		"currency": "EUR",
		"issue_date": "2024-01-01",
		"supplier": {
			"tax_id": {"country": "ES", "code": "B85905495"},
			"name": "Test Supplier"
		},
		"customer": {
			"tax_id": {"country": "ES", "code": "B85905495"},
			"name": "Test Customer"
		},
		"lines": [{
			"quantity": "1",
			"item": {"name": "Test Item", "price": "100.00"},
			"taxes": [{"cat": "VAT", "rate": "standard"}]
		}]
	}`)

	body, _ := json.Marshal(map[string]any{"data": invoice})
	resp, err := http.Post(srv.URL+prefix+"/build", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.NotNil(t, result)
}

func TestServeBuildNoPayload(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	body, _ := json.Marshal(map[string]any{})
	resp, err := http.Post(srv.URL+prefix+"/build", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestServeVerify(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	body, _ := json.Marshal(map[string]any{"data": json.RawMessage(`{}`)})
	resp, err := http.Post(srv.URL+prefix+"/verify", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
}

func TestServeKeygen(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Post(srv.URL+prefix+"/keygen", "application/json", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.NotNil(t, result["private"])
	assert.NotNil(t, result["public"])
}
