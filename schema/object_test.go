package schema_test

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/note"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/invopop/gobl"
)

// See also document tests performed in `gobl` package.

func TestObjectUUID(t *testing.T) {
	tr := &tax.Regime{} // doesn't have a UUID field!
	obj, err := schema.NewObject(tr)
	assert.NoError(t, err)
	assert.Empty(t, obj.UUID())

	msg := &note.Message{
		Title:   "just a test",
		Content: "this is a test message",
	}
	obj, err = schema.NewObject(msg)
	require.NoError(t, err)
	assert.Empty(t, msg.UUID)
	assert.Empty(t, obj.UUID())

	msg = &note.Message{
		Title:   "just a test",
		Content: "this is a test message",
	}
	msg.UUID = uuid.V1()
	obj, err = schema.NewObject(msg)
	require.NoError(t, err)
	assert.Equal(t, msg.UUID, obj.UUID())
}

func TestObjectCalculate(t *testing.T) {
	inv := exampleInvoice()
	obj, err := schema.NewObject(inv)
	require.NoError(t, err)

	assert.Nil(t, inv.Totals)
	assert.Empty(t, inv.UUID)
	require.NoError(t, obj.Calculate())
	assert.NotNil(t, inv.Totals)
	assert.NotEmpty(t, inv.UUID)
	assert.NotEmpty(t, obj.UUID())
	assert.Equal(t, obj.UUID(), inv.UUID)
}

func TestObjectReplicate(t *testing.T) {
	inv := exampleInvoice()
	obj, err := schema.NewObject(inv)
	require.NoError(t, err)
	require.NoError(t, obj.Calculate())
	ou := obj.UUID()
	require.NoError(t, obj.Replicate())
	assert.NotEqual(t, ou, obj.UUID())
	assert.Empty(t, inv.Code, "should remove code")
}

// exampleInvoice defines a simple invoice example pre-calculations.
func exampleInvoice() *bill.Invoice {
	return &bill.Invoice{
		Series: "TEST",
		Code:   "000123",
		Tax: &bill.Tax{
			PricesInclude: tax.CategoryVAT,
		},
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "B98602642",
			},
		},
		Customer: &org.Party{
			Name: "Test Customer",
			TaxID: &tax.Identity{
				Country: "ES",
				Code:    "54387763P",
			},
		},
		IssueDate: cal.MakeDate(2022, 6, 13),
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(1, 0),
				Item: &org.Item{
					Name:  "Item",
					Price: num.MakeAmount(4320, 2),
				},
				Taxes: tax.Set{
					{
						Category: "VAT",
						Percent:  num.NewPercentage(6, 2),
					},
				},
			},
		},
	}
}
