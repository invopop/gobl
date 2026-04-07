package zatca

import (
	"github.com/invopop/gobl/catalogues/cef"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/tax"
)

func taxComboRules() *rules.Set {
	return rules.For(new(tax.Combo),
		rules.When(
			is.Func("is zero, exempt, or outside scope", taxComboRequiresVATEX),
			rules.Field("ext",
				rules.Assert("01", "exempt, zero-rated, or outside-scope tax must have a valid SA VATEX code (BR-KSA-CL-04)",
					is.Func("valid SA VATEX code", taxComboHasValidVATEX),
				),
			),
		),

		rules.When(
			is.Func("is standard rate", taxComboIsStandard),
			rules.Field("ext",
				rules.Assert("02", "standard rate tax must not have a VATEX code (BR-KSA-CL-04)",
					tax.ExtensionsExclude(cef.ExtKeyVATEX),
				),
			),
		),
	)
}

// taxComboIsStandard returns true when the tax combo key is standard rate.
func taxComboIsStandard(val any) bool {
	tc, ok := val.(*tax.Combo)
	if !ok || tc == nil {
		return false
	}
	return tc.Key == tax.KeyStandard
}

// taxComboRequiresVATEX returns true when the tax combo key is zero,
// exempt, or outside-scope and therefore requires a SA VATEX code.
func taxComboRequiresVATEX(val any) bool {
	tc, ok := val.(*tax.Combo)
	if !ok || tc == nil {
		return false
	}
	return tc.Key.In(tax.KeyZero, tax.KeyExempt, tax.KeyOutsideScope)
}

// taxComboHasValidVATEX checks that the VATEX code in the extensions
// is one of the valid SA codes for the combo's tax key.
func taxComboHasValidVATEX(val any) bool {
	ext, ok := val.(tax.Extensions)
	if !ok {
		return false
	}
	code := ext.Get(cef.ExtKeyVATEX)
	if code == cbc.CodeEmpty {
		return false
	}
	return code.In(
		Vatex29,
		Vatex29_7,
		Vatex30,
		Vatex32,
		Vatex33,
		Vatex34_1,
		Vatex34_2,
		Vatex34_3,
		Vatex34_4,
		Vatex34_5,
		Vatex35,
		Vatex36,
		VatexEdu,
		VatexHea,
		VatexMltry,
		VatexOutOfScope,
	)
}
