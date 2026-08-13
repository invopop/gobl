// Package cef provides codes issue by the "Connecting Europe Facility"
// (CEF Digital) initiative.
package cef

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterCatalogueDef("cef.json")
}

const (
	// ExtKeyVATEX is used for the CEF VATEX exemption codes.
	ExtKeyVATEX cbc.Key = "cef-vatex"
)
