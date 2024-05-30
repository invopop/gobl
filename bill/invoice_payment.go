package bill

import (
	"context"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// InvoicePayment contains details as to how the invoice should be paid.
type InvoicePayment struct {
	// The party responsible for receiving payment of the invoice, if not the supplier.
	Payee *org.Party `json:"payee,omitempty" jsonschema:"title=Payee"`
	// Payment terms or conditions.
	Terms *pay.Terms `json:"terms,omitempty" jsonschema:"title=Terms"`
	// Any amounts that have been paid in advance and should be deducted from the amount due.
	Advances []*pay.Advance `json:"advances,omitempty" jsonschema:"title=Advances"`
	// Details on how payment should be made.
	Instructions *pay.Instructions `json:"instructions,omitempty" jsonschema:"title=Instructions"`
}

// ValidateWithContext checks to make sure the payment data looks good
func (p *InvoicePayment) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithRegime(ctx, p,
		validation.Field(&p.Payee),
		validation.Field(&p.Terms),
		validation.Field(&p.Advances),
		validation.Field(&p.Instructions),
	)
}

// ResetAdvances clears the advances list.
func (p *InvoicePayment) ResetAdvances() {
	if p == nil {
		return
	}
	p.Advances = make([]*pay.Advance, 0)
}

func (p *InvoicePayment) calculateAdvances(zero num.Amount, totalWithTax num.Amount) {
	for _, a := range p.Advances {
		a.CalculateFrom(totalWithTax)
		a.Amount = a.Amount.MatchPrecision(zero)
	}
}

func (p *InvoicePayment) totalAdvance(zero num.Amount) *num.Amount {
	if p == nil || len(p.Advances) == 0 {
		return nil
	}
	sum := zero
	for _, a := range p.Advances {
		sum = sum.MatchPrecision(a.Amount)
		sum = sum.Add(a.Amount)
		a.Amount = a.Amount.Rescale(zero.Exp())
	}
	return &sum
}
