package pay_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/norm"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstructionsNormalize(t *testing.T) {
	i := &pay.Instructions{
		Key:    "online",
		Ref:    " fooo ",
		Detail: " Some random payment\t",
		Ext: tax.ExtensionsOf(cbc.CodeMap{
			"random": "",
		}),
	}
	norm.Normalize(i)
	assert.True(t, i.Ext.IsZero())
	assert.Equal(t, "fooo", i.Ref.String())
	assert.Equal(t, "Some random payment", i.Detail)

	i = nil
	assert.NotPanics(t, func() {
		norm.Normalize(i)
	})
}

func TestCreditTransferNormalize(t *testing.T) {
	ct := &pay.CreditTransfer{
		IBAN:     " DK50 0040 0440 1162 43 ",
		BIC:      " DABADKKK ",
		Number:   " 0440-116243 ",
		Clearing: " 1234 ",
		Name:     "Example Bank",
	}

	norm.Normalize(ct)

	assert.Equal(t, cbc.Code("DK50 0040 0440 1162 43"), ct.IBAN)
	assert.Equal(t, cbc.Code("DABADKKK"), ct.BIC)
	assert.Equal(t, cbc.Code("0440-116243"), ct.Number)
	assert.Equal(t, cbc.Code("1234"), ct.Clearing)
	assert.Equal(t, "Example Bank", ct.Name)

	data, err := json.Marshal(ct)
	require.NoError(t, err)
	assert.JSONEq(t, `{
		"iban": "DK50 0040 0440 1162 43",
		"bic": "DABADKKK",
		"number": "0440-116243",
		"clearing": "1234",
		"name": "Example Bank"
	}`, string(data))
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
	require.NoError(t, rules.Validate(instr))
	assert.Equal(t, "Test", instr.Online[0].Label)
	assert.Equal(t, "https://example.com", instr.Online[0].URL)

	inst := &pay.Instructions{}
	data := `{"key":"online","online":[{"name":"Test","addr":"https://example.com"}]}`
	require.NoError(t, json.Unmarshal([]byte(data), inst))

	assert.Equal(t, "Test", inst.Online[0].Label)
	assert.Equal(t, "https://example.com", inst.Online[0].URL)
}
