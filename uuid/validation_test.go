package uuid_test

import (
	"testing"
	"time"

	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUUIDValidation(t *testing.T) {
	base := uuid.UUID("03907310-8daa-11eb-8dcd-0242ac130003")
	tests := []struct {
		name string
		uuid any
		rule uuid.VersionTest
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
			err := tt.rule.Validate(tt.uuid)
			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.err)
			}
		})
	}

	// Timestamp within tests
	id := uuid.V1()
	assert.NoError(t, uuid.Within(1*time.Second).Validate(id))
	time.Sleep(12 * time.Millisecond)
	err := uuid.Within(10 * time.Millisecond).Validate(id)
	assert.ErrorContains(t, err, "timestamp is outside acceptable range")

	id = uuid.V6()
	assert.NoError(t, uuid.Within(1*time.Second).Validate(id))
	time.Sleep(20 * time.Millisecond)
	err = uuid.Within(10 * time.Millisecond).Validate(id)
	assert.ErrorContains(t, err, "timestamp is outside acceptable range")

	id = uuid.V7()
	assert.NoError(t, uuid.Within(1*time.Second).Validate(id))
	time.Sleep(12 * time.Millisecond)
	err = uuid.Within(10 * time.Millisecond).Validate(id)
	assert.ErrorContains(t, err, "timestamp is outside acceptable range")
}

type uuidRulesDoc struct {
	ID uuid.UUID `json:"id"`
}

func TestUUIDRules(t *testing.T) {
	t.Run("valid passes", func(t *testing.T) {
		faults := rules.Validate(uuid.V4())
		assert.NoError(t, faults)
	})

	t.Run("empty passes (AssertIfPresent)", func(t *testing.T) {
		faults := rules.Validate(uuid.UUID(""))
		assert.NoError(t, faults)
	})

	t.Run("invalid UUID produces fault", func(t *testing.T) {
		faults := rules.Validate(uuid.UUID("not-a-uuid"))
		require.NotNil(t, faults)
		assert.Equal(t, rules.Code("GOBL-UUID-UUID-01"), faults.First().Code())
		assert.Equal(t, "invalid UUID", faults.First().Message())
	})

	t.Run("struct field invalid UUID produces fault at path", func(t *testing.T) {
		doc := &uuidRulesDoc{ID: uuid.UUID("not-a-uuid")}
		faults := rules.Validate(doc)
		require.NotNil(t, faults)
		assert.True(t, faults.HasPath("$.id"))
		assert.Equal(t, rules.Code("GOBL-UUID-UUID-01"), faults.First().Code())
	})

	t.Run("struct field valid UUID passes", func(t *testing.T) {
		doc := &uuidRulesDoc{ID: uuid.V4()}
		faults := rules.Validate(doc)
		assert.NoError(t, faults)
	})

	t.Run("struct field empty UUID passes", func(t *testing.T) {
		doc := &uuidRulesDoc{ID: uuid.UUID("")}
		faults := rules.Validate(doc)
		assert.NoError(t, faults)
	})
}

func TestVersionTestCheck(t *testing.T) {
	tests := []struct {
		name string
		rule rules.Test
		uuid uuid.UUID
		want bool
	}{
		{"Valid accepts v4", uuid.Valid, uuid.V4(), true},
		{"Valid accepts v7", uuid.Valid, uuid.V7(), true},
		{"Valid rejects bad format", uuid.Valid, uuid.UUID("bad"), false},
		{"Valid ignores empty", uuid.Valid, uuid.UUID(""), true},
		{"IsV4 accepts v4", uuid.IsV4, uuid.V4(), true},
		{"IsV4 rejects v7", uuid.IsV4, uuid.V7(), false},
		{"IsV7 accepts v7", uuid.IsV7, uuid.V7(), true},
		{"IsV1 accepts v1", uuid.IsV1, uuid.V1(), true},
		{"HasTimestamp accepts v7", uuid.HasTimestamp, uuid.V7(), true},
		{"HasTimestamp rejects v4", uuid.HasTimestamp, uuid.V4(), false},
		{"Timeless accepts v4", uuid.Timeless, uuid.V4(), true},
		{"Timeless rejects v7", uuid.Timeless, uuid.V7(), false},
		{"IsNotZero rejects zero", uuid.IsNotZero, uuid.Zero, false},
		{"IsNotZero accepts v4", uuid.IsNotZero, uuid.V4(), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.rule.Check(tt.uuid))
		})
	}
}
