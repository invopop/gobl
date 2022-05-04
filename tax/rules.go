package tax

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
)

// Action defines what should happen with a rule if it matches
type Action string

const (
	// Allow implies the supplier rates should be applied.
	AllowAction Action = "allow"
	// Deny means that the combination should not be allowed and will raise an error.
	DenyAction Action = "deny"
	// Forward means that the taxes at the destination with a matching category
	// and type should be applied.
	ForwardAction Action = "forward"
	// Skip indicates that no rates in this category should be applied.
	SkipAction Action = "skip"
)

// Rules defines a set of Rule conditions
type Rules []*Rule

// Rule defines a map of conditions
type Rule struct {
	// Since when should this rule be considered valid
	Since *cal.Date `json:"since,omitempty"`

	// TODO!!! Consider using a "tax Identity" for matching conditions

	// The counter party that should be matched.
	Dest []l10n.Code `json:"dest"`

	// When true, rule will only match valid companies
	Company bool `json:"company,omitempty"`

	// What to do if there is a match.
	Action Action `json:"action"`

	// Why was this action chosen?
	Reason *i18n.String `json:"reason,omitempty"`
}
