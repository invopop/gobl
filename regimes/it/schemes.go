package it

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Base scheme keys for this regime.
const (
	SchemeKeyOrdinary  cbc.Key = "ordinary"
	SchemeKeyFreelance cbc.Key = "freelance"
)

var schemes = []*tax.Scheme{
	{
		Key: SchemeKeyOrdinary,
		Name: i18n.String{
			i18n.EN: "Ordinary",
			i18n.IT: "Regime Ordinario",
		},
		InvoiceTypes: []cbc.Key{
			bill.InvoiceTypeDefault.Key(),
		},
		Meta: cbc.Meta{
			KeyFatturaPATipoDocumento: "TD01",
			KeyFatturaPARegimeFiscale: "RF01",
		},
	},
	{
		Key: SchemeKeyOrdinary,
		Name: i18n.String{
			i18n.EN: "Partial invoice",
			i18n.IT: "Anticipo su fattura",
		},
		InvoiceTypes: []cbc.Key{
			bill.InvoiceTypePartial.Key(),
		},
		Meta: cbc.Meta{
			KeyFatturaPATipoDocumento: "TD02",
		},
	},
	{
		Key: SchemeKeyFreelance,
		Name: i18n.String{
			i18n.EN: "Partial for freelancer",
			i18n.IT: "Anticipo su parcella",
		},
		InvoiceTypes: []cbc.Key{
			bill.InvoiceTypePartial.Key(),
		},
		Meta: cbc.Meta{
			KeyFatturaPATipoDocumento: "TD03",
		},
	},
	{
		Key: SchemeKeyOrdinary,
		Name: i18n.String{
			i18n.EN: "Credit Note",
		},
		InvoiceTypes: []cbc.Key{
			bill.InvoiceTypeCreditNote.Key(),
		},
		Meta: cbc.Meta{
			KeyFatturaPATipoDocumento: "TD04",
		},
	},
	{
		Key: SchemeKeyOrdinary,
		Name: i18n.String{
			i18n.EN: "Debit Note",
		},
		InvoiceTypes: []cbc.Key{
			bill.InvoiceTypeDebitNote.Key(),
		},
		Meta: cbc.Meta{
			KeyFatturaPATipoDocumento: "TD05",
		},
	},
	{
		Key: SchemeKeyFreelance,
		Name: i18n.String{
			i18n.EN: "Freelancer",
			i18n.IT: "Parcella",
		},
		Meta: cbc.Meta{
			KeyFatturaPATipoDocumento: "TD06",
		},
	},
}
