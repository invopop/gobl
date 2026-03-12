package bill

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
)

// Event represents a response or acknowledgement to one or more
// previously issued business documents
type Event struct {
	tax.Tags
	uuid.Identify

	Code      cbc.Code `json:"code"`
	Name      string   `json:"name,omitempty"`
	IssueDate cal.Date `json:"issue_date"`
	IssueTime cal.Time `json:"issue_time,omitempty"`

	Sender   *org.Party `json:"sender"`
	Receiver *org.Party `json:"receiver"`

	Docs []*EventDoc `json:"docs"`

	Notes []*org.Note    `json:"notes,omitempty"`
	Ext   tax.Extensions `json:"ext,omitempty"`
	Meta  cbc.Meta       `json:"meta,omitempty"`
}

// EventDoc is the response to a single referenced document.
type EventDoc struct {
	Ref *org.DocumentRef `json:"ref"`

	Status cbc.Code `json:"status"`

	Issuer    *org.Party `json:"issuer,omitempty"`
	Recipient *org.Party `json:"recipient,omitempty"`

	EffectiveDate cal.Date `json:"effective_date,omitempty"`
	EffectiveTime cal.Time `json:"effective_time,omitempty"`

	Statuses []*EventStatus `json:"statuses,omitempty"`

	Ext  tax.Extensions `json:"ext,omitempty"`
	Meta cbc.Meta       `json:"meta,omitempty"`
}

// EventStatus captures a reason, requested action, and optional
// field-level details for a document response.
type EventStatus struct {
	Sequence int `json:"sequence,omitempty"`

	Code   cbc.Code `json:"code,omitempty"`
	Reason string   `json:"reason,omitempty"`

	Fields []*EventField `json:"fields,omitempty"`
}

// EventField captures a field-level or line-level reference within the
// source document, along with its current or expected value.
type EventField struct {
	ID    string   `json:"id,omitempty"`
	Label cbc.Code `json:"label,omitempty"`
	Value string   `json:"value,omitempty"`
}
