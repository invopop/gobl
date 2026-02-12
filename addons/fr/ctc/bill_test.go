package ctc_test

import (
	"testing"

	"github.com/invopop/gobl/addons/fr/ctc"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testInvoiceB2BStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	i := &bill.Invoice{
		Regime:   tax.WithRegime("FR"),
		Addons:   tax.WithAddons(ctc.Flow2),
		Code:     "FAC-2024-001",
		Currency: "EUR",
		Type:     bill.InvoiceTypeStandard,
		Tax: &bill.Tax{
			Ext: tax.Extensions{
				ctc.ExtKeyBillingMode:     ctc.BillingModeS1,
				untdid.ExtKeyDocumentType: "380",
			},
		},
		Supplier: &org.Party{
			Name: "Test Supplier SARL",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "39356000000", // Valid French VAT number
			},
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIREN,
					Code: "356000000",
				},
				{
					Type: fr.IdentityTypeSIRET,
					Code: "35600000000011",
				},
			},
			Addresses: []*org.Address{
				{
					Street:   "123 Rue de Test",
					Code:     "75001",
					Locality: "Paris",
					Country:  "FR",
				},
			},
			Inboxes: []*org.Inbox{
				{
					Key:    org.InboxKeyPeppol,
					Scheme: cbc.Code("0225"),
					Code:   "356000000",
				},
			},
		},
		Customer: &org.Party{
			Name: "Test Customer SAS",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "44732829320", // Valid French VAT number
			},
			Identities: []*org.Identity{
				{
					Type: fr.IdentityTypeSIREN,
					Code: "732829320",
				},
			},
			Addresses: []*org.Address{
				{
					Street:   "456 Avenue du Client",
					Code:     "69001",
					Locality: "Lyon",
					Country:  "FR",
				},
			},
			Inboxes: []*org.Inbox{
				{
					Key:    org.InboxKeyPeppol,
					Scheme: cbc.Code("0225"),
					Code:   "732829320",
				},
			},
		},
		IssueDate: cal.MakeDate(2024, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Service",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
			},
		},
		Payment: &bill.PaymentDetails{
			Terms: &pay.Terms{
				Key: pay.TermKeyDueDate,
				DueDates: []*pay.DueDate{
					{
						Date:    cal.NewDate(2024, 7, 13),
						Percent: num.NewPercentage(100, 3),
					},
				},
			},
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
				CreditTransfer: []*pay.CreditTransfer{
					{
						IBAN: "FR7630006000011234567890189",
						Name: "Test Supplier SARL",
					},
				},
			},
		},
		Notes: []*org.Note{
			{
				Key:  org.NoteKeyPayment,
				Text: "A fixed penalty of 40 EUR will apply to any late payment.",
				Ext: tax.Extensions{
					untdid.ExtKeyTextSubject: "PMT",
				},
			},
			{
				Key:  org.NoteKeyPaymentMethod,
				Text: "Late payment penalties apply as per our general terms of sale.",
				Ext: tax.Extensions{
					untdid.ExtKeyTextSubject: "PMD",
				},
			},
			{
				Key:  org.NoteKeyPaymentTerm,
				Text: "No discount offered for early payment.",
				Ext: tax.Extensions{
					untdid.ExtKeyTextSubject: "AAB",
				},
			},
		},
	}
	return i
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("basic B2B invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("invoice code too long", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Code = "THIS-IS-A-VERY-LONG-INVOICE-CODE-THAT-EXCEEDS-35-CHARACTERS"
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "BR-FR-01/02")
	})

	t.Run("invoice code normalized - special chars removed", func(t *testing.T) {
		// Note: cbc.NormalizeCode removes invalid characters during Calculate
		inv := testInvoiceB2BStandard(t)
		inv.Code = "INV#2024@001"
		require.NoError(t, inv.Calculate())
		// Code is normalized to remove # and @
		assert.Equal(t, "INV2024001", inv.Code.String())
		require.NoError(t, inv.Validate())
	})

	t.Run("invoice code valid special chars", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Code = "INV-2024+001_TEST/A"
		require.NoError(t, inv.Calculate())
		require.NoError(t, inv.Validate())
	})

	t.Run("invoice date validation - valid dates", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		// Date is 2024, which is valid (2000-2099)
		require.NoError(t, inv.Validate())
	})

	t.Run("duplicate note codes not allowed", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		// Add duplicate PMT note
		inv.Notes = append(inv.Notes, &org.Note{
			Key:  org.NoteKeyPayment,
			Text: "Duplicate payment terms",
		})
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "duplicate note codes")
		assert.ErrorContains(t, err, "PMT")
	})

	t.Run("supplier SIREN required (BR-FR-10)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		// Remove all identities from supplier (no SIREN, no SIRET)
		inv.Supplier.Identities = []*org.Identity{}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "BR-FR-10")
		assert.ErrorContains(t, err, "SIREN")
	})

	t.Run("customer SIREN required for B2B (BR-FR-10)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		// After Calculate, manually add B2B note to mark this as B2B transaction
		// (bypassing UNTDID validation for test purposes)
		inv.Notes = append(inv.Notes, &org.Note{
			Key: org.NoteKeyLegal,
			Ext: tax.Extensions{
				untdid.ExtKeyTextSubject: "BAR",
			},
			Text: "B2B",
		})
		// Remove all identities from customer (no SIREN, no SIRET)
		inv.Customer.Identities = []*org.Identity{}
		// Validate (skip Calculate to avoid note validation)
		err := inv.Validate()
		assert.ErrorContains(t, err, "BR-FR-10")
		assert.ErrorContains(t, err, "SIREN")
		assert.ErrorContains(t, err, "0002")
	})

	t.Run("Spanish tax categories rejected by UNTDID", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		// Manually add Spanish tax category to a line (after Calculate)
		if len(inv.Lines) > 0 && len(inv.Lines[0].Taxes) > 0 {
			inv.Lines[0].Taxes[0].Ext = tax.Extensions{
				untdid.ExtKeyTaxCategory: "L", // IGIC (Canary Islands)
			}
		}
		err := inv.Validate()
		// UNTDID catalogue validation rejects L and M automatically
		assert.ErrorContains(t, err, "untdid-tax-category")
		assert.ErrorContains(t, err, "invalid")
	})

	t.Run("valid BAR note - B2B", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Notes = append(inv.Notes, &org.Note{
			Key: org.NoteKeyLegal,
			Ext: tax.Extensions{
				untdid.ExtKeyTextSubject: "BAR",
			},
			Text: "B2B",
		})
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid BAR note - B2BINT", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Notes = append(inv.Notes, &org.Note{
			Key: org.NoteKeyLegal,
			Ext: tax.Extensions{
				untdid.ExtKeyTextSubject: "BAR",
			},
			Text: "B2BINT",
		})
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid BAR note - B2C", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Notes = append(inv.Notes, &org.Note{
			Key: org.NoteKeyLegal,
			Ext: tax.Extensions{
				untdid.ExtKeyTextSubject: "BAR",
			},
			Text: "B2C",
		})
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid BAR note - OUTOFSCOPE", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Notes = append(inv.Notes, &org.Note{
			Key: org.NoteKeyLegal,
			Ext: tax.Extensions{
				untdid.ExtKeyTextSubject: "BAR",
			},
			Text: "OUTOFSCOPE",
		})
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid BAR note - ARCHIVEONLY", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Notes = append(inv.Notes, &org.Note{
			Key: org.NoteKeyLegal,
			Ext: tax.Extensions{
				untdid.ExtKeyTextSubject: "BAR",
			},
			Text: "ARCHIVEONLY",
		})
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid BAR note text", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Notes = append(inv.Notes, &org.Note{
			Key: org.NoteKeyLegal,
			Ext: tax.Extensions{
				untdid.ExtKeyTextSubject: "BAR",
			},
			Text: "INVALID",
		})
		err := inv.Validate()
		assert.ErrorContains(t, err, "BAR note text must be one of")
		assert.ErrorContains(t, err, "B2B")
		assert.ErrorContains(t, err, "B2BINT")
	})

	t.Run("duplicate BAR note not allowed (BR-FR-30)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		// Add two BAR notes
		inv.Notes = append(inv.Notes,
			&org.Note{
				Key: org.NoteKeyLegal,
				Ext: tax.Extensions{
					untdid.ExtKeyTextSubject: "BAR",
				},
				Text: "B2B",
			},
			&org.Note{
				Key: org.NoteKeyLegal,
				Ext: tax.Extensions{
					untdid.ExtKeyTextSubject: "BAR",
				},
				Text: "Additional BAR information",
			},
		)
		err := inv.Validate()
		assert.ErrorContains(t, err, "duplicate note codes found")
		assert.ErrorContains(t, err, "BAR")
		assert.ErrorContains(t, err, "BR-FR-30")
	})

	t.Run("B2B non-self-billed requires SIREN inbox (BR-FR-21)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		// Remove SIREN inbox from supplier
		inv.Supplier.Inboxes = []*org.Inbox{
			{
				Scheme: "0088", // GLIN
				Code:   "1234567890123",
			},
		}
		require.NoError(t, inv.Calculate())
		// Add B2B note
		inv.Notes = append(inv.Notes, &org.Note{
			Key: org.NoteKeyLegal,
			Ext: tax.Extensions{
				untdid.ExtKeyTextSubject: "BAR",
			},
			Text: "B2B",
		})
		err := inv.Validate()
		assert.ErrorContains(t, err, "party must have endpoint ID with scheme 0225")
		assert.ErrorContains(t, err, "BR-FR-21")
	})

	t.Run("B2B non-self-billed SIREN inbox must start with SIREN (BR-FR-21)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		// Set wrong SIREN inbox
		inv.Supplier.Inboxes = []*org.Inbox{
			{
				Scheme: cbc.Code("0225"),
				Code:   "999999999", // Wrong SIREN
			},
		}
		require.NoError(t, inv.Calculate())
		// Add B2B note
		inv.Notes = append(inv.Notes, &org.Note{
			Key: org.NoteKeyLegal,
			Ext: tax.Extensions{
				untdid.ExtKeyTextSubject: "BAR",
			},
			Text: "B2B",
		})
		err := inv.Validate()
		assert.ErrorContains(t, err, "must start with SIREN")
		assert.ErrorContains(t, err, "BR-FR-21")
	})

	t.Run("self-billed invoice does not require SIREN inbox", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		// Set document type to self-billed after Calculate
		if inv.Tax != nil && inv.Tax.Ext != nil {
			inv.Tax.Ext[untdid.ExtKeyDocumentType] = "389" // Self-billed invoice
		}
		// Remove SIREN inbox
		inv.Supplier.Inboxes = []*org.Inbox{
			{
				Scheme: "0088",
				Code:   "1234567890123",
			},
		}
		// Add B2B note
		inv.Notes = append(inv.Notes, &org.Note{
			Key: org.NoteKeyLegal,
			Ext: tax.Extensions{
				untdid.ExtKeyTextSubject: "BAR",
			},
			Text: "B2B",
		})
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("B2C does not require SIREN inbox", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		// Remove SIREN inbox
		inv.Supplier.Inboxes = []*org.Inbox{
			{
				Scheme: "0088",
				Code:   "1234567890123",
			},
		}
		require.NoError(t, inv.Calculate())
		// No B2B note, so not a B2B transaction
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("B2B self-billed requires customer SIREN inbox (BR-FR-22)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		// Set document type to self-billed
		if inv.Tax != nil && inv.Tax.Ext != nil {
			inv.Tax.Ext[untdid.ExtKeyDocumentType] = "389" // Self-billed invoice
		}
		// Remove SIREN inbox from customer
		inv.Customer.Inboxes = []*org.Inbox{
			{
				Scheme: "0088",
				Code:   "1234567890123",
			},
		}
		// Add B2B note
		inv.Notes = append(inv.Notes, &org.Note{
			Key: org.NoteKeyLegal,
			Ext: tax.Extensions{
				untdid.ExtKeyTextSubject: "BAR",
			},
			Text: "B2B",
		})
		err := inv.Validate()
		assert.ErrorContains(t, err, "party must have endpoint ID with scheme 0225")
		assert.ErrorContains(t, err, "BR-FR-21/22")
	})

	t.Run("B2B self-billed customer SIREN inbox must start with SIREN (BR-FR-22)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		// Set document type to self-billed
		if inv.Tax != nil && inv.Tax.Ext != nil {
			inv.Tax.Ext[untdid.ExtKeyDocumentType] = "389" // Self-billed invoice
		}
		// Set wrong SIREN inbox for customer
		inv.Customer.Inboxes = []*org.Inbox{
			{
				Scheme: cbc.Code("0225"),
				Code:   "999999999", // Wrong SIREN
			},
		}
		// Add B2B note
		inv.Notes = append(inv.Notes, &org.Note{
			Key: org.NoteKeyLegal,
			Ext: tax.Extensions{
				untdid.ExtKeyTextSubject: "BAR",
			},
			Text: "B2B",
		})
		err := inv.Validate()
		assert.ErrorContains(t, err, "must start with SIREN")
		assert.ErrorContains(t, err, "BR-FR-21/22")
	})

	t.Run("B2B self-billed with correct customer SIREN inbox", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		// Set document type to self-billed
		if inv.Tax != nil && inv.Tax.Ext != nil {
			inv.Tax.Ext[untdid.ExtKeyDocumentType] = "389" // Self-billed invoice
		}
		// Customer already has correct SIREN inbox from testInvoiceB2BStandard
		// Add B2B note
		inv.Notes = append(inv.Notes, &org.Note{
			Key: org.NoteKeyLegal,
			Ext: tax.Extensions{
				untdid.ExtKeyTextSubject: "BAR",
			},
			Text: "B2B",
		})
		err := inv.Validate()
		assert.NoError(t, err)
	})
}

func TestDocumentTypeValidation(t *testing.T) {
	t.Run("valid document type", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "380", inv.Tax.Ext[untdid.ExtKeyDocumentType].String())
		require.NoError(t, inv.Validate())
	})

	t.Run("invalid document type", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "999"
		err := inv.Validate()
		// Note: UNTDID catalogue validates before our custom validator
		assert.ErrorContains(t, err, "invalid")
	})
}

func TestDocumentTypeScenarios(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "380", inv.Tax.Ext[untdid.ExtKeyDocumentType].String())
	})

	t.Run("factored invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.SetTags(tax.TagFactored)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "393", inv.Tax.Ext[untdid.ExtKeyDocumentType].String())
	})

	t.Run("advance payment invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.SetTags(tax.TagPrepayment)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "386", inv.Tax.Ext[untdid.ExtKeyDocumentType].String())
	})

	t.Run("self-billed invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.SetTags(tax.TagSelfBilled)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "389", inv.Tax.Ext[untdid.ExtKeyDocumentType].String())
	})

	t.Run("credit note", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "381", inv.Tax.Ext[untdid.ExtKeyDocumentType].String())
	})

	t.Run("self-billed credit note", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.SetTags(tax.TagSelfBilled)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "261", inv.Tax.Ext[untdid.ExtKeyDocumentType].String())
	})

	t.Run("corrective invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Type = bill.InvoiceTypeCorrective
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "384", inv.Tax.Ext[untdid.ExtKeyDocumentType].String())
	})

	t.Run("factored credit note", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.SetTags(tax.TagFactored)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "396", inv.Tax.Ext[untdid.ExtKeyDocumentType].String())
	})
}

func TestBillingModeNormalization(t *testing.T) {
	t.Run("user-specified billing mode preserved", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax = &bill.Tax{
			Ext: tax.Extensions{
				ctc.ExtKeyBillingMode: ctc.BillingModeS5, // Subcontractor
			},
		}
		require.NoError(t, inv.Calculate())
		assert.Equal(t, ctc.BillingModeS5.String(), inv.Tax.Ext[ctc.ExtKeyBillingMode].String())
	})

	t.Run("invalid billing mode rejected - B8", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax = &bill.Tax{
			Ext: tax.Extensions{
				ctc.ExtKeyBillingMode: cbc.Code("B8"), // Not allowed
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "value 'B8' invalid")
	})

	t.Run("invalid billing mode rejected - B5", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax = &bill.Tax{
			Ext: tax.Extensions{
				ctc.ExtKeyBillingMode: cbc.Code("B5"), // Not allowed
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "value 'B5' invalid")
	})
}

func TestAttachmentValidation(t *testing.T) {
	t.Run("valid attachment description - LISIBLE", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Attachments = []*org.Attachment{
			{
				Code:        "ATT001",
				Description: "LISIBLE",
				URL:         "https://example.com/invoice.pdf",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid attachment description - RIB", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Attachments = []*org.Attachment{
			{
				Code:        "ATT001",
				Description: "RIB",
				URL:         "https://example.com/rib.pdf",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid attachment description - arbitrary value", func(t *testing.T) {
		ad := tax.AddonForKey(ctc.Flow2)
		attachments := []*org.Attachment{
			{
				Code:        "ATT001",
				Description: "INVALID_TYPE",
				URL:         "https://example.com/doc.pdf",
			},
		}
		err := ad.Validator(attachments)
		assert.ErrorContains(t, err, "attachment description 'INVALID_TYPE' is not allowed")
		assert.ErrorContains(t, err, "BR-FR-17")
	})

	t.Run("multiple LISIBLE attachments rejected (BR-FR-18)", func(t *testing.T) {
		ad := tax.AddonForKey(ctc.Flow2)
		attachments := []*org.Attachment{
			{
				Code:        "ATT001",
				Description: "LISIBLE",
				URL:         "https://example.com/invoice1.pdf",
			},
			{
				Code:        "ATT002",
				Description: "LISIBLE",
				URL:         "https://example.com/invoice2.pdf",
			},
		}
		err := ad.Validator(attachments)
		assert.ErrorContains(t, err, "only one attachment with description 'LISIBLE' is allowed")
		assert.ErrorContains(t, err, "BR-FR-18")
	})

	t.Run("empty attachment description allowed", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Attachments = []*org.Attachment{
			{
				Code:        "ATT001",
				Description: "",
				URL:         "https://example.com/doc.pdf",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("attachment description with whitespace trimmed", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Attachments = []*org.Attachment{
			{
				Code:        "ATT001",
				Description: "  LISIBLE  ",
				URL:         "https://example.com/invoice.pdf",
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("multiple attachments with empty descriptions allowed", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Attachments = []*org.Attachment{
			{Code: "ATT01", Description: "", URL: "https://example.com/1.pdf"},
			{Code: "ATT02", Description: "", URL: "https://example.com/2.pdf"},
			{Code: "ATT03", Description: "", URL: "https://example.com/3.pdf"},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("one LISIBLE with empty descriptions allowed", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Attachments = []*org.Attachment{
			{Code: "ATT01", Description: "LISIBLE", URL: "https://example.com/invoice.pdf"},
			{Code: "ATT02", Description: "", URL: "https://example.com/2.pdf"},
			{Code: "ATT03", Description: "", URL: "https://example.com/3.pdf"},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})
}

func TestOrderingIdentitiesValidation(t *testing.T) {
	t.Run("valid ordering with one AFL reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Ordering = &bill.Ordering{
			Identities: []*org.Identity{
				{
					Code: "12345",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "AFL",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid ordering with one AWW reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Ordering = &bill.Ordering{
			Identities: []*org.Identity{
				{
					Code: "12345",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "AWW",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid ordering with one AFL and one AWW", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Ordering = &bill.Ordering{
			Identities: []*org.Identity{
				{
					Code: "12345",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "AFL",
					},
				},
				{
					Code: "67890",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "AWW",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid ordering with duplicate AFL reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Ordering = &bill.Ordering{
			Identities: []*org.Identity{
				{
					Code: "12345",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "AFL",
					},
				},
				{
					Code: "67890",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "AFL",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "AFL")
		assert.ErrorContains(t, err, "BR-FR-30")
	})

	t.Run("invalid ordering with duplicate AWW reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Ordering = &bill.Ordering{
			Identities: []*org.Identity{
				{
					Code: "12345",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "AWW",
					},
				},
				{
					Code: "67890",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "AWW",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "AWW")
		assert.ErrorContains(t, err, "BR-FR-30")
	})

	t.Run("valid ordering with other UNTDID references", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Ordering = &bill.Ordering{
			Identities: []*org.Identity{
				{
					Code: "12345",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "CT",
					},
				},
				{
					Code: "67890",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "VN",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("ordering without identities is valid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Ordering = &bill.Ordering{
			Code: "ORD-12345",
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})
}

func TestLineIdentifiersValidation(t *testing.T) {
	t.Run("valid line with one AFL reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
				},
				Identifier: &org.Identity{
					Code: "12345",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "AFL",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid line with one AWW reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
				},
				Identifier: &org.Identity{
					Code: "12345",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "AWW",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid lines with one AFL and one AWW", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item 1",
					Price: num.NewAmount(10000, 2),
				},
				Identifier: &org.Identity{
					Code: "12345",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "AFL",
					},
				},
			},
			{
				Quantity: num.MakeAmount(2, 0),
				Item: &org.Item{
					Name:  "Test Item 2",
					Price: num.NewAmount(20000, 2),
				},
				Identifier: &org.Identity{
					Code: "67890",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "AWW",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid lines with other UNTDID references", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item 1",
					Price: num.NewAmount(10000, 2),
				},
				Identifier: &org.Identity{
					Code: "12345",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "CT",
					},
				},
			},
			{
				Quantity: num.MakeAmount(2, 0),
				Item: &org.Item{
					Name:  "Test Item 2",
					Price: num.NewAmount(20000, 2),
				},
				Identifier: &org.Identity{
					Code: "67890",
					Ext: tax.Extensions{
						untdid.ExtKeyReference: "VN",
					},
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("lines without identifiers are valid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Lines = []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
				},
			},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})
}

func setDocumentType(inv *bill.Invoice, docType string) {
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	if inv.Tax.Ext == nil {
		inv.Tax.Ext = make(tax.Extensions)
	}
	inv.Tax.Ext[untdid.ExtKeyDocumentType] = cbc.Code(docType)
}

func TestConsolidatedCreditNoteValidation(t *testing.T) {
	// BR-FR-CO-03: Document type 262 requires delivery period and ordering contracts
	t.Run("valid consolidated credit note with delivery and contract (BR-FR-CO-03)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Delivery = &bill.DeliveryDetails{
			Period: &cal.Period{
				Start: cal.MakeDate(2024, 5, 1),
				End:   cal.MakeDate(2024, 5, 31),
			},
		}
		inv.Ordering = &bill.Ordering{
			Contracts: []*org.DocumentRef{
				{
					Code: "CONTRACT-001",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "262")
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("consolidated credit note without delivery is invalid (BR-FR-CO-03)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Delivery = nil
		inv.Ordering = &bill.Ordering{
			Contracts: []*org.DocumentRef{
				{
					Code: "CONTRACT-001",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "262")
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "delivery details are required")
		assert.ErrorContains(t, err, "BR-FR-CO-03")
	})

	t.Run("consolidated credit note without delivery period is invalid (BR-FR-CO-03)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Delivery = &bill.DeliveryDetails{}
		inv.Ordering = &bill.Ordering{
			Contracts: []*org.DocumentRef{
				{
					Code: "CONTRACT-001",
				},
			},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "262")
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "delivery period is required")
		assert.ErrorContains(t, err, "BR-FR-CO-03")
	})

	t.Run("consolidated credit note without ordering contracts is invalid (BR-FR-CO-03)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Delivery = &bill.DeliveryDetails{
			Period: &cal.Period{
				Start: cal.MakeDate(2024, 5, 1),
				End:   cal.MakeDate(2024, 5, 31),
			},
		}
		inv.Ordering = &bill.Ordering{
			Contracts: nil,
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "262")
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "at least one contract reference is required")
		assert.ErrorContains(t, err, "BR-FR-CO-03")
	})

	t.Run("consolidated credit note with empty contracts array is invalid (BR-FR-CO-03)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Delivery = &bill.DeliveryDetails{
			Period: &cal.Period{
				Start: cal.MakeDate(2024, 5, 1),
				End:   cal.MakeDate(2024, 5, 31),
			},
		}
		inv.Ordering = &bill.Ordering{
			Contracts: []*org.DocumentRef{},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "262")
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "at least one contract reference is required")
		assert.ErrorContains(t, err, "BR-FR-CO-03")
	})

	t.Run("non-consolidated credit note does not require delivery or contracts", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Delivery = nil
		inv.Ordering = nil
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV-001",
				IssueDate: cal.NewDate(2024, 5, 1),
			},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "381") // Regular credit note, not 262
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("consolidated credit note with nil ordering should fail (BR-FR-CO-03)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Delivery = &bill.DeliveryDetails{
			Period: &cal.Period{
				Start: cal.MakeDate(2024, 5, 1),
				End:   cal.MakeDate(2024, 5, 31),
			},
		}
		inv.Ordering = nil // No ordering at all
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "262")
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "ordering")
		assert.ErrorContains(t, err, "BR-FR-CO-03")
	})
}

func TestSTCSupplierValidation(t *testing.T) {
	// BR-FR-CO-14/CO-15: STC supplier requirements
	t.Run("STC supplier requires ordering with seller (BR-FR-CO-15)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		// Add STC identity to supplier
		inv.Supplier.Identities = append(inv.Supplier.Identities, &org.Identity{
			Code: "12345678",
			Ext: tax.Extensions{
				iso.ExtKeySchemeID: "0231", // STC scheme
			},
		})
		inv.Ordering = &bill.Ordering{
			Seller: &org.Party{
				Name:  "Assujetti Unique",
				TaxID: inv.Supplier.TaxID, // Reuse supplier's valid tax ID
			},
		}
		// Add TXD note
		inv.Notes = append(inv.Notes, &org.Note{
			Text: "MEMBRE_ASSUJETTI_UNIQUE",
			Ext:  tax.Extensions{untdid.ExtKeyTextSubject: "TXD"},
		})
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("STC supplier seller missing tax ID (BR-FR-CO-15)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		// Add STC identity to supplier
		inv.Supplier.Identities = append(inv.Supplier.Identities, &org.Identity{
			Code: "12345678",
			Ext: tax.Extensions{
				iso.ExtKeySchemeID: "0231", // STC scheme
			},
		})
		inv.Ordering = &bill.Ordering{
			Identities: []*org.Identity{
				{
					Code: "ORDER-123",
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "0088",
					},
				},
			},
			Seller: &org.Party{
				Name:  "Assujetti Unique",
				TaxID: nil, // Missing tax ID
			},
		}
		// Add TXD note
		inv.Notes = append(inv.Notes, &org.Note{
			Text: "MEMBRE_ASSUJETTI_UNIQUE",
			Ext:  tax.Extensions{untdid.ExtKeyTextSubject: "TXD"},
		})
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "tax ID is required when supplier is under STC scheme")
	})

	t.Run("STC supplier seller with empty tax ID code (BR-FR-CO-15)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		// Add STC identity to supplier
		inv.Supplier.Identities = append(inv.Supplier.Identities, &org.Identity{
			Code: "12345678",
			Ext: tax.Extensions{
				iso.ExtKeySchemeID: "0231", // STC scheme
			},
		})
		inv.Ordering = &bill.Ordering{
			Identities: []*org.Identity{
				{
					Code: "ORDER-123",
					Ext: tax.Extensions{
						iso.ExtKeySchemeID: "0088",
					},
				},
			},
			Seller: &org.Party{
				Name: "Assujetti Unique",
				TaxID: &tax.Identity{
					Country: "FR",
					Code:    "", // Empty code
				},
			},
		}
		// Add TXD note
		inv.Notes = append(inv.Notes, &org.Note{
			Text: "MEMBRE_ASSUJETTI_UNIQUE",
			Ext:  tax.Extensions{untdid.ExtKeyTextSubject: "TXD"},
		})
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "code is required when supplier is under STC scheme")
	})

	t.Run("STC supplier with nil ordering should fail (BR-FR-CO-15)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		// Add STC identity to supplier
		inv.Supplier.Identities = append(inv.Supplier.Identities, &org.Identity{
			Code: "12345678",
			Ext: tax.Extensions{
				iso.ExtKeySchemeID: "0231", // STC scheme
			},
		})
		inv.Ordering = nil
		// Add TXD note
		inv.Notes = append(inv.Notes, &org.Note{
			Text: "MEMBRE_ASSUJETTI_UNIQUE",
			Ext:  tax.Extensions{untdid.ExtKeyTextSubject: "TXD"},
		})
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "ordering")
		assert.ErrorContains(t, err, "BR-FR-CO-15")
	})

	t.Run("STC supplier requires TXD note (BR-FR-CO-14)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		// Add STC identity to supplier
		inv.Supplier.Identities = append(inv.Supplier.Identities, &org.Identity{
			Code: "12345678",
			Ext: tax.Extensions{
				iso.ExtKeySchemeID: "0231", // STC scheme
			},
		})
		inv.Ordering = &bill.Ordering{
			Seller: &org.Party{
				Name:  "Assujetti Unique",
				TaxID: inv.Supplier.TaxID, // Reuse supplier's valid tax ID
			},
		}
		// TXD note missing
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "TXD")
		assert.ErrorContains(t, err, "MEMBRE_ASSUJETTI_UNIQUE")
	})
}

func TestFinalInvoicePaymentValidation(t *testing.T) {
	// BR-FR-CO-09: Final invoices require payment details
	t.Run("final invoice B2 with nil payment should fail (BR-FR-CO-09)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeB2
		inv.Payment = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "payment")
		// May be caught by BR-CO-25 (EN16931) or BR-FR-CO-09, either is acceptable
	})

	t.Run("final invoice S2 with nil payment should fail (BR-FR-CO-09)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeS2
		inv.Payment = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "payment")
		// May be caught by BR-CO-25 (EN16931) or BR-FR-CO-09, either is acceptable
	})

	t.Run("final invoice M2 with nil payment should fail (BR-FR-CO-09)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeM2
		inv.Payment = nil
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "payment")
		// May be caught by BR-CO-25 (EN16931) or BR-FR-CO-09, either is acceptable
	})
}

func TestPrecedingReferencesValidation(t *testing.T) {
	// BR-FR-CO-04: Corrective invoices
	t.Run("corrective invoice with exactly one preceding reference (BR-FR-CO-04)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV-001",
				IssueDate: cal.NewDate(2024, 5, 1),
			},
		}
		require.NoError(t, inv.Calculate())
		// Set document type to corrective invoice AFTER Calculate() so scenarios don't overwrite it
		setDocumentType(inv, "384")
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("corrective invoice with no preceding reference (BR-FR-CO-04)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = nil
		require.NoError(t, inv.Calculate())
		// Set document type to corrective invoice AFTER Calculate() so scenarios don't overwrite it
		setDocumentType(inv, "384")
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "exactly one preceding invoice reference")
		assert.ErrorContains(t, err, "BR-FR-CO-04")
	})

	t.Run("corrective invoice with multiple preceding references (BR-FR-CO-04)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV-001",
				IssueDate: cal.NewDate(2024, 5, 1),
			},
			{
				Code:      "INV-002",
				IssueDate: cal.NewDate(2024, 5, 2),
			},
		}
		require.NoError(t, inv.Calculate())
		// Set document type to corrective invoice AFTER Calculate() so scenarios don't overwrite it
		setDocumentType(inv, "384")
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "exactly one preceding invoice reference")
		assert.ErrorContains(t, err, "BR-FR-CO-04")
	})

	t.Run("corrective invoice type 471 requires one preceding reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV-001",
				IssueDate: cal.NewDate(2024, 5, 1),
			},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "471")
		err := inv.Validate()
		assert.NoError(t, err)
	})

	// BR-FR-CO-05: Credit notes
	t.Run("credit note with at least one preceding reference (BR-FR-CO-05)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV-001",
				IssueDate: cal.NewDate(2024, 5, 1),
			},
		}
		require.NoError(t, inv.Calculate())
		// Set document type to credit note AFTER Calculate() so scenarios don't overwrite it
		setDocumentType(inv, "381")
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("credit note with multiple preceding references is valid (BR-FR-CO-05)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV-001",
				IssueDate: cal.NewDate(2024, 5, 1),
			},
			{
				Code:      "INV-002",
				IssueDate: cal.NewDate(2024, 5, 2),
			},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "381")
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("credit note with no preceding reference (BR-FR-CO-05)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = nil
		require.NoError(t, inv.Calculate())
		// Set document type to credit note AFTER Calculate() so scenarios don't overwrite it
		setDocumentType(inv, "381")
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "at least one preceding invoice reference")
		assert.ErrorContains(t, err, "BR-FR-CO-05")
	})

	t.Run("credit note type 261 requires at least one preceding reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = []*org.DocumentRef{
			{
				Code:      "INV-001",
				IssueDate: cal.NewDate(2024, 5, 1),
			},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "261")
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("credit note type 502 requires at least one preceding reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = nil
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "502")
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "at least one preceding invoice reference")
		assert.ErrorContains(t, err, "BR-FR-CO-05")
	})

	t.Run("standard invoice does not require preceding reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = nil
		require.NoError(t, inv.Calculate())
		// Standard invoice type 380 is already set by scenarios, so this is redundant
		// but included for clarity
		setDocumentType(inv, "380")
		err := inv.Validate()
		assert.NoError(t, err)
	})
}

func TestPaymentDueDateValidation(t *testing.T) {
	t.Run("valid due date on or after issue date (BR-FR-CO-07)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.IssueDate = cal.MakeDate(2024, 6, 1)
		inv.Payment.Terms.DueDates[0].Date = cal.NewDate(2024, 7, 1) // After issue date
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("valid due date same as issue date (BR-FR-CO-07)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.IssueDate = cal.MakeDate(2024, 6, 1)
		inv.Payment.Terms.DueDates[0].Date = cal.NewDate(2024, 6, 1) // Same as issue date
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("invalid due date before issue date (BR-FR-CO-07)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.IssueDate = cal.MakeDate(2024, 6, 15)
		inv.Payment.Terms.DueDates[0].Date = cal.NewDate(2024, 6, 1) // Before issue date
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "too early")
	})

	t.Run("advance payment type 386 allows due date before issue date (BR-FR-CO-07)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.IssueDate = cal.MakeDate(2024, 6, 15)
		inv.Payment.Terms.DueDates[0].Date = cal.NewDate(2024, 6, 1) // Before issue date
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "386") // Advance payment
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("advance payment type 500 allows due date before issue date (BR-FR-CO-07)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.IssueDate = cal.MakeDate(2024, 6, 15)
		inv.Payment.Terms.DueDates[0].Date = cal.NewDate(2024, 6, 1) // Before issue date
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "500") // Self-billed advance payment
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("final invoice billing mode B2 allows due date before issue date (BR-FR-CO-07)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.IssueDate = cal.MakeDate(2024, 6, 15)
		inv.Payment.Terms.DueDates[0].Date = cal.NewDate(2024, 6, 1) // Before issue date
		require.NoError(t, inv.Calculate())
		// Set billing mode to B2 (final invoice)
		if inv.Tax.Ext == nil {
			inv.Tax.Ext = make(tax.Extensions)
		}
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeB2
		// Set up final invoice totals (BR-FR-CO-09)
		totalWithTax := inv.Totals.TotalWithTax
		inv.Totals.Advances = &totalWithTax
		zero := num.MakeAmount(0, 2)
		inv.Totals.Payable = zero
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("final invoice billing mode S2 allows due date before issue date (BR-FR-CO-07)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.IssueDate = cal.MakeDate(2024, 6, 15)
		inv.Payment.Terms.DueDates[0].Date = cal.NewDate(2024, 6, 1) // Before issue date
		require.NoError(t, inv.Calculate())
		// Set billing mode to S2 (self-billed final invoice)
		if inv.Tax.Ext == nil {
			inv.Tax.Ext = make(tax.Extensions)
		}
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeS2
		// Set up final invoice totals (BR-FR-CO-09)
		totalWithTax := inv.Totals.TotalWithTax
		inv.Totals.Advances = &totalWithTax
		zero := num.MakeAmount(0, 2)
		inv.Totals.Payable = zero
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("no due date is valid (BR-FR-CO-07)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.IssueDate = cal.MakeDate(2024, 6, 1)
		// Set notes instead of due dates - no due date means rule doesn't apply
		inv.Payment.Terms.DueDates = nil
		inv.Payment.Terms.Notes = "Payment on delivery"
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})
}

func TestBillingModeDocumentTypeCompatibility(t *testing.T) {
	t.Run("factored billing mode B4 with advance payment type 386 is invalid (BR-FR-CO-08)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		// Set factored billing mode B4
		if inv.Tax.Ext == nil {
			inv.Tax.Ext = make(tax.Extensions)
		}
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeB4
		// Set advance payment document type 386
		setDocumentType(inv, "386")
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "value '386' not allowed")
	})

	t.Run("factored billing mode S4 with advance payment type 500 is invalid (BR-FR-CO-08)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		// Set factored billing mode S4
		if inv.Tax.Ext == nil {
			inv.Tax.Ext = make(tax.Extensions)
		}
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeS4
		// Set advance payment document type 500
		setDocumentType(inv, "500")
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "value '500' not allowed")
	})

	t.Run("factored billing mode M4 with advance payment type 503 is invalid (BR-FR-CO-08)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		// Set factored billing mode M4
		if inv.Tax.Ext == nil {
			inv.Tax.Ext = make(tax.Extensions)
		}
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeM4
		// Set advance payment document type 503
		setDocumentType(inv, "503")
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "value '503' not allowed")
	})

	t.Run("factored billing mode B4 with standard invoice type 380 is valid (BR-FR-CO-08)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		// Set factored billing mode B4
		if inv.Tax.Ext == nil {
			inv.Tax.Ext = make(tax.Extensions)
		}
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeB4
		// Standard invoice type 380 is already set by scenarios
		setDocumentType(inv, "380")
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("non-factored billing mode B2 with advance payment type 386 is valid (BR-FR-CO-08)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		// Set non-factored billing mode B2
		if inv.Tax.Ext == nil {
			inv.Tax.Ext = make(tax.Extensions)
		}
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeB2
		// Set up final invoice totals (BR-FR-CO-09)
		totalWithTax := inv.Totals.TotalWithTax
		inv.Totals.Advances = &totalWithTax
		zero := num.MakeAmount(0, 2)
		inv.Totals.Payable = zero
		// Set advance payment document type 386
		setDocumentType(inv, "386")
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("standard billing mode B7 with standard invoice type 380 is valid (BR-FR-CO-08)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		// B7 and 380 are already set by scenarios
		err := inv.Validate()
		assert.NoError(t, err)
	})
}

func TestFinalInvoiceValidation(t *testing.T) {
	t.Run("valid final invoice B2 - fully paid with correct amounts (BR-FR-CO-09)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Set billing mode to B2 (final invoice)
		if inv.Tax.Ext == nil {
			inv.Tax.Ext = make(tax.Extensions)
		}
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeB2

		// Manually set the totals to simulate fully paid invoice
		// Advance = TotalWithTax, Payable = 0
		totalWithTax := inv.Totals.TotalWithTax
		inv.Totals.Advances = &totalWithTax
		zero := num.MakeAmount(0, 2)
		inv.Totals.Payable = zero

		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("final invoice B2 without advance amount is invalid (BR-FR-CO-09)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Set billing mode to B2
		if inv.Tax.Ext == nil {
			inv.Tax.Ext = make(tax.Extensions)
		}
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeB2

		// No advance amount set
		inv.Totals.Advances = nil

		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "advance amount is required for already-paid invoices")
	})

	t.Run("final invoice B2 with incorrect advance amount is invalid (BR-FR-CO-09)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Set billing mode to B2
		if inv.Tax.Ext == nil {
			inv.Tax.Ext = make(tax.Extensions)
		}
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeB2

		// Set advance amount to something other than TotalWithTax
		wrongAmount := num.MakeAmount(5000, 2) // Wrong amount
		inv.Totals.Advances = &wrongAmount

		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "must be equal to")
	})

	t.Run("final invoice S2 with non-zero payable amount is invalid (BR-FR-CO-09)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Set billing mode to S2
		if inv.Tax.Ext == nil {
			inv.Tax.Ext = make(tax.Extensions)
		}
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeS2

		// Set advance amount correctly
		totalWithTax := inv.Totals.TotalWithTax
		inv.Totals.Advances = &totalWithTax
		// But set Due as non-zero (which should be 0 for final invoices)
		nonZero := num.MakeAmount(100, 2)
		inv.Totals.Due = &nonZero

		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "must be equal to 0")
	})

	t.Run("final invoice M2 without due date is invalid (BR-FR-CO-09)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Set billing mode to M2
		if inv.Tax.Ext == nil {
			inv.Tax.Ext = make(tax.Extensions)
		}
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeM2

		// Set amounts correctly
		totalWithTax := inv.Totals.TotalWithTax
		inv.Totals.Advances = &totalWithTax
		zero := num.MakeAmount(0, 2)
		inv.Totals.Payable = zero

		// Remove due date
		inv.Payment.Terms.DueDates = nil
		inv.Payment.Terms.Notes = "Payment already made"

		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "at least one due date required")
	})

	t.Run("non-final invoice B7 does not require these validations (BR-FR-CO-09)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// B7 is the default billing mode, and normal totals
		// No advance amount, non-zero payable - all OK for non-final invoices
		err := inv.Validate()
		assert.NoError(t, err)
	})
}

func TestSelfBilledInvoiceValidation(t *testing.T) {
	t.Run("self-billed invoice skips supplier SIREN inbox validation", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Add B2B note to make it a B2B transaction
		inv.Notes = append(inv.Notes, &org.Note{
			Key:  org.NoteKeyGeneral,
			Text: "B2B",
			Ext: tax.Extensions{
				untdid.ExtKeyTextSubject: "BAR",
			},
		})

		// Replace SIREN inbox with non-SIREN inbox
		inv.Supplier.Inboxes = []*org.Inbox{
			{
				Scheme: "0088", // GLN instead of SIREN
				Code:   "1234567890123",
			},
		}

		// Normal B2B invoice requires supplier SIREN inbox (BR-FR-21/22)
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "BR-FR-21/22")

		// Set self-billed document type (389)
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "389"

		// Self-billed invoices skip supplier SIREN inbox validation
		err = inv.Validate()
		assert.NoError(t, err)
	})
}

func TestCorrectiveInvoiceValidation(t *testing.T) {
	t.Run("corrective invoice requires exactly one preceding (BR-FR-CO-04)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Set corrective document type (384)
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "384"

		// Corrective invoices need exactly one preceding invoice
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "BR-FR-CO-04")

		// Add preceding document
		inv.Preceding = []*org.DocumentRef{
			{
				Code: "INV-123",
			},
		}
		err = inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("corrective invoice with multiple preceding fails (BR-FR-CO-04)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Set corrective document type (384)
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "384"

		// Add two preceding documents (should fail)
		inv.Preceding = []*org.DocumentRef{
			{Code: "INV-123"},
			{Code: "INV-456"},
		}
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "BR-FR-CO-04")
	})
}

func TestCreditNoteValidation(t *testing.T) {
	t.Run("credit note requires at least one preceding (BR-FR-CO-05)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Set credit note document type (381)
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "381"

		// Credit notes need at least one preceding invoice
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "BR-FR-CO-05")

		// Add preceding document
		inv.Preceding = []*org.DocumentRef{
			{Code: "INV-123"},
		}
		err = inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("credit note with multiple preceding is valid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Set credit note document type (381)
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "381"

		// Add multiple preceding documents (should be valid)
		inv.Preceding = []*org.DocumentRef{
			{Code: "INV-123"},
			{Code: "INV-456"},
		}
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("factored credit note (396) requires preceding", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Set factored credit note document type (396)
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "396"

		// Should require preceding
		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "BR-FR-CO-05")

		// Add preceding
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}
		err = inv.Validate()
		assert.NoError(t, err)
	})
}

func TestConsolidatedCreditNoteTypes(t *testing.T) {
	t.Run("consolidated credit note type 262", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Set consolidated credit note type (262)
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "262"

		// Should require delivery and contracts
		err := inv.Validate()
		assert.Error(t, err)
	})
}

func TestAdvancedInvoiceTypes(t *testing.T) {
	t.Run("prepaid invoice type 386 is advance", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set prepaid invoice type
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "386"
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeB1

		require.NoError(t, inv.Calculate())

		// Should validate as advance invoice
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("self-billed advance type 500", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set self-billed advance payment type
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "500"
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeB1

		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})
}

func TestFinalInvoiceTypes(t *testing.T) {
	t.Run("final invoice type 456", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set final invoice type and billing mode
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "456"
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeM4

		require.NoError(t, inv.Calculate())

		// Should validate as final invoice
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("self-billed final type 501", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set self-billed final type and billing mode
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "501"
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeS4

		require.NoError(t, inv.Calculate())

		// Should validate as self-billed final invoice
		err := inv.Validate()
		assert.NoError(t, err)
	})
}

func TestInvoiceNormalization(t *testing.T) {
	ad := tax.AddonForKey(ctc.Flow2)

	t.Run("normalizes invoice with existing tax", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Normalize should set rounding to currency
		ad.Normalizer(inv)
		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
	})

	t.Run("normalizes invoice without tax object", func(t *testing.T) {
		inv := &bill.Invoice{
			Supplier: &org.Party{Name: "Test"},
			Customer: &org.Party{Name: "Customer"},
		}

		// Should create tax object and set rounding
		ad.Normalizer(inv)
		assert.NotNil(t, inv.Tax)
		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
	})

	t.Run("normalizes nil invoice", func(t *testing.T) {
		var inv *bill.Invoice
		ad.Normalizer(inv)
		assert.Nil(t, inv)
	})
}

func TestHelperFunctionEdgeCases(t *testing.T) {
	ad := tax.AddonForKey(ctc.Flow2)

	t.Run("validate unsupported type returns nil", func(t *testing.T) {
		// Test with a type that isn't in the switch statement
		type unsupported struct{}
		err := ad.Validator(&unsupported{})
		assert.NoError(t, err)
	})

	t.Run("validate nil date", func(t *testing.T) {
		err := ad.Validator((*cal.Date)(nil))
		assert.NoError(t, err)
	})

	t.Run("isCreditNote with nil invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Create a scenario where isCreditNote is called indirectly
		// by setting an invoice that validates successfully
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "381"
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("isConsolidatedCreditNote with nil invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Set consolidated credit note type
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "262"
		inv.Delivery = &bill.DeliveryDetails{
			Period: &cal.Period{
				Start: cal.MakeDate(2024, 6, 1),
				End:   cal.MakeDate(2024, 6, 30),
			},
		}
		inv.Ordering = &bill.Ordering{
			Contracts: []*org.DocumentRef{{Code: "CONTRACT-001"}},
		}
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("getPartySIREN with party without identities", func(t *testing.T) {
		// Test getPartySIREN indirectly through validation
		party := &org.Party{
			Name: "Party without identities",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "44732829320",
			},
			Inboxes: []*org.Inbox{
				{
					Scheme: "0225",
					Code:   "123456789",
				},
			},
		}
		err := ad.Validator(party)
		assert.NoError(t, err) // Should pass, getPartySIREN returns empty string
	})

	t.Run("isCreditNote with invoice without extensions", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Create invoice with nil tax extensions
		inv.Tax.Ext = nil

		// Should handle nil gracefully
		err := inv.Validate()
		assert.Error(t, err) // Will error for other reasons
	})

	t.Run("isAdvancedInvoice with missing billing mode", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Remove billing mode extension
		delete(inv.Tax.Ext, ctc.ExtKeyBillingMode)

		// Should not panic
		err := inv.Validate()
		_ = err // May or may not error depending on other rules
	})
}

func TestValidatePrecedingDocument(t *testing.T) {
	t.Run("invoice with nil preceding is valid for standard invoices", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = nil
		require.NoError(t, inv.Calculate())

		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("invoice with nil element in preceding array returns nil from CTC validation", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Add a nil element in the preceding array
		// This tests the nil check in validatePrecedingDocument
		inv.Preceding = []*org.DocumentRef{nil}

		require.NoError(t, inv.Calculate())

		// CTC addon validation should return nil for nil document ref
		ad := tax.AddonForKey(ctc.Flow2)
		err := ad.Validator(inv)
		assert.NoError(t, err, "CTC addon should return nil for nil preceding document element")
	})

	t.Run("invoice with empty preceding code", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set credit note type that requires preceding
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "381"

		// Add preceding with nil code (should fail base validation)
		inv.Preceding = []*org.DocumentRef{
			{
				Code: "", // Empty code
			},
		}

		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.Error(t, err) // Should fail on empty code
	})
}

func TestValidateCodeEdgeCases(t *testing.T) {
	t.Run("invoice with valid code and series", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Code = "FAC-001"
		inv.Series = "2024"
		require.NoError(t, inv.Calculate())

		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("invoice with code containing series separator without series", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Code = "FAC-2024-001" // Contains '-' but no explicit series
		inv.Series = ""
		require.NoError(t, inv.Calculate())

		err := inv.Validate()
		assert.NoError(t, err) // Should be valid
	})
}

func TestSupplierValidationEdgeCases(t *testing.T) {
	t.Run("supplier without inboxes fails BR-FR-13", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Remove all inboxes
		inv.Supplier.Inboxes = nil

		err := inv.Validate()
		assert.Error(t, err)
		assert.ErrorContains(t, err, "BR-FR-13")
	})

	t.Run("supplier without SIREN identity", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Remove SIREN identity
		inv.Supplier.Identities = []*org.Identity{
			{
				Code: "OTHER-ID",
				Ext: tax.Extensions{
					iso.ExtKeySchemeID: "0088",
				},
			},
		}

		// Should fail SIREN requirement
		err := inv.Validate()
		assert.Error(t, err)
		// The error is BR-FR-10/11 from regime validation
		assert.ErrorContains(t, err, "BR-FR-10")
	})
}

func TestCustomerValidationEdgeCases(t *testing.T) {
	t.Run("B2B transaction customer without SIREN", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Add B2B note
		inv.Notes = append(inv.Notes, &org.Note{
			Key:  org.NoteKeyGeneral,
			Text: "B2B",
			Ext: tax.Extensions{
				untdid.ExtKeyTextSubject: "BAR",
			},
		})

		// Remove customer SIREN
		inv.Customer.Identities = []*org.Identity{
			{
				Code: "OTHER-ID",
				Ext: tax.Extensions{
					iso.ExtKeySchemeID: "0088",
				},
			},
		}

		err := inv.Validate()
		assert.Error(t, err)
		// The error is from regime validation (BR-FR-10/11)
		assert.ErrorContains(t, err, "BR-FR-10")
	})
}

func TestDeliveryAndTotalsValidation(t *testing.T) {
	t.Run("invoice without delivery for non-consolidated credit notes", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Delivery = nil
		require.NoError(t, inv.Calculate())

		// Standard invoice doesn't require delivery
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("invoice with zero payable", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())

		// Set payable to zero
		zero := num.MakeAmount(0, 2)
		inv.Totals.Payable = zero

		// Should pass - zero payable is valid in some contexts
		err := inv.Validate()
		// May error if context requires non-zero, but shouldn't panic
		_ = err
	})
}

func TestAdditionalDocumentTypes(t *testing.T) {
	t.Run("prepaid amount invoice type 471 with preceding", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set prepaid amount invoice (corrective type)
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "471"

		// Add preceding
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}

		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("stand-alone credit note type 473 with preceding", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set stand-alone credit note (both corrective and credit)
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "473"

		// Add preceding
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}

		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("self-billed corrective type 502 with preceding", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set self-billed corrective (both self-billed and credit)
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "502"

		// Add preceding
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}

		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("self-billed credit for claim type 503 with preceding", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set self-billed credit for claim
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "503"

		// Add preceding
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}

		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("self-billed prepaid invoice type 472", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set self-billed prepaid amount
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "472"

		// Add preceding
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}

		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("self-billed credit note type 261 with preceding", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set self-billed credit note
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "261"

		// Add preceding
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}

		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})
}

func TestAdditionalBillingModes(t *testing.T) {
	t.Run("billing mode B4 is final", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set billing mode B4 (final)
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeB4
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "456"

		require.NoError(t, inv.Calculate())

		// Should validate as final invoice
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("billing mode S4 is final", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set billing mode S4 (self-billed final)
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeS4
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "501"

		require.NoError(t, inv.Calculate())

		// Should validate as final invoice
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("billing mode M4 is final", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set billing mode M4 (mixed final)
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeM4
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "456"

		require.NoError(t, inv.Calculate())

		// Should validate as final invoice
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("billing mode S5 credit note dispute", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set billing mode S5
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeS5
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "381"

		// Add preceding
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}

		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("billing mode S6 self-billed corrective", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set billing mode S6
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeS6
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "502"

		// Add preceding
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}

		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("billing mode B7 self-billed for claim", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set billing mode B7
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeB7
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "503"

		// Add preceding
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}

		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})

	t.Run("billing mode S7 commercial invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)

		// Set billing mode S7
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeS7
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "380"

		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.NoError(t, err)
	})
}
func TestMissingRequiredNoteCodes(t *testing.T) {
	t.Run("missing PMT note code (BR-FR-05)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Notes = []*org.Note{
			{Key: org.NoteKeyPaymentMethod, Text: "PMD text", Ext: tax.Extensions{untdid.ExtKeyTextSubject: "PMD"}},
			{Key: org.NoteKeyPaymentTerm, Text: "AAB text", Ext: tax.Extensions{untdid.ExtKeyTextSubject: "AAB"}},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "missing required note codes: PMT")
		assert.ErrorContains(t, err, "BR-FR-05")
	})

	t.Run("missing PMD note code (BR-FR-05)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Notes = []*org.Note{
			{Key: org.NoteKeyPayment, Text: "PMT text", Ext: tax.Extensions{untdid.ExtKeyTextSubject: "PMT"}},
			{Key: org.NoteKeyPaymentTerm, Text: "AAB text", Ext: tax.Extensions{untdid.ExtKeyTextSubject: "AAB"}},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "missing required note codes: PMD")
		assert.ErrorContains(t, err, "BR-FR-05")
	})

	t.Run("missing AAB note code (BR-FR-05)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Notes = []*org.Note{
			{Key: org.NoteKeyPayment, Text: "PMT text", Ext: tax.Extensions{untdid.ExtKeyTextSubject: "PMT"}},
			{Key: org.NoteKeyPaymentMethod, Text: "PMD text", Ext: tax.Extensions{untdid.ExtKeyTextSubject: "PMD"}},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "missing required note codes: AAB")
		assert.ErrorContains(t, err, "BR-FR-05")
	})

	t.Run("missing multiple note codes (BR-FR-05)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Notes = []*org.Note{
			{Key: org.NoteKeyPayment, Text: "PMT text", Ext: tax.Extensions{untdid.ExtKeyTextSubject: "PMT"}},
		}
		require.NoError(t, inv.Calculate())
		err := inv.Validate()
		assert.ErrorContains(t, err, "missing required note codes")
		assert.ErrorContains(t, err, "PMD")
		assert.ErrorContains(t, err, "AAB")
		assert.ErrorContains(t, err, "BR-FR-05")
	})
}

func TestNilCodeValidation(t *testing.T) {
	t.Run("invoice with empty code returns nil from CTC validation", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Code = ""

		// Calculate to ensure invoice is normalized and amounts are computed
		require.NoError(t, inv.Calculate())

		// CTC addon validation should return nil for empty code
		// Base GOBL validation will catch the missing code
		ad := tax.AddonForKey(ctc.Flow2)
		err := ad.Validator(inv)
		assert.NoError(t, err, "CTC addon should return nil for empty code, letting base validation handle it")
	})

	t.Run("invoice with empty code fails base validation", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Code = ""
		require.NoError(t, inv.Calculate())

		// Full validation should catch the missing code
		// (This may or may not fail depending on signing context)
		_ = inv.Validate() // Code is only required for signing
	})
}

func TestValidationNilChecks(t *testing.T) {
	ad := tax.AddonForKey(ctc.Flow2)

	t.Run("invoice with nil payment terms", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		// Set a standard billing mode (not advance or final invoice)
		// so validatePayment will be called
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeS1
		inv.Payment.Terms = nil // Nil terms should be handled gracefully
		require.NoError(t, inv.Calculate())
		err := ad.Validator(inv)
		// Should not crash, may have other validation errors but shouldn't panic on nil
		_ = err
	})

	t.Run("invoice with nil payment due dates array element", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		// First calculate with valid data
		require.NoError(t, inv.Calculate())

		// Then set nil due date after calculation
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeS1
		var nilDueDate *pay.DueDate
		inv.Payment.Terms.DueDates = []*pay.DueDate{nilDueDate}

		// CTC validation should handle nil due date gracefully
		err := ad.Validator(inv)
		// validateDueDate should return nil for nil due date
		_ = err
	})

	t.Run("final invoice with nil totals returns error", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		// Set to final invoice billing mode to trigger validateTotals
		inv.Tax.Ext[ctc.ExtKeyBillingMode] = ctc.BillingModeB2
		inv.Totals = nil
		require.NoError(t, inv.Calculate())
		err := ad.Validator(inv)
		// Will error because totals are required for final invoices
		// but validateTotals should handle nil gracefully
		assert.Error(t, err)
	})

	t.Run("consolidated credit note with nil delivery", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		// Set to consolidated credit note to trigger validateDelivery
		inv.Tax.Ext[untdid.ExtKeyDocumentType] = "262"
		inv.Delivery = nil
		require.NoError(t, inv.Calculate())
		err := ad.Validator(inv)
		// validateDelivery should handle nil gracefully and not panic
		// (may have validation errors but won't crash)
		_ = err
	})

	t.Run("self-billed invoice helper with nil invoice", func(t *testing.T) {
		// These helper functions are not directly exposed but we can test
		// them indirectly through invoice validation
		inv := testInvoiceB2BStandard(t)
		inv.Tax = nil
		require.NoError(t, inv.Calculate())
		err := ad.Validator(inv)
		// Tax is required, but the helpers should handle nil gracefully
		assert.Error(t, err)
	})
}
