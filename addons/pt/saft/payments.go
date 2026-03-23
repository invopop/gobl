package saft

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Add-on custom tags
const (
	TagVATCash cbc.Key = "vat-cash"
)

func billPaymentRules() *rules.Set {
	return rules.For(new(bill.Payment),
		rules.Assert("01", "series format must be valid",
			is.FuncError("series format", paymentSeriesFormatValid),
		),
		rules.Assert("02", "code format must be valid",
			is.FuncError("code format", paymentCodeFormatValid),
		),
		rules.Field("ext",
			rules.Assert("03",
				fmt.Sprintf("'%s' extension is required", ExtKeyPaymentType),
				tax.ExtensionsRequire(ExtKeyPaymentType),
			),
			rules.Assert("04",
				fmt.Sprintf("'%s' extension is required", ExtKeySource),
				tax.ExtensionsRequire(ExtKeySource),
			),
		),
		rules.When(is.Func("source not produced", paymentSourceNotProduced),
			rules.Field("ext",
				rules.Assert("05",
					fmt.Sprintf("'%s' extension is required when source is not produced", ExtKeySourceRef),
					tax.ExtensionsRequire(ExtKeySourceRef),
				),
			),
		),
		rules.Assert("06", "source ref format is invalid",
			is.FuncError("source ref format", paymentSourceRefValid),
		),
		rules.Field("supplier",
			rules.Field("tax_id",
				rules.Assert("07", "supplier tax ID is required", is.Present),
				rules.Field("code",
					rules.Assert("08", "supplier tax ID code is required", is.Present),
				),
			),
		),
		rules.Assert("09", "customer name is required when customer has tax ID code",
			is.Func("customer name present", paymentCustomerNamePresent),
		),
		rules.Field("total",
			rules.Assert("10", "must be no less than 0", num.ZeroOrPositive),
		),
	)
}

func paymentSeriesFormatValid(val any) error {
	pmt, ok := val.(*bill.Payment)
	if !ok || pmt == nil {
		return nil
	}
	return validateSeriesFormat(paymentDocType(pmt)).Validate(pmt.Series)
}

func paymentCodeFormatValid(val any) error {
	pmt, ok := val.(*bill.Payment)
	if !ok || pmt == nil {
		return nil
	}
	dt := paymentDocType(pmt)
	return validateCodeFormat(pmt.Series, dt).Validate(pmt.Code)
}

func paymentSourceNotProduced(val any) bool {
	pmt, ok := val.(*bill.Payment)
	if !ok || pmt == nil {
		return false
	}
	return pmt.Ext != nil && pmt.Ext[ExtKeySource] != "" && pmt.Ext[ExtKeySource] != SourceBillingProduced
}

func paymentSourceRefValid(val any) error {
	pmt, ok := val.(*bill.Payment)
	if !ok || pmt == nil {
		return nil
	}
	return validateSourceRef(paymentDocType(pmt), pmt.Ext)
}

func paymentCustomerNamePresent(val any) bool {
	pmt, ok := val.(*bill.Payment)
	if !ok || pmt == nil || pmt.Customer == nil {
		return true
	}
	if pmt.Customer.TaxID == nil || pmt.Customer.TaxID.Code == cbc.CodeEmpty {
		return true
	}
	return pmt.Customer.Name != ""
}

func paymentDocType(pmt *bill.Payment) cbc.Code {
	if pmt.Ext == nil {
		return cbc.CodeEmpty
	}
	return pmt.Ext[ExtKeyPaymentType]
}

func normalizePayment(pmt *bill.Payment) {
	if pmt.Ext == nil {
		pmt.Ext = tax.Extensions{}
	}

	// TODO: This could be done with scenarios when supported
	if pmt.HasTags(TagVATCash) {
		pmt.Ext[ExtKeyPaymentType] = PaymentTypeCash
	} else {
		pmt.Ext[ExtKeyPaymentType] = PaymentTypeOther
	}

	if !pmt.Ext.Has(ExtKeySource) {
		pmt.Ext[ExtKeySource] = SourceBillingProduced
	}
}

