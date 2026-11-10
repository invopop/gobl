package pay_test

import (
	"strings"
	"testing"

	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/schema"
	"github.com/stretchr/testify/assert"
)

func TestShortSchemas(t *testing.T) {
	tests := map[string]any{
		pay.ShortSchemaInstructions: pay.Instructions{},
		pay.ShortSchemaRecord:       pay.Record{},
		pay.ShortSchemaTerms:        pay.Terms{},
	}
	for short, obj := range tests {
		t.Run(short, func(t *testing.T) {
			id := schema.Lookup(obj)
			assert.True(t, strings.HasSuffix(id.String(), "/"+short), "schema ID %s", id)
		})
	}
}
