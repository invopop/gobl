package sg

import (
	"errors"
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Reference: https://mytax.iras.gov.sg/ESVWeb/default.aspx?target=GSTListingSearch

const (
	// IdentityTypeUEN represents the Unique Entity Number used in Singapore.
	IdentityTypeUEN cbc.Code = "UEN"
)

var (
	taxCodeRegexps = []*regexp.Regexp{
		regexp.MustCompile(`^(19[0-9]{2}|20[0-9]{2})\d{5}[A-Z]$`), // UEN (ROC)
		regexp.MustCompile(`^\d{8}[A-Z]$`),                        // UEN (ROB)
		regexp.MustCompile(`^[TS]\d{2}[A-Z]{2}\d{4}[A-Z]$`),       // UEN (Others)
	}
)

var identityDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeUEN,
		Name: i18n.String{
			i18n.EN: "Unique Entity Number (UEN)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Unique Entity Number (UEN) is a number given to any business entity that is registered for GST with IRAS. Overseas suppliers who register for GST also receive one. GST-registered suppliers are required to print their UEN on every tax invoice and receipt issued.
			`),
		},
	},
}

func normalizeIdentity(id *org.Identity) {
	if id == nil || id.Type != IdentityTypeUEN {
		return
	}
	code := strings.ToUpper(id.Code.String())
	code = tax.IdentityCodeBadCharsRegexp.ReplaceAllString(code, "")
	id.Code = cbc.Code(strings.TrimPrefix(code, string(l10n.SG)))
}

func validateIdentity(id *org.Identity) error {
	if id == nil || id.Type != IdentityTypeUEN {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Required,
			validation.By(validateUEN),
			validation.Skip,
		),
	)
}

func validateUEN(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	val := code.String()
	match := false
	for _, re := range taxCodeRegexps {
		if re.MatchString(val) {
			match = true
			break
		}
	}
	if !match {
		return errors.New("invalid format")
	}

	return nil
}
