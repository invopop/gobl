package bill

import (
	"fmt"
	"strings"

	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
)

// Status type representing the context in which the status was generated.
const (
	StatusTypeResponse cbc.Key = "response"
	StatusTypeUpdate   cbc.Key = "update"
	StatusTypeSystem   cbc.Key = "system"
)

// Status Line keys representing the specific status being reported.
const (
	StatusLineIssued       cbc.Key = "issued"
	StatusLineAcknowledged cbc.Key = "acknowledged"
	StatusLineProcessing   cbc.Key = "processing"
	StatusLineQuerying     cbc.Key = "querying"
	StatusLineRejected     cbc.Key = "rejected"
	StatusLineAccepted     cbc.Key = "accepted"
	StatusLinePaid         cbc.Key = "paid" // only if a bill.Payment cannot be used
	StatusLineError        cbc.Key = "error"
	StatusLineOther        cbc.Key = "other"
)

// Reason keys as used to represent the reason for a specific status line key.
const (
	ReasonKeyNone            cbc.Key = "none"
	ReasonKeyReferences      cbc.Key = "references"
	ReasonKeyLegal           cbc.Key = "legal"
	ReasonKeyUnknownReceiver cbc.Key = "unknown-receiver"
	ReasonKeyQuality         cbc.Key = "quality"
	ReasonKeyDelivery        cbc.Key = "delivery"
	ReasonKeyPrices          cbc.Key = "prices"
	ReasonKeyQuantity        cbc.Key = "quantity"
	ReasonKeyItems           cbc.Key = "items"
	ReasonKeyPaymentTerms    cbc.Key = "payment-terms"
	ReasonKeyNotRecognized   cbc.Key = "not-recognized"
	ReasonKeyFinanceTerms    cbc.Key = "finance-terms"
	ReasonKeyPartial         cbc.Key = "partial"
	ReasonKeyOther           cbc.Key = "other"
)

// Action keys are used to suggest to the recipient what they should do with the next
// document.
const (
	ActionKeyNone          cbc.Key = "none"
	ActionKeyProvide       cbc.Key = "provide"
	ActionKeyReissue       cbc.Key = "reissue"
	ActionKeyCreditFull    cbc.Key = "credit-full"
	ActionKeyCreditPartial cbc.Key = "credit-partial"
	ActionKeyCreditAmount  cbc.Key = "credit-amount"
	ActionKeyOther         cbc.Key = "other"
)

// StatusTypes describes the different types of status messages that can be issued.
var StatusTypes = []*cbc.Definition{
	{
		Key: StatusTypeResponse,
		Name: i18n.String{
			i18n.EN: "Response",
		},
		Desc: i18n.String{
			i18n.EN: "A response to a document that has been submitted, often used in a two way communication between a customer and supplier.",
		},
	},
	{
		Key: StatusTypeUpdate,
		Name: i18n.String{
			i18n.EN: "Update",
		},
		Desc: i18n.String{
			i18n.EN: "Issued by the supplier or seller that will be shared with either the customer, a fifth-corner entity such as a government agency, or both.",
		},
	},
	{
		Key: StatusTypeSystem,
		Name: i18n.String{
			i18n.EN: "System",
		},
		Desc: i18n.String{
			i18n.EN: "A system event that is not directly related to a specific document but needs to be recorded for compliance purposes.",
		},
	},
}

// StatusLineKeys describes the different status line keys that can be reported.
var StatusLineKeys = []*cbc.Definition{
	{
		Key: StatusLineIssued,
		Name: i18n.String{
			i18n.EN: "Issued",
		},
		Desc: i18n.String{
			i18n.EN: "Document has been submitted pending review by the recipient.",
		},
	},
	{
		Key: StatusLineAcknowledged,
		Name: i18n.String{
			i18n.EN: "Acknowledged",
		},
		Desc: i18n.String{
			i18n.EN: "Received a readable invoice message that can be understood and submitted for processing by the Buyer.",
		},
	},
	{
		Key: StatusLineProcessing,
		Name: i18n.String{
			i18n.EN: "In Process",
		},
		Desc: i18n.String{
			i18n.EN: "Indicates that the referenced message or transaction is being processed.",
		},
	},
	{
		Key: StatusLineQuerying,
		Name: i18n.String{
			i18n.EN: "Under Query",
		},
		Desc: i18n.String{
			i18n.EN: "Buyer will not proceed to accept the Invoice without receiving additional information from the Seller.",
		},
	},
	{
		Key: StatusLineRejected,
		Name: i18n.String{
			i18n.EN: "Rejected",
		},
		Desc: i18n.String{
			i18n.EN: "Buyer will not process the referenced Invoice any further. Buyer is rejecting this invoice but not necessarily the commercial transaction. Although it can be used also for rejection for commercial reasons (invoice not corresponding to delivery).",
		},
	},
	{
		Key: StatusLineAccepted,
		Name: i18n.String{
			i18n.EN: "Accepted / Approved",
		},
		Desc: i18n.String{
			i18n.EN: "Buyer has given a final approval of the invoice and the next step is payment.",
		},
	},
	{
		Key: StatusLinePaid,
		Name: i18n.String{
			i18n.EN: "Paid",
		},
		Desc: i18n.String{
			i18n.EN: "Buyer has initiated payment, or the supplier has acknowledged receipt of payment.",
		},
	},
	{
		Key: StatusLineError,
		Name: i18n.String{
			i18n.EN: "Error",
		},
		Desc: i18n.String{
			i18n.EN: "There is a technical issue with the document that has caused it to be rejected.",
		},
	},
	{
		Key: StatusLineOther,
		Name: i18n.String{
			i18n.EN: "Other",
		},
		Desc: i18n.String{
			i18n.EN: "Status to be determined by other codes or details listed in the line.",
		},
	},
}

// ReasonKeys describes the different reasons for a status event.
var ReasonKeys = []*cbc.Definition{
	{
		Key: ReasonKeyNone,
		Name: i18n.String{
			i18n.EN: "No Issue",
		},
		Desc: i18n.String{
			i18n.EN: "Receiver of the documents sends the message just to update the status and there are no problems with document processing.",
		},
	},
	{
		Key: ReasonKeyReferences,
		Name: i18n.String{
			i18n.EN: "References Incorrect",
		},
		Desc: i18n.String{
			i18n.EN: "Received document did not contain references as required by the receiver for correctly routing the document for approval or processing.",
		},
	},
	{
		Key: ReasonKeyLegal,
		Name: i18n.String{
			i18n.EN: "Legal information incorrect",
		},
		Desc: i18n.String{
			i18n.EN: "Information in the received document is not according to legal requirements.",
		},
	},
	{
		Key: ReasonKeyUnknownReceiver,
		Name: i18n.String{
			i18n.EN: "Receiver unknown",
		},
		Desc: i18n.String{
			i18n.EN: "The party to which the document is addressed is not known.",
		},
	},
	{
		Key: ReasonKeyQuality,
		Name: i18n.String{
			i18n.EN: "Item quality issue",
		},
		Desc: i18n.String{
			i18n.EN: "Unacceptable or incorrect quality.",
		},
	},
	{
		Key: ReasonKeyDelivery,
		Name: i18n.String{
			i18n.EN: "Delivery issues",
		},
		Desc: i18n.String{
			i18n.EN: "Delivery proposed or provided is not acceptable.",
		},
	},
	{
		Key: ReasonKeyPrices,
		Name: i18n.String{
			i18n.EN: "Prices incorrect",
		},
		Desc: i18n.String{
			i18n.EN: "Prices not according to previous expectation.",
		},
	},
	{
		Key: ReasonKeyQuantity,
		Name: i18n.String{
			i18n.EN: "Quantity incorrect",
		},
		Desc: i18n.String{
			i18n.EN: "Quantity not according to previous expectation.",
		},
	},
	{
		Key: ReasonKeyItems,
		Name: i18n.String{
			i18n.EN: "Items incorrect",
		},
		Desc: i18n.String{
			i18n.EN: "Items not according to previous expectation.",
		},
	},
	{
		Key: ReasonKeyPaymentTerms,
		Name: i18n.String{
			i18n.EN: "Payment terms incorrect",
		},
		Desc: i18n.String{
			i18n.EN: "Payment terms not according to previous expectation.",
		},
	},
	{
		Key: ReasonKeyNotRecognized,
		Name: i18n.String{
			i18n.EN: "Not recognized",
		},
		Desc: i18n.String{
			i18n.EN: "Commercial transaction not recognized.",
		},
	},
	{
		Key: ReasonKeyFinanceTerms,
		Name: i18n.String{
			i18n.EN: "Finance incorrect",
		},
		Desc: i18n.String{
			i18n.EN: "Finance terms not according to previous expectation.",
		},
	},
	{
		Key: ReasonKeyPartial,
		Name: i18n.String{
			i18n.EN: "Partially paid",
		},
		Desc: i18n.String{
			i18n.EN: "Payment is partially but not fully paid.",
		},
	},
	{
		Key: ReasonKeyOther,
		Name: i18n.String{
			i18n.EN: "Other",
		},
		Desc: i18n.String{
			i18n.EN: "Reason for status is not defined by code.",
		},
	},
}

// ActionKeys describes the different actions that can be suggested to the recipient.
var ActionKeys = []*cbc.Definition{
	{
		Key: ActionKeyNone,
		Name: i18n.String{
			i18n.EN: "No action required",
		},
		Desc: i18n.String{
			i18n.EN: "No action required.",
		},
	},
	{
		Key: ActionKeyProvide,
		Name: i18n.String{
			i18n.EN: "Provide information",
		},
		Desc: i18n.String{
			i18n.EN: "Missing information requested without re-issuing invoice.",
		},
	},
	{
		Key: ActionKeyReissue,
		Name: i18n.String{
			i18n.EN: "Issue new invoice",
		},
		Desc: i18n.String{
			i18n.EN: "Request to re-issue a corrected invoice.",
		},
	},
	{
		Key: ActionKeyCreditFull,
		Name: i18n.String{
			i18n.EN: "Credit fully",
		},
		Desc: i18n.String{
			i18n.EN: "Request to fully cancel the referenced invoice with a credit note.",
		},
	},
	{
		Key: ActionKeyCreditPartial,
		Name: i18n.String{
			i18n.EN: "Credit partially",
		},
		Desc: i18n.String{
			i18n.EN: "Request to issue partial credit note for corrections only.",
		},
	},
	{
		Key: ActionKeyCreditAmount,
		Name: i18n.String{
			i18n.EN: "Credit the amount",
		},
		Desc: i18n.String{
			i18n.EN: "Request to repay the amount paid on the invoice.",
		},
	},
	{
		Key: ActionKeyOther,
		Name: i18n.String{
			i18n.EN: "Other",
		},
		Desc: i18n.String{
			i18n.EN: "Requested action is not defined by code.",
		},
	},
}

var isValidStatusType = cbc.InKeyDefs(StatusTypes)
var isValidStatusLineKey = cbc.InKeyDefs(StatusLineKeys)
var isValidReasonKey = cbc.InKeyDefs(ReasonKeys)
var isValidActionKey = cbc.InKeyDefs(ActionKeys)

// StatusTypeIn returns a test that passes when the Status's Type is
// one of the provided values. Intended as a guard inside rules.When
// when an addon needs per-type rule branches.
func StatusTypeIn(types ...cbc.Key) rules.Test {
	return is.Func(
		fmt.Sprintf("status type in [%s]", strings.Join(cbc.KeyStrings(types), ", ")),
		func(obj any) bool {
			st, ok := obj.(*Status)
			if !ok || st == nil {
				return false
			}
			return st.Type.In(types...)
		},
	)
}

// StatusLineKeyIn returns a test that passes when the StatusLine's
// Key is one of the provided values. Intended as a guard inside
// rules.When inside a rules.Each over Status.Lines.
func StatusLineKeyIn(keys ...cbc.Key) rules.Test {
	return is.Func(
		fmt.Sprintf("status line key in [%s]", strings.Join(cbc.KeyStrings(keys), ", ")),
		func(obj any) bool {
			line, ok := obj.(*StatusLine)
			if !ok || line == nil {
				return false
			}
			return line.Key.In(keys...)
		},
	)
}

// Status represents a system or business event that needs to be recorded
// for compliance purposes. It is intentionally minimal for now; additional
// fields will be added as more use cases are supported.
//
// Much of the status model is based on the Peppol "InvoiceResponse" model,
// and extended potentially with use for reference to other models such as
// bill.Order or bill.Delivery.
type Status struct {
	tax.Regime    `json:",inline"`
	tax.Addons    `json:",inline"`
	uuid.Identify `json:",inline"`

	// Type of status being reported (e.g. "system" for internal events).
	Type cbc.Key `json:"type" jsonschema:"title=Type"`

	// IssueDate is the date when the status is to be considered effective.
	IssueDate cal.Date `json:"issue_date" jsonschema:"title=Issue Date"`

	// IssueTime is used when extra precision is required to determine when exactly
	// the status was issued.
	IssueTime *cal.Time `json:"issue_time,omitempty" jsonschema:"title=Issue Time"`

	// Series is an optional code to group related status events together.
	Series cbc.Code `json:"series,omitempty" jsonschema:"title=Series"`

	// Code provides a way to identify the specific status event being reported.
	Code cbc.Code `json:"code" jsonschema:"title=Code"`

	// Ext provides additional structured data specific to the regime or addon.
	Ext tax.Extensions `json:"ext,omitzero" jsonschema:"title=Extensions"`

	// Supplier represents the entity supplying the goods or services in the
	// original transaction.
	Supplier *org.Party `json:"supplier" jsonschema:"title=Supplier"`

	// Customer is optional and describes the recipient of the original
	// services.
	Customer *org.Party `json:"customer,omitempty" jsonschema:"title=Customer"`

	// Ordering provides links to related documents and additional details about
	// which parties may be involved in the transaction.
	Ordering *Ordering `json:"ordering,omitempty" jsonschema:"title=Ordering"`

	// Lines contain the main payload of the message used to describe individual
	// documents which have a status.
	Lines []*StatusLine `json:"lines" jsonschema:"title=Lines"`

	// Notes for additional details about the event.
	Notes []*org.Note `json:"notes,omitempty" jsonschema:"title=Notes"`

	// Meta contains unstructured data useful for internal tools.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// StatusLine represents a single row inside a message, that describes the
// situation of another business document.
type StatusLine struct {
	// Position of the row inside the message, determined automatically.
	Index int `json:"index" jsonschema:"title=Index"`

	// Status Key indicates the situation of the document
	Key cbc.Key `json:"key" jsonschema:"title=Key"`

	// When this row should be considered effective
	Date *cal.Date `json:"date,omitempty" jsonschema:"title=Date"`

	// Document reference or details about the document that needs to be looked
	// at.
	Doc *org.DocumentRef `json:"doc,omitempty" jsonschema:"title=Document"`

	// Description includes a human readable description that explains the
	// reason for the current status, if necessary.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`

	// Reasons define an array of reason objects that help the recipient
	// determine why the status was provided.
	Reasons []*Reason `json:"reasons,omitempty" jsonschema:"title=Reasons"`

	// Actions contains an array of actions that should be carried out by the
	// recipient of the message. These are suggestions.
	Actions []*Action `json:"actions,omitempty" jsonschema:"title=Actions"`

	// Extensions for local or format focussed data
	Ext tax.Extensions `json:"ext,omitzero" jsonschema:"title=Extensions"`

	// Complements contain regime/addon specific payload data.
	Complements []*schema.Object `json:"complements,omitempty" jsonschema:"title=Complements"`

	// Additional data specific for the source system.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Reason helps the recipient of a message determine why they are receiving it.
type Reason struct {
	// Key helps identify the reason.
	Key cbc.Key `json:"key" jsonschema:"title=Key"`

	// Description contains a simple text that describes the reason why the
	// original document was not processed.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`

	// Faults provides more specific details about what cause the document
	// to be rejected.
	Faults []*Fault `json:"faults,omitempty" jsonschema:"title=Faults"`

	// Extensions for local or format focussed data
	Ext tax.Extensions `json:"ext,omitzero" jsonschema:"title=Extensions"`
}

// Action provides a suggestion about what to do next with the document.
type Action struct {
	// Key helps determine what to do next.
	Key cbc.Key `json:"key" jsonschema:"title=Key"`

	// Description includes human readable details about what steps should be
	// taken next.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`

	// Extensions for local or format focussed data
	Ext tax.Extensions `json:"ext,omitzero" jsonschema:"title=Extensions"`
}

// Fault provides a more formal structure for describing with a specific
// code what has been unacceptable about the source document, including
// potentially references to the fields causing issues.
type Fault struct {
	// Code or business term provided by the system that raised the condition. These
	// should be searchable in order to help users or systems understand what went
	// wrong with the document.
	Code cbc.Code `json:"code" jsonschema:"title=Code"`

	// Message contains human readable details about the specific condition.
	Message string `json:"message,omitempty" jsonschema:"title=Message"`

	// Paths contains an array of JSON paths that maps the GOBL specific error
	// to a field inside the envelope that the condition is applied to.
	Paths []string `json:"paths,omitempty" jsonschema:"title=Paths"`
}

// CanSign returns a boolean indicating whether the status is ready to be signed
// or not.
func (st *Status) CanSign() bool {
	return st != nil && !st.Code.IsEmpty() && !st.IssueDate.IsZero()
}

// Calculate performs all the normalizations and calculations required for
// the status document.
func (st *Status) Calculate() error {
	if st.Regime.IsEmpty() {
		st.SetRegime(partyTaxCountry(st.Supplier))
	}
	// Track which addon normalizers have already been applied so the
	// follow-up passes only run normalizers for newly-added addons.
	seen := make(map[cbc.Key]bool)
	for _, def := range st.AddonDefs() {
		if def != nil {
			seen[def.Key] = true
		}
	}
	st.Normalize(st.normalizers())
	for pass := 0; pass < maxAddonResolutionPasses; pass++ {
		newNorms := tax.ExtractNormalizersForNew(st, seen)
		if len(newNorms) == 0 {
			break
		}
		st.Normalize(newNorms)
	}
	return st.calculate()
}

// Normalize is run as part of the Calculate method to ensure that the status
// is in a consistent state. This will leverage any add-ons alongside the tax
// regime.
func (st *Status) Normalize(normalizers tax.Normalizers) {
	st.Series = cbc.NormalizeCode(st.Series)
	st.Code = cbc.NormalizeCode(st.Code)
	st.Ext = st.Ext.Clean()

	tax.Normalize(normalizers, st.Supplier)
	tax.Normalize(normalizers, st.Customer)
	tax.Normalize(normalizers, st.Ordering)
	tax.Normalize(normalizers, st.Lines)
	tax.Normalize(normalizers, st.Notes)

	normalizers.Each(st)
}

func (st *Status) normalizers() tax.Normalizers {
	normalizers := make(tax.Normalizers, 0)
	if r := st.RegimeDef(); r != nil {
		normalizers = normalizers.Append(r.Normalizer)
	}
	for _, a := range st.AddonDefs() {
		normalizers = normalizers.Append(a.Normalizer)
	}
	return normalizers
}

func (st *Status) calculate() error {
	// Autofill the issue date when not provided. The issue time is optional.
	if st.IssueDate.IsZero() {
		st.IssueDate = cal.Today()
	}

	// Index lines
	for i, l := range st.Lines {
		if l == nil {
			continue
		}
		l.Index = i + 1
		// Complements
		if err := calculateComplements(l.Complements); err != nil {
			return fmt.Errorf("complements: %w", err)
		}
	}

	return nil
}

// FromEndpoint returns the endpoint of the party most likely to be
// sending this status document. A `response` flows from the customer
// back to the supplier; an `update` flows the other way. A `system`
// status (third-party / clearance system) has no inherent direction
// so this returns nil.
func (st *Status) FromEndpoint() *org.Endpoint {
	if st == nil {
		return nil
	}
	switch st.Type {
	case StatusTypeResponse:
		return st.Customer.FirstEndpoint()
	case StatusTypeUpdate:
		return st.Supplier.FirstEndpoint()
	}
	return nil
}

// ToEndpoint returns the endpoint of the party most likely to be
// receiving this status document. See FromEndpoint for direction.
func (st *Status) ToEndpoint() *org.Endpoint {
	if st == nil {
		return nil
	}
	switch st.Type {
	case StatusTypeResponse:
		return st.Supplier.FirstEndpoint()
	case StatusTypeUpdate:
		return st.Customer.FirstEndpoint()
	}
	return nil
}

// Normalize normalizes the status line's sub-objects.
func (sl *StatusLine) Normalize(normalizers tax.Normalizers) {
	if sl == nil {
		return
	}
	sl.Ext = sl.Ext.Clean()
	tax.Normalize(normalizers, sl.Doc)
	tax.Normalize(normalizers, sl.Reasons)
	tax.Normalize(normalizers, sl.Actions)
	normalizers.Each(sl)
}

// Normalize normalizes the reason's sub-objects.
func (r *Reason) Normalize(normalizers tax.Normalizers) {
	if r == nil {
		return
	}
	r.Ext = r.Ext.Clean()
	tax.Normalize(normalizers, r.Faults)
	normalizers.Each(r)
}

// Normalize runs any registered normalizers on the action.
func (a *Action) Normalize(normalizers tax.Normalizers) {
	if a == nil {
		return
	}
	a.Ext = a.Ext.Clean()
	normalizers.Each(a)
}

// Normalize normalizes the fault's code and runs any registered normalizers.
func (f *Fault) Normalize(normalizers tax.Normalizers) {
	if f == nil {
		return
	}
	f.Code = cbc.NormalizeCode(f.Code)
	normalizers.Each(f)
}

// JSONSchemaExtend extends the schema with additional property details
func (Status) JSONSchemaExtend(js *jsonschema.Schema) {
	props := js.Properties
	if its, ok := props.Get("type"); ok {
		its.OneOf = make([]*jsonschema.Schema, len(StatusTypes))
		for i, kd := range StatusTypes {
			its.OneOf[i] = &jsonschema.Schema{
				Const:       kd.Key.String(),
				Title:       kd.Name.String(),
				Description: kd.Desc.String(),
			}
		}
	}
	js.Extras = map[string]any{
		schema.Recommended: []string{
			"$regime",
			"series",
			"lines",
		},
	}
}

// JSONSchemaExtend extends the schema with additional property details
func (StatusLine) JSONSchemaExtend(js *jsonschema.Schema) {
	props := js.Properties
	if its, ok := props.Get("key"); ok {
		// Status line keys are a closed set: only the predefined keys are
		// permitted, so enumerate them with OneOf (no open-ended fallback).
		its.OneOf = make([]*jsonschema.Schema, len(StatusLineKeys))
		for i, kd := range StatusLineKeys {
			its.OneOf[i] = &jsonschema.Schema{
				Const:       kd.Key.String(),
				Title:       kd.Name.String(),
				Description: kd.Desc.String(),
			}
		}
	}
}

// JSONSchemaExtend extends the schema with additional property details
func (Reason) JSONSchemaExtend(js *jsonschema.Schema) {
	props := js.Properties
	if its, ok := props.Get("key"); ok {
		its.OneOf = make([]*jsonschema.Schema, len(ReasonKeys))
		for i, kd := range ReasonKeys {
			its.OneOf[i] = &jsonschema.Schema{
				Const:       kd.Key.String(),
				Title:       kd.Name.String(),
				Description: kd.Desc.String(),
			}
		}
	}
}

// JSONSchemaExtend extends the schema with additional property details
func (Action) JSONSchemaExtend(js *jsonschema.Schema) {
	props := js.Properties
	if its, ok := props.Get("key"); ok {
		its.OneOf = make([]*jsonschema.Schema, len(ActionKeys))
		for i, kd := range ActionKeys {
			its.OneOf[i] = &jsonschema.Schema{
				Const:       kd.Key.String(),
				Title:       kd.Name.String(),
				Description: kd.Desc.String(),
			}
		}
	}
}

func statusRules() *rules.Set {
	return rules.For(new(Status),
		rules.Field("type",
			rules.Assert("01", "status type is required",
				is.Present,
			),
			rules.Assert("02", "status type is not valid",
				isValidStatusType,
			),
		),
		rules.Field("supplier",
			rules.Assert("03", "status supplier is required",
				is.Present,
			),
		),
		rules.Field("lines",
			rules.Assert("04", "status must have at least one line",
				is.Present,
			),
		),
	)
}

func statusLineRules() *rules.Set {
	return rules.For(new(StatusLine),
		rules.Field("key",
			rules.Assert("01", "status line key is required",
				is.Present,
			),
			rules.Assert("02", "status line key is not valid",
				isValidStatusLineKey,
			),
		),
	)
}

func reasonRules() *rules.Set {
	return rules.For(new(Reason),
		rules.Field("key",
			rules.Assert("01", "reason key is required",
				is.Present,
			),
			rules.Assert("02", "reason key is not valid",
				isValidReasonKey,
			),
		),
	)
}

func actionRules() *rules.Set {
	return rules.For(new(Action),
		rules.Field("key",
			rules.Assert("01", "action key is required",
				is.Present,
			),
			rules.Assert("02", "action key is not valid",
				isValidActionKey,
			),
		),
	)
}

func faultRules() *rules.Set {
	return rules.For(new(Fault),
		rules.Field("code",
			rules.Assert("01", "fault code is required",
				is.Present,
			),
		),
	)
}
