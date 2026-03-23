// Package tax encapsulates models related to taxation.
package tax

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
	"github.com/invopop/validation"
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
		scenarioSetRules(),
		setRules(),
		comboRules(),
	)
}

// RequireIdentityCode is deprecated: do not use. Ensure a rule exists to check for the
// tax Identity Code.
var RequireIdentityCode = validation.Skip
