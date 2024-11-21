package nfse

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

const (
	// FiscalIncentiveDefault is the default value for the fiscal incentive extenstion
	FiscalIncentiveDefault = "2" // No incentiva
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
				ExtKeyFiscalIncentive,
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

func normalizeSupplier(sup *org.Party) {
	if sup == nil {
		return
	}

	if !sup.Ext.Has(ExtKeyFiscalIncentive) {
		if sup.Ext == nil {
			sup.Ext = make(tax.Extensions)
		}
		sup.Ext[ExtKeyFiscalIncentive] = FiscalIncentiveDefault
	}
}
