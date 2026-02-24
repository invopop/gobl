# GOBL Romania Tax Regime

### Tax Authority

- [ANAF (Agenția Națională de Administrare Fiscală)](https://www.anaf.ro/) - Romania's National Agency for Fiscal Administration
- [ANAF e-Factura Portal](https://www.anaf.ro/anaf/internet/ANAF/asistenta_contribuabili/e_factura) - Official e-invoicing system

### VAT Rates

- [OECD VAT Rate Database](https://www.oecd.org/tax/tax-policy/tax-database/) - Historical VAT rates
- [Avalara - Romania VAT rate changes 2025](https://www.avalara.com/blog/en/europe/2025/07/blog-romania-vat-rate-changes-2025.html) - Aug 2025 rate reform details
- [Law No. 141/2025](https://www.vatupdate.com/2025/07/31/president-signs-fiscal-law-vat-rises-to-21-reduced-rates-unified/) - VAT increase signed into law

Current VAT rates (effective 1 August 2025):

| Rate | Percentage | Applies to |
|------|-----------|------------|
| Standard | 21% | General goods and services |
| Reduced | 11% | Food, non-alcoholic beverages, hotel accommodation, restaurants, books, newspapers, medical products |
| Super-reduced | 11% | Social housing, cultural events |
| Exempt | — | Medical, educational, financial, and insurance services |
| Zero | 0% | Exports, intra-community supplies |

Previous standard rate was 19% (Jan 2017 - Jul 2025). Reduced was 9%, super-reduced was 5%.

**Limitation**: Historical VAT rates are only tracked back to January 2017. Invoices dated before 2017 (e.g. during the 20% standard rate period of 2016, or the 24% period of 2010-2015) are supported on a best-effort basis for the standard rate only. The reduced and super-reduced rates do not have entries before 2017.

VAT rates are the same across all organization types (SRL, SA, PFA, micro-enterprise) and depend only on the goods/services category. There is a VAT registration threshold (RON 395,000) for exemption, but that is a binary status, not a different rate.

### Tax Identity (CUI/CIF)

- [ONRC (Oficiul Național al Registrului Comerțului)](https://www.onrc.ro/index.php/en/) - National Trade Register
- [Romanian CUI/CIF Validation](https://github.com/mtarnovan/romanianvalidators) - Reference implementation for checksum algorithm

Romanian tax identity codes (CUI/CIF) are 2-10 digit numbers validated with a weighted checksum algorithm using weights `[7, 5, 3, 2, 1, 7, 5, 3, 2]`. The "RO" prefix used for intra-EU VAT identification is automatically stripped during normalization.

Personal identification numbers (CNP, 13 digits) are not tax identity codes and are not accepted. Sole traders (PFA) receive a CUI upon registration.

### e-Invoicing

Romania mandates B2B e-invoicing via the RO e-Factura system (ANAF) since January 2024, aligned with EU standard EN 16931 using UBL 2.1 format.

ANAF validation requires UBL 2.1 XML in the RO_CIUS format — GOBL JSON must first be converted (e.g. via `gobl.ubl`) before submission. For manual validation:

- **Online**: Upload XML at `https://www.anaf.ro/uploadxmi/`
- **API test endpoint**: `POST https://api.anaf.ro/test/FCTEL/rest/validare/FACT1` with `Content-Type: application/xml`
- **Reference Go library**: [`github.com/printesoi/e-factura-go`](https://github.com/printesoi/e-factura-go)

## Romania-specific Requirements

### Corrections

Romania supports both credit notes and debit notes for correcting invoices.

### Simplified Invoices

Simplified invoices (without a customer) are supported using the `simplified` tag.

### Reverse Charge

Reverse charge invoices use the `reverse-charge` tag and automatically include the legal note "Reverse Charge / Taxare inversă."

## Scope and Future Work

- **Self-billing**: Not handled in the base regime; addon territory.
- **e-Factura addon**: CIUS-RO / e-Factura addon for UBL 2.1 XML conversion and ANAF submission is planned as future work.
- **Organization identity (ONRC)**: Romania's trade register number (J-number) is addon territory, following the pattern of DE, FR, IT, SE.
- **CNP (personal ID)**: Not supported as a tax identity. PFAs receive a CUI upon registration.
