package untdid_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	ext := tax.ExtensionForKey("untdid-tax-category")
	assert.NotNil(t, ext)
}
