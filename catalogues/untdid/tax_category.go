package untdid

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

const (
	// ExtKeyTaxCategory is used to identify the UNTDID 5305 duty/tax/fee category code.
	ExtKeyTaxCategory cbc.Key = "untdid-tax-category"
)

var extTaxCategory = &cbc.KeyDefinition{
	Key: ExtKeyTaxCategory,
	Name: i18n.String{
		i18n.EN: "UNTDID 3505 Tax Category",
	},
	Desc: i18n.String{
		i18n.EN: here.Doc(`
				UNTDID 5305 code used to describe the applicable duty/tax/fee category. There are
				multiple versions and subsets of this table so regimes and addons may need to filter
				options for a specific subset of values.

				Data from https://unece.org/fileadmin/DAM/trade/untdid/d16b/tred/tred5305.htm.
			`),
	},
	Values: []*cbc.ValueDefinition{
		{
			Value: "A",
			Name: i18n.String{
				i18n.EN: "Mixed tax rate",
			},
		},
		{
			Value: "AA",
			Name: i18n.String{
				i18n.EN: "Lower rate",
			},
		},
		{
			Value: "AB",
			Name: i18n.String{
				i18n.EN: "Exempt for resale",
			},
		},
		{
			Value: "AC",
			Name: i18n.String{
				i18n.EN: "Exempt for resale",
			},
		},
		{
			Value: "AD",
			Name: i18n.String{
				i18n.EN: "Value Added Tax (VAT) due from a previous invoice",
			},
		},
		{
			Value: "AE",
			Name: i18n.String{
				i18n.EN: "VAT Reverse Charge",
			},
		},
		{
			Value: "B",
			Name: i18n.String{
				i18n.EN: "Transferred (VAT)",
			},
		},
		{
			Value: "C",
			Name: i18n.String{
				i18n.EN: "Duty paid by supplier",
			},
		},
		{
			Value: "D",
			Name: i18n.String{
				i18n.EN: "Value Added Tax (VAT) margin scheme - travel agents",
			},
		},
		{
			Value: "E",
			Name: i18n.String{
				i18n.EN: "Exempt from tax",
			},
		},
		{
			Value: "F",
			Name: i18n.String{
				i18n.EN: "Value Added Tax (VAT) margin scheme - second-hand goods",
			},
		},
		{
			Value: "G",
			Name: i18n.String{
				i18n.EN: "Free export item, tax not charged",
			},
		},
		{
			Value: "H",
			Name: i18n.String{
				i18n.EN: "Higher rate",
			},
		},
		{
			Value: "I",
			Name: i18n.String{
				i18n.EN: "Value Added Tax (VAT) margin scheme - works of art",
			},
		},
		{
			Value: "J",
			Name: i18n.String{
				i18n.EN: "Value Added Tax (VAT) margin scheme - collector's items and antiques",
			},
		},
		{
			Value: "K",
			Name: i18n.String{
				i18n.EN: "VAT exempt for EEA intra-community supply of goods and services",
			},
		},
		{
			Value: "L",
			Name: i18n.String{
				i18n.EN: "Canary Islands general indirect tax",
			},
		},
		{
			Value: "M",
			Name: i18n.String{
				i18n.EN: "Tax for production, services and importation in Ceuta and Melilla",
			},
		},
		{
			Value: "O",
			Name: i18n.String{
				i18n.EN: "Services outside scope of tax",
			},
		},
		{
			Value: "S",
			Name: i18n.String{
				i18n.EN: "Standard Rate",
			},
		},
		{
			Value: "Z",
			Name: i18n.String{
				i18n.EN: "Zero rated goods",
			},
		},
	},
}
