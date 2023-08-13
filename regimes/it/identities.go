package it

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Italian identity keys required by the Italian tax authority (SDI)
// and defined as part of the FatturaPA specification.
const (
	IdentityKeySDIFiscalRegime = "it-sdi-fiscal-regime"
)

var identityKeys = []*tax.KeyDefinition{
	{
		Key: IdentityKeySDIFiscalRegime,
		Name: i18n.String{
			i18n.EN: "Fiscal Regime Code",
			i18n.IT: "Codice Regime Fiscale",
		},
		Codes: []*tax.CodeDefinition{
			{Code: "RF01", Name: i18n.String{i18n.EN: "Ordinary"}},
			{Code: "RF02", Name: i18n.String{i18n.EN: "Minimum taxpayers (Art. 1, section 96-117, Italian Law 244/07)"}},
			{Code: "RF04", Name: i18n.String{i18n.EN: "Agriculture and connected activities and fishing (Arts. 34 and 34-bis, Italian Presidential Decree 633/72)"}},
			{Code: "RF05", Name: i18n.String{i18n.EN: "Sale of salts and tobaccos (Art. 74, section 1, Italian Presidential Decree 633/72)"}},
			{Code: "RF06", Name: i18n.String{i18n.EN: "Match sales (Art. 74, section 1, Italian Presidential Decree 633/72)"}},
			{Code: "RF07", Name: i18n.String{i18n.EN: "Publishing (Art. 74, section 1, Italian Presidential Decree 633/72)"}},
			{Code: "RF08", Name: i18n.String{i18n.EN: "Management of public telephone services (Art. 74, section 1, Italian Presidential Decree 633/72)"}},
			{Code: "RF09", Name: i18n.String{i18n.EN: "Resale of public transport and parking documents (Art. 74, section 1, Italian Presidential Decree 633/72)"}},
			{Code: "RF10", Name: i18n.String{i18n.EN: "Entertainment, gaming and other activities referred to by the tariff attached to Italian Presidential Decree 640/72 (Art. 74, section 6, Italian Presidential Decree 633/72)"}},
			{Code: "RF11", Name: i18n.String{i18n.EN: "Travel and tourism agencies (Art. 74-ter, Italian Presidential Decree 633/72)"}},
			{Code: "RF12", Name: i18n.String{i18n.EN: "Farmhouse accommodation/restaurants (Art. 5, section 2, Italian law 413/91)"}},
			{Code: "RF13", Name: i18n.String{i18n.EN: "Door-to-door sales (Art. 25-bis, section 6, Italian Presidential Decree 600/73)"}},
			{Code: "RF14", Name: i18n.String{i18n.EN: "Resale of used goods, artworks, antiques or collector's items (Art. 36, Italian Decree Law 41/95)"}},
			{Code: "RF15", Name: i18n.String{i18n.EN: "Artwork, antiques or collector's items auction agencies (Art. 40-bis, Italian Decree Law 41/95)"}},
			{Code: "RF16", Name: i18n.String{i18n.EN: "VAT paid in cash by P.A. (Art. 6, section 5, Italian Presidential Decree 633/72)"}},
			{Code: "RF17", Name: i18n.String{i18n.EN: "VAT paid in cash by subjects with business turnover below Euro 200,000 (Art. 7, Italian Decree Law 185/2008)"}},
			{Code: "RF18", Name: i18n.String{i18n.EN: "Other"}},
			{Code: "RF19", Name: i18n.String{i18n.EN: "Flat rate (Art. 1, section 54-89, Italian Law 190/2014)"}},
		},
	},
}
