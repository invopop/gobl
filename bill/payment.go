package bill

import (
	"context"
	"fmt"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/internal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Predefined list of the payment types supported.
const (
	PaymentTypeRequest cbc.Key = "request"
	PaymentTypeReceipt cbc.Key = "receipt"
	PaymentTypeAdvice  cbc.Key = "advice"
)

// PaymentTypes defines the list of potential payment types.
var PaymentTypes = []*cbc.Definition{
	{
		Key: PaymentTypeRequest,
		Name: i18n.String{
			i18n.EN: "Request",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				A payment request sent from the supplier to a customer indicating that they are
				requesting a transfer of funds from the customer directly or a payer.
				This is used to request payment for specific documents and invoices.
			`),
		},
	},
	{
		Key: PaymentTypeAdvice,
		Name: i18n.String{
			i18n.EN: "Advice",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				A remittance advice sent from the customer to the supplier reflecting that payment for
				the referenced documents has been made.
			`),
		},
	},
	{
		Key: PaymentTypeReceipt,
		Name: i18n.String{
			i18n.EN: "Receipt",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				A payment receipt sent from the supplier to a customer indicating that they have
				received a transfer of funds from the customer directly or a payer.
				This is the default payment type and may be required by some tax
				regimes in order to communicate the payment of specific documents and invoices.
			`),
		},
	},
}

var isValidPaymentType = cbc.InKeyDefs(PaymentTypes)

// A Payment is used to link an invoice or invoices with a payment transaction.
type Payment struct {
	tax.Regime
	tax.Addons
	tax.Tags
	uuid.Identify

	// Type of payment document being issued.
	Type cbc.Key `json:"type" jsonschema:"title=Type" jsonschema_extras:"calculated=true"`
	// Details on how the payment was made based on the original instructions.
	Method *pay.Instructions `json:"method,omitempty" jsonschema:"title=Method"`
	// Series is used to identify groups of payments by date, business area, project,
	// type, customer, a combination of any, or other company specific data.
	// If the output format does not support the series as a separate field, it will be
	// prepended to the code for presentation with a dash (`-`) for separation.
	Series cbc.Code `json:"series,omitempty" jsonschema:"title=Series"`
	// Code is a sequential identifier that uniquely identifies the payment. The code can
	// be left empty initially, but is **required** to **sign** the document.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// When the payment was issued.
	IssueDate cal.Date `json:"issue_date" jsonschema:"title=Issue Date" jsonschema_extras:"calculated=true"`
	// IssueTime is an optional field that may be useful to indicate the time of day when
	// the payment was issued.
	IssueTime *cal.Time `json:"issue_time,omitempty" jsonschema:"title=Issue Time" jsonschema_extras:"calculated=true"`
	// When the taxes of this payment become accountable, if none set, the issue date is assumed.
	ValueDate *cal.Date `json:"value_date,omitempty" jsonschema:"title=Value Date"`
	// Currency for all payment totals.
	Currency currency.Code `json:"currency" jsonschema:"title=Currency" jsonschema_extras:"calculated=true"`
	// Exchange rates to be used when converting the payment's monetary values into other currencies.
	ExchangeRates []*currency.ExchangeRate `json:"exchange_rates,omitempty" jsonschema:"title=Exchange Rates"`
	// Extensions for additional codes that may be required.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`

	// Key information regarding previous versions of this document.
	Preceding []*org.DocumentRef `json:"preceding,omitempty" jsonschema:"title=Preceding Details"`

	// The taxable entity who is responsible for supplying goods or services.
	Supplier *org.Party `json:"supplier" jsonschema:"title=Supplier"`
	// Legal entity that receives the goods or services.
	Customer *org.Party `json:"customer,omitempty" jsonschema:"title=Customer"`
	// Legal entity that receives the payment if not the supplier.
	Payee *org.Party `json:"payee,omitempty" jsonschema:"title=Payee"`

	// List of documents that are being paid for.
	Lines []*PaymentLine `json:"lines" jsonschema:"title=Lines"`

	// Ordering allows for additional information about the ordering process including references
	// to other documents and alternative parties involved in the order-to-delivery process.
	Ordering *Ordering `json:"ordering,omitempty" jsonschema:"title=Ordering"`

	// Total advances from any advances registered in the lines. Calculated automatically.
	Advances *num.Amount `json:"advances,omitempty" jsonschema:"title=Advances,calculated=true"`
	// Total amount to be paid in this payment, either positive or negative according to the
	// line types and totals. Calculated automatically.
	Total num.Amount `json:"total" jsonschema:"title=Total,calculated=true"`
	// Due reflects the amount that is still to be paid and will be calculated automatically
	// based on the total and advance amounts, and may be require in some tax regimes or addons.
	Due *num.Amount `json:"due,omitempty" jsonschema:"title=Due,calculated=true"`

	// Summary of the taxes applied to the payment for tax regimes that require
	// this information to be communicated. If payment lines contain `payable` amounts,
	// these will be used to calculate the proportional amount of tax to apply automatically,
	// otherwise the taxes will be added together, assuming 100% is paid.
	Tax *tax.Total `json:"tax,omitempty" jsonschema:"title=Tax,calculated=true"`

	// Unstructured information that is relevant to the payment, such as correction or additional
	// legal details.
	Notes []*org.Note `json:"notes,omitempty" jsonschema:"title=Notes"`

	// Additional complementary objects that add relevant information to the payment.
	Complements []*schema.Object `json:"complements,omitempty" jsonschema:"title=Complements"`

	// Additional semi-structured data that doesn't fit into the body of the invoice.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate runs the validation rules for the payment without the context.
func (pmt *Payment) Validate() error {
	return pmt.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures that the fields contained in the Payment look correct.
func (pmt *Payment) ValidateWithContext(ctx context.Context) error {
	ctx = pmt.validationContext(ctx)
	r := pmt.RegimeDef()
	return tax.ValidateStructWithContext(ctx, pmt,
		validation.Field(&pmt.Regime),
		validation.Field(&pmt.Addons),
		validation.Field(&pmt.UUID),
		validation.Field(&pmt.Type,
			validation.Required,
			isValidPaymentType,
		),
		validation.Field(&pmt.Method, validation.Required),
		validation.Field(&pmt.Series),
		validation.Field(&pmt.Code,
			validation.When(
				internal.IsSigned(ctx),
				validation.Required.Error("required to sign payment"),
			),
		),
		validation.Field(&pmt.IssueDate,
			validation.Required,
			cal.DateNotZero(),
		),
		validation.Field(&pmt.Currency,
			validation.Required,
			currency.CanConvertInto(pmt.ExchangeRates, r.GetCurrency()),
		),
		validation.Field(&pmt.ExchangeRates,
			validation.Each(validation.NotNil),
		),
		validation.Field(&pmt.Ext),
		validation.Field(&pmt.Preceding,
			validation.Each(validation.NotNil),
		),
		validation.Field(&pmt.Supplier, validation.Required),
		validation.Field(&pmt.Customer),
		validation.Field(&pmt.Payee),
		validation.Field(&pmt.Lines,
			validation.Required,
			validation.Each(validation.NotNil),
		),
		validation.Field(&pmt.Ordering),
		validation.Field(&pmt.Tax),
		validation.Field(&pmt.Total, validation.Required),
		validation.Field(&pmt.Notes,
			validation.Each(validation.NotNil),
		),
		validation.Field(&pmt.Complements,
			validation.Each(validation.NotNil),
		),
		validation.Field(&pmt.Meta),
	)
}

// validationContext builds a context with all the validators that the payment might
// need for execution.
func (pmt *Payment) validationContext(ctx context.Context) context.Context {
	if r := pmt.RegimeDef(); r != nil {
		ctx = r.WithContext(ctx)
	}
	for _, a := range pmt.AddonDefs() {
		ctx = a.WithContext(ctx)
	}
	return ctx
}

// Calculate performs all the normalizations and calculations required for the invoice
// totals and taxes. If the original invoice only includes partial calculations, this
// will figure out what's missing.
func (pmt *Payment) Calculate() error {
	// Try to set Regime if not already prepared from the supplier's tax ID
	if pmt.Regime.IsEmpty() {
		pmt.SetRegime(partyTaxCountry(pmt.Supplier))
	}
	pmt.Normalize(pmt.normalizers())
	return pmt.calculate()
}

// Normalize is run as part of the Calculate method to ensure that the invoice
// is in a consistent state before calculations are performed. This will leverage
// any add-ons alongside the tax regime.
func (pmt *Payment) Normalize(normalizers tax.Normalizers) {
	if pmt.Type == cbc.KeyEmpty {
		pmt.Type = PaymentTypeReceipt
	}
	pmt.Series = cbc.NormalizeCode(pmt.Series)
	pmt.Code = cbc.NormalizeCode(pmt.Code)

	normalizers.Each(pmt)

	tax.Normalize(normalizers, pmt.Method)
	tax.Normalize(normalizers, pmt.Tax)
	tax.Normalize(normalizers, pmt.Supplier)
	tax.Normalize(normalizers, pmt.Customer)
	tax.Normalize(normalizers, pmt.Preceding)
	tax.Normalize(normalizers, pmt.Lines)
	tax.Normalize(normalizers, pmt.Ordering)
}

func (pmt *Payment) normalizers() tax.Normalizers {
	normalizers := make(tax.Normalizers, 0)
	if r := pmt.RegimeDef(); r != nil {
		normalizers = normalizers.Append(r.Normalizer)
	}
	for _, a := range pmt.AddonDefs() {
		normalizers = normalizers.Append(a.Normalizer)
	}
	return normalizers
}

func (pmt *Payment) calculate() error {
	var tt *tax.Total
	var total *num.Amount

	r := pmt.RegimeDef()

	// Set the issue date and time
	tz := r.TimeLocation()
	if pmt.IssueTime != nil && pmt.IssueTime.IsZero() {
		// If setting the time, also set the date
		tn := cal.ThisSecondIn(tz)
		hn := tn.Time()
		pmt.IssueDate = tn.Date()
		pmt.IssueTime = &hn
	} else if pmt.IssueDate.IsZero() {
		pmt.IssueDate = cal.TodayIn(tz)
	}

	// Convert empty or invalid currency to the regime's currency
	if pmt.Currency == currency.CodeEmpty && r != nil {
		pmt.Currency = r.Currency
	}
	if pmt.Currency == currency.CodeEmpty {
		return validation.Errors{
			"currency": fmt.Errorf("required, unable to determine"),
		}
	}

	due := pmt.Currency.Def().Zero()
	adv := pmt.Currency.Def().Zero()
	for i, l := range pmt.Lines {
		if l == nil {
			continue
		}
		l.Index = i + 1

		var lt *tax.Total
		if l.Document != nil {
			var er *currency.ExchangeRate
			cur := l.Document.Currency
			if cur == currency.CodeEmpty {
				cur = pmt.Currency
			} else {
				// If the document has a currency, we need to ensure there is an exchange
				// rate so any taxes can be converted correctly.
				if er = currency.MatchExchangeRate(pmt.ExchangeRates, cur, pmt.Currency); er == nil {
					return validation.Errors{
						"exchange_rates": fmt.Errorf("%s to %s missing", cur, pmt.Currency),
					}
				}
			}
			rr := r.GetRoundingRule()
			l.Document.Calculate(cur, rr)
			lt = l.Document.Tax.Clone()
			lt.Exchange(er, rr)

			// Perform extra calculations with the payable amount, if present.
			if p := l.Document.Payable; p != nil {
				// When calculating the taxes, determine if we need to rescale
				// so that the taxes are proportional to the amount paid.
				factor := l.Amount.Upscale(2).Divide(*p)
				lt.Scale(factor, pmt.Currency, rr)
				due = due.Add(*l.Document.Payable)
				if l.Advances != nil {
					due = due.Subtract(*l.Advances)
				}
				due = due.Subtract(l.Amount)
			}

			// Merge the line document taxes
			if l.Document.Tax != nil {
				if tt == nil {
					tt = lt
				} else {
					tt = tt.Merge(lt)
				}
			}
		}

		// Finally add the totals
		a := l.Amount
		if total == nil {
			total = &a
		} else {
			nt := total.Add(a)
			total = &nt
		}
		if l.Advances != nil {
			adv = adv.Add(*l.Advances)
		}
	}
	if total != nil {
		pmt.Total = *total
	}
	if adv.IsZero() {
		pmt.Advances = nil
	} else {
		pmt.Advances = &adv
	}
	if due.IsZero() {
		pmt.Due = nil
	} else {
		pmt.Due = &due
	}
	pmt.Tax = tt

	return nil
}

// JSONSchemaExtend extends the schema with additional property details
func (pmt Payment) JSONSchemaExtend(js *jsonschema.Schema) {
	props := js.Properties
	// Extend type list
	if its, ok := props.Get("type"); ok {
		its.OneOf = make([]*jsonschema.Schema, len(PaymentTypes))
		for i, kd := range PaymentTypes {
			its.OneOf[i] = &jsonschema.Schema{
				Const:       kd.Key.String(),
				Title:       kd.Name.String(),
				Description: kd.Desc.String(),
			}
		}
	}
	pmt.Regime.JSONSchemaExtend(js)
	pmt.Addons.JSONSchemaExtend(js)
	// Recommendations
	js.Extras = map[string]any{
		schema.Recommended: []string{
			"$regime",
			"series",
			"code",
			"lines",
		},
	}
}
