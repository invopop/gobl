// Package ro provides the Romanian tax regime for GOBL.
//
// # Regime Overview
//
// Romania uses a comprehensive VAT system with multiple rates and has implemented
// mandatory e-invoicing through the RO e-Factura platform for B2B (since 2024)
// and B2C transactions (since 2025).
//
// # Tax Identification
//
// Romanian businesses are identified using:
//   - CIF (Cod de Identificare Fiscală) - Tax Identification Code
//   - CUI (Cod Unic de Înregistrare) - Unique Registration Code
//
// Both terms refer to the same identifier and are used interchangeably.
//
// # VAT Rates
//
// Romania applies the following VAT rates (since August 1, 2025):
//   - Standard rate: 21%
//   - Reduced rate: 11% (for specific goods and services)
//   - Super-reduced rate: 5% (for certain essential goods, deprecated in 2025)
//
// # E-Invoicing Requirements
//
// According to Law 296/2023 and subsequent regulations:
//   - B2G: Mandatory since September 2020 (Law 199/2020)
//   - B2B: Mandatory since January 1, 2024
//   - B2C: Mandatory since January 1, 2025
//
// All invoices must be reported via the RO e-Factura platform within 5 calendar
// days of issuance.
//
// # References
//
//   - Romanian Tax Code (Codul Fiscal): https://static.anaf.ro/static/10/Anaf/legislatie/Cod_fiscal_norme_31072017.htm
//   - Law 227/2015: https://static.anaf.ro/static/10/Anaf/Prezentare_R/Law227_11042018.pdf
//   - Law 296/2023 (B2B mandate): https://legislatie.just.ro/Public/DetaliiDocumentAfis/275745
//   - Law 199/2020 (B2G mandate): https://legislatie.just.ro/Public/DetaliiDocumentAfis/229853
//   - Minister of Finance Order 1366/2021 (RO_CIUS): https://legislatie.just.ro/Public/DetaliiDocument/248303
//   - ANAF (Romanian Tax Authority): https://www.anaf.ro/
//   - European Commission eInvoicing Country Sheet: https://ec.europa.eu/digital-building-blocks/sites/spaces/einvoicingCFS/pages/881983595/2025+Romania+2025+eInvoicing+Country+Sheet
package ro

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New instantiates a new Romanian tax regime.
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   l10n.RO.Tax(),
		Currency:  currency.RON,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Romania",
			i18n.RO: "România",
		},
		TimeZone:    "Europe/Bucharest",
		Identities:  identityTypeDefinitions,
		Categories:  taxCategories,
		Corrections: correctionDefinitions(),
		Scenarios: []*tax.ScenarioSet{
			bill.InvoiceScenarios(),
		},
		Validator:  Validate,
		Normalizer: Normalize,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateBillInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateOrgIdentity(obj)
	}
	return nil
}

// Normalize will perform any regime specific normalization.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	case *org.Identity:
		normalizeOrgIdentity(obj)
	}
}
