package xrechnung

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/validation"
)

var validTypes = []cbc.Key{
	bill.InvoiceTypeStandard,
	bill.InvoiceTypeCreditNote,
	bill.InvoiceTypeCorrective,
}

// ValidateInvoice validates the invoice according to the XRechnung standard
func ValidateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		// BR-DE-17
		validation.Field(&inv.Type,
			validation.By(validateInvoiceType),
			validation.Skip,
		),
		// BR-DE-01
		validation.Field(&inv.Payment,
			validation.Required,
			validation.By(validatePayment),
			validation.Skip,
		),
		// BR-DE-15
		validation.Field(&inv.Ordering,
			validation.Required,
			validation.By(validateOrdering),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(validateSupplierTaxInfo),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.By(validateCustomerReceiver),
			validation.Skip,
		),
		validation.Field(&inv.Delivery,
			validation.By(validateDelivery),
			validation.Skip,
		),
		// BR-DE-26
		validation.Field(&inv.Preceding,
			validation.When(inv.Type.In(bill.InvoiceTypeCorrective),
				validation.Required,
			),
		),
	)
}

func validatePayment(value interface{}) error {
	payment, ok := value.(*bill.Payment)
	if !ok || payment == nil {
		return nil
	}
	return validation.ValidateStruct(payment,
		validation.Field(&payment.Instructions,
			validation.Required,
			validation.By(validatePaymentInstructions),
		),
	)
}

func validateOrdering(value interface{}) error {
	ordering, ok := value.(*bill.Ordering)
	if !ok || ordering == nil {
		return nil
	}
	return validation.ValidateStruct(ordering,
		validation.Field(&ordering.Code,
			validation.Required,
		),
	)
}

func validateInvoiceType(value interface{}) error {
	t, ok := value.(cbc.Key)
	if !ok {
		return validation.NewError("type", "invalid invoice type")
	}
	if !t.In(validTypes...) {
		return validation.NewError("invalid", "invalid invoice type")
	}
	return nil
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
			validation.Skip,
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

func validateSupplierTaxInfo(value interface{}) error {
	supplier, ok := value.(*org.Party)
	if !ok || supplier == nil {
		return validation.NewError("invalid_supplier", "Supplier is invalid or nil")
	}

	return validation.ValidateStruct(supplier,
		validation.Field(&supplier.TaxID,
			validation.When(supplier.Identities == nil || org.IdentityForKey(supplier.Identities, "de-tax-number") == nil,
				validation.Required,
			),
		),
		validation.Field(&supplier.Identities,
			validation.When(supplier.TaxID == nil || supplier.TaxID.Code == "",
				validation.Required,
				validation.By(validateTaxNumber),
				validation.Skip,
			),
		),
	)
}

func validateTaxNumber(value interface{}) error {
	identities, ok := value.([]*org.Identity)
	if !ok {
		return validation.NewError("invalid_identities", "identities are invalid")
	}
	if org.IdentityForKey(identities, "de-tax-number") == nil {
		return validation.NewError("missing_tax_identifier", "tax identifier (de-tax-number) is required")
	}
	return nil
}

func validateDelivery(value interface{}) error {
	d, _ := value.(*bill.Delivery)
	if d == nil {
		return nil
	}
	return validation.ValidateStruct(d,
		validation.Field(&d.Receiver,
			validation.By(validateCustomerReceiver),
			validation.Skip,
		),
	)
}

// As the fields for customer and delivery reciver have the same requirements
// they are handled by the same validation function.
func validateCustomerReceiver(value interface{}) error {
	p, _ := value.(*org.Party)
	if p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		// BR-DE-08, BR-DE-09
		validation.Field(&p.Addresses,
			validation.Required,
			validation.Each(validation.By(validatePartyAddress)),
			validation.Skip,
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
