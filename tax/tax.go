// Package tax encapsulates models related to taxation.
package tax

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
)

func init() {
	schema.Register(schema.GOBL.Add("tax"),
		AddonDef{},
		AddonList{},
		CatalogueDef{},
		CorrectionDefinition{},
		CorrectionSet{},
		Extensions{},
		Identity{},
		Note{},
		RegimeCode(""),
		RegimeDef{},
		Set{},
		Scenario{},
		ScenarioSet{},
		TagSet{},
		Total{},
	)
	rules.Register(
		"tax",
		rules.GOBL.Add("TAX"),
		addonRules(),
		addonDefRules(),
		categoryDefRules(),
		comboRules(),
		correctionDefinitionRules(),
		identityRules(),
		rateDefRules(),
		regimeDefRules(),
		scenarioSetRules(),
		setRules(),
		noteRules(),
	)
}
