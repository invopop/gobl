package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
	"github.com/stretchr/testify/assert"
)

func TestLineRequireTaxValidation(t *testing.T) {
	err := validation.Validate(nil, bill.RequireLineTaxCategory("VAT"), validation.Skip)
	assert.NoError(t, err)

	l := &bill.Line{}
	err = validation.Validate(l, bill.RequireLineTaxCategory(""), validation.Skip)
	assert.NoError(t, err)

	err = validation.Validate(l,
		bill.RequireLineTaxCategory("VAT"),
		validation.Skip,
	)
	assert.ErrorContains(t, err, "taxes: missing category VAT.")

	l = &bill.Line{
		Taxes: tax.Set{
			{
				Category: "VAT",
				Percent:  num.NewPercentage(20, 3),
			},
		},
	}
	err = validation.Validate(l,
		bill.RequireLineTaxCategory("VAT"),
		validation.Skip,
	)
	assert.NoError(t, err)

	err = validation.Validate(l,
		bill.RequireLineTaxCategory("IRPEF"),
		validation.Skip,
	)
	assert.ErrorContains(t, err, "taxes: missing category IRPEF.")
}
