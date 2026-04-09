package verifactu

import (
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// InvoiceAnomalyLaunch contains the results of launching an anomaly detection process on
// invoice records. Used with event type 03.
type InvoiceAnomalyLaunch struct {
	// Whether a fingerprint integrity check was performed
	FingerprintCheck bool `json:"fingerprint_check" jsonschema:"title=Fingerprint Integrity Check"`
	// Number of fingerprint records analyzed (required when check is true).
	FingerprintCount *int `json:"fingerprint_count,omitempty" jsonschema:"title=Fingerprint Records Analyzed"`
	// Whether a signature integrity check was performed
	SignatureCheck bool `json:"signature_check" jsonschema:"title=Signature Integrity Check"`
	// Number of signature records analyzed (required when check is true).
	SignatureCount *int `json:"signature_count,omitempty" jsonschema:"title=Signature Records Analyzed"`
	// Whether a chain traceability check was performed
	ChainCheck bool `json:"chain_check" jsonschema:"title=Chain Traceability Check"`
	// Number of chain records analyzed (required when check is true).
	ChainCount *int `json:"chain_count,omitempty" jsonschema:"title=Chain Records Analyzed"`
	// Whether a date traceability check was performed
	DateCheck bool `json:"date_check" jsonschema:"title=Date Traceability Check"`
	// Number of date records analyzed (required when check is true).
	DateCount *int `json:"date_count,omitempty" jsonschema:"title=Date Records Analyzed"`
}

// InvoiceAnomaly describes a detected anomaly in invoice records. Used with event type
// 04.
type InvoiceAnomaly struct {
	// Anomaly type code from L1E list.
	Type cbc.Code `json:"type" jsonschema:"title=Anomaly Type"`
	// Free-text description of the anomaly (max 100 chars).
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// Details of the anomalous invoice, if applicable.
	Invoice *AnomalousInvoice `json:"invoice,omitempty" jsonschema:"title=Anomalous Invoice"`
}

// AnomalousInvoice identifies a specific invoice that triggered an anomaly.
type AnomalousInvoice struct {
	// Tax identification number of the invoice issuer.
	IssuerTaxCode string `json:"issuer_tax_code" jsonschema:"title=Issuer Tax Code"`
	// Full invoice code including series.
	Code string `json:"code" jsonschema:"title=Invoice Number"`
	// Date of issuance.
	IssueDate cal.Date `json:"issue_date" jsonschema:"title=Issue Date"`
}

// EventAnomalyLaunch contains the results of launching an anomaly detection process on
// event records. Used with event type 05.
type EventAnomalyLaunch struct {
	// Whether a fingerprint integrity check was performed
	FingerprintCheck bool `json:"fingerprint_check" jsonschema:"title=Fingerprint Integrity Check"`
	// Number of fingerprint records analyzed (required when check is true).
	FingerprintCount *int `json:"fingerprint_count,omitempty" jsonschema:"title=Fingerprint Records Analyzed"`
	// Whether a signature integrity check was performed
	SignatureCheck bool `json:"signature_check" jsonschema:"title=Signature Integrity Check"`
	// Number of signature records analyzed (required when check is true).
	SignatureCount *int `json:"signature_count,omitempty" jsonschema:"title=Signature Records Analyzed"`
	// Whether a chain traceability check was performed
	ChainCheck bool `json:"chain_check" jsonschema:"title=Chain Traceability Check"`
	// Number of chain records analyzed (required when check is true).
	ChainCount *int `json:"chain_count,omitempty" jsonschema:"title=Chain Records Analyzed"`
	// Whether a date traceability check was performed
	DateCheck bool `json:"date_check" jsonschema:"title=Date Traceability Check"`
	// Number of date records analyzed (required when check is true).
	DateCount *int `json:"date_count,omitempty" jsonschema:"title=Date Records Analyzed"`
}

// EventAnomaly describes a detected anomaly in event records. Used with event type 06.
type EventAnomaly struct {
	// Anomaly type code from L1E list.
	Type cbc.Code `json:"type" jsonschema:"title=Anomaly Type"`
	// Free-text description of the anomaly (max 100 chars).
	Description string `json:"description,omitempty" jsonschema:"title=Description"`
	// Details of the anomalous event, if applicable.
	Event *AnomalousEvent `json:"event,omitempty" jsonschema:"title=Anomalous Event"`
}

// AnomalousEvent identifies a specific event that triggered an anomaly.
type AnomalousEvent struct {
	// Event type code from L2E list.
	Type string `json:"type" jsonschema:"title=Event Type"`
	// Timestamp in ISO 8601 format.
	Timestamp string `json:"timestamp" jsonschema:"title=Timestamp"`
	// SHA-256 fingerprint (64 hex characters).
	Fingerprint string `json:"fingerprint" jsonschema:"title=Fingerprint"`
}

// InvoiceExport contains summary data for an export of invoice records over a period.
// Used with event type 08.
type InvoiceExport struct {
	// Start of the export period in ISO 8601 datetime format.
	Start string `json:"start" jsonschema:"title=Period Start"`
	// End of the export period in ISO 8601 datetime format.
	End string `json:"end" jsonschema:"title=Period End"`
	// First invoice record in the exported period.
	FirstRecord *InvoiceRecord `json:"first" jsonschema:"title=First Invoice"`
	// Last invoice record in the exported period.
	LastRecord *InvoiceRecord `json:"last" jsonschema:"title=Last Invoice"`
	// Number of registration records exported.
	RegistrationCount int `json:"registration_count" jsonschema:"title=Registration Count"`
	// Sum of taxes in the period (Decimal 12,2).
	TaxTotal string `json:"tax_total" jsonschema:"title=Tax Total"`
	// Sum of total amounts in the period (Decimal 12,2).
	AmountTotal string `json:"amount_total" jsonschema:"title=Amount Total"`
	// Number of cancellation records exported.
	CancellationCount int `json:"cancellation_count" jsonschema:"title=Cancellation Count"`
	// Whether any records were discarded ("S" or "N").
	Discarded string `json:"discarded" jsonschema:"title=Records Discarded"`
}

// InvoiceRecord identifies a specific invoice in an export or summary.
type InvoiceRecord struct {
	// Tax identification number of the invoice issuer.
	IssuerTaxCode string `json:"issuer_tax_code" jsonschema:"title=Issuer Tax Code"`
	// Full invoice code including series.
	Code string `json:"code" jsonschema:"title=Invoice Number"`
	// Date of issuance.
	IssueDate cal.Date `json:"issue_date" jsonschema:"title=Issue Date"`
	// SHA-256 fingerprint (64 hex characters).
	Fingerprint string `json:"fingerprint" jsonschema:"title=Fingerprint"`
}

// EventExport contains summary data for an export of event records over a period. Used
// with event type 09.
type EventExport struct {
	// Start of the export period in ISO 8601 datetime format.
	Start string `json:"start" jsonschema:"title=Period Start"`
	// End of the export period in ISO 8601 datetime format.
	End string `json:"end" jsonschema:"title=Period End"`
	// First event record in the exported period.
	FirstRecord *EventRecord `json:"first" jsonschema:"title=First Event"`
	// Last event record in the exported period.
	LastRecord *EventRecord `json:"last" jsonschema:"title=Last Event"`
	// Number of event records exported.
	Count int `json:"count" jsonschema:"title=Event Count"`
	// Whether any records were discarded ("S" or "N").
	Discarded string `json:"discarded" jsonschema:"title=Records Discarded"`
}

// EventRecord identifies a specific event in an export or summary.
type EventRecord struct {
	// Event type code from L2E list.
	Type string `json:"type" jsonschema:"title=Event Type"`
	// Timestamp in ISO 8601 format.
	Timestamp string `json:"timestamp" jsonschema:"title=Timestamp"`
	// SHA-256 fingerprint (64 hex characters).
	Fingerprint string `json:"fingerprint" jsonschema:"title=Fingerprint"`
}

// EventSummary provides a summary of events and invoices over a period. Used with event
// type 10.
type EventSummary struct {
	// Counts per event type.
	Events []*EventTypeCount `json:"events" jsonschema:"title=Event Type Counts"`
	// First invoice record in the period.
	FirstRecord *InvoiceRecord `json:"first_invoice,omitempty" jsonschema:"title=First Invoice"`
	// Last invoice record in the period.
	LastRecord *InvoiceRecord `json:"last_invoice,omitempty" jsonschema:"title=Last Invoice"`
	// Number of registration records in the period.
	RegistrationCount int `json:"registration_count" jsonschema:"title=Registration Count"`
	// Sum of taxes in the period (Decimal 12,2).
	TaxTotal string `json:"tax_total" jsonschema:"title=Tax Total"`
	// Sum of total amounts in the period (Decimal 12,2).
	AmountTotal string `json:"amount_total" jsonschema:"title=Amount Total"`
	// Number of cancellation records in the period.
	CancellationCount int `json:"cancellation_count" jsonschema:"title=Cancellation Count"`
}

// EventTypeCount pairs an event type code with its occurrence count.
type EventTypeCount struct {
	// Event type code from L2E list.
	Type string `json:"type" jsonschema:"title=Event Type"`
	// Number of occurrences.
	Count int `json:"count" jsonschema:"title=Count"`
}

func invoiceAnomalyLaunchRules() *rules.Set {
	return rules.For(new(InvoiceAnomalyLaunch),
		rules.When(is.Func("fingerprint checked", func(v any) bool {
			c, ok := v.(*InvoiceAnomalyLaunch)
			return ok && c != nil && c.FingerprintCheck
		}),
			rules.Field("fingerprint_count",
				rules.Assert("01", "fingerprint count is required when check is enabled", is.Present),
			),
		),
		rules.When(is.Func("signature checked", func(v any) bool {
			c, ok := v.(*InvoiceAnomalyLaunch)
			return ok && c != nil && c.SignatureCheck
		}),
			rules.Field("signature_count",
				rules.Assert("02", "signature count is required when check is enabled", is.Present),
			),
		),
		rules.When(is.Func("chain checked", func(v any) bool {
			c, ok := v.(*InvoiceAnomalyLaunch)
			return ok && c != nil && c.ChainCheck
		}),
			rules.Field("chain_count",
				rules.Assert("03", "chain count is required when check is enabled", is.Present),
			),
		),
		rules.When(is.Func("date checked", func(v any) bool {
			c, ok := v.(*InvoiceAnomalyLaunch)
			return ok && c != nil && c.DateCheck
		}),
			rules.Field("date_count",
				rules.Assert("04", "date count is required when check is enabled", is.Present),
			),
		),
	)
}

func invoiceAnomalyRules() *rules.Set {
	return rules.For(new(InvoiceAnomaly),
		rules.Field("type",
			rules.Assert("01", "anomaly type is required", is.Present),
		),
		rules.Field("description",
			rules.AssertIfPresent("02", "description must be 100 characters or less", is.Length(0, 100)),
		),
		rules.Field("invoice",
			rules.Field("issuer_tax_code",
				rules.Assert("03", "issuer tax code is required", is.Present),
			),
			rules.Field("code",
				rules.Assert("04", "invoice code is required", is.Present),
			),
		),
	)
}

func eventAnomalyLaunchRules() *rules.Set {
	return rules.For(new(EventAnomalyLaunch),
		rules.When(is.Func("fingerprint checked", func(v any) bool {
			c, ok := v.(*EventAnomalyLaunch)
			return ok && c != nil && c.FingerprintCheck
		}),
			rules.Field("fingerprint_count",
				rules.Assert("01", "fingerprint count is required when check is enabled", is.Present),
			),
		),
		rules.When(is.Func("signature checked", func(v any) bool {
			c, ok := v.(*EventAnomalyLaunch)
			return ok && c != nil && c.SignatureCheck
		}),
			rules.Field("signature_count",
				rules.Assert("02", "signature count is required when check is enabled", is.Present),
			),
		),
		rules.When(is.Func("chain checked", func(v any) bool {
			c, ok := v.(*EventAnomalyLaunch)
			return ok && c != nil && c.ChainCheck
		}),
			rules.Field("chain_count",
				rules.Assert("03", "chain count is required when check is enabled", is.Present),
			),
		),
		rules.When(is.Func("date checked", func(v any) bool {
			c, ok := v.(*EventAnomalyLaunch)
			return ok && c != nil && c.DateCheck
		}),
			rules.Field("date_count",
				rules.Assert("04", "date count is required when check is enabled", is.Present),
			),
		),
	)
}

func eventAnomalyRules() *rules.Set {
	return rules.For(new(EventAnomaly),
		rules.Field("type",
			rules.Assert("01", "anomaly type is required", is.Present),
		),
		rules.Field("description",
			rules.AssertIfPresent("02", "description must be 100 characters or less", is.Length(0, 100)),
		),
		rules.Field("event",
			rules.Field("type",
				rules.Assert("03", "event type is required", is.Present),
			),
			rules.Field("timestamp",
				rules.Assert("04", "timestamp is required", is.Present),
			),
			rules.Field("fingerprint",
				rules.Assert("05", "fingerprint is required", is.Present),
			),
		),
	)
}

func invoiceExportRules() *rules.Set {
	return rules.For(new(InvoiceExport),
		rules.Field("start",
			rules.Assert("01", "period start is required", is.Present),
		),
		rules.Field("end",
			rules.Assert("02", "period end is required", is.Present),
		),
		rules.Field("first",
			rules.Assert("03", "first invoice record is required", is.Present),
		),
		rules.Field("last",
			rules.Assert("04", "last invoice record is required", is.Present),
		),
		rules.Field("discarded",
			rules.Assert("05", "discarded flag is required", is.Present),
		),
	)
}

func eventExportRules() *rules.Set {
	return rules.For(new(EventExport),
		rules.Field("start",
			rules.Assert("01", "period start is required", is.Present),
		),
		rules.Field("end",
			rules.Assert("02", "period end is required", is.Present),
		),
		rules.Field("first",
			rules.Assert("03", "first event record is required", is.Present),
		),
		rules.Field("last",
			rules.Assert("04", "last event record is required", is.Present),
		),
		rules.Field("discarded",
			rules.Assert("05", "discarded flag is required", is.Present),
		),
	)
}

func eventSummaryRules() *rules.Set {
	return rules.For(new(EventSummary),
		rules.Field("events",
			rules.Assert("01", "event type counts are required", is.Present),
			rules.Each(
				rules.Field("type",
					rules.Assert("02", "event type is required", is.Present),
				),
			),
		),
	)
}
