package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/invopop/gobl/data"
	"github.com/invopop/gobl/internal/cli"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTools(srv *server.MCPServer) {
	srv.AddTool(buildTool(), handleBuild)
	srv.AddTool(validateTool(), handleValidate)
	srv.AddTool(correctTool(), handleCorrect)
	srv.AddTool(replicateTool(), handleReplicate)
	srv.AddTool(schemaTool(), handleSchema)
	srv.AddTool(regimeTool(), handleRegime)
}

func buildTool() mcp.Tool {
	return mcp.NewTool("build",
		mcp.WithDescription("Calculate and validate a GOBL document. Provide partial or complete JSON/YAML and get back a fully calculated envelope with totals, tax calculations, and validation. This is the primary tool for building GOBL documents iteratively."),
		mcp.WithString("data",
			mcp.Required(),
			mcp.Description("JSON or YAML document data"),
		),
		mcp.WithString("type",
			mcp.Description("Document type hint (e.g. 'bill/invoice', 'org/party')"),
		),
		mcp.WithBoolean("envelop",
			mcp.Description("Wrap the document in a GOBL envelope"),
		),
	)
}

func handleBuild(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	dataStr, _ := args["data"].(string)
	if dataStr == "" {
		return mcp.NewToolResultError("'data' argument is required"), nil
	}

	opts := &cli.BuildOptions{
		ParseOptions: &cli.ParseOptions{
			Input: bytes.NewReader([]byte(dataStr)),
		},
	}

	if docType, ok := args["type"].(string); ok && docType != "" {
		opts.DocType = docType
	}
	if envelop, ok := args["envelop"].(bool); ok {
		opts.Envelop = envelop
	}

	result, err := cli.Build(ctx, opts)
	if err != nil {
		return toolResultFromError(err), nil
	}

	return toolResultFromJSON(result)
}

func validateTool() mcp.Tool {
	return mcp.NewTool("validate",
		mcp.WithDescription("Validate a GOBL document without modifying it. Returns OK or structured validation errors with faults."),
		mcp.WithString("data",
			mcp.Required(),
			mcp.Description("JSON document data to validate"),
		),
	)
}

func handleValidate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	dataStr, _ := args["data"].(string)
	if dataStr == "" {
		return mcp.NewToolResultError("'data' argument is required"), nil
	}

	err := cli.Validate(ctx, bytes.NewReader([]byte(dataStr)))
	if err != nil {
		return toolResultFromError(err), nil
	}

	return mcp.NewToolResultText(`{"ok": true}`), nil
}

func correctTool() mcp.Tool {
	return mcp.NewTool("correct",
		mcp.WithDescription("Create corrective documents (credit notes, debit notes) from an existing invoice. Can also return the available correction options schema for a given document."),
		mcp.WithString("data",
			mcp.Required(),
			mcp.Description("JSON document data of the original invoice"),
		),
		mcp.WithString("options",
			mcp.Description("JSON correction options (e.g. stamps, reason, method)"),
		),
		mcp.WithBoolean("schema",
			mcp.Description("When true, return available correction options instead of correcting"),
		),
	)
}

func handleCorrect(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	dataStr, _ := args["data"].(string)
	if dataStr == "" {
		return mcp.NewToolResultError("'data' argument is required"), nil
	}

	opts := &cli.CorrectOptions{
		ParseOptions: &cli.ParseOptions{
			Input: bytes.NewReader([]byte(dataStr)),
		},
	}

	if optionsStr, ok := args["options"].(string); ok && optionsStr != "" {
		opts.Data = []byte(optionsStr)
	}
	if schema, ok := args["schema"].(bool); ok {
		opts.OptionsSchema = schema
	}

	result, err := cli.Correct(ctx, opts)
	if err != nil {
		return toolResultFromError(err), nil
	}

	return toolResultFromJSON(result)
}

func replicateTool() mcp.Tool {
	return mcp.NewTool("replicate",
		mcp.WithDescription("Clone a GOBL document as a new template with a fresh UUID. Clears stamps and signatures from the original."),
		mcp.WithString("data",
			mcp.Required(),
			mcp.Description("JSON document data to replicate"),
		),
	)
}

func handleReplicate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	dataStr, _ := args["data"].(string)
	if dataStr == "" {
		return mcp.NewToolResultError("'data' argument is required"), nil
	}

	opts := &cli.ReplicateOptions{
		ParseOptions: &cli.ParseOptions{
			Input: bytes.NewReader([]byte(dataStr)),
		},
	}

	result, err := cli.Replicate(ctx, opts)
	if err != nil {
		return toolResultFromError(err), nil
	}

	return toolResultFromJSON(result)
}

func schemaTool() mcp.Tool {
	return mcp.NewTool("schema",
		mcp.WithDescription("Look up the JSON Schema definition for any GOBL type (e.g. 'bill/invoice', 'org/party', 'pay/instructions')."),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Schema path (e.g. 'bill/invoice', 'org/party')"),
		),
	)
}

func handleSchema(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	p, _ := args["path"].(string)
	if p == "" {
		return mcp.NewToolResultError("'path' argument is required"), nil
	}

	ext := filepath.Ext(p)
	if ext == "" {
		p = p + ".json"
	}
	p = path.Join("schemas", p)

	d, err := data.Content.ReadFile(p)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("schema not found: %s", err)), nil
	}

	return mcp.NewToolResultText(string(d)), nil
}

func regimeTool() mcp.Tool {
	return mcp.NewTool("regime",
		mcp.WithDescription("Get the full tax regime definition for a country. Returns tax categories, rates, extensions, scenarios, and correction definitions."),
		mcp.WithString("code",
			mcp.Required(),
			mcp.Description("Country code (e.g. 'ES', 'DE', 'MX', 'GB')"),
		),
	)
}

func handleRegime(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	code, _ := args["code"].(string)
	if code == "" {
		return mcp.NewToolResultError("'code' argument is required"), nil
	}

	p := path.Join("regimes", strings.ToLower(code)+".json")
	d, err := data.Content.ReadFile(p)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("regime not found: %s", err)), nil
	}

	return mcp.NewToolResultText(string(d)), nil
}

// toolResultFromJSON marshals the result as indented JSON and returns it as text content.
func toolResultFromJSON(v any) (*mcp.CallToolResult, error) {
	d, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, fmt.Errorf("marshaling result: %w", err)
	}
	return mcp.NewToolResultText(string(d)), nil
}

// toolResultFromError converts a CLI error into an MCP tool error result,
// preserving structured fault information.
func toolResultFromError(err error) *mcp.CallToolResult {
	if cliErr, ok := err.(*cli.Error); ok {
		d, marshalErr := json.MarshalIndent(cliErr, "", "\t")
		if marshalErr == nil {
			return mcp.NewToolResultError(string(d))
		}
	}
	return mcp.NewToolResultError(err.Error())
}
