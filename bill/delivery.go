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

// Delivery document types.
const (
	DeliveryTypeAdvice  cbc.Key = "advice"
	DeliveryTypeNote    cbc.Key = "note"
	DeliveryTypeWaybill cbc.Key = "waybill"
	DeliveryTypeReceipt cbc.Key = "receipt"
)

// DeliveryTypes provides the list of supported delivery documents in GOBL.
var DeliveryTypes = []*cbc.Definition{
	{
		Key: DeliveryTypeAdvice,
		Name: i18n.String{
			i18n.EN: "Delivery Advice",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				A delivery or despatch advice document send by the supplier to indicate how
				an order is to be delivered.
			`),
		},
	},
	{
		Key: DeliveryTypeNote,
		Name: i18n.String{
			i18n.EN: "Delivery Note",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				A delivery note is a document accompanying a shipment of goods that lists the
				items included in the shipment.
			`),
		},
	},
	{
		Key: DeliveryTypeWaybill,
		Name: i18n.String{
			i18n.EN: "Waybill",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				A waybill is a document issued by a carrier giving details and instructions
				relating to the shipment of a consignment of goods.
			`),
		},
	},
	{
		Key: DeliveryTypeReceipt,
		Name: i18n.String{
			i18n.EN: "Delivery Receipt",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				A delivery receipt is a document that is signed by the person who receives
				goods to confirm that they have been received.
			`),
		},
	},
}

var isValidDeliveryType = cbc.InKeyDefs(DeliveryTypes)

// Delivery document used to describe the delivery of goods or potentially also services.
type Delivery struct {
	tax.Regime
	tax.Addons
	tax.Tags
	uuid.Identify

	// Type of delivery document.
	Type cbc.Key `json:"type" jsonschema:"title=Type" jsonschema_extras:"enum=advice,note,waybill,receipt"`
	// Series is used to identify groups of deliveries by date, business area, project,
	// type, customer, a combination of any, or other company specific data.
	// If the output format does not support the series as a separate field, it will be
	// prepended to the code for presentation with a dash (`-`) for separation.
	Series cbc.Code `json:"series,omitempty" jsonschema:"title=Series"`
	// Code is a sequential identifier that uniquely identifies the delivery. The code can
	// be left empty initially, but is **required** to **sign** the document.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// When the delivery document is to be issued.
	IssueDate cal.Date `json:"issue_date" jsonschema:"title=Issue Date" jsonschema_extras:"calculated=true"`
	// IssueTime is an optional field that may be useful to indicate the time of day when
	// the delivery was issued. Some regions and formats may require this field to be set.
	// An empty string will be automatically updated to reflect the current time, otherwise
	// the field can be left with a nil value.
	IssueTime *cal.Time `json:"issue_time,omitempty" jsonschema:"title=Issue Time" jsonschema_extras:"calculated=true"`
	// When the taxes of this delivery become accountable, if none set, the issue date is used.
	ValueDate *cal.Date `json:"value_date,omitempty" jsonschema:"title=Value Date"`
	// Currency for all delivery totals.
	Currency currency.Code `json:"currency,omitempty" jsonschema:"title=Currency" jsonschema_extras:"calculated=true"`
	// Exchange rates to be used when converting the invoices monetary values into other currencies.
	ExchangeRates []*currency.ExchangeRate `json:"exchange_rates,omitempty" jsonschema:"title=Exchange Rates"`

	// Ordering details for the delivery, including links to other documents.
	Ordering *Ordering `json:"ordering,omitempty" jsonschema:"title=Ordering"`
	// Key information regarding previous delivery documents that this one will either
	// extend or replace.
	Preceding []*org.DocumentRef `json:"preceding,omitempty" jsonschema:"title=Preceding Details"`

	// Tracking is used to define specific codes or IDs that may be used to
	// identify and track delivery.
	Tracking *Tracking `json:"tracking,omitempty" jsonschema:"title=Tracking"`
	// DespatchDate is the date when the goods are expected to be despatched.
	DespatchDate *cal.Date `json:"despatch_date,omitempty" jsonschema:"title=Despatch Date"`
	// ReceiveDate is the date when the goods are expected to be received.
	ReceiveDate *cal.Date `json:"receive_date,omitempty" jsonschema:"title=Receive Date"`

	// Special tax configuration for calculating totals.
	Tax *Tax `json:"tax,omitempty" jsonschema:"title=Tax"`

	// The entity supplying the goods or services and usually responsible for paying taxes.
	Supplier *org.Party `json:"supplier" jsonschema:"title=Supplier"`
	// Legal entity receiving the goods or services, may be nil in certain circumstances such as simplified invoices.
	Customer *org.Party `json:"customer,omitempty" jsonschema:"title=Customer"`
	// The party who will despatch the goods defined in the invoice.
	Despatcher *org.Party `json:"despatcher,omitempty" jsonschema:"title=Despatcher"`
	// The party who will receive delivery of the goods defined in the invoice.
	Receiver *org.Party `json:"receiver,omitempty" jsonschema:"title=Receiver"`
	// The courier responsible for delivering the goods.
	Courier *org.Party `json:"courier,omitempty" jsonschema:"title=Courier"`

	// List of lines representing each of the items to be ordered.
	Lines []*Line `json:"lines,omitempty" jsonschema:"title=Lines"`
	// Discounts or allowances applied to order totals
	Discounts []*Discount `json:"discounts,omitempty" jsonschema:"title=Discounts"`
	// Charges or surcharges applied to order totals
	Charges []*Charge `json:"charges,omitempty" jsonschema:"title=Charges"`

	// Summary of all the order totals, including taxes (calculated).
	Totals *Totals `json:"totals,omitempty" jsonschema:"title=Totals" jsonschema_extras:"calculated=true"`

	// Unstructured information that is relevant to the delivery, such as correction or additional
	// legal details.
	Notes []*org.Note `json:"notes,omitempty" jsonschema:"title=Notes"`

	// Additional complementary objects that add relevant information to the delivery.
	Complements []*schema.Object `json:"complements,omitempty" jsonschema:"title=Complements"`

	// Additional semi-structured data that doesn't fit into the body of the delivery.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`

	// Attachments provide additional information or supporting documents that are not included
	// in the main document. It is important that attachments are not used for alternative
	// versions of the PDF, for that, see "links" inside the envelope headers.
	Attachments []*org.Attachment `json:"attachments,omitempty" jsonschema:"title=Attachments"`
}

// Tracking stores tracking information about a delivery.
type Tracking struct {
	// Code used for tracking
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// Website to access for more tracking details
	Website *org.Website `json:"website,omitempty" jsonschema:"title=Website"`
}

// Validate the delivery document
func (dlv *Delivery) Validate() error {
	return dlv.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures that the fields contained in the Delivery look correct.
func (dlv *Delivery) ValidateWithContext(ctx context.Context) error {
	ctx = dlv.validationContext(ctx)
	r := dlv.RegimeDef()
	return tax.ValidateStructWithContext(ctx, dlv,
		validation.Field(&dlv.Regime),
		validation.Field(&dlv.Addons),
		validation.Field(&dlv.Tags),
		validation.Field(&dlv.UUID),
		validation.Field(&dlv.Type,
			validation.Required,
			isValidDeliveryType,
		),
		validation.Field(&dlv.Series),
		validation.Field(&dlv.Code,
			validation.When(
				internal.IsSigned(ctx),
				validation.Required.Error("required to sign delivery"),
			),
		),
		validation.Field(&dlv.IssueDate,
			validation.Required,
			cal.DateNotZero(),
		),
		validation.Field(&dlv.ValueDate),
		validation.Field(&dlv.Currency,
			currency.CanConvertInto(dlv.ExchangeRates, r.GetCurrency()),
		),
		validation.Field(&dlv.ExchangeRates,
			validation.Each(validation.NotNil),
		),

		validation.Field(&dlv.Ordering),
		validation.Field(&dlv.Preceding,
			validation.Each(validation.NotNil),
		),

		validation.Field(&dlv.Tracking),
		validation.Field(&dlv.DespatchDate),
		validation.Field(&dlv.ReceiveDate),

		validation.Field(&dlv.Tax),

		validation.Field(&dlv.Supplier, validation.Required),
		validation.Field(&dlv.Customer),
		validation.Field(&dlv.Despatcher),
		validation.Field(&dlv.Receiver),
		validation.Field(&dlv.Courier),

		validation.Field(&dlv.Lines,
			validation.Required,
			validation.Each(validation.NotNil),
		),
		validation.Field(&dlv.Discounts,
			validation.Each(validation.NotNil),
		),
		validation.Field(&dlv.Charges,
			validation.Each(validation.NotNil),
		),

		validation.Field(&dlv.Totals),
		validation.Field(&dlv.Notes,
			validation.Each(validation.NotNil),
		),
		validation.Field(&dlv.Complements,
			validation.Each(validation.NotNil),
		),
		validation.Field(&dlv.Meta),
		validation.Field(&dlv.Attachments,
			validation.Each(validation.NotNil),
		),
	)
}

// validationContext builds a context with all the validators that the delivery might
// need for execution.
func (dlv *Delivery) validationContext(ctx context.Context) context.Context {
	if r := dlv.RegimeDef(); r != nil {
		ctx = r.WithContext(ctx)
	}
	for _, a := range dlv.AddonDefs() {
		ctx = a.WithContext(ctx)
	}
	return ctx
}

// Calculate performs all the normalizations and calculations required for the delivery
// totals and taxes. If the original delivery only includes partial calculations, this
// will figure out what's missing.
func (dlv *Delivery) Calculate() error {
	// Try to set Regime if not already prepared from the supplier's tax ID
	if dlv.Regime.IsEmpty() {
		dlv.SetRegime(partyTaxCountry(dlv.Supplier))
	}
	dlv.Normalize(dlv.normalizers())
	return calculate(dlv)
}

// Normalize is run as part of the Calculate method to ensure that the delivery
// is in a consistent state before calculations are performed. This will leverage
// any add-ons alongside the tax regime.
func (dlv *Delivery) Normalize(normalizers tax.Normalizers) {
	if dlv.Type == cbc.KeyEmpty {
		dlv.Type = DeliveryTypeAdvice
	}
	dlv.Series = cbc.NormalizeCode(dlv.Series)
	dlv.Code = cbc.NormalizeCode(dlv.Code)

	tax.Normalize(normalizers, dlv.Tax)
	tax.Normalize(normalizers, dlv.Supplier)
	tax.Normalize(normalizers, dlv.Customer)
	tax.Normalize(normalizers, dlv.Despatcher)
	tax.Normalize(normalizers, dlv.Receiver)
	tax.Normalize(normalizers, dlv.Preceding)
	tax.Normalize(normalizers, dlv.Lines)
	tax.Normalize(normalizers, dlv.Discounts)
	tax.Normalize(normalizers, dlv.Charges)

	normalizers.Each(dlv)
}

// normalizers returns the normalizers for the delivery.
func (dlv *Delivery) normalizers() tax.Normalizers {
	normalizers := make(tax.Normalizers, 0)
	if r := dlv.RegimeDef(); r != nil {
		normalizers = normalizers.Append(r.Normalizer)
	}
	for _, a := range dlv.AddonDefs() {
		normalizers = normalizers.Append(a.Normalizer)
	}
	return normalizers
}

// ConvertInto will use the defined exchange rates in the delivery to convert all the prices
// into the given currency.
//
// The intent of this method is help convert the delivery amounts when the destination is
// unable or unwilling to handle the current currency. This is typically the case
// with tax related reports or declarations.
//
// The method will return a new delivery with all the amounts converted into the given
// currency or an error if the conversion is not possible.
//
// Conversion is done by first exchanging the lowest common amounts to the destination
// currency, then recalculating the totals.
func (dlv *Delivery) ConvertInto(cur currency.Code) (*Delivery, error) {
	// Calculate ensures that all the totals and amounts have been prepared
	// so we can make assumptions about the data that will be available,
	// including the original currency!
	if err := dlv.Calculate(); err != nil {
		return nil, err
	}

	if dlv.Currency == cur {
		return dlv, nil
	}
	ex := currency.MatchExchangeRate(dlv.ExchangeRates, dlv.Currency, cur)
	if ex == nil {
		return nil, fmt.Errorf("no exchange rate defined for '%v' to '%v'", dlv.Currency, cur)
	}

	d2 := *dlv
	d2.Totals = nil
	d2.Lines = convertLinesInto(ex, dlv.Lines)
	d2.Discounts = convertDiscountsInto(ex, dlv.Discounts)
	d2.Charges = convertChargesInto(ex, dlv.Charges)
	d2.Currency = cur

	if err := d2.Calculate(); err != nil {
		return nil, err
	}

	return &d2, nil
}

/** Calculation Interface Methods **/

func (dlv *Delivery) getIssueDate() cal.Date {
	return dlv.IssueDate
}
func (dlv *Delivery) getIssueTime() *cal.Time {
	return dlv.IssueTime
}
func (dlv *Delivery) getValueDate() *cal.Date {
	return dlv.ValueDate
}
func (dlv *Delivery) getTax() *Tax {
	return dlv.Tax
}
func (dlv *Delivery) getPreceding() []*org.DocumentRef {
	return dlv.Preceding
}
func (dlv *Delivery) getCustomer() *org.Party {
	return dlv.Customer
}
func (dlv *Delivery) getCurrency() currency.Code {
	return dlv.Currency
}
func (dlv *Delivery) getExchangeRates() []*currency.ExchangeRate {
	return dlv.ExchangeRates
}
func (dlv *Delivery) getLines() []*Line {
	return dlv.Lines
}
func (dlv *Delivery) getDiscounts() []*Discount {
	return dlv.Discounts
}
func (dlv *Delivery) getCharges() []*Charge {
	return dlv.Charges
}
func (dlv *Delivery) getPaymentDetails() *PaymentDetails {
	return nil // no payment for deliveries
}
func (dlv *Delivery) getTotals() *Totals {
	return dlv.Totals
}
func (dlv *Delivery) getComplements() []*schema.Object {
	return dlv.Complements
}

func (dlv *Delivery) setIssueDate(d cal.Date) {
	dlv.IssueDate = d
}
func (dlv *Delivery) setIssueTime(t *cal.Time) {
	dlv.IssueTime = t
}
func (dlv *Delivery) setCurrency(c currency.Code) {
	dlv.Currency = c
}
func (dlv *Delivery) setTotals(t *Totals) {
	dlv.Totals = t
}

/** ---- **/

// JSONSchemaExtend extends the schema with additional property details
func (dlv Delivery) JSONSchemaExtend(js *jsonschema.Schema) {
	props := js.Properties
	// Extend type list
	if its, ok := props.Get("type"); ok {
		its.OneOf = make([]*jsonschema.Schema, len(DeliveryTypes))
		for i, kd := range DeliveryTypes {
			its.OneOf[i] = &jsonschema.Schema{
				Const:       kd.Key.String(),
				Title:       kd.Name.String(),
				Description: kd.Desc.String(),
			}
		}
	}
	dlv.Regime.JSONSchemaExtend(js)
	dlv.Addons.JSONSchemaExtend(js)
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
