# 🇳🇴 GOBL Norway Tax Regime

This document provides an overview of the Norwegian tax regime as implemented in GOBL.

Find example NO GOBL files in the [`examples`](../../examples/no) (uncalculated documents) and [`examples/out`](../../examples/no/out) (calculated envelopes) subdirectories.

---

## Overview

Norway applies Value Added Tax (VAT), locally known as **Merverdiavgift (MVA)**.

- Country Code: `NO`
- Currency: `NOK`
- Time Zone: `Europe/Oslo`
- Tax Category: `VAT`

Norway is not part of the EU VAT system but operates its own VAT legislation administered by the Norwegian Tax Administration (Skatteetaten).

---

## VAT (Merverdiavgift – MVA)

### Standard VAT Rates

| Rate Type                      | GOBL Rate Key     | Percentage |
|--------------------------------|-------------------|------------|
| Standard Rate                  | `standard`        | 25%        |
| Reduced Rate (Foodstuffs)      | `reduced`         | 15%        |
| Super Reduced (Transport etc.) | `super-reduced`   | 12%        |

These rates are defined in the `VAT` category and handled automatically through standard GOBL rate keys.

### Example

```json
{
  "taxes": [
    {
      "cat": "VAT",
      "rate": "standard"
    }
  ]
}
```

---

## VAT Exemptions

Certain supplies may be exempt or zero-rated under Norwegian VAT law.

In GOBL:

- Use `rate: "exempt"` for VAT-exempt supplies.
- Use `rate: "zero"` where applicable for zero-rated supplies.

### Example

```json
{
  "taxes": [
    {
      "cat": "VAT",
      "rate": "exempt"
    }
  ]
}
```

---

## Norwegian Tax Identities

### Organization Number (Organisasjonsnummer)

Norwegian businesses are identified using a **9-digit organization number**.

Example:

```
974760673
```

---

### VAT Number Representation

When registered for VAT, the organization number is commonly represented as:

```
NO974760673MVA
```

However, in GOBL:

- The canonical internal representation is always the **9-digit organization number**.
- The system automatically normalizes inputs like:
  - `974760673`
  - `974760673MVA`
  - `NO974760673MVA`
  - `974 760 673 mva`

All normalize internally to:

```
974760673
```

---

## Validation Rules

Norwegian organization numbers are validated using the official **MOD11 checksum algorithm**.

Validation includes:

- Exact 9-digit length
- Numeric normalization
- MOD11 checksum verification

Invalid length or checksum values will result in validation errors.

### Example

```json
{
  "tax_id": {
    "country": "NO",
    "code": "NO974760673MVA"
  }
}
```

Internally stored as:

```
974760673
```

---

## Organization Identity Types

The following identity type is defined in the Norwegian regime:

| Code | Description          |
|------|----------------------|
| `ON` | Organization Number  |

This identity is validated using the MOD11 checksum algorithm.

---

## Reverse Charge

In certain cases (such as cross-border services), VAT liability may be shifted to the buyer under the reverse charge mechanism.

In GOBL, this is represented using the `reverse-charge` tax tag.

### Example

```json
{
  "type": "standard",
  "tax": {
    "tags": ["reverse-charge"]
  },
  "lines": [
    {
      "i": 1,
      "quantity": "1",
      "item": {
        "name": "Consulting Services",
        "price": "1000.00"
      },
      "taxes": [
        {
          "cat": "VAT",
          "rate": "exempt"
        }
      ]
    }
  ]
}

## Invoicing Notes

Norway does not currently mandate a single national e-invoicing format comparable to:

- Italy’s FatturaPA
- Greece’s myDATA

However, electronic invoicing using **EHF (Elektronisk Handelsformat)** is widely adopted, especially for B2G (Business-to-Government) transactions, and also in many B2B(Business-to-Business) transaction.

GOBL invoices can be converted to local formats via integration layers where required.

---

## Resources

- [Norwegian Tax Administration (Skatteetaten)](https://www.skatteetaten.no/en/)
- [VAT Overview (MVA)](https://www.skatteetaten.no/en/business-and-organisation/vat-and-duties/vat/)
- [VAT Rates](https://www.skatteetaten.no/en/rates/value-added-tax/)
- [Altinn Business Portal](https://www.altinn.no/en/)
- [Norwegian VAT Act (Merverdiavgiftsloven)](https://lovdata.no/dokument/NL/lov/2009-06-19-58)

---