## Changed

- `net`: **breaking**: verifier and authority countersignatures count only when their `sub` is the envelope's `from` address, so each address needs its own verification. Their `aud` is ignored.
- `net`: **breaking**: verifier countersignatures need an `iat` no later than the present and an `exp` at most a year after it.
- `net`: **breaking**: a signature with a `sub` does not confirm an envelope's `from` and `to`.
- `head`: `WithSubject` sets the `sub` claim of a signature, read from `SigningPayload.Sub`.
- `net`: **breaking**: authority countersignatures need an `iat` no later than the present and an `exp` at most 30 days after it.
- `net`: **breaking**: `Client.VerifyDelivery` and `Client.Who` need an `iat` no later than the present, and a key that is still valid when the envelope is received.
- `net`: **breaking**: `Client.VerifySender(ctx, env)` takes a delivered envelope and requires a verified sender.
- `net`: **breaking**: `Client.VerifyAuthority` is replaced by `Client.VerifyEndorsement`. Authority countersignatures must name a verifier whose countersignature is valid, and there is no registered-only status; `Endorsement.Verified`, `ErrUnknownAuthority` and `ErrNotVerified` are removed in favour of `ErrNotEndorsed`.
- `net`: `WithVerifiers` sets verifiers a client accepts directly, without an authority. There are none by default.
- `net`: `Client.FindVerifications` reads the valid countersignatures of a verifier, for authorities registering a subject.
- `net`: only the first four signatures by one issuer are checked against its keys.
- `net`: countersignatures whose `exp` has passed are skipped before their key is fetched, and do not count towards the four checked per issuer.
- `net`: `Endorsement.VerifiedAt` reports the signing time of the verifier's countersignature.
- `net`: `Client.FetchKey` rejects a `kid` that is not 1 to 64 characters of `A-Z`, `a-z`, `0-9`, `-` and `_`, without fetching it.
- `net`: a countersignature with a `verifier` claim never counts as a verification.

## Fixed

- `net`: `HTTPFetcher` dials the addresses it has checked, rather than resolving the host again.
- `net`: `HTTPFetcher` dial and lookup errors wrap their cause as well as `ErrUnavailable`.
- `net`: a delivery with no valid signature reports why the last candidate failed.
- `net`: `HTTPFetcher` refuses every IANA special-purpose range that is not globally reachable, and the NAT64, 6to4 and Teredo ranges.
- `net`: `HTTPFetcher` follows no redirects, and reports one as a definitive `ErrFetchFailed`.
- `net`: `Client.VerifyDelivery`, `Client.Who`, `Client.VerifyEndorsement` and `Client.FindVerifications` reject an envelope with no document, or whose document does not match its digest.
- `net`: `Client.VerifyEndorsement` and `Client.FindVerifications` report an unreachable verifier key over an expired or invalid verification.
- `net`: deliveries signed as a key becomes valid are accepted from a sender whose clock is up to 30 seconds ahead.
- `net`: one `Client.VerifyEndorsement` call checks at most four signatures per issuer.
- `net`: `HTTPFetcher` tries each resolved address but the last for 2 seconds before the next.
- `net`: `Client.VerifyDelivery`, `Client.Who` and `Client.VerifySender` reject a signature confirming `from` and `to` whose `exp` has passed.
