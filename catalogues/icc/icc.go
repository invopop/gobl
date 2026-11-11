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
	// ExtKeyIncoterms2010 is used to identify the ICC Incoterms® 2010 rule code.
	ExtKeyIncoterms2010 cbc.Key = "icc-incoterms-2010"

	// ExtKeyIncoterms2020 is used to identify the ICC Incoterms® 2020 rule code.
	ExtKeyIncoterms2020 cbc.Key = "icc-incoterms-2020"
)
