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
	getIssueTime() *cal.Time
	getValueDate() *cal.Date
	getTax() *Tax
	getPreceding() []*org.DocumentRef
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
	setIssueTime(*cal.Time)
	setCurrency(currency.Code)
	setTotals(*Totals)
}

func calculate(doc billable) error {
	r := doc.RegimeDef() // may be nil!
	date := calculateIssueDateAndTime(r, doc)

	// Convert empty or invalid currency to the regime's currency
	if doc.getCurrency() == currency.CodeEmpty || doc.getCurrency().Def() == nil {
		if r == nil {
			return validation.Errors{"currency": errors.New("missing")}
		}
		doc.setCurrency(r.Currency)
	}
	cur := doc.getCurrency()

	t := doc.getTotals()
	// Prepare the totals we'll need with amounts based on currency
	if t == nil {
		t = new(Totals)
	}
	zero := cur.Def().Zero()
	t.reset(zero)

	// Figure out rounding rules and if prices include tax early
	var pit cbc.Code
	var rr cbc.Key
	if tx := doc.getTax(); tx != nil {
		if tx.PricesInclude != "" {
			pit = tx.PricesInclude
		}
		if tx.Rounding != "" {
			rr = tx.Rounding
		}
	}
	if rr == "" {
		rr = r.GetRoundingRule()
	}

	// Do we need to deal with the customer-rates tag?
	if doc.HasTags(tax.TagCustomerRates) {
		applyCustomerRates(doc)
	}

	// Complements
	if err := calculateComplements(doc.getComplements()); err != nil {
		return validation.Errors{"complements": err}
	}

	// Preceding
	calculateOrgDocumentRefs(doc.getPreceding(), cur, rr)

	// Lines
	if err := calculateLines(doc.getLines(), cur, doc.getExchangeRates(), rr); err != nil {
		return validation.Errors{"lines": err}
	}
	t.Sum = calculateLineSum(doc.getLines(), cur)
	t.Total = t.Sum

	// Discount Lines
	calculateDiscounts(doc.getDiscounts(), cur, t.Sum, rr)
	if discounts := calculateDiscountSum(doc.getDiscounts(), cur); discounts != nil {
		t.Discount = discounts
		t.Total = t.Total.Subtract(*discounts)
	}

	// Charge Lines
	calculateCharges(doc.getCharges(), cur, t.Sum, rr)
	if charges := calculateChargeSum(doc.getCharges(), cur); charges != nil {
		t.Charge = charges
		t.Total = t.Total.Add(*charges)
	}

	tls := prepareTaxableLines(doc)
	if len(tls) == 0 {
		// This applies for orders and deliveries that might not have
		// any pricing details.
		doc.setTotals(nil)
		return nil
	}

	// Now figure out the tax totals
	t.Taxes = new(tax.Total)
	tc := &tax.TotalCalculator{
		Currency: doc.getCurrency(),
		Rounding: rr,
		Country:  r.GetCountry(),
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
		ti := ct.Amount
		t.TaxIncluded = &ti
		t.Total = t.Total.Subtract(ti)
	}

	// Calculate the total with *all* the taxes.
	t.Tax = t.Taxes.Sum
	t.TotalWithTax = t.Total.Add(t.Tax)
	if t.Taxes.Retained != nil {
		t.RetainedTax = t.Taxes.Retained
	}
	t.Payable = t.TotalWithTax
	if t.RetainedTax != nil {
		t.Payable = t.Payable.Subtract(*t.RetainedTax)
	}
	// Remove taxes object if it doesn't contain any categories
	if len(t.Taxes.Categories) == 0 {
		t.Taxes = nil
	}

	roundTotalsAndPreparePayments(doc, cur, t)

	return nil
}

func roundTotalsAndPreparePayments(doc billable, cur currency.Code, t *Totals) {
	zero := cur.Def().Zero()
	// Before calculating the amount due and advances, we need to round
	// everything. Payments reflect real monetary values and can never
	// be fractions of the currency.
	roundLines(doc.getLines())
	roundDiscounts(doc.getDiscounts(), cur)
	roundCharges(doc.getCharges(), cur)
	t.round(zero)

	if t.Rounding != nil {
		// BT-144 in EN16931
		t.Payable = t.Payable.Add(*t.Rounding)
	}
	if pd := doc.getPaymentDetails(); pd != nil {
		pd.calculateAdvances(zero, t.Payable)
		// Deal with advances, if any. Note that in the current
		// implementation multiple percentage advances are likely to
		// suffer rounding errors. It usually better for users to use
		// fixed payment amounts if possible.
		if t.Advances = pd.totalAdvance(zero); t.Advances != nil {
			v := t.Payable.Subtract(*t.Advances)
			t.Due = &v
		}
		// Calculate any due date amounts
		pd.Terms.CalculateDues(zero, t.Payable)
	}
	doc.setTotals(t)
}

func prepareTaxableLines(doc billable) []tax.TaxableLine {
	// Build list of taxable lines
	tls := make([]tax.TaxableLine, 0)
	for _, l := range doc.getLines() {
		if l != nil && l.Total != nil {
			tls = append(tls, l)
		}
	}
	for _, l := range doc.getDiscounts() {
		if l != nil {
			tls = append(tls, l)
		}
	}
	for _, l := range doc.getCharges() {
		if l != nil {
			tls = append(tls, l)
		}
	}
	return tls
}

func calculateIssueDateAndTime(r *tax.RegimeDef, doc billable) *cal.Date {
	tz := r.TimeLocation()
	if doc.getIssueTime() != nil && doc.getIssueTime().IsZero() {
		dn := cal.ThisSecondIn(tz)
		tn := dn.Time()
		doc.setIssueDate(dn.Date())
		doc.setIssueTime(&tn)
	} else if doc.getIssueDate().IsZero() {
		doc.setIssueDate(cal.TodayIn(tz))
	}

	// Get the date used for tax calculations
	date := doc.getValueDate()
	if date == nil {
		id := doc.getIssueDate()
		date = &id
	}

	return date
}

func calculateOrgDocumentRefs(drs []*org.DocumentRef, cur currency.Code, rr cbc.Key) {
	for _, drs := range drs {
		if drs == nil {
			continue
		}
		if drs.Currency != currency.CodeEmpty {
			cur = drs.Currency
		}
		drs.Calculate(cur, rr)
	}
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
		if c == nil {
			continue
		}
		if err := c.Calculate(); err != nil {
			return err
		}
	}
	return nil
}
