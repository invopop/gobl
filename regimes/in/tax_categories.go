// Package in defines GST tax categories specific to India.
package in

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Tax categories specific for India.
const (
	TaxCategoryCGST  cbc.Code = "CGST"
	TaxCategorySGST  cbc.Code = "SGST"
	TaxCategoryIGST  cbc.Code = "IGST"
	TaxCategoryUTGST cbc.Code = "UTGST"
	TaxCategoryCess  cbc.Code = "CESS"
)

var taxCategories = []*tax.CategoryDef{
	// Central Goods and Services Tax (CGST)
	{
		Code: TaxCategoryCGST,
		Name: i18n.String{
			i18n.EN: "CGST",
			i18n.HI: "सीजीएसटी",
		},
		Title: i18n.String{
			i18n.EN: "Central Goods and Services Tax",
			i18n.HI: "केंद्रीय माल और सेवा कर",
		},
		Sources: []*tax.Source{
			{
				Title: i18n.String{
					i18n.EN: "Central GST Regulations",
					i18n.HI: "केंद्रीय जीएसटी नियमावली",
				},
				URL: "https://gstcouncil.gov.in/central-gst",
			},
		},
	},

	// State Goods and Services Tax (SGST)
	{
		Code: TaxCategorySGST,
		Name: i18n.String{
			i18n.EN: "SGST",
			i18n.HI: "एसजीएसटी",
		},
		Title: i18n.String{
			i18n.EN: "State Goods and Services Tax",
			i18n.HI: "राज्य माल और सेवा कर",
		},
		Sources: []*tax.Source{
			{
				Title: i18n.String{
					i18n.EN: "State GST Regulations",
					i18n.HI: "राज्य जीएसटी नियमावली",
				},
				URL: "https://gstcouncil.gov.in/sgst",
			},
		},
	},

	// Integrated Goods and Services Tax (IGST)
	{
		Code: TaxCategoryIGST,
		Name: i18n.String{
			i18n.EN: "IGST",
			i18n.HI: "आईजीएसटी",
		},
		Title: i18n.String{
			i18n.EN: "Integrated Goods and Services Tax",
			i18n.HI: "एकीकृत माल और सेवा कर",
		},
		Sources: []*tax.Source{
			{
				Title: i18n.String{
					i18n.EN: "Integrated GST Regulations",
					i18n.HI: "एकीकृत जीएसटी नियमावली",
				},
				URL: "https://gstcouncil.gov.in/gst_council/igst",
			},
		},
	},

	// Union Territory Goods and Services Tax (UTGST)
	{
		Code: TaxCategoryUTGST,
		Name: i18n.String{
			i18n.EN: "UTGST",
			i18n.HI: "यूटीजीएसटी",
		},
		Title: i18n.String{
			i18n.EN: "Union Territory Goods and Services Tax",
			i18n.HI: "केंद्र शासित प्रदेश माल और सेवा कर",
		},
		Sources: []*tax.Source{
			{
				Title: i18n.String{
					i18n.EN: "Union Territory GST Regulations",
					i18n.HI: "यूटीजीएसटी नियमावली",
				},
				URL: "https://gstcouncil.gov.in/utgst",
			},
		},
	},

	// Cess (Additional Tax for Luxury or Specific Goods)
	{
		Code: TaxCategoryCess,
		Name: i18n.String{
			i18n.EN: "Cess",
			i18n.HI: "उपकर",
		},
		Title: i18n.String{
			i18n.EN: "GST Compensation Cess on Luxury or Specific Goods",
			i18n.HI: "विलासिता या विशेष वस्तुओं पर जीएसटी मुआवजा उपकर",
		},
		Sources: []*tax.Source{
			{
				Title: i18n.String{
					i18n.EN: "GST Compensation Cess Regulations",
					i18n.HI: "जीएसटी मुआवजा उपकर नियमावली",
				},
				URL: "https://gstcouncil.gov.in",
			},
		},
	},
}
