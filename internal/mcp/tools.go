package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/data"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/internal/cli"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
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
	srv.AddTool(regimeListTool(), handleRegimeListTool)
	srv.AddTool(addonTool(), handleAddonTool)
	srv.AddTool(addonListTool(), handleAddonListTool)
}

func buildTool() mcp.Tool {
	return mcp.NewTool("build",
		mcp.WithDescription(here.Doc(`
			Calculate and validate a GOBL document. Provide partial or complete JSON and get back a fully calculated document with totals, tax calculations, and validation.

			Documents like invoices use the ~$regime~ property with a country code (e.g. "ES", "DE", "MX") to determine which tax rules apply. If no ~$regime~ is available for a given country, GOBL will continue to function. Use 'regime_list' to see available regimes.

			Some documents may use ~$addons~ — an array of addon keys that enable format-specific rules (e.g. ["es-verifactu-v1"] for Spanish Verifactu, ["mx-cfdi-v4"] for Mexican CFDI). Use 'addon_list' to see available addons, then 'addon' to check what extensions, scenarios, and rules will be applied by the addon.

			Minimal invoice example:
			{"$schema":"https://gobl.org/draft-0/bill/invoice","$regime":"ES","$addons":["es-verifactu-v1"],"currency":"EUR","issue_date":"2024-01-15","supplier":{"tax_id":{"country":"ES","code":"B85905495"},"name":"Seller SL"},"customer":{"tax_id":{"country":"ES","code":"B85905495"},"name":"Buyer SL"},"lines":[{"quantity":"1","item":{"name":"Service","price":"100.00"},"taxes":[{"cat":"VAT","rate":"standard"}]}]}

			The "$regime" is usually inferred from the supplier's tax_id country but should be set explicitly. The "$addons" array is only needed when the target format requires specific extensions or scenarios beyond the base regime rules.
		`)),
		mcp.WithString("data",
			mcp.Required(),
			mcp.Description("JSON document data. Must include '$schema' and typically '$regime'. May include '$addons' array for format-specific rules."),
		),
		mcp.WithString("type",
			mcp.Description("Document type hint (e.g. 'bill/invoice', 'org/party'). Only needed when '$schema' is not set in the data."),
		),
		mcp.WithBoolean("envelop",
			mcp.Description("Wrap the document in a GOBL envelope (adds head, signatures support). Default false."),
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
		mcp.WithDescription("Validate a GOBL document without modifying it. Returns OK or structured validation errors with faults. The document must include '$schema' and '$regime', and any required '$addons'."),
		mcp.WithString("data",
			mcp.Required(),
			mcp.Description("JSON document data to validate. Should be a fully built document (output from 'build')."),
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
		mcp.WithDescription("Look up the JSON Schema definition for any GOBL type. Use this to understand the structure and required fields for documents. The '$schema' value in a GOBL document corresponds to 'https://gobl.org/draft-0/' + the path (e.g. 'bill/invoice' -> '$schema': 'https://gobl.org/draft-0/bill/invoice')."),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Schema path (e.g. 'bill/invoice', 'org/party', 'pay/instructions'). Use 'schemas' resource to list all available paths."),
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
		mcp.WithDescription(`Get the full tax regime definition for a country. Returns tax categories (e.g. VAT rates), extensions, scenarios, and correction definitions. Use this to understand what tax categories and rates are available when building invoice lines with the 'taxes' array. The regime's country code is what you set in the document's "$regime" property.`),
		mcp.WithString("code",
			mcp.Required(),
			mcp.Description("Country code (e.g. 'ES', 'DE', 'MX', 'GB'). Use 'regime_list' to see all available codes."),
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

type regimeSummary struct {
	Country     string      `json:"country"`
	Name        i18n.String `json:"name"`
	Description i18n.String `json:"description,omitempty"`
	Currency    string      `json:"currency"`
}

func regimeListTool() mcp.Tool {
	return mcp.NewTool("regime_list",
		mcp.WithDescription(`List all available tax regimes with their country codes, names, and currencies. Each regime's country code can be used as the "$regime" value in a GOBL document. Use the 'regime' tool with a specific country code to see its tax categories, rates, and rules.`),
	)
}

func handleRegimeListTool(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	defs := tax.AllRegimeDefs()
	items := make([]regimeSummary, len(defs))
	for i, r := range defs {
		items[i] = regimeSummary{
			Country:     string(r.Country),
			Name:        r.Name,
			Description: r.Description,
			Currency:    string(r.Currency),
		}
	}
	return toolResultFromJSON(map[string]any{"regimes": items})
}

type addonSummary struct {
	Key         string      `json:"key"`
	Name        i18n.String `json:"name"`
	Description i18n.String `json:"description,omitempty"`
	Requires    []cbc.Key   `json:"requires,omitempty"`
}

func addonTool() mcp.Tool {
	return mcp.NewTool("addon",
		mcp.WithDescription(`Get the full addon definition for a given key. Returns the extensions, scenarios, and validation rules that the addon applies. Use this to understand what extra fields or constraints an addon requires — for example, what extension codes must be set on tax lines or payment instructions. Addon keys are added to a document's "$addons" array (e.g. "$addons": ["es-verifactu-v1"]).`),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("Addon key (e.g. 'es-verifactu-v1', 'mx-cfdi-v4', 'eu-en16931-v2017'). Use 'addon_list' to see all available keys."),
		),
	)
}

func handleAddonTool(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()

	key, _ := args["key"].(string)
	if key == "" {
		return mcp.NewToolResultError("'key' argument is required"), nil
	}

	if !strings.HasSuffix(key, ".json") {
		key = key + ".json"
	}
	p := path.Join("addons", key)

	d, err := data.Content.ReadFile(p)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("addon not found: %s", err)), nil
	}

	return mcp.NewToolResultText(string(d)), nil
}

func addonListTool() mcp.Tool {
	return mcp.NewTool("addon_list",
		mcp.WithDescription(`List all available addons with their keys, names, descriptions, and dependencies. Addons extend a regime with format-specific rules (e.g. Spain's Verifactu, Mexico's CFDI, EU's EN16931). To use an addon, add its key to the document's "$addons" array. Use the 'addon' tool to inspect what a specific addon requires.`),
	)
}

func handleAddonListTool(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	defs := tax.AllAddonDefs()
	items := make([]addonSummary, len(defs))
	for i, a := range defs {
		items[i] = addonSummary{
			Key:         string(a.Key),
			Name:        a.Name,
			Description: a.Description,
			Requires:    a.Requires,
		}
	}
	return toolResultFromJSON(map[string]any{"addons": items})
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
