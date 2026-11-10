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
		key     cbc.Key
		codes   []cbc.Code
		sources []string
	}{
		{
			key:   icc.ExtKeyIncoterms2010Rule,
			codes: []cbc.Code{"EXW", "FCA", "CPT", "CIP", "DAT", "DAP", "DDP", "FAS", "FOB", "CFR", "CIF"},
			sources: []string{
				"https://iccwbo.org/news-publications/icc-rules-guidelines/incoterms-rules-2010/#top",
				"https://www.dhl.com/content/dam/dhl/global/dhl-global-forwarding/documents/pdf/dhl-glo-dgf-incoterms-2010-brochure.pdf",
			},
		},
		{
			key:   icc.ExtKeyIncoterms2020Rule,
			codes: []cbc.Code{"EXW", "FCA", "CPT", "CIP", "DAP", "DPU", "DDP", "FAS", "FOB", "CFR", "CIF"},
			sources: []string{
				"https://iccwbo.org/business-solutions/incoterms-rules/",
				"https://2go.iccwbo.org/incoterms-2020-practical-free-wallchart.html",
				"https://www.dhl.com/content/dam/dhl/global/dhl-global-forwarding/documents/pdf/glo-dgf-incoterms-2020-brochure.pdf",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.key.String(), func(t *testing.T) {
			ext := tax.ExtensionForKey(tt.key)
			require.NotNil(t, ext)
			assert.Equal(t, tt.codes, cbc.DefinitionCodes(ext.Values))
			assert.Equal(t, []string{"bill/delivery-details"}, ext.Usage.Schemas)
			urls := make([]string, len(ext.Sources))
			for i, src := range ext.Sources {
				urls[i] = src.URL
			}
			assert.Equal(t, tt.sources, urls)
		})
	}
}
