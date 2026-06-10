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

func TestDigestEquals(t *testing.T) {
	a := dsig.NewSHA256Digest([]byte("hello"))
	b := dsig.NewSHA256Digest([]byte("hello"))
	c := dsig.NewSHA256Digest([]byte("world"))

	assert.NoError(t, a.Equals(b))

	err := a.Equals(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mismatch")

	wrongAlg := &dsig.Digest{Algorithm: "md5", Value: a.Value}
	err = a.Equals(wrongAlg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "algorithm mismatch")
}

func TestDigestString(t *testing.T) {
	d := &dsig.Digest{Algorithm: dsig.DigestSHA256, Value: "abc"}
	assert.Equal(t, "sha256;abc", d.String())
}
