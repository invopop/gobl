package tax

import (
	"errors"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// ThresholdDef defines a monetary limit that document amounts are compared
// against, such as the total above which an invoice must identify the
// customer. Values change over time, so each one applies from a given date.
type ThresholdDef struct {
	// Key used to identify the threshold in rules.
	Key cbc.Key `json:"key" jsonschema:"title=Key"`
	// Human name of the threshold.
	Name i18n.String `json:"name" jsonschema:"title=Name"`
	// Useful description of the threshold.
	Description i18n.String `json:"desc,omitempty" jsonschema:"title=Description"`
	// Sources of the threshold's values.
	Sources []*cbc.Source `json:"sources,omitempty" jsonschema:"title=Sources"`
	// Currency the threshold amounts are expressed in.
	Currency currency.Code `json:"currency" jsonschema:"title=Currency"`
	// Values contains the current and historical amounts of the threshold,
	// ordered newest first.
	Values []*ThresholdValueDef `json:"values" jsonschema:"title=Values"`
}

// ThresholdValueDef contains the threshold amount from a given date.
type ThresholdValueDef struct {
	// Date from which this value should be applied.
	Since *cal.Date `json:"since,omitempty" jsonschema:"title=Since"`
	// Amount of the threshold.
	Amount num.Amount `json:"amount" jsonschema:"title=Amount"`
}

func thresholdDefRules() *rules.Set {
	return rules.For(new(ThresholdDef),
		rules.Field("key",
			rules.Assert("01", "key is required", is.Present),
		),
		rules.Field("name",
			rules.Assert("02", "name is required", is.Present),
		),
		rules.Field("currency",
			rules.Assert("03", "currency is required", is.Present),
		),
		rules.Field("values",
			rules.Assert("04", "at least one value is required", is.Present),
			rules.Assert("05", "values must be in descending chronological order with only the last missing a since date",
				is.FuncError("date order", checkThresholdValuesOrder),
			),
		),
	)
}

// Value determines the threshold value that applies on the provided date.
func (td *ThresholdDef) Value(date cal.Date) *ThresholdValueDef {
	if td == nil {
		return nil
	}
	for _, tv := range td.Values {
		if tv.Since == nil || !tv.Since.IsValid() || !tv.Since.After(date.Date) {
			return tv
		}
	}
	return nil
}

// Reached returns true when the amount is equal to or greater than the
// threshold value on the provided date.
func (td *ThresholdDef) Reached(date cal.Date, amount num.Amount) bool {
	tv := td.Value(date)
	return tv != nil && amount.Compare(tv.Amount) >= 0
}

// Exceeded returns true when the amount is greater than the threshold
// value on the provided date.
func (td *ThresholdDef) Exceeded(date cal.Date, amount num.Amount) bool {
	tv := td.Value(date)
	return tv != nil && amount.Compare(tv.Amount) > 0
}

func checkThresholdValuesOrder(list any) error {
	values, ok := list.([]*ThresholdValueDef)
	if !ok {
		return errors.New("must be a threshold value array")
	}
	for i, v := range values {
		if v == nil {
			return errors.New("value must not be empty")
		}
		if v.Since == nil {
			if i < len(values)-1 {
				return errors.New("only the last value may have no since date")
			}
			continue
		}
		if i > 0 && !v.Since.Before(values[i-1].Since.Date) {
			return errors.New("invalid date order")
		}
	}
	return nil
}
