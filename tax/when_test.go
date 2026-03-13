package tax_test

import (
	"testing"

	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestWhenDefs(t *testing.T) {
	assert.NotEmpty(t, tax.WhenDefs)
	for _, d := range tax.WhenDefs {
		assert.NotEmpty(t, d.Key, "key must be set")
		assert.NotEmpty(t, d.Name, "name must be set")
		assert.NotEmpty(t, d.Desc, "desc must be set")
	}
}

func TestWhenKeys(t *testing.T) {
	assert.Equal(t, "issue", tax.WhenIssue.String())
	assert.Equal(t, "delivery", tax.WhenDelivery.String())
	assert.Equal(t, "paid", tax.WhenPaid.String())
}
