package cfdi

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func normalizeParty(p *org.Party) {
	if p == nil {
		return
	}
	// 2023-08-25: Migrate identities to extensions
	// Pending removal after migrations completed.
	idents := make([]*org.Identity, 0)
	for _, v := range p.Identities {
		if v.Key.In(migratedExtensionKeys...) {
			if p.Ext == nil {
				p.Ext = make(tax.Extensions)
			}
			p.Ext[v.Key] = tax.ExtValue(v.Code)
		} else {
			idents = append(idents, v)
		}
	}
	p.Identities = idents
}
