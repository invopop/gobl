package no

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func validSupplierSimplified() *org.Party {
	return &org.Party{
		Name: "Supplier",
		TaxID: &tax.Identity{
			Country: l10n.TaxCountryCode(l10n.NO),
			Code:    "974760673MVA",
		},
	}
}

func validSupplierNonSimplified() *org.Party {
	return &org.Party{
		Name: "Supplier",
		TaxID: &tax.Identity{
			Country: l10n.TaxCountryCode(l10n.NO),
			Code:    "974760673MVA",
		},
		Addresses: []*org.Address{{Country: l10n.ISOCountryCode("NO")}},
	}
}

func TestValidateBillInvoiceNilIsNoop(t *testing.T) {
	if err := validateBillInvoice(nil); err != nil {
		t.Fatalf("expected no error for nil invoice, got %v", err)
	}
}

func TestNonSimplifiedInvoiceRequiresCustomer(t *testing.T) {
	inv := &bill.Invoice{
		Supplier: validSupplierNonSimplified(),
		Customer: nil,
	}
	if err := validateBillInvoice(inv); err == nil {
		t.Fatalf("expected error when customer is missing in non-simplified invoice")
	}
}

func TestNonSimplifiedInvoiceCustomerRequiresAddress(t *testing.T) {
	inv := &bill.Invoice{
		Supplier: validSupplierNonSimplified(),
		Customer: &org.Party{
			Name: "Customer",
		},
	}
	if err := validateBillInvoice(inv); err == nil {
		t.Fatalf("expected error when customer address is missing in non-simplified invoice")
	}
}

func TestSimplifiedInvoiceRejectsNonEmptyCustomer(t *testing.T) {
	inv := &bill.Invoice{
		Supplier: validSupplierSimplified(),
		Customer: &org.Party{
			Name: "Customer",
		},
	}
	inv.Tags.SetTags(tax.TagSimplified)

	if err := validateBillInvoice(inv); err == nil {
		t.Fatalf("expected error when simplified invoice includes customer")
	}
}

func TestNonReverseChargeDoesNotRequireCustomerTaxID(t *testing.T) {
	inv := &bill.Invoice{
		Supplier: validSupplierNonSimplified(),
		Customer: &org.Party{
			Name:      "Customer",
			Addresses: []*org.Address{{Country: l10n.ISOCountryCode("NO")}},
		},
	}
	if err := validateBillInvoice(inv); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestReverseChargeWithNilCustomerStillFailsBecauseCustomerRequired(t *testing.T) {
	inv := &bill.Invoice{
		Supplier: validSupplierNonSimplified(),
		Customer: nil,
	}
	inv.Tags.SetTags(tax.TagReverseCharge)

	if err := validateBillInvoice(inv); err == nil {
		t.Fatalf("expected error when customer missing in reverse charge invoice")
	}
}
