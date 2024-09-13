package bill_test

import (
	"context"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaxValidation(t *testing.T) {
	es := es.New()
	ctx := tax.ContextWithValidator(context.Background(), es.Validator)
	tx := &bill.Tax{
		Tags: []cbc.Key{"reverse-charge"},
	}

	err := tx.ValidateWithContext(ctx)
	require.NoError(t, err)

	tx.Tags = []cbc.Key{"invalid-tag"}
	err = tx.ValidateWithContext(ctx)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be a valid value")
}
