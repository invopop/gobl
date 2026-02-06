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

This regime implements synchronous local validation using the algorithm specified in the [Investment Income Reporting File Upload Specification](https://www.ird.govt.nz/-/media/project/ir/home/documents/digital-service-providers/iir-file-upload-specification/investment-income-reporting-file-upload-specification.pdf) (Inland Revenue), instead of relying on asynchronous calls to the government SOAP endpoint.

#### IRD Number Validation Algorithm

1. **Check the valid range**
   - If the IRD number is < 10,000,000 or > 200,000,000 then the number is invalid

2. **Form the eight digit base number**
   - Remove the trailing check digit
   - If the resulting number is seven digits long, pad to eight digits by adding a leading zero

3. **Calculate the check digit**
   - Assign weight factors to each of the 8 base digits (left to right): 3, 2, 7, 6, 5, 4, 3, 2
   - Sum the products of each digit and its weight factor
   - Divide the sum by 11
   - If remainder is 0, the calculated check digit is 0
   - If remainder is not 0, subtract the remainder from 11 to get the calculated check digit
   - If calculated check digit is 0-9, go to step 5
   - If calculated check digit is 10, continue with step 4

4. **Re-calculate the check digit** (only if step 3 yielded 10)
   - Assign secondary weight factors (left to right): 7, 4, 3, 2, 5, 2, 7, 6
   - Sum the products of each digit and its weight factor
   - Divide the sum by 11
   - If remainder is 0, the calculated check digit is 0
   - If remainder is not 0, subtract the remainder from 11
   - If calculated check digit is 10, the IRD number is invalid

5. **Compare the check digit**
   - Compare the calculated check digit to the last digit of the original IRD number
   - If they match, the IRD number is valid

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

NZBN is **not mandatory** for tax invoices. The IRD number remains the primary tax identifier. NZBN is required only for Peppol network participation (using the `0088:` prefix).

Businesses registered on the NZBN can also create [Organization Parts](https://www.nzbn.govt.nz/whats-an-nzbn/identifying-different-parts-of-your-business/) to identify divisions, branches, or departments. Each Organization Part receives its own unique 13-digit GLN.

## Invoice Requirements

New Zealand replaced the term "tax invoice" with **Taxable Supply Information (TSI)** on 1 April 2023. Required information can come from multiple sources rather than a single document.

### Required Information by Transaction Value

#### Supplies $200 or less (inclusive of GST)

| Information Type | Details Required |
|------------------|------------------|
| Seller's details | Name or trade name |
| Buyer's details | Not required |
| Date | Date of invoice, or time of supply if no invoice issued |
| Goods/services | Description of the goods or services |
| Payment | The consideration for the supply |

#### Supplies over $200 and up to $1,000 (inclusive of GST)

| Information Type | Details Required |
|------------------|------------------|
| Seller's details | Name or trade name, GST number |
| Buyer's details | Not required |
| Date | Date of invoice, or time of supply if no invoice issued |
| Goods/services | Description of the goods or services |
| Payment | Either: (1) GST-exclusive amount + GST amount + GST-inclusive amount; OR (2) GST-inclusive amount + statement that GST is included |

#### Supplies over $1,000 (inclusive of GST)

| Information Type | Details Required |
|------------------|------------------|
| Seller's details | Name or trade name, GST number |
| Buyer's details | Name + at least one identifier: address (physical/postal), phone, email, trading name, NZBN, or website URL |
| Date | Date of invoice, or time of supply if no invoice issued |
| Goods/services | Description of the goods or services |
| Payment | Either: (1) GST-exclusive amount + GST amount + GST-inclusive amount; OR (2) GST-inclusive amount + statement that GST is included |

#### Imported Goods and Services

| Information Type | Details Required |
|------------------|------------------|
| Seller's details | Name or trade name, Address |
| Buyer's details | Not required |
| Date | Date of invoice, or time of supply if no invoice issued |
| Goods/services | Description of the goods or services |
| Payment | The consideration + any salary/wages paid to employees of seller (or commonly owned group) + any interest incurred by seller (or commonly owned group) |

#### Secondhand Goods

| Information Type | Details Required |
|------------------|------------------|
| Seller's details | Name or trade name, Address |
| Buyer's details | Not required |
| Date | Date on which goods were supplied |
| Goods/services | Description + quantity or volume of the goods |
| Payment | The consideration for the supply |

Note: Not implemented as an invoice scenario. Secondhand goods purchases from unregistered sellers are a record-keeping requirement for the buyer to claim input tax credits, not a supplier-issued TSI.

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

Disclaimer: based on multiple third-party sources. No official source found yet.

| Date | Requirement |
|------|-------------|
| October 2019 | Peppol framework adopted |
| March 2022 | All central government agencies can receive e-invoices |
| May 2025 | B2G e-invoicing required |
| January 2026 | Agencies processing >2,000 invoices must send AND receive e-invoices |
| January 2027 | Large suppliers (>$33M revenue) must submit e-invoices to government |

B2B e-invoicing remains voluntary with no current mandates.

## Official Sources

### Identities and Tax identities

| Source | URL | Description |
|--------|-----|-------------|
| IRD Numbers Overview | https://www.ird.govt.nz/managing-my-tax/ird-numbers | Official IRD number information |
| IRD Check Digit Algorithm | https://www.ird.govt.nz/-/media/project/ir/home/documents/digital-service-providers/iir-file-upload-specification/investment-income-reporting-file-upload-specification.pdf | Official IRD validation spec (IIR File Upload Specification) |
| IRD Validation Service | https://www.ird.govt.nz/digital-service-providers/services-catalogue/customer-and-account/ird-number-validation | API-based validation (for reference) |
| NZBN Official | https://www.nzbn.govt.nz/whats-an-nzbn/about/ | Official NZBN information |
| GS1 Check Digit | https://www.gs1.org/services/check-digit-calculator | GS1 mod-10 algorithm |

### GST

| Source | URL | Description |
|--------|-----|-------------|
| Taxable Supply Information | https://www.ird.govt.nz/gst/tax-invoices-for-gst/how-tax-invoices-for-gst-work | Official TSI requirements |
| Charging GST | https://www.ird.govt.nz/gst/charging-gst | Official requirements |
| Exempt Supplies | https://www.ird.govt.nz/gst/charging-gst/exempt-supplies | Financial services, residential rent, fine metals |
| Zero-rated Supplies | https://www.ird.govt.nz/gst/charging-gst/zero-rated-supplies | Exports, land transactions, going concerns |
| GST for overseas business | https://www.ird.govt.nz/gst/gst-for-overseas-businesses | Overseas business |

### E-Invoicing / Peppol

| Source | URL | Description |
|--------|-----|-------------|
| NZ E-Invoicing Portal | https://www.einvoicing.govt.nz/peppol | Official government e-invoicing site |
| MBIE Peppol Authority | https://www.mbie.govt.nz | Ministry of Business, Innovation and Employment |
