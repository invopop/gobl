package zatca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/cef"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

// UNTDID 1001 document type codes accepted by ZATCA.
const (
	docTypeTaxInvoice cbc.Code = "388"
	docTypePrepayment cbc.Code = "386"
	docTypeDebitNote  cbc.Code = "383"
	docTypeCreditNote cbc.Code = "381"
)

// Other valid identities used by ZATCA
const (
	IdentityTypeTIN      cbc.Code = "TIN"
	IdentityTypeCRN      cbc.Code = "CRN"
	IdentityTypeMom      cbc.Code = "MOM"
	IdentityTypeMLS      cbc.Code = "MLS"
	IdentityType700      cbc.Code = "700"
	IdentityTypeSAG      cbc.Code = "SAG"
	IdentityTypeNational cbc.Code = "NAT"
	IdentityTypeGcc      cbc.Code = "GCC"
	IdentityTypeIqa      cbc.Code = "IQA"
	IdentityTypePassport cbc.Code = "PAS"
	IdentityTypeOTH      cbc.Code = "OTH"
)

var (
	customerValidIdentities = []cbc.Code{
		IdentityTypeTIN,
		IdentityTypeCRN,
		IdentityTypeMom,
		IdentityTypeMLS,
		IdentityType700,
		IdentityTypeSAG,
		IdentityTypeNational,
		IdentityTypeGcc,
		IdentityTypeIqa,
		IdentityTypePassport,
		IdentityTypeOTH,
	}

	// SupplierValidIdentities holds ZATCA specific supplier identitites
	SupplierValidIdentities = []cbc.Code{
		IdentityTypeCRN,
		IdentityTypeMom,
		IdentityTypeMLS,
		IdentityType700,
		IdentityTypeSAG,
		IdentityTypeOTH,
	}
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
				rules.Assert("03", "untdid document type extension is required",
					tax.ExtensionsRequire(untdid.ExtKeyDocumentType),
				),
				rules.Assert("04", "invoice transaction type extension is required",
					tax.ExtensionsRequire(ExtKeyInvoiceTypeTransactions),
				),
				rules.Assert("05", "document type must be a valid ZATCA type (388, 386, 383, 381) (BR-KSA-05)",
					tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, docTypeTaxInvoice, docTypePrepayment, docTypeDebitNote, docTypeCreditNote),
				),
				rules.Assert("06", "invoice transaction type must be valid (BR-KSA-06)",
					tax.ExtensionsHasCodes(ExtKeyInvoiceTypeTransactions, validTransactionTypes...),
				),
			),
		),

		// Credit or debit note
		rules.When(
			is.Func("credit or debit note", invoiceIsCreditOrDebitNote),
			rules.Field("preceding",
				rules.Assert("07", "credit and debit notes must have a billing reference", is.Present),
				rules.Each(
					rules.Field("code",
						rules.Assert("08", "billing reference must have an identifier (BR-KSA-56)", is.Present),
					),
					rules.Field("reason",
						rules.Assert("09", "credit and debit notes must contain the reason for issuance (BR-KSA-17)",
							is.Present,
						),
					),
				),
			),
		),

		// Standard
		rules.When(
			is.Func("standard tax invoice", invoiceIsStandard),
			rules.Field("customer",
				rules.Field("addresses",
					rules.Each(
						rules.Field("street",
							rules.Assert("10", "customer address must have a street name (BR-KSA-10)", is.Present),
						),
						rules.Field("locality",
							rules.Assert("11", "customer address must have a city name (BR-KSA-10)", is.Present),
						),
						rules.Field("country",
							rules.Assert("12", "customer address must have a country code (BR-KSA-10)", is.Present),
						),
					),
				),
				rules.Assert("13", "customer must have a valid identification scheme for standard invoices",
					is.Func("customer must be either VAT registered or have a valid identification (BR-KSA-14), (BR-KSA-81)", customerValidIdentity),
				),
			),
			rules.Field("lines",
				rules.Each(
					rules.Field("taxes",
						rules.Assert("14", "line taxes are required for standard tax invoices and associated credit notes and debit notes (BR-KSA-52)", is.Present),
					),
				),
			),
			rules.Field("delivery",
				rules.Assert("15", "delivery must be present", is.Present),
				rules.Field("date",
					rules.Assert("16", "delivery must have a supply date (BR-KSA-15)", is.Present),
				),
			),
		),

		// Export invoice
		rules.When(
			is.Func("export invoice", invoiceIsExport),
			rules.Field("customer",
				rules.Field("tax_id",
					rules.Assert("17", "export invoices must not have buyer VAT registration number (BR-KSA-46)",
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
				rules.Assert("18", "delivery must be present for simplified and summary invoices", is.Present),
				rules.Field("period",
					rules.Assert("19", "supply must have a delivery period", is.Present),
					rules.Field("start",
						rules.Assert("20", "delivery start date must be present (BR-KSA-72)", is.Present),
					),
					rules.Field("end",
						rules.Assert("21", "delivery end date must be present (BR-KSA-72)", is.Present),
					),
				),
			),
		),

		// EDU or HEA exemptions
		rules.When(
			is.Func("has EDU or HEA tax exemption", invoiceHasEDUOrHEAExemption),
			rules.Field("customer",
				rules.Field("identities",
					rules.Assert("22", "customer must have a national ID (NAT) when tax exemption is VATEX-SA-EDU or VATEX-SA-HEA (BR-KSA-49)",
						org.IdentitiesTypeIn(IdentityTypeNational),
					),
				),
			),
		),

		// Customer name
		rules.When(
			is.Or(
				is.Func("simplified and (EDU or HEA exemptions)", invoiceIsSimplifiedAndEDUOrHEAExemption),
				is.Func("invoice is simplified and summary", invoiceIsSimplifiedAndSummary),
				is.Func("standard tax invoice", invoiceIsStandard),
			),
			rules.Field("customer",
				rules.Assert("23", "customer must be present", is.Present),
				rules.Field("name",
					rules.Assert("24", "customer name must be present (BR-KSA-71), (BR-KSA-25), (BR-KSA-42)", is.Present),
				),
			),
		),

		// Supplier
		rules.Field("supplier",
			rules.Assert("25", "supplier must have a tax ID code (BR-KSA-39)",
				is.Func("valid VAT code", hasTaxIDCode),
			),
			rules.Field("identities",
				rules.Assert("26", "supplier must have a valid identity (BR-KSA-08)",
					is.Func("identity must be one of: CRN/MOM/MLS/700/SAG/OTH", hasOneSupplierIdentity),
				),
			),
		),

		// Self-billed
		rules.When(
			is.Func("invoice is self-billed", invoiceIsSelfBilled),
			rules.Field("customer",
				rules.Assert("27", "customer must have a tax ID code (BR-KSA-39)",
					is.Func("valid VAT code", hasTaxIDCode),
				),
				rules.Field("identities",
					rules.Assert("28", "customer must have a valid identity (BR-KSA-08)",
						is.Func("identity must be one of: CRN/MOM/MLS/700/SAG/OTH", hasOneSupplierIdentity),
					),
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
	if inv.Tax == nil {
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
	return code != "" && code[5] == '1'
}

func invoiceIsSelfBilled(val any) bool {
	code := getInvTypeCode(val)
	return code != "" && code[6] == '1'
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

func hasTaxIDCode(value any) bool {
	party, _ := value.(*org.Party)
	return party != nil && party.TaxID != nil && party.TaxID.Code != ""
}

func hasOneSupplierIdentity(value any) bool {
	identities, _ := value.([]*org.Identity)
	return len(identities) == 1 && org.IdentitiesTypeIn(SupplierValidIdentities...).Check(identities)
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
	return len(party.Identities) == 1 && org.IdentitiesTypeIn(customerValidIdentities...).Check(party.Identities)
}
