package mx

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

var stamps = []*cbc.Definition{
	{
		Key: StampSATUUID,
		Name: i18n.String{
			i18n.EN: "Fiscal Folio",
			i18n.ES: "Folio Fiscal",
		},
		Desc: i18n.String{
			i18n.EN: "UUID assigned to the CFDI when it is certified by the SAT. Must be copied to preceding document references when correcting.",
		},
	},
	{
		Key: StampSATSignature,
		Name: i18n.String{
			i18n.EN: "SAT Digital Seal",
			i18n.ES: "Sello Digital del SAT",
		},
		Desc: i18n.String{
			i18n.EN: "Digital seal applied by the SAT when certifying the CFDI.",
		},
	},
	{
		Key: StampSATSerial,
		Name: i18n.String{
			i18n.EN: "SAT Certificate Serial",
			i18n.ES: "Número de Certificado SAT",
		},
		Desc: i18n.String{
			i18n.EN: "Serial number of the SAT certificate used to certify the CFDI.",
		},
	},
	{
		Key: StampSATTimestamp,
		Name: i18n.String{
			i18n.EN: "SAT Certification Timestamp",
			i18n.ES: "Fecha y hora de certificación del SAT",
		},
		Desc: i18n.String{
			i18n.EN: "Date and time at which the CFDI was certified.",
		},
	},
	{
		Key: StampSATURL,
		Name: i18n.String{
			i18n.EN: "SAT Verification URL",
		},
		Desc: i18n.String{
			i18n.EN: "URL to include in the QR code used to verify the CFDI.",
		},
	},
	{
		Key: StampSATProviderRFC,
		Name: i18n.String{
			i18n.EN: "Certification Provider RFC",
			i18n.ES: "RFC del Proveedor de Certificación",
		},
		Desc: i18n.String{
			i18n.EN: "RFC of the authorised certification provider (PAC) that certified the CFDI.",
		},
	},
	{
		Key: StampSATChain,
		Name: i18n.String{
			i18n.EN: "SAT Certification Original Chain",
			i18n.ES: "Cadena original del complemento de certificación digital del SAT",
		},
		Desc: i18n.String{
			i18n.EN: "Original chain of the SAT digital certification complement.",
		},
	},
}
