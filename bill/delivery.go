package bill

import (
	"context"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

// Delivery covers the details of the destination for the products described
// in the invoice body.
type Delivery struct {
	// The party who will receive delivery of the goods defined in the invoice and is not responsible for taxes.
	Receiver *org.Party `json:"receiver,omitempty" jsonschema:"title=Receiver"`
	// When the goods should be expected.
	Date *cal.Date `json:"date,omitempty" jsonschema:"title=Date"`
	// Period of time in which to expect delivery if a specific date is not available.
	Period *cal.Period `json:"period,omitempty" jsonschema:"title=Period"`
	// Additional custom data.
	Meta *cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate checks the delivery details
func (d *Delivery) Validate() error {
	return d.ValidateWithContext(context.Background())
}

// ValidateWithContext checks the delivery details
func (d *Delivery) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, d,
		validation.Field(&d.Receiver),
		validation.Field(&d.Date),
		validation.Field(&d.Period),
		validation.Field(&d.Meta),
	)
}
