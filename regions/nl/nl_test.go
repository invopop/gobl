package nl_test

import (
	"testing"

	"github.com/invopop/gobl/regions/nl"
)

func TestTaxRegion(t *testing.T) {
	tr := nl.New()
	if err := tr.Validate(); err != nil {
		t.Errorf("Validation on tax def failed: %v", err.Error())
	}
}
