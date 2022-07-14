package org

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Notes holds an array of Note objects
type Notes []*Note

// Predefined list of supported note keys based on the
// UNTDID 4451 list of text subject qualifiers. We've picked the ones
// which we think are most useful, but if you require an additional
// code, please send a pull request.
const (
	NoteKeyGoods          Key = "goods"           // Goods Description
	NoteKeyPayment        Key = "payment"         // Terms of Payment
	NoteKeyLegal          Key = "legal"           // Legal or regulatory information
	NoteKeyDangerousGoods Key = "dangerous-goods" // Dangerous goods additional information
	NoteKeyAck            Key = "ack"             // Acknowledgement Description
	NoteKeyRate           Key = "rate"            // Rate additional information
	NoteKeyReason         Key = "reason"          // Reason
	NoteKeyDispute        Key = "dispute"         // Dispute
	NoteKeyCustomer       Key = "customer"        // Customer remarks
	NoteKeyGlossary       Key = "glossary"        // Glossary
	NoteKeyCustoms        Key = "customs"         // Customs declaration information
	NoteKeyGeneral        Key = "general"         // General information
	NoteKeyHandling       Key = "handling"        // Handling instructions
	NoteKeyPackaging      Key = "packaging"       // Packaging information
	NoteKeyLoading        Key = "loading"         // Loading instructions
	NoteKeyPrice          Key = "price"           // Price conditions
	NoteKeyPriority       Key = "priority"        // Priority information
	NoteKeyRegulatory     Key = "regulatory"      // Regulatory information
	NoteKeySafety         Key = "safety"          // Safety Instructions
	NoteKeyShipLine       Key = "ship-line"       // Ship Line
	NoteKeySupplier       Key = "supplier"        // Supplier remarks
	NoteKeyTransport      Key = "transport"       // Transportation information
	NoteKeyDelivery       Key = "delivery"        // Delivery Information
	NoteKeyQuarantine     Key = "quarantine"      // Quarantine Information
	NoteKeyTax            Key = "tax"             // Tax declaration
)

var untdid4451NoteKeyMap = map[Key]string{
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
	Key Key `json:"key,omitempty" jsonschema:"title=Key"`
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
		validation.Field(&n.Key, validation.In(validUNTDID4451Keys()...)),
		validation.Field(&n.Text, validation.Required),
	)
}

func validUNTDID4451Keys() []interface{} {
	ks := make([]interface{}, len(untdid4451NoteKeyMap))
	i := 0
	for v := range untdid4451NoteKeyMap {
		ks[i] = v
		i++
	}
	return ks
}

// UNTDID4451 provides the note's UNTDID 4451 equivalent
// value. If not available, returns "NA".
func (n *Note) UNTDID4451() string {
	s, ok := untdid4451NoteKeyMap[n.Key]
	if !ok {
		return "NA"
	}
	return s
}

// WithSrc instantiates a new source instance with the provided
// source property set. This is a useful pattern for regional
// configurations.
func (n *Note) WithSrc(src string) *Note {
	nw := *n // copy
	nw.Src = src
	return &nw
}
