package nl

import (
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	// IdentityTypeKVK represents the Dutch "Kamer van Koophandel" (Chamber of Commerce)
	// registration number used to identify businesses in the Netherlands.
	IdentityTypeKVK cbc.Code = "KVK"

	// IdentityTypeOIN represents the Dutch "Organisatie Identificatie Nummer" used
	// to identify government organizations in the Netherlands.
	IdentityTypeOIN cbc.Code = "OIN"
)

var identityDefinitions = []*cbc.Definition{
	{
		Code: IdentityTypeKVK,
		Name: i18n.String{
			i18n.EN: "KVK Number",
			i18n.NL: "KVK-nummer",
		},
	},
	{
		Code: IdentityTypeOIN,
		Name: i18n.String{
			i18n.EN: "OIN Number",
			i18n.NL: "OIN-nummer",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("OIN Stelsel - Samenstelling OIN"),
				URL:   "https://gitdocumentatie.logius.nl/publicatie/dk/oin/2.2.1/#samenstelling-oin",
			},
		},
	},
}

// oinRegexp validates the OIN format: 6-zero prefix, 2-digit register code, 9-digit identifier, 3-zero suffix (20 digits).
var oinRegexp = regexp.MustCompile(`^0{6}(0[1-9]|10|99)\d{9}0{3}$`)

func validateIdentity(id *org.Identity) error {
	if id == nil {
		return nil
	}
	switch id.Type {
	case IdentityTypeKVK:
		return validation.ValidateStruct(id,
			validation.Field(&id.Code, validation.Length(8, 8)),
		)
	case IdentityTypeOIN:
		return validation.ValidateStruct(id,
			validation.Field(&id.Code,
				validation.Match(oinRegexp),
			),
		)
	}
	return nil
}
