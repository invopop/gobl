package cbc_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
)

type testMetaStruct struct {
	Meta cbc.Meta
}

func (tms *testMetaStruct) Validate() error {
	return validation.ValidateStruct(tms,
		validation.Field(&tms.Meta),
	)
}

func TestMeta(t *testing.T) {
	v := new(testMetaStruct)
	v.Meta = cbc.Meta{
		cbc.Key("test"): "bar",
	}
	err := v.Validate()
	assert.NoError(t, err)

	v.Meta = cbc.Meta{
		cbc.Key("bad_key"): "bar",
	}
	err = v.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Meta: (bad_key: must be in a valid format.)")
}
