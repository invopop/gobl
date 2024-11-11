package it

import (
	"github.com/invopop/gobl/tax"
)

const (
	extKeySDIRetained    = "it-sdi-retained"
	extKeySDIExempt      = "it-sdi-exempt"
	extKeySDIRetainedTax = "it-sdi-retained-tax"
	extKeySDINature      = "it-sdi-nature"
)

func normalizeTaxCombo(tc *tax.Combo) {
	// Migrate tax combos so that even if no Addon is defined for the document
	// yet, it'll still be accepted by GOBL.
	if tc.Ext.Has(extKeySDIRetainedTax) {
		tc.Ext[extKeySDIRetained] = tc.Ext[extKeySDIRetainedTax]
		delete(tc.Ext, extKeySDIRetainedTax)
	}
	if tc.Ext.Has(extKeySDINature) {
		tc.Ext[extKeySDIExempt] = tc.Ext[extKeySDINature]
		delete(tc.Ext, extKeySDINature)
	}
}
