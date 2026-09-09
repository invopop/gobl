// Package iso is used to define ISO/IEC extensions and codes that may be used
// in documents.
package iso

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
)

func init() {
	tax.RegisterCatalogueDef("iso.json")
}

const (
	// ExtKeySchemeID is used by the ISO 6523 scheme identifier.
	ExtKeySchemeID cbc.Key = "iso-scheme-id"
)

// ActorIDScheme is the URI scheme for ISO 6523 participant identifiers, as in
// "iso6523-actorid-upis::0225:356000000". Peppol uses it for participant
// addresses, but the scheme itself is not Peppol-specific.
const ActorIDScheme = "iso6523-actorid-upis"
