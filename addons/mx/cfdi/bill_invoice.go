package cfdi

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func normalizeInvoice(inv *bill.Invoice) {
	normalizeInvoiceIssueDateAndTime(inv)
	if inv.Tags.HasTags(TagGlobal) {
		inv.Customer = nil
	}
}

func normalizeInvoiceIssueDateAndTime(inv *bill.Invoice) {
	// Overwrite the issue date and time to align with
	// CFDI requirements for the emission date, unless the
	// issue time is already set.
	if inv.IssueTime != nil && !inv.IssueTime.IsZero() {
		return
	}
	tz := inv.RegimeDef().TimeLocation()
	dn := cal.ThisSecondIn(tz)
	tn := dn.Time()
	inv.IssueDate = dn.Date()
	inv.IssueTime = &tn
}

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Assert("27", "invoice must be in MXN or provide exchange rate for conversion", currency.CanConvertTo(currency.MXN)),
		// Tax extensions
		rules.Field("tax",
			rules.Field("ext",
				rules.Assert("01",
					fmt.Sprintf("tax requires '%s', '%s', and '%s' extensions", ExtKeyDocType, ExtKeyIssuePlace, ExtKeyPaymentMethod),
					tax.ExtensionsRequire(ExtKeyDocType, ExtKeyIssuePlace, ExtKeyPaymentMethod),
				),
			),
		),
		// Global tag: require global extensions
		rules.When(is.Func("has global tag", invoiceHasGlobalTag),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("02",
						fmt.Sprintf("global invoices require '%s', '%s', and '%s' extensions",
							ExtKeyGlobalPeriod, ExtKeyGlobalMonth, ExtKeyGlobalYear),
						tax.ExtensionsRequire(ExtKeyGlobalPeriod, ExtKeyGlobalMonth, ExtKeyGlobalYear),
					),
				),
			),
			rules.Field("customer",
				rules.Assert("03", "cannot be set with global tag", is.Empty),
			),
			rules.Field("payment",
				rules.Assert("04", "payment is required for global invoices", is.Present),
				rules.Field("advances",
					rules.Assert("05", "advances must be set with global tag", is.Present),
				),
			),
			rules.Field("lines",
				rules.Each(
					rules.Field("item",
						rules.Field("ref",
							rules.Assert("06", "must be set with global tag", is.Present),
						),
					),
				),
			),
		),
		// Non-global tag: global extensions must be all-or-none
		rules.When(is.Func("no global tag", invoiceNoGlobalTag),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("07",
						fmt.Sprintf("'%s', '%s', and '%s' extensions must all be present or all absent",
							ExtKeyGlobalPeriod, ExtKeyGlobalMonth, ExtKeyGlobalYear),
						tax.ExtensionsRequireAllOrNone(ExtKeyGlobalPeriod, ExtKeyGlobalMonth, ExtKeyGlobalYear),
					),
				),
			),
			// Customer validation (non-global only)
			rules.Field("customer",
				rules.Field("tax_id",
					rules.Assert("08", "customer tax ID is required", is.Present),
					rules.Field("code",
						rules.Assert("09", "customer tax ID code is required", is.Present),
					),
				),
				rules.When(is.Func("customer is Mexican", partyIsMexican),
					rules.Field("ext",
						rules.Assert("10",
							fmt.Sprintf("Mexican customer requires '%s' and '%s' extensions", ExtKeyFiscalRegime, ExtKeyUse),
							tax.ExtensionsRequire(ExtKeyFiscalRegime, ExtKeyUse),
						),
					),
					rules.Field("addresses",
						rules.Assert("11", "Mexican customer must have at least one address", is.Present),
						rules.Each(
							rules.Field("code",
								rules.Assert("12", "customer address postal code is required", is.Present),
								rules.Assert("13", "customer address postal code format is invalid",
									is.Matches(PostCodePattern),
								),
							),
						),
					),
				),
			),
			// Line item extensions (non-global only)
			rules.Field("lines",
				rules.Each(
					rules.Field("item",
						rules.Field("ext",
							rules.Assert("14",
								fmt.Sprintf("item requires '%s' extension", ExtKeyProdServ),
								tax.ExtensionsRequire(ExtKeyProdServ),
							),
							rules.Assert("15", "product/service code must have 8 digits",
								is.Func("valid prod-serv", itemExtProdServValid),
							),
						),
					),
				),
			),
		),
		// Preceding validation (always)
		rules.When(is.Func("has preceding", invoiceHasPreceding),
			rules.Field("tax",
				rules.Field("ext",
					rules.Assert("16",
						fmt.Sprintf("tax requires '%s' extension when preceding documents are present", ExtKeyRelType),
						tax.ExtensionsRequire(ExtKeyRelType),
					),
				),
			),
		),
		rules.Field("preceding",
			rules.Each(
				rules.Field("stamps",
					rules.Assert("17", fmt.Sprintf("preceding row is missing '%s' stamp", mx.StampSATUUID),
						head.StampsHas(mx.StampSATUUID),
					),
				),
			),
		),
		// Supplier validation
		rules.Field("supplier",
			rules.Field("tax_id",
				rules.Assert("18", "supplier tax ID is required", is.Present),
				rules.Field("code",
					rules.Assert("19", "supplier tax ID code is required", is.Present),
				),
			),
			rules.Field("ext",
				rules.Assert("20",
					fmt.Sprintf("supplier requires '%s' extension", ExtKeyFiscalRegime),
					tax.ExtensionsRequire(ExtKeyFiscalRegime),
				),
			),
		),
		// Line validation (always)
		rules.Field("lines",
			rules.Each(
				rules.Field("quantity",
					rules.Assert("21", "line quantity must be greater than 0", num.Positive),
				),
				rules.Field("item",
					rules.Field("price",
						rules.Assert("22", "line item price is required", is.Present),
						rules.Assert("23", "line item price must be greater than 0", num.Positive),
					),
				),
				rules.Field("total",
					rules.Assert("24", "line total must not be negative", num.Min(num.AmountZero)),
				),
			),
		),
		// Discounts and charges not supported
		rules.Field("discounts",
			rules.Assert("25", "document level discounts are not supported, use line discounts instead", is.Empty),
		),
		rules.Field("charges",
			rules.Assert("26", "document level charges are not supported", is.Empty),
		),
	)
}

// invoiceHasGlobalTag checks if the invoice has the global tag.
func invoiceHasGlobalTag(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && inv.HasTags(TagGlobal)
}

// invoiceNoGlobalTag checks if the invoice does NOT have the global tag.
func invoiceNoGlobalTag(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && !inv.HasTags(TagGlobal)
}

// invoiceHasPreceding checks if the invoice has preceding documents.
func invoiceHasPreceding(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && len(inv.Preceding) > 0
}

// partyIsMexican checks if a party has a Mexican tax ID.
func partyIsMexican(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil || p.TaxID == nil {
		return false
	}
	return p.TaxID.Country.In("MX")
}

// itemExtProdServValid checks that the ProdServ extension code has 8 digits.
// Skips when the extension is not present (handled by a separate assertion).
func itemExtProdServValid(val any) bool {
	ext, ok := val.(tax.Extensions)
	if !ok {
		return true
	}
	v, has := ext[ExtKeyProdServ]
	if !has {
		return true // not present, other rule handles this
	}
	return itemExtensionValidCodeRegexp.MatchString(string(v))
}
