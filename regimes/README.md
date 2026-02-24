# Tax Regimes

Tax regimes in GOBL define country-specific tax rules, rates, validation logic, and normalization behavior. Each regime is a self-contained Go package inside the `regimes/` directory, named after its [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) country code (e.g., `es` for Spain, `de` for Germany).

## Existing Regimes

| Code | Country |
|------|---------|
| `ae` | United Arab Emirates |
| `ar` | Argentina |
| `at` | Austria |
| `be` | Belgium |
| `br` | Brazil |
| `ca` | Canada |
| `ch` | Switzerland |
| `co` | Colombia |
| `de` | Germany |
| `dk` | Denmark |
| `es` | Spain |
| `fr` | France |
| `gb` | United Kingdom |
| `gr` | Greece |
| `ie` | Ireland |
| `in` | India |
| `it` | Italy |
| `mx` | Mexico |
| `nl` | Netherlands |
| `pl` | Poland |
| `pt` | Portugal |
| `se` | Sweden |
| `sg` | Singapore |
| `us` | United States |

## How Regimes Work

When a regime package is imported, its `init()` function calls `tax.RegisterRegimeDef()` to register the regime definition globally. The `regimes/regimes.go` file uses blank imports to ensure all regimes are loaded:

```go
import (
    _ "github.com/invopop/gobl/regimes/es"
    _ "github.com/invopop/gobl/regimes/de"
    // ...
)
```

Once registered, GOBL automatically applies the correct regime's normalization, validation, tax rates, and scenarios when processing documents for that country.

## Creating a New Regime

### Step 1: Copy the Template

Duplicate the `template/` directory and rename it to the 2-letter country code:

```bash
cp -r regimes/template regimes/xx
```

Replace `xx` with the actual country code (e.g., `jp` for Japan).

### Step 2: Define the Main Regime File

Rename `template.go` to `xx.go` (matching your country code). This file is the entry point for the regime and must contain three things:

**1. An `init()` function that registers the regime:**

```go
func init() {
    tax.RegisterRegimeDef(New())
}
```

**2. A `New()` function that returns a `*tax.RegimeDef`:**

```go
func New() *tax.RegimeDef {
    return &tax.RegimeDef{
        Country:  "XX",
        Currency: currency.XXX,
        Name: i18n.String{
            i18n.EN: "Country Name",
            // i18n.XX: "Local Name",
        },
        TimeZone:   "Region/City",           // IANA Time Zone
        TaxScheme:  tax.CategoryVAT,          // or omit (e.g., US has no scheme)
        Tags:       []*tax.TagSet{...},       // optional
        Identities: identityDefinitions(),    // optional, from identities.go
        Extensions: extensionDefinitions,     // optional, from extensions.go
        Categories: taxCategories(),          // required, from tax_categories.go
        Scenarios:  []*tax.ScenarioSet{...},  // optional, from scenarios.go
        Corrections: correctionDefinitions(), // optional, from corrections.go
        Normalizer: Normalize,
        Validator:  Validate,
    }
}
```

**3. `Normalize` and `Validate` functions** that use type-switching to route to specific handlers:

```go
func Normalize(doc any) {
    switch obj := doc.(type) {
    case *tax.Identity:
        normalizeTaxIdentity(obj)
    case *org.Identity:
        normalizeOrgIdentity(obj)
    }
}

func Validate(doc any) error {
    switch obj := doc.(type) {
    case *tax.Identity:
        return validateTaxIdentity(obj)
    case *bill.Invoice:
        return validateInvoice(obj)
    case *org.Party:
        return validateParty(obj)
    }
    return nil
}
```

The supported types you can normalize and validate include:

- `*tax.Identity` - Tax identification numbers
- `*org.Identity` - Organization identity documents (passport, national ID, etc.)
- `*org.Party` - Supplier/customer party data
- `*bill.Invoice` - Invoice-level rules
- `*bill.Line` - Line item validation
- `*tax.Combo` - Tax combination validation

### Step 3: Define Tax Categories

Create `tax_categories.go` to define all tax categories and their rates. This is the most important data file in a regime:

```go
package xx

import (
    "github.com/invopop/gobl/cal"
    "github.com/invopop/gobl/cbc"
    "github.com/invopop/gobl/i18n"
    "github.com/invopop/gobl/num"
    "github.com/invopop/gobl/tax"
)

func taxCategories() []*tax.CategoryDef {
    return []*tax.CategoryDef{
        {
            Code: tax.CategoryVAT,
            Name: i18n.String{
                i18n.EN: "VAT",
                // i18n.XX: "Local abbreviation",
            },
            Title: i18n.String{
                i18n.EN: "Value Added Tax",
                // i18n.XX: "Local full name",
            },
            Retained: false,
            Keys:     tax.GlobalVATKeys(),
            Rates: []*tax.RateDef{
                {
                    Keys: []cbc.Key{tax.KeyStandard},
                    Rate: tax.RateGeneral,
                    Name: i18n.String{
                        i18n.EN: "Standard rate",
                    },
                    Values: []*tax.RateValueDef{
                        {
                            Since:   cal.NewDate(2024, 1, 1),
                            Percent: num.MakePercentage(20, 2), // 20%
                        },
                    },
                },
                {
                    Keys: []cbc.Key{tax.KeyStandard},
                    Rate: tax.RateReduced,
                    Name: i18n.String{
                        i18n.EN: "Reduced rate",
                    },
                    Values: []*tax.RateValueDef{
                        {
                            Since:   cal.NewDate(2024, 1, 1),
                            Percent: num.MakePercentage(10, 2), // 10%
                        },
                    },
                },
            },
        },
    }
}
```

Key points about tax categories:

- **`Code`**: Use predefined codes like `tax.CategoryVAT`, `tax.CategoryGST`, `tax.CategoryST`, or define custom ones (e.g., `"IRPF"`, `"IGIC"`).
- **`Retained`**: Set to `true` for withholding taxes (e.g., income tax withheld at source).
- **`Keys`**: Use `tax.GlobalVATKeys()` for standard VAT/GST regimes.
- **`Rates`**: Each rate has `Keys` (e.g., `tax.KeyStandard`, `tax.KeyExempt`, `tax.KeySuperReduced`) and historical `Values` with effective dates.
- **Historical rates**: Include multiple `Values` entries with different `Since` dates to support historical tax rate lookups. Order them newest first.

### Step 4: Define Tax Identity Validation

Create `tax_identity.go` to handle normalization and validation of tax identification numbers:

```go
package xx

import (
    "github.com/invopop/gobl/cbc"
    "github.com/invopop/gobl/tax"
    "github.com/invopop/validation"
)

func normalizeTaxIdentity(tID *tax.Identity) {
    if tID == nil {
        return
    }
    // Typical normalizations:
    // - Remove country prefix: tax.NormalizeIdentity(tID)
    // - Uppercase: strings.ToUpper()
    // - Remove whitespace/separators
}

func validateTaxIdentity(tID *tax.Identity) error {
    return validation.ValidateStruct(tID,
        validation.Field(&tID.Code, validation.By(validateTaxCode)),
    )
}

func validateTaxCode(value interface{}) error {
    code, ok := value.(cbc.Code)
    if !ok || code == "" {
        return nil
    }
    // Validate format (regex), length, and checksum
    return nil
}
```

### Step 5: Optional Files

Depending on the complexity of the regime, add any of these files:

#### `scenarios.go` - Invoice Scenarios

Scenarios automatically apply notes to documents based on tags:

```go
package xx

import (
    "github.com/invopop/gobl/bill"
    "github.com/invopop/gobl/cbc"
    "github.com/invopop/gobl/org"
    "github.com/invopop/gobl/tax"
)

func invoiceScenarios() *tax.ScenarioSet {
    return &tax.ScenarioSet{
        Schema: bill.ShortSchemaInvoice,
        List: []*tax.Scenario{
            {
                Tags: []cbc.Key{tax.TagReverseCharge},
                Note: &tax.ScenarioNote{
                    Key:  org.NoteKeyLegal,
                    Src:  tax.TagReverseCharge,
                    Text: "Reverse Charge / Local translation.",
                },
            },
        },
    }
}
```

#### `corrections.go` - Correction Definitions

Define which correction types (credit notes, debit notes, corrective invoices) the regime supports:

```go
package xx

import (
    "github.com/invopop/gobl/bill"
    "github.com/invopop/gobl/cbc"
    "github.com/invopop/gobl/tax"
)

func correctionDefinitions() []*tax.CorrectionDefinition {
    return []*tax.CorrectionDefinition{
        {
            Schema: bill.ShortSchemaInvoice,
            Types: []cbc.Key{
                bill.InvoiceTypeCreditNote,
                bill.InvoiceTypeDebitNote,
                // bill.InvoiceTypeCorrective, // if supported
            },
        },
    }
}
```

#### `identities.go` - Organization Identity Types

Define non-tax identity types recognized by the regime (passports, national IDs, etc.):

```go
package xx

import (
    "github.com/invopop/gobl/cbc"
    "github.com/invopop/gobl/i18n"
    "github.com/invopop/gobl/org"
)

func identityDefinitions() []*cbc.Definition {
    return []*cbc.Definition{
        {
            Key: org.IdentityKeyPassport,
            Name: i18n.String{
                i18n.EN: "Passport",
            },
        },
    }
}
```

#### `invoices.go` - Invoice-Specific Validation

Add custom validation rules for invoices (e.g., requiring supplier tax IDs):

```go
package xx

import (
    "github.com/invopop/gobl/bill"
    "github.com/invopop/validation"
)

func validateInvoice(inv *bill.Invoice) error {
    return validation.ValidateStruct(inv,
        // Add regime-specific invoice validation rules
    )
}
```

#### `extensions.go` - Custom Extension Keys

Define regime-specific extension codes for data that doesn't fit standard GOBL structures:

```go
package xx

import (
    "github.com/invopop/gobl/cbc"
    "github.com/invopop/gobl/i18n"
)

var extensionDefinitions = []*cbc.Definition{
    {
        Key: "xx-custom-code",
        Name: i18n.String{
            i18n.EN: "Custom Code",
        },
        Pattern: `^\d{7}$`, // optional regex for validation
    },
}
```

#### `party.go` - Party-Specific Logic

Normalization and validation for supplier/customer parties.

### Step 6: Register the Regime

Add a blank import for the new regime in `regimes/regimes.go`:

```go
import (
    // ... existing imports ...
    _ "github.com/invopop/gobl/regimes/xx"
)
```

### Step 7: Write Tests

Every `.go` file with logic should have a corresponding `_test.go` file. At minimum:

- `xx_test.go` - Test the `New()` function produces a valid regime definition
- `tax_identity_test.go` - Test normalization and validation of tax IDs

The regime will also be automatically tested by `regimes_test.go`, which validates all registered regime definitions.

### Step 8: Add a README

Create a `README.md` inside your regime directory. Document:

- Overview of the country's tax system
- Tax categories and rates
- Tax identity format and validation rules
- Any special features (scenarios, extensions, corrections)
- Reference links to official government sources

### Step 9: Generate and Test

Run the following commands from the repository root:

```bash
go generate .
go test ./regimes/xx/...
go test ./regimes/...
```

## `RegimeDef` Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `Country` | `l10n.TaxCountryCode` | Yes | ISO 3166-1 alpha-2 country code |
| `Currency` | `currency.Code` | Yes | ISO 4217 currency code |
| `Name` | `i18n.String` | Yes | Localized regime name |
| `TimeZone` | `string` | Yes | IANA Time Zone (e.g., `"Europe/Madrid"`) |
| `TaxScheme` | `cbc.Code` | No | Principal tax scheme (e.g., `tax.CategoryVAT`) |
| `Description` | `i18n.String` | No | Detailed regime description |
| `AltCountryCodes` | `[]l10n.Code` | No | Alternative country codes |
| `Zone` | `l10n.Code` | No | Sub-country locality or region |
| `CalculatorRoundingRule` | `cbc.Key` | No | Rounding rule override (default: `sum-then-round`) |
| `Tags` | `[]*tax.TagSet` | No | Document-level tag definitions |
| `Extensions` | `[]*cbc.Definition` | No | Custom extension key definitions |
| `Identities` | `[]*cbc.Definition` | No | Organization identity type definitions |
| `PaymentMeansKeys` | `[]*cbc.Definition` | No | Regime-specific payment means keys |
| `InboxKeys` | `[]*cbc.Definition` | No | Regime-specific inbox routing keys |
| `Scenarios` | `[]*tax.ScenarioSet` | No | Conditional scenario rules |
| `Corrections` | `[]*tax.CorrectionDefinition` | No | Supported correction types |
| `Categories` | `[]*tax.CategoryDef` | Yes | Tax category definitions with rates |
| `Normalizer` | `tax.Normalizer` | No | Function to normalize regime-specific data |
| `Validator` | `tax.Validator` | No | Function to validate regime-specific data |

## Typical File Structure

Regimes range from minimal to complex depending on the country's tax requirements:

```
regimes/xx/
  xx.go                  # Required: regime definition, init(), New(), Normalize, Validate
  xx_test.go             # Required: tests
  tax_categories.go      # Required: tax category and rate definitions
  tax_identity.go        # Common: tax ID normalization and validation
  tax_identity_test.go   # Common: tax ID tests
  identities.go          # Optional: organization identity type definitions
  scenarios.go           # Optional: tag-based scenario rules
  corrections.go         # Optional: credit/debit note definitions
  invoices.go            # Optional: invoice-specific validation
  party.go               # Optional: party-specific logic
  extensions.go          # Optional: custom extension keys
  README.md              # Recommended: regime documentation
```

## Best Practices

- **Use `i18n.String` for all user-facing text.** Always include at least an English (`i18n.EN`) translation, and add the official local language translation.
- **Include historical tax rates.** Add multiple `Values` entries with `Since` dates to handle rate changes over time, ordered newest first.
- **Include sources.** Reference official government documentation URLs using `cbc.Source` in tax categories and extensions.
- **Use functions to build definitions** (e.g., `taxCategories()` instead of `var taxCategories = ...`). This prevents accidental modification of shared data, although both patterns are used in existing regimes.
- **Keep normalization idempotent.** Running `Normalize` multiple times should produce the same result.
- **Return `nil` for unknown types** in both `Validate` and `Normalize`. Only handle types relevant to the regime.
- **Follow the validation library patterns.** Use `github.com/invopop/validation` for struct validation, consistent with the rest of the codebase.
- **Write comprehensive tests**, especially for tax identity validation with edge cases, checksum verification, and format variations.
