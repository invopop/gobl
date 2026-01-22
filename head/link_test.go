package head_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinkValidation(t *testing.T) {
	t.Run("simple link", func(t *testing.T) {
		l := &head.Link{
			Key:   "test",
			Title: "Test Link Title",
			URL:   "https://example.com",
		}
		assert.NoError(t, l.Validate())
	})
	t.Run("invalid link", func(t *testing.T) {
		l := &head.Link{
			Key:   "test",
			Title: "Test Link Title",
			URL:   "example",
		}
		require.ErrorContains(t, l.Validate(), "url: must be a valid URL")
	})

	t.Run("missing url", func(t *testing.T) {
		l := &head.Link{
			Key:   "test",
			Title: "Test Link Title",
		}
		require.ErrorContains(t, l.Validate(), "url: cannot be blank")
	})

	t.Run("missing key", func(t *testing.T) {
		l := &head.Link{
			Title: "Test Link Title",
			URL:   "https://example.com",
		}
		require.ErrorContains(t, l.Validate(), "key: cannot be blank")
	})

	t.Run("valid MIME types", func(t *testing.T) {
		validMIMEs := []string{
			"application/pdf",
			"image/jpeg",
			"image/png",
			"text/csv",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"application/vnd.oasis.opendocument.spreadsheet",
			"text/html",
			"application/xml",
			"text/xml",
			"application/json",
		}
		for _, mime := range validMIMEs {
			l := &head.Link{
				Key:  "test",
				MIME: mime,
				URL:  "https://example.com",
			}
			assert.NoError(t, l.Validate(), "MIME type %s should be valid", mime)
		}
	})

	t.Run("invalid MIME types", func(t *testing.T) {
		invalidMIMEs := []string{
			"application/octet-stream",
			"text/plain",
			"image/gif",
			"video/mp4",
			"application/zip",
		}
		for _, mime := range invalidMIMEs {
			l := &head.Link{
				Key:  "test",
				MIME: mime,
				URL:  "https://example.com",
			}
			err := l.Validate()
			require.Error(t, err, "MIME type %s should be invalid", mime)
			require.ErrorContains(t, err, "mime:", "Error should mention mime field")
		}
	})

	t.Run("empty MIME type is valid", func(t *testing.T) {
		l := &head.Link{
			Key:  "test",
			MIME: "",
			URL:  "https://example.com",
		}
		assert.NoError(t, l.Validate())
	})
}

func TestLinkDigestValidation(t *testing.T) {
	t.Run("digest with valid MIME type", func(t *testing.T) {
		l := &head.Link{
			Key:  "test",
			MIME: "application/pdf",
			Digest: &dsig.Digest{
				Algorithm: dsig.DigestSHA256,
				Value:     "abc123",
			},
			URL: "https://example.com",
		}
		assert.NoError(t, l.Validate())
	})

	t.Run("digest without MIME type should fail", func(t *testing.T) {
		l := &head.Link{
			Key: "test",
			Digest: &dsig.Digest{
				Algorithm: dsig.DigestSHA256,
				Value:     "abc123",
			},
			URL: "https://example.com",
		}
		err := l.Validate()
		require.Error(t, err)
		require.ErrorContains(t, err, "must be nil when MIME type is not provided")
	})

	t.Run("no digest without MIME type is valid", func(t *testing.T) {
		l := &head.Link{
			Key: "test",
			URL: "https://example.com",
		}
		assert.NoError(t, l.Validate())
	})
}

func TestLinkByKey(t *testing.T) {
	t.Run("find link", func(t *testing.T) {
		l1 := &head.Link{Category: head.LinkCategoryKeyFormat, Key: "test1"}
		l2 := &head.Link{Category: head.LinkCategoryKeyFormat, Key: "test2"}
		l3 := &head.Link{Category: head.LinkCategoryKeyFormat, Key: "test3"}
		list := []*head.Link{l1, l2, l3}

		assert.Equal(t, l2, head.LinkByCategoryAndKey(list, head.LinkCategoryKeyFormat, "test2"))
		assert.Nil(t, head.LinkByCategoryAndKey(list, head.LinkCategoryKeyFormat, "test4"))
	})

	t.Run("find link with different categories", func(t *testing.T) {
		l1 := &head.Link{Category: head.LinkCategoryKeyFormat, Key: "test1"}
		l2 := &head.Link{Category: head.LinkCategoryKeyRequest, Key: "test1"}
		list := []*head.Link{l1, l2}

		assert.Equal(t, l1, head.LinkByCategoryAndKey(list, head.LinkCategoryKeyFormat, "test1"))
		assert.Equal(t, l2, head.LinkByCategoryAndKey(list, head.LinkCategoryKeyRequest, "test1"))
		assert.Nil(t, head.LinkByCategoryAndKey(list, head.LinkCategoryKeyResponse, "test1"))
	})
}

func TestAppendLink(t *testing.T) {
	t.Run("append link", func(t *testing.T) {
		l1 := &head.Link{Key: "test1", URL: "https://example.com/1"}
		l2 := &head.Link{Key: "test2", URL: "https://example.com/2"}
		l3 := &head.Link{Key: "test3", URL: "https://example.com/3"}
		list := []*head.Link{l1, l2}

		list = head.AppendLink(list, l3)
		assert.Len(t, list, 3)
		assert.Equal(t, l3, list[2])

		list = head.AppendLink(list, nil)
		assert.Len(t, list, 3)
		assert.Equal(t, l3, list[2])
	})

	t.Run("update link", func(t *testing.T) {
		l1 := &head.Link{Key: "test1", Title: "Old Title", URL: "https://example.com/1"}
		l2 := &head.Link{Key: "test1", Title: "New Title", URL: "https://example.com/2"}
		list := []*head.Link{l1}

		list = head.AppendLink(list, l2)
		assert.Len(t, list, 1)
		assert.Equal(t, l2, list[0])
	})
}

func TestDetectDuplicateLink(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		l1 := &head.Link{Key: "test1", URL: "https://example.com"}
		l2 := &head.Link{Key: "test2", URL: "https://example.com/2"}
		list := []*head.Link{l1, l2}

		err := validation.Validate(list, head.DetectDuplicateLinks)
		assert.NoError(t, err)
	})
	t.Run("detect duplicate", func(t *testing.T) {
		l1 := &head.Link{Key: "test1", URL: "https://example.com"}
		l2 := &head.Link{Key: "test1", URL: "https://example.com/2"}
		list := []*head.Link{l1, l2}

		err := validation.Validate(list, head.DetectDuplicateLinks)

		require.ErrorContains(t, err, "duplicate key 'test1'")
	})

	t.Run("detect duplicate in category", func(t *testing.T) {
		l1 := &head.Link{Category: head.LinkCategoryKeyFormat, Key: "test1", URL: "https://example.com"}
		l2 := &head.Link{Category: head.LinkCategoryKeyFormat, Key: "test1", URL: "https://example.com/2"}
		list := []*head.Link{l1, l2}

		err := validation.Validate(list, head.DetectDuplicateLinks)

		require.ErrorContains(t, err, "duplicate category 'format' and key 'test1'")
	})
}

func TestLinkExtendJSONSchemas(t *testing.T) {
	base := here.Doc(`
		{
			"properties": {
				"category": {
					"$ref": "https://gobl.org/draft-0/cbc/key",
					"title": "Category"
				}
			}
		}
	`)
	js := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(base), js))
	head.Link{}.JSONSchemaExtend(js)

	prop, ok := js.Properties.Get("category")
	assert.True(t, ok)
	assert.Len(t, prop.OneOf, 7)
	assert.Equal(t, head.LinkCategoryKeyFormat, prop.OneOf[0].Const)
	assert.Equal(t, "Format", prop.OneOf[0].Title)
	assert.Equal(t, head.LinkCategoryKeyPortal, prop.OneOf[1].Const)
	assert.Equal(t, "Portal", prop.OneOf[1].Title)
}
