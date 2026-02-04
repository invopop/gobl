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

## GST Rates

| Rate Key | Percentage | Description | Since |
|----------|------------|-------------|-------|
| `standard` | 15% | Standard rate for most goods and services | 1 October 2010 |
| `accommodation` | 9% | Long-term commercial accommodation (28+ consecutive days) | 1 April 2024 |
| `zero` | 0% | Exports, international services, certain land transactions | - |
| `exempt` | - | Financial services, residential rent, donated goods by non-profits | - |

### Historical Rates

GST records must be retained for 7 years, so rate changes before that period are generally not relevant for current record keeping.

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
| Valid Range | 10,000,000 to 200,000,000 |
| Identity Key | `nz-ird` |

#### Validation Approach

Inland Revenue provides a [SOAP-based IRD Number Validation Service](https://www.ird.govt.nz/digital-service-providers/services-catalogue/customer-and-account/ird-number-validation) for asynchronous validation. However, this requires network calls and API registration.

The validation algorithm is specified in the [Investment Income Reporting File Upload Specification](https://www.ird.govt.nz/-/media/project/ir/home/documents/digital-service-providers/iir-file-upload-specification/investment-income-reporting-file-upload-specification.pdf) (Inland Revenue).

This regime implements synchronous local validation using this algorithm, instead of relying on asynchronous calls to the government SOAP endpoint. Check out IRD Validation Service source.

## Organization Identities

### NZBN (New Zealand Business Number)

The NZBN is a globally unique 13-digit identifier based on the GS1 Global Location Number (GLN) standard. It is automatically assigned to companies registered with the Companies Office; other entities can apply voluntarily.

| Attribute | Specification |
|-----------|---------------|
| Format | 13 digits |
| Prefix | `94` (New Zealand GS1 prefix) |
| Structure | 94 + 10-digit reference + check digit |
| Standard | GS1 GLN (ISO/IEC 6523) |
| Peppol Format | `0088:` + NZBN |
| Identity Key | `gln` |

NZBN is validated using the standard GS1 check digit algorithm provided by the `pkg/gs1` package.

NZBN is **not mandatory** for tax invoices. The IRD number remains the primary tax identifier. NZBN is required only for Peppol network participation (using the `0088:` prefix).

## Invoice Requirements

New Zealand replaced the term "tax invoice" with **Taxable Supply Information (TSI)** on 1 April 2023. Required information can come from multiple sources rather than a single document.

### Required Information by Transaction Value

#### Supplies up to $200 NZD (inclusive of GST)
- Seller's name or trade name
- Date of invoice or supply
- Description of goods/services
- Total amount payable

#### Supplies over $200 up to $1,000 NZD (inclusive of GST)
All of the above, plus:
- **Seller's GST number** (IRD number)
- GST breakdown: either
  - GST-exclusive amount + GST amount + GST-inclusive total, OR
  - GST-inclusive amount + statement "includes GST"

#### Supplies over $1,000 NZD (inclusive of GST)
All of the above, plus:
- **Buyer's name**
- **Buyer's identifier** (at least one of):
  - Address
  - Phone number
  - Email
  - Trading name
  - NZBN
  - Website

#### Exported Goods and Services
For zero-rated exports, the taxable supply information must also include:
- The quantity or volume of goods or services supplied
- The buyer's name and address (for goods) or business details (for services)

#### Secondhand Goods
When purchasing secondhand goods from an unregistered seller, the buyer must record:
- Name and address of the supplier
- Date of purchase
- Description of goods
- Quantity or volume
- Price paid

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

B2B e-invoicing remains voluntary with no current mandates.

## Official Sources

### Tax Identities

| Source | URL | Description |
|--------|-----|-------------|
| IRD Numbers Overview | https://www.ird.govt.nz/managing-my-tax/ird-numbers | Official IRD number information |
| IRD Check Digit Algorithm | https://www.ird.govt.nz/-/media/project/ir/home/documents/digital-service-providers/iir-file-upload-specification/investment-income-reporting-file-upload-specification.pdf | Official IRD validation spec (IIR File Upload Specification) |
| IRD Validation Service | https://www.ird.govt.nz/digital-service-providers/services-catalogue/customer-and-account/ird-number-validation | API-based validation (for reference) |
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

## TODO

- [ ] Find official source for Peppol mandate timeline (current source is third-party)
- [ ] Study Peppol / PINT A-NZ requirements in detail
- [ ] Replace TSI text descriptions with screenshots from official documentation
- [ ] NZBN Organization part
