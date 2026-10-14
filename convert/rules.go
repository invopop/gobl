package convert

import (
	"fmt"
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// Rules provides the set of rules a converter checks documents against
// before exporting them into the format with the key, for use in Accepts.
// Fault codes are prefixed with the key, such as "UBL-PEPPOL-BIS-V3". The
// set is not added to the global rules registry, so it only applies when
// validated directly.
func Rules(key cbc.Key, sets ...*rules.Set) *rules.Set {
	ns := strings.ToUpper(strings.ReplaceAll(key.String(), "+", "-"))
	return rules.NewSet(rules.Code(ns), sets...)
}

// RequireAddons provides a rule for the "$addons" field of a document that
// faults when it does not include every one of the addons.
func RequireAddons(id rules.Code, keys ...cbc.Key) rules.Def {
	list := make([]string, len(keys))
	for i, k := range keys {
		list[i] = k.String()
	}
	desc := "requires addon " + strings.Join(list, ", ")
	if len(keys) > 1 {
		desc = "requires addons " + strings.Join(list, ", ")
	}
	return rules.Field("$addons",
		rules.Assert(id, desc,
			is.Func(fmt.Sprintf("includes %s", strings.Join(list, ", ")), func(value any) bool {
				addons, _ := value.(tax.AddonList)
				for _, k := range keys {
					if !k.In(addons...) {
						return false
					}
				}
				return true
			}),
		),
	)
}
