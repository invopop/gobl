package bill

import (
	"context"
	"fmt"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/internal"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Predefined list of the order types supported.
const (
	OrderTypePurchase cbc.Key = "purchase"
	OrderTypeSale     cbc.Key = "sale"
	OrderTypeQuote    cbc.Key = "quote"
)

// OrderTypes defines the list of order types supported.
var OrderTypes = []*cbc.Definition{
	{
		Key: OrderTypePurchase,
		Name: i18n.String{
			i18n.EN: "Purchase Order",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				A purchase order is a document that a buyer sends to a seller to request goods or services.
			`),
		},
	},
	{
		Key: OrderTypeSale,
		Name: i18n.String{
			i18n.EN: "Sales Order",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				A sales order is a document that a seller sends to a buyer to confirm the sale of goods or services.
			`),
		},
	},
	{
		Key: OrderTypeQuote,
		Name: i18n.String{
			i18n.EN: "Quote",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				A quote is a document that a seller sends to a buyer to provide a price for goods or services.
			`),
		},
	},
}

var isValidOrderType = cbc.InKeyDefs(OrderTypes)

// Order documents are used for the initial part of a order-to-invoice process
// where the buyer requests goods or services from the seller.
type Order struct {
	tax.Regime
	tax.Addons
	tax.Tags
	uuid.Identify

	// Type of the order.
	Type cbc.Key `json:"type,omitempty" jsonschema:"title=Type"`
	// Used as a prefix to group codes.
	Series cbc.Code `json:"series,omitempty" jsonschema:"title=Series"`
	// Sequential code used to identify this invoice in tax declarations.
	Code cbc.Code `json:"code" jsonschema:"title=Code"`
	// When the invoice was created.
	IssueDate cal.Date `json:"issue_date" jsonschema:"title=Issue Date" jsonschema_extras:"calculated=true"`
	// Date when the operation defined by the invoice became effective.
	OperationDate *cal.Date `json:"op_date,omitempty" jsonschema:"title=Operation Date"`
	// When the taxes of this invoice become accountable, if none set, the issue date is used.
	ValueDate *cal.Date `json:"value_date,omitempty" jsonschema:"title=Value Date"`
	// Currency for all invoice totals.
	Currency currency.Code `json:"currency" jsonschema:"title=Currency" jsonschema_extras:"calculated=true"`
	// Exchange rates to be used when converting the invoices monetary values into other currencies.
	ExchangeRates []*currency.ExchangeRate `json:"exchange_rates,omitempty" jsonschema:"title=Exchange Rates"`

	// The identification of contracts.
	Contracts []*org.DocumentRef `json:"contracts,omitempty" jsonschema:"title=Contracts"`
	// Key information regarding previous order documents.
	Preceding []*org.DocumentRef `json:"preceding,omitempty" jsonschema:"title=Preceding Details"`

	// Additional codes, IDs, SKUs, or other regional or custom identifiers that may be used to identify the order.
	Identities []*org.Identity `json:"identities,omitempty" jsonschema:"title=Identities"`
	// Period of time in which the order is valid.
	Period *cal.Period `json:"period,omitempty" jsonschema:"title=Period"`

	// Special tax configuration for billing.
	Tax *Tax `json:"tax,omitempty" jsonschema:"title=Tax"`

	// The entity supplying the goods or services and usually responsible for paying taxes.
	Supplier *org.Party `json:"supplier" jsonschema:"title=Supplier"`
	// Legal entity receiving the goods or services, may be nil in certain circumstances such as simplified invoices.
	Customer *org.Party `json:"customer,omitempty" jsonschema:"title=Customer"`
	// Party who is responsible for issuing payment, if not the same as the customer.
	Buyer *org.Party `json:"buyer,omitempty" jsonschema:"title=Buyer"`
	// Seller is the party liable to pay taxes on the transaction if not the same as the supplier.
	Seller *org.Party `json:"seller,omitempty" jsonschema:"title=Seller"`

	// List of lines representing each of the items to be ordered.
	Lines []*Line `json:"lines,omitempty" jsonschema:"title=Lines"`
	// Discounts or allowances applied to order totals
	Discounts []*Discount `json:"discounts,omitempty" jsonschema:"title=Discounts"`
	// Charges or surcharges applied to order totals
	Charges []*Charge `json:"charges,omitempty" jsonschema:"title=Charges"`

	// Information on when, how, and to whom a final invoice would be paid.
	Payment *PaymentDetails `json:"payment,omitempty" jsonschema:"title=Payment Details"`
	// Specific details on delivery of the goods to copy to the final invoice.
	Delivery *DeliveryDetails `json:"delivery,omitempty" jsonschema:"title=Delivery Details"`

	// Summary of all the order totals, including taxes (calculated).
	Totals *Totals `json:"totals,omitempty" jsonschema:"title=Totals" jsonschema_extras:"calculated=true"`

	// Unstructured information that is relevant to the order, such as correction or additional
	// legal details.
	Notes []*org.Note `json:"notes,omitempty" jsonschema:"title=Notes"`

	// Additional complementary objects that add relevant information to the order.
	Complements []*schema.Object `json:"complements,omitempty" jsonschema:"title=Complements"`

	// Additional semi-structured data that doesn't fit into the body of the order.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`

	// Attachments provide additional information or supporting documents that are not included
	// in the main document. It is important that attachments are not used for alternative
	// versions of the PDF, for that, see "links" inside the envelope headers.
	Attachments []*org.Attachment `json:"attachments,omitempty" jsonschema:"title=Attachments"`
}

// Validate runs the validation rules for the order without the context.
func (ord *Order) Validate() error {
	return ord.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures that the fields contained in the Order look correct.
func (ord *Order) ValidateWithContext(ctx context.Context) error {
	ctx = ord.validationContext(ctx)
	r := ord.RegimeDef()
	return tax.ValidateStructWithContext(ctx, ord,
		validation.Field(&ord.Regime),
		validation.Field(&ord.Addons),
		validation.Field(&ord.UUID),
		validation.Field(&ord.Type,
			validation.Required,
			isValidOrderType,
		),
		validation.Field(&ord.Series),
		validation.Field(&ord.Code,
			validation.When(
				internal.IsSigned(ctx),
				validation.Required.Error("required to sign order"),
			),
		),
		validation.Field(&ord.IssueDate,
			validation.Required,
			cal.DateNotZero(),
		),
		validation.Field(&ord.Currency,
			validation.Required,
			currency.CanConvertInto(ord.ExchangeRates, r.GetCurrency()),
		),
		validation.Field(&ord.ExchangeRates),
		validation.Field(&ord.Contracts),
		validation.Field(&ord.Preceding),
		validation.Field(&ord.Supplier, validation.Required),
		validation.Field(&ord.Customer),
		validation.Field(&ord.Buyer),
		validation.Field(&ord.Seller),
		validation.Field(&ord.Lines, validation.Required),
		validation.Field(&ord.Discounts),
		validation.Field(&ord.Charges),
		validation.Field(&ord.Payment),
		validation.Field(&ord.Delivery),
		validation.Field(&ord.Totals),
		validation.Field(&ord.Notes),
		validation.Field(&ord.Complements),
		validation.Field(&ord.Meta),
		validation.Field(&ord.Attachments),
	)
}

// validationContext builds a context with all the validators that the order might
// need for execution.
func (ord *Order) validationContext(ctx context.Context) context.Context {
	if r := ord.RegimeDef(); r != nil {
		ctx = r.WithContext(ctx)
	}
	for _, a := range ord.AddonDefs() {
		ctx = a.WithContext(ctx)
	}
	return ctx
}

// Calculate performs all the normalizations and calculations required for the order
// totals and taxes. If the original order only includes partial calculations, this
// will figure out what's missing.
func (ord *Order) Calculate() error {
	// Try to set Regime if not already prepared from the supplier's tax ID
	if ord.Regime.IsEmpty() {
		ord.SetRegime(partyTaxCountry(ord.Supplier))
	}
	ord.Normalize(ord.normalizers())
	return calculate(ord)
}

// Normalize is run as part of the Calculate method to ensure that the order
// is in a consistent state before calculations are performed. This will leverage
// any add-ons alongside the tax regime.
func (ord *Order) Normalize(normalizers tax.Normalizers) {
	if ord.Type == cbc.KeyEmpty {
		ord.Type = OrderTypePurchase
	}
	ord.Series = cbc.NormalizeCode(ord.Series)
	ord.Code = cbc.NormalizeCode(ord.Code)

	normalizers.Each(ord)

	tax.Normalize(normalizers, ord.Tax)
	tax.Normalize(normalizers, ord.Supplier)
	tax.Normalize(normalizers, ord.Customer)
	tax.Normalize(normalizers, ord.Buyer)
	tax.Normalize(normalizers, ord.Seller)
	tax.Normalize(normalizers, ord.Preceding)
	tax.Normalize(normalizers, ord.Lines)
	tax.Normalize(normalizers, ord.Discounts)
	tax.Normalize(normalizers, ord.Charges)
	tax.Normalize(normalizers, ord.Payment)
	tax.Normalize(normalizers, ord.Delivery)
}

func (ord *Order) normalizers() tax.Normalizers {
	normalizers := make(tax.Normalizers, 0)
	if r := ord.RegimeDef(); r != nil {
		normalizers = normalizers.Append(r.Normalizer)
	}
	for _, a := range ord.AddonDefs() {
		normalizers = normalizers.Append(a.Normalizer)
	}
	return normalizers
}

// ConvertInto will use the defined exchange rates in the order to convert all the prices
// into the given currency.
//
// The intent of this method is help convert the order amounts when the destination is
// unable or unwilling to handle the current currency. This is typically the case
// with tax related reports or declarations.
//
// The method will return a new order with all the amounts converted into the given
// currency or an error if the conversion is not possible.
//
// Conversion is done by first exchanging the lowest common amounts to the destination
// currency, then recalculating the totals.
func (ord *Order) ConvertInto(cur currency.Code) (*Order, error) {
	// Calculate ensures that all the totals and amounts have been prepared
	// so we can make assumptions about the data that will be available,
	// including the original currency!
	if err := ord.Calculate(); err != nil {
		return nil, err
	}

	if ord.Currency == cur {
		return ord, nil
	}
	ex := currency.MatchExchangeRate(ord.ExchangeRates, ord.Currency, cur)
	if ex == nil {
		return nil, fmt.Errorf("no exchange rate defined for '%v' to '%v'", ord.Currency, cur)
	}

	o2 := *ord
	o2.Totals = nil
	o2.Lines = convertLinesInto(ex, ord.Lines)
	o2.Discounts = convertDiscountsInto(ex, ord.Discounts)
	o2.Charges = convertChargesInto(ex, ord.Charges)
	o2.Payment = convertPaymentDetailsInto(ex, ord.Payment)
	o2.Currency = cur

	if err := o2.Calculate(); err != nil {
		return nil, err
	}

	return &o2, nil
}

/** Calculation Interface Methods **/

func (ord *Order) getIssueDate() cal.Date {
	return ord.IssueDate
}
func (ord *Order) getValueDate() *cal.Date {
	return ord.ValueDate
}
func (ord *Order) getTax() *Tax {
	return ord.Tax
}
func (ord *Order) getCustomer() *org.Party {
	return ord.Customer
}
func (ord *Order) getCurrency() currency.Code {
	return ord.Currency
}
func (ord *Order) getExchangeRates() []*currency.ExchangeRate {
	return ord.ExchangeRates
}
func (ord *Order) getLines() []*Line {
	return ord.Lines
}
func (ord *Order) getDiscounts() []*Discount {
	return ord.Discounts
}
func (ord *Order) getCharges() []*Charge {
	return ord.Charges
}
func (ord *Order) getPaymentDetails() *PaymentDetails {
	return ord.Payment
}
func (ord *Order) getTotals() *Totals {
	return ord.Totals
}
func (ord *Order) getComplements() []*schema.Object {
	return ord.Complements
}

func (ord *Order) setIssueDate(d cal.Date) {
	ord.IssueDate = d
}
func (ord *Order) setCurrency(c currency.Code) {
	ord.Currency = c
}
func (ord *Order) setTotals(t *Totals) {
	ord.Totals = t
}

/** ---- **/

// JSONSchemaExtend extends the schema with additional property details
func (ord Order) JSONSchemaExtend(js *jsonschema.Schema) {
	props := js.Properties
	// Extend type list
	if its, ok := props.Get("type"); ok {
		its.OneOf = make([]*jsonschema.Schema, len(OrderTypes))
		for i, kd := range OrderTypes {
			its.OneOf[i] = &jsonschema.Schema{
				Const:       kd.Key.String(),
				Title:       kd.Name.String(),
				Description: kd.Desc.String(),
			}
		}
	}
	ord.Regime.JSONSchemaExtend(js)
	ord.Addons.JSONSchemaExtend(js)
	// Recommendations
	js.Extras = map[string]any{
		schema.Recommended: []string{
			"$regime",
		},
	}
}
