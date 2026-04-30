# 🇵🇱 GOBL Polish KSeF FA_VAT Addon

Poland uses the FA_VAT format for their e-invoicing system managed by KSeF.

**IMPORTANT**: This addon is currently a work in progress, expect changes.

## Public Documentation

- [XML XSD Source - FA(3)](https://github.com/CIRFMF/ksef-docs/blob/main/faktury/schemy/FA/schemat_FA(3)_v1-0E.xsd)
- [Example Invoices for FA(3)](https://www.gov.pl/attachment/937002fa-c6b5-477d-8b56-22105fa728c2)
- [Information sheet about FA(3) in English](https://www.gov.pl/attachment/52052b6e-b8a9-497e-a0a7-941ae77b8dc8)

Older versions - FA(1), FA(2):
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

The following GOBL maps to the `1` (gotówka = cash) value of the `TFormaPlatnosci` field:

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

### Document Type (TRodzajFaktury)

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

### Self-invoicing (Samofakturowanie)

The `pl-favat-self-billing` extension indicates that the invoice is a self-billing invoice. It's added when using `self-billed` tag in GOBL, e.g.:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  "$tags": ["self-billed"]
  // ...
}
```

### Reverse charge (odwrotne obciążenie)

The `pl-favat-reverse-charge` extension indicates that the invoice is a reverse charge invoice. It's added when using `reverse-charge` tag in GOBL, e.g.:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  "$tags": ["reverse-charge"]
  // ...
}
```

### Margin scheme (procedura marży)

The `pl-favat-margin-scheme` extension indicates that the invoice is subject to the margin scheme. Available values are:

| Value | Meaning |
| ----- | ------- |
| 2     | Travel agency |
| 3.1   | Used goods |
| 3.2   | Works of art |
| 3.3   | Antiques and collectibles |

In GOBL:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  "tax": {
    "ext": {
      "pl-favat-margin-scheme": "2"
    },
  },
  // ...
}
```

### Exemption (zwolnienie z podatku)

The `pl-favat-exemption` extension indicates that the invoice is an exempt invoice. It's **not** automatically added, must be set manually to a value depending on the type of reason for exemption, and the invoice must contain a note describing the legal basis of exemption. Available values are:

| Value | Meaning |
| ----- | ------- |
| A     | Act in Polish law |
| B     | Directive 2006/112/EC |
| C     | Other legal basis |

In GOBL:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  "tax": {
    "ext": {
      "pl-favat-exemption": "A"
    },
  },
  "notes": [
    {
      "key": "legal",
      "code": "A", // must match the code in tax.ext
      "src": "pl-favat-exemption",
      "text": "Art. 25a ust. 1 pkt 9 ustawy o VAT" // text describing the legal basis
    }
  ]
  // ...
}
```

## Not supported features

- Intra-EU sale of new vehicles (`P_22` and `P_42_5` fields)
