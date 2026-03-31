package zatca

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// ZATCA extension keys
const (
	// ExtKeyInvoiceType identifies the ZATCA invoice subtype code (KSA-2).
	// This is a 7-character binary code where:
	//   - Positions 1-2: Invoice type (01 = Standard Tax Invoice, 02 = Simplified Tax Invoice)
	//   - Position 3: Third-party transaction (0 or 1)
	//   - Position 4: Nominal supply transaction (0 or 1)
	//   - Position 5: Export invoice (0 or 1)
	//   - Position 6: Summary invoice (0 or 1)
	//   - Position 7: Self-billed invoice (0 or 1)
	ExtKeyInvoiceType cbc.Key = "sa-zatca-invoice-type"

	// ExtKeySellerIDScheme identifies the ZATCA seller identification scheme (BT-29-1).
	ExtKeySellerIDScheme cbc.Key = "sa-zatca-seller-id-scheme"

	// ExtKeyBuyerIDScheme identifies the ZATCA buyer identification scheme (BT-46-1).
	ExtKeyBuyerIDScheme cbc.Key = "sa-zatca-buyer-id-scheme"
)

var allowedDocumentTypes = []cbc.Code{
	"388", // Tax Invoice
	"386", // Prepayment Invoice
	"383", // Debit Note
	"381", // Credit Note
}

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

var InvTypesSummary = []cbc.Code{
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

// Valid seller identification scheme IDs (BR-KSA-08)
var sellerIDSchemes = []cbc.Code{
	"CRN", // Commercial Registration Number
	"MOM", // MOMRAH license
	"MLS", // MHRSD license
	"700", // 700 Number
	"SAG", // MISA license
	"OTH", // Other ID
}

// Valid buyer identification scheme IDs (BR-KSA-14)
var buyerIDSchemes = []cbc.Code{
	"TIN", // Tax Identification Number
	"CRN", // Commercial Registration Number
	"MOM", // MOMRAH license
	"MLS", // MHRSD license
	"700", // 700 Number
	"SAG", // MISA license
	"NAT", // National ID
	"GCC", // GCC ID
	"IQA", // Iqama Number
	"PAS", // Passport ID
	"OTH", // Other ID
}

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyInvoiceType,
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
	{
		Key: ExtKeySellerIDScheme,
		Name: i18n.String{
			i18n.EN: "ZATCA Seller Identification Scheme",
			i18n.AR: "نظام تعريف البائع",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Scheme ID for the seller's additional identification (BT-29-1).
				One of: CRN (Commercial Registration), MOM (MOMRAH license),
				MLS (MHRSD license), 700 (700 Number), SAG (MISA license),
				OTH (Other ID).
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "CRN",
				Name: i18n.String{
					i18n.EN: "Commercial Registration Number",
					i18n.AR: "رقم السجل التجاري",
				},
			},
			{
				Code: "MOM",
				Name: i18n.String{
					i18n.EN: "MOMRAH License",
					i18n.AR: "ترخيص وزارة الشؤون البلدية",
				},
			},
			{
				Code: "MLS",
				Name: i18n.String{
					i18n.EN: "MHRSD License",
					i18n.AR: "ترخيص وزارة الموارد البشرية",
				},
			},
			{
				Code: "700",
				Name: i18n.String{
					i18n.EN: "700 Number",
					i18n.AR: "رقم 700",
				},
			},
			{
				Code: "SAG",
				Name: i18n.String{
					i18n.EN: "MISA License",
					i18n.AR: "ترخيص وزارة الاستثمار",
				},
			},
			{
				Code: "OTH",
				Name: i18n.String{
					i18n.EN: "Other ID",
					i18n.AR: "معرف آخر",
				},
			},
		},
	},
	{
		Key: ExtKeyBuyerIDScheme,
		Name: i18n.String{
			i18n.EN: "ZATCA Buyer Identification Scheme",
			i18n.AR: "نظام تعريف المشتري",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Scheme ID for the buyer's additional identification (BT-46-1).
				One of: TIN, CRN, MOM, MLS, 700, SAG, NAT, GCC, IQA, PAS, OTH.
			`),
		},
		Values: []*cbc.Definition{
			{
				Code: "TIN",
				Name: i18n.String{
					i18n.EN: "Tax Identification Number",
					i18n.AR: "الرقم الضريبي",
				},
			},
			{
				Code: "CRN",
				Name: i18n.String{
					i18n.EN: "Commercial Registration Number",
					i18n.AR: "رقم السجل التجاري",
				},
			},
			{
				Code: "MOM",
				Name: i18n.String{
					i18n.EN: "MOMRAH License",
					i18n.AR: "ترخيص وزارة الشؤون البلدية",
				},
			},
			{
				Code: "MLS",
				Name: i18n.String{
					i18n.EN: "MHRSD License",
					i18n.AR: "ترخيص وزارة الموارد البشرية",
				},
			},
			{
				Code: "700",
				Name: i18n.String{
					i18n.EN: "700 Number",
					i18n.AR: "رقم 700",
				},
			},
			{
				Code: "SAG",
				Name: i18n.String{
					i18n.EN: "MISA License",
					i18n.AR: "ترخيص وزارة الاستثمار",
				},
			},
			{
				Code: "NAT",
				Name: i18n.String{
					i18n.EN: "National ID",
					i18n.AR: "الهوية الوطنية",
				},
			},
			{
				Code: "GCC",
				Name: i18n.String{
					i18n.EN: "GCC ID",
					i18n.AR: "هوية مجلس التعاون",
				},
			},
			{
				Code: "IQA",
				Name: i18n.String{
					i18n.EN: "Iqama Number",
					i18n.AR: "رقم الإقامة",
				},
			},
			{
				Code: "PAS",
				Name: i18n.String{
					i18n.EN: "Passport ID",
					i18n.AR: "رقم جواز السفر",
				},
			},
			{
				Code: "OTH",
				Name: i18n.String{
					i18n.EN: "Other ID",
					i18n.AR: "معرف آخر",
				},
			},
		},
	},
}
