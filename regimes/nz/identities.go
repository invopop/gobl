package nz

import (
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/validation"
)

// IdentityKeyIRD is the identity key for the NZ Inland Revenue Department number.
const IdentityKeyIRD cbc.Key = "nz-ird"

var identityKeyDefinitions = []*cbc.Definition{
	{
		Key: IdentityKeyIRD,
		Name: i18n.String{
			i18n.EN: "IRD Number",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The IRD number is the primary tax identifier in New Zealand, issued by
				Inland Revenue to individuals, companies, trusts, and other entities.
				Format: 8 or 9 digits (XXX-XXX-XXX when displayed). Valid range:
				10,000,000 to 200,000,000. For GST-registered businesses, the IRD
				number also serves as their GST number.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "IRD Numbers Overview",
				},
				URL: "https://www.ird.govt.nz/managing-my-tax/ird-numbers",
			},
		},
	},
}

func normalizeCodeString(code string) string {
	code = strings.ReplaceAll(code, "-", "")
	code = strings.ReplaceAll(code, " ", "")
	return code
}

// normalizeIdentity handles normalization for org.Identity objects.
func normalizeIdentity(id *org.Identity) {
	if id == nil {
		return
	}
	switch id.Key {
	case IdentityKeyIRD:
		id.Code = cbc.Code(normalizeCodeString(id.Code.String()))
	case org.IdentityKeyGLN:
		normalizeNZBN(id)
	}
}

// validateIdentity checks org.Identity objects for valid IRD or NZBN codes.
func validateIdentity(id *org.Identity) error {
	if id == nil {
		return nil
	}
	switch id.Key {
	case IdentityKeyIRD:
		return validation.ValidateStruct(id,
			validation.Field(&id.Code,
				validation.By(validateTaxCode),
				validation.Skip,
			),
		)
	case org.IdentityKeyGLN:
		return validateNZBNIdentity(id)
	}
	return nil
}
