package pay

// BankTransfer contains fields that can be used for making payments via
// a bank transfer or wire.
type BankTransfer struct {
	IBAN string `json:"iban,omitempty" jsonschema:"title=IBAN,description=International Bank Account Number"`
	BIC  string `json:"bic,omitempty" jsonschema:"title=BIC,description=Bank Identifier Code used for international transfers."`
}
