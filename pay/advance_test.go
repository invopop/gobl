package pay_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/pay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdvanceUnmarshal(t *testing.T) {
	a := new(pay.Advance)
	err := json.Unmarshal([]byte(`{"desc":"foo"}`), a)
	require.NoError(t, err)
	assert.Equal(t, "foo", a.Description)
}
