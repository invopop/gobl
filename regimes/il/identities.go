package il

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/luhn"
	"github.com/invopop/validation"
)

const (
	// IdentityKeyPersonalID represents the Israeli national identity number
	// (Mispar Zehut / מספר זהות), a 9-digit number with a Luhn check digit.
	// Used by sole proprietors (Osek Patur/Zair) who operate below the VAT
	// registration threshold and do not have a Mispar Osek Murshe.
	IdentityKeyPersonalID cbc.Key = "il-personal-id"

	// IdentityKeyCorporationNumber represents the corporation registration
	// number issued by the Corporations Authority (רשות התאגידים).
	// It is a 9-digit number where the prefix indicates entity type: 51 (companies),
	// 50 (public institutions), 56 (foreign non-profit corporations), 58 (associations).
	IdentityKeyCorporationNumber cbc.Key = "il-company-id"
)

var (
	nineDigitsRegex = regexp.MustCompile(`^\d{9}$`)
)

var identityDefinitions = []*cbc.Definition{
	{
		Key: IdentityKeyPersonalID,
		Name: i18n.String{
			i18n.EN: "Personal ID",
			i18n.HE: "מספר זהות",
		},
		Desc: i18n.String{
			i18n.EN: "Israeli national identity number (Mispar Zehut / Teudat Zehut), a 9-digit number validated with the Luhn algorithm. Used by sole proprietors who are not VAT-registered.",
		},
	},
	{
		Key: IdentityKeyCorporationNumber,
		Name: i18n.String{
			i18n.EN: "Corporation Number",
			i18n.HE: "מספר תאגיד",
		},
		Desc: i18n.String{
			i18n.EN: "Corporation registration number issued by the Corporations Authority, a 9-digit number with a prefix indicating entity type.",
		},
	},
}

// normalizeIdentity strips non-numeric characters from identity codes.
func normalizeIdentity(id *org.Identity) {
	if id == nil {
		return
	}
	switch id.Key {
	case IdentityKeyPersonalID, IdentityKeyCorporationNumber:
		id.Code = cbc.NormalizeNumericalCode(id.Code)
	}
}

// validateIdentity checks that the identity code is valid for the given key.
func validateIdentity(id *org.Identity) error {
	if id == nil {
		return nil
	}
	switch id.Key {
	case IdentityKeyPersonalID:
		return validation.ValidateStruct(id,
			validation.Field(&id.Code,
				validation.Required,
				validation.Match(nineDigitsRegex),
				validation.By(checkLuhn),
				validation.Skip,
			),
		)
	case IdentityKeyCorporationNumber:
		return validation.ValidateStruct(id,
			validation.Field(&id.Code,
				validation.Required,
				validation.Match(nineDigitsRegex),
				validation.Skip,
			),
		)
	}
	return nil
}

func checkLuhn(value any) error {
	code, ok := value.(cbc.Code)
	if !ok || code == "" {
		return nil
	}
	if !luhn.Check(code) {
		return errors.New("invalid checksum")
	}
	return nil
}
