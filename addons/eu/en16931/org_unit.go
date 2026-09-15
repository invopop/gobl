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

// UnitToUNTDID converts a GOBL unit key into its corresponding UNTDID unit
// code. It returns an empty code when the unit has no standard mapping. The
// table itself lives in the UNTDID catalogue, shared with every other format
// built on the same codes.
func UnitToUNTDID(unit cbc.Key) cbc.Code {
	return untdid.UnitCode(unit)
}

// UnitFromUNTDID converts a UNTDID unit code into its corresponding GOBL unit
// key. It returns an empty unit when GOBL has no standard mapping.
func UnitFromUNTDID(code cbc.Code) cbc.Key {
	return untdid.UnitKey(code)
}

func normalizeOrgItem(item *org.Item) {
	code := item.Ext.Get(untdid.ExtKeyUnit)
	if unit := UnitFromUNTDID(code); unit != cbc.KeyEmpty {
		item.Unit = unit
	}
	if item.Unit == cbc.KeyEmpty {
		item.Unit = org.UnitOne
	}
	if code == cbc.CodeEmpty {
		if code = UnitToUNTDID(item.Unit); code != cbc.CodeEmpty {
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
