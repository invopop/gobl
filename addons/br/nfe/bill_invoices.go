package nfe

import (
"github.com/invopop/gobl/bill"
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
// Series
rules.Field("series",
rules.Assert("01", "series is required", is.Present),
rules.Assert("02", "series format is invalid", is.Matches(seriesPattern)),
),
// Tax details
rules.Field("tax",
rules.Assert("03", "tax details are required", is.Present),
rules.Field("ext",
rules.Assert("04", "model extension is required",
tax.ExtensionsRequire(ExtKeyModel),
),
rules.Assert("05", "presence extension is required",
tax.ExtensionsRequire(ExtKeyPresence),
),
),
),
// NFe: delivery presence not allowed
rules.When(is.Func("is NFe", invoiceIsNFe),
rules.Field("tax",
rules.Field("ext",
rules.Assert("06", "delivery presence not allowed for NFe",
tax.ExtensionsExcludeCodes(ExtKeyPresence, PresenceDelivery),
),
),
),
),
// NFCe: presence must be in-person or delivery
rules.When(is.Func("is NFCe", invoiceIsNFCe),
rules.Field("tax",
rules.Field("ext",
rules.Assert("07", "NFCe presence must be in-person or delivery",
tax.ExtensionsHasCodes(ExtKeyPresence, PresenceInPerson, PresenceDelivery),
),
),
),
),
// Notes: must have a reason note; each reason note text must be 1-60 chars
rules.Field("notes",
rules.Assert("08", "a reason note is required",
is.Func("has reason note", notesHasReasonNote),
),
rules.Each(
rules.When(is.Func("is reason note", noteIsReason),
rules.Field("text",
rules.Assert("09", "reason note text must be between 1 and 60 characters",
is.RuneLength(1, 60),
),
),
),
),
),
// Payment required when invoice is not fully paid
rules.When(is.Func("invoice not fully paid", invoiceIsNotPaid),
rules.Field("payment",
rules.Assert("10", "payment details are required", is.Present),
rules.Field("instructions",
rules.Assert("11", "payment instructions are required", is.Present),
),
),
),
// Totals: due must be zero or positive
rules.Field("totals",
rules.Field("due",
rules.Assert("12", "due amount must be zero or positive", num.ZeroOrPositive),
),
),
// NFe: customer is required
rules.When(is.Func("is NFe", invoiceIsNFe),
rules.Field("customer",
rules.Assert("13", "customer is required for NFe", is.Present),
),
),
// NFe: customer addresses are required
rules.When(is.Func("is NFe", invoiceIsNFe),
rules.Field("customer",
rules.Field("addresses",
rules.Assert("14", "customer addresses are required for NFe", is.Present),
),
),
),
// Supplier validation
rules.Field("supplier",
rules.Field("name",
rules.Assert("15", "supplier name is required", is.Present),
),
rules.Field("tax_id",
rules.Assert("16", "supplier tax ID is required", is.Present),
rules.Field("code",
rules.Assert("17", "supplier tax ID code is required", is.Present),
),
),
rules.Field("identities",
rules.Assert("18", "supplier state registration identity is required",
org.IdentitiesKeyIn(IdentityKeyStateReg),
),
),
rules.Field("addresses",
rules.Assert("19", "supplier addresses are required", is.Present),
rules.Each(
rules.Assert("20", "address must not be empty", is.Present),
rules.Field("street", rules.Assert("21", "street is required", is.Present)),
rules.Field("num", rules.Assert("22", "number is required", is.Present)),
rules.Field("locality", rules.Assert("23", "locality is required", is.Present)),
rules.Field("state", rules.Assert("24", "state is required", is.Present)),
rules.Field("code", rules.Assert("25", "address code is required", is.Present)),
),
),
rules.When(is.Func("has addresses", partyHasAddresses),
rules.Field("ext",
rules.Assert("26", "municipality extension is required",
tax.ExtensionsRequire(br.ExtKeyMunicipality),
),
),
),
),
// Customer validation (always applied when customer is present)
rules.Field("customer",
rules.Field("tax_id",
rules.Assert("27", "customer tax ID is required", is.Present),
rules.Field("code",
rules.Assert("28", "customer tax ID code is required", is.Present),
),
),
rules.Field("addresses",
rules.Each(
rules.Assert("29", "address must not be empty", is.Present),
rules.Field("street", rules.Assert("30", "street is required", is.Present)),
rules.Field("num", rules.Assert("31", "number is required", is.Present)),
rules.Field("locality", rules.Assert("32", "locality is required", is.Present)),
rules.Field("state", rules.Assert("33", "state is required", is.Present)),
rules.Field("code", rules.Assert("34", "address code is required", is.Present)),
),
),
rules.When(is.Func("has addresses", partyHasAddresses),
rules.Field("ext",
rules.Assert("35", "municipality extension is required",
tax.ExtensionsRequire(br.ExtKeyMunicipality),
),
),
),
),
)
}

func invoiceIsNFe(val any) bool {
inv, ok := val.(*bill.Invoice)
return ok && inv != nil && inv.Tax != nil && inv.Tax.Ext[ExtKeyModel] == ModelNFe
}

func invoiceIsNFCe(val any) bool {
inv, ok := val.(*bill.Invoice)
return ok && inv != nil && inv.Tax != nil && inv.Tax.Ext[ExtKeyModel] == ModelNFCe
}

func notesHasReasonNote(val any) bool {
notes, ok := val.([]*org.Note)
if !ok {
return false
}
for _, n := range notes {
if n != nil && n.Key == org.NoteKeyReason {
return true
}
}
return false
}

func noteIsReason(val any) bool {
n, ok := val.(*org.Note)
return ok && n != nil && n.Key == org.NoteKeyReason
}

func invoiceIsNotPaid(val any) bool {
inv, ok := val.(*bill.Invoice)
return ok && inv != nil && !inv.Totals.Paid()
}

func partyHasAddresses(val any) bool {
p, ok := val.(*org.Party)
return ok && p != nil && len(p.Addresses) > 0
}
