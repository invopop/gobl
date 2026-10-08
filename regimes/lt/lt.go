// Package lt provides tax regime support for Lithuania.
package lt

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the ISO 3166-1 alpha-2 code for Lithuania.
const CountryCode = "LT"

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new Lithuania regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   CountryCode,
		Currency:  currency.EUR,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Lithuania",
			i18n.LT: "Lietuva",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Lithuania's tax system is administered by the State Tax Inspectorate under the
				Ministry of Finance (Valstybinė mokesčių inspekcija, VMI). VAT is governed by the
				Law on Value Added Tax No. IX-751.

				VAT (pridėtinės vertės mokestis, PVM) applies at a standard rate, with reduced
				rates for certain goods and services such as accommodation, passenger transport,
				cultural events, medicines, and books.

				Businesses registered for VAT are identified by their VAT payer code (PVM
				mokėtojo kodas), formed by the prefix LT followed by 9 or 12 digits
				(e.g. LT100538411).
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("e-TAR - Pridėtinės vertės mokesčio įstatymas Nr. IX-751"),
				URL:   "https://www.e-tar.lt/portal/lt/legalAct/TAR.ED68997709F5/asr",
				At:    cal.NewDateTime(2026, 10, 8, 0, 0, 0),
			},
			{
				Title: i18n.NewString("VMI - Standartinis 21 proc. PVM tarifas (19 str.)"),
				URL:   "https://www.vmi.lt/evmi/standartinis-21-proc.-pvm-tarifas-19-str.-",
				At:    cal.NewDateTime(2026, 10, 8, 0, 0, 0),
			},
			{
				Title: i18n.NewString("VMI - Lengvatinis 12 proc. PVM tarifas (19 str.)"),
				URL:   "https://www.vmi.lt/evmi/lengvatinis-12-proc.-pvm-tarifas-19-str.-",
				At:    cal.NewDateTime(2026, 10, 8, 0, 0, 0),
			},
			{
				Title: i18n.NewString("VMI - Lengvatinis 5 proc. PVM tarifas (19 str.)"),
				URL:   "https://www.vmi.lt/evmi/lengvatinis-5-proc.-pvm-tarifas-19-str.-",
				At:    cal.NewDateTime(2026, 10, 8, 0, 0, 0),
			},
			{
				// Former reduced rate, in force until 2025-12-31 and not backfilled
				// (see tax_categories.go).
				Title: i18n.NewString("VMI - Lengvatinis 9 proc. PVM tarifas (19 str.)"),
				URL:   "https://www.vmi.lt/evmi/lengvatinis-9-proc.-pvm-tarifas-19-str.-",
				At:    cal.NewDateTime(2026, 10, 8, 0, 0, 0),
			},
			{
				Title: i18n.NewString("Ministry of Finance of the Republic of Lithuania - Value Added Tax"),
				URL:   "https://finmin.lrv.lt/en/competence-areas/taxation/main-taxes/value-added-tax/",
				At:    cal.NewDateTime(2026, 10, 8, 0, 0, 0),
			},
		},
		TimeZone:   "Europe/Vilnius",
		Categories: taxCategories,
		Scenarios:  []*tax.ScenarioSet{bill.InvoiceScenarios()},
	}
}
