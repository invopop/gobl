// Package net implements GOBL Net identity discovery, authentication,
// envelope delivery, and remote signature verification.
//
// EXPERIMENTAL: GOBL Net is under active development. Its API and wire
// protocol may change without notice and carry no stability guarantee yet.
//
// # Relationship to the specification
//
// The protocol itself is defined in README.md alongside this package, and
// that document is the normative one: it is implementation-neutral and names
// no Go API. This package is a reference implementation of it, in the same
// sense that the c14n package is a reference implementation of GOBL's
// canonical JSON — a library in another language can be written from the
// specification alone, and should be.
//
// What follows maps the specification onto this package. Section numbers
// refer to README.md. Where the two disagree the specification is right and
// this package has a bug.
//
//	§3  addresses            Address, ParseAddress
//	§4  keys                 Client.FetchKey, Client.FlushKeyCache
//	§5  signatures           (gobl/head)
//	§5.1 party identities    Client.VerifyParty, Client.Who
//	§5.2 deliveries          Client.VerifyDelivery
//	§6  request tokens       NewToken, Client.VerifyToken,
//	                         Client.VerifyAuthorization
//	§7  sender policy        WithAuthorities, WithSandbox
//	§7.2 sender status       Client.VerifyAuthority, Endorsement,
//	                         Client.VerifySender
//	§8.2 /who                Client.Who
//	§8.3 /inbox              Client.Send
//	§9  limits               HTTPFetcher, WithKeyCacheTTL
//	§11 failure handling     ErrUnavailable, ErrInboxRejected, and the
//	                         other sentinels below
//
// # Sender status and the Endorsement type
//
// The specification deliberately has no "endorsement" on the wire. It has
// countersignatures — a registration countersignature from an authority, a
// verification countersignature from a verifier — and a sender status that a
// receiver concludes by evaluating them against its own trust list
// (§7.2). Status is never transmitted: two receivers with different trust
// lists can read the same envelope and reach different conclusions, and both
// are correct.
//
// Endorsement is this package's representation of that conclusion, which is
// why it is a return value and never something you construct or send:
//
//	end, err := client.VerifyAuthority(ctx, env)
//	// end.Authority  which trusted authority the conclusion rests on
//	// end.Verifier   set only when verification was confirmed
//	// end.Verified() whether the subject is verified rather than
//	//                merely registered
//
// Verified reports a boolean because the status is one: any verifier the
// authority names satisfies it. A receiver that wants to distinguish between
// verifiers must read end.Verifier and decide for itself which addresses it
// accepts.
//
// Unregistered, registered and verified are names for conclusions, not
// values. Nothing transmits them and no field carries them, so there is no
// status to parse off the wire and none to put back on it.
//
// Client.VerifySender combines the /who lookup with this evaluation, and is
// usually what an application wants:
//
//	party, err := client.VerifySender(ctx, addr, requireVerified)
//
// # Clients
//
// NewClient returns a client using HTTPFetcher, the default live authority
// list, and a five-minute public-key cache. Configure it with:
//
//	WithIdentity      address and key used to mint request tokens
//	WithAuthorities   replace the trusted registration authorities
//	WithSandbox       replace them with the sandbox authorities
//	WithKeyCacheTTL   change the key cache lifetime; zero disables it
//	WithFetcher       replace the HTTP transport
//
// A client without an identity can still verify and fetch; it cannot mint
// request tokens, so it cannot call /who or deliver to an inbox.
//
// The default HTTPFetcher has a 10-second request timeout, a five-second dial
// timeout, and the 1 MiB response limit of §9. Its dialer refuses hosts
// resolving to loopback, private, link-local, multicast or unspecified
// addresses, as §9 requires. A custom Fetcher takes that obligation on:
// implement equivalent protection, or confine it to tests and local services.
//
// Keys are cached by URL and only on success, with a maximum of 1,024
// entries. FlushKeyCache forces subsequent verification to fetch again, which
// is the incident-response lever §12 asks for.
//
// # Errors
//
// Classify these with errors.Is rather than by comparison; they are wrapped
// with context. The three classes of §11 map onto them as follows.
//
// Temporary — retry with backoff:
//
//	ErrUnavailable      retryable network, 429, or 5xx failure
//
// Definitive — do not retry without changing the request:
//
//	ErrAddressEmpty     address is empty
//	ErrAddressInvalid   address is not a valid FQDN
//	ErrFetchFailed      permanent fetch, content, or configuration failure
//	ErrTokenInvalid     request token is malformed or invalid
//	ErrTokenExpired     request token failed freshness checks
//	ErrVerifyFailed     envelope or signature verification failed
//	ErrUnknownAuthority no trusted authority countersignature was found
//	ErrNotVerified      registration is valid but verification was required
//	ErrSignatureExpired an authority countersignature has expired
//	ErrPartyMissing     identity envelope has no org.Party
//	ErrInboxRejected    inbox returned a definitive 4xx
//
// Absent — the peer answered, and the answer is that there is nothing:
//
//	ErrNoContent        endpoint returned 204
//	ErrPending          endpoint returned 202 where a body was expected
//
// ErrUnavailable is the one to be careful with. It means a check could not be
// performed, which §11 distinguishes from a check that failed: surfacing it as
// a rejection turns somebody else's outage into a refusal the caller cannot
// tell apart from a real one.
package net

const (
	// Scheme is the URI scheme labelling a GOBL Net address in
	// multi-scheme contexts (org.Endpoint lists, the envelope
	// header's from/to routing fields), e.g. "gobl:acme.example.com".
	// Signed iss/aud/verifier claims carry bare addresses instead.
	Scheme = "gobl"

	// WellKnownPath is the base path for GOBL Net well-known URLs.
	WellKnownPath = "/.well-known/gobl"
	// KeysPath is the base of the per-key endpoint; the full path for a
	// single key is KeysPath + "/" + kid.
	KeysPath = WellKnownPath + "/keys"
	// WhoPath is the well-known path serving the signed Party envelope.
	WhoPath = WellKnownPath + "/who"
	// InboxPath is the well-known path accepting envelope deliveries.
	InboxPath = WellKnownPath + "/inbox"
	// JWKSPath is the conventional bulk JWK Set endpoint for generic tooling.
	// GOBL Net verification uses the per-key endpoint instead.
	JWKSPath = "/.well-known/jwks.json"
)

// KeyPath returns the well-known path serving a single public key by
// its key ID. Use this to construct lookup URLs.
func KeyPath(kid string) string {
	return KeysPath + "/" + kid
}
