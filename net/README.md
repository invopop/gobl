# GOBL Net

> **Experimental:** the API and wire protocol may change without notice.

**Status:** Draft

## Abstract

GOBL Net lets organizations find each other and exchange signed GOBL
documents, such as invoices, without a central operator. Each participant is
identified by a domain name, its *address*, and serves three resources over
HTTPS at that domain:

| Resource | Purpose |
|---|---|
| `/.well-known/gobl/keys/<kid>` | Publish a public key, so others can check the participant's signatures. |
| `/.well-known/gobl/who` | Publish the participant's identity: a signed `org.Party`. |
| `/.well-known/gobl/inbox` | Accept signed GOBL envelopes. |

Documents travel in GOBL envelopes signed by their sender. A receiver checks
the signature with the key published at the sender's domain, and because
HTTPS authenticates the domain, a valid signature proves the envelope came
from whoever controls that address. It proves nothing more: anyone can
register a domain.

Two further roles tell a receiver who is behind an address:

- A **verifier** checks that a participant's details belong to a real
  organization (on a live network, through a KYC/KYB process) and that this
  organization controls the address. It records the result by adding its own
  signature to the participant's identity.
- A **registration authority** keeps a public index of verified
  participants. It decides which verifiers to rely on, keeps each
  registration current, and adds its own signature naming the verifier.

A receiver that relies on the authority or the verifier treats the sender as
**verified** (§7).

This document defines the protocol, independent of any programming language.
The Go package alongside it is a reference implementation; its documentation
(`go doc github.com/invopop/gobl/net`) maps each section onto the code. The
reference server and CLI live in
[`github.com/invopop/gobl.dev`](https://github.com/invopop/gobl.dev/#gobl-net).

## 1. Conventions

The terms MUST, MUST NOT, SHOULD, SHOULD NOT, and MAY are interpreted as
described in BCP 14 (RFC 2119 and RFC 8174).

## 2. Terminology

| Term | Meaning |
|---|---|
| **Address** | A domain name identifying a participant, such as `supplier.example.com` (§3). |
| **Participant** | Anyone with an address. One participant may play different roles in different exchanges. |
| **Envelope** | A GOBL envelope: a document, a header, and signatures. The header's `from` and `to` fields hold the sending and receiving addresses. |
| **Signature** | A signature over an envelope's document. It also signs a few fields (§5) saying who made it, who it is for, and when. |
| **Issuer** | The participant that made a signature. |
| **Audience** | The participant a signature is for. A signature counts only for its audience. |
| **Sender**, **Receiver** | An envelope's `from` and `to`. They count only when a signature by the sender, with the receiver as its audience, confirms them (§5). |
| **Party envelope** | An envelope whose document is an `org.Party`: a participant's name, tax identities, and other details. |
| **Subject** | The participant a party envelope describes: its `from` address. |
| **Requester** | A participant asking another for its party envelope at `/who` (§6). |
| **Countersignature** | A signature added to a subject's party envelope by another participant, saying something about the subject. Its audience is always the subject. There are two kinds, below. |
| **Verifier** | A participant that confirms who a subject is: that the party details describe a real organization or person (on a live network, by a KYC/KYB process) that controls the subject's address. It says so with a **verification countersignature**. |
| **Authority** | A registration authority: a participant that registers subjects a verifier has confirmed, and publishes them as a public index. It decides which verifiers it will name, follows the subject's key as it changes (§7.3), and re-confirms control of the address at least every 90 days. It says so with a **registration countersignature** naming the verifier. It does not confirm identity itself, and MUST be a different participant from the verifier it names. |
| **Trust list** | The authorities a receiver relies on, and optionally verifiers it relies on directly (§7). |
| **Verified** | The receiver can link the subject to a verifier it relies on, through an authority on its trust list or directly (§7.2). |
| **Unverified** | Anything else. The receiver knows the subject controls its address and key, but the party details are only the subject's own claim. |

Verified and unverified are conclusions a receiver reaches, not values carried
in an envelope. Receivers with different trust lists can reach different
conclusions about the same envelope, and both are correct.

## 3. Addresses

An address is a fully qualified domain name, without a scheme, port, path,
query, or fragment. It is canonicalized by:

1. trimming surrounding whitespace and one trailing dot;
2. converting internationalized names to lowercase ASCII A-labels using the
   IDNA2008 lookup profile;
3. validating DNS label syntax; and
4. requiring at least two labels.

An empty address and one that fails any of these steps are both invalid, and
an implementation SHOULD distinguish them so a caller can tell a missing value
from a malformed one. Protocol fields and comparisons MUST use the canonical
form.

For the address `supplier.example.com`:

```text
supplier.example.com
gobl:supplier.example.com
https://supplier.example.com/.well-known/gobl/keys/<kid>
https://supplier.example.com/.well-known/gobl/who
https://supplier.example.com/.well-known/gobl/inbox
https://supplier.example.com/.well-known/jwks.json
```

The bare form is used inside signatures (§5). The `gobl:` form is used in an
envelope's `from` and `to` fields and in optional `org.Party` endpoints.
Production endpoints MUST use HTTPS.

## 4. Keys

A participant publishes the public keys it signs with at its address. Each key
MUST be available as a single RFC 7517 JSON Web Key at
`/.well-known/gobl/keys/<kid>`, where `kid` is the key's ID; unknown or
removed IDs return `404`. Signatures name their key by ID and do not carry a
`jku` header. The optional `/.well-known/jwks.json` lists all keys for generic
tooling, but is not used by this protocol.

A published key may include `valid_from` and `valid_until`. A signature MUST
be rejected if its signing time (§5) falls outside them; a missing bound is
not checked. Because the signer chooses the signing time, a retired key could
still sign with a backdated one, so a delivery or `/who` response MUST also be
signed with a key whose bounds include the time it is received. Retired keys
SHOULD stay published while countersignatures made with them have not
expired, and for checking historical signatures.

## 5. Envelope signatures

A signature is a JSON Web Signature (RFC 7515) covering the envelope's
document UUID and digest, plus these fields, called claims:

| Claim | Name used in this document | Meaning |
|---|---|---|
| `iss` | issuer | The signer's address. |
| `aud` | audience | The address the signature is for. Required on the signature confirming `from` and `to`, and on countersignatures. |
| `iat` | signing time | An RFC 7519 NumericDate. Required, and MUST NOT be later than the present, allowing the clock skew of §6. |
| `exp` | expiry | Required on registration countersignatures, at most 90 days after the signing time; optional otherwise. |
| `cnf` | bound key | The key a signature applies to, as RFC 7800 `{"jkt": "<thumbprint>"}`: the key's RFC 7638 SHA-256 thumbprint, base64url-encoded. Required on countersignatures, where it is the subject's key, and on the current key's signature in a key change (§7.3), where it is the new key. |
| `verifier` | verifier | The verifier's address. Required on registration countersignatures. |

A signature covers the document only through its digest, so a participant
MUST reject an envelope with no document, or whose document does not match
its digest, before relying on any of its signatures.

Every GOBL Net envelope MUST carry `from` and `to` header fields holding
`gobl:` addresses, and a valid signature whose issuer is `from` and whose
audience is `to`: the signature that **confirms** the envelope. The fields
themselves are not signed, so that signature is what makes them evidence. A
participant MUST reject an envelope missing either field, using another
scheme in either, or with no confirming signature.

A participant passing an envelope on sets `from` to its own address and `to`
to the next recipient, and adds its own confirming signature: nothing is
forwarded on another's behalf. Once an envelope is accepted, its confirming
signature has served its purpose. A holder MAY remove it, and SHOULD before
passing the envelope on, so signatures do not build up towards the limit in
§9. Countersignatures (§7) are the only signatures that matter beyond the
exchange that carries them.

Signature order has no meaning. Implementations MUST search all signatures
and MUST NOT assign roles by position.

### 5.1 Party identities

The subject's address is not part of the party document: it is the envelope's
`from`, confirmed as above. A party can therefore be verified before it has an
address, and one party can serve several addresses, each with its own copy of
the envelope. The party MAY list `gobl:` endpoints, but they are only hints
for display and MUST NOT be used to determine the subject. Nor does the first
signature determine it: anyone holding an envelope can add or reorder
signatures, so the subject MUST NOT be taken from signature order.

The confirming signature shows only that the subject adopted the party details
as its own. Whether anybody vouches for them is §7.

A `/who` response is made for each requester: its `from` MUST be the address
queried and its `to` the requester (the issuer of the request token, §6),
confirmed by the subject's signature. The requester MUST reject a response
from or to anybody else. A server does not keep these per-requester
signatures in the envelope it stores.

### 5.2 Deliveries

A receiver MUST reject a delivered envelope whose `to` is not its own address,
or whose confirming signature was not made with a key valid at the time of
receipt (§4). The sender is the `from` address; signatures by anybody else say
nothing about who sent it.

Envelopes exchanged outside GOBL Net, such as a party envelope handed to a
verifier through its own API, need no audience and are not subject to these
rules.

## 6. Request authentication

Requests to `/who` MUST carry a request token, so the server knows who to
make the response for (§5.1). The token is a compact JSON Web Token
(RFC 7519) signed with ES256, sent as `Authorization: Bearer <token>`
(RFC 6750). Deliveries need no token, because the envelope's own signature
identifies the sender and delivering an envelope twice changes nothing
(§8.3). Key requests need no authentication.

| Claim | Requirement |
|---|---|
| `iss` | The requester's address. |
| `aud` | The address queried. |
| `iat` | Required issue time. |
| `exp` | Required expiry. |
| `jti` | Optional unique identifier, to let a server detect replay. |

A server checks the token against the requester's published key: the
signature, the audience, the key's validity bounds, and freshness. The
freshness parameters are normative, because peers using other values would
reject each other's tokens:

| Parameter | Value |
|---|---|
| Default token lifetime | 60 seconds |
| Minimum accepted lifetime | 30 seconds |
| Maximum accepted lifetime | 5 minutes |
| Permitted clock skew | 30 seconds |
| Maximum age, regardless of `exp` | 5 minutes |

A requester MUST NOT mint a token whose lifetime is outside the minimum and
maximum. A server MUST apply the clock skew to both `iat` and `exp`, and MUST
reject a token issued more than the maximum age ago even if its `exp` has not
passed.

If the requester's key cannot be fetched because of a temporary failure, the
token cannot be checked, which is not the same as being invalid: the server
SHOULD answer `503`, not `401`, so the caller retries.

## 7. Verification and registration

Signatures prove control of an address. Verifiers and authorities establish
who is behind it, by countersigning the subject's party envelope.

**Verifiers.** A verifier confirms that the party details describe a real
organization, and that this organization controls the subject's address. On a
live network that means a KYC/KYB process: the organization exists in a
company registry, the identifiers and legal name in the party match the
registry, and whoever presented them may act for the organization. This is
the only part of the protocol whose evidence comes from outside the network.
A verifier SHOULD publish what it checks, because its address is all a
receiver learns about how thorough the check was (§12).

A verification countersignature MUST have the subject's address as its
audience, and MUST bind the subject's key that the verifier was presented
with (§7.1). The audience stops a verification being reused with a copy of
the envelope at another address; the bound key stops it carrying over to a
new holder of the same address. Its expiry is optional.

**Authorities.** An authority registers only verified subjects, and publishes
them as a public index: the default place for a receiver to look, and the one
that decides which verifiers are worth relying on, so that receivers need
not. A registration countersignature MUST have:

- the subject's address as its audience, once the authority has confirmed the
  subject controls it (§7.1, step 5);
- a signing time, and an expiry no more than 90 days later;
- a bound key: the subject's key that the authority has registered (§7.3);
  and
- a `verifier` naming a different participant from the authority, whose
  valid verification countersignature is on the same envelope.

The bound key means a later holder of the address, after a domain lapses and
is bought by someone else, say, cannot present the registration as its own:
it controls the address's published keys, but not the registered one. The
short expiry limits how long a registration outlives the subject's control of
its address, and how long a leaked, retired authority key could forge one
with a backdated signing time. Like TLS certificates issued through ACME
(RFC 8555, the protocol behind Let's Encrypt), registrations are meant to be
renewed automatically: a subject SHOULD renew once two thirds of the lifetime
has passed.

A subject using one party at several addresses keeps a copy of the envelope
at each, each with its own countersignatures.

**Sandbox networks** exist to exercise the exchange, not to establish
anything about anybody, so a sandbox verifier may check as little as control
of the email address in the party. A live authority MUST NOT name such a
verifier. The authority is the safeguard here: a receiver relying on it takes
its word on which verifiers are worth naming, since a countersignature does
not say how thorough a check was (§12).

**Trust lists.** Each receiver keeps a trust list of the authorities it
relies on. The default live authority is `lookup.gobl.org`, and the default
sandbox authority is `lookup.sandbox.gobl.org`. Live and sandbox lists MUST
NOT be combined: their registrations look the same, so a receiver trusting
both could not tell a test identity from a real one.

A receiver MAY also rely on verifiers directly, without an authority, where it
has its own agreement with one, say. That list is empty by default. An
authority remains the first choice: it follows the subject's keys (§7.3), and
it is where other receivers look.

Sender status is a policy for sending, not a requirement for receiving. A
receive-only address may publish only an inbox and keys, and may answer
`/who` with `204`.

### 7.1 The verification exchange

The subject holds its own party envelope throughout and presents it to the
verifier, then to the authority. Nobody delivers it on the subject's behalf.

1. The subject prepares its party envelope.
2. It asks a verifier of its choosing to verify it for an address. The
   verifier's checks MUST establish that the party describes the organization
   presenting it, and that this organization controls the address. A delivery
   from the address (§5.2) proves control but not who holds it, so the
   verifier MUST tie the two together, for example by having the delivery
   made within a verification session it has already authenticated. The key
   that confirmed that delivery is the one the verifier binds. The checks may
   start before the subject has an address; the countersignature cannot.
3. On success the verifier countersigns that exact envelope, unchanged in
   UUID and digest, as issuer, with the subject's address as audience, the
   key from step 2 as bound key, and a signing time. It MAY add an expiry. It
   does not set a `verifier`: that is the authority's to write.
4. The verifier returns the envelope to the subject. Delivering it to the
   subject's inbox, from the verifier's address to the subject's, is
   RECOMMENDED, and SHOULD be retried on temporary failures (§11). Any other
   route, such as the verifier's own API or a link in an email, also conforms.
   A countersignature is a statement already made: a verifier MUST NOT treat
   an undeliverable result as a failed verification, and MUST NOT sign again
   on a retry, which would extend the countersignature's life. A subject that
   never received one asks again, at step 2.
5. The subject delivers the envelope to the authority's inbox, from its own
   address to the authority's, which shows it controls the address. The
   authority MUST reject the request unless it carries a valid verification
   countersignature from a verifier the authority accepts, binding a key in
   the subject's lineage (§7.3). It then countersigns, binding the key that
   confirmed the delivery and naming the verifier, and delivers the result
   to the subject's inbox. The subject publishes it at `/who`.

Steps 3 and 5 may be far apart. Another address needs its own verification,
though a verifier that has already checked the organization need only confirm
it controls the new address. A subject whose receivers rely on its verifier
directly may stop after step 4 and publish the envelope at `/who`.

**Renewal** repeats step 5 before the registration expires, with the
verification countersignature still on the envelope. An authority MAY limit
how old a verification without an expiry it will accept, using its signing
time, and SHOULD NOT accept one made with a key that is no longer valid.

A verifier needs no relationship with any authority: whether and where the
subject registers is the subject's choice, and whether to accept the verifier
is the authority's. The subject cannot alter a countersigned envelope, since
countersignatures cover its digest; it can only withhold it, which is its
right.

### 7.2 Determining sender status

A receiver decides a subject's status from the countersignatures on its party
envelope and its own trust lists. The procedure is normative: receivers with
the same trust lists MUST reach the same status for the same envelope.

A countersignature is **valid** when it verifies against its issuer's
published key, its audience is the subject, it has a bound key, its signing
time is no later than the present (allowing the clock skew of §6), and any
expiry has not passed.

A countersignature with a `verifier` is a registration. It never counts as a
verification, even when its issuer is also trusted as a verifier.

1. **Through an authority.** Consider each countersignature issued by an
   authority on the trust list. It counts if it is valid, its expiry is no
   more than 90 days after its signing time, its bound key confirmed the
   envelope (§5), and its `verifier` names a different participant whose
   valid countersignature is on the same envelope. Which key the verifier's
   countersignature binds is the authority's concern (§7.3). The first that
   counts makes the subject **verified**. One whose verifier is missing,
   malformed, the authority itself, or not valid does not count.
2. **Directly.** Failing that, a valid countersignature by a verifier on the
   receiver's own trust list, whose bound key confirmed the envelope, makes
   the subject **verified**. With no authority following the subject's key
   changes, a bound key counts for as long as the subject publishes it as
   valid: a subject changing key retires the old one once a verification
   binds the new one.
3. Otherwise the subject is **unverified**. Countersignatures from unknown
   authorities or verifiers, or for other addresses, are not evidence.

A key that cannot be fetched is not a rejection. A receiver MUST NOT report
unverified because of a temporary failure while a countersignature remains
unchecked, and SHOULD handle it as §11 describes.

A status applies to the holder of the bound key. A receiver applying it to a
delivery MUST check that the delivery was confirmed with that same key.

### 7.3 Key lineage

A subject's key changes over time, and its verification should not have to.
An authority follows the subject's key so that it can, while a new holder of
the address cannot take over the old holder's identity.

**What a subject does.** Each registration request is a delivery to the
authority (§7.1, step 5), signed just before it is sent, in one of three
ways:

1. **Renew:** sign with the current key.
2. **Change key:** sign with the current key and the new key, the current
   key's signature naming the new one as its bound key. The verification
   carries over.
3. **Start over:** sign with a new key alone, after losing the current one.
   This needs a verification for the new key (§7.1, step 2).

A subject SHOULD change keys at renewal, when it is already making a request.

**What the authority keeps.** For each address, the subject's **lineage**:
the keys it has bound in turn, by thumbprint, since the lineage last started.
The last is the current key. The lineage is the authority's own record:
requests carry only the keys they are signed with, and envelopes only the
current registration.

**How the authority decides.** Only confirming signatures made within the
past 5 minutes, the maximum age of §6 (allowing its clock skew), count: a
replayed earlier request does not show that anybody holds its keys now. Then:

- **The current key alone:** a renewal.
- **One key not in the lineage:** a first registration, a lost key, or a new
  holder of the address; the authority cannot tell which. A new lineage
  starts with that key.
- **The current key and one not in the lineage, the current key's signature
  naming the other as its bound key:** a key change. The authority adds the
  new key and binds it. The naming shows that the holder of the current key
  chose its successor, as in an ACME key change (RFC 8555 §7.3.5); two
  signatures alone would not, since neither covers the other's key.
- **Anything else** MUST be rejected: among others, two keys at a first
  registration, and a key bound before, which cannot become current again.

An authority MUST NOT register a subject unless a verification
countersignature binds a key in the current lineage. A key change therefore
needs no new verification, while a new lineage needs one for its key, so a
new holder of the address cannot be registered as the old one.

A verifier relied on directly (§7.2) has no lineage behind it, so a subject
relying on one needs a new verification whenever it changes key. A subject
moving to another authority starts a new lineage there, and needs a
verification for its current key.

## 8. Endpoints

### 8.1 `GET /.well-known/gobl/keys/<kid>`

Public. Returns one JSON Web Key as `application/json`.

| Status | Meaning |
|---|---|
| `200` | Key returned. |
| `404` | Key is unknown or removed. |
| `429`, `5xx` | Temporary failure. |

### 8.2 `GET /.well-known/gobl/who`

Requires a request token (§6). Returns the participant's party envelope, made
for the requester (§5.1), or one of these empty responses:

| Status | Meaning |
|---|---|
| `200` | Party envelope returned. |
| `202` | Request recorded for possible later disclosure. |
| `204` | The participant does not publish identity details. |
| `401` | Request token missing or invalid. |
| `404` | Address does not participate. |
| `503` | A required key could not be fetched. |

A requester MUST check a `200` response as §5.1 describes before using it.
Responses SHOULD use `Cache-Control: private`. A requester may cache a
checked identity, but only for a bounded time, so the cache does not outlive
the countersignatures behind it.

### 8.3 `POST /.well-known/gobl/inbox`

Accepts a signed GOBL envelope; no request token. The envelope MUST be
addressed to the inbox's address and confirmed as §5.2 describes. The
receiver may then apply its sender status policy (§7.2).

| Status | Meaning |
|---|---|
| `202` | Envelope accepted and durably recorded. |
| `400` | Body is unreadable or invalid JSON. |
| `401` | The envelope's confirming signature failed. |
| `403` | Sender policy rejected the envelope. |
| `422` | Envelope failed structural validation. |
| `500` | Persistence failed. |
| `503` | A required key could not be fetched. |

Acceptance MUST be idempotent on the envelope UUID and digest: the same
envelope delivered twice is acknowledged twice and stored once, which makes
retries safe (§11). Only `202` means the envelope was accepted; any other
`2xx` MUST NOT be read as success.

## 9. Limits

These protect a participant from a hostile or broken peer. Raising them
weakens the implementation, not the protocol.

| Limit | Value |
|---|---|
| Response body accepted from a peer | 1 MiB |
| Signatures per envelope | 32 |
| Signatures checked against keys, per issuer per envelope | 4 |

An envelope with more signatures than the limit MUST be rejected rather than
partly checked. Signatures by one issuer beyond the second limit are ignored,
so an envelope cannot make a receiver fetch more than a few keys from one
address. Each countersignature adds to an envelope, so a participant about to
countersign SHOULD refuse, and say so, if its signature would reach the
limit.

Fetching a key or an identity means requesting an address a peer chose, which
could be used to reach internal services. An implementation MUST refuse to
fetch from hosts resolving to an address that is not globally reachable
(loopback, private, link-local, multicast, unspecified, or listed as not
globally reachable in the IANA special-purpose address registries), or to a
translation range such as NAT64 that can reach one, unless an operator has deliberately configured
otherwise for testing or a local service.

Redirects MUST NOT be followed. A redirect could move a key fetch, a request
token, or a delivery to plain HTTP or to another host, so it is a definitive
failure (§11), and endpoints MUST answer at the URLs of §3 directly.

Fetched keys MAY be cached, by URL and only on success. A cache SHOULD be
bounded in entries and lifetime (§12).

## 10. Typical delivery flow

1. The sender obtains the receiver's address out of band.
2. The sender may call the receiver's `/who` for its invoicing details.
3. The sender sets `from` and `to` to its own and the receiver's addresses,
   and signs the envelope as issuer with the receiver as audience.
4. The sender posts it to the receiver's inbox.
5. The receiver checks the confirming signature; the sender is the `from`
   address (§5.2).
6. The receiver may fetch the sender's `/who` and decide its status (§7.2).
7. A `202` confirms durable receipt.

A reply is a new delivery, subject to the same rules.

## 11. Failure handling

Every failure falls into one of three classes, and an implementation MUST be
able to tell them apart:

| Class | Examples | Caller behaviour |
|---|---|---|
| **Temporary** | `429`, `5xx`, connection and timeout failures | Retry with backoff. |
| **Definitive** | a `3xx` redirect, a `4xx` other than `429`, a signature that fails to verify | Do not retry without changing the request. |
| **Absent** | `204` from `/who`, `404` for an unknown key or address | Not a failure: the peer answered that there is nothing. |

A check that could not be performed is never reported as one that failed. An
unreachable key endpoint means a signature is unknown, not invalid: a server
answers `503` rather than `401` (§6), and a receiver does not downgrade a
sender's status (§7.2). Otherwise somebody else's outage becomes a rejection
the caller cannot tell from a refusal.

Delivery retries are safe because inbox acceptance is idempotent (§8.3), and
a caller that does not retry can lose a document to a temporary failure.

## 12. Security considerations

- **Web PKI:** whoever controls an address and its TLS certificate controls
  its published keys. An implementation MUST authenticate the HTTPS response
  when fetching a key. Applications must compare the senders they see with
  the ones they expect, or apply a sender status policy (§7.2).
- **Key removal:** a cached key may still be accepted after it is withdrawn,
  until the cache entry expires. Key caches SHOULD be short-lived, and an
  implementation SHOULD offer a way to flush them during an incident.
- **Expiry:** countersignature expiries bound how long they can be relied on.
  Identity caches must be bounded too, or a cached identity outlives the
  countersignatures behind it.
- **Verification is not graded:** a subject is verified or not, and any
  verifier a trusted authority names counts. A receiver wanting to
  distinguish verifiers must read the `verifier` and decide which it accepts.
  An authority that names a weak verifier weakens every subject it registers
  with it.
- **Verifier scope:** a verification describes the subject at its address
  when it was signed. It says nothing about documents the subject sends
  later. An expiry, if present, bounds how long it may be relied on; without
  one, the authority decides how old a verification it accepts.
- **Document content:** this protocol establishes who sent an envelope and
  who that sender is, not that the parties named inside the document are who
  they claim. A verified sender can still issue an invoice naming another
  organization as its supplier, so receivers must reconcile the document's
  parties, such as tax identities and legal names, with the sender's verified
  party before acting on it.
- **Address ownership:** a domain can change hands. A new holder controls the
  published keys but not the registered one, so earlier countersignatures do
  not carry over to it (§7.3). It can be verified and registered afresh, but
  not as the previous holder. An authority MUST re-confirm control of the
  address on every renewal.
- **Hosted addresses:** whoever operates an address holds its keys. A
  participant using an address under a provider's domain trusts that
  provider to act as it.
- **Unsigned header data:** `from` and `to` count only when a signature
  confirms them, and even then show only that the sender addressed the
  envelope to the receiver at some point, not that nobody else holds it.
  Stamps, notes, links, and other header fields are not signed at all, and
  are only the sender's claim.
- **Replay:** a delivery may be posted again, intentionally or not. Inbox
  acceptance is idempotent (§8.3), and a receiver SHOULD recognise documents
  it has already processed by their UUID for as long as it keeps them.
- **Privacy:** request tokens identify requesters but do not make party data
  public. Operators should limit how long they keep request logs, which can
  reveal business relationships.

### 12.1 Bindings

A valid signature only says what it binds. Each fact a participant relies on
is tied to the document, the addresses, the subject's key, and the present by
the checks below. A check with any of these missing lets a valid signature be
reused for something it never said, so a new fact or claim needs its own row,
with each column filled or a reason it need not be.

| Fact | Evidence | Document | Addresses | Key | Time |
|---|---|---|---|---|---|
| `from` sent this envelope to `to` (§5.2) | The confirming signature | UUID and digest, and the digest matches the document | Issuer is `from`, audience is `to`, and `to` is the receiver | Published by `from` | Signing time not in the future; key valid at receipt |
| The party at `/who` is what the subject presents to this requester (§5.1) | The confirming signature on the response | As above | Issuer is the address queried, audience is the requester | Published by the subject | As above |
| A verifier vouches for the party at this address (§7) | A verification countersignature | As above | Audience is the subject | Bound key confirmed the party envelope | Signing time not in the future; expiry, if any, not passed |
| An authority registered the subject (§7.2) | A registration countersignature | As above | Audience is the subject; the named verifier's countersignature has the same audience | Bound key confirmed the party envelope | Signing time not in the future; expiry present, not passed, at most 90 days later |
| The verified subject sent this delivery (§7.2) | The delivery's confirming signature | As in the first row | `from` is the subject | The bound key confirmed the delivery | As in the first row |
| A requester asking `/who` is who it says (§6) | A request token | None | Issuer is the requester, audience is the address queried | Published by the requester | Lifetime and age limits of §6 |

Whose signatures count is never read from the envelope: only the trust lists
decide which authorities and verifiers are evidence, and live and sandbox
lists are never combined (§7).

Nothing binds these, so they are never evidence:

- party endpoints (§5.1), header fields other than `from` and `to`, and
  `from` and `to` without a confirming signature;
- the order of signatures (§5.1);
- the parties named inside a document, which receivers reconcile with the
  sender's verified party (§12, document content); and
- how thorough a verification was, beyond the verifier's address (§12).

## 13. References

- RFC 2119 and RFC 8174: requirement keywords
- RFC 6750: Bearer token usage
- RFC 7515: JSON Web Signature
- RFC 7517: JSON Web Key
- RFC 7518: JSON Web Algorithms
- RFC 7519: JSON Web Token
- RFC 7638: JSON Web Key Thumbprint
- RFC 7800: Proof-of-Possession Key Semantics for JWTs (`cnf`)
- RFC 8555: Automatic Certificate Management Environment (ACME)
- `github.com/invopop/gobl/dsig`: signature and key primitives
- `github.com/invopop/gobl/org`: party schema
