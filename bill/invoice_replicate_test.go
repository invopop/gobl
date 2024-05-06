package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceReplicate(t *testing.T) {
	lines := []*bill.Line{
		{
			Quantity: num.MakeAmount(2, 0),
			Item: &org.Item{
				Name:  "Test Item",
				Price: num.MakeAmount(10, 0),
			},
		},
	}
	inv := baseInvoice(t, lines...)
	inv.ValueDate = cal.NewDate(2022, 6, 1)
	inv.OperationDate = cal.NewDate(2022, 6, 5)
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	assert.Equal(t, "2022-06-13", inv.IssueDate.String())

	require.NoError(t, inv.Replicate())

	assert.Empty(t, inv.UUID)
	assert.Empty(t, inv.Code)
	td := cal.Today()
	assert.Equal(t, inv.IssueDate.String(), td.String())
	assert.Nil(t, inv.ValueDate)
	assert.Nil(t, inv.OperationDate)
}
