package cli

import (
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/invopop/gobl/schema"
)

// FindType can be used to match a partial schema name with a complete schema ID.
func FindType(term string) schema.ID {
	return findType(schema.Types(), term)
}

func findType(types map[reflect.Type]schema.ID, term string) schema.ID {
	schema := toSchema(term)
	for typ, id := range types {
		if term == string(id) {
			return id
		}
		if term == typ.Name() {
			return id
		}
		if term == filepath.Base(typ.PkgPath())+"."+typ.Name() {
			return id
		}
		if strings.HasSuffix(string(id), schema) {
			return id
		}
	}
	return ""
}

var singleCap = regexp.MustCompile("([A-Z])")
var allCaps = regexp.MustCompile("[A-Z]{2,}")

func toSchema(term string) string {
	if strings.HasPrefix(term, "http://") || strings.HasPrefix(term, "https://") {
		return term
	}
	if strings.Contains(term, ".") {
		parts := strings.Split(term, ".")
		for i, part := range parts {
			parts[i] = strings.TrimPrefix(toSchema(part), "-")
		}
		return "/" + strings.Join(parts, "/")
	}
	for _, match := range allCaps.FindAllString(term, -1) {
		match = match[:len(match)-1]
		n := strings.ToLower(match)
		term = strings.Replace(term, match, n, 1)
	}
	term = singleCap.ReplaceAllString(term, "-${1}")
	return strings.ToLower(term)
}
