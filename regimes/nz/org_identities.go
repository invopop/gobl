package nz

import (
	"errors"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/gs1"
	"github.com/invopop/validation"
)

// nzbnGS1Prefix is the GS1 prefix assigned to New Zealand.
const nzbnGS1Prefix = "94"

var orgIdentityDefinitions = []*cbc.Definition{
	{
		Key: org.IdentityKeyGLN,
		Name: i18n.String{
			i18n.EN: "NZ Business Number",
		},
		Desc: i18n.String{
			i18n.EN: "13-digit identifier based on the GS1 Global Location Number (GLN) standard, starting with NZ prefix 94.",
		},
	},
}

func normalizeNZBN(id *org.Identity) {
	code := cbc.NormalizeString(id.Code.String())
	code = strings.ReplaceAll(code, "-", "")
	code = strings.ReplaceAll(code, " ", "")
	id.Code = cbc.Code(code)
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
