package org

import (
	"context"
	"fmt"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// Common identity keys that may be used to identify something, like an item, document,
// person, organisation, or company. Ideally, these will only be used when no other
// more structured properties are available inside GOBL. The keys suggested here are
// non-binding and can be used as a reference for other implementations or mappings to
// scheme identifiers such as UNTDID 1153.
const (
	IdentityKeySKU       cbc.Key = "sku"       // stock code unit ID
	IdentityKeyItem      cbc.Key = "item"      // item number
	IdentityKeyOrder     cbc.Key = "order"     // order number or code
	IdentityKeyAgreement cbc.Key = "agreement" // agreement number
	IdentityKeyContract  cbc.Key = "contract"  // contract number
	IdentityKeyPassport  cbc.Key = "passport"  // Passport number
	IdentityKeyNational  cbc.Key = "national"  // National ID card number
	IdentityKeyForeign   cbc.Key = "foreign"   // Foreigner ID card number
	IdentityKeyResident  cbc.Key = "resident"  // Resident ID card number
	IdentityKeyISBN      cbc.Key = "isbn"      // International Standard Book Number
	IdentityKeyHSN       cbc.Key = "hsn"       // Harmonized System of Nomenclature
	IdentityKeyGTIN      cbc.Key = "gtin"      // GS1 Global Trade Item Number
	IdentityKeyEAN       cbc.Key = "ean"       // European Article Number
	IdentityKeyUPC       cbc.Key = "upc"       // UPC (Universal Product Code)
	IdentityKeyIMEI      cbc.Key = "imei"      // International Mobile Equipment Identity
	IdentityKeyPeppol    cbc.Key = "peppol"    // Peppol Participant Identifier
	IdentityKeyDUNS      cbc.Key = "duns"      // Dun & Bradstreet D-U-N-S Number
	IdentityKeyOther     cbc.Key = "other"
)

// Identity is used to define a code for a specific context. Identities can be used for
// a variety of purposes, such as identifying a person, organisation, item, or document.
type Identity struct {
	uuid.Identify
	// Optional label useful for non-standard identities to give a bit more context.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Country from which the identity was issued.
	Country l10n.ISOCountryCode `json:"country,omitempty" jsonschema:"title=Country"`
	// Uniquely classify this identity using a key instead of a type.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// The type of Code being represented and usually specific for
	// a particular context, country, or tax regime, and cannot be used
	// alongside the key.
	Type cbc.Code `json:"type,omitempty" jsonschema:"title=Type"`
	// The actual value of the identity code.
	Code cbc.Code `json:"code" jsonschema:"title=Code"`
	// Description adds details about what the code could mean or imply
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// Ext provides a way to add additional information to the identity.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
}

// Normalize will try to clean the identity's data.
func (i *Identity) Normalize(normalizers tax.Normalizers) {
	if i == nil {
		return
	}
	uuid.Normalize(&i.UUID)
	i.Ext = tax.CleanExtensions(i.Ext)
	normalizers.Each(i)
}

// Validate ensures the identity looks valid.
func (i *Identity) Validate() error {
	return i.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures the identity looks valid inside the provided context.
func (i *Identity) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, i,
		validation.Field(&i.Label),
		validation.Field(&i.Country),
		validation.Field(&i.Key),
		validation.Field(&i.Type,
			validation.When(i.Key != "",
				validation.Empty.Error("must be empty when key is set"),
			),
		),
		validation.Field(&i.Code,
			validation.Required,
		),
		validation.Field(&i.Ext),
	)
}

// RequireIdentityType provides a validation rule that will determine if at least one
// of the identities defined includes one with the defined type.
func RequireIdentityType(typ cbc.Code) validation.Rule {
	return validateIdentitySet{typ: typ}
}

// RequireIdentityKey provides a validation rule that will determine if at least one
// of the identities defined includes one with one of the defined keys.
func RequireIdentityKey(key ...cbc.Key) validation.Rule {
	return validateIdentitySet{keys: key}
}

type validateIdentitySet struct {
	typ  cbc.Code
	keys []cbc.Key
}

func (v validateIdentitySet) Validate(value interface{}) error {
	ids, ok := value.([]*Identity)
	if !ok {
		return nil
	}
	for _, row := range ids {
		if v.matches(row) {
			return nil
		}
	}

	return fmt.Errorf("missing %s", v)
}

func (v validateIdentitySet) matches(row *Identity) bool {
	return (v.typ == cbc.CodeEmpty || row.Type == v.typ) &&
		(len(v.keys) == 0 || row.Key.In(v.keys...))
}

func (v validateIdentitySet) String() string {
	var parts []string
	if v.typ != cbc.CodeEmpty {
		parts = append(parts, fmt.Sprintf("type '%s'", v.typ))
	}
	if len(v.keys) != 0 {
		parts = append(parts, fmt.Sprintf("key '%s'", strings.Join(cbc.KeyStrings(v.keys), ", ")))
	}
	return strings.Join(parts, ", ")
}

// IdentityForType helps return the identity with a matching type code.
func IdentityForType(in []*Identity, typ cbc.Code) *Identity {
	for _, v := range in {
		if v.Type == typ {
			return v
		}
	}
	return nil
}

// IdentityForKey helps return the identity with the first matching key.
func IdentityForKey(in []*Identity, key ...cbc.Key) *Identity {
	for _, v := range in {
		if v.Key.In(key...) {
			return v
		}
	}
	return nil
}

// AddIdentity makes it easier to add a new identity to a list and replace an
// existing value with a matching type.
func AddIdentity(in []*Identity, i *Identity) []*Identity {
	if in == nil {
		return []*Identity{i}
	}
	for _, v := range in {
		if v.Type == i.Type && v.Key == i.Key {
			*v = *i // copy in place
			return in
		}
	}
	return append(in, i)
}
