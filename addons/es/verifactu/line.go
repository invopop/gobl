package verifactu

import (
	"fmt"
	"slices"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/es"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/validation"
)

func normalizeLine(line *bill.Line) {
	if line == nil {
		return
	}

	// Add exemption notes
	for _, tc := range line.Taxes {
		if tc == nil {
			continue
		}

		ec := tc.Ext.Get(ExtKeyExempt)
		if ec == cbc.CodeEmpty {
			continue
		}

		if hasExemptionNote(line.Notes) {
			continue
		}

		et := exemptionText(ec)
		if et == "" {
			continue
		}

		line.Notes = append(line.Notes, &org.Note{
			Key:  org.NoteKeyLegal,
			Src:  ExtKeyExempt,
			Code: ec,
			Text: et,
		})
	}
}

func validateLine(line *bill.Line) error {
	if line == nil {
		return nil
	}

	return validation.ValidateStruct(line,
		validation.Field(&line.Taxes,
			tax.SetHasOneOf(tax.CategoryVAT, es.TaxCategoryIGIC, es.TaxCategoryIPSI),
			validation.Skip,
		),
		validation.Field(&line.Notes,
			validation.By(validateLineNotes(line)),
			validation.Skip,
		),
	)
}

func validateLineNotes(line *bill.Line) validation.RuleFunc {
	return func(val any) error {
		notes, _ := val.([]*org.Note) //nolint:errcheck
		ec := lineTaxExemptionCode(line)
		return validateExemptionNotes(notes, ec)
	}
}

func validateExemptionNotes(notes []*org.Note, ec cbc.Code) error {
	count := 0
	for i, n := range notes {
		if isExemptionNote(n) {
			if ec == "" {
				return fmt.Errorf("(%d: unexpected exemption note)", i)
			}
			if count > 0 {
				return fmt.Errorf("(%d: too many exemption notes)", i)
			}
			if ec != n.Code {
				return fmt.Errorf("(%d: note code %s must match extension %s)", i, n.Code, ec)
			}
			count++
		}
	}

	if ec != "" && count == 0 {
		return fmt.Errorf("missing exemption note for code %s", ec)
	}

	return nil
}

func hasExemptionNote(notes []*org.Note) bool {
	return slices.ContainsFunc(notes, isExemptionNote)
}

func exemptionText(exemptionCode cbc.Code) string {
	extDef := tax.ExtensionForKey(ExtKeyExempt)
	codeDef := extDef.CodeDef(exemptionCode)
	if codeDef == nil {
		return ""
	}
	return codeDef.Name.In(i18n.ES)
}

func isExemptionNote(n *org.Note) bool {
	return n.Key == org.NoteKeyLegal && n.Src == ExtKeyExempt
}

func lineTaxExemptionCode(line *bill.Line) cbc.Code {
	// Check VAT first, then IGIC, then IPSI
	categories := []cbc.Code{
		tax.CategoryVAT,
		es.TaxCategoryIGIC,
		es.TaxCategoryIPSI,
	}

	for _, cat := range categories {
		tc := line.Taxes.Get(cat)
		if tc != nil {
			if ec := tc.Ext.Get(ExtKeyExempt); ec != cbc.CodeEmpty {
				return ec
			}
		}
	}

	return ""
}
