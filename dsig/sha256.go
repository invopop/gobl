package dsig

import (
	"crypto/sha256"
	"encoding/hex"
)

// NewSHA256Digest creates a SHA256 digest object from the provided byte array.
// We assume the data has already been through a canonicalization (c14n)
// process.
func NewSHA256Digest(data []byte) *Digest {
	sum := sha256.Sum256(data)
	return &Digest{
		Algorithm: DigestSHA256,
		Value:     hex.EncodeToString(sum[:]),
	}
}
