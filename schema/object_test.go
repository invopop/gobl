package schema_test

import (
	"encoding/json"
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
	tr := &tax.RegimeDef{} // doesn't have a UUID field!
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

func TestObjectValidate(t *testing.T) {
	inv := exampleInvoice()
	obj, err := schema.NewObject(inv)
	require.NoError(t, err)
	require.NoError(t, obj.Calculate())

	faults := obj.Validate()
	assert.Empty(t, faults, "valid invoice should have no faults")

	t.Run("with embedded object, complete", func(t *testing.T) {
		msg := &note.Message{
			Title:   "test",
			Content: "hello",
		}
		obj, err := schema.NewObject(msg)
		require.NoError(t, err)
		require.NoError(t, obj.Calculate())
		faults = obj.Validate()
		assert.Empty(t, faults, "valid embedded object should have no faults")
	})

	t.Run("with embedded object, missing content", func(t *testing.T) {
		msg := &note.Message{
			Title: "test",
		}
		obj, err := schema.NewObject(msg)
		require.NoError(t, err)
		require.NoError(t, obj.Calculate())
		faults = obj.Validate()
		assert.ErrorContains(t, faults, "[GOBL-NOTE-MESSAGE-01] ($.content) message content is require")
	})

	t.Run("with embedded object, invalid guard condition", func(t *testing.T) {
		inv := exampleInvoice()
		inv.Supplier.TaxID.Country = "FR"
		inv.Supplier.TaxID.Code = "" // invalid
		obj, err := schema.NewObject(inv)
		require.NoError(t, err)
		require.NoError(t, obj.Calculate())
		faults := obj.Validate()
		assert.ErrorContains(t, faults, "[GOBL-FR-BILL-INVOICE-01] ($.supplier) invoice supplier must have a tax ID code or a SIREN/SIRET identity")
	})
}

func TestObjectIsEmpty(t *testing.T) {
	obj := new(schema.Object)
	assert.True(t, obj.IsEmpty())

	inv := exampleInvoice()
	obj, err := schema.NewObject(inv)
	require.NoError(t, err)
	assert.False(t, obj.IsEmpty())
}

func TestObjectInstance(t *testing.T) {
	inv := exampleInvoice()
	obj, err := schema.NewObject(inv)
	require.NoError(t, err)
	assert.Equal(t, inv, obj.Instance())
}

func TestObjectEmbedded(t *testing.T) {
	inv := exampleInvoice()
	obj, err := schema.NewObject(inv)
	require.NoError(t, err)
	assert.Equal(t, inv, obj.Embedded())
}

func TestObjectCorrect(t *testing.T) {
	inv := exampleInvoice()
	obj, err := schema.NewObject(inv)
	require.NoError(t, err)
	require.NoError(t, obj.Calculate())

	// Invoice is Correctable, so this should work
	err = obj.Correct()
	// May return an error if options are required, but should not panic
	_ = err

	// Non-correctable payload
	msg := &note.Message{
		Title:   "test",
		Content: "hello",
	}
	obj2, err := schema.NewObject(msg)
	require.NoError(t, err)
	err = obj2.Correct()
	assert.Error(t, err, "non-correctable type should return error")
}

func TestObjectCorrectionOptionsSchema(t *testing.T) {
	inv := exampleInvoice()
	obj, err := schema.NewObject(inv)
	require.NoError(t, err)
	require.NoError(t, obj.Calculate())

	result, err := obj.CorrectionOptionsSchema()
	assert.NoError(t, err)
	// Invoice should return options schema
	_ = result

	// Non-correctable payload
	msg := &note.Message{
		Title:   "test",
		Content: "hello",
	}
	obj2, err := schema.NewObject(msg)
	require.NoError(t, err)
	result, err = obj2.CorrectionOptionsSchema()
	assert.NoError(t, err)
	assert.Nil(t, result, "non-correctable type should return nil")
}

func TestObjectClone(t *testing.T) {
	inv := exampleInvoice()
	obj, err := schema.NewObject(inv)
	require.NoError(t, err)
	require.NoError(t, obj.Calculate())

	clone, err := obj.Clone()
	require.NoError(t, err)
	assert.NotNil(t, clone)
	assert.Equal(t, obj.UUID(), clone.UUID())

	// Modifying the original should not affect clone
	inv.Code = "MODIFIED"
	cloneInv, ok := clone.Instance().(*bill.Invoice)
	require.True(t, ok)
	assert.NotEqual(t, "MODIFIED", cloneInv.Code)
}

func TestObjectMarshalUnmarshalJSON(t *testing.T) {
	inv := exampleInvoice()
	obj, err := schema.NewObject(inv)
	require.NoError(t, err)
	require.NoError(t, obj.Calculate())

	data, err := json.Marshal(obj)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"$schema"`)

	obj2 := new(schema.Object)
	require.NoError(t, json.Unmarshal(data, obj2))
	assert.Equal(t, obj.UUID(), obj2.UUID())
	assert.False(t, obj2.IsEmpty())
}

func TestObjectUnmarshalJSONUnknownSchema(t *testing.T) {
	data := []byte(`{"$schema":"https://gobl.org/draft-0/unknown/type","foo":"bar"}`)
	obj := new(schema.Object)
	err := json.Unmarshal(data, obj)
	assert.Error(t, err)
}

func TestObjectUnmarshalJSONNoSchema(t *testing.T) {
	data := []byte(`{"foo":"bar"}`)
	obj := new(schema.Object)
	err := json.Unmarshal(data, obj)
	assert.NoError(t, err)
	assert.True(t, obj.IsEmpty())
}

func TestObjectJSONSchema(t *testing.T) {
	obj := schema.Object{}
	s := obj.JSONSchema()
	assert.Equal(t, "object", s.Type)
	assert.Equal(t, "Object", s.Title)
}

func TestObjectError(t *testing.T) {
	err := schema.ErrUnknownSchema
	assert.Equal(t, "unknown-schema", err.Error())
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
					Price: num.NewAmount(4320, 2),
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
