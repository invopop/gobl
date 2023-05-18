package pay

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/internal/here"
	"github.com/invopop/jsonschema"
	"github.com/invopop/validation"
)

// MeansKey is a pre-defined payment means key, such as "card" or "cash"
// which can be used to describe how a payment was or will be made.
//
// The keys are based on a subset of the UNTDID 4461 code list and simplified
// to cover most use cases. Structures that use the MeansKey type should also
// include a `cbc.Code` field to allow for custom or more precise codes to be
// used according to the tax regime or pre-arranged agreement between parties.
type MeansKey cbc.Key

// Standard payment means codes for instructions.
// If you require more payment means options, please make a pull request
// and try to include references to the use case. All new means keys should
// map to an existing UNTDID 4461 code.
const (
	MeansKeyAny            MeansKey = "any" // Use any method available.
	MeansKeyCard           MeansKey = "card"
	MeansKeyCreditTransfer MeansKey = "credit-transfer"
	MeansKeyDebitTransfer  MeansKey = "debit-transfer"
	MeansKeyCash           MeansKey = "cash"
	MeansKeyPromissoryNote MeansKey = "promissory-note"
	MeansKeyNetting        MeansKey = "netting" // Clearing between parties
	MeansKeyCheque         MeansKey = "cheque"
	MeansKeyBankDraft      MeansKey = "bank-draft"
	MeansKeyDirectDebit    MeansKey = "direct-debit" // aka. Mandate
	MeansKeyOnline         MeansKey = "online"       // Website from which payment can be made
	MeansKeyOther          MeansKey = "other"        // See the code value, if available.
)

// MeansKeyDef is used to define each of the payment means keys
// that can be accepted by GOBL.
type MeansKeyDef struct {
	// Key being described
	Key MeansKey `json:"key" jsonschema:"title=Key"`
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

// IsValidMeansKey returns an error if the means key is not part of the
// pre-defined list.
var IsValidMeansKey = validation.In(validMeansKeys()...)

func validMeansKeys() []interface{} {
	list := make([]interface{}, len(MeansKeyDefinitions))
	for i, v := range MeansKeyDefinitions {
		list[i] = v.Key
	}
	return list
}

// JSONSchemaExtend adds the method key definitions to the schema.
func (MeansKey) JSONSchemaExtend(schema *jsonschema.Schema) {
	val, _ := schema.Properties.Get("key")
	prop, ok := val.(*jsonschema.Schema)
	if ok {
		prop.OneOf = make([]*jsonschema.Schema, len(MeansKeyDefinitions))
		for i, v := range MeansKeyDefinitions {
			prop.OneOf[i] = &jsonschema.Schema{
				Const:       v.Key,
				Title:       v.Title,
				Description: v.Description,
			}
		}
	}
}

// JSONSchema provides a representation of the struct for usage in Schema.
func (MeansKey) JSONSchema() *jsonschema.Schema {
	oneOf := make([]*jsonschema.Schema, len(MeansKeyDefinitions))
	for i, v := range MeansKeyDefinitions {
		oneOf[i] = &jsonschema.Schema{
			Const:       v.Key,
			Title:       v.Title,
			Description: v.Description,
		}
	}
	return &jsonschema.Schema{
		Type:      "string",
		Title:     "Means Key",
		MinLength: cbc.KeyMinLength,
		MaxLength: cbc.KeyMaxLength,
		OneOf:     oneOf,
		Description: here.Doc(`
			MeansKey is a pre-defined payment means key, such as "card" or "cash"
			which can be used to describe how a payment was or will be made.
			
			The keys are based on a subset of the [UNTDID 4461](https://unece.org/fileadmin/DAM/trade/untdid/d16b/tred/tred4461.htm) code list and simplified
			to cover most use cases. Structures that use the MeansKey type should also
			include a cbc.Code field to allow for custom or more precise codes to be
			used according to the tax regime or pre-arranged agreement between parties.`),
	}
}
