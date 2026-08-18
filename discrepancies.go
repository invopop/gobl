package gobl

import (
	"encoding"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// CalculationDiscrepancy describes a single calculated value that was
// explicitly present in the original data with a value different to the
// one GOBL's calculation produced.
type CalculationDiscrepancy struct {
	// Path identifies the field's location using the same "$.foo[0].bar"
	// notation as validation fault paths.
	Path string `json:"path"`
	// Provided is the value exactly as it appeared in the original data.
	Provided json.RawMessage `json:"provided"`
	// Calculated is the value GOBL produced for the same field.
	Calculated json.RawMessage `json:"calculated"`
}

// CalculationDiscrepancies lists calculated values that GOBL replaced
// with a different result during calculation.
type CalculationDiscrepancies []*CalculationDiscrepancy

// FindCalculationDiscrepancies compares data, the original bytes a caller
// submitted, against calculated, the same document or envelope obtained by
// parsing that data and then calling Calculate on it. It reports every
// calculated field — one whose JSON Schema definition or one of its
// ancestors carries the `calculated=true` extension — whose value changed
// during calculation.
//
// data is parsed again to recover the values as they were before
// calculation, then compared against calculated field by field using the
// same reflection GOBL already relies on to walk documents; no separate
// JSON tree walk is needed.
//
// A calculated field left out of data is never reported: GOBL is expected
// to fill those in, and doing so is not a discrepancy. Because a Go zero
// value (an empty string, a zero date, index zero) is indistinguishable
// from an omitted field for the handful of calculated fields that aren't
// pointers, a calculated field explicitly supplied with its zero value is
// treated the same way: not worth reporting, since GOBL is expected to
// replace it regardless of whether it was left out or supplied as such.
//
// An empty, non-nil result means every calculated value the caller
// supplied matched GOBL's calculation.
func FindCalculationDiscrepancies(data []byte, calculated any) (CalculationDiscrepancies, error) {
	provided, err := Parse(data)
	if err != nil {
		return nil, err
	}

	path := "$"
	providedValue := reflect.ValueOf(provided)
	calculatedValue := reflect.ValueOf(calculated)

	if env, ok := calculated.(*Envelope); ok {
		providedEnv, ok := provided.(*Envelope)
		if !ok {
			return nil, ErrInput.WithReason("calculated is an envelope, but data is not")
		}
		providedValue = reflect.ValueOf(providedEnv.Extract())
		calculatedValue = reflect.ValueOf(env.Extract())
		path = "$.doc"
	}

	return findDiscrepancies(providedValue, calculatedValue, false, path), nil
}

// findDiscrepancies walks provided, the value as it was before
// calculation, in step with calculated, the equivalent already-calculated
// value. isCalculated is true once this location is inside a field whose
// schema marks it (or an ancestor) as calculated=true; everything beneath
// such a field is treated as calculated too, since GOBL only annotates the
// outermost calculated field of a subtree such as an invoice's totals.
func findDiscrepancies(provided, calculated reflect.Value, isCalculated bool, path string) CalculationDiscrepancies {
	// A pointer or interface still holding a value once dereferenced was
	// explicitly supplied, even if that value happens to be zero; only a
	// non-pointer field's zero value is ambiguous with an omission.
	nilable := isNilable(provided)

	provided = indirect(provided)
	calculated = indirect(calculated)

	if !provided.IsValid() {
		// Left out of the original data entirely: never a discrepancy.
		return nil
	}
	if !calculated.IsValid() {
		if !isCalculated {
			return nil
		}
		return CalculationDiscrepancies{newDiscrepancy(path, provided, reflect.Value{})}
	}

	// A type such as num.Amount is a Go struct, but marshals to a single
	// JSON string via its own MarshalText/MarshalJSON rather than one
	// object key per (unexported) field; walking it field by field would
	// find nothing; it needs to be treated as a leaf, same as a plain
	// scalar.
	if !hasCustomMarshaling(provided.Type()) {
		switch provided.Kind() {
		case reflect.Struct:
			return findFieldDiscrepancies(provided, calculated, isCalculated, path)
		case reflect.Slice, reflect.Array:
			return findElementDiscrepancies(provided, calculated, isCalculated, path)
		}
	}

	if !isCalculated {
		return nil
	}
	if !nilable && provided.IsZero() {
		return nil
	}
	if equal(provided, calculated) {
		return nil
	}
	return CalculationDiscrepancies{newDiscrepancy(path, provided, calculated)}
}

var (
	jsonMarshalerType = reflect.TypeFor[json.Marshaler]()
	textMarshalerType = reflect.TypeFor[encoding.TextMarshaler]()
)

// hasCustomMarshaling reports whether t controls its own JSON
// representation (e.g. num.Amount marshals to a decimal string) rather
// than being encoded as one object key per field, meaning it should be
// compared as an opaque value instead of walked field by field.
func hasCustomMarshaling(t reflect.Type) bool {
	return t.Implements(jsonMarshalerType) || t.Implements(textMarshalerType) ||
		reflect.PointerTo(t).Implements(jsonMarshalerType) || reflect.PointerTo(t).Implements(textMarshalerType)
}

func findFieldDiscrepancies(provided, calculated reflect.Value, isCalculated bool, path string) CalculationDiscrepancies {
	var discrepancies CalculationDiscrepancies
	for _, field := range reflect.VisibleFields(provided.Type()) {
		// Anonymous (embedded) fields are also returned by VisibleFields
		// alongside the fields they promote, so only the promoted fields
		// need handling here.
		if !field.IsExported() || field.Anonymous {
			continue
		}
		name, ok := fieldName(field)
		if !ok {
			continue
		}
		providedField, err := provided.FieldByIndexErr(field.Index)
		if err != nil {
			continue
		}
		calculatedField, err := calculated.FieldByIndexErr(field.Index)
		if err != nil {
			continue
		}
		discrepancies = append(discrepancies, findDiscrepancies(
			providedField, calculatedField, isCalculated || isCalculatedField(field), path+"."+name,
		)...)
	}
	return discrepancies
}

func findElementDiscrepancies(provided, calculated reflect.Value, isCalculated bool, path string) CalculationDiscrepancies {
	var discrepancies CalculationDiscrepancies
	for i := 0; i < provided.Len(); i++ {
		var calculatedElem reflect.Value
		if i < calculated.Len() {
			calculatedElem = calculated.Index(i)
		}
		discrepancies = append(discrepancies, findDiscrepancies(
			provided.Index(i), calculatedElem, isCalculated, fmt.Sprintf("%s[%d]", path, i),
		)...)
	}
	return discrepancies
}

func newDiscrepancy(path string, provided, calculated reflect.Value) *CalculationDiscrepancy {
	return &CalculationDiscrepancy{Path: path, Provided: marshal(provided), Calculated: marshal(calculated)}
}

func marshal(value reflect.Value) json.RawMessage {
	if !value.IsValid() {
		return json.RawMessage("null")
	}
	data, err := json.Marshal(value.Interface())
	if err != nil {
		return json.RawMessage("null")
	}
	return data
}

// equal compares two values of the same type, preferring an Equals method
// when the type provides one so equivalent representations (e.g. "20" and
// "20.00" for a monetary amount) are not reported as discrepancies.
func equal(provided, calculated reflect.Value) bool {
	if method := provided.MethodByName("Equals"); method.IsValid() {
		t := method.Type()
		if t.NumIn() == 1 && t.In(0) == calculated.Type() && t.NumOut() == 1 && t.Out(0).Kind() == reflect.Bool {
			return method.Call([]reflect.Value{calculated})[0].Bool()
		}
	}
	return reflect.DeepEqual(provided.Interface(), calculated.Interface())
}

// isNilable reports whether value's declared type can be nil, i.e.
// whether it being present with a zero underlying value is distinguishable
// from it being absent altogether.
func isNilable(value reflect.Value) bool {
	return value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface
}

// indirect follows pointers and interfaces down to the concrete value
// they hold, returning the zero Value if any step along the way is nil.
func indirect(value reflect.Value) reflect.Value {
	for value.IsValid() && (value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface) {
		if value.IsNil() {
			return reflect.Value{}
		}
		value = value.Elem()
	}
	return value
}

// isCalculatedField reports whether field's JSON Schema definition carries
// the `calculated=true` extension, which GOBL's schema generator accepts
// either inside `jsonschema_extras` or alongside other options in
// `jsonschema` itself.
func isCalculatedField(field reflect.StructField) bool {
	return hasCalculatedTag(field.Tag.Get("jsonschema_extras")) || hasCalculatedTag(field.Tag.Get("jsonschema"))
}

func hasCalculatedTag(tag string) bool {
	for part := range strings.SplitSeq(tag, ",") {
		if part == "calculated=true" {
			return true
		}
	}
	return false
}

// fieldName returns field's JSON name and whether it's addressable at all
// via JSON, i.e. it isn't tagged `json:"-"`.
func fieldName(field reflect.StructField) (string, bool) {
	name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
	if name == "-" {
		return "", false
	}
	if name == "" {
		name = field.Name
	}
	return name, true
}
