// Package dian provides the DIAN UBL 2.1 extensions used in Colombia.
package dian

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

const (
	// V2 for DIAN UBL 2.1 in Colombia
	V2 cbc.Key = "co-dian-v2"
)

// DIAN official codes to include in stamps.
const (
	StampCUDE cbc.Key = "dian-cude"
	StampQR   cbc.Key = "dian-qr"
)

func init() {
	tax.RegisterAddonDef(newAddon())
}

func newAddon() *tax.AddonDef {
	return &tax.AddonDef{
		Key: V2,
		Name: i18n.String{
			i18n.EN: "Colombia DIAN UBL 2.X",
		},
		Description: i18n.String{
			i18n.EN: here.Doc(`
				Extensions to support the Colombian DIAN (Dirección de Impuestos y Aduanas Nacionales)
				specifications for electronic invoicing based on UBL 2.1.

				## Invoice Series & Code

				Invoices ("Facturas de Venta") as expected in many countries require a series and unique sequential code to be issued for each document. The DIAN in Colombia however have gone a step further and require invoice series to be pre-registered with the government.

				Details on how to do this are available [here](https://www.dian.gov.co/impuestos/sociedades/presentacionclientes/Solicitud_de_Autorizacion_de_Numeracion_de_Facturacion.pdf).

				## Municipality codes

				The DIAN requires that Colombian addresses in the invoice specify its municipality code. The list of municipality codes is available [here](https://www.dian.gov.co/atencionciudadano/formulariosinstructivos/Formularios/2007/Codigos_municipios_2007.pdf).

				In a GOBL invoice, you can provide the supplier's or customer's municipality code using the ~co-dian-municipality~ extension. For example:

				~~~js
				"supplier": {
					"name": "EXAMPLE SUPPLIER S.A.S.",
					"tax_id": {
						"country": "CO",
						"code": "9014514812"
					},
					"ext": {
						"co-dian-municipality": "11001" // Bogotá, D.C.
					},
					// [...]
				},
				"customer": {
					"name": "EXAMPLE CUSTOMER S.A.S.",
					"tax_id": {
						"country": "CO",
						"code": "9014514805"
					},
					"ext": {
						"co-dian-municipality": "05001" // Medellín
					},
					// [...]
				},
				~~~

				## Customer identities

				While the DIAN requires suppliers of invoices to identify themselves using their NIT, in the case of customers, it allows various identification types. Each identification type has a specific code that must be sent to the DIAN.

				### B2B

				In a GOBL invoice, Colombian business customers are required to provide a Tax ID with the company's NIT to be reported with the ~31~ (NIT) DIAN ID type. Foreign business customers must also provide a Tax ID with their country's VAT code, and that will be reported to the DIAN using the ~50~ (NIT de otro país) type.

				For example:

				~~~js
				"customer": {
					"name": "EXAMPLE CUSTOMER S.A.S.",
					"tax_id": {
						"country": "CO",
						"code": "9014514805" // NIT. DIAN type 31
					},
					"ext": {
						"co-dian-municipality": "11001"
					},
					"addresses": [
						{
						"street": "CRA 8 113 31 OF 703",
						"locality": "Bogotá, D.C.",
						"region": "Bogotá",
						"country": "CO"
						}
					]
				}
				~~~

				### B2C

				In the case of non-business customers, the GOBL invoice will need to include the tax tag ~simplified~. That will allow to omit the customer or identify it with any of the other document types accepted by the DIAN. They'll just need to include an Identity object with any of the keys below:

				| GOBL Identity Key      | DIAN ID type | Description                            |
				| ---------------------- | ------------ | -------------------------------------- |
				| ~co-civil-register~    | ~11~         | Registro civil                         |
				| ~co-id-card~           | ~12~         | Tarjeta de identidad                   |
				| ~co-citizen-id~        | ~13~         | Cédula de ciudadanía                   |
				| ~co-foreigner-id-card~ | ~21~         | Tarjeta de extranjería                 |
				| ~co-foreigner-id~      | ~22~         | Cédula de extranjería                  |
				| ~co-passport~          | ~41~         | Pasaporte                              |
				| ~co-foreign-id~        | ~42~         | Documento de identificación extranjero |
				| ~co-pep~               | ~47~         | PEP (Permiso Especial de Permanencia)  |
				| ~co-nuip~              | ~91~         | NUIP                                   |

				For example:

				~~~js
				"tax": {
					"tags": ["simplified"]
				},
				"customer": {
					"name": "Ángel Pérez",
					"identities": [
						{
						"key": "co-passport", // DIAN type 41
						"code": "1234567890"
						}
					]
				}
				~~~

				Note that GOBL ~simplified~ invoices don't require a customer (or its identity) to be present. In the lack of a customer identity, the reserved code for final consumers (~222222222222~) will be automatically sent to the DIAN (i.e., no need to set it explicitly)

				For example:

				~~~js
				"tax": {
					"tags": ["simplified"]
				},
				"customer": { // The customer could be fully omitted as well
					"name": "Juan Fernández"
				}
				~~~

				## Credit and debit correction codes

				The DIAN requires credit and debit notes sent to them to specify a code with a cause of the correction.

				In a GOBL invoice, you'll need to include the extension ~co-dian-credit-code~ (for credit notes) or ~co-dian-debit-code~ (for debit notes) as part of the Preceding section of the document. Each extension allows a different set of values:

				**~co-dian-credit-code~**

				| Code | Description    |
				| ---- | -------------- |
				| 1    | Partial refund |
				| 2    | Revoked        |
				| 3    | Discount       |
				| 4    | Adjustment     |
				| 5    | Other          |

				**~co-dian-debit-code~**

				| Code | Description     |
				| ---- | --------------- |
				| 1    | Interest        |
				| 2    | Pending charges |
				| 3    | Change in value |
				| 4    | Other           |

				For example:

				~~~js
				"preceding": [
					{
						"uuid": "0190e063-7676-7000-8c58-2db7172a4e58",
						"type": "standard",
						"series": "SETT",
						"code": "1010006",
						"issue_date": "2024-07-23",
						"reason": "Reason",
						"stamps": [
						{
							"prv": "dian-cude",
							"val": "57601dd1ab69213ccf8cfd5894f2e9fbfe23643f3a24e2f2526a5bb88d058a0842fffcb339694b6704dc105a9d813327"
						}
						],
						"ext": {
						"co-dian-debit-code": "3"
						}
					}
				],
				~~~
			`),
		},
		Extensions:  extensions,
		Identities:  identities,
		Normalizer:  normalize,
		Validator:   validate,
		Corrections: invoiceCorrectionDefinitions,
	}
}

func normalize(doc any) {
	switch obj := doc.(type) {
	case *bill.Invoice:
		normalizeInvoice(obj)
	}
}

func validate(doc any) error {
	switch obj := doc.(type) {
	case *bill.Invoice:
		return validateInvoice(obj)
	}
	return nil
}
