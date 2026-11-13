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
authenticate the HTTPS response. That establishes control of an address and
nothing more. Two optional roles build on it by countersigning party
envelopes: a registration authority, which binds a party to an address it
controls; and a verifier, which establishes that the party describes a
real-world organization — on a live network, by a KYC/KYB process (§7).

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
| **Subject** | The participant a party envelope belongs to. Its address is named by the envelope header's `from` field and confirmed by signatures (§5.1), never read from the party document. One party may serve several subjects, each holding its own copy of the envelope. |
| **Issuer** | The participant identified by a signature's `iss` claim. This is a signature-level role. |
| **Audience** | The participant identified by a signature's `aud` claim. The signature is bound to this address. |
| **Sender** | The participant named by an envelope's `from`, confirmed by a valid signature with that `iss` and the envelope's `to` as its `aud` (§5). |
| **Receiver** | The participant named by an envelope's `to`: the inbox it is delivered to, or the requester a `/who` response is for. |
| **Requester** | The issuer of a `/who` request token (§6): the participant the response is for. |
| **Authority** | A trusted participant that confirms the subject controls its GOBL Net address and countersigns its party envelope. |
| **Verifier** | A participant, independent of the Authority, that establishes that the details in a party document belong to a real-world organization or person — on a live network, by a KYC/KYB process — and that this organization controls the Subject's address, and countersigns the party envelope for that address to say so. It binds a party document to a legal identity at one address. What it checks, and how, is its own to define and publish; see §7. |
| **Registration countersignature** | A signature issued by an Authority over a party envelope, binding it to the Subject's address with `aud` and to the Subject's key with `cnf`. It may name a Verifier. |
| **Verification countersignature** | A signature issued by a Verifier over the same party envelope, binding it to the Subject's address with `aud`, and naming the key the Subject presented with `cnf`. Its `exp` is optional. |
| **Sender status** | What a receiver concludes about a Subject after evaluating the countersignatures its party envelope carries, by the procedure in §7.2. One of the three below. |
| **Unregistered** | The receiver found no valid registration countersignature from an Authority it trusts. It has only the Subject's own account of itself. This is a statement about what *this* receiver can corroborate, not about whether the Subject registered with somebody. |
| **Registered** | An Authority the receiver trusts has countersigned the party envelope, attesting that the Subject controls its address. Who is behind that address is not established. |
| **Verified** | Registered, and that Authority named a Verifier whose own countersignature the receiver found valid on the same envelope — so a participant independent of the Authority has established who the party is. |
| **Party envelope** | A signed GOBL envelope containing an `org.Party`. It carries the subject's signature for its current recipient, and may carry registration and verification countersignatures. |

A countersignature is a signature on the wire; a sender status is a conclusion
a receiver reaches about one. The distinction matters because **no status is
ever transmitted**. Unregistered, registered and verified name conclusions,
not values: no field carries them, and nothing in an envelope asserts one.
Each receiver derives its own from the signatures in front of it and the
authorities it happens to trust, so two receivers can read the same envelope
and reach different conclusions about it, and both are correct.

Sender and Receiver are the `iss` and `aud` of the signature confirming an
envelope's `from` and `to`; countersignatures on the same envelope have their
own issuers and audiences. Exchanges are peer to peer: nothing is relayed, so
there is no intermediary role. A countersignature's Audience is always the
Subject.

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
`gobl:` URI is used in optional `org.Party` endpoints and optional envelope
routing fields. Production endpoints MUST use HTTPS.

## 4. Keys

Each signing key MUST be available as a single RFC 7517 JWK at
`/.well-known/gobl/keys/<kid>`. Unknown or removed key IDs return `404`. The
optional `/.well-known/jwks.json` endpoint exposes all keys as an RFC 7517 JWK
Set for generic tooling; GOBL Net verification uses the per-key endpoint.

A published key may include `valid_from` and `valid_until`. A signature MUST
be rejected if its `iat`, which every GOBL Net signature carries (§5), falls
outside those bounds; a missing bound skips its comparison. Because the
signer chooses `iat`, a key that has been retired could still sign with a
backdated one, so a delivery or `/who` response MUST also be signed with a key
whose bounds include the time it is received (§5.2). Retired keys SHOULD remain
available while countersignatures made with them are still within their
`exp`, and for checking historical signatures.

Signatures identify keys with `kid`; they do not carry a `jku` header.

## 5. Envelope signatures

The signed payload contains the document `uuid` and digest plus these protocol
claims:

| Claim | Meaning |
|---|---|
| `iss` | Signer's canonical GOBL Net address. |
| `aud` | Address to which the signature is bound; required on the signature confirming an envelope's `from` and `to`, and on countersignatures; optional only on envelopes exchanged outside GOBL Net. |
| `iat` | Signing time, as an RFC 7519 NumericDate. Required, and MUST NOT be later than the present, allowing the clock skew of §6. |
| `exp` | Expiry. Required on registration countersignatures, at most 90 days after `iat`; optional otherwise. |
| `verifier` | Optional verifier named by a registration authority. |
| `cnf` | RFC 7800 confirmation, `{"jkt": "<thumbprint>"}`: the subject key a countersignature is bound to, by its RFC 7638 SHA-256 JWK thumbprint, base64url-encoded. Required on registration and verification countersignatures. |

Every GOBL Net envelope MUST carry header `from` and `to` fields holding
`gobl:` addresses, and a valid signature whose `iss` is the `from` address
and whose `aud` is the `to` address. An envelope missing either field, or
using another transport's scheme in either, is not meant for GOBL Net, and a
participant MUST reject it rather than verify it. An envelope with no
signature confirming the pair MUST be rejected.

Neither field is signed itself; the confirming signature is what makes them
evidence. Other signatures on the envelope, for other pairs or none, change
neither. A participant passing an envelope on does not forward it: it sets
`from` to its own address, `to` to the next recipient, and adds its own
signature for that pair.

Only countersignatures (§7) mean anything beyond the exchange that carries
them. Once an envelope is accepted, any other signature confirming a pair
has served its purpose: a holder MAY remove it, and SHOULD before passing the
envelope on, so envelopes do not accumulate signatures towards the limit in
§9.

Signature order has no meaning. Implementations MUST search all signatures and
MUST NOT assign roles by array position.

### 5.1 Party identities

A party envelope contains an `org.Party`. The subject's address is not part
of the party document: it is the envelope's `from`, confirmed as §5 requires.
Because `from` and `to` are not signed, a party can be verified before it has
any address, and one party can serve several addresses, each holding its own
copy of the envelope. Copies share a digest, so the countersignatures on one
are valid on all.

The party MAY list `gobl:` endpoints. They are hints for display and
discovery only, and MUST NOT be used to determine or check the subject.
Equally, any holder can reorder or add signatures, so the subject MUST NOT be
inferred from whichever one appears first.

The confirming signature establishes only that the subject adopted the party
document as its own; it says nothing about registration, which is §7.

A `/who` response is made for each requester: its `from` MUST be the address
queried, its `to` the requester — the `iss` of the request token (§6) — and
the subject signs it with `aud` equal to the requester. The requester MUST
reject a response from or to anybody else. A server does not keep these
per-requester signatures in the envelope it stores (§5).

Other self-signatures and countersignatures may be present.

### 5.2 Deliveries

For an inbox delivery, the receiver MUST reject an envelope whose `to` is
not its own address, or whose `from` and `to` no signature confirms (§5) with
a key that is valid at the time of receipt (§4). The
sender is the `from` address. Signatures by other issuers whose `aud` is the
receiver say nothing about who sent it.

Outside GOBL Net — a party envelope handed to a verifier through its own API,
say — no `aud` is required and these rules do not apply.

## 6. Request authentication

Requests to `/who` MUST include a compact ES256 request token as
`Authorization: Bearer <token>`, which tells the server who to address and
sign the response for (§5.1). Inbox deliveries carry no token: the envelope's
own signature authenticates the sender (§5.2), and re-posting an envelope
gains nothing because acceptance is idempotent (§8.3). Key discovery is
unauthenticated.

The token contains:

| Claim | Requirement |
|---|---|
| `iss` | Requester's canonical address. |
| `aud` | Address queried. |
| `iat` | Required issue time. |
| `exp` | Required expiry. |
| `jti` | Optional unique identifier, to let a receiver detect replay. |

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

## 7. Authorities, verifiers, and sender policy

A registration authority countersigns a party envelope with:

- `iss` equal to the authority address;
- `aud` equal to the subject's address, which the authority has confirmed the
  subject controls;
- `iat`, and an `exp` no more than 90 days after it;
- `cnf` naming the subject's key that the authority has registered for that
  address (§7.3); and
- optionally, `verifier` naming the address that performed KYC/KYB.

The `cnf` binds the registration to the holder of a key as well as of an
address. Whoever holds an address later — after a domain lapses and is
registered again, say — controls its published keys, but cannot sign with
the registered one, so it cannot present the earlier registration as its
own.

The bounded `exp` limits how long a registration outlives the subject's
control of its address, and how long a leaked, retired authority key can be
used to forge one with a backdated `iat`. Like TLS certificates issued
through ACME (RFC 8555, the protocol behind Let's Encrypt), registrations
are short-lived and meant to be renewed automatically: a
subject SHOULD renew once two thirds of the lifetime has passed (§7.1).

The `aud` is what binds a registration to one address. Without it, any
address could add its own self-signature to another party's published
envelope and present that party's registration as its own. A subject using
one party at several addresses keeps a copy of the envelope at each, with
that address's registration countersignature.

The verifier MUST be a different participant from the registration authority
and MUST countersign the same envelope, with `aud` equal to the subject's
address and `cnf` naming the subject's key it was presented with. That `aud`
is what stops a verification travelling with a copied envelope: without it,
anybody could take a verified party's envelope, register it at an address of
their own, and have an authority name the verifier for it. The `cnf` does the
same across a change in who holds the address, and an authority follows the
subject's key from it as the subject rotates keys (§7.3), so a verification
never needs renewing for a new key. Separating these roles prevents the
authority from attesting to identity checks that no independent verifier
performed.

**A verifier binds a party document to a legal identity.** The other layers
establish control: the keys in §4 are trusted because the TLS certificate at
that address vouches for them, and the authority's countersignature attests
that the subject holds the address and has adopted the party. A domain anyone
can register, with keys anyone can publish, shows only that one participant
holds both. Confirming that the party document describes a real organization
is the verifier's work, and the only part of this protocol whose evidence
comes from outside the network.

On a live network that means a **KYC/KYB process**: confirming the
organization exists in a company registry, that the identifiers and legal name
the party publishes match what the registry holds, and that whoever presented
it is entitled to act for it. A verifier SHOULD publish what it checks,
because its address is the only signal a receiver gets about the strength of
that check (§12).

A sandbox network is deliberately weaker. A verifier there may do something as
slight as confirming control of the email address the party publishes, because
the point is to exercise the exchange rather than to establish anything about
anybody. Such a verifier belongs to a sandbox network only, and a live
authority MUST NOT name one in a `verifier` claim. The authority is the gate
here, not the receiver: §7.2 asks only that a named verifier be independent
and that its countersignature verify, so a receiver takes the authority's word
on which verifiers are worth naming. An endorsement carries no grade (§12), so
it cannot tell a mailbox check from a registry check by looking. The disjoint
trust lists below are the other half of that separation — a receiver trusting
only live authorities never sees a sandbox endorsement at all.

Each receiver holds a **trust list** of the authorities it accepts. The
default live authority is `lookup.gobl.org` and the default sandbox authority
is `lookup.sandbox.gobl.org`. The two lists are disjoint and MUST NOT be
combined: a sandbox registration carries the same shape as a live one, so a
receiver that trusted both could not tell a test identity from a real one.

Sender status is a **policy for sending**, not a prerequisite for receiving. A
receive-only address may publish only an inbox and keys, and may return `204`
from `/who`.

### 7.1 The verification exchange

A verifier establishes who a party is and countersigns the party envelope to
say so. The exchange is **subject-mediated**: the subject holds its own
envelope throughout and presents it to each participant in turn.

1. The subject prepares its party envelope — an `org.Party` with the details
   to be verified.
2. It asks a verifier of its choosing to verify it for an address. The
   verifier performs its checks, by whatever means and over whatever period
   it needs. They MUST establish both that the party document describes the
   organization presenting it, and that this organization controls the
   address. A delivery from the address (§5) proves control, but not who
   holds it: the verifier MUST tie the two together itself, for example by
   having the delivery made within a verification it has already
   authenticated. The key that signed that delivery is the one the verifier
   binds. The checks may start before the subject has an address; the
   countersignature cannot.
3. On success the verifier countersigns **that exact envelope**, unchanged in
   `uuid` and digest, with `iss` equal to its own address, `aud` equal to the
   subject's address, `cnf` naming the subject's key from step 2, and an
   `iat`. It MAY add an `exp`. It does not set a
   `verifier` claim: that pointer is the authority's to write.
4. The verifier returns the countersigned envelope to the subject. Delivering
   it to the subject's inbox (§8.3) is RECOMMENDED: the verifier sets `from` to
   its own address and `to` to the subject's, and its countersignature
   confirms the pair (§5).
5. The subject signs the envelope with `aud` equal to the authority and
   delivers it to the authority's inbox, with `from` its own address and `to`
   the authority's, which shows the authority that the subject controls that
   address. The authority records the key that signed the delivery (§7.3). It
   recognises a countersignature for that address from a verifier on its
   trusted list, whose `cnf` names a key in the subject's lineage, and
   countersigns with `aud` equal to the subject's address, `cnf` naming the
   recorded key and the verifier named in the `verifier` claim. It delivers
   the result back
   with `from` its own address and `to` the subject's, which its
   countersignature confirms. The subject publishes it at `/who` (§5.1).

Signatures are independent of one another and of their order, so step 5 can
happen long after step 3. A further address needs its own verification
countersignature, so its copy of the envelope goes through steps 2 to 5; a
verifier that has already checked the organization need only establish that
the same organization controls the new address.

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
verifier. A verification countersignature with no `exp` remains usable for as
long as the authority is willing to name it; an authority MAY apply its own
limit on its age, using `iat`. An authority SHOULD NOT name a verifier whose
countersignature was made with a key that is no longer valid. Registering, adding verification later, and
renewing are therefore one operation applied at different times, not three.

How the envelope gets back to the subject in step 4 is not constrained. What
the exchange requires is that the subject ends up holding that exact envelope;
the means are between a verifier and its own customers. A verifier that
returns it by some other route — its own API, a link in an email, a console
the subject signs into — conforms, as does one that reports progress or
history however it likes. None of that belongs in this specification.

Inbox delivery is recommended because it costs neither side anything new.
A party envelope moving between two participants is what §8.3 is for, a
subject with an address already operates an inbox by virtue of being on the
network, and it
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

A verifier needs no relationship with any authority. It checks that the
details in the party document belong to the organization presenting them, and
that this organization controls the address, and countersigns on that basis.
It does not need to know which authority the subject will register with, whether the subject has
registered at all, or whether that authority would accept its
countersignature — the last of which is the authority's decision and is made
in §7.2, after the verifier's work is done. A verifier MAY still require the
subject to be registered before it will act, but that is its own admission
policy and not a rule of the protocol.

A verifier returns its result to the subject and never to the authority:
only the subject registers its envelope, and no participant delivers an
envelope on another's behalf (§5).

The subject cannot subvert any of this by carrying the envelope. The
countersignature covers the digest, so altering anything invalidates it, and
the authority verifies it against the verifier's published key. Withholding is
the subject's right — it is their identity.

This exchange requires nothing of an authority beyond what §7 already
specifies.

### 7.2 Determining sender status

A receiver reaches a sender status by evaluating the countersignatures a party
envelope carries against its own trust list. The procedure is normative: two
conforming receivers with the same trust list MUST reach the same status for
the same envelope.

1. Consider every signature whose `iss` is an authority on the trust list and
   whose `aud` equals the subject. An envelope carrying none is
   **unregistered**, whatever else it carries — a countersignature from an
   unknown authority, or one bound to a different address, is not evidence.
2. For each candidate, verify it against that authority's published key, and
   reject it unless it has an `iat` no later than the present, allowing the
   clock skew of §6, an `exp` that has not passed and is no more than 90 days
   after the `iat`, and a `cnf` naming the key of a signature confirming the
   envelope's `from` and `to` (§5).
3. A candidate that verifies makes the subject at least **registered**.
4. If a verifying candidate names a `verifier`, and that verifier is a
   different participant from the authority, and the same envelope carries a
   valid countersignature issued by that verifier with `aud` equal to the
   subject and a `cnf` — with an `iat` no later than the present, allowing
   the clock skew of §6, and no `exp` that has passed — the subject is
   **verified**. Which key the `cnf` names is the authority's to check
   (§7.3); the receiver relies on the authority having named the verifier.
5. Where candidates disagree, the highest status wins: an envelope may carry
   several authority signatures, and one that establishes **verified** is not
   cancelled by another that only establishes **registered**.

A `verifier` claim that is absent, malformed, self-referential, or whose
countersignature is missing, invalid, expired, lacks an `iat` or a `cnf`,
or is for another address reduces the
status to
**registered**. It MUST NOT invalidate an otherwise valid registration — a
claim nobody can corroborate is worth nothing, but it does not make the
authority's own statement worthless.

The one case that is not a downgrade is a verifier key endpoint that cannot be
reached. "Could not check" is not "not verified": a receiver MUST NOT report
**registered** on a temporary failure, and SHOULD treat it as §11 does, so a
verified sender is never quietly demoted by somebody else's outage.

A status describes the holder of the key the registration binds. A receiver
applying it to a delivery MUST check that the delivery is confirmed (§5) by
that same key, or the status does not apply to it.

### 7.3 Key lineage

An authority keeps, for each address it registers, the **lineage** of the
subject's keys: the keys it has bound in turn, by thumbprint, since the
lineage last started. Every registration request is a delivery (§5.2) to the
authority's inbox, and the keys confirming its `from` and `to` decide what
happens:

- **One key, in the lineage's current position or with no lineage yet.** A
  first registration or a renewal: the authority binds that key. A first
  registration starts the lineage with it.
- **One key, not the current one.** The subject has lost its key, or the
  address has changed hands; the authority cannot tell which. It starts a new
  lineage with that key, and registers it.
- **Two keys, one of them the current one.** A rollover: the authority adds
  the other key to the lineage and binds it. Signing with both keys proves
  the holder of the current key chose its successor, as in an ACME key
  change (RFC 8555 §7.3.5).
- **Two keys at a first registration.** The authority starts the lineage
  with whichever key a verification countersignature's `cnf` names, and
  rolls over to the other. If neither is named, it MUST reject the request.
- **Anything else** MUST be rejected.

An authority MUST NOT name a verifier unless that verifier's countersignature
binds, with `cnf`, a key in the current lineage. A rollover therefore keeps a
subject verified with no new verification, and a new lineage drops it: a
subject that has lost its key asks its verifier again (§7.1, step 2), and a
new holder of the address cannot be verified as the old one.

A lineage is the authority's own record, so a subject that moves to another
authority after rotating its key starts a new lineage there, and needs a
verification for its current key. A subject SHOULD rotate keys at renewal,
when it is already making a request.

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

A `200` response MUST be from the queried address to the requester, and MUST
be verified as described in §5.1 before anything in it is used. Authenticated responses SHOULD use `Cache-Control: private`, and a
caller may cache a verified identity for a bounded period — bounded because a
cached identity outliving the countersignature that endorsed it is §12's
problem.

### 8.3 `POST /.well-known/gobl/inbox`

This endpoint requires a signed GOBL envelope and no request token.
The envelope MUST be addressed to the inbox's address and confirmed as
described in §5.2. The receiver may then apply its sender
status policy (§7.2).

| Status | Meaning |
|---|---|
| `202` | Envelope accepted and durably recorded. |
| `400` | Body is unreadable or invalid JSON. |
| `401` | Delivery signature failed. |
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
| Signatures checked against keys per issuer, per envelope | 4 |

An issuer's signatures beyond the second limit are ignored, so an envelope
cannot make a receiver fetch more than a few keys from any one address. An
envelope carrying more signatures than the first limit MUST be rejected rather
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
3. The sender sets `from` and `to` to its own and the receiver's addresses,
   and signs the envelope with `iss=sender` and `aud=receiver`.
4. The sender posts it to the receiver's inbox.
5. The receiver verifies the envelope; the sender is its `from` (§5.2).
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
- **Routing fields:** `from` and `to` are not signed; they are accepted only
  when a valid signature carries the same pair as `iss` and `aud`. That
  shows the sender addressed the document to the receiver at some point, not
  that the receiver is the only one holding it.
- **Verification is not graded.** A subject is verified or it is not, and any
  verifier the authority names satisfies it. A receiver that wants to
  distinguish between verifiers must read the `verifier` claim itself and
  decide which addresses it accepts. It follows that an authority admits a
  weak verifier at the cost of every sender it marks verified, and that a
  verifier's address is the only trust signal a receiver gets: it should name
  what is being vouched for, and a change in what the verifier checks is a
  change in what that name means.
- **Verifier scope:** a verifier attests to a party's identity at the
  address in its `aud`, at the point in time given by its `iat`. It says
  nothing about any document the subject later sends. An `exp`, where
  present, bounds how long the attestation may be relied upon at all; where
  absent, the authority decides how old an attestation it is willing to name.
- **Document content:** this protocol establishes who sent an envelope and
  who that sender is. It does not establish that the parties named inside the
  document are who they claim: a verified sender can still issue an invoice
  naming another organization as its supplier. As on any delivery network,
  receivers must reconcile the document's parties — tax identities, legal
  names — with the sender's verified party before acting on it.
- **Address ownership:** an address is a domain, and a domain can change
  hands. A new holder controls its published keys but not the registered
  one, so registration and verification countersignatures bound with `cnf`
  do not carry over to it (§7.3): it can be registered afresh, but not
  verified as the previous holder. An authority MUST re-confirm control of
  the address on every renewal.
- **Hosted addresses:** whoever operates an address holds its keys. A
  participant using an address under a provider's domain is trusting that
  provider to act as it.
- **Unsigned header data:** stamps, notes, links and other header fields are
  not covered by any signature, and are only the sender's assertion.
- **Replay:** a delivery may be posted again, intentionally or not. Inbox
  acceptance is idempotent (§8.3), and a receiver SHOULD recognise documents
  it has already processed by their UUID for as long as it keeps them.

## 13. References

- RFC 2119 and RFC 8174 — requirement keywords
- RFC 6750 — Bearer token usage
- RFC 7515 — JSON Web Signature
- RFC 7517 — JSON Web Key
- RFC 7518 — JSON Web Algorithms
- RFC 7519 — JSON Web Token
- RFC 7638 — JSON Web Key Thumbprint
- RFC 7800 — Proof-of-Possession Key Semantics for JWTs (`cnf`)
- RFC 8555 — Automatic Certificate Management Environment (ACME)
- `github.com/invopop/gobl/dsig` — signature and key primitives
- `github.com/invopop/gobl/org` — party schema
