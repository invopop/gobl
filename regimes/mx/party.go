package mx

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
)

func normalizeParty(p *org.Party) error {
	// 2023-08-25: Migrate identities to extensions
	// Pending removal after migrations completed.
	idents := make([]*org.Identity, 0)
	for _, v := range p.Identities {
		if v.Key != "" {
			if p.Ext == nil {
				p.Ext = make(cbc.CodeMap)
			}
			p.Ext[v.Key] = v.Code
		} else {
			idents = append(idents, v)
		}
	}
	p.Identities = idents
	return nil
}
