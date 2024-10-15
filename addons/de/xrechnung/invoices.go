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
		validation.Field(&inv.Payment, validation.Required),
		validation.Field(&inv.Payment,
			validation.By(func(value interface{}) error {
				payment, ok := value.(*bill.Payment)
				if !ok || payment == nil {
					return validation.NewError("payment_type", "must be a valid non-empty Payment type")
				}
				return validation.ValidateStruct(payment,
					validation.Field(&payment.Instructions,
						validation.Required,
						validation.By(validatePaymentInstructions),
					),
				)
			}),
		),
		// BR-DE-15
		validation.Field(&inv.Ordering, validation.Required),
		validation.Field(&inv.Ordering,
			validation.By(func(value interface{}) error {
				ordering, ok := value.(*bill.Ordering)
				if !ok || ordering == nil {
					return validation.NewError("ordering_type", "must be a valid Ordering type")
				}
				return validation.ValidateStruct(ordering,
					validation.Field(&ordering.Code, validation.Required),
				)
			}),
		),
		validation.Field(&inv.Supplier,
			validation.By(validateSupplier),
		),
		validation.Field(&inv.Supplier,
			validation.By(validateSellerTaxInfo),
		),
		validation.Field(&inv.Customer,
			validation.By(validateCustomerReceiver),
		),
		validation.Field(&inv.Delivery,
			validation.By(validateDelivery),
		),
		// BR-DE-26
		validation.Field(&inv.Preceding,
			validation.When(inv.Type.In(bill.InvoiceTypeCorrective),
				validation.Required,
			),
		),
	)
}

func validateInvoiceType(value interface{}) error {
	t, ok := value.(cbc.Key)
	if !ok {
		return validation.NewError("type", "Invalid invoice type")
	}
	if !t.In(validTypes...) {
		return validation.NewError("invalid", "Invalid invoice type")
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

// func validateSellerTaxInfo(value interface{}) error {
// 	supplier, ok := value.(*org.Party)
// 	if !ok || supplier == nil {
// 		return validation.NewError("invalid_supplier", "Supplier is invalid or nil")
// 	}
// 	hasVATIdentifier := supplier.TaxID != nil && supplier.TaxID.Code != ""
// 	hasTaxIdentifier := org.IdentityForKey(supplier.Identities, "de-tax-number") != nil

// 	if !hasVATIdentifier && !hasTaxIdentifier {
// 		return validation.NewError(
// 			"missing_seller_tax_info",
// 			"Either Seller VAT identifier or Seller tax identifier must be provided",
// 		)
// 	}

// 	return nil
// }

func validateSellerTaxInfo(value interface{}) error {
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
			),
		),
	)
}

func validateTaxNumber(value interface{}) error {
	identities, ok := value.([]*org.Identity)
	if !ok {
		return validation.NewError("invalid_identities", "Identities are invalid")
	}
	if org.IdentityForKey(identities, "de-tax-number") == nil {
		return validation.NewError("missing_tax_identifier", "German tax identifier (de-tax-number) is required")
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
			validation.By(validateCustomerReceiver)),
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
