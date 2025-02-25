package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestAttachmentValidation(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		a := &org.Attachment{
			Key:  "key",
			Code: "ABC",
			Name: "test.txt",
			URL:  "https://example.com/test.txt",
		}
		err := a.Validate()
		assert.NoError(t, err)
	})

	t.Run("both URL and data", func(t *testing.T) {
		a := &org.Attachment{
			Key:  "key",
			Code: "ABC",
			Name: "test.txt",
			URL:  "https://example.com/test.txt",
			Data: []byte("test"),
		}
		err := a.Validate()
		assert.ErrorContains(t, err, "data: must be blank with url")
	})

	t.Run("missing URL and data", func(t *testing.T) {
		a := &org.Attachment{
			Key:  "key",
			Code: "ABC",
			Name: "test.txt",
		}
		err := a.Validate()
		assert.ErrorContains(t, err, "url: cannot be blank")
	})
}
