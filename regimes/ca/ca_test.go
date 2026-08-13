package ca_test

import (
	"testing"

	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	t.Run("normalize tax IDs", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "CA",
			Code:    "123.456.789",
		}
		norm.Normalize(tID)
		assert.Equal(t, "123456789", tID.Code.String())
	})

}
