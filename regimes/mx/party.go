package mx

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func normalizeParty(p *org.Party) {
	if p == nil {
		return
	}
	// 2024-03-14: Migrate Tax ID Zone to extensions "mx-cfdi-post-code"
	if p.TaxID != nil && p.TaxID.Zone != "" { //nolint:staticcheck
		if p.Ext == nil {
			p.Ext = make(tax.Extensions)
		}
		p.Ext[extKeyPostCode] = tax.ExtValue(p.TaxID.Zone) //nolint:staticcheck
		p.TaxID.Zone = ""                                  //nolint:staticcheck
	}
}
