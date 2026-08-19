package cy_test

import (
	"testing"

	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/cy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvoiceDocumentScenarios(t *testing.T) {
	tests := []struct {
		name string
		tag  cbc.Key
		note string
	}{
		{
			name: "travel agency margin scheme",
			tag:  cy.TagTravelAgencyMargin,
			note: "Margin scheme - Travel agents / Καθεστώς περιθωρίου - Ταξιδιωτικά πρακτορεία",
		},
		{
			name: "cash accounting scheme",
			tag:  cy.TagCashAccounting,
			note: "Cash accounting / Καθεστώς Ταμειακής Λογιστικής",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv := testInvoiceStandard(t)
			inv.SetTags(tt.tag)
			require.NoError(t, inv.Calculate())
			require.NotNil(t, inv.Tax)
			require.Len(t, inv.Tax.Notes, 1)
			assert.Equal(t, tt.note, inv.Tax.Notes[0].Text)
		})
	}
}
