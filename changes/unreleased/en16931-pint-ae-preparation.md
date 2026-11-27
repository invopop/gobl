## Added

- `ae`: identity types for the Tax Identification Number (`TIN`), trade licence (`TL`), Emirates ID (`EID`) and Cabinet Decision (`CD`). TINs must be 10 digits starting with 1 (`GOBL-AE-ORG-IDENTITY-01`) and end in a Luhn check digit (`GOBL-AE-ORG-IDENTITY-02`), as the Federal Tax Authority registered them under ISO 6523. Emirates ID numbers must be 15 digits starting with 784 (`GOBL-AE-ORG-IDENTITY-03`). Shown in the new `invoice-ae-identities` example.
- `untdid`: document type code `480` ("Invoice out of scope of tax") to the `untdid-document-type` catalogue extension.

## Changed

- `ae`: **breaking**: the tax ID code (TRN) must be 15 digits starting with 1 (`GOBL-AE-TAX-IDENTITY-01`), its first 10 digits, the TIN, must end in a Luhn check digit (`GOBL-AE-TAX-IDENTITY-02`), and it must end with `03`, the VAT registration suffix (`GOBL-AE-TAX-IDENTITY-03`). Made-up numbers such as `123456789012345` now fail.

## Removed

- `ae`: the scenarios that added "Reverse Charge" and "Simplified Tax Invoice" notes. PINT AE carries these cases in structured fields, and neither text met the VAT Executive Regulation's wording requirements.
- `ae`: the `invoice-ae-simplified` example and the regime description's paragraph on simplified invoices, which do not apply to PINT AE e-invoices.
