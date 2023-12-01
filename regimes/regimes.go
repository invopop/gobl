// Package regimes simple ensures that each of the individually defined tax regimes
// is loaded correctly and ready to use from other GOBL packages.
package regimes

import (
	// Import all the regime definitions which will automatically
	// add themselves to the tax regime register.
	_ "github.com/invopop/gobl/regimes/ca"
	_ "github.com/invopop/gobl/regimes/co"
	_ "github.com/invopop/gobl/regimes/de"
	_ "github.com/invopop/gobl/regimes/es"
	_ "github.com/invopop/gobl/regimes/fr"
	_ "github.com/invopop/gobl/regimes/gb"
	_ "github.com/invopop/gobl/regimes/it"
	_ "github.com/invopop/gobl/regimes/mx"
	_ "github.com/invopop/gobl/regimes/nl"
	_ "github.com/invopop/gobl/regimes/pl"
	_ "github.com/invopop/gobl/regimes/pt"
	_ "github.com/invopop/gobl/regimes/us"
)
