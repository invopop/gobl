// Package convert provides a central register of converters between GOBL
// and external document formats.
//
// GOBL itself does not convert. Packages that do, such as gobl.ubl or
// gobl.cii, register a Converter from their init functions so that importing
// them makes their formats available here:
//
//	import _ "github.com/invopop/gobl.ubl"
//
//	env, err := convert.Import(data)
//	out, err := convert.Export(env, "ubl+peppol", "ubl+en16931")
package convert

import (
	"slices"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/schema"
)

// Directions of a conversion.
const (
	DirectionImport cbc.Key = "import"
	DirectionExport cbc.Key = "export"
)

var (
	// ErrUnknownFormat is provided when a format key is not registered, or
	// when no converter recognizes the incoming data.
	ErrUnknownFormat = gobl.NewError("unknown-format")

	// ErrAmbiguous is provided when more than one converter recognizes the
	// incoming data.
	ErrAmbiguous = gobl.NewError("ambiguous-format")

	// ErrNotSupported is provided when none of the requested formats accept
	// the envelope for export.
	ErrNotSupported = gobl.NewError("not-supported")

	// ErrConversion wraps errors returned by a converter. Errors that are
	// already GOBL errors, such as validation faults, are returned unchanged.
	ErrConversion = gobl.NewError("conversion")
)

// Format describes one external document format that a converter can import
// and/or export: a syntax together with the specification its documents
// follow, e.g. UBL following Peppol BIS Billing.
type Format struct {
	// Key uniquely identifies the format, e.g. "ubl+peppol".
	Key cbc.Key `json:"key" jsonschema:"title=Key"`
	// Name of the format.
	Name i18n.String `json:"name" jsonschema:"title=Name"`
	// Description of the format.
	Description i18n.String `json:"description,omitempty" jsonschema:"title=Description"`
	// MIME type of the format's documents, e.g. "application/xml".
	MIME string `json:"mime,omitempty" jsonschema:"title=MIME Type"`
	// Syntax the format's documents are written in, e.g. "ubl" or "cii".
	Syntax cbc.Key `json:"syntax,omitempty" jsonschema:"title=Syntax"`
	// Countries where the format applies, including unions such as "EU".
	// Empty means no country restriction.
	Countries []l10n.Code `json:"countries,omitempty" jsonschema:"title=Countries"`
	// Addons the GOBL document must include to be exported into this format.
	Addons []cbc.Key `json:"addons,omitempty" jsonschema:"title=Addons"`
	// Import lists GOBL schemas this format can be converted into.
	Import []schema.ID `json:"import,omitempty" jsonschema:"title=Import Schemas"`
	// Export lists GOBL schemas that can be converted into this format.
	Export []schema.ID `json:"export,omitempty" jsonschema:"title=Export Schemas"`
}

// Converter is implemented by packages that convert between GOBL and one
// or more formats.
type Converter interface {
	// Formats lists the formats handled by the converter.
	Formats() []*Format
	// Detect returns the key of the format the input belongs to, or an
	// empty key if it is not recognized. It should not fully parse the data.
	Detect(in *Input) cbc.Key
	// Import converts data in the given format into a GOBL envelope.
	Import(key cbc.Key, data []byte) (*gobl.Envelope, error)
	// Accepts reports whether the envelope can be exported into the format.
	// It is only called once the document's schema is in the format's Export
	// list and the document includes the format's Addons.
	Accepts(key cbc.Key, env *gobl.Envelope) bool
	// Export converts the envelope into the given format.
	Export(key cbc.Key, env *gobl.Envelope) ([]byte, error)
}

// Output is the result of an export.
type Output struct {
	// Format the envelope was exported into.
	Format *Format
	// Data of the exported document.
	Data []byte
}

// Conversion describes a single 1:1 conversion available in the register.
type Conversion struct {
	// Format key of the external format.
	Format cbc.Key `json:"format" jsonschema:"title=Format"`
	// Schema of the GOBL document.
	Schema schema.ID `json:"schema" jsonschema:"title=Schema"`
	// Direction of the conversion, import into GOBL or export from it.
	Direction cbc.Key `json:"direction" jsonschema:"title=Direction"`
}

// exports checks if the envelope's document has one of the format's export
// schemas and includes all of its addons.
func (c *Format) exports(env *gobl.Envelope) bool {
	if env == nil || env.Document == nil || !slices.Contains(c.Export, env.Document.Schema) {
		return false
	}
	if len(c.Addons) == 0 {
		return true
	}
	doc, ok := env.Extract().(interface{ GetAddons() []cbc.Key })
	if !ok {
		return false
	}
	addons := doc.GetAddons()
	for _, a := range c.Addons {
		if !a.In(addons...) {
			return false
		}
	}
	return true
}

// appliesTo checks if the format can be used in the country, either directly,
// through a union the country belongs to, or because it has no restriction.
func (c *Format) appliesTo(country l10n.Code) bool {
	if len(c.Countries) == 0 {
		return true
	}
	for _, cc := range c.Countries {
		if cc == country {
			return true
		}
		if u := l10n.Union(cc); u != nil && u.HasMember(country) {
			return true
		}
	}
	return false
}

// Register adds the converter and its formats to the global register. This
// is expected to be called from package init functions, and panics if a
// format key is empty or already registered.
func Register(c Converter) {
	converters.add(c)
}

// Formats provides all the registered formats, sorted by key.
func Formats() []*Format {
	return converters.formats()
}

// FormatFor provides the format registered with the key, or nil.
func FormatFor(key cbc.Key) *Format {
	return converters.formatFor(key)
}

// FormatsFor provides the formats that apply to the country, directly or
// through a union it belongs to, along with those without a country restriction.
func FormatsFor(country l10n.Code) []*Format {
	return converters.formatsFor(country)
}

// Conversions lists every 1:1 conversion available from the registered
// formats.
func Conversions() []*Conversion {
	return converters.conversions()
}

// Detect determines the format of the incoming data by asking each
// candidate converter. Keys, if provided, limit the candidates to those
// formats.
func Detect(data []byte, keys ...cbc.Key) (*Format, error) {
	e, err := converters.detect(data, keys)
	if err != nil {
		return nil, err
	}
	return e.format, nil
}

// Import detects the format of the incoming data and converts it into a GOBL
// envelope. Keys, if provided, limit detection to those formats.
func Import(data []byte, keys ...cbc.Key) (*gobl.Envelope, error) {
	e, err := converters.detect(data, keys)
	if err != nil {
		return nil, err
	}
	env, err := e.converter.Import(e.format.Key, data)
	if err != nil {
		return nil, ErrConversion.WithCause(err)
	}
	return env, nil
}

// Export converts the envelope into the first of the keys, in order of
// preference, that exports the document's schema, whose addons the document
// includes, and whose converter accepts it.
func Export(env *gobl.Envelope, keys ...cbc.Key) (*Output, error) {
	if len(keys) == 0 {
		return nil, ErrUnknownFormat.WithReason("no formats provided")
	}
	for _, k := range keys {
		if converters.entryFor(k) == nil {
			return nil, ErrUnknownFormat.WithReason("format %s not registered", k)
		}
	}
	for _, k := range keys {
		e := converters.entryFor(k)
		if !e.format.exports(env) || !e.converter.Accepts(k, env) {
			continue
		}
		data, err := e.converter.Export(k, env)
		if err != nil {
			return nil, ErrConversion.WithCause(err)
		}
		return &Output{Format: e.format, Data: data}, nil
	}
	return nil, ErrNotSupported.WithReason("no format in %v accepts the envelope", keys)
}
