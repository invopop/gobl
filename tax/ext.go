package tax

import (
	"context"
	"errors"
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// ExtMap is a map of extension keys to either a code or a key.
type ExtMap map[cbc.Key]cbc.KeyOrCode

// Validate ensures the extension map data looks correct.
func (em ExtMap) ValidateWithContext(ctx context.Context) error {
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
	for k, kc := range em {
		kd := r.ExtensionDef(k)
		if kd == nil {
			err[k.String()] = errors.New("undefined")
			continue
		}
		if len(kd.Codes) > 0 && !kd.HasCode(kc.Code()) {
			err[k.String()] = fmt.Errorf("code '%s' invalid", kc)
		}
		if len(kd.Keys) > 0 && !kd.HasKey(kc.Key()) {
			err[k.String()] = fmt.Errorf("key '%s' invalid", kc)
		}
	}
	if len(err) > 0 {
		return err
	}
	return nil
}

// Has returns true if the code map has values for all the provided keys.
func (em ExtMap) Has(keys ...cbc.Key) bool {
	for _, k := range keys {
		if _, ok := em[k]; !ok {
			return false
		}
	}
	return true
}

// Equals returns true if the code map has the same keys and values as the provided
// map.
func (em ExtMap) Equals(other ExtMap) bool {
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

// ExtMapHas returns a validation rule that ensures the extension map's
// keys match those provided.
func ExtMapHas(keys ...cbc.Key) validation.Rule {
	return validateCodeMap{keys: keys}
}

// ExtMapRequires returns a validation rule that ensures all the
// extension map's keys match those provided in the list.
func ExtMapRequires(keys ...cbc.Key) validation.Rule {
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
	em, ok := value.(ExtMap)
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
func (ExtMap) JSONSchemaExtend(schema *jsonschema.Schema) {
	prop := schema.AdditionalProperties
	schema.AdditionalProperties = nil
	schema.PatternProperties = map[string]*jsonschema.Schema{
		cbc.KeyPattern: prop,
	}
}
