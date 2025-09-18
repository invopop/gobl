package saft

import (
	"fmt"
	"slices"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
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

		ec := tc.Ext.Get(ExtKeyExemption)
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
			Src:  ExtKeyExemption,
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
		validation.Field(&line.Quantity, num.Positive),
		validation.Field(&line.Sum, num.ZeroOrPositive),
		validation.Field(&line.Total, num.ZeroOrPositive),
		validation.Field(&line.Discounts,
			validation.Each(
				validation.By(validateBillLineDiscount),
				validation.Skip,
			),
		),
		validation.Field(&line.Notes,
			validation.By(validateLineNotes(line)),
			validation.Skip,
		),
	)
}

func validateBillLineDiscount(val any) error {
	disc, _ := val.(*bill.LineDiscount)
	if disc == nil {
		return nil
	}

	return validation.ValidateStruct(disc,
		validation.Field(&disc.Amount, num.ZeroOrPositive),
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
			if len(strings.TrimSpace(n.Text)) < 5 {
				return fmt.Errorf("(%d: note text must be at least 5 characters long)", i)
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
	extDef := tax.ExtensionForKey(ExtKeyExemption)
	codeDef := extDef.CodeDef(exemptionCode)
	if codeDef == nil {
		return ""
	}
	return codeDef.Name.In(i18n.PT)
}

func isExemptionNote(n *org.Note) bool {
	return n.Key == org.NoteKeyLegal && n.Src == ExtKeyExemption
}

func lineTaxExemptionCode(line *bill.Line) cbc.Code {
	vat := line.Taxes.Get(tax.CategoryVAT)
	if vat == nil {
		return ""
	}

	return vat.Ext[ExtKeyExemption]
}
