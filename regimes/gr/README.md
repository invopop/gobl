# 🇬🇷 GOBL Greece Tax Regime

Greece uses the myDATA and Peppol BIS Billing 3.0 formats for their e-invoicing/tax-reporting system.

Example GR GOBL files can be found in the [`examples`](./examples) (YAML uncalculated documents) and [`examples/out`](./examples/out) (JSON calculated envelopes) subdirectories.

## Public Documentation

* [myDATA API Documentation v1.0.7](https://www.aade.gr/sites/default/files/2023-10/myDATA%20API%20Documentation_v1.0.7_eng.pdf)
* [Greek Peppol BIS Billing 3.0](https://www.gsis.gr/sites/default/files/eInvoice/Instructions%20to%20B2G%20Suppliers%20and%20certified%20PEPPOL%20Providers%20for%20the%20Greek%20PEPPOL%20BIS-EN-%20v1.0.pdf)
* [VAT Rates](https://www.gov.gr/en/sdg/taxes/vat/general/basic-vat-rates)

## Greece specifics

### VAT categories

Greece has three VAT rates: standard, reduced and super-reduced. Each of these rates are reduced by a 30% on the islands of Leros, Lesbos, Kos, Samos and Chios. The tax authority identifies each rate with a specific VAT category.

In GOBL, the IAPR VAT category code must be set using the `iapr-gr-vat-category` extension of a line's tax to one of these codes:

| Code | Description            | GOBL Rate       |
| ---- | ---------------------- | --------------- |
| `1`  | VAT rate 24%           | `standard`      |
| `2`  | VAT rate 13%           | `reduced`       |
| `3`  | VAT rate 6%            | `super-reduced` |
| `4`  | VAT rate 17%           |                 |
| `5`  | VAT rate 9%            |                 |
| `6`  | VAT rate 4%            |                 |
| `7`  | Without VAT            | `exempt`        |
| `8`  | Records without VAT    |                 |

Please, note that GOBL will automatically set the proper `gr-iapr-vat-cat` code and tax percent automatically when the line tax uses any of the GOBL rates specified in the table above.

As shown in the table, the islands' reduced version of the rates don't have GOBL rates associated to them, and so the GOBL invoice must provide explicitly the VAT category and the percent in those cases. For example:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  //...
  "tax": {
    "tags": [
      "islands"
    ]
  },
  // ...
  "lines": [
    {
      "i": 1,
      "quantity": "20",
      "item": {
        "name": "Υπηρεσίες Ανάπτυξης",
        "price": "90.00",
      },
      "sum": "1800.00",
      "taxes": [
        {
          "cat": "VAT",
          "percent": "17%",
          "ext": {
            "gr-iapr-vat-cat": "4" // Standard rate in Greek islands
          }
        }
      ],
      "total": "1800.00"
    }
  ],
}
```

### VAT exemptions

Greece invoices can be exempt of VAT for different causes and the tax authority require a specific cause code to be provided.

In a GOBL invoice, the `rate` of a line's tax need to be set to `exempt`, and the `ext` map's `gr-iapr-exemption` property needs to be set to one of these codes:

| Code | Description                                             |
| ---- | ------------------------------------------------------- |
| `1`  | Without VAT - article 3 of the VAT code                 |
| `2`  | Without VAT - article 5 of the VAT code                 |
| `3`  | Without VAT - article 13 of the VAT code                |
| `4`  | Without VAT - article 14 of the VAT code                |
| `5`  | Without VAT - article 16 of the VAT code                |
| `6`  | Without VAT - article 19 of the VAT code                |
| `7`  | Without VAT - article 22 of the VAT code                |
| `8`  | Without VAT - article 24 of the VAT code                |
| `9`  | Without VAT - article 25 of the VAT code                |
| `10` | Without VAT - article 26 of the VAT code                |
| `11` | Without VAT - article 27 of the VAT code                |
| `12` | Without VAT - article 27 - Seagoing Vessels of the VAT code |
| `13` | Without VAT - article 27.1.γ - Seagoing Vessels of the VAT code |
| `14` | Without VAT - article 28 of the VAT code                |
| `15` | Without VAT - article 39 of the VAT code                |
| `16` | Without VAT - article 39a of the VAT code               |
| `17` | Without VAT - article 40 of the VAT code                |
| `18` | Without VAT - article 41 of the VAT code                |
| `19` | Without VAT - article 47 of the VAT code                |
| `20` | VAT included - article 43 of the VAT code               |
| `21` | VAT included - article 44 of the VAT code               |
| `22` | VAT included - article 45 of the VAT code               |
| `23` | VAT included - article 46 of the VAT code               |
| `24` | Without VAT - article 6 of the VAT code                 |
| `25` | Without VAT - ΠΟΛ.1029/1995                             |
| `26` | Without VAT - ΠΟΛ.1167/2015                             |
| `27` | Without VAT - Other VAT exceptions                      |
| `28` | Without VAT - Article 24 (b) (1) of the VAT Code (Tax Free) |
| `29` | Without VAT - Article 47b of the VAT Code (OSS non-EU scheme) |
| `30` | Without VAT - Article 47c of the VAT Code (OSS EU scheme) |
| `31` | Excluding VAT - Article 47d of the VAT Code (IOSS)      |

For example:

```js
"lines": [
  {
    "i": 1,
    "quantity": "20",
    "item": {
      "name": "Υπηρεσίες Ανάπτυξης",
      "price": "90.00",
    },
    "sum": "1800.00",
    "taxes": [
      {
        "cat": "VAT",
        "rate": "exempt",
        "ext": {
          "gr-iapr-exemption": "30"
        }
      }
    ],
    "total": "1800.00"
  }
]
```
