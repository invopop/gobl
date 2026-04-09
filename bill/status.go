package bill

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
)

// Status represents a system or business event that needs to be recorded
// for compliance purposes. It is intentionally minimal for now; additional
// fields will be added as more use cases are supported.
type Status struct {
	tax.Regime    `json:",inline"`
	tax.Addons    `json:",inline"`
	uuid.Identify `json:",inline"`

	// Type of status being reported (e.g. "system" for internal events).
	Type cbc.Key `json:"type" jsonschema:"title=Type"`

	// IssueDate is the date when the status event occurred.
	IssueDate cal.Date `json:"issue_date" jsonschema:"title=Issue Date"`

	// Supplier is the entity that is reporting the event.
	Supplier *org.Party `json:"supplier" jsonschema:"title=Supplier"`

	// Ext provides additional structured data specific to the regime or addon.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`

	// Complements contain regime/addon specific payload data.
	Complements []*schema.Object `json:"complements,omitempty" jsonschema:"title=Complements"`

	// Notes for additional details about the event.
	Notes []*org.Note `json:"notes,omitempty" jsonschema:"title=Notes"`

	// Meta contains unstructured data useful for internal tools.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}
