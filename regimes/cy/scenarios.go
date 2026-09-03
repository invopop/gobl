// Package cy provides a regime definition for Cyprus.
package cy

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Invoice tags specific to Cyprus VAT schemes.
const (
	TagTravelAgencyMargin cbc.Key = "travel-agency-margin"
	TagCashAccounting     cbc.Key = "cash-accounting"
)

var invoiceTags = &tax.TagSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*cbc.Definition{
		{
			Key: TagTravelAgencyMargin,
			Name: i18n.String{
				i18n.EN: "Travel agency margin scheme",
				i18n.EL: "Καθεστώς περιθωρίου ταξιδιωτικών πρακτορείων",
			},
		},
		{
			Key: TagCashAccounting,
			Name: i18n.String{
				i18n.EN: "Cash accounting scheme",
				i18n.EL: "Καθεστώς ταμειακής λογιστικής",
			},
		},
	},
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// ** Special Messages **
		// Reverse Charges
		{
			Tags:       []cbc.Key{tax.TagReverseCharge},
			Categories: []cbc.Code{tax.CategoryVAT},
			Note: &tax.Note{
				Category: tax.CategoryVAT,
				Key:      tax.KeyReverseCharge,
				Text:     "Reverse charge / Αντίστροφη χρέωση",
			},
		},
		// Self-billed invoices
		{
			Tags: []cbc.Key{tax.TagSelfBilled},
			Note: &tax.Note{
				Key:  tax.TagSelfBilled,
				Text: "Self-billing / Αυτοτιμολόγηση",
			},
		},
		// Travel agency margin scheme
		{
			Tags:       []cbc.Key{TagTravelAgencyMargin},
			Categories: []cbc.Code{tax.CategoryVAT},
			Note: &tax.Note{
				Category: tax.CategoryVAT,
				Key:      TagTravelAgencyMargin,
				Text:     "Margin scheme - Travel agents / Καθεστώς περιθωρίου - Ταξιδιωτικά πρακτορεία",
			},
		},
		// Cash accounting scheme
		{
			Tags:       []cbc.Key{TagCashAccounting},
			Categories: []cbc.Code{tax.CategoryVAT},
			Note: &tax.Note{
				Category: tax.CategoryVAT,
				Key:      TagCashAccounting,
				Text:     "Cash accounting / Καθεστώς Ταμειακής Λογιστικής",
			},
		},
	},
}
