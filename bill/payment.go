package bill

import (
	"context"
	"fmt"
	"strconv"

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

var isValidPaymentType = validation.In(validPaymentTypes()...)

func validPaymentTypes() []interface{} {
	list := make([]interface{}, len(PaymentTypes))
	for i, d := range PaymentTypes {
		list[i] = d.Key
	}
	return list
}

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
	// Used as a prefix to group codes.
	Series cbc.Code `json:"series,omitempty" jsonschema:"title=Series"`
	// Sequential code used to identify this payment in tax declarations.
	Code cbc.Code `json:"code" jsonschema:"title=Code"`
	// When the payment was issued.
	IssueDate cal.Date `json:"issue_date" jsonschema:"title=Issue Date" jsonschema_extras:"calculated=true"`
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

	// Summary of the taxes applied to the payment for tax regimes that require
	// this information to be communicated.
	Tax *tax.Total `json:"tax,omitempty" jsonschema:"title=Tax"`
	// Total amount to be paid in this payment, either positive or negative according to the
	// line types and totals.
	Total num.Amount `json:"total" jsonschema:"title=Total"`

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
		validation.Field(&pmt.ExchangeRates),
		validation.Field(&pmt.Ext),
		validation.Field(&pmt.Preceding),
		validation.Field(&pmt.Supplier, validation.Required),
		validation.Field(&pmt.Customer),
		validation.Field(&pmt.Payee),
		validation.Field(&pmt.Lines, validation.Required),
		validation.Field(&pmt.Ordering),
		validation.Field(&pmt.Tax),
		validation.Field(&pmt.Total, validation.Required),
		validation.Field(&pmt.Notes),
		validation.Field(&pmt.Complements),
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

	// Convert empty or invalid currency to the regime's currency
	if pmt.Currency == currency.CodeEmpty && r != nil {
		pmt.Currency = r.Currency
	}
	if pmt.Currency == currency.CodeEmpty {
		return validation.Errors{
			"currency": fmt.Errorf("required, unable to determine"),
		}
	}

	for i, l := range pmt.Lines {
		l.Index = i + 1
		if err := l.calculate(pmt.Currency, pmt.ExchangeRates); err != nil {
			return validation.Errors{
				"lines": validation.Errors{
					strconv.Itoa(i): err,
				},
			}
		}

		l.Tax.Calculate(pmt.Currency, r.GetRoundingRule())
		lt := l.Tax.Clone()
		a := l.Total

		// Merge the line taxes
		if l.Tax != nil {
			if tt == nil {
				tt = lt
			} else {
				tt = tt.Merge(lt)
			}
		}

		// Finally add the totals
		if total == nil {
			total = &a
		} else {
			nt := total.Add(a)
			total = &nt
		}
	}
	pmt.Tax = tt
	if total != nil {
		pmt.Total = *total
	}
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
		},
	}
}
