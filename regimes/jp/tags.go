package jp

import "github.com/invopop/gobl/cbc"

// Invoice tag keys specific to Japan.
const (
	// TagExport identifies a zero-rated export invoice (輸出免税). All JCT lines must use the "zero" rate when this tag
	// is present.
	TagExport cbc.Key = "export"

	// TagSimplified identifies a Simplified Qualified Invoice (簡易適格請求書). Issued by retail, restaurant, taxi, and
	// similar businesses transacting with unspecified numbers of customers. Buyer name is not required.
	TagSimplified cbc.Key = "simplified"

	// TagSelfBilling identifies a buyer-issued invoice (仕入明細書) that substitutes for a qualified invoice once
	// confirmed by the supplier.
	TagSelfBilling cbc.Key = "self-billing"
)
