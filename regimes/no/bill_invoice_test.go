package no

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func TestValidateBillInvoiceSupplierRequiresTaxID(t *testing.T) {
	inv := &bill.Invoice{
		Supplier: &org.Party{
			Name: "Test Supplier",
		},
	}
	inv.Tags.SetTags(tax.TagSimplified)

	if err := validateBillInvoice(inv); err == nil {
		t.Fatalf("expected error when supplier tax id is missing")
	}
}

func TestValidateBillInvoiceSupplierWithTaxIDPasses(t *testing.T) {
	inv := &bill.Invoice{
		Supplier: &org.Party{
			Name: "Test Supplier",
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.NO),
				Code:    "974760673MVA",
			},
		},
	}
	inv.Tags.SetTags(tax.TagSimplified)

	if err := validateBillInvoice(inv); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
