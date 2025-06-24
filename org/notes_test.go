package org_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoteNormalize(t *testing.T) {
	t.Run("accepts nil", func(t *testing.T) {
		var n *org.Note
		assert.NotPanics(t, func() {
			n.Normalize(nil)
		})
	})

	t.Run("accepts empty", func(t *testing.T) {
		n := &org.Note{}
		assert.NotPanics(t, func() {
			n.Normalize(nil)
		})
	})

	t.Run("converts valid", func(t *testing.T) {
		n := &org.Note{
			Key:  org.NoteKeyGeneral,
			Text: "This is a general note test",
			Ext: tax.Extensions{
				"untidid-text-subject": "AAI",
			},
		}
		n.Normalize(nil)
		assert.Equal(t, "AAI", n.Ext.Get("untidid-text-subject").String())
	})

	t.Run("cleans extensions", func(t *testing.T) {
		n := &org.Note{
			Code: " FOO ",
			Text: "This is a general note test",
			Ext: tax.Extensions{
				"missing": "",
			},
		}
		n.Normalize(nil)
		assert.Equal(t, "FOO", n.Code.String())
		assert.Empty(t, n.Ext)
	})
}

func TestNotesValidation(t *testing.T) {
	n := new(org.Note)
	n.Text = "This is a general note test"

	err := n.Validate()
	assert.NoError(t, err) // empty key ok

	n.Key = org.NoteKeyGeneral
	err = n.Validate()
	assert.NoError(t, err)

	n.Key = cbc.Key("fooo")
	err = n.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "key: must be a valid value")
}

func TestNoteFromScenario(t *testing.T) {
	n := org.NoteFromScenario(&tax.ScenarioNote{
		Key:  org.NoteKeyGeneral,
		Code: "note1",
		Src:  "src1",
		Text: "This is a note1",
		Ext: tax.Extensions{
			"untidid-text-subject": "AAI",
		},
	})
	assert.Equal(t, org.NoteKeyGeneral, n.Key)
	assert.Equal(t, "note1", n.Code.String())
	assert.Equal(t, "src1", n.Src.String())
	assert.Equal(t, "This is a note1", n.Text)
	assert.Equal(t, "AAI", n.Ext.Get("untidid-text-subject").String())

	assert.NotPanics(t, func() {
		org.NoteFromScenario(nil)
	})
}

func TestNotesValidationHasKey(t *testing.T) {
	ns := []*org.Note{
		{
			Key:  org.NoteKeyGeneral,
			Text: "This is a general note test",
		},
	}
	err := validation.Validate(ns, org.ValidateNotesHasKey(org.NoteKeyGeneral))
	assert.NoError(t, err)

	err = validation.Validate(ns, org.ValidateNotesHasKey(org.NoteKeyLegal))
	assert.ErrorContains(t, err, "with key 'legal' missin")

	err = validation.Validate(cbc.Key("foo"), org.ValidateNotesHasKey(org.NoteKeyGeneral))
	assert.Nil(t, err, "should ignore invalid type")
}

func TestNoteSameAs(t *testing.T) {
	n := &org.Note{
		Key:  org.NoteKeyGeneral,
		Code: "ABC",
		Text: "This is a test",
	}
	assert.True(t, n.SameAs(&org.Note{
		Key:  org.NoteKeyGeneral,
		Code: "ABC",
		Text: "This is a test ABC",
	}))
	assert.False(t, n.SameAs(&org.Note{
		Key:  org.NoteKeyGeneral,
		Code: "123",
		Text: "This is a test 123",
	}))
	assert.False(t, n.SameAs(&org.Note{
		Key:  org.NoteKeyLegal,
		Code: "ABC",
		Text: "This is a test ABC",
	}))
	t.Run("nils", func(t *testing.T) {
		var n1, n2 *org.Note
		assert.False(t, n1.SameAs(n2), "nil should not match nil")
	})
}

func TestNoteEquals(t *testing.T) {
	n := &org.Note{
		Key:  org.NoteKeyGeneral,
		Code: "ABC",
		Text: "This is a test",
	}
	assert.True(t, n.Equals(&org.Note{
		Key:  org.NoteKeyGeneral,
		Code: "ABC",
		Text: "This is a test",
	}))
	assert.False(t, n.Equals(&org.Note{
		Key:  org.NoteKeyGeneral,
		Code: "123",
		Text: "This is a test",
	}))
	assert.False(t, n.Equals(&org.Note{
		Key:  org.NoteKeyGeneral,
		Code: "ABC",
		Src:  "fooo",
		Text: "This is a test",
	}))
	assert.False(t, n.Equals(&org.Note{
		Key:  org.NoteKeyLegal,
		Code: "ABC",
		Text: "This is a test",
	}))
}

func TestNoteWithSrc(t *testing.T) {
	n := &org.Note{
		Key:  org.NoteKeyGeneral,
		Code: "ABC",
		Text: "This is a test",
	}
	n2 := n.WithSrc("foo")
	assert.Empty(t, n.Src)
	assert.Equal(t, "foo", n2.Src.String())
}

func TestNoteWithCode(t *testing.T) {
	n := &org.Note{
		Key:  org.NoteKeyGeneral,
		Text: "This is a test",
	}
	n2 := n.WithCode("foo")
	assert.Empty(t, n.Code)
	assert.Equal(t, "foo", n2.Code.String())
}

func TestNoteJSONSchemaExtend(t *testing.T) {
	ks := new(jsonschema.Schema)
	require.NoError(t, json.Unmarshal([]byte(`{"properties":{"key":{}}}`), ks))
	n := new(org.Note)
	n.JSONSchemaExtend(ks)
	kp, _ := ks.Properties.Get("key")
	require.NotNil(t, kp)
	assert.Greater(t, len(kp.OneOf), 1)
	first := kp.OneOf[0]
	assert.Equal(t, "goods", first.Const)
	assert.Equal(t, "Goods", first.Title)
	assert.Equal(t, "Goods Description", first.Description)
}
