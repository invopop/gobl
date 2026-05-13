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
	_ "github.com/invopop/gobl/regimes/sa"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Test fixtures
// ============================================================================

// validStandardInvoice returns a fully valid ZATCA standard tax invoice.
// Used as the base for every positive and negative test in this file.
// The KSA-2 invoice-transaction-type (-> "0100000") and the UNTDID document
// type (-> "388") are populated by the ZATCA / EN 16931 scenarios during
// Calculate(), driven by inv.Type = "standard" with no tags.
func validStandardInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime:    tax.WithRegime("SA"),
		Addons:    tax.WithAddons(zatca.V1),
		Type:      bill.InvoiceTypeStandard,
		Currency:  "SAR",
		Code:      "INV-001",
		IssueDate: cal.MakeDate(2024, 6, 15),
		IssueTime: cal.NewTime(12, 0, 0),
		Tax:       &bill.Tax{}, // KSA-2 + UNTDID populated by scenarios during Calculate()
		Supplier:  validSupplier(),
		Customer:  validCustomer(),
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
			{Type: zatca.IdentityTypeCRN, Code: "1234567890"},
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

// validSimplifiedInvoice returns a simplified tax invoice. tax.TagSimplified
// drives the ZATCA scenario which populates KSA-2 = "0200000" during
// Calculate(); no hard-coding needed.
func validSimplifiedInvoice() *bill.Invoice {
	inv := validStandardInvoice()
	inv.SetTags(tax.TagSimplified)
	return inv
}

// validSummaryInvoice returns a standard summary tax invoice — KSA-2 with
// the summary bit set (position 6 = 1). zatca.TagSummary drives the ZATCA
// "Standard + Summary" scenario which populates KSA-2 = "0100010" during
// Calculate(); no hard-coding needed.
func validSummaryInvoice() *bill.Invoice {
	inv := validStandardInvoice()
	inv.SetTags(zatca.TagSummary)
	return inv
}

// validCreditNote returns a credit note that references a preceding invoice
// with reason — both required by ZATCA. The UNTDID document type and KSA-2
// are both populated by scenarios during Calculate(): EN 16931 sets UNTDID
// from inv.Type, and the ZATCA credit/debit-note scenario sets KSA-2 to
// "0100000".
func validCreditNote() *bill.Invoice {
	inv := validStandardInvoice()
	inv.Type = bill.InvoiceTypeCreditNote
	inv.Preceding = []*org.DocumentRef{
		{
			Code:      "INV-001",
			IssueDate: cal.NewDate(2024, 5, 1),
			Reason:    "Return of goods",
		},
	}
	return inv
}

// validDebitNote returns a debit note. UNTDID document type and KSA-2 are
// both populated by scenarios during Calculate() (same path as
// validCreditNote — the ZATCA credit/debit scenario covers both types).
func validDebitNote() *bill.Invoice {
	inv := validCreditNote()
	inv.Type = bill.InvoiceTypeDebitNote
	inv.Preceding[0].Reason = "Price adjustment"
	return inv
}

// validExportInvoice returns a standard tax invoice with the export bit
// set in KSA-2 (position 5 = 1). tax.TagExport drives the ZATCA scenario
// which populates KSA-2 = "0100100" during Calculate(); no hard-coding
// needed. Customer is non-VAT-registered as required by BR-KSA-46.
func validExportInvoice() *bill.Invoice {
	inv := validStandardInvoice()
	inv.SetTags(tax.TagExport)
	inv.Customer.TaxID = nil
	inv.Customer.Identities = []*org.Identity{
		{Type: zatca.IdentityTypeTIN, Code: "123456789012345"},
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
	//
	// The UNTDID document type extension is populated by EN 16931 and ZATCA
	// scenarios from inv.Type during Calculate(), so tests no longer pre-set
	// it or assert on its specific values — that coverage lives in the EN 16931
	// addon's scenarios_test.go.

	t.Run("missing invoice transaction type ext key", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Tax.Ext = inv.Tax.Ext.Delete(zatca.ExtKeyInvoiceTypeTransactions)
		assert.ErrorContains(t, rules.Validate(inv),
			"invoice transaction type extension is required")
	})

	t.Run("invalid invoice transaction type rejected", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Type = bill.InvoiceTypeOther
		inv.Tax.Ext = inv.Tax.Ext.
			Set(untdid.ExtKeyDocumentType, "388").
			Set(zatca.ExtKeyInvoiceTypeTransactions, "9999999")
		require.NoError(t, inv.Calculate())
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
// Group D.2 — Rules 25-26: Supplier tax ID and identity (BR-KSA-39, BR-KSA-08)
// ============================================================================

func TestSupplierTaxIDRequired(t *testing.T) {
	t.Run("valid invoice passes", func(t *testing.T) {
		inv := validStandardInvoice()
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("supplier without tax ID fails", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.TaxID = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have a tax ID code")
	})

	t.Run("supplier with empty tax ID code fails", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.TaxID = &tax.Identity{Country: "SA", Code: ""}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have a tax ID code")
	})
}

func TestSupplierIdentities(t *testing.T) {
	withIdentity := func(idType cbc.Code, code string) *bill.Invoice {
		inv := validStandardInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Type: idType, Code: cbc.Code(code)},
		}
		return inv
	}

	t.Run("supplier with no identities fails", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Identities = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have a valid identity")
	})

	t.Run("CRN identity passes", func(t *testing.T) {
		inv := withIdentity(zatca.IdentityTypeCRN, "1234567890")
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("MOM identity passes", func(t *testing.T) {
		inv := withIdentity(zatca.IdentityTypeMom, "1234567890")
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("MLS identity passes", func(t *testing.T) {
		inv := withIdentity(zatca.IdentityTypeMLS, "1234567890")
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("700 identity passes", func(t *testing.T) {
		inv := withIdentity(zatca.IdentityType700, "1234567890")
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("SAG identity passes", func(t *testing.T) {
		inv := withIdentity(zatca.IdentityTypeSAG, "1234567890")
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("OTH identity passes", func(t *testing.T) {
		inv := withIdentity(zatca.IdentityTypeOTH, "1234567890")
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("invalid identity type fails", func(t *testing.T) {
		inv := withIdentity("INVALID", "1234567890")
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have a valid identity")
	})

	t.Run("two identities fails", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Type: zatca.IdentityTypeCRN, Code: "1234567890"},
			{Type: zatca.IdentityTypeMLS, Code: "1234567890"},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have a valid identity")
	})

	t.Run("NAT identity type fails (not in supplier list)", func(t *testing.T) {
		inv := withIdentity(zatca.IdentityTypeNational, "1234567890")
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have a valid identity")
	})

	t.Run("TIN identity type fails (not in supplier list)", func(t *testing.T) {
		inv := withIdentity(zatca.IdentityTypeTIN, "123456789012345")
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier must have a valid identity")
	})
}

// ============================================================================
// Group E — Rules 10-16: Standard tax invoice cascade
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
			{Type: zatca.IdentityTypeCRN, Code: "1010101010"},
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
	// Build one by deriving from simplified and adding the summary tag.
	// The [TagSimplified, TagSummary] combination drives the ZATCA
	// "Simplified and summary" scenario which sets KSA-2 = "0200010"
	// during Calculate(); no hard-coding needed.
	build := func() *bill.Invoice {
		inv := validSimplifiedInvoice()
		inv.SetTags(tax.TagSimplified, zatca.TagSummary)
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
			{Type: zatca.IdentityTypeNational, Code: "1234567890"},
		}
		withEDUHEACombo(inv, "VATEX-SA-EDU")
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("HEA exemption with NAT identity passes", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Identities = []*org.Identity{
			{Type: zatca.IdentityTypeNational, Code: "1234567890"},
		}
		withEDUHEACombo(inv, "VATEX-SA-HEA")
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("EDU exemption without NAT identity fails (BR-KSA-49)", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Identities = []*org.Identity{
			{Type: zatca.IdentityTypeCRN, Code: "1010101010"},
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
			{Type: zatca.IdentityTypeNational, Code: "1234567890"},
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

// ============================================================================
// Group I — Every code in validTransactionTypes must validate cleanly
// ============================================================================
//
// Codes registered as scenarios in scenarios.go are exercised via $tags so
// that the scenario engine populates KSA-2 during Calculate(). Codes NOT
// registered as scenarios use bill.InvoiceTypeOther, which by contract
// bypasses every regime/addon scenario (see bill/invoice_type.go: "implies
// that any scenarios defined in tax regimes or addons will not be
// applied"). Because no scenario fires, both UNTDID and KSA-2 must be set
// manually on inv.Tax.Ext. Validation still routes the right rule cascade
// because invoiceIsStandard / invoiceIsExport / invoiceIsSimplifiedAndSummary
// inspect the KSA-2 prefix and bits, not inv.Type.
//
// Each KSA-2 code is its own subtest — no for-loops — so a failure points
// at the exact code that broke.

// hasExportBit returns true when KSA-2 position 4 (the export flag) is set,
// which triggers BR-KSA-46 — the customer must NOT carry a VAT registration.
func hasExportBit(code cbc.Code) bool {
	return len(code) >= 5 && code[4] == '1'
}

// clearCustomerVAT removes the customer's VAT ID and replaces it with a
// non-VAT TIN identity so the invoice satisfies BR-KSA-46 (export bit set)
// while still meeting BR-KSA-14/81 when the form is Standard (TT=01).
func clearCustomerVAT(inv *bill.Invoice) {
	inv.Customer.TaxID = nil
	inv.Customer.Identities = []*org.Identity{
		{Type: zatca.IdentityTypeTIN, Code: "123456789012345"},
	}
}

// assertSACodeViaScenario builds a Standard-type invoice with the provided
// tags, runs Calculate (scenarios fire), validates, and asserts the
// resulting KSA-2 equals the expected code. Used for codes scenarios.go
// derives from (Type, $tags).
func assertSACodeViaScenario(t *testing.T, expected cbc.Code, tags ...cbc.Key) {
	t.Helper()
	inv := validStandardInvoice()
	if len(tags) > 0 {
		inv.SetTags(tags...)
	}
	if hasExportBit(expected) {
		clearCustomerVAT(inv)
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
	assert.Equal(t, expected, inv.Tax.Ext.Get(zatca.ExtKeyInvoiceTypeTransactions))
}

// assertSACodeHardCoded builds an "other"-typed invoice — which bypasses
// every regime/addon scenario per bill.InvoiceTypeOther's contract — then
// hardcodes UNTDID = 388 and the target KSA-2 code on inv.Tax.Ext, runs
// Calculate (no scenario overwrites the preset values), and validates.
// Used for the codes in validTransactionTypes that scenarios.go does not
// cover.
func assertSACodeHardCoded(t *testing.T, code cbc.Code) {
	t.Helper()
	inv := validStandardInvoice()
	inv.Type = bill.InvoiceTypeOther
	if hasExportBit(code) {
		clearCustomerVAT(inv)
	}
	inv.Tax.Ext = inv.Tax.Ext.
		Set(untdid.ExtKeyDocumentType, "388").
		Set(zatca.ExtKeyInvoiceTypeTransactions, code)
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestAllValidTransactionTypes(t *testing.T) {
	// --- Scenario-driven codes (via $tags) ---

	t.Run("0100000 standard default", func(t *testing.T) {
		assertSACodeViaScenario(t, "0100000")
	})
	t.Run("0100010 standard+summary", func(t *testing.T) {
		assertSACodeViaScenario(t, "0100010", zatca.TagSummary)
	})
	t.Run("0100100 standard+export", func(t *testing.T) {
		assertSACodeViaScenario(t, "0100100", tax.TagExport)
	})
	t.Run("0200000 simplified default", func(t *testing.T) {
		assertSACodeViaScenario(t, "0200000", tax.TagSimplified)
	})
	t.Run("0200010 simplified+summary", func(t *testing.T) {
		assertSACodeViaScenario(t, "0200010", tax.TagSimplified, zatca.TagSummary)
	})

	// --- Hard-coded Standard codes (TT=01) not covered by scenarios.go ---

	t.Run("0100001 standard+self-billed", func(t *testing.T) {
		assertSACodeHardCoded(t, "0100001")
	})
	t.Run("0100011 standard+summary+self-billed", func(t *testing.T) {
		assertSACodeHardCoded(t, "0100011")
	})
	t.Run("0100110 standard+export+summary", func(t *testing.T) {
		assertSACodeHardCoded(t, "0100110")
	})
	t.Run("0101000 standard+nominal", func(t *testing.T) {
		assertSACodeHardCoded(t, "0101000")
	})
	t.Run("0101001 standard+nominal+self-billed", func(t *testing.T) {
		assertSACodeHardCoded(t, "0101001")
	})
	t.Run("0101010 standard+nominal+summary", func(t *testing.T) {
		assertSACodeHardCoded(t, "0101010")
	})
	t.Run("0101011 standard+nominal+summary+self-billed", func(t *testing.T) {
		assertSACodeHardCoded(t, "0101011")
	})
	t.Run("0101100 standard+nominal+export", func(t *testing.T) {
		assertSACodeHardCoded(t, "0101100")
	})
	t.Run("0101110 standard+nominal+export+summary", func(t *testing.T) {
		assertSACodeHardCoded(t, "0101110")
	})
	t.Run("0110000 standard+third-party", func(t *testing.T) {
		assertSACodeHardCoded(t, "0110000")
	})
	t.Run("0110001 standard+third-party+self-billed", func(t *testing.T) {
		assertSACodeHardCoded(t, "0110001")
	})
	t.Run("0110010 standard+third-party+summary", func(t *testing.T) {
		assertSACodeHardCoded(t, "0110010")
	})
	t.Run("0110011 standard+third-party+summary+self-billed", func(t *testing.T) {
		assertSACodeHardCoded(t, "0110011")
	})
	t.Run("0110100 standard+third-party+export", func(t *testing.T) {
		assertSACodeHardCoded(t, "0110100")
	})
	t.Run("0110110 standard+third-party+export+summary", func(t *testing.T) {
		assertSACodeHardCoded(t, "0110110")
	})
	t.Run("0111000 standard+third-party+nominal", func(t *testing.T) {
		assertSACodeHardCoded(t, "0111000")
	})
	t.Run("0111001 standard+third-party+nominal+self-billed", func(t *testing.T) {
		assertSACodeHardCoded(t, "0111001")
	})
	t.Run("0111010 standard+third-party+nominal+summary", func(t *testing.T) {
		assertSACodeHardCoded(t, "0111010")
	})
	t.Run("0111011 standard+third-party+nominal+summary+self-billed", func(t *testing.T) {
		assertSACodeHardCoded(t, "0111011")
	})
	t.Run("0111100 standard+third-party+nominal+export", func(t *testing.T) {
		assertSACodeHardCoded(t, "0111100")
	})
	t.Run("0111110 standard+third-party+nominal+export+summary", func(t *testing.T) {
		assertSACodeHardCoded(t, "0111110")
	})

	// --- Hard-coded Simplified codes (TT=02) not covered by scenarios.go ---

	t.Run("0201000 simplified+nominal", func(t *testing.T) {
		assertSACodeHardCoded(t, "0201000")
	})
	t.Run("0201010 simplified+nominal+summary", func(t *testing.T) {
		assertSACodeHardCoded(t, "0201010")
	})
	t.Run("0210000 simplified+third-party", func(t *testing.T) {
		assertSACodeHardCoded(t, "0210000")
	})
	t.Run("0210010 simplified+third-party+summary", func(t *testing.T) {
		assertSACodeHardCoded(t, "0210010")
	})
	t.Run("0211000 simplified+third-party+nominal", func(t *testing.T) {
		assertSACodeHardCoded(t, "0211000")
	})
	t.Run("0211010 simplified+third-party+nominal+summary", func(t *testing.T) {
		assertSACodeHardCoded(t, "0211010")
	})
}
