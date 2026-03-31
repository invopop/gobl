package zatca

import (
	"regexp"
	"slices"
	"strings"

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
		// BR-KSA-05
		rules.Field("tax",
			rules.Field("ext",
				rules.Assert("05", "document type must be a valid ZATCA type (388, 386, 383, 381)",
					tax.ExtensionsHasCodes(untdid.ExtKeyDocumentType, allowedDocumentTypes...),
				),
			),
		),

		// BR-KSA-56, BR-KSA-17
		rules.When(
			is.Func("credit or debit note", invoiceIsCreditOrDebitNote),
			rules.Field("preceding",
				rules.Assert("56", "credit and debit notes must have a billing reference",
					is.Present,
				),
				// BR-KSA-17
				rules.Each(
					rules.Field("reason",
						rules.Assert("17", "credit and debit notes must contain the reason for issuance",
							is.Present,
						),
					),
				),
			),
		),

		// KSA-25
		rules.Field("issue_time",
			rules.Assert("15", "issue time must be present", is.Present),
		),

		// Supplier rules
		rules.Field("supplier",
			rules.Assert("04", "supplier is required", is.Present),
			// BR-KSA-39
			rules.Field("tax_id",
				rules.Assert("39", "seller VAT registration number is required", is.Present),
				rules.Field("code",
					rules.Assert("39", "seller VAT registration number code is required", is.Present),
				),
			),
			// BR-KSA-40: Seller VAT number must be 15 digits starting/ending with "3".
			// Already validated by the SA regime's taxIdentityRules().
			// BR-KSA-08
			rules.Field("identities",
				rules.Assert("08", "seller must have exactly one identity with a valid scheme ID (CRN, MOM, MLS, 700, SAG, OTH) and alphanumeric code",
					is.Func("valid seller ID", supplierIdentityValid),
				),
			),
			rules.Field("addresses",
				// KSA-17
				rules.Assert("06", "supplier address is required", is.Present),
				rules.Each(
					// BR-KSA-09, BR-KSA-37
					rules.Field("num",
						rules.Assert("06", "supplier address building number is required", is.Present),
						rules.Assert("07", "supplier address building number must contain 4 digits", is.Matches(`^\d{4}$`)),
					),
					rules.Field("street_extra",
						rules.Assert("08", "supplier address must have a district", is.Present),
					),
					rules.Field("street",
						rules.Assert("10", "supplier address must have a street name", is.Present),
					),
					rules.Field("code",
						// BR-KSA-66
						rules.Assert("65", "supplier postal code is required", is.Present),
						rules.Assert("66", "supplier postal code must be 5 digits", is.Matches(`^\d{5}$`)),
					),
					rules.Field("locality",
						rules.Assert("12", "supplier address must have a city name", is.Present),
					),
					rules.Field("country",
						rules.Assert("13", "supplier address must have a country code", is.Present),
					),
				),
			),
		),

		// Customer rules
		rules.Field("customer",
			rules.Assert("10", "customer is required", is.Present),
			// BR-KSA-14
			rules.Assert("14", "buyer must have a tax ID or an identity with a valid ZATCA scheme",
				is.Func("has tax_id or identity", buyerHasTaxIDOrIdentity),
			),
			// KSA-18
			rules.Field("addresses",
				rules.Assert("09", "customer address is required", is.Present),
			),
		),

		// BR-KSA-10: Customer address details required for standard tax invoices
		rules.When(
			is.Func("standard tax invoice", invoiceIsStandardType),
			rules.Field("customer",
				rules.Field("addresses",
					rules.Each(
						rules.Field("num",
							rules.Assert("10", "customer address must have a building number", is.Present),
						),
						rules.Field("street",
							rules.Assert("10", "customer address must have a street name", is.Present),
						),
						rules.Field("locality",
							rules.Assert("11", "customer address must have a city name", is.Present),
						),
						rules.Field("country",
							rules.Assert("12", "customer address must have a country code", is.Present),
						),
						// mapped to district
						rules.Field("street_extra",
							rules.Assert("08", "customer address must have a district", is.Present),
						),
					),
				),
			),
		),

		// BR-KSA-63: Customer address details required when buyer country is SA
		rules.When(
			is.Func("buyer country code is SA", invoiceBuyerCountryIsSA),
			rules.Field("customer",
				rules.Field("addresses",
					rules.Each(
						rules.Field("street",
							rules.Assert("10", "customer address must have a street name", is.Present),
						),
						rules.Field("num",
							rules.Assert("10", "customer address must have a building number", is.Present),
						),
						rules.Field("code",
							rules.Assert("11", "customer address must have a postal code", is.Present),
							// BR-KSA-67
							rules.Assert("67", "buyer postal code must be 5 digits when country is SA",
								is.Matches(`^\d{5}$`),
							),
						),
						rules.Field("locality",
							rules.Assert("12", "customer address must have a city name", is.Present),
						),
						rules.Field("country",
							rules.Assert("13", "customer address must have a country code", is.Present),
						),
						// mapped to district
						rules.Field("street_extra",
							rules.Assert("13", "customer address must have a district", is.Present),
						),
					),
				),
			),
		),

		// BR-KSA-46
		rules.When(
			is.Func("export invoice", invoiceIsExport),
			rules.Field("customer",
				rules.Field("tax_id",
					rules.Assert("46", "export invoices must not have buyer VAT registration number",
						is.Empty,
					),
				),
			),
		),

		// BR-KSA-49
		rules.Assert("49", "EDU/HEA tax exemption requires buyer with NAT identity",
			is.Func("EDU/HEA buyer NAT", invoiceEDUHEARequiresBuyerNAT),
		),

		// BR-KSA-25
		rules.Assert("25", "simplified invoice with EDU/HEA exemption requires buyer name",
			is.Func("simplified EDU/HEA buyer name", invoiceSimplifiedEDUHEARequiresBuyerName),
		),

		// Delivery rules
		rules.Field("delivery",
			rules.Assert("16", "delivery is required", is.Present),
		),
		// BR-KSA-15/BR-KSA-72
		rules.When(
			is.Or(
				is.Func("invoice is a tax invoice", invoiceIsTax),
				is.Func("invoice is a simplified tax summary invoice", invoiceIsSimplifiedAndSummary),
			),
			rules.Field("delivery",
				rules.Field("date",
					rules.Assert("19", "delivery period must have a supply date", is.Present),
				),
				rules.Field("period",
					rules.Assert("13", "supply must have a delivery period", is.Present),
				),
			),
			// BR-KSA-35, BR-KSA-36
			rules.Assert("35", "delivery period end date must be >= supply date",
				is.Func("supply end date >= supply date", invoiceSupplyEndDateValid)),
		),

		// BR-KSA-71
		rules.When(
			is.Func("simplified summary invoice", invoiceIsSimplifiedAndSummary),
			rules.Field("customer",
				rules.Field("name",
					rules.Assert("71", "buyer name is required for simplified summary invoices",
						is.Present,
					),
				),
			),
		),

		// BR-KSA-31
		rules.Assert("31", "simplified invoices only allow third-party, nominal, and summary flags",
			is.Func("simplified invoice flags", invoiceSimplifiedFlagsValid),
		),

		// BR-KSA-07
		rules.Assert("32", "self-billing is not allowed for export invoices",
			is.Func("no self-billing for exports", invoiceNoSelfBillingForExports),
		),

		// BR-KSA-51: Line amount with VAT (KSA-12) = line net amount (BT-131) + line VAT amount (KSA-11).
		// This identity is guaranteed by GOBL's calculation engine; no additional validation needed.

		// BR-KSA-DEC-04: Line amount with VAT (KSA-12) must have max two decimal places.
		// SAR currency precision is 2 decimal places; GOBL rounds to currency subunits automatically.

		// BR-KSA-F-01: All dates must be formatted as YYYY-MM-DD.
		// GOBL's cal.Date type uses ISO 8601 format by default; no additional validation needed.

		// BR-KSA-16: Payment means code must be from UNTDID 4461.
		// EN16931 already maps GOBL payment keys to UNTDID codes and validates presence.

		// BR-KSA-F-02: Allowance/charge indicator values.
		// GOBL separates allowances (discounts) and charges into distinct types; indicator is implicit.

		// BR-KSA-DEC-01: Allowance percentage must be 0.00-100.00 with max 2 decimals.
		// GOBL's num.Percentage type handles precision; no other addon validates this range.

		// BR-KSA-EN16931-03: Allowance amount = base × percentage / 100.
		// GOBL's calculation engine computes this; no additional validation needed.

		// BR-KSA-EN16931-04: Allowance base required when percentage provided.
		// GOBL's calculation engine handles base/percentage interdependency.

		// BR-KSA-EN16931-05: Allowance percentage required when base provided.
		// GOBL's calculation engine handles base/percentage interdependency.

		// BR-KSA-18: VAT category codes must be S, Z, E, or O.
		// EN16931 validates VAT category codes via UNTDID 5305 mapping.

		// BR-KSA-69: Zero-rated VAT requires exemption reason code and text.
		// EN16931 requires VATEX extension for exempt tax (BR-E-10).

		// BR-KSA-EN16931-11: Line net amount = qty × (price / base qty) − allowances.
		// GOBL's calculation engine computes this; no additional validation needed.

		// BR-KSA-CL-01: Currency code must be per ISO 4217.
		// GOBL's currency.Code type only accepts valid ISO 4217 codes.

		// BR-KSA-EN16931-07: Item net price = gross price − allowance.
		// GOBL's calculation engine computes this; no additional validation needed.

		// BR-KSA-EN16931-06: No charge on price level; only allowance (false) permitted.
		// GOBL does not support charges at the item price level; structural constraint.

		// BR-KSA-DEC-03: VAT amount at line level (KSA-11) must have max two decimal places.
		// SAR currency precision is 2 decimal places; GOBL rounds to currency subunits automatically.

		// BR-KSA-F-04: All document amounts and quantities must be positive.
		// GOBL validates amounts during calculation; no additional validation needed.

		// BR-KSA-EN16931-01: Business process (BT-23) must be "reporting:1.0".
		// BT-23 is a UBL-specific field (cbc:ProfileID) set at the XML conversion layer.

		// BR-KSA-04: Issue date (BT-2) must be <= current date.
		// Enforced at the submission/signing layer; not validated here to avoid
		// time-dependent tests and because no other addon validates this.

		// BR-KSA-CL-02: All currencyID attributes must match invoice currency (BT-5).
		// GOBL uses a single currency per invoice; consistency is guaranteed by design.

		// BR-KSA-68: Tax currency code (BT-6) must exist.
		// GOBL's invoice currency field serves as both BT-5 and BT-6; always present.

		// BR-KSA-44: Buyer VAT number must be 15 digits starting/ending with "3" for non-export invoices.
		// Already validated by the SA regime's taxIdentityRules() for all SA tax identities.

		// BR-KSA-EN16931-02
		rules.Field("currency",
			rules.Assert("02", "invoice currency must be SAR",
				is.In("SAR"),
			),
		),

		// BR-KSA-45, BR-KSA-52, BR-KSA-53
		rules.When(
			is.Func("standard tax invoice", invoiceIsStandardType),
			// BR-KSA-45
			rules.Field("customer",
				rules.Field("name",
					rules.Assert("45", "buyer name is required for standard tax invoices",
						is.Present,
					),
				),
			),
			rules.Field("lines",
				rules.Each(
					// BR-KSA-52
					rules.Field("taxes",
						rules.Assert("52", "line taxes are required for standard tax invoices", is.Present),
					),
					// BR-KSA-53
					rules.Field("total",
						rules.Assert("53", "line total is required for standard tax invoices", is.Present),
					),
				),
			),
		),
	)
}

// invoiceIsStandardType returns true when the invoice's ZATCA type code
// starts with "01" (Standard Tax Invoice), meaning this is NOT a simplified
// tax invoice (02).
func invoiceIsStandardType(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	code := inv.Tax.GetExt(ExtKeyInvoiceType)
	return strings.HasPrefix(code.String(), "01")
}

// invoiceBuyerCountryIsSA returns true when the customer's country
// code is "SA".
func invoiceBuyerCountryIsSA(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Customer == nil {
		return false
	}
	if len(inv.Customer.Addresses) == 0 {
		return false
	}
	return inv.Customer.Addresses[0].Country == l10n.ISOCountryCode("SA")
}

// invvoiceisTax returns true when the invoice is a standard tax invoice
func invoiceIsTax(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}

	invType := inv.Tax.GetExt(ExtKeyInvoiceType)
	return inv.Type == bill.InvoiceTypeStandard && invType.In(InvTypesStandard...)
}

// invvoiceisSimplified returns true when the invoice is a simplified tax invoice
func invoiceIsSimplified(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}

	invType := inv.Tax.GetExt(ExtKeyInvoiceType)
	return inv.Type == bill.InvoiceTypeStandard && invType.In(InvTypesSimplified...)
}

// invvoiceisSummary returns true when the invoice is a summary invoice
func invoiceIsSummary(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}

	invType := inv.Tax.GetExt(ExtKeyInvoiceType)
	return inv.Type == bill.InvoiceTypeStandard && invType.In(InvTypesSummary...)
}

// invoiceIsSimplifiedAndSummary returns true when the invoice is a simplified summary tax invoice
func invoiceIsSimplifiedAndSummary(val any) bool {
	return invoiceIsSummary(val) && invoiceIsSimplified(val)
}

// invoiceSupplyEndDateValid returns true when the supply end date (KSA-24)
// is greater than or equal to the supply date (KSA-5), or when either date
// is not present.
func invoiceSupplyEndDateValid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil || inv.Delivery == nil {
		return true
	}
	if inv.Delivery.Date == nil || inv.Delivery.Period == nil {
		return true
	}
	if inv.Delivery.Date.IsZero() || inv.Delivery.Period.End.IsZero() {
		return true
	}
	return !inv.Delivery.Period.End.Before(inv.Delivery.Date.Date)
}

// invoiceSimplifiedFlagsValid checks that simplified invoices (02) only use
// third-party (pos 3), nominal (pos 4), and summary (pos 6) flags.
// Export (pos 5) and self-billing (pos 7) are not allowed.
func invoiceSimplifiedFlagsValid(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	code := inv.Tax.GetExt(ExtKeyInvoiceType).String()
	if len(code) != 7 || !strings.HasPrefix(code, "02") {
		return true // only applies to simplified invoices
	}
	// Position 5 (export) must be 0
	if code[4] == '1' {
		return false
	}
	// Position 7 (self-billing) must be 0
	if code[6] == '1' {
		return false
	}
	return true
}

// invoiceIsExport returns true when the invoice is an export invoice
// (KSA-2 position 5 = 1).
func invoiceIsExport(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return false
	}
	code := inv.Tax.GetExt(ExtKeyInvoiceType).String()
	return len(code) == 7 && code[4] == '1'
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

// invoiceNoSelfBillingForExports checks that self-billing (pos 7=1) is not
// set when export (pos 5=1) is set.
func invoiceNoSelfBillingForExports(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	code := inv.Tax.GetExt(ExtKeyInvoiceType).String()
	if len(code) != 7 {
		return true
	}
	// If export (pos 5) is set, self-billing (pos 7) must not be set
	if code[4] == '1' && code[6] == '1' {
		return false
	}
	return true
}

var alphanumericRegexp = regexp.MustCompile(`^[A-Za-z0-9]+$`)

// supplierIdentityValid checks that the supplier has exactly one identity
// with a valid ZATCA scheme ID and alphanumeric code.
func supplierIdentityValid(val any) bool {
	identities, ok := val.([]*org.Identity)
	if !ok || len(identities) == 0 {
		return false
	}
	count := 0
	for _, id := range identities {
		if id == nil {
			continue
		}
		scheme := id.Ext.Get(ExtKeySellerIDScheme)
		if scheme == cbc.CodeEmpty {
			continue
		}
		if !slices.Contains(sellerIDSchemes, scheme) {
			return false
		}
		if !alphanumericRegexp.MatchString(string(id.Code)) {
			return false
		}
		count++
	}
	return count == 1
}

// buyerHasTaxIDOrIdentity checks that the buyer has either a tax ID
// or an identity with a valid ZATCA buyer scheme ID.
func buyerHasTaxIDOrIdentity(val any) bool {
	p, ok := val.(*org.Party)
	if !ok || p == nil {
		return true
	}
	return p.TaxID != nil || org.IdentityForExtKey(p.Identities, ExtKeyBuyerIDScheme) != nil
}

// invoiceEDUHEARequiresBuyerNAT checks that when any line has a
// VATEX-SA-EDU or VATEX-SA-HEA exemption, the buyer has a NAT identity.
func invoiceEDUHEARequiresBuyerNAT(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	if !invoiceHasVATEXCode(inv, "VATEX-SA-EDU", "VATEX-SA-HEA") {
		return true
	}
	if inv.Customer == nil {
		return false
	}
	id := org.IdentityForExtKey(inv.Customer.Identities, ExtKeyBuyerIDScheme)
	return id != nil && id.Ext.Get(ExtKeyBuyerIDScheme) == "NAT"
}

// invoiceSimplifiedEDUHEARequiresBuyerName checks that simplified invoices
// with EDU/HEA exemptions have a buyer name.
func invoiceSimplifiedEDUHEARequiresBuyerName(val any) bool {
	inv, ok := val.(*bill.Invoice)
	if !ok || inv == nil {
		return true
	}
	code := inv.Tax.GetExt(ExtKeyInvoiceType).String()
	if !strings.HasPrefix(code, "02") {
		return true
	}
	if !invoiceHasVATEXCode(inv, "VATEX-SA-EDU", "VATEX-SA-HEA") {
		return true
	}
	return inv.Customer != nil && inv.Customer.Name != ""
}

// invoiceHasVATEXCode checks whether any line tax combo has one of the
// given VATEX exemption codes.
func invoiceHasVATEXCode(inv *bill.Invoice, codes ...string) bool {
	for _, line := range inv.Lines {
		if line == nil {
			continue
		}
		for _, combo := range line.Taxes {
			if combo == nil {
				continue
			}
			code := combo.Ext.Get(cef.ExtKeyVATEX).String()
			if slices.Contains(codes, code) {
				return true
			}
		}
	}
	return false
}
