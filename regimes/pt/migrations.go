package pt

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

// Tax Rate Migration
//
// 2023-08-24: Exempt VAT Tax rates have been migrated to use extension codes
// instead of using the tax rate key. This file contains all the data needed to
// automatically migrate the old rate rates to the new format.

// Tax rate exemption tags
const (
	TaxRateExempt        cbc.Key = "exempt"
	TaxRateOutlay        cbc.Key = "outlay"
	TaxRateIntrastate    cbc.Key = "intrastate-export"
	TaxRateImports       cbc.Key = "imports"
	TaxRateExports       cbc.Key = "exports"
	TaxRateSuspension    cbc.Key = "suspension-scheme"
	TaxRateInternalOps   cbc.Key = "internal-operations"
	TaxRateSmallRetail   cbc.Key = "small-retail-scheme"
	TaxRateExemptScheme  cbc.Key = "exempt-scheme"
	TaxRateTobacco       cbc.Key = "tobacco-scheme"
	TaxRateMargin        cbc.Key = "margin-scheme"
	TaxRateTravel        cbc.Key = "travel"
	TaxRateSecondHand    cbc.Key = "second-hand"
	TaxRateArt           cbc.Key = "art"
	TaxRateAntiques      cbc.Key = "antiques"
	TaxRateTransmission  cbc.Key = "goods-transmission"
	TaxRateOther         cbc.Key = "other"
	TaxRateFlatRate      cbc.Key = "flat-rate-scheme"
	TaxRateNonDeductible cbc.Key = "non-deductible"
	TaxRateConsignment   cbc.Key = "consignment-goods"
	TaxRateReverseCharge cbc.Key = "reverse-charge"
	TaxRateWaste         cbc.Key = "waste"
	TaxRateCivilEng      cbc.Key = "civil-eng"
	TaxRateGreenhouse    cbc.Key = "greenhouse"
	TaxRateWoods         cbc.Key = "woods"
	TaxRateB2B           cbc.Key = "b2b"
	TaxRateIntraEU       cbc.Key = "intraeu"
	TaxRateRealEstate    cbc.Key = "real-estate"
	TaxRateGold          cbc.Key = "gold"
	TaxRateNonTaxable    cbc.Key = "non-taxable"
)

func migrateInvoiceRates(inv *bill.Invoice) {
	for _, line := range inv.Lines {
		for _, tax := range line.Taxes {
			migrateInvoiceTaxCombo(tax)
		}
	}
	for _, line := range inv.Discounts {
		for _, tax := range line.Taxes {
			migrateInvoiceTaxCombo(tax)
		}
	}
	for _, line := range inv.Charges {
		for _, tax := range line.Taxes {
			migrateInvoiceTaxCombo(tax)
		}
	}
}

const oldExtKeyExemptionCode cbc.Key = "pt-exemption-code"

func migrateInvoiceTaxCombo(tc *tax.Combo) {
	if tc.Rate.HasPrefix(TaxRateExempt) && tc.Key != TaxRateExempt {
		for _, m := range taxRateVATExemptMigrationMap {
			if m.Rate == tc.Rate {
				tc.Key = tax.KeyExempt
				tc.Rate = cbc.KeyEmpty
				tc.Ext = m.Ext
				break
			}
		}
	}
	// 2024-09-13: Added after move to addons
	if tc.Ext[oldExtKeyExemptionCode] != "" {
		tc.Ext["pt-saft-exemption"] = tc.Ext[oldExtKeyExemptionCode]
		delete(tc.Ext, oldExtKeyExemptionCode)
	}
}

var taxRateVATExemptMigrationMap = []struct {
	Rate cbc.Key
	Ext  tax.Extensions
}{
	{
		Rate: TaxRateExempt.With(TaxRateOutlay),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M01",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateIntrastate),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M02",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateImports),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M04",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateExports),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M05",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateSuspension),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M06",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateInternalOps),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M07",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateSmallRetail),
		Ext: tax.Extensions{
			KeyATTaxCode:           cbc.Code("exempt"),
			oldExtKeyExemptionCode: "M09",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateExemptScheme),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M10",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateTobacco),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M11",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateMargin).With(TaxRateTravel),
		Ext: tax.Extensions{
			KeyATTaxCode:           cbc.Code("exempt"),
			oldExtKeyExemptionCode: "M12",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateMargin).With(TaxRateSecondHand),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M13",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateMargin).With(TaxRateArt),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M14",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateMargin).With(TaxRateAntiques),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M15",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateTransmission),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M16",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateOther),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M19",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateFlatRate),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M20",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateNonDeductible),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M21",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateConsignment),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M25",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateWaste),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M30",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateCivilEng),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M31",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateGreenhouse),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M32",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateWoods),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M33",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateB2B),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M40",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateIntraEU),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M41",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateRealEstate),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M42",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateGold),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M43",
		},
	},
	{
		Rate: TaxRateExempt.With(TaxRateNonTaxable),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M99",
		},
	},
}

// 2024-04-17: Migrate zone to VAT tax rows
func migrateTaxIDZoneToLines(inv *bill.Invoice) {
	if inv.Supplier == nil || inv.Supplier.TaxID == nil || inv.Supplier.TaxID.Zone == "" { //nolint:staticcheck
		return
	}

	ext := make(tax.Extensions)
	zone := inv.Supplier.TaxID.Zone //nolint:staticcheck
	inv.Supplier.TaxID.Zone = ""    //nolint:staticcheck
	switch zone {
	case "20":
		ext[ExtKeyRegion] = RegionAzores
	case "30":
		ext[ExtKeyRegion] = RegionMadeira
	default:
		// nothing to do
		return
	}

	for _, line := range inv.Lines {
		for _, tc := range line.Taxes {
			if tc.Category == tax.CategoryVAT {
				tc.Ext = ext
			}
		}
	}
	for _, line := range inv.Discounts {
		for _, tc := range line.Taxes {
			if tc.Category == tax.CategoryVAT {
				tc.Ext = ext
			}
		}
	}
	for _, line := range inv.Charges {
		for _, tc := range line.Taxes {
			if tc.Category == tax.CategoryVAT {
				tc.Ext = ext
			}
		}
	}
}
