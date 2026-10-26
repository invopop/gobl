package sk

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// billInvoiceRules holds the content that Slovak VAT Act 222/2004 § 74(1) requires
// on a tax invoice. The Act asks for "adresu" of both parties. The Financial
// Administration reads this as street, city and postal code (BT-35, BT-37 and
// BT-38, BT-50, BT-52 and BT-53 "vyžaduje sa" in column R of its BIS mapping),
// although Peppol BIS Billing 3.0 makes these parts optional. They are checked
// here and not left to the format.
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.Field("supplier",
				rules.Assert("01", "invoice SK supplier must have a tax ID code (DIČ or IČ DPH)",
					org.PartyHasTaxIDCode(),
				),
				rules.Assert("02", "invoice SK supplier must have an address with street, locality and code (§ 74 ods. 1 písm. a, column R BT-35, BT-37, BT-38)",
					is.Func("has postal address", partyHasPostalAddress),
				),
			),
			// A simplified invoice does not identify the customer.
			rules.When(
				is.Not(is.Func("is simplified", invoiceIsSimplified)),
				rules.Field("customer",
					rules.Assert("03", "invoice SK customer must have an address with street, locality and code (§ 74 ods. 1 písm. b, column R BT-50, BT-52, BT-53)",
						is.Func("has postal address", partyHasPostalAddress),
					),
				),
			),
			// VAT must be stated in euros (§ 74 ods. 1 písm. i).
			rules.When(
				is.Func("currency is not EUR", invoiceNotInEUR),
				rules.Field("exchange_rates",
					rules.Assert("04", "invoice SK not in EUR must have an exchange rate to EUR to state the VAT in euros (§ 74 ods. 1 písm. i)",
						is.Func("has rate to EUR", hasRateToEUR),
					),
				),
			),
		),
	)
}

// partyHasPostalAddress is true when the party is absent, because the rule only
// applies to a party that exists, or when at least one of its addresses has the
// street, locality and postal code.
func partyHasPostalAddress(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	for _, a := range p.Addresses {
		if a != nil && a.Street != "" && a.Locality != "" && a.Code != "" {
			return true
		}
	}
	return false
}

func invoiceIsSimplified(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && inv.HasTags(tax.TagSimplified)
}

func invoiceNotInEUR(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && inv.Currency != "" && inv.Currency != currency.EUR
}

// hasRateToEUR checks the invoice currency against the rates. The rule only
// runs when the currency is not EUR, and the invoice currency is the "from".
func hasRateToEUR(val any) bool {
	rates, _ := val.([]*currency.ExchangeRate)
	for _, r := range rates {
		if r != nil && r.To == currency.EUR {
			return true
		}
	}
	return false
}
