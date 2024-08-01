package main

import (
	"time"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func main() {
	// Build up a basic invoice document
	inv := &bill.Invoice{
		Series:    "F23",
		Code:      "00010",
		IssueDate: cal.MakeDate(2023, time.May, 11),
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: "US",
			},
			Name:  "Provider One Inc.",
			Alias: "Provider One",
			Emails: []*org.Email{
				{
					Address: "billing@provideone.com",
				},
			},
			Addresses: []*org.Address{
				{
					Number:   "16",
					Street:   "Jessie Street",
					Locality: "San Francisco",
					Region:   "CA",
					Code:     "94105",
					Country:  "US",
				},
			},
		},
		Customer: &org.Party{
			Name: "Sample Customer",
			Emails: []*org.Email{
				{
					Address: "email@sample.com",
				},
			},
		},
		Lines: []*bill.Line{
			{
				Quantity: num.MakeAmount(20, 0),
				Item: &org.Item{
					Name:  "A stylish mug",
					Price: num.MakeAmount(2000, 2),
					Unit:  org.UnitHour,
				},
				Taxes: []*tax.Combo{
					{
						Category: tax.CategoryST,
						Percent:  num.NewPercentage(85, 3),
					},
				},
			},
		},
	}

	// Prepare an "Envelope"
	env := gobl.NewEnvelope()
	if err := env.Insert(inv); err != nil {
		panic(err)
	}

}
