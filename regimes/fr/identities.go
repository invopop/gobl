package fr

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	taxCodeSIRENRegexp = regexp.MustCompile(`^\d{9}$`)
	taxCodeSIRETRegexp = regexp.MustCompile(`^\d{14}$`)
)

func normalizeIdentity(id *org.Identity) {
	if id == nil || id.Type != IdentityTypeSIREN && id.Type != IdentityTypeSIRET {
		return
	}
	code := strings.ToUpper(id.Code.String())
	code = tax.IdentityCodeBadCharsRegexp.ReplaceAllString(code, "")
	code = strings.TrimPrefix(code, string(l10n.FR))
	id.Code = cbc.Code(code)
}

func validateIdentity(id *org.Identity) error {
	if id == nil || id.Type != IdentityTypeSIREN && id.Type != IdentityTypeSIRET {
		return nil
	}

	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Required,
			validation.By(validateSIRENTaxCode),
			validation.Skip,
		),
	)
}
