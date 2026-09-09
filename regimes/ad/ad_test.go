package ad_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/ad"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	r := ad.New()
	assert.Equal(t, l10n.AD.Tax(), r.Country)
	assert.Equal(t, tax.CategoryVAT, r.TaxScheme)
}

func TestInvoiceValidation(t *testing.T) {
	_ = ad.New() // Ensure the package is initialized
	inv := &bill.Invoice{
		Regime: tax.WithRegime(l10n.AD.Tax()),
		Code:   "123",
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: l10n.AD.Tax(),
				Code:    "L123456A",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: l10n.AD.Tax(),
				Code:    "F121212B",
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.NewAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: tax.CategoryVAT,
						Rate:     tax.RateGeneral,
					},
				},
			},
		},
	}

	require.NoError(t, inv.Calculate())
	require.NoError(t, inv.Validate())
}

func TestInvoiceValidationNilSupplier(t *testing.T) {
	_ = ad.New()
	inv := &bill.Invoice{
		Regime: tax.WithRegime(l10n.AD.Tax()),
	}
	require.NoError(t, ad.Validate(inv))
}

func TestTaxCategories(t *testing.T) {
	r := ad.New()
	cat := r.Categories[0]
	assert.Equal(t, tax.CategoryVAT, cat.Code)

	// Test Standard Rate
	rate := cat.RateDef(tax.KeyStandard, tax.RateGeneral)
	require.NotNil(t, rate)
	val := rate.Value(cal.Today(), nil)
	require.NotNil(t, val)
	assert.Equal(t, "4.5%", val.Percent.String())

	// Test Increased Rate
	rate = cat.RateDef(tax.KeyStandard, ad.RateIncreased)
	require.NotNil(t, rate)
	val = rate.Value(cal.Today(), nil)
	require.NotNil(t, val)
	assert.Equal(t, "9.5%", val.Percent.String())
}

func TestValidate(t *testing.T) {
	_ = ad.New()
	tests := []struct {
		name string
		doc  any
	}{
		{
			name: "identity",
			doc: &tax.Identity{
				Country: l10n.AD.Tax(),
				Code:    "F123456A",
			},
		},
		{
			name: "other",
			doc:  &org.Party{Name: "Test"},
		},
		{
			name: "nil",
			doc:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ad.Validate(tt.doc)
			assert.NoError(t, err)
		})
	}
}

func TestNormalize(t *testing.T) {
	_ = ad.New()
	t.Run("identity", func(t *testing.T) {
		tID := &tax.Identity{
			Country: l10n.AD.Tax(),
			Code:    " f123456a ",
		}
		ad.Normalize(tID)
		assert.Equal(t, cbc.Code("F123456A"), tID.Code)
	})
	t.Run("other", func(t *testing.T) {
		p := &org.Party{Name: "Test"}
		ad.Normalize(p)
		assert.Equal(t, "Test", p.Name)
	})
}
