package tax_test

import (
	"testing"

	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestPointDefs(t *testing.T) {
	assert.NotEmpty(t, tax.PointDefs)
	for _, d := range tax.PointDefs {
		assert.NotEmpty(t, d.Key, "key must be set")
		assert.NotEmpty(t, d.Name, "name must be set")
		assert.NotEmpty(t, d.Desc, "desc must be set")
	}
}

func TestPointKeys(t *testing.T) {
	assert.Equal(t, "issue", tax.PointIssue.String())
	assert.Equal(t, "delivery", tax.PointDelivery.String())
	assert.Equal(t, "payment", tax.PointPayment.String())
}
