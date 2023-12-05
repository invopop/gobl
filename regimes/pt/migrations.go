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
	Ext tax.ExtMap
}{
	{
		Key: TaxRateExempt.With(TaxRateOutlay),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M01",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateIntrastate),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M02",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateImports),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M04",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateExports),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M05",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateSuspension),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M06",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateInternalOps),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M07",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateSmallRetail),
		Ext: tax.ExtMap{
			KeyATTaxCode:        cbc.KeyOrCode(TaxCodeExempt),
			ExtKeyExemptionCode: "M09",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateExemptScheme),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M10",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateTobacco),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M11",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateTravel),
		Ext: tax.ExtMap{
			KeyATTaxCode:        cbc.KeyOrCode(TaxCodeExempt),
			ExtKeyExemptionCode: "M12",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateSecondHand),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M13",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateArt),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M14",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateAntiques),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M15",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateTransmission),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M16",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateOther),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M19",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateFlatRate),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M20",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateNonDeductible),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M21",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateConsignment),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M25",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateWaste),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M30",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateCivilEng),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M31",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateGreenhouse),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M32",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateWoods),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M33",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateB2B),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M40",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateIntraEU),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M41",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateRealEstate),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M42",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateGold),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M43",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateNonTaxable),
		Ext: tax.ExtMap{
			ExtKeyExemptionCode: "M99",
		},
	},
}
