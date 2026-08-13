package cbc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello, World!", "Hello, World!"},
		{"Hello, \u0000World!", "Hello, World!"},
		{"\u0000Hello, World!", "Hello, World!"},
		{"Hello, World!\u0000", "Hello, World!"},
		{"\t hello ", "hello"},
		{"   \t\n", ""},
		{"  Hello  World  ", "Hello  World"},
		{"\n\nHello, World!\n", "Hello, World!"},
		{"A\u0000\u0000B", "AB"},
		{string([]byte{0xff, 'H', 'i'}), "?Hi"},    // invalid UTF-8 byte replaced
		{string([]byte("Hola\xC3\x28")), "Hola?("}, // invalid UTF-8 sequence replaced
		{"", ""},
	}

	for _, test := range tests {
		result := cbc.NormalizeString(test.input)
		assert.Equal(t, test.expected, result, test.input)
	}
}
