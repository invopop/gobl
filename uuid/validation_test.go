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
	base := uuid.UUID("03907310-8daa-11eb-8dcd-0242ac130003")
	tests := []struct {
		name string
		uuid any
		rule validation.Rule
		err  string
	}{
		{
			name: "valid v1",
			uuid: uuid.V1(),
			rule: uuid.IsV1,
		},
		{
			name: "valid v1 pointer",
			uuid: &base,
			rule: uuid.IsV1,
		},
		{
			name: "not uuid v1",
			uuid: uuid.V4(),
			rule: uuid.IsV1,
			err:  "invalid version",
		},
		{
			name: "ignore nil",
			uuid: nil,
			rule: uuid.IsV1,
		},
		{
			name: "ignore empty",
			uuid: "",
			rule: uuid.IsV1,
		},
		{
			name: "validate string",
			uuid: uuid.V1().String(),
			rule: uuid.IsV1,
		},
		{
			name: "reject invalid string",
			uuid: uuid.V4().String(),
			rule: uuid.IsV1,
			err:  "invalid version",
		},
		{
			name: "valid v4",
			uuid: uuid.V4(),
			rule: uuid.IsV4,
		},
		{
			name: "not uuid v4",
			uuid: uuid.V1(),
			rule: uuid.IsV4,
			err:  "invalid version",
		},
		{
			name: "valid v3",
			uuid: uuid.V3(base, []byte("test")),
			rule: uuid.IsV3,
		},
		{
			name: "invalid v3",
			uuid: uuid.V5(base, []byte("test")),
			rule: uuid.IsV3,
			err:  "invalid version",
		},
		{
			name: "valid v5",
			uuid: uuid.V5(base, []byte("test")),
			rule: uuid.IsV5,
		},
		{
			name: "valid v7",
			uuid: uuid.V7(),
			rule: uuid.IsV7,
		},
		{
			name: "not uuid v7",
			uuid: uuid.V1(),
			rule: uuid.IsV7,
			err:  "invalid version",
		},
		{
			name: "has timestamp v1",
			uuid: uuid.V1(),
			rule: uuid.HasTimestamp,
		},
		{
			name: "has timestamp v6",
			uuid: uuid.V6(),
			rule: uuid.HasTimestamp,
		},
		{
			name: "has timestamp v7",
			uuid: uuid.V7(),
			rule: uuid.HasTimestamp,
		},
		{
			name: "no timestamp v4",
			uuid: uuid.V4(),
			rule: uuid.HasTimestamp,
			err:  "not timestamped",
		},
		{
			name: "timeless",
			uuid: uuid.V4(),
			rule: uuid.Timeless,
		},
		{
			name: "not timeless",
			uuid: uuid.V7(),
			rule: uuid.Timeless,
			err:  "has timestamp",
		},
		{
			name: "not zero",
			uuid: uuid.V7(),
			rule: uuid.IsNotZero,
		},
		{
			name: "zero",
			uuid: uuid.UUID("00000000-0000-0000-0000-000000000000"),
			rule: uuid.IsNotZero,
			err:  "is zero",
		},
		{
			name: "zero empty",
			uuid: "",
			rule: uuid.IsNotZero,
		},
		{
			name: "zero empty value",
			uuid: uuid.UUID(""),
			rule: uuid.IsNotZero,
		},
		{
			name: "general good v1",
			uuid: uuid.V1(),
			rule: uuid.Valid,
		},
		{
			name: "general good v4",
			uuid: uuid.V4(),
			rule: uuid.Valid,
		},
		{
			name: "general good v7",
			uuid: uuid.V7(),
			rule: uuid.Valid,
		},
		{
			name: "general empty",
			uuid: "",
			rule: uuid.Valid,
		},
		{
			name: "general bad string",
			uuid: "fooo",
			rule: uuid.Valid,
			err:  "invalid UUID length: 4",
		},
		{
			name: "general bad uuid",
			uuid: uuid.UUID("fooo"),
			rule: uuid.Valid,
			err:  "invalid UUID length: 4",
		},
		{
			name: "other type",
			uuid: 123,
			rule: uuid.Valid,
			err:  "not a UUID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.Validate(tt.uuid, tt.rule)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}
		})
	}

	// Timestamp within tests
	id := uuid.V1()
	assert.NoError(t, validation.Validate(id, uuid.Within(1*time.Second)))
	time.Sleep(12 * time.Millisecond)
	err := validation.Validate(id, uuid.Within(10*time.Millisecond))
	assert.ErrorContains(t, err, "timestamp is outside acceptable range")

	id = uuid.V6()
	assert.NoError(t, validation.Validate(id, uuid.Within(1*time.Second)))
	time.Sleep(20 * time.Millisecond)
	err = validation.Validate(id, uuid.Within(10*time.Millisecond))
	assert.ErrorContains(t, err, "timestamp is outside acceptable range")

	id = uuid.V7()
	assert.NoError(t, validation.Validate(id, uuid.Within(1*time.Second)))
	time.Sleep(12 * time.Millisecond)
	err = validation.Validate(id, uuid.Within(10*time.Millisecond))
	assert.ErrorContains(t, err, "timestamp is outside acceptable range")

	sample := new(uuidTestStruct)
	sample.UUID = uuid.NewV1()
	err = sample.Validate()
	assert.NoError(t, err)

}
