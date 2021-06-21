package pay

import (
	"github.com/invopop/gobl/i18n"
)

// Terms defines when we expect the customer to pay, or have paid, for
// the contents of the document.
type Terms struct {
	Code  TermCode     `json:"code" jsonschema:"title=Code,description=Type of terms to be applied."`
	Notes *i18n.String `json:"notes,omitempty" jsonschema:"title=Notes,description=Description of the conditions for payment."`
}

type TermCode string

const (
	TermPIA   TermCode = "PIA"   // Payment in Advance
	TermPOD   TermCode = "POD"   // Payment on Delivery
	TermNet7  TermCode = "NET7"  // within 7 days
	TermNet10 TermCode = "NET10" // within 10 days
	TermNet21 TermCode = "NET21" // within 21 days
	TermNet30 TermCode = "NET30" // within 30 days
	TermNet60 TermCode = "NET60" // within 60 days
	TermNet90 TermCode = "NET90" // within 90 days
	TermEOM   TermCode = "EOM"   // End of Month
	TermOther TermCode = "OTHER"
)
