package uuid

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/invopop/validation"
)

var (
	// IsV1 confirms the UUID is version 1
	IsV1 = versionRule{version: 1}
	// IsV3 confirms the UUID is version 3
	IsV3 = versionRule{version: 3}
	// IsV4 confirms the UUID is version 4
	IsV4 = versionRule{version: 4}
	// IsV5 confirms the UUID is version 5
	IsV5 = versionRule{version: 5}
	// IsNotZero confirms the UUID is not zero
	IsNotZero = versionRule{notZero: true}
)

type versionRule struct {
	version uuid.Version
	notZero bool
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
	if value == nil {
		return nil
	}
	var id UUID
	switch v := value.(type) {
	case UUID:
		id = v
	case *UUID:
		id = *v
	case string:
		if v == "" {
			return nil
		}
		var err error
		id, err = Parse(v)
		if err != nil {
			return err
		}
	default:
		return errors.New("not a UUID")
	}
	if r.notZero {
		if id.IsZero() {
			return errors.New("is zero")
		}
		return nil
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
