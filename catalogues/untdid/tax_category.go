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

var extTaxCategory = &cbc.Definition{
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
	Values: []*cbc.Definition{
		{
			Code: "A",
			Name: i18n.String{
				i18n.EN: "Mixed tax rate",
			},
		},
		{
			Code: "AA",
			Name: i18n.String{
				i18n.EN: "Lower rate",
			},
		},
		{
			Code: "AB",
			Name: i18n.String{
				i18n.EN: "Exempt for resale",
			},
		},
		{
			Code: "AC",
			Name: i18n.String{
				i18n.EN: "Exempt for resale",
			},
		},
		{
			Code: "AD",
			Name: i18n.String{
				i18n.EN: "Value Added Tax (VAT) due from a previous invoice",
			},
		},
		{
			Code: "AE",
			Name: i18n.String{
				i18n.EN: "VAT Reverse Charge",
			},
		},
		{
			Code: "B",
			Name: i18n.String{
				i18n.EN: "Transferred (VAT)",
			},
		},
		{
			Code: "C",
			Name: i18n.String{
				i18n.EN: "Duty paid by supplier",
			},
		},
		{
			Code: "D",
			Name: i18n.String{
				i18n.EN: "Value Added Tax (VAT) margin scheme - travel agents",
			},
		},
		{
			Code: "E",
			Name: i18n.String{
				i18n.EN: "Exempt from tax",
			},
		},
		{
			Code: "F",
			Name: i18n.String{
				i18n.EN: "Value Added Tax (VAT) margin scheme - second-hand goods",
			},
		},
		{
			Code: "G",
			Name: i18n.String{
				i18n.EN: "Free export item, tax not charged",
			},
		},
		{
			Code: "H",
			Name: i18n.String{
				i18n.EN: "Higher rate",
			},
		},
		{
			Code: "I",
			Name: i18n.String{
				i18n.EN: "Value Added Tax (VAT) margin scheme - works of art",
			},
		},
		{
			Code: "J",
			Name: i18n.String{
				i18n.EN: "Value Added Tax (VAT) margin scheme - collector's items and antiques",
			},
		},
		{
			Code: "K",
			Name: i18n.String{
				i18n.EN: "VAT exempt for EEA intra-community supply of goods and services",
			},
		},
		{
			Code: "L",
			Name: i18n.String{
				i18n.EN: "Canary Islands general indirect tax",
			},
		},
		{
			Code: "M",
			Name: i18n.String{
				i18n.EN: "Tax for production, services and importation in Ceuta and Melilla",
			},
		},
		{
			Code: "O",
			Name: i18n.String{
				i18n.EN: "Services outside scope of tax",
			},
		},
		{
			Code: "S",
			Name: i18n.String{
				i18n.EN: "Standard Rate",
			},
		},
		{
			Code: "Z",
			Name: i18n.String{
				i18n.EN: "Zero rated goods",
			},
		},
	},
}
