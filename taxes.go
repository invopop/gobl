package gobl

// TaxCategory defines a grouping of taxes whereby only one
// definition inside a tax category can be applied to a given
// invoice line.
type TaxCategory string

// TaxCode defines a simple code used to describe the tax
// for the given region.
type TaxCode string

// TaxDef defines a tax combination of code and rate.
type TaxDef struct {
	Name     string      `json:"name" jsconschema:"title=Name"`
	Category TaxCategory `json:"category"`
	Code     TaxCode     `json:"code" jsonschema:"title=Code"`
	Rates    []TaxRate   `json:"rates" jsonschema:"title=Rate"`
	Included bool        `json:"included,omitempty"`
}

// TaxRate contains a percentage tax rate for a given date range.
// Fiscal policy changes meen that rates are not fixed so we need to
// be able to apply the correct rate for a given period.
type TaxRate struct {
	From     Date   `json:"from,omitempty"`
	Upto     Date   `json:"upto,omitempty"`
	Amount   Amount `json:"amount"`
	Disabled bool   `json:"disabled,omitempty"`
}

// RateOn determines the tax rate for the provided date.
func (td *TaxDef) RateOn(date Date) Amount {

}

// Tax represents a tax calculation for a given TaxDef and base.
type Tax struct {
	Code   TaxCode `json:"code"`
	Base   Amount  `json:"base"`
	Rate   Amount  `json:"rate"`
	Amount Amount  `json:"amount" jsonschema:"title=Amount,description=Result after applying rate to the base."`
}

// TaxTotal represents the total summary amounts of tax contained
// in the document (usually an Invoice).
type TaxTotal struct {
}
