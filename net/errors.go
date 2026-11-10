package net

import "errors"

var (
	// ErrAddressEmpty is returned when an empty address is provided.
	ErrAddressEmpty = errors.New("net: address is empty")
	// ErrAddressInvalid is returned when the address is not a valid FQDN.
	ErrAddressInvalid = errors.New("net: invalid address")
	// ErrFetchFailed is returned when a well-known resource could not be
	// fetched for a reason retrying will not cure (a definitive non-2xx
	// response, malformed content, or invalid input).
	ErrFetchFailed = errors.New("net: failed to fetch resource")
	// ErrUnavailable is returned when a well-known resource could not be
	// reached at all (network failure, 429, or a 5xx response): a
	// transient condition the caller may retry. Servers verifying
	// request tokens respond 503, not 401, when the issuer's key is
	// unavailable.
	ErrUnavailable = errors.New("net: service unavailable")
	// ErrNoContent is returned when a well-known resource exists but has no
	// content to share (HTTP 204), e.g. a /who endpoint whose owner does not
	// publish identity details.
	ErrNoContent = errors.New("net: no content")
	// ErrPending is returned when a who request was accepted for deferred
	// disclosure (HTTP 202): the owner may deliver its party envelope to the
	// requester's inbox later.
	ErrPending = errors.New("net: request accepted, response deferred")
	// ErrTokenInvalid is returned when a request token is missing or fails
	// verification.
	ErrTokenInvalid = errors.New("net: invalid request token")
	// ErrTokenExpired is returned when a request token fails the freshness
	// check.
	ErrTokenExpired = errors.New("net: request token expired")
	// ErrVerifyFailed is returned when verification fails.
	ErrVerifyFailed = errors.New("net: verification failed")
	// ErrUnknownAuthority is returned when no trusted authority endorsement is found.
	ErrUnknownAuthority = errors.New("net: endorser is not a recognised authority")
	// ErrNotVerified is returned when a sender's endorsement is valid but
	// carries no confirmed verifier and the caller required identity
	// verification.
	ErrNotVerified = errors.New("net: sender identity is not verified")
	// ErrSignatureExpired is returned when an authority countersignature
	// verifies but its exp claim is in the past.
	ErrSignatureExpired = errors.New("net: signature expired")
	// ErrPartyMissing is returned when an identity envelope has no org.Party.
	ErrPartyMissing = errors.New("net: /who response did not contain a party document")
	// ErrInboxRejected is returned when an inbox endpoint rejects an envelope.
	ErrInboxRejected = errors.New("net: inbox rejected envelope")
)
