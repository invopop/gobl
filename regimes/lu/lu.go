// Package lu provides the tax regime definition for Luxembourg.
package lu

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the tax country code for Luxembourg.
const CountryCode = "LU"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("lu", rules.GOBL.Add(CountryCode),
		taxIdentityRules(),
		orgIdentityRules(),
		billInvoiceRules(),
	)
	norm.Register(
		norm.When(tax.IdentityIn(CountryCode), norm.For(normalizeTaxIdentity)),
	)
	norm.RegisterWithGuard(is.InContext(tax.RegimeIn(CountryCode)),
		norm.For(normalizeOrgIdentity),
	)
}

// New instantiates a new Luxembourg tax regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Luxembourg",
			i18n.FR: "Luxembourg",
			i18n.LB: "Lëtzebuerg",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Luxembourg applies VAT (Taxe sur la Valeur Ajoutée, TVA) administered by
				the Administration de l'Enregistrement, des Domaines et de la TVA (AED).
				As an EU member state, Luxembourg follows the EU VAT Directive 2006/112/EC.

				Luxembourg has four non-zero VAT rates: a standard rate, an intermediate
				("parking") rate, a reduced rate, and a super-reduced rate. A temporary
				one-percentage-point reduction was applied across all rates throughout 2023
				as a cost-of-living measure, before the rates were restored to their
				pre-2023 levels on 1 January 2024.

				Businesses are identified by their TVA number (LU followed by 8 digits,
				where the last two digits are a mod-89 check code) and optionally by their
				company registration number (RCS number) from the Registre de Commerce et
				des Sociétés.
			`),
			i18n.FR: here.Doc(`
				Le Luxembourg applique la TVA (Taxe sur la Valeur Ajoutée) administrée par
				l'Administration de l'Enregistrement, des Domaines et de la TVA (AED).
				En tant qu'État membre de l'UE, le Luxembourg suit la Directive TVA
				2006/112/CE.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "AED – VAT rates",
					i18n.FR: "AED – Taux de TVA",
				},
				URL: "https://www.aed.public.lu/en/tva/taux-tva.html",
			},
		},
		TimeZone:   "Europe/Luxembourg",
		Identities: identityTypeDefinitions,
		Categories: taxCategories,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
		},
	}
}
