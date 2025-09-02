package bill_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/addons/co/dian"
	"github.com/invopop/gobl/addons/es/facturae"
	"github.com/invopop/gobl/addons/es/tbai"
	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceCorrect(t *testing.T) {
	i := testInvoicePTForCorrection(t)
	err := i.Correct(bill.Corrective)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid correction type: corrective")

	i = testInvoiceESForCorrection(t)
	require.NoError(t, i.Calculate())
	err = i.Correct(bill.Credit,
		bill.WithReason("test refund"),
		bill.WithExtension(facturae.ExtKeyCorrection, "01"),
	)
	require.NoError(t, err)
	assert.Equal(t, bill.InvoiceTypeCreditNote, i.Type)
	assert.Equal(t, i.Lines[0].Quantity.String(), "10")
	assert.Equal(t, i.IssueDate, cal.Today())
	assert.Equal(t, i.Series, cbc.Code("TEST"))
	assert.Empty(t, i.Code)
	pre := i.Preceding[0]
	assert.Equal(t, pre.Series.String(), "TEST")
	assert.Equal(t, pre.Code.String(), "123")
	assert.Equal(t, pre.IssueDate, cal.NewDate(2022, 6, 13))
	assert.Equal(t, pre.Reason, "test refund")
	assert.Equal(t, i.Totals.Payable.String(), "900.00")
	assert.Nil(t, pre.Tax, "don't copy tax by default")

	// can't run twice
	err = i.Correct(bill.Corrective)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot correct an invoice without a code")

	i = testInvoiceESForCorrection(t)
	err = i.Correct()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing correction type")

	i = testInvoicePTForCorrection(t)
	err = i.Correct(bill.Corrective, bill.WithReason("should fail"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid correction type: corrective")

	i = testInvoiceESForCorrection(t)
	err = i.Correct(
		bill.Corrective,
		bill.WithExtension(facturae.ExtKeyCorrection, "01"),
	)
	require.NoError(t, err)
	assert.Equal(t, i.Type, bill.InvoiceTypeCorrective)

	// With preset date
	i = testInvoiceESForCorrection(t)
	d := cal.MakeDate(2023, 6, 13)
	err = i.Correct(
		bill.Credit,
		bill.WithIssueDate(d),
		bill.WithExtension(facturae.ExtKeyCorrection, "01"),
	)
	require.NoError(t, err)
	assert.Equal(t, i.IssueDate, d)

	t.Run("with series", func(t *testing.T) {
		inv := testInvoiceESForCorrection(t)
		err := inv.Correct(bill.Credit, bill.WithSeries("R-TEST"))
		require.NoError(t, err)
		assert.Equal(t, inv.Series, cbc.Code("R-TEST"))
		assert.Equal(t, inv.Preceding[0].Series.String(), "TEST")
	})

	t.Run("with taxes", func(t *testing.T) {
		inv := testInvoiceESForCorrection(t)
		require.NoError(t, inv.Calculate())
		err := inv.Correct(bill.Credit, bill.WithSeries("R-TEST"), bill.WithCopyTax())
		require.NoError(t, err)
		require.NotNil(t, inv.Preceding[0].Tax)
		assert.Equal(t, "156.20", inv.Preceding[0].Tax.Sum.String())
	})

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
			Provider: dian.StampCUDE,
			Value:    "FOOO",
		},
		{
			Provider: dian.StampQR, // not copied!
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
		bill.WithExtension(dian.ExtKeyCreditCode, "2"),
	)
	require.NoError(t, err)
	assert.Equal(t, i.Type, bill.InvoiceTypeCreditNote)
	pre = i.Preceding[0]
	require.Len(t, pre.Stamps, 1)
	assert.Equal(t, pre.Stamps[0].Provider, dian.StampCUDE)
	// assert.Equal(t, pre.CorrectionMethod, co.CorrectionMethodKeyRevoked)
}

func TestCorrectWithOptions(t *testing.T) {
	i := testInvoiceESForCorrection(t)
	opts := &bill.CorrectionOptions{
		Type:   bill.InvoiceTypeCreditNote,
		Reason: "test refund",
		Series: "R-TEST",
		Ext: tax.Extensions{
			facturae.ExtKeyCorrection: "01",
		},
	}
	err := i.Correct(bill.WithOptions(opts))
	require.NoError(t, err)
	assert.Equal(t, bill.InvoiceTypeCreditNote, i.Type)
	assert.Equal(t, i.Lines[0].Quantity.String(), "10")
	assert.Equal(t, i.IssueDate, cal.Today())
	assert.Equal(t, i.Series.String(), "R-TEST")
	assert.Empty(t, i.Code)
	pre := i.Preceding[0]
	assert.Equal(t, pre.Series.String(), "TEST")
	assert.Equal(t, pre.Code.String(), "123")
	assert.Equal(t, pre.IssueDate, cal.NewDate(2022, 6, 13))
	assert.Equal(t, pre.Reason, "test refund")
	assert.Equal(t, pre.Ext[facturae.ExtKeyCorrection], cbc.Code("01"))
	assert.Equal(t, i.Totals.Payable.String(), "900.00")
}

func TestCorrectionOptionsSchema(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		inv := testInvoiceESForCorrection(t)
		out, err := inv.CorrectionOptionsSchema()
		require.NoError(t, err)

		schema, ok := out.(*jsonschema.Schema)
		require.True(t, ok)

		cos := schema.Definitions["CorrectionOptions"]
		assert.Equal(t, 7, cos.Properties.Len())

		pm, ok := cos.Properties.Get("ext")
		require.True(t, ok)
		_, ok = pm.Properties.Get(string(tbai.ExtKeyCorrection))
		assert.False(t, ok, "should not have tbai key")
		pmp, ok := pm.Properties.Get(string(facturae.ExtKeyCorrection))
		if assert.True(t, ok) {
			assert.Len(t, pmp.OneOf, 22)
		}

		// Sorry, this is copied and pasted from the test output!
		exp := `{"properties":{"type":{"$ref":"https://gobl.org/draft-0/cbc/key","oneOf":[{"const":"credit-note","title":"Credit Note","description":"Reflects a refund either partial or complete of the preceding document. A \ncredit note effectively *extends* the previous document."},{"const":"corrective","title":"Corrective","description":"Corrected invoice that completely *replaces* the preceding document."},{"const":"debit-note","title":"Debit Note","description":"An additional set of charges to be added to the preceding document."}],"title":"Type","description":"The type of corrective invoice to produce.","default":"credit-note"},"issue_date":{"$ref":"https://gobl.org/draft-0/cal/date","title":"Issue Date","description":"When the new corrective invoice's issue date should be set to."},"series":{"$ref":"https://gobl.org/draft-0/cbc/code","title":"Series","description":"Series to assign to the new corrective invoice.","default":"TEST"},"stamps":{"items":{"$ref":"https://gobl.org/draft-0/head/stamp"},"type":"array","title":"Stamps","description":"Stamps of the previous document to include in the preceding data."},"reason":{"type":"string","title":"Reason","description":"Human readable reason for the corrective operation."},"ext":{"properties":{"es-facturae-correction":{"oneOf":[{"const":"01","title":"Invoice code"},{"const":"02","title":"Invoice series"},{"const":"03","title":"Issue date"},{"const":"04","title":"Name and surnames/Corporate name - Issuer (Sender)"},{"const":"05","title":"Name and surnames/Corporate name - Receiver"},{"const":"06","title":"Issuer's Tax Identification Number"},{"const":"07","title":"Receiver's Tax Identification Number"},{"const":"08","title":"Supplier's address"},{"const":"09","title":"Customer's address"},{"const":"10","title":"Item line"},{"const":"11","title":"Applicable Tax Rate"},{"const":"12","title":"Applicable Tax Amount"},{"const":"13","title":"Applicable Date/Period"},{"const":"14","title":"Invoice Class"},{"const":"15","title":"Legal literals"},{"const":"16","title":"Taxable Base"},{"const":"80","title":"Calculation of tax outputs"},{"const":"81","title":"Calculation of tax inputs"},{"const":"82","title":"Taxable Base modified due to return of packages and packaging materials"},{"const":"83","title":"Taxable Base modified due to discounts and rebates"},{"const":"84","title":"Taxable Base modified due to firm court ruling or administrative decision"},{"const":"85","title":"Taxable Base modified due to unpaid outputs where there is a judgement opening insolvency proceedings"}],"type":"string","title":"FacturaE Change","description":"FacturaE requires a specific and single code that explains why the previous invoice is being corrected."}},"type":"object","title":"Extensions","description":"Extensions for region specific requirements that may be added in the preceding\nor at the document level, according to the local rules.","recommended":["es-facturae-correction"]},"copy_tax":{"type":"boolean","title":"Copy Tax Totals","description":"CopyTax when true will copy the tax totals from the previous document to the\npreceding document data."}},"type":"object","required":["type"],"description":"CorrectionOptions defines a structure used to pass configuration options to correct a previous invoice.","recommended":["series","ext"]}`
		data, err := json.Marshal(cos)
		require.NoError(t, err)
		if !assert.JSONEq(t, exp, string(data)) {
			t.Log(string(data))
		}

		data, err = json.Marshal(schema)
		require.NoError(t, err)
		assert.Contains(t, string(data), `"$id":"https://gobl.org/draft-0/bill/correction-options?tax_regime=es"`)
	})
	t.Run("with copy tax", func(t *testing.T) {
		inv := testInvoiceESForCorrection(t)
		// use Verifactu as we know that copies tax
		inv.Addons = tax.WithAddons(verifactu.V1)
		require.NoError(t, inv.Calculate())
		out, err := inv.CorrectionOptionsSchema()
		require.NoError(t, err)

		schema, ok := out.(*jsonschema.Schema)
		require.True(t, ok)

		cos := schema.Definitions["CorrectionOptions"]
		assert.Equal(t, []string{"series", "ext", "copy_tax"}, cos.Extras["recommended"])
	})
}

func TestCorrectWithData(t *testing.T) {
	i := testInvoiceESForCorrection(t)
	data := []byte(`{"type":"credit-note","reason":"test refund"}`)

	err := i.Correct(
		bill.WithData(data),
		bill.WithExtension(facturae.ExtKeyCorrection, "01"),
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
		Regime: tax.WithRegime("ES"),
		Addons: tax.WithAddons(facturae.V3),
		Series: "TEST",
		Code:   "123",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
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

func testInvoicePTForCorrection(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Regime: tax.WithRegime("PT"),
		Series: "TEST",
		Code:   "123",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "PT",
				Code:    "545259045",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "PT",
				Code:    "503504030",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
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
		Regime: tax.WithRegime("FR"),
		Series: "TEST",
		Code:   "123",
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "732829320",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "391838042",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
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
		Regime: tax.WithRegime("CO"),
		Addons: tax.WithAddons(dian.V2),
		Series: "TEST",
		Code:   "123",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "CO",
				Code:    "9014586527",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: "CO",
				Code:    "8001345363",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
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
