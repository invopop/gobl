package untdid

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

const (
	// ExtKeyDocumentType is used to identify the UNTDID 1001 document type code.
	ExtKeyDocumentType cbc.Key = "untdid-document-type"
)

var extDocumentTypes = &cbc.Definition{
	Key: ExtKeyDocumentType,
	Name: i18n.String{
		i18n.EN: "UNTDID 1001 Document Type",
	},
	Desc: i18n.String{
		i18n.EN: here.Doc(`
				UNTDID 1001 code used to describe the type of document. Ths list is based
				on the [EN16931 code list](https://ec.europa.eu/digital-building-blocks/sites/display/DIGITAL/Registry+of+supporting+artefacts+to+implement+EN16931#RegistryofsupportingartefactstoimplementEN16931-Codelists)
				values table which focusses on invoices and payments.

				Other tax regimes and addons may use their own subset of codes.
			`),
	},
	Values: []*cbc.Definition{
		{
			Code: "71",
			Name: i18n.NewString("Request for payment"),
		},
		{
			Code: "80",
			Name: i18n.NewString("Debit note related to goods or services"),
		},
		{
			Code: "81",
			Name: i18n.NewString("Credit note related to goods or services"),
		},
		{
			Code: "82",
			Name: i18n.NewString("Metered services invoice"),
		},
		{
			Code: "83",
			Name: i18n.NewString("Credit note related to financial adjustments"),
		},
		{
			Code: "84",
			Name: i18n.NewString("Debit note related to financial adjustments"),
		},
		{
			Code: "102",
			Name: i18n.NewString("Tax notification"),
		},
		{
			Code: "130",
			Name: i18n.NewString("Invoicing data sheet"),
		},
		{
			Code: "202",
			Name: i18n.NewString("Direct payment valuation"),
		},
		{
			Code: "203",
			Name: i18n.NewString("Provisional payment valuation"),
		},
		{
			Code: "204",
			Name: i18n.NewString("Payment valuation"),
		},
		{
			Code: "211",
			Name: i18n.NewString("Interim application for payment"),
		},
		{
			Code: "218",
			Name: i18n.NewString("Final payment request based on completion of work"),
		},
		{
			Code: "219",
			Name: i18n.NewString("Payment request for completed units"),
		},
		{
			Code: "261",
			Name: i18n.NewString("Self billed credit note"),
		},
		{
			Code: "262",
			Name: i18n.NewString("Consolidated credit note - goods and services"),
		},
		{
			Code: "295",
			Name: i18n.NewString("Price variation invoice"),
		},
		{
			Code: "296",
			Name: i18n.NewString("Credit note for price variation"),
		},
		{
			Code: "308",
			Name: i18n.NewString("Delcredere credit note"),
		},
		{
			Code: "325",
			Name: i18n.NewString("Proforma invoice"),
		},
		{
			Code: "326",
			Name: i18n.NewString("Partial invoice"),
		},
		{
			Code: "380",
			Name: i18n.NewString("Standard Invoice"),
		},
		{
			Code: "381",
			Name: i18n.NewString("Credit note"),
		},
		{
			Code: "382",
			Name: i18n.NewString("Commission note"),
		},
		{
			Code: "383",
			Name: i18n.NewString("Debit note"),
		},
		{
			Code: "384",
			Name: i18n.NewString("Corrected invoice"),
		},
		{
			Code: "385",
			Name: i18n.NewString("Consolidated invoice"),
		},
		{
			Code: "386",
			Name: i18n.NewString("Prepayment invoice"),
		},
		{
			Code: "387",
			Name: i18n.NewString("Hire invoice"),
		},
		{
			Code: "388",
			Name: i18n.NewString("Tax invoice"),
		},
		{
			Code: "389",
			Name: i18n.NewString("Self-billed invoice"),
		},
		{
			Code: "390",
			Name: i18n.NewString("Delcredere invoice"),
		},
		{
			Code: "393",
			Name: i18n.NewString("Factored invoice"),
		},
		{
			Code: "394",
			Name: i18n.NewString("Lease invoice"),
		},
		{
			Code: "395",
			Name: i18n.NewString("Consignment invoice"),
		},
		{
			Code: "396",
			Name: i18n.NewString("Factored credit note"),
		},
		{
			Code: "420",
			Name: i18n.NewString("Optical Character Reading (OCR) payment credit note"),
		},
		{
			Code: "456",
			Name: i18n.NewString("Debit advice"),
		},
		{
			Code: "457",
			Name: i18n.NewString("Reversal of debit"),
		},
		{
			Code: "458",
			Name: i18n.NewString("Reversal of credit"),
		},
		{
			Code: "527",
			Name: i18n.NewString("Self billed debit note"),
		},
		{
			Code: "532",
			Name: i18n.NewString("Forwarder's credit note"),
		},
		{
			Code: "553",
			Name: i18n.NewString("Forwarder's invoice discrepancy report"),
		},
		{
			Code: "575",
			Name: i18n.NewString("Insurer's invoice"),
		},
		{
			Code: "623",
			Name: i18n.NewString("Forwarder's invoice"),
		},
		{
			Code: "633",
			Name: i18n.NewString("Port charges documents"),
		},
		{
			Code: "751",
			Name: i18n.NewString("Invoice information for accounting purposes"),
		},
		{
			Code: "780",
			Name: i18n.NewString("Freight invoice"),
		},
		{
			Code: "817",
			Name: i18n.NewString("Claim notification"),
		},
		{
			Code: "870",
			Name: i18n.NewString("Consular invoice"),
		},
		{
			Code: "875",
			Name: i18n.NewString("Partial construction invoice"),
		},
		{
			Code: "876",
			Name: i18n.NewString("Partial final construction invoice"),
		},
		{
			Code: "877",
			Name: i18n.NewString("Final construction invoice"),
		},
		{
			Code: "935",
			Name: i18n.NewString("Customs invoice"),
		},
	},
}
