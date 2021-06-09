package uuid

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// UUID defines our wrapper for dealing with UUIDs
type UUID struct {
	uuid.UUID
}

// NewV1 generates a version 1 UUID.
func NewV1() UUID {
	return UUID{uuid.Must(uuid.NewUUID())}
}

// NewV4 generates a new completely random UUIDv4.
func NewV4() UUID {
	return UUID{uuid.Must(uuid.NewRandom())}
}

// Timestamp extracts the time.
// Anything other than a v1 UUID will provide zero time without an error,
// so ensure your error checks are performed previously.
func (u UUID) Timestamp() time.Time {
	if u.Version() != 1 {
		return time.Time{}
	}
	return time.Unix(u.Time().UnixTime())
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
