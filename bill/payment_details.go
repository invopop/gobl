package bill

import (
	"context"

	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// PaymentDetails contains details as to how the invoice should be paid.
type PaymentDetails struct {
	// The party responsible for receiving payment of the invoice, if not the supplier.
	Payee *org.Party `json:"payee,omitempty" jsonschema:"title=Payee"`
	// Payment terms or conditions.
	Terms *pay.Terms `json:"terms,omitempty" jsonschema:"title=Terms"`
	// Any amounts that have been paid in advance and should be deducted from the amount due.
	Advances []*pay.Advance `json:"advances,omitempty" jsonschema:"title=Advances"`
	// Details on how payment should be made.
	Instructions *pay.Instructions `json:"instructions,omitempty" jsonschema:"title=Instructions"`
	// Extension fields for additional payment details.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
}

// Normalize will try to normalize the payment's data.
func (p *PaymentDetails) Normalize(normalizers tax.Normalizers) {
	if p == nil {
		return
	}
	p.Ext = tax.CleanExtensions(p.Ext)
	normalizers.Each(p)
	tax.Normalize(normalizers, p.Payee)
	tax.Normalize(normalizers, p.Terms)
	tax.Normalize(normalizers, p.Advances)
	tax.Normalize(normalizers, p.Instructions)
}

// ValidateWithContext checks to make sure the payment data looks good
func (p *PaymentDetails) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, p,
		validation.Field(&p.Payee),
		validation.Field(&p.Terms),
		validation.Field(&p.Advances),
		validation.Field(&p.Instructions),
		validation.Field(&p.Ext),
	)
}

// ResetAdvances clears the advances list.
func (p *PaymentDetails) ResetAdvances() {
	if p == nil {
		return
	}
	p.Advances = make([]*pay.Advance, 0)
}

func (p *PaymentDetails) calculateAdvances(zero num.Amount, payable num.Amount) {
	for _, a := range p.Advances {
		a.CalculateFrom(payable)
		// Payments must always have currency precision
		a.Amount = a.Amount.Rescale(zero.Exp())
	}
}

func (p *PaymentDetails) totalAdvance(zero num.Amount) *num.Amount {
	if p == nil || len(p.Advances) == 0 {
		return nil
	}
	// Payments always maintain currency precision
	sum := zero
	for _, a := range p.Advances {
		sum = sum.Add(a.Amount)
	}
	return &sum
}
