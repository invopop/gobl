package cfdi_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidComplement(t *testing.T) {
	fab := new(cfdi.FuelAccountBalance)

	err := rules.Validate(fab, withAddonContext())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "account number is required")
	assert.Contains(t, err.Error(), "lines are required")
}

func TestFuelAccountInvalidLine(t *testing.T) {
	fab := &cfdi.FuelAccountBalance{Lines: []*cfdi.FuelAccountLine{{}}}

	err := rules.Validate(fab, withAddonContext())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "line e-wallet ID is required")
	assert.Contains(t, err.Error(), "line purchase date and time is required")
	assert.Contains(t, err.Error(), "line vendor tax code is required")
	assert.Contains(t, err.Error(), "line service station code is required")
	assert.Contains(t, err.Error(), "line quantity must be greater than 0")
	assert.Contains(t, err.Error(), "line item is required")
	assert.Contains(t, err.Error(), "line purchase code is required")
	assert.Contains(t, err.Error(), "line taxes are required")

	fab.Lines[0].VendorTaxCode = "1234"
	fab.Lines[0].Quantity = num.MakeAmount(1, 0)
	fab.Lines[0].Item = &cfdi.FuelAccountItem{Price: num.MakeAmount(1, 0)}

	err = rules.Validate(fab, withAddonContext())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "line vendor tax identity code is invalid")
	assert.Contains(t, err.Error(), "line total must be quantity x unit_price")

	fab.Lines[0].VendorTaxCode = "K&A010301I16" // with symbols
	err = rules.Validate(fab, withAddonContext())
	assert.NotContains(t, err.Error(), "vendor_tax_code")
}

func TestInvalidItem(t *testing.T) {
	fab := &cfdi.FuelAccountBalance{Lines: []*cfdi.FuelAccountLine{
		{Item: &cfdi.FuelAccountItem{}}},
	}

	err := rules.Validate(fab, withAddonContext())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "line item type is required")
	assert.Contains(t, err.Error(), "line item name is required")
	assert.Contains(t, err.Error(), "line item price must be greater than 0")
}

func TestInvalidTax(t *testing.T) {
	fab := &cfdi.FuelAccountBalance{Lines: []*cfdi.FuelAccountLine{
		{Taxes: []*cfdi.FuelAccountTax{{}}}},
	}

	err := rules.Validate(fab, withAddonContext())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "tax category is required")
	assert.Contains(t, err.Error(), "tax rate is required when percent is not set")
	assert.Contains(t, err.Error(), "tax amount must be greater than 0")

	fab.Lines[0].Taxes[0].Category = "IRPF"

	err = rules.Validate(fab, withAddonContext())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "tax category must be a valid value")
}

func TestCalculate(t *testing.T) {
	t.Run("example 1", func(t *testing.T) {
		fab := &cfdi.FuelAccountBalance{
			Lines: []*cfdi.FuelAccountLine{
				{
					Quantity: num.MakeAmount(11, 1),
					Item:     &cfdi.FuelAccountItem{Price: num.MakeAmount(9091, 2)},
					Total:    num.MakeAmount(100, 0),
					Taxes: []*cfdi.FuelAccountTax{
						{
							Percent: num.NewPercentage(160, 3),
						},
						{
							Rate: num.NewAmount(56789, 4),
						},
					},
				},
				{
					Total: num.MakeAmount(100009, 3),
					Taxes: []*cfdi.FuelAccountTax{
						{
							Percent: num.NewPercentage(160, 3),
						},
						{
							Rate: num.NewAmount(56789, 4),
						},
					},
				},
			},
		}

		err := fab.Calculate()
		require.NoError(t, err)

		assert.Equal(t, "200.01", fab.Subtotal.String())
		assert.Equal(t, "238.26", fab.Total.String())

		assert.Equal(t, "1.100", fab.Lines[0].Quantity.String())
		assert.Equal(t, "90.910", fab.Lines[0].Item.Price.String())
		assert.Equal(t, "100.00", fab.Lines[0].Total.String())

		assert.Equal(t, "16.0%", fab.Lines[0].Taxes[0].Percent.String())
		assert.Equal(t, "16.00", fab.Lines[0].Taxes[0].Amount.String())
	})

	t.Run("example 2", func(t *testing.T) {
		fab := &cfdi.FuelAccountBalance{
			Lines: []*cfdi.FuelAccountLine{
				{
					Quantity: num.MakeAmount(9661, 3),
					Item:     &cfdi.FuelAccountItem{Price: num.MakeAmount(12743, 3)},
					Taxes: []*cfdi.FuelAccountTax{
						{
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(16, 2),
						},
						{
							Category: mx.TaxCategoryIEPS,
							Rate:     num.NewAmount(59195, 4),
						},
					},
				},
				{
					Quantity: num.MakeAmount(9680, 3),
					Item:     &cfdi.FuelAccountItem{Price: num.MakeAmount(12709, 3)},
					Taxes: []*cfdi.FuelAccountTax{
						{
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(16, 2),
						},
						{
							Category: mx.TaxCategoryIEPS,
							Rate:     num.NewAmount(59195, 4),
						},
					},
				},
			},
		}

		err := fab.Calculate()
		require.NoError(t, err)

		assert.Equal(t, "123.11", fab.Lines[0].Total.String())
		assert.Equal(t, "19.70", fab.Lines[0].Taxes[0].Amount.String())
		assert.Equal(t, "57.19", fab.Lines[0].Taxes[1].Amount.String())

		assert.Equal(t, "57.30", fab.Lines[1].Taxes[1].Amount.String())
		assert.Equal(t, "123.02", fab.Lines[1].Total.String())

		assert.Equal(t, "246.13", fab.Subtotal.String())
		assert.Equal(t, "400.00", fab.Total.String())
	})

	t.Run("example 3", func(t *testing.T) {
		fab := &cfdi.FuelAccountBalance{
			Lines: []*cfdi.FuelAccountLine{
				{
					Quantity: num.MakeAmount(525, 3),
					Item:     &cfdi.FuelAccountItem{Price: num.MakeAmount(19809, 3)},
					Taxes: []*cfdi.FuelAccountTax{
						{
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(16, 2),
						},
						{
							Category: mx.TaxCategoryIEPS,
							Rate:     num.NewAmount(5451, 4),
						},
					},
				},
				{
					Quantity: num.MakeAmount(304, 3),
					Item:     &cfdi.FuelAccountItem{Price: num.MakeAmount(19823, 3)},
					Taxes: []*cfdi.FuelAccountTax{
						{
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(16, 2),
						},
						{
							Category: mx.TaxCategoryIEPS,
							Rate:     num.NewAmount(5451, 4),
						},
					},
				},
			},
		}

		err := fab.Calculate()
		require.NoError(t, err)

		l0 := fab.Lines[0]
		assert.Equal(t, "10.40", l0.Total.String())
		assert.Equal(t, "1.66", l0.Taxes[0].Amount.String())
		assert.Equal(t, "0.29", l0.Taxes[1].Amount.String())
		assert.Equal(t, "12.35", l0.Total.Add(l0.Taxes[0].Amount).Add(l0.Taxes[1].Amount).String())

		assert.Equal(t, "0.17", fab.Lines[1].Taxes[1].Amount.String())
		assert.Equal(t, "6.03", fab.Lines[1].Total.String())

		assert.Equal(t, "16.43", fab.Subtotal.String())
		assert.Equal(t, "19.51", fab.Total.String())
	})

	t.Run("example 4", func(t *testing.T) {
		total := 12.35
		price := 23.774
		ieps := 0.5451
		vat := 0.16

		q := math.Round((total/price)*1000) / 1000
		ip := ((total / q) - ieps) / (1 + vat)

		fab := &cfdi.FuelAccountBalance{
			Lines: []*cfdi.FuelAccountLine{
				{
					Quantity: num.AmountFromFloat64(q, 3),
					Item:     &cfdi.FuelAccountItem{Price: num.AmountFromFloat64(ip, 4)},
					Taxes: []*cfdi.FuelAccountTax{
						{
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(int64(vat*1000), 3),
						},
						{
							Category: mx.TaxCategoryIEPS,
							Rate:     num.NewAmount(int64(ieps*10000), 4),
						},
					},
				},
			},
		}

		err := fab.Calculate()
		require.NoError(t, err)

		data, err := json.MarshalIndent(fab, "", "  ")
		require.NoError(t, err)
		exp := `
			{
				"account_number": "",
				"subtotal": "10.40",
				"total": "12.34",
				"lines": [
				  {
					"i": 1,
					"e_wallet_id": "",
					"purchase_date_time": "0000-00-00T00:00:00",
					"vendor_tax_code": "",
					"service_station_code": "",
					"quantity": "0.519",
					"item": {
					  "type": "",
					  "name": "",
					  "price": "20.0437"
					},
					"purchase_code": "",
					"total": "10.40",
					"taxes": [
					  {
						"cat": "VAT",
						"percent": "16.0%",
						"amount": "1.66"
					  },
					  {
						"cat": "IEPS",
						"rate": "0.5451",
						"amount": "0.28"
					  }
					]
				  }
				]
			  }
		`
		assert.JSONEq(t, exp, string(data))

		assert.NotEqual(t, total, fab.Total.Float64())
	})

	t.Run("example 5", func(t *testing.T) {
		total := 3832.93
		price := 23.774
		ieps := 5.451
		vat := 0.16

		q := math.Round((total/price)*1000) / 1000
		ip := ((total / q) - ieps) / (1 + vat)

		fab := &cfdi.FuelAccountBalance{
			Lines: []*cfdi.FuelAccountLine{
				{
					Quantity: num.AmountFromFloat64(q, 3),
					Item:     &cfdi.FuelAccountItem{Price: num.AmountFromFloat64(ip, 5)},
					Taxes: []*cfdi.FuelAccountTax{
						{
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(int64(vat*1000), 3),
						},
						{
							Category: mx.TaxCategoryIEPS,
							Rate:     num.NewAmount(int64(ieps*1000), 3),
						},
					},
				},
			},
		}

		err := fab.Calculate()
		require.NoError(t, err)

		data, err := json.MarshalIndent(fab, "", "  ")
		require.NoError(t, err)
		exp := `
			{
				"account_number": "",
				"subtotal": "2546.64",
				"total": "3832.93",
				"lines": [
				  {
					"i": 1,
					"e_wallet_id": "",
					"purchase_date_time": "0000-00-00T00:00:00",
					"vendor_tax_code": "",
					"service_station_code": "",
					"quantity": "161.224",
					"item": {
					  "type": "",
					  "name": "",
					  "price": "15.79564"
					},
					"purchase_code": "",
					"total": "2546.64",
					"taxes": [
					  {
						"cat": "VAT",
						"percent": "16.0%",
						"amount": "407.46"
					  },
					  {
						"cat": "IEPS",
						"rate": "5.451",
						"amount": "878.83"
					  }
					]
				  }
				]
			  }
		`
		assert.JSONEq(t, exp, string(data))

		assert.Equal(t, total, fab.Total.Float64())
	})
}
