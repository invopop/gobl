# 🇸🇦 GOBL Saudi Arabia Tax Regime

Saudi Arabia tax regime for GOBL covering VAT, tax identification, and ZATCA e-invoicing requirements.

## VAT Rates

| Rate     | Since      | Percent | Reference                             |
|----------|------------|---------|---------------------------------------|
| Standard | 2020-07-01 | 15%     | Royal Order No. A/638                 |
| Standard | 2018-01-01 | 5%      | GCC VAT Framework Agreement           |

Zero-rated and exempt categories are supported via global VAT keys.

## Identity Types

### Seller (BT-29-1)

| Code | Name                           |
|------|--------------------------------|
| CRN  | Commercial Registration Number |
| MOM  | MOMRAH License                 |
| MLS  | MHRSD License                  |
| 700  | 700 Number (Unified Number)    |
| SAG  | MISA License                   |

### Buyer (BT-46-1)

| Code | Name                        |
|------|-----------------------------|
| TIN  | Tax Identification Number   |
| NAT  | National ID                 |
| IQA  | Iqama                       |
| PAS  | Passport                    |
| GCC  | GCC ID                      |
| OTH  | Other ID                    |

Seller codes (CRN, MOM, MLS, 700, SAG, OTH) are also valid for buyers per BR-KSA-14.

## Invoice Validation

- Supplier VAT registration number required on all invoices (BR-KSA-39)
- Supplier name required on all invoices (BR-06)
- Customer name required on standard invoices (BR-KSA-42)
- Customer identification (TaxID or org identity) required on standard invoices (BR-KSA-81)
- Simplified invoices (B2C) skip customer requirement

## ZATCA Integration

Phase 2 integration (Fatoora platform, XML signing, QR codes) is handled by [gobl.zatca](https://github.com/invopop/gobl.zatca).

## References

- [ZATCA E-Invoicing](https://zatca.gov.sa/en/E-Invoicing/Introduction/Pages/What-is-e-invoicing.aspx)
- [SA TIN Guide](https://lookuptax.com/docs/tax-identification-number/saudi-arabia-tax-id-guide)

Find example SA GOBL files in the [`examples`](../../examples/sa) (uncalculated documents) and [`examples/sa/out`](../../examples/sa/out) (calculated envelopes) subdirectories.
