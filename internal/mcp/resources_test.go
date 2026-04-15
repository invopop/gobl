package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeReadRequest(uri string) mcp.ReadResourceRequest {
	return mcp.ReadResourceRequest{
		Params: mcp.ReadResourceParams{
			URI: uri,
		},
	}
}

func TestHandleSchemaResource(t *testing.T) {
	t.Run("valid schema", func(t *testing.T) {
		result, err := handleSchemaResource(context.Background(), makeReadRequest("gobl://schemas/bill/invoice"))
		require.NoError(t, err)
		require.Len(t, result, 1)

		tc, ok := result[0].(mcp.TextResourceContents)
		require.True(t, ok)
		assert.Equal(t, "gobl://schemas/bill/invoice", tc.URI)
		assert.Equal(t, "application/json", tc.MIMEType)
		assert.Contains(t, tc.Text, "properties")
	})

	t.Run("with .json suffix", func(t *testing.T) {
		result, err := handleSchemaResource(context.Background(), makeReadRequest("gobl://schemas/bill/invoice.json"))
		require.NoError(t, err)
		require.Len(t, result, 1)

		tc, ok := result[0].(mcp.TextResourceContents)
		require.True(t, ok)
		assert.Contains(t, tc.Text, "properties")
	})

	t.Run("missing path", func(t *testing.T) {
		_, err := handleSchemaResource(context.Background(), makeReadRequest("gobl://schemas/"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing schema path")
	})

	t.Run("wrong prefix", func(t *testing.T) {
		_, err := handleSchemaResource(context.Background(), makeReadRequest("gobl://other/bill/invoice"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing schema path")
	})

	t.Run("nonexistent schema", func(t *testing.T) {
		_, err := handleSchemaResource(context.Background(), makeReadRequest("gobl://schemas/nonexistent/type"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "schema not found")
	})
}

func TestHandleRegimeResource(t *testing.T) {
	t.Run("valid regime", func(t *testing.T) {
		result, err := handleRegimeResource(context.Background(), makeReadRequest("gobl://regimes/ES"))
		require.NoError(t, err)
		require.Len(t, result, 1)

		tc, ok := result[0].(mcp.TextResourceContents)
		require.True(t, ok)
		assert.Equal(t, "gobl://regimes/ES", tc.URI)
		assert.Equal(t, "application/json", tc.MIMEType)
		assert.Contains(t, tc.Text, "Spain")
	})

	t.Run("lowercase code", func(t *testing.T) {
		result, err := handleRegimeResource(context.Background(), makeReadRequest("gobl://regimes/de"))
		require.NoError(t, err)
		require.Len(t, result, 1)

		tc, ok := result[0].(mcp.TextResourceContents)
		require.True(t, ok)
		assert.Contains(t, tc.Text, "Germany")
	})

	t.Run("missing code", func(t *testing.T) {
		_, err := handleRegimeResource(context.Background(), makeReadRequest("gobl://regimes/"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing regime code")
	})

	t.Run("wrong prefix", func(t *testing.T) {
		_, err := handleRegimeResource(context.Background(), makeReadRequest("gobl://other/ES"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing regime code")
	})

	t.Run("nonexistent regime", func(t *testing.T) {
		_, err := handleRegimeResource(context.Background(), makeReadRequest("gobl://regimes/ZZ"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "regime not found")
	})
}

func TestHandleAddonResource(t *testing.T) {
	t.Run("valid addon", func(t *testing.T) {
		result, err := handleAddonResource(context.Background(), makeReadRequest("gobl://addons/es-verifactu-v1"))
		require.NoError(t, err)
		require.Len(t, result, 1)

		tc, ok := result[0].(mcp.TextResourceContents)
		require.True(t, ok)
		assert.Equal(t, "gobl://addons/es-verifactu-v1", tc.URI)
		assert.Equal(t, "application/json", tc.MIMEType)
		assert.NotEmpty(t, tc.Text)
	})

	t.Run("with .json suffix", func(t *testing.T) {
		result, err := handleAddonResource(context.Background(), makeReadRequest("gobl://addons/es-verifactu-v1.json"))
		require.NoError(t, err)
		require.Len(t, result, 1)
	})

	t.Run("missing key", func(t *testing.T) {
		_, err := handleAddonResource(context.Background(), makeReadRequest("gobl://addons/"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing addon key")
	})

	t.Run("wrong prefix", func(t *testing.T) {
		_, err := handleAddonResource(context.Background(), makeReadRequest("gobl://other/es-verifactu-v1"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "missing addon key")
	})

	t.Run("nonexistent addon", func(t *testing.T) {
		_, err := handleAddonResource(context.Background(), makeReadRequest("gobl://addons/nonexistent-addon"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "addon not found")
	})
}

func TestHandleSchemaList(t *testing.T) {
	result, err := handleSchemaList(context.Background(), mcp.ReadResourceRequest{})
	require.NoError(t, err)
	require.Len(t, result, 1)

	tc, ok := result[0].(mcp.TextResourceContents)
	require.True(t, ok)
	assert.Equal(t, "gobl://schemas", tc.URI)
	assert.Equal(t, "application/json", tc.MIMEType)

	var data struct {
		Schemas []string `json:"schemas"`
	}
	require.NoError(t, json.Unmarshal([]byte(tc.Text), &data))
	assert.NotEmpty(t, data.Schemas)
	assert.Contains(t, data.Schemas, "https://gobl.org/draft-0/bill/invoice")
}

func TestHandleRegimeList(t *testing.T) {
	result, err := handleRegimeList(context.Background(), mcp.ReadResourceRequest{})
	require.NoError(t, err)
	require.Len(t, result, 1)

	tc, ok := result[0].(mcp.TextResourceContents)
	require.True(t, ok)
	assert.Equal(t, "gobl://regimes", tc.URI)
	assert.Equal(t, "application/json", tc.MIMEType)

	var data struct {
		Regimes []string `json:"regimes"`
	}
	require.NoError(t, json.Unmarshal([]byte(tc.Text), &data))
	assert.NotEmpty(t, data.Regimes)
	assert.Contains(t, data.Regimes, "es")
	assert.Contains(t, data.Regimes, "de")
}

func TestHandleAddonList(t *testing.T) {
	result, err := handleAddonList(context.Background(), mcp.ReadResourceRequest{})
	require.NoError(t, err)
	require.Len(t, result, 1)

	tc, ok := result[0].(mcp.TextResourceContents)
	require.True(t, ok)
	assert.Equal(t, "gobl://addons", tc.URI)
	assert.Equal(t, "application/json", tc.MIMEType)

	var data struct {
		Addons []string `json:"addons"`
	}
	require.NoError(t, json.Unmarshal([]byte(tc.Text), &data))
	assert.NotEmpty(t, data.Addons)
}

func TestListDataDir(t *testing.T) {
	t.Run("regimes", func(t *testing.T) {
		items, err := listDataDir("regimes")
		require.NoError(t, err)
		assert.NotEmpty(t, items)
		assert.Contains(t, items, "es")
	})

	t.Run("addons", func(t *testing.T) {
		items, err := listDataDir("addons")
		require.NoError(t, err)
		assert.NotEmpty(t, items)
	})

	t.Run("nonexistent dir", func(t *testing.T) {
		_, err := listDataDir("nonexistent")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "reading nonexistent directory")
	})
}

func TestExtractParam(t *testing.T) {
	tests := []struct {
		name   string
		uri    string
		prefix string
		want   string
	}{
		{"normal", "gobl://schemas/bill/invoice", "gobl://schemas/", "bill/invoice"},
		{"empty after prefix", "gobl://schemas/", "gobl://schemas/", ""},
		{"wrong prefix", "gobl://other/something", "gobl://schemas/", ""},
		{"empty uri", "", "gobl://schemas/", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractParam(tt.uri, tt.prefix)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	d, err := marshalJSON(map[string]string{"key": "value"})
	require.NoError(t, err)
	assert.Contains(t, string(d), `"key": "value"`)
	// Verify indentation
	assert.Contains(t, string(d), "\t")
}
