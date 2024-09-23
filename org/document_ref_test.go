package org_test

import (
	"context"
	"testing"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestDocumentRefValidation(t *testing.T) {
	dr := new(org.DocumentRef)
	dr.Code = "FOO"
	dr.IssueDate = cal.NewDate(2022, 11, 6)

	err := dr.ValidateWithContext(context.Background())
	assert.NoError(t, err)
}

func TestDocumentNormalize(t *testing.T) {
	dr := &org.DocumentRef{
		Code: " Foo ",
		Ext: tax.Extensions{
			"fooo": "",
		},
	}
	dr.Normalize(nil)
	assert.Equal(t, "FOO", dr.Code.String())
	assert.Empty(t, dr.Ext)
}
