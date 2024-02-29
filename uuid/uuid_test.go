package uuid_test

import (
	"encoding/json"
	"testing"

	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUUIDParsing(t *testing.T) {
	v1s := "03907310-8daa-11eb-8dcd-0242ac130003"
	v4s := "0def554c-54fd-4b3b-9ea0-4f2d288d4435"

	u1, err := uuid.Parse(v1s)
	assert.NoError(t, err)

	if u1.Version() != 1 {
		t.Errorf("did not parse a v1 UUID")
	}

	u4, err := uuid.Parse(v4s)
	if err != nil {
		t.Errorf("did not expect error, got: %v", err.Error())
	}
	if u4.Version() != 4 {
		t.Errorf("did not parse a v4 UUID")
	}

	u1 = uuid.ShouldParse("")
	assert.True(t, u1.IsZero())
	u1 = uuid.ShouldParse("fooo")
	assert.True(t, u1.IsZero())
	u1 = uuid.ShouldParse(v1s)
	assert.Equal(t, v1s, u1.String())
}

func TestUUIDIsZero(t *testing.T) {
	var up1 *uuid.UUID
	assert.True(t, up1.IsZero())
	var u1 uuid.UUID
	assert.True(t, u1.IsZero())

	up1 = uuid.NewV1()
	assert.False(t, up1.IsZero())
	u1 = uuid.MakeV1()
	assert.False(t, u1.IsZero())
}

func TestNormalizeUUID(t *testing.T) {
	var u *uuid.UUID
	u2 := uuid.Normalize(u)
	assert.Nil(t, u2)

	u = &uuid.UUID{}
	assert.Equal(t, "00000000-0000-0000-0000-000000000000", u.String())

	u2 = uuid.Normalize(u)
	assert.Nil(t, u2)

	u3 := uuid.MustParse("03907310-8daa-11eb-8dcd-0242ac130003")
	u2 = uuid.Normalize(&u3)
	assert.Equal(t, u3.String(), u2.String())
}

func TestUUIDJSON(t *testing.T) {
	v1s := "03907310-8daa-11eb-8dcd-0242ac130003"
	type testJSON struct {
		ID    uuid.UUID  `json:"id"`
		OptID *uuid.UUID `json:"opt_id,omitempty"`
	}

	v := testJSON{ID: uuid.MakeV1()}
	if v.ID.IsZero() {
		t.Errorf("did not expect UUID to be zero valued")
	}
	data, err := json.Marshal(v)
	if err != nil {
		t.Errorf("did not expect JSON marshal error")
	}
	if string(data) != `{"id":"`+v.ID.String()+`"}` {
		t.Errorf("did not marshal JSON as expected, got: %v", string(data))
	}

	v2 := testJSON{}
	if !v2.ID.IsZero() {
		t.Errorf("expected ID to be zero")
	}

	b := []byte(`{"id":"` + v1s + `"}`)
	if err := json.Unmarshal(b, &v2); err != nil {
		t.Errorf("did not expect unmarshal to fail, err: %v", err.Error())
	}
	if v2.ID.String() != v1s {
		t.Errorf("did not get same string back, got: %v", v2.ID.String())
	}
}

func TestUUIDv3(t *testing.T) {
	ns := uuid.MustParse("0654a3f4-8ad5-44c8-828e-c25f7ccd6550")
	u := uuid.NewV3(ns, []byte("hello, world"))

	assert.Equal(t, 3, int(u.Version()))
	assert.Equal(t, "61cfb897-b1bb-382b-bab9-a7ba465a27fa", u.String())
}

func TestUUIDv5(t *testing.T) {
	ns := uuid.MustParse("0654a3f4-8ad5-44c8-828e-c25f7ccd6550")
	u := uuid.NewV5(ns, []byte("hello, world"))

	assert.Equal(t, 5, int(u.Version()))
	assert.Equal(t, "1f53a310-2a17-5acb-b76a-c39495e5356f", u.String())
}
