package verifactu_test

import (
	"testing"

	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusValidation(t *testing.T) {
	t.Run("valid status", func(t *testing.T) {
		st := testStatusSystemEvent(t)
		require.NoError(t, rules.Validate(st))
	})

	t.Run("missing supplier", func(t *testing.T) {
		st := testStatusSystemEvent(t)
		st.Supplier = nil
		err := rules.Validate(st)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "supplier is required")
	})

	t.Run("missing supplier tax ID", func(t *testing.T) {
		st := testStatusSystemEvent(t)
		st.Supplier.TaxID = nil
		err := rules.Validate(st)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "supplier tax ID is required")
	})

	t.Run("missing supplier tax code", func(t *testing.T) {
		st := testStatusSystemEvent(t)
		st.Supplier.TaxID.Code = ""
		err := rules.Validate(st)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "supplier tax ID code is required")
	})

	t.Run("missing event type extension", func(t *testing.T) {
		st := testStatusSystemEvent(t)
		delete(st.Ext, verifactu.ExtKeyEventType)
		err := rules.Validate(st)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "event type extension is required")
	})
}

func testStatusSystemEvent(t *testing.T) *bill.Status {
	t.Helper()
	return &bill.Status{
		Addons: tax.WithAddons(verifactu.V1),
		Type:   "system",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Ext: tax.Extensions{
			verifactu.ExtKeyEventType: "01",
		},
	}
}
