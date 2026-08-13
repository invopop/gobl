package rules

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// namedString is a named string type for testing normalizeValue and StringOrBytes.
type namedString string

// namedInt is a named int type.
type namedInt int

// namedUint is a named uint type.
type namedUint uint

// namedFloat is a named float type.
type namedFloat float64

// namedBool is a named bool type.
type namedBool bool

func TestStringOrBytes(t *testing.T) {
	t.Run("plain string", func(t *testing.T) {
		isStr, str, isBytes, bs := StringOrBytes("hello")
		assert.True(t, isStr)
		assert.Equal(t, "hello", str)
		assert.False(t, isBytes)
		assert.Nil(t, bs)
	})

	t.Run("named string type", func(t *testing.T) {
		isStr, str, isBytes, bs := StringOrBytes(namedString("world"))
		assert.True(t, isStr)
		assert.Equal(t, "world", str)
		assert.False(t, isBytes)
		assert.Nil(t, bs)
	})

	t.Run("byte slice", func(t *testing.T) {
		isStr, str, isBytes, bs := StringOrBytes([]byte{0x01, 0x02})
		assert.False(t, isStr)
		assert.Equal(t, "", str)
		assert.True(t, isBytes)
		assert.Equal(t, []byte{0x01, 0x02}, bs)
	})

	t.Run("integer (neither)", func(t *testing.T) {
		isStr, str, isBytes, bs := StringOrBytes(42)
		assert.False(t, isStr)
		assert.Equal(t, "", str)
		assert.False(t, isBytes)
		assert.Nil(t, bs)
	})

	t.Run("empty string", func(t *testing.T) {
		isStr, str, _, _ := StringOrBytes("")
		assert.True(t, isStr)
		assert.Equal(t, "", str)
	})
}

func TestLengthOfValue(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		n, err := LengthOfValue("hello")
		require.NoError(t, err)
		assert.Equal(t, 5, n)
	})

	t.Run("slice", func(t *testing.T) {
		n, err := LengthOfValue([]int{1, 2, 3})
		require.NoError(t, err)
		assert.Equal(t, 3, n)
	})

	t.Run("map", func(t *testing.T) {
		n, err := LengthOfValue(map[string]int{"a": 1})
		require.NoError(t, err)
		assert.Equal(t, 1, n)
	})

	t.Run("array", func(t *testing.T) {
		n, err := LengthOfValue([2]int{1, 2})
		require.NoError(t, err)
		assert.Equal(t, 2, n)
	})

	t.Run("integer returns error", func(t *testing.T) {
		_, err := LengthOfValue(42)
		assert.Error(t, err)
	})
}

func TestIsEmpty(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		assert.True(t, IsEmpty(""))
	})

	t.Run("non-empty string", func(t *testing.T) {
		assert.False(t, IsEmpty("hi"))
	})

	t.Run("zero int", func(t *testing.T) {
		assert.True(t, IsEmpty(0))
	})

	t.Run("non-zero int", func(t *testing.T) {
		assert.False(t, IsEmpty(42))
	})

	t.Run("zero uint", func(t *testing.T) {
		assert.True(t, IsEmpty(uint(0)))
	})

	t.Run("non-zero uint", func(t *testing.T) {
		assert.False(t, IsEmpty(uint(5)))
	})

	t.Run("zero float64", func(t *testing.T) {
		assert.True(t, IsEmpty(0.0))
	})

	t.Run("non-zero float64", func(t *testing.T) {
		assert.False(t, IsEmpty(3.14))
	})

	t.Run("false bool", func(t *testing.T) {
		assert.True(t, IsEmpty(false))
	})

	t.Run("true bool", func(t *testing.T) {
		assert.False(t, IsEmpty(true))
	})

	t.Run("empty slice", func(t *testing.T) {
		assert.True(t, IsEmpty([]int{}))
	})

	t.Run("non-empty slice", func(t *testing.T) {
		assert.False(t, IsEmpty([]int{1}))
	})

	t.Run("empty map", func(t *testing.T) {
		assert.True(t, IsEmpty(map[string]int{}))
	})

	t.Run("nil slice", func(t *testing.T) {
		var s []int
		assert.True(t, IsEmpty(s))
	})

	t.Run("nil pointer", func(t *testing.T) {
		var p *int
		assert.True(t, IsEmpty(p))
	})

	t.Run("non-nil pointer to empty string", func(t *testing.T) {
		s := ""
		assert.True(t, IsEmpty(&s))
	})

	t.Run("non-nil pointer to non-empty string", func(t *testing.T) {
		s := "hi"
		assert.False(t, IsEmpty(&s))
	})

	t.Run("nil interface (reflect.Invalid)", func(t *testing.T) {
		assert.True(t, IsEmpty(nil))
	})

	t.Run("zero time.Time", func(t *testing.T) {
		assert.True(t, IsEmpty(time.Time{}))
	})

	t.Run("non-zero time.Time", func(t *testing.T) {
		assert.False(t, IsEmpty(time.Now()))
	})

	t.Run("non-time struct is not empty", func(t *testing.T) {
		type S struct{ X int }
		assert.False(t, IsEmpty(S{}))
	})
}

// testValuer implements driver.Valuer for testing Indirect.
type testValuer struct {
	val driver.Value
	err error
}

func (v testValuer) Value() (driver.Value, error) {
	return v.val, v.err
}

func TestIndirect(t *testing.T) {
	t.Run("plain value", func(t *testing.T) {
		v, isNil := Indirect("hello")
		assert.Equal(t, "hello", v)
		assert.False(t, isNil)
	})

	t.Run("nil pointer", func(t *testing.T) {
		var p *int
		v, isNil := Indirect(p)
		assert.Nil(t, v)
		assert.True(t, isNil)
	})

	t.Run("non-nil pointer", func(t *testing.T) {
		x := 42
		v, isNil := Indirect(&x)
		assert.Equal(t, 42, v)
		assert.False(t, isNil)
	})

	t.Run("double pointer", func(t *testing.T) {
		x := 42
		p := &x
		v, isNil := Indirect(&p)
		assert.Equal(t, 42, v)
		assert.False(t, isNil)
	})

	t.Run("nil interface (reflect.Invalid)", func(t *testing.T) {
		v, isNil := Indirect(nil)
		assert.Nil(t, v)
		assert.True(t, isNil)
	})

	t.Run("nil slice", func(t *testing.T) {
		var s []int
		v, isNil := Indirect(s)
		assert.Nil(t, v)
		assert.True(t, isNil)
	})

	t.Run("nil map", func(t *testing.T) {
		var m map[string]int
		v, isNil := Indirect(m)
		assert.Nil(t, v)
		assert.True(t, isNil)
	})

	t.Run("nil func", func(t *testing.T) {
		var f func()
		v, isNil := Indirect(f)
		assert.Nil(t, v)
		assert.True(t, isNil)
	})

	t.Run("nil chan", func(t *testing.T) {
		var ch chan int
		v, isNil := Indirect(ch)
		assert.Nil(t, v)
		assert.True(t, isNil)
	})

	t.Run("driver.Valuer with non-nil value", func(t *testing.T) {
		tv := testValuer{val: "from-valuer"}
		v, isNil := Indirect(tv)
		assert.Equal(t, "from-valuer", v)
		assert.False(t, isNil)
	})

	t.Run("driver.Valuer with nil value", func(t *testing.T) {
		tv := testValuer{val: nil}
		v, isNil := Indirect(tv)
		assert.Nil(t, v)
		assert.True(t, isNil)
	})
}

func TestNormalizeValue(t *testing.T) {
	t.Run("named string", func(t *testing.T) {
		v := normalizeValue(namedString("abc"))
		assert.Equal(t, "abc", v)
	})

	t.Run("named int", func(t *testing.T) {
		v := normalizeValue(namedInt(7))
		assert.Equal(t, int64(7), v)
	})

	t.Run("named uint", func(t *testing.T) {
		v := normalizeValue(namedUint(8))
		assert.Equal(t, uint64(8), v)
	})

	t.Run("named float", func(t *testing.T) {
		v := normalizeValue(namedFloat(1.5))
		assert.Equal(t, 1.5, v)
	})

	t.Run("named bool", func(t *testing.T) {
		v := normalizeValue(namedBool(true))
		assert.Equal(t, true, v)
	})

	t.Run("struct passthrough", func(t *testing.T) {
		type S struct{ X int }
		s := S{X: 1}
		v := normalizeValue(s)
		assert.Equal(t, s, v)
	})

	t.Run("plain string", func(t *testing.T) {
		v := normalizeValue("hello")
		assert.Equal(t, "hello", v)
	})
}
