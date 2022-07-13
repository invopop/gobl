package bill

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/org"
)

// Delivery covers the details of the destination for the products described
// in the invoice body.
type Delivery struct {
	// The party who will receive delivery of the goods defined in the invoice and is not responsible for taxes.
	Receiver *org.Party `json:"receiver,omitempty" jsonschema:"title=Receiver"`
	// When the goods should be expected
	Date *cal.Date `json:"date,omitempty" jsonschema:"title=Date"`
	// Start of a n invoicing or delivery period
	StartDate *cal.Date `json:"start_date,omitempty" jsonschema:"title=Start Date"`
	// End of a n invoicing or delivery period
	EndDate *cal.Date `json:"end_date,omitempty" jsonschema:"title=End Date"`
}

// Validate the delivery details
func (d *Delivery) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Receiver),
		validation.Field(&d.Date),
		validation.Field(&d.StartDate),
		validation.Field(&d.EndDate),
	)
}
