package favat_test

import (
	"testing"

	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func creditNote() *bill.Invoice {
	inv := &bill.Invoice{
		Regime:    tax.WithRegime("PL"),
		Addons:    tax.WithAddons(favat.V3),
		Currency:  currency.PLN,
		Code:      "TEST",
		Type:      bill.InvoiceTypeCreditNote,
		IssueDate: cal.MakeDate(2022, 12, 29),
		Preceding: []*org.DocumentRef{
			{
				Code:      "TEST",
				IssueDate: cal.NewDate(2022, 12, 27),
				Ext: tax.Extensions{
					favat.ExtKeyEffectiveDate: "1",
				},
			},
		},
		Supplier: &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "PL",
				Code:    "1111111111",
			},
			Addresses: []*org.Address{
				{
					Locality: "Foo",
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "PL",
				Code:    "2222222222",
			},
			Addresses: []*org.Address{
				{
					Locality: "Foo",
				},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 3),
				Item: &org.Item{
					Name:  "bogus",
					Price: num.NewAmount(1000, 3),
				},
			},
		},
	}
	return inv
}

func standardInvoice() *bill.Invoice {
	inv := &bill.Invoice{
		Regime:    tax.WithRegime("PL"),
		Addons:    tax.WithAddons(favat.V3),
		Currency:  currency.PLN,
		Code:      "STD",
		Type:      bill.InvoiceTypeStandard,
		IssueDate: cal.MakeDate(2022, 12, 29),
		Supplier: &org.Party{
			Name: "Test Party",
			TaxID: &tax.Identity{
				Country: "PL",
				Code:    "1111111111",
			},
			Addresses: []*org.Address{
				{
					Locality: "Foo",
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "PL",
				Code:    "2222222222",
			},
			Addresses: []*org.Address{
				{
					Locality: "Foo",
				},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 3),
				Item: &org.Item{
					Name:  "standard-bogus",
					Price: num.NewAmount(1000, 3),
				},
			},
		},
	}
	return inv
}

func TestBasicCreditNoteValidation(t *testing.T) {
	inv := creditNote()
	inv.Preceding[0].Reason = "Correcting an error"
	err := inv.Calculate()
	require.NoError(t, err)
	err = inv.Validate()
	assert.NoError(t, err)
	assert.Equal(t, inv.Preceding[0].Ext[favat.ExtKeyEffectiveDate], cbc.Code("1"))

	inv.Preceding[0].Ext["foo"] = "bar"
	err = inv.Validate()
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "preceding: (0: (ext: (foo: undefined.).).)")
	}
}

func TestBasicStandardInvoiceValidation(t *testing.T) {
	inv := standardInvoice()
	err := inv.Calculate()
	require.NoError(t, err)
	err = inv.Validate()
	assert.NoError(t, err)
}

func TestExemptStandardInvoiceValidation(t *testing.T) {
	inv := standardInvoice()
	inv.Tags = tax.Tags{List: []cbc.Key{tax.KeyExempt}}
	inv.Tax = &bill.Tax{
		Ext: tax.Extensions{
			favat.ExtKeyExemption: "A",
		},
	}
	inv.Notes = []*org.Note{
		{
			Key:  org.NoteKeyLegal,
			Code: "A",
			Src:  favat.ExtKeyExemption,
			Text: "Art. 25a ust. 1 pkt 9 ustawy o VAT",
		},
	}

	err := inv.Calculate()
	require.NoError(t, err)
	err = inv.Validate()
	assert.NoError(t, err)
}

func TestExemptStandardInvoiceValidationFailsOnMismatchedNoteCode(t *testing.T) {
	inv := standardInvoice()
	inv.Tags = tax.Tags{List: []cbc.Key{tax.KeyExempt}}
	inv.Tax = &bill.Tax{
		Ext: tax.Extensions{
			favat.ExtKeyExemption: "B",
		},
	}
	inv.Notes = []*org.Note{
		{
			Key:  org.NoteKeyLegal,
			Code: "A",
			Src:  favat.ExtKeyExemption,
			Text: "Art. 25a ust. 1 pkt 9 ustawy o VAT",
		},
	}

	err := inv.Calculate()
	require.NoError(t, err)

	err = inv.Validate()
	assert.Error(t, err)
}

func TestExemptStandardInvoiceValidationFailsWithoutNote(t *testing.T) {
	inv := standardInvoice()
	inv.Tags = tax.Tags{List: []cbc.Key{tax.KeyExempt}}
	inv.Tax = &bill.Tax{
		// valid exemption set but no matching note
		Ext: tax.Extensions{
			favat.ExtKeyExemption: "A",
		},
	}

	err := inv.Calculate()
	require.NoError(t, err)

	err = inv.Validate()
	assert.Error(t, err)
}

func TestExemptStandardInvoiceValidationFailsWithoutTax(t *testing.T) {
	inv := standardInvoice()
	inv.Tags = tax.Tags{List: []cbc.Key{tax.KeyExempt}}
	inv.Notes = []*org.Note{
		{
			Key:  org.NoteKeyLegal,
			Code: "A",
			Src:  favat.ExtKeyExemption,
			Text: "Art. 25a ust. 1 pkt 9 ustawy o VAT",
		},
	}

	err := inv.Calculate()
	require.NoError(t, err)

	err = inv.Validate()
	assert.Error(t, err)
}
