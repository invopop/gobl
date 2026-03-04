package org

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Predefined list of supported note keys based on the
// UNTDID 4451 list of text subject qualifiers. We've picked the ones
// which we think are most useful, but if you require an additional
// code, please send a pull request.
const (
	// Goods Description
	NoteKeyGoods cbc.Key = "goods"
	// Terms of Payment
	NoteKeyPayment cbc.Key = "payment"
	// Payment method or remittance information
	NoteKeyPaymentMethod cbc.Key = "payment-method"
	// Payment term details
	NoteKeyPaymentTerm cbc.Key = "payment-term"
	// Legal or regulatory information
	NoteKeyLegal cbc.Key = "legal"
	// Dangerous goods additional information
	NoteKeyDangerousGoods cbc.Key = "dangerous-goods"
	// Acknowledgement Description
	NoteKeyAck cbc.Key = "ack"
	// Rate additional information
	NoteKeyRate cbc.Key = "rate"
	// Reason
	NoteKeyReason cbc.Key = "reason"
	// Dispute
	NoteKeyDispute cbc.Key = "dispute"
	// Customer remarks
	NoteKeyCustomer cbc.Key = "customer"
	// Glossary
	NoteKeyGlossary cbc.Key = "glossary"
	// Customs declaration information
	NoteKeyCustoms cbc.Key = "customs"
	// General information
	NoteKeyGeneral cbc.Key = "general"
	// Handling instructions
	NoteKeyHandling cbc.Key = "handling"
	// Packaging information
	NoteKeyPackaging cbc.Key = "packaging"
	// Loading instructions
	NoteKeyLoading cbc.Key = "loading"
	// Price conditions
	NoteKeyPrice cbc.Key = "price"
	// Priority information
	NoteKeyPriority cbc.Key = "priority"
	// Regulatory information
	NoteKeyRegulatory cbc.Key = "regulatory"
	// Safety Instructions
	NoteKeySafety cbc.Key = "safety"
	// Ship Line
	NoteKeyShipLine cbc.Key = "ship-line"
	// Supplier remarks
	NoteKeySupplier cbc.Key = "supplier"
	// Transportation information
	NoteKeyTransport cbc.Key = "transport"
	// Delivery Information
	NoteKeyDelivery cbc.Key = "delivery"
	// Quarantine Information
	NoteKeyQuarantine cbc.Key = "quarantine"
	// Tax declaration
	NoteKeyTax cbc.Key = "tax"
	// Other
	NoteKeyOther cbc.Key = "other"
)

// NoteKeyDefinitions provides a map of Note Keys to their definitions
// including a description and UNTDID code.
var NoteKeyDefinitions = []cbc.Definition{
	{
		Key:  NoteKeyGoods,
		Name: i18n.NewString("Goods"),
		Desc: i18n.NewString("Goods Description"),
	},
	{
		Key:  NoteKeyPayment,
		Name: i18n.NewString("Payment"),
		Desc: i18n.NewString("Terms of Payment"),
	},
	{
		Key:  NoteKeyPaymentMethod,
		Name: i18n.NewString("Payment Method"),
		Desc: i18n.NewString("Payment method or remittance information"),
	},
	{
		Key:  NoteKeyPaymentTerm,
		Name: i18n.NewString("Payment Term"),
		Desc: i18n.NewString("Payment term details"),
	},
	{
		Key:  NoteKeyLegal,
		Name: i18n.NewString("Legal"),
		Desc: i18n.NewString("Legal or regulatory information"),
	},
	{
		Key:  NoteKeyDangerousGoods,
		Name: i18n.NewString("Dangerous Goods"),
		Desc: i18n.NewString("Dangerous goods additional information"),
	},
	{
		Key:  NoteKeyAck,
		Name: i18n.NewString("Acknowledgement"),
		Desc: i18n.NewString("Acknowledgement Description"),
	},
	{
		Key:  NoteKeyRate,
		Name: i18n.NewString("Rate"),
		Desc: i18n.NewString("Rate additional information"),
	},
	{
		Key:  NoteKeyReason,
		Name: i18n.NewString("Reason"),
		Desc: i18n.NewString("Explanation of something relevant to the document"),
	},
	{
		Key:  NoteKeyDispute,
		Name: i18n.NewString("Dispute"),
		Desc: i18n.NewString("Details on a dispute."),
	},
	{
		Key:  NoteKeyCustomer,
		Name: i18n.NewString("Customer"),
		Desc: i18n.NewString("Customer remarks"),
	},
	{
		Key:  NoteKeyGlossary,
		Name: i18n.NewString("Glossary"),
		Desc: i18n.NewString("Glossary of terms"),
	},
	{
		Key:  NoteKeyCustoms,
		Name: i18n.NewString("Customs"),
		Desc: i18n.NewString("Customs declaration information"),
	},
	{
		Key:  NoteKeyGeneral,
		Name: i18n.NewString("General"),
		Desc: i18n.NewString("General information"),
	},
	{
		Key:  NoteKeyHandling,
		Name: i18n.NewString("Handling"),
		Desc: i18n.NewString("Handling instructions"),
	},
	{
		Key:  NoteKeyPackaging,
		Name: i18n.NewString("Packaging"),
		Desc: i18n.NewString("Packaging information"),
	},
	{
		Key:  NoteKeyLoading,
		Name: i18n.NewString("Loading"),
		Desc: i18n.NewString("Loading instructions"),
	},
	{
		Key:  NoteKeyPrice,
		Name: i18n.NewString("Price"),
		Desc: i18n.NewString("Price conditions"),
	},
	{
		Key:  NoteKeyPriority,
		Name: i18n.NewString("Priority"),
		Desc: i18n.NewString("Priority information"),
	},
	{
		Key:  NoteKeyRegulatory,
		Name: i18n.NewString("Regulatory"),
		Desc: i18n.NewString("Regulatory information"),
	},
	{
		Key:  NoteKeySafety,
		Name: i18n.NewString("Safety"),
		Desc: i18n.NewString("Safety instructions"),
	},
	{
		Key:  NoteKeyShipLine,
		Name: i18n.NewString("Ship Line"),
		Desc: i18n.NewString("Ship line"),
	},
	{
		Key:  NoteKeySupplier,
		Name: i18n.NewString("Supplier"),
		Desc: i18n.NewString("Supplier remarks"),
	},
	{
		Key:  NoteKeyTransport,
		Name: i18n.NewString("Transport"),
		Desc: i18n.NewString("Transportation information"),
	},
	{
		Key:  NoteKeyDelivery,
		Name: i18n.NewString("Delivery"),
		Desc: i18n.NewString("Delivery information"),
	},
	{
		Key:  NoteKeyQuarantine,
		Name: i18n.NewString("Quarantine"),
		Desc: i18n.NewString("Quarantine information"),
	},
	{
		Key:  NoteKeyTax,
		Name: i18n.NewString("Tax"),
		Desc: i18n.NewString("Tax declaration"),
	},
	{
		Key:  NoteKeyOther,
		Name: i18n.NewString("Other"),
		Desc: i18n.NewString("Mutually defined"),
	},
}

// Note represents a free text of additional information that may be
// added to a document.
type Note struct {
	uuid.Identify
	// Key specifying subject of the text
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Code used for additional data that may be required to identify the note.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`
	// Source of this note, especially useful when auto-generated.
	Src cbc.Key `json:"src,omitempty" jsonschema:"title=Source"`
	// The contents of the note
	Text string `json:"text" jsonschema:"title=Text"`
	// Additional information about the note
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
	// Extension data
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
}

// Normalize will perform basic normalization on the Note.
func (n *Note) Normalize() {
	if n == nil {
		return
	}
	uuid.Normalize(&n.UUID)
	n.Code = cbc.NormalizeCode(n.Code)
	n.Text = cbc.NormalizeString(n.Text)
	n.Ext = tax.CleanExtensions(n.Ext)
}

// Validate checks that the note looks okay.
func (n *Note) Validate() error {
	return validation.ValidateStruct(n,
		validation.Field(&n.Key, isValidNoteKey),
		validation.Field(&n.Code),
		validation.Field(&n.Text, validation.Required),
		validation.Field(&n.Src),
		validation.Field(&n.Meta),
		validation.Field(&n.Ext),
	)
}

// NoteFromScenario creates a new Note from a ScenarioNote.
func NoteFromScenario(sn *tax.ScenarioNote) *Note {
	if sn == nil {
		return nil
	}
	return &Note{
		Key:  sn.Key,
		Code: sn.Code,
		Src:  sn.Src,
		Text: sn.Text,
		Ext:  sn.Ext,
	}
}

var isValidNoteKey = validation.In(validNoteKeys()...)

func validNoteKeys() []interface{} {
	ks := make([]interface{}, len(NoteKeyDefinitions))
	for i, v := range NoteKeyDefinitions {
		ks[i] = v.Key
	}
	return ks
}

// WithSrc instantiates a new source instance with the provided
// source property set. This is a useful pattern for regional
// configurations.
func (n *Note) WithSrc(src cbc.Key) *Note {
	nw := *n // copy
	nw.Src = src
	return &nw
}

// WithCode provides a new copy of the note with the code set.
func (n *Note) WithCode(code cbc.Code) *Note {
	nw := *n // copy
	nw.Code = code
	return &nw
}

// SameAs returns true if the provided note is the same as
// the current one. Comparison is only made using the
// Key, Code, and Src properties.
//
// For a more complete comparison, use Equals.
func (n *Note) SameAs(n2 *Note) bool {
	if n == nil || n2 == nil {
		return false
	}
	return n.Key == n2.Key &&
		n.Code == n2.Code &&
		n.Src == n2.Src
}

// Equals returns true if the provided note is the same as the current one.
func (n *Note) Equals(n2 *Note) bool {
	return n.Key == n2.Key &&
		n.Code == n2.Code &&
		n.Src == n2.Src &&
		n.Text == n2.Text &&
		n.Meta.Equals(n2.Meta)
}

type validateNotes struct {
	key cbc.Key
}

// ValidateNotesHasKey returns a validation rule that check that at least one
// of the notes has the provided key.
func ValidateNotesHasKey(key cbc.Key) validation.Rule {
	return &validateNotes{key: key}
}

func (v *validateNotes) Validate(value any) error {
	notes, ok := value.([]*Note)
	if !ok {
		return nil
	}
	for _, n := range notes {
		if n.Key.In(v.key) {
			return nil // match found, this is good
		}
	}
	return fmt.Errorf("with key '%s' missing", v.key.String())
}

// JSONSchemaExtend adds the list of definitions for the notes.
func (Note) JSONSchemaExtend(schema *jsonschema.Schema) {
	ks, _ := schema.Properties.Get("key")
	ks.OneOf = make([]*jsonschema.Schema, len(NoteKeyDefinitions))
	for i, v := range NoteKeyDefinitions {
		ks.OneOf[i] = &jsonschema.Schema{
			Const:       v.Key.String(),
			Title:       v.Name.String(),
			Description: v.Desc.String(),
		}
	}
}
