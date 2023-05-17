package pay

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/validation"
)

// Standard payment means codes for instructions. This is a heavily reduced list
// practical codes which can be linked to UNTDID 4461 counterparts.
// If you require more payment means options, please send your pull requests.
const (
	MeansKeyAny            cbc.Key = "any" // Use any method available.
	MeansKeyCard           cbc.Key = "card"
	MeansKeyCreditTransfer cbc.Key = "credit-transfer"
	MeansKeyDebitTransfer  cbc.Key = "debit-transfer"
	MeansKeyCash           cbc.Key = "cash"
	MeansKeyPromissoryNote cbc.Key = "promissory-note"
	MeansKeyNetting        cbc.Key = "netting"
	MeansKeyCheque         cbc.Key = "cheque"
	MeansKeyBankDraft      cbc.Key = "bank-draft"
	MeansKeyDirectDebit    cbc.Key = "direct-debit" // aka. Mandate
	MeansKeyOnline         cbc.Key = "online"       // Website from which payment can be made
	MeansKeyOther          cbc.Key = "other"        // See the means_code value, if available.
)

// MeansKeyDef is used to define each of the payment means keys
// that can be accepted by GOBL.
type MeansKeyDef struct {
	// Key being described
	Key cbc.Key `json:"key" jsonschema:"title=Key"`
	// Human value of the key
	Title string `json:"title" jsonschema:"title=Title"`
	// Details about the meaning of the key
	Description string `json:"description" jsonschema:"title=Description"`
	// UNTDID 4461 Equivalent Code
	UNTDID4461 cbc.Code `json:"untdid4461" jsonschema:"title=UNTDID 4461 Code"`
}

// MeansKeyDefinitions includes all the payment means keys that
// are accepted by GOBL.
var MeansKeyDefinitions = []MeansKeyDef{
	{MeansKeyAny, "Any", "Any method available, no preference.", "1"},                            // Instrument not defined
	{MeansKeyCard, "Card", "Payment card.", "48"},                                                // Bank card
	{MeansKeyCreditTransfer, "Credit Transfer", "Sender initiated bank or wire transfer.", "30"}, // credit transfer
	{MeansKeyDebitTransfer, "Debit Transfer", "Receiver initiated bank or wire transfer.", "31"}, // debit transfer
	{MeansKeyCash, "Cash", "Cash in hand.", "10"},                                                // in cash
	{MeansKeyCheque, "Cheque", "Cheque from bank.", "20"},                                        // cheque
	{MeansKeyBankDraft, "Draft", "Bankers Draft or Bank Cheque.", "21"},                          // Banker's draft,
	{MeansKeyDirectDebit, "Direct Debit", "Direct debit from the customers bank account.", "49"}, // direct debit
	{MeansKeyOnline, "Online", "Online or web payment.", "68"},                                   // online payment service
	{MeansKeyPromissoryNote, "Promissory Note", "Promissory note contract.", "60"},               // Promissory note
	{MeansKeyNetting, "Netting", "Intercompany clearing or clearing between partners.", "97"},    // Netting
	{MeansKeyOther, "Other", "Other or mutually defined means of payment.", "ZZZ"},               // Other
}

var isValidMeansKey = validation.In(validMeansKeys()...)

func validMeansKeys() []interface{} {
	list := make([]interface{}, len(MeansKeyDefinitions))
	for i, v := range MeansKeyDefinitions {
		list[i] = v.Key
	}
	return list
}
