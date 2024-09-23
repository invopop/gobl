package bill_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestOrderingNormalize(t *testing.T) {
	o := &bill.Ordering{
		Code: " Foo ",
		Projects: []*org.DocumentRef{
			{
				Code: " Bar ",
				Ext: tax.Extensions{
					"missing": "",
				},
			},
		},
	}
	o.Normalize(nil)
	assert.Equal(t, "FOO", o.Code.String())
	assert.Equal(t, "BAR", o.Projects[0].Code.String())
	assert.Empty(t, o.Projects[0].Ext)
}

func TestOrderingValidate(t *testing.T) {
	o := &bill.Ordering{
		Code: "123",
	}
	err := o.Validate()
	assert.NoError(t, err)

	o.Projects = []*org.DocumentRef{
		{},
	}
	err = o.Validate()
	assert.ErrorContains(t, err, "projects: (0: (code: cannot be blank.).)")
}
