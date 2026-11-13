## Changed

- `net`: **breaking**: GOBL Net envelopes must carry header `from` and `to` fields holding `gobl:` addresses, and a valid signature whose `iss` is `from` and whose `aud` is `to`. Envelopes are no longer forwarded on another participant's behalf.
- `net`: **breaking**: `Client.VerifyParty` is removed. `Client.VerifyDelivery` checks any GOBL Net envelope against the `from` and `to` pair, and returns only an error; the sender is the `from` address.
- `net`: **breaking**: a party envelope's subject is its `from` address, instead of the party's `gobl:` endpoint, which is now optional and only a hint. One party envelope can be verified before it has an address and published at several.
- `net`: **breaking**: `Client.Who` requires a client identity, and accepts a response only from the queried address to that identity, signed for it by the subject.
- `net`: **breaking**: `Client.VerifyAuthority` only accepts authority countersignatures whose `aud` is the envelope's `from` address.
- `net`: verifier countersignatures need an `iat` no later than the present; their `aud` is ignored and their `exp` is optional.
