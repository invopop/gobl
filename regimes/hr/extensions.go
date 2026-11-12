package hr

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Regime extension keys
const (
	// ExtKeyPassThrough says which kind of pass-through item an amount left out
	// of the VAT base is. It goes on the exempt VAT combo of the line or charge
	// that carries the amount.
	ExtKeyPassThrough cbc.Key = "hr-pass-through"
)

// Pass-through items, with the codes the Croatian eRačun specification gives
// them in its HR tax category table (HR-TB-2).
const (
	// PassThroughReturnFee is the return fee for single-use packaging
	// (povratna naknada).
	PassThroughReturnFee cbc.Code = "POVNAK"
	// PassThroughMotorVehicleTax is the special tax on motor vehicles
	// (posebni porez na motorna vozila).
	PassThroughMotorVehicleTax cbc.Code = "PPMV"
	// PassThroughNonTaxableFee is any other fee not subject to VAT that is
	// charged together with taxable items.
	PassThroughNonTaxableFee cbc.Code = "N"
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyPassThrough,
		Name: i18n.String{
			i18n.EN: "Pass-through Item",
			i18n.HR: "Prolazna stavka",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Kind of pass-through item: an amount the supplier charges together with
				taxable items that is left out of the VAT base. The amount is a line or a
				charge of its own, whose VAT combo uses the exempt key and carries this
				extension. The codes are those of the HR tax category table (HR-TB-2) of
				the Croatian eRačun specification, without the "HR:" prefix.

				The consumption tax is a pass-through item too, but it is a percentage of
				the price of the items it applies to, so it has its own tax category, PP,
				instead of a code here.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: PassThroughReturnFee,
				Name: i18n.String{
					i18n.EN: "Packaging return fee",
					i18n.HR: "Povratna naknada",
				},
				Desc: i18n.String{
					i18n.EN: "Return fee for single-use packaging.",
				},
			},
			{
				Code: PassThroughMotorVehicleTax,
				Name: i18n.String{
					i18n.EN: "Special tax on motor vehicles",
					i18n.HR: "Posebni porez na motorna vozila",
				},
				Desc: i18n.String{
					i18n.EN: "Special tax on motor vehicles, charged together with the vehicle.",
				},
			},
			{
				Code: PassThroughNonTaxableFee,
				Name: i18n.String{
					i18n.EN: "Non-taxable fee",
					i18n.HR: "Neoporeziva naknada",
				},
				Desc: i18n.String{
					i18n.EN: "Any other fee not subject to VAT that is charged together with taxable items.",
				},
			},
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Ministry of Finance - Core eInvoice usage specification with extensions (HR-TB-2)",
					i18n.HR: "Ministarstvo financija - Specifikacija osnovne uporabe eRačuna s proširenjima (HR-TB-2)",
				},
				URL: "https://porezna.gov.hr/fiskalizacija/api/dokumenti/196",
				At:  cal.NewDateTime(2026, 11, 12, 0, 0, 0),
			},
		},
	},
}
