package uuid_test

import (
	"testing"

	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValue(t *testing.T) {
	u := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	v, err := u.Value()
	require.NoError(t, err)
	assert.Equal(t, u.String(), v)
}

func TestScan(t *testing.T) {
	u := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	t.Run("with UUID", func(t *testing.T) {
		var uu uuid.UUID
		err := uu.Scan(u)
		require.NoError(t, err)
		assert.Equal(t, u, uu)
	})
	t.Run("with string", func(t *testing.T) {
		var uu uuid.UUID
		err := uu.Scan(u.String())
		require.NoError(t, err)
		assert.Equal(t, u, uu)
	})
	t.Run("with []byte text", func(t *testing.T) {
		var uu uuid.UUID
		err := uu.Scan([]byte(u.String()))
		require.NoError(t, err)
		assert.Equal(t, u, uu)
	})
	t.Run("with bytes", func(t *testing.T) {
		var uu uuid.UUID
		err := uu.Scan(u.Bytes())
		require.NoError(t, err)
		assert.Equal(t, u, uu)
	})

	t.Run("with int", func(t *testing.T) {
		var uu uuid.UUID
		err := uu.Scan(42)
		require.ErrorContains(t, err, "cannot convert int to UUI")
	})
}
