package gb_test

import (
	"testing"

	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	t.Run("normalize tax IDs", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "GB",
			Code:    "844 281.425",
		}
		norm.Normalize(tID)
		assert.Equal(t, "844281425", tID.Code.String())
	})

}
