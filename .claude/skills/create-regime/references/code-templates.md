# GOBL Regime Code Templates

Ready-to-adapt Go code templates for every file type in a tax regime. Replace `xx` with the country code, `XX` with the uppercase country code, and fill in country-specific details.

---

## 1. Main Regime File — `xx.go`

### Minimal (like US — no tax identity validation, no scenarios)

```go
// Package xx provides the tax regime definition for [Country Name].
package xx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:  "XX",
		Currency: currency.XXX,
		Name: i18n.String{
			i18n.EN: "Country Name",
		},
		TimeZone:   "Region/City",
		Categories: taxCategories(),
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
					bill.InvoiceTypeDebitNote,
				},
			},
		},
	}
}
```

### Standard (like DE — with tax identity, scenarios, identities)

```go
// Package xx provides the tax regime definition for [Country Name].
package xx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "XX",
		Currency:  currency.XXX,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Country Name",
			// i18n.XX: "Local Name",
		},
		TimeZone:   "Region/City",
		Scenarios:  []*tax.ScenarioSet{invoiceScenarios},
		Identities: identityDefinitions,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
		Validator:  Validate,
		Normalizer: Normalize,
		Categories: taxCategories,
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc interface{}) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	case *org.Identity:
		return validateOrgIdentity(obj)
	}
	return nil
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	case *org.Identity:
		normalizeOrgIdentity(obj)
	}
}
```

### Complex (like ES — with custom tags, custom tax categories, custom corrections)

```go
// Package xx provides tax regime support for [Country Name].
package xx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegimeDef(New())
}

// Local tax category definitions which are not considered standard.
const (
	TaxCategoryXXX cbc.Code = "XXX"
)

// Specific tax rate codes.
const (
	TaxRateCustom cbc.Key = "custom"
)

// New provides the tax regime definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:   "XX",
		Currency:  currency.XXX,
		TaxScheme: tax.CategoryVAT,
		Name: i18n.String{
			i18n.EN: "Country Name",
			// i18n.XX: "Local Name",
		},
		TimeZone: "Region/City",
		Tags: []*tax.TagSet{
			invoiceTags(),
		},
		Identities: identityDefinitions(),
		Categories: taxCategories(),
		Validator:  Validate,
		Normalizer: Normalize,
		Scenarios: []*tax.ScenarioSet{
			invoiceScenarios(),
		},
		Corrections: correctionDefinitions(),
	}
}

// Validate checks the document type and determines if it can be validated.
func Validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *tax.Identity:
		return validateTaxIdentity(obj)
	}
	return nil
}

// Normalize will perform any regime specific normalizations on the data.
func Normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Identity:
		normalizeTaxIdentity(obj)
	}
}
```

---

## 2. Tax Categories — `tax_categories.go`

### Standard VAT regime

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
		//
		// VAT
		//
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
			Sources: []*cbc.Source{
				{
					Title: i18n.String{
						i18n.EN: "Source description",
					},
					URL: "https://example.com/tax-rates",
				},
			},
			Retained: false,
			Keys:     tax.GlobalVATKeys(),
			Rates: []*tax.RateDef{
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateGeneral,
					Name: i18n.String{
						i18n.EN: "Standard rate",
						// i18n.XX: "Local name",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2024, 1, 1),
							Percent: num.MakePercentage(20, 2),
						},
						// Add older rates here, newest first
					},
				},
				{
					Keys: []cbc.Key{tax.KeyStandard},
					Rate: tax.RateReduced,
					Name: i18n.String{
						i18n.EN: "Reduced rate",
						// i18n.XX: "Local name",
					},
					Values: []*tax.RateValueDef{
						{
							Since:   cal.NewDate(2024, 1, 1),
							Percent: num.MakePercentage(10, 2),
						},
					},
				},
				// Add super-reduced, zero, etc. as needed
			},
		},
	}
}
```

### As a package-level variable (like DE)

```go
package xx

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	//
	// VAT
	//
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
					i18n.EN: "General rate",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2024, 1, 1),
						Percent: num.MakePercentage(20, 2),
					},
				},
			},
		},
	},
}
```

### GST regime (like India, Singapore, Australia)

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
			Code: tax.CategoryGST,
			Name: i18n.String{
				i18n.EN: "GST",
			},
			Title: i18n.String{
				i18n.EN: "Goods and Services Tax",
			},
			Retained: false,
			Keys:     tax.GlobalGSTKeys(),
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
							Percent: num.MakePercentage(10, 2),
						},
					},
				},
			},
		},
	}
}
```

### Sales Tax regime (like US)

```go
package xx

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

func taxCategories() []*tax.CategoryDef {
	return []*tax.CategoryDef{
		{
			Code: tax.CategoryST,
			Name: i18n.String{
				i18n.EN: "ST",
			},
			Title: i18n.String{
				i18n.EN: "Sales Tax",
			},
			Retained: false,
			Rates:    []*tax.RateDef{},
		},
	}
}
```

### Retained tax (withholding, like Spain's IRPF)

```go
{
	Code:     TaxCategoryIRPF,
	Retained: true,
	Name: i18n.String{
		i18n.EN: "IRPF",
	},
	Title: i18n.String{
		i18n.EN: "Personal income tax.",
	},
	Rates: []*tax.RateDef{
		{
			Rate: TaxRatePro,
			Name: i18n.String{
				i18n.EN: "Professional Rate",
			},
			Values: []*tax.RateValueDef{
				{
					Since:   cal.NewDate(2015, 7, 12),
					Percent: num.MakePercentage(150, 3),
				},
			},
		},
	},
},
```

---

## 3. Tax Identity — `tax_identity.go`

### Simple (just normalize with NormalizeIdentity)

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
	tax.NormalizeIdentity(tID)
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
	// Add format validation here
	return nil
}
```

### With regex and checksum (like DE)

```go
package xx

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	taxCodeRegexps = []*regexp.Regexp{
		regexp.MustCompile(`^[1-9]\d{8}$`),
	}
)

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
	val := code.String()

	match := false
	for _, re := range taxCodeRegexps {
		if re.MatchString(val) {
			match = true
			break
		}
	}
	if !match {
		return errors.New("invalid format")
	}

	return validateTaxCodeChecksum(val)
}

func validateTaxCodeChecksum(val string) error {
	// Implement checksum algorithm specific to the country
	_ = val
	return nil
}
```

### With type detection (like ES)

```go
package xx

import (
	"errors"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

const (
	TaxIdentityPersonal cbc.Key = "personal"
	TaxIdentityBusiness cbc.Key = "business"
)

var (
	taxCodePersonalRegexp = regexp.MustCompile(`^\d{9}$`)
	taxCodeBusinessRegexp = regexp.MustCompile(`^[A-Z]\d{8}$`)

	errInvalidFormat   = errors.New("invalid format")
	errInvalidChecksum = errors.New("invalid check digit")
)

func normalizeTaxIdentity(tID *tax.Identity) {
	tax.NormalizeIdentity(tID)
}

func validateTaxIdentity(tID *tax.Identity) error {
	if tID == nil || tID.Code == cbc.CodeEmpty {
		return nil
	}
	if err := validateTaxIdentityCode(tID); err != nil {
		return validation.Errors{
			"code": err,
		}
	}
	return nil
}

func validateTaxIdentityCode(tID *tax.Identity) error {
	code := tID.Code.String()
	switch {
	case taxCodePersonalRegexp.MatchString(code):
		return verifyPersonalCode(tID.Code)
	case taxCodeBusinessRegexp.MatchString(code):
		return verifyBusinessCode(tID.Code)
	default:
		return errInvalidFormat
	}
}

func verifyPersonalCode(code cbc.Code) error {
	// Implement personal code checksum
	_ = code
	return nil
}

func verifyBusinessCode(code cbc.Code) error {
	// Implement business code checksum
	_ = code
	return nil
}
```

---

## 4. Tax Identity Tests — `tax_identity_test.go`

```go
package xx_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/xx"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestValidateTaxIdentity(t *testing.T) {
	tests := []struct {
		name string
		code cbc.Code
		err  string
	}{
		{name: "valid 1", code: "123456789"},
		{name: "valid 2", code: "987654321"},
		{
			name: "too short",
			code: "12345",
			err:  "invalid format",
		},
		{
			name: "too long",
			code: "12345678901234",
			err:  "invalid format",
		},
		{
			name: "bad checksum",
			code: "123456780",
			err:  "checksum mismatch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tID := &tax.Identity{Country: "XX", Code: tt.code}
			err := xx.Validate(tID)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				if assert.Error(t, err) {
					assert.Contains(t, err.Error(), tt.err)
				}
			}
		})
	}
}
```

---

## 5. Scenarios — `scenarios.go`

### Simple (tag-based, like DE)

```go
package xx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// Reverse Charges
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
```

### As a function (like ES)

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

---

## 6. Corrections — `corrections.go`

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

---

## 7. Identities — `identities.go`

### As a variable (like DE)

```go
package xx

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

const (
	IdentityKeyCustom cbc.Key = "xx-custom-id"
)

var identityDefinitions = []*cbc.Definition{
	{
		Key: IdentityKeyCustom,
		Name: i18n.String{
			i18n.EN: "Custom Identity",
			// i18n.XX: "Local name",
		},
	},
}
```

### As a function (like ES)

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
				// i18n.XX: "Local name",
			},
		},
		{
			Key: org.IdentityKeyForeign,
			Name: i18n.String{
				i18n.EN: "Foreign ID",
				// i18n.XX: "Local name",
			},
		},
	}
}
```

---

## 8. Invoice Validation — `invoices.go`

### Simple (like DE)

```go
package xx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateInvoiceSupplier),
			validation.Skip,
		),
	)
}

func validateInvoiceSupplier(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}
```

### With struct validator (like ES)

```go
package xx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

type invoiceValidator struct {
	inv *bill.Invoice
}

func validateInvoice(inv *bill.Invoice) error {
	v := &invoiceValidator{inv: inv}
	return v.validate()
}

func (v *invoiceValidator) validate() error {
	inv := v.inv
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(v.supplier),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) supplier(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}
```

---

## 9. Invoice/Regime Tests — `xx_test.go`

```go
package xx_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validInvoice() *bill.Invoice {
	return &bill.Invoice{
		Regime: tax.WithRegime("XX"),
		Series: "TEST",
		Code:   "0001",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "XX",
				Code:    "123456789",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "XX",
				Code:    "987654321",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
					Unit:  org.UnitPackage,
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "general",
					},
				},
			},
		},
	}
}

func TestInvoiceValidation(t *testing.T) {
	t.Run("valid invoice", func(t *testing.T) {
		inv := validInvoice()
		require.NoError(t, inv.Calculate())
		assert.NoError(t, inv.Validate())
	})
}
```

---

## 10. Custom Tags — `scenarios.go` (tags section)

Custom tags are defined together with scenarios in some regimes:

```go
package xx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

const (
	TagCustom cbc.Key = "custom-tag"
)

func invoiceTags() *tax.TagSet {
	return &tax.TagSet{
		Schema: bill.ShortSchemaInvoice,
		List: []*cbc.Definition{
			{
				Key: TagCustom,
				Name: i18n.String{
					i18n.EN: "Custom Tag",
					// i18n.XX: "Local name",
				},
			},
		},
	}
}
```

---

## 11. Registration in `regimes/regimes.go`

Add a blank import in alphabetical order:

```go
import (
	// ... existing imports ...
	_ "github.com/invopop/gobl/regimes/xx"
	// ... existing imports ...
)
```

---

## 12. Org Identity Normalization & Validation

### Template (for regimes with org identity documents)

```go
package xx

import (
	"fmt"
	"regexp"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	IdentityKeyCustom cbc.Key = "xx-custom-number"
)

var customNumberPattern = regexp.MustCompile(`^\d{3}/\d{3}/\d{5}$`)

var identityDefinitions = []*cbc.Definition{
	{
		Key: IdentityKeyCustom,
		Name: i18n.String{
			i18n.EN: "Custom Number",
			// i18n.XX: "Local name",
		},
	},
}

func normalizeOrgIdentity(id *org.Identity) {
	if id == nil || id.Key != IdentityKeyCustom {
		return
	}
	code := cbc.NormalizeNumericalCode(id.Code).String()
	if len(code) == 11 {
		code = fmt.Sprintf("%s/%s/%s", code[:3], code[3:6], code[6:])
	}
	id.Code = cbc.Code(code)
}

func validateOrgIdentity(id *org.Identity) error {
	if id == nil || id.Key != IdentityKeyCustom {
		return nil
	}
	return validation.ValidateStruct(id,
		validation.Field(&id.Code,
			validation.Required,
			validation.Match(customNumberPattern),
			validation.Skip,
		),
	)
}
```
