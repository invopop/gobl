package tax_test

import (
	"testing"

	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestRegimeJSONSchema(t *testing.T) {
	rc := tax.RegimeCode("")
	js := rc.JSONSchema()

	assert.Equal(t, "Tax Regime Code", js.Title)
	assert.Equal(t, "string", js.Type)
	assert.Greater(t, len(js.OneOf), 1)
	rd := tax.AllRegimeDefs()[0]
	assert.Equal(t, rd.Code().String(), js.OneOf[0].Const)
}

func TestRegimeCode(t *testing.T) {
	rc := tax.RegimeCode("US")
	assert.Equal(t, "US", rc.String())
	assert.Equal(t, "US", rc.Code().String())
}
