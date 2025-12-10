package bill_test

import (
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentReplicate(t *testing.T) {
	pmt := testPaymentMinimal(t)
	pmt.IssueTime = cal.NewTime(12, 0, 0)
	pmt.ValueDate = cal.NewDate(2025, 1, 25)
	require.NoError(t, pmt.Calculate())
	require.NoError(t, pmt.Validate())

	assert.Equal(t, "2025-01-24", pmt.IssueDate.String())
	assert.Equal(t, "12:00:00", pmt.IssueTime.String())

	require.NoError(t, pmt.Replicate())

	assert.Empty(t, pmt.UUID)
	assert.Empty(t, pmt.Code)
	assert.True(t, pmt.IssueDate.IsZero())
	assert.Nil(t, pmt.IssueTime)
	assert.Nil(t, pmt.ValueDate)
}
