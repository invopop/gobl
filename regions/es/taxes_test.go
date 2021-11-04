package es_test

import (
	"testing"

	"github.com/invopop/gobl/regions/es"
)

func TestTaxRegion(t *testing.T) {
	r := es.New()
	tr := r.Taxes()
	if err := tr.Validate(); err != nil {
		t.Errorf("Validation on tax def failed: %v", err.Error())
	}
}
