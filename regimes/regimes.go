package regimes

import (
	"github.com/invopop/gobl/regimes/co"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/regimes/fr"
	"github.com/invopop/gobl/regimes/gb"
	"github.com/invopop/gobl/regimes/nl"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterRegime(co.Regime())
	tax.RegisterRegime(es.Regime())
	tax.RegisterRegime(fr.Regime())
	tax.RegisterRegime(gb.Regime())
	tax.RegisterRegime(nl.Regime())
}
