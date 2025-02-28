package bill

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

// DeliveryDetails covers the details of the destination for the products described
// in the invoice body.
type DeliveryDetails struct {
	// The party who will receive delivery of the goods defined in the invoice and is not responsible for taxes.
	Receiver *org.Party `json:"receiver,omitempty" jsonschema:"title=Receiver"`
	// Identities is used to define specific codes or IDs that may be used to
	// identify the delivery.
	Identities []*org.Identity `json:"identities,omitempty" jsonschema:"title=Identities"`
	// When the goods should be expected.
	Date *cal.Date `json:"date,omitempty" jsonschema:"title=Date"`
	// Period of time in which to expect delivery if a specific date is not available.
	Period *cal.Period `json:"period,omitempty" jsonschema:"title=Period"`
	// Additional custom data.
	Meta *cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Validate the delivery details
func (d *DeliveryDetails) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Receiver),
		validation.Field(&d.Identities),
		validation.Field(&d.Date),
		validation.Field(&d.Period),
		validation.Field(&d.Meta),
	)
}
