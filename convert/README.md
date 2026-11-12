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

res, err := convert.Import(data)
out, err := convert.Export(env, "ubl+peppol+bis-v3",
	convert.WithAlternatives("ubl+en16931-v2017"),
)
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
| `Import(data, opts...)` | Detect the format of the data and convert it into a GOBL envelope, with any files embedded in the data. |
| `Export(env, key, opts...)` | Convert the envelope into the format with the key, or the first of the alternatives that can take it. |
| `Detect(data, keys...)` | Detect the format without converting. |
| `Formats()`, `FormatFor(key)` | List all formats, or look one up. |
| `FormatsFor(country)` | List the formats that apply in a country, including through unions such as the EU, and those without a country restriction. |
| `Conversions()` | List every conversion available. |

`Export` takes options:

- `WithAlternatives(keys...)`: further formats to try, in order of preference,
  when the requested format does not take the envelope.
- `WithAttachments(atts...)`: files to embed in the exported document, such as
  a PDF version of the invoice. Only formats with `Attachments` set can embed
  them, and other formats refuse the envelope.

```go
pdf := &convert.Attachment{
	ID:          "INV-001-PDF",
	Name:        "INV-001.pdf",
	Description: "Invoice INV-001",
	MIME:        "application/pdf",
	Data:        data,
}
out, err := convert.Export(env, "ubl+peppol+bis-v3", convert.WithAttachments(pdf))
```

`Export` returns an `ExportResult` with the `Format` used and the exported
`Data`.

`Import` returns an `ImportResult` with the detected `Format`, the
`Envelope`, and the `Attachments` embedded in the data, such as a PDF version
of the invoice. The envelope does not include them, so they can be stored
separately. `Import` takes options:

- `WithFormats(keys...)`: limit detection to these formats.

```go
res, err := convert.Import(data, convert.WithFormats("ubl+peppol+bis-v3"))
for _, a := range res.Attachments {
	// store a.Data as a.Name, of type a.MIME
}
```

Errors are GOBL errors that can be compared with `errors.Is`:

- `ErrUnknownFormat`: a key is not registered, or no converter recognizes the data.
- `ErrAmbiguous`: more than one converter recognizes the data.
- `ErrNotSupported`: none of the requested formats can take the envelope. The
  error carries the reasons each format gave, as rules faults when they all
  are, so they can be shown in the same way as validation faults.
- `ErrConversion`: the converter failed. GOBL errors returned by a converter,
  such as validation faults, are passed on unchanged.

## Implementing a converter

A converter implements the `Converter` interface and registers itself once:

```go
package ubl

var (
	invoiceSchema = schema.Lookup(bill.Invoice{})
	statusSchema  = schema.Lookup(bill.Status{})
)

var formatEN16931 = &convert.Format{
	Key:    "ubl+en16931-v2017",
	Name:   i18n.NewString("UBL EN 16931:2017"),
	MIME:   "application/xml",
	Syntax: "ubl",
	Import: []schema.ID{invoiceSchema},
	Export: []schema.ID{invoiceSchema},
}

// Invoices are written as Peppol BIS Billing, and statuses as Peppol
// Invoice Response.
var formatPeppol = &convert.Format{
	Key:    "ubl+peppol+bis-v3",
	Name:   i18n.NewString("UBL Peppol BIS 3"),
	MIME:   "application/xml",
	Syntax: "ubl",
	Import: []schema.ID{invoiceSchema, statusSchema},
	Export: []schema.ID{invoiceSchema, statusSchema},
}

// Only invoices need the EN 16931 addon.
var formatRules = map[cbc.Key]*rules.Set{
	formatEN16931.Key: convert.Rules(formatEN16931.Key,
		rules.For(new(bill.Invoice), convert.RequireAddons("01", en16931.V2017)),
	),
	formatPeppol.Key: convert.Rules(formatPeppol.Key,
		rules.For(new(bill.Invoice), convert.RequireAddons("01", en16931.V2017)),
	),
}

type converter struct{}

func init() {
	convert.Register(converter{})
}

func (converter) Formats() []*convert.Format {
	return []*convert.Format{formatEN16931, formatPeppol}
}

func (converter) Accepts(key cbc.Key, env *gobl.Envelope) error {
	if fs := formatRules[key].Validate(env.Extract()); fs != nil {
		return fs
	}
	return nil
}

func (converter) Export(key cbc.Key, env *gobl.Envelope, opts *convert.ExportOptions) ([]byte, error) {
	// Convert, embedding opts.Attachments, and encode.
}
```

A format can cover any GOBL document type. The French CTC lifecycle messages,
for example, carry both statuses and payments:

```go
package cdar

var formatCDV = &convert.Format{
	Key:       "cdar+fr-ctc-cdv-v1",
	Name:      i18n.NewString("CDAR France CTC CDV"),
	MIME:      "application/xml",
	Syntax:    "cdar",
	Countries: []l10n.Code{l10n.FR},
	Import:    []schema.ID{schema.Lookup(bill.Status{}), schema.Lookup(bill.Payment{})},
	Export:    []schema.ID{schema.Lookup(bill.Status{}), schema.Lookup(bill.Payment{})},
}

var rulesCDV = convert.Rules(formatCDV.Key,
	rules.For(new(bill.Status), convert.RequireAddons("01", flow6.V1)),
	rules.For(new(bill.Payment), convert.RequireAddons("01", flow6.V1)),
)
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

`Import` receives the detected key, the raw data and the `ImportOptions`, and
returns an `ImportResult` with the envelope and the files embedded in the data
as `Attachments`. The register sets the result's `Format`. It is only called
for a key the converter returned from `Detect`.

The options' `Formats` are those requested, if any. Converters may use them
to import in a format they cannot tell apart from the data alone, such as a
profile only identified by its PDF container, or ignore them.

### Export

Before asking the converter, the register checks that the document's schema is
in the format's `Export` list. `Accepts` then checks every other condition the
format has, such as the addons, document types, or tags each schema needs, and
returns the reason when the document does not meet them. `Export` then returns
the converted data.

Express the conditions as rules where possible, so that the reasons are
returned as faults with codes and paths, just like validation faults:

- `convert.Rules(key, sets...)` provides a rule set whose fault codes are
  prefixed with the format key, such as `UBL-PEPPOL-BIS-V3-BILL-INVOICE-01`.
  The set is not added to the global rules registry, so it only applies when
  the converter validates documents with it.
- `convert.RequireAddons(id, keys...)` provides a rule for the `$addons` field
  that faults when the document does not include every addon.

Return any other error for conditions that rules cannot express.

`Export` receives the options the caller gave. Formats that set `Attachments`
embed `opts.Attachments` in the document. The register only provides
attachments to those formats, so others can ignore them.

Formats expect documents to already include their addons. A converter that
adds missing addons during export does so only to help existing users migrate,
and new converters should not rely on it.

## Regional variants

Most specifications narrow a more general one: Peppol BIS Billing and the
French CTC CIUS are both EN 16931 CIUSs. The recommended structure is a base
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

- `ubl+en16931-v2017`, `cii+en16931-v2017`
- `ubl+peppol+bis-v3`, `cii+peppol+bis-v3`
- `ubl+fr-ctc-cius-v1`, `ubl+fr-ctc-extended-v1`, `cii+fr-ctc-cius-v1`
- `ubl+de-xrechnung-v3`, `cii+de-xrechnung-v3`
- `cii+fr-facturx-v1`, `cii+fr-facturx-v1+basic`, `cii+de-zugferd-v2+extended`
- `ubl+sa-zatca-v1`
- `cdar+fr-ctc-cdv-v1`, `cdar+fr-ctc-cdv-v1+ppf`

### Syntax

The syntax is the base language the document is written in, the same value as the
format's `Syntax` field, such as `ubl` or `cii`. It comes first so that all the
formats of a syntax can be found with `key.HasPrefix`.

### Layers

Each layer names a specification that narrows the one before it:

- **Network layers**, such as `peppol`, name where documents are exchanged
  and have no version.
- **Specification layers** that apply in any country carry their version,
  such as `en16931-v2017` or `bis-v3`, so a new edition can be added
  alongside, as `en16931-v2026` would be.
- **Regional layers** use the key of the GOBL addon that implements the
  specification, with its major version, such as `de-xrechnung-v3` or
  `sa-zatca-v1`. When there is no addon, or one addon covers several
  specifications, use a key in the same style: the addon family, the short
  name of the specification, and its major version, such as `fr-ctc-cius-v1`.
  A specification exchanged over a network but not built on one of its
  specifications has no network layer: `ubl+fr-ctc-cius-v1`, not
  `ubl+peppol+fr-ctc-cius-v1`.

A specification with several profiles, each declared with its own identifier,
adds the profile as a final layer, such as `+basic` or `+extended`, as does a
variant the sender chooses, such as `cdar+fr-ctc-cdv-v1+ppf` for copies sent
to the French PPF.

Only add a layer for the specification a document declares. A layer may be
left out when the specification it would name adds nothing to tell documents
apart, as with `ubl+de-xrechnung-v3`, which is not written as
`ubl+en16931-v2017+de-xrechnung-v3`.

The key describes where the specification comes from. It does not define the
conversion: each format lists the behavior it applies itself.

When a format is its own syntax and has no separate specification, as with
FacturaE or CFDI, use the addon key alone, such as `es-facturae-v3` or
`mx-cfdi-v4`, and set `Syntax` to the format's name.

### Guidelines

- Create **one format per choice the sender makes**. When the GOBL document
  decides which specification identifier to write, such as the UBL
  `CustomizationID` of Peppol BIS Billing, Self-Billing, or Invoice Response,
  use one format and pick the identifier from the document. When the sender
  has to choose, as between EN 16931 and Peppol, or between the BASIC and
  EXTENDED profiles of Factur-X, use one format each.
- List the document types the format supports in `Import` and `Export`, and
  check what each of them needs, such as addons, in `Accepts`.
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
- `Attachments`: whether exported documents can embed attachments.
