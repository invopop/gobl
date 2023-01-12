package bill

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
)

// Payment contains details as to how the invoice should be paid.
type Payment struct {
	// The party responsible for receiving payment of the invoice, if not the supplier.
	Payee *org.Party `json:"payee,omitempty" jsonschema:"title=Payer"`
	// Payment terms or conditions.
	Terms *pay.Terms `json:"terms,omitempty" jsonschema:"title=Terms"`
	// Any amounts that have been paid in advance and should be deducted from the amount due.
	Advances pay.Advances `json:"advances,omitempty" jsonschema:"title=Advances"`
	// Details on how payment should be made.
	Instructions *pay.Instructions `json:"instructions,omitempty" jsonschema:"title=Instructions"`
}

// Validate checks to make sure the payment data looks good
func (p *Payment) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Payee),
		validation.Field(&p.Terms),
		validation.Field(&p.Advances),
		validation.Field(&p.Instructions),
	)
}

func (p *Payment) calculateAdvances(totalWithTax num.Amount) {
	for _, a := range p.Advances {
		a.CalculateFrom(totalWithTax)
	}
}

func (p *Payment) totalAdvance(zero num.Amount) *num.Amount {
	if p == nil || len(p.Advances) == 0 {
		return nil
	}
	sum := zero
	for _, a := range p.Advances {
		sum = sum.Add(a.Amount)
	}
	return &sum
}
