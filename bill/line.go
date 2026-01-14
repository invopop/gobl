package bill

import (
	"context"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// Line is a single row in an invoice.
type Line struct {
	uuid.Identify
	// Line number inside the parent (calculated)
	Index int `json:"i" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`
	// Number of items
	Quantity num.Amount `json:"quantity" jsonschema:"title=Quantity"`
	// Single identifier provided by the supplier for an object on which the
	// line item is based and is not considered a universal identity. Examples
	// include a subscription number, telephone number, meter point, etc.
	// Utilize the label property to provide a description of the identifier.
	Identifier *org.Identity `json:"identifier,omitempty" jsonschema:"title=Identifier"`
	// A period of time relevant to when the service or item is delivered.
	Period *cal.Period `json:"period,omitempty" jsonschema:"title=Period"`
	// Order reference for a specific line within a purchase order provided by the buyer.
	Order cbc.Code `json:"order,omitempty" jsonschema:"title=Order Reference"`
	// Buyer accounting reference cost code to associate with the line.
	Cost cbc.Code `json:"cost,omitempty" jsonschema:"title=Cost Reference"`
	// Details about the item, service or good, that is being sold
	Item *org.Item `json:"item" jsonschema:"title=Item"`
	// Breakdown of the line item for more detailed information. The sum of all lines
	// will be used for the item price.
	Breakdown []*SubLine `json:"breakdown,omitempty" jsonschema:"title=Breakdown"`
	// Result of quantity multiplied by the item's price (calculated)
	Sum *num.Amount `json:"sum,omitempty" jsonschema:"title=Sum" jsonschema_extras:"calculated=true"`
	// Discounts applied to this line
	Discounts []*LineDiscount `json:"discounts,omitempty" jsonschema:"title=Discounts"`
	// Charges applied to this line
	Charges []*LineCharge `json:"charges,omitempty" jsonschema:"title=Charges"`
	// Map of taxes to be applied and used in the invoice totals
	Taxes tax.Set `json:"taxes,omitempty" jsonschema:"title=Taxes"`
	// Total line amount after applying discounts to the sum (calculated).
	Total *num.Amount `json:"total,omitempty" jsonschema:"title=Total"  jsonschema_extras:"calculated=true"`

	// List of substituted lines. Useful for deliveries or corrective documents in order
	// to indicate to the recipient which of the requested lines are being replaced.
	// This is for purely informative purposes, and will not be used for calculations.
	Substituted []*SubLine `json:"substituted,omitempty" jsonschema:"title=Substituted"`

	// Seller of the item if different from the supplier or ordering seller. This can be
	// useful for marketplace or drop-ship scenarios in locations that require the
	// original seller to be indicated.
	Seller *org.Party `json:"seller,omitempty" jsonschema:"title=Seller"`

	// Set of specific notes for this line that may be required for
	// clarification.
	Notes []*org.Note `json:"notes,omitempty" jsonschema:"title=Notes"`

	// Extension codes that apply to the line
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
}

// SubLine provides a simplified line that can be embedded inside other lines
// to provide a more detailed breakdown of the items being sold or substituted.
type SubLine struct {
	uuid.Identify
	// Line number inside the parent (calculated)
	Index int `json:"i" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`
	// Number of items
	Quantity num.Amount `json:"quantity" jsonschema:"title=Quantity"`
	// Single identifier provided by the supplier for an object on which the
	// line item is based and is not considered a universal identity. Examples
	// include a subscription number, telephone number, meter point, etc.
	// Utilize the label property to provide a description of the identifier.
	Identifier *org.Identity `json:"identifier,omitempty" jsonschema:"title=Identifier"`
	// A period of time relevant to when the service or item is delivered.
	Period *cal.Period `json:"period,omitempty" jsonschema:"title=Period"`
	// Order reference for a specific line within a purchase order provided by the buyer.
	Order cbc.Code `json:"order,omitempty" jsonschema:"title=Order Reference"`
	// Buyer accounting reference cost code to associate with the line.
	Cost cbc.Code `json:"cost,omitempty" jsonschema:"title=Cost Reference"`
	// Details about the item, service or good, that is being sold
	Item *org.Item `json:"item" jsonschema:"title=Item"`
	// Result of quantity multiplied by the item's price (calculated)
	Sum *num.Amount `json:"sum,omitempty" jsonschema:"title=Sum" jsonschema_extras:"calculated=true"`
	// Discounts applied to this sub-line
	Discounts []*LineDiscount `json:"discounts,omitempty" jsonschema:"title=Discounts"`
	// Charges applied to this sub-line
	Charges []*LineCharge `json:"charges,omitempty" jsonschema:"title=Charges"`
	// Total sub-line amount after applying discounts to the sum (calculated).
	Total *num.Amount `json:"total,omitempty" jsonschema:"title=Total"  jsonschema_extras:"calculated=true"`
	// Set of specific notes for this sub-line that may be required for
	// clarification.
	Notes []*org.Note `json:"notes,omitempty" jsonschema:"title=Notes"`
}

// GetTaxes responds with the array of tax rates applied to this line.
// This implements the tax.TaxableLine interface.
func (l *Line) GetTaxes() tax.Set {
	return l.Taxes
}

// GetTotal provides the final total for this line, excluding any tax calculations.
// This implements the tax.TaxableLine interface.
func (l *Line) GetTotal() num.Amount {
	if l.Total == nil {
		return num.AmountZero
	}
	return *l.Total
}

// Validate performs a validation check on the line without a context.
func (l *Line) Validate() error {
	return l.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures the line contains everything required using
// the provided context that should include the regime.
func (l *Line) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, l,
		validation.Field(&l.UUID),
		validation.Field(&l.Index, validation.Required),
		validation.Field(&l.Quantity, validation.Required),
		validation.Field(&l.Identifier),
		validation.Field(&l.Period),
		validation.Field(&l.Order),
		validation.Field(&l.Cost),
		validation.Field(&l.Item, validation.Required),
		validation.Field(&l.Breakdown),
		validation.Field(&l.Sum,
			validation.When(
				l.Item != nil && l.Item.Price != nil,
				validation.Required,
			),
		),
		validation.Field(&l.Discounts),
		validation.Field(&l.Charges),
		validation.Field(&l.Taxes),
		validation.Field(&l.Total,
			validation.When(
				l.Item != nil && l.Item.Price != nil,
				validation.Required,
			),
		),
		validation.Field(&l.Substituted),
		validation.Field(&l.Seller),
		validation.Field(&l.Notes),
	)
}

// ValidateWithContext ensures the line contains everything required using
// the provided context that should include the regime.
func (sl *SubLine) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, sl,
		validation.Field(&sl.UUID),
		validation.Field(&sl.Index, validation.Required),
		validation.Field(&sl.Quantity, validation.Required),
		validation.Field(&sl.Identifier),
		validation.Field(&sl.Period),
		validation.Field(&sl.Order),
		validation.Field(&sl.Cost),
		validation.Field(&sl.Item, validation.Required),
		validation.Field(&sl.Sum,
			validation.When(
				sl.Item != nil && sl.Item.Price != nil,
				validation.Required,
			),
		),
		validation.Field(&sl.Discounts),
		validation.Field(&sl.Charges),
		validation.Field(&sl.Total,
			validation.When(
				sl.Item != nil && sl.Item.Price != nil,
				validation.Required,
			),
		),
		validation.Field(&sl.Notes),
	)
}

// Normalize performs normalization on the line and embedded objects using the
// provided list of normalizers.
func (l *Line) Normalize(normalizers tax.Normalizers) {
	if l == nil {
		return
	}
	normalizeLineItemPrice(l)
	l.Taxes = tax.CleanSet(l.Taxes)
	l.Discounts = CleanLineDiscounts(l.Discounts)
	l.Charges = CleanLineCharges(l.Charges)
	tax.Normalize(normalizers, l.Identifier)
	tax.Normalize(normalizers, l.Taxes)
	tax.Normalize(normalizers, l.Item)
	tax.Normalize(normalizers, l.Breakdown)
	tax.Normalize(normalizers, l.Discounts)
	tax.Normalize(normalizers, l.Charges)
	tax.Normalize(normalizers, l.Substituted)
	tax.Normalize(normalizers, l.Seller)
	normalizers.Each(l)
}

// Normalize performs normalization on the subline and embedded objects using the
// provided list of normalizers.
func (sl *SubLine) Normalize(normalizers tax.Normalizers) {
	normalizeSubLineItemPrice(sl)
	sl.Discounts = CleanLineDiscounts(sl.Discounts)
	sl.Charges = CleanLineCharges(sl.Charges)
	tax.Normalize(normalizers, sl.Identifier)
	tax.Normalize(normalizers, sl.Item)
	tax.Normalize(normalizers, sl.Discounts)
	tax.Normalize(normalizers, sl.Charges)
	normalizers.Each(sl)
}

func normalizeLineItemPrice(l *Line) {
	if l == nil || l.Item == nil || l.Item.Price == nil {
		return
	}
	i := l.Item
	if i.Price.IsNegative() {
		p := i.Price.Negate()
		i.Price = &p
		l.Quantity = l.Quantity.Negate()
	}
}

func normalizeSubLineItemPrice(sl *SubLine) {
	if sl == nil || sl.Item == nil || sl.Item.Price == nil {
		return
	}
	i := sl.Item
	if i.Price.IsNegative() {
		p := i.Price.Negate()
		i.Price = &p
		sl.Quantity = sl.Quantity.Negate()
	}
}

func removeLineIncludedTaxes(line *Line, cat cbc.Code) *Line {
	accuracy := defaultTaxRemovalAccuracy
	rate := line.Taxes.Get(cat)
	if rate == nil || rate.Percent == nil {
		return line
	}

	l2 := *line
	l2i := *line.Item

	l2i.AltPrices = nil // empty alternative prices
	price := line.Item.Price.Upscale(accuracy).Remove(*rate.Percent)
	l2i.Price = &price
	// assume sum and total will be calculated automatically

	l2.Breakdown = removeSubLinesIncludedTaxes(line.Breakdown, rate, accuracy)
	l2.Discounts = removeLineDiscountsIncludedTaxes(line.Discounts, rate, accuracy)
	l2.Charges = removeLineChargesIncludedTaxes(line.Charges, rate, accuracy)
	l2.Substituted = removeSubLinesIncludedTaxes(line.Substituted, rate, accuracy)
	l2.Item = &l2i

	return &l2
}

func removeSubLinesIncludedTaxes(sls []*SubLine, tc *tax.Combo, exp uint32) []*SubLine {
	if len(sls) == 0 {
		return nil
	}
	rows := make([]*SubLine, len(sls))
	for i, sl := range sls {
		sl2 := *sl
		sl2i := *sl.Item
		sl2i.AltPrices = nil
		price := sl.Item.Price.Upscale(exp).Remove(*tc.Percent)
		sl2i.Price = &price
		sl2.Discounts = removeLineDiscountsIncludedTaxes(sl.Discounts, tc, exp)
		sl2.Charges = removeLineChargesIncludedTaxes(sl.Charges, tc, exp)
		sl2.Item = &sl2i
		rows[i] = &sl2
	}
	return rows
}

func removeLineDiscountsIncludedTaxes(discounts []*LineDiscount, tc *tax.Combo, exp uint32) []*LineDiscount {
	if len(discounts) == 0 {
		return nil
	}
	rows := make([]*LineDiscount, len(discounts))
	for i, v := range discounts {
		d := *v
		d.Amount = d.Amount.Upscale(exp).Remove(*tc.Percent)
		rows[i] = &d
	}
	return rows
}

func removeLineChargesIncludedTaxes(charges []*LineCharge, tc *tax.Combo, exp uint32) []*LineCharge {
	if len(charges) == 0 {
		return nil
	}
	rows := make([]*LineCharge, len(charges))
	for i, v := range charges {
		d := *v
		d.Amount = d.Amount.Upscale(exp).Remove(*tc.Percent)
		rows[i] = &d
	}
	return rows
}

func lineItemHasPrice(value any) error {
	line, ok := value.(*Line)
	if line == nil || !ok {
		return nil
	}
	return validation.ValidateStruct(line,
		validation.Field(&line.Item,
			org.ItemPriceRequired(),
			validation.Skip,
		),
	)
}
