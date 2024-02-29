package uuid_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
)

type uuidTestStruct struct {
	UUID *uuid.UUID
}

func (ut *uuidTestStruct) Validate() error {
	return validation.ValidateStruct(ut,
		validation.Field(&ut.UUID, uuid.IsV1),
	)
}

func TestUUIDValidation(t *testing.T) {
	u1 := uuid.MakeV1()
	u4 := uuid.MakeV4()
	assert.NoError(t, validation.Validate(u1, uuid.IsV1), "should accept UUIDv1")
	assert.NoError(t, validation.Validate(nil, uuid.IsV1), "should ignore nil")
	assert.NoError(t, validation.Validate("", uuid.IsV1), "should ignore empty string")
	assert.NoError(t, validation.Validate(u1.String(), uuid.IsV1), "should accept string")
	assert.NoError(t, validation.Validate(u4, uuid.IsV4))
	assert.NoError(t, validation.Validate(nil, uuid.IsV4), "should ignore nil")
	assert.NoError(t, validation.Validate("", uuid.IsV4), "should ignore empty string")
	err := validation.Validate(u1, uuid.IsV4)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid version")
	}
	err = validation.Validate(uuid.UUID{}, uuid.IsV1)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid version")
	}
	err = validation.Validate(u4, uuid.IsV1)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "invalid version")
	}
	assert.NoError(t, validation.Validate(u1, uuid.Within(1*time.Second)))
	time.Sleep(11 * time.Millisecond)
	err = validation.Validate(u1, uuid.Within(10*time.Millisecond))
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "timestamp is outside acceptable range")
	}
	assert.NoError(t, validation.Validate(u1, uuid.IsNotZero))
	err = validation.Validate(uuid.UUID{}.String(), uuid.IsNotZero)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "is zero")
	}

	pu1 := uuid.NewV1()
	err = validation.Validate(pu1, uuid.IsV1)
	assert.NoError(t, err, "failed to validate pointer")

	sample := new(uuidTestStruct)
	sample.UUID = uuid.NewV1()
	err = sample.Validate()
	assert.NoError(t, err)

	// Additional checks for other UUID versions
	u3 := uuid.MakeV3(u1, []byte("test"))
	u5 := uuid.MakeV5(u1, []byte("test"))
	assert.NoError(t, validation.Validate(u3, uuid.IsV3))
	assert.NoError(t, validation.Validate(u5, uuid.IsV5))
	assert.Error(t, validation.Validate(u1, uuid.IsV3))
	assert.Error(t, validation.Validate(u1, uuid.IsV5))
}
