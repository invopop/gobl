package ctc

import (
	"testing"

	"github.com/invopop/gobl/addons/fr/ctc/flow10"
	"github.com/invopop/gobl/addons/fr/ctc/flow6"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func frenchByTaxID() *org.Party {
	return &org.Party{TaxID: &tax.Identity{Country: "FR", Code: "732829320"}}
}

func TestPartyIsFrench(t *testing.T) {
	assert.False(t, partyIsFrench(nil))
	assert.False(t, partyIsFrench(&org.Party{}))

	t.Run("French TaxID", func(t *testing.T) {
		assert.True(t, partyIsFrench(frenchByTaxID()))
	})
	t.Run("SIREN identity type", func(t *testing.T) {
		p := &org.Party{Identities: []*org.Identity{nil, {Type: fr.IdentityTypeSIREN, Code: "1"}}}
		assert.True(t, partyIsFrench(p))
	})
	t.Run("SIREN iso scheme 0002", func(t *testing.T) {
		p := &org.Party{Identities: []*org.Identity{
			{Code: "1", Ext: tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: "0002"})},
		}}
		assert.True(t, partyIsFrench(p))
	})
	t.Run("non-French party", func(t *testing.T) {
		p := &org.Party{TaxID: &tax.Identity{Country: "DE", Code: "111111125"}}
		assert.False(t, partyIsFrench(p))
	})
}

func TestDispatchPayment(t *testing.T) {
	assert.Equal(t, flow10.V1, dispatchPayment(nil))

	t.Run("request type goes to flow10", func(t *testing.T) {
		assert.Equal(t, flow10.V1, dispatchPayment(&bill.Payment{Type: bill.PaymentTypeRequest}))
	})
	t.Run("receipt between two French parties goes to flow6", func(t *testing.T) {
		pmt := &bill.Payment{Type: bill.PaymentTypeReceipt, Supplier: frenchByTaxID(), Customer: frenchByTaxID()}
		assert.Equal(t, flow6.V1, dispatchPayment(pmt))
	})
	t.Run("receipt with a non-French party goes to flow10", func(t *testing.T) {
		pmt := &bill.Payment{Type: bill.PaymentTypeReceipt, Supplier: frenchByTaxID()}
		assert.Equal(t, flow10.V1, dispatchPayment(pmt))
	})
}
