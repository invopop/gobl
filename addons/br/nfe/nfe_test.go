package nfe_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfe"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNFeValidation(t *testing.T) {
	addon := tax.AddonForKey(nfe.V4)

	t.Run("unsupported struct", func(t *testing.T) {
		obj := new(bill.SubLine)
		err := addon.Validator(obj)
		assert.NoError(t, err)
	})
}
