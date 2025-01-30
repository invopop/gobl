package saft

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Series and code patterns
const (
	fullCodePattern = "^[^ ]+ [^/^ ]+/[0-9]+$" // extracted from the SAFT-PT XSD to validate the code when the series is not present (e.g. "FT SERIES-A/123")
	seriesPattern   = "^[^ ]+ [^/^ ]+$"        // based on the fullCodePattern, to validate the series when present (e.g. "FT SERIES-A")
	codePattern     = "^[0-9]+$"               // based on the fullCodePattern, to validate the code when the series is present (e.g. "123")
)

// Series and code regexps
var (
	fullCodeRegexp = regexp.MustCompile(fullCodePattern)
	seriesRegexp   = regexp.MustCompile(seriesPattern)
	codeRegexp     = regexp.MustCompile(codePattern)
)

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Tax,
			validation.Required,
			validation.By(validateTax),
			validation.Skip,
		),
		validation.Field(&inv.Series,
			validation.By(validatePrefix(inv)),
			validation.Match(seriesRegexp),
			validation.Skip,
		),
		validation.Field(&inv.Code,
			validation.When(inv.Series != "",
				validation.Match(codeRegexp),
			),
			validation.When(inv.Series == "",
				validation.By(validatePrefix(inv)),
				validation.Match(fullCodeRegexp),
			),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				bill.RequireLineTaxCategory(tax.CategoryVAT),
				validation.Skip,
			),
			validation.Skip,
		),
	)
}

func validateTax(val any) error {
	t, _ := val.(*bill.Tax)
	if t == nil {
		return nil
	}

	return validation.ValidateStruct(t,
		validation.Field(&t.Ext,
			tax.ExtensionsRequire(ExtKeyInvoiceType),
			validation.Skip,
		),
	)
}

func validatePrefix(inv *bill.Invoice) validation.RuleFunc {
	return func(val any) error {
		s, ok := val.(cbc.Code)
		if !ok || s == "" {
			return nil
		}

		if inv == nil || inv.Tax == nil || inv.Tax.Ext == nil || inv.Tax.Ext[ExtKeyInvoiceType] == cbc.CodeEmpty {
			return nil
		}

		prefix := inv.Tax.Ext[ExtKeyInvoiceType].String() + " "
		if !strings.HasPrefix(s.String(), prefix) {
			return fmt.Errorf("must start with '%s'", prefix)
		}

		return nil
	}
}
