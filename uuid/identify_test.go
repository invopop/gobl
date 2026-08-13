package uuid_test

import (
	"testing"

	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIdentify(t *testing.T) {
	type Document struct {
		uuid.Identify
		Name string
	}
	doc := new(Document)
	id := uuid.V1()
	doc.SetUUID(id)
	assert.Equal(t, id, doc.GetUUID())
}

func TestIdentifyParse(t *testing.T) {
	doc := struct {
		uuid.Identify
		Name string
	}{
		Identify: uuid.IdentifyParse("03907310-8daa-11eb-8dcd-0242ac130003"),
		Name:     "test",
	}
	assert.NotEmpty(t, doc.GetUUID())
	assert.Equal(t, uuid.UUID("03907310-8daa-11eb-8dcd-0242ac130003"), doc.GetUUID())
}

func TestIdentifyV1(t *testing.T) {
	doc := struct {
		uuid.Identify
		Name string
	}{
		Identify: uuid.IdentifyV1(),
		Name:     "test",
	}
	assert.NotEmpty(t, doc.GetUUID())
	assert.False(t, doc.UUID.Timestamp().IsZero())
}

func TestIdentifyV4(t *testing.T) {
	doc := struct {
		uuid.Identify
		Name string
	}{
		Identify: uuid.IdentifyV4(),
		Name:     "test",
	}
	assert.NotEmpty(t, doc.GetUUID())
}

func TestIdentifyV3(t *testing.T) {
	ns := uuid.MustParse("0654a3f4-8ad5-44c8-828e-c25f7ccd6550")
	data := []byte("hello, world")
	doc := struct {
		uuid.Identify
		Name string
	}{
		Identify: uuid.IdentifyV3(ns, data),
		Name:     "test",
	}
	assert.NotEmpty(t, doc.GetUUID())
	assert.Equal(t, doc.GetUUID(), uuid.UUID("61cfb897-b1bb-382b-bab9-a7ba465a27fa"))
}

func TestIdentifyV5(t *testing.T) {
	ns := uuid.MustParse("0654a3f4-8ad5-44c8-828e-c25f7ccd6550")
	data := []byte("hello, world")
	doc := struct {
		uuid.Identify
		Name string
	}{
		Identify: uuid.IdentifyV5(ns, data),
		Name:     "test",
	}
	assert.NotEmpty(t, doc.GetUUID())
	assert.Equal(t, doc.GetUUID(), uuid.UUID("1f53a310-2a17-5acb-b76a-c39495e5356f"))
}

func TestIdentifyV6(t *testing.T) {
	doc := struct {
		uuid.Identify
		Name string
	}{
		Identify: uuid.IdentifyV6(),
		Name:     "test",
	}
	assert.NotEmpty(t, doc.GetUUID())
	assert.False(t, doc.UUID.Timestamp().IsZero())
}

func TestIdentifyV7(t *testing.T) {
	doc := struct {
		uuid.Identify
		Name string
	}{
		Identify: uuid.IdentifyV7(),
		Name:     "test",
	}
	assert.NotEmpty(t, doc.GetUUID())
	assert.False(t, doc.UUID.Timestamp().IsZero())
}
