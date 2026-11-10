package gobl_test

import (
	"path"
	"reflect"
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestExtensionDefs(t *testing.T) {
	for _, ed := range tax.AllExtensionDefs() {
		t.Run(ed.Key.String(), func(t *testing.T) {
			assert.NoError(t, rules.Validate(ed))
			for lang, name := range ed.Name {
				assert.Contains(t, name, ": ", "name in '%s' must be prefixed with the issuer", lang)
			}
			if !assert.NotNil(t, ed.Usage) {
				return
			}
			for _, k := range ed.Usage.Keys {
				assert.NotNil(t, tax.AddonForKey(k), "addon '%s' not registered", k)
			}
			for _, c := range ed.Usage.Codes {
				assert.NotNil(t, tax.RegimeDefFor(l10n.Code(c)), "regime '%s' not registered", c)
			}
			for _, s := range ed.Usage.Schemas {
				typ := schema.Type(schema.GOBL.Add(s))
				if assert.NotNil(t, typ, "schema '%s' not registered", s) {
					assert.True(t, hasExtensions(typ), "schema '%s' does not have extensions", s)
				}
			}
		})
	}
}

func TestRegimeExtensionUsage(t *testing.T) {
	for _, rd := range tax.AllRegimeDefs() {
		for _, ed := range rd.Extensions {
			t.Run(ed.Key.String(), func(t *testing.T) {
				if assert.NotNil(t, ed.Usage) {
					assert.Contains(t, ed.Usage.Codes, cbc.Code(rd.Country))
				}
			})
		}
	}
}

func TestAddonExtensionUsage(t *testing.T) {
	for _, ad := range tax.AllAddonDefs() {
		for _, ed := range ad.Extensions {
			t.Run(ed.Key.String(), func(t *testing.T) {
				if assert.NotNil(t, ed.Usage) {
					assert.Contains(t, ed.Usage.Keys, ad.Key)
				}
			})
		}
	}
}

func TestExtensionSchemasRegistered(t *testing.T) {
	// Definitions describe extensions rather than containing them.
	seen := map[reflect.Type]bool{
		reflect.TypeOf(tax.RegimeDef{}):    true,
		reflect.TypeOf(tax.AddonDef{}):     true,
		reflect.TypeOf(tax.CatalogueDef{}): true,
	}
	var walk func(typ reflect.Type)
	walk = func(typ reflect.Type) {
		for typ.Kind() == reflect.Pointer || typ.Kind() == reflect.Slice || typ.Kind() == reflect.Map {
			typ = typ.Elem()
		}
		if typ.Kind() != reflect.Struct || typ == extensionsType || seen[typ] {
			return
		}
		seen[typ] = true
		if hasExtensions(typ) {
			id := schema.Lookup(reflect.New(typ).Interface())
			assert.NotEqual(t, schema.UnknownID, id, "%s.%s has extensions but no schema", path.Base(typ.PkgPath()), typ.Name())
		}
		for i := range typ.NumField() {
			walk(typ.Field(i).Type)
		}
	}
	for typ := range schema.Types() {
		walk(typ)
	}
}

var extensionsType = reflect.TypeOf(tax.Extensions{})

// hasExtensions returns true if the struct type has an extensions field.
func hasExtensions(typ reflect.Type) bool {
	if typ.Kind() != reflect.Struct {
		return false
	}
	for i := range typ.NumField() {
		if typ.Field(i).Type == extensionsType {
			return true
		}
	}
	return false
}
