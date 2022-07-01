package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
)

func TestCodeValidation(t *testing.T) {
	c := org.Code("ABC")
	if err := c.Validate(); err != nil {
		t.Errorf("did not expect error: %v", err)
	}
	c = org.Code("abc")
	if err := c.Validate(); err == nil {
		t.Errorf("expected a validation error")
	}
	c = org.Code("ab")
	if err := c.Validate(); err == nil {
		t.Errorf("expected a validation error")
	}
}
