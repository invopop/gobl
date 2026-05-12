package ctc

import (
	"testing"

	_ "github.com/invopop/gobl/addons/eu/en16931" // Flow 2 ruleset depends on en16931 being registered as an addon.
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// -- Fixtures -------------------------------------------------------------

// testInvoiceB2BStandard mirrors the old flow2 fixture but adds the
// eu-en16931-v2017 addon (now a soft Flow 2 requirement enforced by
// rule, not by addon.Requires) and carries the iso-scheme-id ext on
// identities so org_party rule 03 (identitiesSchemeFormatValid) accepts
// them.
func testInvoiceB2BStandard(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Regime:   tax.WithRegime("FR"),
		Addons:   tax.WithAddons(V1, "eu-en16931-v2017"),
		Code:     "FAC-2024-001",
		Currency: "EUR",
		Type:     bill.InvoiceTypeStandard,
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyBillingMode:         BillingModeS1,
				untdid.ExtKeyDocumentType: "380",
			}),
		},
		Supplier: &org.Party{
			Name: "Test Supplier SARL",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "39356000000",
			},
			Identities: []*org.Identity{
				{
					Type:  fr.IdentityTypeSIREN,
					Code:  "356000000",
					Scope: org.IdentityScopeLegal,
					Ext: tax.ExtensionsOf(cbc.CodeMap{
						iso.ExtKeySchemeID: identitySchemeIDSIREN,
					}),
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
				Code:    "44732829320",
			},
			Identities: []*org.Identity{
				{
					Type:  fr.IdentityTypeSIREN,
					Code:  "732829320",
					Scope: org.IdentityScopeLegal,
					Ext: tax.ExtensionsOf(cbc.CodeMap{
						iso.ExtKeySchemeID: identitySchemeIDSIREN,
					}),
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
						Name: "Test Supplier SARL",
					},
				},
			},
		},
		Notes: []*org.Note{
			{
				Key:  org.NoteKeyPayment,
				Text: "Une penalite fixe de 40 EUR sera appliquee.",
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyTextSubject: "PMT",
				}),
			},
			{
				Key:  org.NoteKeyPaymentMethod,
				Text: "Penalites de retard applicables.",
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyTextSubject: "PMD",
				}),
			},
			{
				Key:  org.NoteKeyPaymentTerm,
				Text: "Aucun escompte pour paiement anticipe.",
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					untdid.ExtKeyTextSubject: "AAB",
				}),
			},
		},
	}
}

// testInvoiceB2BFlow10 returns a Flow 10 cross-border B2B fixture
// (French supplier, German customer) so the Flow 10 dispatcher branch
// fires.
func testInvoiceB2BFlow10(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(V1),
		Code:      "INV-2026-001",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2026, 1, 15),
		Type:      bill.InvoiceTypeStandard,
		Supplier:  frPartyWithSIREN(),
		Customer:  deCustomerWithVATID(),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Product",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{Category: tax.CategoryVAT, Percent: num.NewPercentage(20, 2)},
				},
			},
		},
	}
}

// testInvoiceB2C returns a Flow 10 B2C fixture (no customer).
func testInvoiceB2C(t *testing.T) *bill.Invoice {
	t.Helper()
	return &bill.Invoice{
		Regime:    tax.WithRegime("FR"),
		Addons:    tax.WithAddons(V1),
		Code:      "INV-2026-B2C-001",
		Currency:  "EUR",
		IssueDate: cal.MakeDate(2026, 1, 15),
		Type:      bill.InvoiceTypeStandard,
		Tax: &bill.Tax{
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				ExtKeyB2CCategory: B2CCategoryGoods,
			}),
		},
		Supplier: frPartyWithSIREN(),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Product",
					Price: num.NewAmount(100, 0),
				},
				Taxes: tax.Set{
					{Category: tax.CategoryVAT, Percent: num.NewPercentage(20, 2)},
				},
			},
		},
	}
}

func setDocumentType(inv *bill.Invoice, docType string) {
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	if inv.Tax.Ext.IsZero() {
		inv.Tax.Ext = tax.MakeExtensions()
	}
	inv.Tax.Ext = inv.Tax.Ext.Set(untdid.ExtKeyDocumentType, cbc.Code(docType))
}

// =========================================================================
// Flow 10 invoice tests (ported from addons/fr/ctc/flow10/bill_invoice_test.go)
// =========================================================================

func TestInvoiceB2BFlow10HappyPath(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestInvoiceB2CHappyPath(t *testing.T) {
	inv := testInvoiceB2C(t)
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestInvoiceFlow10CurrencyRequiresEURConversion(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	inv.Currency = "USD"
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "EUR")
}

func TestInvoiceFlow10CurrencyUSDWithExchangeRate(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	inv.Currency = "USD"
	inv.ExchangeRates = []*currency.ExchangeRate{
		{From: "USD", To: "EUR", Amount: num.MakeAmount(875967, 6)},
	}
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
}

func TestInvoiceB2BFlow10DocTypeNotAllowed(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	require.NoError(t, inv.Calculate())
	inv.Tax.Ext = inv.Tax.Ext.Set(untdid.ExtKeyDocumentType, "325") // proforma, not allowed
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "Flow 10 permitted UNTDID 1001 codes")
}

func TestInvoiceB2BFlow10MissingBillingMode(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	require.NoError(t, inv.Calculate())
	inv.Tax.Ext = inv.Tax.Ext.Delete(ExtKeyBillingMode)
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "billing mode")
}

func TestInvoiceB2BFlow10FinalAfterAdvanceRejectsDepositDocType(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	require.NoError(t, inv.Calculate())
	inv.Tax.Ext = inv.Tax.Ext.
		Set(ExtKeyBillingMode, BillingModeM4).
		Set(untdid.ExtKeyDocumentType, "386") // Advance payment invoice
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "G1.60")
}

func TestInvoiceB2BFlow10SupplierRequiresAllowedScheme(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	require.NoError(t, inv.Calculate())
	inv.Supplier.Identities = nil
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "supplier must declare a legal identity")
}

func TestInvoiceB2BFlow10AddressRequiresCountry(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	require.NoError(t, inv.Calculate())
	inv.Supplier.Addresses = []*org.Address{{Street: "No country"}}
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "supplier address must include country")
}

func TestInvoiceB2BFlow10CustomerAddressRequiresCountry(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	require.NoError(t, inv.Calculate())
	inv.Customer.Addresses = []*org.Address{{Street: "No country"}}
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "customer address must include country")
}

func TestInvoiceB2BFlow10ExemptRequiresSellerVATID(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	inv.Lines[0].Taxes = tax.Set{
		{Category: tax.CategoryVAT, Key: tax.KeyExempt},
	}
	require.NoError(t, inv.Calculate())
	inv.Supplier.TaxID = nil
	inv.Ordering = nil
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "supplier VAT ID or ordering.seller")
}

func TestInvoiceB2BFlow10ExemptRequiresExemptTaxNote(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	inv.Lines[0].Taxes = tax.Set{
		{Category: tax.CategoryVAT, Key: tax.KeyExempt},
	}
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "exemption reason")
}

func TestInvoiceB2BFlow10ExemptHappyWithSellerVATAndNote(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	inv.Lines[0].Taxes = tax.Set{
		{Category: tax.CategoryVAT, Key: tax.KeyExempt},
	}
	require.NoError(t, inv.Calculate())
	// supplier already has a TaxID from frPartyWithSIREN. Add an exempt
	// tax note so rule 56 is satisfied.
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	inv.Tax.Notes = []*tax.Note{
		{Key: tax.KeyExempt, Text: "Exempt under VAT directive art. 138"},
	}
	require.NoError(t, rules.Validate(inv))
}

// TestInvoiceB2BFlow10ExemptOrderingSellerHasVATID covers the tax
// representative arrangement: a non-EU supplier (scheme 0227, no
// TaxID required by rule 54) sells exempt goods to a French customer
// and the French tax representative carries the VAT ID via
// ordering.seller. Rule 55 must accept the ordering.seller VAT ID in
// lieu of supplier.TaxID.
func TestInvoiceB2BFlow10ExemptOrderingSellerHasVATID(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	// Replace the French supplier with a non-EU one. Scheme 0227 doesn't
	// trigger rule 54's TaxID requirement, leaving rule 55 (exempt
	// reliance on ordering.seller VAT) as the meaningful check.
	inv.Supplier = &org.Party{
		Name: "Foreign Supplier Inc",
		// Foreign TaxID — the FR regime requires a tax_id code or a
		// SIREN/SIRET identity on the supplier, so we give the foreign
		// company its own (non-French) tax ID. The Flow 10 rule 54
		// requirement to carry a TaxID *when scheme is SIREN/EU-VAT*
		// does not apply here: the legal identity is non-EU (0227).
		TaxID: &tax.Identity{
			Country: "US",
			Code:    "12-3456789",
		},
		Identities: []*org.Identity{
			{
				Code:  "US-12345",
				Scope: org.IdentityScopeLegal,
				Ext: tax.ExtensionsOf(cbc.CodeMap{
					iso.ExtKeySchemeID: "0227",
				}),
			},
		},
		Addresses: []*org.Address{{Country: "US"}},
	}
	inv.Customer = frCustomerWithSIREN()
	inv.Lines[0].Taxes = tax.Set{
		{Category: tax.CategoryVAT, Key: tax.KeyExempt},
	}
	inv.Ordering = &bill.Ordering{
		Seller: &org.Party{
			Name: "Représentant Fiscal SARL",
			TaxID: &tax.Identity{
				Country: "FR",
				Code:    "39356000000",
			},
		},
	}
	require.NoError(t, inv.Calculate())
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	inv.Tax.Notes = []*tax.Note{
		{Key: tax.KeyExempt, Text: "Exempt — represented by French tax rep"},
	}
	require.NoError(t, rules.Validate(inv))
}

func TestInvoiceB2CDefaultsCategoryToTNT1(t *testing.T) {
	inv := testInvoiceB2C(t)
	inv.Tax.Ext = inv.Tax.Ext.Delete(ExtKeyB2CCategory)
	require.NoError(t, inv.Calculate())
	assert.Equal(t, B2CCategoryNotTaxable, inv.Tax.Ext.Get(ExtKeyB2CCategory))
	require.NoError(t, rules.Validate(inv))
}

func TestInvoiceB2CSupplierRequiresSIREN(t *testing.T) {
	inv := testInvoiceB2C(t)
	inv.Supplier.TaxID = nil
	inv.Supplier.Identities = nil
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "SIREN")
}

func TestInvoiceB2CVATRateNotInWhitelist(t *testing.T) {
	inv := testInvoiceB2C(t)
	inv.Lines[0].Taxes = tax.Set{
		{Category: tax.CategoryVAT, Percent: num.NewPercentage(17, 2)},
	}
	require.NoError(t, inv.Calculate())
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "G1.24")
}

func TestNormalizeFlow10DefaultBillingModeM1(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	require.NoError(t, inv.Calculate())
	assert.Equal(t, BillingModeM1, inv.Tax.Ext.Get(ExtKeyBillingMode))
}

func TestNormalizeFlow10TaxCategorySetFromKey(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	inv.Lines[0].Taxes = tax.Set{
		{Category: tax.CategoryVAT, Key: tax.KeyReverseCharge},
	}
	require.NoError(t, inv.Calculate())
	combo := inv.Lines[0].Taxes[0]
	assert.Equal(t, "AE", combo.Ext.Get(untdid.ExtKeyTaxCategory).String())
}

func TestNormalizeFlow10GeneratesSIRENFromFrenchTaxID(t *testing.T) {
	inv := testInvoiceB2BFlow10(t)
	inv.Supplier.Identities = nil
	require.NoError(t, inv.Calculate())
	found := false
	for _, id := range inv.Supplier.Identities {
		if id.Ext.Get(iso.ExtKeySchemeID).String() == "0002" {
			found = true
			assert.Equal(t, "356000000", id.Code.String())
		}
	}
	assert.True(t, found, "expected SIREN identity to be generated from TaxID")
}

func TestNormalizeB2CGeneratesSIRENFromFrenchTaxID(t *testing.T) {
	inv := testInvoiceB2C(t)
	inv.Supplier.Identities = nil
	require.NoError(t, inv.Calculate())
	require.NoError(t, rules.Validate(inv))
	found := false
	for _, id := range inv.Supplier.Identities {
		if id.Ext.Get(iso.ExtKeySchemeID).String() == "0002" {
			found = true
		}
	}
	assert.True(t, found, "expected SIREN identity to be generated from TaxID for B2C")
}

// --- Internal helper coverage (Flow 10) ---------------------------------

func TestPartyHasSIRENWrongType(t *testing.T) {
	assert.False(t, partyHasSIREN("x"))
}

func TestPartyHasAllowedLegalSchemeWrongType(t *testing.T) {
	assert.False(t, partyHasAllowedLegalScheme("x"))
}

func TestPartyHasTaxIDWhenRequiredWrongType(t *testing.T) {
	assert.True(t, partyHasTaxIDWhenRequired("x"))
}

func TestPartyHasTaxIDWhenRequiredNonRequiredScheme(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{{
		Code: "X",
		Ext:  tax.ExtensionsOf(cbc.CodeMap{iso.ExtKeySchemeID: "0227"}),
	}}}
	assert.True(t, partyHasTaxIDWhenRequired(p))
}

func TestInvoiceIsCrossBorderB2BWrongType(t *testing.T) {
	assert.False(t, invoiceIsCrossBorderB2BAny("x"))
}

func TestInvoiceIsB2CWrongType(t *testing.T) {
	assert.False(t, invoiceIsB2CAny("x"))
}

func TestInvoiceDocumentTypeAllowedEmpty(t *testing.T) {
	assert.False(t, invoiceDocumentTypeAllowed(tax.Extensions{}))
}

func TestExtensionsHaveBillingModeMissing(t *testing.T) {
	assert.False(t, extensionsHaveBillingMode(tax.Extensions{}))
}

func TestExtensionsHaveB2CCategoryMissing(t *testing.T) {
	assert.False(t, extensionsHaveB2CCategory(tax.Extensions{}))
}

func TestInvoiceIsFinalAfterAdvanceWrongType(t *testing.T) {
	assert.False(t, invoiceIsFinalAfterAdvance("x"))
}

func TestInvoiceIsFinalAfterAdvanceNoExt(t *testing.T) {
	assert.False(t, invoiceIsFinalAfterAdvance(&bill.Invoice{Tax: &bill.Tax{}}))
}

func TestInvoiceNotAdvancePaymentDocTypeWrongType(t *testing.T) {
	assert.True(t, invoiceNotAdvancePaymentDocType(42))
}

func TestInvoiceHasSellerVATIDForExemptWrongType(t *testing.T) {
	assert.False(t, invoiceHasSellerVATIDForExempt("x"))
}

func TestInvoiceHasExemptComboWrongType(t *testing.T) {
	assert.False(t, invoiceHasExemptCombo("x"))
}

func TestInvoiceHasExemptTaxNoteWrongType(t *testing.T) {
	assert.False(t, invoiceHasExemptTaxNote("x"))
}

func TestInvoiceVATRatesAllowedWrongType(t *testing.T) {
	assert.True(t, invoiceVATRatesAllowed("x"))
}

func TestMustParsePercentagesPanicsOnBadInput(t *testing.T) {
	assert.Panics(t, func() { mustParsePercentages("not-a-percentage") })
}

func TestPercentageInListEmpty(t *testing.T) {
	p := num.MakePercentage(20, 2)
	assert.False(t, percentageInList(p, nil))
}

func TestNormalizeInvoiceNilSafe(t *testing.T) {
	assert.NotPanics(t, func() { normalizeInvoice(nil) })
}

func TestNormalizeBillingModeDefaultsM2WhenPaid(t *testing.T) {
	due := num.MakeAmount(0, 2)
	inv := &bill.Invoice{
		Totals: &bill.Totals{Due: &due},
		Tax:    &bill.Tax{},
	}
	normalizeBillingMode(inv)
	assert.Equal(t, BillingModeM2, inv.Tax.Ext.Get(ExtKeyBillingMode))
}

// =========================================================================
// Flow 2 invoice tests (ported from addons/fr/ctc/flow2/bill_invoice_test.go)
// =========================================================================

func TestInvoiceValidation(t *testing.T) {
	t.Run("basic B2B invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("non-EUR currency without exchange rates", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Currency = "USD"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "EUR")
	})

	t.Run("non-EUR currency with exchange rates", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Currency = "USD"
		inv.ExchangeRates = []*currency.ExchangeRate{
			{From: "USD", To: "EUR", Amount: num.MakeAmount(875967, 6)},
		}
		require.NoError(t, inv.Calculate())
		assert.NoError(t, rules.Validate(inv))
	})

	t.Run("invoice code too long", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Code = "THIS-IS-A-VERY-LONG-INVOICE-CODE-THAT-EXCEEDS-35-CHARACTERS"
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "BR-FR-01/02")
	})

	t.Run("invoice code valid special chars", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Code = "INV-2024+001_TEST/A"
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("duplicate note codes not allowed", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Notes = append(inv.Notes, &org.Note{
			Key:  org.NoteKeyPayment,
			Text: "Duplicate payment terms",
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				untdid.ExtKeyTextSubject: "PMT",
			}),
		})
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "duplicate note codes")
		assert.ErrorContains(t, err, "BR-FR-06")
	})

	t.Run("supplier SIREN required (BR-FR-10)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Supplier.Identities = []*org.Identity{}
		inv.Supplier.TaxID = nil // prevent normalizer from regenerating SIREN
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "SIREN")
	})

	t.Run("B2B non-self-billed requires SIREN inbox (BR-FR-21)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Supplier.Inboxes = []*org.Inbox{
			{Scheme: "0088", Code: "1234567890123"},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "0225")
		assert.ErrorContains(t, err, "BR-FR-21")
	})

	t.Run("B2B non-self-billed SIREN inbox must start with SIREN (BR-FR-21)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Supplier.Inboxes = []*org.Inbox{
			{Scheme: cbc.Code("0225"), Code: "999999999"},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "0225")
	})

	t.Run("self-billed invoice does not require supplier SIREN inbox start match", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.SetTags(tax.TagSelfBilled)
		require.NoError(t, inv.Calculate())
		// docType is 389 now (self-billed). Remove SIREN inbox.
		inv.Supplier.Inboxes = []*org.Inbox{
			{Scheme: "0088", Code: "1234567890123"},
		}
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("B2B self-billed requires customer SIREN inbox", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.SetTags(tax.TagSelfBilled)
		require.NoError(t, inv.Calculate())
		inv.Customer.Inboxes = []*org.Inbox{
			{Scheme: "0088", Code: "1234567890123"},
		}
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "0225")
	})

	t.Run("B2B self-billed customer SIREN inbox must start with SIREN", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.SetTags(tax.TagSelfBilled)
		require.NoError(t, inv.Calculate())
		inv.Customer.Inboxes = []*org.Inbox{
			{Scheme: cbc.Code("0225"), Code: "999999999"},
		}
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "0225")
	})
}

func TestDocumentTypeValidation(t *testing.T) {
	t.Run("valid document type", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "380", inv.Tax.Ext.Get(untdid.ExtKeyDocumentType).String())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("invalid document type", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = inv.Tax.Ext.Set(untdid.ExtKeyDocumentType, "999")
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "BR-FR-04")
	})
}

func TestDocumentTypeScenarios(t *testing.T) {
	t.Run("standard invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "380", inv.Tax.Ext.Get(untdid.ExtKeyDocumentType).String())
	})

	t.Run("factoring invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.SetTags(tax.TagFactoring)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "393", inv.Tax.Ext.Get(untdid.ExtKeyDocumentType).String())
	})

	t.Run("advance payment invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.SetTags(tax.TagPrepayment)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "386", inv.Tax.Ext.Get(untdid.ExtKeyDocumentType).String())
	})

	t.Run("self-billed invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.SetTags(tax.TagSelfBilled)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "389", inv.Tax.Ext.Get(untdid.ExtKeyDocumentType).String())
	})

	t.Run("credit note", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "381", inv.Tax.Ext.Get(untdid.ExtKeyDocumentType).String())
	})

	t.Run("self-billed credit note", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.SetTags(tax.TagSelfBilled)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "261", inv.Tax.Ext.Get(untdid.ExtKeyDocumentType).String())
	})

	t.Run("corrective invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Type = bill.InvoiceTypeCorrective
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "384", inv.Tax.Ext.Get(untdid.ExtKeyDocumentType).String())
	})

	t.Run("factoring credit note", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Type = bill.InvoiceTypeCreditNote
		inv.SetTags(tax.TagFactoring)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, "396", inv.Tax.Ext.Get(untdid.ExtKeyDocumentType).String())
	})
}

func TestBillingModeNormalization(t *testing.T) {
	t.Run("user-specified billing mode preserved", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, BillingModeS5)
		require.NoError(t, inv.Calculate())
		assert.Equal(t, BillingModeS5.String(), inv.Tax.Ext.Get(ExtKeyBillingMode).String())
	})
}

func TestAttachmentValidation(t *testing.T) {
	t.Run("valid attachment description - LISIBLE", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Attachments = []*org.Attachment{
			{Code: "ATT001", Description: "LISIBLE", URL: "https://example.com/invoice.pdf"},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("valid attachment description - RIB", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Attachments = []*org.Attachment{
			{Code: "ATT001", Description: "RIB", URL: "https://example.com/rib.pdf"},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("invalid attachment description", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Attachments = []*org.Attachment{
			{Code: "ATT001", Description: "INVALID_TYPE", URL: "https://example.com/doc.pdf"},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "BR-FR-17")
	})

	t.Run("multiple LISIBLE attachments", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Attachments = []*org.Attachment{
			{Code: "ATT001", Description: "LISIBLE", URL: "https://example.com/invoice1.pdf"},
			{Code: "ATT002", Description: "LISIBLE", URL: "https://example.com/invoice2.pdf"},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "only one attachment with description 'LISIBLE'")
		assert.ErrorContains(t, err, "BR-FR-18")
	})

	t.Run("nil attachments handled gracefully", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		var att *org.Attachment
		inv.Attachments = []*org.Attachment{
			att,
			{Code: "ATT001", Description: "LISIBLE", URL: "https://example.com/invoice.pdf"},
			att,
			{Code: "ATT002", Description: "RIB", URL: "https://example.com/rib.pdf"},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

func TestOrderingIdentitiesValidation(t *testing.T) {
	t.Run("valid ordering with one AFL reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Ordering = &bill.Ordering{
			Identities: []*org.Identity{
				{Code: "12345", Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyReference: "AFL"})},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("valid ordering with one AWW reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Ordering = &bill.Ordering{
			Identities: []*org.Identity{
				{Code: "12345", Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyReference: "AWW"})},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("valid ordering with one AFL and one AWW", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Ordering = &bill.Ordering{
			Identities: []*org.Identity{
				{Code: "12345", Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyReference: "AFL"})},
				{Code: "67890", Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyReference: "AWW"})},
			},
		}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("invalid ordering with duplicate AFL reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Ordering = &bill.Ordering{
			Identities: []*org.Identity{
				{Code: "12345", Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyReference: "AFL"})},
				{Code: "67890", Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyReference: "AFL"})},
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "AFL")
		assert.ErrorContains(t, err, "BR-FR-30")
	})

	t.Run("invalid ordering with duplicate AWW reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Ordering = &bill.Ordering{
			Identities: []*org.Identity{
				{Code: "12345", Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyReference: "AWW"})},
				{Code: "67890", Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyReference: "AWW"})},
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "AWW")
		assert.ErrorContains(t, err, "BR-FR-30")
	})

	t.Run("ordering without identities is valid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Ordering = &bill.Ordering{Code: "ORD-12345"}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

func TestConsolidatedCreditNoteValidation(t *testing.T) {
	t.Run("valid consolidated credit note with delivery and contract", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Delivery = &bill.DeliveryDetails{
			Period: &cal.Period{
				Start: cal.MakeDate(2024, 5, 1),
				End:   cal.MakeDate(2024, 5, 31),
			},
		}
		inv.Ordering = &bill.Ordering{
			Contracts: []*org.DocumentRef{{Code: "CONTRACT-001"}},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "262")
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("consolidated credit note without delivery is invalid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Delivery = nil
		inv.Ordering = &bill.Ordering{
			Contracts: []*org.DocumentRef{{Code: "CONTRACT-001"}},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "262")
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "delivery details are required")
		assert.ErrorContains(t, err, "BR-FR-CO-03")
	})

	t.Run("consolidated credit note without delivery period", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Delivery = &bill.DeliveryDetails{}
		inv.Ordering = &bill.Ordering{
			Contracts: []*org.DocumentRef{{Code: "CONTRACT-001"}},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "262")
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "delivery period is required")
		assert.ErrorContains(t, err, "BR-FR-CO-03")
	})

	t.Run("consolidated credit note without ordering contracts", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Delivery = &bill.DeliveryDetails{
			Period: &cal.Period{
				Start: cal.MakeDate(2024, 5, 1),
				End:   cal.MakeDate(2024, 5, 31),
			},
		}
		inv.Ordering = &bill.Ordering{Contracts: nil}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "262")
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "ordering.contracts")
		assert.ErrorContains(t, err, "BR-FR-CO-03")
	})

	t.Run("consolidated credit note with nil ordering", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Delivery = &bill.DeliveryDetails{
			Period: &cal.Period{
				Start: cal.MakeDate(2024, 5, 1),
				End:   cal.MakeDate(2024, 5, 31),
			},
		}
		inv.Ordering = nil
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "262")
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "BR-FR-CO-03")
	})

	t.Run("non-consolidated credit note does not require delivery or contracts", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Delivery = nil
		inv.Ordering = nil
		inv.Preceding = []*org.DocumentRef{
			{Code: "INV-001", IssueDate: cal.NewDate(2024, 5, 1)},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "381")
		require.NoError(t, rules.Validate(inv))
	})
}

func TestSTCSupplierValidation(t *testing.T) {
	t.Run("STC supplier happy path", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Supplier.Identities = append(inv.Supplier.Identities, &org.Identity{
			Code: "12345678",
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				iso.ExtKeySchemeID: "0231",
			}),
		})
		inv.Ordering = &bill.Ordering{
			Seller: &org.Party{
				Name:  "Assujetti Unique",
				TaxID: inv.Supplier.TaxID,
			},
		}
		require.NoError(t, inv.Calculate())
		// normalizeSTCNote auto-appends the TXD / MEMBRE_ASSUJETTI_UNIQUE note.
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("STC supplier seller missing tax ID", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Supplier.Identities = append(inv.Supplier.Identities, &org.Identity{
			Code: "12345678",
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				iso.ExtKeySchemeID: "0231",
			}),
		})
		inv.Ordering = &bill.Ordering{
			Seller: &org.Party{Name: "Assujetti Unique", TaxID: nil},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "tax ID is required when supplier is under STC scheme")
	})

	t.Run("STC supplier seller with empty tax ID code", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Supplier.Identities = append(inv.Supplier.Identities, &org.Identity{
			Code: "12345678",
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				iso.ExtKeySchemeID: "0231",
			}),
		})
		inv.Ordering = &bill.Ordering{
			Seller: &org.Party{
				Name:  "Assujetti Unique",
				TaxID: &tax.Identity{Country: "FR", Code: ""},
			},
		}
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "code is required when supplier is under STC scheme")
	})

	t.Run("STC supplier with nil ordering", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Supplier.Identities = append(inv.Supplier.Identities, &org.Identity{
			Code: "12345678",
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				iso.ExtKeySchemeID: "0231",
			}),
		})
		inv.Ordering = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "BR-FR-CO-15")
	})

	t.Run("STC supplier requires TXD note", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Supplier.Identities = append(inv.Supplier.Identities, &org.Identity{
			Code: "12345678",
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				iso.ExtKeySchemeID: "0231",
			}),
		})
		inv.Ordering = &bill.Ordering{
			Seller: &org.Party{
				Name:  "Assujetti Unique",
				TaxID: inv.Supplier.TaxID,
			},
		}
		require.NoError(t, inv.Calculate())
		// Strip the auto-added TXD note.
		kept := inv.Notes[:0]
		for _, n := range inv.Notes {
			if n != nil && n.Ext.Get(untdid.ExtKeyTextSubject) == noteSubjectTXD {
				continue
			}
			kept = append(kept, n)
		}
		inv.Notes = kept
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, string(noteSubjectTXD))
		assert.ErrorContains(t, err, stcMembreAssujettiUnique)
	})

	t.Run("STC supplier normalizer auto-fills TXD note", func(t *testing.T) {
		// Exercise normalizeSTCNote directly to avoid the org_party rule 04
		// rejection of the 0231 identity scheme.
		inv := testInvoiceB2BStandard(t)
		inv.Supplier.Identities = append(inv.Supplier.Identities, &org.Identity{
			Code: "12345678",
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				iso.ExtKeySchemeID: "0231",
			}),
		})
		// Drop the auto-added TXD note in case the fixture already has one.
		inv.Notes = inv.Notes[:0]
		normalizeSTCNote(inv)
		var found bool
		for _, n := range inv.Notes {
			if n.Ext.Get(untdid.ExtKeyTextSubject) == noteSubjectTXD && n.Text == stcMembreAssujettiUnique {
				found = true
				break
			}
		}
		assert.True(t, found, "expected normalizer to add TXD note")
	})

	t.Run("normalizeSTCNote is idempotent when TXD already present", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Supplier.Identities = append(inv.Supplier.Identities, &org.Identity{
			Code: "12345678",
			Ext: tax.ExtensionsOf(cbc.CodeMap{
				iso.ExtKeySchemeID: "0231",
			}),
		})
		inv.Notes = []*org.Note{{
			Key:  org.NoteKeyLegal,
			Text: stcMembreAssujettiUnique,
			Ext:  tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: noteSubjectTXD}),
		}}
		normalizeSTCNote(inv)
		assert.Len(t, inv.Notes, 1)
	})

	t.Run("normalizeSTCNote no-op when supplier has no STC scheme", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		before := len(inv.Notes)
		normalizeSTCNote(inv)
		assert.Len(t, inv.Notes, before)
	})
}

func TestFinalInvoicePaymentValidation(t *testing.T) {
	t.Run("final invoice B2 with nil payment should fail", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, BillingModeB2)
		inv.Payment = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "payment")
	})

	t.Run("final invoice S2 with nil payment should fail", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, BillingModeS2)
		inv.Payment = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "payment")
	})

	t.Run("final invoice M2 with nil payment should fail", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, BillingModeM2)
		inv.Payment = nil
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "payment")
	})
}

func TestPrecedingReferencesValidation(t *testing.T) {
	t.Run("corrective invoice with exactly one preceding reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = []*org.DocumentRef{
			{Code: "INV-001", IssueDate: cal.NewDate(2024, 5, 1)},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "384")
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("corrective invoice with no preceding reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = nil
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "384")
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "must reference the original invoice in preceding")
		assert.ErrorContains(t, err, "BR-FR-CO-04")
	})

	t.Run("corrective invoice with multiple preceding references", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = []*org.DocumentRef{
			{Code: "INV-001", IssueDate: cal.NewDate(2024, 5, 1)},
			{Code: "INV-002", IssueDate: cal.NewDate(2024, 5, 2)},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "384")
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "must reference exactly one preceding invoice")
		assert.ErrorContains(t, err, "BR-FR-CO-04")
	})

	t.Run("corrective invoice type 471 requires one preceding reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = []*org.DocumentRef{
			{Code: "INV-001", IssueDate: cal.NewDate(2024, 5, 1)},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "471")
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("credit note with at least one preceding reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = []*org.DocumentRef{
			{Code: "INV-001", IssueDate: cal.NewDate(2024, 5, 1)},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "381")
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("credit note with multiple preceding references is valid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = []*org.DocumentRef{
			{Code: "INV-001", IssueDate: cal.NewDate(2024, 5, 1)},
			{Code: "INV-002", IssueDate: cal.NewDate(2024, 5, 2)},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "381")
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("credit note with no preceding reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = nil
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "381")
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "at least one preceding invoice reference")
		assert.ErrorContains(t, err, "BR-FR-CO-05")
	})

	t.Run("credit note type 261 with preceding", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = []*org.DocumentRef{
			{Code: "INV-001", IssueDate: cal.NewDate(2024, 5, 1)},
		}
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "261")
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("credit note type 502 with no preceding", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = nil
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "502")
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "BR-FR-CO-05")
	})

	t.Run("standard invoice does not require preceding reference", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Preceding = nil
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "380")
		require.NoError(t, rules.Validate(inv))
	})
}

func TestPaymentDueDateValidation(t *testing.T) {
	t.Run("valid due date on or after issue date", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.IssueDate = cal.MakeDate(2024, 6, 1)
		inv.Payment.Terms.DueDates[0].Date = cal.NewDate(2024, 7, 1)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("valid due date same as issue date", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.IssueDate = cal.MakeDate(2024, 6, 1)
		inv.Payment.Terms.DueDates[0].Date = cal.NewDate(2024, 6, 1)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("invalid due date before issue date", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.IssueDate = cal.MakeDate(2024, 6, 15)
		inv.Payment.Terms.DueDates[0].Date = cal.NewDate(2024, 6, 1)
		require.NoError(t, inv.Calculate())
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "due dates must not be before invoice issue date")
	})

	t.Run("advance payment type 386 allows due date before issue date", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.IssueDate = cal.MakeDate(2024, 6, 15)
		inv.Payment.Terms.DueDates[0].Date = cal.NewDate(2024, 6, 1)
		require.NoError(t, inv.Calculate())
		setDocumentType(inv, "386")
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("no due date is valid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.IssueDate = cal.MakeDate(2024, 6, 1)
		inv.Payment.Terms.DueDates = nil
		inv.Payment.Terms.Notes = "Payment on delivery"
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

func TestBillingModeDocumentTypeCompatibility(t *testing.T) {
	t.Run("factoring B4 with advance payment type 386 is invalid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, BillingModeB4)
		setDocumentType(inv, "386")
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "advance payment document types")
	})

	t.Run("factoring S4 with advance payment type 500 is invalid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, BillingModeS4)
		setDocumentType(inv, "500")
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "advance payment document types")
	})

	t.Run("factoring M4 with advance payment type 503 is invalid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, BillingModeM4)
		setDocumentType(inv, "503")
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "advance payment document types")
	})

	t.Run("factoring B4 with standard 380 is valid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, BillingModeB4)
		setDocumentType(inv, "380")
		require.NoError(t, rules.Validate(inv))
	})
}

func TestFinalInvoiceValidation(t *testing.T) {
	t.Run("valid final invoice B2 fully paid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, BillingModeB2)
		totalWithTax := inv.Totals.TotalWithTax
		inv.Totals.Advances = &totalWithTax
		zero := num.MakeAmount(0, 2)
		inv.Totals.Payable = zero
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("final invoice B2 without advance amount is invalid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, BillingModeB2)
		inv.Totals.Advances = nil
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "advance amount is required for already-paid invoices")
	})

	t.Run("final invoice B2 with incorrect advance amount is invalid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, BillingModeB2)
		wrongAmount := num.MakeAmount(5000, 2)
		inv.Totals.Advances = &wrongAmount
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "advance amount must equal total with tax")
	})

	t.Run("final invoice S2 with non-zero payable amount is invalid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, BillingModeS2)
		totalWithTax := inv.Totals.TotalWithTax
		inv.Totals.Advances = &totalWithTax
		nonZero := num.MakeAmount(100, 2)
		inv.Totals.Due = &nonZero
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "payable amount must be zero")
	})

	t.Run("final invoice M2 without due date", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		require.NoError(t, inv.Calculate())
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, BillingModeM2)
		totalWithTax := inv.Totals.TotalWithTax
		inv.Totals.Advances = &totalWithTax
		zero := num.MakeAmount(0, 2)
		inv.Totals.Payable = zero
		inv.Payment.Terms.DueDates = nil
		inv.Payment.Terms.Notes = "Payment already made"
		err := rules.Validate(inv)
		assert.ErrorContains(t, err, "at least one due date required")
	})

	t.Run("non-final invoice B7 does not require these validations", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, BillingModeB7)
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

func TestAdditionalDocumentTypes(t *testing.T) {
	t.Run("471 prepaid amount invoice", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(untdid.ExtKeyDocumentType, "471")
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("473 standalone credit note", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(untdid.ExtKeyDocumentType, "473")
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("502 self-billed corrective", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(untdid.ExtKeyDocumentType, "502")
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("503 self-billed credit for claim", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(untdid.ExtKeyDocumentType, "503")
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("472 self-billed prepaid", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(untdid.ExtKeyDocumentType, "472")
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})

	t.Run("261 self-billed credit note", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		inv.Tax.Ext = inv.Tax.Ext.Set(untdid.ExtKeyDocumentType, "261")
		inv.Preceding = []*org.DocumentRef{{Code: "INV-123"}}
		require.NoError(t, inv.Calculate())
		require.NoError(t, rules.Validate(inv))
	})
}

func TestInvoiceNormalization(t *testing.T) {
	ad := tax.AddonForKey(V1)

	t.Run("normalizes invoice with existing tax sets currency rounding (Flow 2)", func(t *testing.T) {
		inv := testInvoiceB2BStandard(t)
		ad.Normalizer(inv)
		assert.Equal(t, tax.RoundingRuleCurrency, inv.Tax.Rounding)
	})

	t.Run("normalizes nil invoice", func(t *testing.T) {
		var inv *bill.Invoice
		ad.Normalizer(inv)
		assert.Nil(t, inv)
	})
}

// --- Defensive nil / wrong-type branches ---------------------------------

func TestIsSelfBilledInvoiceNilInvoice(t *testing.T) {
	assert.False(t, isSelfBilledInvoice(nil))
}

func TestIsSelfBilledInvoiceMissingDocType(t *testing.T) {
	inv := &bill.Invoice{Tax: &bill.Tax{Ext: tax.ExtensionsOf(cbc.CodeMap{"other": "x"})}}
	assert.False(t, isSelfBilledInvoice(inv))
}

func TestIsCorrectiveInvoiceNilInvoice(t *testing.T) {
	assert.False(t, isCorrectiveInvoice(nil))
}

func TestGetPartySIRENNilParty(t *testing.T) {
	assert.Equal(t, "", getPartySIREN(nil))
}

func TestPrecedingDocCodeValidWrongType(t *testing.T) {
	assert.True(t, precedingDocCodeValid(42))
}

func TestIdentitiesHasLegalSIRENWrongType(t *testing.T) {
	assert.True(t, identitiesHasLegalSIREN(42))
}

func TestPartyHasSIRENInboxWrongType(t *testing.T) {
	assert.True(t, partyHasSIRENInbox(42))
}

func TestOrderingIdentitiesNoDupAFLWrongType(t *testing.T) {
	assert.True(t, orderingIdentitiesNoDupAFL(42))
}

func TestOrderingIdentitiesNoDupAWWWrongType(t *testing.T) {
	assert.True(t, orderingIdentitiesNoDupAWW(42))
}

func TestNotesHaveTXDWrongType(t *testing.T) {
	assert.False(t, notesHaveTXD(42))
}

func TestNotesHaveRequiredWrongType(t *testing.T) {
	assert.False(t, notesHaveRequired(42))
}

func TestNotesNoDuplicatesWrongType(t *testing.T) {
	assert.True(t, notesNoDuplicates(42))
}

func TestInvoiceDueDatesValidWrongType(t *testing.T) {
	assert.True(t, invoiceDueDatesValid(42))
}

func TestFinalInvoicePayableZeroWrongType(t *testing.T) {
	assert.True(t, finalInvoicePayableZero(42))
}

func TestFinalInvoiceAdvancesMatchWrongType(t *testing.T) {
	assert.True(t, finalInvoiceAdvancesMatch(42))
}

func TestAttachmentsUniqueLISIBLEEmpty(t *testing.T) {
	assert.True(t, attachmentsUniqueLISIBLE([]*org.Attachment{}))
}

func TestAttachmentsUniqueLISIBLEWrongType(t *testing.T) {
	assert.True(t, attachmentsUniqueLISIBLE(42))
}

func TestInvoiceIsDomesticFrenchNil(t *testing.T) {
	assert.False(t, invoiceIsDomesticFrench(nil))
}

func TestInvoiceIsDomesticFrenchAnyWrongType(t *testing.T) {
	assert.False(t, invoiceIsDomesticFrenchAny("x"))
}

func TestInvoiceIsNotDomesticFrenchAnyWrongType(t *testing.T) {
	assert.False(t, invoiceIsNotDomesticFrenchAny("x"))
}

func TestInvoiceMissingEN16931OnFlow2(t *testing.T) {
	inv := testInvoiceB2BStandard(t)
	require.NoError(t, inv.Calculate())
	// Drop the en16931 addon to simulate a caller forgetting to declare it.
	inv.Addons = tax.WithAddons(V1)
	err := rules.Validate(inv)
	assert.ErrorContains(t, err, "eu-en16931-v2017")
}

// --- defensive coverage: nil / wrong-type / empty-slice guards --------

func TestInvoiceCodeValidNonInvoice(t *testing.T) {
	assert.True(t, invoiceCodeValid(42))
}

func TestInvoiceCodeValidEmptyCode(t *testing.T) {
	assert.True(t, invoiceCodeValid(&bill.Invoice{}))
}

func TestPrecedingDocCodeValidNonDocumentRef(t *testing.T) {
	assert.True(t, precedingDocCodeValid(42))
}

func TestPrecedingDocCodeValidNil(t *testing.T) {
	assert.True(t, precedingDocCodeValid((*org.DocumentRef)(nil)))
}

func TestInvoiceIsFactoringAnyNonInvoice(t *testing.T) {
	assert.False(t, invoiceIsFactoringAny(42))
}

func TestInvoiceIsFactoringAnyEmptyTax(t *testing.T) {
	assert.False(t, invoiceIsFactoringAny(&bill.Invoice{}))
}

func TestIsCorrectiveInvoiceEmptyExt(t *testing.T) {
	assert.False(t, isCorrectiveInvoice(&bill.Invoice{Tax: &bill.Tax{}}))
}

func TestIsCreditNoteEmptyExt(t *testing.T) {
	assert.False(t, isCreditNote(&bill.Invoice{Tax: &bill.Tax{}}))
}

func TestIsConsolidatedCreditNoteEmptyExt(t *testing.T) {
	assert.False(t, isConsolidatedCreditNote(&bill.Invoice{Tax: &bill.Tax{}}))
}

func TestIsAdvancedInvoiceEmptyExt(t *testing.T) {
	assert.False(t, isAdvancedInvoice(&bill.Invoice{Tax: &bill.Tax{}}))
}

func TestIsFinalInvoiceEmptyExt(t *testing.T) {
	assert.False(t, isFinalInvoice(&bill.Invoice{Tax: &bill.Tax{}}))
}

func TestIsPartyIdentitySTCNilIdentity(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{nil}}
	assert.False(t, isPartyIdentitySTC(p))
}

func TestIsPartyIdentitySTCEmptyExt(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{{Code: "X"}}}
	assert.False(t, isPartyIdentitySTC(p))
}

func TestGetPartySIRENEmpty(t *testing.T) {
	p := &org.Party{Identities: []*org.Identity{{Code: "X"}}}
	assert.Equal(t, "", getPartySIREN(p))
}

func TestIdentitiesHasLegalSIRENNilEntry(t *testing.T) {
	assert.False(t, identitiesHasLegalSIREN([]*org.Identity{nil}))
}

func TestPartyHasSIRENInboxNoSIREN(t *testing.T) {
	p := &org.Party{Inboxes: []*org.Inbox{{Scheme: inboxSchemeSIREN, Code: "X"}}}
	assert.True(t, partyHasSIRENInbox(p))
}

func TestOrderingIdentitiesNoDupWrongType(t *testing.T) {
	assert.True(t, orderingIdentitiesNoDup("x", "AFL"))
}

func TestOrderingIdentitiesNoDupNilEntry(t *testing.T) {
	assert.True(t, orderingIdentitiesNoDup([]*org.Identity{nil}, "AFL"))
}

func TestNotesHaveRequiredEmpty(t *testing.T) {
	assert.False(t, notesHaveRequired([]*org.Note{}))
}

func TestNotesHaveRequiredNilEntry(t *testing.T) {
	assert.False(t, notesHaveRequired([]*org.Note{nil}))
}

func TestInvoiceHasNoteWithSubjectNilNote(t *testing.T) {
	inv := &bill.Invoice{Notes: []*org.Note{nil}}
	assert.False(t, invoiceHasNoteWithSubject(inv, "PMT"))
}

func TestNormalizeRequiredNotesNoOpWhenPresent(t *testing.T) {
	inv := &bill.Invoice{
		Notes: []*org.Note{
			{Key: org.NoteKeyPayment, Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "PMT"})},
			{Key: org.NoteKeyPaymentMethod, Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "PMD"})},
			{Key: org.NoteKeyPaymentTerm, Ext: tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "AAB"})},
		},
	}
	before := len(inv.Notes)
	normalizeRequiredNotes(inv)
	assert.Equal(t, before, len(inv.Notes))
}

func TestNormalizeB2CCategoryOnInvoicePreservesExisting(t *testing.T) {
	inv := &bill.Invoice{Tax: &bill.Tax{
		Ext: tax.ExtensionsOf(cbc.CodeMap{ExtKeyB2CCategory: B2CCategoryGoods}),
	}}
	normalizeB2CCategoryOnInvoice(inv)
	assert.Equal(t, B2CCategoryGoods, inv.Tax.Ext.Get(ExtKeyB2CCategory))
}

func TestNormalizeInvoiceTaxCategoriesNilLine(t *testing.T) {
	inv := &bill.Invoice{Lines: []*bill.Line{nil}}
	assert.NotPanics(t, func() { normalizeInvoiceTaxCategories(inv) })
}

func TestNormalizeInvoiceTaxCategoriesNilCombo(t *testing.T) {
	inv := &bill.Invoice{Lines: []*bill.Line{{Taxes: tax.Set{nil}}}}
	assert.NotPanics(t, func() { normalizeInvoiceTaxCategories(inv) })
}
