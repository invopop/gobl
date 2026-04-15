package mcp

import (
	"context"
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/invopop/gobl/data"
	"github.com/invopop/gobl/schema"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerResources(srv *server.MCPServer) {
	// Register resource templates for browsable access
	srv.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"gobl://schemas/{path}",
			"GOBL JSON Schema",
			mcp.WithTemplateDescription("JSON Schema definitions for all GOBL types"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		handleSchemaResource,
	)

	srv.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"gobl://regimes/{code}",
			"GOBL Tax Regime",
			mcp.WithTemplateDescription("Tax regime definitions by country code"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		handleRegimeResource,
	)

	srv.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"gobl://addons/{key}",
			"GOBL Addon",
			mcp.WithTemplateDescription("Addon definitions for tax and format extensions"),
			mcp.WithTemplateMIMEType("application/json"),
		),
		handleAddonResource,
	)

	// Register static resources that list available items
	srv.AddResource(
		mcp.NewResource(
			"gobl://schemas",
			"GOBL Schema List",
			mcp.WithResourceDescription("List of all registered GOBL schema types"),
			mcp.WithMIMEType("application/json"),
		),
		handleSchemaList,
	)

	srv.AddResource(
		mcp.NewResource(
			"gobl://regimes",
			"GOBL Regime List",
			mcp.WithResourceDescription("List of all available tax regimes"),
			mcp.WithMIMEType("application/json"),
		),
		handleRegimeList,
	)

	srv.AddResource(
		mcp.NewResource(
			"gobl://addons",
			"GOBL Addon List",
			mcp.WithResourceDescription("List of all available addons"),
			mcp.WithMIMEType("application/json"),
		),
		handleAddonList,
	)
}

func handleSchemaResource(_ context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	p := extractParam(request.Params.URI, "gobl://schemas/")
	if p == "" {
		return nil, fmt.Errorf("missing schema path")
	}

	if !strings.HasSuffix(p, ".json") {
		p = p + ".json"
	}
	p = path.Join("schemas", p)

	d, err := data.Content.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("schema not found: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(d),
		},
	}, nil
}

func handleRegimeResource(_ context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	code := extractParam(request.Params.URI, "gobl://regimes/")
	if code == "" {
		return nil, fmt.Errorf("missing regime code")
	}

	p := path.Join("regimes", strings.ToLower(code)+".json")
	d, err := data.Content.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("regime not found: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(d),
		},
	}, nil
}

func handleAddonResource(_ context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	key := extractParam(request.Params.URI, "gobl://addons/")
	if key == "" {
		return nil, fmt.Errorf("missing addon key")
	}

	if !strings.HasSuffix(key, ".json") {
		key = key + ".json"
	}
	p := path.Join("addons", key)

	d, err := data.Content.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("addon not found: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "application/json",
			Text:     string(d),
		},
	}, nil
}

func handleSchemaList(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	list := schema.List()
	items := make([]string, len(list))
	for i, id := range list {
		items[i] = id.String()
	}

	d, err := marshalJSON(map[string]any{"schemas": items})
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "gobl://schemas",
			MIMEType: "application/json",
			Text:     string(d),
		},
	}, nil
}

func handleRegimeList(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	items, err := listDataDir("regimes")
	if err != nil {
		return nil, err
	}

	d, err := marshalJSON(map[string]any{"regimes": items})
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "gobl://regimes",
			MIMEType: "application/json",
			Text:     string(d),
		},
	}, nil
}

func handleAddonList(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	items, err := listDataDir("addons")
	if err != nil {
		return nil, err
	}

	d, err := marshalJSON(map[string]any{"addons": items})
	if err != nil {
		return nil, err
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "gobl://addons",
			MIMEType: "application/json",
			Text:     string(d),
		},
	}, nil
}

// listDataDir reads a directory from the embedded data filesystem and returns
// the file names without extensions.
func listDataDir(dir string) ([]string, error) {
	entries, err := fs.ReadDir(data.Content, dir)
	if err != nil {
		return nil, fmt.Errorf("reading %s directory: %w", dir, err)
	}

	items := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		items = append(items, strings.TrimSuffix(name, ".json"))
	}
	return items, nil
}

// extractParam extracts the parameter value from a URI by stripping the prefix.
func extractParam(uri, prefix string) string {
	if !strings.HasPrefix(uri, prefix) {
		return ""
	}
	return strings.TrimPrefix(uri, prefix)
}
