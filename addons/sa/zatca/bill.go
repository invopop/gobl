package zatca

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/catalogues/cef"
	"github.com/invopop/gobl/catalogues/untdid"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

var taxExemptionReason = map[string]string{
	// Category E — Exempt from VAT
	"VATEX-SA-29":   "Financial services mentioned in Article 29 of the VAT Regulations",
	"VATEX-SA-29-7": "Life insurance services mentioned in Article 29 of the VAT Regulations",
	"VATEX-SA-30":   "Real estate transactions mentioned in Article 30 of the VAT Regulations",

	// Category Z — Zero rated goods
	"VATEX-SA-32":    "Export of goods",
	"VATEX-SA-33":    "Export of services",
	"VATEX-SA-34-1":  "The international transport of Goods",
	"VATEX-SA-34-2":  "International transport of passengers",
	"VATEX-SA-34-3":  "Services directly connected and incidental to a Supply of international passenger transport",
	"VATEX-SA-34-4":  "Supply of a qualifying means of transport",
	"VATEX-SA-34-5":  "Any services relating to Goods or passenger transportation, as defined in article twenty five of these Regulations",
	"VATEX-SA-35":    "Medicines and medical equipment",
	"VATEX-SA-36":    "Qualifying metals",
	"VATEX-SA-EDU":   "Private education to citizen",
	"VATEX-SA-HEA":   "Private healthcare to citizen",
	"VATEX-SA-MLTRY": "Supply of qualified military goods",

	// Category O — Not subject to VAT
	"VATEX-SA-OOS": "Reason is free text, to be provided by the taxpayer on case to case basis",
}

func normalizeInvoice(inv *bill.Invoice) {
	if inv == nil {
		return
	}

	// Ensure Tax object exists
	if inv.Tax == nil {
		inv.Tax = &bill.Tax{}
	}

	// Always set rounding to currency for SA ZATCA
	inv.Tax.Rounding = tax.RoundingRuleCurrency

	// Ensure issue time exists
	if inv.IssueTime == nil {
		inv.IssueTime = &cal.Time{}
	}

	// BR-KSA-O-01: "Not subject to VAT" lines must have a 0% rate.
	for _, line := range inv.Lines {
		vat := line.Taxes.Get(tax.CategoryVAT)
		if vat == nil {
			continue
		}
		if vat.Key == tax.KeyOutsideScope {
			vat.Percent = &num.PercentageZero
		}
	}

	// BR-KSA-83
	for _, line := range inv.Lines {
		vat := line.Taxes.Get(tax.CategoryVAT)
		if vat == nil {
			continue
		}
		ec := vat.Ext.Get(cef.ExtKeyVATEX)
		if ec == cbc.CodeEmpty {
			continue
		}
		reason, ok := taxExemptionReason[string(ec)]
		if !ok {
			continue
		}
		untdidCat := vat.Ext.Get(untdid.ExtKeyTaxCategory)
		if untdidCat == cbc.CodeEmpty || hasTaxNoteForCategory(inv.Tax, untdidCat) {
			continue
		}
		inv.Tax = inv.Tax.MergeNotes(&tax.Note{
			Category: tax.CategoryVAT,
			Key:      vat.Key,
			Text:     reason,
			Ext:      tax.Extensions{untdid.ExtKeyTaxCategory: untdidCat},
		})
	}
}

func hasTaxNoteForCategory(bt *bill.Tax, untdidCat cbc.Code) bool {
	if bt == nil {
		return false
	}
	for _, n := range bt.Notes {
		if n == nil {
			continue
		}
		if n.Category == tax.CategoryVAT && n.Ext.Get(untdid.ExtKeyTaxCategory) == untdidCat {
			return true
		}
	}
	return false
}

func billDiscountRules() *rules.Set {
	return rules.For(new(bill.Discount),
		rules.Field("taxes",
			rules.Assert("01", "taxes are required (BR-32)", is.Present),
		),
	)
}
