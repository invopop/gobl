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
	_ "github.com/invopop/gobl/addons/br/nfse"
	_ "github.com/invopop/gobl/addons/co/dian"
	_ "github.com/invopop/gobl/addons/de/xrechnung"
	_ "github.com/invopop/gobl/addons/de/zugferd"
	_ "github.com/invopop/gobl/addons/es/facturae"
	_ "github.com/invopop/gobl/addons/es/tbai"
	_ "github.com/invopop/gobl/addons/es/verifactu"
	_ "github.com/invopop/gobl/addons/eu/en16931"
	_ "github.com/invopop/gobl/addons/fr/choruspro"
	_ "github.com/invopop/gobl/addons/fr/facturx"
	_ "github.com/invopop/gobl/addons/gr/mydata"
	_ "github.com/invopop/gobl/addons/it/sdi"
	_ "github.com/invopop/gobl/addons/it/ticket"
	_ "github.com/invopop/gobl/addons/mx/cfdi"
	_ "github.com/invopop/gobl/addons/pt/saft"
)
