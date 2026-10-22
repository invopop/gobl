package hr

import (
	"slices"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// Tax tags that can be applied in Croatia.
const (
	// TagCashBasis marks an invoice issued by a supplier that applies the
	// cash accounting scheme (Article 125.i of the VAT Act).
	TagCashBasis cbc.Key = "cash-basis"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.Definition{
		{
			Key: TagCashBasis,
			Name: i18n.String{
				i18n.EN: "Cash accounting scheme",
				i18n.HR: "Postupak oporezivanja prema naplaćenim naknadama",
			},
			Desc: i18n.String{
				i18n.EN: here.Doc(`
					The supplier accounts for VAT when payment is received rather than
					when the invoice is issued, under Article 125.i of the VAT Act. The
					Act requires such invoices to state "obračun prema naplaćenim
					naknadama", which this tag adds as a VAT note.
				`),
			},
		},
	},
}

// invoiceScenarios extends the default invoice scenarios, which provide the
// reverse charge note, with the Croatian ones.
var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: slices.Concat(bill.InvoiceScenarios().List, []*tax.Scenario{
		// Cash accounting scheme
		{
			Tags:       []cbc.Key{TagCashBasis},
			Categories: []cbc.Code{tax.CategoryVAT},
			Note: &tax.Note{
				Category: tax.CategoryVAT,
				Key:      TagCashBasis,
				Text:     "Obračun prema naplaćenim naknadama",
			},
		},
	}),
}
