package nl_test

import (
	"testing"

	"github.com/invopop/gobl/regimes/nl"
)

func TestTaxRegion(t *testing.T) {
	tr := nl.New()
	if err := tr.Validate(); err != nil {
		t.Errorf("Validation on tax def failed: %v", err.Error())
	}
}
