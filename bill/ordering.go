package bill

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/org"
)

// Ordering allows additional order details to be appended
type Ordering struct {
	// Party who is selling the goods and is not responsible for taxes
	Seller *org.Party `json:"seller,omitempty" jsonschema:"title=Seller"`
}

// Validate the ordering details.
func (o *Ordering) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Seller),
	)
}
