package net

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/head"
)

// Authorities is the default list of trusted registration authorities.
// There is no default list of trusted verifiers; see WithVerifiers.
var Authorities = []Address{
	"lookup.gobl.org",
}

// SandboxAuthorities is the trust list selected by WithSandbox. Sandbox
// endorsements must not be accepted in live contexts.
var SandboxAuthorities = []Address{
	"lookup.sandbox.gobl.org",
}

// maxEnvelopeSignatures bounds how many signatures VerifyEndorsement is
// willing to examine: each candidate can cost a network key fetch, so
// an unbounded envelope would be a fetch-amplification gadget.
// Legitimate envelopes carry a handful of signatures.
const maxEnvelopeSignatures = 32

// maxRegistrationLifetime bounds how long after its iat an authority
// countersignature's exp may fall.
const maxRegistrationLifetime = 30 * 24 * time.Hour

// maxVerificationLifetime bounds how long after its iat a verifier
// countersignature's exp may fall: one year.
const maxVerificationLifetime = 366 * 24 * time.Hour

// RegisterAuthority appends addr to the default authority list.
func RegisterAuthority(addr Address) {
	Authorities = append(Authorities, addr)
}

// Endorsement is this package's representation of a verified sender
// status (README.md §7.2): the conclusion a receiver reaches about a
// subject after evaluating the countersignatures its party envelope
// carries against this client's trust lists.
//
// There is no "endorsement" on the wire. The envelope carries
// countersignatures; the status is computed from them, and it is never
// transmitted — two clients with different trust lists can evaluate the
// same envelope and reach different conclusions, and both are right. So
// an Endorsement is always a return value, never something to construct
// or send.
type Endorsement struct {
	// Authority is the trusted authority whose countersignature named
	// the verifier. Empty when the verifier is on the client's own
	// trust list and was accepted directly.
	Authority Address
	// Verifier is the address whose countersignature confirms the
	// subject's identity (KYC/KYB).
	Verifier Address
	// VerifiedAt is the signing time of the verifier's countersignature,
	// the most recent when there are several. A verification lasts up
	// to a year; a receiver may require a more recent one.
	VerifiedAt time.Time
}

// VerifyEndorsement returns the endorsement of the subject of a party
// envelope, the address in its header's from field (README.md §7.2).
// The envelope's document must match its digest, and a valid signature by
// the subject must confirm its from and to.
//
// It first looks for a countersignature from a trusted authority, bound to
// the subject with aud, with an iat no later than now, an exp that has not
// passed and is at most 30 days after it, and a verifier, other than the
// authority, whose own countersignature for the subject is valid: an iat
// no later than now, and an exp that has not passed and is at most a year
// after it. Failing that, it accepts such a countersignature from a
// verifier on the client's own trust list (WithVerifiers).
//
// It returns ErrNotEndorsed when no trusted candidate exists,
// ErrSignatureExpired for expired evidence, ErrUnavailable for transient
// key lookup failures, or ErrVerifyFailed when candidates fail
// verification.
func (c *Client) VerifyEndorsement(ctx context.Context, env *gobl.Envelope) (*Endorsement, error) {
	// A malformed envelope may carry signatures without a header;
	// reject rather than let header verification panic.
	if env == nil || env.Head == nil || len(env.Signatures) == 0 {
		return nil, fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
	}
	if len(c.authorities) == 0 && len(c.verifiers) == 0 {
		return nil, fmt.Errorf("%w: no authorities or verifiers trusted by this client", ErrNotEndorsed)
	}
	if len(env.Signatures) > maxEnvelopeSignatures {
		return nil, fmt.Errorf("%w: envelope carries %d signatures (max %d)",
			ErrVerifyFailed, len(env.Signatures), maxEnvelopeSignatures)
	}
	if err := checkDigest(env); err != nil {
		return nil, err
	}
	subject, to, err := routing(env)
	if err != nil {
		return nil, err
	}
	s := &endorsementCheck{
		c:       c,
		env:     env,
		subject: subject,
		to:      to,
		budget:  make(issuerBudget),
		found:   make(map[Address]foundVerifications),
	}
	if end := s.authorityEndorsement(ctx); end != nil {
		return end, nil
	}
	if end := s.directEndorsement(ctx); end != nil {
		return end, nil
	}
	return nil, s.errs.result()
}

// endorsementCheck is the state of one VerifyEndorsement call.
type endorsementCheck struct {
	c           *Client
	env         *gobl.Envelope
	subject, to Address
	errs        candidateErrs
	budget      issuerBudget
	found       map[Address]foundVerifications

	// confirmedErr holds the result of confirmed once checked.
	confirmedErr error
	checked      bool
}

type foundVerifications struct {
	ps  []*head.SigningPayload
	err error
}

// issuerBudget counts, per issuer, the signatures checked against its
// keys, so no envelope triggers more than maxIssuerCandidates key
// fetches for one issuer.
type issuerBudget map[Address]int

// spend reports whether another signature by issuer may be checked.
func (b issuerBudget) spend(issuer Address) bool {
	b[issuer]++
	return b[issuer] <= maxIssuerCandidates
}

// confirmed checks, once, that the subject's signature confirms the
// envelope's from and to.
func (s *endorsementCheck) confirmed(ctx context.Context) error {
	if !s.checked {
		s.confirmedErr = s.c.confirmed(ctx, s.env, s.subject, s.to)
		s.checked = true
	}
	return s.confirmedErr
}

// verifications returns the valid verifications by verifier for the
// subject, checking them once per call.
func (s *endorsementCheck) verifications(ctx context.Context, verifier Address) ([]*head.SigningPayload, error) {
	if f, ok := s.found[verifier]; ok {
		return f.ps, f.err
	}
	ps, err := s.c.verificationsBy(ctx, s.env, verifier, s.subject, s.budget)
	s.found[verifier] = foundVerifications{ps, err}
	return ps, err
}

// authorityEndorsement returns the first valid endorsement by a trusted
// authority, recording why each other candidate failed.
func (s *endorsementCheck) authorityEndorsement(ctx context.Context) *Endorsement {
	c, env, subject, errs := s.c, s.env, s.subject, &s.errs
	auths := canonical(c.authorities)
	for _, sig := range env.Signatures {
		p, err := head.SignedPayload(sig)
		if err != nil {
			continue
		}
		issuer, err := ParseAddress(p.Iss)
		if err != nil || !auths[issuer] {
			continue
		}
		if aud, err := ParseAddress(p.Aud); err != nil || aud != subject {
			continue
		}
		if expired(p.ExpiresAt) {
			errs.add(errExpired(p.ExpiresAt))
			continue
		}
		if !s.budget.spend(issuer) {
			continue
		}
		pub, err := c.FetchKey(ctx, issuer, sig.KeyID())
		if err != nil {
			errs.add(err)
			continue
		}
		if err := env.Head.Verify(sig, pub); err != nil {
			errs.add(err)
			continue
		}
		// Claims are only trusted after the signature verifies.
		if err := checkRegistration(p); err != nil {
			errs.add(err)
			continue
		}
		if err := s.confirmed(ctx); err != nil {
			errs.add(err)
			continue
		}
		verifier, at, err := s.confirmVerifier(ctx, issuer, p)
		if err != nil {
			errs.add(err)
			continue
		}
		return &Endorsement{Authority: issuer, Verifier: verifier, VerifiedAt: at}
	}
	return nil
}

// directEndorsement returns the first valid endorsement by a verifier on
// the client's own trust list, recording why each other candidate failed.
func (s *endorsementCheck) directEndorsement(ctx context.Context) *Endorsement {
	seen := make(map[Address]bool, len(s.c.verifiers))
	for _, a := range s.c.verifiers {
		verifier, err := ParseAddress(string(a))
		if err != nil || seen[verifier] {
			continue
		}
		seen[verifier] = true
		ps, err := s.verifications(ctx, verifier)
		if err != nil {
			s.errs.add(err)
			continue
		}
		if len(ps) == 0 {
			continue
		}
		if err := s.confirmed(ctx); err != nil {
			s.errs.add(err)
			continue
		}
		return &Endorsement{Verifier: verifier, VerifiedAt: latestIssued(ps)}
	}
	return nil
}

// candidateErrs keeps the most telling reason trusted candidates failed:
// a transient failure, then expired evidence, then any other.
type candidateErrs struct {
	unavailable, expired, other error
}

func (e *candidateErrs) add(err error) {
	switch {
	case errors.Is(err, ErrUnavailable):
		e.unavailable = err
	case errors.Is(err, ErrSignatureExpired):
		e.expired = err
	default:
		e.other = err
	}
}

func (e *candidateErrs) empty() bool {
	return e.unavailable == nil && e.expired == nil && e.other == nil
}

func (e *candidateErrs) result() error {
	switch {
	case e.unavailable != nil:
		return e.unavailable
	case e.expired != nil:
		return e.expired
	case e.other != nil && errors.Is(e.other, ErrVerifyFailed):
		return e.other
	case e.other != nil:
		return fmt.Errorf("%w: %w", ErrVerifyFailed, e.other)
	}
	return ErrNotEndorsed
}

// canonical returns the set of canonical forms of addrs, so U-Label or
// trailing-dot forms — configured or signed — compare equal.
func canonical(addrs []Address) map[Address]bool {
	set := make(map[Address]bool, len(addrs))
	for _, a := range addrs {
		if canon, err := ParseAddress(string(a)); err == nil {
			set[canon] = true
		}
	}
	return set
}

// checkRegistration checks the claims of an authority countersignature:
// an iat no later than now, an exp at most 30 days after it that has not
// passed (per RFC 7519, now must be strictly before it).
func checkRegistration(p *head.SigningPayload) error {
	if err := checkIssued(p); err != nil {
		return err
	}
	if p.ExpiresAt == 0 {
		return fmt.Errorf("%w: signature by %q has no exp", ErrVerifyFailed, p.Iss)
	}
	if time.Unix(p.ExpiresAt, 0).Sub(time.Unix(p.IssuedAt, 0)) > maxRegistrationLifetime {
		return fmt.Errorf("%w: signature by %q has an exp more than 30 days after its iat", ErrVerifyFailed, p.Iss)
	}
	if expired(p.ExpiresAt) {
		return errExpired(p.ExpiresAt)
	}
	return nil
}

// confirmVerifier returns the verifier an authority countersignature
// names, after checking that it is a valid address, distinct from the
// authority, whose own countersignature for subject is valid, and the
// signing time of the most recent such countersignature.
func (s *endorsementCheck) confirmVerifier(ctx context.Context, authority Address, p *head.SigningPayload) (Address, time.Time, error) {
	var none time.Time
	if p.Verifier == "" {
		return "", none, fmt.Errorf("%w: signature by %q names no verifier", ErrVerifyFailed, authority)
	}
	verifier, err := ParseAddress(p.Verifier)
	if err != nil {
		return "", none, fmt.Errorf("%w: signature by %q names an invalid verifier: %v", ErrVerifyFailed, authority, err)
	}
	// Registration and identity verification must be performed by distinct
	// participants. A self-referential claim provides no verifier evidence.
	if verifier == authority {
		return "", none, fmt.Errorf("%w: signature by %q names itself as verifier", ErrVerifyFailed, authority)
	}
	ps, err := s.verifications(ctx, verifier)
	if err != nil {
		return "", none, err
	}
	if len(ps) == 0 {
		return "", none, noVerification(verifier, s.subject)
	}
	return verifier, latestIssued(ps), nil
}

// Verification is a verifier countersignature on a party envelope, as an
// authority reads it before naming the verifier.
type Verification struct {
	// Verifier is the countersignature's issuer.
	Verifier Address
	// IssuedAt and ExpiresAt are its iat and exp.
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// FindVerifications returns the valid countersignatures by verifier for
// subject: each with an iat no later than now (allowing clock skew), and
// an exp that has not passed and is at most a year after it. An authority
// sets its own exp no later than theirs. It returns an error when there are none, or when the envelope's document
// does not match its digest.
func (c *Client) FindVerifications(ctx context.Context, env *gobl.Envelope, verifier, subject Address) ([]*Verification, error) {
	if env == nil || env.Head == nil || !env.Signed() {
		return nil, fmt.Errorf("%w: envelope is not signed", ErrVerifyFailed)
	}
	if len(env.Signatures) > maxEnvelopeSignatures {
		return nil, fmt.Errorf("%w: envelope carries %d signatures (max %d)",
			ErrVerifyFailed, len(env.Signatures), maxEnvelopeSignatures)
	}
	if err := checkDigest(env); err != nil {
		return nil, err
	}
	verifier, err := ParseAddress(string(verifier))
	if err != nil {
		return nil, err
	}
	subject, err = ParseAddress(string(subject))
	if err != nil {
		return nil, err
	}
	ps, err := c.verificationsBy(ctx, env, verifier, subject, make(issuerBudget))
	if err != nil {
		return nil, err
	}
	if len(ps) == 0 {
		return nil, noVerification(verifier, subject)
	}
	vs := make([]*Verification, len(ps))
	for i, p := range ps {
		vs[i] = &Verification{
			Verifier:  verifier,
			IssuedAt:  issuedAt(p),
			ExpiresAt: time.Unix(p.ExpiresAt, 0).UTC(),
		}
	}
	return vs, nil
}

// verificationsBy returns the payloads of the valid verification
// countersignatures by addr for subject, those without a verifier claim,
// each with an iat no later than now (allowing clock skew), and an exp
// that has not passed and is at most a year after it. Each signature
// checked against a key spends budget for addr; an expired one is skipped
// first.
// The error is set only when there are candidates and none is valid,
// and reports a transient failure over expired evidence over any other.
func (c *Client) verificationsBy(ctx context.Context, env *gobl.Envelope, addr, subject Address, budget issuerBudget) ([]*head.SigningPayload, error) {
	var valid []*head.SigningPayload
	var errs candidateErrs
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
		if p.Verifier != "" {
			continue // a registration, not a verification
		}
		if expired(p.ExpiresAt) {
			errs.add(errExpired(p.ExpiresAt))
			continue
		}
		if !budget.spend(issuer) {
			break
		}
		pub, err := c.FetchKey(ctx, issuer, sig.KeyID())
		if err != nil {
			errs.add(err)
			continue
		}
		if err := env.Head.Verify(sig, pub); err != nil {
			errs.add(err)
			continue
		}
		if err := checkIssued(p); err != nil {
			errs.add(err)
			continue
		}
		if p.ExpiresAt == 0 {
			errs.add(fmt.Errorf("%w: signature by %q has no exp", ErrVerifyFailed, addr))
			continue
		}
		if time.Unix(p.ExpiresAt, 0).Sub(time.Unix(p.IssuedAt, 0)) > maxVerificationLifetime {
			errs.add(fmt.Errorf("%w: signature by %q has an exp more than a year after its iat", ErrVerifyFailed, addr))
			continue
		}
		valid = append(valid, p)
	}
	if len(valid) > 0 {
		return valid, nil
	}
	if errs.empty() {
		return nil, nil
	}
	return nil, errs.result()
}

// noVerification reports that an envelope has no countersignature by
// verifier for subject.
func noVerification(verifier, subject Address) error {
	return fmt.Errorf("%w: no signature from %q for %q", ErrVerifyFailed, verifier, subject)
}

// expired reports whether a signed exp claim (zero when absent) has
// passed.
func expired(exp int64) bool {
	return exp != 0 && time.Now().UTC().Unix() >= exp
}

// errExpired reports that a signed exp claim has passed.
func errExpired(exp int64) error {
	return fmt.Errorf("%w: exp %s has passed", ErrSignatureExpired,
		time.Unix(exp, 0).UTC().Format(time.RFC3339))
}

// issuedAt returns a signature's iat as a time.
func issuedAt(p *head.SigningPayload) time.Time {
	return time.Unix(p.IssuedAt, 0).UTC()
}

// latestIssued returns the most recent iat of ps, which is not empty.
func latestIssued(ps []*head.SigningPayload) time.Time {
	at := issuedAt(ps[0])
	for _, p := range ps[1:] {
		if t := issuedAt(p); t.After(at) {
			at = t
		}
	}
	return at
}
