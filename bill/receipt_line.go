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

// ReceiptLine defines the details of a line required in an invoice.
type ReceiptLine struct {
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

// ValidateWithContext ensures that the fields contained in the ReceiptLine look correct.
func (rl *ReceiptLine) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(ctx, rl,
		validation.Field(&rl.Document),
		validation.Field(&rl.Currency),
		validation.Field(&rl.Debit,
			validation.When(
				rl.Credit == nil,
				validation.Required.Error("must have either debit or credit"),
			),
		),
		validation.Field(&rl.Credit),
		validation.Field(&rl.Tax),
		validation.Field(&rl.Total, validation.Required),
		validation.Field(&rl.Notes),
	)
}

// calculate will ensure the total amount is calculated correctly
func (rl *ReceiptLine) calculate(cur currency.Code, rates []*currency.ExchangeRate) error {
	rl.Total = cur.Def().Zero()
	if rl.Debit != nil {
		var a num.Amount
		if rl.Currency != "" {
			na := currency.Convert(rates, rl.Currency, cur, *rl.Debit)
			if na == nil {
				return validation.Errors{
					"currency": fmt.Errorf("no exchange rate found for %s to %s", rl.Currency, cur),
				}
			}
			a = *na
		} else {
			a = *rl.Debit
		}
		rl.Total.MatchPrecision(a)
		rl.Total = rl.Total.Add(a)
	}
	if rl.Credit != nil {
		var a num.Amount
		if rl.Currency != "" {
			na := currency.Convert(rates, rl.Currency, cur, *rl.Credit)
			if na == nil {
				return validation.Errors{
					"currency": fmt.Errorf("no exchange rate found for %s to %s", rl.Currency, cur),
				}
			}
			a = *na
		} else {
			a = *rl.Credit
		}
		rl.Total.MatchPrecision(a)
		rl.Total = rl.Total.Subtract(a)
	}
	return nil
}
