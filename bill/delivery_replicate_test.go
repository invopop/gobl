package bill_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeliveryReplicate(t *testing.T) {
	dlv := baseDeliveryWithLines(t)
	dlv.IssueTime = cal.NewTime(12, 0, 0)
	dlv.ValueDate = cal.NewDate(2022, 6, 1)
	require.NoError(t, dlv.Calculate())
	require.NoError(t, dlv.Validate())

	assert.Equal(t, "2022-06-13", dlv.IssueDate.String())
	assert.Equal(t, "12:00:00", dlv.IssueTime.String())

	require.NoError(t, dlv.Replicate())

	assert.Empty(t, dlv.UUID)
	assert.Empty(t, dlv.Code)
	assert.True(t, dlv.IssueDate.IsZero())
	assert.Nil(t, dlv.IssueTime)
	assert.Nil(t, dlv.ValueDate)
}
