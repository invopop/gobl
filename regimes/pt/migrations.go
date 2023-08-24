package pt

import (
	"errors"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/common"
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
	if tc.Rate.HasPrefix(TaxRateExempt) {
		for _, m := range taxRateVATExemptMigrationMap {
			if m.Key == tc.Rate {
				tc.Rate = common.TaxRateExempt
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
	Ext cbc.CodeMap
}{
	{
		Key: TaxRateExempt.With(TaxRateOutlay),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M01",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateIntrastate),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M02",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateImports),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M04",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateExports),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M05",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateSuspension),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M06",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateInternalOps),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M07",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateSmallRetail),
		Ext: cbc.CodeMap{
			KeyATTaxCode:        TaxCodeExempt,
			ExtKeyExemptionCode: "M09",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateExemptScheme),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M10",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateTobacco),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M11",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateTravel),
		Ext: cbc.CodeMap{
			KeyATTaxCode:        TaxCodeExempt,
			ExtKeyExemptionCode: "M12",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateSecondHand),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M13",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateArt),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M14",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateMargin).With(TaxRateAntiques),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M15",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateTransmission),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M16",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateOther),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M19",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateFlatRate),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M20",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateNonDeductible),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M21",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateConsignment),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M25",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateWaste),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M30",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateCivilEng),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M31",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateGreenhouse),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M32",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateWoods),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M33",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateB2B),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M40",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateIntraEU),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M41",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateRealEstate),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M42",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateReverseCharge).With(TaxRateGold),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M43",
		},
	},
	{
		Key: TaxRateExempt.With(TaxRateNonTaxable),
		Ext: cbc.CodeMap{
			ExtKeyExemptionCode: "M99",
		},
	},
}
