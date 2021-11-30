package pay

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/invopop/gobl/org"
)

// MethodCode defines a standard set of names for payment means.
type MethodCode string

// Standard payment method codes. This is a heavily reduced list of practical
// codes which can be linked to UNTDID 4461 counterparts.
// If you require more payment method options, please send your pull requests.
const (
	MethodCodeAny            MethodCode = "any" // Use any method available.
	MethodCodeCard           MethodCode = "card"
	MethodCodeCreditTransfer MethodCode = "credit_transfer"
	MethodCodeCash           MethodCode = "cash"
	MethodCodeDirectDebit    MethodCode = "direct_debit" // aka. Mandate
	MethodCodeOnline         MethodCode = "online"       // Website from which payment can be made
)

// https://unece.org/fileadmin/DAM/trade/untdid/d16b/tred/tred4461.htm
var untdid4461codes = map[MethodCode]string{
	MethodCodeAny:            "1",  // Instrument not defined
	MethodCodeCard:           "48", // bank card
	MethodCodeCreditTransfer: "30", // credit transfer
	MethodCodeCash:           "10", // in cash
	MethodCodeDirectDebit:    "31", // debit transfer
	MethodCodeOnline:         "68", // online payment service
}

// UNTDID4461 provides the standard UNTDID 4461 code for the payment method.
func (c MethodCode) UNTDID4461() string {
	return untdid4461codes[c]
}

// Instructions holds a set of instructions that determine how the payment has
// or should be made. A single "code" exists in which the preferred payment method
// should be provided. All other details serve as a reference.
type Instructions struct {
	Code           MethodCode        `json:"code" jsonschema:"title=Code,description=How payment is expected or has been arranged to be collected."`
	Detail         string            `json:"detail,omitempty" jsonschema:"title=Detail,description=Optional text description of the payment method."`
	Ref            string            `json:"ref,omitempty" jsonschema:"title=Ref,description=Remittance information, a text value used to link the payment with the invoice."`
	CreditTransfer []*CreditTransfer `json:"credit_transfer,omitempty" jsonschema:"title=Credit Transfer,description=Instructions for sending payment via a bank transfer."`
	Card           *Card             `json:"card,omitempty" jsonschema:"title=Card,description=Details of the payment that will be made via a credit or debit card."`
	DirectDebit    *DirectDebit      `json:"direct_debit,omitempty" jsonschema:"title=Direct Debit,description=A group of terms that can be used by the customer or payer to consolidate direct debit payments."`
	Online         []*Online         `json:"online,omitempty" jsonschema:"title=Online,description=Array of online payment options."`
	Notes          string            `json:"notes,omitempty" jsonschema:"title=Notes,description=Any additional instructions that may be required to make the payment."`
	Meta           map[string]string `json:"meta,omitempty" jsonschema:"title=Meta,description=Non-structured additional data that may be useful."`
}

// Card contains simplified card holder data as a reference for the customer.
type Card struct {
	Last4  string `json:"last4" jsonschema:"title=Last 4,description=Last 4 digits of the card's Primary Account Number (PAN)."`
	Holder string `json:"holder" jsonschema:"title=Holder Name,description=Name of the person whom the card belongs to."`
}

// DirectDebit defines the data that will be used to make the direct debit.
type DirectDebit struct {
	Ref      string `json:"ref,omitempty" jsonschema:"title=Mandate Reference,description=Unique identifier assigned by the payee for referencing the direct debit."`
	Creditor string `json:"creditor,omitempty" jsonschema:"title=Creditor ID,description=Unique banking reference that identifies the payee or seller assigned by the bank."`
	Account  string `json:"account,omitempty" jsonschema:"title=Account,description=Account identifier to be debited by the direct debit."`
}

// CreditTransfer contains fields that can be used for making payments via
// a bank transfer or wire.
type CreditTransfer struct {
	IBAN   string       `json:"iban,omitempty" jsonschema:"title=IBAN,description=International Bank Account Number"`
	BIC    string       `json:"bic,omitempty" jsonschema:"title=BIC,description=Bank Identifier Code used for international transfers."`
	Number string       `json:"number,omitempty" jsonschema:"title=Number,description=Account number, if IBAN not available."`
	Name   string       `json:"name,omitempty" jsonschema:"title=Name,description=Name of the bank."`
	Branch *org.Address `json:"branch,omitempty" jsonschema:"title=Branch,description=Bank office branch address, not normally required."`
}

// Online provides the details required to make a payment online using a website
type Online struct {
	Name    string `json:"name,omitempty" jsonschema:"title=Name,description=Descriptive name given to the online provider."`
	Address string `json:"addr" jsonschema:"title=Address,description=Full URL to be used for payment."`
}

// Validate ensures the Online method details look correct.
func (u *Online) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(u.Address, validation.Required, is.URL),
	)
}
