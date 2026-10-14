## Removed

- **breaking**: `data`: the generated addon definitions under `data/addons` are
  gone, along with the generator that produced them. They could only ever
  describe the addons still compiled into core, so as addons moved to their own
  modules the directory became a shrinking subset that looked complete: of the
  25 addon keys GOBL recognises — 13 built in, 12 on the approved external
  list — only the built-in 13 had a file. Build the definition from the
  registry instead, which covers every addon a binary has loaded:
  `schema.NewObject(tax.AddonForKey(key))` marshalled with two-space indent
  reproduces exactly what these files contained.
- `data`: `data/regimes/gr.json` is gone. Greece registers under `EL`, so this
  was a leftover from before that rename, unreachable from the registry and
  last regenerated in 2024 while `el.json` kept moving. Nothing generates it.
