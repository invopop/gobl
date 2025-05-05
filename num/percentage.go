package num

import (
	"github.com/invopop/jsonschema"
)

var (
	factor1   = MakeAmount(1, 0)
	factor100 = MakeAmount(100, 0)
)

// Percentage wraps around the regular Amount handler to provide support
// for percentage values, especially useful for tax rates.
type Percentage struct {
	amount Amount
}

var (
	// PercentageZero is a convenience variable for testing against zero percentages.
	PercentageZero = MakePercentage(0, 0)
)

// NewPercentage provides a new pointer to a Percentage value.
// Using [MakePercentage] is recommended, but this is useful for handling
// nil values.
func NewPercentage(value int64, exp uint32) *Percentage {
	p := MakePercentage(value, exp)
	return &p
}

// MakePercentage will make a new Percentage instance with the provided
// value and exponent.
//
// Example: To make a 16% rate, you would use
//
//	num.MakePercentage(16, 2)
func MakePercentage(value int64, exp uint32) Percentage {
	return Percentage{Amount{value: value, exp: exp}}
}

// PercentageFromString builds a percentage value from a provided string.
// The % symbol will be removed automatically and rescale the stored amount
// to make future calculations easier. For example, the following two strings
// will be interpreted equally:
//
//   - `0.160`
//   - `16.0%`
func PercentageFromString(str string) (Percentage, error) {
	p := Percentage{}
	l := len(str)
	if l == 0 {
		return p, nil
	}
	rescale := false
	if str[l-1:] == "%" {
		str = str[:l-1]
		rescale = true
	}

	var err error
	p.amount, err = AmountFromString(str)
	if err != nil {
		return p, err
	}
	if rescale {
		return PercentageFromAmount(p.amount), nil
	}

	return p, nil
}

// PercentageFromAmount provides the percentage value of the amount ensuring it
// is correctly scaled.
func PercentageFromAmount(a Amount) Percentage {
	a2 := a.Rescale(a.exp + 2).Divide(factor100)
	return Percentage{amount: a2}
}

// Value provides the percentage amount's value
func (p Percentage) Value() int64 {
	return p.amount.value
}

// Exp provides the percentage amount's exponent value.
func (p Percentage) Exp() uint32 {
	return p.amount.exp
}

// String outputs the percentage value in a human readable way including
// the percentage symbol.
func (p Percentage) String() string {
	return p.StringWithoutSymbol() + "%"
}

// StringWithoutSymbol provides the percent value without a percent symbol.
func (p Percentage) StringWithoutSymbol() string {
	return p.Amount().String()
}

// Base provides the underlying amount value of the percentage which is stored
// internally without any factors applied.
func (p Percentage) Base() Amount {
	return p.amount
}

// Amount provides an amount for the percentage that has been rescaled
// from the underlying value mainly to be used for formatting.
func (p Percentage) Amount() Amount {
	e := int64(p.amount.exp) - 2
	if e < 0 {
		e = 0
	}
	return p.amount.Multiply(factor100).Rescale(uint32(e))
}

// Rescale will rescale the percentage value to the provided exponent.
func (p Percentage) Rescale(exp uint32) Percentage {
	return Percentage{amount: p.amount.Rescale(exp)}
}

// Of calculates the "percent of" the provided amount. The exponent of the
// provided amount is used.
func (p Percentage) Of(a Amount) Amount {
	return a.Multiply(p.amount)
}

// From calculates what "percent from" the provided amount would result
// assuming the rate has already been applied.
func (p Percentage) From(a Amount) Amount {
	x := a.Divide(p.Factor())
	return a.Subtract(x)
}

// Factor returns the percentage amount as a factor, essentially
// adding 1 to the rate.
func (p Percentage) Factor() Amount {
	return p.amount.Add(factor1)
}

// Equals wraps around the amount comparison to see if the two percentages
// have the same value.
func (p Percentage) Equals(p2 Percentage) bool {
	return p.amount.Equals(p2.amount)
}

// Compare two percentages and return an integer value according to the
// sign of the difference:
//
//	-1 if a <  a2
//	 0 if a == a2
//	 1 if a >  a2
func (p Percentage) Compare(p2 Percentage) int {
	return p.amount.Compare(p2.amount)
}

// IsZero checks if the percentage is zero.
func (p Percentage) IsZero() bool {
	return p.amount.IsZero()
}

// IsPositive checks if the percentage is positive.
func (p Percentage) IsPositive() bool {
	return p.amount.IsPositive()
}

// IsNegative checks if the percentage is negative.
func (p Percentage) IsNegative() bool {
	return p.amount.IsNegative()
}

// MarshalText provides the byte value of the amount. See also the
// String() method.
func (p Percentage) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

// Invert provides a new percentage value that is the inverse of the
// current percentage.
// Deprectatd: Use Negate instead.
func (p Percentage) Invert() Percentage {
	return p.Negate()
}

// Negate provides a new percentage value that is the negative inverse of the
// current percentage.
func (p Percentage) Negate() Percentage {
	return Percentage{amount: p.amount.Invert()}
}

// UnmarshalText will decode the percentage value, even if it is quoted
// as a string.
func (p *Percentage) UnmarshalText(value []byte) error {
	if string(value) == "null" {
		return nil
	}
	result, err := PercentageFromString(string(value))
	if err != nil {
		return err
	}
	*p = result
	return nil
}

// UnmarshalJSON ensures percentages will be parsed even if defined as
// numbers in the source JSON.
func (p *Percentage) UnmarshalJSON(value []byte) error {
	return p.UnmarshalText(unquote(value))
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (Percentage) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Pattern:     `^\-?[0-9]+(\.[0-9]+)?%$`,
		Title:       "Percentage",
		Description: "Similar to an Amount, but designed for percentages and includes % symbol in JSON output.",
	}
}
