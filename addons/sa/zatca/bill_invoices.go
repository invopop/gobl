package zatca

import (
	"regexp"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/cef"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),

		rules.Field("issue_time",
			rules.Assert("01", "issue time must be present (BR-KSA-70)", is.Present),
		),

		rules.Field("tax",
			rules.Field("ext",
				rules.Assert("02", "document type must be a valid ZATCA type (388, 386, 383, 381) (BR-KSA-05)",
					tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, allowedDocumentTypes...),
				),
			),
		),

		// Credit or debit note
		rules.When(
			is.Func("credit or debit note", invoiceIsCreditOrDebitNote),
			rules.Field("preceding",
				rules.Assert("03", "credit and debit notes must have a billing reference", is.Present),
				rules.Each(
					rules.Field("code",
						rules.Assert("04", "billing reference must have an identifier (BR-KSA-56)", is.Present),
					),
					rules.Field("reason",
						rules.Assert("05", "credit and debit notes must contain the reason for issuance (BR-KSA-17)",
							is.Present,
						),
					),
				),
			),
		),

		// Supplier
		rules.Field("supplier",
			rules.Assert("06", "supplier is required", is.Present),
			rules.Field("addresses",
				rules.Assert("07", "supplier address is required", is.Present),
				rules.Each(
					rules.Field("num",
						rules.Assert("08", "supplier address building number is required (BR-KSA-09)", is.Present),
						rules.Assert("09", "supplier address building number must contain 4 digits (BR-KSA-37)", is.Matches(`^\d{4}$`)),
					),
					// mapped to district
					rules.Field("street_extra",
						rules.Assert("10", "supplier address must have a district (BR-KSA-09)", is.Present),
					),
					rules.Field("street",
						rules.Assert("11", "supplier address must have a street name (BR-KSA-09)", is.Present),
					),
					rules.Field("code",
						rules.Assert("12", "supplier postal code is required (BR-KSA-09)", is.Present),
						rules.Assert("13", "supplier postal code must be 5 digits (BR-KSA-66)", is.Matches(`^\d{5}$`)),
					),
					rules.Field("locality",
						rules.Assert("14", "supplier address must have a city name (BR-KSA-09)", is.Present),
					),
					rules.Field("country",
						rules.Assert("15", "supplier address must have a country code (BR-KSA-09)", is.Present),
					),
				),
			),

			rules.Field("tax_id",
				rules.Field("code",
					rules.Assert("16", "supplier must have a VAT number (BR-KSA-39)", is.Present),
					rules.Assert("17", "supplier VAT number must be 15 digits starting/ending with 3 (BR-KSA-40)",
						is.Matches(vatIDPattern),
					),
				),
			),

			rules.Assert("18", "supplier identification must be valid",
				is.Func("supplier additional identification must be at most one from: CRN, MOM, MLS, 700, SAG, OTH (BR-KSA-08)", partyHasSingleIdentity),
			),
		),

		// Standard
		rules.When(
			is.Func("standard tax invoice", invoiceIsStandard),
			rules.Field("customer",
				rules.Assert("19", "customer must be present", is.Present),
				rules.Field("name",
					rules.Assert("20", "customer name must be present in the standard tax invoice and associated credit notes and debit notes (BR-KSA-42)", is.Present),
				),
				rules.Field("addresses",
					rules.Each(
						rules.Field("street",
							rules.Assert("21", "customer address must have a street name (BR-KSA-10)", is.Present),
						),
						rules.Field("locality",
							rules.Assert("22", "customer address must have a city name (BR-KSA-10)", is.Present),
						),
						rules.Field("country",
							rules.Assert("23", "customer address must have a country code (BR-KSA-10)", is.Present),
						),
					),
				),
			),
			rules.Field("lines",
				rules.Assert("24", "line must be present in standard or associated debit/credit notes", is.Present),
				rules.Each(
					rules.Field("taxes",
						rules.Assert("25", "line taxes are required for standard tax invoices and associated credit notes and debit notes (BR-KSA-52)", is.Present),
					),
					rules.Field("total",
						rules.Assert("26", "line total amount is required for standard tax invoices and associated credit notes and debit notes (BR-KSA-53)", is.Present),
					),
				),
			),
		),

		// Customer
		rules.Field("customer",
			rules.Assert("27", "customer must be present", is.Present),
			rules.Field("addresses",
				rules.Each(
					rules.When(
						is.Func("customer country code is SA", countryCodeIsSA),
						rules.Field("street",
							rules.Assert("28", "customer address in SA must have a street name (BR-KSA-63)", is.Present),
						),
						rules.Field("num",
							rules.Assert("29", "customer address in SA must have a building number (BR-KSA-63)", is.Present),
						),
						rules.Field("code",
							rules.Assert("30", "customer address in SA must have a postal code (BR-KSA-63)", is.Matches(`^\d{5}$`)),
						),
						rules.Field("locality",
							rules.Assert("31", "customer address in SA must have a city name (BR-KSA-63)", is.Present),
						),
						rules.Field("street_extra",
							rules.Assert("32", "customer address in SA must have a district name (BR-KSA-63)", is.Present),
						),
					),
				),
			),

			rules.Assert("33", "buyer identification is valid",
				is.Func("buyer must be either VAT registered or have a valid identification (BR-KSA-14), (BR-KSA-81)", customerValidIdentity),
			),
		),

		// Export invoice
		rules.When(
			is.Func("export invoice", invoiceIsExport),
			rules.Field("customer",
				rules.Field("tax_id",
					rules.Assert("34", "export invoices must not have buyer VAT registration number (BR-KSA-46)",
						is.Empty,
					),
				),
			),
		),

		// Not export invoice
		rules.When(
			is.Func("not export invoice", invoiceIsNotExport),
			rules.Field("customer",
				rules.Assert("35", "customer must be present", is.Present),
				rules.Field("tax_id",
					rules.Assert("36", "If exists, VAT numbers should be valid (BR-KSA-44)",
						is.Func("valid", vatNumberExistsAndValid),
					),
				),
			),
		),

		// Stardard or (simplified and summary)
		rules.When(
			is.Or(
				is.Func("invoice is standard", invoiceIsStandard),
				is.Func("invoice is simplified and summary", invoiceIsSimplifiedAndSummary),
			),
			rules.Field("delivery",
				rules.Assert("37", "delivery must be present", is.Present),
				rules.Field("date",
					rules.Assert("38", "delivery period must have a supply date (BR-KSA-15)", is.Present),
				),
			),
		),

		// Simplified and summary
		rules.When(
			is.Or(
				is.Func("invoice is simplified and summary", invoiceIsSimplifiedAndSummary),
			),
			rules.Field("delivery",
				rules.Field("period",
					rules.Assert("39", "supply must have a delivery period", is.Present),
					rules.Field("start",
						rules.Assert("40", "delivery start date must be present (BR-KSA-72)", is.Present),
					),
					rules.Field("end",
						rules.Assert("41", "delivery end date must be present (BR-KSA-72)", is.Present),
					),
				),
			),
			rules.Field("customer",
				rules.Assert("42", "customer must be present for simplified, summary invoices", is.Present),
				rules.Field("name",
					rules.Assert("43", "customer name must be present for simplified, summary invoices (BR-KSA-71)", is.Present),
				),
			),
		),

		// EDU or HEA exemptions
		rules.When(
			is.Func("has EDU or HEA tax exemption", invoiceHasEDUOrHEAExemption),
			rules.Field("customer",
				rules.Assert("44", "customer must have a national ID (NAT) when tax exemption is VATEX-SA-EDU or VATEX-SA-HEA (BR-KSA-49)",
					is.Func("customer has NAT identity", customerHasNATIdentity),
				),
			),
		),

		// (Simplified or credit/debit note) and EDU or HEA exemptions
		rules.When(
			is.Func("simplified or associated credit/debit note and EDU or HEA exemptions", invoiceIsSimplifiedAndEDUOrHEAExemption),
			rules.Field("customer",
				rules.Assert("45", "customer must be present for simplified or associated credit/debit note and EDU or HEA exemptions", is.Present),
				rules.Field("name",
					rules.Assert("46", "customer name must be present (BR-KSA-25)", is.Present),
				),
			),
		),
	)
}

// getInvTypeCode extracts the ZATCA invoice type code from the invoice.
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

// invoiceIsStandard returns true when the invoice's ZATCA type code
// starts with "01" (Standard Tax Invoice).
func invoiceIsStandard(val any) bool {
	code := getInvTypeCode(val)
	return code != "" && code[:2] == "01"
}

// invoiceIsExport returns true when the invoice is an export invoice
// (KSA-2 position 5 = 1).
func invoiceIsExport(val any) bool {
	code := getInvTypeCode(val)
	return code != "" && code[InvTypePosExport] == '1'
}

// invoiceIsSummary returns true when the invoice is a summary invoice.
func invoiceIsSummary(val any) bool {
	code := getInvTypeCode(val)
	return code != "" && cbc.Code(code).In(invTypesSummary...)
}

// invoiceIsSimplifiedAndSummary returns true when the invoice is a simplified summary tax invoice
func invoiceIsSimplifiedAndSummary(val any) bool {
	return invoiceIsSummary(val) && !invoiceIsStandard(val)
}

// invoiceIsNotExport returns true when the invoice is not an export
func invoiceIsNotExport(val any) bool {
	return !invoiceIsExport(val)
}

// invoiceIsCreditOrDebitNote returns true when the invoice type is
// credit note (381) or debit note (383).
func invoiceIsCreditOrDebitNote(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	return inv.Type == bill.InvoiceTypeCreditNote || inv.Type == bill.InvoiceTypeDebitNote
}

// invoiceBuyerCountryIsSA returns true when the customer's country
// code is "SA".
// countryCodeIsSA returns true when the address country code is "SA".
func countryCodeIsSA(val any) bool {
	addr, ok := val.(*org.Address)
	if !ok || addr == nil {
		return false
	}
	return addr.Country == l10n.ISOCountryCode("SA")
}

// invoiceHasEDUOrHEAExemption returns true when any invoice line has a
// tax combo with VATEX-SA-EDU or VATEX-SA-HEA exemption code (BR-KSA-49).
func invoiceHasEDUOrHEAExemption(val any) bool {
	return invoiceHasExemption(val, []cbc.Code{VatexEdu, VatexHea})
}

func invoiceIsSimplifiedAndEDUOrHEAExemption(val any) bool {
	return invoiceHasEDUOrHEAExemption(val) && !invoiceIsStandard(val)
}

// invoiceHasExemption returns true when any invoice line has a
// tax combo with the specified exemptions code
func invoiceHasExemption(val any, exemptions []cbc.Code) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	for _, line := range inv.Lines {
		vat := line.GetTaxes().Get("VAT")
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

// customerHasNATIdentity returns true when the customer party has at least
// one identity with type NAT (National ID).
func customerHasNATIdentity(val any) bool {
	party, ok := val.(*org.Party)
	if !ok || party == nil {
		return false
	}
	for _, id := range party.Identities {
		if id.Type == cbc.Code("NAT") {
			return true
		}
	}
	return false
}

// customerValidIdentity returns true when the customer party has either a VAT
// identity or one from: TIN, CRN, MOM, MLS, 700, SAG, NAT, GCC, IQA, OTH)
// customerValidIdentity returns true when the customer has a VAT registration
// (TaxID with code) or exactly one identity. If neither is present, it returns false.
func customerValidIdentity(value any) bool {
	party, _ := value.(*org.Party)
	if party == nil {
		return true
	}
	if party.TaxID != nil && !party.TaxID.Code.IsEmpty() {
		return true
	}
	return len(party.Identities) == 1
}

// partyHasSingleIdentity returns true when the supplier party has at most
// one identity
func partyHasSingleIdentity(value any) bool {
	party, _ := value.(*org.Party)
	if party == nil || len(party.Identities) <= 1 {
		return true
	}
	return false
}

// vatNumberExsitsAndValid validates VAT number against a regular expressin
// only if it exists
func vatNumberExistsAndValid(value any) bool {
	taxID, _ := value.(*tax.Identity)
	if taxID == nil {
		return true
	}
	if taxID.Code.IsEmpty() {
		return true
	}
	match, _ := regexp.MatchString(vatIDPattern, taxID.Code.String())
	return match
}
