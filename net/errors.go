package net

import "errors"

var (
	// ErrAddressEmpty is returned when an empty address is provided.
	ErrAddressEmpty = errors.New("net: address is empty")
	// ErrAddressInvalid is returned when the address is not a valid FQDN.
	ErrAddressInvalid = errors.New("net: invalid address")
	// ErrFetchFailed is returned when a well-known resource could not be fetched.
	ErrFetchFailed = errors.New("net: failed to fetch resource")
	// ErrNoContent is returned when a well-known resource exists but has no
	// content to share (HTTP 204), e.g. a /who endpoint whose owner does not
	// publish identity details.
	ErrNoContent = errors.New("net: no content")
	// ErrVerifyFailed is returned when verification fails.
	ErrVerifyFailed = errors.New("net: verification failed")
	// ErrUnknownAuthority is returned when a /who envelope is signed by an
	// address not in the Authorities list.
	ErrUnknownAuthority = errors.New("net: endorser is not a recognised authority")
	// ErrScopeInsufficient is returned when an authority countersignature
	// verifies but its scope claim does not meet the required minimum.
	ErrScopeInsufficient = errors.New("net: authority scope insufficient")
	// ErrSignatureExpired is returned when an authority countersignature
	// verifies but its exp claim is in the past.
	ErrSignatureExpired = errors.New("net: signature expired")
	// ErrPartyMissing is returned when a /who response does not contain an
	// org.Party document.
	ErrPartyMissing = errors.New("net: /who response did not contain a party document")
	// ErrInboxRejected is returned when an inbox endpoint rejects an envelope.
	ErrInboxRejected = errors.New("net: inbox rejected envelope")
)
