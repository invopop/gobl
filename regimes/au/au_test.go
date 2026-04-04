package au_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/au"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	regime := au.New()

	assert.Equal(t, l10n.AU.Tax(), regime.Country)
	assert.Equal(t, currency.AUD, regime.Currency)
	assert.Equal(t, tax.CategoryGST, regime.TaxScheme)
	assert.Equal(t, "Australia", regime.Name.String())
	assert.NotEmpty(t, regime.Categories)
	assert.NotNil(t, regime.Validator)
	assert.NotNil(t, regime.Normalizer)
	assert.Len(t, regime.Corrections, 1)
	assert.Equal(t, []string{
		bill.InvoiceTypeCreditNote.String(),
		bill.InvoiceTypeDebitNote.String(),
	}, []string{
		regime.Corrections[0].Types[0].String(),
		regime.Corrections[0].Types[1].String(),
	})
	assert.Contains(t, regime.Description.String(), "GST-free")
	assert.Contains(t, regime.Description.String(), "input-taxed")
	assert.True(t, strings.Contains(regime.Description.String(), "exempt key"))
	assert.NotNil(t, regime.CategoryDef(tax.CategoryGST))
	assert.NotNil(t, regime.CategoryDef(tax.CategoryGST).KeyDef(tax.KeyExempt))
}

func TestRegimeValidation(t *testing.T) {
	t.Parallel()

	regime := au.New()
	require.NoError(t, regime.Validate())
}

func TestCorrectionOptionsSchema(t *testing.T) {
	t.Parallel()

	inv := validInvoice()
	require.NoError(t, inv.Calculate())

	out, err := inv.CorrectionOptionsSchema()
	require.NoError(t, err)

	data, err := json.Marshal(out)
	require.NoError(t, err)

	assert.Contains(t, string(data), `"const":"credit-note"`)
	assert.Contains(t, string(data), `"const":"debit-note"`)
}
