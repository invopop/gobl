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
}

func TestTaxCategories(t *testing.T) {
	t.Parallel()
	regime := pe.New()

	require.Len(t, regime.Categories, 1)
	cat := regime.Categories[0]

	assert.Equal(t, tax.CategoryVAT, cat.Code)
	assert.Equal(t, "IGV", cat.Name[i18n.ES])
	assert.NotEmpty(t, cat.Sources)
	require.Len(t, cat.Rates, 2)

	// General rate: 18% (16% IGV + 2% IPM) since March 2011 (Ley 29666),
	// 19% (17% + 2%) before that.
	general := cat.Rates[0]
	assert.Equal(t, tax.RateGeneral, general.Rate)
	require.Len(t, general.Values, 2)
	assert.Equal(t, num.MakePercentage(180, 3), general.Values[0].Percent)
	assert.Equal(t, cal.NewDate(2011, 3, 1), general.Values[0].Since)
	assert.Equal(t, num.MakePercentage(190, 3), general.Values[1].Percent)
	assert.Equal(t, cal.NewDate(2003, 8, 1), general.Values[1].Since)

	// Reduced rate for qualifying MYPEs in restaurants, hotels and tourist
	// accommodation (Ley 31556, extended and modified by Ley 32219; IPM
	// component raised by Ley 32387): 10% from Sept 2022, 10.5% during
	// 2026, 15% during 2027 (the special regime's final year).
	reduced := cat.Rates[1]
	assert.Equal(t, tax.RateReduced, reduced.Rate)
	require.Len(t, reduced.Values, 3)
	assert.Equal(t, num.MakePercentage(150, 3), reduced.Values[0].Percent)
	assert.Equal(t, cal.NewDate(2027, 1, 1), reduced.Values[0].Since)
	assert.Equal(t, num.MakePercentage(105, 3), reduced.Values[1].Percent)
	assert.Equal(t, cal.NewDate(2026, 1, 1), reduced.Values[1].Since)
	assert.Equal(t, num.MakePercentage(100, 3), reduced.Values[2].Percent)
	assert.Equal(t, cal.NewDate(2022, 9, 1), reduced.Values[2].Since)
}

func TestCorrections(t *testing.T) {
	t.Parallel()
	regime := pe.New()

	// Peruvian law only defines the "nota de crédito" (SUNAT Catalogue 09)
	// and "nota de débito" (Catalogue 10) as correction documents.
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
