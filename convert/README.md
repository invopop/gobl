# Convert

The `convert` package is a central register of converters between GOBL and
external document formats. GOBL itself does not convert: packages that do, such
as [gobl.ubl](https://github.com/invopop/gobl.ubl) or
[gobl.cii](https://github.com/invopop/gobl.cii), register a converter from their
`init` function, so a blank import is enough to make their formats available:

```go
import (
	"github.com/invopop/gobl/convert"
	_ "github.com/invopop/gobl.ubl"
)

env, err := convert.Import(data)
out, err := convert.Export(env, "ubl+peppol", "ubl+en16931")
```

## Concepts

- **Context**: one variant of an external format, such as UBL following the
  XRechnung specification. Contexts are what the register lists, detects, and
  exports into, each identified by a unique key.
- **Converter**: the implementation that handles one or more contexts, usually
  every context of a single syntax.
- **Conversion**: a single 1:1 conversion between a context and a GOBL schema in
  one direction, listed by `convert.Conversions`.

## Using the register

| Function | Purpose |
| --- | --- |
| `Import(data, keys...)` | Detect the context of the data and convert it into a GOBL envelope. Keys, if given, limit detection to those contexts. |
| `Export(env, keys...)` | Convert the envelope into the first of the keys, in order of preference, that can take it. |
| `Detect(data, keys...)` | Detect the context without converting. |
| `Contexts()`, `ContextFor(key)` | List all contexts, or look one up. |
| `ContextsFor(country)` | List the contexts that apply in a country, including through unions such as the EU, and those without a country restriction. |
| `Conversions()` | List every conversion available. |

Errors are GOBL errors that can be compared with `errors.Is`:

- `ErrUnknownContext`: a key is not registered, or no converter recognizes the data.
- `ErrAmbiguous`: more than one converter recognizes the data.
- `ErrNotSupported`: none of the requested contexts can take the envelope.
- `ErrConversion`: the converter failed. GOBL errors returned by a converter,
  such as validation faults, are passed on unchanged.

## Implementing a converter

A converter implements the `Converter` interface and registers itself once:

```go
package ubl

var contextPeppol = &convert.Context{
	Key:    "ubl+peppol",
	Name:   i18n.NewString("UBL Peppol BIS Billing 3"),
	MIME:   "application/xml",
	Syntax: "ubl",
	Addons: []cbc.Key{en16931.V2017},
	Import: []schema.ID{schema.Lookup(bill.Invoice{})},
	Export: []schema.ID{schema.Lookup(bill.Invoice{})},
}

type converter struct{}

func init() {
	convert.Register(converter{})
}

func (converter) Contexts() []*convert.Context {
	return []*convert.Context{contextEN16931, contextPeppol}
}
```

### Detection

`Detect` receives an `*convert.Input` and returns the key of one of the
converter's own contexts, or an empty key. Every converter is asked, so that two
claiming the same data is reported as ambiguous rather than resolved by
registration order. Detection must therefore be cheap:

- read only as much of the data as is needed to tell contexts apart, and never
  parse the whole document;
- tell apart your own contexts yourself, so that the register only has to deal
  with clashes between packages;
- return an empty key for anything you are unsure about.

The `Input` is shared between all the detectors asked about the same data.
Packages store the values they extract on it with their own unexported key type,
in the same way as `context.Context`, and export a helper so that other
packages built on the same syntax can reuse them instead of parsing the data
again:

```go
type documentContextKey struct{}

type DocumentContext struct {
	CustomizationID string
	ProfileID       string
	Err             error
}

// ReadDocumentContext provides the UBL specification and business process
// identifiers of the input, reading them only once.
func ReadDocumentContext(in *convert.Input) *DocumentContext {
	if v, ok := in.Get(documentContextKey{}); ok {
		return v.(*DocumentContext)
	}
	dc := readDocumentContext(in.Data)
	in.Set(documentContextKey{}, dc)
	return dc
}
```

Store errors along with values, so that a malformed document is only read once.

### Import

`Import` receives the detected key and the raw data, and returns an envelope.
It is only called for a key the converter returned from `Detect`.

### Export

Before asking the converter, the register checks that the document's schema is
in the context's `Export` list and that the document includes every addon in
`Addons`. `Accepts` only needs to cover any further conditions specific to the
format, such as a document type the context cannot represent. `Export` then
returns the converted data.

Contexts expect documents to already include their addons. A converter that
adds missing addons during export does so only to help existing users migrate,
and new converters should not rely on it.

## Regional variants

Most specifications narrow a more general one: Peppol BIS Billing is an EN 16931
CIUS, and the French CIUS builds on Peppol. The recommended structure is a base
conversion in the syntax package, with no behavior specific to a context, and
layers for each specification that are applied on top of it. The syntax package
registers the contexts that apply in any country. Packages that implement a
regional specification, such as gobl.fr.ctc, build on the syntax package's base
conversion and register their own contexts when imported.

## Naming contexts

Context keys are a public contract: they are stored by users, passed to
`Export`, and shown in user interfaces. They follow the `cbc.Key` rules
(lowercase letters, numbers, `-` and `+`, at most 64 characters), and are built
from the syntax followed by one layer for each specification the context builds
on, from the most general to the most specific:

```
<syntax>[+<layer>...]
```

For example:

- `ubl+en16931`
- `ubl+peppol`, `ubl+peppol+self-billing`, `ubl+peppol+invoice-response`
- `ubl+peppol+fr-cius-v1`, `ubl+peppol+fr-extended-v1`
- `ubl+de-xrechnung-v3`, `cii+de-xrechnung-v3`
- `ubl+sa-zatca-v1`

### Syntax

The syntax is the base format the document is written in, the same value as the
context's `Syntax` field, such as `ubl` or `cii`. It comes first so that all the
contexts of a syntax can be found with `key.HasPrefix`.

### Layers

Each layer names a specification that narrows the one before it:

- **Base and network layers** that apply in any country, such as `en16931`,
  `peppol`, or `self-billing`, have no version.
- **Regional layers** use the key of the GOBL addon that implements the
  specification, with its major version, such as `de-xrechnung-v3` or
  `sa-zatca-v1`. When there is no addon, or one addon covers several
  specifications, use a key in the same style: a country code, the short name
  of the specification, and its major version, such as `fr-cius-v1`.

Only add a layer for the specification a document declares. A layer may be
left out when the specification it would name adds nothing to tell documents
apart, as with `ubl+de-xrechnung-v3`, which is not written as
`ubl+en16931+de-xrechnung-v3`.

The key describes where the specification comes from. It does not define the
conversion: each context lists the behavior it applies itself.

When a format is its own syntax and has no separate specification, as with
FacturaE or CFDI, use the addon key alone, such as `es-facturae-v3` or
`mx-cfdi-v4`, and set `Syntax` to the format's name.

### Guidelines

- Create **one context per specification identifier**, such as the UBL
  `CustomizationID` or the CII guideline ID, rather than one per document type.
  List the document types the context supports in `Import` and `Export`.
- Include only the **major version**. Minor versions of a specification are
  handled inside the converter, and should not change the key users store.
- Use American spelling and the official short name of the specification,
  without words like `format` or `invoice` that every context would share.
- Only register keys for the specifications your package implements.
- Never rename a published key. Add a new context for a new major version, and
  keep the old one for as long as it is supported.

## Describing contexts

The remaining fields help user interfaces list and filter contexts:

- `Name` and `Description`: what the context is, in each supported language.
- `MIME`: the media type of the data, such as `application/xml`.
- `Countries`: where the context applies, including unions such as `EU`. Leave
  it empty for contexts that apply anywhere, such as plain EN 16931.
