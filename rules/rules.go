// Package rules provides a framework for defining and applying validation
// rules to data structures in order to provide consistent error codes
// and messages from GOBL.
package rules

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// GOBL for GOBL rules.
const GOBL Code = "GOBL"

var (
	// coreRegistry holds namespace sets with no guard (always applied).
	coreRegistry = make([]*Set, 0)
	// guardedIndex maps context keys to namespace sets whose guards
	// depend on those keys, enabling O(1) lookup at validation time.
	guardedIndex = make(map[ContextKey][]*Set)
	// unkeyedGuarded holds guarded sets whose guards do not implement
	// ContextKeyable, so they must be evaluated every time.
	unkeyedGuarded = make([]*Set, 0)
)

// Code defines a unique code to use for rules.
type Code string

// Add allows us to create a new code by appending a suffix to the existing code.
func (c Code) Add(code Code) Code {
	return c + "-" + code
}

// Def is a function that modifies a Set during construction.
// Assert, Field, Each, Object, and When all return Def values that compose
// as arguments to For.
type Def func(s *Set)

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
	Compile(val any) error
}

// Registry returns the global registry of rule sets.
func Registry() []*Set {
	return allSets()
}

// Register is used to register a set of rules for a given namespace.
//
// For addon rules, pkg and code should identify the addon *family* rather
// than a specific version (e.g. "ar-arca" and "AR-ARCA", not "ar-arca-v4"
// and "AR-ARCA-V4"). The addon version is expressed through the guard.
// Assertion IDs are a public contract that customers pin for error handling,
// so they must be stable across versions of an addon family: when introducing
// a new version, preserved rules must keep their existing numeric IDs. An ID
// may be retired, but it must never be reassigned to a semantically different
// rule.
func Register(pkg string, code Code, sets ...*Set) {
	RegisterWithGuard(pkg, code, nil, sets...)
}

// RegisterWithGuard is used to register a set of rules for a given namespace
// with an optional guard condition that determines when the rules should be applied.
// See [Register] for guidance on choosing pkg and code for addon families.
func RegisterWithGuard(pkg string, code Code, guard Test, sets ...*Set) {
	sets = cloneSets(sets)
	set := &Set{
		ID:      code,
		Package: pkg,
		Guard:   guard,
		Subsets: sets,
	}
	prependToSets(code, sets)
	buildTypeIndex(set)
	if guard == nil {
		coreRegistry = append(coreRegistry, set)
	} else if ck, ok := guard.(ContextKeyable); ok {
		if keys := ck.ContextKeys(); len(keys) > 0 {
			for _, k := range keys {
				guardedIndex[k] = append(guardedIndex[k], set)
			}
		} else {
			unkeyedGuarded = append(unkeyedGuarded, set)
		}
	} else {
		unkeyedGuarded = append(unkeyedGuarded, set)
	}
}

// NewSet creates a standalone namespace set from the given type-bound subsets.
// It prepends the namespace code to all assertion and set IDs and builds
// the type index for efficient lookup during validation.
//
// Unlike Register, the returned set is NOT added to the global registry.
// Use Set.Validate to validate objects against it directly.
//
// The input sets are cloned internally, so the same output of For can safely
// be passed to multiple NewSet or Register calls.
func NewSet(ns Code, sets ...*Set) *Set {
	sets = cloneSets(sets)
	set := &Set{
		ID:      ns,
		Subsets: sets,
	}
	prependToSets(ns, sets)
	buildTypeIndex(set)
	return set
}

// For creates a new set of rules for the provided object (struct or value type).
// Each Def is applied in order to build up the set's assertions and subsets.
// Assert, Field, Each, Object, and When all return Def values that can be passed here.
//
// We let the compiler know that this function should not be "inlined"
// so that the package the caller is in can be detected reliably at runtime.
//
//go:noinline
func For(obj any, defs ...Def) *Set {
	// Detect whether the direct caller is in the same package as obj. When true,
	// the package segment is omitted from the set ID, since the registration
	// namespace will supply it. This avoids double-encoding, e.g. org.Email
	// registered under GOBL-ORG yields GOBL-ORG-EMAIL, not GOBL-ORG-ORG-EMAIL.
	var callerPkg string
	if pc, _, _, ok := runtime.Caller(1); ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			name := fn.Name() // e.g. "github.com/invopop/gobl/org.init"
			if i := strings.LastIndex(name, "."); i >= 0 {
				callerPkg = name[:i]
			}
		}
	}

	t := reflect.TypeOf(obj)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	normPkg := func(p string) string { return strings.TrimSuffix(p, "_test") }
	samePackage := callerPkg != "" && normPkg(callerPkg) == normPkg(t.PkgPath())
	setID := typeSetID(t, samePackage)
	objName := t.Name()
	if pkg := pkgShortName(t); pkg != "" {
		objName = pkg + "." + objName
	}
	s := &Set{
		ID:      setID,
		Object:  objName,
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

// presentGuard is an internal Test used by AssertIfPresent to skip nil or
// empty values without creating a dependency on the is package.
type presentGuard struct{}

func (presentGuard) Check(value any) bool {
	value, isNil := Indirect(value)
	return !isNil && !IsEmpty(value)
}

func (presentGuard) String() string { return "present" }

// AssertIfPresent returns a Def that adds an assertion that is only evaluated
// when the current scoped value is non-nil and non-empty. Use this for
// optional fields that have format or content constraints.
func AssertIfPresent(id Code, desc string, tests ...Test) Def {
	return func(s *Set) {
		subset := &Set{Guard: presentGuard{}}
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
//	        rules.Assert("02", "line required", is.Present),  // per-element
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
			if err := ct.Compile(obj); err != nil {
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
	if ft.Kind() == reflect.Pointer {
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
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		panic(fmt.Sprintf("rules: Each used on non-slice type %s", t.Name()))
	}
	ft := t.Elem()
	if ft.Kind() == reflect.Pointer {
		ft = ft.Elem()
	}
	ss.objType = ft
	nestedProto := reflect.New(ft).Interface()
	compileAndResolve(ft, nestedProto, ss)
}

// fieldTypeByName returns the reflect.Type for the field with the given JSON
// tag name in struct type t. Anonymous (embedded) struct fields are searched
// recursively so that promoted fields are resolved. Returns nil if not found.
func fieldTypeByName(t reflect.Type, name string) reflect.Type {
	if t.Kind() != reflect.Struct {
		return nil
	}
	for i := range t.NumField() {
		f := t.Field(i)
		if f.Anonymous {
			et := f.Type
			if et.Kind() == reflect.Pointer {
				et = et.Elem()
			}
			if et.Kind() == reflect.Struct {
				if inner := fieldTypeByName(et, name); inner != nil {
					return inner
				}
			}
			continue
		}
		if jsonFieldName(f) == name {
			return f.Type
		}
	}
	return nil
}

// fieldValueByName returns the reflect.Value for the field with the given JSON
// tag name in struct value rv. Anonymous (embedded) struct fields are searched
// recursively. Returns (zero, false) if not found.
func fieldValueByName(rv reflect.Value, name string) (reflect.Value, bool) {
	if rv.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}
	rt := rv.Type()
	for i := range rt.NumField() {
		f := rt.Field(i)
		if f.Anonymous {
			fv := rv.Field(i)
			if fv.Kind() == reflect.Pointer {
				if fv.IsNil() {
					continue
				}
				fv = fv.Elem()
			}
			if fv.Kind() == reflect.Struct {
				if inner, ok := fieldValueByName(fv, name); ok {
					return inner, true
				}
			}
			continue
		}
		if jsonFieldName(f) == name {
			return rv.Field(i), true
		}
	}
	return reflect.Value{}, false
}

func compileAssertions(env any, asserts ...*Assertion) {
	for _, a := range asserts {
		for _, t := range a.Tests {
			if ct, ok := t.(compilableTest); ok {
				if err := ct.Compile(env); err != nil {
					panic(fmt.Sprintf("failed to compile assertion %s: %s", a.ID, err.Error()))
				}
			}
		}
	}
}

// typeSetID derives a set ID from the type. When samePackage is true (the For
// caller is in the same package as the type), the package segment is omitted so
// that the registration namespace supplies it without duplication.
// For example: tax.Identity called from outside → TAX-IDENTITY;
//
//	tax.Identity called from within tax → IDENTITY.
func typeSetID(t reflect.Type, samePackage bool) Code {
	if samePackage {
		return Code(strings.ToUpper(t.Name()))
	}
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

// AllSets returns all rule sets registered in the global registry.
func AllSets() []*Set {
	return allSets()
}

// allSets returns the concatenation of all registry partitions.
func allSets() []*Set {
	all := make([]*Set, 0, len(coreRegistry)+len(unkeyedGuarded)+len(guardedIndex))
	all = append(all, coreRegistry...)
	all = append(all, unkeyedGuarded...)
	for _, sets := range guardedIndex {
		all = append(all, sets...)
	}
	return all
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
	rc := &Context{}
	for _, opt := range opts {
		opt(rc)
	}
	collectContext(rc, obj)
	var faults []*Fault

	// Always apply core (unguarded) rule sets.
	for _, ns := range coreRegistry {
		if fs := ns.validate(rc, obj); fs != nil {
			faults = append(faults, fs.List()...)
		}
	}

	// Apply guarded sets indexed by context key — only iterate sets
	// whose key is present in the current validation context.
	seen := make(map[*Set]struct{})
	for _, key := range rc.Keys() {
		for _, ns := range guardedIndex[key] {
			if _, ok := seen[ns]; ok {
				continue // a set may be indexed under multiple keys
			}
			seen[ns] = struct{}{}
			if fs := ns.validate(rc, obj); fs != nil {
				faults = append(faults, fs.List()...)
			}
		}
	}

	// Apply guarded sets that could not be indexed by key.
	for _, ns := range unkeyedGuarded {
		if fs := ns.validate(rc, obj); fs != nil {
			faults = append(faults, fs.List()...)
		}
	}

	return newFaults(faults...)
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
