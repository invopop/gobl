package nfse

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateInvoice(inv *bill.Invoice) error {
	if inv == nil {
		return nil
	}

	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Charges,
			validation.Empty.Error("not supported by nfse"),
			validation.Skip,
		),
		validation.Field(&inv.Discounts,
			validation.Empty.Error("not supported by nfse"),
			validation.Skip,
		),
	)
}

func validateSupplier(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}

	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
		validation.Field(&obj.Identities,
			org.RequireIdentityKey(IdentityKeyMunicipalReg),
			validation.Skip,
		),
		validation.Field(&obj.Name, validation.Required),
		validation.Field(&obj.Addresses,
			validation.Required,
			validation.Each(
				validation.Required,
				validation.By(validateSupplierAddress),
			),
			validation.Skip,
		),
		validation.Field(&obj.Ext,
			tax.ExtensionsRequires(
				ExtKeySimplesNacional,
				ExtKeyMunicipality,
			),
			validation.Skip,
		),
	)
}

func validateSupplierAddress(value interface{}) error {
	obj, _ := value.(*org.Address)
	if obj == nil {
		return nil
	}

	return validation.ValidateStruct(obj,
		validation.Field(&obj.Street, validation.Required),
		validation.Field(&obj.Number, validation.Required),
		validation.Field(&obj.Locality, validation.Required),
		validation.Field(&obj.State, validation.Required),
		validation.Field(&obj.Code, validation.Required),
	)
}
