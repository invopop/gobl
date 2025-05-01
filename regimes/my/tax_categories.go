package my

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

var taxCategories = []*tax.CategoryDef{
	//
	// Sales Tax (part of SST)
	//
	{
		Code: "sales-tax", // Custom subcategory under SST
		Name: i18n.String{
			i18n.EN: "Sales Tax",
		},
		Title: i18n.String{
			i18n.EN: "Sales and Service Tax (Sales Portion)",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Royal Malaysian Customs Department - Sales and Service Tax",
				},
				URL: "https://mysst.customs.gov.my/",
			},
		},
		Rates: []*tax.RateDef{
			{
				Key: "standard-5",
				Name: i18n.String{
					i18n.EN: "Standard Rate 5%",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2018, 9, 1),
						Percent: num.MakePercentage(5, 2),
					},
				},
			},
			{
				Key: "standard-10",
				Name: i18n.String{
					i18n.EN: "Standard Rate 10%",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2018, 9, 1),
						Percent: num.MakePercentage(10, 2),
					},
				},
			},
		},
	},

	//
	// Service Tax (part of SST)
	//
	{
		Code: "service-tax", // Custom subcategory under SST
		Name: i18n.String{
			i18n.EN: "Service Tax",
		},
		Title: i18n.String{
			i18n.EN: "Sales and Service Tax (Service Portion)",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Royal Malaysian Customs Department - Service Tax Guide",
				},
				URL: "https://mysst.customs.gov.my/",
			},
		},
		Rates: []*tax.RateDef{
			{
				Key: "standard-6",
				Name: i18n.String{
					i18n.EN: "Specific Services Rate 6%",
				},
				Description: i18n.String{
					i18n.EN: "Applies to food and beverage, telecommunications, parking, logistics services only.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2018, 9, 1),
						Percent: num.MakePercentage(6, 2),
					},
				},
			},
			{
				Key: "standard-8",
				Name: i18n.String{
					i18n.EN: "General Services Rate 8%",
				},
				Description: i18n.String{
					i18n.EN: "Applies to general services as of March 1, 2024.",
				},
				Values: []*tax.RateValueDef{
					{
						Since:   cal.NewDate(2024, 3, 1),
						Percent: num.MakePercentage(8, 2),
					},
				},
			},
		},
	},
}
