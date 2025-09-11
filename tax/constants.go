package tax

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
)

// Standard tax categories that may be shared between countries.
const (
	CategoryST  cbc.Code = "ST"  // Sales Tax
	CategoryVAT cbc.Code = "VAT" // Value Added Tax
	CategoryGST cbc.Code = "GST" // Goods and Services Tax
)

// Standard tax combo keys used to identify different tax situations.
const (
	KeyStandard       cbc.Key = "standard"
	KeyZero           cbc.Key = "zero"
	KeyReverseCharge  cbc.Key = "reverse-charge"
	KeyExempt         cbc.Key = "exempt"
	KeyExport         cbc.Key = "export"
	KeyIntraCommunity cbc.Key = "intra-community"
	KeyOutsideScope   cbc.Key = "outside-scope"
)

// Most commonly used rates. Local regions may add their own rate
// keys or extend them.
const (
	RateZero         cbc.Key = "zero"
	RateGeneral      cbc.Key = "general"
	RateIntermediate cbc.Key = "intermediate"
	RateReduced      cbc.Key = "reduced"
	RateSuperReduced cbc.Key = "super-reduced"
	RateSpecial      cbc.Key = "special"
	RateOther        cbc.Key = "other"
)

// Standard tax tags that can be used to indicate special circumstances at the
// level of the document.
const (
	TagSimplified    cbc.Key = "simplified"
	TagReverseCharge cbc.Key = "reverse-charge"
	TagCustomerRates cbc.Key = "customer-rates"
	TagSelfBilled    cbc.Key = "self-billed"
	TagReplacement   cbc.Key = "replacement"
	TagPartial       cbc.Key = "partial"
	TagB2G           cbc.Key = "b2g"
	TagExport        cbc.Key = "export"
	TagEEA           cbc.Key = "eea" // European Economic Area
)

// globalCategories defines the tax categories that can be applied anywhere that
// GOBL does not yet have a specific regime defined for it.
var globalCategories = []*CategoryDef{
	{
		Code: CategoryVAT,
		Name: i18n.String{
			i18n.EN: "VAT",
		},
		Title: i18n.String{
			i18n.EN: "Value Added Tax",
		},
		Retained: false,
		Keys: []*KeyDef{
			{
				Key:  KeyStandard,
				Name: i18n.NewString("Standard"),
			},
			{
				Key:  KeyZero,
				Name: i18n.NewString("Zero"),
			},
			{
				Key:       KeyReverseCharge,
				Name:      i18n.NewString("Reverse charge"),
				NoPercent: true,
			},
			{
				Key:       KeyExempt,
				Name:      i18n.NewString("Exempt"),
				NoPercent: true,
			},
			{
				Key:       KeyExport,
				Name:      i18n.NewString("Export"),
				NoPercent: true,
			},
			{
				Key:       KeyIntraCommunity,
				Name:      i18n.NewString("Intra-community"),
				NoPercent: true,
			},
			{
				Key:       KeyOutsideScope,
				Name:      i18n.NewString("Outside scope"),
				NoPercent: true,
			},
		},
	},
	{
		Code: CategoryGST,
		Name: i18n.String{
			i18n.EN: "GST",
		},
		Title: i18n.String{
			i18n.EN: "Goods and Services Tax",
		},
		Retained: false,
		Keys: []*KeyDef{
			{
				Key:  KeyStandard,
				Name: i18n.NewString("Standard"),
			},
			{
				Key:  KeyZero,
				Name: i18n.NewString("Zero"),
			},
			{
				Key:       KeyExempt,
				Name:      i18n.NewString("Exempt"),
				NoPercent: true,
			},
			{
				Key:       KeyOutsideScope,
				Name:      i18n.NewString("Outside scope"),
				NoPercent: true,
			},
		},
	},
}

// GlobalVAT returns a global VAT category definition that can be use in other regimes.
func GlobalVAT() *CategoryDef {
	return Category(CategoryVAT)
}

// GlobalGST returns a global GST category definition that can be use in other regimes.
func GlobalGST() *CategoryDef {
	return Category(CategoryGST)
}

// GlobalVATKeys returns the keys that are defined for the global VAT category, which can
// be re-used in other regimes subject to VAT.
func GlobalVATKeys() []*KeyDef {
	return GlobalVAT().Keys
}

// GlobalGSTKeys returns the keys that are defined for the global GST category, which can
// be re-used in other regimes subject to GST.
func GlobalGSTKeys() []*KeyDef {
	return GlobalGST().Keys
}

// Category returns a global category definition by its code.
func Category(code cbc.Code) *CategoryDef {
	// Return the category with the given code from the global categories.
	for _, cat := range globalCategories {
		if cat.Code == code {
			return cat
		}
	}
	return nil
}
