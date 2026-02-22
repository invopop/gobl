# 🇯🇵 GOBL Japan (JP) Tax Regime — Research Notes

> Research notes for implementing and maintaining `regimes/jp` in the [invopop/gobl](https://github.com/invopop/gobl) project.
> This document complements the regime's `README.md` by capturing background, rationale, and references based on patterns from existing regimes (`se`, `in`, `pl`, `es`) and Japanese tax law.

---

## 1. Regime Overview

| Field             | Value                                                                                                       |
|-------------------|-------------------------------------------------------------------------------------------------------------|
| **Country**       | Japan (日本)                                                                                                  |
| **ISO Code**      | `JP`                                                                                                        |
| **Currency**      | JPY (Japanese Yen, ¥)                                                                                       |
| **Tax Authority** | National Tax Agency (国税庁, NTA) — https://www.nta.go.jp                                                      |
| **Primary Tax**   | Consumption Tax (消費税, *Shōhizei*) — analogous to VAT                                                        |
| **Key Reform**    | Qualified Invoice System (適格請求書等保存方式, *Tekikaku Seikyūsho-tō Hozon Hōshiki*), effective **October 1, 2023** |

---

## 2. Tax Categories

Japan has a single main indirect tax: **Consumption Tax (CT)**. It is administered at the national level but includes a local portion.

### 2.1 Tax Category Code

Japan's consumption tax (消費税, *Shōhizei* — JCT) is implemented in GOBL using the **VAT** category code, consistent with other non-EU regimes (e.g. UK, Norway, Switzerland) where the local tax has a different name but is functionally a value-added / consumption tax. Documents use `cat: VAT`; the regime's `Name` and `Title` use the local labels "Japanese Consumption Tax" / "JCT" (消費税) for display and exports.

### 2.2 Tax Rates

The regime uses GOBL's standard VAT rate keys (`general`, `reduced`, `zero`; `standard` in input is normalized to `general`).

| Rate Key                      | Rate    | Japanese Term            | Applicable Since | Description                                                                        |
|-------------------------------|---------|--------------------------|------------------|------------------------------------------------------------------------------------|
| `general` (input: `standard`) | **10%** | 標準税率 (*Hyōjun Zeiritsu*) | Oct 1, 2019      | Most goods and services                                                            |
| `reduced`                     | **8%**  | 軽減税率 (*Keigen Zeiritsu*) | Oct 1, 2019      | Food & non-alcoholic beverages (for home consumption), and subscription newspapers |
| `exempt`                      | **0%**  | 非課税 (*Hikazei*)          | —                | Certain medical, welfare, education, financial services (no input credit)          |
| `zero`                        | **0%**  | 輸出免税 (*Yushutsu Menzei*) | —                | Exports (zero-rated; input tax credits *are* recoverable)                          |

> **Note on rates history**: The standard rate was 5% until March 2014, then 8% until September 2019, then raised to 10%. A reduced rate of 8% was introduced simultaneously with the 10% increase to protect essential goods. When defining rate history in GOBL, capture the effective date of **October 1, 2019** as the start of the current dual-rate system.

### 2.3 Tax Rate Breakdown (National vs. Local)

Japan's consumption tax is technically split:

| Component                   | Standard (10%) | Reduced (8%) |
|-----------------------------|----------------|--------------|
| National CT                 | 7.8%           | 6.24%        |
| Local CT (*Chihō Shōhizei*) | 2.2%           | 1.76%        |
| **Total**                   | **10%**        | **8%**       |

For GOBL purposes, this split does **not** need to be modeled as separate categories — invoices display only the combined rate. This is consistent with how Japan's NTA mandates invoice presentation.

---

## 3. Identity Types

### 3.1 Corporate Number (法人番号, *Hōjin Bangō*)

| Property        | Value                                                                                |
|-----------------|--------------------------------------------------------------------------------------|
| **Code**        | `CN` (or `JP-CN` officially)                                                         |
| **Format**      | 13 digits, numeric only                                                              |
| **Structure**   | 1 check digit + 12-digit company registration number                                 |
| **Issuer**      | National Tax Agency                                                                  |
| **Published**   | Yes — public on [NTA Corporate Number Site](https://www.houjin-bangou.nta.go.jp/en/) |
| **Who gets it** | All registered corporations, government bodies, certain associations                 |

**Check digit algorithm** (first digit):
The check digit `p` is derived from the 12-digit base number `n₁₂...n₁` using a weighted sum:

```
Q = Σ (nᵢ × Pᵢ)  where Pᵢ alternates between 2 and 1 (rightmost = 1)
p = 9 - (Q mod 9)   if Q mod 9 ≠ 0
p = 0               if Q mod 9 = 0
```

Reference implementation: https://github.com/kufu/tsubaki (Ruby), which can be ported to Go.

```go
const IdentityTypeCorporateNumber cbc.Code = "CN" // 法人番号 Hōjin Bangō
```

### 3.2 Qualified Invoice Issuer Number (適格請求書発行事業者登録番号)

| Property         | Value                                                                   |
|------------------|-------------------------------------------------------------------------|
| **Format**       | `T` + 13-digit Corporate Number                                         |
| **Example**      | `T1234567890123`                                                        |
| **Who has it**   | CT-registered businesses that have applied as Qualified Invoice Issuers |
| **Mandatory on** | All qualified invoices from Oct 1, 2023                                 |
| **Verification** | NTA Invoice Registration Site: https://www.invoice-kohyo.nta.go.jp      |

This is the **Tax ID** used on invoices (analogous to a VAT number in the EU). It should be modeled as the regime's primary tax identity code (`JP`), stored in the `tax.Identity` with `Code: "T" + corporateNumber`.

**Validation rule**: Must match `^T[1-9]\d{12}$` (T + 13 digits, first digit non-zero).

```go
// Tax identity format: "T" followed by 13-digit Corporate Number
// Pattern: T\d{13}
// Example: T1234567890123
```

### 3.3 Individual Number (個人番号, My Number — *Kojin Bangō*)

| Property        | Value                                                                      |
|-----------------|----------------------------------------------------------------------------|
| **Code**        | `MN`                                                                       |
| **Format**      | 12 digits                                                                  |
| **Check digit** | Last digit, Luhn-like algorithm                                            |
| **Note**        | Highly restricted — disclosure punishable by law; **not used on invoices** |

My Number should **not** be a supported identity type for invoicing in GOBL. It is not permissible to display on business invoices.

---

## 4. Qualified Invoice System (インボイス制度)

### 4.1 Required Invoice Contents

A valid **Qualified Invoice** (*Tekikaku Seikyūsho*, 適格請求書) must include:

1. **Name and Qualified Invoice Issuer Number** of the supplier
2. **Transaction date**
3. **Transaction details**, with clear identification of items subject to the reduced 8% rate
4. **Transaction amount by applicable tax rate** (8% group and 10% group shown separately)
5. **Consumption tax amount** per applicable tax rate (rounding rules apply — see §5)
6. **Name of the counterparty** (buyer)

### 4.2 Simplified Qualified Invoice

A **Simplified Qualified Invoice** (*Kanryaku Tekikaku Seikyūsho*, 簡略適格請求書) may be issued by certain businesses (e.g., retail, restaurants, taxis) and does **not** require the buyer's name. Requirements otherwise the same as above.

### 4.3 Who Must Register

- Any JCT taxpayer wishing to issue qualified invoices must apply to the NTA
- Tax-exempt entities (taxable sales ≤ ¥10M in the base period) cannot be Qualified Invoice Issuers unless they opt into the tax system
- Upon registration, the issuer receives their `T` + 13-digit registration number

### 4.4 Transitional Measures (Partial Credits for Non-Qualified Invoices)

| Period                     | Input CT Credit Allowed                      |
|----------------------------|----------------------------------------------|
| Oct 1, 2023 – Sep 30, 2026 | **80%** of CT on non-qualified invoices      |
| Oct 1, 2026 – Sep 30, 2029 | **50%** of CT on non-qualified invoices      |
| From Oct 1, 2029           | **0%** — no credit on non-qualified invoices |

> GOBL scope: These transitional rules affect the buyer's accounting, not the invoice format itself, so they likely do not need to be modeled in the regime beyond documentation.

---

## 5. Tax Calculation & Rounding Rules

Japan's CT rounding rules are specific and must be implemented correctly:

- Tax amounts are calculated **per tax rate group**, not per line item
- **Rounding is done once per rate per invoice** (not per line)
- Rounding method: **round down** (truncate fractions — 切り捨て) by default per NTA guidance
- The standard method: sum all taxable amounts for each rate, then apply the rate, then round down

```
CT Amount (10%) = floor(Σ taxable_amounts_at_10% × 0.10)
CT Amount (8%)  = floor(Σ taxable_amounts_at_8% × 0.08)
```

> **GOBL consideration**: GOBL's tax calculation engine handles rounding at the regime level. The Japan regime should ensure the `RoundingRule` is set appropriately. Review how the tax calculation precision is set — using an exponent of `0` for JPY (no decimal places) is critical since ¥ has no subdivision.

---

## 6. Currency

- **JPY** has **no decimal places** (exponent = 0)
- All amounts should be expressed as whole numbers
- GOBL must be configured accordingly: `num.MakeAmount` with precision 0

---

## 7. Invoice Types & Corrections

### Standard Invoice Types
- Standard tax invoice (*Tekikaku Seikyūsho*)
- Simplified qualified invoice (*Kanryaku Tekikaku Seikyūsho*)
- Credit note (*返還インボイス* / *Henkan Invoice* — for returns/refunds)

### Correction Documents
Japan does not have a strict equivalent to the EU credit note system, but businesses issue:
- **Return invoices** (*返還インボイス*) for goods/services returned
- **Amendment invoices** for corrections — typically by issuing a corrective document referencing the original

GOBL's `credit-note` type can map to return invoices. Corrections are likely out of MVP scope.

---

## 8. Exemptions & Zero-Rating

| Type                    | CT Rate | Input Credit? | Examples                                                            |
|-------------------------|---------|---------------|---------------------------------------------------------------------|
| **Zero-rated (export)** | 0%      | ✅ Yes         | Exports of goods, cross-border services                             |
| **Exempt**              | 0%      | ❌ No          | Land sales, certain financial services, medical, welfare, education |
| **Outside scope**       | N/A     | ❌ No          | Wages, dividends, insurance claims                                  |

The distinction between zero-rated and exempt is important and should be represented as separate rate types or via a tax exemption extension key.

---

## 9. Files to Create

Following the pattern of `regimes/se` and `regimes/in`, the Japan regime should consist of:

### Required

| File                           | Purpose                                                              |
|--------------------------------|----------------------------------------------------------------------|
| `regimes/jp/jp.go`             | Main regime definition: `New()`, `Normalize()`, `Validate()`         |
| `regimes/jp/tax_categories.go` | VAT category (JCT locally) with general (10%) and reduced (8%) rates |
| `regimes/jp/tax_identities.go` | Normalize/validate the `T` + 13-digit invoice issuer number          |
| `regimes/jp/org_identities.go` | Normalize/validate Corporate Number (13-digit with check digit)      |
| `regimes/jp/README.md`         | Documentation (sources, decisions, Japan-specific notes)             |

### Optional / Deferred

| File                         | Purpose                                                     |
|------------------------------|-------------------------------------------------------------|
| `regimes/jp/bill_invoice.go` | Invoice-level validation (e.g., require tax rate breakdown) |
| `regimes/jp/scenarios.go`    | Standard invoice scenarios                                  |
| `regimes/jp/corrections.go`  | Credit note / return invoice correction options             |

### Supporting Files

| File                              | Purpose                      |
|-----------------------------------|------------------------------|
| `examples/jp/invoice-jp.yaml`     | Example uncalculated invoice |
| `examples/jp/out/invoice-jp.json` | Calculated envelope output   |

### Registration in

- `regimes/regimes.go` — add import and `jp.New()` call

---

## 10. Key Implementation Decisions

### Tax Identity (Tax ID)

The Japan tax identity is the **Qualified Invoice Issuer Number**:
- Format: `T` followed by 13 digits (e.g., `T1234567890123`)
- Stored in `tax.Identity.Code` for the `JP` country
- The `T` prefix is mandatory (per NTA) — keep it normalized with prefix (similar to how SE regime keeps `SE...01`)
- Validation regex: `^T[1-9]\d{12}$`

```go
// Normalize: strip spaces, ensure uppercase T prefix
// Validate: check length = 14, starts with T, digits 2-14 are numeric, first digit non-zero
```

### Corporate Number Check Digit

Port the check digit algorithm to Go. The algorithm (from NTA documentation):

```go
func validateCorporateNumberCheckDigit(cn string) bool {
// cn is 13 digits
digits := parseDigits(cn) // [d0, d1, ..., d12]
// weights alternate 2, 1 from right side (excluding check digit)
sum := 0
for i := 1; i <= 12; i++ {
w := 1
if i % 2 == 1 { w = 2 }  // odd positions (from right) get weight 2
// Actually: positions from right: n1=weight1, n2=weight2, ...
// Simpler: Pᵢ = 2 if i is odd (counting from right), 1 if even
sum += digits[13-i] * w
}
expected := 9 - (sum % 9)
if sum % 9 == 0 { expected = 0 }
return digits[0] == expected
}
```

> **Verify** against https://github.com/kufu/tsubaki/blob/master/lib/tsubaki/corporate_number.rb before finalizing.

### Currency Precision

JPY is a zero-decimal currency. Configure this via the regime or rely on the currency definition already in GOBL. Verify how GOBL's `currency.JPY` handles precision.

### Tax Calculation Rounding

Set up round-down semantics for CT amounts, consistent with NTA rules (per-rate-group, not per-line).

---

## 11. Validation Rules (Official Basis)

The regime enforces the following rules; each is tied to official NTA or Customs guidance.

| Rule                                    | Implementation                                                                            | Official source                                                                                                                                                                                                                    |
|-----------------------------------------|-------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Supplier name and T-number required** | Invoice must have supplier with name and tax identity (Qualified Invoice Issuer Number).  | NTA: Qualified invoice must include "name and Qualified Invoice Issuer Number of the supplier". [Invoice system overview](https://www.nta.go.jp/taxes/shiraberu/zeimokubetsu/shohi/keigenzeiritsu/invoice.htm)                     |
| **Buyer name**                          | Required on standard qualified invoices; optional when tag `simplified` is set (簡易適格請求書). | NTA: Simplified qualified invoice may omit buyer name (e.g. retail, restaurants, taxis). Same source as above.                                                                                                                     |
| **Export = zero-rated only**            | When tag `export` is set, every VAT line must be zero-rated (rate `zero` or key `zero`).  | Japan Customs: "Consumption tax exemption on exports" applies to export of goods. [Customs FAQ 5003](https://www.customs.go.jp/english/c-answer_e/extsukan/5003_e.htm). NTA treats export supplies as zero-rated (仕入税額控除の適用あり).    |
| **Tax identity format**                 | T-number: `T` + 13 digits, first digit non-zero; check digit validated.                   | NTA Corporate Number check digit: [Explanation (PDF)](https://www.houjin-bangou.nta.go.jp/en/setsumei/images/05check_flow_chart.pdf). Invoice issuer registration: [invoice-kohyo.nta.go.jp](https://www.invoice-kohyo.nta.go.jp). |
| **Corporate Number (org identity)**     | 13 digits, check digit validated; hyphens/spaces normalized.                              | NTA: [Corporate Number publication site](https://www.houjin-bangou.nta.go.jp/en/).                                                                                                                                                 |

**Exempt transactions (非課税)**  
Exempt supplies (e.g. interest, insurance premiums, land, certain financial/medical/welfare/education services) use VAT key `exempt` and have no consumption tax amount; input tax credit is not recoverable. NTA: [Consumption tax basic knowledge](https://www.nta.go.jp/english/taxes/consumption_tax/01.htm) and [guides/notifications](https://www.nta.go.jp/english/taxes/consumption_tax/index.htm).

**Edge scenarios**  
- **Export goods**: Use tag `export` and VAT rate `zero` (or key `zero`) on all lines.  
- **Exempt goods/services**: Use VAT key `exempt` (no percent).  
- **Simplified qualified invoice**: Use tag `simplified`; customer name may be omitted.  
- **Self-billing (仕入明細書)**: Use tag `self-billing`; buyer-issued, valid as qualified invoice substitute when supplier confirms.

---

## 12. Sources

| Source                                   | URL                                                                                                                                                                        |
|------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| NTA — Invoice System Overview            | https://www.nta.go.jp/taxes/shiraberu/zeimokubetsu/shohi/keigenzeiritsu/invoice.htm                                                                                        |
| NTA — Corporate Number Publication Site  | https://www.houjin-bangou.nta.go.jp/en/                                                                                                                                    |
| NTA — Invoice Registration Check         | https://www.invoice-kohyo.nta.go.jp                                                                                                                                        |
| EU-Japan Centre — QIS Summary            | https://www.eu-japan.eu/qualified-invoice-system                                                                                                                           |
| EY Japan — Implementation Considerations | https://www.ey.com/en_jp/technical/ey-japan-tax-library/tax-alerts/2022/japan-s-consumption-tax-reform-will-be-effective-from-1-october-2023-implementation-considerations |
| Sovos — QIS Explainer                    | https://sovos.com/blog/vat/what-is-japans-qualified-invoice-system/                                                                                                        |
| Corporate Number Wikipedia               | https://en.wikipedia.org/wiki/Corporate_Number                                                                                                                             |
| SmartStart — Corporate Number Guide      | https://smartstartjapan.com/tax-identification-number-in-japan/                                                                                                            |
| Tsubaki (Ruby validator)                 | https://github.com/kufu/tsubaki                                                                                                                                            |
| Microsoft Dynamics QIS Impl.             | https://learn.microsoft.com/en-us/dynamics365/finance/localizations/japan/apac-jpn-qualified-invoice-system                                                                |
| HLS Global — QIS Transition Measures     | https://hls-global.jp/en/2023/05/17/introduction-to-the-new-japanese-invoice-system-implementation-qualified-invoice-issuers-2/                                            |

---

## 13. Open Questions / TODOs

- [x] **Check digit algorithm**: Implemented in `tax_identities.go` and `org_identities.go` using NTA weighted alternating 1/2 sum mod 9
- [x] **My Number on invoices**: Modeled as `MN` identity type with doc note that it is not used on invoices per privacy law
- [ ] **PEPPOL support**: Japan does not currently use PEPPOL for domestic e-invoicing. The NTA has a separate digital invoice initiative (Peppol-based JP PINT is under development). Consider a `jp-pint` addon later
- [x] **Simplified invoice rules**: `TagSimplified` tag + scenario + customer-not-required validation implemented
- [x] **Reverse charge**: Confirmed Japan has no domestic reverse charge mechanism; not modeled in the regime
- [ ] **Inbound cross-border services**: Foreign digital service providers must register. How should their tax IDs be structured in GOBL?
- [x] **Rate history dates**: Historical rates added in `tax_categories.go`: 3% (1989), 5% (1997), 8% (2014)
- [x] **`regimes.go` registration**: Registered via `init()` in `jp.go`
