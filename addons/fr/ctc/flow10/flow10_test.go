package flow10

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/catalogues/dgfip"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// runNormalize invokes the addon's registered normalizer on the given
// object, matching what tax.Normalize would do during Calculate.
func runNormalize(t *testing.T, doc any) {
	t.Helper()
	tax.Normalize([]tax.Normalizer{tax.AddonForKey(V1).Normalizer}, doc)
}

// frPartyWithSIREN returns a French supplier party with a SIREN
// identity carrying the iso-scheme-id extension.
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
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					iso.ExtKeySchemeID: identitySchemeIDSIREN,
				}),
			},
		},
		Addresses: []*org.Address{{Country: "FR"}},
	}
}

// deCustomerWithVATID returns a German customer party with an EU-VAT
// identity (ICD scheme 0223).
func deCustomerWithVATID() *org.Party {
	return &org.Party{
		Name: "Kunde Deutschland GmbH",
		TaxID: &tax.Identity{
			Country: "DE",
			Code:    "111111125",
		},
		Identities: []*org.Identity{
			{
				Code:  "DE111111125",
				Scope: org.IdentityScopeLegal,
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					iso.ExtKeySchemeID: identitySchemeIDEUVAT,
				}),
			},
		},
		Addresses: []*org.Address{{Country: "DE"}},
	}
}

func testInvoiceB2BCrossBorder(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(V1),
		Code:      "INV-2026-001",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2026, 1, 15),
		Type:      bill.InvoiceTypeStandard,
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				dgfip.ExtKeyBillingMode: dgfip.BillingModeS1,
			}),
		},
		Supplier: frPartyWithSIREN(),
		Customer: deCustomerWithVATID(),
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
	return &bill.Invoice{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(V1),
		Code:      "INV-2026-B2C-001",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2026, 1, 15),
		Type:      bill.InvoiceTypeStandard,
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
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
}

func TestInvoiceB2BCrossBorderHappyPath(t *testing.T) {
	inv := testInvoiceB2BCrossBorder(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestInvoiceB2CHappyPath(t *testing.T) {
	inv := testInvoiceB2C(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestInvoiceB2CDefaultsCategoryToTNT1(t *testing.T) {
	inv := testInvoiceB2C(t)
	inv.Tax.Ext = inv.Tax.Ext.Delete(ExtKeyB2CCategory)
	runNormalize(t, inv)
	assert.Equal(t, B2CCategoryNotTaxable, inv.Tax.Ext.Get(ExtKeyB2CCategory))
}

func TestInvoiceB2CMissingSupplierSIRENFails(t *testing.T) {
	inv := testInvoiceB2C(t)
	inv.Supplier.Identities = nil
	inv.Supplier.TaxID = nil
	assert.Error(t, rules.Validate(inv))
}

func TestInvoiceB2BMissingBillingModeFails(t *testing.T) {
	inv := testInvoiceB2BCrossBorder(t)
	inv.Tax.Ext = inv.Tax.Ext.Delete(dgfip.ExtKeyBillingMode)
	assert.Error(t, rules.Validate(inv))
}

func TestInvoiceB2CVATRateRejectedOutsideWhitelist(t *testing.T) {
	inv := testInvoiceB2C(t)
	inv.Lines[0].Taxes = tax.Set{
		{Category: tax.CategoryVAT, Percent: num.NewPercentage(35, 2)},
	}
	assert.Error(t, rules.Validate(inv))
}

func TestInvoiceB2BMissingSupplierAddressCountryFails(t *testing.T) {
	inv := testInvoiceB2BCrossBorder(t)
	inv.Supplier.Addresses[0].Country = ""
	assert.Error(t, rules.Validate(inv))
}

func testPaymentReceipt(t *testing.T) *bill.Payment {
	t.Helper()
	value := cal.MakeDate(2026, 1, 15)
	return &bill.Payment{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(V1),
		Code:      "PMT-2026-001",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2026, 1, 15),
		Type:      bill.PaymentTypeReceipt,
		ValueDate: &value,
		Supplier:  frPartyWithSIREN(),
		Methods:   []*pay.Record{{Key: pay.MeansKeyCreditTransfer}},
		Lines: []*bill.PaymentLine{
			{
				Amount: num.MakeAmount(12000, 2),
			},
		},
	}
}

func TestPaymentReceiptHappy(t *testing.T) {
	pmt := testPaymentReceipt(t)
	require.NoError(t, rules.Validate(pmt))
}

func TestPaymentRejectsNonReceiptType(t *testing.T) {
	pmt := testPaymentReceipt(t)
	pmt.Type = bill.PaymentTypeRequest
	assert.Error(t, rules.Validate(pmt))
}
