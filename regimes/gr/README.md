# 🇬🇷 GOBL Greece Tax Regime

Greece uses the myDATA and Peppol BIS Billing 3.0 formats for their e-invoicing/tax-reporting system.

Example GR GOBL files can be found in the [`examples`](./examples) (YAML uncalculated documents) and [`examples/out`](./examples/out) (JSON calculated envelopes) subdirectories.

## Public Documentation

* [myDATA API Documentation v1.0.7](https://www.aade.gr/sites/default/files/2023-10/myDATA%20API%20Documentation_v1.0.7_eng.pdf)
* [Greek Peppol BIS Billing 3.0](https://www.gsis.gr/sites/default/files/eInvoice/Instructions%20to%20B2G%20Suppliers%20and%20certified%20PEPPOL%20Providers%20for%20the%20Greek%20PEPPOL%20BIS-EN-%20v1.0.pdf)
* [VAT Rates](https://www.gov.gr/en/sdg/taxes/vat/general/basic-vat-rates)

## Greece specifics

### Invoice Type

The Greek tax authority (IAPR) requires the invoice type to be specified as part of the invoice. In GOBL, this type can be set using the `gr-mydata-invoice-type` extension in the tax section.

Alternatively, GOBL will set the extension for you based on the type and the tax tags you set in your GOBL invoice. The table below shows how this mapping is done:

| Type   | Description                                     | GOBL Type     | GOBL Tags                  |
| ------ | ----------------------------------------------- | ------------- |----------------------------|
| `1.1`  | Sales Invoice                                   | `standard`    | `goods`                    |
| `1.2`  | Sales Invoice/Intra-community Supplies          | `standard`    | `goods`, `export`, `eu`    |
| `1.3`  | Sales Invoice/Third Country Supplies            | `standard`    | `goods`, `export`          |
| `1.4`  | Sales Invoice/Sale on Behalf of Third Parties   | `standard`    | `goods`, `self-billed`     |
| `2.1`  | Service Rendered Invoice                        | `standard`    | `services`                 |
| `2.2`  | Intra-community Service Rendered Invoice        | `standard`    | `services`, `export`, `eu` |
| `2.3`  | Third Country Service Rendered Invoice          | `standard`    | `services`, `export`       |
| `5.1`  | Credit Invoice/Associated                       | `credit-note` |                            |
| `11.1` | Retail Sales Receipt                            | `standard`    | `goods`, `simplified`      |
| `11.2` | Service Rendered Receipt                        | `standard`    | `services`, `simplified`   |
| `11.3` | Simplified Invoice                              | `standard`    | `simplified`               |
| `11.4` | Retail Sales Credit Note                        | `credit-note` | `simplified`               |
| `11.5` | Retail Sales Receipt on Behalf of Third Parties | `credit-note` | `goods`, `simplified`, `self-billed` |

For example, this is how you set the IAPR invoice type explicitly:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  // ...
  "tax": {
    "ext": {
      "gr-mydata-invoice-type": "2.1"
    }
  }
}
```

And this is how you'll get the same result by using the GOBL type and tags:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
  // ...
  "type": "standard",
  "tax": {
    "tags": ["services"]
  }
}
```

### VAT categories

Greece has three VAT rates: standard, reduced and super-reduced. Each of these rates are reduced by a 30% on the islands of Leros, Lesbos, Kos, Samos and Chios. The tax authority identifies each rate with a specific VAT category.

In GOBL, the IAPR VAT category code must be set using the `gr-mydata-vat-cat` extension of a line's tax to one of the codes:

| Code | Description                 | GOBL Rate              |
| ---- | --------------------------- | ---------------------- |
| `1`  | Standard rate               | `standard`             |
| `2`  | Reduced rate                | `reduced`              |
| `3`  | Super-reduced rate          | `super-reduced`        |
| `4`  | Standard rate (Island)      | `standard+island`      |
| `5`  | Reduced rate (Island)       | `reduced+island`       |
| `6`  | Super-reduced rate (Island) | `super-reduced+island` |
| `7`  | Without VAT                 | `exempt`               |
| `8`  | Records without VAT         |                        |


Please, note that GOBL will automatically set the proper `gr-mydata-vat-cat` code and tax percent automatically when the line tax uses any of the GOBL rates specified in the table above. For example:

```js
{
  "$schema": "https://gobl.org/draft-0/bill/invoice",
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
          "rate": "standard+island"
        }
      ],
      "total": "1800.00"
    }
  ],
}
```

### VAT exemptions

Greece invoices can be exempt of VAT for different causes and the tax authority require a specific cause code to be provided.

In a GOBL invoice, the `rate` of a line's tax need to be set to `exempt`, and the `ext` map's `gr-mydata-exemption` property needs to be set to one of these codes:

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
          "gr-mydata-exemption": "30"
        }
      }
    ],
    "total": "1800.00"
  }
]
```

### Payment Methods

The IAPR requires invoices to specify a payment method code. In a GOBL invoice, the payment means is set using the `key` field in the payment instructions. The following table lists all the IAPR payment methods and how GOBL will map from the payment instructions key to each of them:

| Code | Name                             | GOBL Payment Instruction Key |
| ---- | -------------------------------- | ---------------------------- |
| `1`  | Domestic Payments Account Number | `credit-transfer`            |
| `2`  | Foreign Payments Account Number  | `credit-transfer+foreign`    |
| `3`  | Cash                             | `cash`                       |
| `4`  | Check                            | `cheque`                     |
| `5`  | On credit                        | `promissory-note`            |
| `6`  | Web Banking                      | `online`                     |
| `7`  | POS / e-POS                      | `card`                       |

For example:

```js
"payment": {
  "instructions": {
    "key": "credit-transfer+foreign" // Will set the IAPR Payment Method to "2"
  }
}
```

### Income Classification

Invoices reported to the Greek tax authority via myDATA can optionally include information about the income classification of each invoice item.

In a GOBL invoice, the `gr-mydata-income-cat` and `gr-mydata-income-type` extensions can be set at line tax level to any of the values expected by the IAPR:

#### Income Category

| Value          | Description                                            |
| -------------- | ------------------------------------------------------ |
| `category1_1`  | Commodity Sale Income (+)/(-)                          |
| `category1_2`  | Product Sale Income (+)/(-)                            |
| `category1_3`  | Provision of Services Income (+)/(-)                   |
| `category1_4`  | Sale of Fixed Assets Income (+)/(-)                    |
| `category1_5`  | Other Income/Profits (+)/(-)                           |
| `category1_6`  | Self-Deliveries/Self-Supplies (+)/(-)                  |
| `category1_7`  | Income on behalf of Third Parties (+)/(-)              |
| `category1_8`  | Past fiscal years income (+)/(-)                       |
| `category1_9`  | Future fiscal years income (+)/(-)                     |
| `category1_10` | Other Income Adjustment/Regularisation Entries (+)/(-) |
| `category1_95` | Other Income-related Information (+)/(-)               |

#### Income Type

| Value        | Description |
| -------------| ----------- |
| `E3_106`     | Self-Production of Fixed Assets – Self-Deliveries – Destroying in`vntor`y/Com modities |
| `E3_205`     | Self-Production of Fixed Assets – Self-Deliveries – Destroying in`vntor`y/Raw  and other materials |
| `E3_210`     | Self-Production of Fixed Assets – Self-Deliveries – Destroying in`vntor`y/Pro ducts and production in progress |
| `E3_305`     | Self-Production of Fixed Assets – Self-Deliveries – Destroying in`vntor`y/Raw  and other materials |
| `E3_310`     | Self-Production of Fixed Assets – Self-Deliveries – Destroying in`vntor`y/Pro ducts and production in progress |
| `E3_318`     | Self-Production of Fixed Assets – Self-Deliveries – Destroying in`vntory/Pro duction expenses |
| `E3_561_001` | Wholesale Sales of Goods and Services – for Traders |
| `E3_561_002` | Wholesale Sales of Goods and Services pursuant to article 39a pa`ragraph 5 `of the VAT Code (Law 2859/2000) |
| `E3_561_003` | Retail Sales of Goods and Services – Private Clientele |
| `E3_561_004` | Retail Sales of Goods and Services pursuant to article 39a pa`ragraph 5 `of the VAT Code (Law 2859/2000) |
| `E3_561_005` | Intra-Community Foreign Sales of Goods and Services |
| `E3_561_006` | Third Country Foreign Sales of Goods and Services |
| `E3_561_007` | Other Sales of Goods and Services |
| `E3_562`     | Other Ordinary Income |
| `E3_563`     | Credit Interest and Related Income |
| `E3_564`     | Credit Exchange Differences |
| `E3_565`     | Income from Participations |
| `E3_566`     | Profits from Disposing Non-Current Assets |
| `E3_567`     | Profits from the Reversal of Provisions and Impairments |
| `E3_568`     | Profits from Measurement at Fair Value |
| `E3_570`     | Extraordinary income and profits |
| `E3_595`     | Self-Production Expenses |
| `E3_596`     | Subsidies - Grants |
| `E3_597`     | Subsidies – Grants for Investment Purposes – Expense Coverage |
| `E3_880_001` | Wholesale Sales of Fixed Assets |
| `E3_880_002` | Retail Sales of Fixed Assets |
| `E3_880_003` | Intra-Community Foreign Sales of Fixed Assets |
| `E3_880_004` | Third Country Foreign Sales of Fixed Assets |
| `E3_881_001` | Wholesale Sales on behalf of Third Parties |
| `E3_881_002` | Retail Sales on behalf of Third Parties |
| `E3_881_003` | Intra-Community Foreign Sales on behalf of Third Parties |
| `E3_881_004` | Third Country Foreign Sales on behalf of Third Parties |
| `E3_598_001` | Sales of goods belonging to excise duty |
| `E3_598_003` | Sales on behalf of farmers through an agricultural cooperative e.t.c. |

For example:

```js
"lines": [
  {
    "i": 1,
    "quantity": "20",
    "item": {
      "name": "Υπηρεσίες Ανάπτυξης",
      "price": "90.00"
    },
    // ...
    "taxes": [
      {
        "cat": "VAT",
        "rate": "standard",
        "ext": {
          "gr-mydata-income-cat": "category1_1",
          "gr-mydata-income-type": "E3_106",
          // ...
        }
      }
    ]
  }
]
```
