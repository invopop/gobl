package dfe

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
)

func normalizeParty(p *org.Party) {
	if p == nil || p.Ext == nil {
		return
	}

	// migrate old legacy extension keys
	for oldKey, newKey := range map[cbc.Key]cbc.Key{
		"br-nfse-fiscal-incentive": ExtKeyFiscalIncentive,
		"br-nfse-municipality":     ExtKeyMunicipality,
		"br-nfse-simples":          ExtKeySimples,
		"br-nfse-special-regime":   ExtKeySpecialRegime,
	} {
		if value, exists := p.Ext[oldKey]; exists {
			p.Ext[newKey] = value
			delete(p.Ext, oldKey)
		}
	}
}
