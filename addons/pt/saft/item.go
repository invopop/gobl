package saft

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
		validation.Field(&item.Unit, validation.Required),
		validation.Field(&item.Ext,
			tax.ExtensionsRequire(ExtKeyProductType),
			validation.Skip,
		),
	)
}

func normalizeItem(item *org.Item) {
	setDefaultUnit(item)
	setDefaultProductType(item)
}

func setDefaultProductType(item *org.Item) {
	if item == nil {
		return
	}

	if item.Ext == nil {
		item.Ext = make(tax.Extensions)
	}

	if _, ok := item.Ext[ExtKeyProductType]; !ok {
		item.Ext[ExtKeyProductType] = ProductTypeService
	}
}

func setDefaultUnit(item *org.Item) {
	if item == nil {
		return
	}

	if item.Unit == "" {
		item.Unit = org.UnitItem
	}
}
