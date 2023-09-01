// Package complements provides GOBL Complements for Mexican invoices
package complements

import (
	"cloud.google.com/go/civil"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// FuelComplement provides detailed information about fuel purchases made with
// electronic wallets. In Mexico, e-wallet suppliers are required to report this
// information as part of the invoices they issue to their customers.
type FuelComplement struct {
	// Customer's e-wallet account number
	WalletAccountNumber string `json:"wallet-account-number" jsonschema:"title=Wallet Account Number"`
	// List of fuel purchases made with the customer's e-wallets
	Lines []*FuelLine `json:"lines" jsonschema:"title=Lines"`
	// Summary of all the purchases totals, including taxes (calculated)
	Totals bill.Totals `json:"totals" jsonschema:"title=Totals" jsonschema_extras:"calculated=true"`
}

// FuelLine represents a single fuel purchase made with an e-wallet issued by
// the invoice's supplier to the invoice's customer.
type FuelLine struct {
	// Identifier of the e-wallet used to make the purchase
	WalletID cbc.Code `json:"wallet-id" jsonschema:"title=Wallet Identifier"`
	// Date and time of the purchase
	PurchaseDateTime civil.DateTime `json:"purchase-date" jsonschema:"title=Purchase Date"`
	// Tax ID (RFC) of the purchaser
	PurchaserTaxID cbc.Code `json:"purchaser-tax-id" jsonschema:"title=Purchaser's Tax ID"`
	// Code of the service station where the purchase was made
	ServiceStationCode cbc.Code `json:"service-station-code" jsonschema:"title=Service Station Code"`
	// Amount of fuel units purchased
	Quantity num.Amount `json:"quantity" jsonschema:"title=Quantity"`
	// Details about the fuel purchased
	Item FuelItem `json:"item" jsonschema:"title=Item"`
	// Identifier of the purchase (folio)
	PurchaseID cbc.Code `json:"purchase-id" jsonschema:"title=Purchase Identifier"`
	// Result of quantity multiplied by the item's price (calculated)
	Sum num.Amount `json:"sum" jsonschema:"title=Sum" jsonschema_extras:"calculated=true"`
	// Map of taxes applied to the purchase
	Taxes tax.Set `json:"taxes" jsonschema:"title=Taxes"`
}

// FuelItem provides the details of the fuel purchased.
type FuelItem struct {
	// Reference unit of measure used in the price and the quantity
	Unit org.Unit `json:"unit" jsonschema:"title=Unit"`
	// Type of fuel (c_ClaveTipoCombustible codes)
	Type cbc.Code `json:"type,omitempty" jsonschema:"title=Type"`
	// Name of the fuel
	Name string `json:"name" jsonschema:"title=Name"`
	// Base price of a single unit of the fuel
	Price num.Amount `json:"price" jsonschema:"title=Price"`
}
