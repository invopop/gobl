package mx

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// FuelComplement carries the data to produce a CFDI's "Complemento de Estado de
// Cuenta de Combustibles para Monederos Electrónicos" (version 1.2 revision B)
// providing detailed information about fuel purchases made with electronic
// wallets. In Mexico, e-wallet suppliers are required to report this
// complementary information in the invoices they issue to their customers.
type FuelComplement struct {
	// Customer's e-wallet account number.
	WalletAccountNumber string `json:"wallet_account_number" jsonschema:"title=Wallet Account Number"`
	// List of fuel purchases made with the customer's e-wallets.
	Lines []*FuelLine `json:"lines" jsonschema:"title=Lines"`
	// Summary of all the purchases totals, including taxes (calculated).
	Totals *FuelTotals `json:"totals" jsonschema:"title=Totals" jsonschema_extras:"calculated=true"`
}

// FuelLine represents a single fuel purchase made with an e-wallet issued by
// the invoice's supplier.
type FuelLine struct {
	// Identifier of the e-wallet used to make the purchase.
	WalletID cbc.Code `json:"wallet_id" jsonschema:"title=Wallet Identifier"`
	// Date and time of the purchase.
	PurchaseDateTime cal.DateTime `json:"purchase_date_time" jsonschema:"title=Purchase Date and Time"`
	// Tax ID (RFC) of the fuel seller.
	SellerTaxID cbc.Code `json:"seller_tax_id" jsonschema:"title=Seller's Tax ID"`
	// Code of the service station where the purchase was made.
	ServiceStationCode cbc.Code `json:"service_station_code" jsonschema:"title=Service Station Code"`
	// Amount of fuel units purchased
	Quantity num.Amount `json:"quantity" jsonschema:"title=Quantity"`
	// Reference unit of measure used in the price and the quantity.
	Unit org.Unit `json:"unit" jsonschema:"title=Unit"`
	// Type of fuel ("c_ClaveTipoCombustible" codes).
	FuelType cbc.Code `json:"type,omitempty" jsonschema:"title=Type"`
	// Name of the fuel
	FuelName string `json:"name" jsonschema:"title=Name"`
	// Base price of a single unit of the fuel.
	Price num.Amount `json:"price" jsonschema:"title=Price"`
	// Identifier of the purchase ("Folio").
	PurchaseCode cbc.Code `json:"purchase_code" jsonschema:"title=Purchase Code"`
	// Result of quantity multiplied by the item's price (calculated).
	Total num.Amount `json:"total" jsonschema:"title=Total" jsonschema_extras:"calculated=true"`
	// Map of taxes applied to the purchase.
	Taxes tax.Set `json:"taxes" jsonschema:"title=Taxes"`
}

// FuelTotals contains the summaries of all calculations for the fuel purchases.
type FuelTotals struct {
	// Sum of all line sums.
	Total num.Amount `json:"total" jsonschema:"title=Total"`
	// Grand total after taxes have been applied.
	TotalWithTax num.Amount `json:"total_with_tax" jsonschema:"title=Total with Tax"`
}
