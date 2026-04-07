package zatca

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// ZATCA extension keys
const (
	// KSA-2 invoice type code positions (0-indexed).
	// The code is a 7-character string: TTXNESO
	//   - TT (0-1): Invoice type (01=Standard, 02=Simplified)
	//   - X  (2):   Third-party transaction
	//   - N  (3):   Nominal supply transaction
	//   - E  (4):   Export invoice
	//   - S  (5):   Summary invoice
	//   - O  (6):   Self-billed invoice
	InvTypeCodeLen   = 7
	InvTypePosExport = 4

	// ExtKeyInvoiceType identifies the ZATCA invoice subtype code (KSA-2).
	ExtKeyInvoiceTypeTransactions cbc.Key = "sa-zatca-invoice-type"

	// SA VATEX exemption reason codes
	Vatex29         cbc.Code = "VATEX-SA-29"    // Exempt: Financial services (Article 29)
	Vatex29_7       cbc.Code = "VATEX-SA-29-7"  // Exempt: Life insurance (Article 29.7)
	Vatex30         cbc.Code = "VATEX-SA-30"    // Exempt: Real estate (Article 30)
	Vatex32         cbc.Code = "VATEX-SA-32"    // Zero-rated: Export of goods (Article 32)
	Vatex33         cbc.Code = "VATEX-SA-33"    // Zero-rated: Export of services (Article 33)
	Vatex34_1       cbc.Code = "VATEX-SA-34-1"  // Zero-rated: Intra-GCC supply (Article 34.1)
	Vatex34_2       cbc.Code = "VATEX-SA-34-2"  // Zero-rated: Intra-GCC supply (Article 34.2)
	Vatex34_3       cbc.Code = "VATEX-SA-34-3"  // Zero-rated: Intra-GCC supply (Article 34.3)
	Vatex34_4       cbc.Code = "VATEX-SA-34-4"  // Zero-rated: Intra-GCC supply (Article 34.4)
	Vatex34_5       cbc.Code = "VATEX-SA-34-5"  // Zero-rated: Intra-GCC supply (Article 34.5)
	Vatex35         cbc.Code = "VATEX-SA-35"    // Zero-rated: Qualified medicines and medical equipment (Article 35)
	Vatex36         cbc.Code = "VATEX-SA-36"    // Zero-rated: Qualified metals (Article 36)
	VatexEdu        cbc.Code = "VATEX-SA-EDU"   // Zero-rated: Private education
	VatexHea        cbc.Code = "VATEX-SA-HEA"   // Zero-rated: Private healthcare
	VatexMltry      cbc.Code = "VATEX-SA-MLTRY" // Zero-rated: Qualified military goods
	VatexOutOfScope cbc.Code = "VATEX-SA-OOS"   // Out of scope
)

var vatIDPattern = `^3[0-9]{13}3$`

var allowedDocumentTypes = []cbc.Code{
	"388", // Tax Invoice
	"386", // Prepayment Invoice
	"383", // Debit Note
	"381", // Credit Note
}

// InvTypesStandard contains all valid standard tax invoice type codes (KSA-2 starting with 01).
var InvTypesStandard = []cbc.Code{
	"0100000", // Standard Tax Invoice
	"0100001", // Standard Tax Invoice — Self-billed
	"0100010", // Standard Tax Invoice — Summary
	"0100011", // Standard Tax Invoice — Summary, Self-billed
	"0100100", // Standard Tax Invoice — Export
	"0100110", // Standard Tax Invoice — Export, Summary
	"0101000", // Standard Tax Invoice — Nominal
	"0101001", // Standard Tax Invoice — Nominal, Self-billed
	"0101010", // Standard Tax Invoice — Nominal, Summary
	"0101011", // Standard Tax Invoice — Nominal, Summary, Self-billed
	"0101100", // Standard Tax Invoice — Nominal, Export
	"0101110", // Standard Tax Invoice — Nominal, Export, Summary
	"0110000", // Standard Tax Invoice — Third-party
	"0110001", // Standard Tax Invoice — Third-party, Self-billed
	"0110010", // Standard Tax Invoice — Third-party, Summary
	"0110011", // Standard Tax Invoice — Third-party, Summary, Self-billed
	"0110100", // Standard Tax Invoice — Third-party, Export
	"0110110", // Standard Tax Invoice — Third-party, Export, Summary
	"0111000", // Standard Tax Invoice — Third-party, Nominal
	"0111001", // Standard Tax Invoice — Third-party, Nominal, Self-billed
	"0111010", // Standard Tax Invoice — Third-party, Nominal, Summary
	"0111011", // Standard Tax Invoice — Third-party, Nominal, Summary, Self-billed
	"0111100", // Standard Tax Invoice — Third-party, Nominal, Export
	"0111110", // Standard Tax Invoice — Third-party, Nominal, Export, Summary
}

// InvTypesSimplified contains all valid simplified tax invoice type codes (KSA-2 starting with 02).
var InvTypesSimplified = []cbc.Code{
	"0200000", // Simplified Tax Invoice
	"0200010", // Simplified Tax Invoice — Summary
	"0201000", // Simplified Tax Invoice — Nominal
	"0201010", // Simplified Tax Invoice — Nominal, Summary
	"0210000", // Simplified Tax Invoice — Third-party
	"0210010", // Simplified Tax Invoice — Third-party, Summary
	"0211000", // Simplified Tax Invoice — Third-party, Nominal
	"0211010", // Simplified Tax Invoice — Third-party, Nominal, Summary
}

var invTypesSummary = []cbc.Code{
	"0100010", // Standard Tax Invoice — Summary
	"0100011", // Standard Tax Invoice — Summary, Self-billed
	"0100110", // Standard Tax Invoice — Export, Summary
	"0101010", // Standard Tax Invoice — Nominal, Summary
	"0101011", // Standard Tax Invoice — Nominal, Summary, Self-billed
	"0101110", // Standard Tax Invoice — Nominal, Export, Summary
	"0110010", // Standard Tax Invoice — Third-party, Summary
	"0110011", // Standard Tax Invoice — Third-party, Summary, Self-billed
	"0110110", // Standard Tax Invoice — Third-party, Export, Summary
	"0111010", // Standard Tax Invoice — Third-party, Nominal, Summary
	"0111011", // Standard Tax Invoice — Third-party, Nominal, Summary, Self-billed
	"0111110", // Standard Tax Invoice — Third-party, Nominal, Export, Summary
	"0200010", // Simplified Tax Invoice — Summary
	"0201010", // Simplified Tax Invoice — Nominal, Summary
	"0210010", // Simplified Tax Invoice — Third-party, Summary
	"0211010", // Simplified Tax Invoice — Third-party, Nominal, Summary
}

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyInvoiceTypeTransactions,
		Name: i18n.String{
			i18n.EN: "ZATCA Invoice Type",
			i18n.AR: "نوع الفاتورة",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Code used to describe the ZATCA invoice subtype and transaction flags (KSA-2).
				The code is a 7-character string where positions 1-2 indicate the main type
				(01 = Standard Tax Invoice, 02 = Simplified Tax Invoice) and positions 3-7
				are binary flags for third-party, nominal, export, summary, and self-billed
				transactions respectively.
			`),
		},
		Pattern: `^0[12][01]{5}$`,
		Values: []*cbc.Definition{
			{
				Code: "0100000",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice",
					i18n.AR: "فاتورة ضريبية",
				},
			},
			{
				Code: "0100001",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Self-billed",
					i18n.AR: "فاتورة ضريبية — ذاتية الإصدار",
				},
			},
			{
				Code: "0100010",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Summary",
					i18n.AR: "فاتورة ضريبية — ملخص",
				},
			},
			{
				Code: "0100011",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Summary, Self-billed",
					i18n.AR: "فاتورة ضريبية — ملخص، ذاتية الإصدار",
				},
			},
			{
				Code: "0100100",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Export",
					i18n.AR: "فاتورة ضريبية — تصدير",
				},
			},
			{
				Code: "0100110",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Export, Summary",
					i18n.AR: "فاتورة ضريبية — تصدير، ملخص",
				},
			},
			{
				Code: "0101000",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Nominal",
					i18n.AR: "فاتورة ضريبية — اسمية",
				},
			},
			{
				Code: "0101001",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Nominal, Self-billed",
					i18n.AR: "فاتورة ضريبية — اسمية، ذاتية الإصدار",
				},
			},
			{
				Code: "0101010",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Nominal, Summary",
					i18n.AR: "فاتورة ضريبية — اسمية، ملخص",
				},
			},
			{
				Code: "0101011",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Nominal, Summary, Self-billed",
					i18n.AR: "فاتورة ضريبية — اسمية، ملخص، ذاتية الإصدار",
				},
			},
			{
				Code: "0101100",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Nominal, Export",
					i18n.AR: "فاتورة ضريبية — اسمية، تصدير",
				},
			},
			{
				Code: "0101110",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Nominal, Export, Summary",
					i18n.AR: "فاتورة ضريبية — اسمية، تصدير، ملخص",
				},
			},
			{
				Code: "0110000",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Third-party",
					i18n.AR: "فاتورة ضريبية — طرف ثالث",
				},
			},
			{
				Code: "0110001",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Third-party, Self-billed",
					i18n.AR: "فاتورة ضريبية — طرف ثالث، ذاتية الإصدار",
				},
			},
			{
				Code: "0110010",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Third-party, Summary",
					i18n.AR: "فاتورة ضريبية — طرف ثالث، ملخص",
				},
			},
			{
				Code: "0110011",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Third-party, Summary, Self-billed",
					i18n.AR: "فاتورة ضريبية — طرف ثالث، ملخص، ذاتية الإصدار",
				},
			},
			{
				Code: "0110100",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Third-party, Export",
					i18n.AR: "فاتورة ضريبية — طرف ثالث، تصدير",
				},
			},
			{
				Code: "0110110",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Third-party, Export, Summary",
					i18n.AR: "فاتورة ضريبية — طرف ثالث، تصدير، ملخص",
				},
			},
			{
				Code: "0111000",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Third-party, Nominal",
					i18n.AR: "فاتورة ضريبية — طرف ثالث، اسمية",
				},
			},
			{
				Code: "0111001",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Third-party, Nominal, Self-billed",
					i18n.AR: "فاتورة ضريبية — طرف ثالث، اسمية، ذاتية الإصدار",
				},
			},
			{
				Code: "0111010",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Third-party, Nominal, Summary",
					i18n.AR: "فاتورة ضريبية — طرف ثالث، اسمية، ملخص",
				},
			},
			{
				Code: "0111011",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Third-party, Nominal, Summary, Self-billed",
					i18n.AR: "فاتورة ضريبية — طرف ثالث، اسمية، ملخص، ذاتية الإصدار",
				},
			},
			{
				Code: "0111100",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Third-party, Nominal, Export",
					i18n.AR: "فاتورة ضريبية — طرف ثالث، اسمية، تصدير",
				},
			},
			{
				Code: "0111110",
				Name: i18n.String{
					i18n.EN: "Standard Tax Invoice — Third-party, Nominal, Export, Summary",
					i18n.AR: "فاتورة ضريبية — طرف ثالث، اسمية، تصدير، ملخص",
				},
			},
			{
				Code: "0200000",
				Name: i18n.String{
					i18n.EN: "Simplified Tax Invoice",
					i18n.AR: "فاتورة ضريبية مبسطة",
				},
			},
			{
				Code: "0200010",
				Name: i18n.String{
					i18n.EN: "Simplified Tax Invoice — Summary",
					i18n.AR: "فاتورة ضريبية مبسطة — ملخص",
				},
			},
			{
				Code: "0201000",
				Name: i18n.String{
					i18n.EN: "Simplified Tax Invoice — Nominal",
					i18n.AR: "فاتورة ضريبية مبسطة — اسمية",
				},
			},
			{
				Code: "0201010",
				Name: i18n.String{
					i18n.EN: "Simplified Tax Invoice — Nominal, Summary",
					i18n.AR: "فاتورة ضريبية مبسطة — اسمية، ملخص",
				},
			},
			{
				Code: "0210000",
				Name: i18n.String{
					i18n.EN: "Simplified Tax Invoice — Third-party",
					i18n.AR: "فاتورة ضريبية مبسطة — طرف ثالث",
				},
			},
			{
				Code: "0210010",
				Name: i18n.String{
					i18n.EN: "Simplified Tax Invoice — Third-party, Summary",
					i18n.AR: "فاتورة ضريبية مبسطة — طرف ثالث، ملخص",
				},
			},
			{
				Code: "0211000",
				Name: i18n.String{
					i18n.EN: "Simplified Tax Invoice — Third-party, Nominal",
					i18n.AR: "فاتورة ضريبية مبسطة — طرف ثالث، اسمية",
				},
			},
			{
				Code: "0211010",
				Name: i18n.String{
					i18n.EN: "Simplified Tax Invoice — Third-party, Nominal, Summary",
					i18n.AR: "فاتورة ضريبية مبسطة — طرف ثالث، اسمية، ملخص",
				},
			},
		},
	},
}
