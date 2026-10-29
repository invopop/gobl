package bill

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBillableSetTax covers the tax setter on each of the billable documents,
// used to add a rounding rule to documents that carry no tax object.
func TestBillableSetTax(t *testing.T) {
	docs := map[string]billable{
		"invoice":  new(Invoice),
		"order":    new(Order),
		"delivery": new(Delivery),
	}
	for name, doc := range docs {
		t.Run(name, func(t *testing.T) {
			require.Nil(t, doc.getTax())
			tx := new(Tax)
			doc.setTax(tx)
			assert.Same(t, tx, doc.getTax())
		})
	}
}
