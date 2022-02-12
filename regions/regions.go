package regions

import (
	"github.com/invopop/gobl/region"
	"github.com/invopop/gobl/regions/es"
	"github.com/invopop/gobl/regions/gb"
)

// Init is used to ensure all the regions have been prepared
func Init() {
	region.Register(es.New())
	region.Register(gb.New())
	// region.Register(nl.New()) // pending
}
