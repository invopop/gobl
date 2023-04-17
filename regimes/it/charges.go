package it

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// List of charge types specific to the italian regime.
const (
	ChargeKeyStampDuty cbc.Key = "stamp-duty"
)

var chargeKeys = []*tax.KeyDefinition{
	{
		Key: ChargeKeyStampDuty,
		Name: i18n.String{
			i18n.EN: "Stamp Duty",
			i18n.IT: "Imposta di bollo",
		},
		Desc: i18n.String{
			i18n.EN: "A fixed-price tax applied to the production, request or presentation of certain documents: civil, commercial, judicial and extrajudicial documents, on notices, on posters.",
			i18n.IT: "Un'imposta applicata alla produzione, richiesta o presentazione di determinati documenti: atti civili, commerciali, giudiziali ed extragiudiziali, sugli avvisi, sui manifesti.",
		},
	},
}
