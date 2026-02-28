---
name: create-regime
description: Create a new GOBL tax regime for a given country. Accepts a country code and optional research URLs.
argument-hint: "<COUNTRY_CODE> [URL ...]"
---

# create-regime

Create a new GOBL tax regime for a given country.

## Usage

```
/create-regime <COUNTRY_CODE> [URL ...]
```

Where `<COUNTRY_CODE>` is an ISO 3166-1 alpha-2 code (e.g., `JP`, `KR`, `NO`, `ZA`).

Optionally, the user may provide one or more URLs pointing to tax authority pages, rate tables, tax ID format documentation, or other references. These will be fetched and used as primary research sources.

## Arguments

- **First argument** (required): the two-letter country code. If not provided, ask the user for it.
- **Remaining arguments** (optional): URLs to use as research sources. Accept any number of URLs in any position after the country code.

## Workflow

Execute the following phases in order. Do not skip phases. Wait for user confirmation between Phase 2 and Phase 3.

---

### Phase 1: Research

Gather comprehensive information about the country's tax system. There are two input modes — choose based on what the user provided:

#### Mode A: User provided URLs

If the user supplied one or more URLs alongside the country code:

1. **Fetch each URL** using `WebFetch`. For each URL, extract all tax-relevant information: rates, tax ID formats, invoice rules, legal requirements, etc.
2. **Identify gaps** — after processing all provided URLs, check which of the research areas below still have missing information.
3. **Fill gaps with web search** — use `WebSearch` and `WebFetch` to research only the missing areas. Tell the user which areas you're supplementing with web research.

#### Mode B: No URLs provided (autonomous research)

Use `WebSearch` and `WebFetch` to research all areas from scratch. Run multiple searches in parallel where possible to speed things up.

#### Research areas (required regardless of mode)

1. **Tax system type and rates**
   - Primary tax type: VAT, GST, Sales Tax, or other
   - Standard rate and all reduced/super-reduced/zero rates
   - Historical rate changes with effective dates (at least the last 10-15 years)
   - Any regional tax variations (like Spain's IGIC/IPSI)
   - Any retained/withholding taxes (like Spain's IRPF)

2. **Tax identity**
   - Official name of the tax identification number
   - Format: length, character types (numeric, alphanumeric), patterns
   - Any prefix conventions (country code prefix for EU VAT numbers)
   - Validation rules: regex patterns, checksum algorithms
   - Different types of tax IDs (personal vs. business, like Spain's NIF/CIF/NIE)

3. **Invoice requirements**
   - Correction types supported: credit notes, debit notes, corrective invoices
   - Required legal mentions for reverse charge, simplified invoices, etc.
   - Whether supplier tax ID is required on all invoices
   - Special invoice types or tags recognized

4. **Other details**
   - Official timezone (IANA format, e.g., "Asia/Tokyo")
   - Currency code (ISO 4217)
   - Country name in English and the official local language
   - Identity document types recognized (passport, national ID, etc.)
   - i18n language code for the local language

#### Search queries to try (Mode B, or for gap-filling in Mode A)

- `"[country] VAT rates history"` or `"[country] GST rates history"`
- `"[country] tax identification number format validation"`
- `"[country] tax ID checksum algorithm"`
- `"[country] invoice requirements credit note"`
- `"[country] OECD VAT rates"` (oecd.org has comprehensive rate tables)

#### Sourcing notes

- **Prefer official sources**: government tax authority websites, OECD data, EU VIES documentation.
- **Track every URL** you fetch — these will be included as `Sources` in tax category definitions and in the README.
- If a user-provided URL fails to load or redirects, inform the user and fall back to web search for that topic.
- When user URLs conflict with web search results, present both and let the user decide in Phase 2.

---

### Phase 2: Review

Present a structured summary to the user for validation. Use this format:

```
## Tax Regime Summary: [Country Name] ([XX])

### Tax System
- Type: [VAT/GST/Sales Tax]
- TaxScheme: [tax.CategoryVAT / tax.CategoryGST / omit]

### Tax Rates
| Rate | Type | Percentage | Since |
|------|------|-----------|-------|
| Standard | general | X% | YYYY-MM-DD |
| Reduced | reduced | X% | YYYY-MM-DD |
| ... | ... | ... | ... |

### Historical Rate Changes
[List significant rate changes]

### Tax Identity
- Name: [Official name]
- Format: [Description]
- Regex: [Pattern]
- Checksum: [Yes/No — describe algorithm if yes]
- Types: [Personal, Business, etc.]

### Corrections
- Credit Notes: [Yes/No]
- Debit Notes: [Yes/No]
- Corrective Invoices: [Yes/No]

### Scenarios (Legal Notes)
- Reverse Charge: [Required text if applicable]
- [Other scenarios]

### Details
- Timezone: [IANA timezone]
- Currency: [currency.XXX]
- Local language: [i18n.XX]
- Country name (local): [Name]

### Sources
[List every URL consulted — both user-provided and discovered via web search]
- User-provided: [URLs the user gave, or "none"]
- Researched: [URLs found via WebSearch/WebFetch]

### Files to Generate
[List which files will be created based on complexity]
```

**Wait for the user to confirm or correct this information before proceeding.** If the user provides additional URLs or corrections at this stage, fetch the new URLs and update the summary accordingly before moving on.

---

### Phase 3: Generation

Read the reference documents before generating code:
- `.claude/skills/create-regime/references/regime-patterns.md` — constants, types, helpers
- `.claude/skills/create-regime/references/code-templates.md` — Go code templates

Then generate all files under `regimes/xx/` (where `xx` is the lowercase country code).

#### Required files (always generate):

1. **`xx.go`** — Main regime file
   - Package declaration with doc comment
   - `init()` calling `tax.RegisterRegimeDef(New())`
   - `New()` returning `*tax.RegimeDef` with all fields
   - `Normalize()` function with type switch
   - `Validate()` function with type switch
   - Any exported constants for custom tax categories or rate codes
   - Reference the appropriate template from code-templates.md (minimal/standard/complex)

2. **`tax_categories.go`** — Tax category definitions
   - All tax categories with localized names
   - All rates with historical values (newest first)
   - Use `tax.GlobalVATKeys()` for VAT regimes, `tax.GlobalGSTKeys()` for GST
   - Include `Sources` with reference URLs

3. **`tax_identity.go`** — Tax ID normalization and validation
   - `normalizeTaxIdentity()` — at minimum call `tax.NormalizeIdentity(tID)`
   - `validateTaxIdentity()` — format validation with regex
   - Checksum validation if the country uses one
   - Only create if the regime has a TaxScheme (countries without tax schemes like US typically don't need this)

4. **`tax_identity_test.go`** — Tests for tax ID validation
   - Table-driven tests with valid and invalid examples
   - Test format validation, checksum, edge cases
   - Use real or realistic tax ID numbers for test cases
   - Only create if tax_identity.go was created

5. **`xx_test.go`** — Basic regime and invoice tests
   - `validInvoice()` helper function
   - Test that a valid invoice calculates and validates successfully
   - Test key validation rules (e.g., supplier tax ID required)

#### Conditional files (generate based on research):

6. **`scenarios.go`** — Only if the regime has legal note requirements
   - Reverse charge note (most VAT countries)
   - Simplified invoice note
   - Any regime-specific scenarios

7. **`corrections.go`** — Only if correction types need a separate file (otherwise inline in `xx.go`)
   - Define which correction types are supported

8. **`identities.go`** — Only if the regime recognizes specific org identity types
   - Define identity types (passport, national ID, tax number, etc.)
   - Include normalization and validation if applicable

9. **`invoices.go`** — Only if there are invoice-level validation rules
   - Supplier tax ID requirements
   - Conditional validation (e.g., simplified invoices exempt from some rules)

10. **`README.md`** — Regime documentation
    - Overview of the country's tax system
    - Tax categories and current rates
    - Tax identity format and validation
    - Correction types supported
    - Links to official government sources

#### Code generation guidelines:

- **Always use `i18n.String`** for user-facing text with at least `i18n.EN` and the local language
- **Order historical rates newest first** in `Values` arrays
- **Use `num.MakePercentage(value, exp)`** for rate percentages: `num.MakePercentage(19, 2)` = 19%, `num.MakePercentage(210, 3)` = 21%
- **Use `num.NewPercentage(value, exp)`** for pointer percentages (surcharges)
- **Use `cal.NewDate(year, month, day)`** for effective dates
- **Use `validation.Skip`** as the last rule when validating struct fields to prevent recursive validation
- **Return `nil` for unknown types** in Normalize and Validate
- **Follow the variable vs function pattern** based on complexity: simple regimes use variables (`var taxCategories = ...`), complex ones use functions (`func taxCategories() ...`)
- **Use exported constants** for custom tax category codes and rate keys
- **Import paths** must be exact — check `regime-patterns.md` for the list

---

### Phase 4: Registration

Add the blank import to `regimes/regimes.go` in alphabetical order:

```go
_ "github.com/invopop/gobl/regimes/xx"
```

Insert it in the correct alphabetical position among the existing imports.

---

### Phase 5: Verification

Run the following commands and fix any issues:

```bash
go generate .
```

Then run tests:

```bash
go test ./regimes/xx/...
```

If tests fail, fix the issues and re-run. Common problems:
- Missing imports (run `goimports` or fix manually)
- Wrong percentage values (check `num.MakePercentage` arguments)
- Invalid timezone string
- Checksum algorithm bugs (verify with known-valid tax IDs)

Finally, verify the regime integrates correctly:

```bash
go test ./regimes/...
```

This runs `regimes_test.go` which calls `Validate()` on all registered regime definitions.

If all tests pass, report success. If any test fails, diagnose and fix the issue, then re-run.
