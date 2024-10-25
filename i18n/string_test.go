package i18n_test

import (
	"testing"

	"github.com/invopop/gobl/i18n"
	"github.com/stretchr/testify/assert"
)

func TestI18nString(t *testing.T) {
	s := i18n.String{
		"en": "Test",
		"es": "Prueba",
	}

	assert.Equal(t, "Test", s.In("en"))
	assert.Equal(t, "Prueba", s.In("es"))
	assert.Equal(t, "Test", s.In("fo"))
	assert.Equal(t, "Test", s.String())

	snd := i18n.String{
		i18n.AA: "Foo",
	}
	assert.Equal(t, "Foo", snd.In("en"))
	assert.Equal(t, "Foo", snd.String())

	s2 := i18n.NewString("Test")
	assert.Equal(t, "Test", s2.In("en"))
}
