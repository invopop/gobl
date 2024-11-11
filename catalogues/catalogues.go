// Package catalogues provides a set of re-useable extensions, scenarios, and validators
// for specific international standards that can be re-used and incorporated by addons
// or tax regimes.
package catalogues

import (
	// Ensure all the catalogues are registered
	_ "github.com/invopop/gobl/catalogues/iso"
	_ "github.com/invopop/gobl/catalogues/untdid"
)
