package head_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/stretchr/testify/assert"
)

func TestDuplicateStamps(t *testing.T) {
	stamps := []*head.Stamp{
		{
			Provider: cbc.Key("provider"),
			Value:    "value",
		},
		{
			Provider: cbc.Key("provider2"),
			Value:    "value2",
		},
	}

	assert.True(t, head.DetectDuplicateStamps.Check(stamps))

	stamps = append(stamps, &head.Stamp{
		Provider: cbc.Key("provider"),
		Value:    "value3",
	})
	assert.False(t, head.DetectDuplicateStamps.Check(stamps))
}

func TestAddStamp(t *testing.T) {
	st := struct {
		Stamps []*head.Stamp
	}{
		Stamps: []*head.Stamp{
			{
				Provider: cbc.Key("provider"),
				Value:    "value",
			},
		},
	}
	st.Stamps = head.AddStamp(st.Stamps, &head.Stamp{
		Provider: cbc.Key("provider"),
		Value:    "new value",
	})
	assert.Len(t, st.Stamps, 1)
	assert.Equal(t, "new value", st.Stamps[0].Value)
}

func TestGetStamp(t *testing.T) {
	st := struct {
		Stamps []*head.Stamp
	}{
		Stamps: []*head.Stamp{
			{
				Provider: cbc.Key("provider"),
				Value:    "value",
			},
			{
				Provider: cbc.Key("foo"),
				Value:    "bar",
			},
		},
	}
	x := head.GetStamp(st.Stamps, cbc.Key("foo"))
	assert.Equal(t, "bar", x.Value)

	x = head.GetStamp(st.Stamps, cbc.Key("bad"))
	assert.Nil(t, x)
}

func TestNormalizeStamp(t *testing.T) {
	st := head.NormalizeStamps(nil)
	assert.Nil(t, st)

	st = []*head.Stamp{}
	st2 := head.NormalizeStamps(st)
	assert.Nil(t, st2)
	assert.Len(t, st2, 0)

	st = []*head.Stamp{
		{
			Provider: "foo",
			Value:    "",
		},
	}
	st = head.NormalizeStamps(st)
	assert.Len(t, st, 0)

	st = []*head.Stamp{
		{
			Provider: "foo",
			Value:    "bar",
		},
		{
			Provider: "foo2",
			Value:    "",
		},
	}
	st = head.NormalizeStamps(st)
	assert.Len(t, st, 1)
	assert.Equal(t, "bar", st[0].Value)
}

func TestStampsHas(t *testing.T) {
	r := head.StampsHas(cbc.Key("foo"))
	assert.True(t, r.Check([]*head.Stamp{
		{
			Provider: "foo",
			Value:    "bar",
		},
	}))
	assert.False(t, r.Check([]*head.Stamp{
		{
			Provider: "foo2",
			Value:    "bar",
		},
	}))
}
