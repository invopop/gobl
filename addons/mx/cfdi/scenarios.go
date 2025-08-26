package cfdi

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

var scenarios = []*tax.ScenarioSet{
	invoiceScenarios,
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// TipoDeComprobante / TipoRelacion
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Ext: tax.Extensions{
				ExtKeyDocType: "I",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Ext: tax.Extensions{
				ExtKeyDocType: "E",
				ExtKeyRelType: "01",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard, bill.InvoiceTypeCreditNote},
			Filter: func(in any) bool {
				inv, ok := in.(*bill.Invoice)
				if !ok {
					return false
				}
				return inv.Totals.Paid()
			},
			Ext: tax.Extensions{
				ExtKeyPaymentMethod: ExtCodePaymentMethodPUE, // Pago en una sola exhibición
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard, bill.InvoiceTypeCreditNote},
			Filter: func(in any) bool {
				inv, ok := in.(*bill.Invoice)
				if !ok {
					return false
				}
				return !inv.Totals.Paid()
			},
			Ext: tax.Extensions{
				ExtKeyPaymentMethod: ExtCodePaymentMethodPPD, // Pago en parcialidades o diferido
			},
		},
	},
}
