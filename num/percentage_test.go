package num_test

import (
	"testing"

	"github.com/invopop/gobl/num"
)

func TestPercentage(t *testing.T) {
	p := num.Percentage{num.Amount{1600, 4}}
	if p.Value != 1600 {
		t.Errorf("unexpected value, got: %v", p.Value)
	}
	p2 := num.MakePercentage(1600, 4)
	if p.Value != 1600 {
		t.Errorf("unexpected value, got: %v", p.Value)
	}
	if !p.Equals(p2.Amount) {
		t.Errorf("expected percentages to be the same")
	}
}

func TestPercentageFromString(t *testing.T) {
	p, err := num.PercentageFromString("16.0%")
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if p.Value != 160 {
		t.Errorf("unexpected value from string, got: %v", p.Value)
	}
	if p.Exp != 3 {
		t.Errorf("unexpected exponential value, got: %v", p.Exp)
	}
}

func TestPercentageString(t *testing.T) {
	p := num.Percentage{num.Amount{1600, 4}}
	if p.String() != "16.00%" {
		t.Errorf("unexpected string result from percentage, got: %v", p.String())
	}
	if p.StringWithoutSymbol() != "0.1600" {
		t.Errorf("unexpected raw percentage, got: %v", p.StringWithoutSymbol())
	}
}

func TestPercentageOf(t *testing.T) {
	p := num.Percentage{num.Amount{160, 3}}
	a := num.Amount{10000, 2}
	r := p.Of(a)
	if r.String() != "16.00" {
		t.Errorf("unexpected percentage of result, got: %v", r.String())
	}
}
