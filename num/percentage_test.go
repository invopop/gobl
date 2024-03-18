package num_test

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/stretchr/testify/assert"
)

func TestPercentage(t *testing.T) {
	p := num.MakePercentage(1600, 4)
	if p.Value() != 1600 {
		t.Errorf("unexpected value, got: %v", p.Value())
	}
	p2 := num.MakePercentage(1600, 4)
	if p.Value() != 1600 {
		t.Errorf("unexpected value, got: %v", p.Value())
	}
	if !p.Equals(p2) {
		t.Errorf("expected percentages to be the same")
	}
}

func TestPercentageFromString(t *testing.T) {
	p, err := num.PercentageFromString("16.0%")
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if p.Value() != 160 {
		t.Errorf("unexpected value from string, got: %v", p.Value())
	}
	if p.Exp() != 3 {
		t.Errorf("unexpected exponential value, got: %v", p.Exp())
	}
}

func TestPercentageString(t *testing.T) {
	p := num.MakePercentage(1600, 4)
	if p.String() != "16.00%" {
		t.Errorf("unexpected string result from percentage, got: %v", p.String())
	}
	if p.StringWithoutSymbol() != "16.00" {
		t.Errorf("unexpected raw percentage, got: %v", p.StringWithoutSymbol())
	}
	p2 := num.MakePercentage(2000, 2000) // silly number
	if p2.String() != "NA%" {
		t.Errorf("expect invalid exponential to return bad number")
	}
	p3 := num.MakePercentage(200, 1)
	if p3.String() != "2000%" {
		t.Errorf("unexpected percentage string, got: %v", p3.String())
	}
}

func TestPercentageOf(t *testing.T) {
	p := num.MakePercentage(170, 3)
	a := num.MakeAmount(10000, 2)
	r := p.Of(a)
	assert.Equal(t, "17.00", r.String())
}

func TestFactor(t *testing.T) {
	p := num.MakePercentage(160, 3)
	f := p.Factor()
	s := f.String()
	if s != "1.160" {
		t.Errorf("unexpected factor result, got: %v", s)
	}
}

func TestPercentageFrom(t *testing.T) {
	p := num.MakePercentage(160, 3)
	a := num.MakeAmount(11600, 2)
	x := p.From(a)
	assert.Equal(t, "16.00", x.String())
}

func TestPercentageRescale(t *testing.T) {
	p := num.MakePercentage(160, 3)
	x := p.Rescale(4)
	if x.String() != "16.00%" {
		t.Errorf("unexpected percentage from result: %v", x.String())
	}

	p = num.MakePercentage(20, 3)
	x = p.Rescale(2)
	if x.String() != "2%" {
		t.Errorf("unexpected percentage from result: %v", x.String())
	}
}

func TestPercentageInvert(t *testing.T) {
	p := num.MakePercentage(160, 3)
	x := p.Invert()
	assert.Equal(t, "-16.0%", x.String())
}
