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
	if tc.Rate.HasPrefix(TaxRateExempt) && tc.Rate != TaxRateExempt {
		for _, m := range taxRateVATExemptMigrationMap {
			if m.Key == tc.Rate {
				tc.Rate = tax.RateExempt
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
	Key cbc.Key
	Ext tax.Extensions
}{
	{
		Key: TaxRateExempt.With(TaxRateOutlay),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M01",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateIntrastate),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M02",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateImports),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M04",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateExports),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M05",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateSuspension),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M06",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateInternalOps),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M07",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateSmallRetail),
		Ext: tax.Extensions{
			KeyATTaxCode:           tax.ExtValue("exempt"),
			oldExtKeyExemptionCode: "M09",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateExemptScheme),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M10",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateTobacco),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M11",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateTravel),
		Ext: tax.Extensions{
			KeyATTaxCode:           tax.ExtValue("exempt"),
			oldExtKeyExemptionCode: "M12",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateSecondHand),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M13",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateArt),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M14",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateAntiques),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M15",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateTransmission),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M16",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateOther),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M19",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateFlatRate),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M20",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateNonDeductible),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M21",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateConsignment),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M25",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateWaste),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M30",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateCivilEng),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M31",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateGreenhouse),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M32",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateWoods),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M33",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateB2B),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M40",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateIntraEU),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M41",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateRealEstate),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M42",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateGold),
		Ext: tax.Extensions{
			oldExtKeyExemptionCode: "M43",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateNonTaxable),
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
		ext[ExtKeyRegion] = "PT-AC"
	case "30":
		ext[ExtKeyRegion] = "PT-MA"
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
