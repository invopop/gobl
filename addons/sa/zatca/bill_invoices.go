package zatca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/cef"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/sa"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// UNTDID 1001 document type codes accepted by ZATCA.
const (
	DocTypeTaxInvoice cbc.Code = "388"
	DocTypePrepayment cbc.Code = "386"
	DocTypeDebitNote  cbc.Code = "383"
	DocTypeCreditNote cbc.Code = "381"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),

		rules.Field("issue_time",
			rules.Assert("01", "issue time must be present (BR-KSA-70)", is.Present),
		),

		// Tax
		rules.Field("tax",
			rules.Assert("02", "tax must be present", is.Present),
			rules.Field("ext",
				rules.Assert("03", "extensions untdid key document type and key invoice type transaction are required",
					tax.ExtensionsRequire(requiredExtensions...)),
				rules.Assert("04", "document type must be a valid ZATCA type (388, 386, 383, 381) (BR-KSA-05)",
					tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, DocTypeTaxInvoice, DocTypePrepayment, DocTypeDebitNote, DocTypeCreditNote),
				),
				rules.Assert("05", "invoice transaction type must be valid",
					tax.ExtensionsHasCodes(ExtKeyInvoiceTypeTransactions, validTransactionTypes...),
				),
			),
		),

		// Credit or debit note
		rules.When(
			is.Func("credit or debit note", invoiceIsCreditOrDebitNote),
			rules.Field("preceding",
				rules.Assert("06", "credit and debit notes must have a billing reference", is.Present),
				rules.Each(
					rules.Field("code",
						rules.Assert("07", "billing reference must have an identifier (BR-KSA-56)", is.Present),
					),
					rules.Field("reason",
						rules.Assert("08", "credit and debit notes must contain the reason for issuance (BR-KSA-17)",
							is.Present,
						),
					),
				),
			),
		),

		// Supplier
		rules.Field("supplier",
			rules.Field("name",
				rules.Assert("09", "supplier must have a name (BR-06)", is.Present),
			),
			rules.Field("addresses",
				rules.Assert("10", "supplier must have an address", is.Present),
			),
		),

		// Standard
		rules.When(
			is.Func("standard tax invoice", invoiceIsStandard),
			rules.Field("customer",
				rules.Assert("11", "customer must be present", is.Present),
				rules.Field("name",
					rules.Assert("12", "customer name must be present in the standard tax invoice (BR-KSA-42)", is.Present),
				),
				rules.Field("addresses",
					rules.Each(
						rules.Field("street",
							rules.Assert("13", "customer address must have a street name (BR-KSA-10)", is.Present),
						),
						rules.Field("locality",
							rules.Assert("14", "customer address must have a city name (BR-KSA-10)", is.Present),
						),
						rules.Field("country",
							rules.Assert("15", "customer address must have a country code (BR-KSA-10)", is.Present),
						),
					),
				),
				rules.Assert("16", "customer must have a valid identification scheme for standard invoices",
					is.Func("customer must be either VAT registered or have a valid identification (BR-KSA-14), (BR-KSA-81)", customerValidIdentity),
				),
			),
			rules.Field("lines",
				rules.Each(
					rules.Field("taxes",
						rules.Assert("17", "line taxes are required for standard tax invoices and associated credit notes and debit notes (BR-KSA-52)", is.Present),
					),
					rules.Field("total",
						rules.Assert("18", "line total amount is required for standard tax invoices and associated credit notes and debit notes (BR-KSA-53)", is.Present),
					),
				),
			),
			rules.Field("delivery",
				rules.Assert("19", "delivery must be present", is.Present),
				rules.Field("date",
					rules.Assert("20", "delivery must have a supply date (BR-KSA-15)", is.Present),
				),
			),
		),

		// Export invoice
		rules.When(
			is.Func("export invoice", invoiceIsExport),
			rules.Field("customer",
				rules.Field("tax_id",
					rules.Assert("21", "export invoices must not have buyer VAT registration number (BR-KSA-46)",
						is.Empty,
					),
				),
			),
		),

		// Simplified and summary
		rules.When(
			is.Or(
				is.Func("invoice is simplified and summary", invoiceIsSimplifiedAndSummary),
			),
			rules.Field("delivery",
				rules.Assert("22", "delivery must be present for simplified and summary invoices", is.Present),
				rules.Field("period",
					rules.Assert("23", "supply must have a delivery period", is.Present),
					rules.Field("start",
						rules.Assert("24", "delivery start date must be present (BR-KSA-72)", is.Present),
					),
					rules.Field("end",
						rules.Assert("25", "delivery end date must be present (BR-KSA-72)", is.Present),
					),
				),
			),
			rules.Field("customer",
				rules.Assert("26", "customer must be present for simplified, summary invoices", is.Present),
				rules.Field("name",
					rules.Assert("27", "customer name must be present for simplified, summary invoices (BR-KSA-71)", is.Present),
				),
			),
		),

		// EDU or HEA exemptions
		rules.When(
			is.Func("has EDU or HEA tax exemption", invoiceHasEDUOrHEAExemption),
			rules.Field("customer",
				rules.Field("identities",
					rules.Assert("28", "customer must have a national ID (NAT) when tax exemption is VATEX-SA-EDU or VATEX-SA-HEA (BR-KSA-49)",
						org.IdentitiesTypeIn(sa.IdentityTypeNational),
					),
				),
			),
		),

		// Simplified and (EDU or HEA exemptions)
		rules.When(
			is.Func("simplified and (EDU or HEA exemptions)", invoiceIsSimplifiedAndEDUOrHEAExemption),
			rules.Field("customer",
				rules.Assert("29", "customer must be present for simplified or associated credit/debit note and EDU or HEA exemptions", is.Present),
				rules.Field("name",
					rules.Assert("30", "customer name must be present (BR-KSA-25)", is.Present),
				),
			),
		),
	)
}

func getInvTypeCode(val any) string {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return ""
	}
	code := inv.Tax.GetExt(ExtKeyInvoiceTypeTransactions).String()
	if len(code) != InvTypeCodeLen {
		return ""
	}
	return code
}

func invoiceIsStandard(val any) bool {
	code := getInvTypeCode(val)
	return code != "" && code[:2] == "01"
}

func invoiceIsExport(val any) bool {
	code := getInvTypeCode(val)
	return code != "" && code[4] == '1'
}

func invoiceIsSummary(val any) bool {
	code := getInvTypeCode(val)
	return code != "" && cbc.Code(code).In(invTypesSummary...)
}

func invoiceIsSimplifiedAndSummary(val any) bool {
	return invoiceIsSummary(val) && !invoiceIsStandard(val)
}

func invoiceIsCreditOrDebitNote(val any) bool {
	inv, ok := val.(*bill.Invoice)
	return ok && inv != nil && (inv.Type == bill.InvoiceTypeCreditNote || inv.Type == bill.InvoiceTypeDebitNote)
}

func invoiceHasEDUOrHEAExemption(val any) bool {
	return invoiceHasExemption(val, []cbc.Code{VatexPrivateEducation, VatexPrivateHealthcare})
}

func invoiceIsSimplifiedAndEDUOrHEAExemption(val any) bool {
	return invoiceHasEDUOrHEAExemption(val) && !invoiceIsStandard(val)
}

func invoiceHasExemption(val any, exemptions []cbc.Code) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	for _, line := range inv.Lines {
		vat := line.GetTaxes().Get(tax.CategoryVAT)
		if vat == nil {
			continue
		}
		code := vat.Ext.Get(cef.ExtKeyVATEX)
		if code.In(exemptions...) {
			return true
		}
	}
	return false
}

func customerValidIdentity(value any) bool {
	party, _ := value.(*org.Party)
	if party == nil {
		return false
	}
	if party.TaxID != nil && !party.TaxID.Code.IsEmpty() {
		return true
	}
	return len(party.Identities) == 1 && org.IdentitiesTypeIn(sa.CustomerValidIdentities...).Check(party.Identities)
}
