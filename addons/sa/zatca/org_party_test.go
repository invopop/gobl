package zatca_test

import (
	"testing"

	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// nonSAAddress returns an address valid for a non-SA country. The
// SA-ZATCA address rules 02-08 are gated by countryCodeIsSA so the
// SA-specific constraints (4-digit building number, 5-digit postal code,
// district name, etc.) shouldn't fire — but the bill_invoices.go rules
// that require a street name and city for standard tax invoices DO still
// fire, so we satisfy those minimums.
func nonSAAddress(country l10n.ISOCountryCode) *org.Address {
	return &org.Address{
		Street:   "Rue de la Paix",
		Locality: "Paris",
		Country:  country,
	}
}

// ============================================================================
// Rule 01: address must have a country code (always)
// ============================================================================

func TestOrgAddressRule01_CountryRequired(t *testing.T) {
	t.Run("supplier address missing country fails", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Addresses[0].Country = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "address must have a country code")
	})

	t.Run("customer address missing country fails", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Addresses[0].Country = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "address must have a country code")
	})

	t.Run("non-SA country skips SA-specific rules", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Addresses[0] = nonSAAddress("FR")
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

// ============================================================================
// Rules 02-08: SA-specific address constraints
// ============================================================================

func TestOrgAddressRule02_StreetRequired(t *testing.T) {
	t.Run("missing street fails (BR-KSA-09, BR-KSA-63)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Addresses[0].Street = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"address in SA must have a street name")
	})
}

func TestOrgAddressRule03_BuildingNumberRequired(t *testing.T) {
	t.Run("missing building number fails (BR-KSA-09, BR-KSA-63)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Addresses[0].Number = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"address in SA must have a 4 digits building number")
	})
}

func TestOrgAddressRule04_BuildingNumberMustBe4Digits(t *testing.T) {
	t.Run("3-digit building number rejected (BR-KSA-37)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Addresses[0].Number = "123"
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"address in SA must have a 4 digits building number")
	})

	t.Run("non-numeric building number rejected (BR-KSA-37)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Addresses[0].Number = "12A4"
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"address in SA must have a 4 digits building number")
	})

	t.Run("5-digit building number rejected (BR-KSA-37)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Addresses[0].Number = "12345"
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"address in SA must have a 4 digits building number")
	})

	t.Run("4-digit building number accepted", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Addresses[0].Number = "1234"
		require.NoError(t, rules.Validate(calculated(t, inv)))
	})
}

func TestOrgAddressRule05_PostalCodeRequired(t *testing.T) {
	t.Run("missing postal code fails (BR-KSA-09, BR-KSA-63, BR-KSA-67)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Addresses[0].Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"address in SA must have a 5 digits postal code")
	})
}

func TestOrgAddressRule06_PostalCodeMustBe5Digits(t *testing.T) {
	t.Run("4-digit postal code rejected (BR-KSA-66)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Addresses[0].Code = "1234"
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"address in SA must have a 5 digits postal code")
	})

	t.Run("non-numeric postal code rejected (BR-KSA-66)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Addresses[0].Code = "1234A"
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"address in SA must have a 5 digits postal code")
	})

	t.Run("6-digit postal code rejected (BR-KSA-66)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Addresses[0].Code = "123456"
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"address in SA must have a 5 digits postal code")
	})

	t.Run("5-digit postal code accepted", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Addresses[0].Code = "54321"
		require.NoError(t, rules.Validate(calculated(t, inv)))
	})
}

func TestOrgAddressRule07_CityRequired(t *testing.T) {
	t.Run("missing city name fails (BR-KSA-09, BR-KSA-63)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Addresses[0].Locality = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"address in SA must have a city name")
	})
}

func TestOrgAddressRule08_DistrictRequired(t *testing.T) {
	t.Run("missing district (street_extra) fails (BR-KSA-09, BR-KSA-63)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Addresses[0].StreetExtra = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"address in SA must have a district name")
	})
}

// ============================================================================
// Cross-cutting: a fully invalid SA address surfaces every SA-specific error
// ============================================================================

func TestOrgAddress_MultipleViolationsReported(t *testing.T) {
	inv := validStandardInvoice()
	inv.Supplier.Addresses[0].Street = ""
	inv.Supplier.Addresses[0].Number = ""
	inv.Supplier.Addresses[0].StreetExtra = ""
	inv.Supplier.Addresses[0].Code = ""
	inv.Supplier.Addresses[0].Locality = ""
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "address in SA must have a street name")
	assert.ErrorContains(t, err, "address in SA must have a 4 digits building number")
	assert.ErrorContains(t, err, "address in SA must have a 5 digits postal code")
	assert.ErrorContains(t, err, "address in SA must have a city name")
	assert.ErrorContains(t, err, "address in SA must have a district name")
}
