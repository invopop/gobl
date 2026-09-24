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

// alphaConverter handles "alpha:" data, and "alpha-fr:" for the FR variant
// that requires the "test-fr" addon on export.
type alphaConverter struct{}

func (alphaConverter) Contexts() []*convert.Context {
	return []*convert.Context{
		{
			Key:    "test+alpha",
			Syntax: "test",
			Import: []schema.ID{invoiceSchema},
			Export: []schema.ID{invoiceSchema},
		},
		{
			Key:       "test+alpha-fr",
			Syntax:    "test",
			Countries: []l10n.Code{"FR"},
			Addons:    []cbc.Key{"test-fr"},
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

func (alphaConverter) Import(_ cbc.Key, data []byte) (*gobl.Envelope, error) {
	if strings.HasSuffix(string(data), ":fail") {
		return nil, errors.New("boom")
	}
	return gobl.Envelop(&note.Message{Content: string(data)})
}

func (alphaConverter) Accepts(key cbc.Key, env *gobl.Envelope) bool {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return false
	}
	if key == "test+alpha-fr" {
		return cbc.Key("test-fr").In(inv.GetAddons()...)
	}
	return true
}

func (alphaConverter) Export(key cbc.Key, _ *gobl.Envelope) ([]byte, error) {
	return []byte(key.String()), nil
}

// betaConverter handles "beta:" data, and claims "alpha" data as gamma to
// produce an ambiguous match. Gamma accepts nothing on export.
type betaConverter struct{}

func (betaConverter) Contexts() []*convert.Context {
	return []*convert.Context{
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

func (betaConverter) Import(_ cbc.Key, data []byte) (*gobl.Envelope, error) {
	return gobl.Envelop(&note.Message{Content: string(data)})
}

func (betaConverter) Accepts(_ cbc.Key, _ *gobl.Envelope) bool {
	return false
}

func (betaConverter) Export(_ cbc.Key, _ *gobl.Envelope) ([]byte, error) {
	return nil, errors.New("not expected")
}

func init() {
	convert.Register(alphaConverter{})
	convert.Register(betaConverter{})
}

func contextKeys(list []*convert.Context) []cbc.Key {
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

func TestContexts(t *testing.T) {
	assert.Equal(t,
		[]cbc.Key{"test+alpha", "test+alpha-fr", "test+beta", "test+gamma"},
		contextKeys(convert.Contexts()),
	)
	assert.Equal(t, cbc.Key("test+beta"), convert.ContextFor("test+beta").Key)
	assert.Nil(t, convert.ContextFor("test+unknown"))
}

func TestContextsFor(t *testing.T) {
	t.Run("country and union", func(t *testing.T) {
		assert.Equal(t,
			[]cbc.Key{"test+alpha", "test+alpha-fr", "test+beta"},
			contextKeys(convert.ContextsFor("FR")),
		)
	})
	t.Run("outside union", func(t *testing.T) {
		assert.Equal(t,
			[]cbc.Key{"test+alpha", "test+gamma"},
			contextKeys(convert.ContextsFor("US")),
		)
	})
	t.Run("no match", func(t *testing.T) {
		assert.Equal(t,
			[]cbc.Key{"test+alpha"},
			contextKeys(convert.ContextsFor("JP")),
		)
	})
}

func TestConversions(t *testing.T) {
	list := convert.Conversions()
	assert.Contains(t, list, &convert.Conversion{
		Context: "test+alpha", Schema: invoiceSchema, Direction: convert.DirectionImport,
	})
	assert.Contains(t, list, &convert.Conversion{
		Context: "test+alpha-fr", Schema: invoiceSchema, Direction: convert.DirectionExport,
	})
	assert.NotContains(t, list, &convert.Conversion{
		Context: "test+beta", Schema: invoiceSchema, Direction: convert.DirectionExport,
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
		assert.ErrorIs(t, err, convert.ErrUnknownContext)
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
	t.Run("keys exclude a context of the same converter", func(t *testing.T) {
		_, err := convert.Detect([]byte("beta:data"), "test+gamma")
		assert.ErrorIs(t, err, convert.ErrUnknownContext)
	})
	t.Run("unregistered key", func(t *testing.T) {
		_, err := convert.Detect([]byte("beta:data"), "test+unknown")
		assert.ErrorIs(t, err, convert.ErrUnknownContext)
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
		env, err := convert.Import([]byte("alpha-fr:data"))
		require.NoError(t, err)
		msg, ok := env.Extract().(*note.Message)
		require.True(t, ok)
		assert.Equal(t, "alpha-fr:data", msg.Content)
	})
	t.Run("keys", func(t *testing.T) {
		env, err := convert.Import([]byte("alpha:data"), "test+alpha")
		require.NoError(t, err)
		assert.NotNil(t, env)
	})
	t.Run("keys not matching data", func(t *testing.T) {
		_, err := convert.Import([]byte("alpha:data"), "test+beta")
		assert.ErrorIs(t, err, convert.ErrUnknownContext)
	})
	t.Run("converter error", func(t *testing.T) {
		_, err := convert.Import([]byte("alpha-fr:fail"))
		assert.ErrorIs(t, err, convert.ErrConversion)
		assert.ErrorContains(t, err, "conversion: boom")
	})
}

func TestExport(t *testing.T) {
	t.Run("skips context missing addon", func(t *testing.T) {
		out, err := convert.Export(invoiceEnvelope(t), "test+alpha-fr", "test+alpha")
		require.NoError(t, err)
		assert.Equal(t, cbc.Key("test+alpha"), out.Context.Key)
		assert.Equal(t, []byte("test+alpha"), out.Data)
	})
	t.Run("chooses context with addon", func(t *testing.T) {
		out, err := convert.Export(invoiceEnvelope(t, "test-fr"), "test+alpha-fr", "test+alpha")
		require.NoError(t, err)
		assert.Equal(t, cbc.Key("test+alpha-fr"), out.Context.Key)
	})
	t.Run("not supported", func(t *testing.T) {
		_, err := convert.Export(invoiceEnvelope(t), "test+gamma", "test+alpha-fr")
		assert.ErrorIs(t, err, convert.ErrNotSupported)
	})
	t.Run("unknown key", func(t *testing.T) {
		_, err := convert.Export(invoiceEnvelope(t), "test+alpha", "test+unknown")
		assert.ErrorIs(t, err, convert.ErrUnknownContext)
	})
	t.Run("no keys", func(t *testing.T) {
		_, err := convert.Export(invoiceEnvelope(t))
		assert.ErrorIs(t, err, convert.ErrUnknownContext)
	})
}
