package si

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Article 82 of the Slovenian VAT Act (ZDDV-1) requires invoices to state
// the VAT identification number under which the taxable person supplied
// the goods or services.
// Source: https://spot.gov.si/en/info/accountancy/invoicing
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.When(
			is.InContext(tax.RegimeIn(CountryCode)),
			rules.Field("supplier",
				rules.Assert("01", "invoice SI supplier must have a tax ID code",
					is.Func("has tax ID code", hasSupplierTaxIDCode),
				),
			),
		),
	)
}

func hasSupplierTaxIDCode(value any) bool {
	party, _ := value.(*org.Party)
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}
