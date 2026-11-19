// Package ad provides the tax region definition for Andorra.
package ad

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Andorra.
const CountryCode = "AD"

// TaxCategoryIGI is the Andorran Impost General Indirecte tax category.
const TaxCategoryIGI cbc.Code = "IGI"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("ad", rules.GOBL.Add(CountryCode),
		billInvoiceRules(),
		taxIdentityRules(),
	)
}

// New provides the tax region definition for Andorra.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:  CountryCode,
		Currency: currency.EUR,
		Name: i18n.String{
			i18n.EN: "Andorra",
			i18n.ES: "Andorra",
			i18n.CA: "Andorra",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("Llei 11/2012, del 21 de juny, de l'impost general indirecte"),
				URL:   "https://www.consellgeneral.ad/ca/arxiu/arxiu-de-lleis-i-textos-aprovats-en-legislatures-anteriors/vi-legislatura-2011-2015/copy_of_lleis-aprovades/decret-legislatiu-del-23-07-2014-de-publicacio-del-text-refos-de-la-llei-11-2012-del-21-de-juny-de-l2019impost-general-indirecte/at_download/PDF",
			},
		},
		TimeZone: "Europe/Andorra",
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Andorra's tax system is administered by the Departament de Tributs i de
				Fronteres under the Ministeri de Finances. Despite a customs union with
				the European Union, Andorra is not part of the EU VAT area and operates
				its own indirect tax regime.

				The Impost General Indirecte (IGI), introduced on 1 January 2013 by
				Llei 11/2012, replaced the previous sales tax system. IGI applies five
				rate tiers covering general commerce, essential goods, cultural and
				transport services, financial services, and exempt activities. Rates
				are defined in the tax categories below.

				All taxpayers (resident or non-resident, natural or legal) are
				identified by a Número de Registre Tributari (NRT). The NRT consists
				of one letter indicating taxpayer type (F for resident individuals,
				E for non-resident individuals, A for joint-stock companies, L for
				limited liability companies, and others for further entity types),
				followed by six digits and a check letter. The check-letter algorithm
				is not publicly published by the tax authority, so validation here is
				format-only; semantic verification requires the official portal.

				Andorra uses the euro as its official currency under a monetary
				agreement with the European Union.
			`),
		},
		Identities: identityDefinitions,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
		Normalizer: Normalize,
		Categories: taxCategories(),
	}
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	case *org.Identity:
		normalizeIdentity(obj)
	}
}
