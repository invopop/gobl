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
	ext := tax.ExtensionForKey(icc.ExtKeyDeliveryCode)
	require.NotNil(t, ext)
	assert.Equal(t,
		[]cbc.Code{"CFR", "CIF", "CIP", "CPT", "DAP", "DDP", "DPU", "EXW", "FAS", "FCA", "FOB"},
		cbc.DefinitionCodes(ext.Values),
	)
	assert.Equal(t, []string{"bill/delivery-details"}, ext.Usage.Schemas)
	require.Len(t, ext.Sources, 1)
	assert.Equal(t, "https://iccwbo.org/business-solutions/incoterms-rules/", ext.Sources[0].URL)
}
