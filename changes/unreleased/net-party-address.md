## Changed

- `net`: **breaking**: a party envelope's subject is the `gobl:` address in the envelope header's `from`, confirmed by a self-signature, instead of the party's `gobl:` endpoint, which is now optional and only a hint. `from` is not signed, so one party envelope can be verified before it has an address and published at several.
- `net`: **breaking**: `Client.VerifyAuthority` only accepts authority countersignatures whose `aud` is the envelope's `from` address.
- `net`: **breaking**: envelopes must carry header `from` and `to` fields holding `gobl:` addresses, and are rejected otherwise. `Client.VerifyDelivery` also rejects a `from` that is not the sender, or a `to` that is not the receiving inbox.
- `net`: **breaking**: `Client.Who` requires a client identity and rejects a response whose `to` is not that identity.
- `net`: verifier countersignatures need an `iat` no later than the present; their `aud` is ignored and their `exp` is optional.
