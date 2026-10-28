package convert_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/convert"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var invoiceSchema = schema.Lookup(bill.Invoice{})

// headerParses counts how many times the shared header is extracted.
var headerParses int

type headerKey struct{}

// header mimics a format package helper that extracts and caches the
// fields used to tell variants apart. Here that is the text before the
// first ":".
func header(in *convert.Input) string {
	if v, ok := in.Get(headerKey{}); ok {
		return v.(string)
	}
	headerParses++
	h, _, _ := strings.Cut(string(in.Data), ":")
	in.Set(headerKey{}, h)
	return h
}

// alphaConverter handles "alpha:" data, provides an attachment for data
// ending in ":attach", and no result for ":empty", and "alpha-fr:" for the FR variant
// that requires the "test-fr" addon on export. The plain format refuses
// documents with the "test-refuse" addon, and embeds attachments by adding
// their names to the output.
type alphaConverter struct{}

var alphaFRRules = convert.Rules("test+alpha-fr",
	rules.For(new(bill.Invoice),
		convert.RequireAddons("01", "test-fr"),
	),
)

func (alphaConverter) Formats() []*convert.Format {
	return []*convert.Format{
		{
			Key:         "test+alpha",
			Syntax:      "test",
			Import:      []schema.ID{invoiceSchema},
			Export:      []schema.ID{invoiceSchema},
			Attachments: true,
		},
		{
			Key:       "test+alpha-fr",
			Syntax:    "test",
			Countries: []l10n.Code{"FR"},
			Export:    []schema.ID{invoiceSchema},
		},
	}
}

func (alphaConverter) Detect(in *convert.Input) cbc.Key {
	switch header(in) {
	case "alpha":
		return "test+alpha"
	case "alpha-fr":
		return "test+alpha-fr"
	}
	return cbc.KeyEmpty
}

func (alphaConverter) Import(_ cbc.Key, data []byte, opts *convert.ImportOptions) (*convert.ImportResult, error) {
	switch {
	case strings.HasSuffix(string(data), ":fail"):
		return nil, errors.New("boom")
	case strings.HasSuffix(string(data), ":invalid"):
		return nil, gobl.ErrValidation.WithReason("bad document")
	case strings.HasSuffix(string(data), ":empty"):
		return nil, nil
	}
	content := string(data)
	for _, k := range opts.Formats {
		content += ":" + k.String()
	}
	env, err := gobl.Envelop(&note.Message{Content: content})
	if err != nil {
		return nil, err
	}
	res := &convert.ImportResult{Envelope: env}
	if strings.HasSuffix(string(data), ":attach") {
		res.Attachments = []*convert.Attachment{{Name: "invoice.pdf", Data: []byte("%PDF")}}
	}
	return res, nil
}

func (alphaConverter) Accepts(key cbc.Key, env *gobl.Envelope) error {
	inv := env.Extract().(*bill.Invoice)
	if key == "test+alpha-fr" {
		return alphaFRRules.Validate(inv)
	}
	if cbc.Key("test-refuse").In(inv.GetAddons()...) {
		return errors.New("alpha refuses")
	}
	return nil
}

func (alphaConverter) Export(key cbc.Key, _ *gobl.Envelope, opts *convert.ExportOptions) ([]byte, error) {
	out := key.String()
	for _, a := range opts.Attachments {
		out += ":" + a.Name
	}
	return []byte(out), nil
}

// betaConverter handles "beta:" data, and claims "alpha" data as gamma to
// produce an ambiguous match. Beta has no export schemas, and gamma
// requires the "test-gamma" addon on export.
type betaConverter struct{}

var gammaRules = convert.Rules("test+gamma",
	rules.For(new(bill.Invoice),
		convert.RequireAddons("01", "test-gamma"),
	),
)

func (betaConverter) Formats() []*convert.Format {
	return []*convert.Format{
		{
			Key:       "test+beta",
			Countries: []l10n.Code{l10n.EU},
			Import:    []schema.ID{invoiceSchema},
		},
		{
			Key:       "test+gamma",
			Countries: []l10n.Code{"US"},
			Import:    []schema.ID{invoiceSchema},
			Export:    []schema.ID{invoiceSchema},
		},
	}
}

func (betaConverter) Detect(in *convert.Input) cbc.Key {
	switch header(in) {
	case "beta":
		return "test+beta"
	case "alpha":
		return "test+gamma"
	}
	return cbc.KeyEmpty
}

func (betaConverter) Import(_ cbc.Key, data []byte, _ *convert.ImportOptions) (*convert.ImportResult, error) {
	env, err := gobl.Envelop(&note.Message{Content: string(data)})
	if err != nil {
		return nil, err
	}
	return &convert.ImportResult{Envelope: env}, nil
}

func (betaConverter) Accepts(_ cbc.Key, env *gobl.Envelope) error {
	return gammaRules.Validate(env.Extract())
}

func (betaConverter) Export(_ cbc.Key, _ *gobl.Envelope, _ *convert.ExportOptions) ([]byte, error) {
	return nil, errors.New("export failed")
}

func init() {
	convert.Register(alphaConverter{})
	convert.Register(betaConverter{})
}

func formatKeys(list []*convert.Format) []cbc.Key {
	keys := make([]cbc.Key, len(list))
	for i, c := range list {
		keys[i] = c.Key
	}
	return keys
}

func invoiceEnvelope(t *testing.T, addons ...cbc.Key) *gobl.Envelope {
	t.Helper()
	env := gobl.NewEnvelope()
	doc, err := schema.NewObject(&bill.Invoice{Addons: tax.WithAddons(addons...)})
	require.NoError(t, err)
	env.Document = doc
	return env
}

func TestFormats(t *testing.T) {
	assert.Equal(t,
		[]cbc.Key{"test+alpha", "test+alpha-fr", "test+beta", "test+gamma"},
		formatKeys(convert.Formats()),
	)
	assert.Equal(t, cbc.Key("test+beta"), convert.FormatFor("test+beta").Key)
	assert.Nil(t, convert.FormatFor("test+unknown"))
}

func TestFormatsFor(t *testing.T) {
	t.Run("country and union", func(t *testing.T) {
		assert.Equal(t,
			[]cbc.Key{"test+alpha", "test+alpha-fr", "test+beta"},
			formatKeys(convert.FormatsFor("FR")),
		)
	})
	t.Run("outside union", func(t *testing.T) {
		assert.Equal(t,
			[]cbc.Key{"test+alpha", "test+gamma"},
			formatKeys(convert.FormatsFor("US")),
		)
	})
	t.Run("no match", func(t *testing.T) {
		assert.Equal(t,
			[]cbc.Key{"test+alpha"},
			formatKeys(convert.FormatsFor("JP")),
		)
	})
}

func TestConversions(t *testing.T) {
	list := convert.Conversions()
	assert.Contains(t, list, &convert.Conversion{
		Format: "test+alpha", Schema: invoiceSchema, Direction: convert.DirectionImport,
	})
	assert.Contains(t, list, &convert.Conversion{
		Format: "test+alpha-fr", Schema: invoiceSchema, Direction: convert.DirectionExport,
	})
	assert.NotContains(t, list, &convert.Conversion{
		Format: "test+beta", Schema: invoiceSchema, Direction: convert.DirectionExport,
	})
	assert.Len(t, list, 6)
}

func TestDetect(t *testing.T) {
	t.Run("known", func(t *testing.T) {
		ctx, err := convert.Detect([]byte("beta:data"))
		require.NoError(t, err)
		assert.Equal(t, cbc.Key("test+beta"), ctx.Key)
	})
	t.Run("unknown", func(t *testing.T) {
		_, err := convert.Detect([]byte("other:data"))
		assert.ErrorIs(t, err, convert.ErrUnknownFormat)
	})
	t.Run("ambiguous", func(t *testing.T) {
		_, err := convert.Detect([]byte("alpha:data"))
		assert.ErrorIs(t, err, convert.ErrAmbiguous)
		assert.ErrorContains(t, err, "test+alpha and test+gamma")
	})
	t.Run("keys resolve ambiguity", func(t *testing.T) {
		ctx, err := convert.Detect([]byte("alpha:data"), "test+gamma")
		require.NoError(t, err)
		assert.Equal(t, cbc.Key("test+gamma"), ctx.Key)
	})
	t.Run("keys exclude a format of the same converter", func(t *testing.T) {
		_, err := convert.Detect([]byte("beta:data"), "test+gamma")
		assert.ErrorIs(t, err, convert.ErrUnknownFormat)
	})
	t.Run("unregistered key", func(t *testing.T) {
		_, err := convert.Detect([]byte("beta:data"), "test+unknown")
		assert.ErrorIs(t, err, convert.ErrUnknownFormat)
	})
	t.Run("shared values parsed once", func(t *testing.T) {
		headerParses = 0
		_, err := convert.Detect([]byte("beta:data"))
		require.NoError(t, err)
		assert.Equal(t, 1, headerParses)
	})
}

func TestImport(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		res, err := convert.Import([]byte("alpha-fr:data"))
		require.NoError(t, err)
		assert.Equal(t, cbc.Key("test+alpha-fr"), res.Format.Key)
		msg, ok := res.Envelope.Extract().(*note.Message)
		require.True(t, ok)
		assert.Equal(t, "alpha-fr:data", msg.Content)
		assert.Empty(t, res.Attachments)
	})
	t.Run("formats", func(t *testing.T) {
		res, err := convert.Import([]byte("alpha:data"), convert.WithFormats("test+alpha"))
		require.NoError(t, err)
		assert.Equal(t, cbc.Key("test+alpha"), res.Format.Key)
	})
	t.Run("formats passed to the converter", func(t *testing.T) {
		res, err := convert.Import([]byte("alpha-fr:data"), convert.WithFormats("test+alpha-fr", "test+beta"))
		require.NoError(t, err)
		msg := res.Envelope.Extract().(*note.Message)
		assert.Equal(t, "alpha-fr:data:test+alpha-fr:test+beta", msg.Content)
	})
	t.Run("formats not matching data", func(t *testing.T) {
		_, err := convert.Import([]byte("alpha:data"), convert.WithFormats("test+beta"))
		assert.ErrorIs(t, err, convert.ErrUnknownFormat)
	})
	t.Run("attachments", func(t *testing.T) {
		res, err := convert.Import([]byte("alpha-fr:attach"))
		require.NoError(t, err)
		require.Len(t, res.Attachments, 1)
		assert.Equal(t, "invoice.pdf", res.Attachments[0].Name)
		assert.Equal(t, []byte("%PDF"), res.Attachments[0].Data)
	})
	t.Run("no result", func(t *testing.T) {
		_, err := convert.Import([]byte("alpha-fr:empty"))
		assert.ErrorIs(t, err, convert.ErrConversion)
		assert.ErrorContains(t, err, "format test+alpha-fr provided no envelope")
	})
	t.Run("converter error", func(t *testing.T) {
		_, err := convert.Import([]byte("alpha-fr:fail"))
		assert.ErrorIs(t, err, convert.ErrConversion)
		assert.ErrorContains(t, err, "conversion: boom")
	})
	t.Run("gobl error unchanged", func(t *testing.T) {
		_, err := convert.Import([]byte("alpha-fr:invalid"))
		assert.ErrorIs(t, err, gobl.ErrValidation)
	})
}

func TestExport(t *testing.T) {
	t.Run("skips format missing addon", func(t *testing.T) {
		out, err := convert.Export(invoiceEnvelope(t), "test+alpha-fr", convert.WithAlternatives("test+alpha"))
		require.NoError(t, err)
		assert.Equal(t, cbc.Key("test+alpha"), out.Format.Key)
		assert.Equal(t, []byte("test+alpha"), out.Data)
	})
	t.Run("chooses format with addon", func(t *testing.T) {
		out, err := convert.Export(invoiceEnvelope(t, "test-fr"), "test+alpha-fr", convert.WithAlternatives("test+alpha"))
		require.NoError(t, err)
		assert.Equal(t, cbc.Key("test+alpha-fr"), out.Format.Key)
	})
	t.Run("faults", func(t *testing.T) {
		_, err := convert.Export(invoiceEnvelope(t), "test+alpha-fr")
		assert.ErrorIs(t, err, convert.ErrNotSupported)
		var ge *gobl.Error
		require.ErrorAs(t, err, &ge)
		fs := ge.Faults()
		require.NotNil(t, fs)
		assert.Equal(t, rules.Code("TEST-ALPHA-FR-BILL-INVOICE-01"), fs.First().Code())
		assert.Equal(t, []string{"$.$addons"}, fs.First().Paths())
		assert.Equal(t, "requires addon test-fr", fs.First().Message())
	})
	t.Run("faults combined", func(t *testing.T) {
		_, err := convert.Export(invoiceEnvelope(t), "test+beta", convert.WithAlternatives("test+gamma", "test+alpha-fr"))
		assert.ErrorIs(t, err, convert.ErrNotSupported)
		var ge *gobl.Error
		require.ErrorAs(t, err, &ge)
		fs := ge.Faults()
		require.NotNil(t, fs)
		assert.True(t, fs.HasCode("TEST-GAMMA-BILL-INVOICE-01"))
		assert.True(t, fs.HasCode("TEST-ALPHA-FR-BILL-INVOICE-01"))
	})
	t.Run("other reasons", func(t *testing.T) {
		_, err := convert.Export(invoiceEnvelope(t, "test-refuse"), "test+gamma", convert.WithAlternatives("test+alpha"))
		assert.ErrorIs(t, err, convert.ErrNotSupported)
		assert.ErrorContains(t, err, "alpha refuses")
		assert.ErrorContains(t, err, "TEST-GAMMA-BILL-INVOICE-01")
		var ge *gobl.Error
		require.ErrorAs(t, err, &ge)
		assert.Nil(t, ge.Faults())
	})
	t.Run("schema not exported", func(t *testing.T) {
		env, err := gobl.Envelop(&note.Message{Content: "hello"})
		require.NoError(t, err)
		_, err = convert.Export(env, "test+alpha")
		assert.ErrorIs(t, err, convert.ErrNotSupported)
		assert.ErrorContains(t, err, "no format in [test+alpha] exports the document")
	})
	t.Run("converter error", func(t *testing.T) {
		_, err := convert.Export(invoiceEnvelope(t, "test-gamma"), "test+gamma")
		assert.ErrorIs(t, err, convert.ErrConversion)
	})
	t.Run("unknown key", func(t *testing.T) {
		_, err := convert.Export(invoiceEnvelope(t), "test+alpha", convert.WithAlternatives("test+unknown"))
		assert.ErrorIs(t, err, convert.ErrUnknownFormat)
	})
	t.Run("empty key", func(t *testing.T) {
		_, err := convert.Export(invoiceEnvelope(t), "")
		assert.ErrorIs(t, err, convert.ErrUnknownFormat)
	})
	t.Run("alternatives in order", func(t *testing.T) {
		out, err := convert.Export(invoiceEnvelope(t), "test+alpha-fr",
			convert.WithAlternatives("test+gamma"),
			convert.WithAlternatives("test+alpha"),
		)
		require.NoError(t, err)
		assert.Equal(t, cbc.Key("test+alpha"), out.Format.Key)
	})
}

func TestExportAttachments(t *testing.T) {
	pdf := &convert.Attachment{Name: "invoice.pdf", MIME: "application/pdf", Data: []byte("%PDF")}
	csv := &convert.Attachment{Name: "lines.csv", MIME: "text/csv", Data: []byte("a,b")}
	t.Run("embedded", func(t *testing.T) {
		out, err := convert.Export(invoiceEnvelope(t), "test+alpha",
			convert.WithAttachments(pdf),
			convert.WithAttachments(csv),
		)
		require.NoError(t, err)
		assert.Equal(t, []byte("test+alpha:invoice.pdf:lines.csv"), out.Data)
	})
	t.Run("skips formats that cannot embed them", func(t *testing.T) {
		out, err := convert.Export(invoiceEnvelope(t, "test-fr"), "test+alpha-fr",
			convert.WithAlternatives("test+alpha"),
			convert.WithAttachments(pdf),
		)
		require.NoError(t, err)
		assert.Equal(t, cbc.Key("test+alpha"), out.Format.Key)
	})
	t.Run("refused", func(t *testing.T) {
		_, err := convert.Export(invoiceEnvelope(t, "test-fr"), "test+alpha-fr", convert.WithAttachments(pdf))
		assert.ErrorIs(t, err, convert.ErrNotSupported)
		assert.ErrorContains(t, err, "format test+alpha-fr does not embed attachments")
	})
}
