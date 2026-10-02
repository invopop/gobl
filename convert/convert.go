// Package convert provides a central register of converters between GOBL
// and external document formats.
//
// GOBL itself does not convert. Packages that do, such as gobl.ubl or
// gobl.cii, register a Converter from their init functions so that importing
// them makes their contexts available here:
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
	// ErrUnknownContext is provided when a context key is not registered, or
	// when no converter recognizes the incoming data.
	ErrUnknownContext = gobl.NewError("unknown-context")

	// ErrAmbiguous is provided when more than one converter recognizes the
	// incoming data.
	ErrAmbiguous = gobl.NewError("ambiguous-context")

	// ErrNotSupported is provided when none of the requested contexts accept
	// the envelope for export.
	ErrNotSupported = gobl.NewError("not-supported")

	// ErrConversion wraps errors returned by a converter. Errors that are
	// already GOBL errors, such as validation faults, are returned unchanged.
	ErrConversion = gobl.NewError("conversion")
)

// Context describes one external document format variant that a converter
// can import and/or export.
type Context struct {
	// Key uniquely identifies the context, e.g. "ubl+peppol".
	Key cbc.Key `json:"key" jsonschema:"title=Key"`
	// Name of the context.
	Name i18n.String `json:"name" jsonschema:"title=Name"`
	// Description of the context.
	Description i18n.String `json:"description,omitempty" jsonschema:"title=Description"`
	// MIME type of the context's documents, e.g. "application/xml".
	MIME string `json:"mime,omitempty" jsonschema:"title=MIME Type"`
	// Syntax is the base format the context is a variant of, e.g. "ubl" or "cii".
	Syntax cbc.Key `json:"syntax,omitempty" jsonschema:"title=Syntax"`
	// Countries where the context applies, including unions such as "EU".
	// Empty means no country restriction.
	Countries []l10n.Code `json:"countries,omitempty" jsonschema:"title=Countries"`
	// Addons the GOBL document must include to be exported into this context.
	Addons []cbc.Key `json:"addons,omitempty" jsonschema:"title=Addons"`
	// Import lists GOBL schemas this context can be converted into.
	Import []schema.ID `json:"import,omitempty" jsonschema:"title=Import Schemas"`
	// Export lists GOBL schemas that can be converted into this context.
	Export []schema.ID `json:"export,omitempty" jsonschema:"title=Export Schemas"`
}

// Converter is implemented by packages that convert between GOBL and one
// or more contexts.
type Converter interface {
	// Contexts lists the contexts handled by the converter.
	Contexts() []*Context
	// Detect returns the key of the context the input belongs to, or an
	// empty key if it is not recognized. It should not fully parse the data.
	Detect(in *Input) cbc.Key
	// Import converts data in the given context into a GOBL envelope.
	Import(key cbc.Key, data []byte) (*gobl.Envelope, error)
	// Accepts reports whether the envelope can be exported into the context.
	// It is only called once the document's schema is in the context's Export
	// list and the document includes the context's Addons.
	Accepts(key cbc.Key, env *gobl.Envelope) bool
	// Export converts the envelope into the given context.
	Export(key cbc.Key, env *gobl.Envelope) ([]byte, error)
}

// Output is the result of an export.
type Output struct {
	// Context the envelope was exported into.
	Context *Context
	// Data of the exported document.
	Data []byte
}

// Conversion describes a single 1:1 conversion available in the register.
type Conversion struct {
	// Context key of the external format.
	Context cbc.Key `json:"context" jsonschema:"title=Context"`
	// Schema of the GOBL document.
	Schema schema.ID `json:"schema" jsonschema:"title=Schema"`
	// Direction of the conversion, import into GOBL or export from it.
	Direction cbc.Key `json:"direction" jsonschema:"title=Direction"`
}

// exports checks if the envelope's document has one of the context's export
// schemas and includes all of its addons.
func (c *Context) exports(env *gobl.Envelope) bool {
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

// appliesTo checks if the context can be used in the country, either directly,
// through a union the country belongs to, or because it has no restriction.
func (c *Context) appliesTo(country l10n.Code) bool {
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

// Register adds the converter and its contexts to the global register. This
// is expected to be called from package init functions, and panics if a
// context key is empty or already registered.
func Register(c Converter) {
	converters.add(c)
}

// Contexts provides all the registered contexts, sorted by key.
func Contexts() []*Context {
	return converters.contexts()
}

// ContextFor provides the context registered with the key, or nil.
func ContextFor(key cbc.Key) *Context {
	return converters.contextFor(key)
}

// ContextsFor provides the contexts that apply to the country, directly or
// through a union it belongs to, along with those without a country restriction.
func ContextsFor(country l10n.Code) []*Context {
	return converters.contextsFor(country)
}

// Conversions lists every 1:1 conversion available from the registered
// contexts.
func Conversions() []*Conversion {
	return converters.conversions()
}

// Detect determines the context of the incoming data by asking each
// candidate converter. Keys, if provided, limit the candidates to those
// contexts.
func Detect(data []byte, keys ...cbc.Key) (*Context, error) {
	e, err := converters.detect(data, keys)
	if err != nil {
		return nil, err
	}
	return e.context, nil
}

// Import detects the context of the incoming data and converts it into a GOBL
// envelope. Keys, if provided, limit detection to those contexts.
func Import(data []byte, keys ...cbc.Key) (*gobl.Envelope, error) {
	e, err := converters.detect(data, keys)
	if err != nil {
		return nil, err
	}
	env, err := e.converter.Import(e.context.Key, data)
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
		return nil, ErrUnknownContext.WithReason("no contexts provided")
	}
	for _, k := range keys {
		if converters.entryFor(k) == nil {
			return nil, ErrUnknownContext.WithReason("context %s not registered", k)
		}
	}
	for _, k := range keys {
		e := converters.entryFor(k)
		if !e.context.exports(env) || !e.converter.Accepts(k, env) {
			continue
		}
		data, err := e.converter.Export(k, env)
		if err != nil {
			return nil, ErrConversion.WithCause(err)
		}
		return &Output{Context: e.context, Data: data}, nil
	}
	return nil, ErrNotSupported.WithReason("no context in %v accepts the envelope", keys)
}
