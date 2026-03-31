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

func validStandardInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime:    tax.WithRegime("SA"),
		Addons:    tax.WithAddons(zatca.V1),
		Type:      bill.InvoiceTypeStandard,
		Currency:  "SAR",
		Code:      "INV-001",
		IssueDate: cal.MakeDate(2022, 2, 1),
		IssueTime: cal.NewTime(12, 0, 0),
		Tax: &bill.Tax{
			Ext: tax.Extensions{
				zatca.ExtKeyInvoiceType:   "0100000",
				untdid.ExtKeyDocumentType: "388",
			},
		},
		Supplier: &org.Party{
			Name: "Acme Corp Saudi",
			TaxID: &tax.Identity{
				Country: "SA",
				Code:    "312345678912343",
			},
			Identities: []*org.Identity{
				{
					Code: "1234567890",
					Ext: tax.Extensions{
						zatca.ExtKeySellerIDScheme: "CRN",
					},
				},
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
			Date: cal.NewDate(2022, 2, 1),
			Period: &cal.Period{
				Start: cal.MakeDate(2022, 2, 1),
				End:   cal.MakeDate(2022, 2, 28),
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

func validSimplifiedInvoice() *bill.Invoice {
	inv := validStandardInvoice()
	inv.Tax.Ext[zatca.ExtKeyInvoiceType] = "0200000"
	inv.Customer = &org.Party{
		Name: "Sample Consumer LLC",
		TaxID: &tax.Identity{
			Country: "SA",
			Code:    "399999999900003",
		},
		Addresses: []*org.Address{
			{
				Street:      "Olaya Street",
				Country:     "SA",
				Number:      "5678",
				Code:        "54321",
				Locality:    "Riyadh",
				StreetExtra: "Al Malaz",
			},
		},
	}
	return inv
}

func calculatedInvoice(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := validStandardInvoice()
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
	return inv
}

func calculatedSimplifiedInvoice(t *testing.T) *bill.Invoice {
	t.Helper()
	inv := validSimplifiedInvoice()
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
	return inv
}

func TestValidInvoice(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := calculatedInvoice(t)
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified invoice", func(t *testing.T) {
		inv := calculatedSimplifiedInvoice(t)
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-05: Document type validation ---

func TestDocumentTypeValidation(t *testing.T) {
	t.Run("valid 388 tax invoice", func(t *testing.T) {
		inv := calculatedInvoice(t)
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("valid 386 prepayment", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "386"
		inv.Tags = tax.WithTags(tax.TagPrepayment)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("valid 381 credit note", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "381"
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV-001",
				IssueDate: cal.NewDate(2022, 1, 1),
				Reason:    "Return of goods",
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("valid 383 debit note", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "383"
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV-001",
				IssueDate: cal.NewDate(2022, 1, 1),
				Reason:    "Price adjustment",
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("invalid document type", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "380"
		assert.ErrorContains(t, rules.Validate(inv), "document type must be a valid ZATCA type")
	})
}

// --- BR-KSA-56, BR-KSA-17: Credit/debit note rules ---

func TestCreditDebitNoteValidation(t *testing.T) {
	t.Run("credit note missing preceding", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "381"
		inv.Preceding = nil
		assert.ErrorContains(t, rules.Validate(inv), "credit and debit notes must have a billing reference")
	})

	t.Run("credit note missing reason", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "381"
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV-001",
				IssueDate: cal.NewDate(2022, 1, 1),
			},
		}
		assert.ErrorContains(t, rules.Validate(inv), "credit and debit notes must contain the reason for issuance")
	})

	t.Run("debit note missing preceding", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "383"
		inv.Preceding = nil
		assert.ErrorContains(t, rules.Validate(inv), "credit and debit notes must have a billing reference")
	})

	t.Run("debit note missing reason", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "383"
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV-001",
				IssueDate: cal.NewDate(2022, 1, 1),
			},
		}
		assert.ErrorContains(t, rules.Validate(inv), "credit and debit notes must contain the reason for issuance")
	})

	t.Run("standard invoice no preceding required", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Preceding = nil
		require.NoError(t, rules.Validate(inv))
	})
}

// --- KSA-25: Issue time ---

func TestIssueTimeValidation(t *testing.T) {
	t.Run("missing issue time", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.IssueTime = nil
		assert.ErrorContains(t, rules.Validate(inv), "issue time must be present")
	})

	t.Run("present issue time", func(t *testing.T) {
		inv := calculatedInvoice(t)
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-EN16931-02: Currency ---

func TestCurrencyValidation(t *testing.T) {
	t.Run("valid SAR currency", func(t *testing.T) {
		inv := calculatedInvoice(t)
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("invalid currency", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Currency = "USD"
		assert.ErrorContains(t, rules.Validate(inv), "invoice currency must be SAR")
	})
}

// --- Supplier validation ---

func TestSupplierValidation(t *testing.T) {
	t.Run("missing supplier", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier = nil
		assert.ErrorContains(t, rules.Validate(inv), "supplier is required")
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.TaxID = nil
		assert.ErrorContains(t, rules.Validate(inv), "seller VAT registration number is required")
	})

	t.Run("missing supplier tax ID code", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.TaxID.Code = ""
		assert.ErrorContains(t, rules.Validate(inv), "seller VAT registration number code is required")
	})

	// BR-KSA-08: Seller identity schemes
	t.Run("missing supplier identities", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Identities = nil
		assert.ErrorContains(t, rules.Validate(inv), "seller must have exactly one identity with a valid scheme ID")
	})

	t.Run("empty supplier identities", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Identities = []*org.Identity{}
		assert.ErrorContains(t, rules.Validate(inv), "seller must have exactly one identity with a valid scheme ID")
	})

	t.Run("supplier identity with invalid scheme", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Identities = []*org.Identity{
			{
				Code: "1234567890",
				Ext: tax.Extensions{
					zatca.ExtKeySellerIDScheme: "INVALID",
				},
			},
		}
		assert.ErrorContains(t, rules.Validate(inv), "seller must have exactly one identity with a valid scheme ID")
	})

	t.Run("supplier identity with non-alphanumeric code", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Identities = []*org.Identity{
			{
				Code: "123-456",
				Ext: tax.Extensions{
					zatca.ExtKeySellerIDScheme: "CRN",
				},
			},
		}
		assert.ErrorContains(t, rules.Validate(inv), "seller must have exactly one identity with a valid scheme ID")
	})

	t.Run("supplier with multiple ZATCA identities", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Identities = []*org.Identity{
			{
				Code: "1234567890",
				Ext: tax.Extensions{
					zatca.ExtKeySellerIDScheme: "CRN",
				},
			},
			{
				Code: "0987654321",
				Ext: tax.Extensions{
					zatca.ExtKeySellerIDScheme: "MOM",
				},
			},
		}
		assert.ErrorContains(t, rules.Validate(inv), "seller must have exactly one identity with a valid scheme ID")
	})

	t.Run("valid seller identity schemes", func(t *testing.T) {
		for _, scheme := range []cbc.Code{"CRN", "MOM", "MLS", "700", "SAG", "OTH"} {
			t.Run(string(scheme), func(t *testing.T) {
				inv := validStandardInvoice()
				inv.Supplier.Identities = []*org.Identity{
					{
						Code: "ABC123",
						Ext: tax.Extensions{
							zatca.ExtKeySellerIDScheme: scheme,
						},
					},
				}
				require.NoError(t, inv.Calculate())
				require.NoError(t, rules.Validate(inv))
			})
		}
	})
}

// --- Supplier address validation ---

func TestSupplierAddressValidation(t *testing.T) {
	t.Run("missing supplier address", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Addresses = nil
		assert.ErrorContains(t, rules.Validate(inv), "supplier address is required")
	})

	t.Run("missing building number", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Addresses[0].Number = ""
		assert.ErrorContains(t, rules.Validate(inv), "supplier address building number is required")
	})

	t.Run("invalid building number - 3 digits", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Addresses[0].Number = "123"
		assert.ErrorContains(t, rules.Validate(inv), "supplier address building number must contain 4 digits")
	})

	t.Run("invalid building number - 5 digits", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Addresses[0].Number = "12345"
		assert.ErrorContains(t, rules.Validate(inv), "supplier address building number must contain 4 digits")
	})

	t.Run("invalid building number - non-numeric", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Addresses[0].Number = "ABCD"
		assert.ErrorContains(t, rules.Validate(inv), "supplier address building number must contain 4 digits")
	})

	t.Run("missing district (street_extra)", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Addresses[0].StreetExtra = ""
		assert.ErrorContains(t, rules.Validate(inv), "supplier address must have a district")
	})

	t.Run("missing street", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Addresses[0].Street = ""
		assert.ErrorContains(t, rules.Validate(inv), "supplier address must have a street name")
	})

	t.Run("missing postal code", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Addresses[0].Code = ""
		assert.ErrorContains(t, rules.Validate(inv), "supplier postal code is required")
	})

	t.Run("invalid postal code - 4 digits", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Addresses[0].Code = "1234"
		assert.ErrorContains(t, rules.Validate(inv), "supplier postal code must be 5 digits")
	})

	t.Run("invalid postal code - non-numeric", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Addresses[0].Code = "ABCDE"
		assert.ErrorContains(t, rules.Validate(inv), "supplier postal code must be 5 digits")
	})

	t.Run("missing locality", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Addresses[0].Locality = ""
		assert.ErrorContains(t, rules.Validate(inv), "supplier address must have a city name")
	})

	t.Run("missing country", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Supplier.Addresses[0].Country = ""
		assert.ErrorContains(t, rules.Validate(inv), "supplier address must have a country code")
	})
}

// --- Customer validation ---

func TestCustomerValidation(t *testing.T) {
	t.Run("missing customer", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Customer = nil
		assert.ErrorContains(t, rules.Validate(inv), "customer is required")
	})

	// BR-KSA-14: Buyer must have tax ID or identity with valid scheme
	t.Run("customer with tax ID only", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Customer.Identities = nil
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("customer with identity only", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "1234567890",
				Ext: tax.Extensions{
					zatca.ExtKeyBuyerIDScheme: "CRN",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("customer without tax ID or identity", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Customer.TaxID = nil
		inv.Customer.Identities = nil
		assert.ErrorContains(t, rules.Validate(inv), "buyer must have a tax ID or an identity with a valid ZATCA scheme")
	})

	// BR-KSA-45: Standard invoices require buyer name
	t.Run("standard invoice missing buyer name", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Customer.Name = ""
		assert.ErrorContains(t, rules.Validate(inv), "buyer name is required for standard tax invoices")
	})

	t.Run("valid buyer identity schemes", func(t *testing.T) {
		for _, scheme := range []cbc.Code{"TIN", "CRN", "MOM", "MLS", "700", "SAG", "NAT", "GCC", "IQA", "PAS", "OTH"} {
			t.Run(string(scheme), func(t *testing.T) {
				inv := validStandardInvoice()
				inv.Customer.TaxID = nil
				inv.Customer.Identities = []*org.Identity{
					{
						Code: "ABC123",
						Ext: tax.Extensions{
							zatca.ExtKeyBuyerIDScheme: scheme,
						},
					},
				}
				require.NoError(t, inv.Calculate())
				require.NoError(t, rules.Validate(inv))
			})
		}
	})
}

// --- Customer address validation for standard invoices ---

func TestCustomerAddressStandardValidation(t *testing.T) {
	t.Run("missing customer address", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Customer.Addresses = nil
		assert.ErrorContains(t, rules.Validate(inv), "customer address is required")
	})

	t.Run("standard invoice missing customer building number", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Customer.Addresses[0].Number = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer address must have a building number")
	})

	t.Run("standard invoice missing customer street", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Customer.Addresses[0].Street = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer address must have a street name")
	})

	t.Run("standard invoice missing customer locality", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Customer.Addresses[0].Locality = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer address must have a city name")
	})

	t.Run("standard invoice missing customer country", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Customer.Addresses[0].Country = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer address must have a country code")
	})

	t.Run("standard invoice missing customer district", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Customer.Addresses[0].StreetExtra = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer address must have a district")
	})
}

// --- BR-KSA-63/BR-KSA-67: Customer address validation for SA country ---

func TestCustomerAddressSACountryValidation(t *testing.T) {
	t.Run("SA buyer missing postal code", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Customer.Addresses[0].Code = ""
		assert.ErrorContains(t, rules.Validate(inv), "customer address must have a postal code")
	})

	t.Run("SA buyer invalid postal code - 4 digits", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Customer.Addresses[0].Code = "1234"
		assert.ErrorContains(t, rules.Validate(inv), "buyer postal code must be 5 digits when country is SA")
	})

	t.Run("SA buyer invalid postal code - non-numeric", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Customer.Addresses[0].Code = "ABCDE"
		assert.ErrorContains(t, rules.Validate(inv), "buyer postal code must be 5 digits when country is SA")
	})

	t.Run("non-SA buyer address does not require 5-digit postal code", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Customer.Addresses[0].Country = "AE"
		inv.Customer.Addresses[0].Code = "12345678" // non-SA format OK
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

// --- Delivery validation ---

func TestDeliveryValidation(t *testing.T) {
	t.Run("missing delivery", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Delivery = nil
		assert.ErrorContains(t, rules.Validate(inv), "delivery is required")
	})

	// BR-KSA-15/BR-KSA-72: Tax invoices need supply date and delivery period
	t.Run("standard tax invoice missing delivery date", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Delivery.Date = nil
		assert.ErrorContains(t, rules.Validate(inv), "delivery period must have a supply date")
	})

	t.Run("standard tax invoice missing delivery period", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Delivery.Period = nil
		assert.ErrorContains(t, rules.Validate(inv), "supply must have a delivery period")
	})

	// BR-KSA-35, BR-KSA-36: End date >= supply date
	t.Run("supply end date before supply date", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Delivery.Date = cal.NewDate(2022, 3, 1)
		inv.Delivery.Period = &cal.Period{
			Start: cal.MakeDate(2022, 2, 1),
			End:   cal.MakeDate(2022, 2, 15),
		}
		assert.ErrorContains(t, rules.Validate(inv), "delivery period end date must be >= supply date")
	})

	t.Run("supply end date equals supply date", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Delivery.Date = cal.NewDate(2022, 2, 1)
		inv.Delivery.Period = &cal.Period{
			Start: cal.MakeDate(2022, 2, 1),
			End:   cal.MakeDate(2022, 2, 1),
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("supply end date after supply date", func(t *testing.T) {
		inv := calculatedInvoice(t)
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-31: Simplified invoice flag validation ---

func TestSimplifiedInvoiceFlagsValidation(t *testing.T) {
	t.Run("simplified with export flag is invalid", func(t *testing.T) {
		inv := calculatedSimplifiedInvoice(t)
		// 02 + third-party=0 + nominal=0 + export=1 + summary=0 + self-billing=0
		inv.Tax.Ext[zatca.ExtKeyInvoiceType] = "0200100"
		assert.ErrorContains(t, rules.Validate(inv), "simplified invoices only allow third-party, nominal, and summary flags")
	})

	t.Run("simplified with self-billing flag is invalid", func(t *testing.T) {
		inv := calculatedSimplifiedInvoice(t)
		inv.Tax.Ext[zatca.ExtKeyInvoiceType] = "0200001"
		assert.ErrorContains(t, rules.Validate(inv), "simplified invoices only allow third-party, nominal, and summary flags")
	})

	t.Run("simplified with third-party flag is valid", func(t *testing.T) {
		inv := calculatedSimplifiedInvoice(t)
		inv.Tax.Ext[zatca.ExtKeyInvoiceType] = "0210000"
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified with nominal flag is valid", func(t *testing.T) {
		inv := calculatedSimplifiedInvoice(t)
		inv.Tax.Ext[zatca.ExtKeyInvoiceType] = "0201000"
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified with summary flag is valid", func(t *testing.T) {
		inv := calculatedSimplifiedInvoice(t)
		inv.Customer.Name = "Consumer" // summary invoices require buyer name
		inv.Tax.Ext[zatca.ExtKeyInvoiceType] = "0200010"
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("simplified with third-party nominal summary is valid", func(t *testing.T) {
		inv := calculatedSimplifiedInvoice(t)
		inv.Customer.Name = "Consumer"
		inv.Tax.Ext[zatca.ExtKeyInvoiceType] = "0211010"
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("standard with export and self-billing flags is valid", func(t *testing.T) {
		// Standard invoices should NOT be constrained by simplified flag rules
		// But BR-KSA-07 prevents self-billing + export together
		inv := calculatedInvoice(t)
		inv.Tax.Ext[zatca.ExtKeyInvoiceType] = "0100100" // export only
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "ABC123",
				Ext: tax.Extensions{
					zatca.ExtKeyBuyerIDScheme: "CRN",
				},
			},
		}
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-07: No self-billing for exports ---

func TestNoSelfBillingForExports(t *testing.T) {
	t.Run("export with self-billing is invalid", func(t *testing.T) {
		inv := calculatedInvoice(t)
		// export=1, self-billing=1
		inv.Tax.Ext[zatca.ExtKeyInvoiceType] = "0100101"
		assert.ErrorContains(t, rules.Validate(inv), "self-billing is not allowed for export invoices")
	})

	t.Run("export without self-billing is valid", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Tax.Ext[zatca.ExtKeyInvoiceType] = "0100100"
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "ABC123",
				Ext: tax.Extensions{
					zatca.ExtKeyBuyerIDScheme: "CRN",
				},
			},
		}
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("self-billing without export is valid", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Tax.Ext[zatca.ExtKeyInvoiceType] = "0100001"
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-46: Export invoice buyer VAT ---

func TestExportInvoiceBuyerVAT(t *testing.T) {
	t.Run("export invoice with buyer VAT is invalid", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Tax.Ext[zatca.ExtKeyInvoiceType] = "0100100"
		assert.ErrorContains(t, rules.Validate(inv), "export invoices must not have buyer VAT registration number")
	})

	t.Run("export invoice without buyer VAT is valid", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Tax.Ext[zatca.ExtKeyInvoiceType] = "0100100"
		inv.Customer.TaxID = nil
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "ABC123",
				Ext: tax.Extensions{
					zatca.ExtKeyBuyerIDScheme: "PAS",
				},
			},
		}
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("non-export invoice with buyer VAT is valid", func(t *testing.T) {
		inv := calculatedInvoice(t)
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-71: Simplified summary invoice buyer name ---

func TestSimplifiedSummaryBuyerName(t *testing.T) {
	t.Run("simplified summary without buyer name is invalid", func(t *testing.T) {
		inv := calculatedSimplifiedInvoice(t)
		inv.Tax.Ext[zatca.ExtKeyInvoiceType] = "0200010" // simplified summary
		inv.Customer.Name = ""
		assert.ErrorContains(t, rules.Validate(inv), "buyer name is required for simplified summary invoices")
	})

	t.Run("simplified summary with buyer name is valid", func(t *testing.T) {
		inv := calculatedSimplifiedInvoice(t)
		inv.Tax.Ext[zatca.ExtKeyInvoiceType] = "0200010"
		inv.Customer.Name = "Consumer Co"
		require.NoError(t, rules.Validate(inv))
	})
}

// --- BR-KSA-52, BR-KSA-53: Standard invoice line requirements ---

func TestStandardInvoiceLineRequirements(t *testing.T) {
	t.Run("standard invoice missing line taxes", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Lines[0].Taxes = nil
		assert.ErrorContains(t, rules.Validate(inv), "line taxes are required for standard tax invoices")
	})

	t.Run("standard invoice missing line total", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Lines[0].Total = nil
		assert.ErrorContains(t, rules.Validate(inv), "line total is required for standard tax invoices")
	})
}

// --- BR-KSA-49: EDU/HEA exemption requires buyer NAT identity ---

func TestEDUHEABuyerNATValidation(t *testing.T) {
	t.Run("EDU exemption without buyer NAT is invalid", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-EDU",
				},
			},
		}
		assert.ErrorContains(t, rules.Validate(inv), "EDU/HEA tax exemption requires buyer with NAT identity")
	})

	t.Run("HEA exemption without buyer NAT is invalid", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-HEA",
				},
			},
		}
		assert.ErrorContains(t, rules.Validate(inv), "EDU/HEA tax exemption requires buyer with NAT identity")
	})

	t.Run("EDU exemption with buyer NAT is valid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-EDU",
				},
			},
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "1234567890",
				Ext: tax.Extensions{
					zatca.ExtKeyBuyerIDScheme: "NAT",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("HEA exemption with buyer NAT is valid", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-HEA",
				},
			},
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "1234567890",
				Ext: tax.Extensions{
					zatca.ExtKeyBuyerIDScheme: "NAT",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("non-EDU/HEA exemption without NAT is valid", func(t *testing.T) {
		inv := calculatedInvoice(t)
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("EDU with buyer CRN instead of NAT is invalid", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-EDU",
				},
			},
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "1234567890",
				Ext: tax.Extensions{
					zatca.ExtKeyBuyerIDScheme: "CRN",
				},
			},
		}
		assert.ErrorContains(t, rules.Validate(inv), "EDU/HEA tax exemption requires buyer with NAT identity")
	})
}

// --- BR-KSA-25: Simplified EDU/HEA requires buyer name ---

func TestSimplifiedEDUHEABuyerNameValidation(t *testing.T) {
	t.Run("simplified EDU without buyer name is invalid", func(t *testing.T) {
		inv := calculatedSimplifiedInvoice(t)
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-EDU",
				},
			},
		}
		inv.Customer.Name = ""
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "1234567890",
				Ext: tax.Extensions{
					zatca.ExtKeyBuyerIDScheme: "NAT",
				},
			},
		}
		assert.ErrorContains(t, rules.Validate(inv), "simplified invoice with EDU/HEA exemption requires buyer name")
	})

	t.Run("simplified HEA without buyer name is invalid", func(t *testing.T) {
		inv := calculatedSimplifiedInvoice(t)
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-HEA",
				},
			},
		}
		inv.Customer.Name = ""
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "1234567890",
				Ext: tax.Extensions{
					zatca.ExtKeyBuyerIDScheme: "NAT",
				},
			},
		}
		assert.ErrorContains(t, rules.Validate(inv), "simplified invoice with EDU/HEA exemption requires buyer name")
	})

	t.Run("simplified EDU with buyer name is valid", func(t *testing.T) {
		inv := validSimplifiedInvoice()
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-EDU",
				},
			},
		}
		inv.Customer.Name = "Student Name"
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "1234567890",
				Ext: tax.Extensions{
					zatca.ExtKeyBuyerIDScheme: "NAT",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("standard EDU without buyer name triggers standard rule not simplified rule", func(t *testing.T) {
		inv := calculatedInvoice(t)
		inv.Lines[0].Taxes = tax.Set{
			{
				Category: "VAT",
				Ext: tax.Extensions{
					cef.ExtKeyVATEX: "VATEX-SA-EDU",
				},
			},
		}
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "1234567890",
				Ext: tax.Extensions{
					zatca.ExtKeyBuyerIDScheme: "NAT",
				},
			},
		}
		inv.Customer.Name = ""
		// Standard invoices have their own buyer name rule (BR-KSA-45)
		assert.ErrorContains(t, rules.Validate(inv), "buyer name is required for standard tax invoices")
	})
}

// --- Invoice type codes ---

func TestInvoiceTypeCodes(t *testing.T) {
	t.Run("all standard type codes", func(t *testing.T) {
		for _, code := range zatca.InvTypesStandard {
			t.Run(string(code), func(t *testing.T) {
				inv := validStandardInvoice()
				inv.Tax.Ext[zatca.ExtKeyInvoiceType] = code

				// Export invoices need special customer setup
				if len(code) == 7 && code[4] == '1' {
					inv.Customer.TaxID = nil
					inv.Customer.Identities = []*org.Identity{
						{
							Code: "ABC123",
							Ext: tax.Extensions{
								zatca.ExtKeyBuyerIDScheme: "PAS",
							},
						},
					}
				}

				// Summary invoices need delivery period
				require.NoError(t, inv.Calculate())
				// Just check it can calculate; some may fail validation
				// due to combined flag constraints (e.g., export+self-billing)
			})
		}
	})

	t.Run("all simplified type codes", func(t *testing.T) {
		for _, code := range zatca.InvTypesSimplified {
			t.Run(string(code), func(t *testing.T) {
				inv := validSimplifiedInvoice()
				inv.Tax.Ext[zatca.ExtKeyInvoiceType] = code
				// Summary simplified invoices need buyer name
				if len(code) == 7 && code[5] == '1' {
					inv.Customer.Name = "Consumer"
				}
				require.NoError(t, inv.Calculate())
				require.NoError(t, rules.Validate(inv))
			})
		}
	})
}

// --- Scenarios ---

func TestScenarios(t *testing.T) {
	t.Run("standard invoice gets 388", func(t *testing.T) {
		inv := validStandardInvoice()
		// Remove the manually set document type to let scenarios set it
		delete(inv.Tax.Ext, untdid.ExtKeyDocumentType)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, cbc.Code("388"), inv.Tax.Ext[untdid.ExtKeyDocumentType])
	})

	t.Run("prepayment invoice gets 386", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Tags = tax.WithTags(tax.TagPrepayment)
		delete(inv.Tax.Ext, untdid.ExtKeyDocumentType)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, cbc.Code("386"), inv.Tax.Ext[untdid.ExtKeyDocumentType])
	})

	t.Run("credit note gets 381", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Type = bill.InvoiceTypeCreditNote
		inv.Preceding = []*org.DocumentRef{
			{Code: "INV-001", IssueDate: cal.NewDate(2022, 1, 1), Reason: "Return"},
		}
		delete(inv.Tax.Ext, untdid.ExtKeyDocumentType)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, cbc.Code("381"), inv.Tax.Ext[untdid.ExtKeyDocumentType])
	})

	t.Run("debit note gets 383", func(t *testing.T) {
		inv := validStandardInvoice()
		inv.Type = bill.InvoiceTypeDebitNote
		inv.Preceding = []*org.DocumentRef{
			{Code: "INV-001", IssueDate: cal.NewDate(2022, 1, 1), Reason: "Adjustment"},
		}
		delete(inv.Tax.Ext, untdid.ExtKeyDocumentType)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, cbc.Code("383"), inv.Tax.Ext[untdid.ExtKeyDocumentType])
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
}
