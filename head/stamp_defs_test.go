package head_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterStamps(t *testing.T) {
	head.RegisterAddonStamps("test-addon-v1",
		&cbc.Definition{
			Key:  "test-stamp-b",
			Name: i18n.String{i18n.EN: "Test Stamp B"},
		},
		&cbc.Definition{
			Key:  "test-stamp-a",
			Name: i18n.String{i18n.EN: "Test Stamp A"},
		},
	)
	head.RegisterRegimeStamps("XX",
		&cbc.Definition{
			Key:  "test-stamp-regime",
			Name: i18n.String{i18n.EN: "Test Regime Stamp"},
		},
	)

	t.Run("addon stamp", func(t *testing.T) {
		sd := head.StampDefForKey("test-stamp-a")
		require.NotNil(t, sd)
		assert.Equal(t, cbc.Key("test-addon-v1"), sd.Addon)
		assert.Empty(t, sd.Regime)
		assert.Equal(t, "Test Stamp A", sd.Name.In(i18n.EN))
	})

	t.Run("regime stamp", func(t *testing.T) {
		sd := head.StampDefForKey("test-stamp-regime")
		require.NotNil(t, sd)
		assert.Equal(t, "XX", sd.Regime.String())
		assert.Empty(t, sd.Addon)
	})

	t.Run("unknown key", func(t *testing.T) {
		assert.Nil(t, head.StampDefForKey("test-stamp-missing"))
	})

	t.Run("same owner replaces", func(t *testing.T) {
		head.RegisterAddonStamps("test-addon-v1", &cbc.Definition{
			Key:  "test-stamp-b",
			Name: i18n.String{i18n.EN: "Test Stamp B2"},
		})
		assert.Equal(t, "Test Stamp B2", head.StampDefForKey("test-stamp-b").Name.In(i18n.EN))
		count := 0
		for _, sd := range head.AllStampDefs() {
			if sd.Key == "test-stamp-b" {
				count++
			}
		}
		assert.Equal(t, 1, count)
	})

	t.Run("other addon panics", func(t *testing.T) {
		assert.PanicsWithValue(t,
			"stamp 'test-stamp-a' from addon 'other-addon-v1' already registered by addon 'test-addon-v1'",
			func() {
				head.RegisterAddonStamps("other-addon-v1", &cbc.Definition{
					Key:  "test-stamp-a",
					Name: i18n.String{i18n.EN: "Other Stamp A"},
				})
			},
		)
		assert.Equal(t, cbc.Key("test-addon-v1"), head.StampDefForKey("test-stamp-a").Addon)
	})

	t.Run("regime panics on addon stamp", func(t *testing.T) {
		assert.PanicsWithValue(t,
			"stamp 'test-stamp-a' from regime 'XX' already registered by addon 'test-addon-v1'",
			func() {
				head.RegisterRegimeStamps("XX", &cbc.Definition{
					Key:  "test-stamp-a",
					Name: i18n.String{i18n.EN: "Other Stamp A"},
				})
			},
		)
	})

	t.Run("all sorted by key", func(t *testing.T) {
		all := head.AllStampDefs()
		for i := 1; i < len(all); i++ {
			assert.Less(t, all[i-1].Key.String(), all[i].Key.String())
		}
	})

	t.Run("json", func(t *testing.T) {
		data, err := json.Marshal(head.StampDefForKey("test-stamp-a"))
		require.NoError(t, err)
		assert.JSONEq(t, `{"addon":"test-addon-v1","key":"test-stamp-a","name":{"en":"Test Stamp A"}}`, string(data))

		data, err = json.Marshal(head.StampDefForKey("test-stamp-regime"))
		require.NoError(t, err)
		assert.JSONEq(t, `{"regime":"XX","key":"test-stamp-regime","name":{"en":"Test Regime Stamp"}}`, string(data))
	})
}
