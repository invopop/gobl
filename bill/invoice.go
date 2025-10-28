package bill

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/internal"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

const (
	defaultTaxRemovalAccuracy         uint32 = 2
	defaultCurrencyConversionAccuracy uint32 = 2
)

// Invoice represents a payment claim for goods or services supplied under
// conditions agreed between the supplier and the customer. In most cases
// the resulting document describes the actual financial commitment of goods
// or services ordered from the supplier.
type Invoice struct {
	tax.Regime
	tax.Addons
	tax.Tags
	uuid.Identify

	// Type of invoice document. May be restricted by local tax regime requirements.
	Type cbc.Key `json:"type" jsonschema:"title=Type" jsonschema_extras:"calculated=true"`
	// Series is used to identify groups of invoices by date, business area, project,
	// type of document, customer type, a combination of any or other company specific data.
	// If the output format does not support the series as a separate field, it will be
	// prepended to the code for presentation with a dash (`-`) for separation.
	Series cbc.Code `json:"series,omitempty" jsonschema:"title=Series"`
	// Code is a sequential identifier that uniquely identifies the invoice. The code can
	// be left empty initially, but is **required** to **sign** the invoice.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// Issue date for when the invoice was created and issued. Todays date is used if
	// none is set. There are often legal restrictions on how far back or in the future an
	// invoice can be issued.
	IssueDate cal.Date `json:"issue_date" jsonschema:"title=Issue Date" jsonschema_extras:"calculated=true"`
	// IssueTime is an optional field that may be useful to indicate the time of day when
	// the invoice was issued. Some regions and formats may require this field to be set.
	// An empty string will be automatically updated to reflect the current time, otherwise
	// the field can be left with a nil value.
	IssueTime *cal.Time `json:"issue_time,omitempty" jsonschema:"title=Issue Time" jsonschema_extras:"calculated=true"`
	// Date when the operation defined by the invoice became effective.
	OperationDate *cal.Date `json:"op_date,omitempty" jsonschema:"title=Operation Date"`
	// When the taxes of this invoice become accountable, if none set, the issue date is used.
	ValueDate *cal.Date `json:"value_date,omitempty" jsonschema:"title=Value Date"`
	// Currency for all invoice amounts and totals, unless explicitly stated otherwise.
	Currency currency.Code `json:"currency" jsonschema:"title=Currency" jsonschema_extras:"calculated=true"`
	// Exchange rates to be used when converting the invoices monetary values into other currencies.
	ExchangeRates []*currency.ExchangeRate `json:"exchange_rates,omitempty" jsonschema:"title=Exchange Rates"`

	// Document references for previous invoices that this document replaces or extends.
	Preceding []*org.DocumentRef `json:"preceding,omitempty" jsonschema:"title=Preceding Details"`

	// Special billing tax configuration options.
	Tax *Tax `json:"tax,omitempty" jsonschema:"title=Tax"`

	// The entity supplying the goods or services and usually responsible for paying taxes.
	Supplier *org.Party `json:"supplier" jsonschema:"title=Supplier"`
	// Legal entity receiving the goods or services, may be nil in certain circumstances
	// such as simplified invoices.
	Customer *org.Party `json:"customer,omitempty" jsonschema:"title=Customer"`

	// List of invoice lines representing each of the items sold to the customer.
	Lines []*Line `json:"lines,omitempty" jsonschema:"title=Lines"`
	// Discounts or allowances applied to the complete invoice
	Discounts []*Discount `json:"discounts,omitempty" jsonschema:"title=Discounts"`
	// Charges or surcharges applied to the complete invoice
	Charges []*Charge `json:"charges,omitempty" jsonschema:"title=Charges"`

	// Ordering details including document references and buyer or seller parties.
	Ordering *Ordering `json:"ordering,omitempty" jsonschema:"title=Ordering Details"`
	// Information on when, how, and to whom the invoice should be paid.
	Payment *PaymentDetails `json:"payment,omitempty" jsonschema:"title=Payment Details"`
	// Specific details on delivery of the goods referenced in the invoice.
	Delivery *DeliveryDetails `json:"delivery,omitempty" jsonschema:"title=Delivery Details"`

	// Summary of all the invoice totals, including taxes (calculated).
	Totals *Totals `json:"totals" jsonschema:"title=Totals" jsonschema_extras:"calculated=true"`

	// Unstructured information that is relevant to the invoice, such as correction or additional
	// legal details.
	Notes []*org.Note `json:"notes,omitempty" jsonschema:"title=Notes"`

	// Additional complementary objects that add relevant information to the invoice.
	Complements []*schema.Object `json:"complements,omitempty" jsonschema:"title=Complements"`

	// Additional semi-structured data that doesn't fit into the body of the invoice.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`

	// Attachments provide additional information or supporting documents that are not included
	// in the main document. It is important that attachments are not used for alternative
	// versions of the PDF, for that, see "links" inside the envelope headers.
	Attachments []*org.Attachment `json:"attachments,omitempty" jsonschema:"title=Attachments"`
}

// Validate checks to ensure the invoice is valid and contains all the information we need.
func (inv *Invoice) Validate() error {
	return inv.ValidateWithContext(context.Background())
}

// ValidateWithContext checks to ensure the invoice is valid and contains all the
// information we need.
func (inv *Invoice) ValidateWithContext(ctx context.Context) error {
	ctx = inv.validationContext(ctx)

	var exRule validation.Rule
	exRule = validation.Skip
	if r := inv.RegimeDef(); r != nil {
		// regime specific additions for validation
		exRule = currency.CanConvertInto(inv.ExchangeRates, r.Currency)
	}

	return tax.ValidateStructWithContext(ctx, inv,
		validation.Field(&inv.Regime),
		validation.Field(&inv.Addons),
		validation.Field(&inv.Tags.List, tax.TagsIn(inv.supportedTags()...)),
		validation.Field(&inv.UUID),
		validation.Field(&inv.Type,
			validation.Required,
			isValidInvoiceType,
		),
		validation.Field(&inv.Series),
		validation.Field(&inv.Code,
			validation.When(
				internal.IsSigned(ctx),
				validation.Required.Error("required to sign invoice"),
			),
		),
		validation.Field(&inv.IssueDate,
			cal.DateNotZero(),
		),
		validation.Field(&inv.OperationDate),
		validation.Field(&inv.ValueDate),
		validation.Field(&inv.Currency,
			validation.Required,
			exRule,
		),
		validation.Field(&inv.ExchangeRates,
			validation.Each(validation.NotNil),
		),
		validation.Field(&inv.Preceding,
			validation.Each(validation.NotNil),
		),
		validation.Field(&inv.Tax),
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(validateInvoiceSupplier),
		),
		validation.Field(&inv.Customer,
			validation.By(validateInvoiceCustomer),
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.NotNil,
				validation.By(lineItemHasPrice),
			),
			validation.When(
				len(inv.Discounts) == 0 && len(inv.Charges) == 0,
				validation.Required.Error("cannot be empty without discounts or charges"),
			),
		),
		validation.Field(&inv.Discounts,
			validation.Each(validation.NotNil),
		),
		validation.Field(&inv.Charges,
			validation.Each(validation.NotNil),
		),
		validation.Field(&inv.Ordering),
		validation.Field(&inv.Payment),
		validation.Field(&inv.Delivery),
		validation.Field(&inv.Totals,
			validation.Required,
		),
		validation.Field(&inv.Notes,
			validation.Each(validation.NotNil),
		),
		validation.Field(&inv.Complements,
			validation.Each(validation.NotNil),
		),
		validation.Field(&inv.Meta),
		validation.Field(&inv.Attachments,
			validation.Each(validation.NotNil),
		),
	)
}

func validateInvoiceSupplier(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.Name, validation.Required),
	)
}

func validateInvoiceCustomer(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.Name,
			validation.When(
				partyHasTaxIDCode(p),
				validation.Required,
			),
		),
	)
}

func partyHasTaxIDCode(party *org.Party) bool {
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

// Invert effectively reverses the invoice by inverting the sign of all quantity
// or amount values. Caution should be taken when using this method as
// advances will also be inverted, while payment terms will remain the same,
// which could be confusing if no further modifications are made.
// After inverting the invoice is recalculated and any differences will raise
// an error.
func (inv *Invoice) Invert() error {
	payable := inv.Totals.Payable.Invert()

	for _, row := range inv.Lines {
		row.Quantity = row.Quantity.Invert()
		for _, d := range row.Discounts {
			d.Amount = d.Amount.Invert()
		}
		for _, c := range row.Charges {
			c.Amount = c.Amount.Invert()
		}
	}
	for _, row := range inv.Charges {
		row.Amount = row.Amount.Invert()
	}
	for _, row := range inv.Discounts {
		row.Amount = row.Amount.Invert()
	}
	if inv.Payment != nil {
		for _, row := range inv.Payment.Advances {
			row.Amount = row.Amount.Invert()
		}
	}
	inv.Totals = nil

	if err := inv.Calculate(); err != nil {
		return err
	}

	// The following check tries to ensure that any future fields do not cause
	// unexpected results.
	if !payable.Equals(inv.Totals.Payable) {
		return fmt.Errorf("inverted invoice totals do not match %s != %s", payable.String(), inv.Totals.Payable.String())
	}

	return nil
}

// Empty is a convenience method that will empty all the lines and
// related rows.
func (inv *Invoice) Empty() {
	inv.Lines = make([]*Line, 0)
	inv.Charges = make([]*Charge, 0)
	inv.Discounts = make([]*Discount, 0)
	inv.Totals = nil
	inv.Payment.ResetAdvances()
}

// Calculate performs all the normalizations and calculations required for the invoice
// totals and taxes. If the original invoice only includes partial calculations, this
// will figure out what's missing.
func (inv *Invoice) Calculate() error {
	// Try to set Regime if not already prepared from the supplier's tax ID
	if inv.Regime.IsEmpty() {
		inv.SetRegime(partyTaxCountry(inv.Supplier))
	}

	inv.Normalize(tax.ExtractNormalizers(inv))

	if err := calculate(inv); err != nil {
		return err
	}

	if err := inv.prepareScenarios(); err != nil {
		return err
	}

	return nil
}

// Normalize is run as part of the Calculate method to ensure that the invoice
// is in a consistent state before calculations are performed. This will leverage
// any add-ons alongside the tax regime.
func (inv *Invoice) Normalize(normalizers tax.Normalizers) {
	if inv.Type == cbc.KeyEmpty {
		inv.Type = InvoiceTypeStandard
	}
	inv.Series = cbc.NormalizeCode(inv.Series)
	inv.Code = cbc.NormalizeCode(inv.Code)

	tax.Normalize(normalizers, inv.Tax)
	tax.Normalize(normalizers, inv.Supplier)
	tax.Normalize(normalizers, inv.Customer)
	tax.Normalize(normalizers, inv.Preceding)
	tax.Normalize(normalizers, inv.Lines)
	tax.Normalize(normalizers, inv.Discounts)
	tax.Normalize(normalizers, inv.Charges)
	tax.Normalize(normalizers, inv.Ordering)
	tax.Normalize(normalizers, inv.Payment)

	normalizers.Each(inv)
}

func (inv *Invoice) supportedTags() []cbc.Key {
	ts := defaultInvoiceTags
	if r := inv.RegimeDef(); r != nil {
		ts = ts.Merge(tax.TagSetForSchema(r.Tags, ShortSchemaInvoice))
	}
	for _, a := range inv.AddonDefs() {
		ts = ts.Merge(tax.TagSetForSchema(a.Tags, ShortSchemaInvoice))
	}
	return ts.Keys()
}

// validationContext builds a context with all the validators that the invoice might
// need for execution.
func (inv *Invoice) validationContext(ctx context.Context) context.Context {
	if r := inv.RegimeDef(); r != nil {
		ctx = r.WithContext(ctx)
	}
	for _, a := range inv.AddonDefs() {
		ctx = a.WithContext(ctx)
	}
	return ctx
}

// RemoveIncludedTaxes is a special function that will go through all prices which may include
// the tax included in the invoice, and remove them.
//
// This method will call "Calculate" on the invoice automatically after removing the taxes.
//
// If after removing taxes the totals don't match, a rounding error will be added to the
// invoice totals. In most scenarios this shouldn't be more than a cent or two.
//
// This method will replace the invoice contents in place, or return an error.
func (inv *Invoice) RemoveIncludedTaxes() error {
	return removeIncludedTaxes(inv)
}

/** Calculation Interface Methods **/

func (inv *Invoice) getIssueDate() cal.Date {
	return inv.IssueDate
}
func (inv *Invoice) getIssueTime() *cal.Time {
	return inv.IssueTime
}
func (inv *Invoice) getValueDate() *cal.Date {
	return inv.ValueDate
}
func (inv *Invoice) getTax() *Tax {
	return inv.Tax
}
func (inv *Invoice) getPreceding() []*org.DocumentRef {
	return inv.Preceding
}
func (inv *Invoice) getCustomer() *org.Party {
	return inv.Customer
}
func (inv *Invoice) getCurrency() currency.Code {
	return inv.Currency
}
func (inv *Invoice) getExchangeRates() []*currency.ExchangeRate {
	return inv.ExchangeRates
}
func (inv *Invoice) getLines() []*Line {
	return inv.Lines
}
func (inv *Invoice) getDiscounts() []*Discount {
	return inv.Discounts
}
func (inv *Invoice) getCharges() []*Charge {
	return inv.Charges
}
func (inv *Invoice) getPaymentDetails() *PaymentDetails {
	return inv.Payment
}
func (inv *Invoice) getTotals() *Totals {
	return inv.Totals
}
func (inv *Invoice) getComplements() []*schema.Object {
	return inv.Complements
}

func (inv *Invoice) setIssueDate(d cal.Date) {
	inv.IssueDate = d
}
func (inv *Invoice) setIssueTime(t *cal.Time) {
	inv.IssueTime = t
}
func (inv *Invoice) setCurrency(c currency.Code) {
	inv.Currency = c
}
func (inv *Invoice) setTotals(t *Totals) {
	inv.Totals = t
}

/** ---- **/

// UnmarshalJSON implements the json.Unmarshaler interface and provides any
// data migrations that might be required.
func (inv *Invoice) UnmarshalJSON(data []byte) error {
	type Alias *Invoice
	if err := json.Unmarshal(data, (Alias)(inv)); err != nil {
		return err
	}
	// Ensure there is regime set when coming in from a raw JSON source.
	if inv.Regime.IsEmpty() {
		inv.SetRegime(partyTaxCountry(inv.Supplier))
	}
	// Copy the old tags array from the tax object to the invoice's $tags attribute.
	if inv.Tax != nil && len(inv.Tax.tags) > 0 {
		inv.SetTags(inv.Tax.tags...)
	}
	return nil
}

// JSONSchemaExtend extends the schema with additional property details
func (inv Invoice) JSONSchemaExtend(js *jsonschema.Schema) {
	props := js.Properties
	// Extend type list
	if its, ok := props.Get("type"); ok {
		its.OneOf = make([]*jsonschema.Schema, len(InvoiceTypes))
		for i, kd := range InvoiceTypes {
			its.OneOf[i] = &jsonschema.Schema{
				Const:       kd.Key.String(),
				Title:       kd.Name.String(),
				Description: kd.Desc.String(),
			}
		}
	}
	inv.Regime.JSONSchemaExtend(js)
	inv.Addons.JSONSchemaExtend(js)
	inv.Tags.JSONSchemaExtendWithDefs(js, defaultInvoiceTags.List)
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
