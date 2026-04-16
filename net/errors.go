package net

import "errors"

var (
	// ErrAddressEmpty is returned when an empty address is provided.
	ErrAddressEmpty = errors.New("net: address is empty")
	// ErrAddressInvalid is returned when the address is not a valid FQDN.
	ErrAddressInvalid = errors.New("net: invalid address")
	// ErrFetchFailed is returned when the JWKS could not be fetched.
	ErrFetchFailed = errors.New("net: failed to fetch JWKS")
	// ErrKeyNotFound is returned when the requested key ID is not in the JWKS.
	ErrKeyNotFound = errors.New("net: key not found in JWKS")
	// ErrNoGNHeader is returned when a signature does not contain a gn header.
	ErrNoGNHeader = errors.New("net: signature does not contain a gn header")
	// ErrVerifyFailed is returned when verification fails.
	ErrVerifyFailed = errors.New("net: verification failed")
	// ErrKeySetInvalid is returned when a KeySet is malformed.
	ErrKeySetInvalid = errors.New("net: invalid key set")
)
