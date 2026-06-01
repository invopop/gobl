package oioubl_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	oioubl "github.com/invopop/gobl/addons/dk/oioubl-v2-1"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Regime:    tax.WithRegime("DK"),
		Addons:    tax.WithAddons(oioubl.V2_1),
		IssueDate: cal.MakeDate(2026, 1, 1),
		Type:      "standard",
		Currency:  "DKK",
		Series:    "2026",
		Code:      "1000",
		Supplier: &org.Party{
			Name: "Eksempel A/S",
			TaxID: &tax.Identity{
				Country: "DK",
				Code:    "12345674",
			},
			Inboxes: []*org.Inbox{
				{Scheme: "0184", Code: "12345674"},
			},
			Addresses: []*org.Address{
				{Street: "Hovedgaden 1", Locality: "København", Code: "1000", Country: "DK"},
			},
		},
		Customer: &org.Party{
			Name: "Kunde ApS",
			TaxID: &tax.Identity{
				Country: "DK",
				Code:    "88146328",
			},
			Inboxes: []*org.Inbox{
				{Scheme: "0184", Code: "88146328"},
			},
			People: []*org.Person{
				{Name: &org.Name{Given: "Anders", Surname: "Jensen"}},
			},
			Addresses: []*org.Address{
				{Street: "Bygaden 5", Locality: "Aarhus", Code: "8000", Country: "DK"},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Produkt",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{Category: "VAT", Rate: "standard"},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("missing supplier inboxes (F-INV031)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Supplier.Inboxes = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-INV031")
	})

	t.Run("missing customer inboxes (F-INV044)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.Inboxes = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-INV044")
	})

	t.Run("missing customer people (F-INV046)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.People = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-INV046")
	})

	t.Run("customer with two people is allowed (loose vs F-INV046)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Customer.People = append(inv.Customer.People,
			&org.Person{Name: &org.Name{Given: "Mette", Surname: "Hansen"}},
		)
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("ordering present with only accounting cost is allowed", func(t *testing.T) {
		// OIOUBL F-INV024 only constrains cac:OrderReference/ID; an accounting
		// cost emits cbc:AccountingCost, not an OrderReference, so no code is
		// required here.
		inv := testInvoiceStandard(t)
		inv.Ordering = &bill.Ordering{Cost: "5050"}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("missing invoice code fails (F-INV009)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Code = ""
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-INV009")
	})

	t.Run("zero line quantity fails (F-INV147)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Quantity = num.MakeAmount(0, 0)
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-INV147")
	})

	t.Run("line order ref without invoice ordering fails (F-INV142)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Order = "PO-LINE-1"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-INV142")
	})

	t.Run("line order ref with invoice ordering passes", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Order = "PO-LINE-1"
		inv.Ordering = &bill.Ordering{Code: "PO-2026-001"}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("rounding above 10.00 fails (F-INV338)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		excess := num.MakeAmount(1500, 2)
		inv.Totals.Rounding = &excess
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-INV338")
	})

	t.Run("rounding below -10.00 fails (F-INV338)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		excess := num.MakeAmount(-1500, 2)
		inv.Totals.Rounding = &excess
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-INV338")
	})

	t.Run("rounding within range passes", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		require.NoError(t, inv.Calculate())
		amount := num.MakeAmount(500, 2)
		inv.Totals.Rounding = &amount
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("negative line discount fails (F-INV335)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Discounts = []*bill.LineDiscount{
			{Amount: num.MakeAmount(-500, 2)},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-INV335")
	})

	t.Run("negative line charge fails (F-INV335)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Lines[0].Charges = []*bill.LineCharge{
			{Amount: num.MakeAmount(-500, 2)},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-INV335")
	})

	t.Run("delivery with receiver and addresses passes", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Delivery = &bill.DeliveryDetails{
			Receiver: &org.Party{
				Name: "Modtager A/S",
				Addresses: []*org.Address{
					{Street: "Leveringsvej 2", Locality: "Odense", Code: "5000", Country: "DK"},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("delivery with receiver only and no identities fails (F-INV239)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Delivery = &bill.DeliveryDetails{
			Receiver: &org.Party{Name: "Modtager A/S"},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-INV239")
	})

	t.Run("delivery with receiver and identities passes (no addresses)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Delivery = &bill.DeliveryDetails{
			Receiver:   &org.Party{Name: "Modtager A/S"},
			Identities: []*org.Identity{{Code: "DEL-LOC-1"}},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("OIOUBL payment-means code 31 passes (F-LIB100)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
				Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyPaymentMeans: "31"}),
				CreditTransfer: []*pay.CreditTransfer{
					{IBAN: "DK5000400440116243"},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("generic credit-transfer code 30 passes (converter maps it to 31)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key:            pay.MeansKeyCreditTransfer,
				Ext:            tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyPaymentMeans: "30"}),
				CreditTransfer: []*pay.CreditTransfer{{IBAN: "DK5000400440116243"}},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("non-OIOUBL payment-means code fails (F-LIB100)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
				Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyPaymentMeans: "57"}),
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "F-LIB100")
	})

	t.Run("bank-transfer code 42 with account passes", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key:            pay.MeansKeyCreditTransfer,
				Ext:            tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyPaymentMeans: "42"}),
				CreditTransfer: []*pay.CreditTransfer{{Number: "1234567890"}},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("bank-transfer code 42 without account fails (F-LIB126)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
				Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyPaymentMeans: "42"}),
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "F-LIB126")
	})

	t.Run("bank-transfer code 31 without account fails (F-LIB107)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key:            pay.MeansKeyCreditTransfer,
				Ext:            tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyPaymentMeans: "31"}),
				CreditTransfer: []*pay.CreditTransfer{{Name: "Bank, no account number"}},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "F-LIB107")
	})

	t.Run("generic credit-transfer code 30 without account fails (F-LIB107)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
				Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyPaymentMeans: "30"}),
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "F-LIB107")
	})

	t.Run("Giro code 50 with payment id passes", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyPaymentMeans: "50",
					oioubl.ExtKeyPaymentID:    "04",
				}),
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("Giro code 50 without payment id fails (F-LIB144)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
				Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyPaymentMeans: "50"}),
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "F-LIB144")
	})

	t.Run("Giro code 50 with a FIK payment id fails (F-LIB147)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyPaymentMeans: "50",
					oioubl.ExtKeyPaymentID:    "71",
				}),
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "F-LIB147")
	})

	t.Run("FIK code 93 with payment id passes", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyPaymentMeans: "93",
					oioubl.ExtKeyPaymentID:    "73",
				}),
			},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("FIK code 93 without payment id fails (F-LIB152)", func(t *testing.T) {
		inv := testInvoiceStandard(t)
		inv.Payment = &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
				Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyPaymentMeans: "93"}),
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "F-LIB152")
	})
}

func testCreditNoteStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := testInvoiceStandard(t)
	inv.Type = bill.InvoiceTypeCreditNote
	inv.Code = "CN-1000"
	return inv
}

func TestCreditNoteValidation(t *testing.T) {
	t.Run("standard credit note", func(t *testing.T) {
		cn := testCreditNoteStandard(t)
		require.NoError(t, cn.Calculate())
		require.NoError(t, rules.Validate(cn))
	})

	t.Run("missing credit-note code fails (F-CRN006)", func(t *testing.T) {
		cn := testCreditNoteStandard(t)
		cn.Code = ""
		require.NoError(t, cn.Calculate())
		err := rules.Validate(cn)
		assert.ErrorContains(t, err, "F-INV009")
	})

	t.Run("zero credit-note line quantity fails (F-CRN088)", func(t *testing.T) {
		cn := testCreditNoteStandard(t)
		cn.Lines[0].Quantity = num.MakeAmount(0, 0)
		require.NoError(t, cn.Calculate())
		err := rules.Validate(cn)
		assert.ErrorContains(t, err, "F-INV147")
	})

	t.Run("missing supplier inboxes fails (F-CRN028)", func(t *testing.T) {
		cn := testCreditNoteStandard(t)
		cn.Supplier.Inboxes = nil
		require.NoError(t, cn.Calculate())
		err := rules.Validate(cn)
		assert.ErrorContains(t, err, "F-INV031")
	})

	t.Run("credit note with line order ref does not fire F-INV142", func(t *testing.T) {
		cn := testCreditNoteStandard(t)
		cn.Lines[0].Order = "PO-LINE-1"
		require.NoError(t, cn.Calculate())
		assert.NoError(t, rules.Validate(cn))
	})

	t.Run("negative credit-note line discount fails (F-CRN203)", func(t *testing.T) {
		cn := testCreditNoteStandard(t)
		cn.Lines[0].Discounts = []*bill.LineDiscount{
			{Amount: num.MakeAmount(-500, 2)},
		}
		require.NoError(t, cn.Calculate())
		err := rules.Validate(cn)
		assert.ErrorContains(t, err, "F-INV335")
	})
}
