package c14n_test

import (
	"testing"

	"github.com/invopop/gobl/c14n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObjectMarshalJSON(t *testing.T) {
	t.Run("with empty object", func(t *testing.T) {
		o := c14n.Object{}
		d, err := o.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, "{}", string(d))
	})

	t.Run("with single attribute", func(t *testing.T) {
		o := c14n.Object{
			Attributes: []*c14n.Attribute{
				{
					Key:   "name",
					Value: c14n.String("test"),
				},
			},
		}
		d, err := o.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `{"name":"test"}`, string(d))
	})

	t.Run("with multiple attributes", func(t *testing.T) {
		o := c14n.Object{
			Attributes: []*c14n.Attribute{
				{
					Key:   "name",
					Value: c14n.String("test"),
				},
				{
					Key:   "age",
					Value: c14n.Integer(42),
				},
				{
					Key:   "active",
					Value: c14n.Bool(true),
				},
			},
		}
		d, err := o.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `{"name":"test","age":42,"active":true}`, string(d))
	})

	t.Run("with nested object", func(t *testing.T) {
		nested := &c14n.Object{
			Attributes: []*c14n.Attribute{
				{
					Key:   "nested",
					Value: c14n.String("value"),
				},
			},
		}
		o := c14n.Object{
			Attributes: []*c14n.Attribute{
				{
					Key:   "outer",
					Value: c14n.String("value"),
				},
				{
					Key:   "inner",
					Value: nested,
				},
			},
		}
		d, err := o.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `{"outer":"value","inner":{"nested":"value"}}`, string(d))
	})

	t.Run("with null attribute skipped", func(t *testing.T) {
		o := c14n.Object{
			Attributes: []*c14n.Attribute{
				{
					Key:   "name",
					Value: c14n.String("test"),
				},
				{
					Key:   "null_field",
					Value: c14n.Null{},
				},
				{
					Key:   "age",
					Value: c14n.Integer(42),
				},
			},
		}
		d, err := o.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `{"name":"test","age":42}`, string(d))
	})

	t.Run("with marshal error in attribute", func(t *testing.T) {
		o := c14n.Object{
			Attributes: []*c14n.Attribute{
				{
					Key:   "valid",
					Value: c14n.String("test"),
				},
				{
					Key:   "invalid",
					Value: c14n.String("invalid UTF-8: \xff\xfe"),
				},
			},
		}
		_, err := o.MarshalJSON()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "json: unsupported value")
	})
}

func TestArrayMarshalJSON(t *testing.T) {
	t.Run("with empty array", func(t *testing.T) {
		a := c14n.Array{}
		d, err := a.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, "[]", string(d))
	})

	t.Run("with single element", func(t *testing.T) {
		a := c14n.Array{
			Values: []c14n.Canonicalable{
				c14n.String("test"),
			},
		}
		d, err := a.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `["test"]`, string(d))
	})

	t.Run("with multiple elements", func(t *testing.T) {
		a := c14n.Array{
			Values: []c14n.Canonicalable{
				c14n.String("first"),
				c14n.Integer(42),
				c14n.Bool(true),
			},
		}
		d, err := a.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `["first",42,true]`, string(d))
	})

	t.Run("with nested array", func(t *testing.T) {
		nested := &c14n.Array{
			Values: []c14n.Canonicalable{
				c14n.String("nested"),
			},
		}
		a := c14n.Array{
			Values: []c14n.Canonicalable{
				c14n.String("outer"),
				nested,
			},
		}
		d, err := a.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `["outer",["nested"]]`, string(d))
	})

	t.Run("with marshal error", func(t *testing.T) {
		a := c14n.Array{
			Values: []c14n.Canonicalable{
				c14n.String("valid"),
				c14n.String("invalid UTF-8: \xff\xfe"),
			},
		}
		_, err := a.MarshalJSON()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "json: unsupported value")
	})
}

func TestStringMarshalJSON(t *testing.T) {
	s := c14n.String(`This is "a" test with quotes`)
	d, err := s.MarshalJSON()
	if err != nil {
		t.Errorf("unexpected error: %v", err.Error())
	}
	if string(d) != `"This is \"a\" test with quotes"` {
		t.Errorf("unexpected output, got: %v", string(d))
	}
	t.Run("with line feed", func(t *testing.T) {
		s := c14n.String("This is a test\nwith line feed")
		d, err := s.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `"This is a test\nwith line feed"`, string(d))
	})
	t.Run("with line feed and cr", func(t *testing.T) {
		s := c14n.String("This is a test\n\rwith line feed")
		d, err := s.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `"This is a test\n\rwith line feed"`, string(d))
	})
	t.Run("with tab", func(t *testing.T) {
		s := c14n.String("This is a test\twith tab")
		d, err := s.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `"This is a test\twith tab"`, string(d))
	})
	t.Run("with formfeed", func(t *testing.T) {
		s := c14n.String("This is a test\fwith formfeed")
		d, err := s.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `"This is a test\fwith formfeed"`, string(d))
	})
	t.Run("with backspace", func(t *testing.T) {
		s := c14n.String("This is a test\bwith backspace")
		d, err := s.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `"This is a test\bwith backspace"`, string(d))
	})
	t.Run("with unicode", func(t *testing.T) {
		s := c14n.String("This is a test with unicode: \u0001\u001f")
		d, err := s.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `"This is a test with unicode: \u0001\u001F"`, string(d))
	})

	t.Run("with invalid UTF-8", func(t *testing.T) {
		s := c14n.String("This is a test with invalid UTF-8: \xff\xfe")
		_, err := s.MarshalJSON()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "json: unsupported value")
	})
}

func TestNullMarshalJSON(t *testing.T) {
	n := c14n.Null{}

	d, err := n.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if string(d) != "null" {
		t.Errorf("got unexpected result: %v", string(d))
	}
}

func TestAttributeMarshalJSON(t *testing.T) {
	t.Run("with valid attribute", func(t *testing.T) {
		a := c14n.Attribute{
			Key:   "test",
			Value: c14n.String("This is a test"),
		}
		d, err := a.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, `"test":"This is a test"`, string(d))
	})
	t.Run("with encoding error in key", func(t *testing.T) {
		a := c14n.Attribute{
			Key:   "test \xff\xfe",
			Value: c14n.String("This is a test"),
		}
		_, err := a.MarshalJSON()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "json: unsupported value")
	})
	t.Run("with encoding error in value", func(t *testing.T) {
		a := c14n.Attribute{
			Key:   "test",
			Value: c14n.String("This is a test with invalid UTF-8: \xff\xfe"),
		}
		_, err := a.MarshalJSON()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "json: unsupported value")
	})

}

func TestFloatMarshalJSON(t *testing.T) {
	f := c14n.Float(0.0)

	d, err := f.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if string(d) != "0.0E0" {
		t.Errorf("got unexpected result: %v", string(d))
	}

	f = c14n.Float(1.0)
	d, err = f.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if string(d) != "1.0E0" {
		t.Errorf("got unexpected result: %v", string(d))
	}

	f = c14n.Float(123.5)
	d, err = f.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if string(d) != "1.235E2" {
		t.Errorf("got unexpected result: %v", string(d))
	}

	f = c14n.Float(123456789123456.0)
	d, err = f.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if string(d) != "1.23456789123456E14" {
		t.Errorf("got unexpected result: %v", string(d))
	}

	f = c14n.Float(0.000001234567891234560)
	d, err = f.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if string(d) != "1.23456789123456E-6" {
		t.Errorf("got unexpected result: %v", string(d))
	}

	f = c14n.Float(0.0000000000000000001234567891234560)
	d, err = f.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if string(d) != "1.23456789123456E-19" {
		t.Errorf("got unexpected result: %v", string(d))
	}

	f = c14n.Float(1.234567891234560000e-110)
	d, err = f.MarshalJSON()
	if err != nil {
		t.Errorf("did not expect error: %v", err.Error())
	}
	if string(d) != "1.23456789123456E-110" {
		t.Errorf("got unexpected result: %v", string(d))
	}
}
