package net

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/dsig"
	"github.com/invopop/gobl/head"
)

// Authorities is the default list of trusted registration authorities.
var Authorities = []Address{
	"lookup.gobl.org",
}

// SandboxAuthorities is the trust list selected by WithSandbox. Sandbox
// endorsements must not be accepted in live contexts.
var SandboxAuthorities = []Address{
	"lookup.sandbox.gobl.org",
}

// maxEnvelopeSignatures bounds how many signatures VerifyAuthority is
// willing to examine: each candidate can cost a network key fetch, so
// an unbounded envelope would be a fetch-amplification gadget.
// Legitimate envelopes carry a handful of signatures.
const maxEnvelopeSignatures = 32

// maxRegistrationLifetime bounds how long after its iat an authority
// countersignature's exp may fall.
const maxRegistrationLifetime = 90 * 24 * time.Hour

// RegisterAuthority appends addr to the default authority list.
func RegisterAuthority(addr Address) {
	Authorities = append(Authorities, addr)
}

// Endorsement is this package's representation of the sender status
// defined in README.md §7.2: the conclusion a receiver reaches about a
// subject after evaluating the countersignatures its party envelope
// carries against this client's trust list.
//
// There is no "endorsement" on the wire. The envelope carries
// countersignatures; the status is computed from them, and it is never
// transmitted — two clients with different trust lists can evaluate the
// same envelope and reach different conclusions, and both are right. So
// an Endorsement is always a return value, never something to construct
// or send.
type Endorsement struct {
	// Authority is the trusted address whose countersignature
	// endorsed the envelope: the subject is a registered identity.
	Authority Address
	// Verifier is the address the authority named as having performed
	// identity verification (KYC/KYB) of the subject, confirmed by
	// the verifier's own countersignature on the same envelope. Empty
	// when the identity is registered but not verified, including
	// when a named verifier's countersignature is absent, invalid, or
	// expired — verification degrades, registration stands.
	Verifier Address
	// Key is the RFC 7638 thumbprint (see dsig.PublicKey's JKT) of
	// the subject's key that the authority bound the registration to
	// with its cnf claim. The envelope's from and to are confirmed by
	// a signature with this key, and deliveries from the subject must
	// be signed with it too.
	Key string
}

// Verified reports whether the subject is verified rather than merely
// registered (README.md §7.2). Both name conclusions this client reaches,
// not values carried by the envelope.
//
// It is a boolean because the status is one: any verifier the authority
// names satisfies it, and the protocol does not grade them. A receiver
// that wants to distinguish between verifiers must read Verifier itself
// and decide which addresses it accepts.
func (e *Endorsement) Verified() bool {
	return e != nil && e.Verifier != ""
}

// VerifyAuthority verifies countersignatures from the client's trusted
// authorities whose aud binds the envelope to its subject, the address in
// the header's from field. Each needs an iat no later than now, an exp
// that has not passed and is at most 90 days after it, and a cnf naming
// the key of a signature that confirms the envelope's from and to. It returns the
// strongest endorsement found, preferring one with confirmed verifier
// evidence. The verifier must differ from the authority;
// self-referential, missing, or invalid verifier evidence leaves a
// registration-only endorsement valid. It returns ErrUnknownAuthority when no
// trusted candidate exists, ErrSignatureExpired for expired authority evidence,
// ErrUnavailable for transient key lookup failures, or ErrVerifyFailed when
// candidates fail cryptographic verification.
func (c *Client) VerifyAuthority(ctx context.Context, env *gobl.Envelope) (*Endorsement, error) {
	// A malformed envelope may carry signatures without a header;
	// reject rather than let header verification panic.
	if env == nil || env.Head == nil || len(env.Signatures) == 0 {
		return nil, fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
	}
	if len(c.authorities) == 0 {
		return nil, fmt.Errorf("%w: no authorities registered on this client", ErrUnknownAuthority)
	}
	if len(env.Signatures) > maxEnvelopeSignatures {
		return nil, fmt.Errorf("%w: envelope carries %d signatures (max %d)",
			ErrVerifyFailed, len(env.Signatures), maxEnvelopeSignatures)
	}
	subject, to, err := routing(env)
	if err != nil {
		return nil, err
	}
	// Both sides of the lookup are canonicalized so U-Label or
	// trailing-dot forms — configured or signed — compare equal.
	auths := make(map[Address]bool, len(c.authorities))
	for _, a := range c.authorities {
		if canon, err := ParseAddress(string(a)); err == nil {
			auths[canon] = true
		}
	}

	// claimErr records a candidate that verified but carries an
	// expired exp; unavailableErr records a transient key-fetch
	// failure and outranks it. registered holds the first valid
	// registered-only endorsement while the search continues for a
	// confirmed verifier.
	var lastErr, claimErr, unavailableErr error
	var registered *Endorsement
	candidates := make(map[Address]int)
	subjectKeys := &pairKeySet{c: c, env: env, from: subject, to: to}
	for _, sig := range env.Signatures {
		p, err := head.SignedPayload(sig)
		if err != nil {
			continue
		}
		issuer, err := ParseAddress(p.Iss)
		if err != nil {
			continue
		}
		if !auths[issuer] {
			continue
		}
		if aud, err := ParseAddress(p.Aud); err != nil || aud != subject {
			continue
		}
		if candidates[issuer]++; candidates[issuer] > maxIssuerCandidates {
			continue
		}
		// Candidate from a known authority — verify it crypto-wise.
		pub, err := c.FetchKey(ctx, issuer, sig.KeyID())
		if err != nil {
			if errors.Is(err, ErrUnavailable) {
				unavailableErr = err
			} else {
				lastErr = err
			}
			continue
		}
		if err := env.Head.Verify(sig, pub); err != nil {
			lastErr = err
			continue
		}
		// Claims are only trusted after the signature verifies.
		if err := checkRegistration(p); err != nil {
			if errors.Is(err, ErrSignatureExpired) {
				claimErr = err
			} else {
				lastErr = err
			}
			continue
		}
		if err := subjectKeys.require(ctx, p.Confirmation.JKT); err != nil {
			if errors.Is(err, ErrUnavailable) {
				unavailableErr = err
			} else {
				lastErr = err
			}
			continue
		}
		verifier, err := c.confirmVerifier(ctx, env, issuer, subject, p)
		if err != nil {
			// The verifier's key could not be checked (transient):
			// surface the condition rather than silently degrading a
			// genuinely verified endorsement to registered.
			return nil, err
		}
		end := &Endorsement{Authority: issuer, Verifier: verifier, Key: p.Confirmation.JKT}
		if end.Verified() {
			return end, nil
		}
		if registered == nil {
			registered = end
		}
	}
	if registered != nil {
		return registered, nil
	}
	if unavailableErr != nil {
		return nil, unavailableErr
	}
	if claimErr != nil {
		return nil, claimErr
	}
	if lastErr != nil {
		return nil, fmt.Errorf("%w: %v", ErrVerifyFailed, lastErr)
	}
	return nil, ErrUnknownAuthority
}

// pairKeySet fetches, once, the keys confirming an envelope's from and to.
type pairKeySet struct {
	c        *Client
	env      *gobl.Envelope
	from, to Address
	keys     []*dsig.PublicKey
	err      error
	fetched  bool
}

// require fails unless the key with thumbprint jkt confirms the pair.
func (s *pairKeySet) require(ctx context.Context, jkt string) error {
	if !s.fetched {
		s.keys, s.err = s.c.pairKeys(ctx, s.env, s.from, s.to)
		s.fetched = true
	}
	if s.err != nil {
		return s.err
	}
	if !hasKey(s.keys, jkt) {
		return fmt.Errorf("%w: registration is for another key of %q", ErrVerifyFailed, s.from)
	}
	return nil
}

// checkRegistration checks the claims of an authority countersignature:
// an iat no later than now, an exp at most 90 days after it that has not
// passed (per RFC 7519, now must be strictly before it), and a cnf.
func checkRegistration(p *head.SigningPayload) error {
	if err := checkIssued(p); err != nil {
		return err
	}
	if p.ExpiresAt == 0 {
		return fmt.Errorf("%w: signature by %q has no exp", ErrVerifyFailed, p.Iss)
	}
	if time.Unix(p.ExpiresAt, 0).Sub(time.Unix(p.IssuedAt, 0)) > maxRegistrationLifetime {
		return fmt.Errorf("%w: signature by %q has an exp more than 90 days after its iat", ErrVerifyFailed, p.Iss)
	}
	if expired(p.ExpiresAt) {
		return fmt.Errorf("%w: exp %s has passed", ErrSignatureExpired,
			time.Unix(p.ExpiresAt, 0).UTC().Format(time.RFC3339))
	}
	if p.Confirmation == nil || p.Confirmation.JKT == "" {
		return fmt.Errorf("%w: signature by %q has no cnf", ErrVerifyFailed, p.Iss)
	}
	return nil
}

// confirmVerifier resolves the verifier claim on a verified authority
// countersignature. It returns the verifier's address when the claim
// names a valid address, distinct from the authority, whose countersignature
// for subject verifies. An absent, self-referential, malformed, or invalid claim returns
// "", degrading the endorsement to registered. A transient key-fetch failure
// returns ErrUnavailable instead of silently downgrading the endorsement.
func (c *Client) confirmVerifier(ctx context.Context, env *gobl.Envelope, authority, subject Address, p *head.SigningPayload) (Address, error) {
	verifier, err := ParseAddress(p.Verifier)
	if err != nil {
		return "", nil //nolint:nilerr // malformed claim degrades to registered
	}
	// Registration and identity verification must be performed by distinct
	// participants. A self-referential claim provides no verifier evidence.
	if verifier == authority {
		return "", nil
	}
	if _, err := c.verificationBy(ctx, env, verifier, subject); err != nil {
		if errors.Is(err, ErrUnavailable) {
			return "", err
		}
		return "", nil //nolint:nilerr // invalid or absent signature degrades
	}
	return verifier, nil
}

// Verification is a verifier countersignature on a party envelope, as an
// authority reads it before naming the verifier.
type Verification struct {
	// Verifier is the countersignature's issuer.
	Verifier Address
	// Key is the thumbprint in its cnf claim: the subject's key the
	// verifier saw.
	Key string
	// IssuedAt and ExpiresAt are its iat and exp; ExpiresAt is zero
	// when it does not expire.
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// FindVerification returns the valid countersignature by verifier for
// subject, with an iat no later than now (allowing clock skew), no exp
// that has passed, and a cnf naming the subject key the verifier saw.
func (c *Client) FindVerification(ctx context.Context, env *gobl.Envelope, verifier, subject Address) (*Verification, error) {
	if env == nil || env.Head == nil || !env.Signed() {
		return nil, fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
	}
	if len(env.Signatures) > maxEnvelopeSignatures {
		return nil, fmt.Errorf("%w: envelope carries %d signatures (max %d)",
			ErrVerifyFailed, len(env.Signatures), maxEnvelopeSignatures)
	}
	verifier, err := ParseAddress(string(verifier))
	if err != nil {
		return nil, err
	}
	subject, err = ParseAddress(string(subject))
	if err != nil {
		return nil, err
	}
	p, err := c.verificationBy(ctx, env, verifier, subject)
	if err != nil {
		return nil, err
	}
	v := &Verification{
		Verifier: verifier,
		Key:      p.Confirmation.JKT,
		IssuedAt: time.Unix(p.IssuedAt, 0).UTC(),
	}
	if p.ExpiresAt != 0 {
		v.ExpiresAt = time.Unix(p.ExpiresAt, 0).UTC()
	}
	return v, nil
}

// verificationBy returns the payload of a valid countersignature by addr
// for subject, with an iat no later than now (allowing clock skew), no exp
// that has passed, and a cnf.
func (c *Client) verificationBy(ctx context.Context, env *gobl.Envelope, addr, subject Address) (*head.SigningPayload, error) {
	var lastErr error
	candidates := 0
	for _, sig := range env.Signatures {
		p, err := head.SignedPayload(sig)
		if err != nil {
			continue
		}
		issuer, err := ParseAddress(p.Iss)
		if err != nil || issuer != addr {
			continue
		}
		if aud, err := ParseAddress(p.Aud); err != nil || aud != subject {
			continue
		}
		if candidates++; candidates > maxIssuerCandidates {
			break
		}
		pub, err := c.FetchKey(ctx, issuer, sig.KeyID())
		if err != nil {
			lastErr = err
			continue
		}
		if err := env.Head.Verify(sig, pub); err != nil {
			lastErr = err
			continue
		}
		if err := checkIssued(p); err != nil {
			lastErr = err
			continue
		}
		if expired(p.ExpiresAt) {
			lastErr = fmt.Errorf("%w: exp %s has passed", ErrSignatureExpired,
				time.Unix(p.ExpiresAt, 0).UTC().Format(time.RFC3339))
			continue
		}
		if p.Confirmation == nil || p.Confirmation.JKT == "" {
			lastErr = fmt.Errorf("%w: signature by %q has no cnf", ErrVerifyFailed, addr)
			continue
		}
		return p, nil
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("%w: no signature from %q for %q", ErrVerifyFailed, addr, subject)
}

// expired reports whether a signed exp claim (zero when absent) has
// passed.
func expired(exp int64) bool {
	return exp != 0 && time.Now().UTC().Unix() >= exp
}
