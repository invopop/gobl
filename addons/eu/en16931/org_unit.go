package en16931

import (
	"fmt"

	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func normalizeOrgItem(item *org.Item) {
	item.Unit, item.Ext = untdid.NormalizeUnit(item.Unit, item.Ext)
	// BR-23 requires a unit of measure on every line, so stand in the generic
	// unit when the document gives neither a key nor a code.
	if item.Unit == cbc.KeyEmpty && !item.Ext.Has(untdid.ExtKeyUnit) {
		item.Unit = org.UnitOne
	}
}

func normalizeOrgAttribute(a *org.Attribute) {
	a.Unit, a.Ext = untdid.NormalizeUnit(a.Unit, a.Ext)
}

func orgItemRules() *rules.Set {
	return rules.For(new(org.Item),
		rules.Assert("02", fmt.Sprintf("unit code must be determinable from the unit or the `%s` extension (BR-23)", untdid.ExtKeyUnit),
			is.Func("required valid UNTDID unit", func(value any) bool {
				item, ok := value.(*org.Item)
				if !ok || item == nil {
					return false
				}
				if item.Ext.Has(untdid.ExtKeyUnit) {
					return tax.ExtensionHasValidCode(untdid.ExtKeyUnit).Check(item.Ext)
				}
				return untdid.UnitCode(item.Unit) != cbc.CodeEmpty
			}),
		),
	)
}
