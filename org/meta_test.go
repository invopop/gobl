package org_test

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

type testMetaStruct struct {
	Meta org.Meta
}

func (tms *testMetaStruct) Validate() error {
	return validation.ValidateStruct(tms,
		validation.Field(&tms.Meta),
	)
}

func TestMeta(t *testing.T) {
	v := new(testMetaStruct)
	v.Meta = org.Meta{
		org.Key("test"): "bar",
	}
	err := v.Validate()
	assert.NoError(t, err)

	v.Meta = org.Meta{
		org.Key("bad_key"): "bar",
	}
	err = v.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Meta: (bad_key: must be in a valid format.)")
}
