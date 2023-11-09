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
| 1    | gotówka                             | `cash`                        |
| 2    | karta                               | `card`                        |
| 3    | bon                                 | `credit-transfer`             |
| 4    | czek                                | `cheque`                      |
| 5    | kredyt                              | `loan`                        |
| 6    | przelew                             | `credit-transfer`             |
| 7    | mobilna                             | `online`                      |

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
