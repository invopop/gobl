package bill_test

import (
	"testing"

	_ "github.com/invopop/gobl" // load regions
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regions/common"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveIncludedTax(t *testing.T) {
	i := &bill.Invoice{
		Code: "123TEST",
		Tax: &bill.Tax{
			PricesInclude: common.TaxCategoryVAT,
		},
		Supplier: &org.Party{
			TaxID: &org.TaxIdentity{
				Country: l10n.ES,
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			TaxID: &org.TaxIdentity{
				Country: l10n.ES,
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: bill.Lines{
			{
				Quantity: num.MakeAmount(10, 0),
				Item: &org.Item{
					Name:  "Test Item",
					Price: num.MakeAmount(10000, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Rate:     "standard",
					},
				},
				Discounts: []*bill.LineDiscount{
					{
						Reason:  "Testing",
						Percent: num.NewPercentage(10, 2),
					},
				},
			},
		},
	}

	require.NoError(t, i.Calculate())

	i2 := i.RemoveIncludedTaxes(2)

	assert.Empty(t, i2.Tax.PricesInclude)
	l0 := i2.Lines[0]
	assert.Equal(t, "82.6446", l0.Item.Price.String())
	assert.Equal(t, "826.4463", l0.Sum.String())
	assert.Equal(t, "826.4463", l0.Sum.String())
	assert.Equal(t, "82.6446", l0.Discounts[0].Amount.String())
	assert.Equal(t, "743.8017", l0.Total.String())

	assert.Equal(t, "100.00", i.Lines[0].Item.Price.String())
}
