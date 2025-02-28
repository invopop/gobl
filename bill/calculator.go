package bill

import (
	"errors"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// billable defines the methods required to be able to perform calculations and
// other operations on a bill document with a common basic structure.
type billable interface {
	// From tax.Regime
	RegimeDef() *tax.RegimeDef

	// From tax.Tags
	HasTags(tags ...cbc.Key) bool
	GetTags() []cbc.Key

	getIssueDate() cal.Date
	getValueDate() *cal.Date
	getTax() *Tax
	getCustomer() *org.Party
	getCurrency() currency.Code
	getExchangeRates() []*currency.ExchangeRate
	getLines() []*Line
	getDiscounts() []*Discount
	getCharges() []*Charge
	getPaymentDetails() *PaymentDetails
	getTotals() *Totals
	getComplements() []*schema.Object

	setIssueDate(cal.Date)
	setCurrency(currency.Code)
	setTotals(*Totals)
}

func calculate(doc billable) error {
	r := doc.RegimeDef() // may be nil!

	// Normalize data
	if doc.getIssueDate().IsZero() {
		doc.setIssueDate(cal.TodayIn(r.TimeLocation()))
	}
	date := doc.getValueDate()
	if date == nil {
		id := doc.getIssueDate()
		date = &id
	}

	// Convert empty or invalid currency to the regime's currency
	if doc.getCurrency() == currency.CodeEmpty || doc.getCurrency().Def() == nil {
		if r == nil {
			return validation.Errors{"currency": errors.New("missing")}
		}
		doc.setCurrency(r.Currency)
	}

	t := doc.getTotals()
	// Prepare the totals we'll need with amounts based on currency
	if t == nil {
		t = new(Totals)
	}
	zero := doc.getCurrency().Def().Zero()
	t.reset(zero)

	// Do we need to deal with the customer-rates tag?
	if doc.HasTags(tax.TagCustomerRates) {
		applyCustomerRates(doc)
	}

	// Complements
	if err := calculateComplements(doc.getComplements()); err != nil {
		return validation.Errors{"complements": err}
	}

	// Lines
	if err := calculateLines(doc.getLines(), doc.getCurrency(), doc.getExchangeRates()); err != nil {
		return validation.Errors{"lines": err}
	}
	t.Sum = calculateLineSum(doc.getLines(), doc.getCurrency())
	t.Total = t.Sum

	// Discount Lines
	calculateDiscounts(doc.getDiscounts(), t.Sum, zero)
	if discounts := calculateDiscountSum(doc.getDiscounts(), zero); discounts != nil {
		t.Discount = discounts
		t.Total = t.Total.Subtract(*discounts)
	}

	// Charge Lines
	calculateCharges(doc.getCharges(), t.Sum, zero)
	if charges := calculateChargeSum(doc.getCharges(), zero); charges != nil {
		t.Charge = charges
		t.Total = t.Total.Add(*charges)
	}

	// Build list of taxable lines
	tls := make([]tax.TaxableLine, 0)
	for _, l := range doc.getLines() {
		if l.Total != nil {
			tls = append(tls, l)
		}
	}
	for _, l := range doc.getDiscounts() {
		tls = append(tls, l)
	}
	for _, l := range doc.getCharges() {
		tls = append(tls, l)
	}

	if len(tls) == 0 {
		// This applies for orders and deliveries that might not have
		// any pricing details.
		doc.setTotals(nil)
		return nil
	}

	// Now figure out the tax totals
	var pit cbc.Code
	if doc.getTax() != nil && doc.getTax().PricesInclude != "" {
		pit = doc.getTax().PricesInclude
	}
	t.Taxes = new(tax.Total)
	tc := &tax.TotalCalculator{
		Currency: doc.getCurrency(),
		Rounding: r.GetRoundingRule(),
		Country:  r.GetCountry(),
		Tags:     doc.GetTags(),
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

	if pd := doc.getPaymentDetails(); pd != nil {
		pd.calculateAdvances(zero, t.TotalWithTax)

		// Deal with advances, if any
		if t.Advances = pd.totalAdvance(zero); t.Advances != nil {
			v := t.Payable.Subtract(*t.Advances)
			t.Due = &v
		}

		// Calculate any due date amounts
		pd.Terms.CalculateDues(zero, t.Payable)
	}

	t.round(zero)
	doc.setTotals(t)

	return nil
}

func canRemoveIncludedTaxes(doc billable) bool {
	return doc.getTax() != nil && !doc.getTax().PricesInclude.IsEmpty()
}

func removeIncludedTaxes(doc billable) error {
	if !canRemoveIncludedTaxes(doc) {
		return nil
	}
	tpi := doc.getTax().PricesInclude

	totalWithTax := doc.getTotals().TotalWithTax

	doc.setTotals(new(Totals))
	lines := doc.getLines()
	for i, l := range doc.getLines() {
		lines[i] = removeLineIncludedTaxes(l, tpi)
	}

	discounts := doc.getDiscounts()
	if len(discounts) > 0 {
		for i, l := range discounts {
			discounts[i] = l.removeIncludedTaxes(tpi)
		}
	}
	charges := doc.getCharges()
	if len(charges) > 0 {
		for i, l := range charges {
			charges[i] = l.removeIncludedTaxes(tpi)
		}
	}

	tx := doc.getTax()
	tx.PricesInclude = ""

	if err := calculate(doc); err != nil {
		return err
	}

	// Account for any rounding errors that we just can't handle
	t := doc.getTotals()
	if !totalWithTax.Equals(t.TotalWithTax) {
		rnd := totalWithTax.Subtract(t.TotalWithTax)
		t.Rounding = &rnd
		if err := calculate(doc); err != nil {
			return err
		}
	}

	return nil
}

func applyCustomerRates(doc billable) {
	if doc.getCustomer() == nil || doc.getCustomer().TaxID == nil {
		return
	}
	country := doc.getCustomer().TaxID.Country
	for _, l := range doc.getLines() {
		addCountryToTaxes(l.Taxes, country)
	}
	for _, d := range doc.getDiscounts() {
		addCountryToTaxes(d.Taxes, country)
	}
	for _, c := range doc.getCharges() {
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
