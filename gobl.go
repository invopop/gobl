package gobl

import "github.com/invopop/gobl/regions"

func init() {
	// This ensures all the regions are loaded
	regions.Init()
}
