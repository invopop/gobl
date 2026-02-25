# 🇳🇴 Norway (`NO`) Tax Regime for GOBL

This document describes the Norwegian tax regime implementation in GOBL.

Example uncalculated documents can be found in [`examples/no`](../../examples/no),  
and calculated envelopes in [`examples/no/out`](../../examples/no/out).

---

## Overview

Norway applies **Value Added Tax (VAT)**, locally known as **Merverdiavgift (MVA)**.

- **Country Code:** `NO`
- **Currency:** `NOK`
- **Time Zone:** `Europe/Oslo`
- **Primary Tax Category:** `VAT`

Norway is not part of the EU VAT system but operates its own VAT legislation administered by the Norwegian Tax Administration (Skatteetaten).

This regime focuses on:

- VAT rate modeling
- Norwegian organization number validation
- Reverse charge handling
- VAT-exempt supplies
- Simplified invoices (B2C scenarios)

Historical VAT changes and sector-specific edge cases are out of scope.

---

# VAT (Merverdiavgift – MVA)

## VAT Rates

Norway currently applies three main VAT rates:

| Rate Type                                   | GOBL `rate` value | Percentage |
|---------------------------------------------|-------------------|------------|
| Standard Rate                               | `standard`        | 25%        |
| Reduced Rate (Foodstuffs)                   | `reduced`         | 15%        |
| Super Reduced (Transport, cinema, lodging)  | `super-reduced`   | 12%        |

These are implemented using:

- `CategoryVAT`
- Global VAT keys
- Regime-specific rate definitions

### Example (Standard VAT)

```yaml
taxes:
  - cat: VAT
    rate: standard
```

The system automatically resolves to:

```json
{
  "cat": "VAT",
  "key": "standard",
  "rate": "general",
  "percent": "25.0%"
}
```

**Note:**

- `rate: standard` is normalized internally to `rate: general`
- `key: standard` is applied automatically when no key is specified

---

# VAT Exempt & Reverse Charge

## Exempt Supplies

VAT-exempt transactions use the global VAT key:

```yaml
taxes:
  - cat: VAT
    key: exempt
```

Behavior:

- No percentage is applied
- Tax amount = `0.00`
- Totals reflect base only

---

## Reverse Charge

Reverse charge is represented at two levels:

1. **Document-level tag**
2. **Line-level VAT key**

### Example

```yaml
tags:
  - reverse-charge

lines:
  - taxes:
      - cat: VAT
        key: reverse-charge
```

Explanation:

- The `reverse-charge` tag enables invoice-level validation rules.
- The VAT key `reverse-charge` ensures tax is not calculated.
- Tax amount is zero, and the buyer accounts for VAT.

---

# Simplified Invoices

When the invoice includes:

```yaml
tags:
  - simplified
```

The regime allows omission of the customer section.

This models typical Norwegian B2C simplified receipts.

---

# Norwegian Organization Number (Organisasjonsnummer)

## Canonical Representation

Norwegian businesses are identified using a **9-digit organization number**.

Example:

```
974760673
```

If VAT registered, the number is often displayed externally as:

```
NO974760673MVA
```

However, within GOBL:

- The canonical representation is always the 9-digit number.
- Prefixes (`NO`) and suffixes (`MVA`) are removed.
- Spacing and case are normalized.

All of the following normalize to:

```
974760673
```

- `974760673`
- `974760673MVA`
- `NO974760673MVA`
- `974 760 673 mva`

---

## Validation

Organization numbers are validated using:

- Exact 9-digit length
- Numeric normalization
- Official **MOD11 checksum algorithm**

Invalid checksum or incorrect length will result in validation errors.

---

# Identity Types

The following identity type is implemented:

| Code | Description         |
|------|---------------------|
| `ON` | Organization Number |

This identity is validated using the MOD11 checksum.

---

# Scope & Design Decisions

## Included

- VAT modeling with 3 current rates
- Reverse charge handling
- VAT-exempt supplies
- Simplified invoices
- MOD11 validation for organization numbers
- Tax ID normalization
- Regime registration and CLI integration

## Out of Scope

- Historical VAT rate transitions
- Sector-specific VAT exceptions
- Special schemes (e.g., VOEC)
- Sector-specific levies (e.g., the 11.1% raw fish levy)
- EHF schema validation
- Full e-invoicing compliance rules

This implementation aims to provide a robust and extensible baseline rather than exhaustive legal coverage.

The regime follows GOBL global VAT conventions and reuses global VAT keys to maintain cross-regime consistency.

---

# E-Invoicing Context

Norway does not mandate a single centralized invoice reporting platform comparable to:

- Italy (FatturaPA)
- Greece (myDATA)

However:

- **EHF (Elektronisk Handelsformat)** is widely used
- Especially for B2G transactions
- And commonly adopted in B2B

GOBL invoices can be transformed into local formats via integration layers if required.

---

# Resources

- Norwegian Tax Administration (Skatteetaten):  
  https://www.skatteetaten.no/en/

- VAT Rates:  
  https://www.skatteetaten.no/en/rates/value-added-tax/

- VAT Overview:  
  https://www.skatteetaten.no/en/business-and-organisation/vat-and-duties/vat/

- Norwegian VAT Act (Merverdiavgiftsloven):  
  https://lovdata.no/dokument/NL/lov/2009-06-19-58