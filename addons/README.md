# Addons

Addons add normalization, validation rules, and scenarios that adapt a GOBL
document to a specific format or regulatory regime. Each addon is identified by
a versioned key (e.g. `mx-cfdi-v4`) that a document lists under `$addons`.

## In-core addons

Most addons live in this directory and register themselves automatically: the
blank imports in [`addons.go`](./addons.go) run each package's `init()`, which
calls `tax.RegisterAddonDef`. Importing `github.com/invopop/gobl` (which imports
this package) therefore makes every in-core addon available.

## External addons & the approval process

An addon may instead live in its **own Go module** — for example
[`github.com/invopop/gobl.fr.ctc`](https://github.com/invopop/gobl.fr.ctc). This
keeps large, country-specific rule sets out of core and lets a project opt in
only when it needs them. Mechanically an external addon is identical to an
in-core one: its `init()` calls `tax.RegisterAddonDef`, so a consumer just adds a
blank import:

```go
import _ "github.com/invopop/gobl.fr.ctc/addon"
```

### The runtime contract is strict

Listing a key under `$addons` requires the addon to be **actually loaded** at
`Validate`/`Calculate` time. If the module has not been imported, validation
fails with `add-on must be registered`. This is deliberate: a document is never
silently processed without the normalizers and rules its `$addons` promise. Any
service that handles documents for an external addon must import that module.

### The approved-addon registry

So that an external addon's key is still a recognised, schema-valid `$addons`
value even where the implementation is not compiled in, this package keeps a
curated list of **approved** external addons in
[`external.go`](./external.go), alongside the in-core addon imports. Entries on
this list:

- appear in the JSON Schema `$addons` enum (via `AddonList.JSONSchemaExtend`), and
- record the implementing module for provenance.

Being on the approved list is **recognition and governance only** — it does not
relax the strict runtime contract above.

### Adding an approved addon

Approval is a reviewed pull request that adds a `tax.RegisterApprovedAddon` entry
to [`external.go`](./external.go). A new entry should satisfy:

- the implementation is a public module under `github.com/invopop` that
  auto-registers via `init()` + `tax.RegisterAddonDef`;
- the key follows the `<addon>-vN` convention and does not collide with an
  in-core addon key;
- consumers that process documents declaring the key import the module, so the
  strict runtime check still passes.
