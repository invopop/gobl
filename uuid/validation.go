package uuid

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/invopop/gobl/rules"
)

var (
	// Valid confirms the UUID is valid
	Valid = VersionTest{desc: "is valid UUID"}
	// IsV1 confirms the UUID is version 1
	IsV1 = VersionTest{desc: "is UUIDv1", version: 1}
	// IsV3 confirms the UUID is version 3
	IsV3 = VersionTest{desc: "is UUIDv3", version: 3}
	// IsV4 confirms the UUID is version 4
	IsV4 = VersionTest{desc: "is UUIDv4", version: 4}
	// IsV5 confirms the UUID is version 5
	IsV5 = VersionTest{desc: "is UUIDv5", version: 5}
	// IsV6 confirms the UUID is version 6
	IsV6 = VersionTest{desc: "is UUIDv6", version: 6}
	// IsV7 confirms the UUID is version 7
	IsV7 = VersionTest{desc: "is UUIDv7", version: 7}
	// HasTimestamp confirms the UUID is based on a timestamp version
	HasTimestamp = VersionTest{desc: "has timestamp", hasTimestamp: true}
	// Timeless confirms the UUID is not based on a timestamp version
	Timeless = VersionTest{desc: "is timeless", timeless: true}
	// IsNotZero confirms the UUID is not zero
	IsNotZero = VersionTest{desc: "is not zero", notZero: true}
)

func uuidRules() *rules.Set {
	return rules.For(UUID(""),
		rules.AssertIfPresent("01", "invalid UUID", Valid),
	)
}

// VersionTest provides a validation instance that can be used for checking
// UUIDs.
type VersionTest struct {
	desc         string
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
func Within(ttl time.Duration) VersionTest {
	return VersionTest{
		desc:         "is within acceptable time range",
		hasTimestamp: true,
		ttl:          ttl,
	}
}

// Check provides a boolean response.
func (r VersionTest) Check(value any) bool {
	return r.Validate(value) == nil
}

// String provides a string representation of the rule.
func (r VersionTest) String() string {
	return r.desc
}

// Validate checks the value against the rule and returns an error if it does not pass.
func (r VersionTest) Validate(value any) error {
	if value == nil {
		return nil
	}
	var id UUID
	switch v := value.(type) {
	case UUID:
		var err error
		id, err = Parse(string(v))
		if err != nil {
			return err
		}
	case *UUID:
		var err error
		id, err = Parse(string(*v))
		if err != nil {
			return err
		}
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
