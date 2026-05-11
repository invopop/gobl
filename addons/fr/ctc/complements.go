package ctc

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
)

// Characteristic mirrors the CDAR SpecifiedDocumentCharacteristic
// element (MDT-207 and friends) used on Flow 6 lifecycle messages.
// It is attached to a bill.StatusLine via Complements and carries
// either:
//
//  1. A payment-related amount on a paid / partially-accepted /
//     completed line — e.g. TypeCode=MEN with Amount set for the
//     montant encaissé, MPA for amount paid, RAP for remaining.
//
//  2. A field-level correction on a rejected / disputed /
//     partially-accepted line, with ReasonCode pointing at the
//     sibling bill.Reason (via its fr-ctc-reason-code extension).
//
// The shape is intentionally close to CDAR so the converter can
// round-trip losslessly; most fields are optional.
type Characteristic struct {
	// ID optionally identifies the characteristic.
	ID string `json:"id,omitempty" jsonschema:"title=ID"`

	// TypeCode is the CDAR CharacteristicTypeCode.
	TypeCode cbc.Code `json:"type_code,omitempty" jsonschema:"title=Type Code"`

	// ReasonCode links this characteristic to a sibling bill.Reason
	// via its fr-ctc-reason-code extension value.
	ReasonCode cbc.Code `json:"reason_code,omitempty" jsonschema:"title=Reason Code"`

	// Description is a free-form human-readable explanation.
	Description string `json:"description,omitempty" jsonschema:"title=Description"`

	// Changed signals whether the reported value represents a
	// correction (true) or is being reported unchanged (false).
	Changed *bool `json:"changed,omitempty" jsonschema:"title=Changed"`

	// Direction carries the CDAR AdjustmentDirectionCode.
	Direction cbc.Code `json:"direction,omitempty" jsonschema:"title=Direction"`

	// Name is the semantic label of the field the characteristic
	// refers to.
	Name string `json:"name,omitempty" jsonschema:"title=Name"`

	// Location is a locator (XPath, JSON pointer, etc.) into the
	// referenced invoice identifying the specific field.
	Location string `json:"location,omitempty" jsonschema:"title=Location"`

	// Value carries a free-form string value when the field is textual.
	Value string `json:"value,omitempty" jsonschema:"title=Value"`

	// Code carries a coded value when the field is itself a code.
	Code cbc.Code `json:"code,omitempty" jsonschema:"title=Code"`

	// Percent holds a percentage value (e.g. a VAT rate correction).
	Percent *num.Percentage `json:"percent,omitempty" jsonschema:"title=Percent"`

	// Amount holds a monetary value paired with its currency.
	Amount *currency.Amount `json:"amount,omitempty" jsonschema:"title=Amount"`

	// Numeric holds a plain numeric value without currency.
	Numeric *num.Amount `json:"numeric,omitempty" jsonschema:"title=Numeric"`

	// Quantity holds a quantity value, optionally qualified by Measure.
	Quantity *num.Amount `json:"quantity,omitempty" jsonschema:"title=Quantity"`

	// Measure optionally describes the unit of Quantity or Numeric.
	Measure string `json:"measure,omitempty" jsonschema:"title=Measure"`

	// DateTime holds a date-time value.
	DateTime *cal.DateTime `json:"date_time,omitempty" jsonschema:"title=Date Time"`
}

// Characteristic.TypeCode values (MDT-207).
const (
	// Payment-related amounts
	TypeCodeAmountReceived    cbc.Code = "MEN"    // Montant encaissé (TTC)
	TypeCodeAmountPaid        cbc.Code = "MPA"    // Montant payé
	TypeCodeAmountRemaining   cbc.Code = "RAP"    // Reste à payer (paiement partiel)
	TypeCodeDiscount          cbc.Code = "ESC"    // Escompte accordé
	TypeCodeRebate            cbc.Code = "RAB"    // Rabais accordé
	TypeCodeReduction         cbc.Code = "REM"    // Remise accordée
	TypeCodeAmountApproved    cbc.Code = "MAP"    // Montant HT approuvé
	TypeCodeAmountApprovedTTC cbc.Code = "MAPTTC" // Montant TTC approuvé
	TypeCodeAmountRejected    cbc.Code = "MNA"    // Montant HT non approuvé
	TypeCodeAmountRejectedTTC cbc.Code = "MNATTC" // Montant TTC non approuvé

	// Rejection / correction markers
	TypeCodeBankDetailsUpdate cbc.Code = "CBB" // Coordonnées bancaires bénéficiaire à modifier
	TypeCodeInvalidData       cbc.Code = "DIV" // Donnée invalide
	TypeCodeExpectedData      cbc.Code = "DVA" // Donnée valide attendue
	TypeCodeOverrideData      cbc.Code = "MAJ" // Donnée à prendre en compte à la place de celle présente dans la facture
)

// typeCodes lists all accepted Characteristic.TypeCode values; used by
// validation to reject codes outside the controlled MDT-207 set.
var typeCodes = []cbc.Code{
	TypeCodeAmountReceived, TypeCodeAmountPaid, TypeCodeAmountRemaining,
	TypeCodeDiscount, TypeCodeRebate, TypeCodeReduction,
	TypeCodeAmountApproved, TypeCodeAmountApprovedTTC,
	TypeCodeAmountRejected, TypeCodeAmountRejectedTTC,
	TypeCodeBankDetailsUpdate, TypeCodeInvalidData,
	TypeCodeExpectedData, TypeCodeOverrideData,
}
