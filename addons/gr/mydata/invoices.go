package mydata

import (
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/gr"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Series, validation.Required),
		validation.Field(&inv.Tax,
			validation.Required,
			validation.By(validateInvoiceTax),
			validation.Skip,
		),
		validation.Field(&inv.Supplier,
			validation.By(validateBusinessParty),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(
				requiresValidCustomer(inv),
				validation.Required,
				validation.By(validateBusinessParty),
				validation.By(validateBusinessCustomer),
			),
			validation.Skip,
		),
		validation.Field(&inv.Lines,
			validation.Each(
				validation.By(validateInvoiceLine),
				validation.Skip,
			),
			validation.Skip,
		),
		validation.Field(&inv.Discounts,
			validation.Empty.Error("not supported by mydata"),
			validation.Skip,
		),
		validation.Field(&inv.Payment,
			validation.Required,
			validation.By(validateInvoicePaymentDetails),
			validation.Skip,
		),
		validation.Field(&inv.Preceding,
			validation.When(
				inv.Type.In(bill.InvoiceTypeCreditNote),
				validation.Required,
			),
			validation.Each(validation.By(validateInvoicePreceding)),
			validation.Skip,
		),
	)
}

func validateInvoiceTax(value any) error {
	t, ok := value.(*bill.Tax)
	if !ok || t == nil {
		return nil
	}
	return validation.ValidateStruct(t,
		validation.Field(&t.Ext,
			tax.ExtensionsRequire(ExtKeyInvoiceType),
			validation.Skip,
		),
	)
}

func validateBusinessParty(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
	)
}

func validateBusinessCustomer(value any) error {
	p, ok := value.(*org.Party)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.Addresses,
			validation.Required,
			validation.Length(1, 0),
			validation.Each(
				validation.By(validateInvoiceAddress),
			),
			validation.Skip,
		),
	)
}

func validateInvoiceAddress(value any) error {
	a, ok := value.(*org.Address)
	if !ok || a == nil {
		return nil
	}
	return validation.ValidateStruct(a,
		validation.Field(&a.Locality, validation.Required),
		validation.Field(&a.Code, validation.Required),
	)
}

func validateInvoiceLine(value any) error {
	l, ok := value.(*bill.Line)
	if !ok || l == nil {
		return nil
	}
	return validation.ValidateStruct(l,
		validation.Field(&l.Item,
			validation.By(validateInvoiceItem),
			validation.Skip,
		),
		validation.Field(&l.Total,
			num.Positive,
			num.NotZero,
			validation.Skip,
		),
	)
}

func validateInvoiceItem(value any) error {
	i, ok := value.(*org.Item)
	if !ok || i == nil {
		return nil
	}
	return validation.ValidateStruct(i,
		validation.Field(&i.Ext,
			validation.When(
				i.Ext.Has(ExtKeyIncomeCat) || i.Ext.Has(ExtKeyIncomeType),
				tax.ExtensionsRequire(ExtKeyIncomeCat, ExtKeyIncomeType),
			),
			validation.Skip,
		),
	)
}

func validateInvoicePaymentDetails(value any) error {
	p, ok := value.(*bill.PaymentDetails)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.Instructions,
			validation.When(
				len(p.Advances) == 0,
				validation.Required,
			),
			validation.Skip,
		),
	)
}

func validateInvoicePreceding(value any) error {
	p, ok := value.(*org.DocumentRef)
	if !ok || p == nil {
		return nil
	}
	return validation.ValidateStruct(p,
		validation.Field(&p.Stamps,
			head.StampsHas(gr.StampIAPRMark),
			validation.Skip,
		),
	)
}

// requiresValidCustomer returns true if the invoice type requires a customer.
func requiresValidCustomer(inv *bill.Invoice) bool {
	// Invoice type categories that require a valid customer.
	typeCats := []string{"1", "2", "5"}

	it := inv.Tax.Ext[ExtKeyInvoiceType].String()

	for _, prefix := range typeCats {
		if strings.HasPrefix(it, prefix+".") {
			return true
		}
	}

	return false
}
