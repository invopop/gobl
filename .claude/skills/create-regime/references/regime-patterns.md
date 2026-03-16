# GOBL Regime Patterns Reference

Quick-reference for constants, types, and helpers used when building tax regimes.

## RegimeDef Struct

```go
type RegimeDef struct {
    Name                    i18n.String                // Required. Localized name.
    Description             i18n.String                // Optional.
    TimeZone                string                     // Required. IANA timezone (e.g. "Europe/Berlin").
    Country                 l10n.TaxCountryCode        // Required. ISO 3166-1 alpha-2 (e.g. "DE").
    AltCountryCodes         []l10n.Code                // Optional.
    Zone                    l10n.Code                  // Optional. Sub-country code.
    Currency                currency.Code              // Required. ISO 4217 (e.g. currency.EUR).
    TaxScheme               cbc.Code                   // Optional. E.g. tax.CategoryVAT. Omit for regimes without a tax scheme (like US).
    CalculatorRoundingRule  cbc.Key                    // Optional. Default: "sum-then-round".
    Tags                    []*tax.TagSet              // Optional. Document-level tags.
    Extensions              []*cbc.Definition          // Optional. Custom extension keys.
    Identities              []*cbc.Definition          // Optional. Organization identity type definitions.
    PaymentMeansKeys        []*cbc.Definition          // Optional.
    InboxKeys               []*cbc.Definition          // Optional.
    Scenarios               []*tax.ScenarioSet         // Optional.
    Corrections             []*tax.CorrectionDefinition // Optional.
    Categories              []*tax.CategoryDef         // Required. Tax categories with rates.
    Validator               tax.Validator              // Optional. func(doc any) error
    Normalizer              tax.Normalizer             // Optional. func(doc any)
}
```

## Standard Tax Category Codes

```go
tax.CategoryVAT  // cbc.Code = "VAT"  — Value Added Tax
tax.CategoryGST  // cbc.Code = "GST"  — Goods and Services Tax
tax.CategoryST   // cbc.Code = "ST"   — Sales Tax
```

Custom categories can be defined as `cbc.Code` constants (e.g. Spain's `"IRPF"`, `"IGIC"`, `"IPSI"`).

## Standard Tax Combo Keys (used in CategoryDef.Keys)

```go
tax.KeyStandard       // "standard"
tax.KeyZero           // "zero"
tax.KeyReverseCharge  // "reverse-charge"
tax.KeyExempt         // "exempt"
tax.KeyExport         // "export"
tax.KeyIntraCommunity // "intra-community"
tax.KeyOutsideScope   // "outside-scope"
```

Use `tax.GlobalVATKeys()` to get the standard key definitions for VAT/GST regimes.
Use `tax.GlobalGSTKeys()` for GST-specific key definitions.

## Standard Rate Keys (used in RateDef.Rate)

```go
tax.RateZero         // "zero"
tax.RateGeneral      // "general"
tax.RateIntermediate // "intermediate"
tax.RateReduced      // "reduced"
tax.RateSuperReduced // "super-reduced"
tax.RateSpecial      // "special"
tax.RateOther        // "other"
```

Composite rates use `cbc.Key.With()`: e.g. `tax.RateGeneral.With("eqs")` for equivalence surcharge.

## Standard Tags

```go
tax.TagSimplified    // "simplified"
tax.TagReverseCharge // "reverse-charge"
tax.TagCustomerRates // "customer-rates"
tax.TagSelfBilled    // "self-billed"
tax.TagReplacement   // "replacement"
tax.TagPartial       // "partial"
tax.TagBypass        // "bypass"
tax.TagB2G           // "b2g"
tax.TagExport        // "export"
tax.TagEEA           // "eea"
tax.TagPrepayment    // "prepayment"
tax.TagFactoring     // "factoring"
```

Default invoice tags are available via `bill.DefaultInvoiceTags()` which provides definitions for: simplified, reverse-charge, self-billed, customer-rates, partial, bypass.

## CategoryDef Struct

```go
type CategoryDef struct {
    Code        cbc.Code      // Required. E.g. tax.CategoryVAT
    Name        i18n.String   // Required. Short name (e.g. "VAT", "IVA")
    Title       i18n.String   // Optional. Full name (e.g. "Value Added Tax")
    Description *i18n.String  // Optional.
    Retained    bool          // true for withholding taxes (e.g. IRPF)
    Informative bool          // true for informational taxes
    Keys        []*KeyDef     // Optional. Use tax.GlobalVATKeys() for standard VAT.
    Rates       []*RateDef    // Rate definitions with historical values.
    Extensions  []cbc.Key     // Optional.
    Map         cbc.CodeMap   // Optional.
    Sources     []*cbc.Source // Optional. Reference URLs.
    Ext         Extensions    // Optional.
    Meta        cbc.Meta      // Optional.
}
```

## RateDef Struct

```go
type RateDef struct {
    Rate        cbc.Key          // Required. E.g. tax.RateGeneral
    Keys        []cbc.Key        // Optional. E.g. []cbc.Key{tax.KeyStandard}
    Name        i18n.String      // Required. Localized name.
    Description i18n.String      // Optional.
    Values      []*RateValueDef  // Historical rate values, newest first.
    Meta        cbc.Meta         // Optional.
}
```

## RateValueDef Struct

```go
type RateValueDef struct {
    Ext       Extensions      // Optional. Filter by extension.
    Since     *cal.Date       // Optional. Effective date. Use cal.NewDate(year, month, day).
    Percent   num.Percentage  // Required. Use num.MakePercentage(value, exp).
    Surcharge *num.Percentage // Optional. Use num.NewPercentage(value, exp).
    Disabled  bool            // Optional. Mark as no longer active.
}
```

### Percentage helpers

```go
num.MakePercentage(19, 2)   // 0.19 = 19%
num.MakePercentage(210, 3)  // 0.210 = 21%
num.MakePercentage(7, 2)    // 0.07 = 7%
num.MakePercentage(40, 3)   // 0.040 = 4%
num.NewPercentage(52, 3)    // *num.Percentage = 0.052 = 5.2% (pointer)
```

## ScenarioSet / Scenario Structs

```go
type ScenarioSet struct {
    Schema string       // E.g. bill.ShortSchemaInvoice
    List   []*Scenario
}

type Scenario struct {
    Name    i18n.String         // Optional.
    Desc    i18n.String         // Optional.
    Types   []cbc.Key           // Filter: document types.
    Tags    []cbc.Key           // Filter: required tags.
    ExtKey  cbc.Key             // Filter: extension key.
    ExtCode cbc.Code            // Filter: extension code.
    Filter  func(doc any) bool  // Filter: custom function.
    Note    *ScenarioNote       // Output: note to add.
    Codes   cbc.CodeMap         // Output: codes to set.
    Ext     Extensions          // Output: extensions to set.
}

type ScenarioNote struct {
    Key  cbc.Key    // Note type. E.g. org.NoteKeyLegal
    Code cbc.Code   // Optional.
    Src  cbc.Key    // Source tag key.
    Text string     // Note content.
    Ext  Extensions // Optional.
}
```

## CorrectionDefinition Struct

```go
type CorrectionDefinition struct {
    Schema         string     // E.g. bill.ShortSchemaInvoice
    Types          []cbc.Key  // E.g. bill.InvoiceTypeCreditNote, bill.InvoiceTypeDebitNote
    Extensions     []cbc.Key  // Optional.
    ReasonRequired bool       // Optional.
    Stamps         []cbc.Key  // Optional.
    CopyTax        bool       // Optional.
}
```

## Invoice Type Constants

```go
bill.InvoiceTypeStandard   // "standard"
bill.InvoiceTypeProforma   // "proforma"
bill.InvoiceTypeCorrective // "corrective"
bill.InvoiceTypeCreditNote // "credit-note"
bill.InvoiceTypeDebitNote  // "debit-note"
bill.InvoiceTypeOther      // "other"
```

## Helper Functions

### Tax Identity

```go
tax.NormalizeIdentity(tID *tax.Identity)                     // Remove whitespace, uppercase, strip country prefix
tax.RequireIdentityCode                                      // Validation rule: code must be present
tax.ParseIdentity(tin string) (*tax.Identity, error)         // Parse "XXCODE" into Identity
```

### Normalization Helpers (cbc package)

```go
cbc.NormalizeNumericalCode(code cbc.Code) cbc.Code  // Strip non-numeric characters
```

### Registration

```go
tax.RegisterRegimeDef(rd *tax.RegimeDef)  // Call in init() to register
tax.RegimeDefFor(country string) *RegimeDef  // Lookup by country code
tax.AllRegimeDefs() []*RegimeDef             // All registered regimes
```

### Regime Context

```go
tax.WithRegime(country string) *tax.Regime  // For use in bill.Invoice.Regime field in tests
```

## Organization Identity Keys (org package)

```go
org.IdentityKeyPassport  // "passport"
org.IdentityKeyForeign   // "foreign"
org.IdentityKeyResident  // "resident"
org.IdentityKeyOther     // "other"
```

## Note Keys (org package)

```go
org.NoteKeyLegal   // "legal" — for legal/regulatory notes
```

## Language Codes (i18n package)

Common codes used in regimes:

```go
i18n.EN  // English
i18n.ES  // Spanish
i18n.DE  // German
i18n.FR  // French
i18n.IT  // Italian
i18n.PT  // Portuguese
i18n.NL  // Dutch
i18n.PL  // Polish
i18n.SV  // Swedish
i18n.DA  // Danish
i18n.EL  // Greek
i18n.JA  // Japanese
i18n.ZH  // Chinese
i18n.KO  // Korean
i18n.AR  // Arabic
i18n.HI  // Hindi
i18n.CA  // Catalan
i18n.GL  // Galician
i18n.EU  // Basque
```

All ISO 639-1 codes are available as `i18n.XX` constants.

## Currency Codes (currency package)

Common codes:

```go
currency.EUR  // Euro
currency.USD  // US Dollar
currency.GBP  // British Pound
currency.JPY  // Japanese Yen
currency.CHF  // Swiss Franc
currency.CAD  // Canadian Dollar
currency.AUD  // Australian Dollar
currency.INR  // Indian Rupee
currency.BRL  // Brazilian Real
currency.MXN  // Mexican Peso
currency.COP  // Colombian Peso
currency.ARS  // Argentine Peso
currency.SGD  // Singapore Dollar
currency.AED  // UAE Dirham
currency.SEK  // Swedish Krona
currency.DKK  // Danish Krone
currency.PLN  // Polish Zloty
currency.NOK  // Norwegian Krone
```

All ISO 4217 codes are available as `currency.XXX` constants.

## Validation Patterns (invopop/validation)

```go
validation.ValidateStruct(obj,
    validation.Field(&obj.Field, validation.Required),
    validation.Field(&obj.Field, validation.By(customFunc)),
    validation.Field(&obj.Field, validation.Match(regexpPattern)),
    validation.Field(&obj.Field, validation.When(condition, rules...)),
    validation.Field(&obj.Field, validation.Skip),  // Always add as last rule for struct fields
)
```

The `validation.Skip` rule prevents recursive validation of nested structs that have their own validators.

## Import Paths

```go
"github.com/invopop/gobl/bill"
"github.com/invopop/gobl/cal"
"github.com/invopop/gobl/cbc"
"github.com/invopop/gobl/currency"
"github.com/invopop/gobl/i18n"
"github.com/invopop/gobl/l10n"
"github.com/invopop/gobl/num"
"github.com/invopop/gobl/org"
"github.com/invopop/gobl/tax"
"github.com/invopop/validation"
```

## File Naming Conventions

| File | Purpose |
|------|---------|
| `xx.go` | Main regime: `init()`, `New()`, `Normalize`, `Validate` |
| `tax_categories.go` | Tax category and rate definitions |
| `tax_identity.go` | Tax ID normalization and validation |
| `tax_identity_test.go` | Tax ID tests |
| `scenarios.go` | Tag-based scenario rules |
| `corrections.go` | Correction type definitions |
| `identities.go` | Organization identity type definitions |
| `invoices.go` | Invoice-specific validation |
| `extensions.go` | Custom extension key definitions |
| `party.go` | Party-specific normalization/validation |
| `xx_test.go` | Regime and invoice tests |
| `README.md` | Documentation |
