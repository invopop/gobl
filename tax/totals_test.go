package tax_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

func TestTotalAddRate(t *testing.T) {
	zero := num.MakeAmount(0, 2)
	r1 := tax.Rate{
		Category: tax.Code("VAT"),
		Code:     tax.Code("STD"),
		Base:     num.MakeAmount(10000, 2),
		Percent:  num.MakePercentage(200, 3),
		Value:    num.MakeAmount(2000, 2),
		Retained: false,
	}
	r2 := tax.Rate{
		Category: tax.Code("VAT"),
		Code:     tax.Code("RED"),
		Base:     num.MakeAmount(5000, 2),
		Percent:  num.MakePercentage(50, 3),
		Value:    num.MakeAmount(250, 2),
		Retained: false,
	}
	tt := tax.NewTotal(zero)
	if err := tt.AddRate(r1, zero); err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if err := tt.AddRate(r2, zero); err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if len(tt.Categories) != 1 {
		t.Errorf("unexpected catagories list: %v", tt.Categories)
	}
	cat := tt.Categories[0]
	if len(cat.Rates) != 2 {
		t.Errorf("unexpected list of rates")
	}
	if !cat.Base.Equals(num.MakeAmount(15000, 2)) {
		t.Errorf("unexpected base amount, got: %v", cat.Base.String())
	}
	if !cat.Value.Equals(num.MakeAmount(2250, 2)) {
		t.Errorf("unexpected value total, got: %v", cat.Value.String())
	}
	r3 := tax.Rate{
		Category: tax.Code("VAT"),
		Code:     tax.Code("STD"),
		Base:     num.MakeAmount(8000, 2),
		Percent:  num.MakePercentage(200, 3),
		Value:    num.MakeAmount(1600, 2),
		Retained: false,
	}
	if err := tt.AddRate(r3, zero); err != nil {
		t.Errorf("did not expect error adding 3rd rate: %v", err.Error())
	}
	cat = tt.Categories[0]
	if len(cat.Rates) != 2 {
		t.Errorf("unexpected number of rates")
	}
	if !cat.Base.Equals(num.MakeAmount(23000, 2)) {
		t.Errorf("unexpected base amount, got: %v", cat.Base.String())
	}
	if !cat.Value.Equals(num.MakeAmount(3850, 2)) {
		t.Errorf("unexpected value total, got: %v", cat.Value.String())
	}
	if !tt.Sum.Equals(num.MakeAmount(3850, 2)) {
		t.Errorf("unexpected sum: %v", tt.Sum.String())
	}
	r4 := tax.Rate{
		Category: tax.Code("IRPF"),
		Code:     tax.Code("STD"),
		Base:     num.MakeAmount(23000, 2),
		Percent:  num.MakePercentage(150, 3),
		Value:    num.MakeAmount(3450, 2),
		Retained: true,
	}
	if err := tt.AddRate(r4, zero); err != nil {
		t.Errorf("did not expect error adding 4th rate: %v", err.Error())
	}
	if len(tt.Categories) != 2 {
		t.Errorf("unexpected category length")
		return
	}
	cat = tt.Categories[1]
	if len(cat.Rates) != 1 {
		t.Errorf("unexpected number of rates")
	}
	if !cat.Base.Equals(num.MakeAmount(23000, 2)) {
		t.Errorf("unexpected base amount, got: %v", cat.Base.String())
	}
	if !cat.Value.Equals(num.MakeAmount(3450, 2)) {
		t.Errorf("unexpected value total, got: %v", cat.Value.String())
	}
	if !tt.Sum.Equals(num.MakeAmount(400, 2)) {
		t.Errorf("unexpected total sum, got: %v", tt.Sum.String())
	}
	data, _ := json.Marshal(tt)
	t.Logf("TT: %+v", string(data))
}
