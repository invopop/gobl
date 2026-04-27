package flow10

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func frPartyWithSIREN() *org.Party {
	return &org.Party{
		Name: "Supplier SARL",
		TaxID: &tax.Identity{
			Country: "FR",
			Code:    "39356000000",
		},
		Identities: []*org.Identity{
			{
				Type:  fr.IdentityTypeSIREN,
				Code:  "356000000",
				Scope: org.IdentityScopeLegal,
				Ext: tax.ExtensionsOf(tax.ExtMap{
					iso.ExtKeySchemeID: "0002",
				}),
			},
		},
		Addresses: []*org.Address{{Country: "FR"}},
	}
}

func frCustomerWithSIREN() *org.Party {
	return &org.Party{
		Name: "Customer SAS",
		TaxID: &tax.Identity{
			Country: "FR",
			Code:    "44732829320",
		},
		Identities: []*org.Identity{
			{
				Type:  fr.IdentityTypeSIREN,
				Code:  "732829320",
				Scope: org.IdentityScopeLegal,
				Ext: tax.ExtensionsOf(tax.ExtMap{
					iso.ExtKeySchemeID: "0002",
				}),
			},
		},
		Addresses: []*org.Address{{Country: "FR"}},
	}
}

func testInvoiceB2B(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(V1),
		Code:      "INV-2026-001",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2026, 1, 15),
		Type:      bill.InvoiceTypeStandard,
		Supplier:  frPartyWithSIREN(),
		Customer:  frCustomerWithSIREN(),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Product",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{Category: tax.CategoryVAT, Percent: num.NewPercentage(20, 2)},
				},
			},
		},
	}
}

func testInvoiceB2C(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := &bill.Invoice{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(V1),
		Code:      "INV-2026-B2C-001",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2026, 1, 15),
		Type:      bill.InvoiceTypeStandard,
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(tax.ExtMap{
				ExtKeyB2CCategory: B2CCategoryGoods,
			}),
		},
		Supplier: frPartyWithSIREN(),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Product",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{Category: tax.CategoryVAT, Percent: num.NewPercentage(20, 2)},
				},
			},
		},
	}
	return inv
}

func TestInvoiceB2BHappyPath(t *testing.T) {
	inv := testInvoiceB2B(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestInvoiceB2CHappyPath(t *testing.T) {
	inv := testInvoiceB2C(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestInvoiceCurrencyRequiresEURConversion(t *testing.T) {
	inv := testInvoiceB2B(t)
	inv.Currency = "USD"
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "EUR")
}

func TestInvoiceCurrencyUSDWithExchangeRate(t *testing.T) {
	inv := testInvoiceB2B(t)
	inv.Currency = "USD"
	inv.ExchangeRates = []*currency.ExchangeRate{
		{From: "USD", To: "EUR", Amount: num.MakeAmount(875967, 6)},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestInvoiceB2BDocTypeNotAllowed(t *testing.T) {
	inv := testInvoiceB2B(t)
	require.NoError(t, inv.Calculate())
	// Force a document type that is not in the Flow 10 whitelist.
	inv.Tax.Ext = inv.Tax.Ext.Set(untdid.ExtKeyDocumentType, "325") // proforma, not allowed
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "Flow 10 permitted UNTDID 1001 codes")
}

func TestInvoiceB2BMissingBillingMode(t *testing.T) {
	inv := testInvoiceB2B(t)
	require.NoError(t, inv.Calculate())
	// Normalization defaults the billing mode; clear it to simulate a
	// downstream consumer that strips the extension.
	inv.Tax.Ext = inv.Tax.Ext.Delete(ExtKeyBillingMode)
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "billing mode")
}

func TestInvoiceB2BFinalAfterAdvanceRejectsDepositDocType(t *testing.T) {
	inv := testInvoiceB2B(t)
	require.NoError(t, inv.Calculate())
	inv.Tax.Ext = inv.Tax.Ext.
		Set(ExtKeyBillingMode, BillingModeM4).
		Set(untdid.ExtKeyDocumentType, "386") // Advance payment invoice
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "G1.60")
}

func TestInvoiceB2BSupplierRequiresAllowedScheme(t *testing.T) {
	inv := testInvoiceB2B(t)
	require.NoError(t, inv.Calculate())
	// Strip the supplier's identities so no allowed scheme remains.
	inv.Supplier.Identities = nil
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "supplier must declare a legal identity")
}

func TestInvoiceB2BAddressRequiresCountry(t *testing.T) {
	inv := testInvoiceB2B(t)
	require.NoError(t, inv.Calculate())
	inv.Supplier.Addresses = []*org.Address{{Street: "No country"}}
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "supplier address must include country")
}

func TestInvoiceB2BExemptRequiresSellerVATID(t *testing.T) {
	inv := testInvoiceB2B(t)
	inv.Lines[0].Taxes = tax.Set{
		{Category: tax.CategoryVAT, Key: tax.KeyExempt},
	}
	require.NoError(t, inv.Calculate())
	// Drop both potential VAT IDs; no ordering.seller either.
	inv.Supplier.TaxID = nil
	inv.Ordering = nil
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "supplier VAT ID or ordering.seller")
}

func TestInvoiceB2BExemptRequiresExemptTaxNote(t *testing.T) {
	inv := testInvoiceB2B(t)
	inv.Lines[0].Taxes = tax.Set{
		{Category: tax.CategoryVAT, Key: tax.KeyExempt},
	}
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "exemption reason")
}

func TestInvoiceB2CDefaultsCategoryToTNT1(t *testing.T) {
	inv := testInvoiceB2C(t)
	inv.Tax.Ext = inv.Tax.Ext.Delete(ExtKeyB2CCategory)
	require.NoError(t, inv.Calculate())
	assert.Equal(t, B2CCategoryNotTaxable, inv.Tax.Ext.Get(ExtKeyB2CCategory))
	require.NoError(t, rules.Validate(inv))
}

func TestInvoiceB2CSupplierRequiresSIREN(t *testing.T) {
	inv := testInvoiceB2C(t)
	// Clear both TaxID and Identities — party normalization would
	// otherwise regenerate a SIREN from the French TaxID.
	inv.Supplier.TaxID = nil
	inv.Supplier.Identities = nil
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "SIREN")
}

func TestInvoiceB2CVATRateNotInWhitelist(t *testing.T) {
	inv := testInvoiceB2C(t)
	inv.Lines[0].Taxes = tax.Set{
		{Category: tax.CategoryVAT, Percent: num.NewPercentage(17, 2)}, // 17%, not allowed
	}
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "G1.24")
}

func TestNormalizeDefaultBillingModeM1(t *testing.T) {
	inv := testInvoiceB2B(t)
	require.NoError(t, inv.Calculate())
	assert.Equal(t, BillingModeM1, inv.Tax.Ext.Get(ExtKeyBillingMode))
}

func TestNormalizeTaxCategorySetFromKey(t *testing.T) {
	inv := testInvoiceB2B(t)
	inv.Lines[0].Taxes = tax.Set{
		{Category: tax.CategoryVAT, Key: tax.KeyReverseCharge},
	}
	require.NoError(t, inv.Calculate())
	combo := inv.Lines[0].Taxes[0]
	assert.Equal(t, "AE", combo.Ext.Get(untdid.ExtKeyTaxCategory).String())
}

func TestNormalizeGeneratesSIRENFromFrenchTaxID(t *testing.T) {
	inv := testInvoiceB2B(t)
	inv.Supplier.Identities = nil
	require.NoError(t, inv.Calculate())
	// normalizeParty should have injected a SIREN-scheme identity.
	found := false
	for _, id := range inv.Supplier.Identities {
		if id.Ext.Get(iso.ExtKeySchemeID).String() == "0002" {
			found = true
			assert.Equal(t, "356000000", id.Code.String())
		}
	}
	assert.True(t, found, "expected SIREN identity to be generated from TaxID")
}

// --- Internal helper coverage (bill_invoice.go) -------------------------

func TestExtensionsValueNilPointer(t *testing.T) {
	assert.True(t, extensionsValue((*tax.Extensions)(nil)).IsZero())
}

func TestExtensionsValueUnknownType(t *testing.T) {
	assert.True(t, extensionsValue(42).IsZero())
}

func TestExtensionsValueValue(t *testing.T) {
	e := tax.ExtensionsOf(tax.ExtMap{"k": "v"})
	assert.False(t, extensionsValue(e).IsZero())
}

func TestPartyHasSIRENWrongType(t *testing.T) {
	assert.False(t, partyHasSIREN("x"))
}

func TestPartyHasAllowedLegalSchemeWrongType(t *testing.T) {
	assert.False(t, partyHasAllowedLegalScheme("x"))
}

func TestPartyHasTaxIDWhenRequiredWrongType(t *testing.T) {
	assert.True(t, partyHasTaxIDWhenRequired("x"))
}

func TestPartyHasTaxIDWhenRequiredNonRequiredScheme(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{{
		Code: "X",
		Ext:  tax.ExtensionsOf(tax.ExtMap{iso.ExtKeySchemeID: "0227"}),
	}}}
	assert.True(t, partyHasTaxIDWhenRequired(p))
}

func TestInvoiceIsB2BWrongType(t *testing.T) {
	assert.False(t, invoiceIsB2BAny("x"))
}

func TestInvoiceIsB2CWrongType(t *testing.T) {
	assert.False(t, invoiceIsB2CAny("x"))
}

func TestInvoiceDocumentTypeAllowedEmpty(t *testing.T) {
	assert.False(t, invoiceDocumentTypeAllowed(tax.Extensions{}))
}

func TestExtensionsHaveBillingModeMissing(t *testing.T) {
	assert.False(t, extensionsHaveBillingMode(tax.Extensions{}))
}

func TestExtensionsHaveB2CCategoryMissing(t *testing.T) {
	assert.False(t, extensionsHaveB2CCategory(tax.Extensions{}))
}

func TestInvoiceIsFinalAfterAdvanceWrongType(t *testing.T) {
	assert.False(t, invoiceIsFinalAfterAdvance("x"))
}

func TestInvoiceIsFinalAfterAdvanceNoExt(t *testing.T) {
	assert.False(t, invoiceIsFinalAfterAdvance(&bill.Invoice{Tax: &bill.Tax{}}))
}

func TestInvoiceNotAdvancePaymentDocTypeWrongType(t *testing.T) {
	assert.True(t, invoiceNotAdvancePaymentDocType(42))
}

func TestInvoiceHasSellerVATIDForExemptWrongType(t *testing.T) {
	assert.False(t, invoiceHasSellerVATIDForExempt("x"))
}

func TestInvoiceHasExemptComboWrongType(t *testing.T) {
	assert.False(t, invoiceHasExemptCombo("x"))
}

func TestInvoiceHasExemptTaxNoteWrongType(t *testing.T) {
	assert.False(t, invoiceHasExemptTaxNote("x"))
}

func TestInvoiceVATRatesAllowedWrongType(t *testing.T) {
	assert.True(t, invoiceVATRatesAllowed("x"))
}

func TestMustParsePercentagesPanicsOnBadInput(t *testing.T) {
	assert.Panics(t, func() { mustParsePercentages("not-a-percentage") })
}

func TestPercentageInListEmpty(t *testing.T) {
	p := num.MakePercentage(20, 2)
	assert.False(t, percentageInList(p, nil))
}

func TestNormalizeInvoiceNilSafe(t *testing.T) {
	assert.NotPanics(t, func() { normalizeInvoice(nil) })
}

func TestNormalizeInvoiceBillingModeDefaultsM2WhenPaid(t *testing.T) {
	due := num.MakeAmount(0, 2)
	inv := &bill.Invoice{
		Totals: &bill.Totals{Due: &due},
		Tax:    &bill.Tax{},
	}
	normalizeInvoiceBillingMode(inv)
	assert.Equal(t, BillingModeM2, inv.Tax.Ext.Get(ExtKeyBillingMode))
}
