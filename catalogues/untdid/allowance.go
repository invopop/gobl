package untdid

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

const (
	// ExtKeyAllowance is used to identify the UNTDID 5189 allownce codes
	// used in discounts.
	ExtKeyAllowance cbc.Key = "untdid-allowance"
)

var extAllowance = &cbc.KeyDefinition{
	Key: ExtKeyAllowance,
	Name: i18n.String{
		i18n.EN: "UNTDID 5189 Allowance",
	},
	Desc: i18n.String{
		i18n.EN: here.Doc(`
			UNTDID 5189 code used to describe the allowance type. This list is based on the
			[EN16931 code list](https://ec.europa.eu/digital-building-blocks/sites/display/DIGITAL/Registry+of+supporting+artefacts+to+implement+EN16931#RegistryofsupportingartefactstoimplementEN16931-Codelists)
			values table which focusses on invoices and payments.
		`),
	},
	Values: []*cbc.ValueDefinition{
		{
			Value: "41",
			Name:  i18n.NewString("Bonus for works ahead of schedule"),
		},
		{
			Value: "42",
			Name:  i18n.NewString("Other bonus"),
		},
		{
			Value: "60",
			Name:  i18n.NewString("Manufacturer’s consumer discount"),
		},
		{
			Value: "62",
			Name:  i18n.NewString("Due to military status"),
		},
		{
			Value: "63",
			Name:  i18n.NewString("Due to work accident"),
		},
		{
			Value: "64",
			Name:  i18n.NewString("Special agreement"),
		},
		{
			Value: "65",
			Name:  i18n.NewString("Production error discount"),
		},
		{
			Value: "66",
			Name:  i18n.NewString("New outlet discount"),
		},
		{
			Value: "67",
			Name:  i18n.NewString("Sample discount"),
		},
		{
			Value: "68",
			Name:  i18n.NewString("End-of-range discount"),
		},
		{
			Value: "70",
			Name:  i18n.NewString("Incoterm discount"),
		},
		{
			Value: "71",
			Name:  i18n.NewString("Point of sales threshold allowance"),
		},
		{
			Value: "88",
			Name:  i18n.NewString("Material surcharge/deduction"),
		},
		{
			Value: "95",
			Name:  i18n.NewString("Discount"),
		},
		{
			Value: "100",
			Name:  i18n.NewString("Special rebate"),
		},
		{
			Value: "102",
			Name:  i18n.NewString("Fixed long term"),
		},
		{
			Value: "103",
			Name:  i18n.NewString("Temporary"),
		},
		{
			Value: "104",
			Name:  i18n.NewString("Standard"),
		},
		{
			Value: "105",
			Name:  i18n.NewString("Yearly turnover"),
		},
	},
}
