// Package cfdi implements the CFDI (Comprobante Fiscal Digital por Internet) extensions
// and validation rules that need to be applied to GOBL documents
// in order to comply with the Mexican tax authority (SAT).
package cfdi

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
)

// Key to identify the CFDI addon.
const (
	V4 cbc.Key = "mx-cfdi-v4"
)

// Official CFDI codes to include in stamps.
const (
	StampSignature cbc.Key = "cfdi-sig"    // Signature - Sello Digital del CFDI
	StampSerial    cbc.Key = "cfdi-serial" // Cert Serial - Número de Certificado del CFDI
)

func init() {
	tax.RegisterAddonDef(newAddon())

	// TODO: rename complements to use cfdi in schema path.
	schema.Register(schema.GOBL.Add("regimes/mx"),
		FuelAccountBalance{},
		FoodVouchers{},
	)
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V4,
		Name: i18n.String{
			i18n.EN: "Mexican SAT CFDI v4.X",
		},
		Extensions: extensions,
		Normalizer: normalize,
		Scenarios:  scenarios,
		Validator:  validate,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	case *org.Party:
		normalizeParty(obj)
	case *org.Item:
		normalizeItem(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *org.Item:
		return validateItem(obj)
	}
	return nil
}
