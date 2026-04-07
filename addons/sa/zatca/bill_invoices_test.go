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
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- Test fixtures ---

// validStandardInvoice returns a fully valid standard tax invoice (KSA-2 = 01).
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
			Ext: tax.Extensions{
				zatca.ExtKeyInvoiceTypeTransactions: "0100000",
				untdid.ExtKeyDocumentType:           "388",
			},
		},
		Supplier: &org.Party{
			Name: "Acme Corp Saudi",
			TaxID: &tax.Identity{
				Country: "SA",
				Code:    "300000000000003",
			},
			Addresses: []*org.Address{
				{
					Street:      "King Fahd Road",
					StreetExtra: "Al Olaya",
					Number:      "1234",
					Locality:    "Riyadh",
					Code:        "12345",
					Country:     "SA",
				},
			},
		},
		Customer: &org.Party{
			Name: "Sample Consumer LLC",
			TaxID: &tax.Identity{
				Country: "SA",
				Code:    "399999999900003",
			},
			Addresses: []*org.Address{
				{
					Street:      "Olaya Street",
					StreetExtra: "Al Malaz",
					Number:      "5678",
					Locality:    "Riyadh",
					Code:        "54321",
					Country:     "SA",
				},
			},
		},
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
						Category: "VAT",
						Key:      tax.KeyStandard,
					},
				},
			},
		},
	}
}

// validSimplifiedInvoice returns a fully valid simplified tax invoice (KSA-2 = 02).
func validSimplifiedInvoice() *bill.Invoice {
	inv := validStandardInvoice()
	inv.Tax.Ext[zatca.ExtKeyInvoiceTypeTransactions] = "0200000"
	return inv
}

// validCreditNote returns a fully valid standard credit note (381).
func validCreditNote() *bill.Invoice {
	inv := validStandardInvoice()
	inv.Type = bill.InvoiceTypeCreditNote
	inv.Tax.Ext[untdid.ExtKeyDocumentType] = "381"
	inv.Preceding = []*org.DocumentRef{
		{
			Code:      "INV-001",
			IssueDate: cal.NewDate(2024, 5, 1),
			Reason:    "Return of goods",
		},
	}
	return inv
}

// validDebitNote returns a fully valid standard debit note (383).
func validDebitNote() *bill.Invoice {
	inv := validStandardInvoice()
	inv.Type = bill.InvoiceTypeDebitNote
	inv.Tax.Ext[untdid.ExtKeyDocumentType] = "383"
	inv.Preceding = []*org.DocumentRef{
		{
			Code:      "INV-001",
			IssueDate: cal.NewDate(2024, 5, 1),
			Reason:    "Price adjustment",
		},
	}
	return inv
}

// validExportInvoice returns a valid export invoice (KSA-2 position 5 = 1) without buyer VAT.
func validExportInvoice() *bill.Invoice {
	inv := validStandardInvoice()
	inv.Tax.Ext[zatca.ExtKeyInvoiceTypeTransactions] = "0100100"
	inv.Customer.TaxID = nil
	inv.Customer.Identities = []*org.Identity{
		{Type: "TIN", Code: "123456789012345"},
	}
	return inv
}

// validSimplifiedSummaryInvoice returns a valid simplified summary invoice (KSA-2 = 02, position 6 = 1).
func validSimplifiedSummaryInvoice() *bill.Invoice {
	inv := validSimplifiedInvoice()
	inv.Tax.Ext[zatca.ExtKeyInvoiceTypeTransactions] = "0200010"
	return inv
}

func calculated(t *testing.T, inv *bill.Invoice) *bill.Invoice {
	t.Helper()
	require.NoError(t, inv.Calculate())
	return inv
}

// --- Fixture validation ---

func TestValidInvoices(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified invoice", func(t *testing.T) {
		inv := calculated(t, validSimplifiedInvoice())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("credit note", func(t *testing.T) {
		inv := calculated(t, validCreditNote())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("debit note", func(t *testing.T) {
		inv := calculated(t, validDebitNote())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("export invoice", func(t *testing.T) {
		inv := calculated(t, validExportInvoice())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified summary invoice", func(t *testing.T) {
		inv := calculated(t, validSimplifiedSummaryInvoice())
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-70 (rule 01): Invoice must contain issue time (KSA-25) ---

func TestBRKSA70_IssueTime(t *testing.T) {
	t.Run("normalizer auto-creates issue time when nil", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.IssueTime = nil
		require.NoError(t, inv.Calculate())
		assert.NotNil(t, inv.IssueTime, "normalizer should auto-create issue time")
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-05 (rule 02): Document type must be a valid ZATCA type (388, 386, 383, 381) ---

func TestBRKSA05_DocumentType(t *testing.T) {
	t.Run("valid document types", func(t *testing.T) {
		for _, code := range []cbc.Code{"388", "386", "383", "381"} {
			t.Run(string(code), func(t *testing.T) {
				inv := validStandardInvoice()
				inv.Tax.Ext[untdid.ExtKeyDocumentType] = code
				switch code {
				case "381":
					inv.Type = bill.InvoiceTypeCreditNote
					inv.Preceding = []*org.DocumentRef{
						{Code: "INV-001", IssueDate: cal.NewDate(2024, 5, 1), Reason: "Return"},
					}
				case "383":
					inv.Type = bill.InvoiceTypeDebitNote
					inv.Preceding = []*org.DocumentRef{
						{Code: "INV-001", IssueDate: cal.NewDate(2024, 5, 1), Reason: "Adjustment"},
					}
				case "386":
					inv.Tags = tax.WithTags(tax.TagPrepayment)
				}
				require.NoError(t, inv.Calculate())
				require.NoError(t, rules.Validate(inv))
			})
		}
	})

	t.Run("invalid document type", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "999"
		assert.ErrorContains(t, rules.Validate(inv), "document type must be a valid ZATCA type")
	})

	t.Run("empty document type", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = ""
		assert.ErrorContains(t, rules.Validate(inv), "document type must be a valid ZATCA type")
	})
}

// --- BR-KSA-56 (rules 03, 04): Credit/debit notes must have billing reference with identifier ---

func TestBRKSA56_BillingReference(t *testing.T) {
	t.Run("credit note missing preceding", func(t *testing.T) {
		inv := validCreditNote()
		inv.Preceding = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "credit and debit notes must have a billing reference")
	})

	t.Run("debit note missing preceding", func(t *testing.T) {
		inv := validDebitNote()
		inv.Preceding = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "credit and debit notes must have a billing reference")
	})

	t.Run("credit note missing preceding code", func(t *testing.T) {
		inv := validCreditNote()
		inv.Preceding[0].Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "billing reference must have an identifier (BR-KSA-56)")
	})

	t.Run("debit note missing preceding code", func(t *testing.T) {
		inv := validDebitNote()
		inv.Preceding[0].Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "billing reference must have an identifier (BR-KSA-56)")
	})

	t.Run("standard invoice does not require preceding", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-17 (rule 05): Credit/debit notes must contain reason for issuance ---

func TestBRKSA17_CreditDebitNoteReason(t *testing.T) {
	t.Run("credit note missing reason", func(t *testing.T) {
		inv := validCreditNote()
		inv.Preceding[0].Reason = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "credit and debit notes must contain the reason for issuance")
	})

	t.Run("debit note missing reason", func(t *testing.T) {
		inv := validDebitNote()
		inv.Preceding[0].Reason = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "credit and debit notes must contain the reason for issuance")
	})

	t.Run("credit note with reason is valid", func(t *testing.T) {
		inv := calculated(t, validCreditNote())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("debit note with reason is valid", func(t *testing.T) {
		inv := calculated(t, validDebitNote())
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-09 (rules 06-15): Supplier must be present with valid address ---
// Including BR-KSA-37 (building number 4 digits) and BR-KSA-66 (postal code 5 digits).

func TestBRKSA09_SupplierAddress(t *testing.T) {
	t.Run("missing supplier", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier = nil
		assert.ErrorContains(t, rules.Validate(inv), "supplier is required")
	})

	t.Run("missing supplier address", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.Addresses = nil
		assert.ErrorContains(t, rules.Validate(inv), "supplier address is required")
	})

	// BR-KSA-09: Building number (KSA-17) required
	t.Run("missing building number", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.Addresses[0].Number = ""
		assert.ErrorContains(t, rules.Validate(inv), "supplier address building number is required")
	})

	// BR-KSA-37: Building number must be exactly 4 digits
	t.Run("building number 3 digits", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.Addresses[0].Number = "123"
		assert.ErrorContains(t, rules.Validate(inv), "supplier address building number must contain 4 digits")
	})

	t.Run("building number 5 digits", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.Addresses[0].Number = "12345"
		assert.ErrorContains(t, rules.Validate(inv), "supplier address building number must contain 4 digits")
	})

	t.Run("building number non-numeric", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.Addresses[0].Number = "12AB"
		assert.ErrorContains(t, rules.Validate(inv), "supplier address building number must contain 4 digits")
	})

	t.Run("building number valid 4 digits", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.Addresses[0].Number = "0001"
		require.NoError(t, rules.Validate(inv))
	})

	// BR-KSA-09: District (KSA-3) required
	t.Run("missing district", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.Addresses[0].StreetExtra = ""
		assert.ErrorContains(t, rules.Validate(inv), "supplier address must have a district")
	})

	// BR-KSA-09: Street name (BT-35) required
	t.Run("missing street", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.Addresses[0].Street = ""
		assert.ErrorContains(t, rules.Validate(inv), "supplier address must have a street name")
	})

	// BR-KSA-09: Postal code (BT-38) required
	t.Run("missing postal code", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.Addresses[0].Code = ""
		assert.ErrorContains(t, rules.Validate(inv), "supplier postal code is required")
	})

	// BR-KSA-66: Postal code must be exactly 5 digits
	t.Run("postal code 4 digits", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.Addresses[0].Code = "1234"
		assert.ErrorContains(t, rules.Validate(inv), "supplier postal code must be 5 digits")
	})

	t.Run("postal code 6 digits", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.Addresses[0].Code = "123456"
		assert.ErrorContains(t, rules.Validate(inv), "supplier postal code must be 5 digits")
	})

	t.Run("postal code non-numeric", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.Addresses[0].Code = "1234A"
		assert.ErrorContains(t, rules.Validate(inv), "supplier postal code must be 5 digits")
	})

	// BR-KSA-09: City (BT-37) required
	t.Run("missing locality", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.Addresses[0].Locality = ""
		assert.ErrorContains(t, rules.Validate(inv), "supplier address must have a city name")
	})

	// BR-KSA-09: Country code (BT-40) required
	t.Run("missing country", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.Addresses[0].Country = ""
		assert.ErrorContains(t, rules.Validate(inv), "supplier address must have a country code")
	})

	// Multiple violations reported simultaneously
	t.Run("multiple missing fields", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.Addresses[0].Number = ""
		inv.Supplier.Addresses[0].Street = ""
		inv.Supplier.Addresses[0].Code = ""
		inv.Supplier.Addresses[0].Locality = ""
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier address building number is required")
		assert.ErrorContains(t, err, "supplier address must have a street name")
		assert.ErrorContains(t, err, "supplier postal code is required")
		assert.ErrorContains(t, err, "supplier address must have a city name")
	})
}

// --- BR-KSA-39, BR-KSA-40 (rules 16, 17): Supplier VAT registration number ---

func TestBRKSA39_40_SupplierVAT(t *testing.T) {
	// BR-KSA-39: Seller VAT registration number is required
	t.Run("missing VAT number", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.TaxID.Code = ""
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "supplier must have a VAT number")
	})

	// BR-KSA-40: If present, must be 15 digits, first and last digit = "3"
	t.Run("VAT number too short", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.TaxID.Code = "31234567890"
		assert.ErrorContains(t, rules.Validate(inv), "supplier VAT number must be 15 digits starting/ending with 3")
	})

	t.Run("VAT number too long", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.TaxID.Code = "3000000000000003"
		assert.ErrorContains(t, rules.Validate(inv), "supplier VAT number must be 15 digits starting/ending with 3")
	})

	t.Run("VAT number not starting with 3", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.TaxID.Code = "100000000000003"
		assert.ErrorContains(t, rules.Validate(inv), "supplier VAT number must be 15 digits starting/ending with 3")
	})

	t.Run("VAT number not ending with 3", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.TaxID.Code = "300000000000001"
		assert.ErrorContains(t, rules.Validate(inv), "supplier VAT number must be 15 digits starting/ending with 3")
	})

	t.Run("VAT number with non-digit characters", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Supplier.TaxID.Code = "30000000000A003"
		assert.ErrorContains(t, rules.Validate(inv), "supplier VAT number must be 15 digits starting/ending with 3")
	})

	t.Run("valid VAT number", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("valid VAT number edge case 300000000000003", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.TaxID.Code = "300000000000003"
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("valid VAT number edge case 399999999999993", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.TaxID.Code = "399999999999993"
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-08 (rule 18): Supplier must have at most one identity ---

func TestBRKSA08_SupplierIdentity(t *testing.T) {
	t.Run("supplier with one identity is valid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Type: "CRN", Code: "1234567890"},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("supplier with multiple identities is invalid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Supplier.Identities = []*org.Identity{
			{Type: "CRN", Code: "1234567890"},
			{Type: "MOM", Code: "0987654321"},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "supplier identification must be valid")
	})

	t.Run("supplier without identities is valid", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		require.NoError(t, rules.Validate(inv))
	})

	// Test each valid scheme ID individually
	for _, schemeID := range []cbc.Code{"CRN", "MOM", "MLS", "700", "SAG", "OTH"} {
		t.Run("supplier with "+string(schemeID)+" identity", func(t *testing.T) {
			inv := validStandardInvoice()
			inv.Supplier.Identities = []*org.Identity{
				{Type: schemeID, Code: "1234567890"},
			}
			require.NoError(t, inv.Calculate())
			require.NoError(t, rules.Validate(inv))
		})
	}
}

// --- BR-KSA-42 (rules 19, 20): Standard invoices must have customer present with name ---

func TestBRKSA42_StandardBuyerName(t *testing.T) {
	t.Run("standard invoice missing customer", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer = nil
		assert.ErrorContains(t, rules.Validate(inv), "customer must be present")
	})

	t.Run("standard invoice missing buyer name", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer.Name = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer name must be present in the standard tax invoice")
	})

	t.Run("standard credit note missing buyer name", func(t *testing.T) {
		inv := calculated(t, validCreditNote())
		inv.Customer.Name = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer name must be present in the standard tax invoice")
	})

	t.Run("standard debit note missing buyer name", func(t *testing.T) {
		inv := calculated(t, validDebitNote())
		inv.Customer.Name = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer name must be present in the standard tax invoice")
	})

	t.Run("simplified invoice does not require buyer name from standard rule", func(t *testing.T) {
		inv := validSimplifiedInvoice()
		inv.Customer.Name = ""
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{Type: "TIN", Code: "123456789012345"},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-10 (rules 21-23): Standard invoices require customer address with street, city, country ---

func TestBRKSA10_StandardCustomerAddress(t *testing.T) {
	t.Run("missing customer street", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer.Addresses[0].Street = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer address must have a street name (BR-KSA-10)")
	})

	t.Run("missing customer locality", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer.Addresses[0].Locality = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer address must have a city name (BR-KSA-10)")
	})

	t.Run("missing customer country", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer.Addresses[0].Country = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer address must have a country code (BR-KSA-10)")
	})

	t.Run("simplified invoice does not require customer address fields", func(t *testing.T) {
		inv := validSimplifiedInvoice()
		inv.Customer.Addresses[0].Street = ""
		inv.Customer.Addresses[0].Country = "US"
		inv.Customer.Addresses[0].StreetExtra = ""
		inv.Customer.Addresses[0].Number = ""
		inv.Customer.Addresses[0].Code = ""
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	// Standard credit note also requires customer address
	t.Run("standard credit note missing customer street", func(t *testing.T) {
		inv := calculated(t, validCreditNote())
		inv.Customer.Addresses[0].Street = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer address must have a street name (BR-KSA-10)")
	})
}

// --- BR-KSA-52, BR-KSA-53 (rules 24-26): Standard invoices must have lines with taxes and totals ---

func TestBRKSA52_53_StandardInvoiceLines(t *testing.T) {
	// BR-KSA-52: Line taxes required for standard invoices
	t.Run("missing line taxes", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Lines[0].Taxes = nil
		assert.ErrorContains(t, rules.Validate(inv), "line taxes are required for standard tax invoices")
	})

	// BR-KSA-53: Line total required for standard invoices
	t.Run("missing line total", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Lines[0].Total = nil
		assert.ErrorContains(t, rules.Validate(inv), "line total amount is required for standard tax invoices")
	})

	// Standard credit note also requires line taxes
	t.Run("standard credit note missing line taxes", func(t *testing.T) {
		inv := calculated(t, validCreditNote())
		inv.Lines[0].Taxes = nil
		assert.ErrorContains(t, rules.Validate(inv), "line taxes are required for standard tax invoices")
	})

	// Simplified invoices don't have this requirement
	t.Run("simplified invoice does not require line taxes", func(t *testing.T) {
		inv := validSimplifiedInvoice()
		inv.Lines[0].Taxes = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

// --- Rule 27: Customer must be present ---

func TestCustomerPresent(t *testing.T) {
	t.Run("missing customer on standard invoice", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer = nil
		assert.ErrorContains(t, rules.Validate(inv), "customer must be present")
	})

	t.Run("missing customer on simplified invoice", func(t *testing.T) {
		inv := calculated(t, validSimplifiedInvoice())
		inv.Customer = nil
		assert.ErrorContains(t, rules.Validate(inv), "customer must be present")
	})
}

// --- BR-KSA-63 (rules 28-32): SA customer address must have street, building number,
// postal code (5 digits), city, and district ---

func TestBRKSA63_SACustomerAddress(t *testing.T) {
	t.Run("SA buyer missing street", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer.Addresses[0].Street = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer address in SA must have a street name (BR-KSA-63)")
	})

	t.Run("SA buyer missing building number", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer.Addresses[0].Number = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer address in SA must have a building number (BR-KSA-63)")
	})

	// BR-KSA-67: SA buyer postal code must be exactly 5 digits
	t.Run("SA buyer postal code 4 digits", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer.Addresses[0].Code = "1234"
		assert.ErrorContains(t, rules.Validate(inv), "customer address in SA must have a postal code (BR-KSA-63)")
	})

	t.Run("SA buyer postal code 6 digits", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer.Addresses[0].Code = "123456"
		assert.ErrorContains(t, rules.Validate(inv), "customer address in SA must have a postal code (BR-KSA-63)")
	})

	t.Run("SA buyer postal code non-numeric", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer.Addresses[0].Code = "1234A"
		assert.ErrorContains(t, rules.Validate(inv), "customer address in SA must have a postal code (BR-KSA-63)")
	})

	t.Run("SA buyer missing locality", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer.Addresses[0].Locality = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer address in SA must have a city name (BR-KSA-63)")
	})

	t.Run("SA buyer missing district", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer.Addresses[0].StreetExtra = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer address in SA must have a district name (BR-KSA-63)")
	})

	// Non-SA country does not require SA-specific fields
	t.Run("non-SA buyer does not require SA-specific fields", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Addresses[0].Country = "US"
		inv.Customer.Addresses[0].StreetExtra = ""
		inv.Customer.Addresses[0].Number = ""
		inv.Customer.Addresses[0].Code = ""
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	// Multiple SA violations reported simultaneously
	t.Run("multiple SA-specific violations", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer.Addresses[0].Number = ""
		inv.Customer.Addresses[0].StreetExtra = ""
		inv.Customer.Addresses[0].Code = "123"
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer address in SA must have a building number")
		assert.ErrorContains(t, err, "customer address in SA must have a district name")
		assert.ErrorContains(t, err, "customer address in SA must have a postal code")
	})
}

// --- BR-KSA-14 (rule 33): Buyer must be VAT registered or have a valid identification ---

func TestBRKSA14_CustomerIdentification(t *testing.T) {
	t.Run("customer with VAT is valid without identity", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("customer without VAT but with one identity is valid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{Type: "TIN", Code: "123456789012345"},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	// Test each valid buyer scheme ID
	for _, schemeID := range []cbc.Code{"TIN", "CRN", "MOM", "MLS", "700", "SAG", "NAT", "GCC", "IQA", "OTH"} {
		t.Run("buyer with "+string(schemeID)+" identity no VAT", func(t *testing.T) {
			inv := validStandardInvoice()
			inv.Customer.TaxID = nil
			inv.Customer.Identities = []*org.Identity{
				{Type: schemeID, Code: "1234567890"},
			}
			require.NoError(t, inv.Calculate())
			require.NoError(t, rules.Validate(inv))
		})
	}

	t.Run("customer without VAT and without identity is invalid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.TaxID = nil
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "buyer identification is valid")
	})

	t.Run("customer with empty VAT code and no identity is invalid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.TaxID = &tax.Identity{Country: "SA", Code: ""}
		inv.Customer.Identities = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "buyer identification is valid")
	})

	t.Run("customer without VAT and with multiple identities is invalid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{Type: "TIN", Code: "123456789012345"},
			{Type: "CRN", Code: "1234567890"},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "buyer identification is valid")
	})
}

// --- BR-KSA-46 (rule 34): Export invoices must not have buyer VAT registration number ---

func TestBRKSA46_ExportBuyerVAT(t *testing.T) {
	t.Run("export with buyer VAT is invalid", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Tax.Ext[zatca.ExtKeyInvoiceTypeTransactions] = "0100100"
		assert.ErrorContains(t, rules.Validate(inv), "export invoices must not have buyer VAT registration number")
	})

	t.Run("export without buyer VAT is valid", func(t *testing.T) {
		inv := calculated(t, validExportInvoice())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("non-export with buyer VAT is valid", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("export credit note without buyer VAT is valid", func(t *testing.T) {
		inv := validCreditNote()
		inv.Tax.Ext[zatca.ExtKeyInvoiceTypeTransactions] = "0100100"
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{Type: "TIN", Code: "123456789012345"},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("export credit note with buyer VAT is invalid", func(t *testing.T) {
		inv := validCreditNote()
		inv.Tax.Ext[zatca.ExtKeyInvoiceTypeTransactions] = "0100100"
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "export invoices must not have buyer VAT registration number")
	})
}

// --- BR-KSA-44 (rules 35, 36): Non-export invoices must have customer with valid VAT ---

func TestBRKSA44_NonExportBuyerVAT(t *testing.T) {
	t.Run("non-export customer must be present", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer = nil
		assert.ErrorContains(t, rules.Validate(inv), "customer must be present")
	})

	t.Run("non-export buyer VAT valid format", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("non-export buyer VAT invalid format too short", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer.TaxID.Code = "31234567890"
		assert.ErrorContains(t, rules.Validate(inv), "VAT numbers should be valid (BR-KSA-44)")
	})

	t.Run("non-export buyer VAT not starting with 3", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer.TaxID.Code = "100000000000003"
		assert.ErrorContains(t, rules.Validate(inv), "VAT numbers should be valid (BR-KSA-44)")
	})

	t.Run("non-export buyer VAT not ending with 3", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Customer.TaxID.Code = "300000000000001"
		assert.ErrorContains(t, rules.Validate(inv), "VAT numbers should be valid (BR-KSA-44)")
	})

	t.Run("non-export buyer empty VAT passes format check", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{Type: "TIN", Code: "123456789012345"},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("export invoice skips buyer VAT format check", func(t *testing.T) {
		inv := calculated(t, validExportInvoice())
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-15 (rules 37, 38): Standard invoices must have delivery with supply date ---

func TestBRKSA15_SupplyDate(t *testing.T) {
	t.Run("standard invoice missing delivery", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Delivery = nil
		assert.ErrorContains(t, rules.Validate(inv), "delivery must be present")
	})

	t.Run("standard invoice missing supply date", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Delivery.Date = nil
		assert.ErrorContains(t, rules.Validate(inv), "delivery period must have a supply date")
	})

	t.Run("standard invoice does not require delivery period", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Delivery.Period = nil
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified non-summary invoice does not require delivery", func(t *testing.T) {
		inv := validSimplifiedInvoice()
		inv.Delivery = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	// Standard credit note also requires delivery
	t.Run("standard credit note missing delivery", func(t *testing.T) {
		inv := calculated(t, validCreditNote())
		inv.Delivery = nil
		assert.ErrorContains(t, rules.Validate(inv), "delivery must be present")
	})
}

// --- BR-KSA-72 (rules 39-41): Simplified summary invoice delivery period ---

func TestBRKSA72_SimplifiedSummaryDelivery(t *testing.T) {
	// Rule 39: delivery period must be present
	t.Run("simplified summary requires delivery period", func(t *testing.T) {
		inv := validSimplifiedSummaryInvoice()
		inv.Delivery.Period = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "supply must have a delivery period")
	})

	// Rules 40, 41: delivery start and end dates must be present (BR-KSA-72).
	// Core cal.Period validation catches zero dates before ZATCA rules fire.
	t.Run("simplified summary requires delivery start date", func(t *testing.T) {
		inv := validSimplifiedSummaryInvoice()
		inv.Delivery.Period = &cal.Period{
			End: cal.MakeDate(2024, 6, 30),
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "start date cannot be zero")
	})

	t.Run("simplified summary requires delivery end date", func(t *testing.T) {
		inv := validSimplifiedSummaryInvoice()
		inv.Delivery.Period = &cal.Period{
			Start: cal.MakeDate(2024, 6, 1),
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "end date cannot be zero")
	})

	t.Run("simplified summary requires delivery", func(t *testing.T) {
		inv := validSimplifiedSummaryInvoice()
		inv.Delivery = nil
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "delivery must be present")
	})

	t.Run("simplified summary with valid delivery period", func(t *testing.T) {
		inv := calculated(t, validSimplifiedSummaryInvoice())
		require.NoError(t, rules.Validate(inv))
	})

	// Standard summary does NOT require period start/end (only simplified summary)
	t.Run("standard summary does not require delivery period", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Tax.Ext[zatca.ExtKeyInvoiceTypeTransactions] = "0100010" // standard + summary
		inv.Delivery.Period = nil
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-36: Supply end date must be >= supply start date ---

func TestBRKSA36_SupplyEndDate(t *testing.T) {
	t.Run("end date equals start date is valid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Delivery.Date = cal.NewDate(2024, 6, 15)
		inv.Delivery.Period = &cal.Period{
			Start: cal.MakeDate(2024, 6, 15),
			End:   cal.MakeDate(2024, 6, 15),
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("end date after start date is valid", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("end date before start date is invalid", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		inv.Delivery.Period = &cal.Period{
			Start: cal.MakeDate(2024, 6, 1),
			End:   cal.MakeDate(2024, 1, 15),
		}
		assert.ErrorContains(t, rules.Validate(inv), "[GOBL-CAL-PERIOD-10]")
	})
}

// --- BR-KSA-71 (rules 42, 43): Simplified summary invoices must have customer with name ---

func TestBRKSA71_SimplifiedSummaryCustomer(t *testing.T) {
	t.Run("simplified summary missing customer name", func(t *testing.T) {
		inv := validSimplifiedSummaryInvoice()
		inv.Customer.Name = ""
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{Type: "TIN", Code: "123456789012345"},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "customer name must be present for simplified, summary invoices (BR-KSA-71)")
	})

	t.Run("simplified summary with customer name is valid", func(t *testing.T) {
		inv := calculated(t, validSimplifiedSummaryInvoice())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified non-summary does not require customer name", func(t *testing.T) {
		inv := validSimplifiedInvoice()
		inv.Customer.Name = ""
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{Type: "TIN", Code: "123456789012345"},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified summary missing customer", func(t *testing.T) {
		inv := calculated(t, validSimplifiedSummaryInvoice())
		inv.Customer = nil
		assert.ErrorContains(t, rules.Validate(inv), "customer must be present for simplified, summary invoices")
	})
}

// --- BR-KSA-CL-04 (tax combo rule 01): Exempt, zero-rated, and outside-scope tax combos
// must have a valid SA VATEX exemption code ---
//
// Note: EN16931 rule BR-Z-10 prevents setting VATEX extension on zero-rated categories.
// Tests for valid zero-rated VATEX codes use tax.KeyExempt as a workaround. The ZATCA
// spec maps these to "Z" (zero-rated), but in the GOBL model with EN16931, the VATEX
// code itself carries the exemption/zero-rated semantic.

func TestBRKSACL04_VATEXCodes(t *testing.T) {
	t.Run("zero-rated without VATEX is invalid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyZero,
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "must have a valid SA VATEX code")
	})

	t.Run("exempt without VATEX is invalid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "must have a valid SA VATEX code")
	})

	t.Run("exempt with non-SA VATEX is invalid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-EU-79-C",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "must have a valid SA VATEX code")
	})

	t.Run("outside-scope without VATEX is invalid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyOutsideScope,
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "must have a valid SA VATEX code")
	})

	t.Run("standard rate does not require VATEX", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("invalid VATEX code", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-INVALID",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "must have a valid SA VATEX code")
	})

	// Valid exempt VATEX codes (ZATCA spec: E category)
	t.Run("valid exempt VATEX codes", func(t *testing.T) {
		exemptCodes := []cbc.Code{
			"VATEX-SA-29", "VATEX-SA-29-7", "VATEX-SA-30",
		}
		for _, code := range exemptCodes {
			t.Run(string(code), func(t *testing.T) {
				inv := validStandardInvoice()
				inv.Lines[0].Taxes = tax.Set{
					{
						Category: "VAT",
						Key:      tax.KeyExempt,
						Ext: tax.Extensions{
							cef.ExtKeyVATEX: code,
						},
					},
				}
				require.NoError(t, inv.Calculate())
				require.NoError(t, rules.Validate(inv))
			})
		}
	})

	// Valid outside-scope VATEX code (ZATCA spec: O category)
	t.Run("valid outside-scope VATEX code", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyOutsideScope,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-OOS",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	// Valid zero-rated VATEX codes (ZATCA spec: Z category).
	// Note: Uses tax.KeyExempt because EN16931 BR-Z-10 prevents VATEX on zero-rated.
	// The ZATCA-specific VATEX code carries the zero-rated semantic.
	t.Run("valid zero-rated SA VATEX codes via exempt key", func(t *testing.T) {
		zeroCodes := []cbc.Code{
			"VATEX-SA-32", "VATEX-SA-33",
			"VATEX-SA-34-1", "VATEX-SA-34-2", "VATEX-SA-34-3",
			"VATEX-SA-34-4", "VATEX-SA-34-5",
			"VATEX-SA-35", "VATEX-SA-36",
			"VATEX-SA-EDU", "VATEX-SA-HEA", "VATEX-SA-MLTRY",
		}
		for _, code := range zeroCodes {
			t.Run(string(code), func(t *testing.T) {
				inv := validStandardInvoice()
				inv.Lines[0].Taxes = tax.Set{
					{
						Category: "VAT",
						Key:      tax.KeyExempt,
						Ext: tax.Extensions{
							cef.ExtKeyVATEX: code,
						},
					},
				}
				if code == "VATEX-SA-EDU" || code == "VATEX-SA-HEA" {
					inv.Customer.Identities = []*org.Identity{
						{Type: "NAT", Code: "1234567890"},
					}
				}
				require.NoError(t, inv.Calculate())
				require.NoError(t, rules.Validate(inv))
			})
		}
	})
}

// --- BR-KSA-CL-04 (tax combo rule 02): Standard rate must NOT have VATEX code ---

func TestBRKSACL04_StandardNoVATEX(t *testing.T) {
	t.Run("standard rate with VATEX is invalid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyStandard,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-32",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "standard rate tax must not have a VATEX code")
	})

	t.Run("standard rate without VATEX is valid", func(t *testing.T) {
		inv := calculated(t, validStandardInvoice())
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-49 (rule 44): EDU/HEA exemptions require customer NAT identity ---

func TestBRKSA49_EDUHEARequiresNAT(t *testing.T) {
	// Uses tax.KeyExempt to avoid EN16931 BR-Z-10 conflict with VATEX on zero-rated.
	t.Run("EDU without NAT identity is invalid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-EDU",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "customer must have a national ID (NAT)")
	})

	t.Run("HEA without NAT identity is invalid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-HEA",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "customer must have a national ID (NAT)")
	})

	t.Run("EDU with NAT identity is valid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Identities = []*org.Identity{
			{Type: "NAT", Code: "1234567890"},
		}
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-EDU",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("HEA with NAT identity is valid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Identities = []*org.Identity{
			{Type: "NAT", Code: "1234567890"},
		}
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-HEA",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("EDU with non-NAT identity is invalid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Identities = []*org.Identity{
			{Type: "CRN", Code: "1234567890"},
		}
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-EDU",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "customer must have a national ID (NAT)")
	})

	t.Run("non-EDU/HEA exemption does not require NAT", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-32",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-25 (rules 45, 46): Simplified with EDU/HEA exemption requires customer name ---

func TestBRKSA25_SimplifiedEDUHEACustomerName(t *testing.T) {
	t.Run("simplified with EDU missing customer name", func(t *testing.T) {
		inv := validSimplifiedInvoice()
		inv.Customer.Name = ""
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{Type: "NAT", Code: "1234567890"},
		}
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-EDU",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "customer name must be present (BR-KSA-25)")
	})

	t.Run("simplified with HEA missing customer name", func(t *testing.T) {
		inv := validSimplifiedInvoice()
		inv.Customer.Name = ""
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{Type: "NAT", Code: "1234567890"},
		}
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-HEA",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "customer name must be present (BR-KSA-25)")
	})

	t.Run("simplified with EDU and customer name is valid", func(t *testing.T) {
		inv := validSimplifiedInvoice()
		inv.Customer.Identities = []*org.Identity{
			{Type: "NAT", Code: "1234567890"},
		}
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-EDU",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified credit note with EDU missing customer name", func(t *testing.T) {
		inv := validCreditNote()
		inv.Tax.Ext[zatca.ExtKeyInvoiceTypeTransactions] = "0200000"
		inv.Customer.Name = ""
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{Type: "NAT", Code: "1234567890"},
		}
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-EDU",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		assert.ErrorContains(t, rules.Validate(inv), "customer name must be present (BR-KSA-25)")
	})

	t.Run("simplified credit note with HEA and customer name is valid", func(t *testing.T) {
		inv := validCreditNote()
		inv.Tax.Ext[zatca.ExtKeyInvoiceTypeTransactions] = "0200000"
		inv.Customer.Identities = []*org.Identity{
			{Type: "NAT", Code: "1234567890"},
		}
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-HEA",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	// Standard invoice with EDU does NOT trigger BR-KSA-25 (only simplified)
	t.Run("standard with EDU does not trigger BR-KSA-25", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Identities = []*org.Identity{
			{Type: "NAT", Code: "1234567890"},
		}
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-EDU",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

// --- Combined / integration scenarios ---

func TestCombined_StandardCreditNoteRequirements(t *testing.T) {
	t.Run("standard credit note requires both buyer name and reason", func(t *testing.T) {
		inv := validCreditNote()
		inv.Customer.Name = ""
		inv.Preceding[0].Reason = ""
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "customer name must be present in the standard tax invoice")
		assert.ErrorContains(t, err, "credit and debit notes must contain the reason for issuance")
	})
}

func TestCombined_SupplierAddressAndVATEX(t *testing.T) {
	t.Run("supplier address and VATEX errors reported together", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Key:      tax.KeyExempt,
			},
		}
		require.NoError(t, inv.Calculate())
		inv.Supplier.Addresses[0].Street = ""
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "supplier address must have a street name")
		assert.ErrorContains(t, err, "must have a valid SA VATEX code")
	})
}

// --- Invoice type codes ---

func TestInvoiceTypeCodes(t *testing.T) {
	t.Run("all standard type codes calculate", func(t *testing.T) {
		for _, code := range zatca.InvTypesStandard {
			t.Run(string(code), func(t *testing.T) {
				inv := validStandardInvoice()
				inv.Tax.Ext[zatca.ExtKeyInvoiceTypeTransactions] = code
				if len(code) == 7 && code[4] == '1' {
					// Export: remove buyer VAT
					inv.Customer.TaxID = nil
					inv.Customer.Identities = []*org.Identity{
						{Type: "TIN", Code: "123456789012345"},
					}
				}
				require.NoError(t, inv.Calculate())
			})
		}
	})

	t.Run("all simplified type codes calculate", func(t *testing.T) {
		for _, code := range zatca.InvTypesSimplified {
			t.Run(string(code), func(t *testing.T) {
				inv := validSimplifiedInvoice()
				inv.Tax.Ext[zatca.ExtKeyInvoiceTypeTransactions] = code
				require.NoError(t, inv.Calculate())
			})
		}
	})
}

// --- Addon registration ---

func TestAddonRegistration(t *testing.T) {
	t.Run("addon exists", func(t *testing.T) {
		ad := tax.AddonForKey(zatca.V1)
		require.NotNil(t, ad)
		assert.Equal(t, zatca.V1, ad.Key)
	})

	t.Run("addon requires EN16931", func(t *testing.T) {
		ad := tax.AddonForKey(zatca.V1)
		require.NotNil(t, ad)
		assert.Contains(t, ad.Requires, cbc.Key("eu-en16931-v2017"))
	})

	t.Run("addon has extensions", func(t *testing.T) {
		ad := tax.AddonForKey(zatca.V1)
		require.NotNil(t, ad)
		assert.NotEmpty(t, ad.Extensions)
	})

	t.Run("addon has scenarios", func(t *testing.T) {
		ad := tax.AddonForKey(zatca.V1)
		require.NotNil(t, ad)
		assert.NotEmpty(t, ad.Scenarios)
	})

	t.Run("addon has normalizer", func(t *testing.T) {
		ad := tax.AddonForKey(zatca.V1)
		require.NotNil(t, ad)
		assert.NotNil(t, ad.Normalizer)
	})
}
