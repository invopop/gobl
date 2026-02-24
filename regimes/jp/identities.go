package jp

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	// IdentityKeyRegistrationNumber represents the Qualified Invoice Issuer
	// Registration Number (適格請求書発行事業者登録番号) assigned under
	// Japan's Qualified Invoice System (QIS), effective October 1, 2023.
	// Format: "T" followed by 13 digits.
	IdentityKeyRegistrationNumber cbc.Key = "jp-invoice-registration-number"
)

var registrationNumberPattern = regexp.MustCompile(`^T\d{13}$`)

var identityDefinitions = []*cbc.Definition{
	{
		Key: IdentityKeyRegistrationNumber,
		Name: i18n.String{
			i18n.EN: "Invoice Registration Number",
			i18n.JA: "適格請求書発行事業者登録番号",
		},
	},
}

func normalizeRegistrationNumber(id *org.Identity) {
	if id == nil || id.Key != IdentityKeyRegistrationNumber {
		return
	}
	code := strings.ToUpper(strings.TrimSpace(id.Code.String()))
	id.Code = cbc.Code(code)
}

func validateRegistrationNumber(id *org.Identity) error {
	if id == nil || id.Key != IdentityKeyRegistrationNumber {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Required,
			validation.Match(registrationNumberPattern),
			validation.Skip,
		),
	)
}
