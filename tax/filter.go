package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// Filter selects rate totals from a document's tax totals by category and,
// optionally, key and extensions.
type Filter struct {
	// Tax category code the rates must belong to.
	Category cbc.Code `json:"cat" jsonschema:"title=Category"`
	// Tax key the rates must have, if set.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Extensions the rates must include, if set.
	Ext Extensions `json:"ext,omitzero" jsonschema:"title=Extensions"`
}

func filterRules() *rules.Set {
	return rules.For(new(Filter),
		rules.Field("cat",
			rules.Assert("01", "tax filter category is required", is.Present),
		),
	)
}

func (f *Filter) matches(cat cbc.Code, rt *RateTotal) bool {
	if f == nil || f.Category != cat {
		return false
	}
	if f.Key != cbc.KeyEmpty && f.Key != rt.Key {
		return false
	}
	if !f.Ext.IsZero() && !rt.Ext.Contains(f.Ext) {
		return false
	}
	return true
}

// FilteredAmount sums the amounts, excluding surcharges, of the rate totals
// that match any of the filters.
func (t *Total) FilteredAmount(zero num.Amount, filters []*Filter) num.Amount {
	sum := zero
	if t == nil {
		return sum
	}
	for _, ct := range t.Categories {
		for _, rt := range ct.Rates {
			for _, f := range filters {
				if f.matches(ct.Code, rt) {
					sum = sum.Add(rt.Amount)
					break
				}
			}
		}
	}
	return sum
}
