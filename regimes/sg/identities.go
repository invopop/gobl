package sg

import (
	"errors"
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Reference: https://mytax.iras.gov.sg/ESVWeb/default.aspx?target=GSTListingSearch

const (
	IdentityKeyGSTNumber cbc.Key = "sg-gst-number"
)

var GSTNumberRegexPattern = regexp.MustCompile(`^[M][A-Z0-9]\d{7}[A-Z]$`)

var identityDefinitions = []*cbc.Definition{
	{
		Key: IdentityKeyGSTNumber,
		Name: i18n.String{
			i18n.EN: "GST Registration Number",
		},
		Desc: i18n.String{
			i18n.EN: "GST Registration Number is a number given to any business entity that is registered for GST with IRAS. Overseas suppliers who register for GST also receive one. GST-registered suppliers are required to print their GST Registration Number on every tax invoice and receipt issued. In most cases GST registration number is the same as the UEN.",
		},
	},
}

func normalizeIdentity(id *org.Identity) {
	if id == nil || id.Key != IdentityKeyGSTNumber {
		return
	}
	code := strings.ToUpper(id.Code.String())
	code = tax.IdentityCodeBadCharsRegexp.ReplaceAllString(code, "")
	id.Code = cbc.Code(strings.TrimPrefix(code, string(l10n.SG)))

}

func validateIdentity(id *org.Identity) error {
	if id == nil || id.Key != IdentityKeyGSTNumber {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Required,
			validation.By(validateGSTNumber),
			validation.Skip,
		),
	)
}

func validateGSTNumber(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}

	val := code.String()

	if !GSTNumberRegexPattern.MatchString(val) {
		return errors.New("invalid format")
	}

	return nil
}
