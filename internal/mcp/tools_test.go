package mcp_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/invopop/gobl"
	goblmcp "github.com/invopop/gobl/internal/mcp"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func callTool(t *testing.T, srv interface{ HandleMessage(context.Context, json.RawMessage) mcp.JSONRPCMessage }, name string, args map[string]any) *mcp.CallToolResult {
	t.Helper()

	params := map[string]any{
		"name":      name,
		"arguments": args,
	}
	req := map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params":  params,
	}
	raw, err := json.Marshal(req)
	require.NoError(t, err)

	resp := srv.HandleMessage(context.Background(), raw)
	require.NotNil(t, resp)

	// Marshal and re-parse to extract result
	respJSON, err := json.Marshal(resp)
	require.NoError(t, err)

	var rpcResp struct {
		Result mcp.CallToolResult `json:"result"`
		Error  *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal(respJSON, &rpcResp))
	if rpcResp.Error != nil {
		t.Fatalf("RPC error: code=%d message=%s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	return &rpcResp.Result
}

// initServer creates a server and initializes it (required before tool calls).
func initServer(t *testing.T) *initedServer {
	t.Helper()
	srv := goblmcp.NewServer()

	// Send initialize request
	initReq, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      0,
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": "2024-11-05",
			"clientInfo": map[string]any{
				"name":    "test",
				"version": "1.0",
			},
		},
	})
	resp := srv.HandleMessage(context.Background(), initReq)
	require.NotNil(t, resp)

	// Send initialized notification
	initedReq, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"method":  "notifications/initialized",
	})
	srv.HandleMessage(context.Background(), initedReq)

	return &initedServer{srv: srv, t: t}
}

type initedServer struct {
	srv *server.MCPServer
	t   *testing.T
}

// Satisfy the HandleMessage interface used by callTool
func (s *initedServer) HandleMessage(ctx context.Context, msg json.RawMessage) mcp.JSONRPCMessage {
	return s.srv.HandleMessage(ctx, msg)
}

func TestBuildTool(t *testing.T) {
	s := initServer(t)

	invoiceData := `{
		"$schema": "https://gobl.org/draft-0/bill/invoice",
		"currency": "EUR",
		"issue_date": "2024-01-15",
		"type": "standard",
		"supplier": {
			"tax_id": {"country": "ES", "code": "B91983379"},
			"name": "Test Company S.L."
		},
		"customer": {
			"tax_id": {"country": "ES", "code": "B85905495"},
			"name": "Customer S.L."
		},
		"lines": [{
			"quantity": "10",
			"item": {"name": "Test Item", "price": "100.00"},
			"taxes": [{"cat": "VAT", "rate": "standard"}]
		}]
	}`

	result := callTool(t, s, "build", map[string]any{
		"data": invoiceData,
	})

	assert.False(t, result.IsError, "build should succeed")
	require.NotEmpty(t, result.Content)

	text := extractText(t, result)
	assert.Contains(t, text, "bill/invoice")
}

func TestBuildToolWithEnvelop(t *testing.T) {
	s := initServer(t)

	invoiceData := `{
		"$schema": "https://gobl.org/draft-0/bill/invoice",
		"currency": "EUR",
		"issue_date": "2024-01-15",
		"type": "standard",
		"supplier": {
			"tax_id": {"country": "ES", "code": "B91983379"},
			"name": "Test Company S.L."
		},
		"customer": {
			"tax_id": {"country": "ES", "code": "B85905495"},
			"name": "Customer S.L."
		},
		"lines": [{
			"quantity": "10",
			"item": {"name": "Test Item", "price": "100.00"},
			"taxes": [{"cat": "VAT", "rate": "standard"}]
		}]
	}`

	result := callTool(t, s, "build", map[string]any{
		"data":    invoiceData,
		"envelop": true,
	})

	assert.False(t, result.IsError, "build with envelop should succeed")
	text := extractText(t, result)
	assert.Contains(t, text, "envelope")
}

func TestBuildToolWithType(t *testing.T) {
	s := initServer(t)

	invoiceData := `{
		"currency": "EUR",
		"issue_date": "2024-01-15",
		"type": "standard",
		"supplier": {
			"tax_id": {"country": "ES", "code": "B91983379"},
			"name": "Test Company S.L."
		},
		"customer": {
			"tax_id": {"country": "ES", "code": "B85905495"},
			"name": "Customer S.L."
		},
		"lines": [{
			"quantity": "10",
			"item": {"name": "Test Item", "price": "100.00"},
			"taxes": [{"cat": "VAT", "rate": "standard"}]
		}]
	}`

	result := callTool(t, s, "build", map[string]any{
		"data": invoiceData,
		"type": "bill/invoice",
	})

	assert.False(t, result.IsError, "build with type should succeed")
	text := extractText(t, result)
	assert.Contains(t, text, "bill/invoice")
}

func TestBuildToolError(t *testing.T) {
	s := initServer(t)

	result := callTool(t, s, "build", map[string]any{
		"data": `{"$schema": "https://gobl.org/draft-0/bill/invoice"}`,
	})

	assert.True(t, result.IsError, "build with invalid data should fail")
}

func TestBuildToolMissingData(t *testing.T) {
	s := initServer(t)

	result := callTool(t, s, "build", map[string]any{})

	assert.True(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "data")
}

func TestValidateTool(t *testing.T) {
	s := initServer(t)

	// First build a valid envelope
	invoiceData := `{
		"$schema": "https://gobl.org/draft-0/bill/invoice",
		"currency": "EUR",
		"issue_date": "2024-01-15",
		"type": "standard",
		"supplier": {
			"tax_id": {"country": "ES", "code": "B91983379"},
			"name": "Test Company S.L."
		},
		"customer": {
			"tax_id": {"country": "ES", "code": "B85905495"},
			"name": "Customer S.L."
		},
		"lines": [{
			"quantity": "10",
			"item": {"name": "Test Item", "price": "100.00"},
			"taxes": [{"cat": "VAT", "rate": "standard"}]
		}]
	}`

	// Build first to get calculated document
	buildResult := callTool(t, s, "build", map[string]any{
		"data": invoiceData,
	})
	require.False(t, buildResult.IsError)

	// Now validate the built result
	builtText := extractText(t, buildResult)
	result := callTool(t, s, "validate", map[string]any{
		"data": builtText,
	})

	assert.False(t, result.IsError, "validate should succeed on built document")
	text := extractText(t, result)
	assert.Contains(t, text, `"ok": true`)
}

func TestValidateToolError(t *testing.T) {
	s := initServer(t)

	result := callTool(t, s, "validate", map[string]any{
		"data": `{"$schema": "https://gobl.org/draft-0/bill/invoice"}`,
	})

	assert.True(t, result.IsError, "validate with invalid data should fail")
}

func TestSchemaTool(t *testing.T) {
	s := initServer(t)

	result := callTool(t, s, "schema", map[string]any{
		"path": "bill/invoice",
	})

	assert.False(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "properties")
	assert.Contains(t, text, "invoice")
}

func TestSchemaToolNotFound(t *testing.T) {
	s := initServer(t)

	result := callTool(t, s, "schema", map[string]any{
		"path": "nonexistent/type",
	})

	assert.True(t, result.IsError)
}

func TestRegimeTool(t *testing.T) {
	s := initServer(t)

	result := callTool(t, s, "regime", map[string]any{
		"code": "ES",
	})

	assert.False(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "Spain")
}

func TestRegimeToolNotFound(t *testing.T) {
	s := initServer(t)

	result := callTool(t, s, "regime", map[string]any{
		"code": "ZZ",
	})

	assert.True(t, result.IsError)
}

func TestReplicateTool(t *testing.T) {
	s := initServer(t)

	// Build a valid document first
	invoiceData := `{
		"$schema": "https://gobl.org/draft-0/bill/invoice",
		"currency": "EUR",
		"issue_date": "2024-01-15",
		"type": "standard",
		"supplier": {
			"tax_id": {"country": "ES", "code": "B91983379"},
			"name": "Test Company S.L."
		},
		"customer": {
			"tax_id": {"country": "ES", "code": "B85905495"},
			"name": "Customer S.L."
		},
		"lines": [{
			"quantity": "10",
			"item": {"name": "Test Item", "price": "100.00"},
			"taxes": [{"cat": "VAT", "rate": "standard"}]
		}]
	}`

	buildResult := callTool(t, s, "build", map[string]any{
		"data":    invoiceData,
		"envelop": true,
	})
	require.False(t, buildResult.IsError)

	builtText := extractText(t, buildResult)
	result := callTool(t, s, "replicate", map[string]any{
		"data": builtText,
	})

	assert.False(t, result.IsError, "replicate should succeed")
	text := extractText(t, result)
	assert.Contains(t, text, "bill/invoice")
}

func TestCorrectTool(t *testing.T) {
	s := initServer(t)

	// Build a valid envelope first
	invoiceData := `{
		"$schema": "https://gobl.org/draft-0/bill/invoice",
		"currency": "EUR",
		"issue_date": "2024-01-15",
		"type": "standard",
		"supplier": {
			"tax_id": {"country": "ES", "code": "B91983379"},
			"name": "Test Company S.L."
		},
		"customer": {
			"tax_id": {"country": "ES", "code": "B85905495"},
			"name": "Customer S.L."
		},
		"lines": [{
			"quantity": "10",
			"item": {"name": "Test Item", "price": "100.00"},
			"taxes": [{"cat": "VAT", "rate": "standard"}]
		}]
	}`

	buildResult := callTool(t, s, "build", map[string]any{
		"data":    invoiceData,
		"envelop": true,
	})
	require.False(t, buildResult.IsError)
	builtText := extractText(t, buildResult)

	t.Run("correction options schema", func(t *testing.T) {
		result := callTool(t, s, "correct", map[string]any{
			"data":   builtText,
			"schema": true,
		})
		assert.False(t, result.IsError)
		text := extractText(t, result)
		assert.NotEmpty(t, text)
	})

	t.Run("create credit note", func(t *testing.T) {
		opts := `{"credit": true}`
		result := callTool(t, s, "correct", map[string]any{
			"data":    builtText,
			"options": opts,
		})
		// May succeed or fail depending on regime requirements,
		// but should not panic
		_ = result
	})

	t.Run("missing data", func(t *testing.T) {
		result := callTool(t, s, "correct", map[string]any{})
		assert.True(t, result.IsError)
		text := extractText(t, result)
		assert.Contains(t, text, "data")
	})

	t.Run("invalid data", func(t *testing.T) {
		result := callTool(t, s, "correct", map[string]any{
			"data": `{"invalid": true}`,
		})
		assert.True(t, result.IsError)
	})
}

func TestValidateToolMissingData(t *testing.T) {
	s := initServer(t)

	result := callTool(t, s, "validate", map[string]any{})
	assert.True(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "data")
}

func TestReplicateToolMissingData(t *testing.T) {
	s := initServer(t)

	result := callTool(t, s, "replicate", map[string]any{})
	assert.True(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "data")
}

func TestReplicateToolInvalidData(t *testing.T) {
	s := initServer(t)

	result := callTool(t, s, "replicate", map[string]any{
		"data": `{"invalid": true}`,
	})
	assert.True(t, result.IsError)
}

func TestSchemaToolMissingPath(t *testing.T) {
	s := initServer(t)

	result := callTool(t, s, "schema", map[string]any{})
	assert.True(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "path")
}

func TestSchemaToolWithJsonSuffix(t *testing.T) {
	s := initServer(t)

	result := callTool(t, s, "schema", map[string]any{
		"path": "bill/invoice.json",
	})
	assert.False(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "properties")
}

func TestRegimeToolMissingCode(t *testing.T) {
	s := initServer(t)

	result := callTool(t, s, "regime", map[string]any{})
	assert.True(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "code")
}

func TestRegimeToolLowercase(t *testing.T) {
	s := initServer(t)

	result := callTool(t, s, "regime", map[string]any{
		"code": "de",
	})
	assert.False(t, result.IsError)
	text := extractText(t, result)
	assert.Contains(t, text, "Germany")
}

func TestNewServer(t *testing.T) {
	srv := goblmcp.NewServer()
	require.NotNil(t, srv)

	// Verify the server responds to initialize
	initReq, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "initialize",
		"params": map[string]any{
			"protocolVersion": "2024-11-05",
			"clientInfo": map[string]any{
				"name":    "test",
				"version": "1.0",
			},
		},
	})

	resp := srv.HandleMessage(context.Background(), initReq)
	require.NotNil(t, resp)

	respJSON, err := json.Marshal(resp)
	require.NoError(t, err)

	assert.Contains(t, string(respJSON), "gobl")
	assert.Contains(t, string(respJSON), string(gobl.VERSION))
}

func extractText(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	require.NotEmpty(t, result.Content)

	// The content is []Content which can be TextContent
	// Marshal back and re-parse to get text
	d, err := json.Marshal(result.Content[0])
	require.NoError(t, err)

	var tc struct {
		Text string `json:"text"`
	}
	require.NoError(t, json.Unmarshal(d, &tc))
	return tc.Text
}
