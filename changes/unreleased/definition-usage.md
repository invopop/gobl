## Added

- `cbc`: `Usage` on `Definition`, with keys, codes, and JSON paths that
  describe where a definition may be used, such as the addons or regimes that
  offer it.
- `cbc`: `DefinitionRegistry`, an ordered set of definitions by key that can
  be filtered by usage.
- `head`: stamp definitions register, with `RegisterStampDef`, `StampDefFor`,
  `AllStampDefs`, and `FilterStampDefs`. Stamp names are prefixed with the
  issuer, such as "TicketBAI: QR Code".
- `tax`: `AllExtensionDefs` and `FilterExtensionDefs`.
- `pt`, `mx`, `el`: stamp definitions for the stamps issued by the AT, SAT,
  and IAPR.
- `es-tbai-v1`, `it-ticket-v1`: stamp definitions for the TicketBAI and AdE
  stamps.

## Changed

- `tax`: **breaking** registering a different extension definition with a key
  that is already in use panics.
