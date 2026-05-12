package ctc

import (
	"regexp"
	"slices"
	"strings"

	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// -- Constants & lookup tables --------------------------------------------

// invoiceCodeRegexp enforces BR-FR-01/02 invoice-code format: max 35
// characters, alphanumeric plus -+_/.
var invoiceCodeRegexp = regexp.MustCompile(`^[A-Za-z0-9\-\+_/]{1,35}$`)

// Self-billed document types (BR-FR-21/BR-FR-22).
var selfBilledDocumentTypes = []cbc.Code{
	"389", // Self-billed invoice
	"501", // Final invoice (self-billed context)
	"500", // Self-billed advance payment
	"471", // Prepaid amount invoice (self-billed context)
	"473", // Stand-alone credit note (self-billed context)
	"261", // Self-billed credit note
	"502", // Self-billed corrective
}

// Corrective invoice document types (BR-FR-CO-04).
var correctiveInvoiceTypes = []cbc.Code{
	"384", // Corrective invoice
	"471", // Prepaid amount invoice
	"472", // Self-billed prepaid amount
	"473", // Stand-alone credit note
}

// Credit note document types (BR-FR-CO-05).
var creditNoteTypes = []cbc.Code{
	"261", // Self-billed credit note
	"381", // Credit note
	"396", // Factoring credit note
	"502", // Self-billed corrective
	"503", // Self-billed credit for claim
}

// advancePaymentDocumentTypes are the UNTDID 1001 codes representing
// advance-payment invoices (forbidden combined with B4/S4/M4 modes).
var advancePaymentDocumentTypes = []cbc.Code{
	"386", // Advance payment invoice
	"500", // Self-billed advance payment
	"503", // Self-billed credit for claim
}

// finalAfterAdvanceBillingModes are the billing-mode codes that mark
// an invoice as a "final invoice after down payment".
var finalAfterAdvanceBillingModes = []cbc.Code{
	BillingModeB4, BillingModeS4, BillingModeM4,
}

// allowedAttachmentDescriptions enumerates the BR-FR-17 attachment
// descriptions accepted on a French CTC invoice.
var allowedAttachmentDescriptions = []string{
	"RIB",
	"LISIBLE",
	"FEUILLE_DE_STYLE",
	"PJA",
	"BON_LIVRAISON",
	"BON_COMMANDE",
	"DOCUMENT_ANNEXE",
	"BORDEREAU_SUIVI",
	"BORDEREAU_SUIVI_VALIDATION",
	"ETAT_ACOMPTE",
	"FACTURE_PAIEMENT_DIRECT",
	"RECAPITULATIF_COTRAITANCE",
}

// allowedVATRates is the whitelist of VAT percentages authorised on a
// Flow 10 invoice / payment (G1.24).
var allowedVATRates = mustParsePercentages(
	"0%", "0.9%", "1.05%", "1.75%", "2.1%", "5.5%", "7%", "8.5%",
	"9.2%", "9.6%", "10%", "13%", "19.6%", "20%", "20.6%",
)

func mustParsePercentages(values ...string) []num.Percentage {
	out := make([]num.Percentage, len(values))
	for i, v := range values {
		p, err := num.PercentageFromString(v)
		if err != nil {
			panic(err)
		}
		out[i] = p
	}
	return out
}

func percentageInList(p num.Percentage, list []num.Percentage) bool {
	for _, a := range list {
		if p.Compare(a) == 0 {
			return true
		}
	}
	return false
}

// vatKeyToUNTDIDCategory maps each supported GOBL VAT rate key to its
// UNTDID 5305 category code.
var vatKeyToUNTDIDCategory = map[cbc.Key]cbc.Code{
	tax.KeyStandard:       "S",
	tax.KeyZero:           "Z",
	tax.KeyExempt:         "E",
	tax.KeyReverseCharge:  "AE",
	tax.KeyIntraCommunity: "K",
	tax.KeyExport:         "G",
	tax.KeyOutsideScope:   "O",
}

const (
	// attachmentFormatLisible is the attachment format category for BR-FR-18.
	attachmentFormatLisible = "LISIBLE"

	// noteSubjectTXD is the UNTDID 4451 text-subject code carried on the
	// BR-FR-CO-14 STC (single-VAT-group) mention.
	noteSubjectTXD cbc.Code = "TXD"

	// stcMembreAssujettiUnique is the fixed text that pairs with TXD.
	stcMembreAssujettiUnique = "MEMBRE_ASSUJETTI_UNIQUE"
)

// -- Dispatcher -----------------------------------------------------------

// invoiceIsDomesticFrench reports whether both supplier and customer
// resolve as French (SIREN identity or French tax ID). When true, the
// Flow 2 clearance ruleset applies; otherwise Flow 10 reporting does.
func invoiceIsDomesticFrench(inv *bill.Invoice) bool {
	if inv == nil {
		return false
	}
	return partyIsFrench(inv.Supplier) && partyIsFrench(inv.Customer)
}

func invoiceIsDomesticFrenchAny(v any) bool {
	inv, ok := v.(*bill.Invoice)
	return ok && invoiceIsDomesticFrench(inv)
}

func invoiceIsNotDomesticFrenchAny(v any) bool {
	inv, ok := v.(*bill.Invoice)
	return ok && !invoiceIsDomesticFrench(inv)
}

// -- Normalisation --------------------------------------------------------

func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}

	// Map VAT keys to UNTDID 5305 category extensions on every line.
	// Safe to run unconditionally — it only sets the extension when the
	// key is one of the known buckets.
	normalizeInvoiceTaxCategories(inv)

	// Party-level normalisation applies to both supplier and customer
	// (SIREN derivation, peppol-key inbox marking, etc.).
	normalizeParty(inv.Supplier)
	normalizeParty(inv.Customer)

	// Branch on the dispatcher: Flow 2 (French B2B clearance) gets the
	// rounding, billing-mode, note and STC defaults; Flow 10 gets the
	// B2C category default (when applicable) and its own billing-mode
	// default for B2B reporting.
	if invoiceIsDomesticFrench(inv) {
		normalizeFlow2Invoice(inv)
		return
	}
	normalizeFlow10Invoice(inv)
}

func normalizeFlow2Invoice(inv *bill.Invoice) {
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	inv.Tax.Rounding = tax.RoundingRuleCurrency
	normalizeBillingMode(inv)
	normalizeRequiredNotes(inv)
	normalizeSTCNote(inv)
}

func normalizeFlow10Invoice(inv *bill.Invoice) {
	if invoiceIsB2C(inv) {
		normalizeB2CCategoryOnInvoice(inv)
		return
	}
	normalizeBillingMode(inv)
}

// invoiceIsB2C reports whether the invoice is a business-to-consumer
// transaction. Flow 10 distinguishes B2C from B2B by the presence of a
// Customer party.
func invoiceIsB2C(inv *bill.Invoice) bool {
	return inv != nil && inv.Customer == nil
}

// normalizeBillingMode picks a sensible default for the billing-mode
// extension when the caller hasn't supplied one. M2 when the invoice
// is fully paid, M1 otherwise.
func normalizeBillingMode(inv *bill.Invoice) {
	if inv.Tax != nil && !inv.Tax.Ext.IsZero() && inv.Tax.Ext.Get(ExtKeyBillingMode) != "" {
		return
	}
	mode := BillingModeM1
	if inv.Totals != nil && inv.Totals.Paid() {
		mode = BillingModeM2
	}
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyBillingMode, mode)
}

// normalizeSTCNote appends the BR-FR-CO-14 TXD / MEMBRE_ASSUJETTI_UNIQUE
// note when the supplier carries an STC-scheme (0231) identity and no
// such note has been provided yet.
func normalizeSTCNote(inv *bill.Invoice) {
	if !isPartyIdentitySTC(inv.Supplier) {
		return
	}
	for _, n := range inv.Notes {
		if n != nil && n.Ext.Get(untdid.ExtKeyTextSubject) == noteSubjectTXD && n.Text == stcMembreAssujettiUnique {
			return
		}
	}
	inv.Notes = append(inv.Notes, &org.Note{
		Key:  org.NoteKeyLegal,
		Text: stcMembreAssujettiUnique,
		Ext:  tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: noteSubjectTXD}),
	})
}

// defaultRequiredNotes lists the three UNTDID 4451 mentions French CTC
// requires on every B2B invoice (BR-FR-05).
var defaultRequiredNotes = []*org.Note{
	{
		Key:  org.NoteKeyPayment,
		Text: "Conditions de paiement selon les conditions générales de vente.",
		Ext:  tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "PMT"}),
	},
	{
		Key:  org.NoteKeyPaymentMethod,
		Text: "Pénalités et indemnités de retard applicables conformément aux conditions générales de vente.",
		Ext:  tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "PMD"}),
	},
	{
		Key:  org.NoteKeyPaymentTerm,
		Text: "Aucun escompte n'est accordé pour paiement anticipé.",
		Ext:  tax.ExtensionsOf(cbc.CodeMap{untdid.ExtKeyTextSubject: "AAB"}),
	},
}

// normalizeRequiredNotes appends any of the three regulatory PMT / PMD
// / AAB notes that are missing from the invoice.
func normalizeRequiredNotes(inv *bill.Invoice) {
	for _, def := range defaultRequiredNotes {
		want := def.Ext.Get(untdid.ExtKeyTextSubject)
		if invoiceHasNoteWithSubject(inv, want) {
			continue
		}
		clone := *def
		inv.Notes = append(inv.Notes, &clone)
	}
}

func invoiceHasNoteWithSubject(inv *bill.Invoice, subject cbc.Code) bool {
	for _, n := range inv.Notes {
		if n != nil && n.Ext.Get(untdid.ExtKeyTextSubject) == subject {
			return true
		}
	}
	return false
}

// normalizeB2CCategoryOnInvoice defaults the B2C transaction category
// to TNT1 (not subject to French VAT) when the caller has not supplied
// one.
func normalizeB2CCategoryOnInvoice(inv *bill.Invoice) {
	if inv.Tax != nil && inv.Tax.Ext.Get(ExtKeyB2CCategory) != "" {
		return
	}
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}
	inv.Tax.Ext = inv.Tax.Ext.Set(ExtKeyB2CCategory, B2CCategoryNotTaxable)
}

// normalizeInvoiceTaxCategories sets the UNTDID 5305 category extension
// on each VAT combo based on its rate key.
func normalizeInvoiceTaxCategories(inv *bill.Invoice) {
	for _, line := range inv.Lines {
		if line == nil {
			continue
		}
		for _, combo := range line.Taxes {
			if combo == nil || combo.Category != tax.CategoryVAT {
				continue
			}
			if combo.Ext.Get(untdid.ExtKeyTaxCategory) != "" {
				continue
			}
			if code, ok := vatKeyToUNTDIDCategory[combo.Key]; ok {
				combo.Ext = combo.Ext.Set(untdid.ExtKeyTaxCategory, code)
			}
		}
	}
}

// -- Rule set -------------------------------------------------------------

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		// Shared rules — apply to every French CTC invoice regardless of
		// the flow. EUR-convertibility is required by both flows.
		rules.Assert("01", "invoice must be in EUR or provide an exchange rate to EUR",
			currency.CanConvertTo(currency.EUR),
		),

		// Flow 2 — domestic B2B clearance ----------------------------------
		rules.When(
			is.Func("domestic French B2B (Flow 2 clearance)", invoiceIsDomesticFrenchAny),
			flow2InvoiceDefs()...,
		),

		// Flow 10 — e-reporting (non-domestic or B2C) ----------------------
		rules.When(
			is.Func("Flow 10 reporting (cross-border or B2C)", invoiceIsNotDomesticFrenchAny),
			flow10InvoiceDefs()...,
		),
	)
}

// flow2InvoiceDefs returns the rule defs applied to a domestic French
// B2B (Flow 2 clearance) invoice. Pulled out so the rule.When call in
// billInvoiceRules reads cleanly.
func flow2InvoiceDefs() []rules.Def {
	return []rules.Def{
		// EN16931 base profile is mandatory for domestic French B2B
		// (Flow 2 clearance). It is intentionally not a hard Requires
		// on the addon so that pure Flow 10 / Flow 6 callers don't
		// have to drag it in.
		rules.Assert("02", "domestic French B2B invoices must also declare the eu-en16931-v2017 addon",
			tax.HasAddon(en16931.V2017),
		),
		// Invoice code validation (BR-FR-01/02).
		rules.Assert("03", "must be 1-35 characters, alphanumeric plus -+_/ (BR-FR-01/02), including the series",
			is.Func("valid invoice code", invoiceCodeValid),
		),
		// Preceding document code validation.
		rules.Field("preceding",
			rules.Each(
				rules.Assert("04", "preceding code must be 1-35 characters, alphanumeric plus -+_/ (BR-FR-01/02), including the series",
					is.Func("valid preceding code", precedingDocCodeValid),
				),
			),
		),
		rules.When(
			is.Func("corrective invoice", invoiceIsCorrectiveAny),
			rules.Field("preceding",
				rules.Assert("05", "corrective invoices must reference the original invoice in preceding (BR-FR-CO-04)",
					is.Present,
				),
				rules.Assert("06", "corrective invoices must reference exactly one preceding invoice — multiple references are not allowed (BR-FR-CO-04)",
					is.Length(1, 1),
				),
			),
		),
		rules.When(
			is.Func("credit note", invoiceIsCreditNoteAny),
			rules.Field("preceding",
				rules.Assert("07", "credit notes must have at least one preceding invoice reference (BR-FR-CO-05)",
					is.Present,
				),
			),
		),
		rules.Field("tax",
			rules.Assert("08", "tax is required", is.Present),
			rules.Field("ext",
				rules.Assert("09", "UNTDID document type must be valid (BR-FR-04)",
					tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, allowedInvoiceDocumentTypes...),
				),
				rules.Assert("10", "billing mode extension is required",
					tax.ExtensionsRequire(ExtKeyBillingMode),
				),
			),
		),
		rules.When(
			is.Func("factoring mode", invoiceIsFactoringAny),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("11", "advance payment document types (386, 500, 503) are not allowed for factoring billing modes (B4, S4, M4) (BR-FR-CO-08)",
						tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, advancePaymentDocumentTypes...),
					),
				),
			),
		),
		rules.Field("supplier",
			rules.Field("inboxes",
				rules.Assert("12", "seller electronic address required for French B2B invoices (BR-FR-13)",
					is.Present,
				),
			),
			rules.Field("identities",
				rules.Assert("13", "SIREN identity required for French parties with scheme 0002 and scope legal (BR-FR-10/11)",
					is.Func("has SIREN (legal scope)", identitiesHasLegalSIREN),
				),
			),
		),
		rules.When(
			is.Func("not self-billed", invoiceIsNotSelfBilledAny),
			rules.Field("supplier",
				rules.Assert("14", "party must have endpoint ID with scheme 0225 (SIREN) (BR-FR-21/22)",
					is.Func("has SIREN inbox", partyHasSIRENInbox),
				),
			),
		),
		rules.Field("customer",
			rules.Field("inboxes",
				rules.Assert("15", "buyer electronic address required for French B2B invoices (BR-FR-13)",
					is.Present,
				),
			),
			rules.Field("identities",
				rules.Assert("16", "SIREN identity required for French parties with scheme 0002 and scope legal (BR-FR-10/11)",
					is.Func("has SIREN (legal scope)", identitiesHasLegalSIREN),
				),
			),
		),
		rules.When(
			is.Func("self-billed", invoiceIsSelfBilledAny),
			rules.Field("customer",
				rules.Assert("17", "party must have endpoint ID with scheme 0225 (SIREN) (BR-FR-21/22)",
					is.Func("has SIREN inbox", partyHasSIRENInbox),
				),
			),
		),
		rules.Field("ordering",
			rules.Field("identities",
				rules.Assert("18", "only one ordering identity with UNTDID reference 'AFL' is allowed (BR-FR-30)",
					is.Func("no dup AFL", orderingIdentitiesNoDupAFL),
				),
				rules.Assert("19", "only one ordering identity with UNTDID reference 'AWW' is allowed (BR-FR-30)",
					is.Func("no dup AWW", orderingIdentitiesNoDupAWW),
				),
			),
		),
		rules.When(
			is.Func("supplier STC", invoiceSupplierIsSTC),
			rules.Field("ordering",
				rules.Assert("20", "ordering with seller is required when supplier is under STC scheme (BR-FR-CO-15)",
					is.Present,
				),
				rules.Field("seller",
					rules.Assert("21", "seller is required when supplier is under STC scheme (BR-FR-CO-15)",
						is.Present,
					),
					rules.Field("tax_id",
						rules.Assert("22", "tax ID is required when supplier is under STC scheme (BR-FR-CO-15)",
							is.Present,
						),
						rules.Field("code",
							rules.Assert("23", "code is required when supplier is under STC scheme (BR-FR-CO-15)",
								is.Present,
							),
						),
					),
				),
			),
			rules.Field("notes",
				rules.Assert("24", "for sellers with STC scheme (0231), a note with code 'TXD' and text 'MEMBRE_ASSUJETTI_UNIQUE' is required (BR-FR-CO-14)",
					is.Func("has TXD note", notesHaveTXD),
				),
			),
		),
		rules.When(
			is.Func("consolidated credit note", invoiceIsConsolidatedCreditNoteAny),
			rules.Field("ordering",
				rules.Assert("25", "ordering with contracts is required for consolidated credit notes (BR-FR-CO-03)",
					is.Present,
				),
				rules.Field("contracts",
					rules.Assert("26", "ordering.contracts is required for consolidated credit notes (BR-FR-CO-03)",
						is.Present,
					),
					rules.Assert("27", "ordering.contracts must contain at least one entry for consolidated credit notes (BR-FR-CO-03)",
						is.Length(1, 0),
					),
				),
			),
			rules.Field("delivery",
				rules.Assert("28", "delivery details are required for consolidated credit notes (BR-FR-CO-03)",
					is.Present,
				),
				rules.Field("period",
					rules.Assert("29", "delivery period is required for consolidated credit notes (BR-FR-CO-03)",
						is.Present,
					),
				),
			),
		),
		rules.When(
			is.Func("not advance or final", invoiceIsNotAdvanceOrFinalAny),
			rules.Assert("30", "due dates must not be before invoice issue date (BR-FR-CO-07)",
				is.Func("due dates valid", invoiceDueDatesValid),
			),
		),
		rules.When(
			is.Func("final invoice", invoiceIsFinalAny),
			rules.Field("payment",
				rules.Assert("31", "payment details are required for final invoices (BR-FR-CO-09)",
					is.Present,
				),
				rules.Field("terms",
					rules.Assert("32", "payment terms required for final invoices (BR-FR-CO-09)",
						is.Present,
					),
					rules.Field("due_dates",
						rules.Assert("33", "at least one due date required for final invoices (BR-FR-CO-09)",
							is.Present,
						),
					),
				),
			),
			rules.Field("totals",
				rules.Field("advance",
					rules.Assert("34", "advance amount is required for already-paid invoices (BR-FR-CO-09)",
						is.Present,
					),
				),
				rules.Assert("35", "advance amount must equal total with tax for final invoices (BR-FR-CO-09)",
					is.Func("advances match", finalInvoiceAdvancesMatch),
				),
				rules.Assert("36", "payable amount must be zero for final invoices (BR-FR-CO-09)",
					is.Func("payable zero", finalInvoicePayableZero),
				),
			),
		),
		rules.Field("notes",
			rules.Assert("37", "notes are required for French CTC invoices (BR-FR-05)", is.Present),
			rules.Assert("38", "missing required note codes (BR-FR-05)",
				is.Func("has required notes", notesHaveRequired),
			),
			rules.Assert("39", "duplicate note codes found (BR-FR-06/BR-FR-30)",
				is.Func("no duplicate notes", notesNoDuplicates),
			),
		),
		rules.Field("attachments",
			rules.Each(
				rules.Field("description",
					rules.Assert("40", "must be one of the allowed attachment descriptions (BR-FR-17)",
						is.Present,
					),
					rules.Assert("41", "must be one of the allowed attachment descriptions (BR-FR-17)",
						is.In(toAnySlice(allowedAttachmentDescriptions)...),
					),
				),
			),
			rules.Assert("42", "only one attachment with description 'LISIBLE' is allowed per invoice (BR-FR-18)",
				is.Func("unique LISIBLE", attachmentsUniqueLISIBLE),
			),
		),
	}
}

// flow10InvoiceDefs returns the rule defs applied to a non-domestic
// (Flow 10 e-reporting) invoice — B2C or cross-border B2B.
func flow10InvoiceDefs() []rules.Def {
	return []rules.Def{
		// B2C rules: category, supplier SIREN, VAT rate whitelist.
		rules.When(
			is.Func("B2C invoice", invoiceIsB2CAny),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("43", "B2C transaction category extension (fr-ctc-b2c-category) is required on B2C invoices (G1.68)",
						is.Func("has B2C category", extensionsHaveB2CCategory),
					),
				),
			),
			rules.Field("supplier",
				rules.Assert("44", "supplier is required on B2C invoice",
					is.Present,
				),
				rules.Assert("45", "supplier must have a SIREN identity (ISO/IEC 6523 scheme 0002) on a B2C invoice",
					is.Func("party has SIREN", partyHasSIREN),
				),
			),
			rules.Assert("46", "every VAT line rate must be one of the Flow 10 permitted percentages (G1.24): 0, 0.9, 1.05, 1.75, 2.1, 5.5, 7, 8.5, 9.2, 9.6, 10, 13, 19.6, 20, 20.6",
				is.Func("allowed Flow 10 VAT rates", invoiceVATRatesAllowed),
			),
		),
		rules.Field("supplier",
			rules.Field("addresses",
				rules.Each(
					rules.Field("country",
						rules.Assert("47", "supplier address must include country",
							is.Present,
						),
					),
				),
			),
		),
		rules.Field("customer",
			rules.Field("addresses",
				rules.Each(
					rules.Field("country",
						rules.Assert("48", "customer address must include country",
							is.Present,
						),
					),
				),
			),
		),
		// B2B reporting: supplier + customer must each carry a legal
		// identity declaring an allowed ICD 6523 scheme, plus matching
		// TaxID when scheme is SIREN / EU-VAT.
		rules.When(
			is.Func("cross-border B2B invoice", invoiceIsCrossBorderB2BAny),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("49", "invoice document type must be one of the Flow 10 permitted UNTDID 1001 codes",
						is.Func("allowed Flow 10 document type", invoiceDocumentTypeAllowed),
					),
					rules.Assert("50", "billing mode extension (fr-ctc-billing-mode) is required (G1.02)",
						is.Func("has billing mode", extensionsHaveBillingMode),
					),
				),
			),
			rules.When(
				is.Func("billing mode is final-after-advance (B4/S4/M4)", invoiceIsFinalAfterAdvance),
				rules.Field("tax",
					rules.Field("ext",
						rules.Assert("51", "final-after-advance billing mode (B4/S4/M4) cannot be combined with an advance-payment document type (386/500/503) (G1.60)",
							is.Func("not advance-payment doc type", invoiceNotAdvancePaymentDocType),
						),
					),
				),
			),
			rules.Field("supplier",
				rules.Assert("52", "supplier is required for Flow 10 B2B invoice (G2.19)",
					is.Present,
				),
				rules.Assert("53", "supplier must declare a legal identity with an allowed ICD 6523 scheme (G2.19): 0002, 0223, 0227, 0228 or 0229",
					is.Func("party has allowed legal scheme", partyHasAllowedLegalScheme),
				),
				rules.Assert("54", "supplier TaxID is required when legal identity scheme is SIREN (0002) or EU VAT (0223) (G2.33)",
					is.Func("party has TaxID when required", partyHasTaxIDWhenRequired),
				),
			),
			rules.When(
				is.Func("invoice has exempt (E) VAT category", invoiceHasExemptCombo),
				rules.Assert("55", "supplier VAT ID or ordering.seller (tax representative) VAT ID is required when the invoice VAT breakdown contains an exempt (E) category",
					is.Func("supplier or tax rep has VAT ID", invoiceHasSellerVATIDForExempt),
				),
				rules.Assert("56", "invoice with an exempt (E) VAT category must include an exemption reason in tax.notes (key=exempt, non-empty text)",
					is.Func("has exempt tax note", invoiceHasExemptTaxNote),
				),
			),
			rules.Field("customer",
				rules.Assert("57", "customer is required for Flow 10 B2B invoice (G2.19)",
					is.Present,
				),
				rules.Assert("58", "customer must declare a legal identity with an allowed ICD 6523 scheme (G2.19): 0002, 0223, 0227, 0228 or 0229",
					is.Func("party has allowed legal scheme", partyHasAllowedLegalScheme),
				),
				rules.Assert("59", "customer TaxID is required when legal identity scheme is SIREN (0002) or EU VAT (0223) (G2.33)",
					is.Func("party has TaxID when required", partyHasTaxIDWhenRequired),
				),
			),
		),
	}
}

// -- Flow 2 helpers -------------------------------------------------------

func invoiceCodeValid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Code == cbc.CodeEmpty {
		return true
	}
	invoiceID := string(inv.Code)
	if inv.Series != cbc.CodeEmpty {
		invoiceID = string(inv.Series.Join(inv.Code))
	}
	return invoiceCodeRegexp.MatchString(invoiceID)
}

func precedingDocCodeValid(val any) bool {
	docRef, ok := val.(*org.DocumentRef)
	if !ok || docRef == nil || docRef.Code == cbc.CodeEmpty {
		return true
	}
	invoiceID := string(docRef.Code)
	if docRef.Series != cbc.CodeEmpty {
		invoiceID = string(docRef.Series.Join(docRef.Code))
	}
	return invoiceCodeRegexp.MatchString(invoiceID)
}

func invoiceIsCorrectiveAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && isCorrectiveInvoice(inv)
}

func invoiceIsCreditNoteAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && isCreditNote(inv)
}

func invoiceIsFactoringAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	return isFactoringExtension(inv.Tax.Ext.Get(ExtKeyBillingMode))
}

// Within the Flow 2 dispatcher branch the invoice is B2B by
// construction (both parties French), so self-billed vs not is the
// only remaining axis.

func invoiceIsNotSelfBilledAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && !isSelfBilledInvoice(inv)
}

func invoiceIsSelfBilledAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && isSelfBilledInvoice(inv)
}

func invoiceSupplierIsSTC(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && isPartyIdentitySTC(inv.Supplier)
}

func invoiceIsConsolidatedCreditNoteAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && isConsolidatedCreditNote(inv)
}

func invoiceIsNotAdvanceOrFinalAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && !isAdvancedInvoice(inv) && !isFinalInvoice(inv)
}

func invoiceIsFinalAny(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && isFinalInvoice(inv)
}

func isSelfBilledInvoice(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	if docType == "" {
		return false
	}
	return slices.Contains(selfBilledDocumentTypes, docType)
}

func isCorrectiveInvoice(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	if docType == "" {
		return false
	}
	return slices.Contains(correctiveInvoiceTypes, docType)
}

func isPartyIdentitySTC(party *org.Party) bool {
	if party == nil || len(party.Identities) == 0 {
		return false
	}
	for _, id := range party.Identities {
		if id != nil && !id.Ext.IsZero() {
			if code := id.Ext.Get(iso.ExtKeySchemeID); code == "0231" {
				return true
			}
		}
	}
	return false
}

func isCreditNote(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	return slices.Contains(creditNoteTypes, docType)
}

func isConsolidatedCreditNote(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	return docType == "262"
}

func isAdvancedInvoice(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	docType := inv.Tax.Ext.Get(untdid.ExtKeyDocumentType)
	return slices.Contains(advancePaymentDocumentTypes, docType)
}

// isFinalInvoice checks if the invoice is a final invoice based on
// billing mode (B2, S2, M2).
func isFinalInvoice(inv *bill.Invoice) bool {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	bm := inv.Tax.Ext.Get(ExtKeyBillingMode)
	return bm == BillingModeB2 || bm == BillingModeS2 || bm == BillingModeM2
}

func isFactoringExtension(bm cbc.Code) bool {
	return bm == BillingModeB4 || bm == BillingModeS4 || bm == BillingModeM4
}

// getPartySIREN extracts the SIREN string from a party's SIREN identity.
func getPartySIREN(party *org.Party) string {
	if party == nil {
		return ""
	}
	for _, id := range party.Identities {
		if id != nil && (id.Type == fr.IdentityTypeSIREN || (!id.Ext.IsZero() && id.Ext.Get(iso.ExtKeySchemeID) == identitySchemeIDSIREN)) {
			return string(id.Code)
		}
	}
	return ""
}

func identitiesHasLegalSIREN(val any) bool {
	identities, ok := val.([]*org.Identity)
	if !ok {
		return true
	}
	for _, id := range identities {
		if id != nil && !id.Ext.IsZero() {
			if code := id.Ext.Get(iso.ExtKeySchemeID); code == identitySchemeIDSIREN && id.Scope.Has(org.IdentityScopeLegal) {
				return true
			}
		}
	}
	return false
}

func partyHasSIRENInbox(val any) bool {
	party, ok := val.(*org.Party)
	if !ok || party == nil {
		return true
	}
	siren := getPartySIREN(party)
	if siren == "" {
		return true
	}
	for _, inbox := range party.Inboxes {
		if inbox != nil && inbox.Scheme == inboxSchemeSIREN {
			return strings.HasPrefix(string(inbox.Code), siren)
		}
	}
	return false
}

func orderingIdentitiesNoDupAFL(val any) bool {
	return orderingIdentitiesNoDup(val, "AFL")
}

func orderingIdentitiesNoDupAWW(val any) bool {
	return orderingIdentitiesNoDup(val, "AWW")
}

func orderingIdentitiesNoDup(val any, ref string) bool {
	identities, ok := val.([]*org.Identity)
	if !ok {
		return true
	}
	count := 0
	for _, id := range identities {
		if id == nil || id.Ext.IsZero() {
			continue
		}
		if id.Ext.Get(untdid.ExtKeyReference).String() == ref {
			count++
			if count > 1 {
				return false
			}
		}
	}
	return true
}

func notesHaveTXD(val any) bool {
	notes, ok := val.([]*org.Note)
	if !ok || len(notes) == 0 {
		return false
	}
	for _, note := range notes {
		if note != nil && !note.Ext.IsZero() {
			if code := note.Ext.Get(untdid.ExtKeyTextSubject); code == noteSubjectTXD && note.Text == stcMembreAssujettiUnique {
				return true
			}
		}
	}
	return false
}

func notesHaveRequired(val any) bool {
	notes, ok := val.([]*org.Note)
	if !ok || len(notes) == 0 {
		return false
	}
	required := []cbc.Code{"PMT", "PMD", "AAB"}
	counts := make(map[cbc.Code]int)
	for _, note := range notes {
		if note != nil && !note.Ext.IsZero() {
			if code := note.Ext.Get(untdid.ExtKeyTextSubject); code != cbc.CodeEmpty {
				counts[code]++
			}
		}
	}
	for _, code := range required {
		if counts[code] == 0 {
			return false
		}
	}
	return true
}

func notesNoDuplicates(val any) bool {
	notes, ok := val.([]*org.Note)
	if !ok || len(notes) == 0 {
		return true
	}
	counts := make(map[cbc.Code]int)
	for _, note := range notes {
		if note != nil && !note.Ext.IsZero() {
			if code := note.Ext.Get(untdid.ExtKeyTextSubject); code != cbc.CodeEmpty {
				counts[code]++
			}
		}
	}
	checkUnique := []cbc.Code{"PMT", "PMD", "AAB", "TXD"}
	for _, code := range checkUnique {
		if counts[code] > 1 {
			return false
		}
	}
	return true
}

func invoiceDueDatesValid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	if inv.Payment == nil || inv.Payment.Terms == nil || len(inv.Payment.Terms.DueDates) == 0 {
		return true
	}
	for _, dd := range inv.Payment.Terms.DueDates {
		if dd == nil || dd.Date == nil {
			continue
		}
		if inv.IssueDate.DaysSince(dd.Date.Date) > 0 {
			return false
		}
	}
	return true
}

func finalInvoiceAdvancesMatch(val any) bool {
	totals, ok := val.(*bill.Totals)
	if !ok || totals == nil || totals.Advances == nil {
		return true
	}
	return totals.Advances.Equals(totals.TotalWithTax)
}

func finalInvoicePayableZero(val any) bool {
	totals, ok := val.(*bill.Totals)
	if !ok || totals == nil {
		return true
	}
	if totals.Due != nil {
		return totals.Due.Equals(num.AmountZero)
	}
	return totals.Payable.Equals(num.AmountZero)
}

func attachmentsUniqueLISIBLE(val any) bool {
	attachments, ok := val.([]*org.Attachment)
	if !ok || len(attachments) == 0 {
		return true
	}
	count := 0
	for _, att := range attachments {
		if att != nil && att.Description == attachmentFormatLisible {
			count++
		}
	}
	return count <= 1
}

func toAnySlice(ss []string) []any {
	out := make([]any, len(ss))
	for i, s := range ss {
		out[i] = s
	}
	return out
}

// -- Flow 10 helpers ------------------------------------------------------

func invoiceIsB2CAny(v any) bool {
	inv, ok := v.(*bill.Invoice)
	return ok && invoiceIsB2C(inv)
}

// invoiceIsCrossBorderB2BAny reports whether a Flow 10 invoice is a
// cross-border B2B transaction. The parent rules.When already
// constrains the branch to "not domestic French" (at least one party
// is non-French), so "has Customer" is the remaining axis: present →
// cross-border B2B, absent → B2C (handled by invoiceIsB2CAny).
func invoiceIsCrossBorderB2BAny(v any) bool {
	inv, ok := v.(*bill.Invoice)
	return ok && !invoiceIsB2C(inv)
}

func invoiceVATRatesAllowed(v any) bool {
	inv, ok := v.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	for _, line := range inv.Lines {
		if line == nil {
			continue
		}
		for _, combo := range line.Taxes {
			if combo == nil || combo.Category != tax.CategoryVAT || combo.Percent == nil {
				continue
			}
			if !percentageInList(*combo.Percent, allowedVATRates) {
				return false
			}
		}
	}
	return true
}

func extensionsHaveB2CCategory(v any) bool {
	return extValue(v).Get(ExtKeyB2CCategory) != ""
}

func extensionsHaveBillingMode(v any) bool {
	return extValue(v).Get(ExtKeyBillingMode) != ""
}

func partyHasAllowedLegalScheme(v any) bool {
	party, ok := v.(*org.Party)
	if !ok || party == nil {
		return false
	}
	return slices.Contains(allowedPartySchemeIDs, partyLegalSchemeID(party))
}

func invoiceIsFinalAfterAdvance(v any) bool {
	inv, ok := v.(*bill.Invoice)
	if !ok || inv == nil || inv.Tax == nil || inv.Tax.Ext.IsZero() {
		return false
	}
	return slices.Contains(finalAfterAdvanceBillingModes, inv.Tax.Ext.Get(ExtKeyBillingMode))
}

func invoiceNotAdvancePaymentDocType(v any) bool {
	return !slices.Contains(advancePaymentDocumentTypes, extValue(v).Get(untdid.ExtKeyDocumentType))
}

// invoiceDocumentTypeAllowed reads the untdid-document-type extension
// set by the scenarios and confirms it is one of the permitted codes.
func invoiceDocumentTypeAllowed(v any) bool {
	return slices.Contains(allowedInvoiceDocumentTypes, extValue(v).Get(untdid.ExtKeyDocumentType))
}

// invoiceHasSellerVATIDForExempt returns true if either the supplier
// or the ordering.seller (treated as the supplier's tax representative)
// carries a non-empty TaxID code.
func invoiceHasSellerVATIDForExempt(v any) bool {
	inv, ok := v.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	if partyHasVATCode(inv.Supplier) {
		return true
	}
	if inv.Ordering != nil && partyHasVATCode(inv.Ordering.Seller) {
		return true
	}
	return false
}

func partyHasVATCode(p *org.Party) bool {
	return p != nil && p.TaxID != nil && p.TaxID.Code != ""
}

// invoiceHasExemptCombo reports whether the invoice has any VAT combo
// whose UNTDID 5305 tax-category extension is "E" (exempt).
func invoiceHasExemptCombo(v any) bool {
	inv, ok := v.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	for _, line := range inv.Lines {
		if line == nil {
			continue
		}
		for _, combo := range line.Taxes {
			if combo == nil || combo.Category != tax.CategoryVAT {
				continue
			}
			if combo.Ext.Get(untdid.ExtKeyTaxCategory) == "E" {
				return true
			}
		}
	}
	return false
}

// invoiceHasExemptTaxNote checks for at least one tax.Note with
// Key=exempt and non-empty Text.
func invoiceHasExemptTaxNote(v any) bool {
	inv, ok := v.(*bill.Invoice)
	if !ok || inv == nil || inv.Tax == nil {
		return false
	}
	for _, n := range inv.Tax.Notes {
		if n != nil && n.Key == tax.KeyExempt && n.Text != "" {
			return true
		}
	}
	return false
}

func partyHasTaxIDWhenRequired(v any) bool {
	party, ok := v.(*org.Party)
	if !ok || party == nil {
		return true
	}
	scheme := partyLegalSchemeID(party)
	if !slices.Contains(schemeIDsRequiringVAT, scheme) {
		return true
	}
	return party.TaxID != nil && party.TaxID.Code != ""
}
