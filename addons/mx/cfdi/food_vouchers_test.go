package cfdi_test

import (
	"testing"

	"github.com/invopop/gobl/addons/mx/cfdi"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidFoodVouchers(t *testing.T) {
	fvc := &cfdi.FoodVouchers{}

	err := rules.Validate(fvc, withAddonContext())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "account number is required")
	assert.Contains(t, err.Error(), "lines are required")

	fvc.EmployerRegistration = "123456789012345678901"
	fvc.AccountNumber = "012345678901234567891"

	err = rules.Validate(fvc, withAddonContext())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "employer registration must be no more than 20 characters")
	assert.Contains(t, err.Error(), "account number must be no more than 20 characters")
}

func TestInvalidFoodVouchersLine(t *testing.T) {
	fvc := &cfdi.FoodVouchers{Lines: []*cfdi.FoodVouchersLine{{}}}

	err := rules.Validate(fvc, withAddonContext())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "line e-wallet ID is required")
	assert.Contains(t, err.Error(), "line issue date and time is required")
	assert.Contains(t, err.Error(), "line employee is required")

	fvc.Lines[0].EWalletID = "123456789012345678901"

	err = rules.Validate(fvc, withAddonContext())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "line e-wallet ID must be no more than 20 characters")
}

func TestInvalidFoodVouchersEmployee(t *testing.T) {
	fvc := &cfdi.FoodVouchers{Lines: []*cfdi.FoodVouchersLine{{Employee: &cfdi.FoodVouchersEmployee{}}}}

	err := rules.Validate(fvc, withAddonContext())

	require.Error(t, err)
	assert.Contains(t, err.Error(), "employee tax code is required")
	assert.Contains(t, err.Error(), "employee CURP is required")
	assert.Contains(t, err.Error(), "employee name is required")

	fvc.Lines[0].Employee.TaxCode = "INVALID1"
	fvc.Lines[0].Employee.CURP = "INVALID2"
	fvc.Lines[0].Employee.SocialSecurity = "INVALID3"

	err = rules.Validate(fvc, withAddonContext())

	assert.ErrorContains(t, err, "employee tax identity code is invalid")
	assert.ErrorContains(t, err, "employee CURP format is invalid")
	assert.ErrorContains(t, err, "employee social security number format is invalid")
}

func TestCalculateFoodVouchers(t *testing.T) {
	fvc := &cfdi.FoodVouchers{
		Lines: []*cfdi.FoodVouchersLine{
			{Amount: num.MakeAmount(1234, 3)},
			{Amount: num.MakeAmount(4321, 3)},
		},
	}

	err := fvc.Calculate()

	require.NoError(t, err)
	assert.Equal(t, num.MakeAmount(123, 2), fvc.Lines[0].Amount)
	assert.Equal(t, num.MakeAmount(555, 2), fvc.Total)
}
