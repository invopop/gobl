package ctc

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// French CTC-specific tag keys
const (
	// TagFinal indicates a final invoice after advance payments
	TagFinal cbc.Key = "final"

	// TagB2BINT indicates an international B2B invoice requiring e-reporting
	TagB2BINT cbc.Key = "b2b-int"

	// TagArchiveOnly indicates a credit note for internal cancellation only
	TagArchiveOnly cbc.Key = "archive-only"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.Definition{
		{
			Key: TagFinal,
			Name: i18n.String{
				i18n.EN: "Final Invoice",
				i18n.FR: "Facture définitive",
			},
		},
		{
			Key: TagB2BINT,
			Name: i18n.String{
				i18n.EN: "International B2B",
				i18n.FR: "B2B International",
			},
		},
		{
			Key: TagArchiveOnly,
			Name: i18n.String{
				i18n.EN: "Archive Only",
				i18n.FR: "Archivage uniquement",
			},
		},
	},
}
