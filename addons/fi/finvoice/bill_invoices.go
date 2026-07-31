package finvoice

import (
	"regexp"
	"strings"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/pay"
	"github.com/invopop/gobl/rules"
	"github.com/invopop/gobl/rules/is"
)

// Finvoice's EpiDetails payment block is mandatory on every invoice,
// including credit notes, so the payment rules below apply unconditionally
// rather than only when an amount is due (en16931 BR-CO-25).
//
// EpiBfiIdentifier (the BIC) is optional in the Finvoice schema and SEPA
// practice no longer requires it for domestic transfers, so the addon
// recommends but does not require it.
func billInvoiceRules() *rules.Set {
	return rules.For(new(bill.Invoice),
		rules.Field("customer",
			rules.Assert("01", "customer is required (Finvoice BuyerPartyDetails)", is.Present),
			rules.Field("name",
				rules.Assert("02", "customer name is required (Finvoice BuyerOrganisationName)", is.Present),
			),
		),
		rules.Field("payment",
			rules.Assert("03", "payment details are required (Finvoice EpiDetails)", is.Present),
			rules.Field("instructions",
				rules.Assert("04", "payment instructions are required (Finvoice EpiDetails)", is.Present),
				rules.Field("key",
					rules.Assert("05", "payment instructions key must be credit-transfer",
						cbc.HasValidKeyIn(pay.MeansKeyCreditTransfer),
					),
				),
				rules.Field("ref",
					rules.Assert("06", "payment reference is required (Finvoice EpiReference)", is.Present),
				),
				rules.Field("credit_transfer",
					rules.Assert("07", "credit transfer details are required (Finvoice EpiAccountID)", is.Present),
					rules.Assert("08", "first credit transfer IBAN is required (Finvoice EpiAccountID)",
						is.Func("first entry has IBAN", firstCreditTransferHasIBAN),
					),
					rules.Each(
						rules.Field("iban",
							rules.AssertIfPresent("09", "credit transfer IBAN is not valid",
								is.StringFunc("mod-97 checksum", validIBAN),
							),
						),
					),
				),
			),
			rules.Field("terms",
				rules.Assert("10", "payment terms are required (Finvoice EpiDateOptionDate)", is.Present),
				rules.Field("due_dates",
					rules.Assert("11", "at least one due date is required (Finvoice EpiDateOptionDate)", is.Present),
				),
			),
		),
	)
}

// firstCreditTransferHasIBAN checks the entry the UBL converter maps to
// EpiAccountID. An empty list is handled by the presence assertion.
func firstCreditTransferHasIBAN(val any) bool {
	cts, ok := val.([]*pay.CreditTransfer)
	if !ok || len(cts) == 0 {
		return true
	}
	return cts[0] != nil && cts[0].IBAN != ""
}

// ibanPattern covers the ISO 13616 shape: country code, two check digits,
// and up to 30 alphanumeric BBAN characters.
var ibanPattern = regexp.MustCompile(`^[A-Z]{2}\d{2}[A-Z0-9]{1,30}$`)

// validIBAN reports whether the IBAN satisfies the ISO 7064 mod 97-10
// checksum. Country-specific lengths are not checked.
func validIBAN(iban string) bool {
	if !ibanPattern.MatchString(iban) {
		return false
	}
	n := 0
	for _, r := range iban[4:] + iban[:4] {
		if r >= 'A' && r <= 'Z' {
			n = (n*100 + int(r-'A') + 10) % 97
		} else {
			n = (n*10 + int(r-'0')) % 97
		}
	}
	return n == 1
}

// normalizePayInstructions strips the grouping spaces conventionally used
// when displaying Finnish reference numbers and RF references.
func normalizePayInstructions(instr *pay.Instructions) {
	instr.Ref = cbc.Code(strings.ReplaceAll(instr.Ref.String(), " ", ""))
}

// normalizePayCreditTransfer converts the IBAN to its machine form: no
// grouping spaces, upper case.
func normalizePayCreditTransfer(ct *pay.CreditTransfer) {
	ct.IBAN = strings.ToUpper(strings.ReplaceAll(ct.IBAN, " ", ""))
}
