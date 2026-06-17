package nz_test

import (
	"testing"

	"github.com/invopop/gobl/regimes/nz"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	t.Run("normalize tax identity", func(t *testing.T) {
		id := &tax.Identity{
			Country: "NZ",
			Code:    "123-456-785",
		}
		nz.New().Normalizer(id)
		assert.Equal(t, "123456785", id.Code.String())
	})
}
