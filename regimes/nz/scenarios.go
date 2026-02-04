package nz

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/tax"
)

func scenarios() []*tax.ScenarioSet {
	return []*tax.ScenarioSet{bill.InvoiceScenarios()}
}