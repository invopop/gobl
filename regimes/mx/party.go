package mx

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func normalizeParty(p *org.Party) error {
	// 2023-08-25: Migrate identities to extensions
	// Pending removal after migrations completed.
	idents := make([]*org.Identity, 0)
	for _, v := range p.Identities {
		if v.Key.In(migratedExtensionKeys...) {
			if p.Ext == nil {
				p.Ext = make(tax.ExtMap)
			}
			p.Ext[v.Key] = cbc.KeyOrCode(v.Code)
		} else {
			idents = append(idents, v)
		}
	}
	p.Identities = idents
	return nil
}
