package bill_test

import (
	"strings"
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/schema"
	"github.com/stretchr/testify/assert"
)

func TestShortSchemas(t *testing.T) {
	tests := map[string]any{
		bill.ShortSchemaOrder:             bill.Order{},
		bill.ShortSchemaDelivery:          bill.Delivery{},
		bill.ShortSchemaInvoice:           bill.Invoice{},
		bill.ShortSchemaPayment:           bill.Payment{},
		bill.ShortSchemaStatus:            bill.Status{},
		bill.ShortSchemaAction:            bill.Action{},
		bill.ShortSchemaCharge:            bill.Charge{},
		bill.ShortSchemaCorrectionOptions: bill.CorrectionOptions{},
		bill.ShortSchemaDeliveryDetails:   bill.DeliveryDetails{},
		bill.ShortSchemaDiscount:          bill.Discount{},
		bill.ShortSchemaLine:              bill.Line{},
		bill.ShortSchemaLineCharge:        bill.LineCharge{},
		bill.ShortSchemaLineDiscount:      bill.LineDiscount{},
		bill.ShortSchemaReason:            bill.Reason{},
		bill.ShortSchemaStatusLine:        bill.StatusLine{},
		bill.ShortSchemaTax:               bill.Tax{},
	}
	for short, obj := range tests {
		t.Run(short, func(t *testing.T) {
			id := schema.Lookup(obj)
			assert.True(t, strings.HasSuffix(id.String(), "/"+short), "schema ID %s", id)
		})
	}
}
