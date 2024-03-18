package tax

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Extensions is a map of extension keys to values.
type Extensions map[cbc.Key]ExtValue

// ExtValue is a string value that has helper methods to help determine
// if it is a code, key, or regular string.
type ExtValue string

// ValidateWithContext ensures the extension map data looks correct.
func (em Extensions) ValidateWithContext(ctx context.Context) error {
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
	r := RegimeFromContext(ctx)
	if r == nil {
		return nil
	}
	// Validate keys are defined in regime
	for k, ev := range em {
		ks := k.String()
		kd := r.ExtensionDef(k)
		if kd == nil {
			err[ks] = errors.New("undefined")
			continue
		}
		if len(kd.Codes) > 0 && !kd.HasCode(ev.Code()) {
			err[ks] = fmt.Errorf("code '%s' invalid", ev)
		}
		if len(kd.Keys) > 0 && !kd.HasKey(ev.Key()) {
			err[ks] = fmt.Errorf("key '%s' invalid", ev)
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

// Has returns true if the code map has values for all the provided keys.
func (em Extensions) Has(keys ...cbc.Key) bool {
	for _, k := range keys {
		if _, ok := em[k]; !ok {
			return false
		}
	}
	return true
}

// Equals returns true if the code map has the same keys and values as the provided
// map.
func (em Extensions) Equals(other Extensions) bool {
	if len(em) != len(other) {
		return false
	}
	for k, v := range em {
		v2, ok := other[k]
		if !ok {
			return false
		}
		if v2 != v {
			return false
		}
	}
	return true
}

// NormalizeExtensions will try to clean the extension map removing empty values
// and will potentially return a nil if there only keys with no values.
func NormalizeExtensions(em Extensions) Extensions {
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

// ExtensionsHas returns a validation rule that ensures the extension map's
// keys match those provided.
func ExtensionsHas(keys ...cbc.Key) validation.Rule {
	return validateCodeMap{keys: keys}
}

// ExtensionsRequires returns a validation rule that ensures all the
// extension map's keys match those provided in the list.
func ExtensionsRequires(keys ...cbc.Key) validation.Rule {
	return validateCodeMap{
		required: true,
		keys:     keys,
	}
}

type validateCodeMap struct {
	keys     []cbc.Key
	required bool
}

func (v validateCodeMap) Validate(value interface{}) error {
	em, ok := value.(Extensions)
	if !ok {
		return nil
	}
	err := make(validation.Errors)

	if v.required {
		for _, k := range v.keys {
			if _, ok := em[k]; !ok {
				err[k.String()] = errors.New("required")
			}
		}
	} else {
		for k := range em {
			if !k.In(v.keys...) {
				err[k.String()] = errors.New("invalid")
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

// String provides the string representation.
func (ev ExtValue) String() string {
	return string(ev)
}

// Key returns the key value or empty if the value is a Code.
func (ev ExtValue) Key() cbc.Key {
	k := cbc.Key(ev)
	if err := k.Validate(); err == nil {
		return k
	}
	return cbc.KeyEmpty
}

// Code returns the code value or empty if the value is a Key.
func (ev ExtValue) Code() cbc.Code {
	c := cbc.Code(ev)
	if err := c.Validate(); err == nil {
		return c
	}
	return cbc.CodeEmpty
}
