package bill

import (
	"context"
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// PaymentLine defines the details of a line item in a payment document.
type PaymentLine struct {
	uuid.Identify

	// Line number within the parent document (automatically calculated)
	Index int `json:"i" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`

	// Reference to the document being paid
	Document *org.DocumentRef `json:"document,omitempty" jsonschema:"title=Document"`

	// When making multiple payments for a single document, this specifies the
	// installment number for this payment line.
	Installment int `json:"installment,omitempty" jsonschema:"title=Installment"`

	// Payable reflects the amount of the document that is payable. This will be
	// calculated from the embedded document's amount automatically and converted
	// to the currency of the document.
	Payable *num.Amount `json:"payable,omitempty" jsonschema:"title=Payable"`

	// Amount already paid in previous installments, which may be required
	// by some tax regimes or specific use cases.
	Advances *num.Amount `json:"advances,omitempty" jsonschema:"title=Advances"`

	// Amount of the total payment allocated to the referenced document.
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`

	// Due reflects how much still needs to be paid
	Due *num.Amount `json:"due,omitempty" jsonschema:"title=Due,calculated=true"`

	// Tax contains a breakdown of the taxes that will be applied to this payment line
	// after taking into account currency conversion and the relative amounts.
	Tax *tax.Total `json:"tax,omitempty" jsonschema:"title=Tax"`

	// Additional notes specific to this line item for clarification purposes
	Notes []*org.Note `json:"notes,omitempty" jsonschema:"title=Notes"`
}

// ValidateWithContext ensures that the fields contained in the PaymentLine look correct.
func (pl *PaymentLine) ValidateWithContext(ctx context.Context) error {
	mx := num.MakeAmount(0, pl.Amount.Exp())
	if pl.Payable != nil {
		mx = mx.Add(*pl.Payable)
		if pl.Advances != nil {
			mx = mx.Sub(*pl.Advances)
		}
	}
	return validation.ValidateStructWithContext(ctx, pl,
		validation.Field(&pl.Document),
		validation.Field(&pl.Installment, validation.Min(1), validation.Max(999)),
		validation.Field(&pl.Payable,
			num.Positive,
		),
		validation.Field(&pl.Advances,
			num.ZeroOrPositive,
			validation.When(
				pl.Payable != nil,
				num.Max(pl.Payable),
			),
		),
		validation.Field(&pl.Amount,
			num.Positive,
			validation.When(
				mx.Value() > 0,
				num.Max(mx),
			),
		),
		validation.Field(&pl.Due,
			num.ZeroOrPositive,
		),
		validation.Field(&pl.Notes),
	)
}

func (pl *PaymentLine) calculate(rates []*currency.ExchangeRate, cur currency.Code, rr cbc.Key) error {
	var lt *tax.Total

	if pl.Document != nil {
		var er *currency.ExchangeRate
		if dc := pl.Document.Currency; dc != currency.CodeEmpty {
			// If the document has a currency, we need to ensure there is an exchange
			// rate so any taxes can be converted correctly.
			if er = currency.MatchExchangeRate(rates, dc, cur); er == nil {
				return validation.Errors{
					"document": validation.Errors{
						"currency": fmt.Errorf("missing exchange rate from %s to %s", dc, cur),
					},
				}
			}
		}
		pl.Document.Calculate(cur, rr)
		lt = pl.Document.Tax.Clone()
		lt.Exchange(er, rr)

		// Move the document's payable amount to the line in the main
		// currency.
		if p := pl.Document.Payable; p != nil {
			pc := *p
			if er != nil {
				pc = er.Convert(pc)
			}
			pl.Payable = &pc
		}
	}

	// Perform extra calculations with the payable amount, if present.
	if pl.Payable != nil {
		// Do we need to rescale the taxes to fit?
		if lt != nil {
			factor := pl.Amount.Upscale(2).Divide(*pl.Payable)
			lt.Scale(factor, cur, rr)
			pl.Tax = lt
		} else if pl.Tax != nil {
			pl.Tax.Calculate(cur, rr)
		}
		due := *pl.Payable
		if pl.Advances != nil {
			due = due.Subtract(*pl.Advances)
		}
		due = due.Subtract(pl.Amount)
		pl.Due = &due
	}

	return nil

}
