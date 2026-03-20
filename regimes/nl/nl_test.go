package nl_test

import (
	"testing"

	"github.com/invopop/gobl/regimes/nl"
	"github.com/invopop/gobl/rules"
)

func TestTaxRegion(t *testing.T) {
	tr := nl.New()
	if err := rules.Validate(tr); err != nil {
		t.Errorf("Validation on tax def failed: %v", err.Error())
	}
}
