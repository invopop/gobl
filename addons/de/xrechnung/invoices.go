package xrechnung

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

const (
	invoiceTypeSelfBilled               cbc.Key = "self-billed"
	invoiceTypePartial                  cbc.Key = "partial"
	invoiceTypePartialConstruction      cbc.Key = "partial-construction"
	invoiceTypePartialFinalConstruction cbc.Key = "partial-final-construction"
	invoiceTypeFinalConstruction        cbc.Key = "final-construction"
)

var validTypes = []cbc.Key{
	bill.InvoiceTypeStandard,
	bill.InvoiceTypeCreditNote,
	bill.InvoiceTypeCorrective,
	invoiceTypeSelfBilled,
	invoiceTypePartial,
	invoiceTypePartialConstruction,
	invoiceTypePartialFinalConstruction,
	invoiceTypeFinalConstruction,
}

func ValidateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		// BR-DE-17
		validation.Field(&inv.Type,
			validation.By(validateInvoiceType),
		),
		// BR-DE-01
		validation.Field(&inv.Payment.Instructions,
			validation.Required,
			validation.By(validatePaymentInstructions),
		),
		// BR-DE-15
		validation.Field(&inv.Ordering.Code,
			validation.Required,
		),
		validation.Field(&inv.Supplier,
			validation.By(validateSupplier),
		),
		validation.Field(&inv.Customer,
			validation.By(validateCustomer),
		),
		validation.Field(&inv.Delivery,
			validation.When(inv.Delivery != nil,
				validation.By(validateDeliveryParty),
			),
		),
		// BR-DE-26
		validation.Field(&inv,
			validation.By(validateCorrectiveInvoice),
		),
	)
}

func validateInvoiceType(value interface{}) error {
	t, ok := value.(cbc.Key)
	if !ok {
		return validation.NewError("type", "Invalid invoice type")
	}
	if t.In(validTypes...) {
		return nil
	}
	return validation.NewError("invalid", "Invalid invoice type")
}

func validateSupplier(value interface{}) error {
	p, _ := value.(*org.Party)
	if p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		// BR-DE-02
		validation.Field(&p.Name,
			validation.Required,
		),
		// BR-DE-03, BR-DE-04
		validation.Field(&p.Addresses,
			validation.Required,
			validation.Each(validation.By(validatePartyAddress)),
		),
		// BR-DE-06
		validation.Field(&p.People,
			validation.Required,
		),
		// BR-DE-05
		validation.Field(&p.Telephones,
			validation.Required,
		),
		// BR-DE-07
		validation.Field(&p.Emails,
			validation.Required,
		),
	)
}

func validateCustomer(value interface{}) error {
	p, _ := value.(*org.Party)
	if p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		// BR-DE-08, BR-DE-09
		validation.Field(&p.Addresses,
			validation.Required,
			validation.Each(validation.By(validatePartyAddress)),
		),
	)
}

func validatePartyAddress(value interface{}) error {
	addr, _ := value.(*org.Address)
	if addr == nil {
		return nil
	}
	return validation.ValidateStruct(addr,
		validation.Field(&addr.Locality,
			validation.Required,
		),
		validation.Field(&addr.Code,
			validation.Required,
		),
	)
}

func validateDeliveryParty(value interface{}) error {
	p, _ := value.(*org.Party)
	if p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.Addresses,
			validation.Required,
			validation.Each(validation.By(validateDeliveryAddress)),
		),
	)
}

func validateDeliveryAddress(value interface{}) error {
	addr, _ := value.(*org.Address)
	if addr == nil {
		return nil
	}
	return validation.ValidateStruct(addr,
		// BR-DE-10
		validation.Field(&addr.Locality,
			validation.Required,
		),
		// BR-DE-11
		validation.Field(&addr.Code,
			validation.Required,
		),
	)
}

func validateCorrectiveInvoice(value interface{}) error {
	inv, ok := value.(*bill.Invoice)
	if !ok || inv == nil {
		return nil
	}
	if inv.Type.In(bill.InvoiceTypeCorrective) {
		if inv.Preceding == nil {
			return validation.NewError("required", "Preceding invoice details are required for corrective invoices")
		}
	}
	return nil
}
