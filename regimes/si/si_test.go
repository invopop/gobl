package si_test

import (
	"testing"

	"github.com/invopop/gobl/regimes/si"
	"github.com/invopop/gobl/rules"
)

func TestTaxRegion(t *testing.T) {
	tr := si.New()
	if err := rules.Validate(tr); err != nil {
		t.Errorf("Validation on tax def failed: %v", err.Error())
	}
}
