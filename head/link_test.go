package head_test

import (
	"testing"

	"github.com/invopop/gobl/head"
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
}

func TestLinkByKey(t *testing.T) {
	t.Run("find link", func(t *testing.T) {
		l1 := &head.Link{Key: "test1"}
		l2 := &head.Link{Key: "test2"}
		l3 := &head.Link{Key: "test3"}
		list := []*head.Link{l1, l2, l3}

		assert.Equal(t, l2, head.LinkByKey(list, "test2"))
		assert.Nil(t, head.LinkByKey(list, "test4"))
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
}
