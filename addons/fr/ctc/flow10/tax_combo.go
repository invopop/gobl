package flow10

import (
	"github.com/invopop/gobl/addons/eu/en16931"
	"github.com/invopop/gobl/tax"
)

func normalizeTaxCombo(tc *tax.Combo) {
	en16931.NormalizeTaxCombo(tc)
}
