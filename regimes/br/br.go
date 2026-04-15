// Package br provides the tax region definition for Brazil.
package br

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
)

// CountryCode is the ISO 3166-2 code for Brazil.
const CountryCode = "BR"

func init() {
	tax.RegisterRegimeDef(New())
	rules.Register("br", rules.GOBL.Add(CountryCode),
		orgPartyRules(),
		taxIdentityRules(),
	)
}

// New provides the tax region definition
func New() *tax.RegimeDef {
	return &tax.RegimeDef{
		Country:  CountryCode,
		Currency: currency.BRL,
		Name: i18n.String{
			i18n.EN: "Brazil",
			i18n.PT: "Brasil",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Brazil uses Notas Fiscais Eletrônicas (electronic fiscal notes) such as NFSe,
				NFe, and NFCe for reporting tax information to municipal, state, and federal
				authorities. The tax system is administered by the Receita Federal (Federal
				Revenue Service).

				Tax identification is provided through a CNPJ (Cadastro Nacional da Pessoa
				Jurídica) for businesses, consisting of 14 digits, or a CPF (Cadastro de Pessoas
				Físicas) for individuals, consisting of 11 digits. Both types are valid for the
				issuance of NFS-e (electronic service invoices).

				Brazilian addresses have three subdivisions relevant for tax purposes: bairro
				(neighbourhood), município (municipality), and estado (state). Municipality codes
				follow the IBGE coding system.

				Service notes (NFSe) let service providers document and report taxes such as
				ISS (Imposto Sobre Serviços) related to the services they provide. Municipal
				governments regulate them. Special tax regimes include Simples Nacional for
				simplified taxation of micro and small enterprises, and MEI (Micro-Empreendedor
				Individual) for individual micro-entrepreneurs.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.NewString("NFS-e Technical Documentation (ABRASF)"),
				URL:   "https://abrasf.org.br/biblioteca/arquivos-publicos/nfs-e",
			},
			{
				Title: i18n.NewString("IBGE Municipality Codes"),
				URL:   "https://www.ibge.gov.br/explica/codigos-dos-municipios.php",
			},
		},
		TimeZone:   "America/Sao_Paulo",
		Normalizer: Normalize,
		Extensions: extensions,
		Categories: taxCategories,
		Corrections: []*tax.CorrectionDefinition{
			{
				Schema: bill.ShortSchemaInvoice,
				Types: []cbc.Key{
					bill.InvoiceTypeCreditNote,
				},
			},
		},
	}
}

// Normalize will attempt to clean the object passed to it.
func Normalize(doc interface{}) {
	switch obj := doc.(type) {
	case *tax.Identity:
		tax.NormalizeIdentity(obj)
	}
}
