## Added

- `ae`: identity types for the Tax Identification Number (`TIN`), trade licence (`TL`), Emirates ID (`EID`) and Cabinet Decision (`CD`). TINs must be 10 digits starting with 1 (`GOBL-AE-ORG-IDENTITY-01`) and end in a Luhn check digit (`GOBL-AE-ORG-IDENTITY-02`), as the Federal Tax Authority registered them under ISO 6523. Emirates ID numbers must be 15 digits starting with 784 (`GOBL-AE-ORG-IDENTITY-03`).
- `ae`: `margin-scheme` tag for supplies under the profit margin scheme, with a scenario that adds the statement the invoice must carry ("VAT charged under the Profit Margin Scheme") as a VAT note. Lines use the `exempt` key, so the invoice shows no VAT amount. Shown in the new `invoice-ae-margin-scheme` example.

## Changed

- `ae`: **breaking**: the tax ID code (TRN) must be 15 digits starting with 1 (`GOBL-AE-TAX-IDENTITY-01`), and its first 10 digits, the TIN, must end in a Luhn check digit (`GOBL-AE-TAX-IDENTITY-02`). Made-up numbers such as `123456789012345` now fail.
