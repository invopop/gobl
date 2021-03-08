package num

import "fmt"

var factor100 = Amount{Value: 100, Exp: 0}

// Percentage wraps around the regular Amount handler to provide support
// for percentage values, especially useful for tax rates.
type Percentage struct {
	Amount
}

// MakePercentage is a convenience method that will make a new
// Percentage instance with the provided value and exponent.
func MakePercentage(value int64, exp uint32) Percentage {
	return Percentage{Amount{Value: value, Exp: exp}}
}

// PercentageFromString builds a percentage value from a provided string.
// The % symbol will be removed automatically and rescale the stored amount
// to make future calculations easier. For example, the following two strings
// will be interpretted equally:
//
//  * `0.160`
//  * `16.0%`
//
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
	p.Amount, err = AmountFromString(str)
	if err != nil {
		return p, err
	}
	if rescale {
		e := p.Exp
		p.Amount = p.Amount.Rescale(e + 2).Divide(factor100)
	}

	return p, nil
}

// String outputs the percentage value in a human readable way including
// the percentage symbol.
func (p Percentage) String() string {
	e := p.Amount.Exp
	v := p.Amount.Multiply(factor100).Rescale(e - 2)
	return v.String() + "%"
}

// StringWithoutSymbol provides the raw underlying percentage value.
func (p Percentage) StringWithoutSymbol() string {
	return p.Amount.String()
}

// Of calulcates the "percent of" the provided amount. The exponent of the
// provided amount is used.
func (p Percentage) Of(a Amount) Amount {
	return a.Multiply(p.Amount)
}

// MarshalText provides the byte value of the amount. See also the
// String() method.
func (p Percentage) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

// MarshalJSON provides the text value of percentage wrapped in
// quotes ready to be included in a JSON object.
func (p Percentage) MarshalJSON() ([]byte, error) {
	str := "\"" + p.String() + "\""
	return []byte(str), nil
}

// UnmarshalText will decode the percentage value, even if it is quoted
// as a string.
func (p *Percentage) UnmarshalText(value []byte) error {
	if string(value) == "null" {
		return nil
	}

	str, err := unquote(value)
	if err != nil {
		return fmt.Errorf("decoding string `%s`: %w", value, err)
	}

	result, err := PercentageFromString(string(str))
	if err != nil {
		return fmt.Errorf("decoding string `%s`: %w", str, err)
	}
	*p = result

	return nil
}
