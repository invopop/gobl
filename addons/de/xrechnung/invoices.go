package xrechnung

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// BR-DE-14 - set percent in convertor as this rule requires it, even if it is 0. BT-119
// BR-DE-16 - covered by tax ID being required. BT-31
// BR-DE-18 - references format of payment.terms.details. BT-20
// BR-DE-20 - partialy covered by gobl validation of IBAN format. BT-84
// BR-DE-21 - look at BT-24 mapping of gobl. BT-24
// BR-DE-22 - refers to attachments. BG-24
// BR-DE-27 - handled by gobl validation of phone number. BT-42
// BR-DE-28 - handled by gobl validation of email address. BT-43

// BR-DE-17 - restricted subset of UNTDID document type codes
var validInvoiceUNTDIDDocumentTypeValues = []cbc.Code{
	"326", // Partial
	"380", // Commercial
	"384", // Corrected
	"389", // Self-billed
	"381", // Credit note
	"875", // Partial construction invoice
	"876", // Partial Final construction invoice
	"877", // Final construction invoice
}

// validateInvoice validates the invoice according to the XRechnung standard
func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		// BR-DE-17
		validation.Field(&inv.Tax,
			validation.By(validateInvoiceTax),
			validation.Skip,
		),
		// BR-DE-26
		validation.Field(&inv.Preceding,
			validation.When(
				inv.Type.In(
					bill.InvoiceTypeCorrective,
					bill.InvoiceTypeCreditNote,
				),
				validation.Required,
			),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.Required,
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.Required,
			validation.By(validateCustomer),
			validation.Skip,
		),
		// BR-DE-1
		validation.Field(&inv.Payment,
			validation.Required,
			validation.By(validatePayment),
			validation.Skip,
		),
		validation.Field(&inv.Delivery,
			validation.By(validateDelivery),
			validation.Skip,
		),
		validation.Field(&inv.Ordering,
			validation.Required,
			validation.By(validateOrdering),
			validation.Skip,
		),
	)
}

func validateInvoiceTax(value any) error {
	tx, ok := value.(*bill.Tax)
	if !ok || tx == nil {
		return nil
	}
	return validation.ValidateStruct(tx,
		validation.Field(&tx.Ext,
			tax.ExtensionsHasCodes(untdid.ExtKeyTaxCategory, validInvoiceUNTDIDDocumentTypeValues...),
			validation.Skip,
		),
	)
}

func validateDelivery(val any) error {
	delivery, ok := val.(*bill.DeliveryDetails)
	if !ok || delivery == nil {
		return nil
	}

	return validation.ValidateStruct(delivery,
		validation.Field(&delivery.Receiver,
			validation.Required,
			validation.By(validateCustomer),
			validation.Skip,
		),
	)
}

func validateOrdering(val any) error {
	ordering, ok := val.(*bill.Ordering)
	if !ok || ordering == nil {
		return nil
	}

	// BR-DE-15
	if len(ordering.Identities) == 0 && ordering.Code == "" {
		return validation.NewError("required", "either ordering code or identities with codes are required")
	}

	return nil
}

func validatePayment(val any) error {
	payment, ok := val.(*bill.PaymentDetails)
	if !ok || payment == nil {
		return nil
	}

	return validation.ValidateStruct(payment,
		validation.Field(&payment.Instructions,
			validation.Required,
			validation.By(validatePaymentInstructions),
			validation.Skip,
		),
	)
}

func validateSupplier(val any) error {
	party, ok := val.(*org.Party)
	if !ok || party == nil {
		return nil
	}

	// Check if either party or people have telephones/emails
	// BR-DE-6/BR-DE-7
	hasPeopleTelephones := len(party.People) > 0 && len(party.People[0].Telephones) > 0
	hasPeopleEmails := len(party.People) > 0 && len(party.People[0].Emails) > 0

	return validation.ValidateStruct(party,
		// BR-DE-2
		validation.Field(&party.People,
			// BR-DE-5
			validation.Required,
			validation.Skip,
		),
		validation.Field(&party.Addresses,
			validation.Required,
			validation.By(validateAddresses),
			validation.Skip,
		),
		// Check for either party or people telephones
		validation.Field(&party.Telephones,
			validation.When(
				!hasPeopleTelephones,
				validation.Required.Error("either party.telephones or party.people[0].telephones is required"),
			),
			validation.Skip,
		),
		// Check for either party or people emails
		validation.Field(&party.Emails,
			validation.When(
				!hasPeopleEmails,
				validation.Required.Error("either party.emails or party.people[0].emails is required"),
			),
			validation.Skip,
		),
	)
}

func validateCustomer(val any) error {
	party, ok := val.(*org.Party)
	if !ok || party == nil {
		return nil
	}
	return validation.ValidateStruct(party,
		validation.Field(&party.Addresses,
			validation.Required,
			validation.By(validateAddresses),
			validation.Skip,
		),
	)
}

func validateAddresses(val any) error {
	addresses, ok := val.([]*org.Address)
	if !ok || len(addresses) == 0 {
		return nil
	}

	return validation.ValidateStruct(addresses[0],
		// BR-DE-3/BR-DE-8/BR-DE-10
		validation.Field(&addresses[0].Locality,
			validation.Required,
			validation.Skip,
		),
		// BR-DE-4/BR-DE-9/BR-DE-11
		validation.Field(&addresses[0].Code,
			validation.Required,
			validation.Skip,
		),
	)
}
