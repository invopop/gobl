package tax

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"iter"
	"maps"
	"regexp"
	"sort"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/jsonschema"
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
//   - An ISO country code e.g. `mx`, `es`, `gb`, etc.
//   - Short abreviation of the platform or format the extension will be used with,
//     e.g. `cfdi` for Mexico's CFDI defined by the SAT, `facturae` for Spain's
//     FacturaE format, `sdi` for the Italian SDI (document interchange system), etc.
//     This is important, as it helps avoid potential conflicts in the future with
//     new or alternative formats that may appear.
//   - A short descriptive name of the extension, e.g. `exception`, `fiscal-regime`,
//     `vat-cat`, `incoming-typ`, etc. The aim should be to avoid using obvious names
//     like `code` or `key` in the name, as these are already implied through usage.
//
// Please look at the regimes package and other country specific implementations for
// examples of how to define and use extensions.
//
// Extensions is immutable: every mutation method (Set, Delete, Merge, etc.)
// returns a new Extensions instance with its own underlying map. Chaining is
// supported for ergonomic construction:
//
//	ext := tax.MakeExtensions().
//	    Set("es-sii-doc-type", "F1").
//	    Set("untdid-document-type", "380")
//
// Extensions wraps a map of extension keys to values. The internal map is
// unexported so that all access goes through the type's methods, which
// guarantee copy-on-write semantics and deterministic (alphabetical) JSON
// output.
type Extensions struct {
	m ExtMap
}

// ExtMap is a short alias for the bare map type used to construct
// Extensions. It lets callers write
//
//	tax.ExtensionsOf(tax.ExtMap{
//	    "key1": "code1",
//	    "key2": "code2",
//	})
//
// instead of spelling out `map[cbc.Key]cbc.Code` at every site.
type ExtMap = map[cbc.Key]cbc.Code

// MakeExtensions returns an empty Extensions ready to be populated via chained
// Set calls.
func MakeExtensions() Extensions {
	return Extensions{}
}

// ExtensionsOf builds a new Extensions from the provided map. The map is
// copied so later mutations of the source do not affect the Extensions.
func ExtensionsOf(m ExtMap) Extensions {
	if len(m) == 0 {
		return Extensions{}
	}
	return Extensions{m: maps.Clone(m)}
}

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

// IsZero returns true when the Extensions has no entries. Used by
// encoding/json with the "omitzero" tag so that empty extension maps are
// omitted from JSON output.
func (e Extensions) IsZero() bool {
	return len(e.m) == 0
}

// Len returns the number of entries in the Extensions.
func (e Extensions) Len() int {
	return len(e.m)
}

// Clone returns an independent copy of the Extensions. Mutations to the
// returned value will not affect the receiver. A zero Extensions clones to
// another zero Extensions (no allocation).
func (e Extensions) Clone() Extensions {
	if len(e.m) == 0 {
		return Extensions{}
	}
	return Extensions{m: maps.Clone(e.m)}
}

// clone returns a new underlying map with the provided extra capacity
// pre-allocated. Always returns a non-nil map.
func (e Extensions) clone(extra int) ExtMap {
	nm := make(ExtMap, len(e.m)+extra)
	for k, v := range e.m {
		nm[k] = v
	}
	return nm
}

// Set returns a new Extensions with the provided key set to the given code.
// If the code is empty, the key is removed from the returned Extensions.
func (e Extensions) Set(key cbc.Key, code cbc.Code) Extensions {
	if code == cbc.CodeEmpty {
		return e.Delete(key)
	}
	nm := e.clone(1)
	nm[key] = code
	return Extensions{m: nm}
}

// SetOneOf returns a new Extensions where the given key is set to the first
// code unless the current value for that key is already one of the additional
// codes provided.
func (e Extensions) SetOneOf(key cbc.Key, code cbc.Code, codes ...cbc.Code) Extensions {
	if cur, ok := e.m[key]; ok && cur.In(codes...) {
		return e
	}
	nm := e.clone(1)
	nm[key] = code
	return Extensions{m: nm}
}

// SetIfEmpty returns a new Extensions where the given key is set to the
// provided code only if the key was not already present.
func (e Extensions) SetIfEmpty(key cbc.Key, code cbc.Code) Extensions {
	if _, ok := e.m[key]; ok {
		return e
	}
	nm := e.clone(1)
	nm[key] = code
	return Extensions{m: nm}
}

// Delete returns a new Extensions without the provided key. If the key is not
// present, the returned Extensions has the same content as the receiver.
func (e Extensions) Delete(key cbc.Key) Extensions {
	if _, ok := e.m[key]; !ok {
		return e
	}
	nm := e.clone(0)
	delete(nm, key)
	return Extensions{m: nm}
}

// Get returns the value for the provided key or an empty code if not found.
// If the key is composed of sub-keys and no exact match is found, the key
// is progressively shortened until a match is found or the key is empty.
func (e Extensions) Get(k cbc.Key) cbc.Code {
	if len(e.m) == 0 {
		return ""
	}
	for k != cbc.KeyEmpty {
		if v, ok := e.m[k]; ok {
			return v
		}
		k = k.Pop()
	}
	return ""
}

// Has returns true if the Extensions contains entries for all the provided
// keys.
func (e Extensions) Has(keys ...cbc.Key) bool {
	for _, k := range keys {
		if _, ok := e.m[k]; !ok {
			return false
		}
	}
	return true
}

// Equals returns true if the other Extensions has exactly the same keys and
// values as this one.
func (e Extensions) Equals(other Extensions) bool {
	if len(e.m) != len(other.m) {
		return false
	}
	if len(e.m) == 0 {
		return true
	}
	return e.Contains(other)
}

// Contains returns true if this Extensions contains all the key/value pairs
// of the other Extensions. It may contain additional keys that are not in
// the other. Always returns false if the receiver is empty.
func (e Extensions) Contains(other Extensions) bool {
	if len(e.m) == 0 {
		return false
	}
	for k, v := range other.m {
		v2, ok := e.m[k]
		if !ok {
			return false
		}
		if v2 != v {
			return false
		}
	}
	return true
}

// Merge returns a new Extensions combining the entries of the receiver with
// those of the other Extensions. Keys in both maps take the value from
// other.
func (e Extensions) Merge(other Extensions) Extensions {
	if len(e.m) == 0 && len(other.m) == 0 {
		return Extensions{}
	}
	nm := e.clone(len(other.m))
	for k, v := range other.m {
		nm[k] = v
	}
	return Extensions{m: nm}
}

// Lookup returns the first key whose value matches the provided code, or an
// empty key if no match is found. Note: if multiple keys share the value,
// the returned key is not deterministic (depends on map iteration order).
func (e Extensions) Lookup(val cbc.Code) cbc.Key {
	for k, v := range e.m {
		if v == val {
			return k
		}
	}
	return cbc.KeyEmpty
}

// Values returns all the values of the Extensions sorted alphabetically.
func (e Extensions) Values() []cbc.Code {
	if len(e.m) == 0 {
		return nil
	}
	values := make([]cbc.Code, 0, len(e.m))
	for _, v := range e.m {
		values = append(values, v)
	}
	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})
	return values
}

// Keys returns all the keys of the Extensions sorted alphabetically.
func (e Extensions) Keys() []cbc.Key {
	if len(e.m) == 0 {
		return nil
	}
	keys := make([]cbc.Key, 0, len(e.m))
	for k := range e.m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

// All returns an iterator over the Extensions entries in alphabetical order
// of the keys. Intended for use with Go 1.23+ range-over-func:
//
//	for k, v := range ext.All() {
//	    // ...
//	}
func (e Extensions) All() iter.Seq2[cbc.Key, cbc.Code] {
	return func(yield func(cbc.Key, cbc.Code) bool) {
		for _, k := range e.Keys() {
			if !yield(k, e.m[k]) {
				return
			}
		}
	}
}

// Clean returns a new Extensions with empty-code entries removed. If the
// result has no entries, a zero Extensions is returned.
func (e Extensions) Clean() Extensions {
	if len(e.m) == 0 {
		return Extensions{}
	}
	nm := make(map[cbc.Key]cbc.Code, len(e.m))
	for k, v := range e.m {
		if v == cbc.CodeEmpty {
			continue
		}
		nm[k] = v
	}
	if len(nm) == 0 {
		return Extensions{}
	}
	return Extensions{m: nm}
}

// MarshalJSON emits the Extensions as a JSON object with keys sorted
// alphabetically for deterministic output. An empty Extensions marshals to
// "null".
func (e Extensions) MarshalJSON() ([]byte, error) {
	if len(e.m) == 0 {
		return []byte("null"), nil
	}
	keys := e.Keys()
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, k := range keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		kb, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		buf.Write(kb)
		buf.WriteByte(':')
		vb, err := json.Marshal(e.m[k])
		if err != nil {
			return nil, err
		}
		buf.Write(vb)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON reads a JSON object into the Extensions. A JSON null is
// treated as an empty Extensions.
func (e *Extensions) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		e.m = nil
		return nil
	}
	var m map[cbc.Key]cbc.Code
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	e.m = m
	return nil
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

// ExtensionsRule is a validation rule for extension maps. It implements the
// rules.Test interface (Check + String), so it can be used in rules.When
// conditions and rules.Assert tests. It also provides a Validate method
// that returns detailed per-key errors.
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
			em, ok := ExtensionsFromValue(value)
			if !ok {
				return false // only valid for extensions
			}
			ev, ok := em.m[key]
			if !ok {
				return true // if the key is not present, we don't want to fail validation here
			}
			return check.Check(ev)
		})
}

// ExtensionsFromValue extracts an Extensions from an any, accepting either
// an Extensions value or a *Extensions pointer (which the rules framework
// produces when validating struct-typed fields). It is intended for use
// inside custom validation-guard functions such as those passed to is.Func:
//
//	func myGuard(val any) bool {
//	    ext, ok := tax.ExtensionsFromValue(val)
//	    return ok && ext.Has("some-key")
//	}
func ExtensionsFromValue(value any) (Extensions, bool) {
	switch e := value.(type) {
	case Extensions:
		return e, true
	case *Extensions:
		if e == nil {
			return Extensions{}, true
		}
		return *e, true
	default:
		return Extensions{}, false
	}
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

// extErrors is a map of extension keys to errors used by Validate.
//
//nolint:errname
type extErrors map[string]error

func (ee extErrors) Error() string {
	keys := make([]string, 0, len(ee))
	for k := range ee {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for i, k := range keys {
		if i > 0 {
			b.WriteString("; ")
		}
		fmt.Fprintf(&b, "%s: %s", k, ee[k].Error())
	}
	b.WriteString(".")
	return b.String()
}

// Validate returns an error when the extensions map does not satisfy the rule.
//
//nolint:gocyclo
func (v ExtensionsRule) Validate(value any) error {
	em, ok := ExtensionsFromValue(value)
	if !ok {
		return nil
	}
	err := make(extErrors)

	switch v.op {
	case extCodeOpAnd:
		for _, k := range v.keys {
			if _, ok := em.m[k]; !ok {
				err[k.String()] = errors.New("required")
			}
		}
	case extCodeOpNot:
		for _, k := range v.keys {
			if _, ok := em.m[k]; ok {
				err[k.String()] = errors.New("must be blank")
			}
		}
	case extCodeOpXNOr:
		present := 0
		for _, k := range v.keys {
			if _, ok := em.m[k]; ok {
				present++
			}
		}
		if present > 0 && present != len(v.keys) {
			for _, k := range v.keys {
				if _, ok := em.m[k]; !ok {
					err[k.String()] = errors.New("required")
				}
			}
		}
	case extCodeOpOneOf:
		present := false
		for _, k := range v.keys {
			if _, ok := em.m[k]; ok {
				if present {
					for _, k := range v.keys {
						if _, ok := em.m[k]; ok {
							err[k.String()] = errors.New("only one allowed")
						}
					}
					break
				}
				present = true
			}
		}
	case extCodeOpHasCodes:
		if ev, ok := em.m[v.key]; ok {
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
		if ev, ok := em.m[v.key]; ok {
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

// JSONSchemaExtend replaces the struct-reflection artefacts (empty
// `properties`, `additionalProperties: false`) with a pattern-properties
// schema that reproduces the shape Extensions had as a map type:
// keys must match the cbc.Key pattern, and values are references to the
// cbc.Code schema. The ref target is resolved via `schema.Lookup` so we
// don't hardcode the URL.
func (Extensions) JSONSchemaExtend(s *jsonschema.Schema) {
	s.Properties = nil
	s.AdditionalProperties = nil
	s.PatternProperties = map[string]*jsonschema.Schema{
		cbc.KeyPattern: {Ref: schema.Lookup(cbc.Code("")).String()},
	}
}
