## Changed

- `eu-en16931-v2017`: **breaking**: document-level charges must have taxes (BR-37), as discounts already must (BR-32). Charges without taxes now fail with `GOBL-EU-EN16931-BILL-CHARGE-02`. Informative charges fail too, as they can't have taxes and EN 16931 has no way to express them.
- `eu-en16931-v2017`: **breaking**: an invoice customer must have a name (BR-7), with or without a tax ID. Customers without a name now fail with `GOBL-EU-EN16931-BILL-INVOICE-09`.
- `eu-en16931-v2017`: **breaking**: every invoice line must have taxes, which give it the invoiced item VAT category (BR-CO-04). Lines without taxes now fail with `GOBL-EU-EN16931-BILL-INVOICE-10`.

## Fixed

- `eu-en16931-v2017`: the charge and line discount reason rules cite BR-38 and BR-42, not BR-36 (charge amount) and BR-41 (line allowance amount). Fault codes are unchanged.

## Removed

- `eu-en16931-v2017`: payment details and terms are no longer required when an amount is due (`GOBL-EU-EN16931-BILL-INVOICE-06` and `-07`), nor must payment terms have due dates or notes (`GOBL-EU-EN16931-PAY-TERMS-01`). These applied BR-CO-25, which the EN 16931 validation artefacts dropped in version 1.3.16.
