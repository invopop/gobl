package saft

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Series and code patterns
const (
	seriesPattern   = "^[^ ]+ [^/^ ]+$"
	codePattern     = "^[0-9]+$"
	fullCodePattern = "^[^ ]+ [^/^ ]+/[0-9]+$"
)

// Series and code regexps
var (
	seriesRegexp   = regexp.MustCompile(seriesPattern)
	codeRegexp     = regexp.MustCompile(codePattern)
	fullCodeRegexp = regexp.MustCompile(fullCodePattern)
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
		validation.Field(&inv.Currency,
			validation.In(currency.EUR).Error("must be EUR"),
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
	t, ok := val.(*bill.Tax)
	if !ok {
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
