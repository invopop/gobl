## Added

- `ae`: identity keys for the Tax Identification Number (`ae-tin`), trade licence (`ae-trade-license`), Emirates ID (`ae-emirates-id`) and Cabinet Decision (`ae-cabinet-decision`). TINs must be 10 digits starting with 1, with a business identifier that is not all zeros and a Luhn check digit, as the Federal Tax Authority registered them under ISO 6523 (`GOBL-AE-ORG-IDENTITY-01`). Emirates ID numbers must be 15 digits starting with 784 (`GOBL-AE-ORG-IDENTITY-02`).
- `ae`: VAT key `margin-scheme` for supplies under the profit margin scheme. Like `exempt`, it takes no percent, so the invoice shows no VAT amount. Shown in the new `invoice-ae-margin-scheme` example.

## Changed

- `ae`: **breaking**: the tax ID code (TRN) must start with a valid TIN: its first 10 digits start with 1, have a business identifier that is not all zeros, and end in a Luhn check digit. Made-up numbers such as `123456789012345` now fail with `GOBL-AE-TAX-IDENTITY-01`.
- `eu-en16931-v2017`: **breaking**: document-level charges must have taxes (BR-37), as discounts already must (BR-32). Charges without taxes now fail with `GOBL-EU-EN16931-BILL-CHARGE-02`.

## Fixed

- `eu-en16931-v2017`: the charge and line discount reason rules cite BR-38 and BR-42, not BR-36 (charge amount) and BR-41 (line allowance amount). Fault codes are unchanged.
