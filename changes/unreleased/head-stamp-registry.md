## Added

- `head`: stamp registry. Regimes call `RegisterRegimeStamps` with their code and addons call `RegisterAddonStamps` with their key; `StampDefForKey` and `AllStampDefs` return `StampDef`s whose `regime` or `addon` field matches a document's `$regime` or `$addons`. Registering a stamp key from a second regime or addon panics.
- `data/stamps.json`: generated list of all registered stamps with their regime or addon.
- `regimes/mx`, `regimes/pt`, `regimes/gr`, `es-tbai-v1`, `it-ticket-v1`: register definitions for their existing stamp provider keys.
