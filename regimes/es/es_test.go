package es_test

import (
	"testing"

	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/rules"
)

func TestTaxRegion(t *testing.T) {
	if err := rules.Validate(es.New()); err != nil {
		t.Errorf("Validation on tax def failed: %v", err.Error())
	}
}
