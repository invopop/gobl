package bill

import (
	"fmt"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// RequireLineTaxCategory provides a rules.Test that ensures the line has the
// required tax category.
func RequireLineTaxCategory(key cbc.Code) rules.Test {
	return is.Func(
		fmt.Sprintf("line has tax category %s", key),
		func(val any) bool {
			line, ok := val.(*Line)
			if !ok || line == nil || key == cbc.CodeEmpty {
				return true
			}
			return tax.SetHasCategory(key).Check(line.Taxes)
		},
	)
}
