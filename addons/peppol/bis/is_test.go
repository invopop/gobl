package bis

import (
	"testing"

	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValAsParty(t *testing.T) {
	assert.Nil(t, valAsParty(nil))
	assert.Nil(t, valAsParty("string"))
	p := &org.Party{}
	assert.Equal(t, p, valAsParty(p))
}

func TestISDocumentTypeValid(t *testing.T) {
	assert.True(t, isDocumentTypeValid(nil))
	assert.True(t, isDocumentTypeValid(tax.Extensions{}))
	assert.True(t, isDocumentTypeValid(tax.Extensions{untdid.ExtKeyDocumentType: "380"}))
	assert.True(t, isDocumentTypeValid(tax.Extensions{untdid.ExtKeyDocumentType: "381"}))
	assert.False(t, isDocumentTypeValid(tax.Extensions{untdid.ExtKeyDocumentType: "326"}))
}

func TestPartyHasLegalIdentity(t *testing.T) {
	assert.True(t, partyHasLegalIdentity(nil))
	assert.False(t, partyHasLegalIdentity(&org.Party{}))
	assert.True(t, partyHasLegalIdentity(&org.Party{
		Identities: []*org.Identity{{Scope: "legal", Code: "X"}},
	}))
	assert.True(t, partyHasLegalIdentity(&org.Party{
		TaxID: &tax.Identity{Code: "X"},
	}))
	// Identities without legal scope but with TaxID still passes through TaxID branch.
	assert.True(t, partyHasLegalIdentity(&org.Party{
		Identities: []*org.Identity{{Scope: "tax", Code: "X"}},
		TaxID:      &tax.Identity{Code: "Y"},
	}))
}

func TestFirstAddressStreetAndCode(t *testing.T) {
	assert.True(t, firstAddressStreetAndCode(nil))
	assert.True(t, firstAddressStreetAndCode([]*org.Address{}))
	assert.False(t, firstAddressStreetAndCode([]*org.Address{nil}))
	assert.False(t, firstAddressStreetAndCode([]*org.Address{{Street: "X"}}))
	assert.False(t, firstAddressStreetAndCode([]*org.Address{{Code: "1"}}))
	assert.True(t, firstAddressStreetAndCode([]*org.Address{{Street: "X", Code: "1"}}))
}

func TestValidISAccount(t *testing.T) {
	assert.True(t, validISAccount("123456789012"))      // 12-digit domestic
	assert.True(t, validISAccount("IS140159260076545510730339")) // IS IBAN
	assert.True(t, validISAccount("IS14 0159 2600 7654 5510 7303 39")) // IBAN with spaces
	assert.False(t, validISAccount(""))
	assert.False(t, validISAccount("12345"))
	assert.False(t, validISAccount("DE89370400440532013000")) // non-IS IBAN
}

func TestISPaymentCodes(t *testing.T) {
	// Code 9
	assert.True(t, isPaymentCode9Account(nil))
	assert.True(t, isPaymentCode9Account(&pay.Instructions{Ext: payExt("30")}))
	assert.False(t, isPaymentCode9Account(&pay.Instructions{Ext: payExt("9")})) // no transfers
	assert.True(t, isPaymentCode9Account(&pay.Instructions{
		Ext:            payExt("9"),
		CreditTransfer: []*pay.CreditTransfer{{Number: "123456789012"}},
	}))
	assert.True(t, isPaymentCode9Account(&pay.Instructions{
		Ext:            payExt("9"),
		CreditTransfer: []*pay.CreditTransfer{{IBAN: "IS140159260076545510730339"}},
	}))
	assert.False(t, isPaymentCode9Account(&pay.Instructions{
		Ext:            payExt("9"),
		CreditTransfer: []*pay.CreditTransfer{{Number: "12345"}},
	}))

	// Code 42
	assert.True(t, isPaymentCode42Account(nil))
	assert.True(t, isPaymentCode42Account(&pay.Instructions{Ext: payExt("30")}))
	assert.False(t, isPaymentCode42Account(&pay.Instructions{Ext: payExt("42")}))
	assert.True(t, isPaymentCode42Account(&pay.Instructions{
		Ext:            payExt("42"),
		CreditTransfer: []*pay.CreditTransfer{{Number: "123456789012"}},
	}))
	assert.False(t, isPaymentCode42Account(&pay.Instructions{
		Ext:            payExt("42"),
		CreditTransfer: []*pay.CreditTransfer{{Number: "AAA"}},
	}))
}
