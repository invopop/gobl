package tax_test

import (
	"testing"

	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestDateDefs(t *testing.T) {
	assert.NotEmpty(t, tax.DateDefs)
	for _, d := range tax.DateDefs {
		assert.NotEmpty(t, d.Key, "key must be set")
		assert.NotEmpty(t, d.Name, "name must be set")
		assert.NotEmpty(t, d.Desc, "desc must be set")
	}
}

func TestDateKeys(t *testing.T) {
	assert.Equal(t, "issue", tax.DateIssue.String())
	assert.Equal(t, "delivery", tax.DateDelivery.String())
	assert.Equal(t, "paid", tax.DatePaid.String())
}
