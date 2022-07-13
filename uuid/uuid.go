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

// MakeV4 generates a new completely random UUIDv4.
func MakeV4() UUID {
	return UUID{uuid.Must(uuid.NewRandom())}
}

// NewV1 generates a version 1 UUID.
func NewV1() *UUID {
	u := MakeV1()
	return &u
}

// NewV4 creates a pointer a new completely random UUIDv4.
func NewV4() *UUID {
	u := MakeV4()
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
func (u UUID) IsZero() bool {
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

// JSONSchema returns the jsonschema schema object for the UUID.
func (UUID) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "string",
		Format:      "uuid",
		Title:       "UUID",
		Description: "Universally Unique Identifier. We only recommend using versions 1 and 4 within GoBL.",
	}
}
