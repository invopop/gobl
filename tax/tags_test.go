package tax_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestTagSetMerge(t *testing.T) {
	ts1 := &tax.TagSet{
		Schema: "bill/invoice",
		List: []*cbc.KeyDefinition{
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
		List: []*cbc.KeyDefinition{
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
		List: []*cbc.KeyDefinition{
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

func TestTaxSetKeys(t *testing.T) {
	ts1 := &tax.TagSet{
		Schema: "bill/invoice",
		List: []*cbc.KeyDefinition{
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
