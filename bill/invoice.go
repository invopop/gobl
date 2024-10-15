package bill

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/internal"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"

	"github.com/invopop/validation"
)

// Constants used to help identify invoices
const (
	ShortSchemaInvoice = "bill/invoice"
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

	// Type of invoice document subject to the requirements of the local tax regime.
	Type cbc.Key `json:"type" jsonschema:"title=Type" jsonschema_extras:"calculated=true"`
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

	// Key information regarding previous invoices and potentially details as to why they
	// were corrected.
	Preceding []*org.DocumentRef `json:"preceding,omitempty" jsonschema:"title=Preceding Details"`

	// Special tax configuration for billing.
	Tax *Tax `json:"tax,omitempty" jsonschema:"title=Tax"`

	// The taxable entity supplying the goods or services.
	Supplier *org.Party `json:"supplier" jsonschema:"title=Supplier"`
	// Legal entity receiving the goods or services, may be nil in certain circumstances such as simplified invoices.
	Customer *org.Party `json:"customer,omitempty" jsonschema:"title=Customer"`

	// List of invoice lines representing each of the items sold to the customer.
	Lines []*Line `json:"lines,omitempty" jsonschema:"title=Lines"`
	// Discounts or allowances applied to the complete invoice
	Discounts []*Discount `json:"discounts,omitempty" jsonschema:"title=Discounts"`
	// Charges or surcharges applied to the complete invoice
	Charges []*Charge `json:"charges,omitempty" jsonschema:"title=Charges"`
	// Expenses paid for by the supplier but invoiced directly to the customer.
	Outlays []*Outlay `json:"outlays,omitempty" jsonschema:"title=Outlays"`

	// Ordering details including document references and buyer or seller parties.
	Ordering *Ordering `json:"ordering,omitempty" jsonschema:"title=Ordering Details"`
	// Information on when, how, and to whom the invoice should be paid.
	Payment *Payment `json:"payment,omitempty" jsonschema:"title=Payment Details"`
	// Specific details on delivery of the goods referenced in the invoice.
	Delivery *Delivery `json:"delivery,omitempty" jsonschema:"title=Delivery Details"`

	// Summary of all the invoice totals, including taxes (calculated).
	Totals *Totals `json:"totals" jsonschema:"title=Totals" jsonschema_extras:"calculated=true"`

	// Unstructured information that is relevant to the invoice, such as correction or additional
	// legal details.
	Notes []*cbc.Note `json:"notes,omitempty" jsonschema:"title=Notes"`

	// Additional complementary objects that add relevant information to the invoice.
	Complements []*schema.Object `json:"complements,omitempty" jsonschema:"title=Complements"`

	// Additional semi-structured data that doesn't fit into the body of the invoice.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate checks to ensure the invoice is valid and contains all the information we need.
func (inv *Invoice) Validate() error {
	return inv.ValidateWithContext(context.Background())
}

// ValidateWithContext checks to ensure the invoice is valid and contains all the
// information we need.
func (inv *Invoice) ValidateWithContext(ctx context.Context) error {
	ctx = inv.ValidationContext(ctx)

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
		validation.Field(&inv.ExchangeRates),
		validation.Field(&inv.Preceding),
		validation.Field(&inv.Tax),
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(validateInvoiceSupplier),
		),
		validation.Field(&inv.Customer,
			validation.By(validateInvoiceCustomer),
		),
		validation.Field(&inv.Lines,
			validation.Required,
		),
		validation.Field(&inv.Discounts),
		validation.Field(&inv.Charges),
		validation.Field(&inv.Outlays),
		validation.Field(&inv.Ordering),
		validation.Field(&inv.Payment),
		validation.Field(&inv.Delivery),
		validation.Field(&inv.Totals,
			validation.Required,
		),
		validation.Field(&inv.Notes),
		validation.Field(&inv.Complements),
		validation.Field(&inv.Meta),
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
	for _, row := range inv.Outlays {
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
	inv.Outlays = make([]*Outlay, 0)
	inv.Totals = nil
	inv.Payment.ResetAdvances()
}

// Calculate performs all the normalizations and calculations required for the invoice
// totals and taxes. If the original invoice only includes partial calculations, this
// will figure out what's missing.
func (inv *Invoice) Calculate() error {
	// Try to set Regime if not already prepared from the supplier's tax ID
	if inv.Regime.IsEmpty() {
		inv.SetRegime(inv.supplierTaxCountry())
	}

	inv.Normalize(inv.normalizers())

	if err := inv.calculate(); err != nil {
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

	normalizers.Each(inv)

	tax.Normalize(normalizers, inv.Tax)
	tax.Normalize(normalizers, inv.Supplier)
	tax.Normalize(normalizers, inv.Customer)
	tax.Normalize(normalizers, inv.Preceding)
	tax.Normalize(normalizers, inv.Lines)
	tax.Normalize(normalizers, inv.Discounts)
	tax.Normalize(normalizers, inv.Charges)
	tax.Normalize(normalizers, inv.Ordering)
	tax.Normalize(normalizers, inv.Payment)
}

func (inv *Invoice) normalizers() tax.Normalizers {
	normalizers := make(tax.Normalizers, 0)
	if r := inv.RegimeDef(); r != nil {
		normalizers = normalizers.Append(r.Normalizer)
	}
	for _, a := range inv.GetAddonDefs() {
		normalizers = normalizers.Append(a.Normalizer)
	}
	return normalizers
}

func (inv *Invoice) supportedTags() []cbc.Key {
	var ts *tax.TagSet
	if r := inv.RegimeDef(); r != nil {
		ts = ts.Merge(tax.TagSetForSchema(r.Tags, ShortSchemaInvoice))
	}
	for _, a := range inv.GetAddonDefs() {
		ts = ts.Merge(tax.TagSetForSchema(a.Tags, ShortSchemaInvoice))
	}
	return ts.Keys()
}

// ValidationContext builds a context with all the validators that the invoice might
// need for execution.
func (inv *Invoice) ValidationContext(ctx context.Context) context.Context {
	if r := inv.RegimeDef(); r != nil {
		ctx = r.WithContext(ctx)
	}
	for _, a := range inv.GetAddonDefs() {
		ctx = a.WithContext(ctx)
	}
	return ctx
}

// RemoveIncludedTaxes is a special function that will go through all prices which may include
// the tax included in the invoice, and remove them.
//
// This method will call "Calculate" on the invoice automatically both before and after
// to ensure that the data matches.
//
// If after removing taxes the totals don't match, a rounding error will be added to the
// invoice totals. In most scenarios this shouldn't be more than a cent or two.
//
// A new invoice object is returned, leaving the original instance untouched.
func (inv *Invoice) RemoveIncludedTaxes() (*Invoice, error) {
	if inv.Tax == nil || inv.Tax.PricesInclude.IsEmpty() {
		return inv, nil // nothing to do!
	}

	if err := inv.Calculate(); err != nil {
		return nil, err
	}

	i2 := *inv
	i2.Totals = new(Totals)
	i2.Lines = make([]*Line, len(inv.Lines))
	for i, l := range inv.Lines {
		i2.Lines[i] = l.removeIncludedTaxes(inv.Tax.PricesInclude)
	}

	if len(inv.Discounts) > 0 {
		i2.Discounts = make([]*Discount, len(inv.Discounts))
		for i, l := range inv.Discounts {
			i2.Discounts[i] = l.removeIncludedTaxes(inv.Tax.PricesInclude)
		}
	}
	if len(i2.Charges) > 0 {
		i2.Charges = make([]*Charge, len(inv.Charges))
		for i, l := range inv.Charges {
			i2.Charges[i] = l.removeIncludedTaxes(inv.Tax.PricesInclude)
		}
	}

	tx := *i2.Tax
	tx.PricesInclude = ""
	i2.Tax = &tx

	if err := i2.Calculate(); err != nil {
		return nil, err
	}

	// Account for any rounding errors that we just can't handle
	if !inv.Totals.TotalWithTax.Equals(i2.Totals.TotalWithTax) {
		rnd := inv.Totals.TotalWithTax.Subtract(i2.Totals.TotalWithTax)
		i2.Totals.Rounding = &rnd
		if err := i2.Calculate(); err != nil {
			return nil, err
		}
	}

	return &i2, nil
}

// supplierTaxCountry determines the tax country for the invoice based on the supplier tax
// identity.
func (inv *Invoice) supplierTaxCountry() l10n.TaxCountryCode {
	if inv.Supplier == nil {
		return l10n.CodeEmpty.Tax()
	}
	if inv.Supplier.TaxID == nil {
		return l10n.CodeEmpty.Tax()
	}
	return inv.Supplier.TaxID.Country
}

// calculate does not assume that the tax regime is available.
func (inv *Invoice) calculate() error {
	r := inv.RegimeDef() // may be nil!

	// Normalize data
	if inv.IssueDate.IsZero() {
		inv.IssueDate = cal.TodayIn(r.TimeLocation())
	}
	date := inv.ValueDate
	if date == nil {
		date = &inv.IssueDate
	}

	// Convert empty or invalid currency to the regime's currency
	if inv.Currency == currency.CodeEmpty || inv.Currency.Def() == nil {
		if r == nil {
			return validation.Errors{"currency": errors.New("missing")}
		}
		inv.Currency = r.Currency
	}

	// Prepare the totals we'll need with amounts based on currency
	if inv.Totals == nil {
		inv.Totals = new(Totals)
	}
	t := inv.Totals
	zero := inv.Currency.Def().Zero()
	t.reset(zero)

	// Do we need to deal with the customer-rates tag?
	if inv.HasTags(tax.TagCustomerRates) {
		inv.applyCustomerRates()
	}

	// Lines
	if err := calculateLines(inv.Lines, inv.Currency, inv.ExchangeRates); err != nil {
		return validation.Errors{"lines": err}
	}
	t.Sum = calculateLineSum(inv.Lines, inv.Currency)
	t.Total = t.Sum

	// Discount Lines
	calculateDiscounts(inv.Discounts, t.Sum, zero)
	if discounts := calculateDiscountSum(inv.Discounts, zero); discounts != nil {
		t.Discount = discounts
		t.Total = t.Total.Subtract(*discounts)
	}

	// Charge Lines
	calculateCharges(inv.Charges, t.Sum, zero)
	if charges := calculateChargeSum(inv.Charges, zero); charges != nil {
		t.Charge = charges
		t.Total = t.Total.Add(*charges)
	}

	// Build list of taxable lines
	tls := make([]tax.TaxableLine, 0)
	for _, l := range inv.Lines {
		tls = append(tls, l)
	}
	for _, l := range inv.Discounts {
		tls = append(tls, l)
	}
	for _, l := range inv.Charges {
		tls = append(tls, l)
	}

	// Now figure out the tax totals
	var pit cbc.Code
	if inv.Tax != nil && inv.Tax.PricesInclude != "" {
		pit = inv.Tax.PricesInclude
	}
	t.Taxes = new(tax.Total)
	tc := &tax.TotalCalculator{
		Zero:     zero,
		Country:  inv.Regime.Country,
		Tags:     inv.GetTags(),
		Date:     *date,
		Lines:    tls,
		Includes: pit,
	}
	if err := tc.Calculate(t.Taxes); err != nil {
		return err
	}

	// Remove any included taxes from the total.
	ct := t.Taxes.Category(pit)
	if ct != nil {
		ti := ct.PreciseAmount()
		t.TaxIncluded = &ti
		t.Total = t.Total.Subtract(ti)
	}

	// Finally calculate the total with *all* the taxes.
	t.Tax = t.Taxes.PreciseSum()
	t.TotalWithTax = t.Total.Add(t.Tax)
	t.Payable = t.TotalWithTax
	if t.Rounding != nil {
		// BT-144 in EN16931
		t.Payable = t.Payable.Add(*t.Rounding)
	}

	// Remove taxes object if it doesn't contain any categories
	if len(t.Taxes.Categories) == 0 {
		t.Taxes = nil
	}

	// Outlays
	t.Outlays = calculateOutlays(zero, inv.Outlays)
	if t.Outlays != nil {
		t.Payable = t.Payable.Add(*t.Outlays)
	}

	if inv.Payment != nil {
		inv.Payment.calculateAdvances(zero, t.TotalWithTax)

		// Deal with advances, if any
		if t.Advances = inv.Payment.totalAdvance(zero); t.Advances != nil {
			v := t.Payable.Subtract(*t.Advances)
			t.Due = &v
		}

		// Calculate any due date amounts
		inv.Payment.Terms.CalculateDues(zero, t.Payable)
	}

	t.round(zero)

	// Complements
	if err := calculateComplements(inv.Complements); err != nil {
		return validation.Errors{"complements": err}
	}

	return nil
}

func (inv *Invoice) applyCustomerRates() {
	if inv.Customer == nil || inv.Customer.TaxID == nil {
		return
	}
	country := inv.Customer.TaxID.Country
	for _, l := range inv.Lines {
		addCountryToTaxes(l.Taxes, country)
	}
	for _, d := range inv.Discounts {
		addCountryToTaxes(d.Taxes, country)
	}
	for _, c := range inv.Charges {
		addCountryToTaxes(c.Taxes, country)
	}
}

func addCountryToTaxes(ts tax.Set, country l10n.TaxCountryCode) {
	for _, t := range ts {
		t.Country = country
	}
}

func calculateComplements(comps []*schema.Object) error {
	for _, c := range comps {
		if err := c.Calculate(); err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalJSON implements the json.Unmarshaler interface and provides any
// data migrations that might be required.
func (inv *Invoice) UnmarshalJSON(data []byte) error {
	type Alias *Invoice
	if err := json.Unmarshal(data, (Alias)(inv)); err != nil {
		return err
	}
	// Ensure there is regime set when coming in from a raw JSON source.
	if inv.Regime.IsEmpty() {
		inv.SetRegime(inv.supplierTaxCountry())
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
	// Recommendations
	js.Extras = map[string]any{
		schema.Recommended: []string{
			"$regime",
			"lines",
		},
	}
}
