package tax_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagSetForSchema(t *testing.T) {
	ts1 := &tax.TagSet{
		Schema: "bill/invoice",
		List: []*cbc.Definition{
			{
				Key: "test1",
				Name: i18n.String{
					i18n.EN: "Test 1",
				},
			},
		},
	}
	ts2 := &tax.TagSet{
		Schema: "bill/receipt",
		List: []*cbc.Definition{
			{
				Key: "test2",
				Name: i18n.String{
					i18n.EN: "Test 2",
				},
			},
		},
	}
	ts3 := tax.TagSetForSchema([]*tax.TagSet{ts1, ts2}, "bill/invoice")
	assert.Equal(t, ts3, ts1)
	ts4 := tax.TagSetForSchema([]*tax.TagSet{ts1, ts2}, "bill/receipt")
	assert.Equal(t, ts4, ts2)
	ts5 := tax.TagSetForSchema([]*tax.TagSet{ts1, ts2}, "bill/unknown")
	assert.Nil(t, ts5)
}

func TestTagsHasTags(t *testing.T) {
	var tagsNil tax.Tags
	assert.False(t, tagsNil.HasTags("tag1"))

	tags := tax.WithTags("tag1", "tag2")
	assert.True(t, tags.HasTags("tag1"))
	assert.True(t, tags.HasTags("tag2"))
	assert.False(t, tags.HasTags("tag3"))
}

func TestTagSetMerge(t *testing.T) {
	ts1 := &tax.TagSet{
		Schema: "bill/invoice",
		List: []*cbc.Definition{
			{
				Key: "test1",
				Name: i18n.String{
					i18n.EN: "Test 1",
				},
			},
		},
	}
	ts2 := &tax.TagSet{
		Schema: "bill/invoice",
		List: []*cbc.Definition{
			{
				Key: "test1",
				Name: i18n.String{
					i18n.EN: "Test 1 duplicate",
				},
			},
			{
				Key: "test2",
				Name: i18n.String{
					i18n.EN: "Test 2",
				},
			},
		},
	}
	ts3 := ts1.Merge(ts2)
	assert.Len(t, ts1.List, 1, "should not touch original")
	assert.Len(t, ts2.List, 2, "should not touch original")
	assert.Len(t, ts3.List, 2)
	assert.Equal(t, ts3.List[0].Key.String(), "test1")
	assert.Equal(t, ts3.List[0].Name[i18n.EN], "Test 1")
	assert.Equal(t, ts3.List[1].Key.String(), "test2")
	assert.Equal(t, ts3.List[1].Name[i18n.EN], "Test 2")

	ts4 := &tax.TagSet{
		Schema: "note/message",
		List: []*cbc.Definition{
			{
				Key: "test3",
				Name: i18n.String{
					i18n.EN: "Test 3 Note",
				},
			},
		},
	}
	ts5 := ts1.Merge(ts4)
	assert.Equal(t, ts5, ts1)

	var ts6 *tax.TagSet
	ts7 := ts1.Merge(ts6)
	assert.Equal(t, ts7, ts1)
	ts7 = ts6.Merge(ts1)
	assert.Equal(t, ts7, ts1)
}

func TestTagSetValidation(t *testing.T) {
	list := []cbc.Key{"tag1", "tag2"}

	t.Run("valid tag", func(t *testing.T) {
		err := validation.Validate([]cbc.Key{"tag1"}, tax.TagsIn(list...))
		assert.NoError(t, err)
	})
	t.Run("invalid tag", func(t *testing.T) {
		err := validation.Validate([]cbc.Key{"tag3"}, tax.TagsIn(list...))
		assert.ErrorContains(t, err, "0: 'tag3' undefined")
	})
	t.Run("invalid tag 2", func(t *testing.T) {
		err := validation.Validate([]cbc.Key{"tag1", "tag3"}, tax.TagsIn(list...))
		assert.ErrorContains(t, err, "1: 'tag3' undefined")
	})

	t.Run("with tags", func(t *testing.T) {
		set := tax.WithTags("tag1")
		err := validation.Validate(set, tax.TagsIn(list...))
		assert.NoError(t, err)
	})

	t.Run("with something else", func(t *testing.T) {
		codes := []cbc.Code{"FOO"}
		err := validation.Validate(codes, tax.TagsIn(list...))
		assert.NoError(t, err)
	})
}

func TestTaxSetKeys(t *testing.T) {
	ts1 := &tax.TagSet{
		Schema: "bill/invoice",
		List: []*cbc.Definition{
			{
				Key: "test1",
				Name: i18n.String{
					i18n.EN: "Test 1 duplicate",
				},
			},
			{
				Key: "test2",
				Name: i18n.String{
					i18n.EN: "Test 2",
				},
			},
		},
	}
	assert.Equal(t, ts1.Keys(), []cbc.Key{"test1", "test2"})
}

func TestTagsJSONSchemaEmbedWithDefs(t *testing.T) {
	eg := `{
		"properties": {
			"$tags": {
				"items": {
            		"$ref": "https://gobl.org/draft-0/cbc/key",
					"type": "array",
					"title": "Tags"
				}
			}
		}
	}`
	js := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(eg), js))

	ts := tax.Tags{}
	ts.JSONSchemaExtendWithDefs(js, bill.DefaultInvoiceTags().List)

	data, _ := json.Marshal(js)
	t.Logf("JSON Schema: %s", data)

	prop, ok := js.Properties.Get("$tags")
	require.True(t, ok)
	assert.Equal(t, 6, len(prop.Items.OneOf), "should have 5 tags plus 1 catch-all")
	assert.Equal(t, "simplified", prop.Items.OneOf[0].Const)
	assert.Equal(t, "Any", prop.Items.OneOf[5].Title)
}
