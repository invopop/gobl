package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/co"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceCorrect(t *testing.T) {
	i := testInvoiceESForCorrection(t)

	i2, err := i.Correct(bill.Refund, bill.WithReason("test refund"))
	require.NoError(t, err)
	assert.Equal(t, bill.InvoiceTypeCorrective, i2.Type)
	assert.Equal(t, i2.Lines[0].Quantity.String(), "-10")
	pre := i2.Preceding[0]
	assert.Equal(t, pre.Series, i.Series)
	assert.Equal(t, pre.Code, i.Code)
	assert.Equal(t, pre.Reason, "test refund")

	err = i2.Calculate()
	require.NoError(t, err)

	_, err = i.Correct(bill.Append, bill.WithReason("should fail"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "append not supported")

	i2, err = i.Correct()
	require.NoError(t, err)
	assert.Equal(t, i2.Type, bill.InvoiceTypeCorrective)

	i = testInvoiceCOForCorrection(t)
	_, err = i.Correct(bill.Refund)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing stamp")

	stamps := []*cbc.Stamp{
		{
			Provider: co.StampProviderDIANCUDE,
			Value:    "FOOO",
		},
		{
			Provider: co.StampProviderDIANQR,
			Value:    "BARRRR",
		},
	}

	_, err = i.Correct(bill.WithStamps(stamps))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "correction not supported by regime")

	i2, err = i.Correct(bill.Refund, bill.WithStamps(stamps), bill.WithCorrectionMethod(co.CorrectionMethodKeyRevoked))
	require.NoError(t, err)
	assert.Equal(t, i2.Type, bill.InvoiceTypeCreditNote)
	assert.Len(t, i2.Preceding[0].Stamps, 2)
	assert.Equal(t, i2.Preceding[0].CorrectionMethod, co.CorrectionMethodKeyRevoked)
}

func testInvoiceESForCorrection(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Series: "TEST",
		Code:   "123",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.ES,
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
		},
	}
	return i
}

func testInvoiceCOForCorrection(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Series: "TEST",
		Code:   "123",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.CO,
				Code:    "9014586527",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.CO,
				Code:    "8001345363",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
		},
	}
	return i
}
