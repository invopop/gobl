package mx

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Official SAT codes to include in stamps.
const (
	StampSATSignature   cbc.Key = "sat-sig"          // Signature - Sello Digital del SAT (optional)
	StampSATSerial      cbc.Key = "sat-serial"       // Cert Serial - Número de Certificado SAT
	StampSATTimestamp   cbc.Key = "sat-timestamp"    // Timestamp - Fecha y hora de certificación del SAT
	StampSATUUID        cbc.Key = "sat-uuid"         // Folio Fiscal
	StampSATURL         cbc.Key = "sat-url"          // URL QR Code
	StampSATProviderRFC cbc.Key = "sat-provider-rfc" // Provider RFC - RFC del Proveedor de Certificación
	StampSATChain       cbc.Key = "sat-chain"        // Cadena original del complemento de certificación digital del SAT
)

var stampUsage = &cbc.Usage{
	Codes: []cbc.Code{CountryCode},
}

var stamps = []*cbc.Definition{
	{
		Key: StampSATUUID,
		Name: i18n.String{
			i18n.EN: "Fiscal Folio",
			i18n.ES: "Folio Fiscal",
		},
		Desc: i18n.String{
			i18n.EN: "UUID assigned to the CFDI when it is certified.",
			i18n.ES: "UUID asignado al CFDI al ser timbrado.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampSATSignature,
		Name: i18n.String{
			i18n.EN: "SAT Signature",
			i18n.ES: "Sello Digital del SAT",
		},
		Desc: i18n.String{
			i18n.EN: "Digital signature of the SAT included in the certification of the CFDI.",
			i18n.ES: "Sello digital del SAT incluido en el timbrado del CFDI.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampSATSerial,
		Name: i18n.String{
			i18n.EN: "SAT Certificate Number",
			i18n.ES: "Número de Certificado del SAT",
		},
		Desc: i18n.String{
			i18n.EN: "Serial number of the SAT certificate used to certify the CFDI.",
			i18n.ES: "Número de serie del certificado del SAT usado para timbrar el CFDI.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampSATTimestamp,
		Name: i18n.String{
			i18n.EN: "Certification Date",
			i18n.ES: "Fecha y Hora de Certificación",
		},
		Desc: i18n.String{
			i18n.EN: "Date and time the CFDI was certified.",
			i18n.ES: "Fecha y hora en que se timbró el CFDI.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampSATProviderRFC,
		Name: i18n.String{
			i18n.EN: "Certification Provider RFC",
			i18n.ES: "RFC del Proveedor de Certificación",
		},
		Desc: i18n.String{
			i18n.EN: "RFC of the authorized certification provider (PAC) that certified the CFDI.",
			i18n.ES: "RFC del proveedor autorizado de certificación (PAC) que timbró el CFDI.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampSATChain,
		Name: i18n.String{
			i18n.EN: "Original Certification Chain",
			i18n.ES: "Cadena Original del Complemento de Certificación",
		},
		Desc: i18n.String{
			i18n.EN: "Original chain of the SAT digital certification complement.",
			i18n.ES: "Cadena original del complemento de certificación digital del SAT.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampSATURL,
		Name: i18n.String{
			i18n.EN: "QR Code URL",
			i18n.ES: "URL del Código QR",
		},
		Desc: i18n.String{
			i18n.EN: "URL to verify the CFDI with the SAT, included in the QR code.",
			i18n.ES: "URL para verificar el CFDI ante el SAT, incluida en el código QR.",
		},
		Usage: stampUsage,
	},
}
