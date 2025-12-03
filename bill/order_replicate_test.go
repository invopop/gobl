package bill_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderReplicate(t *testing.T) {
	ord := baseOrderWithLines(t)
	ord.IssueTime = cal.NewTime(12, 0, 0)
	ord.ValueDate = cal.NewDate(2022, 6, 1)
	ord.OperationDate = cal.NewDate(2022, 6, 5)
	require.NoError(t, ord.Calculate())
	require.NoError(t, ord.Validate())

	assert.Equal(t, "2022-06-13", ord.IssueDate.String())
	assert.Equal(t, "12:00:00", ord.IssueTime.String())

	require.NoError(t, ord.Replicate())

	assert.Empty(t, ord.UUID)
	assert.Empty(t, ord.Code)
	assert.True(t, ord.IssueDate.IsZero())
	assert.Nil(t, ord.IssueTime)
	assert.Nil(t, ord.ValueDate)
	assert.Nil(t, ord.OperationDate)
}
