package convert_test

import (
	"testing"

	"github.com/invopop/gobl/convert"
	"github.com/stretchr/testify/assert"
)

type keyA struct{}
type keyB struct{}

func TestInput(t *testing.T) {
	in := convert.NewInput([]byte("data"))
	assert.Equal(t, []byte("data"), in.Data)

	t.Run("missing", func(t *testing.T) {
		v, ok := in.Get(keyA{})
		assert.False(t, ok)
		assert.Nil(t, v)
	})
	t.Run("set and get", func(t *testing.T) {
		in.Set(keyA{}, "a")
		v, ok := in.Get(keyA{})
		assert.True(t, ok)
		assert.Equal(t, "a", v)
	})
	t.Run("distinct key types", func(t *testing.T) {
		in.Set(keyB{}, "b")
		va, _ := in.Get(keyA{})
		vb, _ := in.Get(keyB{})
		assert.Equal(t, "a", va)
		assert.Equal(t, "b", vb)
	})
}
