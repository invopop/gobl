package verifactu

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Extension keys for Verifactu
const (
	ExtKeyDocType           cbc.Key = "es-verifactu-doc-type"
	ExtKeyOpClass           cbc.Key = "es-verifactu-op-class"
	ExtKeyCorrectionType    cbc.Key = "es-verifactu-correction-type"
	ExtKeyExempt            cbc.Key = "es-verifactu-exempt"
	ExtKeyRegime            cbc.Key = "es-verifactu-regime"
	ExtKeyIdentityType      cbc.Key = "es-verifactu-identity-type"
	ExtKeySimplifiedArt7273 cbc.Key = "es-verifactu-simplified-art7273"
	ExtKeyIssuerType        cbc.Key = "es-verifactu-issuer-type"
)

// Identity Type Codes - limited subset assigned to identities.
const (
	ExtCodeIdentityTypeVAT      cbc.Code = "02" // NIF-VAT identity
	ExtCodeIdentityTypePassport cbc.Code = "03" // Passport
	ExtCodeIdentityTypeForeign  cbc.Code = "04" // Foreign Identity Document
	ExtCodeIdentityTypeResident cbc.Code = "05" // Spanish Resident Foreigner Identity Card
	ExtCodeIdentityTypeOther    cbc.Code = "06" // Other Identity Document
)

// Issuer Type Codes
const (
	ExtCodeIssuerTypeThirdParty cbc.Code = "T" // Issued by Third Party
	ExtCodeIssuerTypeCustomer   cbc.Code = "D" // Issued by Customer
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyDocType,
		Name: i18n.String{
			i18n.EN: "Invoice Type Code",
			i18n.ES: "Código de Tipo de Factura",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Verifactu Ministerial Order",
					i18n.ES: "Orden Ministerial de Verifactu",
				},
				URL: "https://www.boe.es/diario_boe/txt.php?id=BOE-A-2024-22138",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code used to identify the type of invoice being sent. This will be
				determined automatically by GOBL during normalization according
				to the scenario definitions.

				The codes ~R2~, ~R3~, and ~R4~ are not covered by GOBL's scenarios
				and will need to be set manually if needed.

				Values correspond to L2 list.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "F1",
				Name: i18n.String{
					i18n.EN: "Invoice (Article 6, 7.2 and 7.3 of RD 1619/2012)",
					i18n.ES: "Factura (Art. 6, 7.2 y 7.3 del RD 1619/2012)",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						For regular invoices.
					`),
				},
			},
			{
				Code: "F2",
				Name: i18n.String{
					i18n.EN: "Simplified invoice (Article 6.1.d) of RD 1619/2012)",
					i18n.ES: "Factura Simplificada (Art. 6.1.d) del RD 1619/2012)",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						To use for B2C invoices where details about the customer are not
						normally required.
					`),
				},
			},
			{
				Code: "F3",
				Name: i18n.String{
					i18n.EN: "Invoice issued as a replacement for simplified invoices that have been billed and declared.",
					i18n.ES: "Factura emitida en sustitución de facturas simplificadas facturadas y declaradas.",
				},
			},
			{
				Code: "R1",
				Name: i18n.String{
					i18n.EN: "Rectified invoice: error based on law and Article 80 One, Two and Six LIVA",
					i18n.ES: "Factura rectificativa: error fundado en derecho y Art. 80 Uno, Dos y Seis LIVA",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Use this code when correcting most commercial invoices due to cancellations
						or discounts. This is currently set as the default buy may be overridden if
						needed.
					`),
				},
			},
			{
				Code: "R2",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: artículo 80.3",
					i18n.EN: "Rectified invoice: error based on law and Article 80.3",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						To use for customer declared insolvency proceedings when a court
						is involved.
					`),
				},
			},
			{
				Code: "R3",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: artículo 80.4",
					i18n.EN: "Rectified invoice: error based on law and Article 80.4",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						For unpaid invoices that are not declared as related to insolvency
						and related to bad debt after a 6 or 12 month waiting period.
					`),
				},
			},
			{
				Code: "R4",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: Resto",
					i18n.EN: "Rectified invoice: Other",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Legal or court-imposed corrections that do not fall under any of
						the other corrective reasons.
					`),
				},
			},
			{
				Code: "R5",
				Name: i18n.String{
					i18n.ES: "Factura rectificativa: facturas simplificadas",
					i18n.EN: "Rectified invoice: simplified invoices",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Always used when correcting simplified or B2C invoices.
					`),
				},
			},
		},
	},
	{
		Key: ExtKeyCorrectionType,
		Name: i18n.String{
			i18n.EN: "Verifactu Correction Type Code",
			i18n.ES: "Código de Tipo de Corrección de Verifactu",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Correction type code used to identify the type of correction being
				made. Values map to L3 list.
				
				Code is determined automatically according to the invoice type:

				| Invoice Type		| Code |
				|-------------------|------|
				| ~corrective~		| ~S~  |
				| ~credit-note~		| ~I~  |
				| ~debit-note~		| ~I~  |
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "S",
				Name: i18n.String{
					i18n.EN: "Substitution",
					i18n.ES: "Por Sustitución",
				},
			},
			{
				Code: "I",
				Name: i18n.String{
					i18n.EN: "Differences",
					i18n.ES: "Por Diferencias",
				},
			},
		},
	},
	{
		Key: ExtKeyOpClass,
		Name: i18n.String{
			i18n.EN: "Subject and Not Exempt Operation Class Code",
			i18n.ES: "Clave de la Operación Sujeta y no Exenta o de la Operación no Sujeta.",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Classification code for operations that are subject to tax and not exempt, or for operations not subject to tax.

				VERI*FACTU distinguishes between "subject" and "exempt" operations, while GOBL treats "exempt" as simply having no tax percentage. This difference can make mapping between systems complex.

				GOBL will attempt to automatically assign operation class codes based on tax rates, but if your workflow requires more control, you may prefer to let users select the appropriate operation class and exemption code for each case.

				Automatic mapping of tax rates to operation classes:

				| Tax Rate                | Operation Class |
				|-------------------------|-----------------|
				| ~standard~              | ~S1~            |
				| ~reduced~               | ~S1~            |
				| ~super-reduced~         | ~S1~            |
				| ~zero~                  | ~S1~            |

				S2 should be set for domestic or specific operations where there is a reverse charge of VAT (the customer in Spain assumes the VAT payment). For exports, use exemption codes.

				This extension maps to the ~CalificacionOperacion~ field and must not be used together with the ~es-verifactu-exempt~ extension. Values correspond to the L9 list.

				For further guidance on applying these codes, refer to the AEAT FAQ:
				 * https://sede.agenciatributaria.gob.es/Sede/impuestos-tasas/iva/iva-libros-registro-iva-traves-aeat/preguntas-frecuentes/3-libro-registro-facturas-expedidas.html?faqId=b5556c3d02bc9510VgnVCM100000dc381e0aRCRD
				`),
		},
		Values: []*cbc.Definition{
			{
				Code: "S1",
				Name: i18n.String{
					i18n.EN: "Subject and Not Exempt - Without reverse charge",
					i18n.ES: "Operación Sujeta y No exenta - Sin inversión del sujeto pasivo",
				},
			},
			{
				Code: "S2",
				Name: i18n.String{
					i18n.EN: "Subject and Not Exempt - With reverse charge",
					i18n.ES: "Operación Sujeta y No exenta - Con Inversión del sujeto pasivo",
				},
			},
			{
				Code: "N1",
				Name: i18n.String{
					i18n.EN: "Not Subject - Articles 7, 14, others",
					i18n.ES: "Operación No Sujeta artículo 7, 14, otros",
				},
			},
			{
				Code: "N2",
				Name: i18n.String{
					i18n.EN: "Not Subject - Due to location rules",
					i18n.ES: "Operación No Sujeta por Reglas de localización",
				},
			},
		},
	},
	{
		Key: ExtKeyExempt,
		Name: i18n.String{
			i18n.EN: "Verifactu Exemption Code",
			i18n.ES: "Código de Exención de Impuesto de Verifactu",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Exemption code used to explain why the operation is exempt from taxes.

				This extension maps to the field ~OperacionExenta~, and **cannot** be provided
				alongside the ~es-verifactu-op-class~ extension. Values correspond to the
				L10 list.

				Automatic mapping of tax rates to exemption codes:

				| Tax Rate                | Operation Class |
				|-------------------------|-----------------|
				| ~exempt+export~         | ~E2~            |
				| ~exempt+export+eea~     | ~E5~            |

				The rest of the exemptions will need to be determined and applied on a
				case-by-case basis.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "E1",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 20. Exemptions in internal operations.",
					i18n.ES: "Exenta: por el artículo 20. Exenciones en operaciones interiores.",
				},
			},
			{
				Code: "E2",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 21. Exemptions in exports of goods outside the EU.",
					i18n.ES: "Exenta: por el artículo 21. Exenciones en las exportaciones de bienes fuera de la UE.",
				},
			},
			{
				Code: "E3",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 22. Exemptions in operations asimilated to exports.",
					i18n.ES: "Exenta: por el artículo 22. Exenciones en las operaciones asimiladas a las exportaciones.",
				},
			},
			{
				Code: "E4",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Articles 23 and 24. Exemptions related to temporary deposit, customs and fiscal regimes, and other situations.",
					i18n.ES: "Exenta: por el artículos 23 y 24. Exenciones relativas a las situaciones de depósito temporal, regímenes aduaneros y fiscales, y otras situaciones.",
				},
			},
			{
				Code: "E5",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to Article 25. Exemptions in the delivery of goods destined to another Member State.",
					i18n.ES: "Exenta: por el artículo 25. Exenciones en las entregas de bienes destinados a otro Estado miembro.",
				},
			},
			{
				Code: "E6",
				Name: i18n.String{
					i18n.EN: "Exempt: pursuant to other reasons",
					i18n.ES: "Exenta: por otra causa",
				},
			},
		},
	},
	{
		Key: ExtKeyRegime,
		Name: i18n.String{
			i18n.EN: "VAT/IGIC Regime Code",
			i18n.ES: "Código de Régimen de IVA/IGIC",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Identify the type of VAT or IGIC regime applied to the operation. This list combines
				lists L8A which include values for VAT, and L8B for IGIC.

				Maps to the field ~ClaveRegimen~, and is required for all VAT and IGIC operations.
				Values correspond to L8A (VAT) and L8B (IGIC) lists.

				The regime code must be assigned for each tax combo. If no regime code is provided,
				GOBL will try to assign a code from the following tax combo contexts:

				| Combo Context				| Regime Code |
				|---------------------------|-------------|
				| Rate ~standard~			| ~01~        |
				| Rate has ~export~			| ~02~        |
				| Has surcharge				| ~18~        |
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
					i18n.EN: "Special regime for VAT/IGIC groups (Advanced Level)",
					i18n.ES: "Régimen especial grupo de entidades en IVA/IGIC (Nivel Avanzado)",
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
					i18n.EN: "Operations subject to a different regime",
					i18n.ES: "Operaciones sujetas a un régimen diferente",
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
					i18n.EN: "Collection of professional fees or rights on behalf of third parties",
					i18n.ES: "Cobros por cuenta de terceros de honorarios profesionales o de derechos derivados de la propiedad industrial",
				},
			},
			{
				Code: "11",
				Name: i18n.String{
					i18n.EN: "Business premises rental operations",
					i18n.ES: "Operaciones de arrendamiento de local de negocio",
				},
			},
			{
				Code: "14",
				Name: i18n.String{
					i18n.EN: "Invoice with pending VAT/IGIC accrual in work certifications for Public Administration",
					i18n.ES: "Factura con IVA/IGIC pendiente de devengo en certificaciones de obra cuyo destinatario sea una Administración Pública",
				},
			},
			{
				Code: "15",
				Name: i18n.String{
					i18n.EN: "Invoice with pending VAT/IGIC accrual in successive tract operations",
					i18n.ES: "Factura con IVA/IGIC pendiente de devengo en operaciones de tracto sucesivo",
				},
			},
			{
				Code: "17",
				Name: i18n.String{
					i18n.EN: "Operation under OSS and IOSS regimes (VAT) / Special regime for retail traders. (IGIC)",
					i18n.ES: "Operación acogida a alguno de los regímenes previstos en el capítulo XI del título IX (OSS e IOSS, IVA) / Régimen especial de comerciante minorista. (IGIC)",
				},
			},
			{
				Code: "18",
				Name: i18n.String{
					i18n.EN: "Equivalence surcharge (VAT) / Special regime for small traders or retailers (IGIC)",
					i18n.ES: "Recargo de equivalencia (IVA) / Régimen especial del pequeño comerciante o minorista (IGIC)",
				},
			},
			{
				Code: "19",
				Name: i18n.String{
					i18n.EN: "Operations included in the Special Regime for Agriculture, Livestock and Fisheries",
					i18n.ES: "Operaciones de actividades incluidas en el Régimen Especial de Agricultura, Ganadería y Pesca (REAGYP)",
				},
			},
			{
				Code: "20",
				Name: i18n.String{
					i18n.EN: "Simplified regime (VAT only)",
					i18n.ES: "Régimen simplificado (IVA only)",
				},
			},
		},
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
				defined in the L7 list.

				The regular Party Tax Identity is preferred over using a specific identity type
				code, and will be mapped automatically as follows:
				
				- Spanish Tax IDs will be mapped to the ~NIF~ field.
				- EU Tax IDs will be mapped to the ~IDOtro~ field with code ~02~.
				- Non-EU Tax IDs will be mapped to the ~IDOtro~ field with code ~04~.

				VERI*FACTU will perform validation on both Spanish and EU Tax IDs, so it is important
				to provide the correct details.

				The following identity ~key~ values will be mapped automatically to an extension by the 
				addon for the following keys:

				- ~passport~: ~03~
				- ~foreign~: ~04~
				- ~resident~: ~05~
				- ~other~: ~06~

				The ~07~ "not registered in census" code is not mapped automatically, but
				can be provided directly if needed.

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
								"es-verifactu-identity-type": "03"
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
			{
				// We don't set a constant for this as there are no mappings that
				// will be done automatically.
				Code: "07",
				Name: i18n.String{
					i18n.EN: "Not registered in census",
					i18n.ES: "No censado",
				},
			},
		},
	},
	{
		Key: ExtKeySimplifiedArt7273,
		Name: i18n.String{
			i18n.EN: "Simplified Invoice Art. 7.2 and 7.3, RD 1619/2012",
			i18n.ES: "Factura Simplificada Articulo 7,2 y 7,3 RD 1619/2012",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				This extensions covers a specific use-case when the customer specifically
				requests that the invoice includes their fiscal details, but they are
				not registered for tax.

				Maps to the ~FacturaSimplificadaArt7273~ field in Verifactu documents.

				Can only be true when the invoice type (~TipoFactura~) is one of: ~F1~,
				~F3~, ~R1~, ~R2~, ~R3~, or ~R4~.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "S",
				Name: i18n.String{
					i18n.EN: "Yes",
					i18n.ES: "Sí",
				},
			},
			{
				Code: "N",
				Name: i18n.String{
					i18n.EN: "No",
					i18n.ES: "No",
				},
			},
		},
	},
	{
		Key: ExtKeyIssuerType,
		Name: i18n.String{
			i18n.EN: "Issuer Type Code",
			i18n.ES: "Emitida por Tercero o Destinatario",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Indicates whether the invoice is issued by a third party or by the customer
				themselves.

				Mapped to the field ~EmitidaPorTerceroODestinatario~ in Verifactu documents,
				with list L6.

				The ~self-billed~ tag will automatically be set this extension in the invoice
				to ~D~.

				If the ~issuer~ field is set in the invoice's ordering section, then this
				extension will be set to ~T~.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: ExtCodeIssuerTypeThirdParty,
				Name: i18n.String{
					i18n.EN: "Issued by Third Party",
					i18n.ES: "Emitida por Tercero",
				},
			},
			{
				Code: ExtCodeIssuerTypeCustomer,
				Name: i18n.String{
					i18n.EN: "Issued by Customer",
					i18n.ES: "Emitida por Destinatario",
				},
			},
		},
	},
}
