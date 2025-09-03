package org_test

import (
	"testing"

	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestAttachmentNormalize(t *testing.T) {
	t.Run("trims and normalizes fields", func(t *testing.T) {
		a := &org.Attachment{
			Key:         " key ",
			Code:        " ABC ",
			Name:        "  test.txt  ",
			Description: "  some description \n",
			URL:         "  https://example.com/doc.pdf  ",
			MIME:        "  application/pdf  ",
			Data:        []byte(" keep "),
		}
		a.Normalize(nil)

		assert.Equal(t, "ABC", a.Code.String())
		assert.Equal(t, "test.txt", a.Name)
		assert.Equal(t, "some description", a.Description)
		assert.Equal(t, "https://example.com/doc.pdf", a.URL)
		assert.Equal(t, "application/pdf", a.MIME)
	})

	t.Run("nil receiver no panic", func(t *testing.T) {
		var a *org.Attachment
		assert.NotPanics(t, func() {
			a.Normalize(nil)
		})
	})

	t.Run("blank-only strings become empty", func(t *testing.T) {
		a := &org.Attachment{
			Name:        " name ",
			URL:         "   ",
			Description: "   ",
			MIME:        "   ",
		}
		a.Normalize(nil)
		assert.Equal(t, "name", a.Name)
		assert.Equal(t, "", a.URL)
		assert.Equal(t, "", a.Description)
		assert.Equal(t, "", a.MIME)
	})
}

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
