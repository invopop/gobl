package lu_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/lu"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()
	regime := lu.New()
	require.NotNil(t, regime)

	assert.Equal(t, l10n.LU, regime.Country.Code())
	assert.Equal(t, "Luxembourg", regime.Name.String())
	assert.NotEmpty(t, regime.Description)

	require.Len(t, regime.Categories, 1)
	cat := regime.Categories[0]
	assert.Equal(t, tax.CategoryVAT, cat.Code)
	assert.Len(t, cat.Rates, 4, "expect standard, intermediate, reduced, and super-reduced rates")

	require.Len(t, regime.Identities, 1)
	assert.Equal(t, lu.IdentityTypeRCS, regime.Identities[0].Code)

	require.Len(t, regime.Corrections, 1)
	assert.Equal(t, bill.ShortSchemaInvoice, regime.Corrections[0].Schema)
	assert.Contains(t, regime.Corrections[0].Types, bill.InvoiceTypeCreditNote)
	assert.Contains(t, regime.Corrections[0].Types, bill.InvoiceTypeDebitNote)
}
