# GOBL Net

> **Experimental:** the API and wire protocol may change without notice.

**Status:** Draft

## Abstract

GOBL Net is a decentralized discovery and delivery protocol for signed GOBL
documents. A participant is identified by a fully qualified domain name
(FQDN). HTTPS binds that address to these well-known resources:

| Resource | Purpose |
|---|---|
| `/.well-known/gobl/keys/<kid>` | Publish one verification key. |
| `/.well-known/gobl/who` | Publish a signed `org.Party` identity. |
| `/.well-known/gobl/inbox` | Accept signed GOBL envelopes. |

A signature names its issuer in the signed `iss` claim. The participant
checking it obtains the key from the issuer's domain and relies on Web PKI to
authenticate the HTTPS response. Registration authorities add an optional
trust layer by countersigning party envelopes.

**This document defines the protocol.** It is implementation-neutral and names
no API of any language: every requirement here is on the wire or on the
behaviour of a participant. The Go package alongside it is a reference
implementation from which libraries in other languages can be made; its API is
documented in the package itself (`go doc github.com/invopop/gobl/net`), which
maps each section below onto the types and methods implementing it. The
reference server and CLI live in
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
| **Verifier** | A participant, independent of the Authority, that establishes who is behind an address and countersigns the subject's party envelope. What it checks, and how, is its own to define and publish; see §7.1. |
| **Registration countersignature** | A signature issued by an Authority over a party envelope, binding the Subject. It may name a Verifier. |
| **Verification countersignature** | A signature issued by a Verifier over the same party envelope. |
| **Sender status** | What a receiver concludes about a Subject after evaluating the countersignatures its party envelope carries, by the procedure in §7.2. One of the three below. |
| **Unregistered** | The receiver found no valid registration countersignature from an Authority it trusts. It has only the Subject's own account of itself. This is a statement about what *this* receiver can corroborate, not about whether the Subject registered with somebody. |
| **Registered** | An Authority the receiver trusts has countersigned the party envelope, attesting that the Subject controls its address. Who is behind that address is not established. |
| **Verified** | Registered, and that Authority named a Verifier whose own countersignature the receiver found valid on the same envelope — so a participant independent of the Authority has established who is behind the address. |
| **Party envelope** | A signed GOBL envelope containing an `org.Party`. It identifies its Subject and may carry registration and verification countersignatures. |

A countersignature is a signature on the wire; a sender status is a conclusion
a receiver reaches about one. The distinction matters because **no status is
ever transmitted**. Unregistered, registered and verified name conclusions,
not values: no field carries them, and nothing in an envelope asserts one.
Each receiver derives its own from the signatures in front of it and the
authorities it happens to trust, so two receivers can read the same envelope
and reach different conclusions about it, and both are correct.

Issuer, Sender, and Requester often refer to the same participant in a direct
delivery, but the protocol does not assume they are the same. Audience and
Receiver are equivalent for ordinary document delivery; in other signatures,
such as a registration countersignature, the Audience is the Subject rather
than the recipient of an HTTP request.

An Authority and Verifier MUST be different participants. Other roles may
overlap where their definitions permit it.

## 3. Addresses and discovery

A GOBL Net address is an FQDN without a scheme, port, path, query, or fragment.
An address is canonicalized by:

1. trimming surrounding whitespace and one trailing dot;
2. converting internationalized names to lowercase ASCII A-labels using the
   IDNA2008 lookup profile;
3. validating DNS label syntax; and
4. requiring at least two labels.

An empty address and one that fails any of these steps are both invalid, and
an implementation SHOULD distinguish them so a caller can tell a missing value
from a malformed one. Protocol fields and comparisons MUST use the canonical
ASCII form, so that two participants comparing the same address by different
routes always agree.

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

A published key may include `valid_from` and `valid_until`. When a signature
has an `iat`, it MUST be rejected if the signing time falls outside those
bounds. A missing bound or missing `iat` skips the corresponding
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
| `exp` | Optional expiry, used primarily for countersignatures. |
| `verifier` | Optional verifier named by a registration authority. |

Envelope header `from` and `to` fields are unsigned routing hints and MUST NOT
be used for verification.

Signature order has no meaning. Implementations MUST search all signatures and
MUST NOT assign roles by array position.

### 5.1 Party identities

A party envelope contains an `org.Party` with a `gobl:` endpoint. That endpoint
defines the subject, and it is read from the party document rather than from
any signature: any holder can reorder or add signatures, so the subject MUST
NOT be inferred from whichever one appears first.

To accept a party envelope, a participant MUST find a valid signature whose
`iss` equals the subject. This establishes only that the subject signed its
own party document; it says nothing about registration, which is §7.

A public `/who` response has the following additional requirements:

- the subject MUST equal the address queried; and
- the envelope MUST contain a valid audience-free self-signature.

Other self-signatures and countersignatures may be present.

### 5.2 Deliveries

For an ordinary inbox delivery, the receiver MUST consider every valid
signature whose `aud` equals its own address. All such signatures MUST share
one issuer, and that issuer is the sender. A delivery that no signature binds,
or that two different issuers bind, MUST be rejected: the first has nobody
accountable for it, and the second leaves the receiver to choose, which is not
a choice a receiver may make.

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
| `jti` | Optional unique identifier, to let a receiver detect replay. |

The requester may be an intermediary and need not be the envelope signer.
Authentication of the HTTP request and verification of the document are
independent.

Verifying a token means checking, against the issuer's published key: the
signature, the audience, the key's validity window, and freshness.

The freshness parameters are **normative**, because they decide whether two
independent implementations interoperate at all. A participant that chooses
its own values will have tokens rejected by conforming peers for no visible
reason.

| Parameter | Value |
|---|---|
| Default token lifetime | 60 seconds |
| Minimum accepted lifetime | 30 seconds |
| Maximum accepted lifetime | 5 minutes |
| Permitted clock skew | 30 seconds |
| Maximum age, regardless of `exp` | 5 minutes |

An issuer MUST NOT mint a token whose lifetime falls outside the minimum and
maximum above. A verifier MUST apply the permitted skew to both `iat` and
`exp`, and MUST reject a token issued more than the maximum age ago even when
its `exp` has not passed — an over-long `exp` does not extend a token's life.

If the issuer's key endpoint is temporarily unavailable the token cannot be
checked, which is not the same as being invalid: the server SHOULD answer
`503`, not `401`, so the caller retries rather than re-authenticates.

## 7. Authorities and sender policy

A registration authority countersigns a party envelope with:

- `iss` equal to the authority address;
- `aud` equal to the party address; and
- optionally, `verifier` naming the address that performed KYC/KYB.

The verifier MUST be a different participant from the registration authority
and MUST countersign the same envelope. Separating these roles prevents the
authority from attesting to identity checks that no independent verifier
performed.

Each receiver holds a **trust list** of the authorities it accepts. The
default live authority is `lookup.gobl.org` and the default sandbox authority
is `lookup.sandbox.gobl.org`. The two lists are disjoint and MUST NOT be
combined: a sandbox registration carries the same shape as a live one, so a
receiver that trusted both could not tell a test identity from a real one.

Sender status is a **policy for sending**, not a prerequisite for receiving. A
receive-only address may publish only an inbox and keys, and may return `204`
from `/who`.

### 7.1 The verification exchange

A verifier establishes who is behind an address and countersigns the subject's
party envelope to say so. The exchange is **subject-mediated**: the subject
holds its own envelope throughout and presents it to each participant in turn.

1. The subject prepares its party envelope — an `org.Party` declaring its
   `gobl:` endpoint, self-signed, with the signing key published at that
   address (§4).
2. It asks a verifier of its choosing to verify it. The verifier performs its
   checks, by whatever means and over whatever period it needs.
3. On success the verifier countersigns **that exact envelope**, unchanged in
   `uuid` and digest, with `iss` equal to its own address, `aud` equal to the
   subject, and an `exp`. It does not set a `verifier` claim: that pointer is
   the authority's to write.
4. The verifier returns the countersigned envelope to the subject. Delivering
   it to the subject's inbox (§8.3) is RECOMMENDED.
5. The subject registers that envelope with the authority, which recognises a
   countersignature from an address on its trusted list, countersigns naming
   it in the `verifier` claim, and delivers the result to the subject. The
   subject publishes it at `/who`.

Verification comes before registration because the authority's countersignature
is the last one added. With the verifier's already aboard, a single
registration produces a fully endorsed identity, and the subject's party
document is published once rather than twice.

The reverse order remains valid and needs no separate description: a party
that registered without a verifier, and later wants one, follows exactly these
steps with the envelope it already holds. The authority supersedes its own
earlier countersignature when the result comes back. It simply costs a second
registration to arrive where one would have done.

The same steps also describe keeping an identity current. An authority
countersignature carries an `exp`, and a subject renews by submitting the
endorsed envelope again before it passes — the verification countersignature
travelling with it, since an authority that cannot see one will not name a
verifier. Registering, adding verification later, and renewing are therefore
one operation applied at different times, not three.

How the envelope gets back to the subject in step 4 is not constrained. What
the exchange requires is that the subject ends up holding that exact envelope;
the means are between a verifier and its own customers. A verifier that
returns it by some other route — its own API, a link in an email, a console
the subject signs into — conforms, as does one that reports progress or
history however it likes. None of that belongs in this specification.

Inbox delivery is recommended because it costs neither side anything new.
A party envelope moving between two participants is what §8.3 is for, the
subject already operates an inbox by virtue of being on the network, and it
is the same route the authority itself uses in step 5 — so a verifier that
takes it publishes no endpoint of its own, and a subject that implements
only the inbox can work with any verifier that follows the recommendation.
A verifier choosing another route is choosing to require something of its
own customers, which is its business rather than the network's.

Where the inbox is used, delivery SHOULD be retried on the temporary failures
of §11. However the envelope is returned, a countersignature is a statement
that has already been made: a verifier MUST NOT treat an unreturnable result
as a failed verification, and MUST NOT re-sign on a retry — that would extend
the countersignature's life each time it was re-sent. A subject that never
received one asks for it again, which is step 2.

A verifier needs no relationship with any authority. It verifies the subject's
self-signature against the key published at the subject's own address (§5.1),
which is evidence it can gather alone, and it countersigns on that basis. It
does not need to know which authority the subject will register with, whether
the subject has registered at all, or whether that authority would accept its
countersignature — the last of which is the authority's decision and is made
in §7.2, after the verifier's work is done. A verifier MAY still require the
subject to be registered before it will act, but that is its own admission
policy and not a rule of the protocol.

For an already-registered subject only, a verifier MAY deliver the
countersigned envelope to the authority's inbox instead of returning it, which
re-registers the envelope on the subject's behalf. This is the older
arrangement and remains conforming, which matters for more than history: it is
what lets a network replace a verifier without a flag day, running the two
side by side while countersignatures issued under the previous one are still
within their `exp`. It does not apply before a first registration, where the
subject has not yet chosen an authority to be registered with.

The subject cannot subvert any of this by carrying the envelope. The
countersignature covers the digest, so altering anything invalidates it, and
the authority verifies it against the verifier's published key. Withholding is
the subject's right — it is their identity.

Because the authority looks only at the signatures an envelope carries, and
not at who delivered it, this exchange requires nothing of an authority beyond
what §7 already specifies — and it is the same property that makes the two
delivery paths interchangeable.

### 7.2 Determining sender status

A receiver reaches a sender status by evaluating the countersignatures a party
envelope carries against its own trust list. The procedure is normative: two
conforming receivers with the same trust list MUST reach the same status for
the same envelope.

1. Consider every signature whose `iss` is an authority on the trust list. An
   envelope carrying none is **unregistered**, whatever else it carries — a
   countersignature from an unknown authority is not evidence.
2. For each candidate, verify it against that authority's published key and
   reject it if an `exp` is present and has passed.
3. A candidate that verifies makes the subject at least **registered**.
4. If a verifying candidate names a `verifier`, and that verifier is a
   different participant from the authority, and the same envelope carries a
   valid, unexpired countersignature issued by that verifier, the subject is
   **verified**.
5. Where candidates disagree, the highest status wins: an envelope may carry
   several authority signatures, and one that establishes **verified** is not
   cancelled by another that only establishes **registered**.

A `verifier` claim that is absent, malformed, self-referential, or whose
countersignature is missing, invalid or expired reduces the status to
**registered**. It MUST NOT invalidate an otherwise valid registration — a
claim nobody can corroborate is worth nothing, but it does not make the
authority's own statement worthless.

The one case that is not a downgrade is a verifier key endpoint that cannot be
reached. "Could not check" is not "not verified": a receiver MUST NOT report
**registered** on a temporary failure, and SHOULD treat it as §11 does, so a
verified sender is never quietly demoted by somebody else's outage.

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

A `200` response MUST be verified as described in §5.1 before anything in it
is used. Authenticated responses SHOULD use `Cache-Control: private`, and a
caller may cache a verified identity for a bounded period — bounded because a
cached identity outliving the countersignature that endorsed it is §12's
problem.

### 8.3 `POST /.well-known/gobl/inbox`

This endpoint requires both a valid request token and a signed GOBL envelope.
For ordinary documents, the envelope MUST have exactly one issuer bound to the
inbox address as described in §5.2. The receiver may then apply its sender
status policy (§7.2).

| Status | Meaning |
|---|---|
| `202` | Envelope accepted and durably recorded. |
| `400` | Body is unreadable or invalid JSON. |
| `401` | Authentication or delivery signature failed. |
| `403` | Sender policy rejected the envelope. |
| `422` | Envelope failed structural validation. |
| `500` | Persistence failed. |
| `503` | A required key endpoint is unavailable. |

Acceptance MUST be idempotent on the envelope UUID and digest: the same
envelope delivered twice is acknowledged twice and stored once, which is what
makes the retry rule in §11 safe. Only `202` means the envelope was accepted;
any other `2xx` MUST NOT be read as success.

## 9. Limits

These bound what a participant will do on a peer's behalf. They are a defence
against a hostile or broken peer rather than a tuning parameter, so an
implementation that raises them weakens itself and not the protocol.

| Limit | Value |
|---|---|
| Response body accepted from a peer | 1 MiB |
| Signatures verified per envelope | 32 |

An envelope carrying more signatures than the limit MUST be rejected rather
than partially verified. Because each round trip through an authority or
verifier appends a signature, a party envelope approaching the limit will stop
being endorsable before it stops being valid; a participant that countersigns
SHOULD refuse an envelope close enough to the limit that its own signature
would exhaust it, and say so, rather than produce the signature that breaks
the document.

Fetching a key or an identity means issuing an HTTP request to an address a
peer chose, which is a server-side request forgery primitive. An implementation
MUST refuse to fetch from hosts resolving to loopback, private, link-local,
multicast, or unspecified addresses, except where an operator has deliberately
configured otherwise for testing or a local service.

Fetched keys MAY be cached, keyed by URL and only on success. A cache SHOULD
be bounded in both entries and lifetime; see §12 on what a stale entry costs.

## 10. Typical delivery flow

1. The sender obtains the receiver's address out of band.
2. The sender may call `/who` for the receiver's invoicing details.
3. The sender signs the envelope with `iss=sender` and `aud=receiver`.
4. The sender posts it to the receiver's inbox with a request token.
5. The receiver verifies the token and resolves the sender (§5.2).
6. The receiver may look up that sender's `/who` and apply its status policy
   (§7.2).
7. A `202` confirms durable receipt.

Only step 6 involves sender status, and that requirement falls on the
sender. Any response document is a new delivery and is subject to the same
rules.

## 11. Failure handling

Every failure in this protocol falls into one of three classes, and an
implementation MUST be able to tell them apart — reporting a temporary failure
as a definitive one turns somebody else's outage into a rejection the caller
cannot distinguish from a refusal.

| Class | Examples | Caller behaviour |
|---|---|---|
| **Temporary** | `429`, `5xx`, connection and timeout failures | Retry with backoff. |
| **Definitive** | a `4xx` other than `429`, a signature that fails verification | Do not retry without changing the request. |
| **Absent** | `204` from `/who`, `404` for an unknown key or address | Not a failure. The peer answered, and the answer is that there is nothing. |

Two consequences run through the rest of this document.

A check that could not be *performed* is never reported as one that *failed*.
An unreachable key endpoint means the signature is unknown, not invalid: a
server answers `503` rather than `401` (§6), and a receiver does not downgrade
a sender's status (§7.2).

Delivery retries are safe because inbox acceptance is idempotent on the
envelope UUID and digest (§8.3). A caller that retries a delivery it is unsure
about cannot create a duplicate; one that does not retry can lose a document
to a transient failure. Retrying is the correct behaviour.

## 12. Security considerations

- **Web PKI:** control of an address and its TLS certificate controls its
  published keys. Applications must compare verified issuers with the
  identities they expect, or apply a sender status policy (§7.2). An
  implementation MUST authenticate the HTTPS response when fetching a key; a
  key fetched over an unauthenticated channel proves nothing.
- **Key removal:** a cached key may continue to be accepted after it is
  withdrawn, until the cache entry expires. Key caches SHOULD be short-lived
  and an implementation SHOULD offer a way to flush them during incident
  response, since removing a key from publication is otherwise not an
  immediate revocation.
- **Countersignature freshness:** authority and verifier `exp` claims bound
  their lifetimes. Identity caches must be bounded too, or a cached identity
  outlives the countersignature that justified trusting it.
- **Resource limits:** the limits in §9 are a defence against a hostile peer,
  not a tuning parameter. Servers should impose comparable limits on the
  requests they accept.
- **Privacy:** request tokens identify callers but do not make party data
  public. Operators should limit audit-log retention because requester metadata
  can reveal business relationships.
- **Routing fields:** unsigned `from` and `to` values are hints only. Security
  decisions must use signed claims.
- **Verification is not graded.** A subject is verified or it is not, and any
  verifier the authority names satisfies it. A receiver that wants to
  distinguish between verifiers must read the `verifier` claim itself and
  decide which addresses it accepts. It follows that an authority admits a
  weak verifier at the cost of every sender it marks verified, and that a
  verifier's address is the only
  trust signal a receiver gets: it should name what is being vouched for, and
  a change in what the verifier checks is a change in what that name means.
- **Verifier scope:** a verifier attests to the subject's identity at a point
  in time. It says nothing about any document the subject later sends, and an
  `exp` bounds how long the attestation may be relied upon at all.

## 13. References

- RFC 2119 and RFC 8174 — requirement keywords
- RFC 6750 — Bearer token usage
- RFC 7515 — JSON Web Signature
- RFC 7517 — JSON Web Key
- RFC 7518 — JSON Web Algorithms
- RFC 7519 — JSON Web Token
- `github.com/invopop/gobl/dsig` — signature and key primitives
- `github.com/invopop/gobl/org` — party schema
