// Package dgfip defines the code lists published by the French tax
// authority (Direction Générale des Finances Publiques) used by the
// CTC (Continuous Transaction Control) e-invoicing and e-reporting
// reform. Extensions defined here are shared across the flow-specific
// addons under addons/fr/ctc.
package dgfip

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterCatalogueDef("dgfip.json")
}

const (
	// ExtKeyBillingMode is the DGFiP "Cadre de Facturation" code that
	// describes the nature of the document (Biens / Services / Mixte)
	// and the payment context. Required on Flow 2 clearance invoices
	// and Flow 10 B2B reporting invoices.
	ExtKeyBillingMode cbc.Key = "dgfip-billing-mode"
)

// Billing mode codes. Prefix denotes invoice nature (B = goods, S =
// services, M = mixed); numeric suffix encodes payment context
// (1 = deposit, 2 = already paid, 4 = final after down payment,
// 5 = subcontractor, 6 = co-contractor, 7 = e-reporting).
const (
	BillingModeB1 cbc.Code = "B1"
	BillingModeB2 cbc.Code = "B2"
	BillingModeB4 cbc.Code = "B4"
	BillingModeB7 cbc.Code = "B7"
	BillingModeS1 cbc.Code = "S1"
	BillingModeS2 cbc.Code = "S2"
	BillingModeS4 cbc.Code = "S4"
	BillingModeS5 cbc.Code = "S5"
	BillingModeS6 cbc.Code = "S6"
	BillingModeS7 cbc.Code = "S7"
	BillingModeM1 cbc.Code = "M1"
	BillingModeM2 cbc.Code = "M2"
	BillingModeM4 cbc.Code = "M4"
)
