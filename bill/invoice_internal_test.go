package bill

import (
	"testing"

	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

// TestInvoiceNoLinesHasBypass covers the guard of the rule requiring bypass
// with no-lines, for invoices passed by value or pointer and for other values.
func TestInvoiceNoLinesHasBypass(t *testing.T) {
	noLines := Invoice{Tags: tax.WithTags(tax.TagNoLines)}
	assert.False(t, invoiceNoLinesHasBypass(noLines))
	assert.False(t, invoiceNoLinesHasBypass(&noLines))

	withBypass := Invoice{Tags: tax.WithTags(tax.TagNoLines, tax.TagBypass)}
	assert.True(t, invoiceNoLinesHasBypass(withBypass))
	assert.True(t, invoiceNoLinesHasBypass(&withBypass))

	assert.True(t, invoiceNoLinesHasBypass(Invoice{}))
	assert.True(t, invoiceNoLinesHasBypass((*Invoice)(nil)))
	assert.True(t, invoiceNoLinesHasBypass("not an invoice"))
}
