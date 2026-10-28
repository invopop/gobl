// Package convert provides a central register of converters between GOBL
// and external document formats.
//
// GOBL itself does not convert. Packages that do, such as gobl.ubl or
// gobl.cii, register a Converter from their init functions so that importing
// them makes their formats available here:
//
//	import _ "github.com/invopop/gobl.ubl"
//
//	res, err := convert.Import(data)
//	out, err := convert.Export(env, "ubl+peppol+bis-v3",
//		convert.WithAlternatives("ubl+en16931-v2017"),
//	)
package convert

import (
	"errors"
	"fmt"
	"slices"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/rules"
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
	// Key uniquely identifies the format, e.g. "ubl+peppol+bis-v3".
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
	// Import lists GOBL schemas this format can be converted into.
	Import []schema.ID `json:"import,omitempty" jsonschema:"title=Import Schemas"`
	// Export lists GOBL schemas that can be converted into this format.
	Export []schema.ID `json:"export,omitempty" jsonschema:"title=Export Schemas"`
	// Attachments is true when exported documents can embed attachments.
	Attachments bool `json:"attachments,omitempty" jsonschema:"title=Attachments"`
}

// Converter is implemented by packages that convert between GOBL and one
// or more formats.
type Converter interface {
	// Formats lists the formats handled by the converter.
	Formats() []*Format
	// Detect returns the key of the format the input belongs to, or an
	// empty key if it is not recognized. It should not fully parse the data.
	Detect(in *Input) cbc.Key
	// Import converts data in the given format into a GOBL envelope, and
	// provides the files embedded in it as attachments. Converters that
	// support routing should set the options' From and To on the envelope's
	// head before calculating it. The register sets the result's Format, and
	// any routing the converter did not.
	Import(key cbc.Key, data []byte, opts *ImportOptions) (*ImportResult, error)
	// Accepts checks the envelope can be exported into the format, and
	// otherwise returns the reason, preferably as rules faults from the
	// format's Rules. It is only called once the document's schema is in the
	// format's Export list.
	Accepts(key cbc.Key, env *gobl.Envelope) error
	// Export converts the envelope into the given format, embedding the
	// options' attachments, which are only provided to formats that
	// declare Attachments.
	Export(key cbc.Key, env *gobl.Envelope, opts *ExportOptions) ([]byte, error)
}

// Attachment is a file embedded in a document, such as a PDF version of an
// invoice, given to embed on export and provided on import.
type Attachment struct {
	// ID identifies the attachment in the document.
	ID string `json:"id,omitempty" jsonschema:"title=ID"`
	// Name of the file, such as "invoice.pdf".
	Name string `json:"name,omitempty" jsonschema:"title=Name"`
	// Description of what the file contains.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// MIME type of the file, such as "application/pdf".
	MIME string `json:"mime,omitempty" jsonschema:"title=MIME Type"`
	// Data of the file, encoded in base64 in JSON.
	Data []byte `json:"data" jsonschema:"title=Data"`
}

// ExportOptions configure an export.
type ExportOptions struct {
	// Alternatives are further format keys to try, in order of preference,
	// when the requested format does not take the envelope.
	Alternatives []cbc.Key
	// Attachments to embed in the exported document.
	Attachments []*Attachment
}

// ExportOption sets an export option.
type ExportOption func(*ExportOptions)

// WithAlternatives adds format keys to try, in order of preference, when the
// requested format does not take the envelope.
func WithAlternatives(keys ...cbc.Key) ExportOption {
	return func(o *ExportOptions) {
		o.Alternatives = append(o.Alternatives, keys...)
	}
}

// WithAttachments adds files to embed in the exported document. Formats
// that do not declare Attachments refuse the envelope.
func WithAttachments(atts ...*Attachment) ExportOption {
	return func(o *ExportOptions) {
		o.Attachments = append(o.Attachments, atts...)
	}
}

// ExportResult is the result of an export.
type ExportResult struct {
	// Format the envelope was exported into.
	Format *Format
	// Data of the exported document.
	Data []byte
}

// ImportOptions configure an import.
type ImportOptions struct {
	// Formats limit detection to the formats with these keys.
	Formats []cbc.Key
	// From and To are the transport addresses the document was received
	// with, such as Peppol participant URIs, recorded on the envelope's
	// head.
	From cbc.URI
	To   cbc.URI
}

// ImportOption sets an import option.
type ImportOption func(*ImportOptions)

// WithFormats limits detection to the formats with the keys.
func WithFormats(keys ...cbc.Key) ImportOption {
	return func(o *ImportOptions) {
		o.Formats = append(o.Formats, keys...)
	}
}

// WithRouting sets the transport addresses the document was received with,
// recorded verbatim on the envelope's Head.From and Head.To instead of those
// GOBL would derive from the document, which are those of an outgoing one.
func WithRouting(from, to cbc.URI) ImportOption {
	return func(o *ImportOptions) {
		o.From = from
		o.To = to
	}
}

// ImportResult is the result of an import.
type ImportResult struct {
	// Format the data was detected as.
	Format *Format
	// Envelope with the imported document.
	Envelope *gobl.Envelope
	// Attachments embedded in the data, such as a PDF version of the
	// document, which are not included in the envelope.
	Attachments []*Attachment
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
// schemas.
func (c *Format) exports(env *gobl.Envelope) bool {
	return env != nil && env.Document != nil && slices.Contains(c.Export, env.Document.Schema)
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
// envelope, providing any files embedded in the data as attachments.
func Import(data []byte, opts ...ImportOption) (*ImportResult, error) {
	o := new(ImportOptions)
	for _, opt := range opts {
		opt(o)
	}
	e, err := converters.detect(data, o.Formats)
	if err != nil {
		return nil, err
	}
	res, err := e.converter.Import(e.format.Key, data, o)
	if err != nil {
		return nil, ErrConversion.WithCause(err)
	}
	if res == nil || res.Envelope == nil {
		return nil, ErrConversion.WithReason("format %s provided no envelope", e.format.Key)
	}
	res.Format = e.format
	setRouting(res.Envelope, o)
	return res, nil
}

// setRouting records the requested routing on the envelope, for converters
// that do not. From and To are not covered by the envelope's signatures, so
// they may be set once the envelope is calculated.
func setRouting(env *gobl.Envelope, o *ImportOptions) {
	if o.From == "" && o.To == "" {
		return
	}
	if env.Head == nil {
		env.Head = head.NewHeader()
	}
	if o.From != "" {
		env.Head.From = o.From
	}
	if o.To != "" {
		env.Head.To = o.To
	}
}

// Export converts the envelope into the format with the key or, when it does
// not take the envelope, the first of the alternatives that does. A format
// takes the envelope when it exports the document's schema, can embed any
// attachments, and its converter accepts it. When none does, the error
// includes the reasons given by each converter, with their rules faults
// combined.
func Export(env *gobl.Envelope, key cbc.Key, opts ...ExportOption) (*ExportResult, error) {
	o := new(ExportOptions)
	for _, opt := range opts {
		opt(o)
	}
	keys := append([]cbc.Key{key}, o.Alternatives...)
	for _, k := range keys {
		if converters.entryFor(k) == nil {
			return nil, ErrUnknownFormat.WithReason("format %s not registered", k)
		}
	}
	var refusals []error
	for _, k := range keys {
		e := converters.entryFor(k)
		if !e.format.exports(env) {
			continue
		}
		if len(o.Attachments) > 0 && !e.format.Attachments {
			refusals = append(refusals, fmt.Errorf("format %s does not embed attachments", k))
			continue
		}
		if err := e.converter.Accepts(k, env); err != nil {
			refusals = append(refusals, err)
			continue
		}
		data, err := e.converter.Export(k, env, o)
		if err != nil {
			return nil, ErrConversion.WithCause(err)
		}
		return &ExportResult{Format: e.format, Data: data}, nil
	}
	if len(refusals) == 0 {
		return nil, ErrNotSupported.WithReason("no format in %v exports the document", keys)
	}
	return nil, ErrNotSupported.WithCause(joinRefusals(refusals))
}

// joinRefusals combines the reasons converters gave for refusing an envelope,
// keeping them as rules faults when they all are.
func joinRefusals(errs []error) error {
	list := make([]rules.Faults, 0, len(errs))
	for _, err := range errs {
		fs, ok := err.(rules.Faults)
		if !ok {
			return errors.Join(errs...)
		}
		list = append(list, fs)
	}
	return rules.JoinFaults(list...)
}
