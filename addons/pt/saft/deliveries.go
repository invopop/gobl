package saft

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billDeliveryRules() *rules.Set {
	return rules.For(new(bill.Delivery),
		rules.Assert("01",
			fmt.Sprintf("tax requires '%s' extension", ExtKeyMovementType),
			is.Func("has movement type", deliveryHasMovementType),
		),
		rules.Assert("02", "series format must be valid",
			is.FuncError("series format", deliverySeriesFormatValid),
		),
		rules.Assert("03", "code format must be valid",
			is.FuncError("code format", deliveryCodeFormatValid),
		),
		rules.Field("supplier",
			rules.Field("tax_id",
				rules.Assert("04", "supplier tax ID is required", is.Present),
				rules.Field("code",
					rules.Assert("05", "supplier tax ID code is required", is.Present),
				),
			),
		),
		rules.Field("despatch_date",
			rules.Assert("06", "cannot be blank", is.Present),
		),
	)
}

func deliveryHasMovementType(val any) bool {
	dlv, ok := val.(*bill.Delivery)
	if !ok || dlv == nil {
		return true
	}
	if dlv.Tax == nil || dlv.Tax.Ext == nil {
		return false
	}
	return tax.ExtensionsRequire(ExtKeyMovementType).Check(dlv.Tax.Ext)
}

func deliverySeriesFormatValid(val any) error {
	dlv, ok := val.(*bill.Delivery)
	if !ok || dlv == nil {
		return nil
	}
	return validateSeriesFormat(movementDocType(dlv)).Validate(dlv.Series)
}

func deliveryCodeFormatValid(val any) error {
	dlv, ok := val.(*bill.Delivery)
	if !ok || dlv == nil {
		return nil
	}
	dt := movementDocType(dlv)
	return validateCodeFormat(dlv.Series, dt).Validate(dlv.Code)
}

func movementDocType(dlv *bill.Delivery) cbc.Code {
	if dlv.Tax == nil || dlv.Tax.Ext == nil {
		return cbc.CodeEmpty
	}
	return dlv.Tax.Ext[ExtKeyMovementType]
}

func normalizeDelivery(dlv *bill.Delivery) {
	if dlv.Tax == nil {
		dlv.Tax = new(bill.Tax)
	}

	if dlv.Tax.Ext == nil {
		dlv.Tax.Ext = make(tax.Extensions)
	}

	if !dlv.Tax.Ext.Has(ExtKeyMovementType) {
		// Map delivery types to movement types
		switch dlv.Type {
		case bill.DeliveryTypeNote:
			dlv.Tax.Ext[ExtKeyMovementType] = MovementTypeDeliveryNote
		case bill.DeliveryTypeWaybill:
			dlv.Tax.Ext[ExtKeyMovementType] = MovementTypeWaybill
		}
	}
}
