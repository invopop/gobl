package bill

import (
	"context"
	"fmt"

	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/validation"
)

// PaymentLine defines the details of a line required in an invoice.
type PaymentLine struct {
	uuid.Identify
	// Line number inside the parent (calculated)
	Index int `json:"i" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`

	// The document reference related to the payment.
	Document *org.DocumentRef `json:"document,omitempty" jsonschema:"title=Document"`

	// Currency used for the payment if different from the document currency.
	Currency currency.Code `json:"currency,omitempty" jsonschema:"title=Currency"`

	// Amount received by the supplier for ordinary payments.
	Debit *num.Amount `json:"debit,omitempty" jsonschema:"title=Debit"`
	// Amount received by the customer in case of refunds.
	Credit *num.Amount `json:"credit,omitempty" jsonschema:"title=Credit"`

	// Tax total breakdown from the original document, only if required by a specific tax regime
	// or addon.
	Tax *tax.Total `json:"tax,omitempty" jsonschema:"title=Tax"`

	// Total balance to be paid for this line from the customer to the supplier
	// in the currency of the document.
	Total num.Amount `json:"total" jsonschema:"title=Total" jsonschema_extras:"calculated=true"`

	// Set of specific notes for this line that may be required for
	// clarification.
	Notes []*org.Note `json:"notes,omitempty" jsonschema:"title=Notes"`
}

// ValidateWithContext ensures that the fields contained in the PaymentLine look correct.
func (pl *PaymentLine) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, pl,
		validation.Field(&pl.Document),
		validation.Field(&pl.Currency),
		validation.Field(&pl.Debit,
			validation.When(
				pl.Credit == nil,
				validation.Required.Error("must have either debit or credit"),
			),
		),
		validation.Field(&pl.Credit),
		validation.Field(&pl.Tax),
		validation.Field(&pl.Total, validation.Required),
		validation.Field(&pl.Notes),
	)
}

// calculate will ensure the total amount is calculated correctly
func (pl *PaymentLine) calculate(cur currency.Code, rates []*currency.ExchangeRate) error {
	pl.Total = cur.Def().Zero()
	if pl.Debit != nil {
		var a num.Amount
		if pl.Currency != "" {
			na := currency.Convert(rates, pl.Currency, cur, *pl.Debit)
			if na == nil {
				return validation.Errors{
					"currency": fmt.Errorf("no exchange rate found for %s to %s", pl.Currency, cur),
				}
			}
			a = *na
		} else {
			a = *pl.Debit
		}
		pl.Total.MatchPrecision(a)
		pl.Total = pl.Total.Add(a)
	}
	if pl.Credit != nil {
		var a num.Amount
		if pl.Currency != "" {
			na := currency.Convert(rates, pl.Currency, cur, *pl.Credit)
			if na == nil {
				return validation.Errors{
					"currency": fmt.Errorf("no exchange rate found for %s to %s", pl.Currency, cur),
				}
			}
			a = *na
		} else {
			a = *pl.Credit
		}
		pl.Total.MatchPrecision(a)
		pl.Total = pl.Total.Subtract(a)
	}
	return nil
}
