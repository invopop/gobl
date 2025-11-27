package sii

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Extension keys for SII
const (
	ExtKeyDocType           cbc.Key = "es-sii-doc-type"
	ExtKeyCorrectionType    cbc.Key = "es-sii-correction-type"
	ExtKeyNotExempt         cbc.Key = "es-sii-not-exempt"
	ExtKeyNotSubject        cbc.Key = "es-sii-not-subject"
	ExtKeyExempt            cbc.Key = "es-sii-exempt"
	ExtKeyRegime            cbc.Key = "es-sii-regime"
	ExtKeyIdentityType      cbc.Key = "es-sii-identity-type"
	ExtKeySimplifiedArt7273 cbc.Key = "es-sii-simplified-art7273"
	ExtKeyNonSupplierIssuer cbc.Key = "es-sii-non-supplier-issuer"
	ExtKeyProduct           cbc.Key = "es-sii-product"
)

// Identity Type Codes
const (
	ExtCodeIdentityTypeVAT      cbc.Code = "02" // NIF-VAT identity
	ExtCodeIdentityTypePassport cbc.Code = "03" // Passport
	ExtCodeIdentityTypeForeign  cbc.Code = "04" // Foreign Identity Document
	ExtCodeIdentityTypeResident cbc.Code = "05" // Spanish Resident Foreigner Identity Card
	ExtCodeIdentityTypeOther    cbc.Code = "06" // Other Identity Document
)

// Product Type Codes
const (
	ExtCodeProductGoods    cbc.Code = "goods"    // Delivery of goods
	ExtCodeProductServices cbc.Code = "services" // Provision of services
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyDocType,
		Name: i18n.String{
			i18n.EN: "Invoice Type",
			i18n.ES: "Tipo de Factura",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Ministerial Order HFP/417/2017, of May 12th",
					i18n.ES: "Orden Ministerial HFP/417/2017, de 12 de Mayo",
				},
				URL: "https://www.boe.es/buscar/act.php?id=BOE-A-2017-5312",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code used to identify the type of invoice being sent or received. This will be determined
				automatically by GOBL during normalization according to the scenario definitions.

				The following codes are not covered by GOBL's scenarios and will need to be set manually if
				needed: ~F4~, ~F5~, ~F6~, ~R2~, ~R3~, ~R4~, ~AJ~, ~LC~.

				Maps to the ~TipoFactura~ field. Values correspond to the L2_EMI (issued invoices) and
				L2_RECI (received invoices) lists.
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
						To use for B2C invoices where details about the customer are not normally
						required.
					`),
				},
			},
			{
				Code: "F3",
				Name: i18n.String{
					i18n.EN: "Invoice issued as replacement of simplified invoice",
					i18n.ES: "Factura emitida en sustitución de facturas simplificadas",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						To use when a simplified invoice is replaced by a regular invoice.
					`),
				},
			},
			{
				Code: "F4",
				Name: i18n.String{
					i18n.EN: "Summary entry of invoices",
					i18n.ES: "Asiento resumen de facturas",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Summary entry grouping multiple invoices.
					`),
				},
			},
			{
				Code: "F5",
				Name: i18n.String{
					i18n.EN: "Imports (DUA)",
					i18n.ES: "Importaciones (DUA)",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						For import operations using the Single Administrative Document (DUA). Only
						applicable to received invoices (L2_RECI).
					`),
				},
			},
			{
				Code: "F6",
				Name: i18n.String{
					i18n.EN: "Accounting justifications",
					i18n.ES: "Justificantes contables",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						For accounting justifications. Only applicable to received invoices.
					`),
				},
			},
			{
				Code: "R1",
				Name: i18n.String{
					i18n.EN: "Rectified invoice: articles 80.1, 80.2, and 80.6",
					i18n.ES: "Factura rectificativa: art. 80 Uno, Dos y Seis LIVA",
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
					i18n.EN: "Rectified invoice: article 80.3",
					i18n.ES: "Factura rectificativa: artículo 80.3",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						To use for customer declared insolvency proceedings when a court is
						involved.
					`),
				},
			},
			{
				Code: "R3",
				Name: i18n.String{
					i18n.EN: "Rectified invoice: article 80.4",
					i18n.ES: "Factura rectificativa: artículo 80.4",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						For unpaid invoices that are not declared as related to insolvency and
						related to bad debt after a 6 or 12 month waiting period.
					`),
				},
			},
			{
				Code: "R4",
				Name: i18n.String{
					i18n.EN: "Rectified invoice: other",
					i18n.ES: "Factura rectificativa: resto",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Legal or court-imposed corrections that do not fall under any of the other
						corrective reasons.
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
			{
				Code: "AJ",
				Name: i18n.String{
					i18n.EN: "Profit margin adjustment",
					i18n.ES: "Ajuste del margen de beneficio",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						For profit margin adjustments. Only applicable to issued invoices (L2_EMI).
					`),
				},
			},
			{
				Code: "LC",
				Name: i18n.String{
					i18n.EN: "Customs - Complementary settlement",
					i18n.ES: "Aduanas - Liquidación complementaria",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						For complementary customs settlements. Only applicable to received invoices
						(L2_RECI).
					`),
				},
			},
		},
	},
	{
		Key: ExtKeyCorrectionType,
		Name: i18n.String{
			i18n.EN: "Corrective Invoice Type",
			i18n.ES: "Tipo de Factura Rectificativa",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Ministerial Order HFP/417/2017, of May 12th",
					i18n.ES: "Orden Ministerial HFP/417/2017, de 12 de Mayo",
				},
				URL: "https://www.boe.es/buscar/act.php?id=BOE-A-2017-5312",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Correction type code used to identify the type of correction being made.

				Code is determined automatically according to the invoice type:

				| Invoice Type  | Code |
				| ------------- | ---- |
				| ~corrective~  | ~S~  |
				| ~credit-note~ | ~I~  |
				| ~debit-note~  | ~I~  |

				Maps to the ~TipoRectificativa~ field. Values correspond to the L5 list.
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
		Key: ExtKeyExempt,
		Name: i18n.String{
			i18n.EN: "Exemption Reason",
			i18n.ES: "Causa de Exención",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Ministerial Order HFP/417/2017, of May 12th",
					i18n.ES: "Orden Ministerial HFP/417/2017, de 12 de Mayo",
				},
				URL: "https://www.boe.es/buscar/act.php?id=BOE-A-2017-5312",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Exemption code used to explain why the operation is exempt from taxes.

				The follow mappings will be made automatically by GOBL during normalization:

				| Tax Key           | Exemption Codes            |
				|-------------------|----------------------------|
				| ~exempt~          | ~E1~ (default), ~E6~       |
				| ~export~          | ~E2~ (default), ~E3~, ~E4~ |
				| ~intra-community~ | ~E5~                       |

				Maps to the field ~CausaExencion~. Values correspond to the L9 list.

				Note: This extension is **mutually exclusive** with ~es-sii-not-exempt~ and ~es-sii-not-subject~.
				Only one of these three extensions can be used for a given tax combo.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "E1",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Art. 20 (internal operations).",
					i18n.ES: "Exenta por el art. 20 (operaciones interiores).",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Use for **domestic VAT-exempt operations**: healthcare by medical
						professionals, education by authorised centres, social assistance, certain
						cultural services, financial and insurance services, letting of dwellings,
						etc.
					`),
				},
			},
			{
				Code: "E2",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Art. 21 (exports of goods).",
					i18n.ES: "Exenta por el art. 21 (exportaciones de bienes).",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Use for **exports of goods outside the EU** (including the tax-free
						travellers regime where the sale is exported) and services directly related
						to those exports, under the conditions of the article.
					`),
				},
			},
			{
				Code: "E3",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Art. 22 (operations asimilated to exports).",
					i18n.ES: "Exenta por el art. 22 (operaciones asimiladas a las exportaciones).",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Use for **operations assimilated to exports**, e.g.,
						supplies/repairs/charter of qualifying ships or aircraft, avituallamiento
						(provisioning) of such vessels, and certain services directly connected with
						international transport.
					`),
				},
			},
			{
				Code: "E4",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Art. 23 and 24 (temporary deposit, customs and fiscal regimes, and other situations).",
					i18n.ES: "Exenta por los art. 23 y 24 (situaciones de depósito temporal, regímenes aduaneros y fiscales, y otras situaciones).",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Use for operations linked to **free zones, customs/duty-free warehouses and
						other customs or fiscal regimes** (e.g., depósitos francos, depósito
						temporal, perfeccionamiento activo/pasivo, depósito distinto del aduanero),
						while the goods remain under those regimes.
					`),
				},
			},
			{
				Code: "E5",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to Art. 25 (delivery of goods destined to another Member State).",
					i18n.ES: "Exenta por el art. 25 (entregas de bienes destinados a otro Estado miembro).",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Use for **intra-Community supplies of goods** from Spain to another EU
						Member State when the buyer is VAT-identified in another Member State and
						the goods are shipped there.
					`),
				},
			},
			{
				Code: "E6",
				Name: i18n.String{
					i18n.EN: "Exempt pursuant to other reasons",
					i18n.ES: "Exenta por otra causa",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Use when the operation is exempt for a reason not covered by Articles 20–25,
						and a lawful exemption applies (catch-all "Exenta por otros"). Document the
						legal basis on the invoice text.
					`),
				},
			},
		},
	},
	{
		Key: ExtKeyNotExempt,
		Name: i18n.String{
			i18n.EN: "Type of Operation Subject and Not Exempt",
			i18n.ES: "Tipo de Operación Sujeta y No Exenta",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Ministerial Order HFP/417/2017, of May 12th",
					i18n.ES: "Orden Ministerial HFP/417/2017, de 12 de Mayo",
				},
				URL: "https://www.boe.es/buscar/act.php?id=BOE-A-2017-5312",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Type of operation subject and not exempt for the differentiation of reverse charge.

				GOBL will attempt to automatically assign operation type codes based on tax key:

				| Operation Type | Tax Key                                        |
				| -------------- | ---------------------------------------------- |
				| ~S1~           | ~standard~, ~reduced~, ~super-reduced~, ~zero~ |
				| ~S2~           | ~reverse-charge~                               |

				Maps to the ~TipoNoExenta~ field. Values correspond to the L7 list. The ~S3~ code is not
				meant to be set manually, it will only be set internally when both ~S1~ and ~S2~ are present
				in the same invoice.

				Note: This extension is **mutually exclusive** with ~es-sii-exempt~ and ~es-sii-not-subject~.
				Only one of these three extensions can be used for a given tax combo.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "S1",
				Name: i18n.String{
					i18n.EN: "Non-exempt - Without reverse charge",
					i18n.ES: "No exenta - Sin inversión sujeto pasivo",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Use for operations that are subject to tax and not exempt, without reverse charge.
					`),
				},
			},
			{
				Code: "S2",
				Name: i18n.String{
					i18n.EN: "Non-exempt - With reverse charge",
					i18n.ES: "No exenta - Con Inversión sujeto pasivo",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Use for operations that are subject to tax and not exempt, with reverse charge.
					`),
				},
			},
			{
				Code: "S3",
				Name: i18n.String{
					i18n.EN: "Non-exempt - Without reverse charge and with reverse charge",
					i18n.ES: "No exenta - Sin inversión sujeto pasivo y con Inversión sujeto pasivo",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Use for operations that are subject to tax and not exempt, containing both operations
						without reverse charge and operations with reverse charge.
					`),
				},
			},
		},
	},
	{
		Key: ExtKeyNotSubject,
		Name: i18n.String{
			i18n.EN: "Type of Operation Not Subject",
			i18n.ES: "Tipo de Operación No Sujeta",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Ministerial Order HFP/417/2017, of May 12th",
					i18n.ES: "Orden Ministerial HFP/417/2017, de 12 de Mayo",
				},
				URL: "https://www.boe.es/buscar/act.php?id=BOE-A-2017-5312",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Type of operation not subject to VAT.

				This extension is used to determine which tax amount field should be reported in the SII
				payload:

				- ~N1~ will set the ~ImportePorArticulos7_14_Otros~ field.
				- ~N2~ will set the ~ImporteTAIReglasLocalizacion~ field.

				GOBL will attempt to automatically assign operation type codes based on tax key:

				| Operation Type | Tax Key                   |
				| -------------- | ------------------------- |
				| ~N1~           | ~outside-scope~           |
				| ~N2~           | ~outside-scope~ (default) |

				Doesn't map directly to any field. Used internally to determine how to report the tax amount.

				Note: This extension is **mutually exclusive** with ~es-sii-exempt~ and ~es-sii-not-exempt~.
				Only one of these three extensions can be used for a given tax combo.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "N1",
				Name: i18n.String{
					i18n.EN: "Not Subject - Articles 7, 14, others",
					i18n.ES: "Operación No Sujeta artículo 7, 14, otros",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Use when the operations fall **outside the VAT scope** under Spanish VAT law
						(e.g. transfers of companies, certain exchanges of goods, contributions of
						goods, internal operations of public entities, ...). VAT is **not
						chargeable** or declared.
					`),
				},
			},
			{
				Code: "N2",
				Name: i18n.String{
					i18n.EN: "Not Subject - Due to location rules",
					i18n.ES: "Operación No Sujeta por Reglas de localización",
				},
				Desc: i18n.String{
					i18n.EN: here.Doc(`
						Use when the operation is **not subject to VAT due to its place of supply**
						(e.g. services provided to non‑EU B2B customers that are deemed outside
						Spanish VAT by location rules).
					`),
				},
			},
		},
	},
	{
		Key: ExtKeyRegime,
		Name: i18n.String{
			i18n.EN: "Special Regime or Relevance Key",
			i18n.ES: "Clave de Régimen Especial o Trascendencia",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Ministerial Order HFP/417/2017, of May 12th",
					i18n.ES: "Orden Ministerial HFP/417/2017, de 12 de Mayo",
				},
				URL: "https://www.boe.es/buscar/act.php?id=BOE-A-2017-5312",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Identify the regime applied to the operation.

				The regime key must be assigned for each tax combo. If no regime key is provided, GOBL will
				try to assign a code from the following tax combo contexts:

				| Combo Context  | Regime Code |
				| -------------- | ----------- |
				| Key ~standard~ | ~01~        |
				| Key ~export~   | ~02~        |

				Maps to the field ~ClaveRegimenEspecialOTrascendencia~. Values correspond to the L3.1
				(issued invoices) and L3.2 (received invoices) lists.
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
					i18n.EN: "Export (L3.1) / Operations for which businesses pay compensation in acquisitions from individuals under the special regime for agriculture, livestock, and fishing (L3.2)",
					i18n.ES: "Exportación (L3.1) / Operaciones por las que los empresarios satisfacen compensaciones en las adquisiciones a personas acogidas al Régimen especial de la agricultura, ganadería y pesca (L3.2)",
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
					i18n.EN: "Operations subject to IPSI / IGIC (Tax on Production, Services and Imports / Canary Islands General Indirect Tax)",
					i18n.ES: "Operaciones sujetas al IPSI / IGIC (Impuesto sobre la Producción, los Servicios y la Importación / Impuesto General Indirecto Canario)",
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
				Code: "12",
				Name: i18n.String{
					i18n.EN: "Business premises rental operations not subject to withholding",
					i18n.ES: "Operaciones de arrendamiento de local de negocio no sujetos a retención",
				},
			},
			{
				Code: "13",
				Name: i18n.String{
					i18n.EN: "Business premises rental operations subject and not subject to withholding (L3.1) / Invoice corresponding to an import (reported without being associated with a DUA) (L3.2)",
					i18n.ES: "Operaciones de arrendamiento de local de negocio sujetas y no sujetas a retención (L3.1) / Factura correspondiente a una importación (informada sin asociar a un DUA) (L3.2)",
				},
			},
			{
				Code: "14",
				Name: i18n.String{
					i18n.EN: "Invoice with pending VAT/IGIC accrual in work certifications for Public Administration (L3.1) / First half of 2017 and other invoices prior to inclusion in the SII (L3.2)",
					i18n.ES: "Factura con IVA pendiente de devengo en certificaciones de obra cuyo destinatario sea una Administración Pública (L3.1) / Primer semestre 2017 y otras facturas anteriores a la inclusión en el SII (L3.2)",
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
				Code: "16",
				Name: i18n.String{
					i18n.EN: "First half of 2017 and other invoices prior to inclusion in the SII",
					i18n.ES: "Primer semestre 2017 y otras facturas anteriores a la inclusión en el SII",
				},
			},
			{
				Code: "17",
				Name: i18n.String{
					i18n.EN: "Operation covered by one of the regimes provided for in Chapter XI of Title IX (OSS and IOSS)",
					i18n.ES: "Operación acogida a alguno de los regímenes previstos en el Capítulo XI del Título IX (OSS e IOSS)",
				},
			},
		},
	},
	{
		Key: ExtKeyIdentityType,
		Name: i18n.String{
			i18n.EN: "Identification Type",
			i18n.ES: "Tipos de Identificación",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Ministerial Order HFP/417/2017, of May 12th",
					i18n.ES: "Orden Ministerial HFP/417/2017, de 12 de Mayo",
				},
				URL: "https://www.boe.es/buscar/act.php?id=BOE-A-2017-5312",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Identity code used to identify the type of identity document used by the customer.

				The regular Party Tax Identity is preferred over using a specific identity type code, and
				will be mapped automatically as follows:

				- Spanish Tax IDs will be mapped to the ~NIF~ field.
				- EU Tax IDs will be mapped to the ~IDOtro\IDType~ field with code ~02~.
				- Non-EU Tax IDs will be mapped to the ~IDOtro\IDType~ field with code ~04~.

				SII will perform validation on both Spanish and EU Tax IDs, so it is important to provide
				the correct details.

				The following identity ~key~ values will be mapped automatically to an extension by the
				addon:

				| Identity Key | Code |
				|--------------|------|
				| ~passport~   | ~03~ |
				| ~foreign~    | ~04~ |
				| ~resident~  | ~05~ |
				| ~other~      | ~06~ |

				The ~07~ "not registered in census" code is not mapped automatically, but can be provided
				directly if needed.

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
								"es-sii-identity-type": "03"
							}
						}
					]
				}
				~~~

				Maps to the field ~IDOtro\IDType~. Values correspond to the L4 list.
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
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Ministerial Order HFP/417/2017, of May 12th",
					i18n.ES: "Orden Ministerial HFP/417/2017, de 12 de Mayo",
				},
				URL: "https://www.boe.es/buscar/act.php?id=BOE-A-2017-5312",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				This extensions covers a specific use-case when the customer specifically requests that the
				invoice includes their fiscal details, but they are not registered for tax.

				Can only be true when the invoice type (~TipoFactura~) is one of: ~F1~, ~F3~, ~R1~, ~R2~,
				~R3~, or ~R4~.

				Maps to the ~FacturaSimplificadaArticulos7.2_7.3~ field. Values correspond to the L26 list.
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
		Key: ExtKeyNonSupplierIssuer,
		Name: i18n.String{
			i18n.EN: "Issued by Third Party or Recepient",
			i18n.ES: "Emitida por Tercero o Destinatario",
		},
		Sources: []*cbc.Source{
			{
				Title: i18n.String{
					i18n.EN: "Ministerial Order HFP/417/2017, of May 12th",
					i18n.ES: "Orden Ministerial HFP/417/2017, de 12 de Mayo",
				},
				URL: "https://www.boe.es/buscar/act.php?id=BOE-A-2017-5312",
			},
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Indicates whether the invoice is issued by a third party or by the customer themselves.

				The default value is ~N~ (No).

				The ~self-billed~ tag will automatically be set this extension in the invoice to ~S~ (Yes).

				If the ~issuer~ field is set in the invoice's ordering section, then this extension will be
				set to ~S~ (Yes) too.

				Maps to the field ~EmitidaPorTercerosODestinatarios~. Values correspond to the L10 list.
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
		Key: ExtKeyProduct,
		Name: i18n.String{
			i18n.EN: "Product Type",
			i18n.ES: "Tipo de Producto",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Product type code used to differentiate between goods and services for the purpose of
				reporting breakdowns in the SII format.

				This extension is used to determine the type of operation breakdown when generating the SII
				report. When provided, the value will be used to generate the ~DesgloseTipoOperacion~ field,
				selecting between ~PrestacionServices~ (provision of services) or ~Entrega~ (delivery of
				goods).

				This extension is optional; if not provided, the breakdown will use ~DesgloseFactura~ instead
				of ~DesgloseTipoOperacion~.

				Doesn't map directly to any field. Used internally to structure the breakdown data.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: ExtCodeProductGoods,
				Name: i18n.String{
					i18n.EN: "Delivery of goods",
					i18n.ES: "Entrega de bienes",
				},
			},
			{
				Code: ExtCodeProductServices,
				Name: i18n.String{
					i18n.EN: "Provision of services",
					i18n.ES: "Prestación de servicios",
				},
			},
		},
	},
}
