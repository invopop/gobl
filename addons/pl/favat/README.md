# 🇵🇱 GOBL Polish KSeF FA_VAT Addon

Poland uses the FA_VAT format for their e-invoicing system managed by KSeF.

**IMPORTANT**: This addon is currently a work in progress, expect changes.

## Public Documentation

- [XML XSD Source - KSeF](https://www.podatki.gov.pl/e-deklaracje/dokumentacja-it/struktury-dokumentow-xml/#ksef)
- [Invoice Templates (Wzór faktury) FA(1)](http://crd.gov.pl/wzor/2021/11/29/11089/)
- [Invoice Templates (Wzór faktury) FA(2)](http://crd.gov.pl/wzor/2023/06/29/12648/)

## Poland-specific Requirements

### `TFormaPlatnosci` - Payment Means

The FA_VAT `TFormaPlatnosci` field specifies an invoice's means of payment. The following table lists all the supported values and how GOBL will map them from the invoice's payment instructions key:

| Code | Name    | GOBL Payment Instructions Key |
| ---- | ------- | ----------------------------- |
| 1    | Gotówka | `cash`                        |
| 2    | Karta   | `card`                        |
| 3    | Bon     | `coupon`                      |
| 4    | Czek    | `cheque`                      |
| 5    | Kredyt  | `loan`                        |
| 6    | Przelew | `credit-transfer`             |
| 7    | Mobilna | `mobile`                      |

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

| Code    | Type          | Tax Tags     | Description                        |
| ------- | ------------- | ------------ | ---------------------------------- |
| VAT     | `standard`    |              | Regular invoice                    |
| UPR     | `standard`    | `simplified` | Simplified (no customer)           |
| ZAL     | `standard`    | `partial`    | Advance invioce                    |
| ROZ     | `standard`    | `settlement` | Settlement invoice                 |
| KOR     | `credit-note` |              | Credit note for regular invoice    |
| KOR_ZAL | `credit-note` | `partial`    | Credit note for advance invoice    |
| KOR_ROZ | `credit-note` | `settlement` | Credit note for settlement invoice |
