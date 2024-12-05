package in

import "github.com/invopop/gobl/cbc"

const (
	// ChargeKeyCompensationCess is used for addtional charges added to an invoice for the special
	// compensation "cess" (cess means tax or levy) which may be appled as a percentage or specific
	// amount based on valumes or other criteria.
	//
	// Typically this tariff is applied to luxury goods, tobacco, and other items that are considered
	// harmful to the environment or society.
	ChargeKeyCompensationCess cbc.Key = "compensation-cess"
)
