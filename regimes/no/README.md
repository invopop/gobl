# GOBL Norway Tax Regime

Norway uses VAT (Merverdiavgift, MVA) as its primary indirect tax system. This regime implements Norwegian invoicing requirements based on the Merverdiavgiftsloven (VAT Act) and the Bokføringsloven / Bokføringsforskriften (Bookkeeping Act and Regulation).

Find example GOBL files in the [`examples`](../../examples/no) (uncalculated documents) and [`examples/out`](../../examples/no/out) (calculated envelopes) subdirectories.

## Public Documentation

- [Skatteetaten - VAT rates](https://www.skatteetaten.no/en/rates/value-added-tax/)
- [Lovdata - Merverdiavgiftsloven](https://lovdata.no/dokument/NL/lov/2009-06-19-58)
- [Brønnøysundregistrene - Organisation number](https://www.brreg.no/en/about-us-2/our-registers/about-the-central-coordinating-register-for-legal-entities-ccr/about-the-organisation-number/)
- [Altinn - Invoice requirements](https://info.altinn.no/en/start-and-run-business/accounts-and-auditing/accounting/invoices-sales-documentation/)

## VAT Rates

| Rate | Percent | Description |
|------|---------|-------------|
| General | 25% | Standard rate for most goods and services |
| Reduced | 15% | Food, beverages, water and wastewater services |
| Super-reduced | 12% | Passenger transport, accommodation, cinema, broadcasting |
| Special | 11.11% | Raw fish (wild marine resources via fiskesalgslag) |

## Identification Numbers

### Tax Identity (MVA number)

The Norwegian VAT number is the 9-digit organisation number (organisasjonsnummer) suffixed with "MVA". During normalization, the "NO" prefix and "MVA" suffix are stripped, leaving the raw 9-digit code.

Format: `NO 923 456 783 MVA` → normalized to `923456783`

Validation uses a mod-11 check digit algorithm with weights `[3, 2, 7, 6, 5, 4, 3, 2]`, as specified by Brønnøysundregistrene. The check digit is the only structural constraint; the leading digit is not restricted (the historic 8/9 series is an allocation convention, not a validation rule).

### Organization Identity (Organisasjonsnummer)

Available as org identity type `ON`. Uses the same 9-digit format and mod-11 validation as the tax identity.

## Tags

| Tag | Description |
|-----|-------------|
| `reverse-charge` | Reverse charge / Omvendt avgiftsplikt |
| `simplified` | Simplified invoice (e.g. cash sale) — customer not required |

## Invoice Validation

GOBL core already validates the universal fields (type, dates, supplier presence and name, customer name when a tax ID is present, line prices, totals). On top of that, this regime adds only the Norway-specific rules:

- **Supplier identification**: the supplier must carry a tax ID or an `ON` organisasjonsnummer identity (bokføringsforskriften § 5-1-2). A VAT registration is not required — businesses below the NOK 50,000 threshold are not VAT-registered but may still issue invoices.
- **Customer**: required on standard invoices (§ 5-1-2 first paragraph); relaxed for simplified invoices.
- **Preceding reference**: required on credit notes (§ 5-2-7).

Some genuine legal requirements are intentionally left to the future EHF / SAF-T addon because they depend on data GOBL core does not model (the supplier's legal form) or on currency conversion: the seller's head-office address and the word "Foretaksregisteret" (required only for AS/ASA and foreign branches, § 5-1-2 / foretaksregisterloven § 10-2), and stating the VAT amount in NOK on foreign-currency invoices (§ 5-1-1 nr. 6).

## Correction Types

Only credit notes (kreditnota) are supported as a correction type, requiring a preceding document reference. Norwegian bookkeeping law (bokføringsforskriften § 5-2-7) recognises no debit-note concept.

## Out of Scope

The following are not included in this regime and may be addressed in future work:

- **SAF-T Norway**: Standard Audit File for Tax — a reporting format, not a transaction-level invoice format.
- **EHF / Peppol**: Norwegian e-invoicing format (EHF 3.0, aligned with Peppol BIS Billing 3.0 / EN 16931). B2G e-invoicing via Peppol is already mandatory; B2B mandatory e-invoicing is planned from 1 January 2027 (sending), with electronic bookkeeping/receiving phased in afterwards.
