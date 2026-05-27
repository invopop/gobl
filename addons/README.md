# `addons` — Writing GOBL Addons

Addons layer jurisdiction- or platform-specific behaviour on top of a regime
(or independently of one). They register normalizers and validation rules
against the GOBL data model without changing it — addons add **extensions**
and **rules**, never new fields.

This guide captures the conventions that make addons easy to write, review,
and reason about. The canonical example referenced throughout is
[`addons/fr/ctc/flow6/`](fr/ctc/flow6/).

> **Project-wide guidance** lives in [`CONTRIBUTING.md`](../CONTRIBUTING.md).
> **Validation framework** is documented in [`rules/README.md`](../rules/README.md).
> This document covers the layer above both: how an addon is structured.

## Table of contents

- [Packaging and scope](#packaging-and-scope)
- [Registration](#registration)
- [Extensions](#extensions)
- [Normalization](#normalization)
- [Validation rules](#validation-rules)
- [Rule numbering](#rule-numbering)
- [Assert messages](#assert-messages)
- [Tests](#tests)
- [Generated data files](#generated-data-files)
- [Example documents](#example-documents)
- [Anti-patterns](#anti-patterns)

## Packaging and scope

**One addon = one Go package.** It self-registers via `init()` and exports a
single `Key` and `Namespace` constant pair.

**Split by responsibility, not by document type.** If a regulation has
several independent message flows (e.g. clearance vs. lifecycle vs.
e-reporting), give each flow its own sub-package + addon key. Compose them
with a meta-addon that auto-dispatches based on document content:

```
addons/fr/ctc/
  ctc.go              ← meta-addon: dispatchPayment / dispatchInvoice
  flow2/              ← B2B clearance
  flow6/              ← lifecycle status (CDV)
  flow10/             ← e-reporting
```

The meta-addon uses `tax.ExtractNormalizersForNew` to loop until the addon
set is stable.

**Shared code lists go in `catalogues/<authority>/`.** When several addons
(or an addon + a regime) reference the same external code list, lift it to
a catalogue package — `catalogues/iso`, `catalogues/untdid`,
`catalogues/dgfip`, etc.

**Don't add fields to `bill.*`, `org.*`, `pay.*` for an addon.** If you
need a new dimension on the model, it's either an extension on
`X.Ext` or a discussion about whether the core model needs to evolve.

## Registration

An addon registers in `init()`:

```go
func init() {
    tax.RegisterAddonDef(newV1Addon())
    rules.RegisterWithGuard(
        Key.String(),
        rules.GOBL.Add(Namespace),
        is.InContext(tax.AddonIn(V1)),
        billStatusRules(),
        billPaymentRules(),
        billReasonRules(),
        billActionRules(),
        orgPartyRules(),
    )
}
```

The `AddonDef` declares `Key`, EN+local-language `Name`/`Description`,
`Sources` (links to the underlying spec), `Extensions`, and a single
`Normalizer` function.

Blank-imported aggregation: list the addon package in
[`addons/addons.go`](addons.go).

## Extensions

### One file, declared explicitly

Define every extension key + values inline in `extensions.go`. Don't
generate the Values list programmatically — readers need to see the codes
and their labels in source.

```go
const (
    ExtKeyRole      cbc.Key = "fr-ctc-flow6-role"
    ExtKeyReason    cbc.Key = "fr-ctc-flow6-reason"
    ExtKeyStatus    cbc.Key = "fr-ctc-flow6-status"
    ExtKeyCondition cbc.Key = "fr-ctc-flow6-condition"
    ExtKeyAction    cbc.Key = "fr-ctc-flow6-action"
)

var extensions = []*cbc.Definition{
    {
        Key:  ExtKeyAction,
        Name: i18n.String{i18n.EN: "Requested Action", i18n.FR: "Action demandée"},
        Desc: i18n.String{
            i18n.EN: here.Doc(`CDAR RequestedActionCode (MDT-121) — what the
                upstream party is asking the recipient to do.`),
        },
        Values: []*cbc.Definition{
            {Code: "NOA", Name: i18n.String{i18n.EN: "No action",   i18n.FR: "Aucune action"}},
            {Code: "PIN", Name: i18n.String{i18n.EN: "Provide info", i18n.FR: "Fournir des informations"}},
            // ...
        },
    },
}
```

Every extension key carries a `Desc` that names **which document type(s)
it applies to** — `bill.Status`, `bill.Payment`, etc. This is the user-
facing answer to "where does this key go".

### Per-document-type allow-lists

When one extension applies to multiple document types with different valid
codes, register the full code list in the extension definition but enforce
the per-doc partition at the rules layer with `tax.ExtensionsHasCodes`:

```go
// In extensions.go (or a sibling file):
var (
    statusProcessCodes  = []cbc.Code{"200", "201", /* ... */ "210", "213"}
    paymentProcessCodes = []cbc.Code{"211", "212"}
)

// In billStatusRules():
rules.Field("ext",
    rules.Assert("02", "status ext fr-ctc-flow6-status must be a Status-applicable …",
        tax.ExtensionsHasCodes(ExtKeyStatus, statusProcessCodes...),
    ),
),
```

## Normalization

### One entry point, dispatch on type

```go
func normalize(doc any) {
    switch obj := doc.(type) {
    case *bill.Status:   normalizeStatus(obj)
    case *bill.Payment:  normalizePayment(obj)
    case *bill.Reason:   normalizeReason(obj)
    case *bill.Action:   normalizeAction(obj)
    case *org.Party:     normalizeParty(obj)
    case *org.Identity:  normalizeIdentity(obj)
    }
}
```

Each `normalizeX` handles one document type. Keep them in the file named
after the type (`bill_status.go`, `bill_payment.go`, …).

### Forward + reverse mapping

When an extension is logically equivalent to a `Key` field on the document,
support **both directions**:

- **Forward**: caller sets `Key`; normalizer fills the ext.
- **Reverse**: caller sets only the ext (round-tripping parsed data);
  normalizer fills `Key`.

The canonical shape is `normalizeReason` / `prepareReasonKey`:

```go
func normalizeReason(r *bill.Reason) {
    if r == nil { return }
    prepareReasonKey(r) // reverse step first
    switch r.Key {      // then forward
    case bill.ReasonKeyFinanceTerms:
        r.Ext = r.Ext.
            SetOneOf(ExtKeyReason, "COORD_BANC_ERR").
            SetOneOf(ExtKeyCondition, ConditionBankDetailsUpdate, /* ... */)
    // ...
    }
}

func prepareReasonKey(r *bill.Reason) {
    if !r.Key.IsEmpty() { return }
    switch r.Ext.Get(ExtKeyReason) {
    case "COORD_BANC_ERR": r.Key = bill.ReasonKeyFinanceTerms
    // ...
    }
}
```

### `Set` vs `SetOneOf` vs `SetIfEmpty`

The `tax.Extensions` API is immutable; each call returns a new value.

| Method | Behaviour |
|---|---|
| `Set(key, code)` | Unconditional overwrite. Use when the normalizer is the authoritative source. |
| `SetIfEmpty(key, code)` | Only sets if missing. Use when callers may pin a value. |
| `SetOneOf(key, default, alternatives...)` | If caller's value matches default or any alternative, keeps it; otherwise sets to default. Use when several valid values are acceptable and a sensible default exists. |

If you use unconditional `Set` for an ext that the normalizer derives from
another field, the corresponding rule-layer consistency check becomes
redundant — the normalizer guarantees it. (`normalizeStatusLine` does this:
the ext is always overwritten from the (Type, Key) pair, so no separate
rule checks consistency.)

### Inline the switch — no `XxxCodeFor` helpers

Don't extract a `func XxxCodeFor(key cbc.Key) (string, bool)` helper that
the normalizer then calls. The switch IS the mapping; the helper just
duplicates it. Inline it directly in `normalizeX`.

If a converter package outside the addon needs the same mapping, it
should call `Calculate()` on the document and read the resulting ext —
not import a helper.

## Validation rules

The rules framework is documented in [`rules/README.md`](../rules/README.md);
this section covers conventions specific to addons.

### Compose with built-in testers

Prefer composed `rules.When` / `rules.Field` / `rules.Each` with
built-in testers over wrapping logic in `is.Func`:

| Use | Built-in tester |
|---|---|
| Field is set | `is.Present` |
| Field has one of N values | `is.In(v1, v2, ...)` |
| Slice length | `is.Length(min, max)` |
| String matches regex | `is.MatchesRegexp(re)` |
| Ext is set | `tax.ExtensionsRequire(key)` |
| Ext value is in the extension's registered Values | `tax.ExtensionHasValidCode(key)` |
| Ext value is in an explicit subset | `tax.ExtensionsHasCodes(key, codes...)` |
| Ext value is NOT in a set | `tax.ExtensionsExcludeCodes(key, codes...)` |
| Identity slice has one with ext-value X | `org.IdentitiesExtensionIn(key, codes...)` |

Reach for `is.Func` only when the predicate genuinely needs imperative
logic that no built-in expresses. Even then, the actual assertion should
usually still be a built-in inside a `rules.When` — see "Per-process-code
allow-lists" below.

### Per-document-type guards: typed helpers, not `is.Expr`

Don't gate rules on the document's type or key with `is.Expr` string
expressions. Write a typed helper in the model package and use it as the
guard:

```go
// In bill/payment.go:
func PaymentTypeIn(types ...cbc.Key) rules.Test {
    return is.Func(
        fmt.Sprintf("payment type in [%s]", strings.Join(cbc.KeyStrings(types), ", ")),
        func(obj any) bool {
            pmt, ok := obj.(*Payment)
            return ok && pmt != nil && pmt.Type.In(types...)
        },
    )
}

// In the addon:
rules.When(
    bill.PaymentTypeIn(bill.PaymentTypeAdvice),
    rules.Field("ext",
        rules.Assert("15", "payment ext fr-ctc-flow6-status for an advice payment must be 211 (Paiement transmis)",
            tax.ExtensionsHasCodes(ExtKeyStatus, "211"),
        ),
    ),
),
```

The bill package already exposes `bill.InvoiceTypeIn`, `bill.PaymentTypeIn`,
`bill.StatusTypeIn`, `bill.StatusLineKeyIn`. Add more in that package
(not in your addon) when you need them.

### Gating predicates: tiny `is.Func` returning a `rules.Test`

When the gate genuinely depends on a value that no built-in can read
(e.g. an extension on a per-row basis inside `rules.Each`), write a small
helper that returns a `rules.Test`:

```go
func lineHasStatusCode(code cbc.Code) rules.Test {
    return is.Func(fmt.Sprintf("line status code %s", code), func(v any) bool {
        line, ok := v.(*bill.StatusLine)
        return ok && line != nil && line.Ext.Get(ExtKeyStatus) == code
    })
}
```

Then the actual assertion still leans on a built-in:

```go
rules.When(
    lineHasStatusCode("200"),
    rules.Field("reasons",
        rules.Each(
            rules.Field("ext",
                rules.Assert("17", "status line reason ext fr-ctc-flow6-reason for status code 200 (…) must be NON_TRANSMISE (BR-FR-CDV-CL-09)",
                    tax.ExtensionsHasCodes(ExtKeyReason, "NON_TRANSMISE"),
                ),
            ),
        ),
    ),
),
```

## Rule numbering

Rule codes are part of the addon's **public API**. The full code
(e.g. `GOBL-FR-CTC-FLOW6-BILL-STATUS-13`) appears in validation errors
returned to consumers, who route, log, alert, and write tests against
those codes. Stability matters.

**Conventions**:

- **Sequential** `01..NN` within each `rules.For` base object.
- **No letters** — no `19a`, `19b`, etc. When a rule needs to be split
  into per-value branches (e.g. one assertion per CDAR code), each
  branch gets its own sequential number.
- **Document order during initial development** — while the addon is
  still pre-release, number rules in the order they appear in the
  function body and freely renumber as you reorganise.

### Once an addon is released, numbers are immutable

After the first tagged release that includes the addon, rule numbers
become a **frozen contract**. Downstream consumers may key error
handling, monitoring, or test fixtures on the exact code.

- **Adding a rule**: append the next free number. The new rule lives at
  the end of the function body even if logically it belongs earlier —
  source order no longer dictates numbering once numbers are frozen.
- **Removing a rule**: leave a gap. Don't reuse the number for something
  else; that would silently change the meaning of the same identifier
  for any consumer who already handles it.
- **Replacing a rule's meaning**: don't. Add a new rule with the next
  free number and (if needed) remove the old one. The old number stays
  retired.
- **Splitting a rule** into per-value branches after release: the
  original number is retired; each branch gets a fresh next-free number.
- **Reordering**: only the source order changes; the numbers stay with
  their original assertions.

In other words: after release, treat rule numbers like database
primary keys — append-only, never recycled, never re-pointed.

### When in doubt about "released"

If you're not sure whether an addon has shipped, check `CHANGELOG.md`
for the addon's introduction and whether a tag has been cut since. If
yes, assume frozen numbering. If the addon is brand-new on a feature
branch that hasn't merged or shipped, you can still renumber.

## Assert messages

Format: **`<context/object(s)> <field> <constraint> (<BR-code>)`**

Examples:

| Path | Message |
|---|---|
| `$.type` | `"status type must be one of: response, update"` |
| `$.supplier.ext` | `"status supplier ext fr-ctc-flow6-role is required (BR-FR-CDV-CL-03)"` |
| `$.lines[0].doc.code` | `"status line doc code is required (BR-FR-CDV-10)"` |
| `$.lines[0].reasons[0].ext` | `"status line reason ext fr-ctc-flow6-reason for status code 200 (Déposée …) must be NON_TRANSMISE (BR-FR-CDV-CL-09)"` |

**Conventions**:

- **Lead with the document context** (`status`, `payment`, `reason`,
  `action`, `party`) so the message is decipherable when read in
  isolation — logs, an API error blob, a Slack alert. Don't rely on the
  reader having the JSON path in front of them.
- **Trailing parens for regulatory code**: `(BR-FR-CDV-13)`, `(MDT-121)`,
  `(BR-FR-CO-10)`. Round parens, at the end of the message, single space
  before. Don't prefix the BR code into the body (no
  `"BR-FR-CDV-CL-09: code 200 …"`).
- **Use the actual ext key** (`fr-ctc-flow6-status`), not a paraphrase
  (`status code`).
- **Don't repeat the JSON path** — the rules framework already prefixes
  it (`($.lines[0].doc.code)`).

## Tests

Use **testify** (`assert`, `require`).

```go
require.NoError(t, st.Calculate())
err := rules.Validate(st)
assert.ErrorContains(t, err, "status line key must be a recognised Flow 6 event")
```

**Test against the canonical message phrase**, not whatever substring
happens to appear today. Better still: assert against the **rule code**
(`GOBL-FR-CTC-FLOW6-BILL-STATUS-13`) — it's stable across message
rewordings.

A test fixture (`testStatus(t)`, `testPaymentReceipt(t)`) at the top of
each test file produces a known-valid document; individual tests mutate
one field and assert the targeted rule fires.

Defensive tests for nil-receiver / wrong-type branches of `is.Func`
predicates are fine but optional — they're symptoms of imperative
predicates, not a goal in themselves.

## Generated data files

Two generators regenerate JSON snapshots of the addon and its rules:

```bash
go run addons/generate.go    # → data/addons/<key>.json
go run rules/generate.go     # → data/rules/<key>.json
```

Run **both** after any change to an extension definition, an addon
description, or a rule. Commit the regenerated files alongside the source
change.

## Example documents

End-to-end examples live in `examples/<country>/*.yaml` (or `.json`) and
get regenerated into `examples/<country>/out/*.json`:

```bash
go test . -args --update     # regenerate all example outputs
```

When a normalizer change affects an example's serialized output (different
ext keys, reordered roles, etc.), regenerate and commit. When a snapshot
regenerates for reasons unrelated to your change, **revert that snapshot**
— it belongs to another in-flight refactor.

## Anti-patterns

Patterns to avoid, all encountered (and fixed) in the codebase:

- **`XxxCodeFor` / `YyyKeyFor` helpers** that duplicate a switch the
  normalizer would otherwise own. Inline the switch.
- **`is.Func` predicates that wrap a built-in** (e.g. iterating a slice
  to call `tax.ExtensionsHasCodes`). Express the iteration with
  `rules.Each` and let the built-in do the leaf check.
- **`is.Expr` string expressions for typed guards** (`string(Type) == "advice"`).
  Add a typed helper to the model package instead (`bill.PaymentTypeIn`).
- **Rule-layer consistency checks the normalizer already guarantees.**
  If `normalizeX` unconditionally `Set`s an ext from another field,
  delete the rule that re-checks them — it's unreachable.
- **Rule letters or placeholder numbers** (`19a`, `19b`, `0X`).
  Renumber to a clean sequence — but only if the addon hasn't shipped
  yet. After release, numbers are frozen (see [Rule numbering](#rule-numbering)).
- **BR codes inlined into the message body** (`"BR-FR-CDV-CL-09: code …"`).
  Move them to trailing parens.
- **Messages that hide the document context.** A reader of `("status
  type must be one of …")` doesn't need the path to know we're talking
  about a `bill.Status`; a reader of `("type must be one of …")` does.
- **Adding new `bill.ReasonKey` / `bill.ActionKey` constants for
  addon-specific semantics.** Use the existing Peppol-aligned vocabulary
  and layer your CDAR/CFDI/SDI/etc. codes as extensions on `X.Ext`.
- **Adding fields to `bill.*` / `org.*` / `pay.*` for one addon.** If
  the core model genuinely needs a new dimension, that's a separate
  conversation — not addon work.
