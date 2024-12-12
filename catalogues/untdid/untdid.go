// Package untdid defines the UN/EDIFACT data elements contained in the UNTDID (United Nations Trade Data Interchange Directory).
package untdid

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterCatalogueDef("untdid.json")
}

const (
	// ExtKeyDocumentType is used to identify the UNTDID 1001 document type code.
	ExtKeyDocumentType cbc.Key = "untdid-document-type"

	// ExtKeyReference is used to identify the UNTDID 1153 reference codes
	// qualifiers.
	ExtKeyReference cbc.Key = "untdid-reference"

	// ExtKeyPaymentMeans is used to identify the UNTDID 4461 payment means code.
	ExtKeyPaymentMeans cbc.Key = "untdid-payment-means"

	// ExtKeyAllowance is used to identify the UNTDID 5189 allownce codes
	// used in discounts.
	ExtKeyAllowance cbc.Key = "untdid-allowance"

	// ExtKeyTaxCategory is used to identify the UNTDID 5305 duty/tax/fee category code.
	ExtKeyTaxCategory cbc.Key = "untdid-tax-category"

	// ExtKeyItemType is used to identify the UNTDID 7143 item type code.
	ExtKeyItemType cbc.Key = "untdid-item-type"

	// ExtKeyCharge is used to identify the UNTDID 7161 charge codes.
	ExtKeyCharge cbc.Key = "untdid-charge"
)
