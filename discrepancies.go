package gobl

import (
	"bytes"
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
// ancestors carries the `calculated=true` extension — whose value in data
// differs from the corresponding value in calculated.
//
// A calculated field left out of data is never reported: GOBL is expected
// to fill those in, and doing so is not a discrepancy. Only fields the
// caller supplied with a value GOBL then changed are reported.
//
// An empty, non-nil result means every calculated value the caller
// supplied matched GOBL's calculation.
func FindCalculationDiscrepancies(data []byte, calculated any) (CalculationDiscrepancies, error) {
	raw := json.RawMessage(data)
	path := "$"
	value := reflect.ValueOf(calculated)

	if env, ok := calculated.(*Envelope); ok {
		var envelope struct {
			Doc json.RawMessage `json:"doc"`
		}
		if err := json.Unmarshal(data, &envelope); err != nil {
			return nil, ErrInput.WithCause(err)
		}
		raw = envelope.Doc
		path = "$.doc"
		value = reflect.ValueOf(env.Extract())
	}

	return findDiscrepancies(raw, value, false, path), nil
}

// findDiscrepancies walks raw, the original JSON for this location, in
// step with value, the equivalent already-calculated Go value. calculated
// is true once this location is inside a field whose schema marks it (or
// an ancestor) as calculated=true; everything beneath such a field is
// treated as calculated too, since GOBL only annotates the outermost
// calculated field of a subtree such as an invoice's totals.
func findDiscrepancies(raw json.RawMessage, value reflect.Value, calculated bool, path string) CalculationDiscrepancies {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || isNull(raw) {
		// A missing or explicit null calculated value is never a
		// discrepancy: GOBL is expected to fill it in.
		return nil
	}

	value = indirect(value)

	switch raw[0] {
	case '{':
		return findObjectDiscrepancies(raw, value, calculated, path)
	case '[':
		return findArrayDiscrepancies(raw, value, calculated, path)
	default:
		if !calculated {
			return nil
		}
		return compareValue(raw, value, path)
	}
}

func findObjectDiscrepancies(raw json.RawMessage, value reflect.Value, calculated bool, path string) CalculationDiscrepancies {
	if value.Kind() != reflect.Struct {
		// Not a Go struct we can walk field by field (e.g. a map used for
		// free-form extensions). Compare it as a single opaque value.
		if !calculated {
			return nil
		}
		return compareValue(raw, value, path)
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return nil
	}

	var discrepancies CalculationDiscrepancies
	for _, field := range reflect.VisibleFields(value.Type()) {
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
		fieldRaw, present := fields[name]
		if !present {
			continue
		}
		fieldValue, err := value.FieldByIndexErr(field.Index)
		if err != nil {
			continue
		}
		discrepancies = append(discrepancies, findDiscrepancies(
			fieldRaw, fieldValue, calculated || isCalculatedField(field), path+"."+name,
		)...)
	}
	return discrepancies
}

func findArrayDiscrepancies(raw json.RawMessage, value reflect.Value, calculated bool, path string) CalculationDiscrepancies {
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		if !calculated {
			return nil
		}
		return compareValue(raw, value, path)
	}

	var items []json.RawMessage
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil
	}

	var discrepancies CalculationDiscrepancies
	for i, item := range items {
		var elem reflect.Value
		if i < value.Len() {
			elem = value.Index(i)
		}
		discrepancies = append(discrepancies, findDiscrepancies(
			item, elem, calculated, fmt.Sprintf("%s[%d]", path, i),
		)...)
	}
	return discrepancies
}

// compareValue is reached once raw can no longer be broken down any
// further against value: either it's a JSON scalar, or it's an object or
// array that value's Go type can't be walked field by field or index by
// index. It parses raw into a fresh instance of value's type so the two
// can be compared with the same semantics GOBL itself uses.
func compareValue(raw json.RawMessage, value reflect.Value, path string) CalculationDiscrepancies {
	if !value.IsValid() {
		return CalculationDiscrepancies{{Path: path, Provided: raw, Calculated: json.RawMessage("null")}}
	}

	provided := reflect.New(value.Type())
	if err := json.Unmarshal(raw, provided.Interface()); err != nil {
		// raw was already parsed successfully as part of the original
		// document, so a mismatch here means value's shape no longer
		// matches raw; there's nothing meaningful left to compare.
		return nil
	}
	if equal(provided.Elem(), value) {
		return nil
	}

	calculated, err := json.Marshal(value.Interface())
	if err != nil {
		return nil
	}
	return CalculationDiscrepancies{{Path: path, Provided: raw, Calculated: calculated}}
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

func isNull(raw json.RawMessage) bool {
	return string(raw) == "null"
}
