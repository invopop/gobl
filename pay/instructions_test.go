package pay_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstructionsNormalize(t *testing.T) {
	i := &pay.Instructions{
		Key:    "online",
		Ref:    " fooo ",
		Detail: " Some random payment\t",
		Ext: tax.Extensions{
			"random": "",
		},
	}
	i.Normalize()
	assert.Empty(t, i.Ext)
	assert.Equal(t, "fooo", i.Ref.String())
	assert.Equal(t, "Some random payment", i.Detail)

	i = nil
	assert.NotPanics(t, func() {
		i.Normalize()
	})
}

func TestOnline(t *testing.T) {
	instr := &pay.Instructions{
		Key: pay.MeansKeyOnline,
		Online: []*pay.Online{
			{
				Label: "Test",
				URL:   "https://example.com",
			},
		},
	}
	require.NoError(t, instr.Validate())
	assert.Equal(t, "Test", instr.Online[0].Label)
	assert.Equal(t, "https://example.com", instr.Online[0].URL)

	inst := &pay.Instructions{}
	data := `{"key":"online","online":[{"name":"Test","addr":"https://example.com"}]}`
	require.NoError(t, json.Unmarshal([]byte(data), inst))

	assert.Equal(t, "Test", inst.Online[0].Label)
	assert.Equal(t, "https://example.com", inst.Online[0].URL)
}
