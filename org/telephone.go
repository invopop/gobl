package org

import (
	"strings"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
	"github.com/invopop/gobl/uuid"
)

// Telephone describes what is expected for a telephone number.
type Telephone struct {
	uuid.Identify
	// Identifier for this number.
	Label string `json:"label,omitempty" jsonschema:"title=Label"`
	// Free-text string that represents the telephone number.
	Number string `json:"num" jsonschema:"title=Number"`
}

func telephoneRules() *rules.Set {
	return rules.For(new(Telephone),
		rules.Field("num",
			rules.Assert("01", "telephone number is required", is.Present),
		),
	)
}

func normalizeTelephone(t *Telephone) {
	uuid.Normalize(&t.UUID)
	t.Label = cbc.NormalizeString(t.Label)
	t.Number = strings.TrimSpace(t.Number)
}
