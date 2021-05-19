package c14n

import (
	"encoding/json"
	"errors"
	"io"
)

// UnmarshalJSON expects an io Reader whose data will be parsed using a streaming
// JSON decoder and converted into a "Canonicalable" set of structures. The resulting
// objects can then be re-encoded back into canonical JSON suitable for sending to
// a hashing algorithm.
func UnmarshalJSON(src io.Reader) (Canonicalable, error) {
	dec := json.NewDecoder(src)
	dec.UseNumber()

	res, err := handleNextToken(dec)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CanonicalJSON performs the unmarshal and marshal commands in one go.
func CanonicalJSON(src io.Reader) ([]byte, error) {
	obj, err := UnmarshalJSON(src)
	if err != nil {
		return nil, err
	}
	return obj.MarshalJSON()
}

func handleNextToken(dec *json.Decoder) (Canonicalable, error) {
	t, err := dec.Token()
	if err == io.EOF {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if d, ok := t.(json.Delim); ok {
		switch d {
		case '{':
			return handleObject(dec)
		case '[':
			return handleArray(dec)
		default:
			return nil, nil
		}
	} else {
		// We're dealing with an item, just return it
		return tokenToValue(t)
	}

}

func handleObject(dec *json.Decoder) (*Object, error) {
	obj := new(Object)
	obj.Attributes = make([]*Attribute, 0)
	for {
		a, err := handleAttribute(dec)
		if err != nil {
			return nil, err
		}
		if a == nil { // i.e., no more left
			obj.Sort()
			return obj, nil
		}
		obj.Attributes = append(obj.Attributes, a)
	}
}

func handleArray(dec *json.Decoder) (*Array, error) {
	ary := new(Array)
	ary.Values = make([]Canonicalable, 0)
	for {
		val, err := handleNextToken(dec)
		if err != nil {
			return ary, err
		}
		if val == nil {
			return ary, nil
		}
		ary.Values = append(ary.Values, val)
	}
}

func handleAttribute(dec *json.Decoder) (*Attribute, error) {
	// handle item attempts to get the next two tokens as we
	// know these must form a key-pair.
	key, err := handleNextToken(dec)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, nil
	}
	k, ok := key.(String)
	if !ok {
		return nil, errors.New("item key must be a string")
	}
	a := new(Attribute)
	a.Key = string(k)
	a.Value, err = handleNextToken(dec)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// tokenToValue does all the heavy lifting of ensuring we have a
// usable value.
func tokenToValue(t json.Token) (Canonicalable, error) {
	if s, ok := t.(string); ok {
		return String(s), nil
	}
	if n, ok := t.(json.Number); ok {
		if i, err := n.Int64(); err == nil {
			return Integer(i), nil
		}
		if f, err := n.Float64(); err == nil {
			return Float(f), nil
		}
	}
	if b, ok := t.(bool); ok {
		return Bool(b), nil
	}
	return Null{}, nil
}
