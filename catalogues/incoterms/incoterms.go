// Package incoterms defines the ICC Incoterms 2020 extension used to describe
// the delivery terms of goods in documents.
package incoterms

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterCatalogueDef("incoterms.json")
}

const (
	// ExtKeyIncoterm is used for the ICC Incoterms 2020 delivery terms code.
	ExtKeyIncoterm cbc.Key = "icc-incoterm"
)
