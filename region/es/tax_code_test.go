package es

import "testing"

func TestCleanTaxCode(t *testing.T) {
	tests := []struct {
		Code     string
		Expected string
	}{
		{
			Code:     "93471790-C",
			Expected: "93471790C",
		},
		{
			Code:     " 4359 6386 R ",
			Expected: "43596386R",
		},
		{
			Code:     "Z-8327649-K",
			Expected: "Z8327649K",
		},
	}
	for i, ts := range tests {
		if err := CleanTaxCode(ts.Code); err != ts.Expected {
			t.Errorf("unexpected result: %d: got: %+v", i, err)
		}
	}
}

func TestVerifyNationalCode(t *testing.T) {
	tests := []struct {
		Code     string
		Expected interface{}
		Message  string
	}{
		{
			Code:     "93471790C",
			Expected: nil,
		},
		{
			Code:     "43596386R",
			Expected: nil,
		},
		{
			Code:     "00000010X",
			Expected: nil,
		},
		{
			Code:     "93471790A",
			Expected: ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "00000000A",
			Expected: ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "0111111C",
			Expected: ErrTaxCodeNoMatch,
		},
	}
	for i, ts := range tests {
		if err := verifyNationalCode(ts.Code); err != ts.Expected {
			t.Errorf("unexpected result: %d: got: %+v", i, err)
		}
	}
}

func TestVerifyForeignCode(t *testing.T) {
	tests := []struct {
		Code     string
		Expected interface{}
		Message  string
	}{
		{
			Code:     "X5102754C",
			Expected: nil,
		},
		{
			Code:     "Z8327649K",
			Expected: nil,
		},
		{
			Code:     "Y4174455S",
			Expected: nil,
		},
		{
			Code:     "X5102755C",
			Expected: ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "X111111C",
			Expected: ErrTaxCodeNoMatch,
		},
	}
	for i, ts := range tests {
		if err := verifyForeignCode(ts.Code); err != ts.Expected {
			t.Errorf("unexpected result: %d: got: %+v", i, err)
		}
	}
}

func TestVerifyOrgCode(t *testing.T) {
	tests := []struct {
		Code     string
		Expected interface{}
		Message  string
	}{
		{
			Code:     "A58818501",
			Expected: nil,
		},
		{
			Code:     "B65410011",
			Expected: nil,
		},
		{
			Code:     "V7565938C",
			Expected: nil,
		},
		{
			Code:     "V75659383",
			Expected: nil,
		},
		{
			Code:     "F0605378I",
			Expected: nil,
		},
		{
			Code:     "Q2238877A",
			Expected: nil,
		},
		{
			Code:     "D40022956",
			Expected: nil,
		},
		{
			Code:     "A5881850B",
			Expected: ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "B65410010",
			Expected: ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "V75659382",
			Expected: ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "V7565938B",
			Expected: ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "F06053787",
			Expected: ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "Q22388770",
			Expected: ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "D4002295J",
			Expected: ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "00000000A",
			Expected: ErrTaxCodeNoMatch,
		},
		{
			Code:     "B0111111",
			Expected: ErrTaxCodeNoMatch,
		},
	}
	for i, ts := range tests {
		if err := verifyOrgCode(ts.Code); err != ts.Expected {
			t.Errorf("unexpected result: %d: got: %+v", i, err)
		}
	}
}

func TestVerifyOtherCode(t *testing.T) {
	tests := []struct {
		Code     string
		Expected interface{}
	}{
		{
			Code:     "K9514336H",
			Expected: nil,
		},
		{
			Code:     "K95143363",
			Expected: ErrTaxCodeInvalidCheck,
		},
		{
			Code:     "X111111C",
			Expected: ErrTaxCodeNoMatch,
		},
	}
	for i, ts := range tests {
		if err := verifyOtherCode(ts.Code); err != ts.Expected {
			t.Errorf("unexpected result: %d: got: %+v", i, err)
		}
	}
}
