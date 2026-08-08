# GOBL Net

> ⚠️ **EXPERIMENTAL** — GOBL Net is under active development. The package
> API and the wire protocol may change without notice and are not yet
> covered by any stability guarantee.

**Status:** Draft. This document is the wire-protocol specification for
GOBL Net and describes the current state of the
`github.com/invopop/gobl/net` package and related supporting code in
`dsig`. CLI commands (`gobl init`, `gobl net who/send/serve`, `gobl sign
--domain …`, `gobl verify --remote`) and the reference server live in
[`github.com/invopop/gobl.dev`](https://github.com/invopop/gobl.dev/#gobl-net) —
that README covers on-disk layout, ACME, structured logging, and the
operational stances. The protocol is pre-1.0 and subject to change.

## Abstract

GOBL Net is a decentralized identity and discovery protocol for signed GOBL
documents. It binds a document signature to a fully qualified domain name
(FQDN), and defines a small set of well-known HTTPS endpoints at that domain:

- `/.well-known/gobl/keys/<kid>` — a single public JWK looked up by key
  ID, used to verify signatures.
- `/.well-known/gobl/who` — a signed GOBL Envelope carrying an
  `org.Party` document, endorsing the holder's identity, retrieved
  with an authenticated GET.
- `/.well-known/gobl/inbox` — a write endpoint that accepts signed
  envelopes addressed to the holder.

Trust in an identity is anchored to the TLS certificate for the
Address's FQDN: a signature's verifiable origin lives in its signed
`iss`, the verifier fetches the corresponding public key from
`https://<iss>/.well-known/gobl/keys/<kid>`, and the HTTPS connection
proves the response really came from that FQDN. Key discovery is
unconditional: it works the same for every participant and every
policy below layers on top of it, never replaces it.

The two roles in an exchange carry asymmetric requirements. *Anyone*
may provision a receiving address — publishing an inbox requires no
approval from any third party, and a receive-only participant need
not even publish identity details (§8.2). *Senders*, by contrast, are
expected to carry an endorsement: a countersignature on their `/who`
identity from a KYC vendor ("Authority") that receivers trust.
Receivers accept or reject incoming documents on the strength of that
endorsement (§6.4, §8.4).

Requests to the who and inbox endpoints are themselves
authenticated: every caller presents a short-lived request token
(§5.5) signed with its own published key, so servers always know —
and may record (§11.9) — who is asking.

The protocol layers on top of standards already in use by GOBL:

- **RFC 7515** — JSON Web Signature, for envelope signing.
- **RFC 7517** — JSON Web Key Set, for key discovery.
- **RFC 7519** — JSON Web Token, for request authentication.
- **RFC 3986** — URI syntax, for discovery transport.

GOBL Net does not define new cryptographic primitives. It defines (a) how a
GOBL signature identifies its origin, (b) how to discover the verifying
keys, (c) how to retrieve an endorsed identity for an address, and (d) how
to deliver a signed envelope to its recipient.

## 1. Conventions

The key words "MUST", "MUST NOT", "SHOULD", "SHOULD NOT", and "MAY" in this
document are to be interpreted as described in BCP 14 (RFC 2119, RFC 8174).

## 2. Terminology

- **Address** — An FQDN identifying a GOBL Net participant, e.g.
  `billing.invopop.com`. Represented in code as `net.Address`.
- **Published key** — A `dsig.PublicKey` (an RFC 7517 JWK plus the
  optional `valid_from` / `valid_until` extension members) served at
  `/.well-known/gobl/keys/<kid>`.
- **Party Envelope** — A signed GOBL Envelope whose document is an
  `org.Party`, served at the who endpoint. The first signature is the
  subject's self-signature; Authority countersignatures follow.
- **Sender / Receiver** — The two roles in a document exchange. A
  receiver only needs a domain, TLS, and (if it signs or makes
  authenticated requests — §5.5) published keys. A sender is
  additionally expected to carry an Authority endorsement on its who
  identity (§6.4).
- **Authority / Verifier** — An Authority is an address receivers
  trust to countersign who identities, asserting registration
  (FQDN ownership). A Verifier is an authority named by a
  registration Authority's `verifier` claim as having performed
  identity verification (KYC/KYB) of the subject, confirmed by its
  own countersignature (§5.3).
- **iss / aud** — Fields in a signature's *signed payload* carrying the
  verifiable GOBL Net origin (`iss`) and the address the signature is
  bound to (`aud`), both as bare GOBL Net addresses (FQDNs) — GOBL
  Net is implied inside the signed payload, so no URI scheme is
  carried. These are the authoritative, tamper-proof identities.
- **Request token** — A short-lived JWT presented in the
  `Authorization` header of who and inbox requests, identifying the
  party *making the HTTP request* (§5.5). Its `iss` may name a
  trusted intermediary and is independent of any signature inside
  the request body.
- **header from / to** — Optional unsigned `cbc.URI` fields on the
  envelope header expressing *intent/routing* in any scheme
  (`iso6523-actorid-upis:`, `mailto:`, `gobl:`…). Useful for interop
  with other formats; not used for verification.

## 3. Addressing

### 3.1 Address Format

A GOBL Net Address is an FQDN. It MUST NOT contain a scheme, port, path,
query, or fragment.

Parsing is performed by `net.ParseAddress` and applies the following
normalizations and constraints:

1. Surrounding whitespace is trimmed.
2. A trailing dot, if present, is stripped.
3. The input is normalised to its ASCII (A-Label / Punycode) form via
   `golang.org/x/net/idna` Lookup profile (IDNA2008). This converts
   any U-Labels to A-Labels and lowercases ASCII labels in one step;
   invalid IDN labels MUST be rejected with `ErrAddressInvalid`.
4. The result MUST satisfy `is.DNSName` (RFC 1035 label syntax).
5. The result MUST contain at least one dot (i.e. at least two labels).

Inputs that contain a scheme, port, or path MUST be rejected with
`ErrAddressInvalid`. An empty input MUST be rejected with `ErrAddressEmpty`.

**Canonical form.** Addresses on the wire — in `iss`, `aud`,
well-known URLs, and identity files — MUST be ASCII
A-Labels. Implementations are free to accept U-Label input from
humans (e.g. CLI flags) but MUST normalise to A-Label before
signing, fetching, or comparing.

**Examples.** Accepted (after normalization):

| Input                       | Parsed Address           |
|-----------------------------|--------------------------|
| `billing.invopop.com`       | `billing.invopop.com`    |
| `sub.domain.example.org`    | `sub.domain.example.org` |
| `Billing.Invopop.COM`       | `billing.invopop.com`    |
| `billing.invopop.com.`      | `billing.invopop.com`    |
| `  billing.invopop.com  `   | `billing.invopop.com`    |
| `München.DE`                | `xn--mnchen-3ya.de`      |
| `xn--mnchen-3ya.de`         | `xn--mnchen-3ya.de`      |

Rejected:

| Input                  | Error                |
|------------------------|----------------------|
| `` (empty)             | `ErrAddressEmpty`    |
| `localhost`            | `ErrAddressInvalid` (single label) |
| `http://example.com`   | `ErrAddressInvalid` (scheme) |
| `example.com/path`     | `ErrAddressInvalid` (path)   |
| `example.com:8080`     | `ErrAddressInvalid` (port)   |
| `not valid!.com`       | `ErrAddressInvalid` (illegal characters) |
| `-bad.com`             | `ErrAddressInvalid` (invalid IDN label) |

### 3.2 Well-Known Paths

The following constants are defined in `net/address.go`:

| Constant         | Value                                |
|------------------|--------------------------------------|
| `WellKnownPath`  | `/.well-known/gobl`                  |
| `KeysPath`       | `/.well-known/gobl/keys` (prefix)    |
| `KeyPath(kid)`   | `/.well-known/gobl/keys/<kid>`       |
| `WhoPath`        | `/.well-known/gobl/who`              |
| `InboxPath`      | `/.well-known/gobl/inbox`            |
| `JWKSPath`       | `/.well-known/jwks.json`             |

For an Address `A`, the canonical URI and URLs are

```
<A>                                        ← Address.String() (iss / aud / verifier value)
gobl:<A>                                   ← Address.URI()  (endpoint / header from-to value)
https://<A>/.well-known/gobl/keys/<kid>    ← KeyURL(kid)
https://<A>/.well-known/gobl/who           ← WhoURL()
https://<A>/.well-known/gobl/inbox         ← InboxURL()
https://<A>/.well-known/jwks.json          ← JWKSURL()
```

The protocol exposes two key-discovery surfaces:

1. **Per-kid** at `/.well-known/gobl/keys/<kid>` — single-JWK lookups,
   used by `Client.FetchKey` during verification. Scales as keys
   rotate; an absent or retired kid returns `404`.
2. **Bulk JWKS** at `/.well-known/jwks.json` — standard RFC 7517 JWK
   Set. GOBL Net signatures do **not** carry a `jku` header, so this
   endpoint is for tooling that fetches a JWK Set by convention
   (e.g. derived from the signer's domain) rather than from a JWS
   header reference.

The scheme MUST be `https`. HTTP is not a permitted alternative for
production deployments; client tooling MAY offer an opt-in for plain
HTTP solely for local development.

Responses are always JSON, so file extensions are omitted from the paths.

## 4. Keys

Each domain publishes one or more public keys, each individually
addressable by its `kid` at `/.well-known/gobl/keys/<kid>`. The response
body is a single RFC 7517 JSON Web Key, optionally augmented with a
GOBL Net validity window:

```
{
  "kty": "EC", "crv": "P-256", "kid": "…", "x": "…", "y": "…",
  "valid_from":  "2026-01-01T00:00:00.000Z",
  "valid_until": "2027-01-01T00:00:00.000Z"
}
```

`valid_from` and `valid_until` are *additional JWK members* in the sense
of RFC 7517 §4 — implementations that do not recognise them MUST ignore
them, so the response remains a conformant JWK for any standard JOSE
consumer.

When set, the validity window bounds the signing time that each
signature may carry in its signed payload (`iat`). A verifier rejects
a signature whose `iat` falls outside `[valid_from, valid_until]`.
The checks degrade gracefully: an absent bound on the key, or an
absent `iat` on the signature, simply skips that half of the
comparison.

`valid_from` is stamped automatically when a key is generated.
`valid_until` is left empty and is meant to be set when the operator
rotates the key out — either at retirement, or in advance of a
planned rotation. A retired key remains published so that historical
envelopes signed within its window still verify; only signatures with
an `iat` past `valid_until` are rejected.

Verifiers MUST treat unknown kids as `404 Not Found`; this is how a
domain expresses that a key has been removed entirely (as distinct from
"retired but still serving historical verification"). The per-kid path
is the only one GOBL's own `Client.FetchKey` uses.

### 4a. Bulk JWKS endpoint (jwt.io interop)

In addition to the per-kid endpoint, the server publishes a standard
RFC 7517 JWK Set at `/.well-known/jwks.json`. This is provided as a
convenience for operators inspecting their keys via `curl` and for
third-party JOSE tooling that prefers to fetch a full key set by URL
rather than per-kid. GOBL Net's own verifier
(`Client.FetchKey` → per-kid endpoint) does not use it.

Response shape:

```json
{
  "keys": [
    { "kty": "EC", "kid": "<newest UUIDv7>", "valid_from": "…", … },
    { "kty": "EC", "kid": "<older  UUIDv7>", "valid_until": "…", … }
  ]
}
```

Keys are returned **newest first**, ordered by `valid_from`
descending; entries without a `valid_from` sort last. Since key IDs
are UUIDv7 (time-ordered), kid descending is the deterministic
tie-breaker.

The JOSE header on each signature carries `alg` and `kid` only. The
signed payload's `iss` value names the issuing GOBL Net address; a
verifier resolves that to the per-kid URL via `Client.FetchKey`.

The key type and the validity-window enforcement live in the `dsig`
package: a published key is a `dsig.PublicKey`. `head.Header.Verify`
calls `dsig.PublicKey.Allows` automatically, so every signed-envelope
verifier — whether through GOBL Net or a direct `Envelope.Verify` call
— enforces the window.

A conforming server MUST serve each published key verbatim at
`/.well-known/gobl/keys/<kid>` and the bulk set at
`/.well-known/jwks.json`. The reference implementation in
`gobl.dev`'s `gobl net serve` stores one file per `kid` on disk and
maps 1:1 to a future row-per-`kid` database — but the on-disk layout
is an implementation detail of that server, not part of this
protocol.

Endorsement of a participant happens at the identity layer (see §6) and
does *not* live inside the key material itself.

## 5. Signatures

### 5.1 Signed identities: `iss` / `aud` / `iat`

Each signature signs a payload of `{uuid, dig, iss, aud, iat}`:

- `uuid` + `dig` identify the document (immutable after signing).
- `iss` is the signer's verifiable GOBL Net address as a bare FQDN
  (e.g. `billing.invopop.com`). The verifier reads it to discover
  *which* per-key endpoint to fetch. No URI scheme is carried:
  within the protocol the value can only be a GOBL Net address, and
  since an FQDN can never contain a colon, a future revision could
  admit URI forms without ambiguity.
- `aud` is the GOBL Net address the signature is bound to. It is
  optional at the protocol layer (e.g. a self-signed identity
  document, or an envelope archived for later inspection, may omit
  it), but **inboxes MUST require it** (§8.3): an envelope POSTed
  to a `/inbox` MUST carry `aud == <inbox-owner>`, otherwise the
  request is rejected. Binding the audience inside the signature
  prevents the same valid invoice from being replayed against
  multiple inboxes.
- `iat` is the signing time as a JWT-standard NumericDate (Unix
  seconds, per RFC 7519 §2). It is set automatically by `Sign`;
  verifiers read it for the per-key validity window check but no
  freshness policy is enforced by default — receivers may apply their
  own max-age window when relevant.
- `exp` (optional, JWT-standard per RFC 7519 §4.1.4, set with
  `head.WithExpiration`) is the time after which the signature's
  assertions should no longer be relied upon. Ordinary transport
  signatures omit it — archived envelopes must keep verifying
  indefinitely. Authorities set it on their countersignatures to
  bound the endorsement's lifetime (§5.3).

Because `iss`/`aud`/`iat` are inside the signed payload, the origin,
audience, and signing time are all tamper-proof. Multiple parties may
countersign the same document (shared `uuid`+`dig`) each with their own
`iss`/`aud`/`iat`.

The header `from`/`to` (`cbc.URI`) are a *separate*, unsigned layer for
intent/routing in any scheme; they are never used for verification.

### 5.2 Envelope Signing

`Envelope.Sign(key, opts...)` (→ `head.Header.Sign`) signs the document
identity plus the signer's `iss` and optional `aud`, set with the
`head.WithIssuer` and `head.WithAudience` options. Both may be omitted
for a plain, non-GOBL-Net signature.

### 5.3 Authority countersignatures and verification

A `/who` response (and any envelope an Authority chooses to endorse)
MAY carry one or more countersignatures from addresses in the
receiver's `Authorities` list (§6.3). A countersignature from a
trusted Authority asserts that the subject is **registered**: the
Authority has confirmed the subject controls the named GOBL Net
address (FQDN ownership).

**Verified identities.** When the subject has additionally passed
identity verification (KYC/KYB), the registration Authority's
countersignature carries a `verifier` claim: the GOBL Net address
(a bare FQDN) of the authority that performed the verification.
The named verifier MUST also countersign the same envelope; the
subject is **verified** only when both hold — the registration
Authority's pointer authorizes *which* signature to trust for
verification, and the verifier's own signature is the evidence. An
Authority that performs verification itself names itself, and its
single countersignature serves as both attestations. Both levels
are structural rather than declared: registration is asserted by
the countersignature's presence, verification by the confirmed
`verifier` claim.

```go
// Registration authority (e.g. lookup.gobl.org):
env.Sign(authorityKey,
    head.WithIssuer(authorityAddr.String()),
    head.WithAudience(subjectAddr.String()),
    head.WithVerifier(verifierAddr.String()))
// Verifying authority (KYC/KYB vendor):
env.Sign(verifierKey,
    head.WithIssuer(verifierAddr.String()),
    head.WithAudience(subjectAddr.String()))
```

A named verifier whose countersignature is absent, invalid, or
expired MUST degrade the endorsement to registered rather than
invalidate it: the registration stands on its own, and callers that
require verification reject with `ErrNotVerified`
(`Endorsement.Verified` in code). A `verifier` value that is not a
valid bare address is treated as absent. Adding or revoking
verification is the registration Authority's act: it re-countersigns
the envelope with the pointer added or removed, which makes the
registry the single source of truth for verification state and
requires coordination between the two authorities when KYC
completes.

The two countersignatures carry independent `exp` claims and
deliberately independent lifecycles. A registration
countersignature SHOULD carry an `exp` of at most 90 days, after
which the subject renews its registration (§11.4). A verifier
countersignature MAY be much longer-lived — a year or more, in the
spirit of EV TLS certificates — since the underlying KYC/KYB
process is expensive to repeat; its `exp` bounds the verification
independently of the registration cycle. Verifiers of either kind
MUST treat an expired countersignature as absent:
`Client.VerifyAuthority` rejects expired registration candidates
with `ErrSignatureExpired` and degrades an expired verifier
signature to registered. A countersignature without `exp` does not
expire; verifiers MAY apply their own `iat` max-age policy to such
legacy endorsements.

### 5.4 X.509 evidence (optional, long-term storage)

JWS allows an `x5c` header carrying an X.509 certificate chain (RFC
7515 §4.1.6). GOBL Net does **not** use `x5c` for signature
verification — the trust anchor is the signed `iss` resolved through
the per-key endpoint, secured by the issuer's TLS certificate (see
§11.1). But signers MAY stamp `x5c` on a signature as additional,
self-contained evidence for archival and long-term verification: a
recipient who keeps the envelope for many years, long after the
issuer's `/.well-known/gobl/keys/<kid>` may have stopped responding,
can still verify the signature against the embedded chain and walk
that chain to a CA they trust at audit time.

Treat `x5c` as supplementary evidence, not as a replacement for §11.1:
the chain proves that *some* CA-attested entity signed the document,
but does not bind the signature to a GOBL Net Address. Verifiers MAY
consult `x5c` for archival proofs, but for live verification of an
inbox delivery or `/who` lookup the signed `iss` + Web PKI MUST be
the authoritative chain.

The library does not stamp `x5c` automatically. Operators with an
archival need can add it via a future `dsig.WithX5C` option or by
constructing the JWS with the underlying go-jose options directly.

### 5.5 Request tokens (`Authorization` header)

Requests to the who and inbox endpoints (§8.2, §8.3) MUST carry a
request token: a compact JWS in the JWT form (RFC 7519), sent as an
`Authorization: Bearer` header (RFC 6750) and signed with the same
ES256 key material as any other GOBL Net signature (§4, §5). The
token authenticates the party *making the HTTP request*; it is
independent of, and may name a different party than, any signature
inside the request body.

The signed claims are:

- `iss` — the requester's Address as a bare FQDN. This MAY be a
  trusted intermediary transmitting a document on behalf of its
  signer; the envelope's own `iss` (§5.1) remains the authoritative
  document origin.
- `aud` — the destination Address as a bare FQDN. Verifiers MUST
  reject a token whose `aud` is not their own address: binding the
  audience prevents a captured token from being replayed against a
  different server.
- `iat` — issue time as a NumericDate. REQUIRED.
- `exp` — expiry. REQUIRED. Clients are responsible for keeping the
  window short: `exp − iat` SHOULD be between 30 seconds and 5
  minutes (60 seconds is a sensible default).
- `jti` — OPTIONAL unique token id (UUID) for audit correlation.

The protected header MUST carry the signing key's `kid` and MAY
carry `typ: "JWT"` for generic-tooling interop.

**Verification.** The server reads the unverified `iss` and `kid`,
fetches the key from `https://<iss>/.well-known/gobl/keys/<kid>`
(key endpoints are open — §8.1 — so request authentication never
recurses), verifies the signature, and checks that `aud` equals its
own address and that the key's validity window (§4) allows `iat`.
Freshness: a token MUST be rejected when `exp` has passed, when
`iat` lies in the future, or when `iat` is older than 5 minutes
regardless of `exp`; verifiers SHOULD allow 30 seconds of clock
skew. A request with a missing or invalid token MUST be rejected
with `401 Unauthorized`. A verifier that cannot *reach* the
issuer's key endpoint (`ErrUnavailable`: network failure, 429, 5xx)
MUST NOT treat the token as invalid — the server responds
`503 Service Unavailable` so the client retries, rather than `401`,
which clients treat as a definitive rejection.

Published keys are immutable per `kid` (§4), so verifiers SHOULD
cache fetched keys for a short TTL rather than re-fetch per
request; the reference client caches for 5 minutes (§7.1). When
the request token and the body's signature are made with the same
key — the common case of a sender delivering its own documents —
token and envelope verification then share a single key fetch. The
TTL also bounds how long a key removed from the issuer's endpoint
(§4, 404) keeps verifying; `Client.FlushKeyCache` empties the cache
on demand for operators reacting to a reported compromise.

Any participant with published keys can mint a request token —
endorsement requirements (§6.4) attach to the *document signer*,
not the requester. Servers MAY keep an audit log of authenticated
requests (§11.9) and MAY apply their own policy to which requesters
they serve.

Because tokens are cheap to mint and single-purpose, clients SHOULD
mint a fresh token per request rather than reuse one across
requests; `jti` gives operators a hook for replay suppression
within the freshness window, but the protocol does not require it.

In code: `net.NewToken(key, iss, aud, ttl)` mints a token, a
`Client` configured with `net.WithIdentity(addr, key)` attaches one
to every who and inbox request automatically, and servers verify
inbound tokens with `Client.VerifyToken(ctx, token, aud)`, which
returns the verified requester Address.

## 6. Verification

### 6.1 Envelope Verification Flow

`Client.VerifyEnvelope(ctx, env, expectedAud)` returns the verified
issuer address:

1. The envelope MUST be signed; otherwise `ErrVerifyFailed`.
2. The first signature's signed payload is read; `iss` MUST be a
   valid bare Address (else `ErrVerifyFailed`).
3. `FetchKey(ctx, iss-host, kid)` fetches the issuer's published key
   from `/.well-known/gobl/keys/<kid>` (including its optional
   `valid_from` / `valid_until`).
4. The envelope is verified against that public key.
5. If `expectedAud` is non-empty, the signed `aud` MUST equal it.
6. If the key declares a validity window, the signed `iat` MUST fall
   within `[valid_from, valid_until]` (each bound optional).
7. The verified issuer address is returned.

### 6.2 Identity lookup (`GET /who`)

`/who` is an authenticated GET (see §8.2): the caller presents a
request token (§5.5) identifying itself. The response is the
target's party envelope: document = the target's `org.Party`,
first signature = the target's self-signature with `iss=target`
and no `aud` (the response is the same signed document for every
authorized caller), optionally followed by Authority
countersignatures.

`Client.Who(ctx, addr)` performs the lookup and verifies it:

1. A request token for `aud=<addr>` is minted from the client's
   identity (`WithIdentity`) and sent with
   `GET https://<addr>/.well-known/gobl/who`. A `204` returns
   `ErrNoContent` — the address exists but publishes no identity
   details (a receive-only account). A `202` returns `ErrPending` —
   the request was recorded and the owner may deliver its party
   envelope to the caller's inbox later (§8.2).
2. The response envelope's first signature is verified via
   `VerifyEnvelope` (the signed `iss` resolved to a published key).
3. The verified issuer MUST equal the fetched address — a valid
   envelope for a *different* identity served at this URL is
   rejected.
4. The document MUST be an `org.Party`, else `ErrPartyMissing`.

The response body is still a static signed document — the request
token controls *access*; it does not bind the response to the
caller. The identity's integrity comes from the self-signature and
the TLS origin, so clients MAY cache the verified envelope locally
within a modest TTL (§8.2), which also avoids minting a fresh token
per lookup.

### 6.3 Trusted Authorities

The package-level slice `net.Authorities` holds GOBL Net addresses
treated as trusted registration authorities. The default list
contains the network's default authority, `lookup.gobl.org`, and
`net.RegisterAuthority` appends to it. The `WithAuthorities` client
option *replaces* the list for that client: whoever configures
trust states it fully, so closed deployments can exclude the
default authority entirely (include it explicitly to supplement).

`Client.VerifyAuthority(ctx, env)` returns an `Endorsement` iff the
envelope carries at least one signature whose signed `iss` is in the
client's authorities AND that signature cryptographically verifies
against the authority's published key AND its signed `exp` claim
(if any) has not passed. The endorsement names the authority and,
when the authority's `verifier` claim is confirmed by the named
verifier's own countersignature (§5.3), the verifier;
`Endorsement.Verified` reports the distinction. It returns
`ErrUnknownAuthority` when no candidate signature is from a known
authority (or none has been registered), `ErrSignatureExpired` when
a verified authority signature has expired, and `ErrVerifyFailed`
when a candidate fails its crypto check.

**Sandbox.** The network runs a parallel sandbox environment:
`lookup.sandbox.gobl.org` (the default entry in
`net.SandboxAuthorities`) operates the same registration service as
the live authority, backed by its own database and its own accepted
verification providers — typically relaxed KYB suited to test
identities. The live and sandbox trust lists are disjoint by
construction and MUST stay that way: a live verifier never accepts
a sandbox endorsement, and a sandbox verifier never needs a live
one. Clients opt in with `net.WithSandbox()`, which replaces the
trust list with the sandbox authorities; everything else — request
tokens, endorsement checks, the verifier claim — behaves
identically in both environments.

Note the trust topology this creates: receivers list *registration*
authorities; the verifiers those authorities name are trusted
transitively, because the registry names them inside its signed
payload. Operators wanting tighter control can additionally inspect
`Endorsement.Verifier` against their own policy.

Endorsement requirements attach to the **sending** role. Verifiers
resolving a *receiver's* identity MUST NOT demand an authority
countersignature: receiving addresses are self-provisioned and MAY be
entirely self-signed, or publish nothing at all (204). The trust
anchor in §11.1 — TLS-bound `iss` — does not depend on authorities;
authorities are an additional, opt-in policy layer on top.

### 6.4 Sender verification

`Client.VerifySender(ctx, addr, requireVerified)` combines the two:
it resolves `addr`'s identity via `Who` and requires an Authority
countersignature via `VerifyAuthority`, returning the endorsed
`org.Party`. When `requireVerified` is true the endorsement must
additionally carry a confirmed verifier (§5.3), else
`ErrNotVerified` — registration suffices for most exchanges;
verification may be demanded before acting on inbox deliveries from
new counterparties. Receiving inboxes call it with the verified
issuer of an incoming envelope before accepting the delivery
(§8.4). A receive-only account (204) or a merely self-signed
identity fails this check by construction — such addresses can
receive but cannot act as approved senders.

## 7. Discovery Transport

### 7.1 HTTP Client Defaults

The default `HTTPFetcher` enforces:

| Parameter             | Value              |
|-----------------------|--------------------|
| Request timeout       | 10 seconds         |
| Dial timeout          | 5 seconds          |
| Maximum response size | 1 MiB              |
| Required `Accept`     | `application/json` |
| Required scheme       | `https`            |
| Required status       | `200 OK` (`204` → `ErrNoContent`, `202` → `ErrPending`) |
| Key cache TTL         | 5 minutes          |

Responses larger than 1 MiB are truncated. A `204 No Content`
response causes `ErrNoContent` and a `202 Accepted` response causes
`ErrPending` — distinct sentinels because at the who endpoint both
empty responses are meaningful (§8.2). Transient conditions — 429,
any 5xx, or a transport failure — cause the retryable
`ErrUnavailable`; any other non-200 response causes the permanent
`ErrFetchFailed`.

A `Client` configured with an identity (`net.WithIdentity`) mints a
fresh request token (§5.5) per who or inbox request and sends it as
`Authorization: Bearer`. Key fetches are always sent bare — key
endpoints are open (§8.1).

`Client.FetchKey` caches fetched keys per URL for a short TTL (5
minutes by default, tunable with `net.WithKeyCacheTTL`; zero
disables), since published keys are immutable per `kid` (§5.5).
The cache holds successes only and is capped in size, so hostile
`iss`/`kid` values cannot grow it without bound.

**SSRF defense.** The fetcher's transport refuses to dial any host
whose resolved IP is loopback, private (RFC 1918 / RFC 6598),
link-local, multicast, or unspecified. A signed `iss` URI is
attacker-controlled, so the FQDN it names could resolve to an
internal service (e.g. AWS metadata at `169.254.169.254`, a
container's localhost, a corporate intranet) — refusing those at
dial time prevents the verifier from being used as an SSRF gadget.
There is no public escape hatch; in-process test fixtures (e.g.
`httptest`) should inject their own `Fetcher` via
`net.WithFetcher`.

### 7.2 Pluggable Fetcher

The `Fetcher` interface (`Fetch(ctx, url, header) ([]byte, error)`
and `Post(ctx, url, body, header) error`) allows substituting the
HTTP transport, e.g. for testing, in-process resolution, or
alternative transports; `header` carries any `Authorization`
request token the `Client` has minted for the request, and `Post`
is what `Client.Send` (§8.3) delivers envelopes through.
Verification-only transports may implement `Post` as a plain error
return. Use `net.WithFetcher(f)` when constructing a `Client`.

## 8. Server-Side Endpoints

### 8.1 `GET /.well-known/gobl/keys/<kid>`

Open. Returns a single RFC 7517 JSON Web Key as `application/json` —
the file `<domain>/keys/<kid>.json`, optionally carrying GOBL Net's
`valid_from` / `valid_until` extension members (see §4). Unknown kid
returns `404 Not Found`. No bulk endpoint is exposed.

### 8.2 `GET /.well-known/gobl/who`

Authenticated. The caller MUST present a request token (§5.5); a
request without a valid token is rejected with `401 Unauthorized`.
Returns the domain's party envelope: an `org.Party` document,
self-signed with `iss=self` as the **first** signature (no
`aud`), optionally carrying Authority countersignatures. The
response body is a static signed document — the same bytes for
every authorized caller; the token controls access and feeds the
audit log (§11.9), it does not change the response.

| Status            | Cause                                                |
|-------------------|------------------------------------------------------|
| `200 OK`          | Returns the signed party envelope.                   |
| `202 Accepted`    | Request authenticated and recorded; the owner discloses selectively. If the owner approves, it delivers its party envelope to the requester's inbox (the token's `iss`), signed `iss=owner`, `aud=requester` (§8.3). There is no guarantee and no deadline — requesters MUST NOT wait synchronously and proceed with the details they already hold. |
| `204 No Content`  | The account exists but publishes no identity details. A 204 account is receive-only: it cannot pass sender verification (§6.4), but deliveries *to* it are unaffected. |
| `401 Unauthorized`| Missing or invalid request token (§5.5).             |
| `404 Not Found`   | The address does not participate in GOBL Net.        |
| `503 Service Unavailable` | The requester's key endpoint could not be reached to verify the token (§5.5); retry later. |

`202` and `204` express different stances: `204` means "there is
nothing to share, ever" — the account publishes no identity details
to anyone; `202` means "there may be something to share, but the
owner decides per requester". Deferred disclosure is not the
default — most participants serve `200` to any authenticated
caller — but it suits individuals (B2C) whose party details are
personal data they would rather not hand to every requester
(§11.8).

Because the response requires authentication it MUST NOT be served
from shared public caches: servers SHOULD send
`Cache-Control: private` (adding `Vary: Authorization` where an
intermediary cache is in play). Clients MAY cache the *verified*
party envelope locally instead. Keep the TTL modest (minutes to a
few hours): a cached who bounds how quickly an Authority's
revocation of an endorsement is observed (§11.4). A consequence of
mandatory authentication is that `/who` and `/inbox` can no longer
be served by a static file host; the key endpoints (§8.1) still
can.

Key discovery (§8.1) is independent of who visibility and remains
open: a participant that signs anything MUST serve its published
keys regardless of whether its who returns 200, 202 or 204. Open
key endpoints are also what keep request-token verification from
recursing (§5.5).

> **Note.** Earlier drafts specified `/who` as an authenticated POST
> exchange so the target could pre-approve requesters and log who
> asked for its details; an interim revision replaced it with an
> open GET. This revision restores negotiated disclosure in GET
> form: the request token supplies the requester identity the POST
> body used to carry, and the `202` path supplies the pre-approval
> hook.

### 8.3 `POST /.well-known/gobl/inbox`

Accepts a signed GOBL Envelope. The request MUST carry a request
token (§5.5) identifying the transmitting party; a request without
a valid token is rejected with `401 Unauthorized` before the body
is considered. The token's `iss` MAY differ from the envelope's
signed `iss`: a trusted intermediary (e.g. a hosted provider
transmitting documents its customers signed) authenticates the
*request* with its own identity while the envelope carries the
document signer's. The two layers are independent and both
required.

The envelope layer is unchanged by the token: the signer (`iss`) is
verified against its published key (fetched from
`<iss>/.well-known/gobl/keys/<kid>`); the signed `aud` MUST be
present and MUST equal this inbox's Address — envelopes signed
without an audience, or bound to a different audience, MUST be
rejected. The inbox SHOULD then apply its sender-endorsement
policy: resolve the sender's who (§6.4) and require an Authority
countersignature — optionally with a confirmed verifier (§5.3) when
the operator demands verified identities. Endorsement attaches to
the envelope's signer, never to the token's issuer. Status codes:

| Status                       | Cause                                                |
|------------------------------|------------------------------------------------------|
| `202 Accepted`               | Envelope parsed, validated, signature verified, persisted. |
| `400 Bad Request`            | Body could not be read or did not decode as JSON.    |
| `401 Unauthorized`           | Missing or invalid request token (§5.5); envelope signature did not verify; or `aud` missing / not equal to this inbox. |
| `403 Forbidden`              | Sender (`iss`) is not endorsed by an Authority this inbox trusts, or lacks the verified status the operator requires (§5.3). |
| `422 Unprocessable Entity`   | Envelope failed structural validation.               |
| `500 Internal Server Error`  | Persistence failed.                                  |
| `503 Service Unavailable`    | A key endpoint needed to verify the token or envelope could not be reached (§5.5); retry later. |

The request body size is capped at 1 MiB. The protocol does not
mandate where or how an accepted envelope is persisted — that is an
implementation decision for the operator running the inbox (a
filesystem directory, an object store, a database row, a message
queue, an in-process handler, …). A 202 simply MUST mean the
envelope is durable enough that the server is willing to confirm
receipt. The reference server in
[`gobl.dev`](https://github.com/invopop/gobl.dev/#gobl-net) writes
each accepted envelope to `<config>/<domain>/inbox/<envelope-uuid>.json`;
that layout is one valid choice, not a protocol requirement.

The empty `202 Accepted` body is deliberate: it asserts durable
persistence by the inbox and nothing more. There are no transport
receipts in GOBL Net — experience with receipt-bearing networks
shows they attest to a handoff between intermediaries, not to
anything the recipient actually did. Business-level processing is
signalled instead by follow-up documents (e.g. a `bill.Status`
envelope) sent back through the same inbox mechanism — which, like
any other sending, requires the responding party to be an endorsed
sender (§6.4). Receive-only participants do not send status; a
supplier delivering to consumers should treat delivery as
fire-and-forget beyond the 202.

Inbox delivery MUST be idempotent on the envelope's identity: a
re-POST of an envelope with the same `uuid` and `dig` returns `202`
without creating a duplicate record. Since a status response may
never come, "retry until 202" is the sender's correct recovery
strategy, and idempotency is what makes it safe. (The signed `aud`
already prevents the same envelope being replayed against a
*different* inbox.)

**Deferred who deliveries.** A party envelope (an `org.Party`
document, self-signed by its subject) arriving at an inbox whose
operator has an outstanding who request to that address (§8.2,
`202`) SHOULD be accepted without sender endorsement: it carries
exactly the assertions a `200` who response would, so it needs no
more trust than the lookup it answers. Verification is the same as
`Client.Who` steps 2–4 (§6.2), except that the envelope's signed
`aud` MUST name this inbox. Outside of an outstanding who request,
party envelopes are subject to the normal endorsement policy.

In code, `Client.Send(ctx, addr, env)` performs the delivery: it
mints a request token for the target inbox, POSTs the envelope, and
returns nil on `202`, `ErrInboxRejected` on a definitive 4xx (do
not retry — the inbox has decided), and the retryable
`ErrUnavailable` on 429, 5xx, or transport failure. "Retry on
`ErrUnavailable` until `202`" is the sender's recovery strategy.

### 8.4 End-to-end delivery flow

The canonical invoice exchange between a Supplier (sender) and a
Customer (receiver — a business or an individual consumer):

1. Supplier needs to send an invoice to Customer, whose GOBL Net
   Address it learned out-of-band (checkout, contract, directory).
2. *Optionally*, Supplier performs `GET /who` on Customer — with a
   request token naming itself (§5.5) — to fetch complete invoicing
   details. A `204 No Content` means the account exists but shares
   nothing; a `202 Accepted` means the Customer may deliver its
   details to the Supplier's inbox later. Either way the Supplier
   proceeds with the details it already holds.
3. Supplier signs the invoice envelope with `iss=supplier`,
   `aud=customer` and POSTs it to Customer's `/inbox`, again
   presenting a request token. When a service provider transmits on
   the Supplier's behalf, the token names the provider while the
   envelope keeps `iss=supplier`.
4. Customer verifies the request token, then the envelope signature
   and audience (§6.1), then resolves the Supplier's identity with
   `GET /who` on the verified issuer (a cached copy MAY be used
   within its TTL).
5. Customer accepts iff the Supplier's who carries a countersignature
   from an Authority the Customer trusts (§6.4) — optionally
   requiring a confirmed verifier (§5.3) — i.e. the sender is
   registered, and KYC'd where demanded.
6. `202 Accepted` (empty body) confirms reception to the Supplier.

Only step 5 involves any registration requirement, and it applies
solely to the sending party. The Customer's own address never needs
Authority approval — steps 1–3 work against any self-provisioned
endpoint.

Sending is one role, whoever performs it: any follow-up document
the Customer might send back — a `bill.Status` reporting
acceptance or payment, a corrective exchange — is itself a
delivery under this flow and requires the Customer to be an
endorsed sender. A receive-only Customer sends nothing, and the
Supplier expects nothing: for B2C, the 202 in step 6 is the end of
the exchange. Status flows are a feature of exchanges between
registered participants.

## 9. Reference Implementation

GOBL Net's reference client lives in this package (`net.Client`,
`net.Address`, `net.Authorities`). The reference server and the
operator-facing CLI (`gobl init`, `gobl net who/send/serve`,
`gobl sign --domain …`, `gobl verify --remote`) live in
[`gobl.dev`](https://github.com/invopop/gobl.dev/#gobl-net), which also
documents the server's on-disk layout, ACME setup, multi-tenant
routing, structured logging, and Docker deployment. The network's
default registration Authority, `lookup.gobl.org` (§6.3), is
implemented in
[`gobl.lookup`](https://github.com/invopop/gobl.lookup).

The protocol is transport-defined; any conforming implementation can
serve the well-known endpoints from §8 over HTTPS.

## 10. Errors

The package exports the following sentinel errors:

| Error                  | Cause                                                            |
|------------------------|------------------------------------------------------------------|
| `ErrAddressEmpty`      | Empty input to `ParseAddress`.                                   |
| `ErrAddressInvalid`    | Input is not a valid FQDN per §3.1.                              |
| `ErrFetchFailed`       | Well-known resource fetch failed for a reason retrying will not cure (definitive non-2xx, malformed content, invalid input). |
| `ErrUnavailable`       | A well-known resource could not be reached (network failure, 429, 5xx) — a transient, retryable condition. Servers respond 503, not 401, when token verification hits it. |
| `ErrNoContent`         | The resource exists but has no content (HTTP 204) — a who endpoint that publishes no identity details. |
| `ErrPending`           | A who request was accepted for deferred disclosure (HTTP 202) — the owner may deliver its party envelope to the requester's inbox later. |
| `ErrTokenInvalid`      | A request token failed verification (parse, signature, key fetch, `aud` mismatch, key validity window). |
| `ErrTokenExpired`      | A request token failed the freshness check (`exp` passed, `iat` too old or in the future). |
| `ErrVerifyFailed`      | Envelope verification failed (no signature, invalid `iss`, key fetch failed, signature mismatch, `aud` mismatch, `iat` outside the key's validity window, who issuer/address mismatch). |
| `ErrUnknownAuthority`  | An endorser on a `/who` envelope is not in `Authorities` (only raised by callers that opt into authority enforcement). |
| `ErrNotVerified`       | A sender's endorsement is valid but carries no confirmed verifier (§5.3) and the caller required identity verification. |
| `ErrSignatureExpired`  | An authority countersignature verified but its `exp` claim has passed. |
| `ErrPartyMissing`      | A `/who` response did not contain an `org.Party` document.       |
| `ErrInboxRejected`     | A receiving inbox did not return 202.                            |

All callers using `errors.Is` against these sentinels MUST continue to
work after wrapping with `fmt.Errorf("%w: ...", err)`.

## 11. Security Considerations

### 11.1 Trust Model

A signature's verifiable origin is the signed `iss` URI. Verification
fetches the public key from
`https://<iss-fqdn>/.well-known/gobl/keys/<kid>`, and the HTTPS
connection's certificate proves the response really came from that
FQDN. The trust anchor is therefore the Web PKI binding of the
Address to the entity that controls its TLS certificate.

A forged `iss` pointing to an attacker-controlled host can only
produce a signature that verifies against an attacker-controlled key
served from that host — distinguishable from the expected identity
at the application layer. Callers that already know the expected
Address SHOULD pass it as `expectedAud` to `Client.VerifyEnvelope`,
or compare the returned issuer Address against the expected sender
before acting on the document.

### 11.2 TLS

All well-known endpoints are served over HTTPS in production. Verifiers
rely on the host's TLS certificate to establish that the served content
originates from the named Address. Operators MUST ensure that TLS
certificate issuance is properly controlled for any Address they intend
to use as an identity.

### 11.3 Authority Trust

`Authorities` is an opt-in allowlist of FQDNs (see §6.3). For
verifiers that *do* enforce it, the threat is symmetrical to §11.1
but multiplied: an attacker who gains the ability to issue a TLS
certificate for any address in `Authorities` can serve a forged
`/who` for any participant whose envelopes the verifier accepts on
the strength of that authority's countersignature. Where this hook
is used, the authority list MUST be kept short and reviewed
regularly.

### 11.4 Inbox Authentication and Sender Endorsement

The inbox endpoint authenticates at two independent layers: the
request token (§5.5) identifies the transmitting peer before the
body is read, and the envelope signature identifies the document's
signer. The sender-endorsement policy (§6.4, §8.4) is the layer
that restricts acceptance to Authority-registered participants.
Operators choose the trusted Authority list and whether to demand
verified identities (§5.3).

Endorsement revocation is bounded by two mechanisms. First, the
countersignature's own `exp` claim (§5.3): with the recommended
90-day maximum, a revoked sender's endorsement dies on its own even
if the stale who keeps being served — and the renewal cycle this
forces (re-registering before expiry, in the spirit of Let's
Encrypt certificates) gives the Authority a recurring checkpoint at
which to decline. Second, who caching: once an Authority stops
countersigning a sender (or the sender re-publishes its who without
the countersignature), receivers continue to accept only until
their cached copy expires. Inboxes enforcing endorsement MUST
resolve the sender's who with a bounded TTL rather than trusting a
cached copy indefinitely, and MAY additionally apply a max-age
policy to the `iat` of countersignatures that carry no `exp`.

Verification revocation has a third, faster path: because the
`verifier` pointer lives inside the *registration* countersignature
(§5.3), the registry can withdraw verified status at any renewal —
or immediately, by re-signing without the pointer — without waiting
for the verifier's own longer-lived signature to expire. The
long-lived verifier signature is only evidence when the short-lived
registration signature points at it.

### 11.5 Response Size and Signature Count

The 1 MiB cap on Key, Who, and inbox bodies limits memory amplification
from hostile or misconfigured peers. `VerifyAuthority` additionally
refuses envelopes carrying more than 32 signatures: each candidate
countersignature can cost a network key fetch, so an unbounded
signature list would turn a hostile who envelope into a
fetch-amplification gadget against verifying inboxes.

### 11.6 SSRF / non-public dial targets

A signed `iss` URI is supplied by the signer, so a verifier that
blindly fetches `https://<iss-fqdn>/.well-known/...` could be
tricked into hitting an internal address controlled by the attacker
(loopback, cloud metadata at `169.254.169.254`, an RFC 1918 host on
the verifier's network). The default `HTTPFetcher` defends against
this by resolving the host at dial time and refusing any loopback,
private, link-local, multicast, or unspecified address (see §7.1).
Operators that swap in a custom `Fetcher` SHOULD apply equivalent
checks.

### 11.7 Address Canonicalization

`ParseAddress` lowercases and strips trailing dots before validation, so
two visually distinct strings such as `Example.COM.` and `example.com`
normalize to the same Address. Callers MUST use `ParseAddress` (directly
or via `Address.Validate`) before comparing addresses for equality.

### 11.8 Identity Privacy

`GET /who` requires a request token (§5.5), so published party
details are readable only by callers holding a verifiable GOBL Net
identity — never anonymously. Bulk harvesting therefore carries a
cost the open GET did not: every request names a resolvable domain
identity and MAY be recorded (§11.9). Sending participants
generally publish business details that are public record anyway;
individuals and other receive-only participants SHOULD respond
`204` rather than publish personal data, or use deferred
disclosure (`202`, §8.2) to decide per requester. Servers MAY
additionally rate-limit the endpoint. Hosted providers offering
per-user sub-addresses (`alice.inbox.provider.example`) hold the
keys and the inbox contents for their users; that concentration of
trust is a deployment choice outside this protocol, but users of
such services should understand the provider can read anything
delivered to them.

### 11.9 Request audit logs

Servers MAY keep an audit log of authenticated requests — the
token's `iss`, `jti` and `iat`, the endpoint, and the outcome. The
guarantee that every request is attributable to a requester is a
design goal of §5.5, and the audit log is how operators realize
it. Operators should remember the log is itself personal data:
requester addresses reveal business relationships (who is invoicing
whom, and when). Retention SHOULD be bounded, and logs SHOULD NOT
be published.

## 12. References

- RFC 2119 — Key words for use in RFCs to Indicate Requirement Levels
- RFC 6750 — OAuth 2.0 Bearer Token Usage (`Authorization: Bearer`)
- RFC 7515 — JSON Web Signature (JWS)
- RFC 7517 — JSON Web Key (JWK)
- RFC 7518 — JSON Web Algorithms (JWA)
- RFC 7519 — JSON Web Token (JWT)
- RFC 8174 — Ambiguity of Uppercase vs Lowercase in RFC 2119 Key Words
- `github.com/invopop/gobl/dsig` — Signature, key, and digest primitives
- `github.com/invopop/gobl/org` — `org.Party` schema
