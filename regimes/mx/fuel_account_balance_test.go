package mx_test

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidComplement(t *testing.T) {
	comp := &mx.FuelAccountBalance{}

	err := comp.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "account_number: cannot be blank")
	assert.Contains(t, err.Error(), "lines: cannot be blank")
}

func TestInvalidLine(t *testing.T) {
	comp := &mx.FuelAccountBalance{Lines: []*mx.FuelAccountLine{{}}}

	err := comp.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "e_wallet_id: cannot be blank")
	assert.Contains(t, err.Error(), "purchase_date_time: required")
	assert.Contains(t, err.Error(), "vendor_tax_code: cannot be blank")
	assert.Contains(t, err.Error(), "service_station_code: cannot be blank")
	assert.Contains(t, err.Error(), "quantity: must be greater than 0")
	assert.Contains(t, err.Error(), "fuel_type: cannot be blank")
	assert.Contains(t, err.Error(), "fuel_name: cannot be blank")
	assert.Contains(t, err.Error(), "purchase_code: cannot be blank")
	assert.Contains(t, err.Error(), "unit_price: must be greater than 0")
	assert.Contains(t, err.Error(), "taxes: cannot be blank")

	comp.Lines[0].VendorTaxCode = "1234"
	comp.Lines[0].Quantity = num.MakeAmount(1, 0)
	comp.Lines[0].UnitPrice = num.MakeAmount(1, 0)

	err = comp.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "vendor_tax_code: invalid tax identity code")
	assert.Contains(t, err.Error(), "total: must be quantity x unit_price")
}

func TestInvalidTax(t *testing.T) {
	comp := &mx.FuelAccountBalance{Lines: []*mx.FuelAccountLine{
		{Taxes: []*mx.FuelAccountTax{{}}}},
	}

	err := comp.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "code: cannot be blank")
	assert.Contains(t, err.Error(), "rate: must be greater than 0")
	assert.Contains(t, err.Error(), "amount: must be greater than 0")

	comp.Lines[0].Taxes[0].Code = "IRPF"

	err = comp.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "code: must be a valid value")
}

func TestCalculate(t *testing.T) {
	comp := &mx.FuelAccountBalance{
		Lines: []*mx.FuelAccountLine{
			{
				Quantity:  num.MakeAmount(11, 1),
				UnitPrice: num.MakeAmount(9091, 2),
				Total:     num.MakeAmount(100, 0),
				Taxes: []*mx.FuelAccountTax{
					{
						Rate:   num.MakeAmount(16, 2),
						Amount: num.MakeAmount(16, 0),
					},
					{Amount: num.MakeAmount(56789, 4)},
				},
			},
			{
				Total: num.MakeAmount(100009, 3),
				Taxes: []*mx.FuelAccountTax{
					{Amount: num.MakeAmount(16, 0)},
					{Amount: num.MakeAmount(56789, 4)},
				},
			},
		},
	}

	err := comp.Calculate()

	require.NoError(t, err)
	assert.Equal(t, num.MakeAmount(20001, 2), comp.Subtotal)
	assert.Equal(t, num.MakeAmount(24337, 2), comp.Total)

	assert.Equal(t, num.MakeAmount(1100, 3), comp.Lines[0].Quantity)
	assert.Equal(t, num.MakeAmount(90910, 3), comp.Lines[0].UnitPrice)
	assert.Equal(t, num.MakeAmount(10000, 2), comp.Lines[0].Total)

	assert.Equal(t, num.MakeAmount(160000, 6), comp.Lines[0].Taxes[0].Rate)
	assert.Equal(t, num.MakeAmount(1600, 2), comp.Lines[0].Taxes[0].Amount)
}
