# 🇷🇴 GOBL Romania Tax Regime

Romania uses the RO e-Factura platform for their mandatory e-invoicing system.

## Public Documentation

### Legal Framework

- [Romanian Fiscal Code (Codul Fiscal) - Title VII](https://legislatie.just.ro/Public/DetaliiDocument/171282)
- [Law 199/2020 - B2G E-Invoicing Mandate](https://legislatie.just.ro/Public/DetaliiDocumentAfis/229853)
- [Law 296/2023 - B2B E-Invoicing Mandate](https://legislatie.just.ro/Public/DetaliiDocumentAfis/275745)
- [Minister of Finance Order 1366/2021 - RO_CIUS](https://legislatie.just.ro/Public/DetaliiDocument/248303)

### Tax Authorities

- [ANAF (Romanian National Agency for Fiscal Administration)](https://www.anaf.ro/)
- [European Commission - Romania Country Sheet](https://ec.europa.eu/digital-building-blocks/sites/spaces/einvoicingCFS/pages/881983595/Romania)

## Romania-specific Requirements

### Tax IDs

Romanian businesses are identified using the **CUI** (Cod Unic de Înregistrare) or **CIF** (Cod de
Identificare Fiscală), which are the same identifier. The format is 2-10 digits, optionally prefixed
with "RO" for intra-community transactions.

GOBL validates the CUI/CIF using the official modulo 11 checksum algorithm with weights
`[7, 5, 3, 2, 1, 7, 5, 3, 2]`:

1. Take all digits except the last (check digit)
2. Pad with leading zeros to 9 digits
3. Multiply each digit (left-to-right) by the corresponding weight
4. Sum all products
5. Calculate: `(sum × 10) mod 11`
6. If the remainder is 10, the check digit is 0; otherwise, the check digit equals the remainder

Example valid CUI/CIF numbers:

- `18547290` (8 digits)
- `RO18547290` (with RO prefix for VAT purposes)
- `27` (2 digits)

For individuals, the **CNP** (Cod Numeric Personal) is used:

- 13 digits in format `SYYMMDDJJNNNC`
- `S` = sex and century (1-8 for residents, 9 for foreign residents)
- Validated with checksum using weights `[2, 7, 9, 1, 4, 6, 3, 5, 8, 2, 7, 9]`

### VAT Rates

Romania applies three VAT rates as defined in Law 227/2015 (updated 2025):

| Rate       | Percentage | Applicable Since | Description                                       |
| ---------- | ---------- | ---------------- | ------------------------------------------------- |
| `standard` | 21%        | August 1, 2025   | General goods and services                        |
| `reduced`  | 11%        | August 1, 2025   | Food, pharmaceuticals, books, hotels, restaurants |

The regime also includes historical rates for accurate tax calculations on older documents:

- **2017-2025 (July)**: Standard 19%, Reduced 9%, Super-reduced 5%
- **2016**: Standard 20%
- **2010-2015**: Standard 24%

### E-Invoicing Requirements

According to Romanian legislation, all invoices must be reported via the **RO e-Factura** platform:

- **B2G** (Business-to-Government): Mandatory since September 2020 (Law 199/2020)
- **B2B** (Business-to-Business): Mandatory since January 1, 2024 (Law 296/2023)
- **B2C** (Business-to-Consumer): Mandatory since January 1, 2025

Key requirements:

- Invoices must be submitted within **5 calendar days** of issuance
- Compliance with European Standard EN 16931
- Use of RO_CIUS (Romanian Core Invoice Usage Specification) per Order 1366/2021
- XML format (UBL 2.1 or Cross-Industry Invoice)

### Corrective Invoices

According to Romanian invoicing regulations (Fiscal Code Art. 330), corrections to previously issued
invoices can only be made through credit notes or debit notes. All correction types require a
reference to the original invoice.

GOBL supports the following correction types:

- `credit-note` - For partial or full refunds (negative amounts)
- `debit-note` - For additional charges (positive amounts)
- `corrective` - For general corrections replacing a previous document

When issuing corrections, you must reference the original invoice using the `preceding` field:

```json
{
  "type": "credit-note",
  "code": "CN-001",
  "preceding": [
    {
      "series": "INV",
      "code": "12345",
      "issue_date": "2024-01-15"
    }
  ]
}
```

### Validation Rules

#### Supplier Requirements

- **MUST** have a valid Romanian tax ID (CUI/CIF)
- Tax ID must pass checksum validation

#### Customer Requirements

- **SHOULD** have identification for B2B transactions
- B2C transactions (to individuals) may omit tax ID, but CNP is recommended if available

#### Correction Documents

- **MUST** include at least one preceding document reference
- Credit notes, debit notes, and corrective invoices all require the `preceding` field
