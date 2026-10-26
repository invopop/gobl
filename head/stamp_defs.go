package head

import "github.com/invopop/gobl/cbc"

// stampDefs is a global register of all the stamp definitions that have
// been registered via regimes and addons.
var stampDefs = cbc.NewDefinitionRegistry()

// RegisterStampDef adds the stamp definitions to the global register, and
// is expected to be called from the init functions of the regimes and addons
// that issue them. Names are prefixed with the issuer, such as
// "TicketBAI: QR Code", and the usage identifies the addon keys or regime
// codes.
func RegisterStampDef(defs ...*cbc.Definition) {
	stampDefs.Add(defs...)
}

// StampDefFor provides the stamp definition for the provider key, or nil.
func StampDefFor(provider cbc.Key) *cbc.Definition {
	return stampDefs.For(provider)
}

// AllStampDefs provides all the registered stamp definitions ordered by key.
func AllStampDefs() []*cbc.Definition {
	return stampDefs.All()
}

// FilterStampDefs provides the registered stamp definitions whose usage
// matches the keys, such as addons, or codes, such as tax regime countries.
func FilterStampDefs(keys []cbc.Key, codes []cbc.Code) []*cbc.Definition {
	return stampDefs.Filter(keys, codes)
}
