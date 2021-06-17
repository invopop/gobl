package tax_test

import (
	"testing"

	"github.com/invopop/gobl/tax"
)

func TestCodeValidation(t *testing.T) {
	c := tax.Code("ABC")
	if err := c.Validate(); err != nil {
		t.Errorf("did not expect error: %v", err)
	}
	c = tax.Code("abc")
	if err := c.Validate(); err == nil {
		t.Errorf("expected a validation error")
	}
	c = tax.Code("ab")
	if err := c.Validate(); err == nil {
		t.Errorf("expected a validation error")
	}
}
