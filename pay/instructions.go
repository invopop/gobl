package pay

import (
	"context"
	"encoding/json"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
)

// Instructions determine how the payment has or should be made. A
// single "key" exists in which the preferred payment method
// should be provided, all other details serve as a reference.
type Instructions struct {
	// The payment means expected or that have been arranged to be used to make the payment.
	Key cbc.Key `json:"key" jsonschema:"title=Key"`
	// Optional text description of the payment method
	Detail string `json:"detail,omitempty" jsonschema:"title=Detail"`
	// Remittance information or concept, a code value used to link the payment with the invoice.
	Ref cbc.Code `json:"ref,omitempty" jsonschema:"title=Reference"`
	// Instructions for sending payment via a bank transfer.
	CreditTransfer []*CreditTransfer `json:"credit_transfer,omitempty" jsonschema:"title=Credit Transfer"`
	// Details of the payment that will be made via a credit or debit card.
	Card *Card `json:"card,omitempty" jsonschema:"title=Card"`
	// A group of terms that can be used by the customer or payer to consolidate direct debit payments.
	DirectDebit *DirectDebit `json:"direct_debit,omitempty" jsonschema:"title=Direct Debit"`
	// Array of online payment options
	Online []*Online `json:"online,omitempty" jsonschema:"title=Online"`
	// Any additional instructions that may be required to make the payment.
	Notes string `json:"notes,omitempty" jsonschema:"title=Notes"`
	// Extension key-pairs values defined by a tax regime.
	Ext tax.Extensions `json:"ext,omitempty" jsonschema:"title=Extensions"`
	// Non-structured additional data that may be useful.
	Meta cbc.Meta `json:"meta,omitempty" jsonschema:"title=Meta"`
}

// Card contains simplified card holder data as a reference for the customer.
// PCI compliance requires only the first 6 and last 4 digits of the card number
// to be stored openly.
type Card struct {
	// First 6 digits of the card's Primary Account Number (PAN).
	First6 string `json:"first6" jsonschema:"title=First 6"`
	// Last 4 digits of the card's Primary Account Number (PAN).
	Last4 string `json:"last4" jsonschema:"title=Last 4"`
	// Name of the person whom the card belongs to.
	Holder string `json:"holder" jsonschema:"title=Holder Name"`
}

// DirectDebit defines the data that will be used to make the direct debit.
type DirectDebit struct {
	// Unique identifier assigned by the payee for referencing the direct debit.
	Ref string `json:"ref,omitempty" jsonschema:"title=Mandate Reference"`
	// Unique banking reference that identifies the payee or seller assigned by the bank.
	Creditor string `json:"creditor,omitempty" jsonschema:"title=Creditor ID"`
	// Account identifier to be debited by the direct debit.
	Account string `json:"account,omitempty" jsonschema:"title=Account"`
}

// CreditTransfer contains fields that can be used for making payments via
// a bank transfer or wire.
type CreditTransfer struct {
	// International Bank Account Number
	IBAN string `json:"iban,omitempty" jsonschema:"title=IBAN"`
	// Bank Identifier Code used for international transfers.
	BIC string `json:"bic,omitempty" jsonschema:"title=BIC"`
	// Account number, if IBAN not available.
	Number string `json:"number,omitempty" jsonschema:"title=Number"`
	// Name of the bank.
	Name string `json:"name,omitempty" jsonschema:"title=Name"`
	// Bank office branch address, not normally required.
	Branch *org.Address `json:"branch,omitempty" jsonschema:"title=Branch"`
}

// Online provides the details required to make a payment online using a website
type Online struct {
	// Key identifier for this online payment method.
	Key cbc.Key `json:"key,omitempty" jsonschema:"title=Key"`
	// Descriptive label for the online provider.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// URL to be used for payment.
	URL string `json:"url" jsonschema:"title=URL"`
}

// Normalize will try to normalize the instructions.
func (i *Instructions) Normalize(normalizers tax.Normalizers) {
	if i == nil {
		return
	}
	i.Ref = cbc.NormalizeCode(i.Ref)
	i.Ext = tax.CleanExtensions(i.Ext)
	normalizers.Each(i)
}

// Validate ensures the Online method details look correct.
func (u *Online) Validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Key),
		validation.Field(&u.Label),
		validation.Field(&u.URL, validation.Required, is.URL),
	)
}

// UnmarshalJSON is used to handle the migration of the Online's
// properties to Label and URL.
func (u *Online) UnmarshalJSON(data []byte) error {
	type Alias Online
	aux := struct {
		*Alias
		Name    string `json:"name"`
		Address string `json:"addr"`
	}{
		Alias: (*Alias)(u),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.Name != "" {
		u.Label = aux.Name
	}
	if aux.Address != "" {
		u.URL = aux.Address
	}
	return nil
}

// Validate ensures the fields provided in the instructions are valid.
func (i *Instructions) Validate() error {
	return i.ValidateWithContext(context.Background())
}

// ValidateWithContext ensures the fields provided in the instructions are valid.
func (i *Instructions) ValidateWithContext(ctx context.Context) error {
	return tax.ValidateStructWithContext(ctx, i,
		validation.Field(&i.Key, validation.Required, HasValidMeansKey),
		validation.Field(&i.Ref),
		validation.Field(&i.CreditTransfer),
		validation.Field(&i.DirectDebit),
		validation.Field(&i.Online),
		validation.Field(&i.Ext),
	)
}

// JSONSchemaExtend extends the JSONSchema for the Instructions type.
func (Instructions) JSONSchemaExtend(schema *jsonschema.Schema) {
	extendJSONSchemaWithMeansKey(schema, "key")
}
