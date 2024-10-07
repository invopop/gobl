package co

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

const (
	// copied from DIAN addon
	extKeyDIANMunicipality = "co-dian-municipality"
)

func normalizeParty(p *org.Party) {
	if p == nil {
		return
	}
	// 2024-03-14: Migrate Tax ID Zone to extensions "co-dian-municipality"
	if p.TaxID != nil && p.TaxID.Zone != "" { //nolint:staticcheck
		if p.Ext == nil {
			p.Ext = make(tax.Extensions)
		}
		p.Ext[extKeyDIANMunicipality] = tax.ExtValue(p.TaxID.Zone) //nolint:staticcheck
		p.TaxID.Zone = ""                                          //nolint:staticcheck
	}
}
