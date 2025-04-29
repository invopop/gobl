package saft

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
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

var invoiceWorkTypes = []cbc.Code{
	WorkTypeProforma,
	WorkTypeConsignmentInv,
	WorkTypeConsignmentCredit,
}

func validateInvoice(inv *bill.Invoice) error {
	dt := invoiceDocType(inv)

	return validation.ValidateStruct(inv,
		validation.Field(&inv.Tax,
			validation.By(validateTax),
			validation.Skip,
		),
		validation.Field(&inv.Series,
			validateSeriesFormat(dt),
			validation.Skip,
		),
		validation.Field(&inv.Code,
			validateCodeFormat(inv.Series, dt),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				bill.RequireLineTaxCategory(tax.CategoryVAT),
				validation.Skip,
			),
			validation.Skip,
		),
		validation.Field(&inv.Payment,
			validation.By(validatePaymentDetails(inv)),
			validation.Skip,
		),
		validation.Field(&inv.Totals,
			validation.By(validateTotals(inv)),
			validation.Skip,
		),
	)
}

func validatePaymentDetails(inv *bill.Invoice) validation.RuleFunc {
	return func(val any) error {
		pay, _ := val.(*bill.PaymentDetails)
		if pay == nil {
			return nil
		}

		return validation.ValidateStruct(pay,
			validation.Field(&pay.Advances,
				validation.Each(
					validation.By(validateAdvance(inv)),
					validation.Skip,
				),
				validation.Skip,
			),
		)
	}
}

func validateAdvance(inv *bill.Invoice) validation.RuleFunc {
	return func(val any) error {
		adv, _ := val.(*pay.Advance)
		if adv == nil {
			return nil
		}

		return validation.ValidateStruct(adv,
			validation.Field(&adv.Date,
				validation.Required,
				validation.In(inv.IssueDate).Error("must be the same as the invoice issue date"),
				validation.Skip,
			),
		)
	}
}

func validateTotals(inv *bill.Invoice) validation.RuleFunc {
	return func(val any) error {
		tot, _ := val.(*bill.Totals)
		if tot == nil {
			return nil
		}

		return validation.ValidateStruct(tot,
			validation.Field(&tot.Due,
				validation.When(isInvoiceReceipt(inv), num.Max(num.AmountZero)),
				validation.Skip,
			),
		)
	}
}

func invoiceDocType(inv *bill.Invoice) cbc.Code {
	if inv.Tax == nil || inv.Tax.Ext == nil {
		return cbc.CodeEmpty
	}
	if inv.Tax.Ext.Has(ExtKeyInvoiceType) {
		return inv.Tax.Ext[ExtKeyInvoiceType]
	}
	return inv.Tax.Ext[ExtKeyWorkType]
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

	if wt, ok := ext[ExtKeyWorkType]; ok {
		if !slices.Contains(invoiceWorkTypes, wt) {
			return validation.Errors{
				ExtKeyWorkType.String(): fmt.Errorf("value '%s' invalid", wt),
			}
		}
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

func isInvoiceReceipt(inv *bill.Invoice) bool {
	return invoiceDocType(inv) == InvoiceTypeInvoiceReceipt
}

func normalizeInvoice(inv *bill.Invoice) {
	if inv.Payment == nil {
		return
	}

	// Set the issue date as the default date for advances
	for _, adv := range inv.Payment.Advances {
		if adv.Date == nil {
			date := inv.IssueDate
			adv.Date = &date
		}
	}
}
