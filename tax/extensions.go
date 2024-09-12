package tax

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Extensions are a key component of GOBL that are used to include additional
// structured data in documents that doesn't fit into any of the common or
// universal fields. They're typically defined by local tax agencies that will
// use the data for tax reports or classification. Civil law countries
// have a far greater tendancy to require these than common law countries.
//
// Naming of extension keys is important and should be kept short and descriptive.
// There are three key components to an extension key separated by dashes:
//
// - An ISO country code e.g. `mx`, `es`, `gb`, etc.
// - Short abreviation of the platform or format the extension will be used with,
//   e.g. `cfdi` for Mexico's CFDI defined by the SAT, `facturae` for Spain's
//   FacturaE format, `sdi` for the Italian SDI (document interchange system), etc.
//   This is important, as it helps avoid potential conflicts in the future with
//   new or alternative formats that may appear.
// - A short descriptive name of the extension, e.g. `exception`, `fiscal-regime`,
//   `vat-cat`, `incoming-typ`, etc. The aim should be to avoid using obvious names
//   like `code` or `key` in the name, as these are already implied through usage.
//
// Please look at the regimes package and othe country specific implementations for
// examples of how to define and use extensions.

// Extensions is a map of extension keys to values.
type Extensions map[cbc.Key]ExtValue

// ExtValue is a string value that has helper methods to help determine
// if it is a code, key, or regular string.
type ExtValue string

type extensionCollection struct {
	list map[cbc.Key]*cbc.KeyDefinition
}

// extensionsDefs is a global register of all extension definitions
// that have been registered via regimes and addons.
var extensionDefs = newExtensionCollection()

func newExtensionCollection() *extensionCollection {
	return &extensionCollection{
		list: make(map[cbc.Key]*cbc.KeyDefinition),
	}
}

func (c *extensionCollection) add(kd *cbc.KeyDefinition) {
	c.list[kd.Key] = kd
}

// RegisterExtension is used to add any extension definitions to the global
// register. This is not expected to be called directly, but rather will
// be used by the regimes and addons during their registration processes.
func RegisterExtension(kd *cbc.KeyDefinition) {
	extensionDefs.add(kd)
}

// ExtensionForKey returns the extension definition for the given key or nil.
func ExtensionForKey(key cbc.Key) *cbc.KeyDefinition {
	return extensionDefs.list[key]
}

// Validate ensures the extension map data looks correct and that all keys
// have been registered globally.
func (em Extensions) Validate() error {
	err := make(validation.Errors)
	// Validate key format
	for k := range em {
		if e := k.Validate(); e != nil {
			err[k.String()] = e
		}
	}
	if len(err) > 0 {
		return err
	}
	// Validate keys are defined and correct
	for k, ev := range em {
		ks := k.String()
		kd := ExtensionForKey(k)
		if kd == nil {
			err[ks] = errors.New("undefined")
			continue
		}
		if len(kd.Values) > 0 && !kd.HasValue(ev.String()) {
			err[ks] = fmt.Errorf("value '%s' invalid", ev)
		}
		if kd.Pattern != "" {
			re, rerr := regexp.Compile(kd.Pattern)
			if rerr != nil {
				err[ks] = rerr
				continue
			}
			if !re.MatchString(string(ev)) {
				err[ks] = errors.New("does not match pattern")
			}
		}
	}
	if len(err) > 0 {
		return err
	}
	return nil
}

// Has returns true if the code map has values for all the provided keys.
func (em Extensions) Has(keys ...cbc.Key) bool {
	for _, k := range keys {
		if _, ok := em[k]; !ok {
			return false
		}
	}
	return true
}

// Equals returns true if the extension map has the same keys and values as the provided
// map.
func (em Extensions) Equals(other Extensions) bool {
	if len(em) != len(other) {
		return false
	}
	if len(em) == 0 {
		return true // empty extensions are equal!
	}
	return em.Contains(other)
}

// Contains returns true if the extension map contains the same keys and values as the provided
// map, but may have additional keys.
func (em Extensions) Contains(other Extensions) bool {
	if len(em) == 0 {
		return false
	}
	for k, v := range other {
		v2, ok := em[k]
		if !ok {
			return false
		}
		if v2 != v {
			return false
		}
	}
	return true
}

// Merge will merge the provided extensions map with the current one generating
// a new map. Duplicate keys will be overwritten by the other map's values.
func (em Extensions) Merge(other Extensions) Extensions {
	if em == nil {
		return other
	}
	if other == nil {
		return em
	}
	nem := make(Extensions)
	for k, v := range em {
		nem[k] = v
	}
	for k, v := range other {
		nem[k] = v
	}
	return nem
}

// NormalizeExtensions will try to clean the extension map removing empty values
// and will potentially return a nil if there only keys with no values.
func NormalizeExtensions(em Extensions) Extensions {
	if em == nil {
		return nil
	}
	nem := make(Extensions)
	for k, v := range em {
		if v == "" {
			continue
		}
		nem[k] = v
	}
	if len(nem) == 0 {
		return nil
	}
	return nem
}

// ExtensionsHas returns a validation rule that ensures the extension map's
// keys match those provided.
func ExtensionsHas(keys ...cbc.Key) validation.Rule {
	return validateCodeMap{keys: keys}
}

// ExtensionsRequires returns a validation rule that ensures all the
// extension map's keys match those provided in the list.
func ExtensionsRequires(keys ...cbc.Key) validation.Rule {
	return validateCodeMap{
		required: true,
		keys:     keys,
	}
}

type validateCodeMap struct {
	keys     []cbc.Key
	required bool
}

func (v validateCodeMap) Validate(value interface{}) error {
	em, ok := value.(Extensions)
	if !ok {
		return nil
	}
	err := make(validation.Errors)

	if v.required {
		for _, k := range v.keys {
			if _, ok := em[k]; !ok {
				err[k.String()] = errors.New("required")
			}
		}
	} else {
		for k := range em {
			if !k.In(v.keys...) {
				err[k.String()] = errors.New("invalid")
			}
		}
	}

	if len(err) > 0 {
		return err
	}
	return nil
}

// JSONSchemaExtend provides extra details about the extension map which are
// not automatically determined. In this case we add validation for the map's
// keys.
func (Extensions) JSONSchemaExtend(schema *jsonschema.Schema) {
	prop := schema.AdditionalProperties
	schema.AdditionalProperties = nil
	schema.PatternProperties = map[string]*jsonschema.Schema{
		cbc.KeyPattern: prop,
	}
}

// String provides the string representation.
func (ev ExtValue) String() string {
	return string(ev)
}

// Key returns the key value or empty if the value is a Code.
func (ev ExtValue) Key() cbc.Key {
	k := cbc.Key(ev)
	if err := k.Validate(); err == nil {
		return k
	}
	return cbc.KeyEmpty
}

// Code returns the code value or empty if the value is a Key.
func (ev ExtValue) Code() cbc.Code {
	c := cbc.Code(ev)
	if err := c.Validate(); err == nil {
		return c
	}
	return cbc.CodeEmpty
}
