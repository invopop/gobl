package zatca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// ZATCA custom tags to be used to
// determine invoice transaction type
const (
	TagSummary    cbc.Key = "summary"
	TagThirdParty cbc.Key = "third-party"
	TagNominal    cbc.Key = "nominal"
)

var tags = []*tax.TagSet{
	{
		Schema: bill.ShortSchemaInvoice,
		List: []*cbc.Definition{
			{
				Key: TagSummary,
				Name: i18n.String{
					i18n.EN: "Summary",
					i18n.AR: "ملخص",
				},
			},
			{
				Key: TagThirdParty,
				Name: i18n.String{
					i18n.EN: "Third-party",
					i18n.AR: "طرف ثالث",
				},
			},
			{
				Key: TagNominal,
				Name: i18n.String{
					i18n.EN: "Nominal",
					i18n.AR: "اسمية",
				},
			},
		},
	},
}
