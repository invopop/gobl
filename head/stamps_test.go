package head_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
)

func TestDuplicateStamps(t *testing.T) {
	st := struct {
		Stamps []*head.Stamp
	}{
		Stamps: []*head.Stamp{
			{
				Provider: cbc.Key("provider"),
				Value:    "value",
			},
			{
				Provider: cbc.Key("provider2"),
				Value:    "value2",
			},
		},
	}

	err := validation.ValidateStruct(&st,
		validation.Field(&st.Stamps, head.DetectDuplicateStamps),
	)
	assert.NoError(t, err)

	st.Stamps = append(st.Stamps, &head.Stamp{
		Provider: cbc.Key("provider"),
		Value:    "value3",
	})
	err = validation.ValidateStruct(&st,
		validation.Field(&st.Stamps, head.DetectDuplicateStamps),
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate stamp 'provider'")
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
	err := r.Validate([]*head.Stamp{
		{
			Provider: "foo",
			Value:    "bar",
		},
	})
	assert.NoError(t, err)

	err = r.Validate([]*head.Stamp{
		{
			Provider: "foo2",
			Value:    "bar",
		},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing foo stamp")
}
