package flow2

import (
	"testing"

	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/catalogues/dgfip"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// frPartyWithSIREN returns a French party with a SIREN identity.
func frPartyWithSIREN(name, taxCode, siren string) *org.Party {
	return &org.Party{
		Name: name,
		TaxID: &tax.Identity{
			Country: "FR",
			Code:    cbc.Code(taxCode),
		},
		Identities: []*org.Identity{
			{
				Type:  fr.IdentityTypeSIREN,
				Code:  cbc.Code(siren),
				Scope: org.IdentityScopeLegal,
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					iso.ExtKeySchemeID: identitySchemeIDSIREN,
				}),
			},
		},
		Addresses: []*org.Address{
			{
				Street:   "1 Rue",
				Code:     "75001",
				Locality: "Paris",
				Country:  "FR",
			},
		},
		Inboxes: []*org.Inbox{
			{
				Key:    org.InboxKeyPeppol,
				Scheme: cbc.Code("0225"),
				Code:   cbc.Code(siren),
			},
		},
	}
}

func testInvoiceB2BStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Regime:   tax.WithRegime("FR"),
		Addons:   tax.WithAddons(V1, en16931.V2017),
		Code:     "FAC-2024-001",
		Currency: "EUR",
		Type:     bill.InvoiceTypeStandard,
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				dgfip.ExtKeyBillingMode:   dgfip.BillingModeS1,
				untdid.ExtKeyDocumentType: "380",
			}),
		},
		Supplier:  frPartyWithSIREN("Supplier SARL", "39356000000", "356000000"),
		Customer:  frPartyWithSIREN("Customer SAS", "44732829320", "732829320"),
		IssueDate: cal.MakeDate(2024, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Service",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{Category: "VAT", Rate: "standard"},
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
						Name: "Supplier SARL",
					},
				},
			},
		},
		Notes: []*org.Note{
			{
				Key:  org.NoteKeyPayment,
				Text: "Conditions.",
				Ext:  tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "PMT"}),
			},
			{
				Key:  org.NoteKeyPaymentMethod,
				Text: "Penalties.",
				Ext:  tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "PMD"}),
			},
			{
				Key:  org.NoteKeyPaymentTerm,
				Text: "No early discount.",
				Ext:  tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "AAB"}),
			},
		},
	}
}

func TestInvoiceB2BHappyPath(t *testing.T) {
	inv := testInvoiceB2BStandard(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestInvoiceCodeFormatRejectsBadChars(t *testing.T) {
	inv := testInvoiceB2BStandard(t)
	inv.Code = "INVALID CODE WITH SPACE"
	assert.Error(t, rules.Validate(inv))
}

func TestInvoiceMissingNotesFails(t *testing.T) {
	inv := testInvoiceB2BStandard(t)
	inv.Notes = nil
	assert.Error(t, rules.Validate(inv))
}

func TestInvoiceMissingBillingModeFails(t *testing.T) {
	inv := testInvoiceB2BStandard(t)
	inv.Tax.Ext = inv.Tax.Ext.Delete(dgfip.ExtKeyBillingMode)
	assert.Error(t, rules.Validate(inv))
}

func TestNormalizeAddsRequiredNotes(t *testing.T) {
	inv := testInvoiceB2BStandard(t)
	inv.Notes = nil
	tax.Normalize([]tax.Normalizer{tax.AddonForKey(V1).Normalizer}, inv)
	assert.GreaterOrEqual(t, len(inv.Notes), 3)
}
