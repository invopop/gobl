// Package icc defines the International Chamber of Commerce (ICC) code lists, such
// as the Incoterms® delivery codes, that may be used in documents.
package icc

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterCatalogueDef("icc.json")
}

const (
	// ExtKeyDeliveryCode is used to identify the ICC Incoterms® 2020 delivery code.
	ExtKeyDeliveryCode cbc.Key = "icc-delivery-code"
)
