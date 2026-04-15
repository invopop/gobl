package nfse

import (
	"fmt"
	"regexp"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/br"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

const (
	// FiscalIncentiveDefault is the default value for the fiscal incentive extenstion
	FiscalIncentiveDefault = "2" // No incentiva
)

var (
	// CodeRegexp is the regular expression used to validate the invoice code
	CodeRegexp = regexp.MustCompile(`^[1-9][0-9]*$`)
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Field("series",
			rules.Assert("01", "series is required", is.Present),
		),
		rules.Field("code",
			rules.Assert("02", "code must be a positive integer", is.Matches(`^[1-9][0-9]*$`)),
		),
		rules.Field("supplier",
			rules.Assert("03", "supplier is required", is.Present),
			rules.Field("tax_id",
				rules.Assert("04", "supplier tax ID is required", is.Present),
				rules.Field("code",
					rules.Assert("05", "supplier tax ID code is required", is.Present),
				),
			),
			rules.Field("name",
				rules.Assert("06", "supplier name is required", is.Present),
			),
			rules.Field("addresses",
				rules.Assert("07", "supplier must have at least one address", is.Present),
				rules.Each(
					rules.Assert("08", "supplier address must not be empty", is.Present),
					rules.Field("street",
						rules.Assert("09", "supplier address requires a street", is.Present),
					),
					rules.Field("num",
						rules.Assert("10", "supplier address requires a number", is.Present),
					),
					rules.Field("locality",
						rules.Assert("11", "supplier address requires a locality", is.Present),
					),
					rules.Field("state",
						rules.Assert("12", "supplier address requires a state", is.Present),
					),
					rules.Field("code",
						rules.Assert("13", "supplier address requires a postal code", is.Present),
					),
				),
			),
			rules.Field("ext",
				rules.Assert("14", fmt.Sprintf("supplier requires '%s', '%s', and '%s' extensions", br.ExtKeyMunicipality, ExtKeySimples, ExtKeyFiscalIncentive),
					tax.ExtensionsRequire(
						br.ExtKeyMunicipality,
						ExtKeySimples,
						ExtKeyFiscalIncentive,
					),
				),
			),
		),
		rules.Field("charges",
			rules.Assert("15", "charges are not supported by NFS-e", is.Empty),
		),
		rules.Field("discounts",
			rules.Assert("16", "discounts are not supported by NFS-e", is.Empty),
		),
	)
}

func normalizeSupplier(sup *org.Party) {
	if sup == nil {
		return
	}

	if !sup.Ext.Has(ExtKeyFiscalIncentive) {
		if sup.Ext == nil {
			sup.Ext = make(tax.Extensions)
		}
		sup.Ext[ExtKeyFiscalIncentive] = FiscalIncentiveDefault
	}
}
