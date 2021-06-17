package pay

import "errors"

// MethodCode defines a standard name for a given payment method.
type MethodCode string

// Acceptable Methods of payment
const (
	MethodCodeCard     MethodCode = "CARD"
	MethodCodeTransfer MethodCode = "XFER"
	MethodCodeCash     MethodCode = "CASH"
	MethodCodeMandate  MethodCode = "MAN"    // aka. Direct Debit
	MethodCodeOnline   MethodCode = "ONLINE" // e.g. PayPal, GoogleCheckout, Stripe, etc.
)

// methodCodes defines the list of acceptable payment method types.
var methodCodes = []MethodCode{
	MethodCodeCard,
	MethodCodeTransfer,
	MethodCodeCash,
	MethodCodeMandate,
	MethodCodeOnline,
}

// Method describes how payment is expected to be made and under what conditions.
type Method struct {
	Code  MethodCode        `json:"code" jsonschema:"title=Code,description=Code for the method type that can be used."`
	Notes string            `json:"notes,omitempty" jsonschema:"title=Notes,description=Additional details related to this payment method."`
	Meta  map[string]string `json:"meta,omitempty" jsonschema:"title=Meta,description=Additional non-structure data."`
}

// Validate ensures the method code is valid according
// to the GoBL list of acceptable codes.
func (c MethodCode) Validate() error {
	if string(c) == "" {
		return nil
	}
	for _, v := range methodCodes {
		if v == c {
			return nil
		}
	}
	return errors.New("invalid payment method code")
}
