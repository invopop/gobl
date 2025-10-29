package br

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Regime extension keys
const (
	ExtKeyMunicipality = "br-ibge-municipality"
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyMunicipality,
		Name: i18n.String{
			i18n.EN: "IBGE Municipality Code",
			i18n.PT: "Código do Município do IBGE",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The municipality code as defined by the IBGE (Brazilian Institute of Geography and
				Statistics).
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "IBGE - Municipalities Codes",
					i18n.PT: "IBGE - Códigos dos Municípios",
				},
				URL: "https://www.ibge.gov.br/explica/codigos-dos-municipios.php",
			},
		},
		Pattern: `^\d{7}$`,
	},
}
