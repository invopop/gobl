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

var extDocumentTypes = &cbc.KeyDefinition{
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
	Values: []*cbc.ValueDefinition{
		{
			Value: "71",
			Name:  i18n.NewString("Request for payment"),
		},
		{
			Value: "80",
			Name:  i18n.NewString("Debit note related to goods or services"),
		},
		{
			Value: "81",
			Name:  i18n.NewString("Credit note related to goods or services"),
		},
		{
			Value: "82",
			Name:  i18n.NewString("Metered services invoice"),
		},
		{
			Value: "83",
			Name:  i18n.NewString("Credit note related to financial adjustments"),
		},
		{
			Value: "84",
			Name:  i18n.NewString("Debit note related to financial adjustments"),
		},
		{
			Value: "102",
			Name:  i18n.NewString("Tax notification"),
		},
		{
			Value: "130",
			Name:  i18n.NewString("Invoicing data sheet"),
		},
		{
			Value: "202",
			Name:  i18n.NewString("Direct payment valuation"),
		},
		{
			Value: "203",
			Name:  i18n.NewString("Provisional payment valuation"),
		},
		{
			Value: "204",
			Name:  i18n.NewString("Payment valuation"),
		},
		{
			Value: "211",
			Name:  i18n.NewString("Interim application for payment"),
		},
		{
			Value: "218",
			Name:  i18n.NewString("Final payment request based on completion of work"),
		},
		{
			Value: "219",
			Name:  i18n.NewString("Payment request for completed units"),
		},
		{
			Value: "261",
			Name:  i18n.NewString("Self billed credit note"),
		},
		{
			Value: "262",
			Name:  i18n.NewString("Consolidated credit note - goods and services"),
		},
		{
			Value: "295",
			Name:  i18n.NewString("Price variation invoice"),
		},
		{
			Value: "296",
			Name:  i18n.NewString("Credit note for price variation"),
		},
		{
			Value: "308",
			Name:  i18n.NewString("Delcredere credit note"),
		},
		{
			Value: "325",
			Name:  i18n.NewString("Proforma invoice"),
		},
		{
			Value: "326",
			Name:  i18n.NewString("Partial invoice"),
		},
		{
			Value: "380",
			Name:  i18n.NewString("Standard Invoice"),
		},
		{
			Value: "381",
			Name:  i18n.NewString("Credit note"),
		},
		{
			Value: "382",
			Name:  i18n.NewString("Commission note"),
		},
		{
			Value: "383",
			Name:  i18n.NewString("Debit note"),
		},
		{
			Value: "384",
			Name:  i18n.NewString("Corrected invoice"),
		},
		{
			Value: "385",
			Name:  i18n.NewString("Consolidated invoice"),
		},
		{
			Value: "386",
			Name:  i18n.NewString("Prepayment invoice"),
		},
		{
			Value: "387",
			Name:  i18n.NewString("Hire invoice"),
		},
		{
			Value: "388",
			Name:  i18n.NewString("Tax invoice"),
		},
		{
			Value: "389",
			Name:  i18n.NewString("Self-billed invoice"),
		},
		{
			Value: "390",
			Name:  i18n.NewString("Delcredere invoice"),
		},
		{
			Value: "393",
			Name:  i18n.NewString("Factored invoice"),
		},
		{
			Value: "394",
			Name:  i18n.NewString("Lease invoice"),
		},
		{
			Value: "395",
			Name:  i18n.NewString("Consignment invoice"),
		},
		{
			Value: "396",
			Name:  i18n.NewString("Factored credit note"),
		},
		{
			Value: "420",
			Name:  i18n.NewString("Optical Character Reading (OCR) payment credit note"),
		},
		{
			Value: "456",
			Name:  i18n.NewString("Debit advice"),
		},
		{
			Value: "457",
			Name:  i18n.NewString("Reversal of debit"),
		},
		{
			Value: "458",
			Name:  i18n.NewString("Reversal of credit"),
		},
		{
			Value: "527",
			Name:  i18n.NewString("Self billed debit note"),
		},
		{
			Value: "532",
			Name:  i18n.NewString("Forwarder's credit note"),
		},
		{
			Value: "553",
			Name:  i18n.NewString("Forwarder's invoice discrepancy report"),
		},
		{
			Value: "575",
			Name:  i18n.NewString("Insurer's invoice"),
		},
		{
			Value: "623",
			Name:  i18n.NewString("Forwarder's invoice"),
		},
		{
			Value: "633",
			Name:  i18n.NewString("Port charges documents"),
		},
		{
			Value: "751",
			Name:  i18n.NewString("Invoice information for accounting purposes"),
		},
		{
			Value: "780",
			Name:  i18n.NewString("Freight invoice"),
		},
		{
			Value: "817",
			Name:  i18n.NewString("Claim notification"),
		},
		{
			Value: "870",
			Name:  i18n.NewString("Consular invoice"),
		},
		{
			Value: "875",
			Name:  i18n.NewString("Partial construction invoice"),
		},
		{
			Value: "876",
			Name:  i18n.NewString("Partial final construction invoice"),
		},
		{
			Value: "877",
			Name:  i18n.NewString("Final construction invoice"),
		},
		{
			Value: "935",
			Name:  i18n.NewString("Customs invoice"),
		},
	},
}
