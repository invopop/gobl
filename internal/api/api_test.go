package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/internal/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const prefix = "/v0"

var testInvoice = json.RawMessage(`{
	"$schema": "https://gobl.org/draft-0/bill/invoice",
	"currency": "EUR",
	"issue_date": "2024-01-01",
	"supplier": {
		"tax_id": { "country": "ES", "code": "B85905495" },
		"name": "Test Supplier"
	},
	"customer": {
		"tax_id": { "country": "ES", "code": "B85905495" },
		"name": "Test Customer"
	},
	"lines": [
		{
			"quantity": "1",
			"item": {
				"name": "Test Item",
				"price": "100.00"
			},
			"taxes": [
				{ "cat": "VAT", "rate": "standard" }
			]
		}
	]
}`)

func postJSON(t *testing.T, url string, v any) *http.Response {
	t.Helper()
	body, err := json.Marshal(v)
	require.NoError(t, err)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	return resp
}

func TestVersionEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + prefix + "/")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	assert.Equal(t, "Welcome", body["gobl"])
	assert.NotEmpty(t, body["version"])
}

func TestBuildEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp := postJSON(t, srv.URL+prefix+"/build", map[string]any{"data": testInvoice})
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.NotNil(t, result)
}

func TestBuildEndpointNoPayload(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp := postJSON(t, srv.URL+prefix+"/build", map[string]any{})
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestValidateEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	// First build to get a valid document
	buildResp := postJSON(t, srv.URL+prefix+"/build", map[string]any{"data": testInvoice})
	require.Equal(t, http.StatusOK, buildResp.StatusCode)
	builtData, _ := io.ReadAll(buildResp.Body)
	buildResp.Body.Close() //nolint:errcheck

	// Now validate it
	resp := postJSON(t, srv.URL+prefix+"/validate", map[string]any{"data": json.RawMessage(builtData)})
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.Equal(t, true, result["ok"])
}

func TestSchemaListEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + prefix + "/schemas")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	schemas, ok := result["schemas"].([]any)
	assert.True(t, ok)
	assert.Greater(t, len(schemas), 0)
}

func TestSchemaEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + prefix + "/schemas/bill/invoice")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}

func TestSchemaEndpointNotFound(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + prefix + "/schemas/nonexistent/thing")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestSchemaBundleQueryParam(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + prefix + "/schemas/bill/invoice?bundle")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

	var result map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	// Should still have $id and $ref for the root type
	assert.Equal(t, "https://gobl.org/draft-0/bill/invoice", result["$id"])
	assert.Contains(t, result["$ref"], "#/$defs/")

	// $defs should contain the root type plus all dependencies
	defs, ok := result["$defs"].(map[string]any)
	require.True(t, ok)
	assert.Greater(t, len(defs), 10)

	// All $ref values should be internal (#/$defs/...), not external URLs.
	raw, _ := json.Marshal(result)
	assert.NotContains(t, string(raw), `"$ref":"https://gobl.org/`)
	assert.NotContains(t, string(raw), `"$ref": "https://gobl.org/`)
}

func TestRegimeEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + prefix + "/regimes/ES")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}

func TestRegimeEndpointNotFound(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + prefix + "/regimes/XX")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestServerTimingHeader(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + prefix + "/")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	timing := resp.Header.Get("Server-Timing")
	assert.Contains(t, timing, "total;dur=")
}

func TestVersionHeader(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + prefix + "/")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, string(gobl.VERSION), resp.Header.Get("GOBL-Version"))
}

func TestETagHeader(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + prefix + "/schemas")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	etag := resp.Header.Get("ETag")
	assert.Equal(t, `"`+string(gobl.VERSION)+`"`, etag)
}

func TestETagNotModified(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	etag := `"` + string(gobl.VERSION) + `"`
	req, _ := http.NewRequest(http.MethodGet, srv.URL+prefix+"/schemas", nil)
	req.Header.Set("If-None-Match", etag)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusNotModified, resp.StatusCode)
}

func TestCORSHeaders(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + prefix + "/")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
}

func TestCORSPreflight(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodOptions, srv.URL+prefix+"/build", nil)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
}
