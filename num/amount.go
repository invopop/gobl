package num

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/invopop/jsonschema"
)

// Amount represents a quantity with decimal places that will not suffer
// rounding errors like traditional floats.
// Use cases are assumed to be within the "human manageable domain", i.e.
// for dealing with counts, money, rates, short distances, etc.
// Implementation is inspired by https://github.com/shopspring/decimal, but
// simplified to account for the expectations of GOBL.
type Amount struct {
	value int64
	exp   uint32
}

// NewAmount provides a pointer to an Amount instance. Normally we'd recommend
// using the `MakeAmount` method.
func NewAmount(val int64, exp uint32) *Amount {
	a := MakeAmount(val, exp)
	return &a
}

// MakeAmount is a helper to make it a little easier to build a new Amount
// instance. We use "Make" instead of "New" as there are no pointers.
func MakeAmount(val int64, exp uint32) Amount {
	return Amount{value: val, exp: exp}
}

// AmountFromString takes the provided string and tries to convert it
// into an amount object. Strings must be in a simplified format with no
// commas and a single `.` to separate the decimal places. Numbers are
// expected to have a fixed number of decimal places, so if your dealing
// with a string like `"12.000"`, the accuracy will be assumed to be 3
// decimal places.
//
// If you're dealing with numbers from humans which may contain symbols,
// commas, european style fullstops, underscores, etc. then you should use
// the `AmountFromHumanString` method.
func AmountFromString(val string) (Amount, error) {
	a := Amount{}
	n := strings.HasPrefix(val, "-")
	x := strings.Split(strings.TrimPrefix(val, "-"), ".")
	l := len(x)
	if l > 2 {
		return a, fmt.Errorf("amount must contain 0 or 1 decimal separators: %v", val)
	}

	// Parse the "major" part
	v, err := strconv.ParseInt(x[0], 10, 64)
	if err != nil {
		return a, fmt.Errorf("invalid major number '%v', %w", val, err)
	}
	e := uint32(0)
	v2 := int64(0)

	// Parse the decimal places (if present)
	if l == 2 {
		v2, err = strconv.ParseInt(x[1], 10, 64)
		if err != nil {
			return a, fmt.Errorf("invalid decimal number '%v', %w", val, err)
		}
		e = uint32(len(x[1]))
		v = v * intPow(10, e)
		v += v2
	}

	// Prepare the result
	if n {
		a.value = -v
	} else {
		a.value = v
	}
	a.exp = e
	return a, nil
}

// AmountFromHumanString removes an excess decimal places, commas, or
// other symbols so that we end up with a simple string that can be parsed.
func AmountFromHumanString(val string) (Amount, error) {
	return Amount{}, errors.New("not yet implemented")
}

// IsZero returns true if the value of the amount is 0.
func (a Amount) IsZero() bool {
	return a.value == 0
}

// Add will add the two amounts together using the base's exponential
// value for the resulting new amount.
func (a Amount) Add(a2 Amount) Amount {
	a2 = a2.Rescale(a.exp)
	return Amount{a.value + a2.value, a.exp}
}

// Subtract takes away the amount provided from the base.
func (a Amount) Subtract(a2 Amount) Amount {
	a2 = a2.Rescale(a.exp)
	return Amount{value: a.value - a2.value, exp: a.exp}
}

// Multiply our base amount by the provided amount.
func (a Amount) Multiply(a2 Amount) Amount {
	return Amount{
		value: (a.value * a2.value) / intPow(10, a2.exp),
		exp:   a.exp,
	}
}

// Divide our base amount by the provided amount. We use floating point to do the actual division
// and then round again to get an int. This prevents rounding errors, but if you want true division
// with a base and a remainder, use the Split method.
func (a Amount) Divide(a2 Amount) Amount {
	v := float64(a.value*intPow(10, a2.exp)) / float64(a2.value)
	return Amount{
		value: int64(math.Round(v)),
		exp:   a.exp,
	}
}

// Split divides the amount by x, like Divide, but also provides an
// additional amount with a remainder so that we avoid rounding
// errors.
func (a Amount) Split(x int) (Amount, Amount) {
	a2 := a.Divide(MakeAmount(int64(x), 0))
	a3 := a2.Multiply(MakeAmount(int64(x-1), 0))
	a3 = a.Subtract(a3)
	return a2, a3
}

// Compare two amounts and return an integer value according to the
// sign of the difference:
//
//	-1 if a <  a2
//	 0 if a == a2
//	 1 if a >  a2
func (a Amount) Compare(a2 Amount) int {
	a, a2 = rescaleAmountPair(a, a2)
	if a.value < a2.value {
		return -1
	}
	if a.value > a2.value {
		return 1
	}
	return 0
}

// Equals returns true if the two amounts represent the same value,
// regardless of the exponential.
func (a Amount) Equals(a2 Amount) bool {
	return a.Compare(a2) == 0
}

// Rescale will multiply or divide the amount's value to match the
// provided exponential. This method will round values in the case of
// reducing the exponent.
func (a Amount) Rescale(exp uint32) Amount {
	if a.exp > exp {
		// need to divide
		e := a.exp - exp
		v := float64(a.value) / float64(intPow(10, e))
		return Amount{int64(math.Round(v)), exp}
	}
	if a.exp < exp {
		// need to multiply
		e := exp - a.exp
		v := a.value * intPow(10, e)
		return Amount{v, exp}
	}
	return a
}

// MatchPrecision will rescale the exponent value of the amount so that it
// matches the scale of the provided amount, but *only* if it is higher.
func (a Amount) MatchPrecision(a2 Amount) Amount {
	if a2.exp > a.exp {
		return a.Rescale(a2.exp)
	}
	return a
}

// Upscale increases the accuracy of the amount by rescaling the exponent
// by the provided amount.
func (a Amount) Upscale(accuracy uint32) Amount {
	return a.Rescale(a.Exp() + accuracy)
}

// Downscale decreases the amount's exponent by the provided accuracy.
func (a Amount) Downscale(accuracy uint32) Amount {
	var x uint32
	if accuracy > a.Exp() {
		x = 0
	} else {
		x = a.Exp() - accuracy
	}
	return a.Rescale(x)
}

// Remove takes the provided percentage away from the amount assuming it was
// already applied previously.
func (a Amount) Remove(percent Percentage) Amount {
	p := percent.From(a)
	return a.Subtract(p)
}

// Invert the value.
func (a Amount) Invert() Amount {
	return Amount{value: -a.value, exp: a.exp}
}

// Value provides the amount's value
func (a Amount) Value() int64 {
	return a.value
}

// Exp provides the amount's exponent value.
func (a Amount) Exp() uint32 {
	return a.exp
}

// String returns the simplified string amount.
func (a Amount) String() string {
	if a.exp == 0 {
		return fmt.Sprintf("%d", a.value)
	}
	if a.exp > 1000 {
		return "NA"
	}
	p := intPow(10, a.exp)
	v1 := a.value / p
	v2 := a.value - (v1 * p)
	if v2 < 0 {
		v2 = -v2
	}
	return fmt.Sprintf("%d.%0*d", v1, a.exp, v2)
}

// MinimalString provides the amount without any tailing 0s or '.'
// if one is left over.
func (a Amount) MinimalString() string {
	s := a.String()
	if !strings.Contains(s, ".") {
		return s
	}
	s = strings.TrimRight(s, "0")
	return strings.TrimSuffix(s, ".")
}

// Float64 provides the amount as a float64 value which should be used
// with caution!
func (a Amount) Float64() float64 {
	return float64(a.value) / float64(intPow(10, a.exp))
}

// MarshalText provides the byte value of the amount. See also the
// String() method.
// We always add quotes around values as number representations do not
// guarantee that tailing 0s will be maintained. It's important
// to remember that amounts are typically for humans, and thus
// it makes sense to consider them as strings.
func (a Amount) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

// MarshalJSON takes string value of the text and adds quotes around
// it ready to be used in a JSON object.
func (a Amount) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteByte('"')
	buf.WriteString(a.String())
	buf.WriteByte('"')
	return buf.Bytes(), nil
}

// UnmarshalText will decode the amount value, even if it is quoted
// as a string and will be used for JSON, XML, or any other text
// unmarshaling.
func (a *Amount) UnmarshalText(value []byte) error {
	if string(value) == "null" {
		return nil
	}

	str := unquote(value)
	amount, err := AmountFromString(string(str))
	if err != nil {
		return fmt.Errorf("decoding string `%s`: %w", str, err)
	}
	*a = amount

	return nil
}

// UnmarshalJSON ensures amounts will be parsed even if defined as
// numbers in the source JSON.
func (a *Amount) UnmarshalJSON(value []byte) error {
	return a.UnmarshalText(value)
}

func unquote(value []byte) []byte {
	// If the amount is quoted, strip the quotes
	if len(value) > 2 && value[0] == '"' && value[len(value)-1] == '"' {
		value = value[1 : len(value)-1]
	}
	return value
}

func rescaleAmountPair(a, a2 Amount) (Amount, Amount) {
	// Take the largest exp
	exp := a.exp
	if a2.exp > exp {
		exp = a2.exp
	}
	return a.Rescale(exp), a2.Rescale(exp)
}

func intPow(base int, exp uint32) int64 { // nolint:unparam
	out := int64(1)
	for exp != 0 {
		out *= int64(base)
		exp--
	}
	return out
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (Amount) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Pattern:     `^\-?[0-9]+(\.[0-9]+)?$`,
		Title:       "Amount",
		Description: "Quantity with optional decimal places that determine accuracy.",
	}
}
