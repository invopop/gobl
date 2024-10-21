package tax_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComboUnmarshal(t *testing.T) {
	data := []byte(`{"cat":"VAT","tags":["standard"],"percent":"20%"}`)
	var c tax.Combo
	err := json.Unmarshal(data, &c)
	require.NoError(t, err)
	assert.Equal(t, c.Category, cbc.Code("VAT"))
	assert.Equal(t, c.Rate, cbc.Key("standard"))
}
