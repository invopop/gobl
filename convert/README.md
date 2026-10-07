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

- **Format**: a syntax together with the specification its documents follow,
  such as UBL following XRechnung. It is the identity a document declares about
  itself, and what the register lists, detects, and exports into, each
  identified by a unique key. A format covers every document type its
  specification defines, listed by the GOBL schemas in `Import` and `Export`.
  Containers that carry a document, such as Factur-X's PDF/A-3 or a Peppol
  SBDH envelope, are not part of the format.
- **Converter**: the implementation that handles one or more formats, usually
  every format of a single syntax.
- **Conversion**: a single 1:1 conversion between a format and a GOBL schema in
  one direction, listed by `convert.Conversions`.

## Naming directions

The same terms are used by the register and recommended for converter
packages, so that each name always means one direction:

| Term | Direction |
| --- | --- |
| **Import** | From an external format into GOBL. |
| **Export** | From GOBL into an external format. |
| **Decode** | From bytes into the syntax's document structures, such as a UBL invoice. |
| **Encode** | From the syntax's document structures into bytes. |

Decoding and encoding are lossless serialization, while importing and exporting
map between two data models. A converter package typically offers all four, for
example:

```go
doc, err := ubl.Decode(data)
env, err := ubl.Import(doc)

doc, err := ubl.Export(env, ubl.WithFormat(ubl.FormatPeppol))
data, err := ubl.Encode(doc)
```

The register's `Import` is then a decode followed by an import, and its
`Export` an export followed by an encode. "Convert" is kept for the general
act and this package's name, as it implies no direction.

## Using the register

| Function | Purpose |
| --- | --- |
| `Import(data, keys...)` | Detect the format of the data and convert it into a GOBL envelope. Keys, if given, limit detection to those formats. |
| `Export(env, keys...)` | Convert the envelope into the first of the keys, in order of preference, that can take it. |
| `Detect(data, keys...)` | Detect the format without converting. |
| `Formats()`, `FormatFor(key)` | List all formats, or look one up. |
| `FormatsFor(country)` | List the formats that apply in a country, including through unions such as the EU, and those without a country restriction. |
| `Conversions()` | List every conversion available. |

Errors are GOBL errors that can be compared with `errors.Is`:

- `ErrUnknownFormat`: a key is not registered, or no converter recognizes the data.
- `ErrAmbiguous`: more than one converter recognizes the data.
- `ErrNotSupported`: none of the requested formats can take the envelope.
- `ErrConversion`: the converter failed. GOBL errors returned by a converter,
  such as validation faults, are passed on unchanged.

## Implementing a converter

A converter implements the `Converter` interface and registers itself once:

```go
package ubl

var formatPeppol = &convert.Format{
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

func (converter) Formats() []*convert.Format {
	return []*convert.Format{formatEN16931, formatPeppol}
}
```

### Detection

`Detect` receives an `*convert.Input` and returns the key of one of the
converter's own formats, or an empty key. Every converter is asked, so that two
claiming the same data is reported as ambiguous rather than resolved by
registration order. Detection must therefore be cheap:

- read only as much of the data as is needed to tell formats apart, and never
  parse the whole document;
- tell apart your own formats yourself, so that the register only has to deal
  with clashes between packages;
- return an empty key for anything you are unsure about.

The `Input` is shared between all the detectors asked about the same data.
Packages store the values they extract on it with their own unexported key type,
in the same way as `context.Context`, and export a helper so that other
packages built on the same syntax can reuse them instead of parsing the data
again:

```go
type headerKey struct{}

type Header struct {
	CustomizationID string
	ProfileID       string
	Err             error
}

// ReadHeader provides the UBL specification and business process
// identifiers of the input, reading them only once.
func ReadHeader(in *convert.Input) *Header {
	if v, ok := in.Get(headerKey{}); ok {
		return v.(*Header)
	}
	h := readHeader(in.Data)
	in.Set(headerKey{}, h)
	return h
}
```

Store errors along with values, so that a malformed document is only read once.

### Import

`Import` receives the detected key and the raw data, and returns an envelope.
It is only called for a key the converter returned from `Detect`.

### Export

Before asking the converter, the register checks that the document's schema is
in the format's `Export` list and that the document includes every addon in
`Addons`. `Accepts` only needs to cover any further conditions specific to the
format, such as a document type the format cannot represent. `Export` then
returns the converted data.

Formats expect documents to already include their addons. A converter that
adds missing addons during export does so only to help existing users migrate,
and new converters should not rely on it.

## Regional variants

Most specifications narrow a more general one: Peppol BIS Billing is an EN 16931
CIUS, and the French CIUS builds on Peppol. The recommended structure is a base
import and export in the syntax package, with no behavior specific to a format,
and functions for each format that adjust the documents they produce. The syntax
package registers the formats that apply in any country. Packages that implement
a regional specification, such as gobl.fr.ctc, build on the syntax package's
base and register their own formats when imported.

## Naming formats

Format keys are a public contract: they are stored by users, passed to
`Export`, and shown in user interfaces. They follow the `cbc.Key` rules
(lowercase letters, numbers, `-` and `+`, at most 64 characters), and are built
from the syntax followed by one layer for each specification the format builds
on, from the most general to the most specific:

```
<syntax>[+<layer>...]
```

For example:

- `ubl+en16931`
- `ubl+peppol`, `ubl+peppol+self-billing`, `ubl+peppol+invoice-response`
- `ubl+peppol+fr-cius-v1`, `ubl+peppol+fr-extended-v1`, `cii+peppol+fr-cius-v1`
- `ubl+de-xrechnung-v3`, `cii+de-xrechnung-v3`
- `cii+fr-facturx-v1`, `cii+fr-facturx-v1+basic`, `cii+de-zugferd-v2+extended`
- `ubl+sa-zatca-v1`
- `cdar+peppol+fr-cdv-v1`

### Syntax

The syntax is the base language the document is written in, the same value as the
format's `Syntax` field, such as `ubl` or `cii`. It comes first so that all the
formats of a syntax can be found with `key.HasPrefix`.

### Layers

Each layer names a specification that narrows the one before it:

- **Base and network layers** that apply in any country, such as `en16931`,
  `peppol`, or `self-billing`, have no version.
- **Regional layers** use the key of the GOBL addon that implements the
  specification, with its major version, such as `de-xrechnung-v3` or
  `sa-zatca-v1`. When there is no addon, or one addon covers several
  specifications, use a key in the same style: a country code, the short name
  of the specification, and its major version, such as `fr-cius-v1`.

A specification with several profiles, each declared with its own identifier,
adds the profile as a final layer, such as `+basic` or `+extended`.

Only add a layer for the specification a document declares. A layer may be
left out when the specification it would name adds nothing to tell documents
apart, as with `ubl+de-xrechnung-v3`, which is not written as
`ubl+en16931+de-xrechnung-v3`.

The key describes where the specification comes from. It does not define the
conversion: each format lists the behavior it applies itself.

When a format is its own syntax and has no separate specification, as with
FacturaE or CFDI, use the addon key alone, such as `es-facturae-v3` or
`mx-cfdi-v4`, and set `Syntax` to the format's name.

### Guidelines

- Create **one format per specification identifier**, such as the UBL
  `CustomizationID` or the CII guideline ID, rather than one per document type.
  List the document types the format supports in `Import` and `Export`.
- Include only the **major version**. Minor versions of a specification are
  handled inside the converter, and should not change the key users store.
- Use American spelling and the official short name of the specification,
  without words like `format` or `invoice` that every format would share.
- Only register keys for the specifications your package implements.
- Never rename a published key. Add a new format for a new major version, and
  keep the old one for as long as it is supported.

## Describing formats

The remaining fields help user interfaces list and filter formats:

- `Name` and `Description`: what the format is, in each supported language.
- `MIME`: the media type of the data, such as `application/xml`.
- `Countries`: where the format applies, including unions such as `EU`. Leave
  it empty for formats that apply anywhere, such as plain EN 16931.
