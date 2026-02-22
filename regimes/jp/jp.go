// Package jp provides the Japanese tax regime for the Qualified Invoice System (適格請求書等保存方式, QIS) effective since
// October 2023.
//
// Japan levies a consumption tax (消費税, JCT) modeled in GOBL under the standard VAT category. The regime supports four
// rate tiers: standard (10%), reduced (8%), zero (exports), and exempt.
// It validates Qualified Invoice Issuer Registration Numbers (T-numbers), Corporate Numbers (法人番号), and enforces
// invoice-level rules mandated by the National Tax Agency (NTA).
//
// Key references:
//   - NTA Invoice System overview: https://www.nta.go.jp/taxes/shiraberu/zeimokubetsu/shohi/keigenzeiritsu/invoice.htm
//   - NTA Consumption Tax: https://www.nta.go.jp/english/taxes/consumption_tax/
//   - Japan Customs export exemption: https://www.customs.go.jp/english/c-answer_e/extsukan/5003_e.htm
//   - NTA Qualified Invoice Issuer registration: https://www.invoice-kohyo.nta.go.jp
//   - NTA Corporate Number publication: https://www.houjin-bangou.nta.go.jp/en/
package jp

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New returns a new Japanese tax regime definition covering JCT categories, QIS tags and scenarios, reduced-rate
// extensions, and credit-note corrections.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "JP", // l10n.TaxCountryCode(l10n.JP),
		Currency:  currency.JPY,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Japan",
			i18n.JA: "日本",
		},
		TimeZone:   "Asia/Tokyo",
		Validator:  Validate,
		Normalizer: Normalize,
		Extensions: extensions,
		Identities: orgIdentityDefs,
		Categories: taxCategories,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
		Tags:      []*tax.TagSet{invoiceTags()},
		Scenarios: []*tax.ScenarioSet{invoiceScenarios()},
	}
}

// Validate checks the document type and dispatches to the appropriate validator.
func Validate(doc any) error {
	switch d := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(d)
	case *tax.Identity:
		return validateTaxIdentity(d)
	case *org.Identity:
		return validateOrgIdentity(d)
	}
	return nil
}

// Normalize performs any regime-specific normalisation.
func Normalize(doc any) {
	switch d := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(d)
	case *org.Identity:
		normalizeOrgIdentity(d)
	}
}
