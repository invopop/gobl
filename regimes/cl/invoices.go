package cl

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// invoiceValidator validates Chilean invoices according to SII requirements.
//
// Validation enforces:
//   - Supplier: RUT and address required for all DTEs
//   - Customer: RUT and address required for Facturas (B2B), optional for Boletas (B2C)
//
// Use tax.TagSimplified to indicate a Boleta Electrónica (B2C receipt).
//
// Additional SII requirements (documented but not validated): Giro Comercial, Comuna,
// payment terms, and item descriptions per Resolution 36/2024.
//
// References:
//   - https://www.sii.cl/factura_electronica/factura_mercado/formato_dte.pdf
//   - https://www.sii.cl/factura_electronica/factura_mercado/formato_boletas_elec_202306.pdf
type invoiceValidator struct {
	inv *bill.Invoice
}

func validateInvoice(inv *bill.Invoice) error {
	if inv == nil {
		return nil
	}
	v := &invoiceValidator{inv: inv}
	return v.validate()
}

func (v *invoiceValidator) validate() error {
	inv := v.inv
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(v.supplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(v.customer),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) supplier(value interface{}) error {
	obj, _ := value.(*org.Party)
	if obj == nil {
		return nil
	}
	return validation.ValidateStruct(obj,
		// RUT required for all DTEs (mandatory since 2018/2021 for B2B/B2C)
		validation.Field(&obj.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
		// Address required (Dirección, Ciudad, Comuna)
		validation.Field(&obj.Addresses,
			validation.Required,
			validation.Length(1, 0),
			validation.Skip,
		),
	)
}

func (v *invoiceValidator) customer(value interface{}) error {
	inv := v.inv
	// Facturas (B2B): RUT and address required
	// Boletas (B2C with tax.TagSimplified): RUT and address optional
	// Note: High-value boletas (>135 UF) will require RUT from Sep 2025
	isB2B := !inv.Tags.HasTags(tax.TagSimplified)

	obj, _ := value.(*org.Party)
	if obj == nil {
		if isB2B {
			return validation.NewError("validation_required", "customer is required for B2B invoices")
		}
		return nil
	}

	return validation.ValidateStruct(obj,
		validation.Field(&obj.TaxID,
			validation.When(
				isB2B,
				validation.Required,
				tax.RequireIdentityCode,
			),
			validation.Skip,
		),
		validation.Field(&obj.Addresses,
			validation.When(
				isB2B,
				validation.Required,
				validation.Length(1, 0),
			),
			validation.Skip,
		),
	)
}
