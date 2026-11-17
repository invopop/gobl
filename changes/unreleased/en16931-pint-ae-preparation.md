## Added

- `ae`: identity types for the Tax Identification Number (`TIN`), trade licence (`TL`), Emirates ID (`EID`) and Cabinet Decision (`CD`). TINs must be 10 digits starting with 1 (`GOBL-AE-ORG-IDENTITY-01`) and end in a Luhn check digit (`GOBL-AE-ORG-IDENTITY-02`), with a business identifier that is not all zeros (`GOBL-AE-ORG-IDENTITY-04`), as the Federal Tax Authority registered them under ISO 6523. Emirates ID numbers must be 15 digits starting with 784 (`GOBL-AE-ORG-IDENTITY-03`).
- `ae`: `margin-scheme` tag for supplies under the profit margin scheme, with a scenario that adds the statement the invoice must carry ("VAT charged under the Profit Margin Scheme") as a VAT note. Lines use the `exempt` key, so the invoice shows no VAT amount. Shown in the new `invoice-ae-margin-scheme` example.

## Changed

- `ae`: **breaking**: the tax ID code (TRN) must be 15 digits starting with 1 (`GOBL-AE-TAX-IDENTITY-01`), and its first 10 digits, the TIN, must end in a Luhn check digit (`GOBL-AE-TAX-IDENTITY-02`) and have a business identifier that is not all zeros (`GOBL-AE-TAX-IDENTITY-03`). Made-up numbers such as `123456789012345` now fail.
- `eu-en16931-v2017`: **breaking**: document-level charges must have taxes (BR-37), as discounts already must (BR-32). Charges without taxes now fail with `GOBL-EU-EN16931-BILL-CHARGE-02`. Informative charges are left out, as they can't have taxes.

## Fixed

- `eu-en16931-v2017`: the charge and line discount reason rules cite BR-38 and BR-42, not BR-36 (charge amount) and BR-41 (line allowance amount). Fault codes are unchanged.
