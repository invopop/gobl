package nfse

import (
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateItem(value any) error {
	item, _ := value.(*org.Item)
	if item == nil {
		return nil
	}

	return validation.ValidateStruct(item,
		validation.Field(&item.Ext,
			tax.ExtensionsRequires(ExtKeyService),
			validation.Skip,
		),
	)
}
