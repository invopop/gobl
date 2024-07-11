package gr

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

func validateAddress(value any) error {
	a, ok := value.(*org.Address)
	if !ok || a == nil {
		return nil
	}
	return validation.ValidateStruct(a,
		validation.Field(&a.Locality, validation.Required),
		validation.Field(&a.Code, validation.Required),
	)
}
