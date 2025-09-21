# Australia (AU) Tax Regime

This document provides an overview of the tax regime in Australia.

## Goods and Services Tax (GST)

Australia operates a single-rate GST system with two main categories:

- **Standard Rate (10%)**: Applies to most goods and services in Australia since July 1, 2000.
- **GST-Free (0%)**: Applies to essential goods and services, exports, and specific exempt categories.

Unlike many countries with multiple VAT rates, Australia maintains a simplified single standard rate structure.

## GST Registration Requirements

Businesses in Australia must register for GST based on their annual turnover:

- **Mandatory Registration**: Businesses must register for GST if their annual turnover is **AUD $75,000** or more.
- **Non-Profit Organizations**: Must register if their annual turnover is **AUD $150,000** or more.
- **Voluntary Registration**: Businesses below these thresholds may choose to register voluntarily.

Registered businesses receive an Australian Business Number (ABN) which serves as their tax identification.

For more information, visit the [Australian Taxation Office website](https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst).

## GST-Free Supplies

GST-free supplies include a wide range of essential goods and services:

- Most basic food items
- Some education courses and materials
- Medical, health and care services
- Menstrual products
- Medical aids and medicines
- Some childcare and religious services
- Water, sewerage and drainage services
- Precious metals
- Exports of goods and services
- Sales of businesses as going concerns
- Cars for people with disabilities (when requirements are met)
- Farmland
- International transport services
- Eligible emissions units

For a comprehensive list, see the [ATO GST-free sales guide](https://www.ato.gov.au/businesses-and-organisations/gst-excise-and-indirect-taxes/gst/when-to-charge-gst-and-when-not-to/gst-free-sales).

## Australian Business Number (ABN)

The ABN is an 11-digit unique identifier issued to all entities registered in the Australian Business Register (ABR).

### ABN Format and Validation

The ABN uses a modulus 89 checksum algorithm:

1. Subtract 1 from the first (leftmost) digit
2. Multiply each digit by its position weighting factor: [10, 1, 3, 5, 7, 9, 11, 13, 15, 17, 19]
3. Sum all the products
4. Divide by 89 - if the remainder is zero, the ABN is valid

**Example**: ABN 51 824 753 556
- Modified: 41 824 753 556
- Calculation: (4×10)+(1×1)+(8×3)+(2×5)+(4×7)+(7×9)+(5×11)+(3×13)+(5×15)+(5×17)+(6×19) = 534
- Validation: 534 ÷ 89 = 6 remainder 0 ✓

For more details, refer to the [ABN validation guide](https://abr.business.gov.au/Help/AbnFormat).

### ABN Verification

ABNs can be verified through the [ABN Lookup service](https://abr.business.gov.au/) provided by the Australian Business Register.
