package flow6

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/catalogues/iso"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// PaymentCDARCodeFor returns the CDAR ProcessConditionCode that a
// bill.Payment maps to in Flow 6: advice → 211 ("Paiement transmis"),
// receipt → 212 ("Encaissée"). Any other type returns false.
func PaymentCDARCodeFor(pmtType cbc.Key) (string, bool) {
	switch pmtType {
	case bill.PaymentTypeAdvice:
		return "211", true
	case bill.PaymentTypeReceipt:
		return "212", true
	}
	return "", false
}

// normalizePayment surfaces the CDAR ProcessConditionCode for the
// payment on the fr-ctc-flow6-status-code extension and defaults the
// roles on the payment's parties — mirrors what normalizeStatus does
// for bill.Status. Advice payments are issued by the payer (BY → SE);
// receipt payments by the payee (SE → BY).
func normalizePayment(pmt *bill.Payment) {
	if pmt == nil {
		return
	}
	code, ok := PaymentCDARCodeFor(pmt.Type)
	if !ok {
		return
	}
	if pmt.Ext.Get(ExtKeyStatusCode) == "" {
		pmt.Ext = pmt.Ext.Set(ExtKeyStatusCode, cbc.Code(code))
	}
	// Roles per CDV side: 211 is buyer-issued (payer → payee), 212 is
	// seller-issued (payee → payer).
	switch pmt.Type {
	case bill.PaymentTypeAdvice:
		setPartyRoleDefault(pmt.Customer, RoleSE)
		setPartyRoleDefault(pmt.Supplier, RoleBY)
	case bill.PaymentTypeReceipt:
		setPartyRoleDefault(pmt.Supplier, RoleSE)
		setPartyRoleDefault(pmt.Customer, RoleBY)
	}
}

func billPaymentRules() *rules.Set {
	return rules.For(new(bill.Payment),
		rules.Field("type",
			rules.Assert("01", "payment type must be 'advice' (CDAR 211) or 'receipt' (CDAR 212) for a Flow 6 CDV message — 'request' is not a CDV event",
				is.In(bill.PaymentTypeAdvice, bill.PaymentTypeReceipt),
			),
		),
		rules.Field("supplier",
			rules.Assert("02", "supplier is required (BR-FR-CDV-13)",
				is.Present,
			),
			rules.Assert("03", "supplier must have an identity with ISO/IEC 6523 scheme 0002 (SIREN)",
				is.Func("supplier has SIREN", paymentPartyHasSIRENIdentity),
			),
		),
		rules.Field("customer",
			rules.Assert("04", "customer is required (BR-FR-CDV-CL-04)",
				is.Present,
			),
			rules.Assert("05", "customer must have an identity with ISO/IEC 6523 scheme 0002 (SIREN)",
				is.Func("customer has SIREN", paymentPartyHasSIRENIdentity),
			),
		),
		rules.Field("lines",
			rules.Assert("06", "exactly one payment line is required (a CDV references a single invoice)",
				is.Func("exactly one line", paymentHasExactlyOneLine),
			),
			rules.Each(
				rules.Field("document",
					rules.Assert("07", "payment line must reference the underlying invoice (BR-FR-CDV-10)",
						is.Present,
					),
					rules.Field("code",
						rules.Assert("08", "referenced invoice code is required (BR-FR-CDV-10)",
							is.Present,
						),
					),
					rules.Field("issue_date",
						rules.Assert("09", "referenced invoice issue date is required (BR-FR-CDV-11)",
							is.Present,
						),
					),
				),
			),
		),
		rules.Assert("10", "ext.fr-ctc-flow6-status-code must match the CDAR ProcessConditionCode implied by the payment type",
			is.Func("status code matches type", paymentStatusCodeMatchesType),
		),
		rules.Assert("11", "every CDV party identity scheme must be in the Flow 6 allow-list (STC 0231 is rejected — it is a Flow 2 invoice concept)",
			is.Func("CDV identity schemes allowed", paymentPartiesIdentitySchemesAllowed),
		),
	)
}

// paymentPartyHasSIRENIdentity reports whether the party carries at
// least one identity scoped with iso-scheme-id=0002.
func paymentPartyHasSIRENIdentity(v any) bool {
	p, ok := v.(*org.Party)
	if !ok || p == nil {
		return false
	}
	for _, id := range p.Identities {
		if id == nil || id.Ext.IsZero() {
			continue
		}
		if id.Ext.Get(iso.ExtKeySchemeID).String() == identitySchemeIDSIREN {
			return true
		}
	}
	return false
}

func paymentHasExactlyOneLine(v any) bool {
	lines, ok := v.([]*bill.PaymentLine)
	if !ok {
		return false
	}
	return len(lines) == 1
}

func paymentStatusCodeMatchesType(v any) bool {
	pmt, ok := v.(*bill.Payment)
	if !ok || pmt == nil {
		return true
	}
	code := pmt.Ext.Get(ExtKeyStatusCode).String()
	if code == "" {
		return true
	}
	expected, ok := PaymentCDARCodeFor(pmt.Type)
	if !ok {
		return true
	}
	return code == expected
}

// paymentPartiesIdentitySchemesAllowed rejects any identity whose
// iso-scheme-id falls outside allowedFlow6IdentitySchemes on either
// the supplier or customer of the payment.
func paymentPartiesIdentitySchemesAllowed(v any) bool {
	pmt, ok := v.(*bill.Payment)
	if !ok || pmt == nil {
		return true
	}
	for _, p := range []*org.Party{pmt.Supplier, pmt.Customer, pmt.Payee} {
		if p == nil {
			continue
		}
		for _, id := range p.Identities {
			if id == nil || id.Ext.IsZero() {
				continue
			}
			scheme := id.Ext.Get(iso.ExtKeySchemeID).String()
			if scheme != "" && !containsString(allowedFlow6IdentitySchemes, scheme) {
				return false
			}
		}
	}
	return true
}

func containsString(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
