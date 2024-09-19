// Package addons contains all the base addon packages offerd by GOBL
// to provide additional tax handling.
//
// Addons are designed to use normalization and validation rules,
// alongside scenarios, to add support for fields and values that
// cannot be determined automatically or reliably from a GOBL
// document into the destination format, and vice versa.
package addons

import (
	// Import all the addons to ensure they're ready to use.
	_ "github.com/invopop/gobl/addons/es/facturae"
	_ "github.com/invopop/gobl/addons/es/tbai"
	_ "github.com/invopop/gobl/addons/gr/mydata"
	_ "github.com/invopop/gobl/addons/it/sdi"
	_ "github.com/invopop/gobl/addons/mx/cfdi"
	_ "github.com/invopop/gobl/addons/pt/saft"
)
