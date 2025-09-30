package gb_test

import (
	"testing"

	"github.com/invopop/gobl/regimes/gb"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	t.Run("normalize tax IDs", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "GB",
			Code:    "844 281.425",
		}
		gb.New().Normalizer(tID)
		assert.Equal(t, "844281425", tID.Code.String())
	})

}
