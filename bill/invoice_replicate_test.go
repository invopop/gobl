package bill_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceReplicate(t *testing.T) {
	inv := baseInvoiceWithLines(t)
	inv.IssueTime = cal.NewTime(12, 0, 0)
	inv.ValueDate = cal.NewDate(2022, 6, 1)
	inv.OperationDate = cal.NewDate(2022, 6, 5)
	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())

	assert.Equal(t, "2022-06-13", inv.IssueDate.String())
	assert.Equal(t, "12:00:00", inv.IssueTime.String())

	require.NoError(t, inv.Replicate())

	assert.Empty(t, inv.UUID)
	assert.Empty(t, inv.Code)
	assert.True(t, inv.IssueDate.IsZero())
	assert.Nil(t, inv.IssueTime)
	assert.Nil(t, inv.ValueDate)
	assert.Nil(t, inv.OperationDate)
}
