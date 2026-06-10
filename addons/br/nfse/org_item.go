package nfse

import (
	"fmt"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

func orgItemRules() *rules.Set {
	return rules.For(new(org.Item),
		rules.Field("ext",
			rules.Assert("01", fmt.Sprintf("item requires '%s' extension", ExtKeyService),
				tax.ExtensionsRequire(ExtKeyService),
			),
			rules.Assert("02", fmt.Sprintf("item extensions '%s', '%s', and '%s' must all be present or all absent", ExtKeyOperation, ExtKeyTaxStatus, ExtKeyTaxClass),
				tax.ExtensionsRequireAllOrNone(
					ExtKeyOperation,
					ExtKeyTaxStatus,
					ExtKeyTaxClass,
				),
			),
		),
	)
}
