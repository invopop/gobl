package mx

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

// Constants for the precision of complement's amounts
const (
	FuelAccountInterimPrecision = 3
	FuelAccountFinalPrecision   = 2
	FuelAccountRatePrecision    = 6
)

// Constants for the complement's allowed tax codes
const (
	FuelAccountTaxCodeVAT  = cbc.Code("IVA")
	FuelAccountTaxCodeIEPS = cbc.Code("IEPS")
)

var validTaxCodes = []interface{}{
	FuelAccountTaxCodeVAT,
	FuelAccountTaxCodeIEPS,
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
	// Type of fuel (one of `c_ClaveTipoCombustible` codes, maps to `TipoCombustible`).
	FuelType cbc.Code `json:"fuel_type" jsonschema:"title=Fuel Type"`
	// Reference unit of measure used in the price and the quantity (maps to `Unidad`).
	Unit org.Unit `json:"unit,omitempty" jsonschema:"title=Unit"`
	// Name of the fuel (maps to `NombreCombustible`).
	FuelName string `json:"fuel_name" jsonschema:"title=Fuel Name"`
	// Base price of a single unit of the fuel without taxes (maps to `ValorUnitario`).
	UnitPrice num.Amount `json:"unit_price" jsonschema:"title=Unit Price"`
	// Identifier of the purchase (maps to `FolioOperacion`).
	PurchaseCode cbc.Code `json:"purchase_code" jsonschema:"title=Purchase Code"`
	// Result of quantity multiplied by the unit price (maps to `Importe`).
	Total num.Amount `json:"total" jsonschema:"title=Total"`
	// Map of taxes applied to the purchase (maps to `Traslados`).
	Taxes []*FuelAccountTax `json:"taxes" jsonschema:"title=Taxes"`
}

// FuelAccountTax represents a single tax applied to a fuel purchase. It maps to
// one `Traslado` node in the CFDI's complement.
type FuelAccountTax struct {
	// Code that identifies the tax ("IVA" or "IEPS", maps to `Impuesto`)
	Code cbc.Code `json:"code" jsonschema:"title=Code"`
	// Rate applicable to either the line total (tasa) or the line quantity (cuota) (maps to `TasaOCuota`).
	Rate num.Amount `json:"rate" jsonschema:"title=Rate"`
	// Total amount of the tax once the rate has been applied (maps to `Importe`).
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
}

// Validate ensures that the complement's data is valid.
func (comp *FuelAccountBalance) Validate() error {
	return validation.ValidateStruct(comp,
		validation.Field(&comp.AccountNumber,
			validation.Required,
			validation.Length(1, 50),
		),
		validation.Field(&comp.Subtotal, validation.Required),
		validation.Field(&comp.Total, validation.Required),
		validation.Field(&comp.Lines,
			validation.Required,
			validation.Each(validation.By(validateFuelAccountLine)),
		),
	)
}

func validateFuelAccountLine(value interface{}) error {
	line, _ := value.(*FuelAccountLine)
	if line == nil {
		return nil
	}

	return validation.ValidateStruct(line,
		validation.Field(&line.EWalletID, validation.Required),
		validation.Field(&line.PurchaseDateTime, cal.DateTimeNotZero()),
		validation.Field(&line.VendorTaxCode,
			validation.Required,
			validation.By(validateTaxCode),
		),
		validation.Field(&line.ServiceStationCode,
			validation.Required,
			validation.Length(1, 20),
		),
		validation.Field(&line.Quantity, num.Positive),
		validation.Field(&line.FuelType, validation.Required),
		validation.Field(&line.FuelName,
			validation.Required,
			validation.Length(1, 300),
		),
		validation.Field(&line.PurchaseCode,
			validation.Required,
			validation.Length(1, 50),
		),
		validation.Field(&line.UnitPrice, num.Positive),
		validation.Field(&line.Total, isValidLineTotal(line)),
		validation.Field(&line.Taxes,
			validation.Required,
			validation.Each(validation.By(validateFuelAccountTax)),
		),
	)
}

func validateFuelAccountTax(value interface{}) error {
	tax, _ := value.(*FuelAccountTax)
	if tax == nil {
		return nil
	}

	return validation.ValidateStruct(tax,
		validation.Field(&tax.Code,
			validation.Required,
			validation.In(validTaxCodes...),
		),
		validation.Field(&tax.Rate, num.Positive),
		validation.Field(&tax.Amount, num.Positive),
	)
}

func isValidLineTotal(line *FuelAccountLine) validation.Rule {
	expected := line.Quantity.Multiply(line.UnitPrice).Rescale(2)

	return validation.In(expected).Error("must be quantity x unit_price")
}

// Calculate performs the complement's calculations and normalisations.
func (comp *FuelAccountBalance) Calculate() error {
	var subtotal, taxtotal num.Amount

	for _, line := range comp.Lines {
		// Normalise amounts to the expected precision
		line.Quantity = line.Quantity.Rescale(FuelAccountInterimPrecision)
		line.UnitPrice = line.UnitPrice.Rescale(FuelAccountInterimPrecision)
		line.Total = line.Total.Rescale(FuelAccountFinalPrecision)

		subtotal = line.Total.Add(subtotal)

		for _, tax := range line.Taxes {
			// Normalise amounts to the expected precision
			tax.Rate = tax.Rate.Rescale(FuelAccountRatePrecision)
			tax.Amount = tax.Amount.Rescale(FuelAccountFinalPrecision)

			taxtotal = tax.Amount.Add(taxtotal)
		}
	}

	comp.Subtotal = subtotal.Rescale(FuelAccountFinalPrecision)
	comp.Total = subtotal.Add(taxtotal).Rescale(FuelAccountFinalPrecision)

	return nil
}
