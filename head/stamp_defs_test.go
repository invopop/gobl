package head_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterStamps(t *testing.T) {
	stampA := &cbc.Definition{
		Key:  "test-stamp-a",
		Name: i18n.String{i18n.EN: "Test Stamp A"},
	}
	head.RegisterStamps(
		&cbc.Definition{
			Key:  "test-stamp-b",
			Name: i18n.String{i18n.EN: "Test Stamp B"},
		},
		stampA,
	)

	t.Run("lookup by key", func(t *testing.T) {
		sd := head.StampDefForKey("test-stamp-a")
		require.NotNil(t, sd)
		assert.Equal(t, "Test Stamp A", sd.Name.In(i18n.EN))
	})

	t.Run("unknown key", func(t *testing.T) {
		assert.Nil(t, head.StampDefForKey("test-stamp-missing"))
	})

	t.Run("same definition is ignored", func(t *testing.T) {
		assert.NotPanics(t, func() {
			head.RegisterStamps(stampA)
		})
		count := 0
		for _, sd := range head.AllStampDefs() {
			if sd.Key == "test-stamp-a" {
				count++
			}
		}
		assert.Equal(t, 1, count)
	})

	t.Run("different definition panics", func(t *testing.T) {
		assert.PanicsWithValue(t, "stamp 'test-stamp-a' already registered", func() {
			head.RegisterStamps(&cbc.Definition{
				Key:  "test-stamp-a",
				Name: i18n.String{i18n.EN: "Other Stamp A"},
			})
		})
		assert.Same(t, stampA, head.StampDefForKey("test-stamp-a"))
	})

	t.Run("all sorted by key", func(t *testing.T) {
		all := head.AllStampDefs()
		for i := 1; i < len(all); i++ {
			assert.Less(t, all[i-1].Key.String(), all[i].Key.String())
		}
	})
}
