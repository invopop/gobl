package tax_test

import (
	"encoding/json"
	"testing"

	// this will also prepare registers
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCleanExtensions(t *testing.T) {
	var em tax.Extensions

	em2 := em.Clean()
	assert.True(t, em2.IsZero())

	em = tax.ExtensionsOf(tax.ExtMap{
		"key": "",
	})
	em2 = em.Clean()
	assert.True(t, em2.IsZero())

	em = tax.ExtensionsOf(tax.ExtMap{
		"key": "foo",
		"bar": "",
	})
	em2 = em.Clean()
	assert.False(t, em2.IsZero())
	assert.Equal(t, 1, em2.Len())
	assert.Equal(t, "foo", em2.Get("key").String())
}

func TestExtensionsRequiresValidation(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := tax.ExtensionsRequire(untdid.ExtKeyDocumentType).Validate(nil)
		assert.NoError(t, err)
	})
	t.Run("empty", func(t *testing.T) {
		em := tax.Extensions{}
		err := tax.ExtensionsRequire(untdid.ExtKeyDocumentType).Validate(em)
		assert.ErrorContains(t, err, "untdid-document-type: required")
	})
	t.Run("correct", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			untdid.ExtKeyDocumentType: "326",
		})
		err := tax.ExtensionsRequire(untdid.ExtKeyDocumentType).Validate(em)
		assert.NoError(t, err)
	})
	t.Run("correct with extras", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			untdid.ExtKeyDocumentType: "326",
			iso.ExtKeySchemeID:        "1234",
		})
		err := tax.ExtensionsRequire(untdid.ExtKeyDocumentType).Validate(em)
		assert.NoError(t, err)
	})
	t.Run("missing", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			iso.ExtKeySchemeID: "1234",
		})
		err := tax.ExtensionsRequire(untdid.ExtKeyDocumentType).Validate(em)
		assert.ErrorContains(t, err, "untdid-document-type: required")
	})
}

func TestExtensionsAllOrNoneValidation(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := tax.ExtensionsRequireAllOrNone(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID).Validate(nil)
		assert.NoError(t, err)
	})
	t.Run("empty", func(t *testing.T) {
		em := tax.Extensions{}
		err := tax.ExtensionsRequireAllOrNone(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID).Validate(em)
		assert.NoError(t, err)
	})
	t.Run("all present", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			untdid.ExtKeyDocumentType: "326",
			iso.ExtKeySchemeID:        "1234",
		})
		err := tax.ExtensionsRequireAllOrNone(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID).Validate(em)
		assert.NoError(t, err)
	})
	t.Run("none present", func(t *testing.T) {
		em := tax.Extensions{}
		err := tax.ExtensionsRequireAllOrNone(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID).Validate(em)
		assert.NoError(t, err)
	})
	t.Run("some present", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			untdid.ExtKeyDocumentType: "326",
		})
		err := tax.ExtensionsRequireAllOrNone(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID).Validate(em)
		assert.ErrorContains(t, err, "iso-scheme-id: required")
	})
	t.Run("some present reversed", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			iso.ExtKeySchemeID: "1234",
		})
		err := tax.ExtensionsRequireAllOrNone(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID).Validate(em)
		assert.ErrorContains(t, err, "untdid-document-type: required")
	})
}

func TestExtensionsExcludeValidation(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := tax.ExtensionsExclude(untdid.ExtKeyDocumentType).Validate(nil)
		assert.NoError(t, err)
	})
	t.Run("empty", func(t *testing.T) {
		em := tax.Extensions{}
		err := tax.ExtensionsExclude(untdid.ExtKeyDocumentType).Validate(em)
		assert.NoError(t, err)
	})
	t.Run("correct", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			untdid.ExtKeyDocumentType: "326",
		})
		err := tax.ExtensionsExclude(untdid.ExtKeyDocumentType).Validate(em)
		assert.ErrorContains(t, err, "untdid-document-type: must be blank")
	})
	t.Run("correct with extras", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			untdid.ExtKeyDocumentType: "326",
			iso.ExtKeySchemeID:        "1234",
		})
		err := tax.ExtensionsExclude(untdid.ExtKeyCharge).Validate(em)
		assert.NoError(t, err)
	})
}

func TestExtensionsAllowOneOfValidation(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := tax.ExtensionsAllowOneOf(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID).Validate(nil)
		assert.NoError(t, err)
	})
	t.Run("empty", func(t *testing.T) {
		em := tax.Extensions{}
		err := tax.ExtensionsAllowOneOf(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID).Validate(em)
		assert.NoError(t, err)
	})
	t.Run("one present", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			untdid.ExtKeyDocumentType: "326",
		})
		err := tax.ExtensionsAllowOneOf(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID).Validate(em)
		assert.NoError(t, err)
	})
	t.Run("both present", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			untdid.ExtKeyDocumentType: "326",
			iso.ExtKeySchemeID:        "1234",
		})
		err := tax.ExtensionsAllowOneOf(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID).Validate(em)
		assert.ErrorContains(t, err, "untdid-document-type: only one allowed")
		assert.ErrorContains(t, err, "iso-scheme-id: only one allowed")
	})
}

func TestExtensionsHasValues(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, "326", "389").Validate(nil)
		assert.NoError(t, err)
	})
	t.Run("empty", func(t *testing.T) {
		em := tax.Extensions{}
		err := tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, "326", "389").Validate(em)
		assert.NoError(t, err)
	})
	t.Run("different extensions", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			iso.ExtKeySchemeID: "1234",
		})
		err := tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, "326", "389").Validate(em)
		assert.NoError(t, err)
	})
	t.Run("has codes", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			untdid.ExtKeyDocumentType: "326",
		})
		err := tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, "326", "389").Validate(em)
		assert.NoError(t, err)
	})
	t.Run("invalid code", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			untdid.ExtKeyDocumentType: "102",
		})
		err := tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, "326", "389").Validate(em)
		assert.ErrorContains(t, err, "untdid-document-type: invalid value")
	})
}

func TestExtensionsExcludeCodes(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, "380", "381").Validate(nil)
		assert.NoError(t, err)
	})
	t.Run("empty", func(t *testing.T) {
		em := tax.Extensions{}
		err := tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, "380", "381").Validate(em)
		assert.NoError(t, err)
	})
	t.Run("different extensions", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			iso.ExtKeySchemeID: "1234",
		})
		err := tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, "380", "381").Validate(em)
		assert.NoError(t, err)
	})
	t.Run("allowed code", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			untdid.ExtKeyDocumentType: "326",
		})
		err := tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, "380", "381").Validate(em)
		assert.NoError(t, err)
	})
	t.Run("excluded code", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			untdid.ExtKeyDocumentType: "380",
		})
		err := tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, "380", "381").Validate(em)
		assert.ErrorContains(t, err, "untdid-document-type: value '380' not allowed")
	})
	t.Run("another excluded code", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			untdid.ExtKeyDocumentType: "381",
		})
		err := tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, "380", "381").Validate(em)
		assert.ErrorContains(t, err, "untdid-document-type: value '381' not allowed")
	})
}

func TestExtensionsHas(t *testing.T) {
	em := tax.ExtensionsOf(tax.ExtMap{
		"key": "value",
	})
	assert.True(t, em.Has("key"))
	assert.False(t, em.Has("invalid"))
}

func TestExtensionsValues(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var em tax.Extensions
		assert.Empty(t, em.Values())
	})
	t.Run("with values", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			"key1": "value1",
			"key2": "value2",
		})
		assert.Equal(t, []cbc.Code{"value1", "value2"}, em.Values())
	})
	t.Run("sorted output", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			"a": "cherry",
			"b": "apple",
			"c": "banana",
		})
		assert.Equal(t, []cbc.Code{"apple", "banana", "cherry"}, em.Values())
	})
}

func TestExtensionsEquals(t *testing.T) {
	tests := []struct {
		name string
		em1  tax.Extensions
		em2  tax.Extensions
		want bool
	}{
		{
			name: "empty",
			em1:  tax.Extensions{},
			em2:  tax.Extensions{},
			want: true,
		},
		{
			name: "same",
			em1:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			em2:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			want: true,
		},
		{
			name: "different",
			em1:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			em2:  tax.ExtensionsOf(tax.ExtMap{"key": "value2"}),
			want: false,
		},
		{
			name: "different keys",
			em1:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			em2:  tax.ExtensionsOf(tax.ExtMap{"key2": "value"}),
			want: false,
		},
		{
			name: "different lengths",
			em1:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			em2:  tax.ExtensionsOf(tax.ExtMap{"key": "value", "key2": "value"}),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.em1.Equals(tt.em2))
		})
	}
}

func TestExtensionsContains(t *testing.T) {
	tests := []struct {
		name string
		em1  tax.Extensions
		em2  tax.Extensions
		want bool
	}{
		{
			name: "empty",
			em1:  tax.Extensions{},
			em2:  tax.Extensions{},
			want: false,
		},
		{
			name: "same",
			em1:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			em2:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			want: true,
		},
		{
			name: "different",
			em1:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			em2:  tax.ExtensionsOf(tax.ExtMap{"key": "value2"}),
			want: false,
		},
		{
			name: "different keys",
			em1:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			em2:  tax.ExtensionsOf(tax.ExtMap{"key2": "value"}),
			want: false,
		},
		{
			name: "different lengths",
			em1:  tax.ExtensionsOf(tax.ExtMap{"key": "value", "key2": "value"}),
			em2:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.em1.Contains(tt.em2))
		})
	}
}

func TestExtensionsMerge(t *testing.T) {
	tests := []struct {
		name string
		em1  tax.Extensions
		em2  tax.Extensions
		want tax.Extensions
	}{
		{
			name: "empty",
			em1:  tax.Extensions{},
			em2:  tax.Extensions{},
			want: tax.Extensions{},
		},
		{
			name: "same",
			em1:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			em2:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			want: tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
		},
		{
			name: "zero source",
			em1:  tax.Extensions{},
			em2:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			want: tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
		},
		{
			name: "zero destination",
			em1:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			em2:  tax.Extensions{},
			want: tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
		},
		{
			name: "different",
			em1:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			em2:  tax.ExtensionsOf(tax.ExtMap{"key": "value2"}),
			want: tax.ExtensionsOf(tax.ExtMap{"key": "value2"}),
		},
		{
			name: "different keys",
			em1:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			em2:  tax.ExtensionsOf(tax.ExtMap{"key2": "value"}),
			want: tax.ExtensionsOf(tax.ExtMap{"key": "value", "key2": "value"}),
		},
		{
			name: "different lengths",
			em1:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			em2:  tax.ExtensionsOf(tax.ExtMap{"key": "value", "key2": "value"}),
			want: tax.ExtensionsOf(tax.ExtMap{"key": "value", "key2": "value"}),
		},
		{
			name: "different lengths 2",
			em1:  tax.ExtensionsOf(tax.ExtMap{"key": "value", "key2": "value"}),
			em2:  tax.ExtensionsOf(tax.ExtMap{"key": "value"}),
			want: tax.ExtensionsOf(tax.ExtMap{"key": "value", "key2": "value"}),
		},
		{
			name: "different lengths 3",
			em1:  tax.ExtensionsOf(tax.ExtMap{"key": "value2"}),
			em2:  tax.ExtensionsOf(tax.ExtMap{"key": "value", "key2": "value"}),
			want: tax.ExtensionsOf(tax.ExtMap{"key": "value", "key2": "value"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.em1.Merge(tt.em2))
		})
	}
}

func TestExtensionLookup(t *testing.T) {
	em := tax.ExtensionsOf(tax.ExtMap{
		"key1": "foo",
		"key2": "bar",
	})
	assert.Equal(t, cbc.Key("key1"), em.Lookup("foo"))
	assert.Equal(t, cbc.Key("key2"), em.Lookup("bar"))
	assert.Equal(t, cbc.KeyEmpty, em.Lookup("missing"))
}

func TestExtensionGet(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var em tax.Extensions
		assert.Equal(t, "", em.Get("key").String())
	})
	t.Run("with value", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			"key": "value",
		})
		assert.Equal(t, "value", em.Get("key").String())
	})
	t.Run("missing", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			"key": "value",
		})
		assert.Equal(t, "", em.Get("missing").String())
	})
	t.Run("with sub-keys", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{
			"key": "value",
		})
		assert.Equal(t, "value", em.Get("key+foo").String())
	})
}

func TestExtensionsSet(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var em tax.Extensions
		em = em.Set("key", "value")
		assert.Equal(t, tax.ExtensionsOf(tax.ExtMap{"key": "value"}), em)
	})

	t.Run("immutable: discarded result does not mutate", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{"key": "value1"})
		em.Set("key", "value2") // result discarded
		assert.Equal(t, tax.ExtensionsOf(tax.ExtMap{"key": "value1"}), em)
	})

	t.Run("with new value", func(t *testing.T) {
		em := tax.Extensions{}
		em = em.Set("key", "value1")
		assert.Equal(t, tax.ExtensionsOf(tax.ExtMap{"key": "value1"}), em)
	})

	t.Run("with empty value removes key", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{"key": "value1"})
		em = em.Set("key", "")
		assert.True(t, em.IsZero())
	})
}

func TestExtensionsSetIfEmpty(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var em tax.Extensions
		em = em.SetIfEmpty("key", "value")
		assert.Equal(t, tax.ExtensionsOf(tax.ExtMap{"key": "value"}), em)
	})

	t.Run("with existing value stays the same", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{"key": "value1"})
		em = em.SetIfEmpty("key", "value2")
		assert.Equal(t, tax.ExtensionsOf(tax.ExtMap{"key": "value1"}), em)
	})

	t.Run("with new value", func(t *testing.T) {
		em := tax.Extensions{}
		em = em.SetIfEmpty("key", "value1")
		assert.Equal(t, tax.ExtensionsOf(tax.ExtMap{"key": "value1"}), em)
	})
}

func TestExtensionsSetOneOf(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var em tax.Extensions
		em = em.SetOneOf("key", "value1", "value2")
		assert.Equal(t, tax.ExtensionsOf(tax.ExtMap{"key": "value1"}), em)
	})

	t.Run("with existing primary value and no alternatives replaces it", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{"key": "value1"})
		em = em.SetOneOf("key", "value2", "value3")
		assert.Equal(t, tax.ExtensionsOf(tax.ExtMap{"key": "value2"}), em)
	})

	t.Run("with existing alternative value keeps it", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{"key": "value3"})
		em = em.SetOneOf("key", "value2", "value3")
		assert.Equal(t, tax.ExtensionsOf(tax.ExtMap{"key": "value3"}), em)
	})

	t.Run("with no existing value", func(t *testing.T) {
		em := tax.Extensions{}
		em = em.SetOneOf("key", "value1", "value2")
		assert.Equal(t, tax.ExtensionsOf(tax.ExtMap{"key": "value1"}), em)
	})
}

func TestExtensionsDelete(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var em tax.Extensions
		em = em.Delete("key")
		assert.True(t, em.IsZero())
	})

	t.Run("with value", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{"key": "value"})
		em = em.Delete("key")
		assert.True(t, em.IsZero())
	})

	t.Run("with missing value", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{"key": "value"})
		em = em.Delete("missing")
		assert.Equal(t, tax.ExtensionsOf(tax.ExtMap{"key": "value"}), em)
	})
}

func init() {
	tax.RegisterExtension(&cbc.Definition{
		Key:  "test-regime-color",
		Name: i18n.String{"en": "Color"},
		Values: []*cbc.Definition{
			{Code: "red", Name: i18n.String{"en": "Red"}},
			{Code: "green", Name: i18n.String{"en": "Green"}},
			{Code: "blue", Name: i18n.String{"en": "Blue"}},
		},
	})
	tax.RegisterExtension(&cbc.Definition{
		Key:     "test-regime-postal-code",
		Name:    i18n.String{"en": "Postal Code"},
		Pattern: `^\d{5}$`,
	})
	tax.RegisterExtension(&cbc.Definition{
		Key:  "test-regime-bare",
		Name: i18n.String{"en": "Bare"},
	})
}

func TestExtensionHasValidCode(t *testing.T) {
	t.Run("panic on unknown key", func(t *testing.T) {
		assert.Panics(t, func() {
			tax.ExtensionHasValidCode("test-regime-unknown")
		})
	})

	t.Run("panic on definition with no values or pattern", func(t *testing.T) {
		assert.Panics(t, func() {
			tax.ExtensionHasValidCode("test-regime-bare")
		})
	})

	t.Run("values-based: key absent passes", func(t *testing.T) {
		rule := tax.ExtensionHasValidCode("test-regime-color")
		em := tax.Extensions{}
		assert.True(t, rule.Check(em))
	})

	t.Run("values-based: nil extensions fails type check", func(t *testing.T) {
		rule := tax.ExtensionHasValidCode("test-regime-color")
		assert.False(t, rule.Check(nil))
	})

	t.Run("values-based: non-extensions value fails", func(t *testing.T) {
		rule := tax.ExtensionHasValidCode("test-regime-color")
		assert.False(t, rule.Check("not an extensions map"))
	})

	t.Run("values-based: valid code passes", func(t *testing.T) {
		rule := tax.ExtensionHasValidCode("test-regime-color")
		em := tax.ExtensionsOf(tax.ExtMap{"test-regime-color": "red"})
		assert.True(t, rule.Check(em))
	})

	t.Run("values-based: another valid code passes", func(t *testing.T) {
		rule := tax.ExtensionHasValidCode("test-regime-color")
		em := tax.ExtensionsOf(tax.ExtMap{"test-regime-color": "blue"})
		assert.True(t, rule.Check(em))
	})

	t.Run("values-based: invalid code fails", func(t *testing.T) {
		rule := tax.ExtensionHasValidCode("test-regime-color")
		em := tax.ExtensionsOf(tax.ExtMap{"test-regime-color": "purple"})
		assert.False(t, rule.Check(em))
	})

	t.Run("values-based: other keys not checked", func(t *testing.T) {
		rule := tax.ExtensionHasValidCode("test-regime-color")
		em := tax.ExtensionsOf(tax.ExtMap{
			"test-regime-postal-code": "not-a-number",
			"test-regime-color":       "green",
		})
		assert.True(t, rule.Check(em))
	})

	t.Run("values-based: string description", func(t *testing.T) {
		rule := tax.ExtensionHasValidCode("test-regime-color")
		assert.Contains(t, rule.String(), "test-regime-color")
		assert.Contains(t, rule.String(), "red")
		assert.Contains(t, rule.String(), "green")
		assert.Contains(t, rule.String(), "blue")
	})

	t.Run("pattern-based: key absent passes", func(t *testing.T) {
		rule := tax.ExtensionHasValidCode("test-regime-postal-code")
		em := tax.Extensions{}
		assert.True(t, rule.Check(em))
	})

	t.Run("pattern-based: valid code passes", func(t *testing.T) {
		rule := tax.ExtensionHasValidCode("test-regime-postal-code")
		em := tax.ExtensionsOf(tax.ExtMap{"test-regime-postal-code": "12345"})
		assert.True(t, rule.Check(em))
	})

	t.Run("pattern-based: invalid code fails", func(t *testing.T) {
		rule := tax.ExtensionHasValidCode("test-regime-postal-code")
		em := tax.ExtensionsOf(tax.ExtMap{"test-regime-postal-code": "1234"})
		assert.False(t, rule.Check(em))
	})

	t.Run("pattern-based: non-matching value fails", func(t *testing.T) {
		rule := tax.ExtensionHasValidCode("test-regime-postal-code")
		em := tax.ExtensionsOf(tax.ExtMap{"test-regime-postal-code": "abcde"})
		assert.False(t, rule.Check(em))
	})

	t.Run("pattern-based: string description", func(t *testing.T) {
		rule := tax.ExtensionHasValidCode("test-regime-postal-code")
		assert.Contains(t, rule.String(), "test-regime-postal-code")
		assert.Contains(t, rule.String(), `^\d{5}$`)
	})
}

func TestJSONSchemaExtend(t *testing.T) {
	in := here.Doc(`
		{
			"additionalProperties": {
				"^(?:[a-z]|[a-z0-9][a-z0-9-+]*[a-z0-9])$": {
					"$ref": "https://gobl.org/draft-0/cbc/code"
				}
			},
			"type": "object",
			"description": "Extensions is a map of extension keys to values."
		}
	`)
	schema := new(jsonschema.Schema)
	assert.NoError(t, json.Unmarshal([]byte(in), schema))
	assert.NotNil(t, schema.AdditionalProperties)
	var es tax.Extensions
	es.JSONSchemaExtend(schema)
	assert.Nil(t, schema.AdditionalProperties)
	assert.NotEmpty(t, schema.PatternProperties)
}

func TestExtensionsOfEmptyReturnsZero(t *testing.T) {
	assert.True(t, tax.ExtensionsOf(nil).IsZero())
	assert.True(t, tax.ExtensionsOf(tax.ExtMap{}).IsZero())
}

func TestExtensionsClone(t *testing.T) {
	t.Run("zero value clones to zero", func(t *testing.T) {
		var em tax.Extensions
		c := em.Clone()
		assert.True(t, c.IsZero())
	})
	t.Run("clone is independent", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{"a": "1", "b": "2"})
		c := em.Clone()
		assert.True(t, c.Equals(em))
		// Mutating the clone via Set must not affect the original.
		c = c.Set("a", "changed")
		assert.Equal(t, "1", em.Get("a").String())
		assert.Equal(t, "changed", c.Get("a").String())
	})
	t.Run("clone preserves all entries", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{"a": "1", "b": "2", "c": "3"})
		c := em.Clone()
		assert.Equal(t, 3, c.Len())
		assert.Equal(t, em.Keys(), c.Keys())
	})
}

func TestExtensionsAllIterator(t *testing.T) {
	t.Run("iterates in alphabetical key order", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{"c": "cherry", "a": "apple", "b": "banana"})
		var keys []cbc.Key
		var vals []cbc.Code
		for k, v := range em.All() {
			keys = append(keys, k)
			vals = append(vals, v)
		}
		assert.Equal(t, []cbc.Key{"a", "b", "c"}, keys)
		assert.Equal(t, []cbc.Code{"apple", "banana", "cherry"}, vals)
	})
	t.Run("early break stops iteration", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{"a": "1", "b": "2", "c": "3"})
		var count int
		for range em.All() {
			count++
			if count == 2 {
				break
			}
		}
		assert.Equal(t, 2, count)
	})
	t.Run("empty extensions yields nothing", func(t *testing.T) {
		var em tax.Extensions
		var count int
		for range em.All() {
			count++
		}
		assert.Equal(t, 0, count)
	})
}

func TestExtensionsMarshalJSON(t *testing.T) {
	t.Run("zero marshals to null", func(t *testing.T) {
		var em tax.Extensions
		data, err := json.Marshal(em)
		assert.NoError(t, err)
		assert.Equal(t, "null", string(data))
	})
	t.Run("keys are sorted alphabetically", func(t *testing.T) {
		// Build the same Extensions from different insertion orders;
		// JSON output must be byte-identical.
		em1 := tax.MakeExtensions().
			Set("zeta", "z").
			Set("alpha", "a").
			Set("mu", "m")
		em2 := tax.MakeExtensions().
			Set("mu", "m").
			Set("zeta", "z").
			Set("alpha", "a")
		b1, err := json.Marshal(em1)
		require.NoError(t, err)
		b2, err := json.Marshal(em2)
		require.NoError(t, err)
		assert.Equal(t, string(b1), string(b2))
		assert.Equal(t, `{"alpha":"a","mu":"m","zeta":"z"}`, string(b1))
	})
	t.Run("round-trip preserves entries", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{"a": "1", "b": "2"})
		data, err := json.Marshal(em)
		require.NoError(t, err)
		var out tax.Extensions
		require.NoError(t, json.Unmarshal(data, &out))
		assert.True(t, em.Equals(out))
	})
	t.Run("omitzero skips field in struct", func(t *testing.T) {
		type wrapper struct {
			Ext tax.Extensions `json:"ext,omitzero"`
		}
		w := wrapper{}
		data, err := json.Marshal(w)
		require.NoError(t, err)
		assert.Equal(t, `{}`, string(data))
	})
}

func TestExtensionsUnmarshalJSON(t *testing.T) {
	t.Run("null unmarshals to zero", func(t *testing.T) {
		var em tax.Extensions
		require.NoError(t, json.Unmarshal([]byte("null"), &em))
		assert.True(t, em.IsZero())
	})
	t.Run("empty object unmarshals to zero", func(t *testing.T) {
		var em tax.Extensions
		require.NoError(t, json.Unmarshal([]byte("{}"), &em))
		assert.True(t, em.IsZero())
	})
	t.Run("object populates entries", func(t *testing.T) {
		var em tax.Extensions
		require.NoError(t, json.Unmarshal([]byte(`{"a":"1","b":"2"}`), &em))
		assert.Equal(t, 2, em.Len())
		assert.Equal(t, "1", em.Get("a").String())
		assert.Equal(t, "2", em.Get("b").String())
	})
	t.Run("wrong-shape JSON returns error", func(t *testing.T) {
		// A JSON array is syntactically valid, so our UnmarshalJSON is
		// called and fails the inner unmarshal-to-map step.
		var em tax.Extensions
		assert.Error(t, em.UnmarshalJSON([]byte("[1,2,3]")))
	})
}

func TestExtensionsFromValue(t *testing.T) {
	t.Run("value type", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{"a": "1"})
		out, ok := tax.ExtensionsFromValue(em)
		assert.True(t, ok)
		assert.True(t, em.Equals(out))
	})
	t.Run("pointer type", func(t *testing.T) {
		em := tax.ExtensionsOf(tax.ExtMap{"a": "1"})
		out, ok := tax.ExtensionsFromValue(&em)
		assert.True(t, ok)
		assert.True(t, em.Equals(out))
	})
	t.Run("nil pointer", func(t *testing.T) {
		var em *tax.Extensions
		out, ok := tax.ExtensionsFromValue(em)
		assert.True(t, ok)
		assert.True(t, out.IsZero())
	})
	t.Run("unrelated value", func(t *testing.T) {
		out, ok := tax.ExtensionsFromValue(42)
		assert.False(t, ok)
		assert.True(t, out.IsZero())
	})
}

func TestExtensionsRuleTestInterface(t *testing.T) {
	// ExtensionsRule also implements the rules.Test interface via Check and
	// String, which are used when the rule is embedded in rules.When conditions.
	rule := tax.ExtensionsRequire(untdid.ExtKeyDocumentType)
	assert.Equal(t, "ext require [untdid-document-type]", rule.String())

	em := tax.ExtensionsOf(tax.ExtMap{untdid.ExtKeyDocumentType: "326"})
	assert.True(t, rule.Check(em))
	assert.True(t, rule.Check(&em))

	assert.False(t, rule.Check(tax.Extensions{}))
	// Non-extensions values short-circuit to pass (Validate returns nil).
	assert.True(t, rule.Check("not an extensions"))
}

func TestExtensionsJSONSchemaExtendShape(t *testing.T) {
	// JSONSchemaExtend is responsible for converting the struct reflection
	// (empty properties / additionalProperties: false) into a
	// patternProperties schema that mirrors the old map-based shape.
	s := &jsonschema.Schema{
		Type:                 "object",
		Properties:           jsonschema.NewProperties(),
		AdditionalProperties: jsonschema.FalseSchema,
	}
	var em tax.Extensions
	em.JSONSchemaExtend(s)
	assert.Nil(t, s.Properties, "Properties should be cleared")
	assert.Nil(t, s.AdditionalProperties, "AdditionalProperties should be cleared")
	require.Len(t, s.PatternProperties, 1)
	pp, ok := s.PatternProperties[cbc.KeyPattern]
	require.True(t, ok, "PatternProperties must be keyed by cbc.KeyPattern")
	assert.Equal(t, "https://gobl.org/draft-0/cbc/code", pp.Ref,
		"pattern value must reference cbc.Code via schema.Lookup")
}
