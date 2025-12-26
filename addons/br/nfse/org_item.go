package nfse

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateItem(item *org.Item) error {
	if item == nil {
		return nil
	}

	return validation.ValidateStruct(item,
		validation.Field(&item.Ext,
			tax.ExtensionsRequire(ExtKeyService),
			tax.ExtensionsRequireAllOrNone(
				ExtKeyOperation,
				ExtKeyTaxStatus,
				ExtKeyTaxClass,
			),
			validation.Skip,
		),
	)
}
