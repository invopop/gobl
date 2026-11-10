// Package tax encapsulates models related to taxation.
package tax

import (
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/schema"
)

// Short schema IDs of objects used inside documents
const (
	ShortSchemaCombo     = "tax/combo"
	ShortSchemaNote      = "tax/note"
	ShortSchemaRateTotal = "tax/rate-total"
)

func init() {
	schema.Register(schema.GOBL.Add("tax"),
		AddonDef{},
		AddonList{},
		CatalogueDef{},
		Combo{},
		CorrectionDefinition{},
		CorrectionSet{},
		Extensions{},
		Filter{},
		Identity{},
		Note{},
		RateTotal{},
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
		filterRules(),
	)
}
