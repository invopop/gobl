package nfe_test

import (
"testing"

"github.com/invopop/gobl/addons/br/nfe"
"github.com/invopop/gobl/bill"
"github.com/invopop/gobl/num"
"github.com/invopop/gobl/org"
"github.com/invopop/gobl/rules"
"github.com/invopop/gobl/tax"
"github.com/stretchr/testify/assert"
)

func TestNFeValidation(t *testing.T) {
t.Run("unsupported struct - no NFe-specific errors", func(t *testing.T) {
// A SubLine should not trigger any NFe-specific validation rules
obj := &bill.SubLine{
Index:    1,
Quantity: num.MakeAmount(1, 0),
Item:     &org.Item{Name: "Test Item"},
}
err := rules.Validate(obj, tax.AddonContext(nfe.V4))
assert.NoError(t, err)
})
}
