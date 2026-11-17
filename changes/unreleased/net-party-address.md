## Changed

- `net`: **breaking**: GOBL Net envelopes must carry header `from` and `to` fields holding `gobl:` addresses, and a valid signature whose `iss` is `from` and whose `aud` is `to`. Envelopes are no longer forwarded on another participant's behalf.
- `net`: **breaking**: `Client.VerifyParty` is removed. `Client.VerifyDelivery` checks any GOBL Net envelope against the `from` and `to` pair, and returns only an error; the sender is the `from` address.
- `net`: **breaking**: a party envelope's subject is its `from` address, instead of the party's `gobl:` endpoint, which is now optional and only a hint. A party can be verified before it has an address, and used at several addresses, each holding its own copy of the envelope.
- `net`: `Sign` sets an envelope's `from` and `to` and signs it for that pair.
- `net`: **breaking**: `Client.Send(ctx, env)` posts to the envelope's `to` address and sends no request token; inbox deliveries are authenticated by the envelope's signature alone. `Fetcher.Post` no longer takes a header.
- `net`: **breaking**: `Client.Who` requires a client identity, and accepts a response only from the queried address to that identity, signed for it by the subject.
- `net`: **breaking**: `Client.VerifyAuthority` only accepts authority countersignatures whose `aud` is the envelope's `from` address.
- `net`: verifier countersignatures need an `iat` no later than the present; their `aud` is ignored and their `exp` is optional.
