package org

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Notes holds an array of Note objects
type Notes []*Note

// NoteCode contains a code for the type of note.
type NoteCode string

// Predefined list of supported note code types based on the
// UNTDID 4451 list of text subject qualifiers. We've picked the ones
// which we think are most useful, but if you require an additional
// code, please send a pull request.
const (
	GoodsNoteCode          NoteCode = "goods"           // Goods Description
	PaymentNoteCode        NoteCode = "payment"         // Terms of Payment
	LegalNoteCode          NoteCode = "legal"           // Legal or regulatory information
	DangerousGoodsNoteCode NoteCode = "dangerous-goods" // Dangerous goods additional information
	AckNoteCode            NoteCode = "ack"             // Acknowledgement Description
	RateNoteCode           NoteCode = "rate"            // Rate additional information
	ReasonNoteCode         NoteCode = "reason"          // Reason
	DisputeNoteCode        NoteCode = "dispute"         // Dispute
	CustomerNoteCode       NoteCode = "customer"        // Customer remarks
	GlossaryNoteCode       NoteCode = "glossary"        // Glossary
	CustomsNoteCode        NoteCode = "customs"         // Customs declaration information
	GeneralNoteCode        NoteCode = "general"         // General information
	HandlingNoteCode       NoteCode = "handling"        // Handling instructions
	PackagingNoteCode      NoteCode = "packaging"       // Packaging information
	LoadingNoteCode        NoteCode = "loading"         // Loading instructions
	PriceNoteCode          NoteCode = "price"           // Price conditions
	PriorityNoteCode       NoteCode = "priority"        // Priority information
	RegulatoryNoteCode     NoteCode = "regulatory"      // Regulatory information
	SafetyNoteCode         NoteCode = "safety"          // Safety Instructions
	ShipLineNoteCode       NoteCode = "ship-line"       // Ship Line
	SupplierNoteCode       NoteCode = "supplier"        // Supplier remarks
	TransportNoteCode      NoteCode = "transport"       // Transportation information
)

// UNTDID4451NoteCodeMap used to convert note codes into their official
// representation.
var UNTDID4451NoteCodeMap = map[NoteCode]string{
	GoodsNoteCode:          "AAA",
	DangerousGoodsNoteCode: "AAC",
	AckNoteCode:            "AAE",
	RateNoteCode:           "AAF",
	LegalNoteCode:          "ABY", // Regulatory information
	ReasonNoteCode:         "ACD",
	DisputeNoteCode:        "ACE",
	CustomerNoteCode:       "CUR",
	CustomsNoteCode:        "CUS",
	GlossaryNoteCode:       "ACZ",
	GeneralNoteCode:        "AAI",
	HandlingNoteCode:       "HAN",
	PriceNoteCode:          "AAK",
	LoadingNoteCode:        "LOI",
	PackagingNoteCode:      "PKG",
	PaymentNoteCode:        "PMT",
	PriorityNoteCode:       "PRI",
	RegulatoryNoteCode:     "REG",
	SafetyNoteCode:         "SAF",
	ShipLineNoteCode:       "SLR",
	SupplierNoteCode:       "SUR",
	TransportNoteCode:      "TRA",
}

// Note represents a free text of additional information that may be
// added to a document.
type Note struct {
	// Code specifying subject of the text
	Code NoteCode `json:"code,omitempty" jsonschema:"title=Code"`
	// The contents of the note
	Text string `json:"text" jsonschema:"title=Text"`
}

// Validate checks that the note looks okay.
func (n *Note) Validate() error {
	return validation.ValidateStruct(n,
		validation.Field(&n.Code),
		validation.Field(&n.Text, validation.Required),
	)
}

// Validate checks to ensure the note code is part of the list of
// accepted values.
func (c NoteCode) Validate() error {
	_, ok := UNTDID4451NoteCodeMap[c]
	if !ok {
		return errors.New("invalid")
	}
	return nil
}

// UNTDID4451 returns the official type code, or "NA" if none
// is set.
func (c NoteCode) UNTDID4451() string {
	s, ok := UNTDID4451NoteCodeMap[c]
	if !ok {
		return "NA"
	}
	return s
}
