package rules

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
)

type exprTest struct {
	expr *vm.Program
	test string
}

// Expr returns a new expression-base test
func Expr(test string, args ...any) Test {
	return &exprTest{
		test: fmt.Sprintf(test, args...),
		expr: nil, // compiled lazily when embedded in a Set
	}
}

func (et *exprTest) String() string {
	return et.test
}

func (et *exprTest) Check(val any) bool {
	if et.expr == nil {
		panic("expression test was not compiled; wrap it in ForStruct() or ForValue()")
	}
	env := et.buildEnv(val)
	result, err := expr.Run(et.expr, env)
	if err != nil {
		panic("expression runtime error: " + err.Error())
	}
	return result.(bool)
}

func (et *exprTest) compile(val any) error {
	env := et.buildEnv(val)

	// Backslashes in tests must be doubled for embedding inside an expr
	// double-quoted string literal (expr interprets \\ as a single \).
	etStr := strings.ReplaceAll(et.test, `\`, `\\`)
	prog, err := expr.Compile(etStr,
		expr.AsBool(),
		expr.WithTag("json"),
		expr.Env(env),
	)
	if err != nil {
		return fmt.Errorf("invalid test expression %q: %w", et.test, err)
	}
	et.expr = prog
	return nil
}

// buildEnv constructs the expression environment for the given type and value.
// Struct values are returned directly so that expressions can reference fields
// by their json tag names. Non-struct types are exposed as a map with a single
// "this" key holding the underlying primitive value, so that assertion
// expressions can reference the value uniformly regardless of type.
func (et *exprTest) buildEnv(val any) any {
	rv := reflect.ValueOf(val)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	t := rv.Type()
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
