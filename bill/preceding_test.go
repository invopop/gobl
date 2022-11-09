package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/stretchr/testify/assert"
)

func TestPrecedingValidation(t *testing.T) {
	p := new(bill.Preceding)
	p.Code = "FOO"
	p.IssueDate = cal.MakeDate(2022, 11, 6)
	p.Corrections = []bill.CorrectionKey{bill.CodeCorrectionKey, bill.LineCorrectionKey}
	p.CorrectionMethod = bill.CompleteCorrectionMethodKey

	err := p.Validate()
	assert.NoError(t, err)
}
