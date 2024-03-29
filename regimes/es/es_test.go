package es_test

import (
	"testing"

	"github.com/invopop/gobl/regimes/es"
)

func TestTaxRegion(t *testing.T) {
	if err := es.New().Validate(); err != nil {
		t.Errorf("Validation on tax def failed: %v", err.Error())
	}
}
