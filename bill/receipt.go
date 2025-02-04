package bill

import (
	"context"
	"fmt"
	"strconv"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
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

// Predefined list of the receipt types supported.
const (
	ReceiptTypePayment    cbc.Key = "payment"
	ReceiptTypeRemittance cbc.Key = "remittance"
)

// ReceiptTypes defines the list of potential payment types.
var ReceiptTypes = []*cbc.Definition{
	{
		Key: ReceiptTypePayment,
		Name: i18n.String{
			i18n.EN: "Payment",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				A payment receipt sent from the supplier to a customer indicating that they have
				received a transfer of funds from the customer directly or a payer.
				This is the default receipt type and may be required by some tax
				regimes in order to communicate the payment of specific documents and invoices.
			`),
		},
	},
	{
		Key: ReceiptTypeRemittance,
		Name: i18n.String{
			i18n.EN: "Remittance",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				A remittance advice sent from the customer to the supplier reflecting that payment for
				the referenced documents has been made.
			`),
		},
	},
}

var isValidReceiptType = validation.In(validReceiptTypes()...)

func validReceiptTypes() []interface{} {
	list := make([]interface{}, len(ReceiptTypes))
	for i, d := range ReceiptTypes {
		list[i] = d.Key
	}
	return list
}

// A Receipt is used to link an invoice or invoices with a payment transaction.
type Receipt struct {
	tax.Regime
	tax.Addons
	tax.Tags
	uuid.Identify

	// Type of receipt document being issued.
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
	Lines []*ReceiptLine `json:"lines" jsonschema:"title=Lines"`

	// Ordering allows for additional information about the ordering process including references
	// to other documents and alternative parties involved in the order-to-delivery process.
	Ordering *Ordering `json:"ordering,omitempty" jsonschema:"title=Ordering"`

	// Summary of the taxes applied to the payment for tax regimes that require
	// this information to be communicated.
	Tax *tax.Total `json:"tax,omitempty" jsonschema:"title=Tax"`
	// Total amount to be paid in this receipt, either positive or negative according to the
	// line types and totals.
	Total num.Amount `json:"total" jsonschema:"title=Total"`

	// Unstructured information that is relevant to the receipt, such as correction or additional
	// legal details.
	Notes []*org.Note `json:"notes,omitempty" jsonschema:"title=Notes"`

	// Additional complementary objects that add relevant information to the receipt.
	Complements []*schema.Object `json:"complements,omitempty" jsonschema:"title=Complements"`

	// Additional semi-structured data that doesn't fit into the body of the invoice.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate runs the validation rules for the receipt without the context.
func (rct *Receipt) Validate() error {
	return rct.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures that the fields contained in the Receipt look correct.
func (rct *Receipt) ValidateWithContext(ctx context.Context) error {
	ctx = rct.validationContext(ctx)
	r := rct.RegimeDef()
	return tax.ValidateStructWithContext(ctx, rct,
		validation.Field(&rct.Regime),
		validation.Field(&rct.Addons),
		validation.Field(&rct.UUID),
		validation.Field(&rct.Type,
			validation.Required,
			isValidReceiptType,
		),
		validation.Field(&rct.Method, validation.Required),
		validation.Field(&rct.Series),
		validation.Field(&rct.Code, validation.Required),
		validation.Field(&rct.IssueDate,
			validation.Required,
			cal.DateNotZero(),
		),
		validation.Field(&rct.Currency,
			validation.Required,
			currency.CanConvertInto(rct.ExchangeRates, r.GetCurrency()),
		),
		validation.Field(&rct.ExchangeRates),
		validation.Field(&rct.Preceding),
		validation.Field(&rct.Supplier, validation.Required),
		validation.Field(&rct.Customer),
		validation.Field(&rct.Payee),
		validation.Field(&rct.Lines, validation.Required),
		validation.Field(&rct.Ordering),
		validation.Field(&rct.Tax),
		validation.Field(&rct.Total, validation.Required),
		validation.Field(&rct.Notes),
		validation.Field(&rct.Complements),
		validation.Field(&rct.Meta),
	)
}

// validationContext builds a context with all the validators that the receipt might
// need for execution.
func (rct *Receipt) validationContext(ctx context.Context) context.Context {
	if r := rct.RegimeDef(); r != nil {
		ctx = r.WithContext(ctx)
	}
	for _, a := range rct.AddonDefs() {
		ctx = a.WithContext(ctx)
	}
	return ctx
}

// Calculate performs all the normalizations and calculations required for the invoice
// totals and taxes. If the original invoice only includes partial calculations, this
// will figure out what's missing.
func (rct *Receipt) Calculate() error {
	// Try to set Regime if not already prepared from the supplier's tax ID
	if rct.Regime.IsEmpty() {
		rct.SetRegime(partyTaxCountry(rct.Supplier))
	}
	rct.Normalize(rct.normalizers())
	return rct.calculate()
}

// Normalize is run as part of the Calculate method to ensure that the invoice
// is in a consistent state before calculations are performed. This will leverage
// any add-ons alongside the tax regime.
func (rct *Receipt) Normalize(normalizers tax.Normalizers) {
	if rct.Type == cbc.KeyEmpty {
		rct.Type = ReceiptTypePayment
	}
	rct.Series = cbc.NormalizeCode(rct.Series)
	rct.Code = cbc.NormalizeCode(rct.Code)

	normalizers.Each(rct)

	tax.Normalize(normalizers, rct.Method)
	tax.Normalize(normalizers, rct.Tax)
	tax.Normalize(normalizers, rct.Supplier)
	tax.Normalize(normalizers, rct.Customer)
	tax.Normalize(normalizers, rct.Preceding)
	tax.Normalize(normalizers, rct.Lines)
	tax.Normalize(normalizers, rct.Ordering)
}

func (rct *Receipt) normalizers() tax.Normalizers {
	normalizers := make(tax.Normalizers, 0)
	if r := rct.RegimeDef(); r != nil {
		normalizers = normalizers.Append(r.Normalizer)
	}
	for _, a := range rct.AddonDefs() {
		normalizers = normalizers.Append(a.Normalizer)
	}
	return normalizers
}

func (rct *Receipt) calculate() error {
	var tt *tax.Total
	var total *num.Amount

	r := rct.RegimeDef()

	// Convert empty or invalid currency to the regime's currency
	if rct.Currency == currency.CodeEmpty && r != nil {
		rct.Currency = r.Currency
	}
	if rct.Currency == currency.CodeEmpty {
		return validation.Errors{
			"currency": fmt.Errorf("required, unable to determine"),
		}
	}

	for i, l := range rct.Lines {
		if err := l.calculate(rct.Currency, rct.ExchangeRates); err != nil {
			return validation.Errors{
				"lines": validation.Errors{
					strconv.Itoa(i): err,
				},
			}
		}

		l.Tax.Calculate(rct.Currency, r.GetRoundingRule())
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
	rct.Tax = tt
	if total != nil {
		rct.Total = *total // TODO: Review and test
	}
	return nil
}

// JSONSchemaExtend extends the schema with additional property details
func (rct Receipt) JSONSchemaExtend(js *jsonschema.Schema) {
	props := js.Properties
	// Extend type list
	if its, ok := props.Get("type"); ok {
		its.OneOf = make([]*jsonschema.Schema, len(ReceiptTypes))
		for i, kd := range ReceiptTypes {
			its.OneOf[i] = &jsonschema.Schema{
				Const:       kd.Key.String(),
				Title:       kd.Name.String(),
				Description: kd.Desc.String(),
			}
		}
	}
	rct.Regime.JSONSchemaExtend(js)
	rct.Addons.JSONSchemaExtend(js)
	// Recommendations
	js.Extras = map[string]any{
		schema.Recommended: []string{
			"$regime",
		},
	}
}
