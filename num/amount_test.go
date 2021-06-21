package num_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/num"
)

func TestAmountAdd(t *testing.T) {
	a := num.MakeAmount(200, 2)
	a2 := num.MakeAmount(1000, 3)
	r := a.Add(a2)
	e := num.MakeAmount(300, 2)
	if !r.Equals(e) {
		t.Errorf("did not add amounts correctly, got: %v", r)
	}
}

func TestAmountSubtract(t *testing.T) {
	a := num.MakeAmount(200, 2)
	a2 := num.MakeAmount(1000, 3)
	r := a.Subtract(a2)
	e := num.MakeAmount(100, 2)
	if !r.Equals(e) {
		t.Errorf("did not add amounts correctly, got: %v", r)
	}
	a = num.MakeAmount(200, 2)
	a2 = num.MakeAmount(1000, 2)
	r = a.Subtract(a2)
	e = num.MakeAmount(-800, 2)
	if !r.Equals(e) {
		t.Errorf("did not add amounts correctly, got: %v", r)
	}
}

func TestAmountCompare(t *testing.T) {
	a := num.MakeAmount(1000, 2)
	b := num.MakeAmount(2000, 2)
	if a.Compare(a) != 0 {
		t.Errorf("expected 0")
	}
	if a.Compare(b) != -1 {
		t.Errorf("expected -1")
	}
	if b.Compare(a) != 1 {
		t.Errorf("expected 1")
	}
}

func TestAmountNewFromString(t *testing.T) {
	a, err := num.AmountFromString("245.890")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	e := num.MakeAmount(245890, 3)
	if !a.Equals(e) {
		t.Errorf("unexpected parsed value, got: %v", a)
	}
	a, err = num.AmountFromString("245")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	e = num.MakeAmount(245, 0)
	if !a.Equals(e) {
		t.Errorf("unexpected parsed value, got: %v", a)
	}
	a, err = num.AmountFromString("23.433.00")
	if err == nil {
		t.Errorf("expected error, got: %v", a)
	}
	a, err = num.AmountFromString("23,433.00")
	if err == nil {
		t.Errorf("expected error, got: %v", a)
	}
	a, err = num.AmountFromString("1234.bar")
	if err == nil {
		t.Errorf("expected error, got: %v", a)
	}
}

func TestMultiply(t *testing.T) {
	a := num.MakeAmount(10010, 2)
	x := num.MakeAmount(21, 1)
	e := num.MakeAmount(21021, 2)
	r := a.Multiply(x)
	if !r.Equals(e) {
		t.Errorf("failed to multiply, expected: %v, got: %v", e, r)
	}
	if a.String() != "100.10" {
		t.Errorf("base was modified")
	}
	a = num.MakeAmount(200, 0)
	x = num.MakeAmount(21, 2)
	e = num.MakeAmount(42, 0)
	r = a.Multiply(x)
	if !r.Equals(e) {
		t.Errorf("failed to multiply, expected: %v, got: %v", e, r)
	}
	a = num.MakeAmount(1002002, 4)
	x = num.MakeAmount(150, 2)
	e = num.MakeAmount(1503003, 4)
	r = a.Multiply(x)
	if !r.Equals(e) {
		t.Errorf("failed to multiply, expected: %v, got: %v", e, r)
	}
}

func TestDivide(t *testing.T) {
	a := num.MakeAmount(10010, 2)
	x := num.MakeAmount(22, 1)
	e := num.MakeAmount(4550, 2)
	r := a.Divide(x)
	if !r.Equals(e) {
		t.Errorf("failed to divide, expected: %v, got: %v", e, r)
	}
	if a.String() != "100.10" {
		t.Errorf("base was modified")
	}
	a = num.MakeAmount(200, 0)
	x = num.MakeAmount(21, 2)
	e = num.MakeAmount(952, 0)
	r = a.Divide(x)
	if !r.Equals(e) {
		t.Errorf("failed to divide, expected: %v, got: %v", e, r)
	}
	a = num.MakeAmount(1000, 2)
	x = num.MakeAmount(11, 0)
	e = num.MakeAmount(91, 2)
	r = a.Divide(x)
	if !r.Equals(e) {
		t.Errorf("unexpected division result, expected: %v, got %v", e, r)
	}
}

func TestSplit(t *testing.T) {
	a := num.MakeAmount(1000, 2)
	x := 11
	r, r2 := a.Split(x)
	if r.String() != "0.91" {
		t.Errorf("unexpected result, got: %v", r.String())
	}
	if r2.String() != "0.90" {
		t.Errorf("unexpected modulus, got: %v", r2.String())
	}
}

func TestAmountString(t *testing.T) {
	a := num.MakeAmount(12345670, 3)
	if a.String() != "12345.670" {
		t.Errorf("unexpected string result, got: %v", a.String())
	}
	a = num.MakeAmount(2, 0)
	if a.String() != "2" {
		t.Errorf("unexpected string result, got: %v", a.String())
	}
	a = num.MakeAmount(50, 0)
	if a.String() != "50" {
		t.Errorf("unexpected string result, got: %v", a.String())
	}
	a = num.MakeAmount(-5025, 2)
	if a.String() != "-50.25" {
		t.Errorf("unexpected string result, got: %v", a.String())
	}
}

func TestAmountRescale(t *testing.T) {
	a := num.MakeAmount(123456, 2)
	r := a.Rescale(2)
	if r.String() != "1234.56" {
		t.Errorf("unexpected rescale result: %v", r.String())
	}
	r = a.Rescale(4)
	if r.String() != "1234.5600" {
		t.Errorf("unexpected rescale result: %v", r.String())
	}
	r = a.Rescale(0)
	if r.String() != "1234" {
		t.Errorf("unexpected rescale result: %v", r.String())
	}
	a = num.MakeAmount(21, 1)
	r = a.Rescale(2)
	if r.String() != "2.10" {
		t.Errorf("unexpected rescale result, got: %v", r.String())
	}
	a = num.MakeAmount(56, 1)
	r = a.Rescale(4)
	if r.String() != "5.6000" {
		t.Errorf("unexpected rescale result, got: %v", r.String())
	}
}

func TestAmountUnmarshalJSON(t *testing.T) {
	d := []byte(`{"amount":"12.43"}`)
	o := struct {
		Amount num.Amount
	}{}
	if err := json.Unmarshal(d, &o); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if o.Amount.Compare(num.MakeAmount(1243, 2)) != 0 {
		t.Errorf("got back unexpected response: %+v", o)
	}
}

func TestAmountMarshalJSON(t *testing.T) {
	o := struct {
		Amount num.Amount `json:"amount"`
	}{
		Amount: num.MakeAmount(1267, 2),
	}
	d, err := json.Marshal(o)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	de := []byte(`{"amount":"12.67"}`)
	if bytes.Compare(d, de) != 0 {
		t.Errorf("results don't match, got: %s", d)
	}
}
