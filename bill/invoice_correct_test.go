package bill_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/co"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceCorrect(t *testing.T) {
	// Spanish Case (only corrective)

	// debit note not supported in Spain
	i := testInvoiceESForCorrection(t)
	err := i.Correct(bill.Debit)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid correction type: debit-note")

	i = testInvoiceESForCorrection(t)
	err = i.Correct(bill.Credit,
		bill.WithReason("test refund"),
		bill.WithChanges(es.CorrectionKeyLine),
	)
	require.NoError(t, err)
	assert.Equal(t, bill.InvoiceTypeCreditNote, i.Type)
	assert.Equal(t, i.Lines[0].Quantity.String(), "10")
	assert.Equal(t, i.IssueDate, cal.Today())
	pre := i.Preceding[0]
	assert.Equal(t, pre.Series, "TEST")
	assert.Equal(t, pre.Code, "123")
	assert.Equal(t, pre.IssueDate, cal.NewDate(2022, 6, 13))
	assert.Equal(t, pre.Reason, "test refund")
	assert.Equal(t, i.Totals.Payable.String(), "900.00")

	// can't run twice
	err = i.Correct(bill.Corrective)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot correct an invoice without a code")

	i = testInvoiceESForCorrection(t)
	err = i.Correct()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing correction type")

	i = testInvoiceESForCorrection(t)
	err = i.Correct(bill.Debit, bill.WithReason("should fail"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid correction type: debit-note")

	i = testInvoiceESForCorrection(t)
	err = i.Correct(
		bill.Corrective,
		bill.WithChanges(es.CorrectionKeyLine),
	)
	require.NoError(t, err)
	assert.Equal(t, i.Type, bill.InvoiceTypeCorrective)

	// With preset date
	i = testInvoiceESForCorrection(t)
	d := cal.MakeDate(2023, 6, 13)
	err = i.Correct(
		bill.Credit,
		bill.WithIssueDate(d),
		bill.WithChanges(es.CorrectionKeyLine),
	)
	require.NoError(t, err)
	assert.Equal(t, i.IssueDate, d)

	// France case (both corrective and credit note)
	i = testInvoiceFRForCorrection(t)
	err = i.Correct(bill.Corrective)
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
	err = i.Correct(bill.Corrective, bill.WithStamps(stamps))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid correction type: corrective")

	i = testInvoiceCOForCorrection(t)
	err = i.Correct(
		bill.Credit,
		bill.WithStamps(stamps),
		bill.WithReason("test refund"),
		bill.WithChanges(co.CorrectionKeyRevoked),
	)
	require.NoError(t, err)
	assert.Equal(t, i.Type, bill.InvoiceTypeCreditNote)
	pre = i.Preceding[0]
	require.Len(t, pre.Stamps, 1)
	assert.Equal(t, pre.Stamps[0].Provider, co.StampProviderDIANCUDE)
	// assert.Equal(t, pre.CorrectionMethod, co.CorrectionMethodKeyRevoked)
}

func TestCorrectWithOptions(t *testing.T) {
	i := testInvoiceESForCorrection(t)
	opts := &bill.CorrectionOptions{
		Type:    bill.InvoiceTypeCreditNote,
		Reason:  "test refund",
		Changes: []cbc.Key{es.CorrectionKeyLine},
	}
	err := i.Correct(bill.WithOptions(opts))
	require.NoError(t, err)
	assert.Equal(t, bill.InvoiceTypeCreditNote, i.Type)
	assert.Equal(t, i.Lines[0].Quantity.String(), "10")
	assert.Equal(t, i.IssueDate, cal.Today())
	pre := i.Preceding[0]
	assert.Equal(t, pre.Series, "TEST")
	assert.Equal(t, pre.Code, "123")
	assert.Equal(t, pre.IssueDate, cal.NewDate(2022, 6, 13))
	assert.Equal(t, pre.Reason, "test refund")
	assert.Equal(t, i.Totals.Payable.String(), "900.00")
}

func TestCorrectionOptionsSchema(t *testing.T) {
	inv := testInvoiceESForCorrection(t)
	out, err := inv.CorrectionOptionsSchema()
	require.NoError(t, err)

	schema, ok := out.(*jsonschema.Schema)
	require.True(t, ok)

	cos := schema.Definitions["CorrectionOptions"]
	assert.Equal(t, cos.Properties.Len(), 5)

	pm, ok := cos.Properties.Get("changes")
	require.True(t, ok)
	assert.Len(t, pm.Items.OneOf, 22)

	exp := `{"items":{"$ref":"https://gobl.org/draft-0/cbc/key","oneOf":[{"const":"code","title":"Invoice code"},{"const":"series","title":"Invoice series"},{"const":"issue-date","title":"Issue date"},{"const":"supplier-name","title":"Name and surnames/Corporate name - Issuer (Sender)"},{"const":"customer-name","title":"Name and surnames/Corporate name - Receiver"},{"const":"supplier-tax-id","title":"Issuer's Tax Identification Number"},{"const":"customer-tax-id","title":"Receiver's Tax Identification Number"},{"const":"supplier-addr","title":"Issuer's address"},{"const":"customer-addr","title":"Receiver's address"},{"const":"line","title":"Item line"},{"const":"tax-rate","title":"Applicable Tax Rate"},{"const":"tax-amount","title":"Applicable Tax Amount"},{"const":"period","title":"Applicable Date/Period"},{"const":"type","title":"Invoice Class"},{"const":"legal-details","title":"Legal literals"},{"const":"tax-base","title":"Taxable Base"},{"const":"tax","title":"Calculation of tax outputs"},{"const":"tax-retained","title":"Calculation of tax inputs"},{"const":"refund","title":"Taxable Base modified due to return of packages and packaging materials"},{"const":"discount","title":"Taxable Base modified due to discounts and rebates"},{"const":"judicial","title":"Taxable Base modified due to firm court ruling or administrative decision"},{"const":"insolvency","title":"Taxable Base modified due to unpaid outputs where there is a judgement opening insolvency proceedings"}]},"type":"array","title":"Changes","description":"Changes keys that describe the specific changes according to the tax regime."}`
	data, err := json.Marshal(pm)
	require.NoError(t, err)
	if !assert.JSONEq(t, exp, string(data)) {
		t.Log(string(data))
	}

	data, err = json.Marshal(schema)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"$id":"https://gobl.org/draft-0/bill/correction-options?tax_regime=es"`)
}

func TestCorrectWithData(t *testing.T) {
	i := testInvoiceESForCorrection(t)
	data := []byte(`{"type":"credit-note","reason":"test refund"}`)

	err := i.Correct(
		bill.WithData(data),
		bill.WithChanges(es.CorrectionKeyLine),
	)
	assert.NoError(t, err)
	assert.Equal(t, i.Type, bill.InvoiceTypeCreditNote)
	assert.Equal(t, i.Lines[0].Quantity.String(), "10") // implies credit was made

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
			PricesInclude: tax.CategoryVAT,
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
			PricesInclude: tax.CategoryVAT,
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
