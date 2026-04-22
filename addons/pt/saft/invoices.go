package saft

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

var corrections = tax.CorrectionSet{
	{
		Schema: bill.ShortSchemaInvoice,
		Types: []cbc.Key{
			bill.InvoiceTypeCreditNote,
			bill.InvoiceTypeDebitNote,
		},
		ReasonRequired: true,
	},
}

// Series and code patterns
const (
	fullCodePattern  = "^[^ ]+ [^/^ ]+/[0-9]+$" // extracted from the SAFT-PT XSD to validate the code when the series is not present (e.g. "FT SERIES-A/123")
	seriesPattern    = "^[^ ]+ [^/^ ]+$"        // based on the fullCodePattern, to validate the series when present (e.g. "FT SERIES-A")
	codePattern      = "^[0-9]+$"               // based on the fullCodePattern, to validate the code when the series is present (e.g. "123")
	sourceRefPattern = "^([^ ]+)(?:M|D ([^ ]+)) [^/^ ]+/[$0-9]+$"
)

// Series and code regexps
var (
	fullCodeRegexp  = regexp.MustCompile(fullCodePattern)
	seriesRegexp    = regexp.MustCompile(seriesPattern)
	codeRegexp      = regexp.MustCompile(codePattern)
	sourceRefRegexp = regexp.MustCompile(sourceRefPattern)
)

var invoiceWorkTypes = []cbc.Code{
	WorkTypeProforma,
	WorkTypeConsignmentInv,
	WorkTypeConsignmentCredit,
}

func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}

	normalizeInvoiceTax(inv)
	normalizeInvoiceAdvances(inv)
	normalizeInvoiceValueDate(inv)
}

func normalizeInvoiceTax(inv *bill.Invoice) {
	if inv.Tax == nil {
		inv.Tax = new(bill.Tax)
	}

	if inv.Tax.Ext == nil {
		inv.Tax.Ext = make(tax.Extensions)
	}

	if !inv.Tax.Ext.Has(ExtKeySource) {
		inv.Tax.Ext[ExtKeySource] = SourceBillingProduced
	}
}

func normalizeInvoiceAdvances(inv *bill.Invoice) {
	if inv.Payment == nil {
		return
	}

	// Set the issue date as the default date for advances
	for _, adv := range inv.Payment.Advances {
		if adv.Date == nil {
			adv.Date = issueDate(inv)
		}
	}
}

func normalizeInvoiceValueDate(inv *bill.Invoice) {
	inv.ValueDate = determineValueDate(
		issueDate(inv),
		inv.OperationDate,
		inv.ValueDate,
	)
}

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Assert("15", "invoice must be in EUR or provide exchange rate for conversion", currency.CanConvertTo(currency.EUR)),
		// Tax ext must have either work type or invoice type
		rules.Assert("01",
			fmt.Sprintf("either '%s' or '%s' must be set", ExtKeyWorkType, ExtKeyInvoiceType),
			is.Func("has doc type", invoiceHasDocType),
		),
		// Tax ext must not have both work type and invoice type
		rules.Assert("02",
			fmt.Sprintf("either '%s' or '%s' must be set, but not both", ExtKeyWorkType, ExtKeyInvoiceType),
			is.Func("not both doc types", invoiceNotBothDocTypes),
		),
		// If work type is set, it must be a valid invoice work type
		rules.Assert("03", "invoice work type is not valid",
			is.FuncError("valid invoice work type", invoiceWorkTypeValid),
		),
		// Series format depends on doc type
		rules.Assert("04", "series format must be valid",
			is.FuncError("series format", invoiceSeriesFormatValid),
		),
		// Code format depends on series and doc type
		rules.Assert("05", "code format must be valid",
			is.FuncError("code format", invoiceCodeFormatValid),
		),
		rules.Field("value_date",
			rules.Assert("06", "cannot be blank", is.Present),
		),
		// Lines need VAT category
		rules.Field("lines",
			rules.Each(
				rules.Assert("07", "line taxes must include VAT category",
					bill.RequireLineTaxCategory(tax.CategoryVAT),
				),
			),
		),
		rules.Field("payment",
			rules.Field("advances",
				rules.Each(
					rules.Field("date",
						rules.Assert("08", "cannot be blank", is.Present),
					),
				),
			),
		),
		rules.Field("totals",
			rules.Field("payable",
				rules.Assert("09", "must be no less than 0", num.ZeroOrPositive),
			),
		),
		// Due must be zero for invoice-receipt
		rules.When(is.Func("is invoice-receipt", invoiceIsReceipt),
			rules.Field("totals",
				rules.Field("due",
					rules.Assert("10", "must be equal to 0", num.Equals(num.AmountZero)),
				),
			),
		),
		rules.Field("preceding",
			rules.Assert("11", "the length must be no more than 1", is.Length(0, 1)),
		),
		// Tax ext requires source
		rules.Assert("12", fmt.Sprintf("tax requires '%s' extension", ExtKeySource),
			is.Func("tax has source", invoiceTaxHasSource),
		),
		// Tax ext requires sourceRef when source != produced
		rules.When(is.Func("source not produced", invoiceSourceNotProduced),
			rules.Assert("13", fmt.Sprintf("tax requires '%s' extension when source is not produced", ExtKeySourceRef),
				is.Func("tax has source ref", invoiceTaxHasSourceRef),
			),
		),
		// Source ref format validation
		rules.Assert("14", "source ref format is invalid",
			is.FuncError("source ref format", invoiceSourceRefValid),
		),
	)
}

// invoiceHasDocType checks that either work type or invoice type is set.
func invoiceHasDocType(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	ext := invoiceTaxExt(inv)
	return ext.Has(ExtKeyWorkType) || ext.Has(ExtKeyInvoiceType)
}

// invoiceNotBothDocTypes checks that work type and invoice type are not both set.
func invoiceNotBothDocTypes(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	ext := invoiceTaxExt(inv)
	return !ext.Has(ExtKeyWorkType, ExtKeyInvoiceType)
}

// invoiceWorkTypeValid checks that if the work type is present, it's a valid invoice work type.
func invoiceWorkTypeValid(val any) error {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return nil
	}
	ext := invoiceTaxExt(inv)
	wt, ok := ext[ExtKeyWorkType]
	if !ok {
		return nil
	}
	if !slices.Contains(invoiceWorkTypes, wt) {
		return fmt.Errorf("value '%s' invalid", wt)
	}
	return nil
}

// invoiceSeriesFormatValid validates the series format against the doc type.
func invoiceSeriesFormatValid(val any) error {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return nil
	}
	return validateSeriesFormat(invoiceDocType(inv)).Validate(inv.Series)
}

// invoiceCodeFormatValid validates the code format against the series and doc type.
func invoiceCodeFormatValid(val any) error {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return nil
	}
	return validateCodeFormat(inv.Series, invoiceDocType(inv)).Validate(inv.Code)
}

// invoiceIsReceipt returns true if the invoice is an invoice-receipt.
func invoiceIsReceipt(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && isInvoiceReceipt(inv)
}

// invoiceTaxHasSource checks that the tax extensions include the source key.
func invoiceTaxHasSource(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Tax == nil {
		return false
	}
	return tax.ExtensionsRequire(ExtKeySource).Check(inv.Tax.Ext)
}

// invoiceSourceNotProduced returns true when the invoice source is not "produced".
func invoiceSourceNotProduced(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	ext := invoiceTaxExt(inv)
	return ext[ExtKeySource] != "" && ext[ExtKeySource] != SourceBillingProduced
}

// invoiceTaxHasSourceRef checks that the tax extensions include the source ref key.
func invoiceTaxHasSourceRef(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Tax == nil {
		return false
	}
	return tax.ExtensionsRequire(ExtKeySourceRef).Check(inv.Tax.Ext)
}

// invoiceSourceRefValid validates the source ref format for the invoice.
func invoiceSourceRefValid(val any) error {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return nil
	}
	ext := invoiceTaxExt(inv)
	return validateSourceRef(invoiceDocType(inv), ext)
}

// invoiceTaxExt safely returns the invoice's tax extensions.
func invoiceTaxExt(inv *bill.Invoice) tax.Extensions {
	if inv.Tax == nil || inv.Tax.Ext == nil {
		return nil
	}
	return inv.Tax.Ext
}

func validateSourceRef(docType cbc.Code, ext tax.Extensions) error {
	if ext == nil {
		return nil
	}

	if ext[ExtKeySource] != SourceBillingManual {
		// source ref format only validated for manual documents
		return nil
	}

	ref := ext[ExtKeySourceRef].String()
	if ref == "" || docType == "" {
		return nil
	}

	matches := sourceRefRegexp.FindStringSubmatch(ref)
	if len(matches) == 0 {
		return errors.New("must be in valid format")
	}
	if matches[1] != docType.String() {
		return fmt.Errorf("must start with the document type '%s' not '%s'", docType, matches[1])
	}
	if matches[2] != "" && matches[2] != docType.String() {
		return fmt.Errorf("must refer to an original document '%s' not '%s'", docType, matches[2])
	}

	return nil
}

// validateSeriesFormat validates the format of the series to meet the requirements of the
// AT (e.g. "FT SERIES-A"). The series is allowed to be empty, in which case the code is
// expected to be a full code (e.g. "FT SERIES-A/123") (see `validateCodeFormat`).
func validateSeriesFormat(docType cbc.Code) seriesFormatRule {
	return seriesFormatRule{docType: docType}
}

type seriesFormatRule struct {
	docType cbc.Code
}

func (r seriesFormatRule) Validate(val any) error {
	s, ok := val.(cbc.Code)
	if !ok || s == "" {
		return nil
	}

	if r.docType != cbc.CodeEmpty {
		prefix := r.docType.String() + " "
		if !strings.HasPrefix(s.String(), prefix) {
			return fmt.Errorf("must start with '%s'", prefix)
		}
	}

	if !seriesRegexp.MatchString(s.String()) {
		return fmt.Errorf("must be in a valid format")
	}

	return nil
}

// validateCodeFormat validates the format of the code to meet the requirements of the
// AT. If the series is present, the code must be a valid number (e.g. 123). If the series
// is missing, the code is expected to be a full code (e.g. "FT SERIES-A/123").
func validateCodeFormat(series cbc.Code, docType cbc.Code) codeFormatRule {
	return codeFormatRule{series: series, docType: docType}
}

type codeFormatRule struct {
	series  cbc.Code
	docType cbc.Code
}

func (r codeFormatRule) Validate(val any) error {
	c, ok := val.(cbc.Code)
	if !ok || c == "" {
		return nil
	}

	if r.series != cbc.CodeEmpty {
		if !codeRegexp.MatchString(c.String()) {
			return fmt.Errorf("must be in a valid format")
		}
		return nil
	}

	if r.docType != cbc.CodeEmpty {
		prefix := r.docType.String() + " "
		if !strings.HasPrefix(c.String(), prefix) {
			return fmt.Errorf("must start with '%s'", prefix)
		}
	}

	if !fullCodeRegexp.MatchString(c.String()) {
		return fmt.Errorf("must be in a valid format")
	}
	return nil
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

func isInvoiceReceipt(inv *bill.Invoice) bool {
	return invoiceDocType(inv) == InvoiceTypeInvoiceReceipt
}

func determineValueDate(idate, odate, vdate *cal.Date) *cal.Date {
	if vdate != nil {
		return vdate
	}
	if odate != nil {
		return odate
	}
	return idate
}

func issueDate(inv *bill.Invoice) *cal.Date {
	return dateOrToday(&inv.IssueDate, inv.Regime)
}

func dateOrToday(date *cal.Date, reg tax.Regime) *cal.Date {
	if date != nil && !date.IsZero() {
		return date
	}

	rd := reg.RegimeDef()
	loc := rd.TimeLocation()
	today := cal.TodayIn(loc)
	return &today
}
