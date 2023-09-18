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
	p.IssueDate = cal.NewDate(2022, 11, 6)

	err := p.Validate()
	assert.NoError(t, err)
}

func TestPrecedingJSONMigration(t *testing.T) {
	data := []byte(`{"correction_method":"foo","corrections":["bar"]}`)
	p := new(bill.Preceding)
	err := p.UnmarshalJSON(data)
	assert.NoError(t, err)
	assert.Equal(t, "foo", p.CorrectionMethod.String())
	assert.Equal(t, "bar", p.Changes[0].String())
}
