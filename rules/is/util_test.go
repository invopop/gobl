package is

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// utilNamedString is a named string type for testing normalizeValue.
type utilNamedString string

// utilNamedInt is a named int type.
type utilNamedInt int

// utilNamedUint is a named uint type.
type utilNamedUint uint

// utilNamedFloat is a named float type.
type utilNamedFloat float64

// utilNamedBool is a named bool type.
type utilNamedBool bool

func TestUtilNormalizeValue(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		assert.Equal(t, "abc", normalizeValue("abc"))
	})

	t.Run("named string", func(t *testing.T) {
		assert.Equal(t, "abc", normalizeValue(utilNamedString("abc")))
	})

	t.Run("int", func(t *testing.T) {
		assert.Equal(t, int64(5), normalizeValue(utilNamedInt(5)))
	})

	t.Run("uint", func(t *testing.T) {
		assert.Equal(t, uint64(8), normalizeValue(utilNamedUint(8)))
	})

	t.Run("float", func(t *testing.T) {
		assert.Equal(t, 1.5, normalizeValue(utilNamedFloat(1.5)))
	})

	t.Run("bool", func(t *testing.T) {
		assert.Equal(t, true, normalizeValue(utilNamedBool(true)))
	})

	t.Run("struct passthrough", func(t *testing.T) {
		type S struct{ X int }
		s := S{X: 1}
		assert.Equal(t, s, normalizeValue(s))
	})
}

// utilValuer implements driver.Valuer for testing indirectValuer.
type utilValuer struct {
	val driver.Value
	err error
}

func (v utilValuer) Value() (driver.Value, error) {
	return v.val, v.err
}

func TestUtilLengthOfValue(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		n, err := LengthOfValue("hello")
		assert.NoError(t, err)
		assert.Equal(t, 5, n)
	})

	t.Run("slice", func(t *testing.T) {
		n, err := LengthOfValue([]int{1, 2})
		assert.NoError(t, err)
		assert.Equal(t, 2, n)
	})

	t.Run("map", func(t *testing.T) {
		n, err := LengthOfValue(map[string]int{"a": 1})
		assert.NoError(t, err)
		assert.Equal(t, 1, n)
	})

	t.Run("array", func(t *testing.T) {
		n, err := LengthOfValue([3]int{1, 2, 3})
		assert.NoError(t, err)
		assert.Equal(t, 3, n)
	})

	t.Run("int returns error", func(t *testing.T) {
		_, err := LengthOfValue(42)
		assert.Error(t, err)
	})
}

func TestUtilStringOrBytes(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		isStr, str, isBytes, _ := StringOrBytes("hello")
		assert.True(t, isStr)
		assert.Equal(t, "hello", str)
		assert.False(t, isBytes)
	})

	t.Run("bytes", func(t *testing.T) {
		isStr, _, isBytes, bs := StringOrBytes([]byte{1, 2})
		assert.False(t, isStr)
		assert.True(t, isBytes)
		assert.Equal(t, []byte{1, 2}, bs)
	})

	t.Run("int neither", func(t *testing.T) {
		isStr, _, isBytes, _ := StringOrBytes(42)
		assert.False(t, isStr)
		assert.False(t, isBytes)
	})
}

func TestUtilIndirect(t *testing.T) {
	t.Run("driver.Valuer with non-nil value", func(t *testing.T) {
		v, isNil := Indirect(utilValuer{val: "hello"})
		assert.Equal(t, "hello", v)
		assert.False(t, isNil)
	})

	t.Run("driver.Valuer with nil value", func(t *testing.T) {
		v, isNil := Indirect(utilValuer{val: nil})
		assert.Nil(t, v)
		assert.True(t, isNil)
	})

	t.Run("driver.Valuer with error", func(t *testing.T) {
		v, isNil := Indirect(utilValuer{val: nil, err: assert.AnError})
		assert.Nil(t, v)
		assert.True(t, isNil)
	})

	t.Run("plain value", func(t *testing.T) {
		v, isNil := Indirect("world")
		assert.Equal(t, "world", v)
		assert.False(t, isNil)
	})

	t.Run("nil pointer", func(t *testing.T) {
		var p *int
		v, isNil := Indirect(p)
		assert.Nil(t, v)
		assert.True(t, isNil)
	})
}

func TestUtilEmptyValue(t *testing.T) {
	t.Run("zero uint", func(t *testing.T) {
		assert.True(t, emptyValue(uint(0)))
	})

	t.Run("non-zero uint", func(t *testing.T) {
		assert.False(t, emptyValue(uint(5)))
	})

	t.Run("zero float", func(t *testing.T) {
		assert.True(t, emptyValue(float64(0)))
	})

	t.Run("non-zero float", func(t *testing.T) {
		assert.False(t, emptyValue(float64(3.14)))
	})

	t.Run("non-time struct is not empty", func(t *testing.T) {
		type S struct{ X int }
		assert.False(t, emptyValue(S{}))
	})

	t.Run("zero time is empty", func(t *testing.T) {
		assert.True(t, emptyValue(time.Time{}))
	})

	t.Run("nil interface", func(t *testing.T) {
		assert.True(t, emptyValue(nil))
	})

	t.Run("nil pointer", func(t *testing.T) {
		var p *int
		assert.True(t, emptyValue(p))
	})

	t.Run("pointer to empty string", func(t *testing.T) {
		s := ""
		assert.True(t, emptyValue(&s))
	})
}
