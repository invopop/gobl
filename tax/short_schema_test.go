package tax_test

import (
	"strings"
	"testing"

	"github.com/invopop/gobl/schema"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestShortSchemas(t *testing.T) {
	tests := map[string]any{
		tax.ShortSchemaCombo:     tax.Combo{},
		tax.ShortSchemaNote:      tax.Note{},
		tax.ShortSchemaRateTotal: tax.RateTotal{},
	}
	for short, obj := range tests {
		t.Run(short, func(t *testing.T) {
			id := schema.Lookup(obj)
			assert.True(t, strings.HasSuffix(id.String(), "/"+short), "schema ID %s", id)
		})
	}
}
