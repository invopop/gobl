package cfdi_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/mx/sat"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidComplement(t *testing.T) {
	fab := new(cfdi.FuelAccountBalance)

	err := fab.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "account_number: cannot be blank")
	assert.Contains(t, err.Error(), "lines: cannot be blank")
}

func TestInvalidLine(t *testing.T) {
	fab := &cfdi.FuelAccountBalance{Lines: []*cfdi.FuelAccountLine{{}}}

	err := fab.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "e_wallet_id: cannot be blank")
	assert.Contains(t, err.Error(), "purchase_date_time: required")
	assert.Contains(t, err.Error(), "vendor_tax_code: cannot be blank")
	assert.Contains(t, err.Error(), "service_station_code: cannot be blank")
	assert.Contains(t, err.Error(), "quantity: must be greater than 0")
	assert.Contains(t, err.Error(), "item: cannot be blank")
	assert.Contains(t, err.Error(), "purchase_code: cannot be blank")
	assert.Contains(t, err.Error(), "taxes: cannot be blank")

	fab.Lines[0].VendorTaxCode = "1234"
	fab.Lines[0].Quantity = num.MakeAmount(1, 0)
	fab.Lines[0].Item = &cfdi.FuelAccountItem{Price: num.MakeAmount(1, 0)}

	err = fab.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "vendor_tax_code: invalid tax identity code")
	assert.Contains(t, err.Error(), "total: must be quantity x unit_price")
}

func TestInvalidItem(t *testing.T) {
	fab := &cfdi.FuelAccountBalance{Lines: []*cfdi.FuelAccountLine{
		{Item: &cfdi.FuelAccountItem{}}},
	}

	err := fab.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "type: cannot be blank")
	assert.Contains(t, err.Error(), "name: cannot be blank")
	assert.Contains(t, err.Error(), "price: must be greater than 0")
}

func TestInvalidTax(t *testing.T) {
	fab := &cfdi.FuelAccountBalance{Lines: []*cfdi.FuelAccountLine{
		{Taxes: []*cfdi.FuelAccountTax{{}}}},
	}

	err := fab.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cat: cannot be blank")
	assert.Contains(t, err.Error(), "rate: cannot be blank")
	assert.Contains(t, err.Error(), "amount: must be greater than 0")

	fab.Lines[0].Taxes[0].Category = "IRPF"

	err = fab.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cat: must be a valid value")
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
							Category: sat.TaxCategoryIEPS,
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
							Category: sat.TaxCategoryIEPS,
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
							Category: sat.TaxCategoryIEPS,
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
							Category: sat.TaxCategoryIEPS,
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
		// reverse calculate the item price based on the expected total and
		// price per litre on the day.

		total := 12.35
		price := 23.774
		ieps := 0.5451
		vat := 0.16

		// quantity = total_with_taxes / precio_litro_with_tax
		// item_price = ((total_with_taxes / quantity) - ieps_rate) / (1 + vat_rate)
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
							Category: sat.TaxCategoryIEPS,
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

		// Loosing precision means that this calculation does not
		// work.
		assert.NotEqual(t, total, fab.Total.Float64())
	})

	t.Run("example 5", func(t *testing.T) {
		// reverse calculate the item price based on the expected total and
		// price per litre on the day.

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
					// This case needs 5 decimal places to work due to large quantity:
					Item: &cfdi.FuelAccountItem{Price: num.AmountFromFloat64(ip, 5)},
					Taxes: []*cfdi.FuelAccountTax{
						{
							Category: tax.CategoryVAT,
							Percent:  num.NewPercentage(int64(vat*1000), 3),
						},
						{
							Category: sat.TaxCategoryIEPS,
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
