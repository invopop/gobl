package uuid

import (
	"errors"
	"time"

	"github.com/alecthomas/jsonschema"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
)

var (
	// IsV1 is used to ensure value is a UUIDv1
	IsV1 = versionRule{version: 1}
	// IsV4 confirms the UUID is version 4
	IsV4 = versionRule{version: 4}
)

type versionRule struct {
	version uuid.Version
	ttl     time.Duration
}

// Within is a validation method that can be used to determine if the UUID
// corresponds do UUIDv1 standard and was timestamped within the acceptable
// time to live from now. If time checks are enabled, future UUIDs will not
// be allowed, this could be a problem.
func Within(ttl time.Duration) validation.Rule {
	return versionRule{
		version: 1,
		ttl:     ttl,
	}
}

func (r versionRule) Validate(value interface{}) error {
	id, ok := value.(UUID)
	if !ok {
		return errors.New("not a UUID")
	}
	if id.Version() != r.version {
		return errors.New("invalid version")
	}
	if r.ttl == 0 {
		// don't check empty duration
		return nil
	}

	// check the time range and allow defined ttl margin
	tn := time.Now()
	ti := id.Timestamp()
	d := tn.Sub(ti)
	if d < 0 {
		return errors.New("timestamp cannot be in the future")
	}
	if d > r.ttl {
		return errors.New("timestamp is outside acceptable range")
	}

	return nil
}

// JSONSchemaType returns the jsonschema type object.
func (UUID) JSONSchemaType() *jsonschema.Type {
	return &jsonschema.Type{
		Type:        "string",
		Format:      "uuid",
		Title:       "UUID",
		Description: "Universally Unique Identifier. We only recommend using versions 1 and 4 within GoBL.",
	}
}
