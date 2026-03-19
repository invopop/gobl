package tax

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
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
// Please look at the regimes package and other country specific implementations for
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
/*
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
*/

// Set will update the extension map with the provided key and value, and
// return the updated map. If the map is nil, it will be created. If the
// provided code is empty, the key will be removed from the map.
func (em Extensions) Set(key cbc.Key, code cbc.Code) Extensions {
	if em == nil {
		em = make(Extensions)
	}
	if code == cbc.CodeEmpty {
		delete(em, key)
		return em
	}
	em[key] = code
	return em
}

// SetOneOf sets the value to the first code provided, unless it is already
// set as one of the other codes in the list, and return the updated map.
// If the map is nil, it will be created.
func (em Extensions) SetOneOf(key cbc.Key, code cbc.Code, codes ...cbc.Code) Extensions {
	if em == nil {
		em = make(Extensions)
	}
	cur, ok := em[key]
	if !ok || !cur.In(codes...) {
		em[key] = code
	}
	return em
}

// SetIfEmpty sets the value for the key only if it is not already set.
func (em Extensions) SetIfEmpty(key cbc.Key, code cbc.Code) Extensions {
	if em == nil {
		em = make(Extensions)
	}
	if _, ok := em[key]; !ok {
		em[key] = code
	}
	return em
}

// Delete safely removes the key from the extensions map. Returns the
// extension for chaining.
func (em Extensions) Delete(k cbc.Key) Extensions {
	if em == nil {
		return nil
	}
	delete(em, k)
	return em
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

// Values provides an array of all the accepted codes for the extensions
// defined.
func (em Extensions) Values() []cbc.Code {
	if len(em) == 0 {
		return nil
	}
	values := make([]cbc.Code, 0, len(em))
	for _, v := range em {
		values = append(values, v)
	}
	return values
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

type extCodeOp int

const (
	_                      = iota
	extCodeOpAnd extCodeOp = 1 + iota
	extCodeOpNot
	extCodeOpXNOr
	extCodeOpOneOf
	extCodeOpHasCodes
	extCodeOpExcludeCodes
)

// ExtensionsRule is a validation rule for extension maps. It implements both
// validation.Rule (for use with the invopop/validation package) and the
// rules.Test interface (Check + String), so it can be used in rules.When
// conditions and rules.Assert tests as well as regular validation.
type ExtensionsRule struct {
	desc   string
	op     extCodeOp
	keys   []cbc.Key  // used by key-based operators
	key    cbc.Key    // used by code-based operators
	values []cbc.Code // used by code-based operators
}

// ExtensionHasValidCode returns a validation rule that ensures that if the provided key is present
// in the extensions map, that it's code matches the underlying extension's definition. Unlike other
// tests, if the extension key is not present, the test will still pass.
func ExtensionHasValidCode(key cbc.Key) rules.Test {
	ed := ExtensionForKey(key)
	if ed == nil {
		panic("invalid ext key '" + key.String() + "' provided to ExtensionHasValidCode rule: no definition found")
	}
	desc := "ext '" + key.String() + "' "
	var check rules.Test
	if len(ed.Values) > 0 {
		codes := cbc.DefinitionCodes(ed.Values)
		desc = desc + "in [" + strings.Join(cbc.CodeStrings(codes), ", ") + "]"
		check = cbc.InCodes(codes...)
	} else if ed.Pattern != "" {
		desc = desc + "matches pattern '" + ed.Pattern + "'"
		re := regexp.MustCompile(ed.Pattern)
		check = is.MatchesRegexp(re)
	} else {
		panic("invalid ext definition for key '" + key.String() + "': no values or pattern defined")
	}
	return is.Func(
		desc,
		func(value any) bool {
			em, ok := value.(Extensions)
			if !ok {
				return false // only valid for extensions
			}
			ev, ok := em[key]
			if !ok {
				return true // if the key is not present, we don't want to fail validation here
			}
			return check.Check(ev)
		})
}

// ExtensionsRequire returns a validation rule that ensures that all of
// the provided keys are present.
func ExtensionsRequire(keys ...cbc.Key) ExtensionsRule {
	return ExtensionsRule{
		op:   extCodeOpAnd,
		keys: keys,
		desc: "ext require " + extKeyList(keys),
	}
}

// ExtensionsRequireAllOrNone returns a validation rule that performs an XNOR
// operation on the provided keys. If one of the keys is present, then
// all of them must be present. If none of the keys are present,
// then all of them must be absent.
func ExtensionsRequireAllOrNone(keys ...cbc.Key) ExtensionsRule {
	return ExtensionsRule{
		op:   extCodeOpXNOr,
		keys: keys,
		desc: "ext require all or none of " + extKeyList(keys),
	}
}

// ExtensionsExclude returns a validation rule that ensures that
// an extensions map does **not** include the provided keys.
func ExtensionsExclude(keys ...cbc.Key) ExtensionsRule {
	return ExtensionsRule{
		op:   extCodeOpNot,
		keys: keys,
		desc: "ext exclude " + extKeyList(keys),
	}
}

// ExtensionsAllowOneOf returns a validation rule that ensures at most
// one of the provided keys is present in the extensions map. This is useful
// for mutually exclusive options where none or one is allowed.
func ExtensionsAllowOneOf(keys ...cbc.Key) ExtensionsRule {
	return ExtensionsRule{
		op:   extCodeOpOneOf,
		keys: keys,
		desc: "ext allow one of " + extKeyList(keys),
	}
}

// ExtensionsHasCodes returns a validation rule that ensures the extension map's
// key has one of the provided **codes**.
func ExtensionsHasCodes(key cbc.Key, codes ...cbc.Code) ExtensionsRule {
	return ExtensionsRule{
		op:     extCodeOpHasCodes,
		key:    key,
		values: codes,
		desc:   "ext '" + key.String() + "' in " + extCodeList(codes),
	}
}

// ExtensionsExcludeCodes returns a validation rule that ensures the extension map's
// key does not have any of the provided **codes**.
func ExtensionsExcludeCodes(key cbc.Key, codes ...cbc.Code) ExtensionsRule {
	return ExtensionsRule{
		op:     extCodeOpExcludeCodes,
		key:    key,
		values: codes,
		desc:   "ext '" + key.String() + "' not in " + extCodeList(codes),
	}
}

// Validate implements the validation.Rule interface. It returns an error when
// the extensions map does not satisfy the rule.
func (v ExtensionsRule) Validate(value any) error {
	em, ok := value.(Extensions)
	if !ok {
		return nil
	}
	err := make(validation.Errors)

	switch v.op {
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
	case extCodeOpOneOf:
		present := false
		for _, k := range v.keys {
			if _, ok := em[k]; ok {
				if present {
					for _, k := range v.keys {
						if _, ok := em[k]; ok {
							err[k.String()] = errors.New("only one allowed")
						}
					}
					break
				}
				present = true
			}
		}
	case extCodeOpHasCodes:
		if ev, ok := em[v.key]; ok {
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
		}
	case extCodeOpExcludeCodes:
		if ev, ok := em[v.key]; ok {
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

// Check implements the rules.Test interface. It returns true when the
// extensions map satisfies the rule (i.e. validation passes).
func (v ExtensionsRule) Check(val any) bool {
	return v.Validate(val) == nil
}

// String implements the rules.Test interface, returning the human-readable
// description set when the rule was constructed.
func (v ExtensionsRule) String() string {
	return v.desc
}

func extKeyList(keys []cbc.Key) string {
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = k.String()
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func extCodeList(codes []cbc.Code) string {
	parts := make([]string, len(codes))
	for i, c := range codes {
		parts[i] = c.String()
	}
	return "[" + strings.Join(parts, ", ") + "]"
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
