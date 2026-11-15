package legal

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/l10n"
)

// Effect describes when a contract starts and ceases to have effect. A
// narrative Trigger supports agreements whose effect depends on an event that
// cannot be reduced to a date.
type Effect struct {
	Start   *cal.Date `json:"start,omitempty" jsonschema:"title=Start Date"`
	End     *cal.Date `json:"end,omitempty" jsonschema:"title=End Date"`
	Trigger string    `json:"trigger,omitempty" jsonschema:"title=Trigger"`
}

// Recital records non-operative background or purpose text.
type Recital struct {
	Anchor  string `json:"$anchor,omitempty" jsonschema:"title=Anchor"`
	Content string `json:"content" jsonschema:"title=Content"`
}

// Definition binds a term used in the contract to its authoritative meaning.
type Definition struct {
	Anchor  string `json:"$anchor,omitempty" jsonschema:"title=Anchor"`
	Term    string `json:"term" jsonschema:"title=Term"`
	Meaning string `json:"meaning" jsonschema:"title=Meaning"`
}

// GoverningLaw records an express choice of law. Jurisdiction is deliberately
// free text because legal jurisdictions do not map cleanly to countries.
type GoverningLaw struct {
	Country      l10n.ISOCountryCode `json:"country,omitempty" jsonschema:"title=Country"`
	Jurisdiction string              `json:"jurisdiction,omitempty" jsonschema:"title=Jurisdiction"`
	Instrument   string              `json:"instrument,omitempty" jsonschema:"title=Legal Instrument"`
}

// Dispute resolution method keys.
const (
	DisputeMethodCourt       cbc.Key = "court"
	DisputeMethodArbitration cbc.Key = "arbitration"
	DisputeMethodMediation   cbc.Key = "mediation"
	DisputeMethodNegotiation cbc.Key = "negotiation"
)

// DisputeResolution describes the agreed process for resolving disputes.
type DisputeResolution struct {
	Method   cbc.Key    `json:"method" jsonschema:"title=Method"`
	Forum    string     `json:"forum,omitempty" jsonschema:"title=Forum"`
	Seat     string     `json:"seat,omitempty" jsonschema:"title=Seat"`
	Rules    string     `json:"rules,omitempty" jsonschema:"title=Rules"`
	Language *i18n.Lang `json:"language,omitempty" jsonschema:"title=Language"`
}
