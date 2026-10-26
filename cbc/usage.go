package cbc

import (
	"slices"
	"strings"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// UsageSchemaPattern defines the format of the schema IDs used in a Usage,
// relative to the GOBL base, such as "tax/combo" or "org/party".
const UsageSchemaPattern = `^[a-z][a-z0-9-]*/[a-z][a-z0-9-]*$`

// UsagePathPattern defines the format of the JSON paths used in a Usage:
// a "$" root followed by property names, with "[*]" to match any array index.
const UsagePathPattern = `^\$(\.[$a-zA-Z0-9_-]+(\[\*\])*)*$`

// Usage describes the context in which a definition may be used. Keys and
// codes are generic filters, such as the keys of the addons or the country
// codes of the tax regimes that offer the definition. Empty fields place no
// restriction. Usage describes context for filtering and does not restrict
// validation.
type Usage struct {
	// Keys used for filtering, such as the addons that offer the definition.
	Keys []Key `json:"keys,omitempty" jsonschema:"title=Keys"`
	// Codes used for filtering, such as the tax regime country codes that offer the
	// definition.
	Codes []Code `json:"codes,omitempty" jsonschema:"title=Codes"`
	// Schema IDs, relative to the GOBL base, of the objects in which the definition
	// may be used, such as "tax/combo" or "org/party".
	Schemas []string `json:"schemas,omitempty" jsonschema:"title=Schemas"`
	// JSON paths inside a document where the definition may be used, with "[*]"
	// matching any array index, for when only some of the schema's locations apply.
	Paths []string `json:"paths,omitempty" jsonschema:"title=Paths"`
}

func usageRules() *rules.Set {
	return rules.For(new(Usage),
		rules.Field("schemas",
			rules.Each(
				rules.Assert("01", "schema must be a relative schema ID", is.Matches(UsageSchemaPattern)),
			),
		),
		rules.Field("paths",
			rules.Each(
				rules.Assert("02", "path must be a JSON path", is.Matches(UsagePathPattern)),
			),
		),
	)
}

// Matches returns true if the usage has no keys or codes, or if any of the
// provided keys or codes is listed.
func (u *Usage) Matches(keys []Key, codes []Code) bool {
	if u == nil || (len(u.Keys) == 0 && len(u.Codes) == 0) {
		return true
	}
	for _, k := range keys {
		if k.In(u.Keys...) {
			return true
		}
	}
	for _, c := range codes {
		if c.In(u.Codes...) {
			return true
		}
	}
	return false
}

// MatchesSchema returns true if the usage has no schemas, or if the provided
// schema ID, either relative like "tax/combo" or complete, matches one of them.
func (u *Usage) MatchesSchema(schema string) bool {
	if u == nil || len(u.Schemas) == 0 {
		return true
	}
	return slices.ContainsFunc(u.Schemas, func(s string) bool {
		return schema == s || strings.HasSuffix(schema, "/"+s)
	})
}

// MatchesPath returns true if the usage has no paths, or if the provided
// JSON path, such as "$.lines[0].taxes[1].ext", matches one of them.
func (u *Usage) MatchesPath(path string) bool {
	if u == nil || len(u.Paths) == 0 {
		return true
	}
	return slices.ContainsFunc(u.Paths, func(p string) bool {
		return matchUsagePath(p, path)
	})
}

// matchUsagePath compares the path with the pattern, where each "[*]" in the
// pattern matches any array index, or "[*]", in the path.
func matchUsagePath(pattern, path string) bool {
	for {
		i := strings.Index(pattern, "[*]")
		if i < 0 {
			return pattern == path
		}
		if !strings.HasPrefix(path, pattern[:i+1]) {
			return false
		}
		path = path[i+1:]
		j := strings.IndexByte(path, ']')
		if j < 0 || !validUsageIndex(path[:j]) {
			return false
		}
		pattern = pattern[i+3:]
		path = path[j+1:]
	}
}

func validUsageIndex(s string) bool {
	if s == "*" {
		return true
	}
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
