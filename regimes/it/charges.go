package it

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

const (
	ChargeTypeBollo cbc.Key = "bollo"
)

var chargeTypes = []*tax.KeyDefinition{
	{
		Key: ChargeTypeBollo,
		Name: i18n.String{
			i18n.EN: "Duty Stamp",
			i18n.IT: "Bollo",
		},
		Desc: i18n.String{
			i18n.EN: "A fixed-price tax applied to the production, request or presentation of certain documents: civil, commercial, judicial and extrajudicial documents, on notices, on posters.",
			i18n.IT: "Un'imposta applicata alla produzione, richiesta o presentazione di determinati documenti: atti civili, commerciali, giudiziali ed extragiudiziali, sugli avvisi, sui manifesti.",
		},
	},
}
