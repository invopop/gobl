package nfse

import (
	"regexp"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

var (
	validStates = []cbc.Code{
		"AC", "AL", "AM", "AP", "BA", "CE", "DF", "ES", "GO",
		"MA", "MG", "MS", "MT", "PA", "PB", "PE", "PI", "PR",
		"RJ", "RN", "RO", "RR", "RS", "SC", "SE", "SP", "TO",
	}

	validAddressCode = regexp.MustCompile(`^(?:\D*\d){8}\D*$`)
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
		),
		validation.Field(&obj.Identities,
			org.RequireIdentityKey(IdentityKeyMunicipalReg),
		),
		validation.Field(&obj.Name, validation.Required),
		validation.Field(&obj.Addresses,
			validation.Required,
			validation.Each(
				validation.Required,
				validation.By(validateSupplierAddress),
			),
		),
		validation.Field(&obj.Ext,
			tax.ExtensionsRequires(
				ExtKeySimplesNacional,
				ExtKeyMunicipality,
			),
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
		validation.Field(&obj.State,
			validation.Required,
			validation.In(validStates...),
		),
		validation.Field(&obj.Code,
			validation.Required,
			validation.Match(validAddressCode),
		),
	)
}
