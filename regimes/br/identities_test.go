package br_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/br"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeIdentities(t *testing.T) {
	t.Run("nil identity", func(t *testing.T) {
		var ident *org.Identity
		br.Normalize(ident)
		assert.Nil(t, ident)
	})

	t.Run("migrates old key br-nfse-municipal-reg", func(t *testing.T) {
		ident := &org.Identity{
			Key:  "br-nfse-municipal-reg",
			Code: "1234567890",
		}
		br.Normalize(ident)
		assert.Equal(t, cbc.Key("br-municipal-reg"), ident.Key)
		assert.Equal(t, cbc.Code("1234567890"), ident.Code)
	})

	t.Run("migrates old key br-nfse-national-reg", func(t *testing.T) {
		ident := &org.Identity{
			Key:  "br-nfse-national-reg",
			Code: "1234567890",
		}
		br.Normalize(ident)
		assert.Equal(t, cbc.Key("br-state-reg"), ident.Key)
		assert.Equal(t, cbc.Code("1234567890"), ident.Code)
	})
}
