package nfe

import (
	"fmt"
	"regexp"

	"github.com/invopop/gobl/addons/br/dfe"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

// Validation patterns
const (
	seriesPattern = `^(?:0|[1-9]{1}[0-9]{0,2})$` // extracted from the NFe XSD to validate the series
)

// Validation regexps
var (
	seriesRegex = regexp.MustCompile(seriesPattern)
)

func validateInvoice(inv *bill.Invoice) error {
	return validation.ValidateStruct(inv,
		validation.Field(&inv.Supplier,
			validation.By(validateParty),
			validation.By(validateSupplier),
			validation.Skip,
		),
		validation.Field(&inv.Customer,
			validation.When(isNFe(inv.Tax),
				validation.Required,
				validation.Skip,
			),
			validation.By(validateParty),
			validation.By(validateCustomer(inv)),
			validation.Skip,
		),
		validation.Field(&inv.Series,
			validation.Required,
			validation.Match(seriesRegex),
			validation.Skip,
		),
		validation.Field(&inv.Tax,
			validation.By(validateInvoiceTax),
			validation.Required,
			validation.Skip,
		),
		validation.Field(&inv.Notes,
			validation.Each(validation.By(validateInvoiceNote)),
			validation.By(validateInvoiceNotes),
			validation.Skip,
		),
		validation.Field(&inv.Payment,
			validation.When(
				!inv.Totals.Paid(),
				validation.By(validateUnpaidInvoicePayment),
				validation.Required,
				validation.Skip,
			),
			validation.Skip,
		),
		validation.Field(&inv.Totals,
			validation.By(validateInvoiceTotals),
			validation.Skip,
		),
	)
}

// validateParty validates rules common to both supplier and customer
func validateParty(value any) error {
	p, _ := value.(*org.Party)
	if p == nil {
		return nil
	}

	return validation.ValidateStruct(p,
		validation.Field(&p.TaxID,
			validation.Required,
			tax.RequireIdentityCode,
			validation.Skip,
		),
		validation.Field(&p.Addresses,
			validation.Each(
				validation.By(validateAddress),
				validation.Required,
				validation.Skip,
			),
			validation.Skip,
		),
		validation.Field(&p.Ext,
			validation.When(len(p.Addresses) > 0,
				tax.ExtensionsRequire(dfe.ExtKeyMunicipality),
			),
			validation.Skip,
		),
	)
}

func validateSupplier(value any) error {
	s, _ := value.(*org.Party)
	if s == nil {
		return nil
	}

	return validation.ValidateStruct(s,
		validation.Field(&s.Identities,
			org.RequireIdentityKey(dfe.IdentityKeyStateReg),
			validation.Skip,
		),
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.Addresses, validation.Required),
	)
}

func validateCustomer(inv *bill.Invoice) validation.RuleFunc {
	return func(value any) error {
		c, _ := value.(*org.Party)
		if c == nil {
			return nil
		}

		return validation.ValidateStruct(c,
			validation.Field(&c.Addresses,
				validation.When(isNFe(inv.Tax), validation.Required),
			),
		)
	}
}

func validateAddress(value any) error {
	a, _ := value.(*org.Address)
	if a == nil {
		return nil
	}

	return validation.ValidateStruct(a,
		validation.Field(&a.Street, validation.Required),
		validation.Field(&a.Number, validation.Required),
		validation.Field(&a.Locality, validation.Required),
		validation.Field(&a.State, validation.Required),
		validation.Field(&a.Code, validation.Required),
	)
}

func validateInvoiceTax(value any) error {
	t, _ := value.(*bill.Tax)
	if t == nil {
		return nil
	}

	return validation.ValidateStruct(t,
		validation.Field(&t.Ext,
			tax.ExtensionsRequire(ExtKeyModel, ExtKeyPresence),
			validation.When(
				isNFe(t),
				tax.ExtensionsExcludeCodes(ExtKeyPresence, PresenceDelivery),
			),
			validation.When(
				isNFCe(t),
				tax.ExtensionsHasCodes(ExtKeyPresence,
					PresenceInPerson, PresenceDelivery),
			),
			validation.Skip,
		),
	)
}

func validateInvoiceNotes(value any) error {
	nts, _ := value.([]*org.Note)

	for _, n := range nts {
		if n == nil {
			continue
		}

		if n.Key == org.NoteKeyReason {
			return nil
		}
	}

	return fmt.Errorf("note with key `%s` required. It must describe the nature of the operation (natOp)", org.NoteKeyReason)
}

func validateInvoiceNote(value any) error {
	n, _ := value.(*org.Note)
	if n == nil || n.Key != org.NoteKeyReason {
		return nil
	}

	return validation.ValidateStruct(n,
		validation.Field(&n.Text,
			validation.Length(1, 60),
			validation.Skip,
		),
	)
}

func validateUnpaidInvoicePayment(value any) error {
	p, _ := value.(*bill.PaymentDetails)
	if p == nil {
		return nil
	}

	return validation.ValidateStruct(p,
		validation.Field(&p.Instructions,
			validation.Required,
			validation.Skip,
		),
	)
}

func validateInvoiceTotals(value any) error {
	t, _ := value.(*bill.Totals)
	if t == nil {
		return nil
	}

	return validation.ValidateStruct(t,
		validation.Field(&t.Due,
			num.ZeroOrPositive,
			validation.Skip,
		),
	)
}

func isNFCe(t *bill.Tax) bool {
	return t != nil && t.Ext[ExtKeyModel] == ModelNFCe
}

func isNFe(t *bill.Tax) bool {
	return t != nil && t.Ext[ExtKeyModel] == ModelNFe
}
