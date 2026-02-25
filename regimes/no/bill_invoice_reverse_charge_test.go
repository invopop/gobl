package no

import (
	"testing"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

func TestReverseChargeRequiresCustomerTaxID(t *testing.T) {
	inv := &bill.Invoice{
		Supplier: &org.Party{
			Name: "Supplier",
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.NO),
				Code:    "974760673MVA",
			},
			Addresses: []*org.Address{{Country: l10n.ISOCountryCode("NO")}},
		},
		Customer: &org.Party{
			Name:      "Customer",
			Addresses: []*org.Address{{Country: l10n.ISOCountryCode("NO")}},
		},
	}
	inv.Tags.SetTags(tax.TagReverseCharge)

	if err := validateBillInvoice(inv); err == nil {
		t.Fatalf("expected error when reverse charge invoice customer has no tax id")
	}
}

func TestReverseChargeWithCustomerTaxIDPasses(t *testing.T) {
	inv := &bill.Invoice{
		Supplier: &org.Party{
			Name: "Supplier",
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.NO),
				Code:    "974760673MVA",
			},
			Addresses: []*org.Address{{Country: l10n.ISOCountryCode("NO")}},
		},
		Customer: &org.Party{
			Name: "Customer",
			TaxID: &tax.Identity{
				Country: l10n.TaxCountryCode(l10n.NO),
				Code:    "974760673",
			},
			Addresses: []*org.Address{{Country: l10n.ISOCountryCode("NO")}},
		},
	}
	inv.Tags.SetTags(tax.TagReverseCharge)

	if err := validateBillInvoice(inv); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
