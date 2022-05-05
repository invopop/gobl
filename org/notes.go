package org

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Notes holds an array of Note objects
type Notes []*Note

// NoteKey contains a code for the type of note.
type NoteKey string

// Predefined list of supported note keys based on the
// UNTDID 4451 list of text subject qualifiers. We've picked the ones
// which we think are most useful, but if you require an additional
// code, please send a pull request.
const (
	NoteKeyGoods          NoteKey = "goods"           // Goods Description
	NoteKeyPayment        NoteKey = "payment"         // Terms of Payment
	NoteKeyLegal          NoteKey = "legal"           // Legal or regulatory information
	NoteKeyDangerousGoods NoteKey = "dangerous-goods" // Dangerous goods additional information
	NoteKeyAck            NoteKey = "ack"             // Acknowledgement Description
	NoteKeyRate           NoteKey = "rate"            // Rate additional information
	NoteKeyReason         NoteKey = "reason"          // Reason
	NoteKeyDispute        NoteKey = "dispute"         // Dispute
	NoteKeyCustomer       NoteKey = "customer"        // Customer remarks
	NoteKeyGlossary       NoteKey = "glossary"        // Glossary
	NoteKeyCustoms        NoteKey = "customs"         // Customs declaration information
	NoteKeyGeneral        NoteKey = "general"         // General information
	NoteKeyHandling       NoteKey = "handling"        // Handling instructions
	NoteKeyPackaging      NoteKey = "packaging"       // Packaging information
	NoteKeyLoading        NoteKey = "loading"         // Loading instructions
	NoteKeyPrice          NoteKey = "price"           // Price conditions
	NoteKeyPriority       NoteKey = "priority"        // Priority information
	NoteKeyRegulatory     NoteKey = "regulatory"      // Regulatory information
	NoteKeySafety         NoteKey = "safety"          // Safety Instructions
	NoteKeyShipLine       NoteKey = "ship-line"       // Ship Line
	NoteKeySupplier       NoteKey = "supplier"        // Supplier remarks
	NoteKeyTransport      NoteKey = "transport"       // Transportation information
	NoteKeyDelivery       NoteKey = "delivery"        // Delivery Information
	NoteKeyQuarantine     NoteKey = "quarantine"      // Quarantine Information
	NoteKeyTax            NoteKey = "tax"             // Tax declaration
)

// UNTDID4451NoteKeyMap used to convert note codes into their official
// representation.
var UNTDID4451NoteKeyMap = map[NoteKey]string{
	NoteKeyGoods:          "AAA",
	NoteKeyDangerousGoods: "AAC",
	NoteKeyAck:            "AAE",
	NoteKeyRate:           "AAF",
	NoteKeyLegal:          "ABY", // Regulatory information
	NoteKeyReason:         "ACD",
	NoteKeyDispute:        "ACE",
	NoteKeyCustomer:       "CUR",
	NoteKeyCustoms:        "CUS",
	NoteKeyGlossary:       "ACZ",
	NoteKeyGeneral:        "AAI",
	NoteKeyHandling:       "HAN",
	NoteKeyPrice:          "AAK",
	NoteKeyLoading:        "LOI",
	NoteKeyPackaging:      "PKG",
	NoteKeyPayment:        "PMT",
	NoteKeyPriority:       "PRI",
	NoteKeyRegulatory:     "REG",
	NoteKeySafety:         "SAF",
	NoteKeyShipLine:       "SLR",
	NoteKeySupplier:       "SUR",
	NoteKeyTransport:      "TRA",
	NoteKeyDelivery:       "DEL",
	NoteKeyQuarantine:     "QIN",
	NoteKeyTax:            "TXD",
}

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
		validation.Field(&n.Key),
		validation.Field(&n.Text, validation.Required),
	)
}

// Validate checks to ensure the note code is part of the list of
// accepted values.
func (c NoteKey) Validate() error {
	_, ok := UNTDID4451NoteKeyMap[c]
	if !ok {
		return errors.New("invalid")
	}
	return nil
}

// WithSrc instantiates a new source instance with the provided
// source property set. This is a useful pattern for regional
// configurations.
func (n *Note) WithSrc(src string) *Note {
	nw := *n // copy
	nw.Src = src
	return &nw
}

// UNTDID4451 returns the official type code, or "NA" if none
// is set.
func (c NoteKey) UNTDID4451() string {
	s, ok := UNTDID4451NoteKeyMap[c]
	if !ok {
		return "NA"
	}
	return s
}
