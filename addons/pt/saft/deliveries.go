package saft

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// validateDelivery ensures that the delivery has the required movement type extension.
func validateDelivery(dlv *bill.Delivery) error {
	dt := movementDocType(dlv)

	return validation.ValidateStruct(dlv,
		validation.Field(&dlv.Supplier,
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&dlv.Series,
			validateSeriesFormat(dt),
			validation.Skip,
		),
		validation.Field(&dlv.Code,
			validateCodeFormat(dlv.Series, dt),
			validation.Skip,
		),
		validation.Field(&dlv.Tax,
			validation.By(validateDeliveryTax),
			validation.Skip,
		),
		validation.Field(&dlv.DespatchDate, validation.Required),
	)
}

func movementDocType(dlv *bill.Delivery) cbc.Code {
	if dlv.Tax == nil || dlv.Tax.Ext == nil {
		return cbc.CodeEmpty
	}
	return dlv.Tax.Ext[ExtKeyMovementType]
}

func validateDeliveryTax(val any) error {
	t, _ := val.(*bill.Tax)
	if t == nil {
		// If no tax is given, init a blank one so that we can return meaningful
		// validation errors. The blank tax object is not assigned to the invoice
		// and so the original document doesn't actually change.
		t = new(bill.Tax)
	}

	return validation.ValidateStruct(t,
		validation.Field(&t.Ext,
			tax.ExtensionsRequire(ExtKeyMovementType),
			validation.Skip,
		),
	)
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
