package nfe_test

import (
	"testing"

	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
)

func TestNFeValidation(t *testing.T) {
	t.Run("unsupported struct", func(t *testing.T) {
		// A simple struct with no NFe rules should not produce NFe-specific errors
		type simpleStruct struct {
			Name string `json:"name"`
		}
		obj := &simpleStruct{Name: "test"}
		err := rules.Validate(obj)
		assert.NoError(t, err)
	})
}
