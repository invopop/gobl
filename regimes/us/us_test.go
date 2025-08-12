package us_test

import (
	"testing"

	"github.com/invopop/gobl/regimes/us"
)

func TestNew(t *testing.T) {
	t.Run("should create a new US regime", func(t *testing.T) {
		regime := us.New()
		if regime == nil {
			t.Fatal("expected non-nil regime")
		}
		if regime.Country != "US" {
			t.Errorf("expected country code 'US', got '%s'", regime.Country)
		}
		if regime.Name["en"] != "United States of America" {
			t.Errorf("expected name 'United States of America', got '%s'", regime.Name["en"])
		}
	})
}
