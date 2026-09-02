package pe_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/pe"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()
	regime := pe.New()
	require.NotNil(t, regime)

	assert.Equal(t, l10n.PE, regime.Country.Code())
	assert.Equal(t, currency.PEN, regime.Currency)
	assert.Equal(t, tax.CategoryVAT, regime.TaxScheme)
	assert.Equal(t, "Peru", regime.Name.String())
	assert.Equal(t, "Perú", regime.Name[i18n.ES])
	assert.Equal(t, "America/Lima", regime.TimeZone)
	assert.NotEmpty(t, regime.Description)
	assert.NotEmpty(t, regime.Sources)
}

func TestTaxCategories(t *testing.T) {
	t.Parallel()
	regime := pe.New()

	require.Len(t, regime.Categories, 1)
	cat := regime.Categories[0]

	assert.Equal(t, tax.CategoryVAT, cat.Code)
	assert.Equal(t, "IGV", cat.Name[i18n.ES])
	assert.NotEmpty(t, cat.Sources)
	require.Len(t, cat.Rates, 1)

	// General rate: 18% (16% IGV + 2% IPM) since March 2011 (Ley 29666),
	// 19% (17% + 2%) before that.
	general := cat.Rates[0]
	assert.Equal(t, tax.RateGeneral, general.Rate)
	require.Len(t, general.Values, 2)
	assert.Equal(t, num.MakePercentage(180, 3), general.Values[0].Percent)
	assert.Equal(t, cal.NewDate(2011, 3, 1), general.Values[0].Since)
	assert.Equal(t, num.MakePercentage(190, 3), general.Values[1].Percent)
	assert.Equal(t, cal.NewDate(2003, 8, 1), general.Values[1].Since)
}

func TestSources(t *testing.T) {
	t.Parallel()
	regime := pe.New()
	require.Len(t, regime.Categories, 1)

	assert.ElementsMatch(t, []string{
		"https://centrovirtual.sunat.gob.pe/tramites/inscribete-ruc",
		"https://www.sunat.gob.pe/legislacion/superin/2021/anexo-026-2021.pdf",
		"https://www.sunat.gob.pe/legislacion/comprob/regla/capituloIII.pdf",
	}, sourceURLs(t, regime.Sources))
	assert.ElementsMatch(t, []string{
		"https://www.sunat.gob.pe/legislacion/igv/ley/",
		"https://orientacion.sunat.gob.pe/3053-concepto-tasa-y-operaciones-gravadas-igv-empresas",
	}, sourceURLs(t, regime.Categories[0].Sources))
}

func TestRateDateResolution(t *testing.T) {
	t.Parallel()
	regime := pe.New()
	require.Len(t, regime.Categories, 1)
	require.Len(t, regime.Categories[0].Rates, 1)
	general := regime.Categories[0].Rates[0]

	tests := []struct {
		name    string
		rate    *tax.RateDef
		date    cal.Date
		percent num.Percentage
	}{
		{name: "general before 2011 change", rate: general, date: cal.MakeDate(2011, 2, 28), percent: num.MakePercentage(190, 3)},
		{name: "general from 2011 change", rate: general, date: cal.MakeDate(2011, 3, 1), percent: num.MakePercentage(180, 3)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rv := tt.rate.Value(tt.date, tax.Extensions{})
			require.NotNil(t, rv)
			assert.Equal(t, tt.percent, rv.Percent)
		})
	}
}

func TestCorrections(t *testing.T) {
	t.Parallel()
	regime := pe.New()

	// Article 10 of the Payment Voucher Regulations covers credit and debit
	// notes as documents that modify previously issued payment documents.
	require.Len(t, regime.Corrections, 1)
	correction := regime.Corrections[0]
	assert.Equal(t, bill.ShortSchemaInvoice, correction.Schema)
	assert.Equal(t, []cbc.Key{
		bill.InvoiceTypeCreditNote,
		bill.InvoiceTypeDebitNote,
	}, correction.Types)
}

func TestIdentities(t *testing.T) {
	t.Parallel()
	regime := pe.New()

	require.Len(t, regime.Identities, 3)
	codes := make([]cbc.Code, len(regime.Identities))
	for i, def := range regime.Identities {
		codes[i] = def.Code
	}
	assert.Equal(t, []cbc.Code{
		pe.IdentityTypeDNI,
		pe.IdentityTypeCE,
		pe.IdentityTypePassport,
	}, codes)
}

func sourceURLs(t *testing.T, sources []*cbc.Source) []string {
	t.Helper()
	urls := make([]string, len(sources))
	for i, source := range sources {
		require.NotNil(t, source)
		urls[i] = source.URL
	}
	return urls
}
