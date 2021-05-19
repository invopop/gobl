package c14n_test

import (
	"testing"

	"github.com/invopop/gobl/c14n"
)

func TestStringMarshalJSON(t *testing.T) {
	s := c14n.String(`This is "a" test with quotes`)
	d, err := s.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err.Error())
	}
	if string(d) != `"This is \"a\" test with quotes"` {
		t.Errorf("unexpected output, got: %v", string(d))
	}

}

func TestFloatMarshalJSON(t *testing.T) {
	f := c14n.Float(0.0)

	d, err := f.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if string(d) != "0.0E0" {
		t.Errorf("got unexpected result: %v", string(d))
	}

	f = c14n.Float(1.0)
	d, err = f.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if string(d) != "1.0E0" {
		t.Errorf("got unexpected result: %v", string(d))
	}

	f = c14n.Float(123.5)
	d, err = f.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if string(d) != "1.235E2" {
		t.Errorf("got unexpected result: %v", string(d))
	}

	f = c14n.Float(123456789123456.0)
	d, err = f.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if string(d) != "1.23456789123456E14" {
		t.Errorf("got unexpected result: %v", string(d))
	}

	f = c14n.Float(0.000001234567891234560)
	d, err = f.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if string(d) != "1.23456789123456E-6" {
		t.Errorf("got unexpected result: %v", string(d))
	}

	f = c14n.Float(0.0000000000000000001234567891234560)
	d, err = f.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if string(d) != "1.23456789123456E-19" {
		t.Errorf("got unexpected result: %v", string(d))
	}

	f = c14n.Float(1.234567891234560000e-110)
	d, err = f.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if string(d) != "1.23456789123456E-110" {
		t.Errorf("got unexpected result: %v", string(d))
	}
}
