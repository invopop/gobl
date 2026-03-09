// Package rules provides a framework for defining and applying validation
// rules to data structures in order to provide consistent error codes
// and messages from GOBL.
package rules

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/invopop/gobl/schema"
)

// GOBL for GOBL rules.
const GOBL Code = "GOBL"

var globalRegistry = make([]*Set, 0)

// Code defines a unique code to use for rules.
type Code string

// Set represents a collection of rules grouped by a namespace
// an associated with a specific struct.
type Set struct {
	// ID is the namespace for this set of rules, typically a package-level code like "GOBL" or "GOBL-ORG".
	ID Code `json:"id,omitempty"`
	// Name is the name of the struct type this set of rules applies to. It is used for informational purposes and is not required to be unique.
	Name string `json:"name,omitempty"`
	// Schema identifies the schema that this set of rules applies to. It is optional and can be used to further specify the context of the rules, but it is not required for validation to work.
	Schema schema.ID `json:"schema,omitempty"`
	// Test is an optional expression that determines when this set of rules should be applied. If provided, the set will only be applied when the expression evaluates to true. The expression can reference any exported field from the struct associated with this set of rules.
	Test string `json:"test,omitempty"`
	// Assert is a list of assertions to evaluate directly on the struct associated with this set of rules.
	Assert []*Assertion `json:"assert,omitempty"`
	// Subsets are additional sets of rules to apply recursively to the struct associated with this set of rules. They will be applied in order, and their assertions will be evaluated after the assertions in this set. Subsets can also have their own Test conditions, which will be evaluated independently.
	Subsets []*Set `json:"subsets,omitempty"`

	expr    *vm.Program
	objType reflect.Type
}

// Assertion represents a single validation rule definition.
type Assertion struct {
	// ID defines a globally unique code for this assertion.
	ID Code `json:"id"`
	// Name of the field to test
	Name string `json:"name,omitempty"`
	// Test is the expression to evaluate for this assertion. A true result indicates a failure.
	Test string `json:"test"`
	// Desc is the human-readable message to include in faults when this assertion fails.
	Desc string `json:"desc,omitempty"`

	field any
	expr  *vm.Program
}

// Rulable is an interface that structs can implement to provide their own validation rules.
type Rulable interface {
	Rules() *Set
}

// Register is used to register a set of rules for a given namespace. Namespaces
func Register(name string, pkg Code, objs ...Rulable) {
	set := &Set{
		ID:   pkg,
		Name: name,
	}
	for _, obj := range objs {
		subset := obj.Rules()
		prependToAssertions(pkg, []*Set{subset})
		set.Subsets = append(set.Subsets, subset)
	}
	globalRegistry = append(globalRegistry, set)
}

// When allows us to create a new set of rules that will only be applied when the provided test expression evaluates to true. The test expression can reference any exported field from the struct associated with this set of rules.
func When(obj any, test string, sets ...*Set) *Set {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	env := reflect.New(t).Elem().Interface()
	prog, err := expr.Compile(test, expr.AsBool(), expr.WithTag("json"), expr.Env(env))
	if err != nil {
		panic("invalid rules condition: " + err.Error())
	}
	compileSetAssertions(env, sets...)
	return &Set{
		Test:    test,
		Subsets: sets,
		expr:    prog,
		objType: t,
	}
}

// ForValue creates a new set of rules for the provided value. Each assertion will be
// evaluated with a `this` variable in the expression environment that holds the
// underlying primitive value of obj.
func ForValue(obj any, asserts ...*Assertion) *Set {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	setID := typeSetID(t)
	env := buildEnv(t, reflect.New(t).Elem())
	compileAssertions(env, asserts...)
	for _, a := range asserts {
		// Prepend the type name to the assertion ID, mirroring ForStruct behaviour.
		if a.ID != "" {
			a.ID = setID.Add(a.ID)
		}
	}
	return &Set{
		ID:      setID,
		Name:    t.Name(),
		Schema:  schema.Lookup(obj),
		objType: t,
		Assert:  asserts,
	}
}

// ForStruct creates a new set of rules for the provided struct, attaching the provided
// subsets (from Field or When) and resolving assertion field names. obj must be
// a pointer. Field pointer assertions are resolved by byte offset.
func ForStruct(obj any, subsets ...*Set) *Set {
	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	setID := typeSetID(t)
	prependToAssertions(setID, subsets)
	out := &Set{
		ID:      setID,
		Name:    t.Name(),
		Schema:  schema.Lookup(obj),
		objType: t,
		Subsets: subsets,
	}
	if t.Kind() == reflect.Struct {
		// Resolve the JSON field name for each assertion that carries a field identifier.
		for _, ss := range subsets {
			for _, a := range ss.Assert {
				if a.field != nil {
					a.Name = resolveAttributeName(t, obj, a.field)
				}
			}
		}
	}
	env := buildEnv(t, reflect.New(t).Elem())
	compileSetAssertions(env, subsets...)
	return out
}

// compileSubAssertions compiles any uncompiled assertions in the direct Assert
// slices of each provided set, using env for type-aware compilation. It does not
// recurse into sub-subsets: When and For each compile their own immediate assertions.
func compileSetAssertions(env any, set ...*Set) {
	for _, ss := range set {
		compileAssertions(env, ss.Assert...)
	}
}

func compileAssertions(env any, asserts ...*Assertion) {
	for _, a := range asserts {
		if a.expr != nil {
			continue
		}
		// Backslashes in tests must be doubled for embedding inside an expr
		// double-quoted string literal (expr interprets \\ as a single \).
		et := strings.ReplaceAll(a.Test, `\`, `\\`)
		opts := append(isHelpers, expr.AsBool(), expr.WithTag("json"), expr.Env(env))
		prog, err := expr.Compile(et, opts...)
		if err != nil {
			panic(fmt.Sprintf("invalid assertion %s (%q): %s", a.ID, a.Test, err.Error()))
		}
		a.expr = prog
	}
}

// typeSetID derives a set ID from the type name, converted to upper case.
// For example, Email becomes EMAIL. The package namespace is contributed
// separately by Register, so only the type name is used here.
func typeSetID(t reflect.Type) Code {
	return Code(strings.ToUpper(t.Name()))
}

// prependToAssertions recursively prepends code to all assertion IDs within the
// provided sets and their subsets.
func prependToAssertions(code Code, subsets []*Set) {
	for _, ss := range subsets {
		for _, a := range ss.Assert {
			if a.ID != "" {
				a.ID = code.Add(a.ID)
			}
		}
		prependToAssertions(code, ss.Subsets)
	}
}

// Field helps define a set of assertions associated with a specific field
// embedded inside a struct. field may be either a Go or JSON field name string,
// or a pointer to the field (e.g. &myStruct.FieldName) for compile-time typo
// prevention. The field name is resolved to the JSON tag name when the set is
// embedded in a For call.
func Field(field any, assert ...*Assertion) *Set {
	rules := make([]*Assertion, len(assert))
	for i, a := range assert {
		a.field = field
		rules[i] = a
	}
	return &Set{
		Assert: rules,
	}
}

// Assert creates a new assertion with the provided code, test expression,
// and description. Positive test results will cause this assertion to be
// considered a failure, and the provided code and description will be included
// in validation faults. The expression is compiled lazily when the assertion
// is embedded in a For or When call, which provides the parent struct type
// for field validation.
func Assert(id Code, test string, desc string) *Assertion {
	return &Assertion{
		ID:   id,
		Test: test,
		Desc: desc,
	}
}

// Add allows us to create a new code by appending a suffix to the existing code.
func (c Code) Add(code Code) Code {
	return c + "-" + code
}

// AllSets returns all rule sets registered in the global registry.
func AllSets() []*Set {
	return globalRegistry
}

// Validate uses the global registry of rule sets to test the provided object against
// all available assertions, then recursively validates all exported struct fields.
// Returns nil when no faults are found.
func Validate(obj any) Faults {
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	objType := rv.Type()
	var faults []*Fault

	// Find and apply all matching sets from the global registry.
	for _, ns := range globalRegistry {
		for _, subset := range ns.Subsets {
			if subset.objType == objType {
				if fs := subset.Validate(obj); fs != nil {
					faults = append(faults, fs.List()...)
				}
			}
		}
	}

	if rv.Kind() != reflect.Struct {
		return newFaults(faults)
	}

	// Recurse into exported fields.
	rt := rv.Type()
	for i := range rv.NumField() {
		sf := rt.Field(i)
		if !sf.IsExported() {
			continue
		}
		fv := rv.Field(i)
		fs := validateFieldValue(fv)
		if len(fs) == 0 {
			continue
		}
		if sf.Anonymous {
			// Promote embedded struct faults to parent level.
			faults = append(faults, fs...)
			continue
		}
		name := jsonFieldName(sf)
		if name != "" {
			faults = append(faults, prependPath(name, fs)...)
		}
	}

	return newFaults(faults)
}

// Validate validates an object against the set's rules. If the set has a test
// condition (from When), it is evaluated first and the set is skipped when false.
// Expressions reference exported Go field names on the struct directly.
// Returns nil when no faults are found.
func (s *Set) Validate(obj any) Faults {
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	// Skip if this set is bound to a different type.
	if s.objType != nil && rv.Type() != s.objType {
		return nil
	}

	env := buildEnv(rv.Type(), rv)

	// Evaluate the When condition; skip the set if it doesn't match.
	if s.expr != nil {
		result, err := expr.Run(s.expr, env)
		if err != nil {
			panic("expression runtime error: " + err.Error())
		}
		if !result.(bool) {
			return nil
		}
	}

	var faults []*Fault

	// Run assertions: a positive (true) result indicates a failure.
	for _, a := range s.Assert {
		if a.expr == nil {
			panic(fmt.Sprintf("assertion %s (%q) was not compiled; wrap it in For() or When()", a.ID, a.Test))
		}
		result, err := expr.Run(a.expr, env)
		if err != nil {
			panic(fmt.Sprintf("assertion %s (%q) runtime error: %s", a.ID, a.Test, err.Error()))
		}
		if result.(bool) {
			faults = append(faults, newFault(a.Name, a.ID, a.Desc))
		}
	}

	// Recurse into subsets.
	for _, ss := range s.Subsets {
		if fs := ss.Validate(obj); fs != nil {
			faults = append(faults, fs.List()...)
		}
	}

	return newFaults(faults)
}

// validateFieldValue recursively validates a field value, handling pointers,
// structs, slices, and arrays.
func validateFieldValue(fv reflect.Value) []*Fault {
	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			return nil
		}
		fv = fv.Elem()
	}
	switch fv.Kind() {
	case reflect.Struct:
		if fs := Validate(fv.Interface()); fs != nil {
			return fs.List()
		}
		return nil
	case reflect.Slice, reflect.Array:
		var faults []*Fault
		for i := range fv.Len() {
			ev := fv.Index(i)
			if fs := validateFieldValue(ev); len(fs) > 0 {
				faults = append(faults, prependPath("["+strconv.Itoa(i)+"]", fs)...)
			}
		}
		return faults
	default:
		// For named non-struct types (e.g. cbc.Code), check the global registry.
		if fv.Type().PkgPath() != "" {
			if fs := Validate(fv.Interface()); fs != nil {
				return fs.List()
			}
		}
	}
	return nil
}

// buildEnv constructs the expression environment for the given type and value.
// Struct values are returned directly so that expressions can reference fields
// by their json tag names. Non-struct types are exposed as a map with a single
// "this" key holding the underlying primitive value, so that assertion
// expressions can reference the value uniformly regardless of type.
func buildEnv(t reflect.Type, rv reflect.Value) any {
	if t.Kind() == reflect.Struct {
		return rv.Interface()
	}
	_, fieldVal := underlyingTypeAndValue(t, rv)
	return map[string]any{"this": fieldVal.Interface()}
}

// underlyingTypeAndValue converts a named type to its underlying Go primitive
// type so that expr operators work correctly. For example, cbc.Code (a named
// string type) becomes plain string. Unnamed types are returned unchanged.
func underlyingTypeAndValue(t reflect.Type, rv reflect.Value) (reflect.Type, reflect.Value) {
	if t.PkgPath() == "" {
		return t, rv // already an unnamed/built-in type
	}
	switch t.Kind() {
	case reflect.String:
		ut := reflect.TypeOf("")
		return ut, rv.Convert(ut)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ut := reflect.TypeOf(int64(0))
		return ut, rv.Convert(ut)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ut := reflect.TypeOf(uint64(0))
		return ut, rv.Convert(ut)
	case reflect.Float32, reflect.Float64:
		ut := reflect.TypeOf(float64(0))
		return ut, rv.Convert(ut)
	case reflect.Bool:
		ut := reflect.TypeOf(false)
		return ut, rv.Convert(ut)
	}
	return t, rv
}

// resolveAttributeName returns the JSON field name for field within the struct
// type t. field may be a Go or JSON field name string, or a pointer to a field
// of obj (which must be a pointer to a struct of type t). In the pointer case
// the field is identified by its byte offset from the struct base.
func resolveAttributeName(t reflect.Type, obj any, field any) string {
	if name, ok := field.(string); ok {
		for i := range t.NumField() {
			f := t.Field(i)
			if f.Name == name || jsonFieldName(f) == name {
				return jsonFieldName(f)
			}
		}
		return name
	}
	baseAddr := reflect.ValueOf(obj).Pointer()
	valAddr := reflect.ValueOf(field).Pointer()
	offset := valAddr - baseAddr
	for i := range t.NumField() {
		if t.Field(i).Offset == offset {
			return jsonFieldName(t.Field(i))
		}
	}
	return ""
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
