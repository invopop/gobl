# GOBL New Zealand Tax Regime

New Zealand tax regime for the [GOBL](https://github.com/invopop/gobl) library.

## Overview

This package provides support for New Zealand's Goods and Services Tax (GST) system, including tax identity validation (IRD numbers and NZBN), tax categories, and invoice requirements.

## Tax System Summary

| Attribute | Value |
|-----------|-------|
| Country Code | `NZ` |
| Currency | `NZD` |
| Timezone | `Pacific/Auckland` |
| Tax Authority | Inland Revenue (IR) |
| Primary Tax | GST (Goods and Services Tax) |
| Standard Rate | 15% |

## GST Rates

| Rate Key | Percentage | Description | Since |
|----------|------------|-------------|-------|
| `standard` | 15% | Standard rate for most goods and services | 1 October 2010 |
| `accommodation` | 9% | Long-term commercial accommodation (28+ consecutive days) | 1 April 2024 |
| `zero` | 0% | Exports, international services, certain land transactions | - |
| `exempt` | - | Financial services, residential rent, donated goods by non-profits | - |

### Historical Rates

| Period | Standard Rate |
|--------|---------------|
| 1 October 2010 - present | 15% |
| 1 July 1989 - 30 September 2010 | 12.5% |
| 1 October 1986 - 30 June 1989 | 10% |

## Tax Identities

### IRD Number

The IRD number is the primary tax identifier in New Zealand, issued by Inland Revenue to individuals, companies, trusts, and other entities. For GST-registered businesses, the IRD number also serves as their GST number.

| Attribute | Specification |
|-----------|---------------|
| Format | 8 or 9 digits |
| Display | `XXX-XXX-XXX` (with hyphens) |
| Valid Range | 10,000,000 to 150,000,000 |
| Check Digit | Modulo 11 algorithm |
| Identity Key | `nz-ird` |

#### Validation Algorithm

The IRD number validation uses a modulo-11 check digit algorithm with weighted positions. The algorithm is documented in the "Non-Resident Withholding Tax and Resident Withholding Tax Specification Document" issued by Inland Revenue (31 March 2016).

**Steps:**

1. **Range Check**: If the IRD number (as integer) is < 10,000,000 or > 150,000,000, it is invalid.

2. **Form Base Number**: Remove the trailing check digit. If 7 digits remain, pad to 8 digits with a leading zero.

3. **Primary Weight Calculation**:
   - Apply weights `[3, 2, 7, 6, 5, 4, 3, 2]` to each digit (left to right)
   - Sum the products
   - Calculate: `remainder = sum mod 11`
   - If remainder = 0, check digit = 0
   - Otherwise, check digit = 11 - remainder

4. **Secondary Weight Calculation** (only if primary check digit = 10):
   - Apply weights `[7, 4, 3, 2, 5, 2, 7, 6]` to each digit
   - Sum the products
   - Calculate: `remainder = sum mod 11`
   - If remainder = 0, check digit = 0
   - Otherwise, check digit = 11 - remainder
   - If check digit is still 10, the IRD number is invalid

5. **Verification**: Compare calculated check digit with the actual last digit.

### NZBN (New Zealand Business Number)

The NZBN is a globally unique 13-digit identifier based on the GS1 Global Location Number (GLN) standard. It is automatically assigned to companies registered with the Companies Office; other entities can apply voluntarily.

| Attribute | Specification |
|-----------|---------------|
| Format | 13 digits |
| Prefix | `94` (New Zealand GS1 prefix) |
| Structure | 94 + 10-digit reference + check digit |
| Standard | ISO/IEC 6523 compliant |
| Peppol Format | `0088:` + NZBN |
| Identity Key | `nz-nzbn` |

#### Validation Algorithm

The NZBN uses the standard GS1 Modulo-10 check digit algorithm:

1. Starting from the rightmost digit (excluding check digit), assign alternating weights of 3 and 1
2. Multiply each digit by its weight and sum all products
3. Calculate: `check digit = (10 - (sum mod 10)) mod 10`

#### NZBN Adoption Status

As of the latest government reports, NZBN adoption is voluntary for most business types:

- **Mandatory**: Companies registered with the Companies Office receive an NZBN automatically
- **Voluntary**: Sole traders, partnerships, and other entity types can opt-in
- **E-invoicing**: Required for Peppol network participation (using `0088:` prefix)

**Important**: NZBN is NOT mandatory for tax invoices. The IRD number remains the primary tax identifier.

## Invoice Requirements

New Zealand replaced the term "tax invoice" with **Taxable Supply Information (TSI)** on 1 April 2023. This provides flexibility—required information can come from multiple sources.

### Required Information by Transaction Value

#### Supplies ≤ $200 NZD
- Seller's name or trade name
- Date of invoice or supply
- Description of goods/services
- Total amount payable

#### Supplies > $200 up to $1,000 NZD
All of the above, plus:
- **Seller's GST number** (IRD number)
- GST breakdown: either
  - GST-exclusive amount + GST amount + GST-inclusive total, OR
  - GST-inclusive amount + statement "includes GST"

#### Supplies > $1,000 NZD
All of the above, plus:
- **Buyer's name**
- **Buyer's identifier** (at least one of):
  - Address
  - Phone number
  - Email
  - Trading name
  - NZBN
  - Website

### Record Retention

GST records must be retained for **7 years**.

## E-Invoicing (Peppol)

New Zealand adopted the Peppol framework in October 2019, with MBIE (Ministry of Business, Innovation and Employment) serving as the Peppol Authority.

### Format

- **Standard**: PINT A-NZ (Peppol International Invoice Template for Australia-New Zealand)
- **Based on**: Peppol BIS Billing 3.0 / UBL 2.1

### Participant Identifier

New Zealand entities are identified on the Peppol network using:
```
0088:{NZBN}
```

Example: `0088:9429041234567`

### Mandate Timeline

| Date | Requirement |
|------|-------------|
| October 2019 | Peppol framework adopted |
| March 2022 | All central government agencies can receive e-invoices |
| May 2025 | B2G e-invoicing required |
| January 2026 | Agencies processing >2,000 invoices must send AND receive e-invoices |
| January 2027 | Large suppliers (>$33M revenue) must submit e-invoices to government |

**Note**: B2B e-invoicing remains voluntary with no current mandates.

## Implementation Notes

### Reference Implementation

Since Australia (`au`) is not yet implemented in GOBL, this regime uses structural patterns from:
- Italy (`it`) - Similar GST/VAT structure with single primary tax category
- Spain (`es`) - Well-documented regime with comprehensive validation

### Key Concerns Addressed

Following the GOBL contribution guidelines (see [gobl.ksef](https://github.com/invopop/gobl.ksef)):

1. **Basic B2B invoices support** ✓
2. **Tax ID validation as per local rules** ✓ (IRD mod-11, NZBN mod-10)
3. **Support for simplified invoices** ✓ (supplies ≤ $1,000 NZD)
4. **Credit notes / corrective invoices** ✓
5. **Additional field validation** ✓ (GST number requirements by threshold)

### What This Package Does NOT Implement

- Peppol message conversion (handled by separate `gobl.ubl` or similar packages)
- Real-time IRD validation service integration (async API calls are out of scope)
- Historical rate lookups before October 2010

## Official Sources

### Tax Rates and GST System

| Source | URL | Description |
|--------|-----|-------------|
| Inland Revenue - GST | https://www.ird.govt.nz/gst | Official GST overview |
| GST Rates (Wikipedia) | https://en.wikipedia.org/wiki/Goods_and_Services_Tax_(New_Zealand) | Historical rate changes |
| Taxually GST Guide | https://www.taxually.com/manuals/new-zealand | Comprehensive GST guide |

### Zero-Rated and Exempt Supplies

| Source | URL | Description |
|--------|-----|-------------|
| IRD Zero-Rated Supplies | https://www.classic.ird.govt.nz/gst/additional-calcs/calc-spec-supplies/calc-zero/calc-zero.html | Official zero-rating rules |
| Tax Accountant NZ | https://taxaccountant.kiwi.nz/gst-zero-rated-supplies | Zero-rated supply categories |

### Tax Identity Validation

| Source | URL | Description |
|--------|-----|-------------|
| IRD Numbers Overview | https://www.ird.govt.nz/managing-my-tax/ird-numbers | Official IRD number information |
| IRD Validation Service | https://www.ird.govt.nz/digital-service-providers/services-catalogue/customer-and-account/ird-number-validation | API-based validation (for reference) |
| nz-ird-validator (npm) | https://github.com/jarden-digital/nz-ird-validator | Algorithm implementation citing IRD spec document |
| Spectrum.Ird (.NET) | https://github.com/twoteesbrett/Spectrum.Ird | Algorithm implementation with link to IRD spec |
| NZBN Official | https://www.nzbn.govt.nz/whats-an-nzbn/about/ | Official NZBN information |
| GS1 Check Digit | https://www.gs1.org/services/check-digit-calculator | GS1 mod-10 algorithm |

### Invoice Requirements

| Source | URL | Description |
|--------|-----|-------------|
| Taxable Supply Information | https://www.ird.govt.nz/gst/tax-invoices-for-gst/how-tax-invoices-for-gst-work | Official TSI requirements |

### E-Invoicing / Peppol

| Source | URL | Description |
|--------|-----|-------------|
| NZ E-Invoicing Portal | https://www.einvoicing.govt.nz/peppol | Official government e-invoicing site |
| MBIE Peppol Authority | https://www.mbie.govt.nz | Ministry of Business, Innovation and Employment |
| Global VAT Compliance | https://www.globalvatcompliance.com/globalvatnews/new-zealand-mandates-peppol-e-invoicing-for-fovernment-agencies-2026/ | Mandate timeline |

### Algorithm Specification Document

The IRD number validation algorithm is specified in:

> **"Non-Resident Withholding Tax and Resident Withholding Tax Specification Document"**
> Inland Revenue, 31 March 2016

This document is not publicly available online but is distributed to registered Digital Service Providers. The algorithm has been independently implemented and verified by multiple open-source projects (see links above).

To request the official specification document, contact: **sdlu@ird.govt.nz**

## File Structure

```
regimes/nz/
├── nz.go              # Main regime definition
├── tax_identity.go    # IRD and NZBN validation
├── tax_categories.go  # GST categories and rates
├── scenarios.go       # Invoice scenarios (if needed)
├── examples/
│   ├── invoice-nz-nz.yaml    # Sample B2B invoice
│   └── out/
│       └── invoice-nz-nz.json # Calculated envelope
└── README.md          # This file
```

## Usage Examples

### Basic Invoice

```yaml
$schema: "https://gobl.org/draft-0/bill/invoice"
currency: NZD
issue_date: "2024-01-15"
supplier:
  name: "Kiwi Services Ltd"
  tax_id:
    country: NZ
    code: "123-456-789"  # IRD number
  addresses:
    - locality: Auckland
      country: NZ
customer:
  name: "Wellington Enterprises"
  tax_id:
    country: NZ
    code: "987-654-321"
lines:
  - quantity: 10
    item:
      name: "Consulting Services"
      price: "150.00"
    taxes:
      - cat: GST
        rate: standard
```

### Zero-Rated Export

```yaml
lines:
  - quantity: 1
    item:
      name: "Software License (Export)"
      price: "5000.00"
    taxes:
      - cat: GST
        rate: zero
```

## Contributing

Contributions are welcome! Please ensure all changes include:

1. Unit tests for validation logic
2. Updates to this README if behavior changes
3. References to official sources for any new rules

## License

Released under the Apache License 2.0. See the main GOBL repository for full license terms.

---

*Last updated: February 2026*
*GOBL regime code: `nz`*
