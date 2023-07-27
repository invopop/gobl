package mx

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
)

// Regime Specific Tags (UsoCFDI codes)
const (
	TagUse cbc.Key = "use" // UsoCFDI codes "namespace"

	TagGoodsAcquisition     cbc.Key = "goods-acquisition"
	TagReturns              cbc.Key = "returns"
	TagGeneralExpenses      cbc.Key = "general-expenses"
	TagConstruction         cbc.Key = "construction"
	TagOfficeEquipment      cbc.Key = "office-equipment"
	TagTransportEquipment   cbc.Key = "transport-equipment"
	TagComputerEquipment    cbc.Key = "computer-equipment"
	TagManufacturingTooling cbc.Key = "manufacturing-tooling"
	TagTelephoneComms       cbc.Key = "telephone-comms"
	TagSatelliteComms       cbc.Key = "satellite-comms"
	TagOtherMachinery       cbc.Key = "other-machinery"
	TagMedicalExpenses      cbc.Key = "medical-expenses"
	TagDisability           cbc.Key = "disability"
	TagFuneralExpenses      cbc.Key = "funeral-expenses"
	TagDonation             cbc.Key = "donation"
	TagMortgageInterest     cbc.Key = "mortgage-interest"
	TagSARContribution      cbc.Key = "sar-contribution"
	TagMedicalInsurance     cbc.Key = "medical-insurance"
	TagSchoolTransportation cbc.Key = "school-transportation"
	TagSavingsDeposit       cbc.Key = "savings-deposit"
	TagSchoolFees           cbc.Key = "school-fees"
	TagNoTaxEffects         cbc.Key = "no-tax-effects"
	TagSuplementaryPayment  cbc.Key = "suplementary-payment"
	TagPayroll              cbc.Key = "payroll"
)

var invoiceTags = []*tax.KeyDefinition{

	{
		Key: TagUse.With(TagGoodsAcquisition),
		Name: i18n.String{
			i18n.EN: "Acquisition of goods",
			i18n.ES: "Adquisición de mercancías",
		},
	},
	{
		Key: TagUse.With(TagReturns),
		Name: i18n.String{
			i18n.EN: "Returns, discounts or rebates",
			i18n.ES: "Devoluciones, descuentos o bonificaciones",
		},
	},
	{
		Key: TagUse.With(TagGeneralExpenses),
		Name: i18n.String{
			i18n.EN: "General expenses",
			i18n.ES: "Gastos en general",
		},
	},
	{
		Key: TagUse.With(TagConstruction),
		Name: i18n.String{
			i18n.EN: "Construction",
			i18n.ES: "Construcciones",
		},
	},
	{
		Key: TagUse.With(TagOfficeEquipment),
		Name: i18n.String{
			i18n.EN: "Office furniture and equipment as investmen",
			i18n.ES: "Mobiliario y equipo de oficina por inversiones",
		},
	},
	{
		Key: TagUse.With(TagTransportEquipment),
		Name: i18n.String{
			i18n.EN: "Transport equipment",
			i18n.ES: "Equipo de transporte",
		},
	},
	{
		Key: TagUse.With(TagComputerEquipment),
		Name: i18n.String{
			i18n.EN: "Computer equipment and accessories",
			i18n.ES: "Equipo de computo y accesorios",
		},
	},
	{
		Key: TagUse.With(TagManufacturingTooling),
		Name: i18n.String{
			i18n.EN: "Dies, punches, molds, matrices and other toolin",
			i18n.ES: "Dados, troqueles, moldes, matrices y herramental",
		},
	},
	{
		Key: TagUse.With(TagTelephoneComms),
		Name: i18n.String{
			i18n.EN: "Telephone communications",
			i18n.ES: "Comunicaciones telefónicas",
		},
	},
	{
		Key: TagUse.With(TagSatelliteComms),
		Name: i18n.String{
			i18n.EN: "Satellite communications",
			i18n.ES: "Comunicaciones satelitales",
		},
	},
	{
		Key: TagUse.With(TagOtherMachinery),
		Name: i18n.String{
			i18n.EN: "Other machinery and equipment",
			i18n.ES: "Otra maquinaria y equipo",
		},
	},
	{
		Key: TagUse.With(TagMedicalExpenses),
		Name: i18n.String{
			i18n.EN: "Medical and dental fees and hospital expenses",
			i18n.ES: "Honorarios médicos, dentales y gastos hospitalarios",
		},
	},
	{
		Key: TagUse.With(TagMedicalExpenses).With(TagDisability),
		Name: i18n.String{
			i18n.EN: "Medical expenses for disability or incapacity",
			i18n.ES: "Gastos médicos por incapacidad o discapacidad",
		},
	},
	{
		Key: TagUse.With(TagFuneralExpenses),
		Name: i18n.String{
			i18n.EN: "Funeral expenses",
			i18n.ES: "Gastos funerales",
		},
	},
	{
		Key: TagUse.With(TagDonation),
		Name: i18n.String{
			i18n.EN: "Donations",
			i18n.ES: "Donativos",
		},
	},
	{
		Key: TagUse.With(TagMortgageInterest),
		Name: i18n.String{
			i18n.EN: "Interest actually paid on mortgage loans (housing)",
			i18n.ES: "Intereses reales efectivamente pagados por créditos hipotecarios (casa habitación)",
		},
	},
	{
		Key: TagUse.With(TagSARContribution),
		Name: i18n.String{
			i18n.EN: "Voluntary contributions to the SAR",
			i18n.ES: "Aportaciones voluntarias al SAR",
		},
	},
	{
		Key: TagUse.With(TagMedicalInsurance),
		Name: i18n.String{
			i18n.EN: "Medical insurance premiums",
			i18n.ES: "Primas por seguros de gastos médicos",
		},
	},
	{
		Key: TagUse.With(TagSchoolTransportation),
		Name: i18n.String{
			i18n.EN: "Mandatory school transportation expenses",
			i18n.ES: "Gastos de transportación escolar obligatoria",
		},
	},
	{
		Key: TagUse.With(TagSavingsDeposit),
		Name: i18n.String{
			i18n.EN: "Deposits in savings accounts, pension plans premiums",
			i18n.ES: "Depósitos en cuentas para el ahorro, primas que tengan como base planes de pensiones",
		},
	},
	{
		Key: TagUse.With(TagSchoolFees),
		Name: i18n.String{
			i18n.EN: "Payments for educational services (school fees)",
			i18n.ES: "Pagos por servicios educativos (colegiaturas)",
		},
	},
	{
		Key: TagUse.With(TagNoTaxEffects),
		Name: i18n.String{
			i18n.EN: "Without tax effects",
			i18n.ES: "Sin efectos fiscales",
		},
	},
	{
		Key: TagUse.With(TagSuplementaryPayment),
		Name: i18n.String{
			i18n.EN: "Payments",
			i18n.ES: "Pagos",
		},
	},
	{
		Key: TagUse.With(TagPayroll),
		Name: i18n.String{
			i18n.EN: "Payroll",
			i18n.ES: "Nómina",
		},
	},
}

var scenarios = []*tax.ScenarioSet{
	invoiceScenarios,
}

var invoiceScenarios = &tax.ScenarioSet{
	Schema: bill.ShortSchemaInvoice,
	List: []*tax.Scenario{
		// TipoDeComprobante / TipoRelacion
		{
			Types: []cbc.Key{bill.InvoiceTypeStandard},
			Codes: cbc.CodeSet{
				KeySATTipoDeComprobante: "I",
			},
		},
		{
			Types: []cbc.Key{bill.InvoiceTypeCreditNote},
			Codes: cbc.CodeSet{
				KeySATTipoDeComprobante: "E",
				KeySATTipoRelacion:      "01",
			},
		},

		// UsoCFDI
		{
			Tags: []cbc.Key{TagUse.With(TagGoodsAcquisition)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "G01",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagReturns)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "G02",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagGeneralExpenses)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "G03",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagConstruction)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "I01",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagOfficeEquipment)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "I02",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagTransportEquipment)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "I03",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagComputerEquipment)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "I04",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagManufacturingTooling)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "I05",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagTelephoneComms)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "I06",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagSatelliteComms)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "I07",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagOtherMachinery)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "I08",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagMedicalExpenses)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "D01",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagMedicalExpenses).With(TagDisability)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "D02",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagFuneralExpenses)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "D03",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagDonation)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "D04",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagMortgageInterest)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "D05",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagSARContribution)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "D06",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagMedicalInsurance)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "D07",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagSchoolTransportation)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "D08",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagSavingsDeposit)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "D09",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagSchoolFees)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "D10",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagNoTaxEffects)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "S01",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagSuplementaryPayment)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "CP01",
			},
		},

		{
			Tags: []cbc.Key{TagUse.With(TagPayroll)},
			Codes: cbc.CodeSet{
				KeySATUsoCFDI: "CN01",
			},
		},
	},
}
