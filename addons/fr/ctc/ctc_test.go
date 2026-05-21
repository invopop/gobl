package ctc_test

import (
	"testing"

	"github.com/invopop/gobl/addons/fr/ctc"
	"github.com/invopop/gobl/addons/fr/ctc/flow10"
	"github.com/invopop/gobl/addons/fr/ctc/flow2"
	"github.com/invopop/gobl/addons/fr/ctc/flow6"
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
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func frParty(name, taxCode, siren string) *org.Party {
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
					iso.ExtKeySchemeID: "0002",
				}),
			},
		},
		Addresses: []*org.Address{{
			Street: "1 Rue", Code: "75001", Locality: "Paris", Country: "FR",
		}},
		Inboxes: []*org.Inbox{{
			Key: org.InboxKeyPeppol, Scheme: cbc.Code("0225"), Code: cbc.Code(siren),
		}},
	}
}

func deParty() *org.Party {
	return &org.Party{
		Name: "Kunde GmbH",
		TaxID: &tax.Identity{
			Country: "DE",
			Code:    "111111125",
		},
		Identities: []*org.Identity{{
			Code:  "DE111111125",
			Scope: org.IdentityScopeLegal,
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				iso.ExtKeySchemeID: "0223",
			}),
		}},
		Addresses: []*org.Address{{Country: "DE"}},
	}
}

func TestInvoiceTwoFrenchPartiesDispatchesFlow2(t *testing.T) {
	inv := &bill.Invoice{
		Regime:   tax.WithRegime("FR"),
		Addons:   tax.WithAddons(ctc.V1),
		Code:     "FAC-2024-001",
		Currency: "EUR",
		Type:     bill.InvoiceTypeStandard,
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				dgfip.ExtKeyBillingMode:   dgfip.BillingModeS1,
				untdid.ExtKeyDocumentType: "380",
			}),
		},
		Supplier:  frParty("Supplier SARL", "39356000000", "356000000"),
		Customer:  frParty("Customer SAS", "44732829320", "732829320"),
		IssueDate: cal.MakeDate(2024, 6, 13),
		Lines: []*bill.Line{{
			Quantity: num.MakeAmount(1, 0),
			Item:     &org.Item{Name: "Service", Price: num.NewAmount(100, 0)},
			Taxes:    tax.Set{{Category: "VAT", Rate: "standard"}},
		}},
		Payment: &bill.PaymentDetails{
			Terms: &pay.Terms{
				Key: pay.TermKeyDueDate,
				DueDates: []*pay.DueDate{{
					Date:    cal.NewDate(2024, 7, 13),
					Percent: num.NewPercentage(100, 3),
				}},
			},
			Instructions: &pay.Instructions{
				Key: pay.MeansKeyCreditTransfer,
				CreditTransfer: []*pay.CreditTransfer{{
					IBAN: "FR7630006000011234567890189",
					Name: "Supplier",
				}},
			},
		},
	}
	require.NoError(t, inv.Calculate())
	// Meta-addon should have appended flow2 (and via Requires also
	// eu-en16931-v2017).
	assert.Contains(t, inv.Addons.List, flow2.V1)
	// Flow2's normalizer adds the default required notes — verify it
	// actually ran in the same Calculate() pass.
	assert.GreaterOrEqual(t, len(inv.Notes), 3, "flow2 normalizer should have added default required notes")
}

func TestInvoiceCrossBorderDispatchesFlow10(t *testing.T) {
	inv := &bill.Invoice{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(ctc.V1),
		Code:      "INV-2026-001",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2026, 1, 15),
		Type:      bill.InvoiceTypeStandard,
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				dgfip.ExtKeyBillingMode: dgfip.BillingModeS1,
			}),
		},
		Supplier: frParty("Fournisseur", "39356000000", "356000000"),
		Customer: deParty(),
		Lines: []*bill.Line{{
			Quantity: num.MakeAmount(1, 0),
			Item:     &org.Item{Name: "Product", Price: num.NewAmount(100, 0)},
			Taxes:    tax.Set{{Category: tax.CategoryVAT, Percent: num.NewPercentage(20, 2)}},
		}},
	}
	require.NoError(t, inv.Calculate())
	assert.Contains(t, inv.Addons.List, flow10.V1)
	assert.NotContains(t, inv.Addons.List, flow2.V1, "cross-border invoice should not be flow2")
}

func TestInvoiceB2CDispatchesFlow10(t *testing.T) {
	inv := &bill.Invoice{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(ctc.V1),
		Code:      "INV-2026-B2C-001",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2026, 1, 15),
		Type:      bill.InvoiceTypeStandard,
		Supplier:  frParty("Fournisseur", "39356000000", "356000000"),
		Lines: []*bill.Line{{
			Quantity: num.MakeAmount(1, 0),
			Item:     &org.Item{Name: "Product", Price: num.NewAmount(100, 0)},
			Taxes:    tax.Set{{Category: tax.CategoryVAT, Percent: num.NewPercentage(20, 2)}},
		}},
	}
	require.NoError(t, inv.Calculate())
	assert.Contains(t, inv.Addons.List, flow10.V1)
	// Flow 10 normalizer should have defaulted the B2C category to TNT1
	// in the same Calculate() pass.
	assert.Equal(t, flow10.B2CCategoryNotTaxable, inv.Tax.Ext.Get(flow10.ExtKeyB2CCategory))
}

func TestStatusDispatchesFlow6(t *testing.T) {
	issued := cal.MakeDate(2026, 2, 1)
	st := &bill.Status{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(ctc.V1),
		IssueDate: cal.MakeDate(2026, 2, 2),
		Code:      "STA-2026-0001",
		Supplier: &org.Party{
			Name: "Platform",
			Identities: []*org.Identity{{
				Code: "356000000",
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					iso.ExtKeySchemeID: "0002",
				}),
			}},
		},
		Issuer: &org.Party{
			Name: "Acheteur",
			Identities: []*org.Identity{{
				Code: "200000008",
				Ext:  tax.MakeExtensions().Set(iso.ExtKeySchemeID, "0002"),
			}},
			Inboxes: []*org.Inbox{{Scheme: "0225", Code: "200000008_PEP"}},
		},
		Recipient: &org.Party{
			Name: "Vendeur",
			Identities: []*org.Identity{{
				Code: "356000000",
				Ext:  tax.MakeExtensions().Set(iso.ExtKeySchemeID, "0002"),
			}},
			Inboxes: []*org.Inbox{{Scheme: "0225", Code: "356000000_PEP"}},
		},
		Lines: []*bill.StatusLine{{
			Key:  bill.StatusEventAccepted,
			Date: &issued,
			Doc:  &org.DocumentRef{Code: "INV-2026-001", IssueDate: &issued},
		}},
	}
	require.NoError(t, st.Calculate())
	assert.Contains(t, st.Addons.List, flow6.V1)
	// Flow 6 normalizer should have derived the status type from the
	// line key (StatusEventAccepted → response).
	assert.Equal(t, bill.StatusTypeResponse, st.Type)
}

func TestPaymentB2CDispatchesFlow10(t *testing.T) {
	value := cal.MakeDate(2026, 1, 15)
	pmt := &bill.Payment{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(ctc.V1),
		Code:      "PMT-2026-001",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2026, 1, 15),
		Type:      bill.PaymentTypeReceipt,
		ValueDate: &value,
		Supplier:  frParty("Supplier", "39356000000", "356000000"),
		Methods:   []*pay.Record{{Key: pay.MeansKeyCreditTransfer}},
		Lines: []*bill.PaymentLine{{
			Amount: num.MakeAmount(10000, 2),
		}},
	}
	require.NoError(t, pmt.Calculate())
	assert.Contains(t, pmt.Addons.List, flow10.V1)
}

func TestPaymentReceiptBetweenFrenchPartiesDispatchesFlow6(t *testing.T) {
	value := cal.MakeDate(2026, 1, 15)
	pmt := &bill.Payment{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(ctc.V1),
		Code:      "PMT-2026-002",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2026, 1, 15),
		Type:      bill.PaymentTypeReceipt,
		ValueDate: &value,
		Supplier:  frParty("Vendeur", "39356000000", "356000000"),
		Customer:  frParty("Acheteur", "44732829320", "732829320"),
		Methods:   []*pay.Record{{Key: pay.MeansKeyCreditTransfer}},
		Lines: []*bill.PaymentLine{{
			Amount: num.MakeAmount(120000, 2),
			Document: &org.DocumentRef{
				Code:      "2026-00042",
				IssueDate: cal.NewDate(2026, 1, 1),
			},
		}},
	}
	require.NoError(t, pmt.Calculate())
	assert.Contains(t, pmt.Addons.List, flow6.V1)
	assert.NotContains(t, pmt.Addons.List, flow10.V1, "domestic French payment should not be flow10")
	// flow6 normalizer should have set the CDAR status code to 212.
	assert.Equal(t, cbc.Code("212"), pmt.Ext.Get("fr-ctc-flow6-status-code"))
}

func TestPaymentRequestBetweenFrenchPartiesDispatchesFlow10(t *testing.T) {
	// A "request" payment is not a CDV event — even between two
	// French parties it routes to flow10.
	value := cal.MakeDate(2026, 1, 15)
	pmt := &bill.Payment{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(ctc.V1),
		Code:      "PMT-2026-003",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2026, 1, 15),
		Type:      bill.PaymentTypeRequest,
		ValueDate: &value,
		Supplier:  frParty("Vendeur", "39356000000", "356000000"),
		Customer:  frParty("Acheteur", "44732829320", "732829320"),
		Methods:   []*pay.Record{{Key: pay.MeansKeyCreditTransfer}},
		Lines: []*bill.PaymentLine{{
			Amount: num.MakeAmount(120000, 2),
		}},
	}
	require.NoError(t, pmt.Calculate())
	assert.Contains(t, pmt.Addons.List, flow10.V1)
	assert.NotContains(t, pmt.Addons.List, flow6.V1, "flow6 should not handle request payments")
}
