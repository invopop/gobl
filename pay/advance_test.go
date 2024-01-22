package pay_test

import (
	"testing"

	"github.com/invopop/gobl/pay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdvanceUnmarshal(t *testing.T) {
	a := new(pay.Advance)
	err := a.JSONUnmarshal([]byte(`{"desc":"foo"}`))
	require.NoError(t, err)
	assert.Equal(t, "foo", a.Description)
}
