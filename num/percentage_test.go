package num_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPercentage(t *testing.T) {
	p := num.MakePercentage(1600, 4)
	assert.Equal(t, int64(1600), p.Value())
	p2 := num.MakePercentage(1600, 4)
	assert.True(t, p.Equals(p2))
	pp := num.NewPercentage(1600, 4)
	assert.Equal(t, "16.00%", pp.String())
	assert.True(t, pp.Equals(p2))
	assert.True(t, p2.Equals(*pp))
}

func TestPercentageFromString(t *testing.T) {
	p, err := num.PercentageFromString("")
	require.NoError(t, err)
	assert.Equal(t, int64(0), p.Value())

	p, err = num.PercentageFromString("16.0%")
	require.NoError(t, err)
	assert.Equal(t, int64(160), p.Value())
	assert.Equal(t, uint32(3), p.Exp())

	_, err = num.PercentageFromString("bad")
	assert.ErrorContains(t, err, "invalid major number 'bad', strconv.ParseInt: parsing \"bad\": invalid syntax")

	p, err = num.PercentageFromString("0.160")
	require.NoError(t, err)
	assert.Equal(t, int64(160), p.Value())
	assert.Equal(t, uint32(3), p.Exp())
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

func TestPercentageIsZero(t *testing.T) {
	p := num.MakePercentage(0, 0)
	assert.True(t, p.IsZero())
	p = num.MakePercentage(160, 0)
	assert.False(t, p.IsZero())
	p = num.MakePercentage(-160, 0)
	assert.False(t, p.IsZero())
}

func TestPercentageIsNegative(t *testing.T) {
	p := num.MakePercentage(-160, 0)
	assert.True(t, p.IsNegative())
	p = num.MakePercentage(160, 0)
	assert.False(t, p.IsNegative())
	p = num.MakePercentage(0, 0)
	assert.False(t, p.IsNegative())
}

func TestPercentageIsPositive(t *testing.T) {
	p := num.MakePercentage(160, 0)
	assert.True(t, p.IsPositive())
	p = num.MakePercentage(-160, 0)
	assert.False(t, p.IsPositive())
	p = num.MakePercentage(0, 0)
	assert.False(t, p.IsPositive())
}

func TestPercentageUnmarshalJSONBasic(t *testing.T) {
	d := []byte(`{"percent":"16.0%"}`)
	o := struct {
		Percent num.Percentage
	}{}
	require.NoError(t, json.Unmarshal(d, &o))
	assert.Equal(t, 0, o.Percent.Compare(num.MakePercentage(160, 3)))

	d = []byte(`{"percent":0.10}`)
	require.NoError(t, json.Unmarshal(d, &o))
	assert.Equal(t, int64(10), o.Percent.Value())
	assert.Equal(t, uint32(2), o.Percent.Exp())

	o.Percent = num.MakePercentage(0, 0)
	d = []byte(`{"percent":null}`)
	require.NoError(t, json.Unmarshal(d, &o))
	assert.Equal(t, int64(0), o.Percent.Value())

	d = []byte(`{"percent":"bad"}`)
	require.ErrorContains(t, json.Unmarshal(d, &o), "invalid major number 'bad', strconv.ParseInt: parsing \"bad\": invalid syntax")
}

func TestPercentageUnmarshalJSONPointer(t *testing.T) {
	d := []byte(`{"percent":"16.0%"}`)
	o := struct {
		Percent *num.Percentage
	}{}
	require.NoError(t, json.Unmarshal(d, &o))
	assert.Equal(t, 0, o.Percent.Compare(num.MakePercentage(160, 3)))
	d = []byte(`{"percent":null}`)
	require.NoError(t, json.Unmarshal(d, &o))
	assert.Nil(t, o.Percent)

	d = []byte(`{"percent":"bad"}`)
	require.ErrorContains(t, json.Unmarshal(d, &o), "invalid major number 'bad', strconv.ParseInt: parsing \"bad\": invalid syntax")
}

func TestPercentageMarshalJSON(t *testing.T) {
	o := struct {
		Percent num.Percentage `json:"percent"`
	}{
		Percent: num.MakePercentage(160, 3),
	}
	d, err := json.Marshal(o)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	de := []byte(`{"percent":"16.0%"}`)
	assert.Equal(t, de, d)
}
