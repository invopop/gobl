package zatca_test

import (
	"testing"

	"github.com/invopop/gobl/addons/sa/zatca"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/catalogues/cef"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/sa"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Test fixtures
// ============================================================================

// validStandardInvoice returns a fully valid ZATCA standard tax invoice
// (KSA-2 = "0100000", document type "388"). Used as the base for every
// positive and negative test in this file.
func validStandardInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime:    tax.WithRegime("SA"),
		Addons:    tax.WithAddons(zatca.V1),
		Type:      bill.InvoiceTypeStandard,
		Currency:  "SAR",
		Code:      "INV-001",
		IssueDate: cal.MakeDate(2024, 6, 15),
		IssueTime: cal.NewTime(12, 0, 0),
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				zatca.ExtKeyInvoiceTypeTransactions: "0100000",
				untdid.ExtKeyDocumentType:           "388",
			}),
		},
		Supplier: validSupplier(),
		Customer: validCustomer(),
		Delivery: &bill.DeliveryDetails{
			Date: cal.NewDate(2024, 6, 15),
			Period: &cal.Period{
				Start: cal.MakeDate(2024, 6, 1),
				End:   cal.MakeDate(2024, 6, 30),
			},
		},
		Payment: &bill.PaymentDetails{
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
			},
			Terms: &pay.Terms{
				Notes: "Payment due within 30 days",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Development services",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateGeneral,
					},
				},
			},
		},
	}
}

func validSupplier() *org.Party {
	return &org.Party{
		Name: "Acme Corp Saudi",
		TaxID: &tax.Identity{
			Country: "SA",
			Code:    "300000000000003",
		},
		Identities: []*org.Identity{
			{Type: sa.IdentityTypeCRN, Code: "1234567890"},
		},
		Addresses: []*org.Address{validAddress()},
	}
}

func validCustomer() *org.Party {
	return &org.Party{
		Name: "Sample Consumer LLC",
		TaxID: &tax.Identity{
			Country: "SA",
			Code:    "399999999900003",
		},
		Addresses: []*org.Address{validAddress()},
	}
}

func validAddress() *org.Address {
	return &org.Address{
		Street:      "King Fahd Road",
		StreetExtra: "Al Olaya",
		Number:      "1234",
		Locality:    "Riyadh",
		Code:        "12345",
		Country:     "SA",
	}
}

// validSimplifiedInvoice returns a simplified tax invoice (KSA-2 starts with "02").
func validSimplifiedInvoice() *bill.Invoice {
	inv := validStandardInvoice()
	inv.Tax.Ext = inv.Tax.Ext.Set(zatca.ExtKeyInvoiceTypeTransactions, "0200000")
	return inv
}

// validSummaryInvoice returns a standard summary tax invoice — KSA-2 with
// the summary bit set (position 6 = 1).
func validSummaryInvoice() *bill.Invoice {
	inv := validStandardInvoice()
	inv.Tax.Ext = inv.Tax.Ext.Set(zatca.ExtKeyInvoiceTypeTransactions, "0100010")
	return inv
}

// validCreditNote returns a credit note (document type 381) that references
// a preceding invoice with reason — both required by ZATCA.
func validCreditNote() *bill.Invoice {
	inv := validStandardInvoice()
	inv.Type = bill.InvoiceTypeCreditNote
	inv.Tax.Ext = inv.Tax.Ext.Set(untdid.ExtKeyDocumentType, "381")
	inv.Preceding = []*org.DocumentRef{
		{
			Code:      "INV-001",
			IssueDate: cal.NewDate(2024, 5, 1),
			Reason:    "Return of goods",
		},
	}
	return inv
}

// validDebitNote returns a debit note (document type 383).
func validDebitNote() *bill.Invoice {
	inv := validCreditNote()
	inv.Type = bill.InvoiceTypeDebitNote
	inv.Tax.Ext = inv.Tax.Ext.Set(untdid.ExtKeyDocumentType, "383")
	inv.Preceding[0].Reason = "Price adjustment"
	return inv
}

// validExportInvoice returns a standard tax invoice with the export bit
// set in KSA-2 (position 5 = 1) and a non-VAT-registered customer, as
// required by BR-KSA-46.
func validExportInvoice() *bill.Invoice {
	inv := validStandardInvoice()
	inv.Tax.Ext = inv.Tax.Ext.Set(zatca.ExtKeyInvoiceTypeTransactions, "0100100")
	inv.Customer.TaxID = nil
	inv.Customer.Identities = []*org.Identity{
		{Type: sa.IdentityTypeTIN, Code: "123456789012345"},
	}
	return inv
}

// calculated runs Calculate on the invoice and returns it, failing the
// test if calculation fails. Convenience wrapper used in the tests.
func calculated(t *testing.T, inv *bill.Invoice) *bill.Invoice {
	t.Helper()
	require.NoError(t, inv.Calculate())
	return inv
}

// ============================================================================
// Sanity: every fixture must validate clean
// ============================================================================

func TestValidFixtures(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		require.NoError(t, rules.Validate(calculated(t, validStandardInvoice())))
	})
	t.Run("simplified invoice", func(t *testing.T) {
		require.NoError(t, rules.Validate(calculated(t, validSimplifiedInvoice())))
	})
	t.Run("summary invoice", func(t *testing.T) {
		require.NoError(t, rules.Validate(calculated(t, validSummaryInvoice())))
	})
	t.Run("credit note", func(t *testing.T) {
		require.NoError(t, rules.Validate(calculated(t, validCreditNote())))
	})
	t.Run("debit note", func(t *testing.T) {
		require.NoError(t, rules.Validate(calculated(t, validDebitNote())))
	})
	t.Run("export invoice", func(t *testing.T) {
		require.NoError(t, rules.Validate(calculated(t, validExportInvoice())))
	})
}

// ============================================================================
// Group A — Rule 01: Issue time required (KSA-25, BR-KSA-70)
// ============================================================================

func TestBRKSA70_IssueTime(t *testing.T) {
	t.Run("present is valid", func(t *testing.T) {
		inv := validStandardInvoice()
		require.NotNil(t, inv.IssueTime)
		require.NoError(t, rules.Validate(calculated(t, inv)))
	})

	t.Run("nil triggers normaliser to auto-fill", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.IssueTime = nil
		require.NoError(t, inv.Calculate())
		assert.NotNil(t, inv.IssueTime, "the SA-ZATCA normaliser should auto-create issue time")
		require.NoError(t, rules.Validate(inv))
	})
}

// ============================================================================
// Group B — Rules 02-05: Tax block extensions
// ============================================================================

func TestTaxBlockExtensions(t *testing.T) {
	// Note: rule 02 ("tax must be present") is unreachable from invoice tests
	// because Calculate() auto-creates an empty Tax block before validation.
	// Coverage of that path would require directly invoking the rule chain.

	t.Run("required ext keys missing", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Tax.Ext = tax.Extensions{}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"extensions keys untdid document type and invoice type transaction are required")
	})

	t.Run("invalid document type code rejected", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Tax.Ext = inv.Tax.Ext.Set(untdid.ExtKeyDocumentType, "999")
		assert.ErrorContains(t, rules.Validate(inv),
			"document type must be a valid ZATCA type")
	})

	t.Run("invalid invoice transaction type rejected", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Tax.Ext = inv.Tax.Ext.Set(zatca.ExtKeyInvoiceTypeTransactions, "9999999")
		assert.ErrorContains(t, rules.Validate(inv),
			"invoice transaction type must be valid")
	})
}

// ============================================================================
// Group C — Rules 06-08: Credit/debit notes need preceding billing reference
//                       (BR-KSA-17, BR-KSA-56)
// ============================================================================

func TestBRKSA17_56_CreditDebitBillingReference(t *testing.T) {
	t.Run("credit note without preceding fails", func(t *testing.T) {
		inv := validCreditNote()
		inv.Preceding = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"credit and debit notes must have a billing reference")
	})

	t.Run("credit note preceding without code fails (BR-KSA-56)", func(t *testing.T) {
		inv := validCreditNote()
		inv.Preceding[0].Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"billing reference must have an identifier")
	})

	t.Run("credit note preceding without reason fails (BR-KSA-17)", func(t *testing.T) {
		inv := validCreditNote()
		inv.Preceding[0].Reason = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"credit and debit notes must contain the reason for issuance")
	})

	t.Run("debit note preceding without reason fails (BR-KSA-17)", func(t *testing.T) {
		inv := validDebitNote()
		inv.Preceding[0].Reason = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"credit and debit notes must contain the reason for issuance")
	})
}

// ============================================================================
// Group D — Rules 09-10: Supplier name and address (BR-06, BR-KSA-09)
// ============================================================================

func TestSupplierRequirements(t *testing.T) {
	t.Run("missing name fails (BR-06)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Name = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "supplier name is required")
	})

	t.Run("missing addresses fails (BR-KSA-09)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Addresses = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "supplier addresses are required")
	})
}

// ============================================================================
// Group E — Rules 11-20: Standard tax invoice cascade
//   (customer presence/name/address, identification, lines, delivery)
// ============================================================================

func TestStandardInvoiceRequirements(t *testing.T) {
	t.Run("customer missing fails (rule 11)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "customer must be present")
	})

	t.Run("customer name missing fails (BR-KSA-42)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Name = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"customer name must be present (BR-KSA-71), (BR-KSA-25), (BR-KSA-42)")
	})

	t.Run("customer address street missing fails (BR-KSA-10)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Addresses[0].Street = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"customer address must have a street name")
	})

	t.Run("customer address city missing fails (BR-KSA-10)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Addresses[0].Locality = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"customer address must have a city name")
	})

	t.Run("customer address country missing fails (BR-KSA-10)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Addresses[0].Country = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"customer address must have a country code")
	})

	t.Run("customer with valid identity but no VAT passes (BR-KSA-14, BR-KSA-81)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{Type: sa.IdentityTypeCRN, Code: "1010101010"},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("customer with no VAT and no identity fails (BR-KSA-14, BR-KSA-81)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.TaxID = nil
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"customer must have a valid identification scheme")
	})

	t.Run("line missing taxes fails (BR-KSA-52)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines[0].Taxes = nil
		_ = inv.Calculate()
		assert.ErrorContains(t, rules.Validate(inv),
			"line taxes are required")
	})

	t.Run("delivery missing fails", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Delivery = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "delivery must be present")
	})

	t.Run("delivery date missing fails (BR-KSA-15)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Delivery.Date = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"delivery must have a supply date")
	})
}

// ============================================================================
// Group F — Rule 21: Export invoices must NOT carry buyer VAT (BR-KSA-46)
// ============================================================================

func TestBRKSA46_ExportInvoiceCustomerVAT(t *testing.T) {
	t.Run("export without buyer VAT is valid", func(t *testing.T) {
		require.NoError(t, rules.Validate(calculated(t, validExportInvoice())))
	})

	t.Run("export with buyer VAT registration fails", func(t *testing.T) {
		inv := validExportInvoice()
		inv.Customer.TaxID = &tax.Identity{Country: "SA", Code: "399999999900003"}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"export invoices must not have buyer VAT registration number")
	})
}

// ============================================================================
// Group G — Rules 22-27: Simplified-summary invoices need delivery period
//                         and customer name (BR-KSA-71, BR-KSA-72)
// ============================================================================

func TestSimplifiedSummaryRequirements(t *testing.T) {
	// "Simplified-and-summary" means an invoice with KSA-2 of the form 02....1.
	// Build one by deriving from simplified and adding the summary bit.
	build := func() *bill.Invoice {
		inv := validSimplifiedInvoice()
		inv.Tax.Ext = inv.Tax.Ext.Set(zatca.ExtKeyInvoiceTypeTransactions, "0200010")
		return inv
	}

	t.Run("baseline simplified+summary is valid", func(t *testing.T) {
		require.NoError(t, rules.Validate(calculated(t, build())))
	})

	t.Run("delivery missing fails", func(t *testing.T) {
		inv := build()
		inv.Delivery = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"delivery must be present for simplified and summary invoices")
	})

	t.Run("delivery period missing fails", func(t *testing.T) {
		inv := build()
		inv.Delivery.Period = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"supply must have a delivery period")
	})

	// Period start/end zero values are first caught by the global cal.Period
	// validator (CAL-PERIOD-01 / CAL-PERIOD-02), so the SA-specific message
	// from rules 24/25 may not surface verbatim. We just assert that
	// validation fails when either bound is missing.
	t.Run("delivery period start missing fails (BR-KSA-72)", func(t *testing.T) {
		inv := build()
		inv.Delivery.Period.Start = cal.Date{}
		require.NoError(t, inv.Calculate())
		assert.Error(t, rules.Validate(inv))
	})

	t.Run("delivery period end missing fails (BR-KSA-72)", func(t *testing.T) {
		inv := build()
		inv.Delivery.Period.End = cal.Date{}
		require.NoError(t, inv.Calculate())
		assert.Error(t, rules.Validate(inv))
	})

	t.Run("customer name missing fails (BR-KSA-71)", func(t *testing.T) {
		inv := build()
		inv.Customer.Name = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"customer name must be present (BR-KSA-71), (BR-KSA-25), (BR-KSA-42)")
	})
}

// ============================================================================
// Group H — Rules 28-30: EDU/HEA exemptions require customer NAT identity
//                        and a customer name on simplified docs
//                        (BR-KSA-25, BR-KSA-49)
// ============================================================================

func TestBRKSA49_25_EDUHEAExemption(t *testing.T) {
	// Helper: replace line tax with a Z-rated combo carrying the given VATEX.
	withEDUHEACombo := func(inv *bill.Invoice, vatex cbc.Code) {
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: tax.CategoryVAT,
				Key:      tax.KeyZero,
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					cef.ExtKeyVATEX: vatex,
				}),
			},
		}
	}

	t.Run("EDU exemption with NAT identity passes", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Identities = []*org.Identity{
			{Type: sa.IdentityTypeNational, Code: "1234567890"},
		}
		withEDUHEACombo(inv, "VATEX-SA-EDU")
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("HEA exemption with NAT identity passes", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Identities = []*org.Identity{
			{Type: sa.IdentityTypeNational, Code: "1234567890"},
		}
		withEDUHEACombo(inv, "VATEX-SA-HEA")
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("EDU exemption without NAT identity fails (BR-KSA-49)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Identities = []*org.Identity{
			{Type: sa.IdentityTypeCRN, Code: "1010101010"},
		}
		withEDUHEACombo(inv, "VATEX-SA-EDU")
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"customer must have a national ID (NAT)")
	})

	t.Run("non-EDU/HEA exemption does NOT require NAT", func(t *testing.T) {
		inv := validStandardInvoice()
		withEDUHEACombo(inv, "VATEX-SA-32")
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified+EDU without customer name fails (BR-KSA-25)", func(t *testing.T) {
		inv := validSimplifiedInvoice()
		inv.Customer.Name = ""
		inv.Customer.Identities = []*org.Identity{
			{Type: sa.IdentityTypeNational, Code: "1234567890"},
		}
		withEDUHEACombo(inv, "VATEX-SA-EDU")
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv),
			"customer name must be present (BR-KSA-71), (BR-KSA-25), (BR-KSA-42)")
	})
}

// ============================================================================
// Defensive nil guards coverage
// ============================================================================

func TestNilGuards(t *testing.T) {
	t.Run("standard invoice with nil customer fails", func(t *testing.T) {
		inv := validStandardInvoice()
		require.NoError(t, inv.Calculate())
		inv.Customer = nil
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer must be present")
	})

	t.Run("standard invoice with nil tax after calculate", func(t *testing.T) {
		inv := validStandardInvoice()
		require.NoError(t, inv.Calculate())
		inv.Tax = nil
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "tax must be present")
	})
}
