package uuid_test

import (
	"testing"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/uuid"
)

func TestUUIDValidation(t *testing.T) {
	u1 := uuid.NewV1()
	u4 := uuid.NewV4()
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
}
