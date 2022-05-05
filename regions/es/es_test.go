package es_test

import (
	"testing"

	"github.com/invopop/gobl/regions/es"
)

func TestTaxRegion(t *testing.T) {
	tr := es.New()
	if err := tr.Validate(); err != nil {
		t.Errorf("Validation on tax def failed: %v", err.Error())
	}
}
