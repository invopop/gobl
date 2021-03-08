package es

import (
	"cloud.google.com/go/civil"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/tax"
)

// Tax Category definitions.
const (
	TaxCategoryVAT  tax.Category = "VAT"
	TaxCategoryIRPF tax.Category = "IRPF"
	TaxCategoryIGIC tax.Category = "IGIC"
	TaxCategoryIPSI tax.Category = "IPSI"
)

// Tax Code definitions for each tax within a category.
const (
	// VAT Category
	TaxCodeVATZero         tax.Code = "00"
	TaxCodeVATStandard     tax.Code = "01"
	TaxCodeVATReduced      tax.Code = "02"
	TaxCodeVATSuperReduced tax.Code = "03"

	// IRPF Category

)

var taxDefs = tax.Defs{
	/*
	 * VAT
	 */
	{
		Category: TaxCategoryVAT,
		Code:     TaxCodeVATZero,
		Name:     "VAT Zero Rate",
		Rates: []tax.Rate{
			{
				Value: num.MakePercentage(0, 3),
			},
		},
	},
	{
		Category: TaxCategoryVAT,
		Code:     TaxCodeVATStandard,
		Name:     "VAT Standard Rate",
		Rates: []tax.Rate{
			{
				Since: civil.Date{Year: 2012, Month: 9, Day: 1},
				Value: num.MakePercentage(210, 3),
			},
			{
				Since: civil.Date{Year: 2010, Month: 7, Day: 1},
				Value: num.MakePercentage(180, 3),
			},
			{
				Since: civil.Date{Year: 2007, Month: 1, Day: 1},
				Value: num.MakePercentage(160, 3),
			},
		},
	},
	{
		Category: TaxCategoryVAT,
		Code:     TaxCodeVATReduced,
		Name:     "VAT Reduced Rate",
		Rates: []tax.Rate{
			{
				Since: civil.Date{Year: 2012, Month: 9, Day: 1},
				Value: num.MakePercentage(100, 3),
			},
			{
				Since: civil.Date{Year: 2010, Month: 7, Day: 1},
				Value: num.MakePercentage(80, 3),
			},
			{
				Since: civil.Date{Year: 2007, Month: 1, Day: 1},
				Value: num.MakePercentage(70, 3),
			},
		},
	},
	{
		Category: TaxCategoryVAT,
		Code:     TaxCodeVATSuperReduced,
		Name:     "VAT Super-Reduced Rate",
		Rates: []tax.Rate{
			{
				Since: civil.Date{Year: 2007, Month: 1, Day: 1},
				Value: num.MakePercentage(40, 3),
			},
		},
	},

	/*
	 * IRPF
	 */
	{
		Category: TaxCategoryIRPF,
		Name:     "IRPF Max",
		Code:     "01",
		Rates: []tax.Rate{
			{
				Since: civil.Date{Year: 2015, Month: 7, Day: 12},
				Value: num.MakePercentage(150, 3),
			},
			{
				Since: civil.Date{Year: 2015, Month: 1, Day: 1},
				Value: num.MakePercentage(190, 3),
			},
			{
				Since: civil.Date{Year: 2012, Month: 9, Day: 1},
				Value: num.MakePercentage(210, 3),
			},
			{
				Since: civil.Date{Year: 2007, Month: 1, Day: 1},
				Value: num.MakePercentage(150, 3),
			},
		},
		Retained: true,
	},
	{
		Category: TaxCategoryIRPF,
		Name:     "IRPF Low",
		Code:     "02",
		Rates: []tax.Rate{
			{
				Since: civil.Date{Year: 2007, Month: 1, Day: 1},
				Value: num.MakePercentage(70, 3),
			},
		},
		Retained: true,
	},
	{
		Category: TaxCategoryIRPF,
		Name:     "IRPF Min",
		Code:     "03",
		Rates: []tax.Rate{
			{
				Since: civil.Date{Year: 2007, Month: 1, Day: 1},
				Value: num.MakePercentage(10, 3),
			},
		},
		Retained: true,
	},
}
