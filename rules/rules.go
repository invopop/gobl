// Package rules provides a framework for defining and applying validation
// rules to data structures in order to provide consistent error codes
// and messages from GOBL.
package rules

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// GOBL for GOBL rules.
const GOBL Code = "GOBL"

var globalRegistry = make([]*Set, 0)

// Code defines a unique code to use for rules.
type Code string

// Def is a function that modifies a Set during construction.
// Assert, Field, Each, Object, and When all return Def values that compose
// as arguments to For.
type Def func(s *Set)

// Set represents a collection of rules grouped by a namespace
// an associated with a specific struct.
type Set struct {
	// ID is the namespace for this set of rules, typically a package-level code like "GOBL" or "GOBL-ORG".
	ID Code `json:"id,omitempty"`
	// Name is the name of the struct type this set of rules applies to. It is used for informational purposes and is not required to be unique.
	Name string `json:"name,omitempty"`
	// FieldName is the JSON tag name of the field this subset is scoped to. When non-empty, Validate extracts this field from the parent object and delegates to it.
	FieldName string `json:"field,omitempty"`
	// Each when true causes Validate to iterate over the slice elements of the field named by FieldName.
	Each bool `json:"each,omitempty"`
	// Guard is an optional expression that determines when this set of rules should be applied. If provided, the set will only be applied when the expression evaluates to true. The expression can reference any exported field from the struct associated with this set of rules.
	Guard Test `json:"guard,omitempty"`
	// Assert is a list of assertions to evaluate directly on the struct associated with this set of rules.
	Assert []*Assertion `json:"assert,omitempty"`
	// Subsets are additional sets of rules to apply recursively to the struct associated with this set of rules. They will be applied in order, and their assertions will be evaluated after the assertions in this set. Subsets can also have their own Test conditions, which will be evaluated independently.
	Subsets []*Set `json:"subsets,omitempty"`

	objType reflect.Type
}

// Assertion represents a single validation rule definition.
type Assertion struct {
	// ID defines a globally unique code for this assertion.
	ID Code `json:"id"`
	// Desc is the human-readable message to include in faults when this assertion fails.
	Desc string `json:"desc,omitempty"`
	// Tests is a list of tests to evaluate for this assertion. A false result indicates a failure.
	Tests []Test `json:"tests"`
}

// Test defines an interface expected for a test condition.
type Test interface {
	Check(val any) bool
	String() string
}

// Embeddable is implemented by types that wrap a private payload whose rules
// should be validated as if the payload were at the same JSON level as the wrapper.
// When the rules traversal encounters a struct that implements Embeddable, it calls
// Embedded() and recurses into the result without adding any path prefix.
type Embeddable interface {
	Embedded() any
}

type compilableTest interface {
	compile(val any) error
}

// Registry returns the global registry of rule sets.
func Registry() []*Set {
	return globalRegistry
}

// Register is used to register a set of rules for a given namespace.
func Register(name string, pkg Code, sets ...*Set) {
	RegisterWithGuard(name, pkg, nil, sets...)
}

// RegisterWithGuard is used to register a set of rules for a given namespace
// with an optional guard condition that determines when the rules should be applied.
func RegisterWithGuard(name string, pkg Code, guard Test, sets ...*Set) {
	set := &Set{
		ID:      pkg,
		Name:    name,
		Guard:   guard,
		Subsets: sets,
	}
	prependToSets(pkg, sets)
	globalRegistry = append(globalRegistry, set)
}

// For creates a new set of rules for the provided object (struct or value type).
// Each Def is applied in order to build up the set's assertions and subsets.
// Assert, Field, Each, Object, and When all return Def values that can be passed here.
func For(obj any, defs ...Def) *Set {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	setID := typeSetID(t)
	name := t.Name()
	if pkg := pkgShortName(t); pkg != "" {
		name = pkg + "." + name
	}
	s := &Set{
		ID:      setID,
		Name:    name,
		objType: t,
	}
	for _, def := range defs {
		def(s)
	}
	prependToAssertions(setID, s.Assert)
	prependToSets(setID, s.Subsets)
	compileAndResolve(t, obj, s)
	return s
}

// Assert returns a Def that adds a single assertion to the parent set.
// The assertion is evaluated against the parent object (or extracted field
// value when used inside Field or Each).
func Assert(id Code, desc string, tests ...Test) Def {
	a := &Assertion{
		ID:    id,
		Desc:  desc,
		Tests: tests,
	}
	return func(s *Set) {
		s.Assert = append(s.Assert, a)
	}
}

// AssertIfPresent returns a Def that adds an assertion that is only evaluated
// when the current scoped value is non-nil and non-empty. Use this for
// optional fields that have format or content constraints.
func AssertIfPresent(id Code, desc string, tests ...Test) Def {
	return func(s *Set) {
		subset := &Set{Guard: Present}
		subset.Assert = append(subset.Assert, &Assertion{
			ID:    id,
			Desc:  desc,
			Tests: tests,
		})
		s.Subsets = append(s.Subsets, subset)
	}
}

// Object returns a Def that groups assertions evaluated against the whole
// object. It is equivalent to passing the assertions directly to For or When,
// and exists for organisational clarity.
func Object(defs ...Def) Def {
	return func(s *Set) {
		for _, def := range defs {
			def(s)
		}
	}
}

// Field returns a Def that creates a field-scoped subset. name must be the
// JSON tag name of a field in the parent struct. All assertions and subsets
// inside Field receive the extracted field value when validating.
func Field(name string, defs ...Def) Def {
	return func(s *Set) {
		subset := &Set{FieldName: name}
		for _, def := range defs {
			def(subset)
		}
		s.Subsets = append(s.Subsets, subset)
	}
}

// Each returns a Def that iterates over the elements of the current context,
// which must be a slice or array. It is intended to be used inside a Field
// that targets a slice field. All assertions and subsets inside Each are
// applied to each element individually.
//
// Usage:
//
//	rules.Field("lines",
//	    rules.Assert("01", "no duplicates", checkNoDups),  // whole-slice assertion
//	    rules.Each(
//	        rules.Assert("02", "line required", rules.Present),  // per-element
//	    ),
//	)
//
// Each panics at initialisation time if the parent context is not a slice or array.
func Each(defs ...Def) Def {
	return func(s *Set) {
		subset := &Set{Each: true}
		for _, def := range defs {
			def(subset)
		}
		s.Subsets = append(s.Subsets, subset)
	}
}

// When returns a Def that conditionally applies its sub-definitions only when
// test evaluates to true. The test expression is compiled by the parent For call.
func When(guard Test, defs ...Def) Def {
	return func(s *Set) {
		subset := &Set{Guard: guard}
		for _, def := range defs {
			def(subset)
		}
		s.Subsets = append(s.Subsets, subset)
	}
}

// compileAndResolve recursively compiles assertions and test conditions
// throughout the set tree using obj as the prototype environment.
func compileAndResolve(t reflect.Type, obj any, s *Set) {
	compileAssertions(obj, s.Assert...)
	if s.Guard != nil {
		if ct, ok := s.Guard.(compilableTest); ok {
			if err := ct.compile(obj); err != nil {
				panic("invalid rules condition: " + err.Error())
			}
		}
	}
	for _, ss := range s.Subsets {
		if ss.FieldName != "" {
			compileFieldSubset(t, ss)
		} else if ss.Each {
			compileEachSubset(t, ss)
		} else {
			compileAndResolve(t, obj, ss)
		}
	}
}

// compileFieldSubset infers the nested type for a field-scoped subset by looking
// up the field in the parent struct's reflect type, then recursively compiles it.
func compileFieldSubset(t reflect.Type, ss *Set) {
	ft := fieldTypeByName(t, ss.FieldName)
	if ft == nil {
		panic(fmt.Sprintf("rules: field %q not found in type %s", ss.FieldName, t.Name()))
	}
	if ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
	}
	ss.objType = ft
	nestedProto := reflect.New(ft).Interface()
	compileAndResolve(ft, nestedProto, ss)
}

// compileEachSubset infers the element type from the parent slice/array type and
// recursively compiles the subset for that element type. It panics if the parent
// type is not a slice or array.
func compileEachSubset(t reflect.Type, ss *Set) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		panic(fmt.Sprintf("rules: Each used on non-slice type %s", t.Name()))
	}
	ft := t.Elem()
	if ft.Kind() == reflect.Ptr {
		ft = ft.Elem()
	}
	ss.objType = ft
	nestedProto := reflect.New(ft).Interface()
	compileAndResolve(ft, nestedProto, ss)
}

// fieldTypeByName returns the reflect.Type for the field with the given JSON
// tag name in struct type t. Returns nil if not found.
func fieldTypeByName(t reflect.Type, name string) reflect.Type {
	if t.Kind() != reflect.Struct {
		return nil
	}
	for i := range t.NumField() {
		f := t.Field(i)
		if jsonFieldName(f) == name {
			return f.Type
		}
	}
	return nil
}

// fieldValueByName returns the reflect.Value for the field with the given JSON
// tag name in struct value rv. Returns (zero, false) if not found.
func fieldValueByName(rv reflect.Value, name string) (reflect.Value, bool) {
	if rv.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}
	rt := rv.Type()
	for i := range rt.NumField() {
		if jsonFieldName(rt.Field(i)) == name {
			return rv.Field(i), true
		}
	}
	return reflect.Value{}, false
}

// validateEachValue validates each element of a slice/array value against the
// given subset. Fault paths are reported as [0], [1], etc. (no field-name
// prefix; the caller's Field already contributes that).
func validateEachValue(rc *RunCtx, fv reflect.Value, ss *Set) []*Fault {
	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			return nil
		}
		fv = fv.Elem()
	}
	if fv.Kind() != reflect.Slice && fv.Kind() != reflect.Array {
		return nil
	}
	var faults []*Fault
	for i := range fv.Len() {
		ev := fv.Index(i)
		if fs := ss.validate(rc, ev.Interface()); fs != nil {
			faults = append(faults, prependPath("["+strconv.Itoa(i)+"]", fs.List())...)
		}
	}
	return faults
}

func compileAssertions(env any, asserts ...*Assertion) {
	for _, a := range asserts {
		for _, t := range a.Tests {
			if ct, ok := t.(compilableTest); ok {
				if err := ct.compile(env); err != nil {
					panic(fmt.Sprintf("failed to compile assertion %s: %s", a.ID, err.Error()))
				}
			}
		}
	}
}

// typeSetID derives a set ID from the type, including the Go package short name
// when present. For example, tax.Identity becomes TAX-IDENTITY and Email (no
// package) becomes EMAIL. The GOBL prefix and registry namespace are contributed
// by Register, which prepends its code avoiding duplication.
func typeSetID(t reflect.Type) Code {
	pkg := pkgShortName(t)
	if pkg == "" {
		return Code(strings.ToUpper(t.Name()))
	}
	return Code(strings.ToUpper(pkg)).Add(Code(strings.ToUpper(t.Name())))
}

// pkgShortName returns the short package name for a type, stripping any "_test"
// suffix added by Go's external test packages. Returns "" for built-in types.
func pkgShortName(t reflect.Type) string {
	path := t.PkgPath()
	if path == "" {
		return ""
	}
	name := path
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		name = path[idx+1:]
	}
	return strings.TrimSuffix(name, "_test")
}

func prependToSets(code Code, sets []*Set) {
	for _, s := range sets {
		if s.ID != "" {
			s.ID = code.Prepend(s.ID)
		}
		prependToAssertions(code, s.Assert)
		prependToSets(code, s.Subsets)
	}
}

// prependToAssertions recursively prepends code to all assertion IDs within the
// provided sets and their subsets.
func prependToAssertions(code Code, asserts []*Assertion) {
	for _, a := range asserts {
		if a.ID != "" {
			a.ID = code.Prepend(a.ID)
		}
	}
}

// Add allows us to create a new code by appending a suffix to the existing code.
func (c Code) Add(code Code) Code {
	return c + "-" + code
}

// Prepend prepends c to id, deduplicating the last segment of c when it
// matches the first segment of id. This avoids double-encoding the package
// name when the registry namespace already contains it.
// For example: Code("GOBL-ORG").Prepend("ORG-EMAIL") → "GOBL-ORG-EMAIL"
// but: Code("GOBL-GB").Prepend("TAX-IDENTITY") → "GOBL-GB-TAX-IDENTITY"
func (c Code) Prepend(id Code) Code {
	codeStr := string(c)
	idStr := string(id)
	// Extract last segment of c.
	suffix := codeStr
	if i := strings.LastIndex(codeStr, "-"); i >= 0 {
		suffix = codeStr[i+1:]
	}
	// Drop the leading segment of id when it duplicates the suffix.
	if strings.HasPrefix(idStr, suffix+"-") {
		return c + "-" + Code(idStr[len(suffix)+1:])
	}
	return c.Add(id)
}

// AllSets returns all rule sets registered in the global registry.
func AllSets() []*Set {
	return globalRegistry
}

// Validate uses the global registry of rule sets to validate the provided object.
// Each registered namespace set is applied in order; the Set.Validate method is
// responsible for matching the object type, evaluating guard conditions, running
// assertions, and recursively iterating exported struct fields.
// Returns nil when no faults are found.
//
// Optional WithContext values inject additional context into the validation
// session. Context is also collected automatically from the root object's
// exported fields that implement ContextAdder (e.g. tax.Regime, tax.Addons).
func Validate(obj any, opts ...WithContext) Faults {
	rc := &RunCtx{}
	for _, opt := range opts {
		opt(rc)
	}
	collectContext(rc, obj)
	var faults []*Fault
	for _, ns := range globalRegistry {
		if fs := ns.validate(rc, obj); fs != nil {
			faults = append(faults, fs.List()...)
		}
	}
	return newFaults(faults)
}

// Validate validates an object against the set's rules. If the set has a test
// condition (from When), it is evaluated first and the set is skipped when false.
// Returns nil when no faults are found.
func (s *Set) Validate(obj any) Faults {
	return s.validate(nil, obj)
}

// validate is the internal context-aware implementation of Validate.
func (s *Set) validate(rc *RunCtx, obj any) Faults {
	rv := reflect.ValueOf(obj)
	if !rv.IsValid() {
		return nil
	}
	isNil := false
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			isNil = true
		} else {
			rv = rv.Elem()
		}
	}

	// Skip if this set is bound to a different type.
	if s.objType != nil {
		var checkType reflect.Type
		if isNil {
			checkType = rv.Type().Elem()
		} else {
			checkType = rv.Type()
		}
		if checkType != s.objType {
			return nil
		}
	}

	// Normalize obj to a pointer for consistent test calling. When the caller
	// passes a plain struct value (e.g. via fv.Interface() on a non-pointer
	// field), By-style tests that assert value.(*T) would otherwise fail.
	callObj := obj
	if !isNil && rv.Kind() == reflect.Struct {
		if reflect.TypeOf(obj).Kind() != reflect.Ptr {
			ptr := reflect.New(rv.Type())
			ptr.Elem().Set(rv)
			callObj = ptr.Interface()
		}
	}

	// Evaluate the When condition; skip the set if it doesn't match.
	if s.Guard != nil && !runTest(rc, s.Guard, callObj) {
		return nil
	}

	var faults []*Fault

	// Run assertions. For nil pointer objects, assertions still run so that
	// Required can detect the missing value.
	for _, a := range s.Assert {
		if len(a.Tests) == 0 {
			panic(fmt.Sprintf("assertion %s (%q) tests missing", a.ID, a.Tests))
		}
		for _, t := range a.Tests {
			if !runTest(rc, t, callObj) {
				faults = append(faults, newFault("", a.ID, a.Desc))
				break
			}
		}
	}

	// Sub-subsets: skip field extraction when the object is nil.
	if !isNil {
		for _, ss := range s.Subsets {
			if ss.FieldName == "" {
				if ss.Each {
					faults = append(faults, validateEachValue(rc, rv, ss)...)
				} else {
					if fs := ss.validate(rc, obj); fs != nil {
						faults = append(faults, fs.List()...)
					}
				}
			} else {
				fv, ok := fieldValueByName(rv, ss.FieldName)
				if !ok {
					continue
				}
				if fs := ss.validate(rc, fv.Interface()); fs != nil {
					faults = append(faults, prependPath(ss.FieldName, fs.List())...)
				}
			}
		}

		// For namespace-level sets that own type-specific rules, iterate all
		// exported struct fields and apply this set's rules to each nested value.
		// This keeps field iteration scoped to the namespace so that guard tests
		// are not re-evaluated against nested objects, and type-specific rules
		// registered under this namespace (e.g. taxComboRules) are still applied
		// to nested types discovered during traversal.
		if s.isNamespace() && rv.Kind() == reflect.Struct {
			hasTypeRules := false
			for _, ss := range s.Subsets {
				if ss.objType != nil {
					hasTypeRules = true
					break
				}
			}
			if hasTypeRules {
				rt := rv.Type()
				for i := range rv.NumField() {
					sf := rt.Field(i)
					if !sf.IsExported() {
						continue
					}
					fv := rv.Field(i)
					fs := s.validateNestedFieldValue(rc, fv)
					if len(fs) == 0 {
						continue
					}
					if sf.Anonymous {
						faults = append(faults, fs...)
						continue
					}
					name := jsonFieldName(sf)
					if name != "" {
						faults = append(faults, prependPath(name, fs)...)
					}
				}
			}
		}
	}

	return newFaults(faults)
}

// validateNestedValue applies this set's type-specific subsets to obj and
// recursively processes its exported struct fields. It is used during
// namespace-level field iteration and does not re-check the namespace guard.
// isNamespace reports whether s is a registered namespace set — one that has an
// ID, is not bound to a specific type, field, or iteration, and therefore
// owns the struct-field traversal used to apply type-specific subsets (e.g.
// taxComboRules) to nested values discovered at runtime.
func (s *Set) isNamespace() bool {
	return s.ID != "" && s.objType == nil && s.FieldName == "" && !s.Each
}

func (s *Set) validateNestedValue(rc *RunCtx, obj any) []*Fault {
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return nil
	}

	objType := rv.Type()
	callObj := obj
	if rv.Kind() == reflect.Struct && reflect.TypeOf(obj).Kind() != reflect.Ptr {
		ptr := reflect.New(rv.Type())
		ptr.Elem().Set(rv)
		callObj = ptr.Interface()
	}

	var faults []*Fault

	// Apply matching type-specific subsets from this namespace.
	for _, ss := range s.Subsets {
		if ss.objType != objType {
			continue
		}
		if fs := ss.validate(rc, callObj); fs != nil {
			faults = append(faults, fs.List()...)
		}
	}

	// Recurse into struct fields.
	if rv.Kind() == reflect.Struct {
		rt := rv.Type()
		for i := range rv.NumField() {
			sf := rt.Field(i)
			if !sf.IsExported() {
				continue
			}
			fv := rv.Field(i)
			fs := s.validateNestedFieldValue(rc, fv)
			if len(fs) == 0 {
				continue
			}
			if sf.Anonymous {
				faults = append(faults, fs...)
				continue
			}
			name := jsonFieldName(sf)
			if name != "" {
				faults = append(faults, prependPath(name, fs)...)
			}
		}
	}

	// If the object exposes an embedded payload, validate it at the same path level.
	if emb, ok := callObj.(Embeddable); ok {
		if inner := emb.Embedded(); inner != nil {
			faults = append(faults, s.validateNestedValue(rc, inner)...)
		}
	}

	return faults
}

// validateNestedFieldValue handles pointers, structs, slices, and named types
// during namespace-internal field iteration.
func (s *Set) validateNestedFieldValue(rc *RunCtx, fv reflect.Value) []*Fault {
	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			return nil
		}
		fv = fv.Elem()
	}
	switch fv.Kind() {
	case reflect.Struct:
		return s.validateNestedValue(rc, fv.Interface())
	case reflect.Slice, reflect.Array:
		var faults []*Fault
		for i := range fv.Len() {
			ev := fv.Index(i)
			if fs := s.validateNestedFieldValue(rc, ev); len(fs) > 0 {
				faults = append(faults, prependPath("["+strconv.Itoa(i)+"]", fs)...)
			}
		}
		return faults
	case reflect.Map:
		if fv.IsNil() {
			return nil
		}
		// Collect and sort keys for deterministic output.
		keys := fv.MapKeys()
		sorted := make([]string, len(keys))
		for i, k := range keys {
			sorted[i] = fmt.Sprintf("%v", k.Interface())
		}
		sort.Strings(sorted)
		keyByStr := make(map[string]reflect.Value, len(keys))
		for _, k := range keys {
			keyByStr[fmt.Sprintf("%v", k.Interface())] = k
		}
		var faults []*Fault
		for _, ks := range sorted {
			k := keyByStr[ks]
			// Validate named key types (e.g. cbc.Key).
			if k.Type().PkgPath() != "" {
				if fs := s.validateNestedValue(rc, k.Interface()); len(fs) > 0 {
					faults = append(faults, prependPath(ks, fs)...)
				}
			}
			// Validate map values recursively.
			ev := fv.MapIndex(k)
			if fs := s.validateNestedFieldValue(rc, ev); len(fs) > 0 {
				faults = append(faults, prependPath(ks, fs)...)
			}
		}
		return faults
	default:
		// For named non-struct types (e.g. cbc.Code), check this namespace's rules.
		if fv.Type().PkgPath() != "" {
			return s.validateNestedValue(rc, fv.Interface())
		}
	}
	return nil
}

func jsonFieldName(f reflect.StructField) string {
	tag := f.Tag.Get("json")
	if tag == "" {
		return f.Name
	}
	if tag == "-" {
		return ""
	}
	name, _, _ := strings.Cut(tag, ",")
	if name == "" {
		return f.Name
	}
	return name
}

// MarshalJSON serializes Set to JSON, converting the Test field to its string representation.
func (s Set) MarshalJSON() ([]byte, error) {
	type alias struct {
		ID        Code         `json:"id,omitempty"`
		Name      string       `json:"name,omitempty"`
		FieldName string       `json:"field,omitempty"`
		Each      bool         `json:"each,omitempty"`
		Guard     string       `json:"guard,omitempty"`
		Assert    []*Assertion `json:"assert,omitempty"`
		Subsets   []*Set       `json:"subsets,omitempty"`
	}
	a := alias{
		ID:        s.ID,
		Name:      s.Name,
		FieldName: s.FieldName,
		Each:      s.Each,
		Assert:    s.Assert,
		Subsets:   s.Subsets,
	}
	if s.Guard != nil {
		a.Guard = s.Guard.String()
	}
	return json.Marshal(a)
}

// MarshalJSON serializes Assertion to JSON, converting Tests to a comma-joined string.
func (a Assertion) MarshalJSON() ([]byte, error) {
	type alias struct {
		ID    Code   `json:"id"`
		Desc  string `json:"desc,omitempty"`
		Tests string `json:"tests,omitempty"`
	}
	parts := make([]string, len(a.Tests))
	for i, t := range a.Tests {
		parts[i] = t.String()
	}
	return json.Marshal(alias{
		ID:    a.ID,
		Desc:  a.Desc,
		Tests: strings.Join(parts, ", "),
	})
}
