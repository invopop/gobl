package sdi

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/it"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func normalizeInvoice(inv *bill.Invoice) {
	normalizeSupplier(inv.Supplier)
}

func normalizeSupplier(party *org.Party) {
	if party == nil {
		return
	}
	if party.Ext == nil || party.Ext[ExtKeyFiscalRegime] == "" {
		if party.Ext == nil {
			party.Ext = make(tax.Extensions)
		}
		party.Ext[ExtKeyFiscalRegime] = "RF01" // Ordinary regime is default
	}

	// Normalize Italian supplier telephone numbers by stripping '+39' prefix
	if isItalianParty(party) && len(party.Telephones) > 0 {
		for _, tel := range party.Telephones {
			if tel != nil && len(tel.Number) >= 3 && tel.Number[:3] == "+39" {
				tel.Number = tel.Number[3:]
			}
		}
	}
}

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Field("tax",
			rules.Assert("01", "tax is required", is.Present),
			rules.Field("ext",
				rules.Assert("02",
					fmt.Sprintf("tax requires '%s' and '%s' extensions", ExtKeyDocumentType, ExtKeyFormat),
					tax.ExtensionsRequire(
						ExtKeyFormat,
						ExtKeyDocumentType,
					),
				),
			),
		),
		rules.Field("supplier",
			rules.Field("name",
				rules.Assert("03", "supplier name must use Latin-1 characters",
					is.FuncError("latin1", validateLatin1String),
				),
			),
			rules.Field("addresses",
				rules.Assert("04", "supplier addresses are required", is.Present),
			),
			rules.Field("ext",
				rules.Assert("05",
					fmt.Sprintf("supplier requires '%s' extension", ExtKeyFiscalRegime),
					tax.ExtensionsRequire(ExtKeyFiscalRegime),
				),
			),
			rules.Field("registration",
				rules.Field("entry",
					rules.Assert("06", "supplier registration entry is required when registration is present",
						is.Present,
					),
				),
				rules.Field("office",
					rules.Assert("07", "supplier registration office is required when registration is present",
						is.Present,
					),
				),
			),
		),
		rules.When(is.Func("supplier is Italian", invoiceSupplierIsItalian),
			rules.Field("supplier",
				rules.Field("telephones",
					rules.Each(
						rules.Field("num",
							rules.Assert("08", "Italian telephone number length must be between 5 and 12",
								is.Length(5, 12),
							),
						),
					),
				),
			),
		),
		rules.Field("customer",
			rules.Assert("09", "customer is required", is.Present),
			rules.Field("name",
				rules.Assert("10", "customer name must use Latin-1 characters",
					is.FuncError("latin1", validateLatin1String),
				),
			),
			rules.Field("tax_id",
				rules.Assert("11", "customer tax ID is required", is.Present),
			),
			rules.Field("addresses",
				rules.Assert("12", "customer addresses are required", is.Present),
			),
		),
		// Customer name required when tax_id code is present or people is nil
		rules.Assert("13", "customer name is required",
			is.Func("customer name check", invoiceCustomerHasNameOrPeople),
		),
		// Customer people required when name is empty
		rules.Assert("14", "customer people are required when name is empty",
			is.Func("customer people check", invoiceCustomerHasPeopleOrName),
		),
		// Customer tax_id code required for Italian parties without fiscal code
		rules.When(is.Func("Italian customer without fiscal code", invoiceCustomerIsItalianWithoutFiscalCode),
			rules.Field("customer",
				rules.Field("tax_id",
					rules.Field("code",
						rules.Assert("15", "customer tax ID code is required for Italian parties without fiscal code",
							is.Present,
						),
					),
				),
			),
		),
		// Customer identity required for Italian parties without tax ID code
		rules.When(is.Func("Italian customer without tax ID code", invoiceCustomerIsItalianWithoutTaxIDCode),
			rules.Field("customer",
				rules.Field("identities",
					rules.Assert("16",
						fmt.Sprintf("customer requires identity with key '%s'", it.IdentityKeyFiscalCode),
						is.Func("has fiscal code", invoiceCustomerHasFiscalCodeIdentity),
					),
				),
			),
		),
		rules.Field("lines",
			rules.Each(
				rules.Assert("17", "line must have VAT tax category",
					bill.RequireLineTaxCategory(tax.CategoryVAT),
				),
				rules.Field("item",
					rules.Field("name",
						rules.Assert("18", "item name must use Latin-1 characters",
							is.FuncError("latin1", validateLatin1String),
						),
					),
				),
			),
		),
		// Ordering despatch validation
		rules.When(is.Func("has deferred tag", invoiceHasDeferredTag),
			rules.Field("ordering",
				rules.Field("despatch",
					rules.Each(
						rules.Field("issue_date",
							rules.Assert("19", "despatch issue date is required", is.Present),
						),
					),
				),
			),
		),
		rules.When(is.Func("no deferred tag", invoiceDoesNotHaveDeferredTag),
			rules.Field("ordering",
				rules.Field("despatch",
					rules.Assert("20", "despatch can only be set when invoice has deferred tag", is.Empty),
				),
			),
		),
		// Payment: instructions required when terms have due dates
		rules.Assert("21", "payment instructions are required when terms with due dates are present",
			is.Func("payment instructions check", invoicePaymentInstructionsPresent),
		),
	)
}

func billChargeRules() *rules.Set {
	return rules.For(new(bill.Charge),
		rules.When(is.Func("is fund contribution", chargeIsFundContribution),
			rules.Field("percent",
				rules.Assert("01", "fund contribution charge requires a percentage", is.Present),
			),
			rules.Field("ext",
				rules.Assert("02",
					fmt.Sprintf("fund contribution charge requires '%s' extension", ExtKeyFundType),
					tax.ExtensionsRequire(ExtKeyFundType),
				),
			),
			rules.Field("taxes",
				rules.Assert("03", "fund contribution charge must have VAT tax category",
					tax.SetHasCategory(tax.CategoryVAT),
				),
			),
		),
	)
}

// --- Helper functions for rules ---

func invoiceSupplierIsItalian(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	return isItalianParty(inv.Supplier)
}

func invoiceCustomerHasNameOrPeople(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Customer == nil {
		return true
	}
	c := inv.Customer
	// Name required when tax_id code is present or people is nil
	if (c.TaxID != nil && c.TaxID.Code != cbc.CodeEmpty) || c.People == nil {
		return c.Name != ""
	}
	return true
}

func invoiceCustomerHasPeopleOrName(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Customer == nil {
		return true
	}
	c := inv.Customer
	if c.Name == "" {
		return len(c.People) > 0
	}
	return true
}

func invoiceCustomerIsItalianWithoutFiscalCode(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Customer == nil {
		return false
	}
	return isItalianParty(inv.Customer) && !hasFiscalCode(inv.Customer)
}

func invoiceCustomerIsItalianWithoutTaxIDCode(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Customer == nil {
		return false
	}
	return isItalianParty(inv.Customer) && !hasTaxIDCode(inv.Customer)
}

func invoiceCustomerHasFiscalCodeIdentity(val any) bool {
	ids, ok := val.([]*org.Identity)
	if !ok {
		return false
	}
	return org.IdentityForKey(ids, it.IdentityKeyFiscalCode) != nil
}

func invoiceHasDeferredTag(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	return inv.HasTags(TagDeferred)
}

func invoiceDoesNotHaveDeferredTag(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	return !inv.HasTags(TagDeferred)
}

func invoicePaymentInstructionsPresent(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Payment == nil {
		return true
	}
	p := inv.Payment
	if p.Terms != nil && len(p.Terms.DueDates) > 0 {
		return p.Instructions != nil
	}
	return true
}

func chargeIsFundContribution(val any) bool {
	c, ok := val.(*bill.Charge)
	if !ok || c == nil {
		return false
	}
	return c.Key.Has(KeyFundContribution)
}

func hasTaxIDCode(party *org.Party) bool {
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

func hasFiscalCode(party *org.Party) bool {
	if party == nil {
		return false
	}
	return org.IdentityForKey(party.Identities, it.IdentityKeyFiscalCode) != nil

}

func isItalianParty(party *org.Party) bool {
	if party == nil || party.TaxID == nil {
		return false
	}
	return party.TaxID.Country.In("IT")
}
