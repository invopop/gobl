# GOBL Net

> ⚠️ **EXPERIMENTAL** — GOBL Net is under active development. The package
> API, CLI commands (`gobl net`, `gobl init`), on-disk layout, and the wire
> protocol may change without notice and are not yet covered by any
> stability guarantee.

**Status:** Draft. This document describes the current state of the implementation
in the `github.com/invopop/gobl/net` package, the `internal/ops` helpers that
back the `gobl net` CLI, and related supporting code in `dsig`. The protocol
is pre-1.0 and subject to change.

## Abstract

GOBL Net is a decentralized identity and discovery protocol for signed GOBL
documents. It binds a document signature to a fully qualified domain name
(FQDN), and defines a small set of well-known HTTPS endpoints at that domain:

- `/.well-known/gobl/keys/<kid>` — a single public JWK looked up by key
  ID, used to verify signatures.
- `/.well-known/gobl/who` — a signed GOBL Envelope carrying an
  `org.Party` document, endorsing the holder's identity.
- `/.well-known/gobl/inbox` — a write endpoint that accepts signed
  envelopes addressed to the holder.

Trust in an identity is anchored to a list of approved *KYC vendor*
addresses. Trust is not anchored to specific public keys, so a vendor may
rotate keys at will without breaking endorsements.

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
- **Published key** — A `dsig.PublicKey` (an RFC 7517 JWK plus the optional `valid_from` / `valid_until` extension members) served at `/.well-known/gobl/keys/<kid>`.
  Represented in code as `net.KeySet`.
- **Party Envelope** — A signed GOBL Envelope whose document is an
  `org.Party`, exchanged at the who endpoint.
- **iss / aud** — Fields in a signature's *signed payload* carrying the
  verifiable GOBL Net origin (`iss`) and the address the signature is
  bound to (`aud`), both as `gobl:` `cbc.URI` values. These are the
  authoritative, tamper-proof identities.
- **header from / to** — Optional unsigned `cbc.URI` fields on the
  envelope header expressing *intent/routing* in any scheme (`peppol:`,
  `mailto:`, `gobl:`…). Useful for interop with other formats; not used
  for verification.
- **allow-list** — Optional `<domain>/allow.json` (array of addresses)
  restricting which callers a domain accepts on `/who` and `/inbox`.

## 3. Addressing

### 3.1 Address Format

A GOBL Net Address is an FQDN. It MUST NOT contain a scheme, port, path,
query, or fragment.

Parsing is performed by `net.ParseAddress` and applies the following
normalizations and constraints:

1. Surrounding whitespace is trimmed.
2. The address is lowercased.
3. A trailing dot, if present, is stripped.
4. The result MUST satisfy `is.DNSName` (RFC 1035 label syntax).
5. The result MUST contain at least one dot (i.e. at least two labels).

Inputs that contain a scheme, port, or path MUST be rejected with
`ErrAddressInvalid`. An empty input MUST be rejected with `ErrAddressEmpty`.

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
   Set, used by generic JWT tooling (jwt.io, OIDC-style verifiers)
   that resolves the `jku` header on each signature.

The scheme MUST be `https`. HTTP is not a permitted alternative for
production deployments; the `gobl net send --insecure` flag exists solely
for local development.

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

`valid_from` is stamped automatically when a key is generated (by
`gobl init` or `gobl net serve`'s autogeneration). `valid_until` is
left empty and is meant to be set when the operator rotates the key out
— either at retirement, or in advance of a planned rotation. A retired
key remains published so that historical envelopes signed within its
window still verify; only signatures with an `iat` past `valid_until`
are rejected.

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

On-disk, each domain stores its published keys as individual files
under `<config-dir>/<domain>/keys/`, named by the key's `kid`:

```
~/.config/gobl/billing.invopop.com/
├── private.jwk                              ← the active signing key
├── keys/
│   ├── 9f8b…json                            ← currently published key
│   └── 4c10…json                            ← retired key (valid_until set)
├── party.json
├── allow.json
└── inbox/
```

The server reads the directory on startup, validates that every
filename matches the `kid` in its JWK, and serves each file verbatim
at `/.well-known/gobl/keys/<kid>`. Rotation is just file ops: drop a
new `<kid>.json` to publish a key, edit one to set `valid_until` to
retire it, or `rm` it entirely to drop it from the published set
(future requests for that `kid` return `404`). The model is a 1:1
mapping to a future database table — one row per `kid`.

Endorsement of a participant happens at the identity layer (see §6) and
does *not* live inside the key material itself.

## 5. Signatures

### 5.1 Signed identities: `iss` / `aud` / `iat`

Each signature signs a payload of `{uuid, dig, iss, aud, iat}`:

- `uuid` + `dig` identify the document (immutable after signing).
- `iss` is the signer's verifiable GOBL Net address as a `gobl:` URI
  (e.g. `gobl:billing.invopop.com`). The verifier reads it to
  discover *which* per-key endpoint to fetch.
- `aud` is the optional GOBL Net address the signature is bound to. When
  present, the recipient checks `aud == self` to reject misrouted or
  replayed envelopes.
- `iat` is the signing time as a JWT-standard NumericDate (Unix
  seconds, per RFC 7519 §2). It is set automatically by `Sign`;
  verifiers read it for the per-key validity window check but no
  freshness policy is enforced by default — receivers may apply their
  own max-age window when relevant.

Because `iss`/`aud`/`iat` are inside the signed payload, the origin,
audience, and signing time are all tamper-proof. Multiple parties may
countersign the same document (shared `uuid`+`dig`) each with their own
`iss`/`aud`/`iat`.

The header `from`/`to` (`cbc.URI`) are a *separate*, unsigned layer for
intent/routing in any scheme; they are never used for verification.

### 5.2 Envelope Signing

`Envelope.Sign(key, iss, aud)` (→ `head.Header.Sign`) signs the document
identity plus the signer's `iss` and optional `aud`. Both may be empty
for a plain, non-GOBL-Net signature.

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
6. If the key declares a validity window, the signed `ts` MUST fall
   within `[valid_from, valid_until]` (each bound optional).
7. The verified issuer address is returned.

### 6.2 Identity exchange (`/who`)

`/who` is an authenticated, mutual exchange (see §8.2). The caller POSTs
a signed envelope (`iss=gobl:caller`, `aud=gobl:target`, document = the
caller's `org.Party`); the target verifies it, applies its allow-list,
and responds with its own party envelope signed `iss=gobl:target`,
`aud=gobl:caller`. The `ops.NetWho` helper / `gobl net who` performs the
exchange and verifies the response is signed by the target and bound to
the caller.

### 6.3 Trusted Authorities (optional)

The package-level slice `net.Authorities` holds GOBL Net addresses
treated as trusted KYC vendors. Endorsement is optional and best-effort:
a party may carry an authority's countersignature in addition to its own
self-signature, but no endpoint *requires* it.

## 7. Discovery Transport

### 7.1 HTTP Client Defaults

The default `HTTPFetcher` enforces:

| Parameter             | Value              |
|-----------------------|--------------------|
| Request timeout       | 10 seconds         |
| Maximum response size | 1 MiB              |
| Required `Accept`     | `application/json` |
| Required scheme       | `https`            |
| Required status       | `200 OK`           |

Responses larger than 1 MiB are truncated. Any non-200 response causes
`ErrFetchFailed`.

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

### 8.2 `POST /.well-known/gobl/who`

Authenticated party exchange. The caller POSTs a signed envelope
(`iss=gobl:caller`, `aud=gobl:self`, document = the caller's
`org.Party`). The server verifies the signature against the caller's
published key (fetched from the caller's per-kid endpoint), requires
`aud == self`, and applies the allow-list to `iss`. It responds `200`
with its own party envelope signed
`iss=gobl:self`, `aud=gobl:caller`.

| Status            | Cause                                              |
|-------------------|----------------------------------------------------|
| `200 OK`          | Verified; returns the signed party envelope.       |
| `400 Bad Request` | Body did not decode as an envelope.                |
| `401 Unauthorized`| Signature/issuer/audience verification failed.     |
| `403 Forbidden`   | Caller (`iss`) not on the domain's allow-list.     |

### 8.3 `POST /.well-known/gobl/inbox`

Accepts a signed GOBL Envelope. The signer (`iss`) is verified against
its published key (fetched from `<iss>/.well-known/gobl/keys/<kid>`);
if the envelope carries an `aud` it MUST equal this inbox; the
allow-list (if present) is applied to `iss`. Status codes:

| Status                       | Cause                                                |
|------------------------------|------------------------------------------------------|
| `202 Accepted`               | Envelope parsed, validated, signature verified, persisted. |
| `400 Bad Request`            | Body could not be read or did not decode as JSON.    |
| `401 Unauthorized`           | Envelope signature did not verify.                   |
| `422 Unprocessable Entity`   | Envelope failed structural validation.               |
| `500 Internal Server Error`  | Persistence failed.                                  |

The request body size is capped at 1 MiB. On 202 the envelope is written
to the configured inbox directory under `<envelope-uuid>.json`.

In this release the server does not return a signed receipt on 202; the
response body is empty. A future `net.Response` type may carry an
acknowledgement body.

## 9. Command Line

### 9.1 `gobl verify`

Supports remote verification via two flags:

- `-a, --address <fqdn>` — Require the verified issuer (`iss`) to equal
  this address.
- `-r, --remote` — Fetch the verifying key from the issuer published in
  the signature's signed `iss`.

### 9.2 `gobl net serve`

Starts the HTTP server described in §8. The server always listens on an
HTTP port (default 80). When a TLS source is configured, it additionally
listens on an HTTPS port (default 443) and serves identical content
there — there is no redirect from HTTP to HTTPS; senders choose the
scheme they want and verify the certificate themselves.

**Config directory.** All on-disk state defaults to a single base
directory (`--config-dir`, default `~/.config/gobl/`):

| File                              | Purpose                                          |
|-----------------------------------|--------------------------------------------------|
| `<config-dir>/keys/<kid>.json`    | One public JWK per `kid`, served at `/keys/<kid>`. Override the directory with `--keys-dir`. |
| `<config-dir>/private.jwk`        | Active signing key paired with one of the published keys. Override with `--private-key`. |
| `<config-dir>/party.json`         | Raw `org.Party` (or a signed envelope) served, signed per-request, at `/who`. Override with `--party`. |
| `<config-dir>/allow.json`         | Optional accept-list (array of addresses) gating `/who` and `/inbox`; absent = accept any verified caller. |
| `<config-dir>/inbox/`             | Directory accepted envelopes are written into. Override with `--inbox`. |
| `<config-dir>/certs/`             | ACME certificate cache. Override with `--cert-dir`. |

(In multi-domain mode these live under `<config-dir>/<domain>/`.)

**Startup behavior:**

- **Auto-key-generation.** If neither `keys/` nor `private.jwk` exists,
  the server generates an ECDSA P-256 keypair, writes `private.jwk`
  (0600) and `keys/<kid>.json` (with `valid_from = now`), and logs the
  new kid plus the paths. If only one of the two exists, startup fails
  with a clear error — the operator's setup is inconsistent.
- **Filename ↔ kid invariant.** Every file in `keys/` MUST be named
  `<kid>.json` where `kid` equals the JWK's `kid` field. A mismatch is
  a startup error. Non-`.json` entries and subdirectories are ignored.
- **Active key check.** The `private.jwk`'s `kid` MUST be one of the
  published kids in `keys/`. Otherwise startup fails.
- **Party self-signature.** The party envelope MUST contain at least
  one signature whose `kid` is present in `keys/` and which verifies
  against that key. Endorser signatures (from KYC vendors) are allowed
  alongside and are served verbatim — they simply don't satisfy this
  check on their own. Startup fails with a clear error if no
  self-signature is present.
- **Missing party.** If `party.env.json` is absent, the server fails
  to start with bootstrap instructions naming both the expected path
  and the `gobl sign` invocation needed to produce the file from a
  Party JSON.

**Optional port overrides:**

- `--http-port <int>` — HTTP listen port (default 80).
- `--https-port <int>` — HTTPS listen port (default 443). Only used when
  a TLS source is configured.

**TLS sources** (mutually exclusive):

- `--acme-live` — Activate HTTPS via Let's Encrypt (production directory).
  Requires `--domain`. Also RECOMMENDED: `--acme-email`.
- `--acme-test` — Same as `--acme-live` but uses Let's Encrypt's staging
  directory. Use this while iterating to avoid the production rate
  limits. Issued certificates are not trusted by browsers.
- `--tls-cert <path>` + `--tls-key <path>` — Operator-supplied PEM
  cert + key (corporate CA, externally-managed cert lifecycle,
  self-signed for testing).

**ACME-specific flags:**

- `--domain <fqdn>` — Hostname the ACME client is allowed to issue for.
  MUST match the participant's GOBL Net address. **Optional**: when
  omitted, it is derived from the party's GOBL inbox (see below).
- `--acme-email <email>` — Account email registered with the ACME
  directory (RECOMMENDED by LE). **Optional**: when omitted, it is
  derived from the party's first email address.
- `--cert-dir <path>` — Directory used to cache ACME-issued
  certificates. Default `<config-dir>/certs/`.

**Derivation from the party.** The served party envelope is the source
of truth for the participant's identity, so the ACME domain and account
email do not need to be passed as flags:

- **Domain** — taken from the party's GOBL inbox: an
  `org.Party.inboxes` entry with `key: "gobl"` whose `code` is the
  participant's GOBL Net address (FQDN). Example:

  ```json
  {"inboxes": [{"key": "gobl", "code": "samlown.example.com"}]}
  ```

  If neither `--domain` nor a GOBL inbox is present, ACME startup fails
  with a clear error.
- **Email** — taken from the party's first `org.Party.emails` entry.
  Email is optional for ACME, so an absent email is not an error.

An explicit `--domain` / `--acme-email` flag always overrides the
party-derived value.

**Three operational stances:**

| Stance                                          | Listens on             | Use when                                          |
|-------------------------------------------------|------------------------|---------------------------------------------------|
| default (no TLS flags)                          | HTTP only              | Behind a reverse proxy that terminates TLS upstream. |
| `--acme-live` / `--acme-test` + `--domain`      | HTTP + HTTPS           | Direct internet exposure; LE manages the cert.    |
| `--tls-cert` + `--tls-key`                      | HTTP + HTTPS           | Cert is sourced elsewhere (corporate CA, etc.).   |

**Docker deployment:**

Containers can bind privileged ports inside their own network namespace
without host-level privilege, so the simplest deployment is:

```
docker run \
    -p 80:80 -p 443:443 \
    -v gobl-config:/root/.config/gobl \
    gobl net serve
```

When the container itself runs unprivileged (e.g. `USER 1000` images),
pick high ports inside and remap externally:

```
docker run \
    -p 80:8080 -p 443:8443 \
    -v gobl-config:/home/gobl/.config/gobl \
    gobl net serve --http-port 8080 --https-port 8443
```

**ACME operational sequence:** start the server → ACME challenge
(HTTP-01 on the HTTP port, with TLS-ALPN-01 fallback on the HTTPS port)
→ cert issued and cached → ready. If the public internet cannot reach
the HTTP or HTTPS port for the configured `--domain`, the challenge
fails and the server logs a clear error. Successful issuance
double-checks that the deployment is actually reachable.

### 9.3 `gobl net send`

Reads a signed envelope from a file or stdin and POSTs it to the
destination's inbox endpoint:

- `-t, --to <fqdn>` — Destination Address.
- `--insecure` — Use `http://` and permit `host:port` form in `--to`
  (development only).

Exit status is 0 on a 202 response; any other status returns
`ErrInboxRejected`.

## 10. Errors

The package exports the following sentinel errors:

| Error                  | Cause                                                            |
|------------------------|------------------------------------------------------------------|
| `ErrAddressEmpty`      | Empty input to `ParseAddress`.                                   |
| `ErrAddressInvalid`    | Input is not a valid FQDN per §3.1.                              |
| `ErrFetchFailed`       | Well-known resource fetch failed (network, non-200, oversize, malformed). |
| `ErrKeyNotFound`       | The requested `kid` was not found in the KeySet.                 |
| `ErrNoGNHeader`        | Signature has no `gn` header and none was provided.              |
| `ErrVerifyFailed`      | Envelope verification failed.                                    |
| `ErrKeySetInvalid`     | KeySet is malformed.                                             |
| `ErrUnknownAuthority`  | A `/who` envelope is signed by an address not in `Authorities`.  |
| `ErrPartyMissing`      | A `/who` response did not contain an `org.Party` document.       |
| `ErrInboxRejected`     | A receiving inbox did not return 202.                            |

All callers using `errors.Is` against these sentinels MUST continue to
work after wrapping with `fmt.Errorf("%w: ...", err)`.

## 10a. Logging

All operator-facing log output goes through `log/slog` and is written
to **stderr**. Result output (signed envelopes, the `/who` party JSON,
`gobl version`'s JSON) stays on **stdout**, so a pipeline like
`gobl sign … | gobl net send …` is unaffected by logging.

The format is controlled by the top-level `--json` flag:

| flag         | stderr format        | example                                          |
|--------------|----------------------|--------------------------------------------------|
| (default)    | slog text            | `time=… level=INFO msg=listening scheme=http addr=:8080` |
| `--json`     | slog JSON-per-line   | `{"time":"…","level":"INFO","msg":"listening","scheme":"http","addr":":8080"}` |

### Startup messages (`gobl net serve`)

- `generated keypair` — emitted on first boot (auto-keygen). Fields:
  `kid`, `private`, `key_file`.
- `initialised domain` — emitted by `gobl init`. Fields: `domain`,
  `party`, `inbox`.
- `GOBL Net listening` — once per listener. Fields: `scheme`
  (`http`/`https`), `addr`.
- `ACME enabled` — when ACME is configured. Field: `domains`.
- `Shutting down` — emitted on graceful shutdown.

### Access logs

Every HTTP request emits one **baseline** entry plus one or more
**handler-specific** entries with high-signal fields:

| msg              | level | fields                                                                |
|------------------|-------|-----------------------------------------------------------------------|
| `http_request`   | INFO  | `method`, `path`, `host`, `remote`, `status`, `duration_ms`           |
| `keys.lookup`    | INFO  | `kid`, `found`                                                        |
| `who.exchange`   | INFO  | `caller` (verified `iss` as FQDN)                                     |
| `who.rejected`   | WARN  | `reason` (`bad_body`/`verify_failed`/`not_allowed`), `remote`/`caller`/`error` |
| `inbox.accepted` | INFO  | `caller`, `envelope` (UUID)                                           |
| `inbox.rejected` | WARN  | `reason` (`bad_body`/`validation`/`verify_failed`/`aud_mismatch`/`not_allowed`), `caller`/`aud`/`error` |
| `inbox.write_failed` | ERROR | `caller`, `envelope`, `error`                                       |
| `who.party_load_failed` / `who.sign_failed` / `who.encode_failed` | ERROR | `caller`, `error` |

### Error reporting

A CLI command that fails emits a single `command failed` entry on
stderr with `key=<gobl-error-key>` and (when present) `message=…` and
`faults=…`. With `--json` the same fields appear as a JSON object.
Successful commands write no log output and their result still lands
on stdout.

## 11. Security Considerations

### 11.1 Trust Model

The `gn` header is a discovery hint and carries no inherent trust. A
forged `gn` header pointing to an attacker-controlled host can only
produce a signature that verifies against an attacker-controlled key —
which is then distinguishable from the expected identity at the
application layer. Callers that already know the expected Address SHOULD
supply it explicitly to `VerifyEnvelope` rather than relying on the
embedded hint.

### 11.2 TLS

All well-known endpoints are served over HTTPS in production. Verifiers
rely on the host's TLS certificate to establish that the served content
originates from the named Address. Operators MUST ensure that TLS
certificate issuance is properly controlled for any Address they intend
to use as an identity.

When the server is started with `--acme-live` (or `--acme-test`), the
successful issuance of a Let's Encrypt certificate doubles as a
reachability check: it proves that the public internet can reach the
configured `--domain` on the HTTP/HTTPS ports. In the default plain-HTTP
mode no such guarantee exists, and operators MUST verify their upstream
TLS proxy independently.

### 11.3 Authority Trust

`Authorities` is a small allowlist of FQDNs. An attacker who gains the
ability to issue a TLS certificate for any address in `Authorities` can
serve a forged `/who` for any participant and pass `WhoIs` verification.
The KYC vendor list MUST therefore be kept short and reviewed regularly.

### 11.4 Inbox Authentication

The inbox endpoint verifies the sender's signature before persisting an
envelope, but does *not* require the sender to be a known KYC-endorsed
participant. Operators that need to restrict inbox acceptance to known
correspondents MUST apply additional filtering (e.g. wrap the
`Client.VerifyEnvelope` call with a `WhoIs` check against the sender's
address).

### 11.5 Response Size

The 1 MiB cap on Key, Who, and inbox bodies limits memory amplification
from hostile or misconfigured peers.

### 11.6 Address Canonicalization

`ParseAddress` lowercases and strips trailing dots before validation, so
two visually distinct strings such as `Example.COM.` and `example.com`
normalize to the same Address. Callers MUST use `ParseAddress` (directly
or via `Address.Validate`) before comparing addresses for equality.

## 12. References

- RFC 2119 — Key words for use in RFCs to Indicate Requirement Levels
- RFC 7515 — JSON Web Signature (JWS)
- RFC 7517 — JSON Web Key (JWK)
- RFC 7518 — JSON Web Algorithms (JWA)
- RFC 8174 — Ambiguity of Uppercase vs Lowercase in RFC 2119 Key Words
- `github.com/invopop/gobl/dsig` — Signature, key, and digest primitives
- `github.com/invopop/gobl/org` — `org.Party` schema
