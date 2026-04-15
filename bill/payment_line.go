package bill

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
)

// PaymentLine defines the details of a line item in a payment document.
type PaymentLine struct {
	uuid.Identify

	// Line number within the parent document (automatically calculated)
	Index int `json:"i" jsonschema:"title=Index" jsonschema_extras:"calculated=true"`

	// Indicates whether this payment is a refund of a previous payment, effectively reversing
	// the flow of funds between the supplier and customer or their representatives.
	Refund bool `json:"refund,omitempty" jsonschema:"title=Refund"`

	// Reference to the document being paid
	Document *org.DocumentRef `json:"document,omitempty" jsonschema:"title=Document"`

	// When making multiple payments for a single document, this specifies the
	// installment number for this payment line.
	Installment int `json:"installment,omitempty" jsonschema:"title=Installment"`

	// Additional human readable description of the payment line which may be useful for
	// explaining the purpose or special conditions. Notes should be used for more
	// formal comments.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`

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

func paymentLineRules() *rules.Set {
	return rules.For(new(PaymentLine),
		rules.Field("installment",
			rules.Assert("01", "installment must be between 1 and 999",
				is.Min(1), is.Max(999),
			),
		),
		rules.Field("payable",
			rules.Assert("02", "payable must be positive",
				num.Positive,
			),
		),
		rules.Field("advances",
			rules.Assert("03", "advances must be zero or positive",
				num.ZeroOrPositive,
			),
		),
		rules.Field("amount",
			rules.Assert("04", "amount must be positive",
				num.Positive,
			),
		),
		rules.Field("due",
			rules.Assert("05", "due must be zero or positive",
				num.ZeroOrPositive,
			),
		),
		rules.Assert("11", "advances must not exceed payable",
			is.Func("advances within payable", paymentLineAdvancesWithinPayable),
		),
		rules.Assert("12", "amount must not exceed payable less advances",
			is.Func("amount within limit", paymentLineAmountWithinLimit),
		),
	)
}

func paymentLineAdvancesWithinPayable(val any) bool {
	pl, ok := val.(*PaymentLine)
	if !ok {
		return false
	}
	if pl.Payable == nil || pl.Advances == nil {
		// nothing to compare, assume okay
		return true
	}
	return pl.Advances.Compare(*pl.Payable) <= 0
}

func paymentLineAmountWithinLimit(val any) bool {
	pl, ok := val.(*PaymentLine)
	if !ok {
		return false
	}
	if pl == nil || pl.Payable == nil {
		return true
	}
	mx := *pl.Payable
	if pl.Advances != nil {
		mx = mx.Subtract(*pl.Advances)
	}
	if !mx.IsPositive() {
		return true
	}
	return pl.Amount.Compare(mx) <= 0
}

func (pl *PaymentLine) calculate(rates []*currency.ExchangeRate, cur currency.Code, rr cbc.Key) error {
	var lt *tax.Total

	if pl.Document != nil {
		var er *currency.ExchangeRate
		if dc := pl.Document.Currency; dc != currency.CodeEmpty {
			// If the document has a currency, we need to ensure there is an exchange
			// rate so any taxes can be converted correctly.
			if er = currency.MatchExchangeRate(rates, dc, cur); er == nil {
				return fmt.Errorf("document: currency: missing exchange rate from %s to %s", dc, cur)
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
	} else if pl.Tax != nil {
		// Ensure the taxes are calculated correctly.
		pl.Tax.Calculate(cur, rr)
	}

	return nil
}
