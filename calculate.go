package gobl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/invopop/gobl/schema"
)

// CalculateOption configures the package-level Calculate operation.
type CalculateOption func(*calculateOptions)

type calculateOptions struct {
	discrepancies bool
}

// WithDiscrepancies asks Calculate to reject explicitly supplied calculated
// fields that do not match the values produced by GOBL. The returned error can
// be matched to CalculationDiscrepancies with errors.As.
func WithDiscrepancies() CalculateOption {
	return func(opts *calculateOptions) {
		opts.discrepancies = true
	}
}

// Calculate parses and calculates a GOBL document or envelope. By default it
// performs the same calculation as the value's Calculate method.
//
// WithDiscrepancies additionally rejects explicitly supplied calculated
// fields that do not match the values produced by GOBL. The returned error can
// be matched to CalculationDiscrepancies with errors.As.
func Calculate(data []byte, options ...CalculateOption) (any, error) {
	opts := new(calculateOptions)
	for _, option := range options {
		if option != nil {
			option(opts)
		}
	}

	if opts.discrepancies {
		value, discrepancies, err := calculateWithDiscrepancies(data)
		if err != nil {
			return nil, err
		}
		if len(discrepancies) > 0 {
			return nil, discrepancies
		}
		return value, nil
	}

	value, err := Parse(data)
	if err != nil {
		return nil, err
	}
	if err := calculateParsed(value); err != nil {
		return nil, err
	}
	return value, nil
}

// CalculationDiscrepancy describes a calculated value supplied by the caller
// that did not match the value produced by GOBL.
type CalculationDiscrepancy struct {
	Path       string          `json:"path"`
	Provided   json.RawMessage `json:"provided"`
	Calculated json.RawMessage `json:"calculated"`
}

// CalculationDiscrepancies contains the calculated values that GOBL changed.
// An empty list means that every calculated value explicitly present in the
// input matched the result. Calculated values omitted from the input are not
// discrepancies.
type CalculationDiscrepancies []*CalculationDiscrepancy

// Error provides a human-readable summary suitable for logs. Callers should
// use the structured discrepancy fields when preparing an API response.
func (ds CalculationDiscrepancies) Error() string {
	if len(ds) == 0 {
		return ""
	}
	if len(ds) == 1 {
		return ds[0].Error()
	}
	return fmt.Sprintf("%d calculation discrepancies", len(ds))
}

// Error provides a human-readable description of the discrepancy.
func (d *CalculationDiscrepancy) Error() string {
	return fmt.Sprintf(
		"%s: provided %s does not match calculated %s",
		d.Path,
		string(d.Provided),
		string(d.Calculated),
	)
}

// CalculateWithDiscrepancies parses and calculates a GOBL document or
// envelope, returning both the calculated value and any calculated fields that
// were explicitly supplied with a different value.
//
// Fields are considered calculated when they or one of their parents has the
// `calculated=true` JSON Schema annotation. Values are compared semantically
// when their type provides an Equals method, so equivalent representations
// such as monetary amounts with different precision do not produce a
// discrepancy.
func calculateWithDiscrepancies(data []byte) (any, CalculationDiscrepancies, error) {
	provided, err := Parse(data)
	if err != nil {
		return nil, nil, err
	}
	calculated, err := Parse(data)
	if err != nil {
		return nil, nil, err
	}

	if err := calculateParsed(calculated); err != nil {
		return nil, nil, err
	}
	result := calculated

	var raw json.RawMessage = data
	path := "$"
	if env, ok := calculated.(*Envelope); ok {
		var root map[string]json.RawMessage
		if err := json.Unmarshal(data, &root); err != nil {
			return nil, nil, ErrInput.WithCause(err)
		}
		raw = root["doc"]
		path = "$.doc"
		provided = calculationPayload(provided)
		calculated = calculationPayload(env)
	}

	ds := make(CalculationDiscrepancies, 0)
	compareCalculated(raw, reflect.ValueOf(provided), reflect.ValueOf(calculated), path, false, &ds)
	return result, ds, nil
}

func calculateParsed(obj any) error {
	if env, ok := obj.(*Envelope); ok {
		return env.Calculate()
	}
	doc, err := schema.NewObject(obj)
	if err != nil {
		return wrapError(err)
	}
	if err := doc.Calculate(); err != nil {
		return ErrCalculation.WithCause(err)
	}
	return nil
}

func calculationPayload(obj any) any {
	if env, ok := obj.(*Envelope); ok {
		return env.Extract()
	}
	return obj
}

func compareCalculated(raw json.RawMessage, before, after reflect.Value, path string, inherited bool, out *CalculationDiscrepancies) {
	if isJSONNull(raw) {
		return
	}

	before = indirectValue(before)
	after = indirectValue(after)

	switch firstJSONByte(raw) {
	case '{':
		compareCalculatedObject(raw, before, after, path, inherited, out)
	case '[':
		compareCalculatedArray(raw, before, after, path, inherited, out)
	default:
		if inherited && !valuesEqual(before, after) {
			*out = append(*out, newCalculationDiscrepancy(path, before, after))
		}
	}
}

func compareCalculatedObject(raw json.RawMessage, before, after reflect.Value, path string, inherited bool, out *CalculationDiscrepancies) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return
	}

	typ := valueType(before, after)
	if typ == nil || typ.Kind() != reflect.Struct {
		if inherited && !valuesEqual(before, after) {
			*out = append(*out, newCalculationDiscrepancy(path, before, after))
		}
		return
	}

	compareStructFields(fields, before, after, typ, path, inherited, out)
}

func compareStructFields(raw map[string]json.RawMessage, before, after reflect.Value, typ reflect.Type, path string, inherited bool, out *CalculationDiscrepancies) {
	for i := range typ.NumField() {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}
		name, embedded := jsonFieldName(field)
		if name == "-" {
			continue
		}
		if embedded {
			compareStructFields(raw, fieldValue(before, i), fieldValue(after, i), indirectType(field.Type), path, inherited || isCalculatedField(field), out)
			continue
		}
		value, present := raw[name]
		if !present || isJSONNull(value) {
			continue
		}
		compareCalculated(
			value,
			fieldValue(before, i),
			fieldValue(after, i),
			appendPath(path, name),
			inherited || isCalculatedField(field),
			out,
		)
	}
}

func compareCalculatedArray(raw json.RawMessage, before, after reflect.Value, path string, inherited bool, out *CalculationDiscrepancies) {
	var values []json.RawMessage
	if err := json.Unmarshal(raw, &values); err != nil {
		return
	}
	for i, value := range values {
		compareCalculated(
			value,
			indexValue(before, i),
			indexValue(after, i),
			fmt.Sprintf("%s[%d]", path, i),
			inherited,
			out,
		)
	}
}

func valuesEqual(before, after reflect.Value) bool {
	if !before.IsValid() || !after.IsValid() {
		return before.IsValid() == after.IsValid()
	}
	before = indirectValue(before)
	after = indirectValue(after)
	if !before.IsValid() || !after.IsValid() {
		return before.IsValid() == after.IsValid()
	}
	if before.Type() != after.Type() {
		return false
	}
	if equal, ok := callEquals(before, after); ok {
		return equal
	}
	return reflect.DeepEqual(before.Interface(), after.Interface())
}

func callEquals(before, after reflect.Value) (bool, bool) {
	method := before.MethodByName("Equals")
	if !method.IsValid() {
		return false, false
	}
	typ := method.Type()
	if typ.NumIn() != 1 || typ.In(0) != after.Type() || typ.NumOut() != 1 || typ.Out(0).Kind() != reflect.Bool {
		return false, false
	}
	result := method.Call([]reflect.Value{after})
	return result[0].Bool(), true
}

func newCalculationDiscrepancy(path string, before, after reflect.Value) *CalculationDiscrepancy {
	return &CalculationDiscrepancy{
		Path:       path,
		Provided:   marshalValue(before),
		Calculated: marshalValue(after),
	}
}

func marshalValue(value reflect.Value) json.RawMessage {
	value = indirectValue(value)
	if !value.IsValid() {
		return json.RawMessage("null")
	}
	data, err := json.Marshal(value.Interface())
	if err != nil {
		return json.RawMessage("null")
	}
	return data
}

func valueType(values ...reflect.Value) reflect.Type {
	for _, value := range values {
		value = indirectValue(value)
		if value.IsValid() {
			return value.Type()
		}
	}
	return nil
}

func indirectType(typ reflect.Type) reflect.Type {
	for typ != nil && (typ.Kind() == reflect.Pointer || typ.Kind() == reflect.Interface) {
		typ = typ.Elem()
	}
	return typ
}

func indirectValue(value reflect.Value) reflect.Value {
	for value.IsValid() && (value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface) {
		if value.IsNil() {
			return reflect.Value{}
		}
		value = value.Elem()
	}
	return value
}

func fieldValue(value reflect.Value, index int) reflect.Value {
	value = indirectValue(value)
	if !value.IsValid() || value.Kind() != reflect.Struct || index >= value.NumField() {
		return reflect.Value{}
	}
	return value.Field(index)
}

func indexValue(value reflect.Value, index int) reflect.Value {
	value = indirectValue(value)
	if !value.IsValid() || (value.Kind() != reflect.Slice && value.Kind() != reflect.Array) || index >= value.Len() {
		return reflect.Value{}
	}
	return value.Index(index)
}

func jsonFieldName(field reflect.StructField) (string, bool) {
	tag := field.Tag.Get("json")
	name := strings.Split(tag, ",")[0]
	if name == "-" {
		return "-", false
	}
	if name != "" {
		return name, false
	}
	if field.Anonymous {
		return "", true
	}
	return field.Name, false
}

func isCalculatedField(field reflect.StructField) bool {
	return hasCalculatedAnnotation(field.Tag.Get("jsonschema_extras")) ||
		hasCalculatedAnnotation(field.Tag.Get("jsonschema"))
}

func hasCalculatedAnnotation(tag string) bool {
	parts := strings.Split(tag, ",")
	sort.Strings(parts)
	i := sort.SearchStrings(parts, "calculated=true")
	return i < len(parts) && parts[i] == "calculated=true"
}

func appendPath(path, name string) string {
	if path == "$" {
		return "$." + name
	}
	return path + "." + name
}

func firstJSONByte(raw json.RawMessage) byte {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return 0
	}
	return trimmed[0]
}

func isJSONNull(raw json.RawMessage) bool {
	return bytes.Equal(bytes.TrimSpace(raw), []byte("null"))
}
