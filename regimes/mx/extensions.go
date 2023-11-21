package mx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Mexican CFDI extension keys required by the SAT (tax authority in Mexico) in all
// invoices and cannot be determined automatically.
const (
	ExtKeyCFDIFiscalRegime = "mx-cfdi-fiscal-regime"
	ExtKeyCFDIUse          = "mx-cfdi-use"
	ExtKeyCFDIProdServ     = "mx-cfdi-prod-serv" // name from XML field: ClaveProdServ
)

// Mabe Addendum extension keys
const (
	ExtKeyMabeProviderCode  = "mx-mabe-provider-code"
	ExtKeyMabeDeliveryPlant = "mx-mabe-delivery-plant"
	ExtKeyMabeItemCode      = "mx-mabe-item-code"
	ExtKeyMabeReference1    = "mx-mabe-reference-1"
	ExtKeyMabeReference2    = "mx-mabe-reference-2"
)

var extensionKeys = []*tax.KeyDefinition{
	{
		Key: ExtKeyCFDIProdServ,
		Name: i18n.String{
			i18n.EN: "Product or Service Code",
			i18n.ES: "Clave de Producto o Servicio", //nolint:misspell
		},
		Desc: i18n.String{
			i18n.EN: "Code defined in the CFDI catalogue used to identify a product or service.",
			i18n.ES: "Código definido en el catálogo del CFDI utilizado para identificar un producto o servicio.", //nolint:misspell
		},
	},
	{
		Key: ExtKeyCFDIFiscalRegime,
		Name: i18n.String{
			i18n.EN: "Fiscal Regime Code",
			i18n.ES: "Código de Régimen Fiscal",
		},
		Desc: i18n.String{
			i18n.EN: "Fiscal regime associated with suppliers and customers.",
			i18n.ES: "Régimen fiscal asociado con el emisor y receptor.",
		},
		Codes: []*tax.CodeDefinition{
			{
				Code: "601",
				Name: i18n.String{
					i18n.ES: "General de Ley Personas Morales",
				},
			},
			{
				Code: "603",
				Name: i18n.String{
					i18n.ES: "Personas Morales con Fines no Lucrativos",
				},
			},
			{
				Code: "605",
				Name: i18n.String{
					i18n.ES: "Sueldos y Salarios e Ingresos Asimilados a Salarios",
				},
			},
			{
				Code: "606",
				Name: i18n.String{
					i18n.ES: "Arrendamiento",
				},
			},
			{
				Code: "607",
				Name: i18n.String{
					i18n.ES: "Régimen de Enajenación o Adquisición de Bienes",
				},
			},
			{
				Code: "608",
				Name: i18n.String{
					i18n.ES: "Demás ingresos",
				},
			},
			{
				Code: "610",
				Name: i18n.String{
					i18n.ES: "Residentes en el Extranjero sin Establecimiento Permanente en México",
				},
			},
			{
				Code: "611",
				Name: i18n.String{
					i18n.ES: "Ingresos por Dividendos (socios y accionistas)", //nolint:misspell
				},
			},
			{
				Code: "612",
				Name: i18n.String{
					i18n.ES: "Personas Físicas con Actividades Empresariales y Profesionales",
				},
			},
			{
				Code: "614",
				Name: i18n.String{
					i18n.ES: "Ingresos por intereses",
				},
			},
			{
				Code: "615",
				Name: i18n.String{
					i18n.ES: "Régimen de los ingresos por obtención de premios",
				},
			},
			{
				Code: "616",
				Name: i18n.String{
					i18n.ES: "Sin obligaciones fiscales",
				},
			},
			{
				Code: "620",
				Name: i18n.String{
					i18n.ES: "Sociedades Cooperativas de Producción que optan por diferir sus ingresos",
				},
			},
			{
				Code: "621",
				Name: i18n.String{
					i18n.ES: "Incorporación Fiscal",
				},
			},
			{
				Code: "622",
				Name: i18n.String{
					i18n.ES: "Actividades Agrícolas, Ganaderas, Silvícolas y Pesqueras",
				},
			},
			{
				Code: "623",
				Name: i18n.String{
					i18n.ES: "Opcional para Grupos de Sociedades",
				},
			},
			{
				Code: "624",
				Name: i18n.String{
					i18n.ES: "Coordinados",
				},
			},
			{
				Code: "625",
				Name: i18n.String{
					i18n.ES: "Régimen de las Actividades Empresariales con ingresos a través de Plataformas Tecnológicas",
				},
			},
			{
				Code: "626",
				Name: i18n.String{
					i18n.ES: "Régimen Simplificado de Confianza",
				},
			},
		},
	},
	{
		Key: ExtKeyCFDIUse,
		Name: i18n.String{
			i18n.EN: "CFDI Use Code",
			i18n.ES: "Código de Uso CFDI",
		},
		Desc: i18n.String{
			i18n.EN: "Chosen by the customer to indicate the purpose of an invoice.",
			i18n.ES: "Elegido por el cliente para indicar el propósito de una factura.",
		},
		Codes: []*tax.CodeDefinition{
			{
				Code: "G01",
				Name: i18n.String{
					i18n.EN: "Acquisition of goods",
					i18n.ES: "Adquisición de mercancías",
				},
			},
			{
				Code: "G02",
				Name: i18n.String{
					i18n.EN: "Returns, discounts or rebates",
					i18n.ES: "Devoluciones, descuentos o bonificaciones",
				},
			},
			{
				Code: "G03",
				Name: i18n.String{
					i18n.EN: "General expenses",
					i18n.ES: "Gastos en general",
				},
			},
			{
				Code: "I01",
				Name: i18n.String{
					i18n.EN: "Construction",
					i18n.ES: "Construcciones",
				},
			},
			{
				Code: "I02",
				Name: i18n.String{
					i18n.EN: "Office furniture and equipment as investmen",
					i18n.ES: "Mobiliario y equipo de oficina por inversiones",
				},
			},
			{
				Code: "I03",
				Name: i18n.String{
					i18n.EN: "Transport equipment",
					i18n.ES: "Equipo de transporte",
				},
			},
			{
				Code: "I04",
				Name: i18n.String{
					i18n.EN: "Computer equipment and accessories",
					i18n.ES: "Equipo de computo y accesorios",
				},
			},
			{
				Code: "I05",
				Name: i18n.String{
					i18n.EN: "Dies, punches, molds, matrices and other toolin",
					i18n.ES: "Dados, troqueles, moldes, matrices y herramental",
				},
			},
			{
				Code: "I06",
				Name: i18n.String{
					i18n.EN: "Telephone communications",
					i18n.ES: "Comunicaciones telefónicas",
				},
			},
			{
				Code: "I07",
				Name: i18n.String{
					i18n.EN: "Satellite communications",
					i18n.ES: "Comunicaciones satelitales",
				},
			},
			{
				Code: "I08",
				Name: i18n.String{
					i18n.EN: "Other machinery and equipment",
					i18n.ES: "Otra maquinaria y equipo",
				},
			},
			{
				Code: "D01",
				Name: i18n.String{
					i18n.EN: "Medical and dental fees and hospital expenses",
					i18n.ES: "Honorarios médicos, dentales y gastos hospitalarios",
				},
			},
			{
				Code: "D02",
				Name: i18n.String{
					i18n.EN: "Medical expenses for disability or incapacity",
					i18n.ES: "Gastos médicos por incapacidad o discapacidad",
				},
			},
			{
				Code: "D03",
				Name: i18n.String{
					i18n.EN: "Funeral expenses",
					i18n.ES: "Gastos funerales",
				},
			},
			{
				Code: "D04",
				Name: i18n.String{
					i18n.EN: "Donations",
					i18n.ES: "Donativos",
				},
			},
			{
				Code: "D05",
				Name: i18n.String{
					i18n.EN: "Interest actually paid on mortgage loans (housing)",
					i18n.ES: "Intereses reales efectivamente pagados por créditos hipotecarios (casa habitación)",
				},
			},
			{
				Code: "D06",
				Name: i18n.String{
					i18n.EN: "Voluntary contributions to the SAR",
					i18n.ES: "Aportaciones voluntarias al SAR",
				},
			},
			{
				Code: "D07",
				Name: i18n.String{
					i18n.EN: "Medical insurance premiums",
					i18n.ES: "Primas por seguros de gastos médicos",
				},
			},
			{
				Code: "D08",
				Name: i18n.String{
					i18n.EN: "Mandatory school transportation expenses",
					i18n.ES: "Gastos de transportación escolar obligatoria",
				},
			},
			{
				Code: "D09",
				Name: i18n.String{
					i18n.EN: "Deposits in savings accounts, pension plans premiums",
					i18n.ES: "Depósitos en cuentas para el ahorro, primas que tengan como base planes de pensiones",
				},
			},
			{
				Code: "D10",
				Name: i18n.String{
					i18n.EN: "Payments for educational services (school fees)",
					i18n.ES: "Pagos por servicios educativos (colegiaturas)",
				},
			},
			{
				Code: "S01",
				Name: i18n.String{
					i18n.EN: "Without tax effects",
					i18n.ES: "Sin efectos fiscales",
				},
			},
			{
				Code: "CP01",
				Name: i18n.String{
					i18n.EN: "Payments",
					i18n.ES: "Pagos",
				},
			},
			{
				Code: "CN01",
				Name: i18n.String{
					i18n.EN: "Payroll",
					i18n.ES: "Nómina",
				},
			},
		},
	},

	// Mabe Addendum Extension Key Definitions
	{
		Key: ExtKeyMabeProviderCode,
		Name: i18n.String{
			i18n.EN: "Mabe’s Provider Code",
			i18n.ES: "Código de Proveedor de Mabe", //nolint:misspell
		},
		Desc: i18n.String{
			i18n.EN: "Code used by Mabe to identify the supplier.",
			i18n.ES: "Código con el que Mabe identifica al proveedor.", //nolint:misspell
		},
	},
	{
		Key: ExtKeyMabeDeliveryPlant,
		Name: i18n.String{
			i18n.EN: "Mabe’s Delivery Plant ID",
			i18n.ES: "ID de Planta de Entrega de Mabe", //nolint:misspell
		},
		Desc: i18n.String{
			i18n.EN: "ID identifying the plant to which the invoice is addressed.",
			i18n.ES: "ID con el que se identifica la planta a la cual está dirigida la factura.", //nolint:misspell
		},
	},
	{
		Key: ExtKeyMabeItemCode,
		Name: i18n.String{
			i18n.EN: "Mabe’s Item Code",
			i18n.ES: "Código de Artículo de Mabe", //nolint:misspell
		},
		Desc: i18n.String{
			i18n.EN: "Code used by Mabe to identify the item.",
			i18n.ES: "Código con el que Mabe identifica el artículo.", //nolint:misspell
		},
	},
	{
		Key: ExtKeyMabeReference1,
		Name: i18n.String{
			i18n.EN: "Mabe’s Reference 1",
			i18n.ES: "Referencia 1 de Mabe", //nolint:misspell
		},
		Desc: i18n.String{
			i18n.EN: "Reference 1 used by Mabe.",
			i18n.ES: "Referencia 1 utilizada por Mabe.", //nolint:misspell
		},
	},
	{
		Key: ExtKeyMabeReference2,
		Name: i18n.String{
			i18n.EN: "Mabe’s Reference 2",
			i18n.ES: "Referencia 2 de Mabe", //nolint:misspell
		},
		Desc: i18n.String{
			i18n.EN: "Reference 2 used by Mabe.",
			i18n.ES: "Referencia 2 utilizada por Mabe.", //nolint:misspell
		},
	},
}

// IsMabeSupplier returns true if the invoice is from a Mabe supplier.
func IsMabeSupplier(inv *bill.Invoice) bool {
	return inv != nil && inv.Supplier != nil && inv.Supplier.Ext.Has(ExtKeyMabeProviderCode)
}
