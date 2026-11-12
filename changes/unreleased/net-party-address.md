## Changed

- `net`: **breaking**: a party envelope's address comes from its signatures, not from the party's `gobl:` endpoint, which is now optional and only a hint. `Client.VerifyParty` takes the subject address and returns only an error, so one party envelope can be verified before it has an address and used at several.
- `net`: **breaking**: `Client.VerifyAuthority` takes the subject address and only accepts authority countersignatures whose `aud` is that address.
- `net`: verifier countersignatures need an `iat` no later than the present; their `aud` is ignored and their `exp` is optional.
