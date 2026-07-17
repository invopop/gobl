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
  with a plain GET.
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

The protocol layers on top of standards already in use by GOBL:

- **RFC 7515** — JSON Web Signature, for envelope signing.
- **RFC 7517** — JSON Web Key Set, for key discovery.
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
  receiver only needs a domain, TLS, and (if it signs) published
  keys. A sender is additionally expected to carry an Authority
  endorsement on its who identity (§6.4).
- **iss / aud** — Fields in a signature's *signed payload* carrying the
  verifiable GOBL Net origin (`iss`) and the address the signature is
  bound to (`aud`), both as `gobl:` `cbc.URI` values. These are the
  authoritative, tamper-proof identities.
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
gobl:<A>                                   ← Address.URI()  (iss / aud value)
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

### 3.3 Topic Derivation

`Address.Topic()` returns the address with its labels reversed and joined
by dots: `billing.invopop.com` becomes `com.invopop.billing`. The topic
form is provided for use by notification fan-out implementations; it is
not consumed by any code in this package.

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
- `iss` is the signer's verifiable GOBL Net address as a `gobl:` URI
  (e.g. `gobl:billing.invopop.com`). The verifier reads it to
  discover *which* per-key endpoint to fetch.
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

### 5.3 Authority countersignatures and scope

A `/who` response (and any envelope an Authority chooses to endorse)
MAY carry one or more countersignatures from addresses in the
verifier's `Authorities` list (§6.3). An Authority countersignature
SHOULD include a `scope` claim in its signed payload declaring the
level of confidence the Authority is asserting:

| Scope (`cbc.Key`) | Meaning                                                                                        |
|-------------------|------------------------------------------------------------------------------------------------|
| (omitted)         | No assertion beyond the cryptographic facts of the signature itself.                           |
| `registered`      | The Authority has confirmed the subject controls the named GOBL Net address (FQDN ownership).  |
| `verified`        | The Authority has performed full identity verification (e.g. KYC) in addition to registration. |

Scope is set with the `head.WithScope` signing option:

```go
env.Sign(authorityKey,
    head.WithIssuer(authorityAddr.URI()),
    head.WithAudience(subjectAddr.URI()),
    head.WithScope(head.ScopeVerified))
```

Verifiers MAY require a minimum scope per use case (`registered`
suffices for `/who` discovery; `verified` may be required before
acting on inbox deliveries from new counterparties). The minimum is
passed to `Client.VerifyAuthorityWithScope` / `Client.VerifySender`:
`registered` is satisfied by `registered` or `verified`, and a
custom scope key is satisfied only by an exact match; plain
`Client.VerifyAuthority` accepts any authority signature. Operators MAY define additional scope keys;
`registered` and `verified` are the baseline values the protocol
recognises.

An Authority countersignature SHOULD carry an `exp` claim bounding
the endorsement's lifetime — 90 days is the recommended maximum,
after which the subject renews its registration (§11.4). Verifiers
MUST treat an expired countersignature as absent: `VerifyAuthority`
and `VerifyAuthorityWithScope` reject candidates whose `exp` has
passed with `ErrSignatureExpired`. A countersignature without `exp`
does not expire; verifiers MAY apply their own `iat` max-age policy
to such legacy endorsements.

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

## 6. Verification

### 6.1 Envelope Verification Flow

`Client.VerifyEnvelope(ctx, env, expectedAud)` returns the verified
issuer address:

1. The envelope MUST be signed; otherwise `ErrVerifyFailed`.
2. The first signature's signed payload is read; `iss` MUST be a `gobl:`
   URI (else `ErrVerifyFailed`).
3. `FetchKey(ctx, iss-host, kid)` fetches the issuer's published key
   from `/.well-known/gobl/keys/<kid>` (including its optional
   `valid_from` / `valid_until`).
4. The envelope is verified against that public key.
5. If `expectedAud` is non-empty, the signed `aud` MUST equal it.
6. If the key declares a validity window, the signed `iat` MUST fall
   within `[valid_from, valid_until]` (each bound optional).
7. The verified issuer address is returned.

### 6.2 Identity lookup (`GET /who`)

`/who` is an open, unauthenticated GET (see §8.2). The response is
the target's party envelope: document = the target's `org.Party`,
first signature = the target's self-signature with `iss=gobl:target`
and no `aud` (a GET response has no caller to bind to), optionally
followed by Authority countersignatures.

`Client.Who(ctx, addr)` performs the lookup and verifies it:

1. `GET https://<addr>/.well-known/gobl/who`. A `204` returns
   `ErrNoContent` — the address exists but publishes no identity
   details (a receive-only account).
2. The response envelope's first signature is verified via
   `VerifyEnvelope` (the signed `iss` resolved to a published key).
3. The verified issuer MUST equal the fetched address — a valid
   envelope for a *different* identity served at this URL is
   rejected.
4. The document MUST be an `org.Party`, else `ErrPartyMissing`.

Because the response is a static signed document, it is cacheable
(§8.2). The identity's integrity comes from the self-signature and
the TLS origin, not from any per-request binding.

### 6.3 Trusted Authorities

The package-level slice `net.Authorities` holds GOBL Net addresses
treated as trusted KYC vendors. The default list is empty;
`net.RegisterAuthority` or the `WithAuthorities` client option add to
it.

`Client.VerifyAuthority(ctx, env)` returns nil iff the envelope
carries at least one signature whose signed `iss` is in the client's
authorities AND that signature cryptographically verifies against
the authority's published key AND its signed `exp` claim (if any)
has not passed. `Client.VerifyAuthorityWithScope(ctx, env, minScope)`
additionally requires the signed scope claim to satisfy `minScope`
(§5.3). They return `ErrUnknownAuthority` when no candidate
signature is from a known authority (or none has been registered),
`ErrSignatureExpired` when a verified authority signature has
expired, `ErrScopeInsufficient` when it falls short of the minimum
scope, and `ErrVerifyFailed` when a candidate fails its crypto
check.

Endorsement requirements attach to the **sending** role. Verifiers
resolving a *receiver's* identity MUST NOT demand an authority
countersignature: receiving addresses are self-provisioned and MAY be
entirely self-signed, or publish nothing at all (204). The trust
anchor in §11.1 — TLS-bound `iss` — does not depend on authorities;
authorities are an additional, opt-in policy layer on top.

### 6.4 Sender verification

`Client.VerifySender(ctx, addr, minScope)` combines the two: it
resolves `addr`'s identity via `Who` and requires an Authority
countersignature satisfying `minScope` via
`VerifyAuthorityWithScope`, returning the endorsed `org.Party`. Receiving inboxes call it with
the verified issuer of an incoming envelope before accepting the
delivery (§8.4). A receive-only account (204) or a merely self-signed
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
| Required status       | `200 OK` (`204` → `ErrNoContent`) |

Responses larger than 1 MiB are truncated. A `204 No Content`
response causes `ErrNoContent` — a distinct sentinel because at the
who endpoint an empty response is meaningful (§8.2). Any other
non-200 response causes `ErrFetchFailed`.

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

The `Fetcher` interface (`Fetch(ctx, url) ([]byte, error)`) allows
substituting the HTTP transport, e.g. for testing, in-process resolution,
or alternative transports. Use `net.WithFetcher(f)` when constructing a
`Client`.

## 8. Server-Side Endpoints

### 8.1 `GET /.well-known/gobl/keys/<kid>`

Open. Returns a single RFC 7517 JSON Web Key as `application/json` —
the file `<domain>/keys/<kid>.json`, optionally carrying GOBL Net's
`valid_from` / `valid_until` extension members (see §4). Unknown kid
returns `404 Not Found`. No bulk endpoint is exposed.

### 8.2 `GET /.well-known/gobl/who`

Open. Returns the domain's party envelope: an `org.Party` document,
self-signed with `iss=gobl:self` as the **first** signature (no
`aud`), optionally carrying Authority countersignatures. The response
is a static signed document — the same bytes for every caller.

| Status            | Cause                                                |
|-------------------|------------------------------------------------------|
| `200 OK`          | Returns the signed party envelope.                   |
| `204 No Content`  | The account exists but publishes no identity details. A 204 account is receive-only: it cannot pass sender verification (§6.4), but deliveries *to* it are unaffected. |
| `404 Not Found`   | The address does not participate in GOBL Net.        |

Because the response is static, servers SHOULD set standard HTTP
caching headers (`Cache-Control`, `ETag`) and clients MAY cache
accordingly. Keep the TTL modest (minutes to a few hours): a cached
who bounds how quickly an Authority's revocation of an endorsement is
observed (§11.4).

Key discovery (§8.1) is independent of who visibility: a participant
that signs anything MUST serve its published keys regardless of
whether its who returns 200 or 204.

> **Note.** Earlier drafts specified `/who` as an authenticated POST
> exchange so the target could pre-approve requesters and log who
> asked for its details. That negotiated disclosure MAY return in a
> future revision as an optional POST alongside the GET; the GET is
> the baseline every participant serves.

### 8.3 `POST /.well-known/gobl/inbox`

Accepts a signed GOBL Envelope. The signer (`iss`) is verified against
its published key (fetched from `<iss>/.well-known/gobl/keys/<kid>`);
the signed `aud` MUST be present and MUST equal this inbox's
Address — envelopes signed without an audience, or bound to a
different audience, MUST be rejected. The inbox SHOULD then apply its
sender-endorsement policy: resolve the sender's who (§6.4) and
require an Authority countersignature at the operator's minimum
scope. Status codes:

| Status                       | Cause                                                |
|------------------------------|------------------------------------------------------|
| `202 Accepted`               | Envelope parsed, validated, signature verified, persisted. |
| `400 Bad Request`            | Body could not be read or did not decode as JSON.    |
| `401 Unauthorized`           | Envelope signature did not verify, or `aud` missing / not equal to this inbox. |
| `403 Forbidden`              | Sender (`iss`) is not endorsed by an Authority this inbox trusts, or falls short of its minimum scope. |
| `422 Unprocessable Entity`   | Envelope failed structural validation.               |
| `500 Internal Server Error`  | Persistence failed.                                  |

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

### 8.4 End-to-end delivery flow

The canonical invoice exchange between a Supplier (sender) and a
Customer (receiver — a business or an individual consumer):

1. Supplier needs to send an invoice to Customer, whose GOBL Net
   Address it learned out-of-band (checkout, contract, directory).
2. *Optionally*, Supplier performs `GET /who` on Customer to fetch
   complete invoicing details. A `204 No Content` means the account
   exists but shares nothing — the Supplier proceeds with the details
   it already holds.
3. Supplier signs the invoice envelope with `iss=gobl:supplier`,
   `aud=gobl:customer` and POSTs it to Customer's `/inbox`.
4. Customer verifies the envelope signature and audience (§6.1), then
   resolves the Supplier's identity with `GET /who` on the verified
   issuer (a cached copy MAY be used within its TTL).
5. Customer accepts iff the Supplier's who carries a countersignature
   from an Authority the Customer trusts, at the Customer's minimum
   scope (§6.4) — i.e. the sender is registered/KYC'd.
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
routing, structured logging, and Docker deployment.

The protocol is transport-defined; any conforming implementation can
serve the well-known endpoints from §8 over HTTPS.

## 10. Errors

The package exports the following sentinel errors:

| Error                  | Cause                                                            |
|------------------------|------------------------------------------------------------------|
| `ErrAddressEmpty`      | Empty input to `ParseAddress`.                                   |
| `ErrAddressInvalid`    | Input is not a valid FQDN per §3.1.                              |
| `ErrFetchFailed`       | Well-known resource fetch failed (network, non-200/204, malformed). |
| `ErrNoContent`         | The resource exists but has no content (HTTP 204) — a who endpoint that publishes no identity details. |
| `ErrVerifyFailed`      | Envelope verification failed (no signature, non-`gobl:` `iss`, key fetch failed, signature mismatch, `aud` mismatch, `iat` outside the key's validity window, who issuer/address mismatch). |
| `ErrUnknownAuthority`  | An endorser on a `/who` envelope is not in `Authorities` (only raised by callers that opt into authority enforcement). |
| `ErrScopeInsufficient` | An authority countersignature verified but its scope claim falls short of the required minimum. |
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
is used, the KYC vendor list MUST be kept short and reviewed
regularly.

### 11.4 Inbox Authentication and Sender Endorsement

The inbox endpoint verifies the sender's signature before persisting
an envelope; the sender-endorsement policy (§6.4, §8.4) is the layer
that restricts acceptance to Authority-registered participants.
Operators choose the minimum scope and the trusted Authority list.

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

### 11.5 Response Size

The 1 MiB cap on Key, Who, and inbox bodies limits memory amplification
from hostile or misconfigured peers.

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

`GET /who` makes the published party details readable by anyone who
knows the address. Sending participants generally publish business
details that are public record anyway; individuals and other
receive-only participants SHOULD respond `204` rather than publish
personal data. Servers MAY rate-limit the endpoint to slow bulk
harvesting. Hosted providers offering per-user sub-addresses
(`alice.inbox.provider.example`) hold the keys and the inbox contents
for their users; that concentration of trust is a deployment choice
outside this protocol, but users of such services should understand
the provider can read anything delivered to them.

## 12. References

- RFC 2119 — Key words for use in RFCs to Indicate Requirement Levels
- RFC 7515 — JSON Web Signature (JWS)
- RFC 7517 — JSON Web Key (JWK)
- RFC 7518 — JSON Web Algorithms (JWA)
- RFC 8174 — Ambiguity of Uppercase vs Lowercase in RFC 2119 Key Words
- `github.com/invopop/gobl/dsig` — Signature, key, and digest primitives
- `github.com/invopop/gobl/org` — `org.Party` schema
