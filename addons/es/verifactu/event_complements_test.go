package verifactu_test

import (
	"testing"

	"github.com/invopop/gobl/addons/es/verifactu"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withAddonContext() rules.WithContext {
	return func(rc *rules.Context) {
		rc.Set(rules.ContextKey(verifactu.V1), tax.AddonForKey(verifactu.V1))
	}
}

func TestInvoiceAnomalyLaunchRules(t *testing.T) {
	t.Run("valid with all checks enabled", func(t *testing.T) {
		c := validInvoiceAnomalyLaunch()
		err := rules.Validate(c, withAddonContext())
		require.NoError(t, err)
	})

	t.Run("valid with no checks enabled", func(t *testing.T) {
		c := &verifactu.InvoiceAnomalyLaunch{}
		err := rules.Validate(c, withAddonContext())
		require.NoError(t, err)
	})

	t.Run("missing fingerprint count when check enabled", func(t *testing.T) {
		c := validInvoiceAnomalyLaunch()
		c.FingerprintCount = nil
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "fingerprint count is required when check is enabled")
	})

	t.Run("missing signature count when check enabled", func(t *testing.T) {
		c := validInvoiceAnomalyLaunch()
		c.SignatureCount = nil
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "signature count is required when check is enabled")
	})

	t.Run("missing chain count when check enabled", func(t *testing.T) {
		c := validInvoiceAnomalyLaunch()
		c.ChainCount = nil
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "chain count is required when check is enabled")
	})

	t.Run("missing date count when check enabled", func(t *testing.T) {
		c := validInvoiceAnomalyLaunch()
		c.DateCount = nil
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "date count is required when check is enabled")
	})
}

func TestInvoiceAnomalyRules(t *testing.T) {
	t.Run("empty complement", func(t *testing.T) {
		c := new(verifactu.InvoiceAnomaly)
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "anomaly type is required")
	})

	t.Run("description too long", func(t *testing.T) {
		c := &verifactu.InvoiceAnomaly{
			Type:        "01",
			Description: string(make([]byte, 101)),
		}
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "description must be 100 characters or less")
	})

	t.Run("valid with invoice", func(t *testing.T) {
		c := &verifactu.InvoiceAnomaly{
			Type: "01",
			Invoice: &verifactu.AnomalousInvoice{
				IssuerTaxCode: "B85905495",
				Code:          "SAMPLE-001",
				IssueDate:     cal.MakeDate(2024, 11, 15),
			},
		}
		err := rules.Validate(c, withAddonContext())
		require.NoError(t, err)
	})

	t.Run("invoice missing required fields", func(t *testing.T) {
		c := &verifactu.InvoiceAnomaly{
			Type:    "01",
			Invoice: &verifactu.AnomalousInvoice{},
		}
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "issuer tax code is required")
		assert.Contains(t, err.Error(), "invoice code is required")
	})
}

func TestEventAnomalyLaunchRules(t *testing.T) {
	t.Run("missing count when check enabled", func(t *testing.T) {
		c := &verifactu.EventAnomalyLaunch{
			FingerprintCheck: true,
			SignatureCheck:   true,
		}
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "fingerprint count is required when check is enabled")
		assert.Contains(t, err.Error(), "signature count is required when check is enabled")
	})
}

func TestEventAnomalyRules(t *testing.T) {
	t.Run("empty complement", func(t *testing.T) {
		c := new(verifactu.EventAnomaly)
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "anomaly type is required")
	})

	t.Run("event missing required fields", func(t *testing.T) {
		c := &verifactu.EventAnomaly{
			Type:  "07",
			Event: &verifactu.AnomalousEvent{},
		}
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "event type is required")
		assert.Contains(t, err.Error(), "timestamp is required")
		assert.Contains(t, err.Error(), "fingerprint is required")
	})
}

func TestInvoiceExportRules(t *testing.T) {
	t.Run("empty complement", func(t *testing.T) {
		c := new(verifactu.InvoiceExport)
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "period start is required")
		assert.Contains(t, err.Error(), "period end is required")
		assert.Contains(t, err.Error(), "first invoice record is required")
		assert.Contains(t, err.Error(), "last invoice record is required")
		assert.Contains(t, err.Error(), "discarded flag is required")
	})
}

func TestEventExportRules(t *testing.T) {
	t.Run("empty complement", func(t *testing.T) {
		c := new(verifactu.EventExport)
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "period start is required")
		assert.Contains(t, err.Error(), "period end is required")
		assert.Contains(t, err.Error(), "first event record is required")
		assert.Contains(t, err.Error(), "last event record is required")
		assert.Contains(t, err.Error(), "discarded flag is required")
	})
}

func TestEventSummaryRules(t *testing.T) {
	t.Run("empty complement", func(t *testing.T) {
		c := new(verifactu.EventSummary)
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "event type counts are required")
	})

	t.Run("event type entry missing type", func(t *testing.T) {
		c := &verifactu.EventSummary{
			Events: []*verifactu.EventTypeCount{
				{Count: 5},
			},
		}
		err := rules.Validate(c, withAddonContext())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "event type is required")
	})

	t.Run("valid summary", func(t *testing.T) {
		c := &verifactu.EventSummary{
			Events: []*verifactu.EventTypeCount{
				{Type: "01", Count: 2},
				{Type: "10", Count: 4},
			},
			TaxTotal:    "3780.00",
			AmountTotal: "21780.00",
		}
		err := rules.Validate(c, withAddonContext())
		require.NoError(t, err)
	})
}

func validInvoiceAnomalyLaunch() *verifactu.InvoiceAnomalyLaunch {
	count := 150
	return &verifactu.InvoiceAnomalyLaunch{
		FingerprintCheck: true,
		FingerprintCount: &count,
		SignatureCheck:   true,
		SignatureCount:   &count,
		ChainCheck:       true,
		ChainCount:       &count,
		DateCheck:        true,
		DateCount:        &count,
	}
}
