# GOBL Net

> **Experimental:** the API and wire protocol may change without notice.

**Status:** Draft

## Abstract

GOBL Net is a decentralized discovery and delivery protocol for signed GOBL
documents. A participant is identified by a fully qualified domain name
(FQDN). HTTPS binds that address to three well-known resources:

| Resource | Purpose |
|---|---|
| `/.well-known/gobl/keys/<kid>` | Publish one verification key. |
| `/.well-known/gobl/who` | Publish a signed `org.Party` identity. |
| `/.well-known/gobl/inbox` | Accept signed GOBL envelopes. |

A signature names its issuer in the signed `iss` claim. A verifier obtains the
key from the issuer's domain and relies on Web PKI to authenticate the HTTPS
response. Registration authorities add an optional trust layer by
countersigning party envelopes.

This package implements the client, address, authentication, and verification
parts of the protocol. The reference server and CLI live in
[`github.com/invopop/gobl.dev`](https://github.com/invopop/gobl.dev/#gobl-net).

## 1. Conventions

The terms MUST, MUST NOT, SHOULD, SHOULD NOT, and MAY are interpreted as
described in BCP 14 (RFC 2119 and RFC 8174).

## 2. Terminology

| Term | Meaning |
|---|---|
| **Participant** | An entity identified by a GOBL Net address. A participant may hold different roles in different exchanges. |
| **Subject** | The participant described by a party envelope. Its address comes from the party's `gobl:` endpoint. |
| **Issuer** | The participant identified by a signature's `iss` claim. This is a signature-level role. |
| **Audience** | The participant identified by a signature's `aud` claim. The signature is bound to this address. |
| **Sender** | The issuer of a valid document signature whose audience is the receiving inbox. |
| **Receiver** | The participant operating the inbox to which a document is delivered. For ordinary deliveries, the receiver is also the signature's audience. |
| **Requester** | The issuer of an HTTP request token. A requester may be an intermediary and need not be the document sender. |
| **Authority** | A trusted participant that confirms the subject controls its GOBL Net address and countersigns its party envelope. |
| **Verifier** | A participant, independent of the Authority, that performs KYC/KYB and countersigns the subject's party envelope. |
| **Endorsement** | A valid Authority countersignature, optionally supported by a confirmed Verifier countersignature. |
| **Party envelope** | A signed GOBL envelope containing an `org.Party`. It identifies its Subject and may carry Authority and Verifier countersignatures. |

Issuer, Sender, and Requester often refer to the same participant in a direct
delivery, but the protocol does not assume they are the same. Audience and
Receiver are equivalent for ordinary document delivery; in other signatures,
such as an Authority endorsement, the Audience is the Subject rather than the
recipient of an HTTP request.

An Authority and Verifier MUST be different participants. Other roles may
overlap where their definitions permit it.

## 3. Addresses and discovery

A GOBL Net address is an FQDN without a scheme, port, path, query, or fragment.
`ParseAddress`:

1. trims surrounding whitespace and one trailing dot;
2. converts internationalized names to lowercase ASCII A-labels using the
   IDNA2008 lookup profile;
3. validates DNS label syntax; and
4. requires at least two labels.

An empty value returns `ErrAddressEmpty`; other invalid values return
`ErrAddressInvalid`. Protocol fields and comparisons MUST use the canonical
ASCII form.

For address `supplier.example.com`:

```text
supplier.example.com
gobl:supplier.example.com
https://supplier.example.com/.well-known/gobl/keys/<kid>
https://supplier.example.com/.well-known/gobl/who
https://supplier.example.com/.well-known/gobl/inbox
https://supplier.example.com/.well-known/jwks.json
```

The first form is used in signed `iss`, `aud`, and `verifier` claims. The
`gobl:` URI is used in `org.Party` endpoints and optional envelope routing
fields. Production endpoints MUST use HTTPS.

## 4. Keys

Each signing key MUST be available as a single RFC 7517 JWK at
`/.well-known/gobl/keys/<kid>`. Unknown or removed key IDs return `404`. The
optional `/.well-known/jwks.json` endpoint exposes all keys as an RFC 7517 JWK
Set for generic tooling; GOBL Net verification uses the per-key endpoint.

A published `dsig.PublicKey` may include `valid_from` and `valid_until`.
When a signature has an `iat`, verification rejects it if the signing time is
outside those bounds. A missing bound or missing `iat` skips the corresponding
comparison. Retired keys SHOULD remain available when historical signatures
must remain verifiable.

Signatures identify keys with `kid`; they do not carry a `jku` header.

## 5. Envelope signatures

The signed payload contains the document `uuid` and digest plus these protocol
claims:

| Claim | Meaning |
|---|---|
| `iss` | Signer's canonical GOBL Net address. |
| `aud` | Address to which the signature is bound; optional except for ordinary inbox deliveries. |
| `iat` | Signing time, as an RFC 7519 NumericDate. |
| `exp` | Optional expiry, used primarily for endorsements. |
| `verifier` | Optional verifier named by a registration authority. |

Envelope header `from` and `to` fields are unsigned routing hints and MUST NOT
be used for verification.

Signature order has no meaning. Implementations MUST search all signatures and
MUST NOT assign roles by array position.

### 5.1 Party identities

A party envelope contains an `org.Party` with a `gobl:` endpoint. That endpoint
defines the subject. `Client.VerifyParty` requires a valid signature whose
`iss` equals the subject and returns the canonical subject address. It does not
check authority endorsements.

A public `/who` response has the following additional requirements:

- the subject MUST equal the address queried; and
- the envelope MUST contain a valid audience-free self-signature.

Other self-signatures and countersignatures may be present.

### 5.2 Deliveries

For an ordinary inbox delivery, `Client.VerifyDelivery(ctx, env, receiver)`
searches for valid signatures with `aud` equal to `receiver`. All matching
signatures MUST have the same issuer. The method returns that issuer, or fails
if no issuer or multiple issuers bind the delivery.

Registration, verification, and deferred identity workflows may carry party
envelopes without an audience binding. Their receiver-specific rules must be
applied separately.

## 6. Request authentication

Requests to `/who` and `/inbox` MUST include a compact ES256 request token as
`Authorization: Bearer <token>`. Key discovery is unauthenticated.

The token contains:

| Claim | Requirement |
|---|---|
| `iss` | Requester's canonical address. |
| `aud` | Destination address. |
| `iat` | Required issue time. |
| `exp` | Required expiry. |
| `jti` | Optional identifier; generated by `NewToken`. |

The requester may be an intermediary and need not be the envelope signer.
Authentication of the HTTP request and verification of the document are
independent.

`NewToken` uses a one-minute lifetime by default and clamps explicit lifetimes
to 30 seconds through five minutes. `VerifyToken` verifies the issuer's
published key, audience, signature, key validity window, and freshness. It
allows 30 seconds of clock skew and rejects tokens issued more than five
minutes ago regardless of `exp`. `VerifyAuthorization` additionally parses the
Bearer header.

If the issuer's key endpoint is temporarily unavailable, verification returns
`ErrUnavailable`; a server SHOULD answer `503`, not `401`.

## 7. Authorities and sender policy

A registration authority countersigns a party envelope with:

- `iss` equal to the authority address;
- `aud` equal to the party address; and
- optionally, `verifier` naming the address that performed KYC/KYB.

The verifier MUST be a different participant from the registration authority
and MUST countersign the same envelope. Separating these roles prevents the
authority from attesting to identity checks that no independent verifier
performed. A self-referential, missing, invalid, or expired verifier claim
reduces the result to registered; it does not invalidate a valid registration
endorsement.

`Client.VerifyAuthority` checks all candidate authority signatures and prefers
a confirmed verifier over a registration-only result. An authority signature's
`exp`, when present, is enforced. The method returns an `Endorsement`:

- `Authority` identifies the trusted registration authority;
- `Verifier` is set only when verifier evidence is confirmed; and
- `Endorsement.Verified()` reports whether `Verifier` is set.

The default live authority is `lookup.gobl.org`. `WithAuthorities` replaces the
client's trust list. `WithSandbox` replaces it with
`lookup.sandbox.gobl.org`; live and sandbox trust lists are separate.

`Client.VerifySender(ctx, addr, requireVerified)` combines `/who` lookup and
authority verification. It always requires registration and, when
`requireVerified` is true, confirmed identity verification.

Authority endorsement is a sender policy, not a prerequisite for receiving.
A receive-only address may publish only an inbox and keys, and may return `204`
from `/who`.

## 8. Endpoints

### 8.1 `GET /.well-known/gobl/keys/<kid>`

This endpoint is public and returns one JWK as `application/json`.

| Status | Meaning |
|---|---|
| `200` | Key returned. |
| `404` | Key is unknown or removed. |
| `429`, `5xx` | Temporary failure. |

### 8.2 `GET /.well-known/gobl/who`

This endpoint requires a valid request token. It returns a signed party
envelope or one of these empty responses:

| Status | Meaning |
|---|---|
| `200` | Party envelope returned. |
| `202` | Request recorded for possible deferred disclosure. |
| `204` | The participant does not publish identity details. |
| `401` | Request token missing or invalid. |
| `404` | Address does not participate. |
| `503` | A required key endpoint is unavailable. |

`Client.Who` maps `202` to `ErrPending` and `204` to `ErrNoContent`, then
verifies a `200` response as described in section 5.1. Authenticated responses
SHOULD use `Cache-Control: private`. Clients may cache verified identities for
a bounded period.

### 8.3 `POST /.well-known/gobl/inbox`

This endpoint requires both a valid request token and a signed GOBL envelope.
For ordinary documents, the envelope MUST have exactly one issuer bound to the
inbox address as described in section 5.2. The receiver may then apply its
sender endorsement policy.

| Status | Meaning |
|---|---|
| `202` | Envelope accepted and durably recorded. |
| `400` | Body is unreadable or invalid JSON. |
| `401` | Authentication or delivery signature failed. |
| `403` | Sender policy rejected the envelope. |
| `422` | Envelope failed structural validation. |
| `500` | Persistence failed. |
| `503` | A required key endpoint is unavailable. |

Acceptance MUST be idempotent on the envelope UUID and digest. Clients may
retry `ErrUnavailable`; they SHOULD NOT retry `ErrInboxRejected` without
changing the request.

`Client.Send` serializes the envelope, adds a fresh request token when the
client has an identity, and succeeds only on `202`.

## 9. Client and transport

`NewClient` uses `HTTPFetcher`, the default live authority list, and a
five-minute public-key cache. Configure it with:

| Option | Effect |
|---|---|
| `WithIdentity` | Sets the address and key used for request tokens. |
| `WithAuthorities` | Replaces trusted registration authorities. |
| `WithSandbox` | Replaces live authorities with sandbox authorities. |
| `WithKeyCacheTTL` | Changes the key cache lifetime; zero disables it. |
| `WithFetcher` | Replaces the HTTP transport. |

The default `HTTPFetcher` has a 10-second request timeout, a five-second dial
timeout, and a 1 MiB response limit. It classifies `429`, `5xx`, and transport
failures as `ErrUnavailable`. Other unsuccessful fetch responses are
`ErrFetchFailed`; definitive inbox `4xx` responses are `ErrInboxRejected`.

To reduce SSRF risk, the default dialer rejects hosts resolving to loopback,
private, link-local, multicast, or unspecified addresses. Custom `Fetcher`
implementations SHOULD provide equivalent protection. Tests and local services
can inject a custom fetcher.

Fetched keys are cached by URL, only on success, with a maximum of 1,024
entries. `FlushKeyCache` forces subsequent verification to fetch keys again.

## 10. Typical delivery flow

1. The sender obtains the receiver's address out of band.
2. The sender may call `/who` for the receiver's invoicing details.
3. The sender signs the envelope with `iss=sender` and `aud=receiver`.
4. The sender posts it to the receiver's inbox with a request token.
5. The receiver verifies the token and calls `VerifyDelivery`.
6. The receiver may call `VerifySender` on the returned issuer.
7. A `202` confirms durable receipt.

Only step 6 requires authority endorsement, and that requirement applies to
the sender. Any response document is a new delivery and is subject to the same
rules.

## 11. Errors

| Error | Meaning |
|---|---|
| `ErrAddressEmpty` | Address is empty. |
| `ErrAddressInvalid` | Address is not a valid FQDN. |
| `ErrFetchFailed` | Permanent fetch, content, or configuration failure. |
| `ErrUnavailable` | Retryable network, `429`, or `5xx` failure. |
| `ErrNoContent` | Endpoint returned `204`. |
| `ErrPending` | Endpoint returned `202` where a body was expected. |
| `ErrTokenInvalid` | Request token is malformed or invalid. |
| `ErrTokenExpired` | Request token failed freshness checks. |
| `ErrVerifyFailed` | Envelope or signature verification failed. |
| `ErrUnknownAuthority` | No trusted authority endorsement was found. |
| `ErrNotVerified` | Registration is valid but verification was required. |
| `ErrSignatureExpired` | Authority endorsement has expired. |
| `ErrPartyMissing` | Identity envelope has no `org.Party`. |
| `ErrInboxRejected` | Inbox returned a definitive `4xx`. |

Callers should classify wrapped errors with `errors.Is`.

## 12. Security considerations

- **Web PKI:** control of an address and its TLS certificate controls its
  published keys. Applications must compare verified issuers with the
  identities they expect or apply authority policy.
- **Key removal:** the client cache may continue to accept a removed key until
  its TTL expires. Use short TTLs and `FlushKeyCache` during incident response.
- **Endorsement freshness:** authority and verifier `exp` claims bound their
  lifetimes. Identity caches must also have bounded lifetimes.
- **Resource limits:** clients limit response bodies to 1 MiB and verification
  to 32 signatures per envelope. Servers should impose comparable request
  limits.
- **Privacy:** request tokens identify callers but do not make party data
  public. Operators should limit audit-log retention because requester metadata
  can reveal business relationships.
- **Routing fields:** unsigned `from` and `to` values are hints only. Security
  decisions must use signed claims.

## 13. References

- RFC 2119 and RFC 8174 — requirement keywords
- RFC 6750 — Bearer token usage
- RFC 7515 — JSON Web Signature
- RFC 7517 — JSON Web Key
- RFC 7518 — JSON Web Algorithms
- RFC 7519 — JSON Web Token
- `github.com/invopop/gobl/dsig` — signature and key primitives
- `github.com/invopop/gobl/org` — party schema
