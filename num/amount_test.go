package num_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAmountAdd(t *testing.T) {
	// Use table driven tests to test multiple scenarios
	tests := []struct {
		a, b, e num.Amount
	}{
		{num.MakeAmount(200, 2), num.MakeAmount(1000, 3), num.MakeAmount(300, 2)},
		{num.MakeAmount(2000, 2), num.MakeAmount(100, 2), num.MakeAmount(2100, 2)},
		{num.MakeAmount(299, 3), num.MakeAmount(1000, 2), num.MakeAmount(10299, 3)},
		{num.MakeAmount(2000, 2), num.MakeAmount(-1000, 2), num.MakeAmount(1000, 2)},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v + %v = %v", test.a.String(), test.b.String(), test.e.String()), func(t *testing.T) {
			r := test.a.Add(test.b)
			assert.True(t, r.Equals(test.e))
			assert.Equal(t, test.e.String(), r.String())
		})
	}
}

func TestAmountSubtract(t *testing.T) {
	// Use table driven tests to test multiple scenarios
	tests := []struct {
		a, b, e num.Amount
	}{
		{num.MakeAmount(200, 2), num.MakeAmount(1000, 3), num.MakeAmount(100, 2)},
		{num.MakeAmount(200, 2), num.MakeAmount(1000, 2), num.MakeAmount(-800, 2)},
		{num.MakeAmount(299, 3), num.MakeAmount(1000, 2), num.MakeAmount(-9701, 3)},
		{num.MakeAmount(2000, 2), num.MakeAmount(-1000, 2), num.MakeAmount(3000, 2)},
		{num.MakeAmount(1890000, 2), num.MakeAmount(1890002, 2), num.MakeAmount(-2, 2)},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v - %v = %v", test.a.String(), test.b.String(), test.e.String()), func(t *testing.T) {
			r := test.a.Subtract(test.b)
			assert.True(t, r.Equals(test.e))
			assert.Equal(t, test.e.String(), r.String())
		})
	}
}

func TestAmountCompare(t *testing.T) {
	a := num.MakeAmount(1000, 2)
	b := num.MakeAmount(2000, 2)
	assert.Equal(t, 0, a.Compare(a))
	assert.Equal(t, -1, a.Compare(b))
	assert.Equal(t, 1, b.Compare(a))
}

func TestAmountNewFromString(t *testing.T) {
	a, err := num.AmountFromString("245.890")
	require.NoError(t, err)

	e := num.MakeAmount(245890, 3)
	assert.True(t, a.Equals(e))

	a, err = num.AmountFromString("245")
	assert.NoError(t, err)
	e = num.MakeAmount(245, 0)
	assert.True(t, a.Equals(e))

	a, err = num.AmountFromString("-245.00")
	assert.NoError(t, err)
	assert.EqualValues(t, a.Value(), -24500)

	a, err = num.AmountFromString("-245.12")
	assert.NoError(t, err)
	assert.EqualValues(t, a.Value(), -24512)

	a, err = num.AmountFromString("0.022")
	assert.NoError(t, err)
	assert.EqualValues(t, a.Value(), 22)

	a, err = num.AmountFromString("-0.022")
	assert.NoError(t, err)
	assert.EqualValues(t, a.Value(), -22)

	_, err = num.AmountFromString("23.433.00")
	assert.Error(t, err)
	_, err = num.AmountFromString("23,433.00")
	assert.Error(t, err)
	_, err = num.AmountFromString("1234.bar")
	assert.Error(t, err)
}

func TestMultiply(t *testing.T) {
	a := num.MakeAmount(10010, 2)
	x := num.MakeAmount(21, 1)
	a.Multiply(x)
	assert.Equal(t, "100.10", a.String(), "should not modify original amount")

	tests := []struct {
		a, b, e num.Amount
	}{
		{num.MakeAmount(10010, 2), num.MakeAmount(21, 1), num.MakeAmount(21021, 2)},
		{num.MakeAmount(200, 0), num.MakeAmount(21, 2), num.MakeAmount(42, 0)},
		{num.MakeAmount(1002002, 4), num.MakeAmount(150, 2), num.MakeAmount(1503003, 4)},
		{num.MakeAmount(669099, 2), num.MakeAmount(23, 2), num.MakeAmount(153893, 2)},
		{num.MakeAmount(101, 2), num.MakeAmount(101, 2), num.MakeAmount(102, 2)},
		{num.MakeAmount(133, 2), num.MakeAmount(133, 2), num.MakeAmount(177, 2)},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v x %v = %v", test.a.String(), test.b.String(), test.e.String()), func(t *testing.T) {
			r := test.a.Multiply(test.b)
			assert.True(t, r.Equals(test.e))
			assert.Equal(t, test.e.String(), r.String())
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		a, b, e num.Amount
	}{
		{num.MakeAmount(10010, 2), num.MakeAmount(22, 1), num.MakeAmount(4550, 2)},
		{num.MakeAmount(20000, 2), num.MakeAmount(21, 2), num.MakeAmount(95238, 2)},
		{num.MakeAmount(200, 0), num.MakeAmount(21, 2), num.MakeAmount(952, 0)},
		{num.MakeAmount(1000, 2), num.MakeAmount(11, 0), num.MakeAmount(91, 2)},
		{num.MakeAmount(1000, 0), num.MakeAmount(16, 0), num.MakeAmount(63, 0)},
		{num.MakeAmount(1000, 0), num.MakeAmount(14, 0), num.MakeAmount(71, 0)},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v / %v = %v", test.a.String(), test.b.String(), test.e.String()), func(t *testing.T) {
			r := test.a.Divide(test.b)
			assert.True(t, r.Equals(test.e))
			assert.Equal(t, test.e.String(), r.String())
		})
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
	a = num.MakeAmount(-2, 2)
	assert.Equal(t, "-0.02", a.String())
}

func TestAmountMinimalString(t *testing.T) {
	a := num.MakeAmount(123000, 3)
	assert.Equal(t, "123", a.MinimalString())
	a = num.MakeAmount(123000, 5)
	assert.Equal(t, "1.23", a.MinimalString())
	a = num.MakeAmount(123000, 0)
	assert.Equal(t, "123000", a.MinimalString())
}

func TestAmountFloat64(t *testing.T) {
	a := num.MakeAmount(123123, 3)
	assert.Equal(t, 123.123, a.Float64())
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
	if r.String() != "1235" {
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
	a = num.MakeAmount(12345678, 4)
	r = a.Rescale(2)
	assert.Equal(t, "1234.57", r.String(), "rounded number")
}

func TestAmountRemove(t *testing.T) {
	a := num.MakeAmount(20000, 2)
	p := num.MakePercentage(10, 2)
	b := a.Remove(p)
	assert.Equal(t, "181.82", b.String())
}

func TestAmountUpscale(t *testing.T) {
	a := num.MakeAmount(2123, 2)
	b := a.Upscale(2)
	assert.Equal(t, "21.2300", b.String())
}

func TestAmountDownscale(t *testing.T) {
	a := num.MakeAmount(2183, 2)
	b := a.Downscale(2)
	assert.Equal(t, "22", b.String())
	b = a.Downscale(5)
	assert.Equal(t, "22", b.String())
}

func TestAmountRescaleUp(t *testing.T) {
	a := num.MakeAmount(123456, 2)
	r := a.RescaleUp(4)
	assert.Equal(t, uint32(4), r.Exp(), "expected precision match")
	r = a.RescaleUp(2)
	assert.Equal(t, uint32(2), r.Exp(), "expected no precision change")
}

func TestAmountMatchPrecision(t *testing.T) {
	a := num.MakeAmount(123456, 2)
	a2 := num.MakeAmount(12345678, 4)
	a3 := num.MakeAmount(1234, 0)
	a4 := num.MakeAmount(456789, 2)
	r := a.MatchPrecision(a2)
	assert.Equal(t, a2.Exp(), r.Exp(), "expected precision match")
	r = a.MatchPrecision(a3)
	assert.Equal(t, a.Exp(), r.Exp(), "expected no precision change")
	r = a.MatchPrecision(a4)
	assert.Equal(t, a.Exp(), r.Exp(), "expected no precision change")
}

func TestAmountInvert(t *testing.T) {
	a := num.MakeAmount(1234, 2)
	a = a.Invert()
	assert.Equal(t, "-12.34", a.String())
	a = num.MakeAmount(-1234, 2)
	a = a.Invert()
	assert.Equal(t, "12.34", a.String())
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

func TestNegativeAmountUnmarshalJSON(t *testing.T) {
	d := []byte(`{"amount":"-12.43"}`)
	o := struct {
		Amount num.Amount
	}{}
	if err := json.Unmarshal(d, &o); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if o.Amount.Compare(num.MakeAmount(-1243, 2)) != 0 {
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
	if !bytes.Equal(d, de) {
		t.Errorf("results don't match, got: %s", d)
	}
}

func TestNegativeAmountMarshalJSON(t *testing.T) {
	o := struct {
		Amount num.Amount `json:"amount"`
	}{
		Amount: num.MakeAmount(-1267, 2),
	}
	d, err := json.Marshal(o)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	de := []byte(`{"amount":"-12.67"}`)
	if !bytes.Equal(d, de) {
		t.Errorf("results don't match, got: %s", d)
	}
}
