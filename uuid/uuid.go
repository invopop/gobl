// Package uuid provides a wrapper for handling UUID codes.
package uuid

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/jsonschema"
)

func init() {
	schema.Register(schema.GOBL.Add("uuid"), UUID{})
}

// UUID defines our wrapper for dealing with UUIDs
type UUID struct {
	uuid.UUID
}

// MakeV1 generates a version 1 UUID.
func MakeV1() UUID {
	return UUID{uuid.Must(uuid.NewUUID())}
}

// MakeV3 generates a new UUIDv3 using the provided namespace and data. The behaviour is
// deterministic, that is, the same inputs will always generate the same UUID. This is handy to
// transform any other types of IDs into UUIDs, among other uses.
//
// In UUIDv3, the data is hashed using MD5 which is a performant algorithm, but it can be broken. If
// you need to ensure an attacker can't reverse the UUID to the original data used to generate it,
// use UUIDv5 instead.
func MakeV3(space UUID, data []byte) UUID {
	return UUID{uuid.NewMD5(space.UUID, data)}
}

// MakeV4 generates a new completely random UUIDv4.
func MakeV4() UUID {
	return UUID{uuid.Must(uuid.NewRandom())}
}

// MakeV5 generates a new UUIDv5 using the provided namespace and data. The behaviour is
// deterministic, that is, the same inputs will always generate the same UUID. This is handy to
// transform any other types of IDs into UUIDs, among other uses.
//
// In UUIDv5, the data is hashed using SHA1, a secure algorithm but slower than MD5. If you
// don't need the security of SHA1 and perfomance is a concern, use UUIDv3 instead.
func MakeV5(space UUID, data []byte) UUID {
	return UUID{uuid.NewSHA1(space.UUID, data)}
}

// NewV1 generates a version 1 UUID.
func NewV1() *UUID {
	u := MakeV1()
	return &u
}

// NewV3 creates a new MD5 UUID using the provided namespace and data. See MakeV3 for more details.
func NewV3(space UUID, data []byte) *UUID {
	u := MakeV3(space, data)
	return &u
}

// NewV4 creates a pointer a new completely random UUIDv4.
func NewV4() *UUID {
	u := MakeV4()
	return &u
}

// NewV5 creates a new SHA1 UUID using the provided namespace and data. See MakeV5 for more details.
func NewV5(space UUID, data []byte) *UUID {
	u := MakeV5(space, data)
	return &u
}

// Timestamp extracts the time.
// Anything other than a v1 UUID will provide zero time without an error,
// so ensure your error checks are performed previously.
func (u UUID) Timestamp() time.Time {
	if u.UUID.Version() != 1 {
		return time.Time{}
	}
	return time.Unix(u.UUID.Time().UnixTime())
}

// IsZero returns true if the UUID is all zeros.
func (u *UUID) IsZero() bool {
	if u == nil {
		return true
	}
	for _, v := range u.UUID {
		if v != 0 {
			return false
		}
	}
	return true
}

// Parse decodes s into a UUID or provides an error.
func Parse(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	return UUID{id}, err
}

// ShouldParse will return a UUID if the string is valid, otherwise it will
// provide a zero UUID.
func ShouldParse(s string) UUID {
	id, err := Parse(s)
	if err != nil {
		return UUID{}
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
// indending to generate V1 UUIDs to get a node ID that won't change.
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

// Normalize looks at the provided UUID and tries to return a consistent
// value, which may be nil. Only works with pointers to UUID.
func Normalize(u *UUID) *UUID {
	if u == nil {
		return nil
	}
	if u.IsZero() {
		return nil
	}
	return u
}

// JSONSchema returns the jsonschema schema object for the UUID.
func (UUID) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Format:      "uuid",
		Title:       "UUID",
		Description: "Universally Unique Identifier. We only recommend using versions 1 and 4 within GOBL.",
	}
}
