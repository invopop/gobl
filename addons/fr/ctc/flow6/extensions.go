package flow6

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
	"github.com/invopop/gobl/tax"
)

// Flow 6 extension keys.
const (
	// ExtKeyRole carries the CDAR RoleCode for a party (UNCL 3035 subset).
	// Applied per populated party (Supplier / Customer / Issuer / Recipient)
	// on a bill.Status message.
	ExtKeyRole cbc.Key = "fr-ctc-role"

	// ExtKeyReasonCode pins the exact CDAR ReasonCode for a bill.Reason.
	// When set, takes precedence over the default_for_key lookup that the
	// converter would otherwise perform from Reason.Key.
	ExtKeyReasonCode cbc.Key = "fr-ctc-reason-code"
)

// Flow 6 party role codes (UNCL 3035 subset accepted by CDAR).
const (
	RoleSE  cbc.Code = "SE"  // Seller
	RoleBY  cbc.Code = "BY"  // Buyer
	RoleWK  cbc.Code = "WK"  // Work/Service receiver
	RoleDFH cbc.Code = "DFH" // Delivery from
	RoleAB  cbc.Code = "AB"  // Bank
	RoleSR  cbc.Code = "SR"  // Sender / issuer on behalf of
	RoleDL  cbc.Code = "DL"  // Dealer / intermediary
	RolePE  cbc.Code = "PE"  // Payee
	RolePR  cbc.Code = "PR"  // Payer
	RoleII  cbc.Code = "II"  // Issuer of invoice
	RoleIV  cbc.Code = "IV"  // Invoicee
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyRole,
		Name: i18n.String{
			i18n.EN: "Party Role Code",
			i18n.FR: "Code rôle partie",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				UNCL 3035 role code carried as the CDAR RoleCode for each
				populated party on a Flow 6 lifecycle message. The normalizer
				fills the obvious defaults (Supplier → SE, Customer → BY)
				and leaves the rest for the caller to set explicitly.
			`),
		},
		Values: []*cbc.Definition{
			{Code: RoleSE, Name: i18n.String{i18n.EN: "Seller"}},
			{Code: RoleBY, Name: i18n.String{i18n.EN: "Buyer"}},
			{Code: RoleWK, Name: i18n.String{i18n.EN: "Work / Service Receiver"}},
			{Code: RoleDFH, Name: i18n.String{i18n.EN: "Delivery From"}},
			{Code: RoleAB, Name: i18n.String{i18n.EN: "Bank"}},
			{Code: RoleSR, Name: i18n.String{i18n.EN: "Sender"}},
			{Code: RoleDL, Name: i18n.String{i18n.EN: "Dealer"}},
			{Code: RolePE, Name: i18n.String{i18n.EN: "Payee"}},
			{Code: RolePR, Name: i18n.String{i18n.EN: "Payer"}},
			{Code: RoleII, Name: i18n.String{i18n.EN: "Issuer of Invoice"}},
			{Code: RoleIV, Name: i18n.String{i18n.EN: "Invoicee"}},
		},
	},
	{
		Key: ExtKeyReasonCode,
		Name: i18n.String{
			i18n.EN: "CDAR Reason Code",
			i18n.FR: "Code motif CDAR",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Exact CDAR ReasonCode pinned on a bill.Reason for Flow 6
				lifecycle messages. The CDAR ReasonCode dimension is 1:N
				with bill.Reason.Key: this extension lets the caller pick
				the precise code within a bucket. When absent, the
				converter falls back to the default_for_key code for
				Reason.Key.
			`),
		},
		Values: reasonCodeDefinitions(),
	},
}

// extValue unwraps a tax.Extensions value whether the rules engine has
// passed it to us by value or by pointer.
func extValue(v any) tax.Extensions {
	switch e := v.(type) {
	case tax.Extensions:
		return e
	case *tax.Extensions:
		if e == nil {
			return tax.Extensions{}
		}
		return *e
	}
	return tax.Extensions{}
}

// reasonCodeDefinitions builds the value list for the fr-ctc-reason-code
// extension from the authoritative reasonTable — avoids drift between
// the helper table and the extension's accepted value set.
func reasonCodeDefinitions() []*cbc.Definition {
	out := make([]*cbc.Definition, len(reasonTable))
	for i, e := range reasonTable {
		out[i] = &cbc.Definition{
			Code: cbc.Code(e.Code),
			Name: i18n.String{i18n.EN: string(e.Key)},
		}
	}
	return out
}
