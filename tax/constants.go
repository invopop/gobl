package tax

import "github.com/invopop/gobl/cbc"

// Standard tax categories that may be shared between countries.
const (
	CategoryST  cbc.Code = "ST"  // Sales Tax
	CategoryVAT cbc.Code = "VAT" // Value Added Tax
	CategoryGST cbc.Code = "GST" // Goods and Services Tax
)

// Most commonly used keys. Local regions may add their own rate
// keys or extend them.
const (
	RateExempt       cbc.Key = "exempt"
	RateZero         cbc.Key = "zero"
	RateStandard     cbc.Key = "standard"
	RateIntermediate cbc.Key = "intermediate"
	RateReduced      cbc.Key = "reduced"
	RateSuperReduced cbc.Key = "super-reduced"
	RateSpecial      cbc.Key = "special"
	RateOther        cbc.Key = "other"
)

// Standard tax tags
const (
	TagSimplified    cbc.Key = "simplified"
	TagReverseCharge cbc.Key = "reverse-charge"
	TagCustomerRates cbc.Key = "customer-rates"
	TagSelfBilled    cbc.Key = "self-billed"
	TagPartial       cbc.Key = "partial"
	TagB2G           cbc.Key = "b2g"
	TagExport        cbc.Key = "export"
	TagEEA           cbc.Key = "eea" // European Economic Area, used with exports
)
