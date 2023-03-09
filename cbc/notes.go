package cbc

import (
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// Predefined list of supported note keys based on the
// UNTDID 4451 list of text subject qualifiers. We've picked the ones
// which we think are most useful, but if you require an additional
// code, please send a pull request.
const (
	// Goods Description
	NoteKeyGoods Key = "goods"
	// Terms of Payment
	NoteKeyPayment Key = "payment"
	// Legal or regulatory information
	NoteKeyLegal Key = "legal"
	// Dangerous goods additional information
	NoteKeyDangerousGoods Key = "dangerous-goods"
	// Acknowledgement Description
	NoteKeyAck Key = "ack"
	// Rate additional information
	NoteKeyRate Key = "rate"
	// Reason
	NoteKeyReason Key = "reason"
	// Dispute
	NoteKeyDispute Key = "dispute"
	// Customer remarks
	NoteKeyCustomer Key = "customer"
	// Glossary
	NoteKeyGlossary Key = "glossary"
	// Customs declaration information
	NoteKeyCustoms Key = "customs"
	// General information
	NoteKeyGeneral Key = "general"
	// Handling instructions
	NoteKeyHandling Key = "handling"
	// Packaging information
	NoteKeyPackaging Key = "packaging"
	// Loading instructions
	NoteKeyLoading Key = "loading"
	// Price conditions
	NoteKeyPrice Key = "price"
	// Priority information
	NoteKeyPriority Key = "priority"
	// Regulatory information
	NoteKeyRegulatory Key = "regulatory"
	// Safety Instructions
	NoteKeySafety Key = "safety"
	// Ship Line
	NoteKeyShipLine Key = "ship-line"
	// Supplier remarks
	NoteKeySupplier Key = "supplier"
	// Transportation information
	NoteKeyTransport Key = "transport"
	// Delivery Information
	NoteKeyDelivery Key = "delivery"
	// Quarantine Information
	NoteKeyQuarantine Key = "quarantine"
	// Tax declaration
	NoteKeyTax Key = "tax"
)

// DefNoteKey holds a note key definition
type DefNoteKey struct {
	// Key to match against
	Key Key `json:"key" jsonschema:"title=Key"`
	// Description of the Note Key
	Description string `json:"description" jsonschema:"title=Description"`
	// UNTDID 4451 code
	UNTDID4451 Code `json:"untdid4451" jsonschema:"title=UNTDID4451 Code"`
}

// NoteKeyDefinitions provides a map of Note Keys to their definitions
// including a description and UNTDID code.
var NoteKeyDefinitions = []DefNoteKey{
	{
		Key:         NoteKeyGoods,
		Description: "Goods Description",
		UNTDID4451:  "AAA",
	},
	{
		Key:         NoteKeyPayment,
		Description: "Terms of Payment",
		UNTDID4451:  "PMT",
	},
	{
		Key:         NoteKeyLegal,
		Description: "Legal or regulatory information",
		UNTDID4451:  "ABY", // Regulatory information
	},
	{
		Key:         NoteKeyDangerousGoods,
		Description: "Dangerous goods additional information",
		UNTDID4451:  "AAC",
	},
	{
		Key:         NoteKeyAck,
		Description: "Acknowledgement Description",
		UNTDID4451:  "AAE",
	},
	{
		Key:         NoteKeyRate,
		Description: "Rate additional information",
		UNTDID4451:  "AAF",
	},
	{
		Key:         NoteKeyReason,
		Description: "Reason",
		UNTDID4451:  "ACD",
	},
	{
		Key:         NoteKeyDispute,
		Description: "Dispute",
		UNTDID4451:  "ACE",
	},
	{
		Key:         NoteKeyCustomer,
		Description: "Customer remarks",
		UNTDID4451:  "CUR",
	},
	{
		Key:         NoteKeyGlossary,
		Description: "Glossary",
		UNTDID4451:  "ACZ",
	},
	{
		Key:         NoteKeyCustoms,
		Description: "Customs declaration information",
		UNTDID4451:  "CUS",
	},
	{
		Key:         NoteKeyGeneral,
		Description: "General information",
		UNTDID4451:  "AAI",
	},
	{
		Key:         NoteKeyHandling,
		Description: "Handling instructions",
		UNTDID4451:  "HAN",
	},
	{
		Key:         NoteKeyPackaging,
		Description: "Packaging information",
		UNTDID4451:  "PKG",
	},
	{
		Key:         NoteKeyLoading,
		Description: "Loading instructions",
		UNTDID4451:  "LOI",
	},
	{
		Key:         NoteKeyPrice,
		Description: "Price conditions",
		UNTDID4451:  "AAK",
	},
	{
		Key:         NoteKeyPriority,
		Description: "Priority information",
		UNTDID4451:  "PRI",
	},
	{
		Key:         NoteKeyRegulatory,
		Description: "Regulatory information",
		UNTDID4451:  "REG",
	},
	{
		Key:         NoteKeySafety,
		Description: "Safety instructions",
		UNTDID4451:  "SAF",
	},
	{
		Key:         NoteKeyShipLine,
		Description: "Ship line",
		UNTDID4451:  "SLR",
	},
	{
		Key:         NoteKeySupplier,
		Description: "Supplier remarks",
		UNTDID4451:  "SUR",
	},
	{
		Key:         NoteKeyTransport,
		Description: "Transportation information",
		UNTDID4451:  "TRA",
	},
	{
		Key:         NoteKeyDelivery,
		Description: "Delivery information",
		UNTDID4451:  "DEL",
	},
	{
		Key:         NoteKeyQuarantine,
		Description: "Quarantine information",
		UNTDID4451:  "QIN",
	},
	{
		Key:         NoteKeyTax,
		Description: "Tax declaration",
		UNTDID4451:  "TXD",
	},
}

// Note represents a free text of additional information that may be
// added to a document.
type Note struct {
	// Key specifying subject of the text
	Key Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Code used for additional data that may be required to identify the note.
	Code string `json:"code,omitempty" jsonschema:"title=Code"`
	// Source of this note, especially useful when auto-generated.
	Src Key `json:"src,omitempty" jsonschema:"title=Source"`
	// The contents of the note
	Text string `json:"text" jsonschema:"title=Text"`
}

// Validate checks that the note looks okay.
func (n *Note) Validate() error {
	return validation.ValidateStruct(n,
		validation.Field(&n.Key, isValidNoteKey),
		validation.Field(&n.Text, validation.Required),
	)
}

var isValidNoteKey = validation.In(validNoteKeys()...)

func validNoteKeys() []interface{} {
	ks := make([]interface{}, len(NoteKeyDefinitions))
	for i, v := range NoteKeyDefinitions {
		ks[i] = v.Key
	}
	return ks
}

// UNTDID4451 provides the note's UNTDID 4451 equivalent
// value. If not available, returns CodeEmpty.
func (n *Note) UNTDID4451() Code {
	for _, v := range NoteKeyDefinitions {
		if v.Key == n.Key {
			return v.UNTDID4451
		}
	}
	return CodeEmpty
}

// WithSrc instantiates a new source instance with the provided
// source property set. This is a useful pattern for regional
// configurations.
func (n *Note) WithSrc(src Key) *Note {
	nw := *n // copy
	nw.Src = src
	return &nw
}

// JSONSchemaExtend adds the list of definitions for the notes.
func (Note) JSONSchemaExtend(schema *jsonschema.Schema) {
	schema.OneOf = make([]*jsonschema.Schema, len(NoteKeyDefinitions))
	for i, v := range NoteKeyDefinitions {
		schema.OneOf[i] = &jsonschema.Schema{
			Const:       v.Key.String(),
			Description: v.Description,
		}
	}
}
