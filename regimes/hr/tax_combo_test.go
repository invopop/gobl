package hr_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/hr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// passThroughLine is a line for a pass-through item worth 5.00.
func passThroughLine(code cbc.Code) *bill.Line {
	return &bill.Line{
		Quantity: num.MakeAmount(1, 0),
		Item: &org.Item{
			Name:  "Povratna naknada",
			Price: num.NewAmount(500, 2),
		},
		Taxes: tax.Set{
			{
				Category: tax.CategoryVAT,
				Key:      tax.KeyExempt,
				Ext:      tax.ExtensionsOf(cbc.CodeMap{hr.ExtKeyPassThrough: code}),
			},
		},
	}
}

func TestPassThroughExtension(t *testing.T) {
	t.Run("defined by the regime", func(t *testing.T) {
		ed := tax.ExtensionForKey(hr.ExtKeyPassThrough)
		require.NotNil(t, ed)
		assert.Equal(t,
			[]cbc.Code{hr.PassThroughReturnFee, hr.PassThroughMotorVehicleTax, hr.PassThroughNonTaxableFee},
			cbc.DefinitionCodes(ed.Values),
		)
	})

	t.Run("line kept out of the VAT base", func(t *testing.T) {
		inv := testInvoice(t)
		inv.Lines = append(inv.Lines, passThroughLine(hr.PassThroughReturnFee))
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))

		vat := inv.Totals.Taxes.Category(tax.CategoryVAT)
		require.NotNil(t, vat)
		require.Len(t, vat.Rates, 2)
		assert.Equal(t, "100.00", vat.Rates[0].Base.String())
		pt := vat.Rates[1]
		assert.Equal(t, tax.KeyExempt, pt.Key)
		assert.Equal(t, hr.PassThroughReturnFee, pt.Ext.Get(hr.ExtKeyPassThrough))
		assert.Equal(t, "5.00", pt.Base.String())
		assert.Equal(t, "0.00", pt.Amount.String())
		assert.Equal(t, "25.00", vat.Amount.String())
		assert.Equal(t, "130.00", inv.Totals.Payable.String())
	})

	t.Run("document charge", func(t *testing.T) {
		inv := testInvoice(t)
		inv.Charges = []*bill.Charge{
			{
				Reason: "Turistička pristojba",
				Amount: num.MakeAmount(1000, 2),
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Key:      tax.KeyExempt,
						Ext:      tax.ExtensionsOf(cbc.CodeMap{hr.ExtKeyPassThrough: hr.PassThroughNonTaxableFee}),
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
		assert.Equal(t, "135.00", inv.Totals.Payable.String())
	})

	t.Run("unknown code", func(t *testing.T) {
		inv := testInvoice(t)
		inv.Lines = append(inv.Lines, passThroughLine("XYZ"))
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-HR-TAX-COMBO-01]")
	})

	t.Run("on another tax category", func(t *testing.T) {
		inv := testInvoice(t)
		inv.Lines[0].Taxes = append(inv.Lines[0].Taxes, &tax.Combo{
			Category: hr.TaxCategoryPP,
			Percent:  num.NewPercentage(3, 2),
			Ext:      tax.ExtensionsOf(cbc.CodeMap{hr.ExtKeyPassThrough: hr.PassThroughNonTaxableFee}),
		})
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-HR-TAX-COMBO-02]")
	})

	t.Run("with VAT charged", func(t *testing.T) {
		inv := testInvoice(t)
		line := passThroughLine(hr.PassThroughMotorVehicleTax)
		line.Taxes[0].Key = tax.KeyStandard
		line.Taxes[0].Rate = tax.RateGeneral
		inv.Lines = append(inv.Lines, line)
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "[GOBL-HR-TAX-COMBO-02]")
	})
}
