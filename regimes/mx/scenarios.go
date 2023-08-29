package mx

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
			Codes: cbc.CodeMap{
				KeySATTipoDeComprobante: "I",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Codes: cbc.CodeMap{
				KeySATTipoDeComprobante: "E",
				KeySATTipoRelacion:      "01",
			},
		},
	},
}
