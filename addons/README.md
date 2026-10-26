# Addons

Addons add normalization, validation rules, and scenarios that adapt a GOBL
document to a specific format or regulatory regime. Each addon is identified by
a versioned key (e.g. `mx-cfdi-v4`) that a document lists under `$addons`.

## In-core addons

Most addons live in this directory and register themselves automatically: the
blank imports in [`addons.go`](./addons.go) run each package's `init()`, which
calls `tax.RegisterAddonDef`. Importing `github.com/invopop/gobl` (which imports
this package) therefore makes every in-core addon available.

## Extensions and stamps

Addons describe the extensions and stamps they offer with `cbc.Definition`s so
that tools, such as document editors, can list and label the fields available
for a document. Each definition must have:

- **A name prefixed with the issuer** in every language, such as
  `"SII: Invoice Type"` or `"TicketBAI: QR Code"`, so the field can be told apart
  from similar fields offered by other addons. Don't repeat the issuer in the
  rest of the name.
- **A `Usage`** describing where the definition applies:
  - `Keys` with the addon's key, such as `[]cbc.Key{V1}`. Tax regimes use
    `Codes` with their country code instead.
  - `Schemas` (extensions only) with the schema IDs, relative to the GOBL base,
    of the objects whose `ext` holds the extension, such as `tax/combo`,
    `bill/tax`, or `org/party`. Each must be a registered schema that has
    extensions.
  - `Paths` (optional) with JSON paths, only when the extension applies to some
    of a schema's locations. A path also covers everything under it, so
    `$.supplier` limits an `org/party` extension to the supplier, and
    `$.lines[*].taxes[*]` limits a `tax/combo` extension to line taxes.

```go
{
	Key: ExtKeyBIActivity,
	Name: i18n.String{
		i18n.EN: "TicketBAI: Activity Code (Bizkaia)",
		i18n.ES: "TicketBAI: Código de Actividad (Bizkaia)",
	},
	Usage: &cbc.Usage{
		Keys:    []cbc.Key{V1},
		Schemas: []string{"org/party"},
		Paths:   []string{"$.supplier"},
	},
	// ...
}
```

Extensions and stamps are listed in the `AddonDef`'s `Extensions` and `Stamps`,
and `tax.RegisterAddonDef` registers them in the `tax` and `head` packages.
Registering a different definition with a key already in use panics, so keys
must be unique across all regimes, addons, and catalogues.

`Usage` is descriptive: it helps tools filter and present definitions, but it is
not used to validate documents. Rules and normalizers still decide where an
extension is required.

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

### Expectations of an external addon

An external addon is held to the same standards as an in-core one, plus the
requirements that come with living in its own repository. Before its key is
approved, the module must satisfy all of the following.

**The standard addon requirements** — everything expected of any GOBL addon:

- a clear scope and a versioned key (`<addon>-vN`) that does not collide with an
  in-core addon key;
- extensions, normalizers, scenarios, and validation rules defined in idiomatic
  GOBL style, with stable fault codes;
- extension and stamp definitions with prefixed names and a `Usage`, as
  described in [Extensions and stamps](#extensions-and-stamps);
- self-registration via `init()` + `tax.RegisterAddonDef` so a blank import is
  all a consumer needs;
- a `README` documenting the addon's purpose, the flows/documents it covers, and
  how to use it.

**Additional requirements for an external module:**

- **Public module** under `github.com/invopop`, importable on its own.
- **Automated linting and testing are mandatory.** The repository must run
  `golangci-lint` and `go test` in CI (GitHub Actions) on every pull request,
  mirroring this repository's `lint` and `test` workflows, and the default
  branch must be green.
- **Status badges must be present** in the module `README` — at least build
  (lint/test), code coverage, latest release tag, and Go reference — so its
  health is visible at a glance.
- **At least 90% test coverage** is required, reported through the CI coverage
  pipeline (e.g. Codecov) and kept at or above that threshold.
- **Example documents** for each supported flow/document, verified with the
  shared [`pkg/examples`](../pkg/examples) helpers (`examples.Run`), so behaviour
  is exercised end to end and stays stable.

**Approval** is then a reviewed pull request that adds a
`tax.RegisterApprovedAddon` entry to [`external.go`](./external.go), once the
above are met and consumers that process documents declaring the key import the
module (so the strict runtime check still passes).

[`github.com/invopop/gobl.fr.ctc`](https://github.com/invopop/gobl.fr.ctc) is the
reference implementation of an external addon following these expectations.
