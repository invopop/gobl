package saft

import (
	"slices"

	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// List of units typically used for services. Used to infer a default product type.
var serviceUnits = []org.Unit{
	org.UnitEmpty,
	org.UnitKilometre,
	org.UnitWatt,
	org.UnitKilowatt,
	org.UnitKilowattHour,
	org.UnitRate,
	org.UnitMonth,
	org.UnitDay,
	org.UnitSecond,
	org.UnitHour,
	org.UnitMinute,
	org.UnitService,
	org.UnitJob,
	org.UnitActivity,
	org.UnitTrip,

	// Ambivalent units but we default to service
	org.UnitItem,
	org.UnitOne,
}

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
		if slices.Contains(serviceUnits, item.Unit) {
			item.Ext[ExtKeyProductType] = ProductTypeService
		} else {
			item.Ext[ExtKeyProductType] = ProductTypeGoods
		}
	}
}

func setDefaultUnit(item *org.Item) {
	if item == nil {
		return
	}

	if item.Unit == "" {
		item.Unit = org.UnitOne
	}
}
