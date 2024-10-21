package pay_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdvanceNormalize(t *testing.T) {
	a := &pay.Advance{
		Identify:    uuid.Identify{UUID: uuid.Zero},
		Description: "Test advance",
		Percent:     num.NewPercentage(100, 2),
		Ext: tax.Extensions{
			"random": "",
		},
	}
	a.Normalize(nil)
	assert.Empty(t, a.UUID)
	assert.Empty(t, a.Ext)

	a = nil
	assert.NotPanics(t, func() {
		a.Normalize(nil)
	})

}

func TestAdvanceUnmarshal(t *testing.T) {
	a := new(pay.Advance)
	err := json.Unmarshal([]byte(`{"desc":"foo"}`), a)
	require.NoError(t, err)
	assert.Equal(t, "foo", a.Description)
}
