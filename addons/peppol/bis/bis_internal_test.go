package bis

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPartyCountry(t *testing.T) {
	t.Run("nil party", func(t *testing.T) {
		assert.Equal(t, l10n.Code(""), partyCountry(nil))
	})
	t.Run("from TaxID", func(t *testing.T) {
		p := &org.Party{TaxID: &tax.Identity{Country: "DE"}}
		assert.Equal(t, l10n.DE, partyCountry(p))
	})
	t.Run("falls back to address country", func(t *testing.T) {
		p := &org.Party{Addresses: []*org.Address{{Country: "ES"}}}
		assert.Equal(t, l10n.ES, partyCountry(p))
	})
	t.Run("empty when neither", func(t *testing.T) {
		assert.Equal(t, l10n.Code(""), partyCountry(&org.Party{}))
	})
	t.Run("nil address ignored", func(t *testing.T) {
		p := &org.Party{Addresses: []*org.Address{nil}}
		assert.Equal(t, l10n.Code(""), partyCountry(p))
	})
}

func TestInvoiceSupplierCountry(t *testing.T) {
	assert.Equal(t, l10n.Code(""), invoiceSupplierCountry(nil))
	assert.Equal(t, l10n.Code(""), invoiceSupplierCountry("not an invoice"))
	assert.Equal(t, l10n.Code(""), invoiceSupplierCountry(&bill.Invoice{}))
	assert.Equal(t, l10n.DE, invoiceSupplierCountry(&bill.Invoice{
		Supplier: &org.Party{TaxID: &tax.Identity{Country: "DE"}},
	}))
}

func TestSupplierCountryIs(t *testing.T) {
	t.Run("matches", func(t *testing.T) {
		test := supplierCountryIs(l10n.DE)
		inv := &bill.Invoice{Supplier: &org.Party{TaxID: &tax.Identity{Country: "DE"}}}
		assert.True(t, test.Check(inv))
	})
	t.Run("does not match", func(t *testing.T) {
		test := supplierCountryIs(l10n.DE)
		inv := &bill.Invoice{Supplier: &org.Party{TaxID: &tax.Identity{Country: "FR"}}}
		assert.False(t, test.Check(inv))
	})
}

func TestNormalizeIdentity(t *testing.T) {
	t.Run("nil identity", func(t *testing.T) {
		assert.NotPanics(t, func() { normalizeIdentity(nil) })
	})
	t.Run("non-fskatt key untouched", func(t *testing.T) {
		id := &org.Identity{Key: "other"}
		normalizeIdentity(id)
		assert.Equal(t, cbc.Code(""), id.Code)
	})
	t.Run("fskatt fills code", func(t *testing.T) {
		id := &org.Identity{Key: IdentityKeyFSkatt}
		normalizeIdentity(id)
		assert.Equal(t, cbc.Code(FSkattText), id.Code)
	})
	t.Run("dispatcher routes identities", func(t *testing.T) {
		id := &org.Identity{Key: IdentityKeyFSkatt}
		normalize(id)
		assert.Equal(t, cbc.Code(FSkattText), id.Code)
	})
	t.Run("dispatcher ignores non-identities", func(t *testing.T) {
		assert.NotPanics(t, func() { normalize("anything") })
		assert.NotPanics(t, func() { normalize(nil) })
	})
}
