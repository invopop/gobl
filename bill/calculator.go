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

func calculate(doc billable) error {
	r := doc.RegimeDef() // may be nil!
	calculateIssueDateAndTime(r, doc)

	// Get the date used for tax calculations
	date := doc.GetValueDate()
	if date == nil {
		id := doc.GetIssueDate()
		date = &id
	}

	// Convert empty or invalid currency to the regime's currency
	if doc.GetCurrency() == currency.CodeEmpty || doc.GetCurrency().Def() == nil {
		if r == nil {
			return validation.Errors{"currency": errors.New("missing")}
		}
		doc.SetCurrency(r.Currency)
	}
	cur := doc.GetCurrency()

	t := doc.GetTotals()
	// Prepare the totals we'll need with amounts based on currency
	if t == nil {
		t = new(Totals)
	}
	zero := cur.Def().Zero()
	t.reset(zero)

	// Figure out rounding rules and if prices include tax early
	var pit cbc.Code
	var rr cbc.Key
	if tx := doc.GetTax(); tx != nil {
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
	if err := calculateComplements(doc.GetComplements()); err != nil {
		return validation.Errors{"complements": err}
	}

	// Preceding
	calculateOrgDocumentRefs(doc.GetPreceding(), cur, rr)

	// Lines
	if err := calculateLines(doc.GetLines(), cur, doc.GetExchangeRates(), rr); err != nil {
		return validation.Errors{"lines": err}
	}
	t.Sum = calculateLineSum(doc.GetLines(), cur)
	t.Total = t.Sum

	// Discount Lines
	calculateDiscounts(doc.GetDiscounts(), cur, t.Sum, rr)
	if discounts := calculateDiscountSum(doc.GetDiscounts(), cur); discounts != nil {
		t.Discount = discounts
		t.Total = t.Total.Subtract(*discounts)
	}

	// Charge Lines
	calculateCharges(doc.GetCharges(), cur, t.Sum, rr)
	if charges := calculateChargeSum(doc.GetCharges(), cur); charges != nil {
		t.Charge = charges
		t.Total = t.Total.Add(*charges)
	}

	// Build list of taxable lines
	tls := make([]tax.TaxableLine, 0)
	for _, l := range doc.GetLines() {
		if l.Total != nil {
			tls = append(tls, l)
		}
	}
	for _, l := range doc.GetDiscounts() {
		tls = append(tls, l)
	}
	for _, l := range doc.GetCharges() {
		tls = append(tls, l)
	}

	if len(tls) == 0 {
		// This applies for orders and deliveries that might not have
		// any pricing details.
		doc.SetTotals(nil)
		return nil
	}

	// Now figure out the tax totals
	t.Taxes = new(tax.Total)
	tc := &tax.TotalCalculator{
		Currency: doc.GetCurrency(),
		Rounding: rr,
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

	// Before calculating the amount due and advances, we need to round
	// everything. Payments reflect real monetary values and can never
	// be fractions of the currency.
	roundLines(doc.GetLines())
	roundDiscounts(doc.GetDiscounts(), cur)
	roundCharges(doc.GetCharges(), cur)
	t.round(zero)

	if t.Rounding != nil {
		// BT-144 in EN16931
		t.Payable = t.Payable.Add(*t.Rounding)
	}
	if pd := doc.GetPaymentDetails(); pd != nil {
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
	doc.SetTotals(t)

	return nil
}

func calculateIssueDateAndTime(r *tax.RegimeDef, doc billable) {
	tz := r.TimeLocation()
	if doc.GetIssueTime() != nil && doc.GetIssueTime().IsZero() {
		dn := cal.ThisSecondIn(tz)
		tn := dn.Time()
		doc.SetIssueDate(dn.Date())
		doc.SetIssueTime(&tn)
	} else if doc.GetIssueDate().IsZero() {
		doc.SetIssueDate(cal.TodayIn(tz))
	}
}

func calculateOrgDocumentRefs(drs []*org.DocumentRef, cur currency.Code, rr cbc.Key) {
	for _, drs := range drs {
		if drs.Currency != currency.CodeEmpty {
			cur = drs.Currency
		}
		drs.Calculate(cur, rr)
	}
}

func canRemoveIncludedTaxes(doc billable) bool {
	return doc.GetTax() != nil && !doc.GetTax().PricesInclude.IsEmpty()
}

func removeIncludedTaxes(doc billable) error {
	if !canRemoveIncludedTaxes(doc) {
		return nil
	}
	tpi := doc.GetTax().PricesInclude

	totalWithTax := doc.GetTotals().TotalWithTax

	doc.SetTotals(new(Totals))
	lines := doc.GetLines()
	for i, l := range doc.GetLines() {
		lines[i] = removeLineIncludedTaxes(l, tpi)
	}

	discounts := doc.GetDiscounts()
	if len(discounts) > 0 {
		for i, l := range discounts {
			discounts[i] = l.removeIncludedTaxes(tpi)
		}
	}
	charges := doc.GetCharges()
	if len(charges) > 0 {
		for i, l := range charges {
			charges[i] = l.removeIncludedTaxes(tpi)
		}
	}

	tx := doc.GetTax()
	tx.PricesInclude = ""

	if err := calculate(doc); err != nil {
		return err
	}

	// Account for any rounding errors that we just can't handle
	t := doc.GetTotals()
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
	if doc.GetCustomer() == nil || doc.GetCustomer().TaxID == nil {
		return
	}
	country := doc.GetCustomer().TaxID.Country
	for _, l := range doc.GetLines() {
		addCountryToTaxes(l.Taxes, country)
	}
	for _, d := range doc.GetDiscounts() {
		addCountryToTaxes(d.Taxes, country)
	}
	for _, c := range doc.GetCharges() {
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
