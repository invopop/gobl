package flow10

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Flow 10 extension keys.
const (
	// ExtKeyB2CCategory classifies the B2C transaction for PPF reporting
	// (G1.68). Required on Flow 10 B2C invoices.
	ExtKeyB2CCategory cbc.Key = "fr-ctc-flow10-b2c-category"
)

// B2C transaction category codes (G1.68).
const (
	// B2CCategoryGoods — deliveries of goods subject to VAT.
	B2CCategoryGoods cbc.Code = "TLB1"
	// B2CCategoryServices — services subject to VAT.
	B2CCategoryServices cbc.Code = "TPS1"
	// B2CCategoryNotTaxable — deliveries / services not subject to VAT
	// in France, including intra-EU distance sales under CGI art. 258 A
	// / 259 B.
	B2CCategoryNotTaxable cbc.Code = "TNT1"
	// B2CCategoryMargin — operations under the VAT-on-margin regime
	// (CGI art. 266-1-e, 268, 297 A).
	B2CCategoryMargin cbc.Code = "TMA1"
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyB2CCategory,
		Name: i18n.String{
			i18n.EN: "B2C Transaction Category",
			i18n.FR: "Catégorie de transaction B2C",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Classifies a B2C transaction for French e-reporting to the PPF
				(G1.68). Required on Flow 10 B2C invoices.

				- TLB1: Goods deliveries subject to VAT.
				- TPS1: Services subject to VAT.
				- TNT1: Goods / services not subject to French VAT,
				  including intra-EU distance sales per CGI articles 258 A
				  and 259 B.
				- TMA1: Operations under the VAT-on-margin regime
				  (CGI articles 266-1-e, 268, 297 A).
			`),
			i18n.FR: here.Doc(`
				Catégorie de transaction pour le e-reporting au PPF (G1.68).
				Obligatoire sur les factures B2C en Flux 10.
			`),
		},
		Values: []*cbc.Definition{
			{Code: B2CCategoryGoods, Name: i18n.String{i18n.EN: "Goods subject to VAT", i18n.FR: "Livraisons de biens soumises à la TVA"}},
			{Code: B2CCategoryServices, Name: i18n.String{i18n.EN: "Services subject to VAT", i18n.FR: "Prestations de services soumises à la TVA"}},
			{Code: B2CCategoryNotTaxable, Name: i18n.String{i18n.EN: "Not subject to French VAT", i18n.FR: "Non soumis à la TVA en France"}},
			{Code: B2CCategoryMargin, Name: i18n.String{i18n.EN: "VAT-on-margin regime", i18n.FR: "Régime de TVA sur la marge"}},
		},
	},
}
