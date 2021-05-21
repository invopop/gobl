package i18n_test

import (
	"testing"

	"github.com/invopop/gobl/i18n"
)

func TestI18nString(t *testing.T) {
	s := i18n.String{
		"en": "Test",
		"es": "Prueba",
	}

	if x := s.String("en"); x != "Test" {
		t.Errorf("Unexpected string result: %v", x)
	}
	if x := s.String("es"); x != "Prueba" {
		t.Errorf("Unexpected string result: %v", x)
	}
	if x := s.String("fo"); x != "Test" {
		t.Errorf("Unexpected string result: %v", x)
	}

	snd := i18n.String{
		i18n.AA: "Foo",
	}
	if x := snd.String("en"); x != "Foo" {
		t.Errorf("Unexpected string result: %v", x)
	}
}
