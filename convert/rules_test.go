package convert_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/convert"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequireAddons(t *testing.T) {
	set := convert.Rules("test+multi",
		rules.For(new(bill.Invoice),
			convert.RequireAddons("01", "test-a", "test-b"),
		),
	)
	t.Run("missing", func(t *testing.T) {
		fs := set.Validate(&bill.Invoice{Addons: tax.WithAddons("test-a")})
		require.Error(t, fs)
		assert.Equal(t, rules.Code("TEST-MULTI-BILL-INVOICE-01"), fs.First().Code())
		assert.Equal(t, "requires addons test-a, test-b", fs.First().Message())
	})
	t.Run("included", func(t *testing.T) {
		inv := &bill.Invoice{Addons: tax.WithAddons(cbc.Key("test-b"), "test-a")}
		assert.Nil(t, set.Validate(inv))
	})
}
