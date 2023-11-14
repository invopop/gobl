# pl GOBL Polish Tax Regime

Poland uses the FA_VAT format for their e-invoicing system.

Example PL GOBL files can be found in the [`examples`](./examples) (YAML uncalculated documents) and [`examples/out`](./examples/out) (JSON calculated envelopes) subdirectories.

## Table of contents

* [Public Documentation](#public-documentation)
* [Zones](#zones)
* [Local Codes](#local-codes)
* [Complements](#complements)

## Public Documentation

- [Wzór faktury)](http://crd.gov.pl/wzor/2021/11/29/11089/)

### `TFormaPlatnosci` - Payment Means

The FA_VAT `TFormaPlatnosci` field specifies an invoice's means of payment. The following table lists all the supported values and how GOBL will map them from the invoice's payment instructions key:

| Code | Name                                | GOBL Payment Instructions Key |
| ---- | ----------------------------------- | ----------------------------- |
| 1    | Gotówka                             | `cash`                        |
| 2    | Karta                               | `card`                        |
| 3    | Bon                                 | `coupon`                      |
| 4    | Czek                                | `cheque`                      |
| 5    | Kredyt                              | `loan`                        |
| 6    | Przelew                             | `credit-transfer`             |
| 7    | Mobilna                             | `mobile`                      |

#### Example

The following GOBL maps to the `1` (gotówka) value of the `TFormaPlatnosci` field:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",

  // [...]

  "payment": {
    "instructions": {
      "key": "cash"
    }
  }
}
```

#### Document Type (TRodzajFaktury)

All Polish invoices must be identified with a specific type code defined by the FA_VAT format. The following table helps identify how GOBL will map the expected Polish code with a combination of the Invoice Type and tax tags.

| Code    | Type          | Tax Tags                            | Description                                           |
| ------- | ------------- | ----------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| VAT     | `standard`    |                                     | Regular invoice                                       |
| UPR     | `standard`    | `simplified`                        | Simplified (no customer)                              |
| ZAL     | `standard`    | `partial`                           | Advance invioce                                       |
| ROZ     | `standard`    | `settlement`                        | Settlement invoice                                    |
| KOR     | `corrective`  |                                     | Corrective (regular)                                  |
| KOR_ZAL | `corrective`  | `partial`                           | Corrective (advance)                                  |
| KOR_ROZ | `corrective`  | `settlement`                        | Corrective (settlement)                               |
