package mx

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Constants for the precision of complement's amounts
const (
	FuelAccountPriceMinimumPrecision = 3
	FuelAccountTotalsPrecision       = 2
)

// FuelAccountValidTaxCodes lists of the complement's allowed tax codes
var FuelAccountValidTaxCodes = []any{
	tax.CategoryVAT,
	TaxCategoryIEPS,
}

// FuelAccountBalance carries the data to produce a CFDI's "Complemento de
// Estado de Cuenta de Combustibles para Monederos Electrónicos" (version 1.2
// revision B) providing detailed information about fuel purchases made with
// electronic wallets. In Mexico, e-wallet suppliers are required to report this
// complementary information in the invoices they issue to their customers.
//
// This struct maps to the `EstadoDeCuentaCombustible` root node in the CFDI's
// complement.
type FuelAccountBalance struct {
	// Customer's account number (maps to `NumeroDeCuenta`).
	AccountNumber string `json:"account_number" jsonschema:"title=Account Number"`
	// Sum of all line totals (i.e. taxes not included) (calculated, maps to `SubTotal`).
	Subtotal num.Amount `json:"subtotal" jsonschema:"title=Subtotal" jsonschema_extras:"calculated=true"`
	// Grand total after taxes have been applied (calculated, maps to `Total`).
	Total num.Amount `json:"total" jsonschema:"title=Total" jsonschema_extras:"calculated=true"`
	// List of fuel purchases made with the customer's e-wallets (maps to `Conceptos`).
	Lines []*FuelAccountLine `json:"lines" jsonschema:"title=Lines"`
}

// FuelAccountLine represents a single fuel purchase made with an e-wallet
// issued by the invoice's supplier. It maps to one
// `ConceptoEstadoDeCuentaCombustible` node in the CFDI's complement.
type FuelAccountLine struct {
	// Identifier of the e-wallet used to make the purchase (maps to `Identificador`).
	EWalletID cbc.Code `json:"e_wallet_id" jsonschema:"title=E-wallet Identifier"`
	// Date and time of the purchase (maps to `Fecha`).
	PurchaseDateTime cal.DateTime `json:"purchase_date_time" jsonschema:"title=Purchase Date and Time"`
	// Tax Identity Code of the fuel's vendor (maps to `Rfc`)
	VendorTaxCode cbc.Code `json:"vendor_tax_code" jsonschema:"title=Vendor's Tax Identity Code"`
	// Code of the service station where the purchase was made (maps to `ClaveEstacion`).
	ServiceStationCode cbc.Code `json:"service_station_code" jsonschema:"title=Service Station Code"`
	// Amount of fuel units purchased (maps to `Cantidad`)
	Quantity num.Amount `json:"quantity" jsonschema:"title=Quantity"`
	// Details of the fuel purchased.
	Item *FuelAccountItem `json:"item" jsonschema:"title=Item"`
	// Identifier of the purchase (maps to `FolioOperacion`).
	PurchaseCode cbc.Code `json:"purchase_code" jsonschema:"title=Purchase Code"`
	// Result of quantity multiplied by the unit price (maps to `Importe`).
	Total num.Amount `json:"total" jsonschema:"title=Total" jsonschema_extras:"calculated=true"`
	// Map of taxes applied to the purchase (maps to `Traslados`).
	Taxes []*FuelAccountTax `json:"taxes" jsonschema:"title=Taxes"`
}

// FuelAccountItem provides the details of a fuel purchase. Its fields map to
// attributes of the `ConceptoEstadoDeCuentaCombustible` node in the CFDI's
// complement.
type FuelAccountItem struct {
	// Type of fuel (one of `c_ClaveTipoCombustible` codes, maps to `TipoCombustible`).
	Type cbc.Code `json:"type" jsonschema:"title=Type"`
	// Reference unit of measure used in the price and the quantity (maps to `Unidad`).
	Unit org.Unit `json:"unit,omitempty" jsonschema:"title=Unit"`
	// Name of the fuel (maps to `NombreCombustible`).
	Name string `json:"name" jsonschema:"title=Name"`
	// Base price of a single unit of the fuel without taxes (maps to `ValorUnitario`).
	Price num.Amount `json:"price" jsonschema:"title=Price"`
}

// FuelAccountTax represents a single tax applied to a fuel purchase. It maps to
// one `Traslado` node in the CFDI's complement.
type FuelAccountTax struct {
	// Category that identifies the tax ("VAT" or "IEPS", maps to `Impuesto`)
	Category cbc.Code `json:"cat" jsonschema:"title=Category"`
	// Percent applicable to the line total (tasa) to use instead of Rate (maps to `TasaoCuota`)
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent"`
	// Rate is a fixed fee to apply to the line quantity (cuota) (maps to `TasaOCuota`)
	Rate *num.Amount `json:"rate,omitempty" jsonschema:"title=Rate"`
	// Total amount of the tax once the percent or rate has been applied (maps to `Importe`).
	Amount num.Amount `json:"amount" jsonschema:"title=Amount" jsonschema_extras:"calculated=true"`
}

// Validate ensures that the complement's data is valid.
func (fab *FuelAccountBalance) Validate() error {
	return validation.ValidateStruct(fab,
		validation.Field(&fab.AccountNumber,
			validation.Required,
			validation.Length(1, 50),
		),
		validation.Field(&fab.Subtotal, validation.Required),
		validation.Field(&fab.Total, validation.Required),
		validation.Field(&fab.Lines, validation.Required),
	)
}

// Validate ensures that the line's data is valid.
func (fal *FuelAccountLine) Validate() error {
	return validation.ValidateStruct(fal,
		validation.Field(&fal.EWalletID, validation.Required),
		validation.Field(&fal.PurchaseDateTime, cal.DateTimeNotZero()),
		validation.Field(&fal.VendorTaxCode,
			validation.Required,
			validation.By(validateTaxCode),
		),
		validation.Field(&fal.ServiceStationCode,
			validation.Required,
			validation.Length(1, 20),
		),
		validation.Field(&fal.Quantity, num.Positive),
		validation.Field(&fal.Item, validation.Required),

		validation.Field(&fal.PurchaseCode,
			validation.Required,
			validation.Length(1, 50),
		),
		validation.Field(&fal.Total, isValidLineTotal(fal)),
		validation.Field(&fal.Taxes, validation.Required),
	)
}

// Validate ensures that the item's data is valid.
func (fai *FuelAccountItem) Validate() error {
	return validation.ValidateStruct(fai,
		validation.Field(&fai.Type, validation.Required),
		validation.Field(&fai.Name,
			validation.Required,
			validation.Length(1, 300),
		),
		validation.Field(&fai.Price, num.Positive),
	)
}

// Validate ensures that the tax's data is valid.
func (fat *FuelAccountTax) Validate() error {
	return validation.ValidateStruct(fat,
		validation.Field(&fat.Category,
			validation.Required,
			validation.In(FuelAccountValidTaxCodes...),
		),
		validation.Field(&fat.Rate,
			num.Positive,
			validation.When(
				fat.Percent == nil,
				validation.Required,
			),
		),
		validation.Field(&fat.Percent),
		validation.Field(&fat.Amount, num.Positive),
	)
}

func isValidLineTotal(line *FuelAccountLine) validation.Rule {
	if line.Item == nil {
		return validation.Skip
	}

	expected := line.Quantity.Multiply(line.Item.Price).Rescale(2)

	return validation.In(expected).Error("must be quantity x unit_price")
}

// Calculate performs the complement's calculations and normalisations.
func (fab *FuelAccountBalance) Calculate() error {
	// Subtotal an tax total need to be calculated using the expected
	// precision for SAT as the PACs recalculate them as part of the
	// validation process. Inevitably this means precision loss.
	taxtotal := num.MakeAmount(0, FuelAccountTotalsPrecision)
	fab.Subtotal = num.MakeAmount(0, FuelAccountTotalsPrecision)

	for _, l := range fab.Lines {
		// Normalise amounts to the expected precision
		l.Quantity = l.Quantity.RescaleUp(FuelAccountPriceMinimumPrecision)
		if l.Item != nil {
			l.Item.Price = l.Item.Price.RescaleUp(FuelAccountPriceMinimumPrecision)
			l.Total = l.Item.Price.Multiply(l.Quantity)
		}

		for _, t := range l.Taxes {
			// Always calculate totals for each tax
			if t.Percent != nil {
				t.Amount = t.Percent.Of(l.Total)
			} else if t.Rate != nil {
				t.Amount = l.Quantity.Multiply(*t.Rate)
			}
			t.Amount = t.Amount.Rescale(FuelAccountTotalsPrecision)
			taxtotal = taxtotal.Add(t.Amount)
		}

		l.Total = l.Total.Rescale(FuelAccountTotalsPrecision)
		fab.Subtotal = fab.Subtotal.Add(l.Total)
	}

	fab.Total = fab.Subtotal.Add(taxtotal)

	return nil
}
