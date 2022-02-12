package regions

import (
	"github.com/invopop/gobl/region"
	"github.com/invopop/gobl/regions/es"
	"github.com/invopop/gobl/regions/gb"
)

func init() {
	region.Register(es.New())
	region.Register(gb.New())
	// region.Register(nl.New()) // pending
}
