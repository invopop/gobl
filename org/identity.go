package org

import (
	"fmt"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
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
	IdentityKeyGLN       cbc.Key = "gln"       // GS1 Global Location Number
	IdentityKeyGTIN      cbc.Key = "gtin"      // GS1 Global Trade Item Number
	IdentityKeyEAN       cbc.Key = "ean"       // European Article Number
	IdentityKeyUPC       cbc.Key = "upc"       // UPC (Universal Product Code)
	IdentityKeyIMEI      cbc.Key = "imei"      // International Mobile Equipment Identity
	IdentityKeyDUNS      cbc.Key = "duns"      // Dun & Bradstreet D-U-N-S Number
	IdentityKeyNCM       cbc.Key = "ncm"       // Mercosur Common Nomenclature
	IdentityKeyOther     cbc.Key = "other"
)

// Identity scopes that may be used to further classify an identity's intended use.
const (
	IdentityScopeTax   cbc.Key = "tax"
	IdentityScopeLegal cbc.Key = "legal"
)

// Identity is used to define a code for a specific context. Identities can be used for
// a variety of purposes, such as identifying a person, organisation, item, or document.
type Identity struct {
	uuid.Identify
	// Optional label useful for non-standard identities to give a bit more context.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Scope defines the context in which this identity is meant to be used.
	Scope cbc.Key `json:"scope,omitempty" jsonschema:"title=Scope"`
	// Country from which the identity was issued.
	Country l10n.ISOCountryCode `json:"country,omitempty" jsonschema:"title=Country"`
	// Uniquely classify this identity using a key instead of a Type.
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
func (i *Identity) Normalize() {
	if i == nil {
		return
	}
	uuid.Normalize(&i.UUID)
	i.Label = cbc.NormalizeString(i.Label)
	i.Type = cbc.NormalizeCode(i.Type)
	i.Code = cbc.NormalizeCode(i.Code)
	i.Description = cbc.NormalizeString(i.Description)
	i.Ext = tax.CleanExtensions(i.Ext)
}

func identityRules() *rules.Set {
	return rules.For(new(Identity),
		rules.Field("code",
			rules.Assert("01", "identity code must be provided", is.Present),
		),
		rules.Field("scope",
			rules.AssertIfPresent("02", "identity scope when provided must be either 'tax' or 'legal'",
				is.In(IdentityScopeTax, IdentityScopeLegal),
			),
		),
		rules.Assert("03", "identity must have either a key or type defined, but not both",
			identityHasKeyOrTypeNotBoth(),
		),
	)
}

func identityHasKeyOrTypeNotBoth() rules.Test {
	return is.Func("key and type must not be used together", func(value any) bool {
		id, _ := value.(*Identity)
		if id == nil {
			return false
		}
		return id.Key == "" || id.Type == ""
	})
}

// IdentityTypeIn provides a test that will determine if the identity defined has a type with one of the defined codes.
func IdentityTypeIn(typ ...cbc.Code) rules.Test {
	return identitiesTest{
		desc: fmt.Sprintf("type in [%s]", strings.Join(cbc.CodeStrings(typ), ", ")),
		typs: typ,
	}
}

// IdentityKeyIn provides a test that will determine if the identity defined has a key with one of the defined keys.
func IdentityKeyIn(key ...cbc.Key) rules.Test {
	return identitiesTest{
		desc: fmt.Sprintf("key in [%s]", strings.Join(cbc.KeyStrings(key), ", ")),
		keys: key,
	}
}

// IdentitiesTypeIn provides a test that will determine if at least one
// of the identities defined includes one with the defined type.
func IdentitiesTypeIn(typ ...cbc.Code) rules.Test {
	return identitiesTest{
		desc: fmt.Sprintf("has a type in [%s]", strings.Join(cbc.CodeStrings(typ), ", ")),
		typs: typ,
	}
}

// IdentitiesKeyIn provides a test that will determine if at least one
// of the identities defined includes one with one of the defined keys.
func IdentitiesKeyIn(key ...cbc.Key) rules.Test {
	return identitiesTest{
		desc: fmt.Sprintf("has a key in [%s]", strings.Join(cbc.KeyStrings(key), ", ")),
		keys: key,
	}
}

type identitiesTest struct {
	desc string
	typs []cbc.Code
	keys []cbc.Key
}

// Check will determine if the provided identities array complies with the criteria.
func (v identitiesTest) Check(value any) bool {
	switch obj := value.(type) {
	case *Identity:
		return v.matches(obj)
	case []*Identity:
		for _, row := range obj {
			if v.matches(row) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// String provides a description of the test being performed.
func (v identitiesTest) String() string {
	return v.desc
}

func (v identitiesTest) matches(row *Identity) bool {
	return (len(v.typs) == 0 || row.Type.In(v.typs...)) &&
		(len(v.keys) == 0 || row.Key.In(v.keys...))
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

// IdentityForExtKey helps return the identity with the first matching extension key.
func IdentityForExtKey(in []*Identity, key cbc.Key) *Identity {
	for _, v := range in {
		if v.Ext.Get(key) != cbc.CodeEmpty {
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

// JSONSchemaExtend adds extra details to the schema.
func (Identity) JSONSchemaExtend(js *jsonschema.Schema) {
	prop, ok := js.Properties.Get("scope")
	if ok {
		prop.OneOf = []*jsonschema.Schema{
			{
				Const: IdentityScopeTax,
				Title: "Tax",
			},
			{
				Const: IdentityScopeLegal,
				Title: "Legal",
			},
		}
	}
}
