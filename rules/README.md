# `rules` — GOBL Validation Framework

The `rules` package provides structured validation for GOBL types. Rules produce
machine-readable fault codes (e.g. `GOBL-HEAD-HEADER-02`) alongside human-readable
messages, making errors testable by stable code rather than fragile string matching
and suitable for export as structured data.

## Core concepts

### `For` — define a rule set for a type

```go
func myRules() *rules.Set {
    return rules.For(new(MyStruct),
        // ... Defs
    )
}
```

`rules.For` accepts either a struct pointer or a named value type (e.g.
`rules.For(MyCode(""), ...)`). The prototype value is used for type inference
and to validate field names at initialisation time.

### `Field` — scope assertions to a field

```go
rules.Field("name",
    rules.Assert("01", "name is required", is.Present),
)
```

The name must match the JSON tag of a field in the parent struct. It is
validated at initialisation — an unknown name panics immediately. All
assertions inside `Field` receive the extracted field value.

### `Assert` — a single validation assertion

```go
rules.Assert("01", "description", test1, test2, ...)
```

All tests must pass. The first failure short-circuits the assertion and emits a
fault with the given description. Assertion codes are prefixed automatically
by `Register` or `NewSet` to form globally unique codes like `GOBL-ORG-EMAIL-01`.

Use `AssertIfPresent` when the assertion should be skipped for nil or empty values:

```go
rules.Field("code",
    rules.AssertIfPresent("01", "code format invalid",
        is.Func("valid", isValidCode),
    ),
)
```

### `Each` — per-element assertions on a slice field

`Each` is a nameless `Def` used inside `Field` to apply assertions to each
element of a slice. It does not take a field name — it operates on the slice
already extracted by the enclosing `Field`:

```go
rules.Field("lines",
    rules.Assert("01", "no duplicate line codes",
        is.Func("no duplicates", hasNoDuplicateLineCodes),
    ),
    rules.Each(
        rules.Field("code",
            rules.Assert("02", "line code is required", is.Present),
        ),
    ),
)
```

Faults carry indexed paths: `lines[0].code`, `lines[1].code`, etc. `Each` panics
at initialisation time when used outside a slice field. If the element type has
its own registered rule set, those rules are applied automatically during
recursive validation — `Each` is only needed when adding assertions from the
parent's perspective that don't belong on the element type itself.

### `When` — conditional rule subsets

```go
rules.When(conditionTest,
    rules.Field("code", rules.Assert("01", "code is required", is.Present)),
)
```

The subset is only evaluated when `conditionTest.Check(obj)` returns `true`.
The condition receives the full parent struct. Use `is.Expr(...)`,
`is.Func(...)`, or any `Test` implementation.

### `Object` — group object-level assertions

```go
rules.Object(
    rules.Assert("10", "cross-field constraint", is.Expr(`FieldA != "" || FieldB == nil`)),
)
```

`Object` is sugar for passing assertions directly to `For` or `When`. Use it
for organisational clarity when mixing field and object-level assertions.

### `Register` — add rules to the global registry

In the package `init()` function (typically `mypkg.go`):

```go
func init() {
    schema.Register(schema.GOBL.Add("mypkg"), MyStruct{})
    rules.Register(
        "mypkg",
        rules.GOBL.Add("MYPKG"),
        myStructRules(),
        anotherStructRules(),
    )
}
```

The first argument is the package name (used for generation), and the second is
the namespace code prepended to all assertion IDs. Rules registered here are
automatically applied by `rules.Validate(obj)` to any matching object type
anywhere in the object graph.

### `NewSet` — standalone namespace sets

When using rule sets outside the GOBL global registry (e.g. in a separate
application for validating request bodies), use `NewSet` to create a namespace
set that can be validated directly:

```go
const MyApp rules.Code = "MYAPP"

personSet := rules.For(new(Person),
    rules.Assert("01", "name required", is.Present),
)
emailSet := rules.For(new(Email),
    rules.Field("addr",
        rules.Assert("01", "email required", is.Present),
    ),
)

validator := rules.NewSet(MyApp, personSet, emailSet)
faults := validator.Validate(person)
// Codes: MYAPP-PERSON-01, MYAPP-EMAIL-01
```

Unlike `Register`, the returned set is NOT added to the global registry and
will not be applied by `rules.Validate(obj)`. Like `Register`, the namespace
code is prepended to all assertion and set IDs during construction.

Input sets are cloned internally, so the same output of `For` can safely be
reused across multiple `NewSet` or `Register` calls.

`Set.Validate` also accepts optional `WithContext` values to inject context
for context-aware guards:

```go
faults := validator.Validate(obj, func(rc *rules.Context) {
    rc.Set("country", "ES")
})
```

## Available tests

All tests live in the `github.com/invopop/gobl/rules/is` package. Import it alongside `rules`:

```go
import (
    "github.com/invopop/gobl/rules"
    "github.com/invopop/gobl/rules/is"
)
```

| Test                                                  | Notes                                                                  |
| ----------------------------------------------------- | ---------------------------------------------------------------------- |
| `is.Present`                                          | Fails if nil, zero, or empty                                           |
| `is.NilOrNotEmpty`                                    | Passes if nil pointer or non-empty                                     |
| `is.Empty`                                            | Passes if nil or empty; fails if a value is present                    |
| `is.Nil`                                              | Passes only for a nil pointer; fails for any non-nil value, even empty |
| `is.In(vals...)`                                      | Skips nil; works with named types                                      |
| `is.NotIn(vals...)`                                   | Skips nil; works with named types                                      |
| `is.Matches(pattern)`                                 | Skips nil/empty strings                                                |
| `is.Length(min, max)`                                 | `max=0` means no upper bound                                           |
| `is.RuneLength(min, max)`                             | Unicode-aware                                                          |
| `is.Min(v)` / `is.Max(v)`                             | int, uint, float, time                                                 |
| `is.Expr(expr)`                                       | CEL-like expression; fields accessed by Go field name                  |
| `is.Func(desc, func(any) bool)`                       | Custom boolean function                                                |
| `is.StringFunc(desc, func(string) bool)`              | Convenience for string-typed fields                                    |
| `is.FuncError(desc, func(any) error)`                 | Error message is discarded; use `desc`                                 |
| `is.FuncContext(desc, func(rules.Context, any) bool)` | Context-aware custom function                                          |
| `is.Or(tests...)`                                     | Passes if any test passes                                              |
| `is.InContext(test)`                                  | Passes when any context value satisfies the inner test                 |

The `rules/is` package also re-exports all format tests from `github.com/invopop/validation/is`
(e.g. `is.URL`, `is.EmailFormat`, `is.Alphanumeric`).

## Common patterns

### Required field with format check

Split presence and format into separate assertions so callers can distinguish a
missing value from a malformed one:

```go
rules.Field("addr",
    rules.Assert("01", "email address is required", is.Present),
    rules.Assert("02", "email address must be valid", is.EmailFormat),
)
```

### Allowed-values check

```go
rules.Field("category",
    rules.Assert("02", "category is not valid", is.In("a", "b", "c")),
)
```

`is.In` normalises named string/int types so `In("a", "b")` matches both
`string("a")` and `MyType("a")`. For non-pointer named types like `cbc.Key`
where `In` cannot distinguish absent from invalid, use `AssertIfPresent` with a
custom validator instead.

### Custom validation logic

Extract logic into named private functions and use `is.Func`, `is.StringFunc`,
or `is.FuncError`. **Prefer private named functions over inline anonymous
functions** — they are easier to test in isolation, appear in stack traces, and
keep the rule set readable at a glance.

```go
func myCodeChecksumValid(code string) bool {
    return isValidChecksum(code)
}

rules.Field("code",
    rules.Assert("03", "code checksum mismatch",
        is.StringFunc("checksum", myCodeChecksumValid),
    ),
)
```

### Object-level (cross-field) assertions

Without `Field`, an assertion receives the full object. Prefer `is.Expr` for
simple comparisons; use `is.Func` when the logic is more involved or when you
want a named, testable function:

```go
// Simple cross-field check
rules.Assert("10", "digest must be nil when MIME type is not provided",
    is.Expr(`MIME != "" || Digest == nil`),
)

// More complex logic
func digestHasMIME(val any) bool {
    obj, ok := val.(*MyStruct)
    if !ok || obj == nil {
        return false
    }
    return obj.MIME != "" || obj.Digest == nil
}

rules.Assert("10", "digest must be nil when MIME type is not provided",
    is.Func("no digest without MIME", digestHasMIME),
)
```

> **Note on receiver shape:** `rules.Validate(obj)` may pass either `*T` or `T`
> to an object-level `Func` depending on how the object was reached. Always
> handle both in object-level helpers. `Expr` handles this automatically.

### Conditional validation

```go
func envelopeNotSigned(val any) bool {
    e, ok := val.(*Envelope)
    return ok && len(e.Signatures) == 0
}

rules.When(is.Func("not signed", envelopeNotSigned),
    rules.Field("stamps",
        rules.Assert("12", "stamps not allowed before signing", is.Length(0, 0)),
    ),
)
```

Rules have no access to `context.Context`. Conditions that depend on runtime
context (e.g. "is signed?") must be derived from the object's own state.

### Nested struct fields

Define rules for each type independently and register them all (or group them
with `NewSet`). Both `rules.Validate` and `Set.Validate` recurse into every
exported field automatically — no wiring is needed between parent and child.

When you need to add constraints on a nested type from the **parent's
perspective** (e.g. regime-specific rules that don't belong on the child type),
nest `rules.Field` calls to drill down the path:

```go
func invoiceRules() *rules.Set {
    return rules.For(new(Invoice),
        rules.When(is.InContext(tax.RegimeIn("XX")),
            rules.Field("supplier",
                rules.Assert("01", "supplier is required", is.Present),
                rules.Field("tax_id",
                    rules.Assert("02", "supplier tax ID is required", is.Present),
                    rules.Field("code",
                        rules.Assert("03", "supplier tax ID must have a code", is.Present),
                    ),
                ),
            ),
        ),
    )
}
```

Each `rules.Field` in the chain constrains the context for its children, so
assertions inside `rules.Field("tax_id", ...)` operate on the `TaxIdentity`
struct, not the outer `Invoice`.

### Named value types (e.g. `cbc.Code`, `tax.Rate`)

`rules.For` works with named non-struct types:

```go
func myCodeRules() *rules.Set {
    return rules.For(MyCode(""),
        rules.Assert("01", "code must not be empty", is.Present),
        rules.Assert("02", "code must be alphanumeric", is.Alphanumeric),
    )
}
```

Inside `Expr`, the value is exposed as `this`:

```go
rules.Assert("02", "code must not exceed 10 characters",
    is.Expr(`len(this) <= 10`),
)
```

## Testing

Call `rules.Validate(obj)` (global registry) or `set.Validate(obj)` (standalone)
and assert on the returned `rules.Faults` value:

```go
import "github.com/invopop/gobl/rules"

// Global registry validation:
err := rules.Validate(obj)
assert.NoError(t, err)

faults := rules.Validate(obj)
require.NotNil(t, faults)
assert.True(t, faults.HasPath("field"))
assert.True(t, faults.HasCode("GOBL-PKG-STRUCT-01"))
assert.Equal(t, "assertion description", faults.First().Message())

// Standalone set validation:
faults = mySet.Validate(obj)
require.NotNil(t, faults)
assert.True(t, faults.HasCode("MYAPP-STRUCT-01"))
```

`rules.Faults` implements `error`. A nil return means no faults. The full error
string format is:

```
[GOBL-PKG-STRUCT-01] field: assertion description
```

## Assertion code conventions

Codes within a set are short local identifiers (e.g. `"01"`, `"02"`). They are
prefixed by `Register` or `NewSet` to form globally unique codes. Follow this
convention:

- `01`–`09`: field-level assertions, in the order fields appear in the struct
- `10`–`19`: object-level (cross-field) assertions
- `20`+: reserved for `When` conditional subsets if needed

The fully-qualified code is constructed as:
`{NAMESPACE}-{STRUCT}-{LOCAL_CODE}`

For example, a `"03"` assertion on `head.Header` registered under `GOBL-HEAD`
becomes `GOBL-HEAD-HEADER-03`.

## Assertion message conventions

Write messages so they are self-explanatory without inspecting the fault path or
source code:

1. **Include the parent context** for nested fields — write
   `"supplier tax ID is required"`, not `"tax ID is required"`.
2. **Include extension keys in single quotes** using `fmt.Sprintf` — write
   `fmt.Sprintf("tax requires '%s' extension", ExtKeyModel)`, not
   `"tax requires a model extension"`.
3. **Include extension values** when the code alone is ambiguous — write
   `fmt.Sprintf("NF-e does not support '%s' for '%s'", PresenceDelivery, ExtKeyPresence)`.
4. **Preserve business rule codes** (e.g. `BR-FR-30`) in messages when the
   original validation spec includes them — they are the primary reference for
   compliance.
