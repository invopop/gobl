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
	if err := validation.Validate(u1, uuid.IsV1); err != nil {
		t.Errorf("did not expect an error: %v", err)
	}
	if err := validation.Validate(u4, uuid.IsV4); err != nil {
		t.Errorf("did not expect an error: %v", err)
	}
	if err := validation.Validate(u1, uuid.IsV4); err == nil {
		t.Errorf("expected an error")
	}
	if err := validation.Validate(u4, uuid.IsV1); err == nil {
		t.Errorf("expected an error")
	}
	if err := validation.Validate(u1, uuid.Within(1*time.Second)); err != nil {
		t.Errorf("did not expect an error so soon, got: %v", err.Error())
	}
	time.Sleep(11 * time.Millisecond)
	if err := validation.Validate(u1, uuid.Within(10*time.Millisecond)); err == nil {
		t.Errorf("expected an error")
	}

	pu1 := uuid.NewV1()
	err := validation.Validate(pu1, uuid.IsV1)
	assert.NoError(t, err, "failed to validate pointer")

	sample := new(uuidTestStruct)
	sample.UUID = uuid.NewV1()
	err = sample.Validate()
	assert.NoError(t, err)
}
