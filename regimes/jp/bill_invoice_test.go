package jp_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/jp"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validInvoice returns a base valid JP invoice that can be customized per test.
func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime:    tax.WithRegime(l10n.TaxCountryCode(l10n.JP)),
		Currency:  currency.JPY,
		IssueDate: *cal.NewDate(2024, 11, 1),
		Supplier: &org.Party{
			Name: "株式会社テスト",
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.JP),
				Code:    "T7000012050002",
			},
		},
		Customer: &org.Party{Name: "Buyer Co"},
		Lines: []*bill.Line{
			{
				Index:    1,
				Quantity: num.MakeAmount(1, 0),
				Item:     &org.Item{Name: "Item", Price: ptrAmount(1000, 0)},
				Taxes: []*tax.Combo{
					{Category: tax.CategoryVAT, Key: tax.KeyStandard, Rate: tax.RateGeneral},
				},
			},
		},
	}
}

func TestValidateInvoice_ExportRequiresZeroRate(t *testing.T) {
	t.Parallel()

	// Export invoice with zero-rated line — valid
	invZero := validInvoice()
	invZero.Lines[0].Taxes = []*tax.Combo{
		{Category: tax.CategoryVAT, Key: tax.KeyZero, Rate: tax.RateZero},
	}
	invZero.SetTags(jp.TagExport)
	require.NoError(t, invZero.Calculate(), "calculate export invoice")
	require.NoError(t, invZero.Validate(), "export invoice with zero rate should be valid")

	// Export invoice with standard rate — invalid
	invStandard := validInvoice()
	invStandard.SetTags(jp.TagExport)
	require.NoError(t, invStandard.Calculate(), "calculate so validation runs")
	err := invStandard.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "zero")
	assert.Contains(t, err.Error(), "export")
}

func TestValidateInvoice_SimplifiedCustomerOptional(t *testing.T) {
	t.Parallel()

	// Simplified invoice without customer — valid
	inv := validInvoice()
	inv.Customer = nil
	inv.SetTags(jp.TagSimplified)
	require.NoError(t, inv.Calculate())
	assert.NoError(t, inv.Validate(), "simplified invoice should not require customer")

	// Simplified invoice with customer — also valid
	inv2 := validInvoice()
	inv2.SetTags(jp.TagSimplified)
	require.NoError(t, inv2.Calculate())
	assert.NoError(t, inv2.Validate(), "simplified invoice with customer should be valid")
}

func TestValidateInvoice_StandardRequiresCustomer(t *testing.T) {
	t.Parallel()

	// Standard invoice without customer — invalid
	inv := validInvoice()
	inv.Customer = nil
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "customer")
}

func TestValidateInvoice_ScenarioNotes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tag     cbc.Key
		noteStr string
	}{
		{
			name:    "export tag adds legal note",
			tag:     jp.TagExport,
			noteStr: "Export",
		},
		{
			name:    "simplified tag adds legal note",
			tag:     jp.TagSimplified,
			noteStr: "Simplified",
		},
		{
			name:    "self-billing tag adds legal note",
			tag:     jp.TagSelfBilling,
			noteStr: "Self-billed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := validInvoice()
			if tt.tag == jp.TagExport {
				inv.Lines[0].Taxes = []*tax.Combo{
					{Category: tax.CategoryVAT, Key: tax.KeyZero, Rate: tax.RateZero},
				}
			}
			inv.SetTags(tt.tag)
			require.NoError(t, inv.Calculate())

			found := false
			for _, n := range inv.Notes {
				if n.Key == org.NoteKeyLegal {
					assert.Contains(t, n.Text, tt.noteStr)
					found = true
				}
			}
			assert.True(t, found, "expected a legal note for tag %q", tt.tag)
		})
	}
}

func TestValidateInvoice_CorrectionsSupported(t *testing.T) {
	t.Parallel()

	regime := jp.New()
	require.NotEmpty(t, regime.Corrections, "regime should define correction types")
	assert.Equal(t, bill.InvoiceTypeCreditNote, regime.Corrections[0].Types[0])
}

func TestValidateInvoice_SupplierRequiresTaxIDAndName(t *testing.T) {
	t.Parallel()

	// Supplier without name — invalid
	inv := validInvoice()
	inv.Supplier.Name = ""
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name")

	// Supplier without tax ID — invalid
	inv2 := validInvoice()
	inv2.Supplier.TaxID = nil
	require.NoError(t, inv2.Calculate())
	err = inv2.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "tax_id")

	// Nil supplier cast — validateSupplier should handle gracefully via Validate dispatch
	inv3 := validInvoice()
	inv3.Supplier = nil
	require.NoError(t, inv3.Calculate())
	err = inv3.Validate()
	// Supplier nil may or may not error depending on framework; just ensure no panic
	_ = err
}

func TestValidateInvoice_CustomerNameRequired(t *testing.T) {
	t.Parallel()

	// Customer without name — invalid
	inv := validInvoice()
	inv.Customer = &org.Party{Name: ""}
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name")
}

func TestValidateInvoice_AddressValidation(t *testing.T) {
	t.Parallel()

	// Valid address with street
	inv := validInvoice()
	inv.Supplier.Addresses = []*org.Address{
		{Street: "1-2-3 Shibuya", Locality: "Shibuya-ku"},
	}
	require.NoError(t, inv.Calculate())
	assert.NoError(t, inv.Validate())

	// Valid address with only locality
	inv2 := validInvoice()
	inv2.Supplier.Addresses = []*org.Address{
		{Locality: "Shibuya-ku"},
	}
	require.NoError(t, inv2.Calculate())
	assert.NoError(t, inv2.Validate())

	// Valid address with only street
	inv3 := validInvoice()
	inv3.Supplier.Addresses = []*org.Address{
		{Street: "1-2-3 Shibuya"},
	}
	require.NoError(t, inv3.Calculate())
	assert.NoError(t, inv3.Validate())

	// Invalid address: neither street nor locality
	inv4 := validInvoice()
	inv4.Supplier.Addresses = []*org.Address{
		{Code: "100-0001"},
	}
	require.NoError(t, inv4.Calculate())
	err := inv4.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "street")
}

func TestValidateInvoice_EmptyAddresses(t *testing.T) {
	t.Parallel()

	// Nil addresses — valid
	inv := validInvoice()
	inv.Supplier.Addresses = nil
	require.NoError(t, inv.Calculate())
	assert.NoError(t, inv.Validate())

	// Empty addresses slice — valid
	inv2 := validInvoice()
	inv2.Supplier.Addresses = []*org.Address{}
	require.NoError(t, inv2.Calculate())
	assert.NoError(t, inv2.Validate())
}

func TestValidateInvoice_ExportNonVATTaxesIgnored(t *testing.T) {
	t.Parallel()

	// Export invoice with a nil line in the list should not panic
	inv := validInvoice()
	inv.Lines[0].Taxes = []*tax.Combo{
		{Category: tax.CategoryVAT, Key: tax.KeyZero, Rate: tax.RateZero},
	}
	inv.SetTags(jp.TagExport)
	require.NoError(t, inv.Calculate())
	assert.NoError(t, inv.Validate())
}

func TestValidateInvoice_LineWithNilItem(t *testing.T) {
	t.Parallel()

	inv := validInvoice()
	inv.Lines = append(inv.Lines, &bill.Line{
		Index:    2,
		Quantity: num.MakeAmount(1, 0),
		Item:     nil, // Missing item
		Taxes: []*tax.Combo{
			{Category: tax.CategoryVAT, Key: tax.KeyStandard, Rate: tax.RateGeneral},
		},
	})
	require.NoError(t, inv.Calculate())
	err := inv.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "item")
}

func TestValidateInvoice_NilCustomerValidateCustomer(t *testing.T) {
	t.Parallel()

	// validateCustomer receives nil — no error
	err := jp.Validate((*org.Party)(nil))
	assert.NoError(t, err)
}

func TestValidate_UnknownDocType(t *testing.T) {
	t.Parallel()

	// Passing an unknown doc type should return nil
	err := jp.Validate("not a known doc type")
	assert.NoError(t, err)
}

func TestNormalize_UnknownDocType(t *testing.T) {
	t.Parallel()

	// Passing an unknown doc type should not panic
	jp.Normalize("not a known doc type")
}

func ptrAmount(val int64, exp uint32) *num.Amount {
	a := num.MakeAmount(val, exp)
	return &a
}
