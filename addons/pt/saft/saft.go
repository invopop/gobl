// Package saft provides the SAF-T addon for Portuguese invoices.
package saft

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

const (
	// V1 for SAF-T (PT) versions 1.x
	V1 cbc.Key = "pt-saft-v1"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V1,
		Name: i18n.String{
			i18n.EN: "Portugal SAF-T",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Portugal doesn't have an e-invoicing format per se. Tax information is reported
				electronically to the AT (Autoridade Tributária e Aduaneira) either periodically in
				batches via a SAF-T (PT) report or individually in real time via a web service. This addon
				ensures that the GOBL documents have all the required fields to be able to be reported to
				the AT.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title:       i18n.NewString("Portaria n.o 302/2016 – SAF-T Data Structure & Taxonomies"),
				URL:         "https://info.portaldasfinancas.gov.pt/pt/informacao_fiscal/legislacao/diplomas_legislativos/Documents/Portaria_302_2016.pdf",
				ContentType: "application/pdf",
			},
			{
				Title:       i18n.NewString("Portaria n.o 195/2020 – Comunicação de Séries Documentais, Aspetos Específicos"),
				URL:         "https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/Faturacao/Comunicacao_Series_ATCUD/Documents/Comunicacao_de_Series_Documentais_Manual_de_Integracao_de_SW_Aspetos_Genericos.pdf",
				ContentType: "application/pdf",
			},
			{
				Title:       i18n.NewString("Portaria n.o 195/2020 – Especificações Técnicas Código QR"),
				URL:         "https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/Novas_regras_faturacao/Documents/Especificacoes_Tecnicas_Codigo_QR.pdf",
				ContentType: "application/pdf",
			},
			{
				Title:       i18n.NewString("Comunicação dos elementos dos documentos de faturação à AT, por webservice"),
				URL:         "https://info.portaldasfinancas.gov.pt/pt/apoio_contribuinte/Faturacao/Fatcorews/Documents/Comunicacao_dos_elementos_dos_documentos_de_faturacao.pdf",
				ContentType: "application/pdf",
			},
		},
		Extensions: extensions,
		Normalizer: normalize,
		Scenarios:  scenarios,
		Validator:  validate,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *tax.Combo:
		normalizeTaxCombo(obj)
	case *org.Item:
		normalizeItem(obj)
	case *pay.Instructions:
		normalizePayInstructions(obj)
	case *pay.Advance:
		normalizePayAdvance(obj)
	case *bill.Payment:
		normalizePayment(obj)
	case *bill.Delivery:
		normalizeDelivery(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	case *bill.Payment:
		return validatePayment(obj)
	case *bill.Delivery:
		return validateDelivery(obj)
	case *tax.Combo:
		return validateTaxCombo(obj)
	case *org.Item:
		return validateItem(obj)
	}
	return nil
}
