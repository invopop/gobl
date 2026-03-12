# Migrating from `validation` to `rules`

This document describes how to migrate packages from the
`github.com/invopop/validation` library to the `github.com/invopop/gobl/rules`
framework.

## Why migrate?

The `rules` framework produces machine-readable fault codes (e.g.
`GOBL-HEAD-HEADER-02`) in addition to human-readable messages. This makes
errors easier to handle programmatically, testable by stable code rather than
fragile string matching, and suitable for export as structured data alongside
the schemas they validate.

## Key concepts

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
    rules.Assert("01", "name is required", rules.Required),
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
during `Register` to form globally unique codes like `GOBL-ORG-EMAIL-01`.

### `Each` — per-element assertions on a slice field

`Each` is a nameless `Def` used inside `Field` to apply assertions to each
element of a slice. It does not take a field name — it operates on the slice
already extracted by the enclosing `Field`:

```go
rules.Field("lines",
    rules.Assert("01", "no duplicate codes", rules.By("no dups", hasNoDuplicateCodes)),
    rules.Each(
        rules.Field("code",
            rules.Assert("02", "line code is required", rules.Required),
        ),
    ),
)
```

Faults carry indexed paths: `lines[0].code`, `lines[1].code`, etc. Using `Each`
directly inside `For` (rather than inside a `Field` on a slice) panics at init time.

### `When` — conditional rule subsets

```go
rules.When(conditionTest,
    rules.Field("code", rules.Assert("01", "code is required", rules.Required)),
)
```

The subset is only evaluated when `conditionTest.Check(obj)` returns `true`.
The condition receives the full parent struct. Use `rules.Expr(...)`,
`rules.By(...)`, or any `Test` implementation.

### `Object` — group object-level assertions

```go
rules.Object(
    rules.Assert("10", "cross-field constraint", rules.Expr(`field_a != "" || field_b == nil`)),
)
```

`Object` is sugar for passing assertions directly to `For` or `When`. Use it
for organisational clarity when mixing field and object-level assertions.

### `Register` — add rules to the global registry

In the package `init()` function (typically `mypkg.go`):

```go
rules.Register(
    "mypkg",           // human-readable namespace name
    rules.GOBL.Add("MYPKG"), // code prefix, e.g. "GOBL-MYPKG"
    myStructRules(),
    anotherStructRules(),
)
```

Rules registered here are automatically applied by `rules.Validate(obj)` to
any matching object type anywhere in the object graph.

## Available tests

| Test | Replaces | Notes |
|------|----------|-------|
| `rules.Required` | `validation.Required` | Fails if nil, zero, or empty |
| `rules.NilOrNotEmpty` | `validation.NilOrNotEmpty` | Passes if nil pointer or non-empty |
| `rules.Empty` | `validation.Empty` | Passes if nil or empty; fails if a value is present |
| `rules.Nil` | `validation.Nil` | Passes only for a nil pointer; fails for any non-nil value, even empty |
| `rules.In(vals...)` | `validation.In(vals...)` | Skips nil; works with named types |
| `rules.NotIn(vals...)` | `validation.NotIn(vals...)` | Skips nil; works with named types |
| `rules.Matches(pattern)` | `validation.Match(re)` | Skips nil/empty strings |
| `rules.Length(min, max)` | `validation.Length(min, max)` | `max=0` means no upper bound |
| `rules.RuneLength(min, max)` | `validation.RuneLength(min, max)` | Unicode-aware |
| `rules.Min(v)` / `rules.Max(v)` | `validation.Min(v)` / `validation.Max(v)` | int, uint, float, time |
| `rules.Expr(expr)` | — | CEL-like expression; fields accessed by JSON name |
| `rules.By(desc, func(any) bool)` | `validation.By(func)` | Custom boolean function |
| `rules.ByString(desc, func(string) bool)` | — | Convenience for string-typed fields |
| `rules.ByError(desc, func(any) error)` | `validation.By(func)` (error variant) | Error message is discarded; use `desc` |
| `rules.Func(desc, func(any) bool)` | — | Equivalent to `By`; prefer `By` for consistency |

### The `rules/is` package

`github.com/invopop/gobl/rules/is` mirrors `github.com/invopop/validation/is`:

```go
// Before
import "github.com/invopop/validation/is"
validation.Field(&obj.URL, is.URL)

// After
import "github.com/invopop/gobl/rules/is"
rules.Field("url",
    rules.Assert("03", "URL must be valid", is.URL),
)
```

## Migration patterns

### Simple required field

```go
// Before
validation.Field(&obj.Name, validation.Required)

// After
rules.Field("name",
    rules.Assert("01", "name is required", rules.Required),
)
```

### Required field with format check

```go
// Before
validation.Field(&obj.Address, validation.Required, is.Email)

// After
rules.Field("addr",
    rules.Assert("01", "email address is required", rules.Required),
    rules.Assert("02", "email address must be valid", is.EmailFormat),
)
```

Note that `Required` and the format check are **separate assertions** with
separate codes, so callers can distinguish a missing value from a malformed one.

### Optional field with format check

Leave out `Required`. All built-in tests skip nil/empty values automatically:

```go
// Before
validation.Field(&obj.URL, is.URL)

// After
rules.Field("url",
    rules.Assert("05", "URL must be valid", is.URL),
)
```

### Allowed-values check (`In`)

```go
// Before
validation.Field(&obj.Category, validation.In("a", "b", "c"))

// After
rules.Field("category",
    rules.Assert("02", "category is not valid", rules.In("a", "b", "c")),
)
```

`rules.In` normalises named string/int types so `In("a", "b")` matches both
`string("a")` and `MyType("a")`.

To allow an optional field to be empty *or* one of the valid values, either
omit `Required` (the `In` test skips nil pointers automatically) or, for
non-pointer named types like `cbc.Key`, extract the logic into a helper:

```go
func isValidCategory(val any) bool {
    key, ok := val.(cbc.Key)
    if !ok || key == "" {
        return true // optional
    }
    for _, def := range validDefs {
        if def.Key == key {
            return true
        }
    }
    return false
}

rules.Field("category",
    rules.Assert("02", "category is not valid",
        rules.By("valid or empty", isValidCategory),
    ),
)
```

### Regex pattern match

```go
// Before
validation.Field(&obj.Code, validation.Match(regexp.MustCompile(`^\d{9}$`)))

// After
rules.Field("code",
    rules.Assert("01", "invalid format", rules.Matches(`^\d{9}$`)),
)
```

### Custom validation logic

Extract the logic into a named private function and use `rules.By`,
`rules.ByString`, or `rules.ByError`:

```go
// Before
validation.Field(&obj.Code,
    validation.By(func(v any) error {
        code, _ := v.(string)
        if !isValidChecksum(code) {
            return errors.New("checksum mismatch")
        }
        return nil
    }),
)

// After
func myCodeChecksumValid(code string) bool {
    return isValidChecksum(code)
}

rules.Field("code",
    rules.Assert("03", "code checksum mismatch",
        rules.ByString("checksum", myCodeChecksumValid),
    ),
)
```

Keep helper functions short and named descriptively — they double as
documentation and make `linkRules()` / `headerRules()` easy to read at a
glance.

### Field must not be set (`Empty` / `Nil`)

Use `rules.Empty` when a field must be absent or zero — the inverse of `Required`:

```go
// Before
validation.Field(&obj.Discount, validation.Empty)

// After
rules.Field("discount",
    rules.Assert("05", "discount must not be set", rules.Empty),
)
```

Use `rules.Nil` when the field must be a nil pointer specifically. Unlike `Empty`,
it fails even when the pointer is non-nil but points to a zero/empty value:

```go
// Before
validation.Field(&obj.Digest, validation.Nil)

// After
rules.Field("digest",
    rules.Assert("06", "digest must not be set", rules.Nil),
)
```

### Object-level (cross-field) assertion

Without `Field`, an assertion receives the full object. This is useful for
cross-field constraints:

```go
// Before
validation.Field(&obj.Digest,
    validation.When(obj.MIME == "", validation.Nil.Error("must be nil when MIME not set")),
)

// After (using Expr — field names by JSON tag)
rules.Assert("06", "digest must be nil when MIME type is not provided",
    rules.Expr(`mime != "" || digest == nil`),
)

// After (using By — handles both *T and T receiver shapes)
func digestRequiresMIME(val any) bool {
    var obj *MyStruct
    switch v := val.(type) {
    case *MyStruct:
        obj = v
    case MyStruct:
        obj = &v
    default:
        return true
    }
    return obj == nil || obj.MIME != "" || obj.Digest == nil
}

rules.Assert("06", "digest must be nil when MIME type is not provided",
    rules.By("no digest without MIME", digestRequiresMIME),
)
```

Prefer `Expr` for simple field comparisons. Use `By` when the logic is more
involved or when you want an explicitly named and testable function.

> **Note on receiver shape:** `rules.Validate(obj)` may pass either `*T` or
> `T` to an object-level `By` function depending on how the object was reached
> (top-level call vs. recursive field traversal). Always handle both in `By`
> helpers for object-level assertions. `Expr` handles this automatically via
> its `buildEnv` logic.

### Conditional validation (`When`)

Replace `validation.When(condition, ...)` with `rules.When(test, ...)`. The
condition is a `Test` and receives the full parent struct.

```go
// Before — context-aware condition
validation.Field(&obj.Stamps,
    validation.When(!internal.IsSigned(ctx), validation.Empty),
)

// After — condition derived from the object itself
rules.When(rules.By("not signed", func(val any) bool {
    e, ok := val.(*Envelope)
    return ok && len(e.Signatures) == 0
}),
    rules.Field("stamps",
        rules.Assert("12", "stamps not allowed before signing",
            rules.Length(0, 0),
        ),
    ),
)
```

> **Context-dependent rules:** The old system supported `context.Context`
> threading through `ValidateWithContext`. Rules have no context. Conditions
> that previously depended on context (e.g. "is signed?") must instead be
> derived from the object's own state. Move such checks to the outermost type
> that carries the relevant state (typically `Envelope`), and use
> `rules.When(...)` there.

### Nested struct fields

The preferred approach is to define rules for each type independently and
register them all. `rules.Validate` recurses into every exported field
automatically, so there is no wiring required between parent and child:

```go
// address_rules.go
func addressRules() *rules.Set {
    return rules.For(new(Address),
        rules.Field("city",
            rules.Assert("01", "city is required", rules.Required),
        ),
    )
}

// person_rules.go
func personRules() *rules.Set {
    return rules.For(new(Person),
        rules.Field("name",
            rules.Assert("01", "name is required", rules.Required),
        ),
        // No wiring for Address — addressRules() is registered separately
        // and applied automatically when rules.Validate recurses into the field.
    )
}

func init() {
    rules.Register("mypkg", rules.GOBL.Add("MYPKG"),
        addressRules(),
        personRules(),
    )
}
```

When you need to add rules about a nested type from the **parent's perspective**
(e.g. regime-specific constraints that don't belong on the child type itself),
nest `rules.Field` calls to drill down the path:

```go
// Before — regime-specific Validate method on the parent
func (inv *Invoice) Validate() error {
    return validation.ValidateStruct(inv,
        validation.Field(&inv.Supplier, validation.Required),
        // inside supplier.Validate(), further checks on TaxID...
    )
}

// After — regime rule set drilling into nested fields
func invoiceRules() *rules.Set {
    return rules.For(new(Invoice),
        rules.When(tax.RegimeIn("XX"),
            rules.Field("supplier",
                rules.Assert("01", "supplier is required", rules.Required),
                rules.Field("tax_id",
                    rules.Assert("02", "supplier tax ID is required", rules.Required),
                    rules.Field("code",
                        rules.Assert("03", "supplier tax ID must have a code", rules.Required),
                    ),
                ),
            ),
        ),
    )
}
```

Each `rules.Field` in the chain constrains the context for its children, so
assertions and tests inside `rules.Field("tax_id", ...)` operate on the
`TaxIdentity` struct, not the outer `Invoice`.

### Slice fields (`Each`)

`rules.Each` is a nameless `Def` that iterates over the elements of the current
context. It is used **inside** a `rules.Field` that targets a slice field:

```go
// Before
func (obj *MyStruct) Validate() error {
    return validation.ValidateStruct(obj,
        validation.Field(&obj.Lines,
            validation.Each(validation.Required, validation.By(lineIsValid)),
        ),
    )
}

// After
func myStructRules() *rules.Set {
    return rules.For(new(MyStruct),
        rules.Field("lines",
            rules.Each(
                rules.Assert("01", "line must not be empty", rules.Required),
                rules.Assert("02", "line must be valid", rules.By("valid", lineIsValid)),
            ),
        ),
    )
}
```

Faults from `Each` carry a path like `lines[0]`, `lines[1]`, etc.

Because `Each` is just a `Def` inside `Field`, whole-slice and per-element
assertions can coexist on the same field naturally:

```go
rules.Field("lines",
    rules.Assert("01", "no duplicate line codes",
        rules.By("no duplicates", hasNoDuplicateLineCodes),
    ),
    rules.Each(
        rules.Field("code",
            rules.Assert("02", "line code is required", rules.Required),
        ),
    ),
)
```

If the element type has its own registered rule set, those rules are applied
automatically during recursive validation — `Each` is only needed when you want
to add **additional** assertions from the parent's perspective that don't belong
on the element type itself.

`rules.Each` panics at initialisation time when used outside a slice field.

### Named value types (e.g. `cbc.Code`, `tax.Rate`)

`rules.For` works with named non-struct types too:

```go
func myCodeRules() *rules.Set {
    return rules.For(MyCode(""),
        rules.Assert("01", "code must not be empty", rules.Required),
        rules.Assert("02", "code must be alphanumeric", is.Alphanumeric),
    )
}
```

Inside `Expr`, the value is exposed as `this`:

```go
rules.Assert("02", "code must not exceed 10 characters",
    rules.Expr(`len(this) <= 10`),
)
```

## Wiring it up

### In the package's `init` file

```go
// mypkg.go
func init() {
    schema.Register(schema.GOBL.Add("mypkg"),
        MyStruct{},
    )
    rules.Register(
        "mypkg",
        rules.GOBL.Add("MYPKG"),
        myStructRules(),
    )
}
```

### Removing `Validate` / `ValidateWithContext`

Once all rules are expressed in the `rules.Set`, delete the `Validate()` and
`ValidateWithContext()` methods from the struct. `rules.Validate(obj)` recurses
into all exported fields automatically — no wiring is needed on the types
themselves.

Remove unused imports: `"context"`, `"github.com/invopop/gobl/internal"`, and
`"github.com/invopop/validation"`.

Keep any exported `validation.Rule` helpers (e.g. `StampsHas`) that are still
consumed by other packages — those will be migrated separately.

## Updating tests

Replace `.Validate()` calls with `rules.Validate()`:

```go
// Before
err := obj.Validate()
assert.NoError(t, err)
assert.ErrorContains(t, err, "field: cannot be blank")

// After
import "github.com/invopop/gobl/rules"

err := rules.Validate(obj)
assert.NoError(t, err)
assert.ErrorContains(t, err, "field: description from assertion")
```

`rules.Validate` returns `rules.Faults`, which implements `error`. A nil return
means no faults. The full error string format is:

```
[GOBL-PKG-STRUCT-01] field: assertion description
```

You can also assert on the structured `Faults` value directly:

```go
faults := rules.Validate(obj)
require.NotNil(t, faults)
assert.True(t, faults.HasPath("field"))
assert.True(t, faults.HasCode("GOBL-PKG-STRUCT-01"))
assert.Equal(t, "assertion description", faults.First().Message())
```

## Assertion code conventions

Codes within a set are short local identifiers (e.g. `"01"`, `"02"`). They are
prefixed during `Register` to form globally unique codes. Follow this
convention:

- `01`–`09`: field-level assertions, in the order fields appear in the struct
- `10`–`19`: object-level (cross-field) assertions
- `20`+: reserved for `When` conditional subsets if needed

The fully-qualified code is constructed as:
`{REGISTER_PREFIX}-{PKG}-{STRUCT}-{LOCAL_CODE}`

For example, a `"03"` assertion on `head.Header` registered under `GOBL-HEAD`
becomes `GOBL-HEAD-HEADER-03`.
