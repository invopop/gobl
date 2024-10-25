package iso_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	ext := tax.ExtensionForKey("iso-scheme-id")
	assert.NotNil(t, ext)
}
