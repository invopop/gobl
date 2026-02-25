# GOBL Norway Tax Regime

Norway uses VAT (Merverdiavgift, MVA) as its primary indirect tax system. This regime implements Norwegian invoicing requirements based on the Merverdiavgiftsloven (VAT Act) and Regnskapsloven (Accounting Act).

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

Validation uses a mod-11 check digit algorithm with weights `[3, 2, 7, 6, 5, 4, 3, 2]`, as specified by Brønnøysundregistrene. The first digit must be 8 or 9.

### Organization Identity (Organisasjonsnummer)

Available as org identity type `ON`. Uses the same 9-digit format and mod-11 validation as the tax identity.

## Tags

| Tag | Description |
|-----|-------------|
| `reverse-charge` | Reverse charge / Omvendt avgiftsplikt |
| `simplified` | Simplified invoice — customer not required, supplier address relaxed |

## Correction Types

Credit notes and debit notes are supported as correction types. Both require a preceding document reference.

## Out of Scope

The following are not included in this regime and may be addressed in future work:

- **SAF-T Norway**: Standard Audit File for Tax — a reporting format, not a transaction-level invoice format.
- **EHF / Peppol**: Norwegian e-invoicing format (EHF, based on Peppol BIS 3.0). B2B mandatory e-invoicing is proposed from 2028 but not yet law.
