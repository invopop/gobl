package currency_test

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCode(t *testing.T) {
	c := currency.EUR
	assert.NoError(t, c.Validate())

	assert.Equal(t, "EUR", c.String())

	d := c.Def()
	assert.Equal(t, d.Name, "Euro")

	c = currency.CodeEmpty
	assert.NoError(t, c.Validate())

	c = currency.Code("FOOO")
	err := c.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "currency code FOOO not defined")
}

func TestCodeJSONSchema(t *testing.T) {
	var c currency.Code
	s := c.JSONSchema()
	assert.Equal(t, "Currency Code", s.Title)
	assert.Equal(t, "string", s.Type)
	assert.Len(t, s.OneOf, 180)
	assert.Equal(t, currency.USD, s.OneOf[0].Const)
}
