package bis

import (
	"testing"
	"time"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/assert"
)

// d builds a test date in the fixed test year 2026.
func d(month time.Month, day int) cal.Date {
	return cal.MakeDate(2026, month, day)
}

func TestLineStartsWithinInvoice(t *testing.T) {
	t.Run("nil/wrong type passes", func(t *testing.T) {
		assert.True(t, lineStartsWithinInvoice(nil))
	})
	t.Run("no ordering passes", func(t *testing.T) {
		assert.True(t, lineStartsWithinInvoice(&bill.Invoice{}))
	})
	t.Run("zero invoice start passes", func(t *testing.T) {
		assert.True(t, lineStartsWithinInvoice(&bill.Invoice{
			Ordering: &bill.Ordering{Period: &cal.Period{}},
		}))
	})
	t.Run("line start before invoice start fails", func(t *testing.T) {
		inv := &bill.Invoice{
			Ordering: &bill.Ordering{Period: &cal.Period{Start: d(4, 10)}},
			Lines:    []*bill.Line{{Period: &cal.Period{Start: d(4, 1)}}},
		}
		assert.False(t, lineStartsWithinInvoice(inv))
	})
	t.Run("line start equal/after invoice start passes", func(t *testing.T) {
		inv := &bill.Invoice{
			Ordering: &bill.Ordering{Period: &cal.Period{Start: d(4, 10)}},
			Lines:    []*bill.Line{{Period: &cal.Period{Start: d(4, 15)}}},
		}
		assert.True(t, lineStartsWithinInvoice(inv))
	})
	t.Run("nil line and zero line period skipped", func(t *testing.T) {
		inv := &bill.Invoice{
			Ordering: &bill.Ordering{Period: &cal.Period{Start: d(4, 10)}},
			Lines:    []*bill.Line{nil, {Period: nil}, {Period: &cal.Period{}}},
		}
		assert.True(t, lineStartsWithinInvoice(inv))
	})
}

func TestLineEndsWithinInvoice(t *testing.T) {
	t.Run("nil/wrong type passes", func(t *testing.T) {
		assert.True(t, lineEndsWithinInvoice(nil))
	})
	t.Run("no ordering passes", func(t *testing.T) {
		assert.True(t, lineEndsWithinInvoice(&bill.Invoice{}))
	})
	t.Run("zero invoice end passes", func(t *testing.T) {
		assert.True(t, lineEndsWithinInvoice(&bill.Invoice{
			Ordering: &bill.Ordering{Period: &cal.Period{}},
		}))
	})
	t.Run("line end after invoice end fails", func(t *testing.T) {
		inv := &bill.Invoice{
			Ordering: &bill.Ordering{Period: &cal.Period{End: d(4, 30)}},
			Lines:    []*bill.Line{{Period: &cal.Period{End: d(5, 5)}}},
		}
		assert.False(t, lineEndsWithinInvoice(inv))
	})
	t.Run("line end on/before invoice end passes", func(t *testing.T) {
		inv := &bill.Invoice{
			Ordering: &bill.Ordering{Period: &cal.Period{End: d(4, 30)}},
			Lines:    []*bill.Line{{Period: &cal.Period{End: d(4, 25)}}},
		}
		assert.True(t, lineEndsWithinInvoice(inv))
	})
}
