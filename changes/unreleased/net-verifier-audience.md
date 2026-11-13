## Changed

- `net`: **breaking**: verifier countersignatures count only when their `aud` is the envelope's `from` address, so each address needs its own verification.
- `net`: **breaking**: authority countersignatures need an `iat` no later than the present and an `exp` at most 90 days after it.
- `net`: **breaking**: `Client.VerifyDelivery` and `Client.Who` need an `iat` no later than the present, and a key that is still valid when the envelope is received.
- `net`: only the first four signatures by one issuer are checked against its keys.

## Fixed

- `net`: `HTTPFetcher` dials the addresses it has checked, rather than resolving the host again.
