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
type Extensions map[cbc.Key]cbc.Code

type extensionCollection struct {
	list map[cbc.Key]*cbc.Definition
}

// extensionsDefs is a global register of all extension definitions
// that have been registered via regimes and addons.
var extensionDefs = newExtensionCollection()

func newExtensionCollection() *extensionCollection {
	return &extensionCollection{
		list: make(map[cbc.Key]*cbc.Definition),
	}
}

func (c *extensionCollection) add(kd *cbc.Definition) {
	c.list[kd.Key] = kd
}

// RegisterExtension is used to add any extension definitions to the global
// register. This is not expected to be called directly, but rather will
// be used by the regimes and addons during their registration processes.
func RegisterExtension(kd *cbc.Definition) {
	extensionDefs.add(kd)
}

// ExtensionForKey returns the extension definition for the given key or nil.
func ExtensionForKey(key cbc.Key) *cbc.Definition {
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
		if len(kd.Values) > 0 && !kd.HasCode(ev) {
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

// Get returns the value for the provided key or an empty string if not found
// or the extensions map is nil. If the key is composed of sub-keys and
// no precise match is found, the key will be split until one of the sub
// components is found.
func (em Extensions) Get(k cbc.Key) cbc.Code {
	if len(em) == 0 {
		return ""
	}
	// while k is not empty, pop the last key and check if it exists
	for k != cbc.KeyEmpty {
		if v, ok := em[k]; ok {
			return v
		}
		k = k.Pop()
	}
	return ""
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

// Lookup returns the key for the provided value or an empty
// key if not found. This is useful for reverse lookups.
func (em Extensions) Lookup(val cbc.Code) cbc.Key {
	for k, v := range em {
		if v == val {
			return k
		}
	}
	return cbc.KeyEmpty
}

// CleanExtensions will try to clean the extension map removing empty values
// and will potentially return a nil if there only keys with no values.
func CleanExtensions(em Extensions) Extensions {
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

// ExtensionsRequire returns a validation rule that ensures that all of
// the provided keys are present.
func ExtensionsRequire(keys ...cbc.Key) validation.Rule {
	return validateExtCodeMap{
		operator: extCodeOpAnd,
		keys:     keys,
	}
}

// ExtensionsRequireAllOrNone returns a validation rule that performs an XNOR
// operation on the provided keys. If one of the keys is present, then
// all of them must be present. If none of the keys are present,
// then all of them must be absent.
func ExtensionsRequireAllOrNone(keys ...cbc.Key) validation.Rule {
	return validateExtCodeMap{
		operator: extCodeOpXNOr,
		keys:     keys,
	}
}

// ExtensionsExclude returns a validation rule that ensures that
// an extensions map does **not** include the provided keys.
func ExtensionsExclude(keys ...cbc.Key) validation.Rule {
	return validateExtCodeMap{
		operator: extCodeOpNot,
		keys:     keys,
	}
}

type extCodeOp int

const (
	_                      = iota
	extCodeOpAnd extCodeOp = 1 + iota
	extCodeOpNot
	extCodeOpXNOr
)

type validateExtCodeMap struct {
	keys     []cbc.Key
	operator extCodeOp
}

func (v validateExtCodeMap) Validate(value interface{}) error {
	em, ok := value.(Extensions)
	if !ok {
		return nil
	}
	err := make(validation.Errors)

	switch v.operator {
	case extCodeOpAnd:
		for _, k := range v.keys {
			if _, ok := em[k]; !ok {
				err[k.String()] = errors.New("required")
			}
		}
	case extCodeOpNot:
		for _, k := range v.keys {
			if _, ok := em[k]; ok {
				err[k.String()] = errors.New("must be blank")
			}
		}
	case extCodeOpXNOr:
		present := 0
		for _, k := range v.keys {
			if _, ok := em[k]; ok {
				present++
			}
		}
		if present > 0 && present != len(v.keys) {
			for _, k := range v.keys {
				if _, ok := em[k]; !ok {
					err[k.String()] = errors.New("required")
				}
			}
		}
	}

	if len(err) > 0 {
		return err
	}
	return nil
}

// ExtensionsHasCodes returns a validation rule that ensures the extension map's
// key has one of the provided **codes**.
func ExtensionsHasCodes(key cbc.Key, codes ...cbc.Code) validation.Rule {
	return validateExtCodes{
		key:       key,
		values:    codes,
		inclusion: true,
	}
}

// ExtensionsExcludeCodes returns a validation rule that ensures the extension map's
// key does not have any of the provided **codes**.
func ExtensionsExcludeCodes(key cbc.Key, codes ...cbc.Code) validation.Rule {
	return validateExtCodes{
		key:       key,
		values:    codes,
		inclusion: false,
	}
}

type validateExtCodes struct {
	key       cbc.Key
	values    []cbc.Code
	inclusion bool
}

func (v validateExtCodes) Validate(value interface{}) error {
	em, ok := value.(Extensions)
	if !ok {
		return nil
	}
	err := make(validation.Errors)

	if ev, ok := em[v.key]; ok {
		if v.inclusion {
			// Inclusion mode: value must be in the list
			match := false
			for _, val := range v.values {
				if ev == val {
					match = true
					break
				}
			}
			if !match {
				err[v.key.String()] = errors.New("invalid value")
			}
		} else {
			// Exclusion mode: value must NOT be in the list
			for _, val := range v.values {
				if ev == val {
					err[v.key.String()] = fmt.Errorf("value '%s' not allowed", ev)
					break
				}
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
