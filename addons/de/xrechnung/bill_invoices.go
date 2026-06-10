package xrechnung

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// BR-DE-17 - restricted subset of UNTDID document type codes
var validInvoiceUNTDIDDocumentTypeValues = []cbc.Code{
	"326", // Partial
	"380", // Commercial
	"384", // Corrected
	"389", // Self-billed
	"381", // Credit note
	"875", // Partial construction invoice
	"876", // Partial Final construction invoice
	"877", // Final construction invoice
}

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Field("tax",
			rules.Field("ext",
				rules.Assert("01", "tax ext must have a valid UNTDID document type code (BR-DE-17)",
					tax.ExtensionsHasCodes(untdid.ExtKeyTaxCategory, validInvoiceUNTDIDDocumentTypeValues...),
				),
			),
		),
		rules.When(
			bill.InvoiceTypeIn(bill.InvoiceTypeCorrective, bill.InvoiceTypeCreditNote),
			rules.Field("preceding",
				rules.Assert("02", "preceding documents are required for corrective and credit note invoices (BR-DE-26)", is.Present),
			),
		),
		rules.Field("supplier",
			rules.Assert("03", "supplier is required", is.Present),
			rules.Field("people",
				rules.Assert("04", "supplier people is required (BR-DE-5)", is.Present),
			),
			rules.Field("addresses",
				rules.Assert("05", "supplier addresses is required (BR-DE-2)", is.Present),
				rules.Assert("06", "supplier first address must have a locality (BR-DE-3)",
					is.Func("has locality", firstAddressHasLocality),
				),
				rules.Assert("07", "supplier first address must have a postal code (BR-DE-4)",
					is.Func("has code", firstAddressHasCode),
				),
			),
			rules.Field("inboxes",
				rules.Assert("08", "supplier inboxes are required (PEPPOL-EN16931-R020)", is.Present),
			),
			rules.Assert("09", "either party.telephones or party.people[0].telephones is required (BR-DE-6)",
				is.Func("has telephones", partyHasTelephones),
			),
			rules.Assert("10", "either party.emails or party.people[0].emails is required (BR-DE-7)",
				is.Func("has emails", partyHasEmails),
			),
		),
		rules.Field("customer",
			rules.Assert("11", "customer is required", is.Present),
			rules.Field("addresses",
				rules.Assert("12", "customer addresses are required (BR-DE-8)", is.Present),
				rules.Assert("13", "customer first address must have a locality (BR-DE-9)",
					is.Func("has locality", firstAddressHasLocality),
				),
				rules.Assert("14", "customer first address must have a postal code (BR-DE-10)",
					is.Func("has code", firstAddressHasCode),
				),
			),
			rules.Field("inboxes",
				rules.Assert("15", "customer inboxes are required (PEPPOL-EN16931-R010)", is.Present),
			),
		),
		rules.Field("payment",
			rules.Assert("16", "payment is required (BR-DE-1)", is.Present),
			rules.Field("instructions",
				rules.Assert("17", "payment instructions are required (BR-DE-1)", is.Present),
				rules.When(
					is.Func("has credit-transfer key", instructionsHasCreditTransferKey),
					rules.Field("credit_transfer",
						rules.Assert("18", "credit transfer is required for credit-transfer payments (BR-DE-23)", is.Present),
						rules.Each(
							rules.When(
								is.Func("no IBAN", creditTransferHasNoIBAN),
								rules.Field("number",
									rules.Assert("19", "account number is required when IBAN is not provided (BR-DE-19)", is.Present),
								),
							),
						),
					),
				),
				rules.When(
					is.Func("has card key", instructionsHasCardKey),
					rules.Field("card",
						rules.Assert("20", "card details are required for card payments (BR-DE-24)", is.Present),
					),
				),
				rules.When(
					is.Func("has direct-debit key", instructionsHasDirectDebitKey),
					rules.Field("direct_debit",
						rules.Assert("21", "direct debit details are required for direct-debit payments (BR-DE-25)", is.Present),
						rules.Field("ref",
							rules.Assert("22", "direct debit mandate reference is required (BR-DE-29)", is.Present),
						),
						rules.Field("creditor",
							rules.Assert("23", "direct debit creditor identifier is required (BR-DE-30)", is.Present),
						),
						rules.Field("account",
							rules.Assert("24", "direct debit account is required (BR-DE-31)", is.Present),
						),
					),
				),
			),
		),
		rules.Field("delivery",
			rules.Field("receiver",
				rules.Assert("25", "delivery receiver is required", is.Present),
				rules.Field("addresses",
					rules.Assert("26", "delivery receiver addresses are required", is.Present),
					rules.Assert("27", "delivery receiver first address must have a locality (BR-DE-11)",
						is.Func("has locality", firstAddressHasLocality),
					),
					rules.Assert("28", "delivery receiver first address must have a postal code (BR-DE-12)",
						is.Func("has code", firstAddressHasCode),
					),
				),
			),
		),
		rules.Field("ordering",
			rules.Assert("29", "ordering is required (BR-DE-15)", is.Present),
			rules.Field("code",
				rules.Assert("30", "ordering code is required (BR-DE-15)", is.Present),
			),
		),
	)
}

// Address helpers

func firstAddressHasLocality(val any) bool {
	addrs, ok := val.([]*org.Address)
	if !ok || len(addrs) == 0 {
		return true // handled by Required check
	}
	return addrs[0].Locality != ""
}

func firstAddressHasCode(val any) bool {
	addrs, ok := val.([]*org.Address)
	if !ok || len(addrs) == 0 {
		return true // handled by Required check
	}
	return addrs[0].Code != ""
}

// Supplier helpers

func partyHasTelephones(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	if len(p.Telephones) > 0 {
		return true
	}
	if len(p.People) > 0 && len(p.People[0].Telephones) > 0 {
		return true
	}
	return false
}

func partyHasEmails(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	if len(p.Emails) > 0 {
		return true
	}
	if len(p.People) > 0 && len(p.People[0].Emails) > 0 {
		return true
	}
	return false
}

// Payment instruction helpers

func instructionsHasCreditTransferKey(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return false
	}
	return instr.Key.Has("credit-transfer")
}

func instructionsHasCardKey(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return false
	}
	return instr.Key.Has("card")
}

func instructionsHasDirectDebitKey(val any) bool {
	instr, ok := val.(*pay.Instructions)
	if !ok || instr == nil {
		return false
	}
	return instr.Key.Has("direct-debit")
}

func creditTransferHasNoIBAN(val any) bool {
	ct, ok := val.(*pay.CreditTransfer)
	return ok && ct != nil && ct.IBAN == ""
}
