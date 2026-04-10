package bill

import (
	"fmt"

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

// Status event keys representing the specific status being reported.
const (
	StatusEventIssued       cbc.Key = "issued"
	StatusEventAcknowledged cbc.Key = "acknowledged"
	StatusEventProcessing   cbc.Key = "processing"
	StatusEventQuerying     cbc.Key = "querying"
	StatusEventRejected     cbc.Key = "rejected"
	StatusEventAccepted     cbc.Key = "accepted"
	StatusEventPaid         cbc.Key = "paid"
	StatusEventError        cbc.Key = "error"
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
			i18n.EN: "Issued by the supplier (or issuer) that may be shared with the customer, but more likely with a fifth-corner, like a government agency.",
		},
	},
	{
		Key: StatusTypeSystem,
		Name: i18n.String{
			i18n.EN: "System",
		},
		Desc: i18n.String{
			i18n.EN: "A system event that is not directly related to a specific document, but needs to be recorded for compliance purposes.",
		},
	},
}

// StatusEvents describes the different status events that can be reported.
var StatusEvents = []*cbc.Definition{
	{
		Key: StatusEventIssued,
		Name: i18n.String{
			i18n.EN: "Issued",
		},
		Desc: i18n.String{
			i18n.EN: "Document has been submitted pending review by the recipient.",
		},
	},
	{
		Key: StatusEventAcknowledged,
		Name: i18n.String{
			i18n.EN: "Acknowledged",
		},
		Desc: i18n.String{
			i18n.EN: "Received a readable invoice message that can be understood and submitted for processing by the Buyer.",
		},
	},
	{
		Key: StatusEventProcessing,
		Name: i18n.String{
			i18n.EN: "In Process",
		},
		Desc: i18n.String{
			i18n.EN: "Indicates that the referenced message or transaction is being processed.",
		},
	},
	{
		Key: StatusEventQuerying,
		Name: i18n.String{
			i18n.EN: "Under Query",
		},
		Desc: i18n.String{
			i18n.EN: "Buyer will not proceed to accept the Invoice without receiving additional information from the Seller.",
		},
	},
	{
		Key: StatusEventRejected,
		Name: i18n.String{
			i18n.EN: "Rejected",
		},
		Desc: i18n.String{
			i18n.EN: "Buyer will not process the referenced Invoice any further. Buyer is rejecting this invoice but not necessarily the commercial transaction. Although it can be used also for rejection for commercial reasons (invoice not corresponding to delivery).",
		},
	},
	{
		Key: StatusEventAccepted,
		Name: i18n.String{
			i18n.EN: "Accepted / Approved",
		},
		Desc: i18n.String{
			i18n.EN: "Buyer has given a final approval of the invoice and the next step is payment.",
		},
	},
	{
		Key: StatusEventPaid,
		Name: i18n.String{
			i18n.EN: "Paid",
		},
		Desc: i18n.String{
			i18n.EN: "Buyer has initiated payment, or the supplier has acknowledged receipt of payment.",
		},
	},
	{
		Key: StatusEventError,
		Name: i18n.String{
			i18n.EN: "Error",
		},
		Desc: i18n.String{
			i18n.EN: "There is a technical issue with the document that has caused it to be rejected.",
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
var isValidStatusEvent = cbc.InKeyDefs(StatusEvents)
var isValidReasonKey = cbc.InKeyDefs(ReasonKeys)
var isValidActionKey = cbc.InKeyDefs(ActionKeys)

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

	// IssueTime is the time when the status event occurred.
	IssueTime *cal.Time `json:"issue_time,omitempty" jsonschema:"title=Issue Time"`

	// Series is an optional code to group related status events together.
	Series cbc.Code `json:"series,omitempty" jsonschema:"title=Series"`

	// Code provides a way to identify the specific status event being reported.
	Code cbc.Code `json:"code" jsonschema:"title=Code"`

	// Ext provides additional structured data specific to the regime or addon.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`

	// Ordering provides links to related documents and details that may have occurred
	// before this status was created.
	Ordering *Ordering `json:"ordering,omitempty" jsonschema:"title=Ordering"`

	// Supplier represents the entity supplying the goods or services in the
	// original transaction and may not be the issuer of the document.
	Supplier *org.Party `json:"supplier" jsonschema:"title=Supplier"`

	// Customer is optional and describes the recipient of the original
	// services. In the case of a local or system event, this will be empty.
	Customer *org.Party `json:"customer,omitempty" jsonschema:"title=Customer"`

	// Issuer represents an intermediary acting on behalf of either the supplier
	// or customer in order to provide a status update, this is useful
	// specifically if only the third party is registered on a network.
	// (Optional).
	Issuer *org.Party `json:"issuer,omitempty" jsonschema:"title=Issuer"`

	// Recipient represents another intermediary responsible for receiving the
	// event when the supplier or customer do not have networking capabilities.
	// May also be a tax agency in a five corner model.
	// (Optional).
	Recipient *org.Party `json:"recipient,omitempty" jsonschema:"title=Recipient"`

	// Lines contain the main payload of the message used to describe individual
	// documents which have a status. A message may not have any lines.
	// (Optional).
	Lines []*StatusLine `json:"lines,omitempty" jsonschema:"title=Lines"`

	// Complements contain regime/addon specific payload data.
	Complements []*schema.Object `json:"complements,omitempty" jsonschema:"title=Complements"`

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
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`

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

	// Condition provides additional details with codes about what has gone
	// wrong with the incoming document.
	Conditions []*Condition `json:"conditions,omitempty" jsonschema:"title=Conditions"`
}

// Action provides a suggestion about what to do next with the document.
type Action struct {
	// Key helps determine what to do next.
	Key cbc.Key `json:"key" jsonschema:"title=Key"`

	// Description includes human readable details about what steps should be
	// take next.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
}

// Condition provides a more formal structure for describing with a specific
// code what has been unacceptable about the source document, including
// potentially references to the fields causing issues.
type Condition struct {
	// Code is generated by the system that raised the condition.
	Code cbc.Code `json:"code" jsonschema:"title=Code"`

	// Paths contains an array of JSON paths that maps the GOBL specific error
	// to a field inside the envelope that the condition is applied to.
	Paths []string `json:"paths,omitempty" jsonschema:"title=Paths"`

	// Message contains human readable details about the specific condition.
	Message string `json:"message,omitempty" jsonschema:"title=Message"`
}

// Calculate performs all the normalizations and calculations required for
// the status document.
func (st *Status) Calculate() error {
	if st.Regime.IsEmpty() {
		st.SetRegime(partyTaxCountry(st.Supplier))
	}
	st.Normalize(st.normalizers())
	return st.calculate()
}

// Normalize is run as part of the Calculate method to ensure that the status
// is in a consistent state. This will leverage any add-ons alongside the tax
// regime.
func (st *Status) Normalize(normalizers tax.Normalizers) {
	st.Series = cbc.NormalizeCode(st.Series)
	st.Code = cbc.NormalizeCode(st.Code)

	tax.Normalize(normalizers, st.Supplier)
	tax.Normalize(normalizers, st.Customer)
	tax.Normalize(normalizers, st.Issuer)
	tax.Normalize(normalizers, st.Recipient)
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
	r := st.RegimeDef()

	// Set the issue date and time
	tz := r.TimeLocation()
	if st.IssueTime != nil && st.IssueTime.IsZero() {
		tn := cal.ThisSecondIn(tz)
		hn := tn.Time()
		st.IssueDate = tn.Date()
		st.IssueTime = &hn
	} else if st.IssueDate.IsZero() {
		st.IssueDate = cal.TodayIn(tz)
	}

	// Index lines
	for i, l := range st.Lines {
		if l == nil {
			continue
		}
		l.Index = i + 1
	}

	// Complements
	if err := calculateComplements(st.Complements); err != nil {
		return fmt.Errorf("complements: %w", err)
	}

	return nil
}

// Normalize normalizes the status line's sub-objects.
func (sl *StatusLine) Normalize(normalizers tax.Normalizers) {
	if sl == nil {
		return
	}
	tax.Normalize(normalizers, sl.Doc)
	normalizers.Each(sl)
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
		its.OneOf = make([]*jsonschema.Schema, len(StatusEvents))
		for i, kd := range StatusEvents {
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
			rules.Assert("01", "status type is required", is.Present),
			rules.Assert("02", "status type is not valid", isValidStatusType),
		),
		rules.Field("issue_date",
			rules.Assert("03", "status issue date is required", cal.DateNotZero()),
		),
		rules.Field("code",
			rules.Assert("04", "status code is required", is.Present),
		),
		rules.Field("supplier",
			rules.Assert("05", "status supplier is required", is.Present),
		),
	)
}

func statusLineRules() *rules.Set {
	return rules.For(new(StatusLine),
		rules.Field("key",
			rules.Assert("01", "status line key is required", is.Present),
			rules.Assert("02", "status line key is not valid", isValidStatusEvent),
		),
	)
}

func reasonRules() *rules.Set {
	return rules.For(new(Reason),
		rules.Field("key",
			rules.Assert("01", "reason key is required", is.Present),
			rules.Assert("02", "reason key is not valid", isValidReasonKey),
		),
	)
}

func actionRules() *rules.Set {
	return rules.For(new(Action),
		rules.Field("key",
			rules.Assert("01", "action key is required", is.Present),
			rules.Assert("02", "action key is not valid", isValidActionKey),
		),
	)
}

func conditionRules() *rules.Set {
	return rules.For(new(Condition),
		rules.Field("code",
			rules.Assert("01", "condition code is required", is.Present),
		),
	)
}
