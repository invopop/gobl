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
	it := invoiceType(inv)

	return validation.ValidateStruct(inv,
		validation.Field(&inv.Tax,
			validation.By(validateTax),
			validation.Skip,
		),
		validation.Field(&inv.Series,
			validateSeriesFormat(it),
			validation.Skip,
		),
		validation.Field(&inv.Code,
			validateCodeFormat(inv.Series, it),
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
		// If no tax is given, init a blank one so that we can return meaningful
		// validation errors. The blank tax object is not assigned to the invoice
		// and so the original document doesn't actually change.
		t = new(bill.Tax)
	}

	return validation.ValidateStruct(t,
		validation.Field(&t.Ext,
			validation.By(validateTaxExt),
			validation.Skip,
		),
	)
}

func validateTaxExt(val any) error {
	ext, _ := val.(tax.Extensions)
	if ext == nil {
		ext = make(tax.Extensions) // Empty temporary map to return meaningful errors
	}

	msg := fmt.Sprintf("either `%s` or `%s` must be set", ExtKeyWorkType, ExtKeyInvoiceType)

	if !ext.Has(ExtKeyWorkType) && !ext.Has(ExtKeyInvoiceType) {
		return validation.NewError("invalid", msg)
	}

	if ext.Has(ExtKeyWorkType, ExtKeyInvoiceType) {
		return validation.NewError("invalid", msg+", but not both")
	}

	return nil
}

// validateSeriesFormat validates the format of the series to meet the requirements of the
// AT (e.g. "FT SERIES-A"). The series is allowed to be empty, in which case the code is
// expected to be a full code (e.g. "FT SERIES-A/123") (see `validateCodeFormat`).
func validateSeriesFormat(docType cbc.Code) validation.Rule {
	return validation.By(func(val any) error {
		s, ok := val.(cbc.Code)
		if !ok || s == "" {
			return nil
		}

		if docType != cbc.CodeEmpty {
			prefix := docType.String() + " "
			if !strings.HasPrefix(s.String(), prefix) {
				return fmt.Errorf("must start with '%s'", prefix)
			}
		}

		if !seriesRegexp.MatchString(s.String()) {
			return fmt.Errorf("must be in a valid format")
		}

		return nil
	})
}

// validateCodeFormat validates the format of the code to meet the requirements of the
// AT. If the series is present, the code must be a valid number (e.g. 123). If the series
// is missing, the code is expected to be a full code (e.g. "FT SERIES-A/123").
func validateCodeFormat(series cbc.Code, docType cbc.Code) validation.Rule {
	return validation.By(func(val any) error {
		c, ok := val.(cbc.Code)
		if !ok || c == "" {
			return nil
		}

		if series != cbc.CodeEmpty {
			if !codeRegexp.MatchString(c.String()) {
				return fmt.Errorf("must be in a valid format")
			}
			return nil
		}

		if docType != cbc.CodeEmpty {
			prefix := docType.String() + " "
			if !strings.HasPrefix(c.String(), prefix) {
				return fmt.Errorf("must start with '%s'", prefix)
			}
		}

		if !fullCodeRegexp.MatchString(c.String()) {
			return fmt.Errorf("must be in a valid format")
		}
		return nil
	})
}

func invoiceType(inv *bill.Invoice) cbc.Code {
	if inv == nil || inv.Tax == nil || inv.Tax.Ext == nil {
		return cbc.CodeEmpty
	}

	return inv.Tax.Ext[ExtKeyInvoiceType]
}
