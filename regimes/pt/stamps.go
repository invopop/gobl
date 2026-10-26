package pt

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// AT official codes to include in stamps.
const (
	StampProviderATATCUD     cbc.Key = "at-atcud"
	StampProviderATQR        cbc.Key = "at-qr"
	StampProviderATHash      cbc.Key = "at-hash"
	StampProviderATHashFull  cbc.Key = "at-hash-full"
	StampProviderATAppID     cbc.Key = "at-app-id"
	StampProviderATTimestamp cbc.Key = "at-ts"
	StampProviderATDocCode   cbc.Key = "at-doc-code"
)

var stampUsage = &cbc.Usage{
	Codes: []cbc.Code{CountryCode},
}

var stamps = []*cbc.Definition{
	{
		Key: StampProviderATATCUD,
		Name: i18n.String{
			i18n.EN: "ATCUD",
			i18n.PT: "ATCUD",
		},
		Desc: i18n.String{
			i18n.EN: "Unique document code formed by the series validation code assigned by the AT and the document's sequential number.",
			i18n.PT: "Código único do documento formado pelo código de validação da série atribuído pela AT e o número sequencial do documento.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampProviderATQR,
		Name: i18n.String{
			i18n.EN: "QR Code",
			i18n.PT: "Código QR",
		},
		Desc: i18n.String{
			i18n.EN: "Contents of the QR code that must be printed on the document.",
			i18n.PT: "Conteúdo do código QR que deve ser impresso no documento.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampProviderATHash,
		Name: i18n.String{
			i18n.EN: "Hash Characters",
			i18n.PT: "Caracteres do Hash",
		},
		Desc: i18n.String{
			i18n.EN: "Characters extracted from the document's signature that must be printed on the document.",
			i18n.PT: "Caracteres extraídos da assinatura do documento que devem ser impressos no documento.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampProviderATHashFull,
		Name: i18n.String{
			i18n.EN: "Hash",
			i18n.PT: "Hash",
		},
		Desc: i18n.String{
			i18n.EN: "Complete signature of the document, chained to the signature of the previous document in the series.",
			i18n.PT: "Assinatura completa do documento, encadeada com a assinatura do documento anterior da série.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampProviderATAppID,
		Name: i18n.String{
			i18n.EN: "Software Certificate Number",
			i18n.PT: "Número de Certificado do Programa",
		},
		Desc: i18n.String{
			i18n.EN: "Number of the certificate assigned by the AT to the invoicing software that issued the document.",
			i18n.PT: "Número do certificado atribuído pela AT ao programa de faturação que emitiu o documento.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampProviderATTimestamp,
		Name: i18n.String{
			i18n.EN: "System Entry Date",
			i18n.PT: "Data de Gravação",
		},
		Desc: i18n.String{
			i18n.EN: "Date and time the document was signed and recorded by the invoicing software.",
			i18n.PT: "Data e hora em que o documento foi assinado e gravado pelo programa de faturação.",
		},
		Usage: stampUsage,
	},
	{
		Key: StampProviderATDocCode,
		Name: i18n.String{
			i18n.EN: "AT Document Code",
			i18n.PT: "Código de Identificação AT",
		},
		Desc: i18n.String{
			i18n.EN: "Identification code returned by the AT when a transport document is communicated.",
			i18n.PT: "Código de identificação devolvido pela AT na comunicação de um documento de transporte.",
		},
		Usage: stampUsage,
	},
}
