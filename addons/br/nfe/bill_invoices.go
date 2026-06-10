package nfe

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Validation patterns
const (
	seriesPattern = `^(?:0|[1-9]{1}[0-9]{0,2})$` // extracted from the NFe XSD to validate the series
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Assert("34", "invoice currency must be BRL or provide exchange rate for conversion", currency.CanConvertTo(currency.BRL)),
		// Supplier validation
		rules.Field("supplier",
			rules.Field("tax_id",
				rules.Assert("01", "invoice supplier tax ID is required", is.Present),
				rules.Field("code",
					rules.Assert("02", "invoice supplier tax ID code is required", is.Present),
				),
			),
			rules.Field("addresses",
				rules.Each(
					rules.Assert("03", "invoice supplier address must not be empty", is.Present),
					rules.Field("street",
						rules.Assert("04", "invoice supplier address requires a street", is.Present),
					),
					rules.Field("num",
						rules.Assert("05", "invoice supplier address requires a number", is.Present),
					),
					rules.Field("locality",
						rules.Assert("06", "invoice supplier address requires a locality", is.Present),
					),
					rules.Field("state",
						rules.Assert("07", "invoice supplier address requires a state", is.Present),
					),
					rules.Field("code",
						rules.Assert("08", "invoice supplier address requires a postal code", is.Present),
					),
				),
			),
			rules.When(
				is.Func("has addresses", partyHasAddresses),
				rules.Field("ext",
					rules.Assert("09", fmt.Sprintf("invoice supplier requires '%s' extension when addresses are present", br.ExtKeyMunicipality),
						tax.ExtensionsRequire(br.ExtKeyMunicipality),
					),
				),
			),
			rules.Field("name",
				rules.Assert("10", "invoice supplier name is required", is.Present),
			),
			rules.Field("addresses",
				rules.Assert("11", "invoice supplier must have at least one address", is.Present),
			),
		),
		// Customer: required for NF-e model
		rules.When(
			is.Func("invoice is NFe", invoiceIsNFe),
			rules.Field("customer",
				rules.Assert("12", "invoice customer is required for NF-e invoices", is.Present),
				rules.Field("addresses",
					rules.Assert("13", "invoice customer must have at least one address for NF-e invoices", is.Present),
				),
			),
		),
		// Customer: general party validation when present
		rules.Field("customer",
			rules.Field("tax_id",
				rules.Assert("14", "invoice customer tax ID is required", is.Present),
				rules.Field("code",
					rules.Assert("15", "invoice customer tax ID code is required", is.Present),
				),
			),
			rules.Field("addresses",
				rules.Each(
					rules.Assert("16", "invoice customer address must not be empty", is.Present),
					rules.Field("street",
						rules.Assert("17", "invoice customer address requires a street", is.Present),
					),
					rules.Field("num",
						rules.Assert("18", "invoice customer address requires a number", is.Present),
					),
					rules.Field("locality",
						rules.Assert("19", "invoice customer address requires a locality", is.Present),
					),
					rules.Field("state",
						rules.Assert("20", "invoice customer address requires a state", is.Present),
					),
					rules.Field("code",
						rules.Assert("21", "invoice customer address requires a postal code", is.Present),
					),
				),
			),
			rules.When(
				is.Func("has addresses", partyHasAddresses),
				rules.Field("ext",
					rules.Assert("22", fmt.Sprintf("invoice customer requires '%s' extension when addresses are present", br.ExtKeyMunicipality),
						tax.ExtensionsRequire(br.ExtKeyMunicipality),
					),
				),
			),
		),
		// Series
		rules.Field("series",
			rules.Assert("23", "invoice series is required", is.Present),
			rules.Assert("24", "invoice series format is invalid; must be 0 or 1-999", is.Matches(seriesPattern)),
		),
		// Tax
		rules.Field("tax",
			rules.Assert("25", "invoice tax is required", is.Present),
			rules.Field("ext",
				rules.Assert("26", fmt.Sprintf("invoice tax requires '%s' and '%s' extensions", ExtKeyModel, ExtKeyPresence),
					tax.ExtensionsRequire(ExtKeyModel, ExtKeyPresence),
				),
				rules.When(
					is.Func("NFe model", taxExtIsNFe),
					rules.Assert("27", fmt.Sprintf("NF-e invoices do not support '%s' for '%s'", PresenceDelivery, ExtKeyPresence),
						tax.ExtensionsExcludeCodes(ExtKeyPresence, PresenceDelivery),
					),
				),
				rules.When(
					is.Func("NFCe model", taxExtIsNFCe),
					rules.Assert("28", fmt.Sprintf("NFC-e invoices require in-person or delivery for '%s'", ExtKeyPresence),
						tax.ExtensionsHasCodes(ExtKeyPresence, PresenceInPerson, PresenceDelivery),
					),
				),
			),
		),
		// Notes
		rules.Field("notes",
			rules.Each(
				rules.When(
					is.Func("reason note", isReasonNote),
					rules.Field("text",
						rules.Assert("29", "invoice reason note text must be between 1 and 60 characters", is.Length(1, 60)),
					),
				),
			),
		),
		rules.Assert("30", "invoice requires a note with key 'reason' to describe the nature of the operation (natOp)",
			is.Func("has reason note", invoiceHasReasonNote),
		),
		// Payment: required when unpaid
		rules.When(
			is.Func("invoice not paid", invoiceNotPaid),
			rules.Field("payment",
				rules.Assert("31", "invoice payment is required when invoice is unpaid", is.Present),
				rules.Field("instructions",
					rules.Assert("32", "invoice payment instructions are required when invoice is unpaid", is.Present),
				),
			),
		),
		// Totals
		rules.Field("totals",
			rules.Field("due",
				rules.Assert("33", "invoice due amount must not be negative",
					is.Func("zero or positive", amountZeroOrPositive),
				),
			),
		),
		// Lines: NF-e requires CFOP on each line
		rules.When(
			is.Func("invoice is NFe", invoiceIsNFe),
			rules.Field("lines",
				rules.Each(
					rules.Field("ext",
						rules.Assert("35", fmt.Sprintf("NF-e invoice lines require '%s' extension", ExtKeyCFOP),
							tax.ExtensionsRequire(ExtKeyCFOP),
						),
					),
				),
			),
		),
	)
}

// partyHasAddresses checks if a party has addresses (used in When conditions
// on both supplier and customer).
func partyHasAddresses(val any) bool {
	p, ok := val.(*org.Party)
	return ok && p != nil && len(p.Addresses) > 0
}

// invoiceIsNFe checks if the invoice's tax model is NF-e.
func invoiceIsNFe(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && isNFe(inv.Tax)
}

// taxExtIsNFe checks if the tax extensions indicate NF-e model.
func taxExtIsNFe(val any) bool {
	ext, ok := tax.ExtensionsFromValue(val)
	return ok && ext.Get(ExtKeyModel) == ModelNFe
}

// taxExtIsNFCe checks if the tax extensions indicate NFC-e model.
func taxExtIsNFCe(val any) bool {
	ext, ok := tax.ExtensionsFromValue(val)
	return ok && ext.Get(ExtKeyModel) == ModelNFCe
}

// isReasonNote checks if a note has the reason key.
func isReasonNote(val any) bool {
	note, ok := val.(*org.Note)
	return ok && note != nil && note.Key == org.NoteKeyReason
}

// invoiceHasReasonNote checks if the invoice has at least one note with the reason key.
func invoiceHasReasonNote(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	for _, n := range inv.Notes {
		if n != nil && n.Key == org.NoteKeyReason {
			return true
		}
	}
	return false
}

// invoiceNotPaid checks if the invoice totals indicate the invoice is not paid.
func invoiceNotPaid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	return !inv.Totals.Paid()
}

// amountZeroOrPositive checks that a num.Amount pointer is nil, zero, or positive.
func amountZeroOrPositive(val any) bool {
	amt, ok := val.(*num.Amount)
	if !ok || amt == nil {
		return true
	}
	return !amt.IsNegative()
}

func isNFe(t *bill.Tax) bool {
	return t != nil && t.Ext.Get(ExtKeyModel) == ModelNFe
}
