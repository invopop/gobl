package bill

import (
	"errors"
	"fmt"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"

	"github.com/invopop/validation"
)

// Constants used to help identify invoices
const (
	ShortSchemaInvoice = "bill/invoice"
)

// Invoice represents a payment claim for goods or services supplied under
// conditions agreed between the supplier and the customer. In most cases
// the resulting document describes the actual financial commitment of goods
// or services ordered from the supplier.
type Invoice struct {
	// Unique document ID. Not required, but always recommended in addition to the Code.
	UUID *uuid.UUID `json:"uuid,omitempty" jsonschema:"title=UUID"`
	// Used as a prefix to group codes.
	Series string `json:"series,omitempty" jsonschema:"title=Series"`
	// Sequential code used to identify this invoice in tax declarations.
	Code string `json:"code" jsonschema:"title=Code"`
	// Type of invoice document subject to the requirements of the local tax regime.
	Type cbc.Key `json:"type" jsonschema:"title=Type"`
	// Currency for all invoice totals.
	Currency currency.Code `json:"currency" jsonschema:"title=Currency"`
	// Exchange rates to be used when converting the invoices monetary values into other currencies.
	ExchangeRates []*currency.ExchangeRate `json:"exchange_rates,omitempty" jsonschema:"title=Exchange Rates"`
	// Special tax configuration for billing.
	Tax *Tax `json:"tax,omitempty" jsonschema:"title=Tax"`

	// Key information regarding previous invoices and potentially details as to why they
	// were corrected.
	Preceding []*Preceding `json:"preceding,omitempty" jsonschema:"title=Preceding Details"`

	// When the invoice was created.
	IssueDate cal.Date `json:"issue_date" jsonschema:"title=Issue Date"`
	// Date when the operation defined by the invoice became effective.
	OperationDate *cal.Date `json:"op_date,omitempty" jsonschema:"title=Operation Date"`
	// When the taxes of this invoice become accountable, if none set, the issue date is used.
	ValueDate *cal.Date `json:"value_date,omitempty" jsonschema:"title=Value Date"`

	// The taxable entity supplying the goods or services.
	Supplier *org.Party `json:"supplier" jsonschema:"title=Supplier"`
	// Legal entity receiving the goods or services, may be empty in certain circumstances such as simplified invoices.
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

	// The EN 16931-1:2017 standard recognises a need to be able to attach additional
	// documents to an invoice. We don't support this yet, but this is where
	// it could go.
	//Attachments Attachments `json:"attachments,omitempty" jsonschema:"title=Attachments"`

	// Unstructured information that is relevant to the invoice, such as correction or additional
	// legal details.
	Notes []*cbc.Note `json:"notes,omitempty" jsonschema:"title=Notes"`
	// Additional semi-structured data that doesn't fit into the body of the invoice.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate checks to ensure the invoice is valid and contains all the information we need.
func (inv *Invoice) Validate() error {
	err := validation.ValidateStruct(inv,
		validation.Field(&inv.UUID),
		validation.Field(&inv.Code, validation.Required),
		validation.Field(&inv.Type, validation.Required, isValidInvoiceType),
		validation.Field(&inv.Currency, validation.Required),
		validation.Field(&inv.ExchangeRates),
		validation.Field(&inv.Tax),

		validation.Field(&inv.Preceding),

		validation.Field(&inv.IssueDate, cal.DateNotZero()),
		validation.Field(&inv.OperationDate),
		validation.Field(&inv.ValueDate),

		validation.Field(&inv.Supplier, validation.Required),
		validation.Field(&inv.Customer),

		validation.Field(&inv.Lines, validation.Required),
		validation.Field(&inv.Discounts),
		validation.Field(&inv.Charges),
		validation.Field(&inv.Outlays),

		validation.Field(&inv.Ordering),
		validation.Field(&inv.Payment),
		validation.Field(&inv.Delivery),

		validation.Field(&inv.Totals, validation.Required),

		validation.Field(&inv.Notes),
		validation.Field(&inv.Meta),
	)
	if err == nil && inv.Supplier != nil {
		// Always validate contents using supplier's tax
		// identity.
		tID := inv.Supplier.TaxID
		if tID == nil {
			return errors.New("supplier: missing tax identity")
		}
		r := tax.Regimes().For(tID.Country, tID.Zone)
		err = r.ValidateObject(inv)
	}
	return err
}

// Totals contains the summaries of all calculations for the invoice.
type Totals struct {
	// Sum of all line item sums
	Sum num.Amount `json:"sum" jsonschema:"title=Sum"`
	// Sum of all document level discounts
	Discount *num.Amount `json:"discount,omitempty" jsonschema:"title=Discount"`
	// Sum of all document level charges
	Charge *num.Amount `json:"charge,omitempty" jsonschema:"title=Charge"`
	// If prices include tax, this is the total tax included in the price.
	TaxIncluded *num.Amount `json:"tax_included,omitempty" jsonschema:"title=Tax Included"`
	// Sum of all line sums minus the discounts, plus the charges, without tax.
	Total num.Amount `json:"total" jsonschema:"title=Total"`
	// Summary of all the taxes included in the invoice.
	Taxes *tax.Total `json:"taxes,omitempty" jsonschema:"title=Tax Totals"`
	// Total amount of tax to apply to the invoice.
	Tax num.Amount `json:"tax,omitempty" jsonschema:"title=Tax"`
	// Grand total after all taxes have been applied.
	TotalWithTax num.Amount `json:"total_with_tax" jsonschema:"title=Total with Tax"`
	// Total paid in outlays that need to be reimbursed
	Outlays *num.Amount `json:"outlays,omitempty" jsonschema:"title=Outlay Totals"`
	// Total amount to be paid after applying taxes and outlays.
	Payable num.Amount `json:"payable" jsonschema:"title=Payable"`
	// Total amount already paid in advance.
	Advances *num.Amount `json:"advance,omitempty" jsonschema:"title=Advance"`
	// How much actually needs to be paid now.
	Due *num.Amount `json:"due,omitempty" jsonschema:"title=Due"`
}

// Validate the totals used in invoice.
func (t *Totals) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.Sum, validation.Required),
		validation.Field(&t.Discount),
		validation.Field(&t.Charge),
		validation.Field(&t.TaxIncluded),
		validation.Field(&t.Total, validation.Required),
		validation.Field(&t.Taxes),
		validation.Field(&t.Tax),
		validation.Field(&t.TotalWithTax),
		validation.Field(&t.Outlays),
		validation.Field(&t.Payable),
		validation.Field(&t.Advances),
		validation.Field(&t.Due),
	)
}

// Calculate performs all the calculations required for the invoice totals and taxes. If the original
// invoice only includes partial calculations, this will figure out what's missing.
func (inv *Invoice) Calculate() error {
	if inv.Type == cbc.KeyEmpty {
		inv.Type = InvoiceTypeStandard
	}
	if inv.Supplier == nil {
		return errors.New("missing or invalid supplier tax identity")
	}
	if err := inv.Supplier.Calculate(); err != nil {
		return fmt.Errorf("supplier: %w", err)
	}
	if inv.Customer != nil {
		if err := inv.Customer.Calculate(); err != nil {
			return fmt.Errorf("customer: %w", err)
		}
	}

	if err := inv.prepareTagsAndScenarios(); err != nil {
		return err
	}

	// Should we use the customers identity for calculations?
	tID := inv.determineTaxIdentity()
	if tID == nil {
		return errors.New("unable to determine tax identity")
	}
	r := tax.RegimeFor(tID.Country, tID.Zone)
	if r == nil {
		return fmt.Errorf("no tax regime for %v", tID.Country)
	}

	return inv.calculate(r, tID)
}

// RemoveIncludedTaxes is a special function that will go through all prices which may include
// the tax included in the invoice, and remove them. The accuracy parameter is used to determine
// the additional exponent that will be added to prices before calculation with the aim of reducing
// rounding errors. An accuracy value of 2 is recommended.
//
// A new invoice object is returned, leaving the original objects untouched.
func (inv *Invoice) RemoveIncludedTaxes(accuracy uint32) *Invoice {
	if inv.Tax == nil || inv.Tax.PricesInclude.IsEmpty() {
		return inv // nothing to do!
	}

	i2 := *inv
	i2.Lines = make([]*Line, len(inv.Lines))
	for i, l := range inv.Lines {
		i2.Lines[i] = l.removeIncludedTaxes(inv.Tax.PricesInclude, accuracy)
	}

	if len(inv.Discounts) > 0 {
		i2.Discounts = make([]*Discount, len(inv.Discounts))
		for i, l := range inv.Discounts {
			i2.Discounts[i] = l.removeIncludedTaxes(inv.Tax.PricesInclude, accuracy)
		}
	}
	if len(i2.Charges) > 0 {
		i2.Charges = make([]*Charge, len(inv.Charges))
		for i, l := range inv.Charges {
			i2.Charges[i] = l.removeIncludedTaxes(inv.Tax.PricesInclude, accuracy)
		}
	}

	tx := *i2.Tax
	tx.PricesInclude = ""
	i2.Tax = &tx

	return &i2
}

// ScenarioSummary determines a summary of the tax scenario for the invoice based on
// the document type and tax tags.
func (inv *Invoice) ScenarioSummary() *tax.ScenarioSummary {
	r := taxRegimeFor(inv.Supplier)
	if r == nil {
		return nil
	}
	return inv.scenarioSummary(r)
}

func (inv *Invoice) scenarioSummary(r *tax.Regime) *tax.ScenarioSummary {
	ss := r.ScenarioSet(ShortSchemaInvoice)
	if ss == nil {
		return nil
	}
	tags := []cbc.Key{}
	if inv.Tax != nil {
		tags = inv.Tax.Tags
	}
	return ss.SummaryFor(inv.Type, tags)
}

func (inv *Invoice) prepareTagsAndScenarios() error {
	r := taxRegimeFor(inv.Supplier)
	if r == nil {
		return nil
	}
	if inv.Tax == nil {
		return nil
	}

	// First check the tags are all valid
	for _, k := range inv.Tax.Tags {
		if t := r.Tag(k); t == nil {
			return fmt.Errorf("invalid document tag: %v", k)
		}
	}

	// Use the scenario summary to add any notes to the invoice
	ss := inv.scenarioSummary(r)
	if ss == nil {
		return nil
	}
	for _, n := range ss.Notes {
		// make sure we don't already have the same note in the invoice
		var en *cbc.Note
		for _, n2 := range inv.Notes {
			if n.Src == n2.Src {
				en = n
				break
			}
		}
		if en == nil {
			inv.Notes = append(inv.Notes, n)
		}
	}

	return nil
}

func (inv *Invoice) calculate(r *tax.Regime, tID *tax.Identity) error {
	date := inv.ValueDate
	if date == nil {
		date = &inv.IssueDate
	}
	if date == nil {
		return errors.New("issue date cannot be empty")
	}

	// Prepare the totals we'll need with amounts based on currency
	t := new(Totals)
	zero := r.CurrencyDef().BaseAmount()
	t.reset(zero)

	tls := make([]tax.TaxableLine, 0)

	// Ensure all the lines are up to date first
	for i, l := range inv.Lines {
		l.Index = i + 1
		l.calculate()

		// Basic sum
		t.Sum = t.Sum.Add(l.Total)
		tls = append(tls, l)
	}
	t.Total = t.Sum

	// Subtract discounts
	discounts := zero
	for i, l := range inv.Discounts {
		l.Index = i + 1
		if l.Percent != nil && !l.Percent.IsZero() {
			if l.Base == nil {
				l.Base = &t.Sum
			}
			l.Amount = l.Percent.Of(*l.Base)
		}
		discounts = discounts.Add(l.Amount)
		tls = append(tls, l)
	}
	if !discounts.IsZero() {
		t.Discount = &discounts
		t.Total = t.Total.Subtract(discounts)
	}

	// Add charges
	charges := zero
	for i, l := range inv.Charges {
		l.Index = i + 1
		if l.Percent != nil && !l.Percent.IsZero() {
			if l.Base == nil {
				l.Base = &t.Sum
			}
			l.Amount = l.Percent.Of(*l.Base)
		}
		charges = charges.Add(l.Amount)
		tls = append(tls, l)
	}
	if !charges.IsZero() {
		t.Charge = &charges
		t.Total = t.Total.Add(charges)
	}

	// Now figure out the tax totals (with some interface conversion)
	var pit cbc.Code
	if inv.Tax != nil && inv.Tax.PricesInclude != "" {
		pit = inv.Tax.PricesInclude
	}
	t.Taxes = new(tax.Total)
	tc := &tax.TotalCalculator{
		Zero:     zero,
		Regime:   r,
		Zone:     tID.Zone,
		Date:     *date,
		Includes: pit,
		Lines:    tls,
	}
	if err := tc.Calculate(t.Taxes); err != nil {
		return err
	}

	// Remove any included taxes from the total.
	ct := t.Taxes.Category(pit)
	if ct != nil {
		t.TaxIncluded = &ct.Amount
		t.Total = t.Total.Subtract(ct.Amount)
	}

	// Finally calculate the total with *all* the taxes.
	if inv.Tax != nil && inv.Tax.ContainsTag(common.TagReverseCharge) {
		t.Tax = zero
	} else {
		t.Tax = t.Taxes.Sum
	}
	t.TotalWithTax = t.Total.Add(t.Tax)
	t.Payable = t.TotalWithTax

	// Outlays
	if len(inv.Outlays) > 0 {
		t.Outlays = &zero
		for i, o := range inv.Outlays {
			o.Index = i + 1
			v := t.Outlays.Add(o.Amount)
			t.Outlays = &v
		}
		t.Payable = t.Payable.Add(*t.Outlays)
	}

	if inv.Payment != nil {
		inv.Payment.calculateAdvances(t.TotalWithTax)

		// Deal with advances, if any
		if t.Advances = inv.Payment.totalAdvance(zero); t.Advances != nil {
			v := t.Payable.Subtract(*t.Advances)
			t.Due = &v
		}

		// Calculate any due date amounts
		inv.Payment.Terms.CalculateDues(t.Payable)
	}

	inv.Totals = t
	return nil
}

func (inv *Invoice) determineTaxIdentity() *tax.Identity {
	if inv.Tax != nil {
		if inv.Tax.ContainsTag(common.TagCustomerRates) {
			if inv.Customer == nil {
				return nil
			}
			return inv.Customer.TaxID
		}
	}
	if inv.Supplier == nil {
		return nil
	}
	return inv.Supplier.TaxID
}

func taxRegimeFor(party *org.Party) *tax.Regime {
	if party == nil {
		return nil
	}
	tID := party.TaxID
	if tID == nil {
		return nil
	}
	return tax.RegimeFor(tID.Country, tID.Zone)
}

// Reset sets all the totals to the provided zero amount with the correct
// decimal places.
func (t *Totals) reset(zero num.Amount) {
	t.Sum = zero
	t.Discount = nil
	t.Charge = nil
	t.TaxIncluded = nil
	t.Total = zero
	t.Taxes = nil
	t.Tax = zero
	t.TotalWithTax = zero
	t.Outlays = nil
	t.Payable = zero
	t.Advances = nil
	t.Due = nil
}

// JSONSchemaExtend extends the schema with additional property details
func (Invoice) JSONSchemaExtend(schema *jsonschema.Schema) {
	props := schema.Properties
	if val, ok := props.Get("type"); ok {
		its := val.(*jsonschema.Schema)
		its.OneOf = make([]*jsonschema.Schema, len(InvoiceTypes))
		for i, v := range InvoiceTypes {
			its.OneOf[i] = &jsonschema.Schema{
				Const:       v.Key.String(),
				Description: v.Description,
			}
		}
	}
}
