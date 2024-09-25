package au

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(supplier),
			validation.Skip,
		),
	)
}

func supplier(val any) error {
	obj, _ := val.(*org.Party)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.By(checkSupplierCountry),
			validation.Skip,
		),
		validation.Field(&obj.Name,
			validation.Required,
			validation.Skip,
		),
	)
}

// Supplier must have ABN, therefore be Australian
func checkSupplierCountry(value interface{}) error {
	obj, _ := value.(*tax.Identity)
	if obj == nil {
		return nil
	}

	return validation.ValidateStruct(obj,
		validation.Field(&obj.Country,
			validation.In(l10n.AU.Tax()),
			validation.Skip,
		),
	)
}
