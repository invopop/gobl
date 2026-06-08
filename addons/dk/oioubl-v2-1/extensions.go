package oioubl

import (
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/i18n"
	"github.com/invopop/gobl/pkg/here"
)

// Extension keys for OIOUBL 2.1.
const (
	// ExtKeyPaymentID carries the OIOUBL cbc:PaymentID "kortart" code used with
	// the Giro (PaymentMeansCode 50) and FIK (PaymentMeansCode 93) payment
	// methods. The payment reference itself is carried separately in the GOBL
	// payment instruction's Ref (emitted as cbc:InstructionID).
	ExtKeyPaymentID cbc.Key = "dk-oioubl-payment-id"

	// ExtKeyTaxCategory carries the OIOUBL taxcategoryid-1.1 category code
	// emitted as cac:TaxCategory/cbc:ID. The addon normalizer derives it from
	// the EN 16931 UNTDID tax category so the gobl.ubl serializer emits it
	// directly instead of mapping the codes itself.
	ExtKeyTaxCategory cbc.Key = "dk-oioubl-tax-category"

	// ExtKeyPaymentChannel carries the OIOUBL paymentchannelcode-1.1 value
	// emitted as cbc:PaymentChannelCode. The addon normalizer derives it from
	// the payment means so the gobl.ubl serializer emits it directly.
	ExtKeyPaymentChannel cbc.Key = "dk-oioubl-payment-channel"

	// ExtKeyResponseCode carries the OIOUBL responsecode-1.1 value for an
	// ApplicationResponse (Invoice Response) status line, emitted as
	// cac:Response/cbc:ResponseCode. The addon normalizer derives it from the
	// GOBL status event (and conversely derives the event from a parsed value)
	// so the gobl.ubl serializer emits it directly instead of mapping the codes.
	ExtKeyResponseCode cbc.Key = "dk-oioubl-response-code"
)

// OIOUBL taxcategoryid-1.1 category codes.
const (
	ExtValueTaxCategoryStandardRated cbc.Code = "StandardRated"
	ExtValueTaxCategoryZeroRated     cbc.Code = "ZeroRated"
	ExtValueTaxCategoryReverseCharge cbc.Code = "ReverseCharge"
)

// OIOUBL paymentchannelcode-1.1 values.
const (
	ExtValuePaymentChannelIBAN cbc.Code = "IBAN"
	ExtValuePaymentChannelGiro cbc.Code = "DK:GIRO"
	ExtValuePaymentChannelFIK  cbc.Code = "DK:FIK"
)

// OIOUBL responsecode-1.1 values accepted by the ApplicationResponse schematron
// (F-APR018 allows five of the six codelist values; ProfileAccept is rejected).
const (
	ExtValueResponseCodeBusinessAccept  cbc.Code = "BusinessAccept"
	ExtValueResponseCodeBusinessReject  cbc.Code = "BusinessReject"
	ExtValueResponseCodeTechnicalAccept cbc.Code = "TechnicalAccept"
	ExtValueResponseCodeTechnicalReject cbc.Code = "TechnicalReject"
	ExtValueResponseCodeProfileReject   cbc.Code = "ProfileReject"
)

// OIOUBL Giro (code 50) PaymentID values.
const (
	ExtValuePaymentIDGiro01 cbc.Code = "01"
	ExtValuePaymentIDGiro04 cbc.Code = "04"
	ExtValuePaymentIDGiro15 cbc.Code = "15"
)

// OIOUBL FIK (code 93) PaymentID values.
const (
	ExtValuePaymentIDFIK71 cbc.Code = "71"
	ExtValuePaymentIDFIK73 cbc.Code = "73"
	ExtValuePaymentIDFIK75 cbc.Code = "75"
)

// giroPaymentIDs and fikPaymentIDs are the PaymentID values OIOUBL allows for
// each method (F-LIB147 for Giro, F-LIB152 family for FIK).
var (
	giroPaymentIDs = []cbc.Code{ExtValuePaymentIDGiro01, ExtValuePaymentIDGiro04, ExtValuePaymentIDGiro15}
	fikPaymentIDs  = []cbc.Code{ExtValuePaymentIDFIK71, ExtValuePaymentIDFIK73, ExtValuePaymentIDFIK75}
)

var extensions = []*cbc.Definition{
	{
		Key: ExtKeyPaymentID,
		Name: i18n.String{
			i18n.EN: "OIOUBL Payment ID (Giro/FIK kortart)",
			i18n.DA: "OIOUBL Betalings-ID (Giro/FIK kortart)",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				Identifies the OIOUBL ` + "`cbc:PaymentID`" + ` "kortart" code that
				accompanies the Danish Giro and FIK payment methods. It is mandatory
				for ` + "`PaymentMeansCode`" + ` 50 (Giro, values 01/04/15) and 93
				(FIK, values 71/73/75), per the OIOUBL Common schematron.
			`),
		},
		Values: []*cbc.Definition{
			{Code: ExtValuePaymentIDGiro01, Name: i18n.String{i18n.EN: "Giro payment type 01"}},
			{Code: ExtValuePaymentIDGiro04, Name: i18n.String{i18n.EN: "Giro payment type 04"}},
			{Code: ExtValuePaymentIDGiro15, Name: i18n.String{i18n.EN: "Giro payment type 15"}},
			{Code: ExtValuePaymentIDFIK71, Name: i18n.String{i18n.EN: "FIK payment type 71"}},
			{Code: ExtValuePaymentIDFIK73, Name: i18n.String{i18n.EN: "FIK payment type 73"}},
			{Code: ExtValuePaymentIDFIK75, Name: i18n.String{i18n.EN: "FIK payment type 75"}},
		},
	},
	{
		Key: ExtKeyTaxCategory,
		Name: i18n.String{
			i18n.EN: "OIOUBL Tax Category",
			i18n.DA: "OIOUBL Momskategori",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The OIOUBL ` + "`taxcategoryid-1.1`" + ` category code emitted as
				` + "`cac:TaxCategory/cbc:ID`" + `. Derived from the EN 16931 UNTDID
				tax category during normalization (S → StandardRated, Z → ZeroRated,
				AE → ReverseCharge; VAT-exempt is reported as ZeroRated).
			`),
		},
		Values: []*cbc.Definition{
			{Code: ExtValueTaxCategoryStandardRated, Name: i18n.String{i18n.EN: "Standard rated"}},
			{Code: ExtValueTaxCategoryZeroRated, Name: i18n.String{i18n.EN: "Zero rated"}},
			{Code: ExtValueTaxCategoryReverseCharge, Name: i18n.String{i18n.EN: "Reverse charge"}},
		},
	},
	{
		Key: ExtKeyPaymentChannel,
		Name: i18n.String{
			i18n.EN: "OIOUBL Payment Channel",
			i18n.DA: "OIOUBL Betalingskanal",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The OIOUBL ` + "`paymentchannelcode-1.1`" + ` value emitted as
				` + "`cbc:PaymentChannelCode`" + `. Derived from the payment means
				during normalization: Giro (50) → DK:GIRO, FIK (93) → DK:FIK, and
				other settled means → IBAN (direct debit carries no channel).
			`),
		},
		Values: []*cbc.Definition{
			{Code: ExtValuePaymentChannelIBAN, Name: i18n.String{i18n.EN: "IBAN bank transfer"}},
			{Code: ExtValuePaymentChannelGiro, Name: i18n.String{i18n.EN: "Danish Giro"}},
			{Code: ExtValuePaymentChannelFIK, Name: i18n.String{i18n.EN: "Danish FIK"}},
		},
	},
	{
		Key: ExtKeyResponseCode,
		Name: i18n.String{
			i18n.EN: "OIOUBL Response Code",
			i18n.DA: "OIOUBL Svarkode",
		},
		Desc: i18n.String{
			i18n.EN: here.Doc(`
				The OIOUBL ` + "`responsecode-1.1`" + ` value emitted as
				` + "`cac:Response/cbc:ResponseCode`" + ` on an Invoice Response.
				Derived from the GOBL status event during normalization (accepted →
				BusinessAccept, rejected → BusinessReject, acknowledged →
				TechnicalAccept, error → TechnicalReject); the reverse is applied
				when parsing an inbound document.
			`),
		},
		Values: []*cbc.Definition{
			{Code: ExtValueResponseCodeBusinessAccept, Name: i18n.String{i18n.EN: "Business accept"}},
			{Code: ExtValueResponseCodeBusinessReject, Name: i18n.String{i18n.EN: "Business reject"}},
			{Code: ExtValueResponseCodeTechnicalAccept, Name: i18n.String{i18n.EN: "Technical accept"}},
			{Code: ExtValueResponseCodeTechnicalReject, Name: i18n.String{i18n.EN: "Technical reject"}},
			{Code: ExtValueResponseCodeProfileReject, Name: i18n.String{i18n.EN: "Profile reject"}},
		},
	},
}
