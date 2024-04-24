// Package uuid provides a wrapper for handling UUID codes.
package uuid

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// UUID defines a string wrapper for dealing with UUIDs using the google uuid
// package for parsing and specific method support. This implementation has
// been optimized for convenience and JSON conversion as opposed to performance.
//
// Note that this package is not registered inside its own schema. We
// instead rely on the `"format"` parameter already included in JSON Schema.
type UUID string

// Version represents the version number of the UUID
type Version byte

// Constants for empty and zero states.
const (
	Empty UUID = ""
	Zero  UUID = "00000000-0000-0000-0000-000000000000"
)

// V1 generates a version 1 UUID. We strongly recommend using V7 now as an alternative
// as it provides the same functionality with a more secure random node ID.
func V1() UUID {
	return UUID(uuid.Must(uuid.NewUUID()).String())
}

// V3 generates a new UUIDv3 using the provided namespace and data. The behavior is
// deterministic, that is, the same inputs will always generate the same UUID. This is handy to
// transform any other types of IDs into UUIDs, among other uses.
//
// In UUIDv3, the data is hashed using MD5 which is a performant algorithm, but is subject to
// collision attacks and other vulnerabilities. If security is a concern, use UUIDv5 instead.
func V3(space UUID, data []byte) UUID {
	return UUID(uuid.NewMD5(parse(space), data).String())
}

// V4 generates a new completely random UUIDv4.
func V4() UUID {
	return UUID(uuid.Must(uuid.NewRandom()).String())
}

// V5 generates a new UUIDv5 using the provided namespace and data. The behavior is
// deterministic, that is, the same inputs will always generate the same UUID. This is handy to
// transform any other types of IDs into UUIDs, among other uses.
//
// In UUIDv5, the data is hashed using SHA1, a secure algorithm but slower than MD5. If you
// don't need the security of SHA1 and performance is a concern, use UUIDv3 instead.
func V5(space UUID, data []byte) UUID {
	return UUID(uuid.NewSHA1(parse(space), data).String())
}

// V6 generates a version 6 UUID, a drop-in replaced for V1 UUIDs that uses random
// data instead of node. It maintains a similar structure to V1 UUIDs with the
// timestamp, so no ordering is maintained.
func V6() UUID {
	return UUID(uuid.Must(uuid.NewV6()).String())
}

// V7 generates a new UUIDv7, a replacement for V1 or V6 UUIDs which
// combines the Unix timestamp with millisecond precision and random data.
// An important difference with other versions is that order is maintained,
// making this a great option for primary keys in databases.
func V7() UUID {
	return UUID(uuid.Must(uuid.NewV7()).String())
}

// MakeV1 generates a version 1 UUID.
//
// Deprecated: use V1() instead.
func MakeV1() UUID {
	return V1()
}

// MakeV3 generates a new UUIDv3 using the provided namespace and data. The behavior is
// deterministic, that is, the same inputs will always generate the same UUID. This is handy to
// transform any other types of IDs into UUIDs, among other uses.
//
// In UUIDv3, the data is hashed using MD5 which is a performant algorithm, but is subject to
// collision attacks and other vulnerabilities. If security is a concern, use UUIDv5 instead.
//
// Deprecated: use V3() instead.
func MakeV3(space UUID, data []byte) UUID {
	return V3(space, data)
}

// MakeV4 generates a new completely random UUIDv4.
//
// Deprecated: use V4() instead.
func MakeV4() UUID {
	return V4()
}

// MakeV5 generates a new UUIDv5 using the provided namespace and data. The behavior is
// deterministic, that is, the same inputs will always generate the same UUID. This is handy to
// transform any other types of IDs into UUIDs, among other uses.
//
// In UUIDv5, the data is hashed using SHA1, a secure algorithm but slower than MD5. If you
// don't need the security of SHA1 and performance is a concern, use UUIDv3 instead.
//
// Deprecated: use V5() instead.
func MakeV5(space UUID, data []byte) UUID {
	return V5(space, data)
}

// NewV1 generates a version 1 UUID.
//
// Deprecated: Use V1() instead.
func NewV1() *UUID {
	u := MakeV1()
	return &u
}

// NewV3 creates a new MD5 UUID using the provided namespace and data. See MakeV3 for more details.
//
// Deprecated: Use V3() instead.
func NewV3(space UUID, data []byte) *UUID {
	u := MakeV3(space, data)
	return &u
}

// NewV4 creates a pointer a new completely random UUIDv4.
//
// Deprecated: Use V4() instead.
func NewV4() *UUID {
	u := MakeV4()
	return &u
}

// NewV5 creates a new SHA1 UUID using the provided namespace and data. See MakeV5 for more details.
//
// Deprecated: Use V5() instead
func NewV5(space UUID, data []byte) *UUID {
	u := MakeV5(space, data)
	return &u
}

// Timestamp extracts the time.
// Anything other than a version 1, 6, or 7 UUID will provide zero time without an error,
// so ensure your error checks are performed previously.
func (u UUID) Timestamp() time.Time {
	id := parse(u)
	switch id.Version() {
	case 1, 6, 7:
		// good
	default:
		return time.Time{}
	}
	return time.Unix(id.Time().UnixTime())
}

// Version returns the version number of the UUID.
func (u UUID) Version() Version {
	return Version(parse(u).Version())
}

// IsZero returns true if the UUID is all zeros or empty.
func (u *UUID) IsZero() bool {
	return u == nil || *u == Empty || *u == Zero
}

// String provides the string representation of the UUID.
func (u UUID) String() string {
	return string(u)
}

// Validate checks to ensure the value is a UUID
func (u UUID) Validate() error {
	return validation.Validate(string(u), Valid)
}

// Parse decodes s into a UUID or provides an error.
func Parse(s string) (UUID, error) {
	if s == "" {
		return Empty, nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return Empty, err
	}
	return UUID(id.String()), err
}

func parse(s UUID) uuid.UUID {
	return uuid.MustParse(string(s))
}

// ShouldParse will return a UUID if the string is valid, otherwise it will
// provide a zero UUID.
func ShouldParse(s string) UUID {
	id, err := Parse(s)
	if err != nil {
		return Zero
	}
	return id
}

// MustParse will panic if the UUID does not look good.
func MustParse(s string) UUID {
	id, err := Parse(s)
	if err != nil {
		panic(err.Error())
	}
	return id
}

// SetRandomNodeID is used to generate a random host ID to be used in V1 UUIDs
// instead of the MAC address. This is stored in the uuid library as a global
// constant, so should be called just once when starting the application if you're
// intending to generate V1 UUIDs to get a node ID that won't change.
func SetRandomNodeID() {
	id := make([]byte, 6)
	_, err := rand.Read(id)
	if err != nil {
		panic(err.Error())
	}
	uuid.SetNodeID(id)
}

// NodeID returns the hex representation of the current host bytes
func NodeID() string {
	return fmt.Sprintf("%x", uuid.NodeID())
}

// Normalize will ensure that zero value UUIDs will be empty strings
// instead of zeros.
func Normalize(u *UUID) {
	if u == nil {
		return
	}
	if u.IsZero() {
		*u = ""
	}
}

// UnmarshalText will ensure the UUID is always a valid UUID when unmarshalling
// and just return an empty value if incorrect.
// TODO: Remove this and instead depend on validation to provide more readable errors.
func (u *UUID) UnmarshalText(txt []byte) error {
	id, err := Parse(string(txt))
	if err != nil {
		return err
	}
	*u = id
	return nil
}

// JSONSchemaExtend ensures the schema contains the additional UUID format
// wherever it is used.
func (UUID) JSONSchemaExtend(s *jsonschema.Schema) {
	s.Format = "uuid"
}
