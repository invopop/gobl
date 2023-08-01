package pay

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/jsonschema"
)

// Standard payment means codes for instructions.
// If you require more payment means options, please make a pull request
// and try to include references to the use case. All new means keys should
// map to an existing UNTDID 4461 code.
const (
	MeansKeyAny            cbc.Key = "any" // Use any method available.
	MeansKeyCard           cbc.Key = "card"
	MeansKeyCreditTransfer cbc.Key = "credit-transfer"
	MeansKeyDebitTransfer  cbc.Key = "debit-transfer"
	MeansKeyCash           cbc.Key = "cash"
	MeansKeyPromissoryNote cbc.Key = "promissory-note"
	MeansKeyNetting        cbc.Key = "netting" // Clearing between parties
	MeansKeyCheque         cbc.Key = "cheque"
	MeansKeyBankDraft      cbc.Key = "bank-draft"
	MeansKeyDirectDebit    cbc.Key = "direct-debit" // aka. Mandate
	MeansKeyOnline         cbc.Key = "online"       // Website from which payment can be made
	MeansKeyOther          cbc.Key = "other"
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

// HasValidMeansKey provides a usable validator for the means key
// to ensure it is at least *based* on one of the primary keys.
// This allows means keys to be extended or customised.
var HasValidMeansKey = cbc.HasValidKeyIn(validBaseMeansKeys()...)

func validBaseMeansKeys() []cbc.Key {
	list := make([]cbc.Key, len(MeansKeyDefinitions))
	for i, v := range MeansKeyDefinitions {
		list[i] = v.Key
	}
	return list
}

func extendJSONSchemaWithMeansKey(schema *jsonschema.Schema, property string) {
	val, _ := schema.Properties.Get(property)
	prop, ok := val.(*jsonschema.Schema)
	if ok {
		anyOf := make([]*jsonschema.Schema, len(MeansKeyDefinitions))
		for i, v := range MeansKeyDefinitions {
			anyOf[i] = &jsonschema.Schema{
				Const:       v.Key,
				Title:       v.Title,
				Description: v.Description,
			}
		}
		anyOf = append(anyOf, &jsonschema.Schema{
			Title:   "Regime Specific Key",
			Pattern: cbc.KeyPattern, // Allow custom keys
		})
		prop.Pattern = ""
		prop.AnyOf = anyOf
	}
}
