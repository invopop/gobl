package org

import "github.com/invopop/gobl/num"

// Discount represents an amount and/or percentage that can be applied to
// a given price.
//
// TODO: use the UNTDID 5189 code list for Discount Reason Code.
//
type Discount struct {
	Rate   num.Percentage `json:"rate,omitempty" jsonschema:"title=Rate"`
	Value  num.Amount     `json:"val" jsonschema:"title=Value,description=How much to deduct"`
	Reason string         `json:"reason,omitempty" jsonschema:"title=Reason,description=Description as to why this discount was applied."`
	Code   string         `json:"code,omitempty" jsonschema:"title=Code,description=Reason Code"`
}
