package mx_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidComplement(t *testing.T) {
	fab := &mx.FuelAccountBalance{}

	err := fab.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "account_number: cannot be blank")
	assert.Contains(t, err.Error(), "lines: cannot be blank")
}

func TestInvalidLine(t *testing.T) {
	fab := &mx.FuelAccountBalance{Lines: []*mx.FuelAccountLine{{}}}

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
	fab.Lines[0].Item = &mx.FuelAccountItem{Price: num.MakeAmount(1, 0)}

	err = fab.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "vendor_tax_code: invalid tax identity code")
	assert.Contains(t, err.Error(), "total: must be quantity x unit_price")
}

func TestInvalidItem(t *testing.T) {
	fab := &mx.FuelAccountBalance{Lines: []*mx.FuelAccountLine{
		{Item: &mx.FuelAccountItem{}}},
	}

	err := fab.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "type: cannot be blank")
	assert.Contains(t, err.Error(), "name: cannot be blank")
	assert.Contains(t, err.Error(), "price: must be greater than 0")
}

func TestInvalidTax(t *testing.T) {
	fab := &mx.FuelAccountBalance{Lines: []*mx.FuelAccountLine{
		{Taxes: []*mx.FuelAccountTax{{}}}},
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
		fab := &mx.FuelAccountBalance{
			Lines: []*mx.FuelAccountLine{
				{
					Quantity: num.MakeAmount(11, 1),
					Item:     &mx.FuelAccountItem{Price: num.MakeAmount(9091, 2)},
					Total:    num.MakeAmount(100, 0),
					Taxes: []*mx.FuelAccountTax{
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
					Taxes: []*mx.FuelAccountTax{
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
		fab := &mx.FuelAccountBalance{
			Lines: []*mx.FuelAccountLine{
				{
					Quantity: num.MakeAmount(9661, 3),
					Item:     &mx.FuelAccountItem{Price: num.MakeAmount(12743, 3)},
					Taxes: []*mx.FuelAccountTax{
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
					Item:     &mx.FuelAccountItem{Price: num.MakeAmount(12709, 3)},
					Taxes: []*mx.FuelAccountTax{
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
		fab := &mx.FuelAccountBalance{
			Lines: []*mx.FuelAccountLine{
				{
					Quantity: num.MakeAmount(525, 3),
					Item:     &mx.FuelAccountItem{Price: num.MakeAmount(19809, 3)},
					Taxes: []*mx.FuelAccountTax{
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
					Item:     &mx.FuelAccountItem{Price: num.MakeAmount(19823, 3)},
					Taxes: []*mx.FuelAccountTax{
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

		fab := &mx.FuelAccountBalance{
			Lines: []*mx.FuelAccountLine{
				{
					Quantity: num.MakeAmount(int64(q*1000), 3),
					Item:     &mx.FuelAccountItem{Price: num.MakeAmount(int64(ip*10000), 4)},
					Taxes: []*mx.FuelAccountTax{
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
				"total": "12.35",
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
					  "price": "20.0436"
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

		assert.Equal(t, total, fab.Total.Float64())
	})

	t.Run("example 5", func(t *testing.T) {
		// reverse calculate the item price based on the expected total and
		// price per litre on the day.

		total := 3832.93
		price := 23.774
		ieps := 0.5451
		vat := 0.16

		q := math.Round((total/price)*1000) / 1000
		ip := ((total / q) - ieps) / (1 + vat)

		fab := &mx.FuelAccountBalance{
			Lines: []*mx.FuelAccountLine{
				{
					Quantity: num.MakeAmount(int64(q*1000), 3),
					// This case needs 5 decimal places to work due to large quantity:
					Item: &mx.FuelAccountItem{Price: num.MakeAmount(int64(ip*100000), 5)},
					Taxes: []*mx.FuelAccountTax{
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
				"subtotal": "3228.49",
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
					  "price": "20.02486"
					},
					"purchase_code": "",
					"total": "3228.49",
					"taxes": [
					  {
						"cat": "VAT",
						"percent": "16.0%",
						"amount": "516.56"
					  },
					  {
						"cat": "IEPS",
						"rate": "0.5451",
						"amount": "87.88"
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
