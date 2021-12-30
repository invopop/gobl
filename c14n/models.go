package c14n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"unicode/utf8"
)

// Canonicalable defines what we expect from objects that need to be converted
// into our standardized JSON. All structures that comply with this interface
// are expected to contain data that was already sourced from a JSON document,
// so we don't need to worry too much about the conversion process.
type Canonicalable interface {
	MarshalJSON() ([]byte, error)
}

// Object contains a simple list of items, which are in essence key-value pairs.
// The item array means that attributes can be ordered by their key.
type Object struct {
	Attributes []*Attribute
}

// Array contains a list of canonicable values, as opposed to the objects key-value
// pairs.
type Array struct {
	Values []Canonicalable
}

// String is our representation of a regular string, prepared for JSON marshalling.
type String string

// Integer numbers have no decimal places and are limited to 64 bits.
type Integer int64

// Float numbers must be represented by a 64-bit signed integer and exponential
// that reflects the position of the decimal place. We're not going to support
// numbers whose signifcant digits do not fit inside an int64, for big numbers,
// use an alternative method of serialization such as Base64.
type Float float64

// Bool handles binary true or false.
type Bool bool

// Null wraps around a null value
type Null struct{}

// Attribute represents a key-value pair used in objects. Using an array guarantees ordering
// of keys, which is one of the fundamental requirements for canonicalization.
type Attribute struct {
	Key   string
	Value Canonicalable
}

// hex defines the list of hex characters used to escape unicodes that are
// not considered safe.
var hex = "0123456789ABCDEF"

// Sort ensures all the object's attributes are ordered according to the key.
func (o *Object) Sort() {
	sort.SliceStable(o.Attributes, func(i, j int) bool {
		return o.Attributes[i].Key < o.Attributes[j].Key
	})
}

// MarshalJSON combines all the objects elements into an ordered
// key-value list of marshalled attributes.
func (o *Object) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, v := range o.Attributes {
		a, err := v.MarshalJSON()
		if err != nil {
			return nil, err
		}
		if len(a) == 0 { // as per spec, skip empty attributes
			continue
		}
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.Write(a)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// MarshalJSON recursively marshals all of the arrays items
// and joins the results together to form a JSON byte array.
func (a *Array) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i, v := range a.Values {
		if i > 0 {
			buf.WriteByte(',')
		}
		data, err := v.MarshalJSON()
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

// MarshalJSON provides byte array of the string inside
// quotes.
func (o String) MarshalJSON() ([]byte, error) {
	return encodeString(string(o))
}

// MarshalJSON provides string representation of integer.
func (i Integer) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i), 10)), nil
}

// MarshalJSON for floats uses the strconv library with some output
// hacks to ensure we're in a valid JSON format. This is also actually what the
// fmt package and Printf related methods do to get around all the complexities
// of float conversion.
func (f Float) MarshalJSON() ([]byte, error) {
	num := []byte{}
	num = strconv.AppendFloat(num, float64(f), 'E', -1, 64)

	// When decimal place is missing, add it. This only happens
	// when the number is 0.
	if num[1] != '.' {
		num = append(num[0:3], num[1:]...)
		num[1] = '.'
		num[2] = '0'
	}

	// Split into two parts
	i := bytes.IndexByte(num, 'E')
	exp := make([]byte, len(num)-i-1)
	copy(exp, num[i+1:])
	num = num[:i+1] // shorten

	// Remove + in exponent
	if exp[0] == '+' {
		exp = exp[1:]
	}

	// Remove excess exponential 0s
	j := 0
	k := 0
	for i, v := range exp {
		if v == '-' || v == '+' {
			j = 1
		} else if v == '0' && (i+1) < len(exp) {
			k = i + 1
		} else {
			break // first non-zero
		}
	}
	if k != 0 {
		exp = append(exp[:j], exp[k:]...)
	}

	num = append(num, exp...)
	return num, nil
}

// MarshalJSON provides the null string.
func (n Null) MarshalJSON() ([]byte, error) {
	return []byte(`null`), nil
}

// MarshalJSON provides the JSON standard true or false response.
func (b Bool) MarshalJSON() ([]byte, error) {
	if b {
		return []byte(`true`), nil
	}
	return []byte(`false`), nil
}

// MarshalJSON creates a key-value pair in JSON format. A null value
// in an attribute will return an empty byte array.
func (a *Attribute) MarshalJSON() ([]byte, error) {
	if _, ok := a.Value.(Null); ok {
		return nil, nil
	}
	key, err := encodeString(a.Key)
	if err != nil {
		return nil, err
	}
	val, err := a.Value.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	buf.Write(key)
	buf.WriteByte(':')
	buf.Write(val)
	return buf.Bytes(), nil
}

// encodeString is inspired by the golang encoding/json package,
// with a few modifications for canonicalisation.
func encodeString(s string) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if safeSet[b] {
				i++
				continue
			}
			if start < i {
				buf.WriteString(s[start:i])
			}
			buf.WriteByte('\\')
			// Always prefer short escape codes
			switch b {
			case '\\', '"':
				buf.WriteByte(b)
			case '\n': // Line Feed
				buf.WriteByte('n')
			case '\r': // Carriage Return
				buf.WriteByte('r')
			case '\t': // Tab
				buf.WriteByte('t')
			case '\f': // Form Feed
				buf.WriteByte('f')
			case '\b': // Backspace
				buf.WriteByte('b')
			default:
				// This encodes bytes < 0x20 unless already handled.
				buf.WriteString(`u00`)
				buf.WriteByte(hex[b>>4])
				buf.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError {
			// don't accept anything that isn't valid UTF-8, no exceptions.
			return nil, &json.UnsupportedValueError{Value: reflect.ValueOf(s), Str: fmt.Sprintf("%q", s)}
		}
		i += size
	}
	if start < len(s) {
		buf.WriteString(s[start:])
	}
	buf.WriteByte('"')
	return buf.Bytes(), nil
}
