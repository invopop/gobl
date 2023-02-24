package it

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// invoiceValidator adds validation checks to invoices which are relevant
// for the region.
type invoiceValidator struct {
	inv *bill.Invoice
}

func validateInvoice(inv *bill.Invoice) error {
	v := &invoiceValidator{inv: inv}
	return v.validate()
}

func (v *invoiceValidator) validate() error {
	inv := v.inv
	return validation.ValidateStruct(inv,
		// Currently only supporting invoices and credit notes
		validation.Field(&inv.Type, validation.In(
			bill.InvoiceTypeNone,
			bill.InvoiceTypeCreditNote,
		)),
		validation.Field(&inv.Supplier, validation.Required, validation.By(v.supplier)),
		validation.Field(&inv.Customer, validation.Required, validation.By(v.customer)),
		validation.Field(&inv.Meta, validation.Required, validation.By(v.meta)),
	)
}

func (v *invoiceValidator) supplier(value interface{}) error {
	sup, _ := value.(*org.Party)
	if sup == nil {
		return nil
	}
	// Suppliers must have a VAT ID (Partita IVA)
	return validation.ValidateStruct(sup,
		validation.Field(&sup.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
		),
	)
}

func (v *invoiceValidator) customer(value interface{}) error {
	cus, _ := value.(*org.Party)
	if cus == nil {
		return nil
	}
	if cus.TaxID == nil {
		return nil // validation already handled, this prevents panics
	}

	var fiscalCode *org.Identity
	for _, identity := range cus.Identities {
		if identity.Type == IdentityTypeFiscalCode {
			fiscalCode = identity
			break
		}
	}

	// Customers must have either a VAT ID (Partita IVA) or a fiscal code (Codice
	// Fiscale)
	return validation.ValidateStruct(cus,
		validation.Field(fiscalCode, validation.When(
			cus.TaxID == nil,
			validation.Required,
		)),
		validation.Field(&cus.TaxID, validation.When(
			fiscalCode == nil,
			validation.Required,
			tax.RequireIdentityCode,
		)),
	)
}

// Meta data of the invoice must contain valid FatturaPA codes
func (v *invoiceValidator) meta(value interface{}) error {
	meta, _ := value.(cbc.Meta)
	if meta == nil {
		return nil
	}

	// Validate Tax System FPA code. Currently only supporting
	// RF01 "Ordinary" tax system
	err := validation.Validate(meta[FPACodeTypeTaxSystem],
		validation.Required,
		validation.In(FPACodeTaxSystemOrdinary),
	)
	if err != nil {
		return err
	}

	// Validate Document Type FPA code
	invoiceType := v.inv.Type
	codeGroup := InvoiceTypeMap[invoiceType]
	err = validateMetaFPACode(meta, codeGroup)
	if err != nil {
		return err
	}

	// Validate Payment Method FPA code
	payMethod := v.inv.Payment.Instructions.Key
	codeGroup = PaymentMethodMap[payMethod]
	err = validateMetaFPACode(meta, codeGroup)
	if err != nil {
		return err
	}

	// Validate Tax Scheme FPA code
	// Currently only reverse charge schemes are supported.
	for _, scheme := range v.inv.Tax.Schemes {
		codeGroup = SchemeMap[scheme]
		err = validateMetaFPACode(meta, codeGroup)
		if err != nil {
			return err
		}
	}

	// Validate Withholding Tax FPA code
	// Only retained taxes are considered.
	for _, catTotal := range v.inv.Totals.Taxes.Categories {
		if !catTotal.Retained {
			break
		}
		codeGroup = TaxCategoryMap[catTotal.Code]
		err = validateMetaFPACode(meta, codeGroup)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateMetaFPACode(meta cbc.Meta, codeGroup *FPACodeGroup) error {
	if codeGroup == nil {
		return nil
	}

	key := codeGroup.Type
	fpaCodes := codeGroup.Codes
	codes := make([]interface{}, len(fpaCodes))
	for i, code := range fpaCodes {
		codes[i] = code
	}

	return validation.Validate(meta[key],
		validation.Required,
		validation.In(codes...),
	)
}
