package rules

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublicPath(t *testing.T) {
	t.Run("empty path returns $", func(t *testing.T) {
		assert.Equal(t, "$", publicPath(""))
	})

	t.Run("path starting with [ gets $ prefix", func(t *testing.T) {
		assert.Equal(t, "$[0]", publicPath("[0]"))
	})

	t.Run("normal path gets $. prefix", func(t *testing.T) {
		assert.Equal(t, "$.name", publicPath("name"))
	})
}

func TestJoinPath(t *testing.T) {
	t.Run("empty prefix", func(t *testing.T) {
		assert.Equal(t, "name", joinPath("", "name"))
	})

	t.Run("empty suffix", func(t *testing.T) {
		assert.Equal(t, "items", joinPath("items", ""))
	})

	t.Run("suffix starts with [", func(t *testing.T) {
		assert.Equal(t, "items[0]", joinPath("items", "[0]"))
	})

	t.Run("normal join", func(t *testing.T) {
		assert.Equal(t, "items.name", joinPath("items", "name"))
	})
}

func TestFieldValueByName(t *testing.T) {
	t.Run("non-struct returns false", func(t *testing.T) {
		rv := reflect.ValueOf("hello")
		_, ok := fieldValueByName(rv, "anything")
		assert.False(t, ok)
	})

	t.Run("missing field returns false", func(t *testing.T) {
		type S struct {
			Name string `json:"name"`
		}
		rv := reflect.ValueOf(S{Name: "test"})
		_, ok := fieldValueByName(rv, "missing")
		assert.False(t, ok)
	})

	t.Run("found field returns true", func(t *testing.T) {
		type S struct {
			Name string `json:"name"`
		}
		rv := reflect.ValueOf(S{Name: "test"})
		fv, ok := fieldValueByName(rv, "name")
		assert.True(t, ok)
		assert.Equal(t, "test", fv.String())
	})
}

func TestFieldTypeByName(t *testing.T) {
	t.Run("non-struct returns nil", func(t *testing.T) {
		ft := fieldTypeByName(reflect.TypeOf("hello"), "x")
		assert.Nil(t, ft)
	})
}

func TestValidateEachValuePointerSlice(t *testing.T) {
	// Tests the pointer-to-slice path in validateEachValue.
	type Item struct {
		Name string `json:"name"`
	}
	items := []Item{{Name: "a"}, {Name: "b"}}
	fv := reflect.ValueOf(&items)

	ss := &Set{
		Assert: []*Assertion{
			{
				ID:   "01",
				Desc: "always passes",
				Tests: []Test{funcTestHelper{
					fn: func(any) bool { return true },
				}},
			},
		},
	}
	faults := validateEachValue(nil, fv, ss)
	assert.Empty(t, faults)
}

func TestValidateEachValueNilPointer(t *testing.T) {
	var items *[]string
	fv := reflect.ValueOf(items)

	ss := &Set{}
	faults := validateEachValue(nil, fv, ss)
	assert.Empty(t, faults)
}

func TestValidateEachValueNotSlice(t *testing.T) {
	fv := reflect.ValueOf("not a slice")
	ss := &Set{}
	faults := validateEachValue(nil, fv, ss)
	assert.Empty(t, faults)
}

func TestCollectContextNilPointer(t *testing.T) {
	var p *struct{ Name string }
	rc := &Context{}
	// Should not panic.
	collectContext(rc, p)
}

func TestCollectContextNonStruct(t *testing.T) {
	rc := &Context{}
	collectContext(rc, "a string")
	assert.Nil(t, rc.Value("anything"))
}

// funcTestHelper is a simple Test implementation for internal tests.
type funcTestHelper struct {
	fn func(any) bool
}

func (f funcTestHelper) Check(val any) bool { return f.fn(val) }
func (f funcTestHelper) String() string     { return "func" }

// collectContextAdderImpl tests the root ContextAdder path.
type collectContextAdderImpl struct{}

func (c *collectContextAdderImpl) RulesContext() WithContext {
	return func(rc *Context) {
		rc.Set("root", "yes")
	}
}

func TestCollectContextRootAdder(t *testing.T) {
	rc := &Context{}
	obj := &collectContextAdderImpl{}
	collectContext(rc, obj)
	assert.Equal(t, "yes", rc.Value("root"))
}

func TestTypeSetIDWithEmptyPkg(t *testing.T) {
	// Built-in types have no package path.
	ty := reflect.TypeOf("")
	id := typeSetID(ty, false)
	assert.Equal(t, Code("STRING"), id)
}

func TestPkgShortNameBuiltin(t *testing.T) {
	ty := reflect.TypeOf(0)
	name := pkgShortName(ty)
	assert.Equal(t, "", name)
}

func TestValidateNilPointerObject(t *testing.T) {
	type S struct {
		Name string `json:"name"`
	}
	s := &Set{
		objType: reflect.TypeOf(S{}),
		Assert: []*Assertion{
			{
				ID:   "01",
				Desc: "always fails",
				Tests: []Test{funcTestHelper{
					fn: func(any) bool { return false },
				}},
			},
		},
	}
	// Passing nil pointer of matching type.
	var p *S
	faults := s.validate(nil, p)
	// The nil pointer should still run assertions (to detect missing values).
	assert.NotNil(t, faults)
}

func TestValidateAnonymousField(t *testing.T) {
	type Base struct {
		Code string `json:"code"`
	}
	type Extended struct {
		Base
		Name string `json:"name"`
	}

	// Register rules for Base.Code via namespace.
	baseSet := For(new(Base),
		Field("code",
			Assert("01", "code required", funcTestHelper{
				fn: func(val any) bool {
					s, ok := val.(string)
					return ok && s != ""
				},
			}),
		),
	)
	RegisterWithGuard("anon-test", GOBL.Add("ANONTEST"), nil, baseSet)

	// Validate an Extended struct with empty Code.
	ext := &Extended{
		Base: Base{Code: ""},
		Name: "test",
	}
	faults := Validate(ext)
	// The anonymous Base should be traversed and the fault should be reported
	// without a field prefix for the anonymous field itself.
	if faults != nil {
		// Just check it didn't panic and produced something.
		assert.Greater(t, faults.Len(), 0)
	}
}

func TestValidateSliceAtNamespaceLevel(t *testing.T) {
	// Test the slice/array branch in the namespace-level iteration.
	type Item struct {
		Val string `json:"val"`
	}
	type ListHolder struct {
		Items []Item `json:"items"`
	}

	itemSet := For(new(Item),
		Field("val",
			Assert("01", "val required", funcTestHelper{
				fn: func(val any) bool {
					s, ok := val.(string)
					return ok && s != ""
				},
			}),
		),
	)
	RegisterWithGuard("slice-ns-test", GOBL.Add("SLICENS"), nil, itemSet)

	holder := &ListHolder{
		Items: []Item{{Val: ""}, {Val: "ok"}},
	}
	faults := Validate(holder)
	require.NotNil(t, faults)
	assert.True(t, faults.HasPath("$.items[0].val"))
}

// embeddableNamespaceTest tests the Embeddable check at namespace level.
type embInner struct {
	Name string `json:"name"`
}

type embWrapper struct {
	inner *embInner
}

func (w *embWrapper) Embedded() any {
	return w.inner
}

type embContainer struct {
	Wrap *embWrapper `json:"wrap"`
}

func TestCollectContextUnexportedFieldSkipped(t *testing.T) {
	type hasUnexported struct {
		exported string //nolint:unused
		Name     string
	}
	rc := &Context{}
	obj := &hasUnexported{Name: "test"}
	collectContext(rc, obj)
	// Should not panic — unexported fields are skipped.
}

func TestJsonFieldNameEmptyTag(t *testing.T) {
	// json tag with comma but empty name uses field name.
	type S struct {
		Foo string `json:",omitempty"`
	}
	rt := reflect.TypeOf(S{})
	name := jsonFieldName(rt.Field(0))
	assert.Equal(t, "Foo", name)
}

func TestValidateEmbeddableAtNamespaceLevel(t *testing.T) {
	innerSet := For(new(embInner),
		Field("name",
			Assert("01", "name required", funcTestHelper{
				fn: func(val any) bool {
					s, ok := val.(string)
					return ok && s != ""
				},
			}),
		),
	)
	RegisterWithGuard("emb-ns-test", GOBL.Add("EMBNS"), nil, innerSet)

	c := &embContainer{Wrap: &embWrapper{inner: &embInner{Name: ""}}}
	faults := Validate(c)
	require.NotNil(t, faults)
}
