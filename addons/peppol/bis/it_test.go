package bis

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestItalianTaxIDLength(t *testing.T) {
	assert.True(t, italianTaxIDLength(nil))
	assert.True(t, italianTaxIDLength(&org.Party{}))
	assert.True(t, italianTaxIDLength(&org.Party{TaxID: &tax.Identity{}}))
	assert.True(t, italianTaxIDLength(&org.Party{TaxID: &tax.Identity{Code: "12345678901"}}))        // 11
	assert.True(t, italianTaxIDLength(&org.Party{TaxID: &tax.Identity{Code: "1234567890123456"}}))   // 16
	assert.False(t, italianTaxIDLength(&org.Party{TaxID: &tax.Identity{Code: "1234567890"}}))        // 10
	assert.False(t, italianTaxIDLength(&org.Party{TaxID: &tax.Identity{Code: "12345678901234567"}})) // 17
}

func TestFirstAddressHasStreet(t *testing.T) {
	assert.True(t, firstAddressHasStreet(nil))
	assert.True(t, firstAddressHasStreet([]*org.Address{}))
	assert.False(t, firstAddressHasStreet([]*org.Address{nil}))
	assert.False(t, firstAddressHasStreet([]*org.Address{{}}))
	assert.True(t, firstAddressHasStreet([]*org.Address{{Street: "Via Roma"}}))
}
