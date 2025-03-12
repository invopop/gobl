package dsig_test

import (
	"testing"

	"github.com/invopop/gobl/dsig"
	"github.com/stretchr/testify/assert"
)

func TestNewSHA256Digest(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		d := dsig.NewSHA256Digest([]byte("hello"))
		assert.NotNil(t, d)
		assert.Equal(t, dsig.DigestSHA256, d.Algorithm)
		assert.Equal(t, "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824", d.Value)
	})
}
