package nz

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Thresholds for NZ GST invoice requirements
var (
	threshold200  = num.MakeAmount(20000, 2)  // $200.00
	threshold1000 = num.MakeAmount(100000, 2) // $1000.00
)

// validateInvoice checks NZ-specific invoice requirements based on GST thresholds.
func validateInvoice(inv *bill.Invoice) error {
	if inv == nil {
		return nil
	}
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateInvoiceSupplier(inv)),
		),
		validation.Field(&inv.Customer,
			validation.By(validateInvoiceCustomer(inv)),
		),
	)
}

func validateInvoiceSupplier(inv *bill.Invoice) validation.RuleFunc {
	return func(value any) error {
		nzdTotal, err := getNZDTotal(inv)
		if err != nil {
			return err
		}
		p, ok := value.(*org.Party)
		if !ok || p == nil {
			if nzdTotal.IsZero() || nzdTotal.Compare(threshold200) <= 0 {
				return nil
			}
			return validation.NewError("validation_supplier_missing", "supplier details required for invoices > $200")
		}
		if nzdTotal.IsZero() || nzdTotal.Compare(threshold200) <= 0 {
			return nil
		}

		return validation.ValidateStruct(p,
			validation.Field(&p.TaxID,
				validation.Required.Error("supplier must have GST number for invoices > $200"),
				tax.RequireIdentityCode,
			),
		)
	}
}

func validateInvoiceCustomer(inv *bill.Invoice) validation.RuleFunc {
	return func(value any) error {
		nzdTotal, err := getNZDTotal(inv)
		if err != nil {
			return err
		}
		p, ok := value.(*org.Party)
		if !ok || p == nil {
			if nzdTotal.IsZero() || nzdTotal.Compare(threshold1000) <= 0 {
				return nil
			}
			return validation.NewError("validation_customer_missing", "customer details required for invoices > $1,000")
		}
		if nzdTotal.IsZero() || nzdTotal.Compare(threshold1000) <= 0 {
			return nil
		}

		if err := validation.ValidateStruct(p,
			validation.Field(&p.Name,
				validation.Required.Error("customer name required for invoices > $1,000"),
			),
		); err != nil {
			return err
		}

		return requireCustomerIdentifier(p)
	}
}

func requireCustomerIdentifier(value any) error {
	// No need to check type again since this is only called after validating the customer is a *org.Party
	p, _ := value.(*org.Party)

	if hasCustomerIdentifier(p) {
		return nil
	}

	return validation.NewError(
		"validation_customer_identifier",
		"customer must have at least one identifier (address, phone, email, tax ID, or identities) for invoices > $1,000",
	)
}

func hasCustomerIdentifier(p *org.Party) bool {
	if len(p.Addresses) > 0 {
		return true
	}
	if len(p.Emails) > 0 {
		return true
	}
	if len(p.Telephones) > 0 {
		return true
	}
	if p.TaxID != nil && p.TaxID.Code != "" {
		return true
	}
	if len(p.Identities) > 0 {
		return true
	}
	return false
}

func getNZDTotal(inv *bill.Invoice) (*num.Amount, error) {
	if inv.Totals == nil {
		return nil, validation.NewError("validation_totals_missing", "invoice totals must be calculated before NZ regime validation")
	}
	total := inv.Totals.TotalWithTax
	if inv.Currency == "NZD" {
		return &total, nil
	}
	nzdTotal := currency.Convert(inv.ExchangeRates, inv.Currency, currency.NZD, total)
	if nzdTotal == nil {
		return nil, validation.NewError("validation_currency_conversion", "cannot convert invoice total to NZD for threshold validation")
	}
	return nzdTotal, nil
}
