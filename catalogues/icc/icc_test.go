package icc_test

import (
	"testing"

	_ "github.com/invopop/gobl"
	"github.com/invopop/gobl/catalogues/icc"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	tests := []struct {
		key   cbc.Key
		codes []cbc.Code
	}{
		{
			key:   icc.ExtKeyIncoterms2010Rule,
			codes: []cbc.Code{"EXW", "FCA", "CPT", "CIP", "DAT", "DAP", "DDP", "FAS", "FOB", "CFR", "CIF"},
		},
		{
			key:   icc.ExtKeyIncoterms2020Rule,
			codes: []cbc.Code{"EXW", "FCA", "CPT", "CIP", "DAP", "DPU", "DDP", "FAS", "FOB", "CFR", "CIF"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.key.String(), func(t *testing.T) {
			ext := tax.ExtensionForKey(tt.key)
			require.NotNil(t, ext)
			assert.Equal(t, tt.codes, cbc.DefinitionCodes(ext.Values))
			assert.Equal(t, []string{"bill/delivery-details"}, ext.Usage.Schemas)
			require.Len(t, ext.Sources, 1)
			assert.Equal(t, "https://iccwbo.org/business-solutions/incoterms-rules/", ext.Sources[0].URL)
		})
	}
}
