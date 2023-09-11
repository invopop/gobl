package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/co"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceCorrect(t *testing.T) {
	// Spanish Case (only corrective)
	i := testInvoiceESForCorrection(t)
	err := i.Correct(bill.Credit, bill.Debit)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot use both credit and debit options")

	i = testInvoiceESForCorrection(t)
	err = i.Correct(bill.Credit,
		bill.WithReason("test refund"),
		bill.WithCorrectionMethod(es.CorrectionMethodKeyComplete),
		bill.WithCorrection(es.CorrectionKeyLine),
	)
	require.NoError(t, err)
	assert.Equal(t, bill.InvoiceTypeCorrective, i.Type)
	assert.Equal(t, i.Lines[0].Quantity.String(), "-10")
	assert.Equal(t, i.IssueDate, cal.Today())
	pre := i.Preceding[0]
	assert.Equal(t, pre.Series, "TEST")
	assert.Equal(t, pre.Code, "123")
	assert.Equal(t, pre.IssueDate, cal.NewDate(2022, 6, 13))
	assert.Equal(t, pre.Reason, "test refund")
	assert.Equal(t, i.Totals.Payable.String(), "-900.00")

	// can't run twice
	err = i.Correct()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot correct an invoice without a code")

	i = testInvoiceESForCorrection(t)
	err = i.Correct(bill.Debit, bill.WithReason("should fail"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "debit note not supported by regime")

	i = testInvoiceESForCorrection(t)
	err = i.Correct(
		bill.WithCorrection(es.CorrectionKeyLine),
		bill.WithCorrectionMethod(es.CorrectionMethodKeyComplete),
	)
	require.NoError(t, err)
	assert.Equal(t, i.Type, bill.InvoiceTypeCorrective)

	// With preset date
	i = testInvoiceESForCorrection(t)
	d := cal.MakeDate(2023, 6, 13)
	err = i.Correct(
		bill.WithIssueDate(d),
		bill.WithCorrection(es.CorrectionKeyLine),
		bill.WithCorrectionMethod(es.CorrectionMethodKeyComplete),
	)
	require.NoError(t, err)
	assert.Equal(t, i.IssueDate, d)

	// France case (both corrective and credit note)
	i = testInvoiceFRForCorrection(t)
	err = i.Correct()
	require.NoError(t, err)
	assert.Equal(t, i.Type, bill.InvoiceTypeCorrective)

	i = testInvoiceFRForCorrection(t)
	err = i.Correct(bill.Credit)
	require.NoError(t, err)
	assert.Equal(t, i.Type, bill.InvoiceTypeCreditNote)

	// Colombia case (only credit note)

	i = testInvoiceCOForCorrection(t)
	err = i.Correct(bill.Credit)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing stamp")

	stamps := []*head.Stamp{
		{
			Provider: co.StampProviderDIANCUDE,
			Value:    "FOOO",
		},
		{
			Provider: co.StampProviderDIANQR, // not copied!
			Value:    "BARRRR",
		},
	}

	i = testInvoiceCOForCorrection(t)
	err = i.Correct(bill.WithStamps(stamps))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "corrective invoice type not supported by regime, try credit or debit")

	i = testInvoiceCOForCorrection(t)
	err = i.Correct(
		bill.Credit,
		bill.WithStamps(stamps),
		bill.WithCorrectionMethod(co.CorrectionMethodKeyRevoked),
		bill.WithReason("test refund"),
	)
	require.NoError(t, err)
	assert.Equal(t, i.Type, bill.InvoiceTypeCreditNote)
	pre = i.Preceding[0]
	require.Len(t, pre.Stamps, 1)
	assert.Equal(t, pre.Stamps[0].Provider, co.StampProviderDIANCUDE)
	assert.Equal(t, pre.CorrectionMethod, co.CorrectionMethodKeyRevoked)
}

func TestCorrectWithOptions(t *testing.T) {
	i := testInvoiceESForCorrection(t)
	opts := &bill.CorrectionOptions{
		Credit: true,
		Reason: "test refund",
	}
	err := i.Correct(bill.WithOptions(opts))
	require.NoError(t, err)
	assert.Equal(t, bill.InvoiceTypeCorrective, i.Type)
	assert.Equal(t, i.Lines[0].Quantity.String(), "-10")
	assert.Equal(t, i.IssueDate, cal.Today())
	pre := i.Preceding[0]
	assert.Equal(t, pre.Series, "TEST")
	assert.Equal(t, pre.Code, "123")
	assert.Equal(t, pre.IssueDate, cal.NewDate(2022, 6, 13))
	assert.Equal(t, pre.Reason, "test refund")
	assert.Equal(t, i.Totals.Payable.String(), "-900.00")
}

func TestCorrectWithData(t *testing.T) {
	i := testInvoiceESForCorrection(t)
	data := []byte(`{"credit":true,"reason":"test refund"}`)

	err := i.Correct(bill.WithData(data))
	assert.NoError(t, err)
	assert.Equal(t, i.Lines[0].Quantity.String(), "-10") // implies credit was made

	data = []byte(`{"credit": true`) // invalid json
	err = i.Correct(bill.WithData(data))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected end of JSON input")
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

func testInvoiceFRForCorrection(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Series: "TEST",
		Code:   "123",
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.FR,
				Code:    "732829320",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.FR,
				Code:    "391838042",
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
