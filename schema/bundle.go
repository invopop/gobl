package schema

import (
	"encoding/json"
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/invopop/gobl/data"
)

// schemaFile is the minimal structure parsed from each schema JSON file.
type schemaFile struct {
	Schema string                     `json:"$schema,omitempty"`
	ID     string                     `json:"$id,omitempty"`
	Ref    string                     `json:"$ref,omitempty"`
	Defs   map[string]json.RawMessage `json:"$defs,omitempty"`
}

// refPattern matches "$ref": "https://gobl.org/..." values.
var refPattern = regexp.MustCompile(`"\$ref"\s*:\s*"(` + regexp.QuoteMeta(GOBL.String()) + `[^"]*)"`)

// BundleSchema reads the schema at the given embedded path (e.g.
// "schemas/bill/invoice.json") and returns a self-contained version
// with all transitive external dependencies inlined into $defs.
func BundleSchema(schemaPath string) ([]byte, error) {
	// Load the root schema file.
	root, err := loadSchemaFile(schemaPath)
	if err != nil {
		return nil, err
	}

	// Collect all transitive dependencies.
	// urlToDef maps external $id URLs to definition key names.
	urlToDef := map[string]string{root.ID: strings.TrimPrefix(root.Ref, "#/$defs/")}
	// allDefs accumulates every definition we need.
	allDefs := make(map[string]json.RawMessage)
	for name, def := range root.Defs {
		allDefs[name] = def
	}

	// Seed the work queue with external refs from the root definitions.
	queue := findExternalRefs(root.Defs)
	visited := map[string]bool{root.ID: true}

	for len(queue) > 0 {
		url := queue[0]
		queue = queue[1:]
		if visited[url] {
			continue
		}
		visited[url] = true

		dep, err := loadSchemaByURL(url)
		if err != nil {
			continue // skip unresolvable refs
		}

		mainDef := strings.TrimPrefix(dep.Ref, "#/$defs/")
		if dep.ID != "" && mainDef != "" {
			urlToDef[dep.ID] = mainDef
		}

		for name, def := range dep.Defs {
			if _, exists := allDefs[name]; !exists {
				allDefs[name] = def
			}
		}

		// Find further external refs from the newly added definitions.
		queue = append(queue, findExternalRefs(dep.Defs)...)
	}

	// Rewrite all external $ref URLs to internal #/$defs/ references.
	rewritten := make(map[string]json.RawMessage, len(allDefs))
	for name, raw := range allDefs {
		rewritten[name] = rewriteRefs(raw, urlToDef)
	}

	// Rebuild the root schema with all definitions merged.
	root.Defs = rewritten

	return json.MarshalIndent(root, "", "  ")
}

// loadSchemaFile reads and parses a schema from the embedded FS.
func loadSchemaFile(p string) (*schemaFile, error) {
	raw, err := data.Content.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", p, err)
	}
	var s schemaFile
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", p, err)
	}
	return &s, nil
}

// loadSchemaByURL resolves an external GOBL schema URL to a file and loads it.
// e.g. "https://gobl.org/draft-0/bill/line" -> "schemas/bill/line.json"
func loadSchemaByURL(url string) (*schemaFile, error) {
	suffix := strings.TrimPrefix(url, GOBL.String())
	p := path.Join("schemas", suffix+".json")
	return loadSchemaFile(p)
}

// findExternalRefs scans all definitions for external $ref URLs.
func findExternalRefs(defs map[string]json.RawMessage) []string {
	var refs []string
	prefix := GOBL.String()
	for _, raw := range defs {
		for _, m := range refPattern.FindAllSubmatch(raw, -1) {
			url := string(m[1])
			if strings.HasPrefix(url, prefix) {
				refs = append(refs, url)
			}
		}
	}
	return refs
}

// rewriteRefs replaces external $ref URLs with internal #/$defs/ references.
func rewriteRefs(raw json.RawMessage, urlToDef map[string]string) json.RawMessage {
	s := string(raw)
	for url, defName := range urlToDef {
		s = strings.ReplaceAll(s, `"`+url+`"`, `"#/$defs/`+defName+`"`)
	}
	return json.RawMessage(s)
}
