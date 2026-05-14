package tbai

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Extension keys for TicketBAI
const (
	ExtKeyRegion       cbc.Key = "es-tbai-region"
	ExtKeyExempt       cbc.Key = "es-tbai-exemption"
	ExtKeyProduct      cbc.Key = "es-tbai-product"
	ExtKeyCorrection   cbc.Key = "es-tbai-correction"
	ExtKeyBIActivity   cbc.Key = "es-tbai-bi-activity"
	ExtKeyRegime       cbc.Key = "es-tbai-regime"
	ExtKeyIdentityType cbc.Key = "es-tbai-identity-type"
)

// Identity Type Codes - subset of the L7 list assigned to identities.
const (
	ExtCodeIdentityTypeVAT      cbc.Code = "02" // NIF-VAT identity (VIES)
	ExtCodeIdentityTypePassport cbc.Code = "03" // Passport
	ExtCodeIdentityTypeForeign  cbc.Code = "04" // Foreign Identity Document
	ExtCodeIdentityTypeResident cbc.Code = "05" // Residential Certificate
	ExtCodeIdentityTypeOther    cbc.Code = "06" // Other Identity Document
)

// Extension values for product key.
const (
	ExtValueProductGoods    cbc.Code = "goods"
	ExtValueProductServices cbc.Code = "services"
	ExtValueProductResale   cbc.Code = "resale"
)

// Extension values for region key.
const (
	ExtValueRegionVI cbc.Code = "VI" // Araba
	ExtValueRegionBI cbc.Code = "BI" // Bizkaia
	ExtValueRegionSS cbc.Code = "SS" // Gipuzkoa
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyRegion,
		Name: i18n.String{
			i18n.EN: "TicketBAI Region Code",
			i18n.ES: "Código de Región TicketBAI",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Region codes are used by TicketBAI to differentiate between the different
				subdivisions of the Basque Country. This is used to determine the correct
				API endpoint to use when submitting documents.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: ExtValueRegionVI,
				Name: i18n.String{
					i18n.EN: "Araba",
					i18n.ES: "Álava",
				},
			},
			{
				Code: ExtValueRegionBI,
				Name: i18n.String{
					i18n.EN: "Bizkaia",
					i18n.ES: "Vizcaya",
				},
			},
			{
				Code: ExtValueRegionSS,
				Name: i18n.String{
					i18n.EN: "Gipuzkoa",
					i18n.ES: "Guipúzcoa",
				},
			},
		},
	},
	{
		Key: ExtKeyProduct,
		Name: i18n.String{
			i18n.EN: "TicketBAI Product Key",
			i18n.ES: "Clave de Producto TicketBAI",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Product keys are used by TicketBAI to differentiate between -exported- goods
				and services. It may be useful to classify all products regardless of wether
				they are exported or not.

				There is an additional exception case for goods that are resold without modification
				when the supplier is in the simplified tax regime. For must purposes this special
				case can be ignored.

				If no product key is provided, the default is "services".
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: ExtValueProductGoods,
				Name: i18n.String{
					i18n.EN: "Delivery of goods",
					i18n.ES: "Entrega de bienes",
				},
			},
			{
				Code: ExtValueProductServices,
				Name: i18n.String{
					i18n.EN: "Provision of services",
					i18n.ES: "Prestacion de servicios",
				},
			},
			{
				Code: ExtValueProductResale,
				Name: i18n.String{
					i18n.EN: "Resale of goods without modification by vendor in the simplified regime",
					i18n.ES: "Reventa de bienes sin modificación por vendedor en regimen simplificado",
				},
			},
		},
	},
	{
		Key: ExtKeyExempt,
		Name: i18n.String{
			i18n.EN: "TicketBAI Exemption code",
			i18n.ES: "Código de Exención de TicketBAI",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Codes used by TicketBAI for both "exempt", "not-subject", and reverse
				charge transactions. In the TicketBAI format these are separated,
				but in order to simplify GOBL and be more closely aligned with
				other countries we've combined them into one.

				The follow mappings will be made automatically by GOBL during normalization.

				| Tax Key           | Exemption Codes            |
				|-------------------|----------------------------|
				| ~exempt~          | ~E1~ (default), ~E6~       |
				| ~export~          | ~E2~ (default), ~E3~, ~E4~ |
				| ~intra-community~ | ~E5~                       |
				| ~reverse-charge~  | ~S2~                       |
				| ~outside-scope~   | ~OT~, ~RL~, ~VT~, ~IE~     |
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "E1",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 20 of the Foral VAT Law",
					i18n.ES: "Exenta: por el artículo 20 de la Norma Foral del IVA",
				},
			},
			{
				Code: "E2",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 21 of the Foral VAT Law",
					i18n.ES: "Exenta: por el artículo 21 de la Norma Foral del IVA",
				},
			},
			{
				Code: "E3",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 22 of the Foral VAT Law",
					i18n.ES: "Exenta: por el artículo 22 de la Norma Foral del IVA",
				},
			},
			{
				Code: "E4",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Articles 23 and 24 of the Foral VAT Law",
					i18n.ES: "Exenta: por el artículos 23 y 24 de la Norma Foral del IVA",
				},
			},
			{
				Code: "E5",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 25 of the Foral VAT law",
					i18n.ES: "Exenta: por el artículo 25 de la Norma Foral del IVA",
				},
			},
			{
				Code: "E6",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to other reasons",
					i18n.ES: "Exenta: por otra causa",
				},
			},
			{
				Code: "OT",
				Name: i18n.String{
					i18n.EN: "Not subject: pursuant to Article 7 of the VAT Law - other cases of non-subject",
					i18n.ES: "No sujeto: por el artículo 7 de la Ley del IVA - otros supuestos de no sujeción",
				},
			},
			{
				Code: "RL",
				Name: i18n.String{
					i18n.EN: "Not subject: pursuant to localization rules",
					i18n.ES: "No sujeto: por reglas de localización",
				},
			},
			{
				Code: "VT",
				Name: i18n.String{
					i18n.EN: "Not subject: sales made on behalf of third parties (amount not computable for VAT or IRPF purposes)",
					i18n.ES: "No sujeto: ventas realizadas por cuenta de terceros (importe no computable a efectos de IVA ni de IRPF)",
				},
			},
			{
				Code: "IE",
				Name: i18n.String{
					i18n.EN: "Not subject in the TAI due to localization rules, but foreign tax, IPS/IGIC or VAT from another EU member state is passed on",
					i18n.ES: "No sujeto en el TAI por reglas de localización, pero repercute impuesto extranjero, IPS/IGIC o IVA de otro estado miembro UE",
				},
			},
			/*
				// S1 is the default value for regular invoices, so we don't need to include it here
				// alongside the exemption codes.
				{
					Code: "S1",
					Name: i18n.String{
						i18n.EN: "Subject and not exempt: without reverse charge",
						i18n.ES: "Sujeto y no exenta: sin inversión del sujeto pasivo",
					},
				},
			*/
			{
				Code: "S2",
				Name: i18n.String{
					i18n.EN: "Subject and not exempt: with reverse charge",
					i18n.ES: "Sujeto y no exenta: con inversión del sujeto pasivo",
				},
			},
		},
	},
	{
		Key: ExtKeyCorrection,
		Name: i18n.String{
			i18n.EN: "TicketBAI Rectification Type Code",
			i18n.ES: "TicketBAI Código de Factura Rectificativa",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Corrected or rectified invoices that need to be sent in the TicketBAI format
				require a specific type code to be defined alongside the preceding invoice
				data.
			`),
		},
		// Codes taken from TicketBAI XSD
		Values: []*cbc.Definition{
			{
				Code: "R1",
				Name: i18n.String{
					i18n.EN: "Rectified invoice: error based on law and Article 80 One, Two and Six of the Provincial Tax Law of VAT",
					i18n.ES: "Factura rectificativa: error fundado en derecho y Art. 80 Uno, Dos y Seis de la Norma Foral del IVA",
					i18n.EU: "Faktura zuzentzailea: zuzenbidean oinarritutako akatsa eta BEZaren Foru Arauaren 80.artikuluko Bat, Bi eta Sei",
				},
			},
			{
				Code: "R2",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: artículo 80 Tres de la Norma Foral del IVA",
					i18n.EN: "Rectified invoice: error based on law and Article 80 Three of the Provincial Tax Law of VAT",
					i18n.EU: "Faktura zuzentzailea: BEZari buruzko Foru Arauko 80. artikuluko Hiru",
				},
			},
			{
				Code: "R3",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: artículo 80 Cuatro de la Norma Foral del IVA",
					i18n.EN: "Rectified invoice: error based on law and Article 80 Four of the Provincial Tax Law of VAT",
					i18n.EU: "Faktura zuzentzailea: BEZari buruzko Foru Arauko 80. artikuluko Lau",
				},
			},
			{
				Code: "R4",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: Resto",
					i18n.EN: "Rectified invoice: Other",
					i18n.EU: "Faktura zuzentzailea: Gainerakoak",
				},
			},
			{
				Code: "R5",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: facturas simplificadas",
					i18n.EN: "Rectified invoice: simplified invoices",
					i18n.EU: "Faktura zuzentzaile: faktura erraztuetan",
				},
			},
		},
	},
	{
		Key: ExtKeyRegime,
		Name: i18n.String{
			i18n.EN: "TicketBAI VAT Regime Key",
			i18n.ES: "Clave de Régimen de IVA TicketBAI",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Identifies the VAT regime or operation classification applied to the
				transaction. Maps to the ~ClaveRegimenIvaOpTrascendencia~ field, with
				values taken from the TicketBAI XSD.

				The regime code is normally assigned per tax combo. If no regime code
				is provided, GOBL will try to assign a code from the following contexts:

				| Combo Context              | Regime Code |
				|----------------------------|-------------|
				| Key ~export~               | ~02~        |
				| Has surcharge              | ~51~        |
				| Invoice tag ~simplified-scheme~ | ~52~   |
				| Otherwise                  | ~01~        |
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "01",
				Name: i18n.String{
					i18n.EN: "General regime operation",
					i18n.ES: "Operación de régimen general",
				},
			},
			{
				Code: "02",
				Name: i18n.String{
					i18n.EN: "Export",
					i18n.ES: "Exportación",
				},
			},
			{
				Code: "03",
				Name: i18n.String{
					i18n.EN: "Special regime for used goods, art objects, antiques and collectibles",
					i18n.ES: "Operaciones a las que se aplique el régimen especial de bienes usados, objetos de arte, antigüedades y objetos de colección",
				},
			},
			{
				Code: "04",
				Name: i18n.String{
					i18n.EN: "Special regime for investment gold",
					i18n.ES: "Régimen especial del oro de inversión",
				},
			},
			{
				Code: "05",
				Name: i18n.String{
					i18n.EN: "Special regime for travel agencies",
					i18n.ES: "Régimen especial de las agencias de viajes",
				},
			},
			{
				Code: "06",
				Name: i18n.String{
					i18n.EN: "Special regime for VAT groups (Advanced Level)",
					i18n.ES: "Régimen especial grupo de entidades en IVA (Nivel Avanzado)",
				},
			},
			{
				Code: "07",
				Name: i18n.String{
					i18n.EN: "Special cash accounting regime",
					i18n.ES: "Régimen especial del criterio de caja",
				},
			},
			{
				Code: "08",
				Name: i18n.String{
					i18n.EN: "Operations subject to IPSI/IGIC",
					i18n.ES: "Operaciones sujetas al IPSI/IGIC",
				},
			},
			{
				Code: "09",
				Name: i18n.String{
					i18n.EN: "Billing of travel agency services acting as mediators in name and on behalf of others",
					i18n.ES: "Facturación de las prestaciones de servicios de agencias de viaje que actúan como mediadoras en nombre y por cuenta ajena",
				},
			},
			{
				Code: "10",
				Name: i18n.String{
					i18n.EN: "Collection of professional fees or industrial property rights on behalf of third parties",
					i18n.ES: "Cobros por cuenta de terceros de honorarios profesionales o de derechos derivados de la propiedad industrial",
				},
			},
			{
				Code: "11",
				Name: i18n.String{
					i18n.EN: "Business premises rental operations subject to withholding",
					i18n.ES: "Operaciones de arrendamiento de local de negocio sujetos a retención",
				},
			},
			{
				Code: "12",
				Name: i18n.String{
					i18n.EN: "Business premises rental operations not subject to withholding",
					i18n.ES: "Operaciones de arrendamiento de local de negocio no sujetos a retención",
				},
			},
			{
				Code: "13",
				Name: i18n.String{
					i18n.EN: "Business premises rental operations, both subject and not subject to withholding",
					i18n.ES: "Operaciones de arrendamiento de local de negocio sujetas y no sujetas a retención",
				},
			},
			{
				Code: "14",
				Name: i18n.String{
					i18n.EN: "Invoice with pending VAT accrual in work certifications for Public Administration",
					i18n.ES: "Factura con IVA pendiente de devengo en certificaciones de obra cuyo destinatario sea una Administración Pública",
				},
			},
			{
				Code: "15",
				Name: i18n.String{
					i18n.EN: "Invoice with pending VAT accrual in successive tract operations",
					i18n.ES: "Factura con IVA pendiente de devengo en operaciones de tracto sucesivo",
				},
			},
			{
				Code: "17",
				Name: i18n.String{
					i18n.EN: "Operation under OSS or IOSS regimes",
					i18n.ES: "Operación acogida a alguno de los regímenes previstos en el Capítulo XI del Título IX (OSS e IOSS)",
				},
			},
			{
				Code: "51",
				Name: i18n.String{
					i18n.EN: "Operations under the equivalence surcharge regime",
					i18n.ES: "Operaciones en recargo de equivalencia",
				},
			},
			{
				Code: "52",
				Name: i18n.String{
					i18n.EN: "Operations under the simplified VAT regime",
					i18n.ES: "Operaciones en régimen simplificado",
				},
			},
			{
				Code: "53",
				Name: i18n.String{
					i18n.EN: "Operations carried out by or for entities without a permanent establishment",
					i18n.ES: "Operaciones realizadas por o para entidades sin establecimiento permanente",
				},
			},
		},
	},
	{
		Key: ExtKeyBIActivity,
		Name: i18n.String{
			i18n.EN: "Activity Code (Bizkaia)",
			i18n.ES: "Código de Actividad (Bizkaia)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Economic activity code (epígrafe) for individual issuers submitting through
				Bizkaia's LROE Modelo 140 register. Not required for organisations, who
				file through Modelo 240.
			`),
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Batuz LROE list of activity codes",
					i18n.ES: "Lista de epígrafes LROE Batuz",
				},
				URL:         "https://www.batuz.eus/fitxategiak/batuz/lroe/batuz_lroe_lista_epigrafes_v1_0_4.xlsx",
				ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			},
			{
				Title: i18n.String{
					i18n.EN: "LROE Modelo 140 specification",
					i18n.ES: "Especificación LROE Modelo 140",
				},
				URL:         "https://www.batuz.eus/fitxategiak/batuz/lroe/lroe_140_v_1_0.pdf",
				ContentType: "application/pdf",
			},
		},
		Pattern: `^\d{1,7}$`,
	},
	{
		Key: ExtKeyIdentityType,
		Name: i18n.String{
			i18n.EN: "Identity Type Code",
			i18n.ES: "Código de Tipo de Identidad",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Identity code used to identify the type of identity document used by the customer
				when there is no tax identifier available. Maps to the ~IDType~ value under
				~IDOtro~ in the TicketBAI XML.

				The regular Party Tax Identity is preferred over using a specific identity type
				code, and will be mapped automatically as follows:

				- Spanish Tax IDs will be mapped to the ~NIF~ field.
				- EU Tax IDs will be mapped to the ~IDOtro~ field with code ~02~.
				- Non-EU Tax IDs will be mapped to the ~IDOtro~ field with code ~04~.

				The following identity ~key~ values will be mapped automatically to an extension
				by the addon:

				| Identity Key | Extension Code |
				|--------------|----------------|
				| ~passport~   | ~03~           |
				| ~foreign~    | ~04~           |
				| ~resident~   | ~05~           |
				| ~other~      | ~06~           |

				Example identity of a UK passport:

				~~~
				{
					"identities": [
						{
							"key": "passport",
							"country": "GB",
							"code": "123456789"
						}
					]
				}
				~~~

				Will be normalized to:

				~~~
				{
					"identities": [
						{
							"key": "passport",
							"country": "GB",
							"code": "123456789",
							"ext": {
								"es-tbai-identity-type": "03"
							}
						}
					]
				}
				~~~
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: ExtCodeIdentityTypeVAT,
				Name: i18n.String{
					i18n.EN: "NIF-VAT Identity (VIES)",
					i18n.ES: "NIF-VAT (VIES)",
				},
			},
			{
				Code: ExtCodeIdentityTypePassport,
				Name: i18n.String{
					i18n.EN: "Passport",
					i18n.ES: "Pasaporte",
				},
			},
			{
				Code: ExtCodeIdentityTypeForeign,
				Name: i18n.String{
					i18n.EN: "Foreign Identity Document",
					i18n.ES: "Documento de Identidad Extranjero",
				},
			},
			{
				Code: ExtCodeIdentityTypeResident,
				Name: i18n.String{
					i18n.EN: "Residential Certificate",
					i18n.ES: "Certificado Residencia",
				},
			},
			{
				Code: ExtCodeIdentityTypeOther,
				Name: i18n.String{
					i18n.EN: "Other Identity Document",
					i18n.ES: "Otro Documento Probatorio",
				},
			},
		},
	},
}
