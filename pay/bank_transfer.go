package pay

import (
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
)

// BankTransfer contains fields that can be used for making payments via
// a bank transfer or wire.
type BankTransfer struct {
	IBAN   string            `json:"iban,omitempty" jsonschema:"title=IBAN,description=International Bank Account Number"`
	BIC    string            `json:"bic,omitempty" jsonschema:"title=BIC,description=Bank Identifier Code used for international transfers."`
	Number string            `json:"number,omitempty" jsonschema:"title=Number,description=Account number, if IBAN not available."`
	Name   string            `json:"name,omitempty" jsonschema:"title=Name,description=Name of the bank."`
	Branch *org.Address      `json:"branch,omitempty" jsonschema:"title=Branch,description=Bank office branch address, very rarely needed."`
	Notes  *i18n.String      `json:"notes,omitempty" jsonschema:"title=Notes,description=Any additional instructions that may be required to make the transfer."`
	Meta   map[string]string `json:"meta,omitempty" jsonschema:"title=Meta,description=Non-structured additional data that may be useful."`
}
