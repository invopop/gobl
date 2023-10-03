package mx_test

import (
	"testing"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/mx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvalidFoodVouchersComplement(t *testing.T) {
	fvc := &mx.FoodVouchersComplement{}

	err := fvc.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "account_number: cannot be blank")
	assert.Contains(t, err.Error(), "lines: cannot be blank")

	fvc.EmployerRegistration = "123456789012345678901"
	fvc.AccountNumber = "012345678901234567891"

	err = fvc.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "employer_registration: the length must be no more than 20")
	assert.Contains(t, err.Error(), "account_number: the length must be no more than 20")
}

func TestInvalidFoodVouchersLine(t *testing.T) {
	fvc := &mx.FoodVouchersComplement{Lines: []*mx.FoodVouchersLine{{}}}

	err := fvc.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "e_wallet_id: cannot be blank")
	assert.Contains(t, err.Error(), "issue_date_time: required")
	assert.Contains(t, err.Error(), "employee_tax_code: cannot be blank")
	assert.Contains(t, err.Error(), "employee_curp: cannot be blank")
	assert.Contains(t, err.Error(), "employee_name: cannot be blank")

	fvc.Lines[0].EWalletID = "123456789012345678901"
	fvc.Lines[0].EmployeeTaxCode = "INVALID1"
	fvc.Lines[0].EmployeeCURP = "INVALID2"
	fvc.Lines[0].EmployeeSocialSecurity = "INVALID3"

	err = fvc.Validate()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "e_wallet_id: the length must be no more than 20")
	assert.Contains(t, err.Error(), "employee_tax_code: invalid tax identity code")
	assert.Contains(t, err.Error(), "employee_curp: must be in a valid format")
	assert.Contains(t, err.Error(), "employee_social_security: must be in a valid format")
}

func TestCalculateFoodVouchersComplement(t *testing.T) {
	fvc := &mx.FoodVouchersComplement{
		Lines: []*mx.FoodVouchersLine{
			{Amount: num.MakeAmount(1234, 3)},
			{Amount: num.MakeAmount(4321, 3)},
		},
	}

	err := fvc.Calculate()

	require.NoError(t, err)
	assert.Equal(t, num.MakeAmount(123, 2), fvc.Lines[0].Amount)
	assert.Equal(t, num.MakeAmount(555, 2), fvc.Total)
}
