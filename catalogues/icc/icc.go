// Package icc defines the International Chamber of Commerce (ICC) code lists.
package icc

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterCatalogueDef("icc.json")
}

const (
	// ExtKeyIncoterms2010Rule is used to identify the ICC Incoterms® 2010 rule code.
	ExtKeyIncoterms2010Rule cbc.Key = "icc-incoterms-2010-rule"

	// ExtKeyIncoterms2020Rule is used to identify the ICC Incoterms® 2020 rule code.
	ExtKeyIncoterms2020Rule cbc.Key = "icc-incoterms-2020-rule"
)
