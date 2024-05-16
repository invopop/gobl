package bill_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
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

	i := testInvoicePTForCorrection(t)
	err := i.Correct(bill.Corrective)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid correction type: corrective")

	i = testInvoiceESForCorrection(t)
	err = i.Correct(bill.Credit,
		bill.WithReason("test refund"),
		bill.WithExtension(es.ExtKeyFacturaECorrection, "01"),
	)
	require.NoError(t, err)
	assert.Equal(t, bill.InvoiceTypeCreditNote, i.Type)
	assert.Equal(t, i.Lines[0].Quantity.String(), "10")
	assert.Equal(t, i.IssueDate, cal.Today())
	assert.Equal(t, i.Series, "TEST")
	assert.Empty(t, i.Code)
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

	i = testInvoicePTForCorrection(t)
	err = i.Correct(bill.Corrective, bill.WithReason("should fail"))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid correction type: corrective")

	i = testInvoiceESForCorrection(t)
	err = i.Correct(
		bill.Corrective,
		bill.WithExtension(es.ExtKeyFacturaECorrection, "01"),
	)
	require.NoError(t, err)
	assert.Equal(t, i.Type, bill.InvoiceTypeCorrective)

	// With preset date
	i = testInvoiceESForCorrection(t)
	d := cal.MakeDate(2023, 6, 13)
	err = i.Correct(
		bill.Credit,
		bill.WithIssueDate(d),
		bill.WithExtension(es.ExtKeyFacturaECorrection, "01"),
	)
	require.NoError(t, err)
	assert.Equal(t, i.IssueDate, d)

	t.Run("with series", func(t *testing.T) {
		inv := testInvoiceESForCorrection(t)
		err := inv.Correct(bill.Credit, bill.WithSeries("R-TEST"))
		require.NoError(t, err)
		assert.Equal(t, inv.Series, "R-TEST")
		assert.Equal(t, inv.Preceding[0].Series, "TEST")
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
		bill.WithExtension(co.ExtKeyDIANCorrection, "2"),
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
		Type:   bill.InvoiceTypeCreditNote,
		Reason: "test refund",
		Ext: tax.Extensions{
			es.ExtKeyFacturaECorrection: "01",
		},
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
	assert.Equal(t, pre.Ext[es.ExtKeyFacturaECorrection], tax.ExtValue("01"))
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

	pm, ok := cos.Properties.Get("ext")
	require.True(t, ok)
	pmp, ok := pm.Properties.Get(string(es.ExtKeyFacturaECorrection))
	require.True(t, ok)
	assert.Len(t, pmp.OneOf, 22)

	// Sorry, this is copied and pasted from the test output!
	exp := `{"properties":{"type":{"$ref":"https://gobl.org/draft-0/cbc/key","oneOf":[{"const":"corrective","title":"Corrective","description":"Corrected invoice that completely *replaces* the preceding document."},{"const":"credit-note","title":"Credit Note","description":"Reflects a refund either partial or complete of the preceding document. A \ncredit note effectively *extends* the previous document."},{"const":"debit-note","title":"Debit Note","description":"An additional set of charges to be added to the preceding document."}],"title":"Type","description":"The type of corrective invoice to produce."},"issue_date":{"$ref":"https://gobl.org/draft-0/cal/date","title":"Issue Date","description":"When the new corrective invoice's issue date should be set to."},"series":{"type":"string","title":"Series","description":"Series to assign to the new corrective invoice."},"stamps":{"items":{"$ref":"https://gobl.org/draft-0/head/stamp"},"type":"array","title":"Stamps","description":"Stamps of the previous document to include in the preceding data."},"reason":{"type":"string","title":"Reason","description":"Human readable reason for the corrective operation."},"ext":{"properties":{"es-facturae-correction":{"oneOf":[{"const":"01","title":"Invoice code"},{"const":"02","title":"Invoice series"},{"const":"03","title":"Issue date"},{"const":"04","title":"Name and surnames/Corporate name - Issuer (Sender)"},{"const":"05","title":"Name and surnames/Corporate name - Receiver"},{"const":"06","title":"Issuer's Tax Identification Number"},{"const":"07","title":"Receiver's Tax Identification Number"},{"const":"08","title":"Supplier's address"},{"const":"09","title":"Customer's address"},{"const":"10","title":"Item line"},{"const":"11","title":"Applicable Tax Rate"},{"const":"12","title":"Applicable Tax Amount"},{"const":"13","title":"Applicable Date/Period"},{"const":"14","title":"Invoice Class"},{"const":"15","title":"Legal literals"},{"const":"16","title":"Taxable Base"},{"const":"80","title":"Calculation of tax outputs"},{"const":"81","title":"Calculation of tax inputs"},{"const":"82","title":"Taxable Base modified due to return of packages and packaging materials"},{"const":"83","title":"Taxable Base modified due to discounts and rebates"},{"const":"84","title":"Taxable Base modified due to firm court ruling or administrative decision"},{"const":"85","title":"Taxable Base modified due to unpaid outputs where there is a judgement opening insolvency proceedings"}],"type":"string","title":"FacturaE Change","description":"FacturaE requires a specific and single code that explains why the previous invoice is being corrected."},"es-tbai-correction":{"oneOf":[{"const":"R1","title":"Rectified invoice: error based on law and Article 80 One, Two and Six of the Provincial Tax Law of VAT"},{"const":"R2","title":"Rectified invoice: error based on law and Article 80 Three of the Provincial Tax Law of VAT"},{"const":"R3","title":"Rectified invoice: error based on law and Article 80 Four of the Provincial Tax Law of VAT"},{"const":"R4","title":"Rectified invoice: Other"},{"const":"R5","title":"Rectified invoice: simplified invoices"}],"type":"string","title":"TicketBAI Rectification Type Code","description":"Corrected or rectified invoices that need to be sent in the TicketBAI format\nrequire a specific type code to be defined alongside the preceding invoice\ndata."}},"type":"object","title":"Extensions","description":"Extensions for region specific requirements.","recommended":["es-facturae-correction","es-tbai-correction"]}},"type":"object","required":["type"],"description":"CorrectionOptions defines a structure used to pass configuration options to correct a previous invoice.","recommended":["series","ext"]}`
	data, err := json.Marshal(cos)
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
		bill.WithExtension(es.ExtKeyFacturaECorrection, "01"),
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

func testInvoicePTForCorrection(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Series: "TEST",
		Code:   "123",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.PT,
				Code:    "545259045",
			},
		},
		Customer: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.PT,
				Code:    "503504030",
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
