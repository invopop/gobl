package uuid

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/invopop/validation"
)

var (
	// Valid confirms the UUID is valid
	Valid = versionRule{}
	// IsV1 confirms the UUID is version 1
	IsV1 = versionRule{version: 1}
	// IsV3 confirms the UUID is version 3
	IsV3 = versionRule{version: 3}
	// IsV4 confirms the UUID is version 4
	IsV4 = versionRule{version: 4}
	// IsV5 confirms the UUID is version 5
	IsV5 = versionRule{version: 5}
	// IsV6 confirms the UUID is version 6
	IsV6 = versionRule{version: 6}
	// IsV7 confirms the UUID is version 7
	IsV7 = versionRule{version: 7}
	// HasTimestamp confirms the UUID is based on a timestamp version
	HasTimestamp = versionRule{hasTimestamp: true}
	// Timeless confirms the UUID is not based on a timestamp version
	Timeless = versionRule{timeless: true}
	// IsNotZero confirms the UUID is not zero
	IsNotZero = versionRule{notZero: true}
)

type versionRule struct {
	version      uuid.Version
	hasTimestamp bool
	timeless     bool
	notZero      bool
	ttl          time.Duration
}

const (
	maxFutureDuration = -10 * time.Second
)

// Within is a validation method that can be used to determine if the UUID
// is version 1, 6, or 7 and contains a timestamp that is greater than the
// current time minus the ttl. A tolerance is allowed for future timestamps.
func Within(ttl time.Duration) validation.Rule {
	return versionRule{
		hasTimestamp: true,
		ttl:          ttl,
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
		var err error
		id, err = Parse(v)
		if err != nil {
			return err
		}
	default:
		return errors.New("not a UUID")
	}
	// always ignore empty
	if id == Empty {
		return nil
	}
	if r.notZero {
		if id.IsZero() {
			return errors.New("is zero")
		}
		return nil
	}
	if r.version != 0 {
		if id.Version() != Version(r.version) {
			return errors.New("invalid version")
		}
	}
	if r.hasTimestamp {
		switch id.Version() {
		case 1, 6, 7:
			// good
		default:
			return errors.New("not timestamped")
		}
	}
	if r.timeless {
		switch id.Version() {
		case 3, 4, 5:
			// good
		default:
			return errors.New("has timestamp")
		}
	}
	if r.ttl == 0 {
		// don't check empty duration
		return nil
	}

	// check the time range and allow defined ttl margin
	tn := time.Now()
	ti := id.Timestamp()
	d := tn.Sub(ti)
	if d < maxFutureDuration {
		return errors.New("timestamp cannot be in the future")
	}
	if d > r.ttl {
		return errors.New("timestamp is outside acceptable range")
	}

	return nil
}
