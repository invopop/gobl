package cbc

import (
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// NoteKey is used to describe the key used for identifying
// the type of note.
type NoteKey Key

// Predefined list of supported note keys based on the
// UNTDID 4451 list of text subject qualifiers. We've picked the ones
// which we think are most useful, but if you require an additional
// code, please send a pull request.
const (
	// Goods Description
	NoteKeyGoods NoteKey = "goods"
	// Terms of Payment
	NoteKeyPayment NoteKey = "payment"
	// Legal or regulatory information
	NoteKeyLegal NoteKey = "legal"
	// Dangerous goods additional information
	NoteKeyDangerousGoods NoteKey = "dangerous-goods"
	// Acknowledgement Description
	NoteKeyAck NoteKey = "ack"
	// Rate additional information
	NoteKeyRate NoteKey = "rate"
	// Reason
	NoteKeyReason NoteKey = "reason"
	// Dispute
	NoteKeyDispute NoteKey = "dispute"
	// Customer remarks
	NoteKeyCustomer NoteKey = "customer"
	// Glossary
	NoteKeyGlossary NoteKey = "glossary"
	// Customs declaration information
	NoteKeyCustoms NoteKey = "customs"
	// General information
	NoteKeyGeneral NoteKey = "general"
	// Handling instructions
	NoteKeyHandling NoteKey = "handling"
	// Packaging information
	NoteKeyPackaging NoteKey = "packaging"
	// Loading instructions
	NoteKeyLoading NoteKey = "loading"
	// Price conditions
	NoteKeyPrice NoteKey = "price"
	// Priority information
	NoteKeyPriority NoteKey = "priority"
	// Regulatory information
	NoteKeyRegulatory NoteKey = "regulatory"
	// Safety Instructions
	NoteKeySafety NoteKey = "safety"
	// Ship Line
	NoteKeyShipLine NoteKey = "ship-line"
	// Supplier remarks
	NoteKeySupplier NoteKey = "supplier"
	// Transportation information
	NoteKeyTransport NoteKey = "transport"
	// Delivery Information
	NoteKeyDelivery NoteKey = "delivery"
	// Quarantine Information
	NoteKeyQuarantine NoteKey = "quarantine"
	// Tax declaration
	NoteKeyTax NoteKey = "tax"
)

// DefNoteKey holds a note key definition
type DefNoteKey struct {
	// Key to match against
	Key NoteKey `json:"key" jsonschema:"title=Key"`
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

// Notes holds an array of Note objects
// type Notes []*Note

// Note represents a free text of additional information that may be
// added to a document.
type Note struct {
	// Key specifying subject of the text
	Key NoteKey `json:"key,omitempty" jsonschema:"title=Key"`
	// Code used for additional data that may be required to identify the note.
	Code string `json:"code,omitempty" jsonschema:"title=Code"`
	// Source of this note, especially useful when auto-generated.
	Src string `json:"src,omitempty" jsonschema:"title=Source"`
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
func (n *Note) WithSrc(src string) *Note {
	nw := *n // copy
	nw.Src = src
	return &nw
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (k NoteKey) JSONSchema() *jsonschema.Schema {
	s := &jsonschema.Schema{
		Title:       "Note Key",
		Type:        "string", // they're all strings
		OneOf:       make([]*jsonschema.Schema, len(NoteKeyDefinitions)),
		Description: "NoteKey identifies the type of note being edited",
	}
	for i, v := range NoteKeyDefinitions {
		s.OneOf[i] = &jsonschema.Schema{
			Const:       Key(v.Key).String(),
			Description: v.Description,
		}
	}
	return s
}
