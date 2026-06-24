package nfe_test

import (
	"testing"

	"github.com/invopop/gobl/addons/br/nfe"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// addonExtension finds the extension definition with the given key registered
// in the NF-e addon.
func addonExtension(t *testing.T, key cbc.Key) *cbc.Definition {
	t.Helper()
	ad := tax.AddonForKey(nfe.V4)
	require.NotNil(t, ad, "br-nfe-v4 addon must be registered")
	for _, ext := range ad.Extensions {
		if ext.Key == key {
			return ext
		}
	}
	t.Fatalf("extension %q not found in br-nfe-v4 addon", key)
	return nil
}

func TestPurposeExtension(t *testing.T) {
	def := addonExtension(t, nfe.ExtKeyPurpose)

	tests := []struct {
		code cbc.Code
		name string
	}{
		{nfe.PurposeNormal, "Normal"},
		{nfe.PurposeComplementary, "Complementary"},
		{nfe.PurposeAdjustment, "Adjustment"},
		{nfe.PurposeGoodsReturn, "Goods Return"},
		{nfe.PurposeCreditNote, "Credit Note"},
		{nfe.PurposeDebitNote, "Debit Note"},
	}
	for _, tt := range tests {
		t.Run(tt.code.String(), func(t *testing.T) {
			cd := def.CodeDef(tt.code)
			require.NotNil(t, cd, "purpose code %q must resolve", tt.code)
			assert.Equal(t, tt.name, cd.Name[i18n.EN])
		})
	}
}

func TestOperationTypeExtension(t *testing.T) {
	def := addonExtension(t, nfe.ExtKeyOperationType)

	tests := []struct {
		code cbc.Code
		name string
	}{
		{nfe.OperationInbound, "Inbound"},
		{nfe.OperationOutbound, "Outbound"},
	}
	for _, tt := range tests {
		t.Run(tt.code.String(), func(t *testing.T) {
			cd := def.CodeDef(tt.code)
			require.NotNil(t, cd, "operation-type code %q must resolve", tt.code)
			assert.Equal(t, tt.name, cd.Name[i18n.EN])
		})
	}
}

func TestCreditNoteTypeExtension(t *testing.T) {
	def := addonExtension(t, nfe.ExtKeyCreditNoteType)

	// All five credit-note type codes (01..05) must resolve to a value name.
	for _, code := range []cbc.Code{"01", "02", "03", "04", "05"} {
		t.Run(code.String(), func(t *testing.T) {
			cd := def.CodeDef(code)
			require.NotNil(t, cd, "credit-note-type code %q must resolve", code)
			assert.NotEmpty(t, cd.Name[i18n.EN])
		})
	}
	assert.Equal(t, "Penalty and interest", def.CodeDef("01").Name[i18n.EN])
	assert.Nil(t, def.CodeDef("06"), "credit-note-type 06 must not exist")
}

func TestDebitNoteTypeExtension(t *testing.T) {
	def := addonExtension(t, nfe.ExtKeyDebitNoteType)

	// All eight debit-note type codes (01..08) must resolve to a value name.
	for _, code := range []cbc.Code{"01", "02", "03", "04", "05", "06", "07", "08"} {
		t.Run(code.String(), func(t *testing.T) {
			cd := def.CodeDef(code)
			require.NotNil(t, cd, "debit-note-type code %q must resolve", code)
			assert.NotEmpty(t, cd.Name[i18n.EN])
		})
	}
	assert.Equal(t, "Transfer of credits from Cooperatives", def.CodeDef("01").Name[i18n.EN])
	assert.Nil(t, def.CodeDef("09"), "debit-note-type 09 must not exist")
}
