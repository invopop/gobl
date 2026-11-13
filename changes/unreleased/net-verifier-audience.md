## Changed

- `net`: **breaking**: verifier countersignatures count only when their `aud` is the envelope's `from` address, so each address needs its own verification, and they need a `cnf` naming the subject key the verifier saw.
- `net`: **breaking**: authority countersignatures need an `iat` no later than the present and an `exp` at most 90 days after it.
- `net`: **breaking**: `Client.VerifyDelivery` and `Client.Who` need an `iat` no later than the present, and a key that is still valid when the envelope is received.
- `net`: **breaking**: authority countersignatures need a `cnf` naming the key of a signature that confirms the envelope's `from` and `to`; `Endorsement.Key` reports it.
- `net`: **breaking**: `Client.VerifySender(ctx, env, requireVerified)` takes a delivered envelope, and requires it to be signed with the sender's registered key.
- `net`: `Client.DeliveryKeys` returns the keys confirming a delivery, and `Client.FindVerification` reads a verifier countersignature, for authorities tracking key lineage.
- `head`: `WithConfirmation` sets an RFC 7800 `cnf` claim, read from `SigningPayload.Confirmation`.
- `dsig`: `JKT` returns a key's RFC 7638 thumbprint in base64url form.
- `net`: only the first four signatures by one issuer are checked against its keys.

## Fixed

- `net`: `HTTPFetcher` dials the addresses it has checked, rather than resolving the host again.
