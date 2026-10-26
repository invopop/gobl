package sk

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Tax tags that can be applied in Slovakia. The self-billing and reverse charge
// tags come from the base invoice tags.
const (
	TagTravelAgency    cbc.Key = "travel-agency"
	TagSecondHandGoods cbc.Key = "second-hand-goods"
	TagArt             cbc.Key = "art"
	TagAntiques        cbc.Key = "antiques"
	TagCashBasis       cbc.Key = "cash-basis"
)

func invoiceTags() *tax.TagSet {
	return &tax.TagSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*cbc.Definition{
			{
				Key: TagTravelAgency,
				Name: i18n.String{
					i18n.EN: "Special scheme for travel agencies",
					i18n.SK: "Úprava zdaňovania prirážky - cestovné kancelárie",
				},
			},
			{
				Key: TagSecondHandGoods,
				Name: i18n.String{
					i18n.EN: "Special scheme for second-hand goods",
					i18n.SK: "Úprava zdaňovania prirážky - použitý tovar",
				},
			},
			{
				Key: TagArt,
				Name: i18n.String{
					i18n.EN: "Special scheme for works of art",
					i18n.SK: "Úprava zdaňovania prirážky - umelecké diela",
				},
			},
			{
				Key: TagAntiques,
				Name: i18n.String{
					i18n.EN: "Special scheme for collectors' items and antiques",
					i18n.SK: "Úprava zdaňovania prirážky - zberateľské predmety a starožitnosti",
				},
			},
			{
				Key: TagCashBasis,
				Name: i18n.String{
					i18n.EN: "VAT on a cash basis",
					i18n.SK: "Daň sa uplatňuje na základe prijatia platby",
				},
			},
		},
	}
}

// invoiceScenarios adds the wording that VAT Act 222/2004 § 74(1) j) to n) and
// § 68d(4) require on the invoice in each special situation. The text is the
// phrase of the law. For exempt supplies, § 74(1) h) also accepts a reference to
// the provision, which formats carry as the exemption reason (e.g. VATEX in
// Peppol BIS), so no note is added.
func invoiceScenarios() *tax.ScenarioSet {
	return &tax.ScenarioSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*tax.Scenario{
			// § 74 ods. 1 písm. k)
			{
				Tags:       []cbc.Key{tax.TagReverseCharge},
				Categories: []cbc.Code{tax.CategoryVAT},
				Note: &tax.Note{
					Category: tax.CategoryVAT,
					Key:      tax.KeyReverseCharge,
					Text:     "Prenesenie daňovej povinnosti",
				},
			},
			// § 74 ods. 1 písm. j)
			{
				Tags: []cbc.Key{tax.TagSelfBilled},
				Note: &tax.Note{
					Key:  tax.TagSelfBilled,
					Text: "Vyhotovenie faktúry odberateľom",
				},
			},
			// § 74 ods. 1 písm. m)
			{
				Tags: []cbc.Key{TagTravelAgency},
				Note: &tax.Note{
					Key:  TagTravelAgency,
					Text: "Úprava zdaňovania prirážky - cestovné kancelárie",
				},
			},
			// § 74 ods. 1 písm. n)
			{
				Tags: []cbc.Key{TagSecondHandGoods},
				Note: &tax.Note{
					Key:  TagSecondHandGoods,
					Text: "Úprava zdaňovania prirážky - použitý tovar",
				},
			},
			{
				Tags: []cbc.Key{TagArt},
				Note: &tax.Note{
					Key:  TagArt,
					Text: "Úprava zdaňovania prirážky - umelecké diela",
				},
			},
			{
				Tags: []cbc.Key{TagAntiques},
				Note: &tax.Note{
					Key:  TagAntiques,
					Text: "Úprava zdaňovania prirážky - zberateľské predmety a starožitnosti",
				},
			},
			// § 68d ods. 4: without this phrase the tax liability arises under § 19.
			{
				Tags: []cbc.Key{TagCashBasis},
				Note: &tax.Note{
					Key:  TagCashBasis,
					Text: "Daň sa uplatňuje na základe prijatia platby",
				},
			},
		},
	}
}
