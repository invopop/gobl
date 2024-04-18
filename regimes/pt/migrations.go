package pt

import (
	"errors"

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

func migrateInvoiceRates(inv *bill.Invoice) error {
	for _, line := range inv.Lines {
		for _, tax := range line.Taxes {
			if err := migrateInvoiceTaxCombo(tax); err != nil {
				return err
			}
		}
	}
	for _, line := range inv.Discounts {
		for _, tax := range line.Taxes {
			if err := migrateInvoiceTaxCombo(tax); err != nil {
				return err
			}
		}
	}
	for _, line := range inv.Charges {
		for _, tax := range line.Taxes {
			if err := migrateInvoiceTaxCombo(tax); err != nil {
				return err
			}
		}
	}
	return nil
}

func migrateInvoiceTaxCombo(tc *tax.Combo) error {
	if tc.Rate.HasPrefix(TaxRateExempt) && tc.Rate != TaxRateExempt {
		for _, m := range taxRateVATExemptMigrationMap {
			if m.Key == tc.Rate {
				tc.Rate = tax.RateExempt
				tc.Ext = m.Ext
				return nil
			}
		}
		return errors.New("invalid tax rate")
	}
	return nil
}

var taxRateVATExemptMigrationMap = []struct {
	Key cbc.Key
	Ext tax.Extensions
}{
	{
		Key: TaxRateExempt.With(TaxRateOutlay),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M01",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateIntrastate),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M02",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateImports),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M04",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateExports),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M05",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateSuspension),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M06",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateInternalOps),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M07",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateSmallRetail),
		Ext: tax.Extensions{
			KeyATTaxCode:        tax.ExtValue(TaxCodeExempt),
			ExtKeyExemptionCode: "M09",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateExemptScheme),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M10",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateTobacco),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M11",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateTravel),
		Ext: tax.Extensions{
			KeyATTaxCode:        tax.ExtValue(TaxCodeExempt),
			ExtKeyExemptionCode: "M12",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateSecondHand),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M13",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateArt),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M14",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateAntiques),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M15",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateTransmission),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M16",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateOther),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M19",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateFlatRate),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M20",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateNonDeductible),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M21",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateConsignment),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M25",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateWaste),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M30",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateCivilEng),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M31",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateGreenhouse),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M32",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateWoods),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M33",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateB2B),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M40",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateIntraEU),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M41",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateRealEstate),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M42",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateGold),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M43",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateNonTaxable),
		Ext: tax.Extensions{
			ExtKeyExemptionCode: "M99",
		},
	},
}

// 2024-04-17: Migrate zone to VAT tax rows
func migrateTaxIDZoneToLines(inv *bill.Invoice) error {
	if inv.Supplier == nil || inv.Supplier.TaxID == nil || inv.Supplier.TaxID.Zone == "" {
		return nil
	}

	ext := make(tax.Extensions)
	zone := inv.Supplier.TaxID.Zone
	inv.Supplier.TaxID.Zone = ""
	switch zone {
	case "20":
		ext[ExtKeyRegion] = "PT-AC"
	case "30":
		ext[ExtKeyRegion] = "PT-MA"
	default:
		// nothing to do
		return nil
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
	return nil
}
