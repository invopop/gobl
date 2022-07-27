package currency_test

import (
	"testing"

	"github.com/invopop/gobl/currency"
	"github.com/stretchr/testify/assert"
)

func TestCode(t *testing.T) {
	c := currency.EUR
	assert.NoError(t, c.Validate())

	d := c.Def()
	assert.Equal(t, d.Name, "Euro")
}
