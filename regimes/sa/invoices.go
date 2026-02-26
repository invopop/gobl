package sa

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(
				!isSimplified(inv),
				validation.Required,
				validation.By(validateCustomer),
			),
			// NOTE: BR-KSA-25 requires buyer name and National ID on simplified
			// invoices for education (VATEX-SA-EDU) and healthcare (VATEX-SA-HEA)
			// exemptions. This is handled by the ZATCA addon which defines
			// the VATEX-SA extension codes.
			validation.Skip,
		),
	)
}

func validateSupplier(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		// BR-06: Seller name required on all invoices.
		validation.Field(&p.Name,
			validation.Required,
		),
	)
}

func validateCustomer(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		// BR-KSA-42: Standard invoices require a customer name.
		validation.Field(&p.Name,
			validation.Required,
		),
		// BR-KSA-81: Buyer must have either a VAT registration number (TaxID)
		// or an alternative buyer identification (org.Identity).
		validation.Field(&p.TaxID,
			validation.When(
				!hasIdentities(p),
				validation.Required,
				tax.RequireIdentityCode,
			),
			validation.Skip,
		),
		validation.Field(&p.Identities,
			validation.When(
				!hasTaxIDCode(p),
				validation.Required,
			),
			validation.Skip,
		),
	)
}

func isSimplified(inv *bill.Invoice) bool {
	return inv.HasTags(tax.TagSimplified)
}

func hasTaxIDCode(party *org.Party) bool {
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

func hasIdentities(party *org.Party) bool {
	return party != nil && len(party.Identities) > 0
}
