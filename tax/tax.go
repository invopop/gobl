// Package tax encapsulates models related to taxation.
package tax

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("tax"),
		Identity{},
		Set{},
		Extensions{},
		Total{},
		RegimeCode(""),
		AddonList{},
		RegimeDef{},
		AddonDef{},
		CatalogueDef{},
	)
	rules.Register(
		"tax",
		rules.GOBL.Add("TAX"),
		addonRules(),
		addonDefRules(),
		identityRules(),
		categoryDefRules(),
		correctionDefinitionRules(),
		rateDefRules(),
		regimeDefRules(),
		setRules(),
		comboRules(),
	)
}
