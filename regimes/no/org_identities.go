//go:build ignore

package template

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	// IdentityTypeTemplate defines the key for the template identity.
	IdentityTypeTemplate cbc.Code = "Template" // Find official international keys
)

func identityTypeDefinitions() []*cbc.Definition {
	return []*cbc.Definition{
		{
			Code: IdentityTypeTemplate,
			Name: i18n.String{
				i18n.EN: "Template",
				// Add official local name here.
				// i18n.XX: "Template",
			},
			Desc: i18n.String{
				i18n.EN: "A template identity description",
				// Add official local description here.
				// i18n.XX: "A template identity description",
			},
		},
	}
}

// normalizeOrgIdentity performs normalization specific to the regime.
//
//   - Explanation 1
//   - Explanation 2
//
// Some edge cases.
func normalizeOrgIdentity(id *org.Identity) {
	if id == nil {
		return
	}

	switch id.Type {
	case IdentityTypeTemplate:
		// Handle normalization here for each Identity type. This is just an example.
		// cbc has extra methods to help with this.
		code := cbc.NormalizeNumericalCode(id.Code).String()

		id.Code = cbc.Code(code)

	default:
		return
	}
}

// validateOrgIdentity performs validation for the regime.
// Assumes the code has already been normalized.
//
//   - Explanation 1
//   - Explanation 2
//
// If the number is not valid, it returns an error.
//
// If the identity type is not valid, it returns nil.
func validateOrgIdentity(id *org.Identity) error {
	if id == nil {
		return nil
	}

	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.By(func(value any) error {
				code, ok := value.(cbc.Code)
				if !ok || code == "" {
					return nil
				}

				// Handle validation here for each Identity type. This is just an example.

				switch id.Type {
				case IdentityTypeTemplate:
					// ...
				default:
					return nil
				}

				return nil
			}),
			validation.Skip,
		),
	)
}
