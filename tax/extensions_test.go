package tax_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/addons/es/tbai"
	"github.com/invopop/gobl/addons/gr/mydata"
	"github.com/invopop/gobl/addons/mx/cfdi" // this will also prepare registers
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
)

func TestCleanExtensions(t *testing.T) {
	var em tax.Extensions

	em2 := tax.CleanExtensions(em)
	assert.Nil(t, em2)

	em = tax.Extensions{
		"key": "",
	}
	em2 = tax.CleanExtensions(em)
	assert.Nil(t, em2)

	em = tax.Extensions{
		"key": "foo",
		"bar": "",
	}
	em2 = tax.CleanExtensions(em)
	assert.NotNil(t, em2)
	assert.Len(t, em2, 1)
	assert.Equal(t, "foo", em2["key"].String())
}

func TestExtValidation(t *testing.T) {
	t.Run("with mexico", func(t *testing.T) {
		t.Run("test patterns", func(t *testing.T) {
			em := tax.Extensions{
				cfdi.ExtKeyIssuePlace: "12345",
			}
			err := em.Validate()
			assert.NoError(t, err)

			em = tax.Extensions{
				cfdi.ExtKeyIssuePlace: "123457",
			}
			err = em.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "mx-cfdi-issue-place: does not match pattern")

			kd := tax.ExtensionForKey(cfdi.ExtKeyIssuePlace)
			pt := kd.Pattern
			kd.Pattern = "[][" // invalid
			err = em.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "mx-cfdi-issue-place: error parsing regexp: missing closing ]: `[][`")
			kd.Pattern = pt // put back!
		})

		t.Run("test codes", func(t *testing.T) {
			em := tax.Extensions{
				cfdi.ExtKeyFiscalRegime: "601",
			}
			err := em.Validate()
			assert.NoError(t, err)

			em = tax.Extensions{
				cfdi.ExtKeyFiscalRegime: "000",
			}
			err = em.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "mx-cfdi-fiscal-regime: value '000' invalid")
		})
	})

	t.Run("with spain", func(t *testing.T) {
		t.Run("test good key", func(t *testing.T) {
			em := tax.Extensions{
				tbai.ExtKeyProduct: "goods",
			}
			err := em.Validate()
			assert.NoError(t, err)
		})

		t.Run("test bad key", func(t *testing.T) {
			em := tax.Extensions{
				tbai.ExtKeyProduct: "bads",
			}
			err := em.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "es-tbai-product: value 'bads' invalid")
		})

		t.Run("missing extension", func(t *testing.T) {
			em := tax.Extensions{
				"random-key": "type",
			}
			err := em.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "random-key: undefined")
		})

		t.Run("invalid key", func(t *testing.T) {
			em := tax.Extensions{
				"INVALID": "value",
			}
			err := em.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "INVALID: must be in a valid format")
		})
	})

	t.Run("with greece", func(t *testing.T) {
		t.Run("test good value", func(t *testing.T) {
			em := tax.Extensions{
				mydata.ExtKeyIncomeCat: "category1_1",
			}
			err := em.Validate()
			assert.NoError(t, err)
		})

		t.Run("test bad value", func(t *testing.T) {
			em := tax.Extensions{
				mydata.ExtKeyIncomeCat: "xxx",
			}
			err := em.Validate()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "gr-mydata-income-cat: value 'xxx' invalid")
		})
	})
}

func TestExtensionsRequiresValidation(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := validation.Validate(nil,
			tax.ExtensionsRequire(untdid.ExtKeyDocumentType),
		)
		assert.NoError(t, err)
	})
	t.Run("empty", func(t *testing.T) {
		em := tax.Extensions{}
		err := validation.Validate(em,
			tax.ExtensionsRequire(untdid.ExtKeyDocumentType),
		)
		assert.ErrorContains(t, err, "untdid-document-type: required")
	})
	t.Run("correct", func(t *testing.T) {
		em := tax.Extensions{
			untdid.ExtKeyDocumentType: "326",
		}
		err := validation.Validate(em,
			tax.ExtensionsRequire(untdid.ExtKeyDocumentType),
		)
		assert.NoError(t, err)
	})
	t.Run("correct with extras", func(t *testing.T) {
		em := tax.Extensions{
			untdid.ExtKeyDocumentType: "326",
			iso.ExtKeySchemeID:        "1234",
		}
		err := validation.Validate(em,
			tax.ExtensionsRequire(untdid.ExtKeyDocumentType),
		)
		assert.NoError(t, err)
	})
	t.Run("missing", func(t *testing.T) {
		em := tax.Extensions{
			iso.ExtKeySchemeID: "1234",
		}
		err := validation.Validate(em,
			tax.ExtensionsRequire(untdid.ExtKeyDocumentType),
		)
		assert.ErrorContains(t, err, "untdid-document-type: required")
	})
}

func TestExtensionsAllOrNoneValidation(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := validation.Validate(nil,
			tax.ExtensionsRequireAllOrNone(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID),
		)
		assert.NoError(t, err)
	})
	t.Run("empty", func(t *testing.T) {
		em := tax.Extensions{}
		err := validation.Validate(em,
			tax.ExtensionsRequireAllOrNone(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID),
		)
		assert.NoError(t, err)
	})
	t.Run("all present", func(t *testing.T) {
		em := tax.Extensions{
			untdid.ExtKeyDocumentType: "326",
			iso.ExtKeySchemeID:        "1234",
		}
		err := validation.Validate(em,
			tax.ExtensionsRequireAllOrNone(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID),
		)
		assert.NoError(t, err)
	})
	t.Run("none present", func(t *testing.T) {
		em := tax.Extensions{}
		err := validation.Validate(em,
			tax.ExtensionsRequireAllOrNone(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID),
		)
		assert.NoError(t, err)
	})
	t.Run("some present", func(t *testing.T) {
		em := tax.Extensions{
			untdid.ExtKeyDocumentType: "326",
		}
		err := validation.Validate(em,
			tax.ExtensionsRequireAllOrNone(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID),
		)
		assert.ErrorContains(t, err, "iso-scheme-id: required")
	})
	t.Run("some present reversed", func(t *testing.T) {
		em := tax.Extensions{
			iso.ExtKeySchemeID: "1234",
		}
		err := validation.Validate(em,
			tax.ExtensionsRequireAllOrNone(untdid.ExtKeyDocumentType, iso.ExtKeySchemeID),
		)
		assert.ErrorContains(t, err, "untdid-document-type: required")
	})
}

func TestExtensionsExcludeValidation(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := validation.Validate(nil,
			tax.ExtensionsExclude(untdid.ExtKeyDocumentType),
		)
		assert.NoError(t, err)
	})
	t.Run("empty", func(t *testing.T) {
		em := tax.Extensions{}
		err := validation.Validate(em,
			tax.ExtensionsExclude(untdid.ExtKeyDocumentType),
		)
		assert.NoError(t, err)
	})
	t.Run("correct", func(t *testing.T) {
		em := tax.Extensions{
			untdid.ExtKeyDocumentType: "326",
		}
		err := validation.Validate(em,
			tax.ExtensionsExclude(untdid.ExtKeyDocumentType),
		)
		assert.ErrorContains(t, err, "untdid-document-type: must be blank")
	})
	t.Run("correct with extras", func(t *testing.T) {
		em := tax.Extensions{
			untdid.ExtKeyDocumentType: "326",
			iso.ExtKeySchemeID:        "1234",
		}
		err := validation.Validate(em,
			tax.ExtensionsExclude(untdid.ExtKeyCharge),
		)
		assert.NoError(t, err)
	})
}

func TestExtensionsHasValues(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := validation.Validate(nil,
			tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, "326", "389"),
		)
		assert.NoError(t, err)
	})
	t.Run("empty", func(t *testing.T) {
		em := tax.Extensions{}
		err := validation.Validate(em,
			tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, "326", "389"),
		)
		assert.NoError(t, err)
	})
	t.Run("different extensions", func(t *testing.T) {
		em := tax.Extensions{
			iso.ExtKeySchemeID: "1234",
		}
		err := validation.Validate(em,
			tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, "326", "389"),
		)
		assert.NoError(t, err)
	})
	t.Run("has codes", func(t *testing.T) {
		em := tax.Extensions{
			untdid.ExtKeyDocumentType: "326",
		}
		err := validation.Validate(em,
			tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, "326", "389"),
		)
		assert.NoError(t, err)
	})
	t.Run("invalid code", func(t *testing.T) {
		em := tax.Extensions{
			untdid.ExtKeyDocumentType: "102",
		}
		err := validation.Validate(em,
			tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, "326", "389"),
		)
		assert.ErrorContains(t, err, "untdid-document-type: invalid value")
	})
}

func TestExtensionsExcludeCodes(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		err := validation.Validate(nil,
			tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, "380", "381"),
		)
		assert.NoError(t, err)
	})
	t.Run("empty", func(t *testing.T) {
		em := tax.Extensions{}
		err := validation.Validate(em,
			tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, "380", "381"),
		)
		assert.NoError(t, err)
	})
	t.Run("different extensions", func(t *testing.T) {
		em := tax.Extensions{
			iso.ExtKeySchemeID: "1234",
		}
		err := validation.Validate(em,
			tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, "380", "381"),
		)
		assert.NoError(t, err)
	})
	t.Run("allowed code", func(t *testing.T) {
		em := tax.Extensions{
			untdid.ExtKeyDocumentType: "326",
		}
		err := validation.Validate(em,
			tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, "380", "381"),
		)
		assert.NoError(t, err)
	})
	t.Run("excluded code", func(t *testing.T) {
		em := tax.Extensions{
			untdid.ExtKeyDocumentType: "380",
		}
		err := validation.Validate(em,
			tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, "380", "381"),
		)
		assert.ErrorContains(t, err, "untdid-document-type: value '380' not allowed")
	})
	t.Run("another excluded code", func(t *testing.T) {
		em := tax.Extensions{
			untdid.ExtKeyDocumentType: "381",
		}
		err := validation.Validate(em,
			tax.ExtensionsExcludeCodes(untdid.ExtKeyDocumentType, "380", "381"),
		)
		assert.ErrorContains(t, err, "untdid-document-type: value '381' not allowed")
	})
}

func TestExtensionsHas(t *testing.T) {
	em := tax.Extensions{
		"key": "value",
	}
	assert.True(t, em.Has("key"))
	assert.False(t, em.Has("invalid"))
}

func TestExtensionsValues(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var em tax.Extensions
		assert.Empty(t, em.Values())
	})
	t.Run("with values", func(t *testing.T) {
		em := tax.Extensions{
			"key1": "value1",
			"key2": "value2",
		}
		assert.ElementsMatch(t, []cbc.Code{"value1", "value2"}, em.Values())
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
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value"},
			want: true,
		},
		{
			name: "different",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value2"},
			want: false,
		},
		{
			name: "different keys",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key2": "value"},
			want: false,
		},
		{
			name: "different lengths",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value", "key2": "value"},
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
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value"},
			want: true,
		},
		{
			name: "different",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value2"},
			want: false,
		},
		{
			name: "different keys",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key2": "value"},
			want: false,
		},
		{
			name: "different lengths",
			em1:  tax.Extensions{"key": "value", "key2": "value"},
			em2:  tax.Extensions{"key": "value"},
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
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value"},
			want: tax.Extensions{"key": "value"},
		},
		{
			name: "nil source",
			em1:  nil,
			em2:  tax.Extensions{"key": "value"},
			want: tax.Extensions{"key": "value"},
		},
		{
			name: "nil destination",
			em1:  tax.Extensions{"key": "value"},
			em2:  nil,
			want: tax.Extensions{"key": "value"},
		},
		{
			name: "different",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value2"},
			want: tax.Extensions{"key": "value2"},
		},
		{
			name: "different keys",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key2": "value"},
			want: tax.Extensions{"key": "value", "key2": "value"},
		},
		{
			name: "different lengths",
			em1:  tax.Extensions{"key": "value"},
			em2:  tax.Extensions{"key": "value", "key2": "value"},
			want: tax.Extensions{"key": "value", "key2": "value"},
		},
		{
			name: "different lengths 2",
			em1:  tax.Extensions{"key": "value", "key2": "value"},
			em2:  tax.Extensions{"key": "value"},
			want: tax.Extensions{"key": "value", "key2": "value"},
		},
		{
			name: "different lengths 3",
			em1:  tax.Extensions{"key": "value2"},
			em2:  tax.Extensions{"key": "value", "key2": "value"},
			want: tax.Extensions{"key": "value", "key2": "value"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.em1.Merge(tt.em2))
		})
	}
}

func TestExtensionLookup(t *testing.T) {
	em := tax.Extensions{
		"key1": "foo",
		"key2": "bar",
	}
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
		em := tax.Extensions{
			"key": "value",
		}
		assert.Equal(t, "value", em.Get("key").String())
	})
	t.Run("missing", func(t *testing.T) {
		em := tax.Extensions{
			"key": "value",
		}
		assert.Equal(t, "", em.Get("missing").String())
	})
	t.Run("with sub-keys", func(t *testing.T) {
		em := tax.Extensions{
			"key": "value",
		}
		assert.Equal(t, "value", em.Get("key+foo").String())
	})
}

func TestExtensionsSet(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var em tax.Extensions
		em = em.Set("key", "value")
		assert.Equal(t, tax.Extensions{"key": "value"}, em)
	})

	t.Run("with existing value", func(t *testing.T) {
		em := tax.Extensions{"key": "value1"}
		em.Set("key", "value2")
		assert.Equal(t, tax.Extensions{"key": "value2"}, em)
	})

	t.Run("with new value", func(t *testing.T) {
		em := tax.Extensions{}
		em.Set("key", "value1")
		assert.Equal(t, tax.Extensions{"key": "value1"}, em)
	})

	t.Run("with empty value", func(t *testing.T) {
		em := tax.Extensions{"key": "value1"}
		em.Set("key", "")
		assert.Empty(t, em)
	})
}

func TestExtensionsSetIfEmpty(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var em tax.Extensions
		em = em.SetIfEmpty("key", "value")
		assert.Equal(t, tax.Extensions{"key": "value"}, em)
	})

	t.Run("with existing value", func(t *testing.T) {
		em := tax.Extensions{"key": "value1"}
		em.SetIfEmpty("key", "value2")
		assert.Equal(t, tax.Extensions{"key": "value1"}, em)
	})

	t.Run("with new value", func(t *testing.T) {
		em := tax.Extensions{}
		em.SetIfEmpty("key", "value1")
		assert.Equal(t, tax.Extensions{"key": "value1"}, em)
	})
}

func TestExtensionsSetOneOf(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var em tax.Extensions
		em = em.SetOneOf("key", "value1", "value2")
		assert.Equal(t, tax.Extensions{"key": "value1"}, em)
	})

	t.Run("with existing value", func(t *testing.T) {
		em := tax.Extensions{"key": "value1"}
		em.SetOneOf("key", "value2", "value3")
		assert.Equal(t, tax.Extensions{"key": "value2"}, em)
	})
	t.Run("with existing value and output", func(t *testing.T) {
		em := tax.Extensions{"key": "value1"}
		em = em.SetOneOf("key", "value2", "value3")
		assert.Equal(t, tax.Extensions{"key": "value2"}, em)
	})

	t.Run("with existing secondary value", func(t *testing.T) {
		em := tax.Extensions{"key": "value3"}
		em = em.SetOneOf("key", "value2", "value3")
		assert.Equal(t, tax.Extensions{"key": "value3"}, em)
	})

	t.Run("with no existing value", func(t *testing.T) {
		em := tax.Extensions{}
		em = em.SetOneOf("key", "value1", "value2")
		assert.Equal(t, tax.Extensions{"key": "value1"}, em)
	})
}

func TestExtensionsDelete(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var em tax.Extensions
		em = em.Delete("key")
		assert.Nil(t, em)
	})

	t.Run("with value", func(t *testing.T) {
		em := tax.Extensions{"key": "value"}
		em = em.Delete("key")
		assert.Empty(t, em)
	})

	t.Run("with missing value", func(t *testing.T) {
		em := tax.Extensions{"key": "value"}
		em = em.Delete("missing")
		assert.Equal(t, tax.Extensions{"key": "value"}, em)
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
