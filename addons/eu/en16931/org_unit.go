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
	code := item.Ext.Get(untdid.ExtKeyUnit)
	if unit := untdid.UnitKey(code); unit != cbc.KeyEmpty {
		item.Unit = unit
	}
	if item.Unit == cbc.KeyEmpty {
		item.Unit = org.UnitOne
	}
	if code == cbc.CodeEmpty {
		if code = untdid.UnitCode(item.Unit); code != cbc.CodeEmpty {
			item.Ext = item.Ext.Set(untdid.ExtKeyUnit, code)
		}
	}
}

func orgItemRules() *rules.Set {
	return rules.For(new(org.Item),
		rules.Assert("02", fmt.Sprintf("UNTDID unit `%s` must be present and valid (BR-23)", untdid.ExtKeyUnit),
			is.Func("required valid UNTDID unit", func(value any) bool {
				item, ok := value.(*org.Item)
				return ok && item != nil &&
					item.Ext.Has(untdid.ExtKeyUnit) &&
					tax.ExtensionHasValidCode(untdid.ExtKeyUnit).Check(item.Ext)
			}),
		),
	)
}
