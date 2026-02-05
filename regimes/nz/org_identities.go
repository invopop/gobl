package nz

import (
	"errors"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/gs1"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/validation"
)

const nzbnGS1Prefix = "94"

var orgIdentityDefinitions = []*cbc.Definition{
	{
		Key: org.IdentityKeyGLN,
		Name: i18n.String{
			i18n.EN: "NZ Business Number",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The NZBN is a globally unique 13-digit identifier automatically assigned to
				companies registered with the New Zealand Companies Office. Other entities
				can apply voluntarily. Based on the GS1 Global Location Number (GLN) standard,
				it always starts with the New Zealand GS1 prefix '94'. For Peppol e-invoicing,
				the NZBN is used with the '0088:' scheme identifier (e.g., 0088:9429041234567).
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "NZBN Official",
				},
				URL: "https://www.nzbn.govt.nz/whats-an-nzbn/about/",
			},
			{
				Title: i18n.String{
					i18n.EN: "GS1 Check Digit Calculator",
				},
				URL: "https://www.gs1.org/services/check-digit-calculator",
			},
		},
	},
}

func normalizeNZBN(id *org.Identity) {
	id.Code = cbc.Code(normalizeCodeString(cbc.NormalizeString(id.Code.String())))
}

func validateNZBNIdentity(id *org.Identity) error {
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Required,
			validation.By(checkNZBN),
			validation.Skip,
		),
	)
}

func checkNZBN(value interface{}) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	if !gs1.CheckGLN(code) {
		return errors.New("invalid NZBN: must be a valid 13-digit GLN")
	}
	if !gs1.HasPrefix(code, nzbnGS1Prefix) {
		return errors.New("NZBN must start with '94' (New Zealand GS1 prefix)")
	}
	return nil
}
