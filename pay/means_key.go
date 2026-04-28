package pay

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/jsonschema"
)

// Standard payment means codes for instructions.
// If you require more payment means options, please make a pull request
// and try to include references to the use case.
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
	MeansKeySEPA           cbc.Key = "sepa"         // extension for SEPA payments
	MeansKeyOther          cbc.Key = "other"

	// Compositional sub-keys to qualify a payment means. Combine with
	// `With`, e.g. `MeansKeyCard.With(MeansKeyDebit)`.
	MeansKeyCredit cbc.Key = "credit"
	MeansKeyDebit  cbc.Key = "debit"
)

// MeansKeyDefinitions includes all the payment means keys that
// are accepted by GOBL.
var MeansKeyDefinitions = []*cbc.Definition{
	{
		Key:  MeansKeyAny,
		Name: i18n.NewString("Any"),
		Desc: i18n.NewString("Any method available, no preference."),
	},
	{
		Key:  MeansKeyCard,
		Name: i18n.NewString("Card"),
		Desc: i18n.NewString("Payment card. For backwards compatibility, a bare `card` key is treated as equivalent to `card+credit` by addons that distinguish credit and debit cards."),
	},
	{
		Key:  MeansKeyCard.With(MeansKeyCredit),
		Name: i18n.NewString("Credit Card"),
		Desc: i18n.NewString("Payment by credit card."),
	},
	{
		Key:  MeansKeyCard.With(MeansKeyDebit),
		Name: i18n.NewString("Debit Card"),
		Desc: i18n.NewString("Payment by debit card."),
	},
	{
		Key:  MeansKeyCreditTransfer,
		Name: i18n.NewString("Credit Transfer"),
		Desc: i18n.NewString("Sender initiated bank or wire transfer."),
	},
	{
		Key:  MeansKeyCreditTransfer.With(MeansKeySEPA),
		Name: i18n.NewString("SEPA Credit Transfer"),
		Desc: i18n.NewString("Sender initiated bank or wire transfer via SEPA."),
	},
	{
		Key:  MeansKeyDebitTransfer,
		Name: i18n.NewString("Debit Transfer"),
		Desc: i18n.NewString("Receiver initiated bank or wire transfer."),
	},
	{
		Key:  MeansKeyCash,
		Name: i18n.NewString("Cash"),
		Desc: i18n.NewString("Cash in hand."),
	},
	{
		Key:  MeansKeyCheque,
		Name: i18n.NewString("Cheque"),
		Desc: i18n.NewString("Cheque from bank."),
	},
	{
		Key:  MeansKeyBankDraft,
		Name: i18n.NewString("Draft"),
		Desc: i18n.NewString("Bankers Draft or Bank Cheque."),
	},
	{
		Key:  MeansKeyDirectDebit,
		Name: i18n.NewString("Direct Debit"),
		Desc: i18n.NewString("Direct debit from the customers bank account."),
	},
	{
		Key:  MeansKeyDirectDebit.With(MeansKeySEPA),
		Name: i18n.NewString("SEPA Direct Debit"),
		Desc: i18n.NewString("Direct debit from the customers bank account via SEPA."),
	},
	{
		Key:  MeansKeyOnline,
		Name: i18n.NewString("Online"),
		Desc: i18n.NewString("Online or web payment."),
	},
	{
		Key:  MeansKeyPromissoryNote,
		Name: i18n.NewString("Promissory Note"),
		Desc: i18n.NewString("Promissory note contract."),
	},
	{
		Key:  MeansKeyNetting,
		Name: i18n.NewString("Netting"),
		Desc: i18n.NewString("Intercompany clearing or clearing between partners."),
	},
	{
		Key:  MeansKeyOther,
		Name: i18n.NewString("Other"),
		Desc: i18n.NewString("Other or mutually defined means of payment."),
	},
}

// LookupMeansCode resolves the most specific code in m for the given
// payment means key. If key has no exact entry, the helper progressively
// pops sub-keys (`card+debit` → `card` → empty) until a match is found.
// This preserves backwards compatibility with addons that only registered
// the bare `card`, while letting addons that distinguish credit vs debit
// register explicit `card+credit` / `card+debit` entries that win.
func LookupMeansCode(m map[cbc.Key]cbc.Code, key cbc.Key) cbc.Code {
	for k := key; !k.IsEmpty(); k = k.Pop() {
		if code, ok := m[k]; ok {
			return code
		}
	}
	return ""
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
	prop, ok := schema.Properties.Get(property)
	if ok {
		anyOf := make([]*jsonschema.Schema, len(MeansKeyDefinitions))
		for i, v := range MeansKeyDefinitions {
			anyOf[i] = &jsonschema.Schema{
				Const:       v.Key,
				Title:       v.Name.String(),
				Description: v.Desc.String(),
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
