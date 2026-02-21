package au_test

import (
	"testing"

	"github.com/invopop/gobl/regimes/au"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	t.Run("normalize tax IDs", func(t *testing.T) {
		tID := &tax.Identity{
			Country: "AU",
			Code:    "51 824 753 556",
		}
		au.New().Normalizer(tID)
		assert.Equal(t, "51824753556", tID.Code.String())
	})
}
