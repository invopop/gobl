package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const prefix = "/v0"

// testSignableInvoice includes series+code so the envelope is ready to be signed.
var testSignableInvoice = json.RawMessage(`{
	"$schema": "https://gobl.org/draft-0/bill/invoice",
	"series": "TEST",
	"code": "001",
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

func TestBuildEndpointInvalidJSON(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Post(srv.URL+prefix+"/build", "application/json", bytes.NewReader([]byte("not-json")))
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestBuildEndpointWithTemplate(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	tmpl := json.RawMessage(`{"supplier":{"name":"Override Supplier"}}`)
	resp := postJSON(t, srv.URL+prefix+"/build", map[string]any{
		"data":     testInvoice,
		"template": tmpl,
	})
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestBuildEndpointBuildError(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	// Valid JSON but invalid invoice data that will cause a build error
	resp := postJSON(t, srv.URL+prefix+"/build", map[string]any{
		"data": json.RawMessage(`{"$schema":"https://gobl.org/draft-0/bill/invoice"}`),
	})
	defer resp.Body.Close() //nolint:errcheck
	assert.NotEqual(t, http.StatusOK, resp.StatusCode)
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

func TestSignEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	t.Run("success with key", func(t *testing.T) {
		key := dsig.NewES256Key()
		keyJSON, err := json.Marshal(key)
		require.NoError(t, err)
		resp := postJSON(t, srv.URL+prefix+"/sign", map[string]any{
			"data":       testSignableInvoice,
			"privatekey": json.RawMessage(keyJSON),
		})
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.NotNil(t, result)
	})

	t.Run("with template", func(t *testing.T) {
		key := dsig.NewES256Key()
		keyJSON, err := json.Marshal(key)
		require.NoError(t, err)
		tmpl := json.RawMessage(`{"supplier":{"name":"Override Supplier"}}`)
		resp := postJSON(t, srv.URL+prefix+"/sign", map[string]any{
			"data":       testSignableInvoice,
			"privatekey": json.RawMessage(keyJSON),
			"template":   tmpl,
		})
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("no payload", func(t *testing.T) {
		resp := postJSON(t, srv.URL+prefix+"/sign", map[string]any{})
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		resp, err := http.Post(srv.URL+prefix+"/sign", "application/json", bytes.NewReader([]byte("not-json")))
		require.NoError(t, err)
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestVerifyEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	t.Run("success", func(t *testing.T) {
		// Sign a document first
		key := dsig.NewES256Key()
		keyJSON, err := json.Marshal(key)
		require.NoError(t, err)
		pubKeyJSON, err := json.Marshal(key.Public())
		require.NoError(t, err)
		signResp := postJSON(t, srv.URL+prefix+"/sign", map[string]any{
			"data":       testSignableInvoice,
			"privatekey": json.RawMessage(keyJSON),
			"envelop":    true,
		})
		require.Equal(t, http.StatusOK, signResp.StatusCode)
		signedData, _ := io.ReadAll(signResp.Body)
		signResp.Body.Close() //nolint:errcheck

		resp := postJSON(t, srv.URL+prefix+"/verify", map[string]any{
			"data":      json.RawMessage(signedData),
			"publickey": json.RawMessage(pubKeyJSON),
		})
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Equal(t, true, result["ok"])
	})

	t.Run("no payload", func(t *testing.T) {
		resp := postJSON(t, srv.URL+prefix+"/verify", map[string]any{})
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		resp, err := http.Post(srv.URL+prefix+"/verify", "application/json", bytes.NewReader([]byte("bad")))
		require.NoError(t, err)
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("unsigned document", func(t *testing.T) {
		buildResp := postJSON(t, srv.URL+prefix+"/build", map[string]any{"data": testInvoice, "envelop": true})
		require.Equal(t, http.StatusOK, buildResp.StatusCode)
		builtData, _ := io.ReadAll(buildResp.Body)
		buildResp.Body.Close() //nolint:errcheck

		resp := postJSON(t, srv.URL+prefix+"/verify", map[string]any{"data": json.RawMessage(builtData)})
		defer resp.Body.Close() //nolint:errcheck
		assert.NotEqual(t, http.StatusOK, resp.StatusCode)
	})
}

func TestCorrectEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	// Build an envelope first
	buildResp := postJSON(t, srv.URL+prefix+"/build", map[string]any{"data": testInvoice, "envelop": true})
	require.Equal(t, http.StatusOK, buildResp.StatusCode)
	builtData, _ := io.ReadAll(buildResp.Body)
	buildResp.Body.Close() //nolint:errcheck

	t.Run("correction options schema", func(t *testing.T) {
		resp := postJSON(t, srv.URL+prefix+"/correct", map[string]any{
			"data":   json.RawMessage(builtData),
			"schema": true,
		})
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("no payload", func(t *testing.T) {
		resp := postJSON(t, srv.URL+prefix+"/correct", map[string]any{})
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		resp, err := http.Post(srv.URL+prefix+"/correct", "application/json", bytes.NewReader([]byte("bad")))
		require.NoError(t, err)
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestReplicateEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	// Build a valid document first
	buildResp := postJSON(t, srv.URL+prefix+"/build", map[string]any{"data": testInvoice})
	require.Equal(t, http.StatusOK, buildResp.StatusCode)
	builtData, _ := io.ReadAll(buildResp.Body)
	buildResp.Body.Close() //nolint:errcheck

	t.Run("success", func(t *testing.T) {
		resp := postJSON(t, srv.URL+prefix+"/replicate", map[string]any{"data": json.RawMessage(builtData)})
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]any
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.NotNil(t, result)
	})

	t.Run("no payload", func(t *testing.T) {
		resp := postJSON(t, srv.URL+prefix+"/replicate", map[string]any{})
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		resp, err := http.Post(srv.URL+prefix+"/replicate", "application/json", bytes.NewReader([]byte("bad")))
		require.NoError(t, err)
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestAddonListEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + prefix + "/addons")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	addons, ok := result["addons"].([]any)
	assert.True(t, ok)
	assert.Greater(t, len(addons), 0)
}

func TestAddonEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	t.Run("success", func(t *testing.T) {
		resp, err := http.Get(srv.URL + prefix + "/addons/es-verifactu-v1")
		require.NoError(t, err)
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	})

	t.Run("with json suffix", func(t *testing.T) {
		resp, err := http.Get(srv.URL + prefix + "/addons/es-verifactu-v1.json")
		require.NoError(t, err)
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("not found", func(t *testing.T) {
		resp, err := http.Get(srv.URL + prefix + "/addons/nonexistent-addon")
		require.NoError(t, err)
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestRegimeListEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Get(srv.URL + prefix + "/regimes")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	regimes, ok := result["regimes"].([]any)
	assert.True(t, ok)
	assert.Greater(t, len(regimes), 0)
}

func TestKeygenEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Post(srv.URL+prefix+"/keygen", "application/json", nil)
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
	assert.NotEmpty(t, result["private"])
	assert.NotEmpty(t, result["public"])
}

func TestFaviconEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler(api.WithFavicon()))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/favicon.svg")
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "image/svg+xml", resp.Header.Get("Content-Type"))
}

func TestOpenAPIEndpoint(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	t.Run("success", func(t *testing.T) {
		resp, err := http.Get(srv.URL + prefix + "/openapi.json")
		require.NoError(t, err)
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	})

	t.Run("not modified", func(t *testing.T) {
		etag := `"` + string(gobl.VERSION) + `"`
		req, _ := http.NewRequest(http.MethodGet, srv.URL+prefix+"/openapi.json", nil)
		req.Header.Set("If-None-Match", etag)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close() //nolint:errcheck
		assert.Equal(t, http.StatusNotModified, resp.StatusCode)
	})
}

func TestValidateEndpointNoPayload(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp := postJSON(t, srv.URL+prefix+"/validate", map[string]any{})
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestValidateEndpointInvalidJSON(t *testing.T) {
	srv := httptest.NewServer(api.NewHandler())
	defer srv.Close()

	resp, err := http.Post(srv.URL+prefix+"/validate", "application/json", bytes.NewReader([]byte("bad")))
	require.NoError(t, err)
	defer resp.Body.Close() //nolint:errcheck
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
