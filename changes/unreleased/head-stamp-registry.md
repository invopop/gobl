## Added

- `head`: stamp registry. Packages call `RegisterStamps` with the definitions of the stamp provider keys they produce; `StampDefForKey` and `AllStampDefs` look them up. Registering a different definition for an existing key panics.
- `tax`: `Stamps` field on `RegimeDef` and `AddonDef` listing the keys of the stamps they produce.
- `data/stamps.json`: generated list of all registered stamp definitions.
- `regimes/mx`, `regimes/pt`, `regimes/gr`, `es-tbai-v1`, `it-ticket-v1`: register definitions for their existing stamp provider keys and list them in their definitions.
